const URL_CHAT_REGISTRATION = `http://localhost:8080/api/chat/register`;

interface UserRegistrationResponse {
  id: string;
  accessToken: string;
  displayName: string;
}

class UserService {
  registerUser(username: string): Promise<UserRegistrationResponse> {
    const options = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ displayName: username }),
    };
  
    try {
      const response = await fetch(URL_CHAT_REGISTRATION, options);
      const result = await response.json();
      return result;
    } catch (e) {
      console.error(e);
    }

    return null;
  }
}