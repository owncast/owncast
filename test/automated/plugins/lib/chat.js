// Minimal chat helpers for the plugin integration tests: register an
// anonymous chat user and send a message over the websocket. (Mirrors
// test/automated/api/lib/chat.js, trimmed to what these tests need.)
var request = require('supertest');
request = request('http://127.0.0.1:8080');
const WebSocket = require('ws');

async function registerChat() {
	const response = await request.post('/api/chat/register');
	return response.body;
}

function sendChatMessage(message, accessToken) {
	const ws = new WebSocket(`ws://localhost:8080/ws?accessToken=${accessToken}`);
	ws.on('open', function () {
		ws.send(JSON.stringify(message), function () {
			ws.close();
		});
	});
}

module.exports.registerChat = registerChat;
module.exports.sendChatMessage = sendChatMessage;
