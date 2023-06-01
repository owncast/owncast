import { MessageType, SocketEvent } from '../interfaces/socket-events';

export interface SocketMessage {
  type: MessageType;
  data: any;
}

export default class WebsocketService {
  websocket: WebSocket;

  accessToken: string;

  host: string;

  path: string;

  websocketReconnectTimer: ReturnType<typeof setTimeout>;

  isShutdown = false;

  backOff = 0;

  handleMessage?: (message: SocketEvent) => void;

  socketConnected: () => void;

  socketDisconnected: () => void;

  constructor(accessToken, path, host) {
    this.accessToken = accessToken;
    this.path = path;
    this.websocketReconnectTimer = null;
    this.isShutdown = false;
    this.host = host;

    this.createAndConnect = this.createAndConnect.bind(this);
    this.shutdown = this.shutdown.bind(this);

    this.createAndConnect();
  }

  createAndConnect() {
    if (!this.host) {
      return;
    }

    if (this.isShutdown) {
      return;
    }

    const url = new URL(this.host);
    url.protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    url.pathname = '/ws';
    url.port = window.location.port === '3000' ? '8080' : window.location.port;
    url.searchParams.append('accessToken', this.accessToken);

    const ws = new WebSocket(url.toString());
    ws.onopen = this.onOpen.bind(this);
    ws.onerror = this.onError.bind(this);
    ws.onmessage = this.onMessage.bind(this);

    this.websocket = ws;
  }

  onOpen() {
    if (this.websocketReconnectTimer) {
      clearTimeout(this.websocketReconnectTimer);
    }
    this.socketConnected();
    this.backOff = 0;
  }

  // On ws error just close the socket and let it re-connect again for now.
  onError() {
    handleNetworkingError();
    this.socketDisconnected();
    this.websocket.close();
    if (!this.isShutdown) {
      this.scheduleReconnect();
    }
  }

  scheduleReconnect() {
    if (this.isShutdown) {
      return;
    }

    if (this.websocketReconnectTimer) {
      clearTimeout(this.websocketReconnectTimer);
    }
    this.websocketReconnectTimer = setTimeout(
      this.createAndConnect,
      Math.min(this.backOff, 10_000),
    );
    this.backOff += 1000;
  }

  shutdown() {
    this.isShutdown = true;
    this.websocket.close();
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
    let socketEvent: SocketEvent;

    // eslint-disable-next-line no-plusplus
    for (let i = 0; i < messages.length; i++) {
      try {
        socketEvent = JSON.parse(messages[i]);
        if (this.handleMessage) {
          this.handleMessage(socketEvent);
        }
      } catch (err) {
        console.error(err, err.data);
        return;
      }

      if (!socketEvent.type) {
        console.error('No type provided', socketEvent);
        return;
      }

      // Send PONGs
      if (socketEvent.type === MessageType.PING) {
        this.sendPong();
        return;
      }
    }
  }

  isConnected(): boolean {
    return this.websocket?.readyState === this.websocket?.OPEN;
  }

  // Outbound: Other components can pass an object to `send`.
  send(socketEvent: any) {
    // Sanity check that what we're sending is a valid type.
    if (!socketEvent.type || !MessageType[socketEvent.type]) {
      console.warn(`Outbound message: Unknown socket message type: "${socketEvent.type}" sent.`);
    }

    const messageJSON = JSON.stringify(socketEvent);
    this.websocket.send(messageJSON);
  }

  // Reply to a PING as a keep alive.
  sendPong() {
    const pong = { type: MessageType.PONG };
    this.send(pong);
  }
}

function handleNetworkingError() {
  console.error(
    `Chat has been disconnected and is likely not working for you. It's possible you were removed from chat. If this is a server configuration issue, visit troubleshooting steps to resolve. https://owncast.online/docs/troubleshooting/#chat-is-disabled`,
  );
}
