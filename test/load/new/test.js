import http from 'k6/http';
import ws from 'k6/ws';

import { check, sleep } from 'k6';
import { randomIntBetween, randomString } from "https://jslib.k6.io/k6-utils/1.0.0/index.js";

const baseUserAgent = 'Owncast LoadTest/1.0';

function randomNumber() {
  return randomIntBetween(1, 10000);
}

function randomSleep() {
  sleep(randomIntBetween(10, 60));
}

function pingViewerAPI() {
  // Fake the user-agent so the server side mapping to agent+ip
  // sees each user as unique.
  var params = {
    headers: {
      'User-Agent': baseUserAgent + ' ' + randomNumber(),
    },
  };

  const response = http.get('http://localhost:8080/api/ping', params);
  check(response, { 'status was 200': (r) => r.status == 200 });
}

function fetchHLSPlaylist() {
  const response = http.get('http://localhost:8080/hls/stream.m3u8');
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

function connectToChat(accessToken) {
  // Connect to the chat server via web socket.
  const params = {
    headers: {
      'User-Agent': `${baseUserAgent} (iteration ${__ITER}; virtual-user ${__VU})`,
    },
  };

  var wsResponse = ws.connect(
    `ws://127.0.0.1:8080/ws?accessToken=${accessToken}`,
    params,
    function (socket) {
      socket.on('open', function (data) {
        const testMessage = {
          body: `${randomString(randomIntBetween(10, 1000))} ${randomNumber()}`,
          type: 'CHAT',
        };

        randomSleep(); // After a user joins they wait to send a message
        socket.send(JSON.stringify(testMessage));
        randomSleep(); // The user waits after sending a message to leave.
        socket.close();
      });
    }
  );
  check(wsResponse, { 'status was 200': (r) => r.status == 200 });
}

export default function () {
  const accessToken = registerUser();
  // Fetch chat history once you have an access token.
  fetchChatHistory(accessToken);

  // A client hits the ping API every once in a while to
  // keep track of the number of viewers. So we can emulate
  // that.
  pingViewerAPI();

  // Emulate loading the master HLS playlist
  fetchHLSPlaylist();

  // Register as a new chat user and connect.
  connectToChat(accessToken);

  sleep(8); // Emulate the ping timer on the client.
  pingViewerAPI();
}

export let options = {
  userAgent: baseUserAgent,
  scenarios: {
    loadstages: {
      executor: 'ramping-vus',
      startVUs: 0,
      gracefulStop: '120s',
      stages: [
        { duration: '10s', target: 20 },
        { duration: '120s', target: 1000 },
        { duration: '300s', target: 4000 },
      ],
      gracefulRampDown: '10s',
    },
  },
};
