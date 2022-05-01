import http from 'k6/http';
import ws from 'k6/ws';

import { check, sleep } from 'k6';
import {
  randomIntBetween,
  randomString,
} from 'https://jslib.k6.io/k6-utils/1.0.0/index.js';

const baseUserAgent = 'Owncast LoadTest/1.0';

function randomNumber() {
  return randomIntBetween(1, 10000);
}

function pingViewerAPI(params) {
  const response = http.get('http://localhost:8080/api/ping', params);
  check(response, { 'status was 200': (r) => r.status == 200 });
}

function fetchStatusAPI(params) {
  const response = http.get('http://localhost:8080/api/status', params);
  check(response, { 'status was 200': (r) => r.status == 200 });
}

function fetchHLSPlaylist(params) {
  const response = http.get('http://localhost:8080/hls/0/stream.m3u8', params);
  check(response, { 'status was 200': (r) => r.status == 200 });
}

function fetchChatHistory(accessToken) {
  const response = http.get(
    `http://localhost:8080/api/chat?accessToken=${accessToken}`
  );
  check(response, { 'status was 200': (r) => r.status == 200 });
}

function registerUser() {
  const response = http.post('http://localhost:8080/api/chat/register');
  check(response, { 'status was 200': (r) => r.status == 200 });
  const accessToken = response.json('accessToken');

  return accessToken;
}

function connectToChat(accessToken, params, callback) {
  // Connect to the chat server via web socket.
  var chatSocket = ws.connect(
    `ws://127.0.0.1:8080/ws?accessToken=${accessToken}`,
    params,
    callback
  );
  check(chatSocket, { 'status was 200': (r) => r.status == 200 });
  return chatSocket;
}

export default function () {
  // Fake the user-agent so the server side mapping to agent+ip
  // sees each user as unique.
  var params = {
    headers: {
      'User-Agent': baseUserAgent + ' ' + randomNumber(),
    },
  };

  // Register a new chat user.
  const accessToken = registerUser(params);

  // Fetch chat history once you have an access token.
  fetchChatHistory(accessToken);

  // Connect to websocket and send messages
  const callback = (chatSocket) => {
    chatSocket.on('open', function (data) {
      const testMessage = {
        body: `${randomString(randomIntBetween(10, 1000))} ${randomNumber()}`,
        type: 'CHAT',
      };

      chatSocket.send(JSON.stringify(testMessage));
      sleep(4);
      chatSocket.send(JSON.stringify(testMessage));
      sleep(4);
      chatSocket.send(JSON.stringify(testMessage));
      sleep(4);

      chatSocket.close();
    });
  };
  connectToChat(accessToken, params, callback);

  // Emulate a user playing back video and hitting the ping api.
  pingViewerAPI(params);
  fetchHLSPlaylist(params);
  fetchStatusAPI(params);
}

export let options = {
  userAgent: baseUserAgent,
  scenarios: {
    loadstages: {
      executor: 'ramping-vus',
      gracefulStop: '10s',
      stages: [
        { duration: '10s', target: 200 },
        { duration: '25s', target: 3000 },
        { duration: '10s', target: 0 },
      ],
    },
  },
};
