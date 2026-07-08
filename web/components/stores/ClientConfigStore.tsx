import { FC, useContext, useEffect } from 'react';
import { atom, useAtom, useSetAtom } from 'jotai';
import { useMachine } from '@xstate/react';
import { makeEmptyClientConfig, ClientConfig } from '../../interfaces/client-config.model';
import { ClientConfigServiceContext } from '../../services/client-config-service';
import { ChatServiceContext } from '../../services/chat-service';
import WebsocketService from '../../services/websocket-service';
import { ChatMessage } from '../../interfaces/chat-message.model';
import { CurrentUser } from '../../interfaces/current-user';
import { ServerStatus, makeEmptyServerStatus } from '../../interfaces/server-status.model';
import appStateModel, {
  AppStateEvent,
  AppStateOptions,
  makeEmptyAppState,
} from './application-state';
import { setLocalStorage, getLocalStorage } from '../../utils/localStorage';
import {
  ConnectedClientInfoEvent,
  MessageType,
  ChatEvent,
  NameChangeEvent,
  MessageVisibilityEvent,
  SocketEvent,
} from '../../interfaces/socket-events';
import { mergeMeta } from '../../utils/helpers';
import { handleConnectedClientInfoMessage } from './eventhandlers/connected-client-info-handler';
import { ServerStatusServiceContext } from '../../services/status-service';
import { handleNameChangeEvent } from './eventhandlers/handleNameChangeEvent';
import { DisplayableError } from '../../types/displayable-error';

const SERVER_STATUS_POLL_DURATION = 5000;

// Helper to safely parse hydration data injected into the page by the Owncast
// server (see ServerRenderedHydration). Read from a mount effect, never during
// render: the statically-exported HTML was built with empty config/status, so
// initializing state from this data during hydration makes React's first
// render differ from the server HTML (React errors #418/#423/#425).
// Returns both the config and whether parsing succeeded.
const getInitialConfig = (): { config: ClientConfig; success: boolean } => {
  if (typeof window !== 'undefined' && (window as any).configHydration) {
    try {
      const parsed = JSON.parse((window as any).configHydration);
      if (parsed) {
        return { config: parsed, success: true };
      }
    } catch (e) {
      console.error('Error parsing config hydration during init', e);
    }
  }
  return { config: makeEmptyClientConfig(), success: false };
};

const getInitialStatus = (): { status: ServerStatus; success: boolean } => {
  if (typeof window !== 'undefined' && (window as any).statusHydration) {
    try {
      const parsed = JSON.parse((window as any).statusHydration);
      if (parsed) {
        return { status: parsed, success: true };
      }
    } catch (e) {
      console.error('Error parsing status hydration during init', e);
    }
  }
  return { status: makeEmptyServerStatus(), success: false };
};

const ACCESS_TOKEN_KEY = 'accessToken';

let serverStatusRefreshPoll: ReturnType<typeof setInterval>;
let hasBeenModeratorNotified = false;
let hasWebsocketDisconnected = false;

const serverConnectivityError = `Cannot connect to the Owncast service. Please check your internet connection and verify this Owncast server is running.`;

// Server status is what gets updated such as viewer count, durations,
// stream title, online/offline state, etc.
// Starts empty to match the statically-exported HTML; hydration data and API
// polls fill it in after mount.
export const serverStatusState = atom<ServerStatus>(makeEmptyServerStatus());

// The config that comes from the API.
// Starts empty to match the statically-exported HTML; hydration data or the
// API fills it in after mount.
export const clientConfigStateAtom = atom(makeEmptyClientConfig());

// Whether the client config has been populated, via hydration or the API.
// Consumers that must not act on default config values (like the player,
// whose video.js options are init-only) gate on this.
export const isClientConfigLoadedAtom = atom<boolean>(false);

// The `null as T` casts below matter: this project compiles without
// strictNullChecks, so a bare null would match jotai's read-function
// overload and produce a read-only atom instead of a writable one.
export const accessTokenAtom = atom<string>(null as string);

export const currentUserAtom = atom<CurrentUser>(null as CurrentUser);

export const chatMessagesAtom = atom<ChatMessage[]>([]);

