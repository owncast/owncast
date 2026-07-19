import { render, waitFor, act } from '@testing-library/react';
import { Provider, createStore } from 'jotai';
import {
  ClientConfigStore,
  clientConfigStateAtom,
  isClientConfigLoadedAtom,
  serverStatusState,
  accessTokenAtom,
  currentUserAtom,
  chatMessagesAtom,
  visibleChatMessagesSelector,
  appStateAtom,
  websocketServiceAtom,
  fatalErrorStateAtom,
} from '../components/stores/ClientConfigStore';
import { ClientConfigServiceContext } from '../services/client-config-service';
import { ChatServiceContext } from '../services/chat-service';
import { ServerStatusServiceContext } from '../services/status-service';
import { makeEmptyClientConfig } from '../interfaces/client-config.model';
import { makeEmptyServerStatus } from '../interfaces/server-status.model';
import { MessageType, SocketEvent } from '../interfaces/socket-events';

type MockSocket = {
  accessToken: string;
  shutdown: jest.Mock;
  handleMessage?: (message: SocketEvent) => void;
  socketConnected?: () => void;
  socketDisconnected?: () => void;
};

// Capture every websocket the store constructs so tests can drive inbound
// events without a real socket.
const mockSockets: MockSocket[] = [];
jest.mock('../services/websocket-service', () => ({
  __esModule: true,
  default: jest.fn().mockImplementation(function fakeSocket(this: MockSocket, accessToken: string) {
    this.accessToken = accessToken;
    this.shutdown = jest.fn();
    mockSockets.push(this);
  }),
}));

const hydrationWindow = window as unknown as {
  configHydration?: string;
  statusHydration?: string;
};

const testConfig = { ...makeEmptyClientConfig(), name: 'test server', summary: 'from api' };
const onlineStatus = { ...makeEmptyServerStatus(), online: true, viewerCount: 3 };

type Services = {
  configService: { getConfig: jest.Mock };
  statusService: { getStatus: jest.Mock };
  chatService: { registerUser: jest.Mock; getChatHistory: jest.Mock };
};

const makeServices = (): Services => ({
  configService: { getConfig: jest.fn().mockResolvedValue(testConfig) },
  statusService: { getStatus: jest.fn().mockResolvedValue(onlineStatus) },
  chatService: {
    registerUser: jest.fn().mockResolvedValue({
      id: 'u1',
      accessToken: 'new-token',
      displayName: 'rando',
      displayColor: 2,
    }),
    getChatHistory: jest.fn().mockResolvedValue([]),
  },
});

const renderStore = (services: Services) => {
  const store = createStore();
  render(
    <Provider store={store}>
      <ClientConfigServiceContext.Provider value={services.configService}>
        <ChatServiceContext.Provider value={services.chatService}>
          <ServerStatusServiceContext.Provider value={services.statusService}>
            <ClientConfigStore />
          </ServerStatusServiceContext.Provider>
        </ChatServiceContext.Provider>
      </ClientConfigServiceContext.Provider>
    </Provider>,
  );
  return store;
};

const chatEvent = (id: string, body: string): SocketEvent =>
  ({
    type: MessageType.CHAT,
    id,
    timestamp: new Date(),
    body,
    user: { id: 'u2', displayName: 'someone', displayColor: 1 },
  }) as unknown as SocketEvent;

