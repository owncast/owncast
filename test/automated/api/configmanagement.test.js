var request = require('supertest');

const Random = require('crypto-random');

const sendConfigChangeRequest = require('./lib/config').sendConfigChangeRequest;
const getAdminConfig = require('./lib/config').getAdminConfig;
const getAdminStatus = require('./lib/config').getAdminStatus;

request = request('http://127.0.0.1:8080');


// initial configuration of server
const defaultAdminPassword = 'abc123';
const defaultStreamKeys = [{ key: defaultAdminPassword, comment: 'Default stream key' }];
const defaultYPConfig = {
	enabled: false
};
const defaultS3Config = {
	enabled: false,
	forcePathStyle: false
};
const defaultFederationConfig = {
	enabled: false,
	isPrivate: false,
	showEngagement: true,
	goLiveMessage: "I've gone live!",
	blockedDomains: []
};
const defaultHideViewerCount = false;

// new configuration for testing
const newServerName = randomString();
const newStreamTitle = randomString();
const newServerSummary = randomString();
const newOfflineMessage = randomString();
const newPageContent = `<p>${randomString()}</p>`;
const newTags = [randomString(), randomString(), randomString()];
const newStreamKeys = [
	{ key: randomString(), comment: 'test key 1' },
	{ key: randomString(), comment: 'test key 2' },
	{ key: randomString(), comment: 'test key 3' },
];
const newAdminPassword = randomString();

const latencyLevel = Math.floor(Math.random() * 4);
const appearanceValues = {
	variable1: randomString(),
	variable2: randomString(),
	variable3: randomString(),
};

const streamOutputVariants = {
	videoBitrate: randomNumber() * 100,
	framerate: 42,
	cpuUsageLevel: 2,
	scaledHeight: randomNumber() * 100,
	scaledWidth: randomNumber() * 100,
};
const newSocialHandles = [
	{
		url: 'http://facebook.org/' + randomString(),
		platform: randomString(),
	},
];

const newS3Config = {
	enabled: !defaultS3Config.enabled,
	endpoint: 'http://' + randomString() + ".tld",
	accessKey: randomString(),
	secret: randomString(),
	bucket: randomString(),
	region: randomString(),
	forcePathStyle: !defaultS3Config.forcePathStyle,
};

const newForbiddenUsernames = [randomString(), randomString(), randomString()];

const newYPConfig = {
	enabled: !defaultYPConfig.enabled,
	instanceUrl: 'http://' + randomString() + ".tld"
};

const newFederationConfig = {
	enabled: !defaultFederationConfig.enabled,
	isPrivate: !defaultFederationConfig.isPrivate,
	username: randomString(),
	goLiveMessage: randomString(),
	showEngagement: !defaultFederationConfig.showEngagement,
	blockedDomains: [randomString() + ".tld", randomString() + ".tld"],
};

const newHideViewerCount = !defaultHideViewerCount;

test('verify default streamKeys', async (done) => {
	const res = await getAdminConfig();

	expect(res.body.streamKeys).toStrictEqual(defaultStreamKeys);
	done();
});

test('verify default adminPassword', async (done) => {
	const res = await getAdminConfig();

	expect(res.body.adminPassword).toBe(defaultAdminPassword);
	done();
});

test('verify default directory configurations', async (done) => {
	const res = await getAdminConfig();

	expect(res.body.yp.enabled).toBe(defaultYPConfig.enabled);
	done();
});

test('verify default federation configurations', async (done) => {
	const res = await getAdminConfig();

	expect(res.body.federation.enabled).toBe(defaultFederationConfig.enabled);
	expect(res.body.federation.isPrivate).toBe(defaultFederationConfig.isPrivate);
	expect(res.body.federation.showEngagement).toBe(defaultFederationConfig.showEngagement);
	expect(res.body.federation.goLiveMessage).toBe(defaultFederationConfig.goLiveMessage);
	expect(res.body.federation.blockedDomains).toStrictEqual(defaultFederationConfig.blockedDomains);
	done();

});

