import { FC, useContext, useEffect, useState } from 'react';
import { atom, selector, useRecoilState, useSetRecoilState, RecoilEnv } from 'recoil';
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
  FediverseEvent,
} from '../../interfaces/socket-events';
import { mergeMeta } from '../../utils/helpers';
import { handleConnectedClientInfoMessage } from './eventhandlers/connected-client-info-handler';
import { ServerStatusServiceContext } from '../../services/status-service';
import { handleNameChangeEvent } from './eventhandlers/handleNameChangeEvent';
import { DisplayableError } from '../../types/displayable-error';

RecoilEnv.RECOIL_DUPLICATE_ATOM_KEY_CHECKING_ENABLED = false;

const SERVER_STATUS_POLL_DURATION = 5000;
const ACCESS_TOKEN_KEY = 'accessToken';

let serverStatusRefreshPoll: ReturnType<typeof setInterval>;
let hasBeenModeratorNotified = false;
let hasWebsocketDisconnected = false;

const serverConnectivityError = `Cannot connect to the Owncast service. Please check your internet connection and verify this Owncast server is running.`;

// Server status is what gets updated such as viewer count, durations,
// stream title, online/offline state, etc.
export const serverStatusState = atom<ServerStatus>({
  key: 'serverStatusState',
  default: makeEmptyServerStatus(),
});

// The config that comes from the API.
export const clientConfigStateAtom = atom({
  key: 'clientConfigState',
  default: makeEmptyClientConfig(),
});

export const accessTokenAtom = atom<string>({
  key: 'accessTokenAtom',
  default: null,
});

export const currentUserAtom = atom<CurrentUser>({
  key: 'currentUserAtom',
  default: null,
});

export const chatMessagesAtom = atom<ChatMessage[]>({
  key: 'chatMessages',
  default: [] as ChatMessage[],
});

export const chatAuthenticatedAtom = atom<boolean>({
  key: 'chatAuthenticatedAtom',
  default: false,
});

export const websocketServiceAtom = atom<WebsocketService>({
  key: 'websocketServiceAtom',
  default: null,
  dangerouslyAllowMutability: true,
});

export const appStateAtom = atom<AppStateOptions>({
  key: 'appState',
  default: makeEmptyAppState(),
});

export const isMobileAtom = atom<boolean | undefined>({
  key: 'isMobileAtom',
  default: undefined,
});

export const isVideoPlayingAtom = atom<boolean>({
  key: 'isVideoPlayingAtom',
  default: false,
});

export const fatalErrorStateAtom = atom<DisplayableError>({
  key: 'fatalErrorStateAtom',
  default: null,
});

export const clockSkewAtom = atom<Number>({
  key: 'clockSkewAtom',
  default: 0.0,
});

const removedMessageIdsAtom = atom<string[]>({
  key: 'removedMessageIds',
  default: [],
});

export const isChatAvailableSelector = selector({
  key: 'isChatAvailableSelector',
  get: ({ get }) => {
    const state: AppStateOptions = get(appStateAtom);
    const accessToken: string = get(accessTokenAtom);
    return accessToken && state.chatAvailable && !hasWebsocketDisconnected;
  },
});

// The requested state of chat in the UI
export enum ChatState {
  VISIBLE, // Chat is open (the default state when the stream is online)
  HIDDEN, // Chat is hidden
  POPPED_OUT, // Chat is playing in a popout window
  EMBEDDED, // This window is opened at /embed/chat/readwrite/
}

export const chatStateAtom = atom<ChatState>({
  key: 'chatState',
  default: (() => {
    // XXX Somehow, `window` is undefined here, even though this runs in client
    const window = globalThis;
    return window?.location?.pathname === '/embed/chat/readwrite/'
      ? ChatState.EMBEDDED
      : ChatState.VISIBLE;
  })(),
});

// We display in an "online/live" state as long as video is actively playing.
// Even during the time where technically the server has said it's no longer
// live, however the last few seconds of video playback is still taking place.
export const isOnlineSelector = selector({
  key: 'isOnlineSelector',
  get: ({ get }) => {
    const state: AppStateOptions = get(appStateAtom);
    const isVideoPlaying: boolean = get(isVideoPlayingAtom);
    return state.videoAvailable || isVideoPlaying;
  },
});

export const visibleChatMessagesSelector = selector<ChatMessage[]>({
  key: 'visibleChatMessagesSelector',
  get: ({ get }) => {
    const messages: ChatMessage[] = get(chatMessagesAtom);
    const removedIds: string[] = get(removedMessageIdsAtom);
    return messages.filter(message => !removedIds.includes(message.id));
  },
});

export const ClientConfigStore: FC = () => {
  const ClientConfigService = useContext(ClientConfigServiceContext);
  const ChatService = useContext(ChatServiceContext);
  const ServerStatusService = useContext(ServerStatusServiceContext);

  const [appState, appStateSend, appStateService] = useMachine(appStateModel);
  const [currentUser, setCurrentUser] = useRecoilState(currentUserAtom);
  const setChatAuthenticated = useSetRecoilState<boolean>(chatAuthenticatedAtom);
  const [clientConfig, setClientConfig] = useRecoilState<ClientConfig>(clientConfigStateAtom);
  const setServerStatus = useSetRecoilState<ServerStatus>(serverStatusState);
  const setClockSkew = useSetRecoilState<Number>(clockSkewAtom);
  const setChatMessages = useSetRecoilState<SocketEvent[]>(chatMessagesAtom);
  const [accessToken, setAccessToken] = useRecoilState<string>(accessTokenAtom);
  const setAppState = useSetRecoilState<AppStateOptions>(appStateAtom);
  const setGlobalFatalErrorMessage = useSetRecoilState<DisplayableError>(fatalErrorStateAtom);
  const setWebsocketService = useSetRecoilState<WebsocketService>(websocketServiceAtom);
  const setHiddenMessageIds = useSetRecoilState<string[]>(removedMessageIdsAtom);
  const [hasLoadedConfig, setHasLoadedConfig] = useState(false);

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
        setChatMessages(currentState => [...currentState, message as FediverseEvent]);
        break;
      case MessageType.FEDIVERSE_ENGAGEMENT_LIKE:
        setChatMessages(currentState => [...currentState, message as FediverseEvent]);
        break;
      case MessageType.FEDIVERSE_ENGAGEMENT_REPOST:
        setChatMessages(currentState => [...currentState, message as FediverseEvent]);
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

  // Read the config and status on initial load from a JSON string that lives
  // in window. This is placed there server-side and allows for fast initial
  // load times because we don't have to wait for the API calls to complete.
  useEffect(() => {
    try {
      if ((window as any).configHydration) {
        const config = JSON.parse((window as any).configHydration);
        setClientConfig(config);
        setHasLoadedConfig(true);
      }
    } catch (e) {
      console.error('Error parsing config hydration', e);
    }

    try {
      if ((window as any).statusHydration) {
        const status = JSON.parse((window as any).statusHydration);
        setServerStatus(status);
        handleStatusChange(status);
      }
    } catch (e) {
      console.error('error parsing status hydration', e);
    }

    try {
      if ((window as any).configHydration && (window as any).statusHydration) {
        sendEvent([AppStateEvent.Loaded]);
      }
    } catch (e) {
      console.error('error sending loaded event', e);
    }
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
    if (!(window as any).configHydration) {
      updateClientConfig();
    }
    handleUserRegistration();
    if (!(window as any).statusHydration) {
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
