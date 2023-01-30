import { ChatMessage } from '../interfaces/chat-message.model';
import { getUnauthedData } from '../utils/apis';

const ENDPOINT = `/api/chat`;
const URL_CHAT_REGISTRATION = `/api/chat/register`;

interface UserRegistrationResponse {
  id: string;
  accessToken: string;
  displayName: string;
  displayColor: number;
}

class ChatService {
  public static async getChatHistory(accessToken: string): Promise<ChatMessage[]> {
    const response = await getUnauthedData(`${ENDPOINT}?accessToken=${accessToken}`);
    return response;
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

export default ChatService;
