import { createContext } from 'react';
import { ChatMessage } from '../interfaces/chat-message.model';
import { getUnauthedData } from '../utils/apis';

const ENDPOINT = `/api/chat`;
const URL_CHAT_REGISTRATION = `/api/chat/register`;

export interface UserRegistrationResponse {
  id: string;
  accessToken: string;
  displayName: string;
  displayColor: number;
}

export interface ChatStaticService {
  getChatHistory(accessToken: string): Promise<ChatMessage[]>;
  registerUser(username: string): Promise<UserRegistrationResponse>;
}

class ChatService {
  public static async getChatHistory(accessToken: string): Promise<ChatMessage[]> {
    try {
      const response = await getUnauthedData(`${ENDPOINT}?accessToken=${accessToken}`);
      return response;
    } catch (e) {
      return [];
    }
  }

  public static async registerUser(username: string): Promise<UserRegistrationResponse> {
    const options = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ displayName: username }),
    };

    const response = await getUnauthedData(URL_CHAT_REGISTRATION, options);
    return response;
  }
}

export const ChatServiceContext = createContext<ChatStaticService>(ChatService);
