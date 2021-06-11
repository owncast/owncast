const WebSocket = require('ws');
const fetch = require('node-fetch');

const messages = [
  'I am a test message',
  'this is fake',
  'i write emoji 😀',
  'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.',
  'Sed pulvinar proin gravida hendrerit. Mauris in aliquam sem fringilla ut morbi tincidunt augue. In cursus turpis massa tincidunt dui.',
  'Feugiat in ante metus dictum at tempor commodo ullamcorper. Nunc aliquet bibendum enim facilisis gravida neque convallis a. Vitae tortor condimentum lacinia quis vel eros donec ac odio.',
];

var availableMessages = messages.slice();


async function registerChat() {
  const options = {
      method: 'POST',
      headers: {
          'Content-Type': 'application/json'
      }
  }

  try {
      const response = await fetch('http://localhost:8080/api/chat/register', options);
      const result = await response.json();
      return result;
  } catch(e) {
      console.error(e);
  }
}

async function sendMessage() {
  const registration = await registerChat();
  const accessToken = registration.accessToken;

  function send() {
    if (availableMessages.length == 0) {
      availableMessages = messages.slice();
    }
  
    const messageIndex = Math.floor(Math.random() * availableMessages.length);
    const message = availableMessages[messageIndex];
    availableMessages.splice(messageIndex, 1);
  
    const testMessage = {
      body: message,
      type: 'CHAT',
    };
  
    ws.send(JSON.stringify(testMessage));
  
    const nextMessageTimeout = (Math.floor(Math.random() * (25 + 10) * 10));
    setTimeout(sendMessage, nextMessageTimeout);
  }

  const ws = new WebSocket(`ws://localhost:8080/ws?accessToken=${accessToken}`, {
    origin: 'http://localhost:8080',
  });

  ws.on('open', function open() {
    setTimeout(send, 1000);
  });

  ws.on('error', function incoming(data) {
    console.log(data);
  });
}

sendMessage();
