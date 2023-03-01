const fs = require('fs');
const BrowserStack = require('browserstack');
const https = require('https');

let portraitJobId = null;
let landscapeJobId = null;
let finished = false;

const testBrowsersPortrait = [
	{
		os: 'OS X',
		os_version: 'Ventura',
		browser: 'chrome',
		device: null,
		browser_version: '71.0',
		real_mobile: null,
	},
	{
		os: 'OS X',
		os_version: 'Ventura',
		browser: 'safari',
		device: null,
		browser_version: '16.0',
		real_mobile: null,
	},
	{
		os: 'OS X',
		os_version: 'Ventura',
		browser: 'firefox',
		device: null,
		browser_version: '89.0',
		real_mobile: null,
	},
	{
		os: 'Windows',
		os_version: '11',
		browser: 'firefox',
		device: null,
		browser_version: '89.0',
		real_mobile: null,
	},
	{
		os: 'Windows',
		os_version: '11',
		browser: 'chrome',
		device: null,
		browser_version: '71.0',
		real_mobile: null,
	},
	{
		os: 'ios',
		os_version: '16',
		browser: 'Mobile Safari',
		device: 'iPhone 14 Pro',
		browser_version: null,
		real_mobile: true,
	},
	{
		os: 'ios',
		os_version: '16',
		browser: 'Mobile Safari',
		device: 'iPad Pro 11 2022',
		browser_version: null,
		real_mobile: true,
	},
	{
		os: 'android',
		os_version: '13.0',
		browser: 'Android Browser',
		device: 'Google Pixel 7 Pro',
		browser_version: null,
		real_mobile: true,
	},
];

const testBrowsersLandscape = [
	{
		os: 'android',
		os_version: '10.0',
		browser: 'Android Browser',
		device: 'Samsung Galaxy S20 Ultra',
		browser_version: null,
		real_mobile: true,
	},
	{
		os: 'ios',
		os_version: '16',
		browser: 'Mobile Safari',
		device: 'iPad Pro 11 2022',
		browser_version: null,
		real_mobile: true,
	},
	// {
	// 	os: 'ios',
	// 	os_version: '16',
	// 	browser: 'Mobile Safari',
	// 	device: 'iPhone 14',
	// 	browser_version: null,
	// 	real_mobile: true,
	// },
];

const USERNAME = process.env.BROWSERSTACK_USERNAME;
const PASSWORD = process.env.BROWSERSTACK_PASSWORD;
const URL = process.env.TEST_URL;
const FILE_SUFFIX = process.env.FILE_SUFFIX;

if (!URL) {
	console.error('Missing URL');
	process.exit(1);
}

if (!USERNAME || !PASSWORD) {
	console.error('Missing BrowserStack credentials');
	process.exit(1);
}

if (!FILE_SUFFIX) {
	console.error('Missing file suffix');
	process.exit(1);
}

const { exit } = require('process');
var browserStackCredentials = {
	username: USERNAME,
	password: PASSWORD,
};

let screenshotClient = BrowserStack.createScreenshotClient(
	browserStackCredentials
);

function start() {
	console.log(
		`creating screenshots for ${URL} using ${
			[...testBrowsersPortrait, ...testBrowsersLandscape].length
		} browsers...`
	);

	const portraitOptions = {
		url: URL,
		orientation: 'portrait',
		mac_res: '1920x1080',
		win_res: '1280x1024',
		browsers: testBrowsersPortrait,
		wait_time: 20,
		local: true,
	};

	const landscapeOptions = {
		url: URL,
		orientation: 'landscape',
		mac_res: '1920x1080',
		win_res: '1280x1024',
		browsers: testBrowsersLandscape,
		wait_time: 20,
		local: true,
	};

	screenshotClient.generateScreenshots(
		portraitOptions,
		function (error, response) {
			if (error) {
				console.error(error);
				exit(0);
			} else {
				const { job_id } = response;
				portraitJobId = job_id;
			}
		}
	);

	screenshotClient.generateScreenshots(
		landscapeOptions,
		function (error, response) {
			if (error) {
				console.error(error);
				exit(0);
			} else {
				const { job_id } = response;
				landscapeJobId = job_id;
			}
		}
	);

	setTimeout(check, 30000);
}

async function checkJob(jobId) {
	return new Promise((resolve, reject) => {
		screenshotClient.getJob(jobId, function (error, job) {
			if (error) {
				return reject(error);
			}

			return resolve(job);
		});
	});
}

async function check() {
	if (finished) {
		exit(0);
	}

	if (!portraitJobId || !landscapeJobId) {
		setTimeout(check, 5000);
		return;
	}

	try {
		const portraitStatus = await checkJob(portraitJobId);
		const landscapeStatus = await checkJob(landscapeJobId);

		const { screenshots: portraitScreenshots } = portraitStatus;
		const { screenshots: landscapeScreenshots } = landscapeStatus;
		const screenshots = [...portraitScreenshots, ...landscapeScreenshots];

		const completed = screenshots.filter((s) => s.state === 'done');
		console.log(`completed ${completed.length} of ${screenshots.length}...`);
		console.log(screenshots.filter((s) => s.state !== 'done'));

		if (completed.length === screenshots.length) {
			complete(screenshots);
		}

		const timedOut = screenshots.filter((s) => s.state === 'timed-out');
		if (timedOut.length > 0) {
			console.log('timed out');
			exit(1);
		}

		setTimeout(check, 5000);
	} catch (e) {
		console.error(e);
		exit(1);
	}
}

start();

function complete(screenshots) {
	console.log('downloading screenshots...');
	var saveCounter = 0;
	for (const screenshot of screenshots) {
		const { os, os_version, browser, device, orientation, image_url } =
			screenshot;
		const name = `${os} ${os_version} ${browser} ${device || 'desktop'} ${
			orientation || 'default'
		} ${FILE_SUFFIX}`;

		const filename = name.replace(/ /g, '-').toLowerCase();

		const file = fs.createWriteStream(`./screenshots/${filename}.png`);
		https
			.get(image_url, function (response) {
				response.pipe(file);

				file.on('finish', () => {
					file.close();
					saveCounter++;
					if (saveCounter === screenshots.length) {
						finished = true;
					}
				});
			})
			.on('error', function (err) {
				console.error(err);
			});
	}
}
