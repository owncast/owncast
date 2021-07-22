import http from 'k6/http';
import ws from 'k6/ws';

import { check, sleep } from 'k6';

const baseUserAgent = 'Owncast LoadTest/1.0';

function randomNumber() {
  return Math.floor(Math.random() * 10000) + 1
}

function pingViewerAPI() {
  // Fake the user-agent so the server side mapping to agent+ip
  // sees each user as unique.
  var params = {
    headers: {
      'User-Agent': 'test-client-' + randomNumber(),
    },
  };

  http.get('http://localhost:8080/api/ping', params);
}

function fetchHLSPlaylist() {
  http.get('http://localhost:8080/hls/stream.m3u8');
}

function connectToChat() {
    const response = http.post('http://localhost:8080/api/chat/register');
    check(response, { 'status was 200': (r) => r.status == 200 });
    const accessToken = response.json('accessToken');

    const params = {
      headers: {
        "User-Agent": `${baseUserAgent} (iteration ${__ITER}; virtual-user ${__VU})`,
      }
    };  
      
  var res = ws.connect(
    `ws://127.0.0.1:8080/ws?accessToken=${accessToken}`,
    params,
    function (socket) {
      socket.on('open', function (data) {
        const testMessage = {
          body: `Test message ${randomNumber()}`,
          type: 'CHAT',
        };
        sleep(10); // After a user joins they wait 10 seconds to send a message
        socket.send(JSON.stringify(testMessage));
        sleep(60); // The user waits a minute after sending a message to leave.
        socket.close();
      });
    }
  );
}

export default function () {
  pingViewerAPI();
  fetchHLSPlaylist();
  connectToChat();
}

export let options = {
    userAgent: baseUserAgent,
    scenarios: {
        // starting: {
        //     executor: 'shared-iterations',
        //     gracefulStop: '5s',
        //     vus: 10,
        //     iterations: 100,
        //     env: { SEND_MESSAGES: "true" },
        // },
        loadstages: {
          executor: 'ramping-vus',
          startVUs: 0,
          gracefulStop: '5s',
          stages: [
            { duration: '10s', target: 5 },
            { duration: '30s', target: 100 },
            { duration: '120s', target: 1000 },
            { duration: '300s', target: 5000 },
          ],
          gracefulRampDown: '10s',
          env: { SEND_MESSAGES: "false" },
      }
    }
  };
