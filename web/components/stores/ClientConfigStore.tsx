/* eslint-disable no-case-declarations */
import { useEffect, useLayoutEffect } from 'react';
import { atom, useRecoilState, useSetRecoilState } from 'recoil';
import { makeEmptyClientConfig, ClientConfig } from '../../interfaces/client-config.model';
import ClientConfigService from '../../services/client-config-service';
import ChatService from '../../services/chat-service';
import WebsocketService from '../../services/websocket-service';
import { ChatMessage } from '../../interfaces/chat-message.model';
import { getLocalStorage, setLocalStorage } from '../../utils/helpers';
import {
  AppState,
  ChatState,
  ChatVisibilityState,
  getChatState,
  getChatVisibilityState,
} from '../../interfaces/application-state';
import {
  SocketEvent,
  ConnectedClientInfoEvent,
  MessageType,
  ChatEvent,
} from '../../interfaces/socket-events';
import handleConnectedClientInfoMessage from './eventhandlers/connectedclientinfo';
import handleChatMessage from './eventhandlers/handleChatMessage';

// The config that comes from the API.
export const clientConfigStateAtom = atom({
  key: 'clientConfigState',
  default: makeEmptyClientConfig(),
});

export const appStateAtom = atom<AppState>({
  key: 'appStateAtom',
  default: AppState.Loading,
});

export const chatStateAtom = atom<ChatState>({
  key: 'chatStateAtom',
  default: ChatState.Offline,
});

export const chatVisibilityAtom = atom<ChatVisibilityState>({
  key: 'chatVisibility',
  default: ChatVisibilityState.Visible,
});

export const chatDisplayNameAtom = atom<string>({
  key: 'chatDisplayName',
  default: null,
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

export function ClientConfigStore() {
  const setClientConfig = useSetRecoilState<ClientConfig>(clientConfigStateAtom);
  const setChatVisibility = useSetRecoilState<ChatVisibilityState>(chatVisibilityAtom);
  const setChatState = useSetRecoilState<ChatState>(chatStateAtom);
  const [chatMessages, setChatMessages] = useRecoilState<ChatMessage[]>(chatMessagesAtom);
  const setChatDisplayName = useSetRecoilState<string>(chatDisplayNameAtom);
  const [appState, setAppState] = useRecoilState<AppState>(appStateAtom);
  const [accessToken, setAccessToken] = useRecoilState<string>(accessTokenAtom);
  const [websocketService, setWebsocketService] =
    useRecoilState<WebsocketService>(websocketServiceAtom);

  let ws: WebsocketService;

  const updateClientConfig = async () => {
    try {
      const config = await ClientConfigService.getConfig();
      setClientConfig(config);
      setAppState(AppState.Online);
    } catch (error) {
      console.error(`ClientConfigService -> getConfig() ERROR: \n${error}`);
    }
  };

  const handleUserRegistration = async (optionalDisplayName?: string) => {
    try {
      setAppState(AppState.Registering);
      const response = await ChatService.registerUser(optionalDisplayName);
      console.log(`ChatService -> registerUser() response: \n${response}`);
      const { accessToken: newAccessToken, displayName: newDisplayName } = response;
      if (!newAccessToken) {
        return;
      }
      console.log('setting access token', newAccessToken);
      setAccessToken(newAccessToken);
      // setLocalStorage('accessToken', newAccessToken);
      setChatDisplayName(newDisplayName);
      setAppState(AppState.Online);
    } catch (e) {
      console.error(`ChatService -> registerUser() ERROR: \n${e}`);
    }
  };

  const handleMessage = (message: SocketEvent) => {
    switch (message.type) {
      case MessageType.CONNECTED_USER_INFO:
        handleConnectedClientInfoMessage(message as ConnectedClientInfoEvent);
        break;
      case MessageType.CHAT:
        handleChatMessage(message as ChatEvent, chatMessages, setChatMessages);
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
    setChatState(ChatState.Loading);
    try {
      ws = new WebsocketService(accessToken, '/ws');
      ws.handleMessage = handleMessage;
      setWebsocketService(ws);
    } catch (error) {
      console.error(`ChatService -> startChat() ERROR: \n${error}`);
    }
    setChatState(ChatState.Available);
  };

  useEffect(() => {
    updateClientConfig();
    handleUserRegistration();
  }, []);

  useLayoutEffect(() => {
    if (!accessToken) {
      return;
    }

    getChatHistory();
    startChat();
  }, [accessToken]);

  useEffect(() => {
    const updatedChatState = getChatState(appState);
    setChatState(updatedChatState);
    const updatedChatVisibility = getChatVisibilityState(appState);
    console.log(
      'app state: ',
      AppState[appState],
      'chat state:',
      ChatState[updatedChatState],
      'chat visibility:',
      ChatVisibilityState[updatedChatVisibility],
    );
    setChatVisibility(updatedChatVisibility);
  }, [appState]);

  return null;
}
