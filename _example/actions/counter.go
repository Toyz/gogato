package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Toyz/gogato"
)

type CounterAction struct {
	Count int `json:"count"`
}

func (c *CounterAction) ActionID() string {
	return "com.helba.counter.action1"
}

func (c *CounterAction) WillAppear(context string, gogato *gogato.Gogato, payload gogato.EventAppearPayLoad) error {
	err := json.Unmarshal(*payload.Settings, c)
	if err != nil {
		return err
	}

	return gogato.SetTitle(context, fmt.Sprintf("%d", c.Count))
}

func (c *CounterAction) KeyDown(context string, gogato *gogato.Gogato, payload gogato.EventKeyPayLoad) error {
	return nil
}

func (c *CounterAction) KeyUp(context string, gogato *gogato.Gogato, payload gogato.EventKeyPayLoad) error {
	c.Count = c.Count + 1

	err := gogato.SetSettings(context, c)
	if err != nil {
		return err
	}

	return gogato.SetTitle(context, fmt.Sprintf("%d", c.Count))
}

func (c *CounterAction) WillDisappear(context string, gogato *gogato.Gogato, payload gogato.EventAppearPayLoad) error {
	return nil
}

func (c *CounterAction) ReceivedSettings(context string, gogato *gogato.Gogato, settings gogato.EventSettingsPayLoad) error {
	return nil
}

func (c *CounterAction) ReceivedGlobalSettings(context string, gogato *gogato.Gogato, settings gogato.EventGlobalSettingsPayLoad) error {
	return nil
}

func (c *CounterAction) FromPropertyInspector(context string, gogato *gogato.Gogato, settings *json.RawMessage) error {
	var data map[string]*json.RawMessage
	json.Unmarshal(*settings, &data)
	if len(data) == 0 {
		return errors.New("from property inspector map was nil")
	}

	var method string
	json.Unmarshal(*data["property_inspector"], &method)
	if method == "" {
		return errors.New("property_inspector was empty")
	}

	switch method {
	case "propertyInspectorConnected":
		return gogato.SendToPropertyInspector(context, c.ActionID(), c)
	case "updateSettings":
		json.Unmarshal(*settings, c)
		_ = gogato.SetTitle(context, fmt.Sprintf("%d", c.Count))
		return gogato.SetSettings(context, c)
	}

	return nil
}