import { useEffect } from 'react';
import { atom, selector, useRecoilState, useSetRecoilState } from 'recoil';
import { useMachine } from '@xstate/react';
import { makeEmptyClientConfig, ClientConfig } from '../../interfaces/client-config.model';
import ClientConfigService from '../../services/client-config-service';
import ChatService from '../../services/chat-service';
import WebsocketService from '../../services/websocket-service';
import { ChatMessage } from '../../interfaces/chat-message.model';
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
  SocketEvent,
} from '../../interfaces/socket-events';

import handleChatMessage from './eventhandlers/handleChatMessage';
import handleConnectedClientInfoMessage from './eventhandlers/connected-client-info-handler';
import ServerStatusService from '../../services/status-service';
import handleNameChangeEvent from './eventhandlers/handleNameChangeEvent';
import { DisplayableError } from '../../types/displayable-error';

const SERVER_STATUS_POLL_DURATION = 5000;
const ACCESS_TOKEN_KEY = 'accessToken';

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

export const chatDisplayNameAtom = atom<string>({
  key: 'chatDisplayName',
  default: null,
});

export const chatUserIdAtom = atom<string>({
  key: 'chatUserIdAtom',
  default: null,
});

export const isChatModeratorAtom = atom<boolean>({
  key: 'isModeratorAtom',
  default: false,
});

export const accessTokenAtom = atom<string>({
  key: 'accessTokenAtom',
  default: null,
});

export const chatMessagesAtom = atom<ChatMessage[]>({
  key: 'chatMessages',
  default: [] as ChatMessage[],
});

export const websocketServiceAtom = atom<WebsocketService>({
  key: 'websocketServiceAtom',
  default: null,
});

export const appStateAtom = atom<AppStateOptions>({
  key: 'appState',
  default: makeEmptyAppState(),
});

