const { test } = require('@jest/globals');
var request = require('supertest');
request = request('http://127.0.0.1:8080');

const WebSocket = require('ws');
var ws;

const testVisibilityMessage = {
    author: "username",
    body: "message " + Math.floor(Math.random() * 100),
    type: 'CHAT',
    visible: true,
    timestamp: new Date().toISOString()
};

test('can send a chat message', (done) => {
    ws = new WebSocket('ws://127.0.0.1:8080/entry', {
        origin: 'http://localhost',
    });

    function onOpen() {
        ws.send(JSON.stringify(testVisibilityMessage), function () {
            ws.close();
            done();
        });
    }

    ws.on('open', onOpen);
});

var messageId;

test('verify we can make API call to mark message as hidden', async (done) => {
    const res = await request.get('/api/admin/chat/messages').auth('admin', 'abc123').expect(200)
    const message = res.body[0];
    messageId = message.id;
    await request.post('/api/admin/chat/updatemessagevisibility')
        .auth('admin', 'abc123')
        .send({ "idArray": [messageId], "visible": false }).expect(200);
    done();
});

test('verify message has become hidden', async (done) => {
    const res = await request.get('/api/admin/chat/messages')
        .expect(200)
        .auth('admin', 'abc123')

    const message = res.body.filter(obj => {
        return obj.id === messageId;
    });
    expect(message.length).toBe(1);
    expect(message[0].visible).toBe(false);
    done();
});