export const chatAuthenticatedAtom = atom<boolean>(false);

// Stores chat input draft to preserve text across mobile/desktop mode switches
export const chatInputDraftAtom = atom<string>('');

export const websocketServiceAtom = atom<WebsocketService>(null as WebsocketService);

// Starts in the "loading" app state to match the statically-exported HTML;
// the hydration mount effect in ClientConfigStore transitions it to
// online/offline immediately after mount.
export const appStateAtom = atom<AppStateOptions>(makeEmptyAppState());

export const isMobileAtom = atom<boolean | undefined>(undefined as boolean | undefined);

export const isVideoPlayingAtom = atom<boolean>(false);

export const fatalErrorStateAtom = atom<DisplayableError>(null as DisplayableError);

export const clockSkewAtom = atom<Number>(0.0);

const removedMessageIdsAtom = atom<string[]>([]);

export const isChatAvailableSelector = atom(get => {
  const state: AppStateOptions = get(appStateAtom);
  const accessToken: string = get(accessTokenAtom);
  return Boolean(accessToken && state.chatAvailable && !hasWebsocketDisconnected);
});

// The requested state of chat in the UI
export enum ChatState {
  VISIBLE, // Chat is open (the default state when the stream is online)
  HIDDEN, // Chat is hidden
  POPPED_OUT, // Chat is playing in a popout window
  EMBEDDED, // This window is opened at /embed/chat/readwrite/
}

export const chatStateAtom = atom<ChatState>(
  (() => {
    // XXX Somehow, `window` is undefined here, even though this runs in client
    const window = globalThis;
    return window?.location?.pathname === '/embed/chat/readwrite/'
      ? ChatState.EMBEDDED
      : ChatState.VISIBLE;
  })(),
);

// We display in an "online/live" state as long as video is actively playing.
// Even during the time where technically the server has said it's no longer
// live, however the last few seconds of video playback is still taking place.
export const isOnlineSelector = atom(get => {
  const state: AppStateOptions = get(appStateAtom);
  const isVideoPlaying: boolean = get(isVideoPlayingAtom);
  return state.videoAvailable || isVideoPlaying;
});

export const visibleChatMessagesSelector = atom<ChatMessage[]>(get => {
  const messages: ChatMessage[] = get(chatMessagesAtom);
  const removedIds: string[] = get(removedMessageIdsAtom);
  return messages.filter(message => !removedIds.includes(message.id));
});

