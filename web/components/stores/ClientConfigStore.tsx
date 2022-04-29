import { useEffect } from 'react';
import { atom, useRecoilState } from 'recoil';
import { makeEmptyClientConfig, ClientConfig } from '../../interfaces/client-config.model';
import ClientConfigService from '../../services/client-config-service';
import ChatService from '../../services/chat-service';
import { ChatMessage } from '../../interfaces/chat-message.model';
import { getLocalStorage, setLocalStorage } from '../../utils/helpers';
import { ChatVisibilityState } from '../../interfaces/application-state';

// The config that comes from the API.
export const clientConfigState = atom({
  key: 'clientConfigState',
  default: makeEmptyClientConfig(),
});

export const chatVisibility = atom<ChatVisibilityState>({
  key: 'chatVisibility',
  default: ChatVisibilityState.Hidden,
});

export const chatDisplayName = atom({
  key: 'chatDisplayName',
  default: null,
});

export const accessTokenAtom = atom({
  key: 'accessToken',
  default: null,
});

export const chatMessages = atom({
  key: 'chatMessages',
  default: [] as ChatMessage[],
});

export function ClientConfigStore() {
  const [, setClientConfig] = useRecoilState<ClientConfig>(clientConfigState);
  const [, setChatMessages] = useRecoilState<ChatMessage[]>(chatMessages);
  const [accessToken, setAccessToken] = useRecoilState<string>(accessTokenAtom);
  const [, setChatDisplayName] = useRecoilState<string>(chatDisplayName);

  const updateClientConfig = async () => {
    try {
      const config = await ClientConfigService.getConfig();
      console.log(`ClientConfig: ${JSON.stringify(config)}`);
      setClientConfig(config);
    } catch (error) {
      console.error(`ClientConfigService -> getConfig() ERROR: \n${error}`);
    }
  };

  const handleUserRegistration = async (optionalDisplayName: string) => {
    try {
      const response = await ChatService.registerUser(optionalDisplayName);
      console.log(`ChatService -> registerUser() response: \n${JSON.stringify(response)}`);
      const { accessToken: newAccessToken, displayName } = response;
      if (!newAccessToken) {
        return;
      }
      setAccessToken(accessToken);
      setLocalStorage('accessToken', newAccessToken);
      setChatDisplayName(displayName);
    } catch (e) {
      console.error(`ChatService -> registerUser() ERROR: \n${e}`);
    }
  };

  // TODO: Requires access token.
  const getChatHistory = async () => {
    try {
      const messages = await ChatService.getChatHistory(accessToken);
      console.log(`ChatService -> getChatHistory() messages: \n${JSON.stringify(messages)}`);
      setChatMessages(messages);
    } catch (error) {
      console.error(`ChatService -> getChatHistory() ERROR: \n${error}`);
    }
  };

  useEffect(() => {
    updateClientConfig();
    handleUserRegistration();
  }, []);

  useEffect(() => {
    getChatHistory();
  }, [accessToken]);

  return null;
}
