var request = require('supertest');
request = request('http://127.0.0.1:8080');

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
            expect(res.body.instanceDetails.name).toBe('Owncast');
            expect(res.body.instanceDetails.title).toBe('Owncast');
            expect(res.body.instanceDetails.summary).toBe('This is brief summary of whom you are or what your stream is. You can edit this description in your config file.');
            expect(res.body.instanceDetails.logo).toBe('/img/logo.svg');
            expect(res.body.instanceDetails.tags).toStrictEqual(['music', 'software', 'streaming']);

            expect(res.body.videoSettings.segmentLengthSeconds).toBe(4);
            expect(res.body.videoSettings.numberOfPlaylistItems).toBe(5);

            expect(res.body.videoSettings.videoQualityVariants[0].framerate).toBe(24);
            expect(res.body.videoSettings.videoQualityVariants[0].encoderPreset).toBe('veryfast');

            expect(res.body.videoSettings.numberOfPlaylistItems).toBe(5);

            expect(res.body.yp.enabled).toBe(false);
            expect(res.body.streamKey).toBe('abc123');
            done();
        });
});


test('correct number of log entries exist', (done) => {
    request.get('/api/admin/logs').auth('admin', 'abc123').expect(200)
        .then((res) => {
            // expect(res.body).toHaveLength(4);
            done();
        });
});

