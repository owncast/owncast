import { ChatMessage } from '../interfaces/chat-message.model';
import { ChatStaticService, UserRegistrationResponse } from './chat-service';

export const chatServiceMockOf = (
  chatHistory: ChatMessage[],
  userRegistrationResponse: UserRegistrationResponse,
): ChatStaticService =>
  class ChatServiceMock {
    public static async getChatHistory(): Promise<ChatMessage[]> {
      return chatHistory;
    }

    public static async registerUser(): Promise<UserRegistrationResponse> {
      return userRegistrationResponse;
    }
  };

export default chatServiceMockOf;