test('set server name', async (done) => {
	const res = await sendConfigChangeRequest('name', newServerName);
	done();
});

test('set stream title', async (done) => {
	const res = await sendConfigChangeRequest('streamtitle', newStreamTitle);
	done();
});

test('set server summary', async (done) => {
	const res = await sendConfigChangeRequest('serversummary', newServerSummary);
	done();
});

test('set extra page content', async (done) => {
	const res = await sendConfigChangeRequest('pagecontent', newPageContent);
	done();
});

test('set tags', async (done) => {
	const res = await sendConfigChangeRequest('tags', newTags);
	done();
});

test('set stream keys', async (done) => {
	const res = await sendConfigChangeRequest('streamkeys', newStreamKeys);
	done();
});

test('set latency level', async (done) => {
	const res = await sendConfigChangeRequest(
		'video/streamlatencylevel',
		latencyLevel
	);
	done();
});

test('set video stream output variants', async (done) => {
	const res = await sendConfigChangeRequest('video/streamoutputvariants', [
		streamOutputVariants,
	]);
	done();
});

test('set social handles', async (done) => {
	const res = await sendConfigChangeRequest('socialhandles', newSocialHandles);
	done();
});

test('set s3 configuration', async (done) => {
	const res = await sendConfigChangeRequest('s3', newS3Config);
	done();
});

test('set forbidden usernames', async (done) => {
	const res = await sendConfigChangeRequest(
		'chat/forbiddenusernames',
		newForbiddenUsernames
	);
	done();
});

test('set server url', async (done) => {
	const res = await sendConfigChangeRequest('serverurl', newYPConfig.instanceUrl);
	done();
});

test('set federation username', async (done) => {
	const res = await sendConfigChangeRequest('federation/username', newFederationConfig.username);
	done();
});

test('set federation goLiveMessage', async (done) => {
	const res = await sendConfigChangeRequest('federation/livemessage', newFederationConfig.goLiveMessage);
	done();
});

test('set hide viewer count', async (done) => {
	const res = await sendConfigChangeRequest('hideviewercount', newHideViewerCount);
	done();
});

test('toggle private federation mode', async (done) => {
	const res = await sendConfigChangeRequest('federation/private', newFederationConfig.isPrivate);
	done();
});

test('toggle federation engagement', async (done) => {
	const res = await sendConfigChangeRequest('federation/showengagement', newFederationConfig.showEngagement);
	done();
});

test('set federation blocked domains', async (done) => {
	const res = await sendConfigChangeRequest('federation/blockdomains', newFederationConfig.blockedDomains);
	done();
});


test('set offline message', async (done) => {
	const res = await sendConfigChangeRequest('offlinemessage', newOfflineMessage);
	done();
});

test('set custom style values', async (done) => {
	const res = await sendConfigChangeRequest('appearance', appearanceValues);
	done();
});

test('enable directory', async (done) => {
	const res = await sendConfigChangeRequest('directoryenabled', true);
	done();
});

test('enable federation', async (done) => {
	const res = await sendConfigChangeRequest('federation/enable', newFederationConfig.enabled);
	done();
});

test('change admin password', async (done) => {
	const res = await sendConfigChangeRequest('adminpass', newAdminPassword);
	done();
});

test('verify admin password change', async (done) => {
	const res = await getAdminConfig(adminPassword = newAdminPassword);

	expect(res.body.adminPassword).toBe(newAdminPassword);
	done();
});

test('reset admin password', async (done) => {
	const res = await sendConfigChangeRequest('adminpass', defaultAdminPassword, adminPassword = newAdminPassword);
	done();
});

test('verify updated config values', async (done) => {
	const res = await request.get('/api/config');
	expect(res.body.name).toBe(newServerName);
	expect(res.body.streamTitle).toBe(newStreamTitle);
	expect(res.body.summary).toBe(`${newServerSummary}`);
	expect(res.body.extraPageContent).toBe(newPageContent);
	expect(res.body.offlineMessage).toBe(newOfflineMessage);
	expect(res.body.logo).toBe('/logo');
	expect(res.body.socialHandles).toStrictEqual(newSocialHandles);
	done();
});

