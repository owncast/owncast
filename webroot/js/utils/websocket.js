import { URL_WEBSOCKET } from './constants.js';
/**
 * These are the types of messages that we can handle with the websocket.
 * Mostly used by `websocket.js` but if other components need to handle
 * different types then it can import this file.
 */
export const SOCKET_MESSAGE_TYPES = {
  CHAT: 'CHAT',
  PING: 'PING',
  NAME_CHANGE: 'NAME_CHANGE',
  PONG: 'PONG',
  SYSTEM: 'SYSTEM'
};

export const CALLBACKS = {
  RAW_WEBSOCKET_MESSAGE_RECEIVED: 'rawWebsocketMessageReceived',
  WEBSOCKET_CONNECTED: 'websocketConnected',
  WEBSOCKET_DISCONNECTED: 'websocketDisconnected',
}

const TIMER_WEBSOCKET_RECONNECT = 5000; // ms

export default class Websocket {
  constructor() {
    this.websocket = null;
    this.websocketReconnectTimer = null;

    this.websocketConnectedListeners = [];
    this.websocketDisconnectListeners = [];
    this.rawMessageListeners = [];

    this.send = this.send.bind(this);
    this.createAndConnect = this.createAndConnect.bind(this);
    this.scheduleReconnect = this.scheduleReconnect.bind(this);

    this.createAndConnect();
  }

  createAndConnect() {
    const ws = new WebSocket(URL_WEBSOCKET);
    ws.onopen = this.onOpen.bind(this);
    ws.onclose = this.onClose.bind(this);
    ws.onerror = this.onError.bind(this);
    ws.onmessage = this.onMessage.bind(this);

    this.websocket = ws;
  }

  // Other components should register for websocket callbacks.
  addListener(type, callback) {
    if (type == CALLBACKS.WEBSOCKET_CONNECTED) {
      this.websocketConnectedListeners.push(callback);
    } else if (type == CALLBACKS.WEBSOCKET_DISCONNECTED) {
      this.websocketDisconnectListeners.push(callback);
    } else if (type == CALLBACKS.RAW_WEBSOCKET_MESSAGE_RECEIVED) {
      this.rawMessageListeners.push(callback);
    }
  }

  // Interface with other components

  // Outbound: Other components can pass an object to `send`.
  send(message) {
    // Sanity check that what we're sending is a valid type.
    if (!message.type || !SOCKET_MESSAGE_TYPES[message.type]) {
      console.warn(
        `Outbound message: Unknown socket message type: "${message.type}" sent.`
      );
    }

    const messageJSON = JSON.stringify(message);
    this.websocket.send(messageJSON);
  }

  // Private methods

  // Fire the callbacks of the listeners.

  notifyWebsocketConnectedListeners(message) {
    this.websocketConnectedListeners.forEach(function (callback) {
      callback(message);
    });
  }

  notifyWebsocketDisconnectedListeners(message) {
    this.websocketDisconnectListeners.forEach(function (callback) {
      callback(message);
    });
  }

  notifyRawMessageListeners(message) {
    this.rawMessageListeners.forEach(function (callback) {
      callback(message);
    });
  }

  // Internal websocket callbacks

  onOpen(e) {
    if (this.websocketReconnectTimer) {
      clearTimeout(this.websocketReconnectTimer);
    }

    this.notifyWebsocketConnectedListeners();
  }

  onClose(e) {
    // connection closed, discard old websocket and create a new one in 5s
    this.websocket = null;
    this.notifyWebsocketDisconnectedListeners();
    this.handleNetworkingError('Websocket closed.');
    this.scheduleReconnect();
  }

  // On ws error just close the socket and let it re-connect again for now.
  onError(e) {
    this.handleNetworkingError(`Socket error: ${JSON.parse(e)}`);
    this.websocket.close();
    this.scheduleReconnect();
  }

  scheduleReconnect() {
    this.websocketReconnectTimer = setTimeout(
      this.createAndConnect,
      TIMER_WEBSOCKET_RECONNECT
    );
  }

  /*
  onMessage is fired when an inbound object comes across the websocket.
  If the message is of type `PING` we send a `PONG` back and do not
  pass it along to listeners.
  */
  onMessage(e) {
    try {
      var model = JSON.parse(e.data);
    } catch (e) {
      console.log(e);
    }

    // Send PONGs
    if (model.type === SOCKET_MESSAGE_TYPES.PING) {
      this.sendPong();
      return;
    }

    // Notify any of the listeners via the raw socket message callback.
    this.notifyRawMessageListeners(model);
  }

  // Reply to a PING as a keep alive.
  sendPong() {
    const pong = { type: SOCKET_MESSAGE_TYPES.PONG };
    this.send(pong);
  }

  handleNetworkingError(error) {
    console.error(`Websocket Error. Chat is likely not working. Visit troubleshooting steps to resolve. https://owncast.online/docs/troubleshooting/#chat-is-disabled: ${error}`);
  }
}