export const ClientConfigStore: FC = () => {
  const ClientConfigService = useContext(ClientConfigServiceContext);
  const ChatService = useContext(ChatServiceContext);
  const ServerStatusService = useContext(ServerStatusServiceContext);

  const [appState, appStateSend, appStateService] = useMachine(appStateModel);
  const [currentUser, setCurrentUser] = useAtom(currentUserAtom);
  const setChatAuthenticated = useSetAtom(chatAuthenticatedAtom);
  const [clientConfig, setClientConfig] = useAtom(clientConfigStateAtom);
  const setServerStatus = useSetAtom(serverStatusState);
  const setClockSkew = useSetAtom(clockSkewAtom);
  const setChatMessages = useSetAtom(chatMessagesAtom);
  const [accessToken, setAccessToken] = useAtom(accessTokenAtom);
  const setAppState = useSetAtom(appStateAtom);
  const setGlobalFatalErrorMessage = useSetAtom(fatalErrorStateAtom);
  const setWebsocketService = useSetAtom(websocketServiceAtom);
  const setHiddenMessageIds = useSetAtom(removedMessageIdsAtom);
  const [hasLoadedConfig, setHasLoadedConfig] = useAtom(isClientConfigLoadedAtom);

  let ws: WebsocketService;

  const setGlobalFatalError = (title: string, message: string) => {
    setGlobalFatalErrorMessage({
      title,
      message,
    });
  };
  const sendEvent = (events: string[]) => {
    // console.debug('---- sending event:', event);
    appStateSend(events);
  };

  const handleStatusChange = (status: ServerStatus) => {
    if (appState.matches('loading')) {
      const events = [AppStateEvent.Loaded];
      if (status.online) {
        events.push(AppStateEvent.Online);
      } else {
        events.push(AppStateEvent.Offline);
      }
      sendEvent(events);
      return;
    }

    if (status.online && appState.matches('ready')) {
      sendEvent([AppStateEvent.Online]);
    } else if (!status.online && !appState.matches('ready.offline')) {
      sendEvent([AppStateEvent.Offline]);
    }
  };

  const updateClientConfig = async () => {
    try {
      const config = await ClientConfigService.getConfig();
      setClientConfig(config);
      setGlobalFatalErrorMessage(null);
      setHasLoadedConfig(true);
    } catch (error) {
      setGlobalFatalError('Unable to reach Owncast server', serverConnectivityError);
      console.error(`ClientConfigService -> getConfig() ERROR: \n`, error);
    }
  };

  const updateServerStatus = async () => {
    try {
      const status = await ServerStatusService.getStatus();
      handleStatusChange(status);
      setServerStatus(status);

      const { serverTime } = status;

      const clockSkew = new Date(serverTime).getTime() - Date.now();
      setClockSkew(clockSkew);

      setGlobalFatalErrorMessage(null);
    } catch (error) {
      sendEvent([AppStateEvent.Fail]);
      setGlobalFatalError('Unable to reach Owncast server', serverConnectivityError);
      console.error(`serverStatusState -> getStatus() ERROR: \n`, error);
    }
  };

  const handleUserRegistration = async (optionalDisplayName?: string) => {
    const savedAccessToken = getLocalStorage(ACCESS_TOKEN_KEY);
    if (savedAccessToken) {
      setAccessToken(savedAccessToken);

      return;
    }

    try {
      sendEvent([AppStateEvent.NeedsRegister]);
      const response = await ChatService.registerUser(optionalDisplayName);
      const { accessToken: newAccessToken, displayName: newDisplayName, displayColor } = response;
      if (!newAccessToken) {
        return;
      }

      setCurrentUser({
        ...currentUser,
        displayName: newDisplayName,
        displayColor,
      });
      setAccessToken(newAccessToken);
      setLocalStorage(ACCESS_TOKEN_KEY, newAccessToken);
    } catch (e) {
      sendEvent([AppStateEvent.Fail]);
      console.error(`ChatService -> registerUser() ERROR: \n${e}`);
    }
  };

  const resetAndReAuth = () => {
    setLocalStorage(ACCESS_TOKEN_KEY, '');
    setAccessToken(null);
    ws?.shutdown();
    handleUserRegistration();
  };

  const handleMessageVisibilityChange = (message: MessageVisibilityEvent) => {
    const { ids, visible } = message;
    if (visible) {
      setHiddenMessageIds(currentState => currentState.filter(id => !ids.includes(id)));
    } else {
      setHiddenMessageIds(currentState => [...currentState, ...ids]);
    }
  };

  const handleSocketDisconnect = () => {
    hasWebsocketDisconnected = true;
  };

  const handleSocketConnected = () => {
    hasWebsocketDisconnected = false;
  };

  const handleMessage = (message: SocketEvent) => {
    switch (message.type) {
      case MessageType.ERROR_NEEDS_REGISTRATION:
        resetAndReAuth();
        break;
      case MessageType.CONNECTED_USER_INFO:
        handleConnectedClientInfoMessage(
          message as ConnectedClientInfoEvent,
          setChatAuthenticated,
          setCurrentUser,
        );
        if (message as ChatEvent) {
          const m = new ChatEvent(message);
          if (!hasBeenModeratorNotified && m.user?.isModerator) {
            setChatMessages(currentState => [...currentState, message as ChatEvent]);
            hasBeenModeratorNotified = true;
          }
        }

        break;
      case MessageType.CHAT:
        setChatMessages(currentState => [...currentState, message as ChatEvent]);
        break;
      case MessageType.NAME_CHANGE:
        handleNameChangeEvent(message as NameChangeEvent, setChatMessages, setCurrentUser);
        break;
      case MessageType.USER_JOINED:
        setChatMessages(currentState => [...currentState, message as ChatEvent]);
        break;
      case MessageType.USER_PARTED:
        setChatMessages(currentState => [...currentState, message as ChatEvent]);
        break;
      case MessageType.SYSTEM:
        setChatMessages(currentState => [...currentState, message as ChatEvent]);
        break;
      case MessageType.CHAT_ACTION:
        setChatMessages(currentState => [...currentState, message as ChatEvent]);
        break;
      case MessageType.FEDIVERSE_ENGAGEMENT_FOLLOW:
        setChatMessages(currentState => [...currentState, message as unknown as ChatMessage]);
        break;
      case MessageType.FEDIVERSE_ENGAGEMENT_LIKE:
        setChatMessages(currentState => [...currentState, message as unknown as ChatMessage]);
        break;
      case MessageType.FEDIVERSE_ENGAGEMENT_REPOST:
        setChatMessages(currentState => [...currentState, message as unknown as ChatMessage]);
        break;
      case MessageType.VISIBILITY_UPDATE:
        handleMessageVisibilityChange(message as MessageVisibilityEvent);
        break;
      case MessageType.ERROR_USER_DISABLED:
        console.log('User has been disabled');
        sendEvent([AppStateEvent.ChatUserDisabled]);
        break;
      default:
        console.error('Unknown socket message type: ', message.type);
    }
  };

  const getChatHistory = async () => {
    try {
      const messages = await ChatService.getChatHistory(accessToken);
      if (messages) {
        setChatMessages(currentState => [...currentState, ...messages]);
      }
    } catch (error) {
      console.error(`ChatService -> getChatHistory() ERROR: \n${error}`);
    }
  };

  const startChat = async () => {
    try {
      if (ws) {
        ws?.shutdown();
        setWebsocketService(null);
        ws = null;
      }

      const { socketHostOverride } = clientConfig;

      // Get a copy of the browser location without #fragments.
      const location = window.location.origin + window.location.pathname;
      const host = socketHostOverride || location;

      ws = new WebsocketService(accessToken, '/ws', host);
      ws.handleMessage = handleMessage;
      ws.socketDisconnected = handleSocketDisconnect;
      ws.socketConnected = handleSocketConnected;
      setWebsocketService(ws);
    } catch (error) {
      console.error(`ChatService -> startChat() ERROR: \n${error}`);
      sendEvent([AppStateEvent.ChatUserDisabled]);
    }
  };

  // Apply the server-injected hydration data (window.configHydration /
  // window.statusHydration) after mount. This fills in real config and
  // status without waiting on an API round trip, while keeping React's
  // hydration pass identical to the statically-exported HTML.
  useEffect(() => {
    const { config, success: hasHydratedConfig } = getInitialConfig();
    const { status, success: hasHydratedStatus } = getInitialStatus();

    if (hasHydratedConfig) {
      setClientConfig(config);
      setHasLoadedConfig(true);
    } else {
      updateClientConfig();
    }

    handleUserRegistration();

    if (hasHydratedStatus) {
      handleStatusChange(status);
      setServerStatus(status);
      if (status.serverTime) {
        const clockSkew = new Date(status.serverTime).getTime() - Date.now();
        setClockSkew(clockSkew);
      }
    } else {
      updateServerStatus();
    }

    clearInterval(serverStatusRefreshPoll);
    serverStatusRefreshPoll = setInterval(() => {
      updateServerStatus();
    }, SERVER_STATUS_POLL_DURATION);

    return () => {
      clearInterval(serverStatusRefreshPoll);
    };
  }, []);

  useEffect(() => {
    if (clientConfig.chatDisabled) {
      return;
    }

    if (!accessToken) {
      return;
    }

    if (!hasLoadedConfig) {
      return;
    }

    if (ws) {
      return;
    }

    startChat();
  }, [hasLoadedConfig, accessToken]);

  useEffect(() => {
    if (accessToken) {
      getChatHistory();
    }
  }, [accessToken]);

  useEffect(() => {
    appStateService.onTransition(state => {
      const metadata = mergeMeta(state.meta) as AppStateOptions;

      // console.debug('--- APP STATE: ', state.value);
      // console.debug('--- APP META: ', metadata);

      setAppState(metadata);
    });
  }, []);

  return null;
};
