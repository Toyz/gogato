var websocket = null,
    uuid = null,
    actionInfo = {},
    isQT = navigator.appVersion.includes('QtWebEngine');

function connectElgatoStreamDeckSocket(inPort, inUUID, inRegisterEvent, inInfo, inActionInfo) {
    connectSocket(inPort, inUUID, inRegisterEvent, inActionInfo);
}

function connectSocket(inPort, inUUID, inRegisterEvent, inActionInfo) {
    uuid = inUUID;
    actionInfo = JSON.parse(inActionInfo); // cache the info
    websocket = new WebSocket('ws://localhost:' + inPort);

    websocket.onopen = function () {
        var json = {
            event: inRegisterEvent,
            uuid: inUUID
        };

        websocket.send(JSON.stringify(json));

        // Notify the plugin that we are connected
        sendValueToPlugin('propertyInspectorConnected', 'property_inspector');
    };

    websocket.onmessage = function (evt) {
        // Received message from Stream Deck
        var jsonObj = JSON.parse(evt.data);

        console.log(evt.data);

        if (jsonObj.event === 'sendToPropertyInspector') {
            var payload = jsonObj.payload;

            var sheetID = document.getElementById('count');
            sheetID.value = payload['count'];
        }
    };
}

function updateSettings() {
    var sheetID = document.getElementById('count');
    var payload = {};

    payload.property_inspector = 'updateSettings';
    payload.count = parseInt(sheetID.value) || 0;

    sendPayloadToPlugin(payload);
}

function sendPayloadToPlugin(payload) {
    if (websocket && (websocket.readyState === 1)) {
        const json = {
            'action': actionInfo['action'],
            'event': 'sendToPlugin',
            'context': uuid,
            'payload': payload
        };
        websocket.send(JSON.stringify(json));
    }
}

function openWebsite() {
    if (websocket && (websocket.readyState === 1)) {
        const json = {
            'event': 'openUrl',
            'payload': {
                'url': 'https://blizztrack.com'
            }
        };
        websocket.send(JSON.stringify(json));
    }
}

// our method to pass values to the plugin
function sendValueToPlugin(value, param) {
    if (websocket && (websocket.readyState === 1)) {
        const json = {
            'action': actionInfo['action'],
            'event': 'sendToPlugin',
            'context': uuid,
            'payload': {
                [param]: value
            }
        };
        websocket.send(JSON.stringify(json));
    }
}

if (!isQT) {
    document.addEventListener('DOMContentLoaded', function () {
        initPropertyInspector();
    });
}

window.addEventListener('beforeunload', function (e) {
    e.preventDefault();

    // Notify the plugin we are about to leave
    sendValueToPlugin('propertyInspectorWillDisappear', 'property_inspector');

    // Don't set a returnValue to the event, otherwise Chromium with throw an error.
});

function initPropertyInspector() {
    // Place to add functions
}

/** UTILITIES */

/*
    Quick utility to lighten or darken a color (doesn't take color-drifting, etc. into account)
    Usage:
    fadeColor('#061261', 100); // will lighten the color
    fadeColor('#200867'), -100); // will darken the color
*/
function fadeColor(col, amt) {
    const min = Math.min, max = Math.max;
    const num = parseInt(col.replace(/#/g, ''), 16);
    const r = min(255, max((num >> 16) + amt, 0));
    const g = min(255, max((num & 0x0000FF) + amt, 0));
    const b = min(255, max(((num >> 8) & 0x00FF) + amt, 0));
    return '#' + (g | (b << 8) | (r << 16)).toString(16).padStart(6, 0);
}
