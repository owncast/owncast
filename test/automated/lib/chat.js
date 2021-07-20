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

function sendChatMessage(message, accessToken, done) {
  const ws = new WebSocket(
    `ws://localhost:8080/ws?accessToken=${accessToken}`,
    {
      origin: 'http://localhost:8080',
    }
  );

  function onOpen() {
    ws.send(JSON.stringify(message), function () {
      ws.close();
      done();
    });
  }

  ws.on('open', onOpen);
}

module.exports.sendChatMessage = sendChatMessage;
module.exports.registerChat = registerChat;