export const chatVisibleToggleAtom = atom<boolean>({
  key: 'chatVisibilityToggleAtom',
  default: true,
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

// Chat is visible if the user wishes it to be visible AND the required
// chat state is set.
export const isChatVisibleSelector = selector({
  key: 'isChatVisibleSelector',
  get: ({ get }) => {
    const state: AppStateOptions = get(appStateAtom);
    const userVisibleToggle: boolean = get(chatVisibleToggleAtom);
    const accessToken: String = get(accessTokenAtom);
    return accessToken && state.chatAvailable && userVisibleToggle;
  },
});

export const isChatAvailableSelector = selector({
  key: 'isChatAvailableSelector',
  get: ({ get }) => {
    const state: AppStateOptions = get(appStateAtom);
    const accessToken: String = get(accessTokenAtom);
    return accessToken && state.chatAvailable;
  },
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

// Take a nested object of state metadata and merge it into
// a single flattened node.
function mergeMeta(meta) {
  return Object.keys(meta).reduce((acc, key) => {
    const value = meta[key];
    Object.assign(acc, value);

    return acc;
  }, {});
}

export function ClientConfigStore() {
  const [appState, appStateSend, appStateService] = useMachine(appStateModel);

  const setChatDisplayName = useSetRecoilState<string>(chatDisplayNameAtom);
  const setChatUserId = useSetRecoilState<string>(chatUserIdAtom);
  const setIsChatModerator = useSetRecoilState<boolean>(isChatModeratorAtom);
  const setClientConfig = useSetRecoilState<ClientConfig>(clientConfigStateAtom);
  const setServerStatus = useSetRecoilState<ServerStatus>(serverStatusState);
  const setClockSkew = useSetRecoilState<Number>(clockSkewAtom);
  const [chatMessages, setChatMessages] = useRecoilState<ChatMessage[]>(chatMessagesAtom);
  const [accessToken, setAccessToken] = useRecoilState<string>(accessTokenAtom);
  const setAppState = useSetRecoilState<AppStateOptions>(appStateAtom);
  const setGlobalFatalErrorMessage = useSetRecoilState<DisplayableError>(fatalErrorStateAtom);
  const setWebsocketService = useSetRecoilState<WebsocketService>(websocketServiceAtom);

  let ws: WebsocketService;

  const setGlobalFatalError = (title: string, message: string) => {
    setGlobalFatalErrorMessage({
      title,
      message,
    });
  };
  const sendEvent = (event: string) => {
    // console.log('---- sending event:', event);
    appStateSend({ type: event });
  };

  const updateClientConfig = async () => {
    try {
      const config = await ClientConfigService.getConfig();
      setClientConfig(config);
      sendEvent('LOADED');
      setGlobalFatalErrorMessage(null);
    } catch (error) {
      setGlobalFatalError(
        'Unable to reach Owncast server',
        `Owncast cannot launch. Please make sure the Owncast server is running. ${error}`,
      );
      console.error(`ClientConfigService -> getConfig() ERROR: \n${error}`);
    }
  };

  const updateServerStatus = async () => {
    try {
      const status = await ServerStatusService.getStatus();
      setServerStatus(status);
      const { serverTime } = status;

      const clockSkew = new Date(serverTime).getTime() - Date.now();
      setClockSkew(clockSkew);

      if (status.online) {
        sendEvent(AppStateEvent.Online);
      } else if (!status.online) {
        sendEvent(AppStateEvent.Offline);
      }
      setGlobalFatalErrorMessage(null);
    } catch (error) {
      sendEvent(AppStateEvent.Fail);
      setGlobalFatalError(
        'Unable to reach Owncast server',
        `Owncast cannot launch. Please make sure the Owncast server is running. ${error}`,
      );
      console.error(`serverStatusState -> getStatus() ERROR: \n${error}`);
    }
    return null;
  };

  const handleUserRegistration = async (optionalDisplayName?: string) => {
    const savedAccessToken = getLocalStorage(ACCESS_TOKEN_KEY);
    if (savedAccessToken) {
      setAccessToken(savedAccessToken);
      return;
    }

    try {
      sendEvent(AppStateEvent.NeedsRegister);
      const response = await ChatService.registerUser(optionalDisplayName);
      console.log(`ChatService -> registerUser() response: \n${response}`);
      const { accessToken: newAccessToken, displayName: newDisplayName } = response;
      if (!newAccessToken) {
        return;
      }

      console.log('setting access token', newAccessToken);
      setAccessToken(newAccessToken);
      setLocalStorage(ACCESS_TOKEN_KEY, newAccessToken);
      setChatDisplayName(newDisplayName);
    } catch (e) {
      sendEvent(AppStateEvent.Fail);
      console.error(`ChatService -> registerUser() ERROR: \n${e}`);
    }
  };

  const resetAndReAuth = () => {
    setLocalStorage(ACCESS_TOKEN_KEY, '');
    setAccessToken('');
    handleUserRegistration();
  };

  const handleMessage = (message: SocketEvent) => {
    switch (message.type) {
      case MessageType.ERROR_NEEDS_REGISTRATION:
        resetAndReAuth();
        break;
      case MessageType.CONNECTED_USER_INFO:
        handleConnectedClientInfoMessage(
          message as ConnectedClientInfoEvent,
          setChatDisplayName,
          setChatUserId,
          setIsChatModerator,
        );
        break;
      case MessageType.CHAT:
        handleChatMessage(message as ChatEvent, chatMessages, setChatMessages);
        break;
      case MessageType.NAME_CHANGE:
        handleNameChangeEvent(message as ChatEvent, chatMessages, setChatMessages);
        break;
      default:
        console.error('Unknown socket message type: ', message.type);
    }
  };

  const getChatHistory = async () => {
    try {
      const messages = await ChatService.getChatHistory(accessToken);
      const updatedChatMessages = [...messages, ...chatMessages];
      setChatMessages(updatedChatMessages);
    } catch (error) {
      console.error(`ChatService -> getChatHistory() ERROR: \n${error}`);
    }
  };

  const startChat = async () => {
    sendEvent(AppStateEvent.Loading);
    try {
      ws = new WebsocketService(accessToken, '/ws');
      ws.handleMessage = handleMessage;
      setWebsocketService(ws);
      sendEvent(AppStateEvent.Loaded);
    } catch (error) {
      console.error(`ChatService -> startChat() ERROR: \n${error}`);
    }
  };

  useEffect(() => {
    updateClientConfig();
    handleUserRegistration();
    updateServerStatus();
    setInterval(() => {
      updateServerStatus();
    }, SERVER_STATUS_POLL_DURATION);
  }, [appState]);

  useEffect(() => {
    if (!accessToken) {
      return;
    }

    getChatHistory();
    startChat();
  }, [accessToken]);

  useEffect(() => {
    appStateService.onTransition(state => {
      if (!state.changed) {
        return;
      }

      const metadata = mergeMeta(state.meta) as AppStateOptions;

      console.debug('--- APP STATE: ', state.value);
      console.debug('--- APP META: ', metadata);

      setAppState(metadata);
    });
  });

  return null;
}
