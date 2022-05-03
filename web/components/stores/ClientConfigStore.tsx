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

export function ClientConfigStore() {
  const setClientConfig = useSetRecoilState<ClientConfig>(clientConfigStateAtom);
  const [appState, setAppState] = useRecoilState<AppState>(appStateAtom);
  const setChatVisibility = useSetRecoilState<ChatVisibilityState>(chatVisibilityAtom);
  const [chatState, setChatState] = useRecoilState<ChatState>(chatStateAtom);
  const setChatMessages = useSetRecoilState<ChatMessage[]>(chatMessagesAtom);
  const [accessToken, setAccessToken] = useRecoilState<string>(accessTokenAtom);
  const setChatDisplayName = useSetRecoilState<string>(chatDisplayNameAtom);

  const updateClientConfig = async () => {
    try {
      const config = await ClientConfigService.getConfig();
      // console.log(`ClientConfig: ${JSON.stringify(config)}`);
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

  const getChatHistory = async () => {
    setChatState(ChatState.Loading);
    try {
      const messages = await ChatService.getChatHistory(accessToken);
      // console.log(`ChatService -> getChatHistory() messages: \n${JSON.stringify(messages)}`);
      setChatMessages(messages);
    } catch (error) {
      console.error(`ChatService -> getChatHistory() ERROR: \n${error}`);
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

    console.log('access token changed', accessToken);
    getChatHistory();
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
