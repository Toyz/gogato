package gogato

import "encoding/json"

/* Event payloads */
type EventBasePayload struct {
	Action  string           `json:"action"`
	Event   elgatoEvent      `json:"event"`
	Context string           `json:"context"`
	Device  string           `json:"device,omitempty"`
	Payload *json.RawMessage `json:"payload"`
}

type EventKeyPayLoad struct {
	Settings         *json.RawMessage `json:"settings"`
	Coordinates      EventCoordinates `json:"coordinates"`
	UserDesiredState int              `json:"userDesiredState"`
	IsInMultiAction  bool             `json:"isInMultiAction"`
}

type EventAppearPayLoad struct {
	Settings        *json.RawMessage `json:"settings"`
	Coordinates     EventCoordinates `json:"coordinates"`
	State           int              `json:"state"`
	IsInMultiAction bool             `json:"isInMultiAction"`
}

type EventGlobalSettingsPayLoad struct {
	Settings *json.RawMessage `json:"settings"`
}

type EventSettingsPayLoad struct {
	Settings        *json.RawMessage `json:"settings"`
	Coordinates     EventCoordinates `json:"coordinates"`
	IsInMultiAction bool             `json:"isInMultiAction"`
}

type EventCoordinates struct {
	Column int `json:"column"`
	Row    int `json:"row"`
}

/* API Payloads */
type register struct {
	Event string `json:"event"`
	UUID  string `json:"uuid"`
}

type sendEvent struct {
	Action  string      `json:"action,omitempty"`
	Event   elgatoEvent `json:"event"`
	Context string      `json:"context"`
	Payload interface{} `json:"payload,omitempty"`
}

/* Send event payloads */
type SetTitlePayload struct {
	Title  string       `json:"title"`
	Target elgatoTarget `json:"target"`
}
