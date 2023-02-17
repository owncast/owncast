// This will send raw, unformatted strings to the websocket to make sure the socket server
// is handling bad data.

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
  
  const ws = new WebSocket('ws://localhost:8080/ws', {
    origin: 'http://watch.owncast.online',
  });
  
  ws.on('open', function open() {
    setTimeout(sendMessage, 100);
  });
  
  ws.on('error', function incoming(data) {
    console.log(data);
  });
  
  function sendMessage() {
    const messageIndex = Math.floor(Math.random() * availableMessages.length);
    ws.send(JSON.stringify(availableMessages[messageIndex]));

    setTimeout(sendMessage, 100);
  }
  
  
  
  