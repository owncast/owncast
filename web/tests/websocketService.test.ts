import WebsocketService from '../services/websocket-service';
import { MessageType, SocketEvent } from '../interfaces/socket-events';

// A stand-in for the browser WebSocket: records constructed instances and
// outbound frames, and lets tests fire the event handlers the service wires.
class FakeWebSocket {
  static instances: FakeWebSocket[] = [];

  url: string;

  readyState = 0;

  OPEN = 1;

  sent: string[] = [];

  onopen: () => void;

  onclose: () => void;

  onmessage: (e: { data: string }) => void;

  constructor(url: string) {
    this.url = url;
    FakeWebSocket.instances.push(this);
  }

  send(data: string) {
    this.sent.push(data);
  }

  close() {
    this.readyState = 3;
  }
}

const lastSocket = () => FakeWebSocket.instances[FakeWebSocket.instances.length - 1];

const connect = () => {
  const service = new WebsocketService('test-token', '/ws', 'http://localhost');
  service.handleMessage = jest.fn();
  service.socketConnected = jest.fn();
  service.socketDisconnected = jest.fn();
  return service;
};

describe('WebsocketService', () => {
  beforeEach(() => {
    // jsdom's WebSocket would open real connections; swap in the fake.
    globalThis.WebSocket = FakeWebSocket as unknown as typeof WebSocket;
    FakeWebSocket.instances = [];
    jest.spyOn(console, 'error').mockImplementation(() => {});
    jest.spyOn(console, 'warn').mockImplementation(() => {});
  });

  afterEach(() => {
    jest.useRealTimers();
    jest.restoreAllMocks();
  });

  it('connects to /ws on the given host with the access token', () => {
    connect();

    expect(FakeWebSocket.instances).toHaveLength(1);
    const url = new URL(lastSocket().url);
    expect(url.protocol).toBe('ws:');
    expect(url.pathname).toBe('/ws');
    expect(url.searchParams.get('accessToken')).toBe('test-token');
  });

  it('notifies socketConnected and resets the backoff on open', () => {
    const service = connect();
    service.backOff = 3000;

    lastSocket().onopen();

    expect(service.socketConnected).toHaveBeenCalledTimes(1);
    expect(service.backOff).toBe(0);
  });

  it('parses an incoming message and hands it to handleMessage', () => {
    const service = connect();
    const event = { type: MessageType.CHAT, id: 'msg-1', body: 'hello' };

    lastSocket().onmessage({ data: JSON.stringify(event) });

    expect(service.handleMessage).toHaveBeenCalledTimes(1);
    expect(service.handleMessage).toHaveBeenCalledWith(event);
  });

  it('splits newline-batched messages into individual events', () => {
    const service = connect();
    const first = { type: MessageType.CHAT, id: 'msg-1', body: 'one' };
    const second = { type: MessageType.SYSTEM, id: 'msg-2', body: 'two' };

    lastSocket().onmessage({ data: `${JSON.stringify(first)}\n${JSON.stringify(second)}` });

    expect(service.handleMessage).toHaveBeenCalledTimes(2);
    expect(service.handleMessage).toHaveBeenNthCalledWith(1, first);
    expect(service.handleMessage).toHaveBeenNthCalledWith(2, second);
  });

  it('replies to a PING with a PONG', () => {
    connect();
    const socket = lastSocket();

    socket.onmessage({ data: JSON.stringify({ type: MessageType.PING }) });

    expect(socket.sent).toHaveLength(1);
    expect(JSON.parse(socket.sent[0])).toEqual({ type: MessageType.PONG });
  });

  it('survives malformed JSON without throwing or delivering', () => {
    const service = connect();

    expect(() => lastSocket().onmessage({ data: 'not json at all' })).not.toThrow();

    expect(service.handleMessage).not.toHaveBeenCalled();
  });

  it('reconnects after an unexpected close with a growing backoff', () => {
    jest.useFakeTimers();
    const service = connect();

    // First drop reconnects immediately (backoff starts at 0).
    lastSocket().onclose();
    expect(service.socketDisconnected).toHaveBeenCalledTimes(1);
    jest.advanceTimersByTime(0);
    expect(FakeWebSocket.instances).toHaveLength(2);

    // Second drop waits a full second.
    lastSocket().onclose();
    jest.advanceTimersByTime(999);
    expect(FakeWebSocket.instances).toHaveLength(2);
    jest.advanceTimersByTime(1);
    expect(FakeWebSocket.instances).toHaveLength(3);
  });

  it('caps the reconnect delay at ten seconds', () => {
    jest.useFakeTimers();
    const service = connect();
    service.backOff = 60_000;

    lastSocket().onclose();
    jest.advanceTimersByTime(9_999);
    expect(FakeWebSocket.instances).toHaveLength(1);
    jest.advanceTimersByTime(1);
    expect(FakeWebSocket.instances).toHaveLength(2);
  });

  it('does not reconnect or report a disconnect after shutdown', () => {
    jest.useFakeTimers();
    const service = connect();
    const socket = lastSocket();

    service.shutdown();
    expect(socket.readyState).toBe(3);
    // The browser fires the close event after close(); the service must
    // treat it as intentional.
    socket.onclose();

    expect(service.socketDisconnected).not.toHaveBeenCalled();
    jest.runOnlyPendingTimers();
    expect(FakeWebSocket.instances).toHaveLength(1);
  });

  it('serializes outbound messages and warns on unknown types', () => {
    const service = connect();
    const socket = lastSocket();
    const outbound = { type: MessageType.CHAT, body: 'hi' } as unknown as SocketEvent;

    service.send(outbound);
    expect(JSON.parse(socket.sent[0])).toEqual({ type: 'CHAT', body: 'hi' });
    expect(console.warn).not.toHaveBeenCalled();

    service.send({ type: 'BOGUS' });
    expect(console.warn).toHaveBeenCalled();
  });
});
