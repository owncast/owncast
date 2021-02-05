var request = require('supertest');
request = request('http://127.0.0.1:8080');

const serverName = randomString();
const streamTitle = randomString();
const serverSummary = randomString();
const pageContent = `<p>${randomString()}</p>`;
const logo = '/img/smiley.png';
const tags = [randomString(), randomString(), randomString()];
const latencyLevel = Math.floor((Math.random() * 5) + 1);

const streamOutputVariants = 
    {
        videoBitrate: randomNumber() * 100,
        framerate: 42,
        cpuUsageLevel: 2,
        scaledHeight: randomNumber() * 100,
        scaledWidth: randomNumber() * 100,
    };
const socialHandles = [
    {
        url: 'http://facebook.org/' + randomString(),
        platform: randomString()
    }
];

const s3Config = {
    enabled: true,
    endpoint: 'http://' + randomString(),
    accessKey: randomString(),
    secret: randomString(),
    bucket: randomString(),
    region: randomString()
};

test('set server name', async (done) => {
    const res = await sendConfigChangeRequest('name', serverName);
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

test('set latency level', async (done) => {
    const res = await sendConfigChangeRequest('video/streamlatencylevel', latencyLevel);
    done();
});

test('set video stream output variants', async (done) => {
    const res = await sendConfigChangeRequest('video/streamoutputvariants', [streamOutputVariants]);
    done();
});

test('set social handles', async (done) => {
    const res = await sendConfigChangeRequest('socialhandles', socialHandles);
    done();
});

test('set s3 configuration', async (done) => {
    const res = await sendConfigChangeRequest('s3', s3Config);
    done();
});

test('verify updated config values', async (done) => {
    const res = await request.get('/api/config');
    expect(res.body.name).toBe(serverName);
    expect(res.body.streamTitle).toBe(streamTitle);
    expect(res.body.summary).toBe(serverSummary);
    expect(res.body.extraPageContent).toBe(pageContent);
    expect(res.body.logo).toBe(logo);
    expect(res.body.socialHandles).toStrictEqual(socialHandles);
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
            expect(res.body.instanceDetails.summary).toBe(serverSummary);
            expect(res.body.instanceDetails.logo).toBe(logo);
            expect(res.body.instanceDetails.tags).toStrictEqual(tags);
            expect(res.body.instanceDetails.socialHandles).toStrictEqual(socialHandles);

            expect(res.body.videoSettings.latencyLevel).toBe(latencyLevel);
            expect(res.body.videoSettings.videoQualityVariants[0].framerate).toBe(streamOutputVariants.framerate);
            expect(res.body.videoSettings.videoQualityVariants[0].cpuUsageLevel).toBe(streamOutputVariants.cpuUsageLevel);

            expect(res.body.yp.enabled).toBe(false);
            expect(res.body.streamKey).toBe('abc123');
            
            expect(res.body.s3.enabled).toBe(s3Config.enabled);
            expect(res.body.s3.endpoint).toBe(s3Config.endpoint);
            expect(res.body.s3.accessKey).toBe(s3Config.accessKey);
            expect(res.body.s3.secret).toBe(s3Config.secret);
            expect(res.body.s3.bucket).toBe(s3Config.bucket);
            expect(res.body.s3.region).toBe(s3Config.region);

            done();
        });
});

test('frontend configuration is correct', (done) => {
    request.get('/api/config').expect(200)
        .then((res) => {
            expect(res.body.name).toBe(serverName);
            expect(res.body.logo).toBe(logo);
            expect(res.body.socialHandles).toStrictEqual(socialHandles);
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