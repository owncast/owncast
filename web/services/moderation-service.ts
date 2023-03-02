const HIDE_MESSAGE_ENDPOINT = `/api/chat/messagevisibility`;
const BAN_USER_ENDPOINT = `/api/chat/users/setenabled`;

class ChatModerationService {
  public static async removeMessage(id: string, accessToken: string): Promise<any> {
    const url = new URL(HIDE_MESSAGE_ENDPOINT, window.location.toString());
    url.searchParams.append('accessToken', accessToken);
    const hideMessageUrl = url.toString();

    const options = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ idArray: [id] }),
    };

    await fetch(hideMessageUrl, options);
  }

  public static async banUser(id: string, accessToken: string): Promise<any> {
    const url = new URL(BAN_USER_ENDPOINT, window.location.toString());
    url.searchParams.append('accessToken', accessToken);
    const hideMessageUrl = url.toString();

    const options = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ userId: id }),
    };

    await fetch(hideMessageUrl, options);
  }
}

export default ChatModerationService;