// Test that the raw video details being broadcasted are coming through
test('verify admin stream details', async (done) => {

	const res = await getAdminStatus();

	expect(res.body.broadcaster.streamDetails.width).toBe(320);
	expect(res.body.broadcaster.streamDetails.height).toBe(180);
	expect(res.body.broadcaster.streamDetails.framerate).toBe(24);
	expect(res.body.broadcaster.streamDetails.videoBitrate).toBe(1269);
	expect(res.body.broadcaster.streamDetails.videoCodec).toBe('H.264');
	expect(res.body.broadcaster.streamDetails.audioCodec).toBe('AAC');
	expect(res.body.online).toBe(true);
	done();
});

test('verify updated admin configuration', async (done) => {
	const res = await getAdminConfig();

	expect(res.body.instanceDetails.name).toBe(newServerName);
	expect(res.body.instanceDetails.summary).toBe(newServerSummary);
	expect(res.body.instanceDetails.offlineMessage).toBe(newOfflineMessage);
	expect(res.body.instanceDetails.tags).toStrictEqual(newTags);
	expect(res.body.instanceDetails.socialHandles).toStrictEqual(
		newSocialHandles
	);
	expect(res.body.forbiddenUsernames).toStrictEqual(newForbiddenUsernames);
	expect(res.body.streamKeys).toStrictEqual(newStreamKeys);

	expect(res.body.videoSettings.latencyLevel).toBe(latencyLevel);
	expect(res.body.videoSettings.videoQualityVariants[0].framerate).toBe(
		streamOutputVariants.framerate
	);
	expect(res.body.videoSettings.videoQualityVariants[0].cpuUsageLevel).toBe(
		streamOutputVariants.cpuUsageLevel
	);

	expect(res.body.yp.enabled).toBe(true);
	expect(res.body.yp.instanceUrl).toBe(newYPConfig.instanceUrl);

	expect(res.body.adminPassword).toBe(defaultAdminPassword);

	expect(res.body.s3.enabled).toBe(newS3Config.enabled);
	expect(res.body.s3.endpoint).toBe(newS3Config.endpoint);
	expect(res.body.s3.accessKey).toBe(newS3Config.accessKey);
	expect(res.body.s3.secret).toBe(newS3Config.secret);
	expect(res.body.s3.bucket).toBe(newS3Config.bucket);
	expect(res.body.s3.region).toBe(newS3Config.region);
	expect(res.body.s3.forcePathStyle).toBe(newS3Config.forcePathStyle);
	expect(res.body.hideViewerCount).toBe(newHideViewerCount);

	expect(res.body.federation.enabled).toBe(newFederationConfig.enabled);
	expect(res.body.federation.isPrivate).toBe(newFederationConfig.isPrivate);
	expect(res.body.federation.username).toBe(newFederationConfig.username);
	expect(res.body.federation.goLiveMessage).toBe(newFederationConfig.goLiveMessage);
	expect(res.body.federation.showEngagement).toBe(newFederationConfig.showEngagement);
	expect(res.body.federation.blockedDomains).toStrictEqual(newFederationConfig.blockedDomains);
	done();

});

test('verify updated frontend configuration', (done) => {
	request
		.get('/api/config')
		.expect(200)
		.then((res) => {
			expect(res.body.name).toBe(newServerName);
			expect(res.body.logo).toBe('/logo');
			expect(res.body.socialHandles).toStrictEqual(newSocialHandles);
			done();
		});
});

test('verify frontend status', (done) => {
	request
		.get('/api/status')
		.expect(200)
		.then((res) => {
			expect(res.body.viewerCount).toBe(undefined);
			done();
		});
});


function randomString(length = 20) {
	return Random.value().toString(16).substr(2, length);
}

function randomNumber() {
	return Random.range(0, 5);
}
