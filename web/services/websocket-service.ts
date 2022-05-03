import { message } from 'antd';
import { SocketMessageType } from '../interfaces/socket-events';

interface SocketMessage {
  type: SocketMessageType;
  data: any;
}

export default class WebsocketService {
  websocket: WebSocket;

  accessToken: string;

  path: string;

  websocketReconnectTimer: ReturnType<typeof setTimeout>;

  handleMessage?: (message: SocketMessage) => void;

  constructor(accessToken, path) {
    this.accessToken = accessToken;
    this.path = path;
    // this.websocketReconnectTimer = null;
    // this.accessToken = accessToken;

    // this.websocketConnectedListeners = [];
    // this.websocketDisconnectListeners = [];
    // this.rawMessageListeners = [];

    // this.send = this.send.bind(this);
    // this.createAndConnect = this.createAndConnect.bind(this);
    // this.scheduleReconnect = this.scheduleReconnect.bind(this);
    // this.shutdown = this.shutdown.bind(this);

    // this.isShutdown = false;

    this.createAndConnect();
  }

  createAndConnect() {
    const url = new URL('ws://localhost:8080/ws');
    url.searchParams.append('accessToken', this.accessToken);

    console.log('connecting to ', url.toString());
    const ws = new WebSocket(url.toString());
    ws.onopen = this.onOpen.bind(this);
    // ws.onclose = this.onClose.bind(this);
    ws.onerror = this.onError.bind(this);
    ws.onmessage = this.onMessage.bind(this);

    this.websocket = ws;
  }

  onOpen() {
    if (this.websocketReconnectTimer) {
      clearTimeout(this.websocketReconnectTimer);
    }
  }

  // On ws error just close the socket and let it re-connect again for now.
  onError(e) {
    console.error(e);
    handleNetworkingError(`Socket error: ${e}`);
    this.websocket.close();
    // if (!this.isShutdown) {
    //   this.scheduleReconnect();
    // }
  }

  /*
  onMessage is fired when an inbound object comes across the websocket.
  If the message is of type `PING` we send a `PONG` back and do not
  pass it along to listeners.
  */
  onMessage(e: SocketMessage) {
    // Optimization where multiple events can be sent within a
    // single websocket message. So split them if needed.
    const messages = e.data.split('\n');
    let message: SocketMessage;

    // eslint-disable-next-line no-plusplus
    for (let i = 0; i < messages.length; i++) {
      try {
        message = JSON.parse(messages[i]);
        if (this.handleMessage) {
          this.handleMessage(message);
        }
      } catch (e) {
        console.error(e, e.data);
        return;
      }

      if (!message.type) {
        console.error('No type provided', message);
        return;
      }

      // Send PONGs
      if (message.type === SocketMessageType.PING) {
        this.sendPong();
        return;
      }
    }
  }

  // Outbound: Other components can pass an object to `send`.
  send(message: any) {
    // Sanity check that what we're sending is a valid type.
    if (!message.type || !SocketMessageType[message.type]) {
      console.warn(`Outbound message: Unknown socket message type: "${message.type}" sent.`);
    }

    const messageJSON = JSON.stringify(message);
    this.websocket.send(messageJSON);
  }

  // Reply to a PING as a keep alive.
  sendPong() {
    const pong = { type: SocketMessageType.PONG };
    this.send(pong);
  }
}

function handleNetworkingError(error) {
  console.error(
    `Chat has been disconnected and is likely not working for you. It's possible you were removed from chat. If this is a server configuration issue, visit troubleshooting steps to resolve. https://owncast.online/docs/troubleshooting/#chat-is-disabled: ${error}`,
  );
}
