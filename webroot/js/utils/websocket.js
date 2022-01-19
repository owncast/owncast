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
  SYSTEM: 'SYSTEM',
  USER_JOINED: 'USER_JOINED',
  CHAT_ACTION: 'CHAT_ACTION',
  FEDIVERSE_ENGAGEMENT_FOLLOW: 'FEDIVERSE_ENGAGEMENT_FOLLOW',
  FEDIVERSE_ENGAGEMENT_LIKE: 'FEDIVERSE_ENGAGEMENT_LIKE',
  FEDIVERSE_ENGAGEMENT_REPOST: 'FEDIVERSE_ENGAGEMENT_REPOST',
  CONNECTED_USER_INFO: 'CONNECTED_USER_INFO',
  ERROR_USER_DISABLED: 'ERROR_USER_DISABLED',
  ERROR_NEEDS_REGISTRATION: 'ERROR_NEEDS_REGISTRATION',
  ERROR_MAX_CONNECTIONS_EXCEEDED: 'ERROR_MAX_CONNECTIONS_EXCEEDED',
  VISIBILITY_UPDATE: 'VISIBILITY-UPDATE',
};

export const CALLBACKS = {
  RAW_WEBSOCKET_MESSAGE_RECEIVED: 'rawWebsocketMessageReceived',
  WEBSOCKET_CONNECTED: 'websocketConnected',
};

const TIMER_WEBSOCKET_RECONNECT = 5000; // ms

export default class Websocket {
  constructor(accessToken, path) {
    this.websocket = null;
    this.path = path;
    this.websocketReconnectTimer = null;
    this.accessToken = accessToken;

    this.websocketConnectedListeners = [];
    this.websocketDisconnectListeners = [];
    this.rawMessageListeners = [];

    this.send = this.send.bind(this);
    this.createAndConnect = this.createAndConnect.bind(this);
    this.scheduleReconnect = this.scheduleReconnect.bind(this);
    this.shutdown = this.shutdown.bind(this);

    this.isShutdown = false;

    this.createAndConnect();
  }

  createAndConnect() {
    const url = new URL(this.path);
    url.searchParams.append('accessToken', this.accessToken);

    const ws = new WebSocket(url.toString());
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

  shutdown() {
    this.isShutdown = true;
    this.websocket.close();
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
    if (!this.isShutdown) {
      this.scheduleReconnect();
    }
  }

  // On ws error just close the socket and let it re-connect again for now.
  onError(e) {
    this.handleNetworkingError(`Socket error: ${JSON.parse(e)}`);
    this.websocket.close();
    if (!this.isShutdown) {
      this.scheduleReconnect();
    }
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
    // Optimization where multiple events can be sent within a
    // single websocket message. So split them if needed.
    var messages = e.data.split('\n');
    for (var i = 0; i < messages.length; i++) {
      try {
        var model = JSON.parse(messages[i]);
      } catch (e) {
        console.error(e, e.data);
        return;
      }

      if (!model.type) {
        console.error('No type provided', model);
        return;
      }

      // Send PONGs
      if (model.type === SOCKET_MESSAGE_TYPES.PING) {
        this.sendPong();
        return;
      }

      // Notify any of the listeners via the raw socket message callback.
      this.notifyRawMessageListeners(model);
    }
  }

  // Reply to a PING as a keep alive.
  sendPong() {
    const pong = { type: SOCKET_MESSAGE_TYPES.PONG };
    this.send(pong);
  }

  handleNetworkingError(error) {
    console.error(
      `Chat has been disconnected and is likely not working for you. It's possible you were removed from chat. If this is a server configuration issue, visit troubleshooting steps to resolve. https://owncast.online/docs/troubleshooting/#chat-is-disabled: ${error}`
    );
  }
}
