package gogato

/* Elgato Target enum */
type elgatoTarget int

const (
	hardwareAndSoftware elgatoTarget = iota
	hardwareOnly
	softwareOnly
)

type elgatoEvent string

const (
	/* Sent Events */
	setTitleEvent     elgatoEvent = "setTitle"
	setSettingsEvent  elgatoEvent = "setSettings"
	getSettingsEvent  elgatoEvent = "getSettings"
	setGlobalSettings elgatoEvent = "setGlobalSettings"
	getGlobalSettings elgatoEvent = "getGlobalSettings"
	sendToPropertyInspector elgatoEvent = "sendToPropertyInspector"

	/* Received Events */
	willAppearEvent    elgatoEvent = "willAppear"
	willDisappearEvent elgatoEvent = "willDisappear"
	onKeyDownEvent     elgatoEvent = "keyDown"
	onKeyUpEvent       elgatoEvent = "keyUp"
	sendToPlugin       elgatoEvent = "sendToPlugin"

	/* Setting events */
	didReceiveSettings       elgatoEvent = "didReceiveSettings"
	didReceiveGlobalSettings elgatoEvent = "didReceiveGlobalSettings"

	/* Show events */
	showOkEvent    elgatoEvent = "showOk"
	showAlertEvent elgatoEvent = "showAlert"
)
