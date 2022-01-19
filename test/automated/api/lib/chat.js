var request = require('supertest');
request = request('http://127.0.0.1:8080');
const WebSocket = require('ws');

async function registerChat() {
  try {
    const response = await request.post('/api/chat/register');
    return response.body;
  } catch (e) {
    console.error(e);
  }
}

async function sendChatMessage(message, accessToken, done) {
  const ws = new WebSocket(
    `ws://localhost:8080/ws?accessToken=${accessToken}`,
    {
      origin: 'http://localhost:8080',
    }
  );

  async function onOpen() {
    ws.send(JSON.stringify(message), async function () {
      ws.close();
      await new Promise((r) => setTimeout(r, 2000));
      done();
    });
  }

  ws.on('open', onOpen);
}

async function listenForEvent(name, accessToken, done) {
  const ws = new WebSocket(
    `ws://localhost:8080/ws?accessToken=${accessToken}`,
    {
      origin: 'http://localhost:8080',
    }
  );

  ws.on('message', function incoming(message) {
    const messages = message.split('\n');
    messages.forEach(function (message) {
      const event = JSON.parse(message);

      if (event.type === name) {
        done();
        ws.close();
      }
    });
  });
}

module.exports.sendChatMessage = sendChatMessage;
module.exports.registerChat = registerChat;
module.exports.listenForEvent = listenForEvent;
