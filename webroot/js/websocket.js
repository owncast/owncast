/**
 * Websocket.
 * This is a standalone websocket interface that passes messages to `listeners`, and accepts
 * messages to be sent via `send`.
 * If you have a component that would like to be notified about websocket activity you can
 * pass it to `addListener(object)` and implement any of the following methods:
 * 
 * rawWebsocketMessageReceived(msg): Will send all raw JSON objects that come across the socket.
 * websocketConnected/websocketDisconnected: Be notified on the connectivity state of the socket.
 * 
 * This class should not know anything about views.  Note that error states are not being passed to
 * listeners because listeners shouldn't know or care about specific socket-level errors.
 */

import SOCKET_MESSAGE_TYPES from './chat/socketMessageTypes.js';

const LOCAL_TEST = window.location.href.indexOf('localhost:') >= 0;
const URL_WEBSOCKET = LOCAL_TEST
  ? 'wss://goth.land/entry'
  : `${location.protocol === 'https:' ? 'wss' : 'ws'}://${location.host}/entry`;

const TIMER_WEBSOCKET_RECONNECT = 5000; // ms

const CALLBACKS = {
  RAW_WEBSOCKET_MESSAGE_RECEIVED: 'rawWebsocketMessageReceived',
  WEBSOCKET_CONNECTED: 'websocketConnected',
  WEBSOCKET_DISCONNECTED: 'websocketDisconnected',
}

class Websocket {

  constructor() {
    this.websocket = null;
    this.websocketReconnectTimer = null;
    this.listeners = [];

    this.send = this.send.bind(this);

    const ws = new WebSocket(URL_WEBSOCKET);
    ws.onopen = this.onOpen.bind(this);
    ws.onclose = this.onClose.bind(this);
    ws.onerror = this.onError.bind(this);
    ws.onmessage = this.onMessage.bind(this);

    this.websocket = ws;
  }

  // Other components should register for websocket callbacks.
  addListener(listener) {
    this.listeners.push(listener);
  }

  
  // Interface with other components

  // Outbound: Other components can pass an object to `send`.
  send(message) {
    // Sanity check that what we're sending is a valid type.
    if (!message.type || !SOCKET_MESSAGE_TYPES[message.type]) {
      console.warn(`Outbound message: Unknown socket message type: "${message.type}" sent.`);
    }
    
    const messageJSON = JSON.stringify(message);
    this.websocket.send(messageJSON);
  }

  // Private methods

  // Fire the callbacks of the listeners.
  notifyListeners(callbackName, message) {
    this.listeners.forEach(function (listener) {
      if (doesObjectSupportFunction(listener, callbackName)) {
        listener[callbackName](message);
      }
    });
  }

  // Internal websocket callbacks

  onOpen(e) {
    if (this.websocketReconnectTimer) {
      clearTimeout(this.websocketReconnectTimer);
    }

    this.notifyListeners(CALLBACKS.WEBSOCKET_CONNECTED);
  }

  onClose(e) {
    // connection closed, discard old websocket and create a new one in 5s
    this.websocket = null;
    this.notifyListeners(CALLBACKS.WEBSOCKET_DISCONNECTED);
    this.handleNetworkingError('Websocket closed.');
    this.websocketReconnectTimer = setTimeout(this.setupWebsocket, TIMER_WEBSOCKET_RECONNECT);
  }

  // On ws error just close the socket and let it re-connect again for now.
  onError(e) {
    this.handleNetworkingError(`Socket error: ${JSON.parse(e)}`);
    this.websocket.close();
  }

  /*
  onMessage is fired when an inbound object comes across the websocket.
  We pass these messages along to listeners by firing the
  RAW_WEBSOCKET_MESSAGE_RECEIVED type callback.
  If the message is of type `PING` we send a `PONG` back and do not
  pass it along to listeners.
  */
  onMessage(e) {
    const model = JSON.parse(e.data);
    
    // Send PONGs
    if (model.type === SOCKET_MESSAGE_TYPES.PING) {
      this.sendPong();
      return;
    }

    // Notify any of the listeners via the raw socket message callback.
    this.notifyListeners(CALLBACKS.RAW_WEBSOCKET_MESSAGE_RECEIVED, model);
  }

  // Reply to a PING as a keep alive.
  sendPong() {
    try {
      const pong = { type: SOCKET_MESSAGE_TYPES.PONG };
      this.websocket.send(JSON.stringify(pong));
    } catch (e) {
      console.log('PONG error:', e);
    }
  }

  handleNetworkingError(error) {
    console.error(`Websocket Error: ${error}`)
  };
}

export default Websocket;