describe('ClientConfigStore', () => {
  beforeEach(() => {
    localStorage.clear();
    mockSockets.length = 0;
    delete hydrationWindow.configHydration;
    delete hydrationWindow.statusHydration;
    jest.spyOn(console, 'error').mockImplementation(() => {});
  });

  afterEach(() => {
    jest.restoreAllMocks();
  });

  it('fetches config and status from the API after mount and applies them', async () => {
    const services = makeServices();
    const store = renderStore(services);

    await waitFor(() => expect(store.get(isClientConfigLoadedAtom)).toBe(true));
    expect(services.configService.getConfig).toHaveBeenCalled();
    expect(store.get(clientConfigStateAtom).summary).toBe('from api');

    await waitFor(() => expect(store.get(serverStatusState).online).toBe(true));
    // The app state machine must leave "loading" and reflect the live stream.
    await waitFor(() => expect(store.get(appStateAtom).videoAvailable).toBe(true));
  });

  it('applies server-injected hydration data from the mount effect without hitting the API', async () => {
    // A saved token keeps the registration flow from fetching config too.
    localStorage.setItem('accessToken', 'saved-token');
    hydrationWindow.configHydration = JSON.stringify({ ...testConfig, summary: 'from hydration' });
    hydrationWindow.statusHydration = JSON.stringify(onlineStatus);
    const services = makeServices();
    const store = renderStore(services);

    await waitFor(() => expect(store.get(isClientConfigLoadedAtom)).toBe(true));
    expect(store.get(clientConfigStateAtom).summary).toBe('from hydration');
    expect(store.get(serverStatusState).online).toBe(true);
    await waitFor(() => expect(store.get(appStateAtom).videoAvailable).toBe(true));
    expect(services.configService.getConfig).not.toHaveBeenCalled();
    expect(services.statusService.getStatus).not.toHaveBeenCalled();
  });

  it('falls back to the API when hydration data is malformed', async () => {
    localStorage.setItem('accessToken', 'saved-token');
    hydrationWindow.configHydration = '{this is not json';
    const services = makeServices();
    const store = renderStore(services);

    await waitFor(() => expect(store.get(isClientConfigLoadedAtom)).toBe(true));
    expect(services.configService.getConfig).toHaveBeenCalled();
    expect(store.get(clientConfigStateAtom).summary).toBe('from api');
  });

  it('registers a new chat user and stores the access token', async () => {
    const services = makeServices();
    const store = renderStore(services);

    await waitFor(() => expect(store.get(accessTokenAtom)).toBe('new-token'));
    expect(services.chatService.registerUser).toHaveBeenCalled();
    expect(localStorage.getItem('accessToken')).toBe('new-token');
    expect(store.get(currentUserAtom).displayName).toBe('rando');
  });

  it('refreshes personalized page content after first-time registration when chat is disabled', async () => {
    // Server-rendered hydration delivers the anonymous visitor's config, so
    // the only API config fetch is the post-registration refresh that picks
    // up viewer-personalized plugin content (e.g. a viewer gate's page
    // content keyed to the new access token).
    hydrationWindow.configHydration = JSON.stringify({
      ...testConfig,
      chatDisabled: true,
      extraPageContent: '<p>visitor</p>',
    });
    hydrationWindow.statusHydration = JSON.stringify(makeEmptyServerStatus());
    const services = makeServices();
    services.configService.getConfig.mockResolvedValue({
      ...testConfig,
      chatDisabled: true,
      extraPageContent: '<p>registered viewer</p>',
    });
    const store = renderStore(services);

    await waitFor(() => {
      expect(store.get(clientConfigStateAtom).extraPageContent).toBe('<p>registered viewer</p>');
    });
    expect(services.chatService.registerUser).toHaveBeenCalledTimes(1);
    expect(services.configService.getConfig).toHaveBeenCalledTimes(1);
    expect(localStorage.getItem('accessToken')).toBe('new-token');
  });

  it('reuses a saved access token instead of re-registering', async () => {
    localStorage.setItem('accessToken', 'saved-token');
    const services = makeServices();
    const store = renderStore(services);

    await waitFor(() => expect(store.get(accessTokenAtom)).toBe('saved-token'));
    expect(services.chatService.registerUser).not.toHaveBeenCalled();
  });

  it('starts the websocket with the access token once config is loaded', async () => {
    const services = makeServices();
    const store = renderStore(services);

    await waitFor(() => expect(mockSockets).toHaveLength(1));
    const socket = mockSockets[0];
    expect(socket.accessToken).toBe('new-token');
    expect(store.get(websocketServiceAtom)).toBe(socket);
    expect(typeof socket.handleMessage).toBe('function');
  });

  it('appends incoming chat messages to chat state', async () => {
    const services = makeServices();
    const store = renderStore(services);
    await waitFor(() => expect(mockSockets).toHaveLength(1));

    const message = chatEvent('msg-1', 'hello');
    act(() => mockSockets[0].handleMessage(message));

    expect(store.get(chatMessagesAtom)).toHaveLength(1);
    expect(store.get(chatMessagesAtom)[0]).toBe(message);
  });

  it('hides and re-shows messages on visibility updates', async () => {
    const services = makeServices();
    const store = renderStore(services);
    await waitFor(() => expect(mockSockets).toHaveLength(1));
    const socket = mockSockets[0];

    act(() => {
      socket.handleMessage(chatEvent('msg-1', 'one'));
      socket.handleMessage(chatEvent('msg-2', 'two'));
    });

    const visibility = (visible: boolean) =>
      ({
        type: MessageType.VISIBILITY_UPDATE,
        id: 'vis-1',
        timestamp: new Date(),
        visible,
        ids: ['msg-1'],
      }) as unknown as SocketEvent;

    act(() => socket.handleMessage(visibility(false)));
    expect(store.get(visibleChatMessagesSelector).map(m => m.id)).toEqual(['msg-2']);

    act(() => socket.handleMessage(visibility(true)));
    expect(store.get(visibleChatMessagesSelector).map(m => m.id)).toEqual(['msg-1', 'msg-2']);
  });

  it('re-registers when the server rejects the saved token', async () => {
    localStorage.setItem('accessToken', 'stale-token');
    const services = makeServices();
    const store = renderStore(services);
    await waitFor(() => expect(mockSockets).toHaveLength(1));
    const socket = mockSockets[0];
    expect(socket.accessToken).toBe('stale-token');

    act(() =>
      socket.handleMessage({
        type: MessageType.ERROR_NEEDS_REGISTRATION,
        id: 'err-1',
        timestamp: new Date(),
      }),
    );

    expect(socket.shutdown).toHaveBeenCalled();
    await waitFor(() => expect(store.get(accessTokenAtom)).toBe('new-token'));
    expect(localStorage.getItem('accessToken')).toBe('new-token');
    expect(services.chatService.registerUser).toHaveBeenCalled();
  });

  it('surfaces a fatal error when the server is unreachable', async () => {
    const services = makeServices();
    services.configService.getConfig.mockRejectedValue(new Error('down'));
    services.statusService.getStatus.mockRejectedValue(new Error('down'));
    const store = renderStore(services);

    await waitFor(() => expect(store.get(fatalErrorStateAtom)).not.toBeNull());
    expect(store.get(fatalErrorStateAtom).title).toBe('Unable to reach Owncast server');
  });
});
