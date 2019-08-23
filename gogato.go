package gogato

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

/* Always return a error so we can have internal logging */
type (
	ElgatoAction interface {
		ActionID() string
		KeyDown(context string, gogato *Gogato, payload EventKeyPayLoad) error
		KeyUp(context string, gogato *Gogato, payload EventKeyPayLoad) error
		WillAppear(context string, gogato *Gogato, payload EventAppearPayLoad) error
		WillDisappear(context string, gogato *Gogato, payload EventAppearPayLoad) error
		ReceivedSettings(context string, gogato *Gogato, settings EventSettingsPayLoad) error
		ReceivedGlobalSettings(context string, gogato *Gogato, settings EventGlobalSettingsPayLoad) error
		FromPropertyInspector(context string, gogato *Gogato, settings *json.RawMessage) error
	}

	Gogato struct {
		actions map[string]ElgatoAction
		conn    *websocket.Conn
	}
)

var (
	port          = kingpin.Flag("port", "Elgato websocket port").String()
	pluginUUID    = kingpin.Flag("pluginUUID", "Elgato plugin UUID").String()
	registerEvent = kingpin.Flag("registerEvent", "Elgato register event").String()
	info          = kingpin.Flag("info", "Elgato json info payload").String()
)

func NewGogato() *Gogato {
	for id := range os.Args {
		item := os.Args[id]
		if item[0] == '-' {
			os.Args[id] = "-" + os.Args[id]
		}
	}

	kingpin.Parse()

	f, err := os.OpenFile("plugin.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Warnf("Failed to create log file: %v", err)
	} else {
		log.SetOutput(f)
	}

	log.Infof("Booted Plugin: %s %s %s %s", *port, *pluginUUID, *registerEvent, *info)

	return &Gogato{
		actions: map[string]ElgatoAction{},
	}
}

func (ga *Gogato) RegisterAction(actions ...ElgatoAction) error {
	for _, action := range actions {
		uuid := action.ActionID()
		if _, ok := ga.actions[uuid]; ok {
			return fmt.Errorf("%s already registered", uuid)
		}

		ga.actions[uuid] = action
	}

	return nil
}

func (ga *Gogato) Run() error {
	if len(ga.actions) == 0 {
		return errors.New("you must register at least one action")
	}

	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("localhost:%s", *port)}
	log.Infof("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	ga.conn = c
	defer ga.conn.Close()

	if err := ga.writeJson(register{
		Event: *registerEvent,
		UUID:  *pluginUUID,
	}); err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Fatalf("read:", err)
				return
			}

			var payload EventBasePayload
			if err := json.Unmarshal(message, &payload); err != nil {
				log.Warnf("unmarshal:", err)
				continue
			}

			log.Infof("Event %s action %s context %s payload %s", payload.Event, payload.Action, payload.Context, message)

			action, ok := ga.actions[payload.Action]
			if !ok {
				if payload.Action == "" {
					payload.Action = "EmptyDeckAction"
				}

				log.Warnf("unknown action: %s", payload.Action)
				continue
			}


			var actionError error
			switch payload.Event {
			case willAppearEvent:
				var appearPayload EventAppearPayLoad
				_ = json.Unmarshal(*payload.Payload, &appearPayload)

				if actionError = action.WillAppear(payload.Context, ga, appearPayload); actionError != nil {
					log.Warnf("%s WillAppear: %v", payload.Action, actionError)
					continue
				}
				break
			case willDisappearEvent:
				var appearPayload EventAppearPayLoad
				_ = json.Unmarshal(*payload.Payload, &appearPayload)

				if actionError = action.WillDisappear(payload.Context, ga, appearPayload); actionError != nil {
					log.Warnf("%s WillAppear: %v", payload.Action, actionError)
					continue
				}
				break
			case onKeyDownEvent:
				var keyPayload EventKeyPayLoad
				_ = json.Unmarshal(*payload.Payload, &keyPayload)

				if actionError = action.KeyDown(payload.Context, ga, keyPayload); actionError != nil {
					log.Warnf("%s KeyDown: %v", payload.Action, actionError)
					continue
				}
				break
			case onKeyUpEvent:
				var keyPayload EventKeyPayLoad
				_ = json.Unmarshal(*payload.Payload, &keyPayload)

				if actionError = action.KeyUp(payload.Context, ga, keyPayload); actionError != nil {
					log.Warnf("%s KeyUP: %v", payload.Action, actionError)
					continue
				}
				break
				/* Setting paylodas */
			case didReceiveSettings:
				var settingPayload EventSettingsPayLoad
				_ = json.Unmarshal(*payload.Payload, &settingPayload)

				if actionError = action.ReceivedSettings(payload.Context, ga, settingPayload); actionError != nil {
					log.Warnf("%s ReceivedSettings: %v", payload.Action, actionError)
					continue
				}
				break
			case didReceiveGlobalSettings:
				var globalSettings EventGlobalSettingsPayLoad
				_ = json.Unmarshal(*payload.Payload, &globalSettings)

				if actionError = action.ReceivedGlobalSettings(payload.Context, ga, globalSettings); actionError != nil {
					log.Warnf("%s ReceivedGlobalSettings: %v", payload.Action, actionError)
					continue
				}
				break
				/* From Insepctor */
			case sendToPlugin:
				if actionError = action.FromPropertyInspector(payload.Context, ga, payload.Payload); actionError != nil {
					log.Warnf("%s FromInspector: %v", payload.Action, actionError)
					continue
				}
				break
			}
		}
	}(&wg)

	wg.Wait()
	return nil
}

func (ga *Gogato) SetTitle(context, title string) error {
	payload := sendEvent{
		Event:   setTitleEvent,
		Context: context,
		Payload: SetTitlePayload{
			Title:  title,
			Target: hardwareAndSoftware,
		},
	}

	return ga.writeJson(payload)
}

func (ga *Gogato) SetSettings(context string, settings interface{}) error {
	payload := sendEvent{
		Event:   setSettingsEvent,
		Context: context,
		Payload: settings,
	}

	return ga.writeJson(payload)
}

func (ga *Gogato) GetSettings(context string) error {
	payload := sendEvent{
		Event:   getSettingsEvent,
		Context: context,
	}

	return ga.writeJson(payload)
}

func (ga *Gogato) SetGlobalSettings(context string, settings interface{}) error {
	payload := sendEvent{
		Event:   setGlobalSettings,
		Context: context,
		Payload: settings,
	}

	return ga.writeJson(payload)
}

func (ga *Gogato) GetGlobalSettings(context string) error {
	payload := sendEvent{
		Event:   getGlobalSettings,
		Context: context,
	}

	return ga.writeJson(payload)
}

/* Error and Ok Alerts */
func (ga *Gogato) ShowOk(context string) error {
	payload := sendEvent{
		Event:   showOkEvent,
		Context: context,
	}

	return ga.writeJson(payload)
}

func (ga *Gogato) ShowAlert(context string) error {
	payload := sendEvent{
		Event:   showAlertEvent,
		Context: context,
	}

	return ga.writeJson(payload)
}

func (ga *Gogato) SendToPropertyInspector(context, id string, data interface{}) error {
	return ga.writeJson(sendEvent{
		Action:  id,
		Event:   sendToPropertyInspector,
		Context: context,
		Payload: data,
	})
}

func (ga *Gogato) writeJson(data interface{}) error {
	return ga.conn.WriteJSON(data)
}
