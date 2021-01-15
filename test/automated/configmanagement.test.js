var request = require('supertest');
request = request('http://127.0.0.1:8080');

const serverName = randomString();
const serverTitle = randomString();
const streamTitle = randomString();
const serverSummary = randomString();
const pageContent = `<p>${randomString()}</p>`;
const logo = '/img/' + randomString();
const tags = [randomString(), randomString(), randomString()];
const segmentLength = randomNumber();
const segmentCount = randomNumber();
const streamOutputVariants = [
    {
        videoBitrate: randomNumber() * 100,
        framerate: randomNumber() * 10,
        encoderPreset: 'fast',
        scaledHeight: randomNumber() * 100,
        scaledWidth: randomNumber() * 100,
    }
];

test('set server name', async (done) => {
    const res = await sendConfigChangeRequest('name', serverName);
    done();
});

test('set server title', async (done) => {
    const res = await sendConfigChangeRequest('servertitle', serverTitle);
    done();
});

test('set stream title', async (done) => {
    const res = await sendConfigChangeRequest('streamtitle', streamTitle);
    done();
});

test('set server summary', async (done) => {
    const res = await sendConfigChangeRequest('serversummary', serverSummary);
    done();
});

test('set extra page content', async (done) => {
    const res = await sendConfigChangeRequest('pagecontent', pageContent);
    done();
});

test('set logo', async (done) => {
    const res = await sendConfigChangeRequest('logo', logo);
    done();
});

test('set tags', async (done) => {
    const res = await sendConfigChangeRequest('tags', tags);
    done();
});

test('set segment duration', async (done) => {
    const res = await sendConfigChangeRequest('video/segmentlength', segmentLength);
    done();
});

test('set segment count', async (done) => {
    const res = await sendConfigChangeRequest('video/segmentcount', segmentCount);
    done();
});

test('set video stream output variants', async (done) => {
    const res = await sendConfigChangePayload('video/streamoutputvariants', streamOutputVariants);
    done();
});

test('verify updated config values', async (done) => {
    const res = await request.get('/api/config');
    expect(res.body.name).toBe(serverName);
    expect(res.body.streamTitle).toBe(streamTitle);
    expect(res.body.summary).toBe(serverSummary);
    expect(res.body.extraPageContent).toBe(pageContent);
    expect(res.body.logo).toBe(logo);

    done();
});

// Test that the raw video details being broadcasted are coming through 
test('stream details are correct', (done) => {
    request.get('/api/admin/status').auth('admin', 'abc123').expect(200)
        .then((res) => {
            expect(res.body.broadcaster.streamDetails.width).toBe(320);
            expect(res.body.broadcaster.streamDetails.height).toBe(180);
            expect(res.body.broadcaster.streamDetails.framerate).toBe(24);
            expect(res.body.broadcaster.streamDetails.videoBitrate).toBe(1269);
            expect(res.body.broadcaster.streamDetails.videoCodec).toBe('H.264');
            expect(res.body.broadcaster.streamDetails.audioCodec).toBe('AAC');
            expect(res.body.online).toBe(true);
            done();
        });
});

test('admin configuration is correct', (done) => {
    request.get('/api/admin/serverconfig').auth('admin', 'abc123').expect(200)
        .then((res) => {
            expect(res.body.instanceDetails.name).toBe(serverName);
            expect(res.body.instanceDetails.title).toBe(serverTitle);
            expect(res.body.instanceDetails.summary).toBe(serverSummary);
            expect(res.body.instanceDetails.logo).toBe(logo);
            expect(res.body.instanceDetails.tags).toStrictEqual(tags);

            expect(res.body.videoSettings.segmentLengthSeconds).toBe(segmentLength);
            expect(res.body.videoSettings.numberOfPlaylistItems).toBe(segmentCount);

            expect(res.body.videoSettings.videoQualityVariants[0].framerate).toBe(streamOutputVariants[0].framerate);
            expect(res.body.videoSettings.videoQualityVariants[0].encoderPreset).toBe(streamOutputVariants[0].encoderPreset);

            expect(res.body.yp.enabled).toBe(false);
            expect(res.body.streamKey).toBe('abc123');
            done();
        });
});

test('frontend configuration is correct', (done) => {
    request.get('/api/config').expect(200)
        .then((res) => {
            expect(res.body.title).toBe(serverTitle);
            expect(res.body.logo).toBe(logo);
            expect(res.body.socialHandles[0].platform).toBe('github');
            expect(res.body.socialHandles[0].url).toBe('https://github.com/owncast/owncast');
            done();
        });
});


async function sendConfigChangeRequest(endpoint, value) {
    const url = '/api/admin/config/' + endpoint;
    const res = await request.post(url)
        .auth('admin', 'abc123')
        .send({"value": value}).expect(200);

    expect(res.body.success).toBe(true);

    return res
}

async function sendConfigChangePayload(endpoint, payload) {
    const url = '/api/admin/config/' + endpoint;
    const res = await request.post(url)
        .auth('admin', 'abc123')
        .send(payload).expect(200);

    expect(res.body.success).toBe(true);

    return res
}


function randomString(length = 20) {
    return Math.random().toString(16).substr(2, length);
}

function randomNumber() {
    return Math.floor(Math.random() * 5); 
}