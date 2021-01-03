const usernames = [
  'User ' + Math.floor(Math.random() * 100),
  'User ' + Math.floor(Math.random() * 100),
  'User ' + Math.floor(Math.random() * 100),
  'User ' + Math.floor(Math.random() * 100),
];

const messages = [
  'I am a test message',
  'this is fake',
  'i write emoji 😀',
  'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.',
  'Sed pulvinar proin gravida hendrerit. Mauris in aliquam sem fringilla ut morbi tincidunt augue. In cursus turpis massa tincidunt dui.',
  'Feugiat in ante metus dictum at tempor commodo ullamcorper. Nunc aliquet bibendum enim facilisis gravida neque convallis a. Vitae tortor condimentum lacinia quis vel eros donec ac odio.',
];

var availableMessages = messages.slice();

const WebSocket = require('ws');

const ws = new WebSocket('ws://localhost:8080/entry', {
  origin: 'http://watch.owncast.online',
});

ws.on('open', function open() {
  setTimeout(sendMessage, 15000);
});

ws.on('error', function incoming(data) {
  console.log(data);
});

function sendMessage() {
  if (availableMessages.length == 0) {
    availableMessages = messages.slice();
  }

  const id = Math.random().toString(36).substring(7);
  const username = usernames[Math.floor(Math.random() * usernames.length)];
  const messageIndex = Math.floor(Math.random() * availableMessages.length);
  const message = availableMessages[messageIndex];
  availableMessages.splice(messageIndex, 1);

  const testMessage = {
    author: username,
    body: message,
    image: 'https://robohash.org/' + username,
    id: id,
    type: 'CHAT',
    visible: true,
    timestamp: new Date().toISOString(),
  };

  ws.send(JSON.stringify(testMessage));

  const nextMessageTimeout = (Math.floor(Math.random() * (25 - 10)) + 10) * 100;
  setTimeout(sendMessage, nextMessageTimeout);
}



