const listenForErrors = require('./lib/errors.js').listenForErrors;
const interactiveChatTest = require('./tests/chat.js').interactiveChatTest;
const videoTest = require('./tests/video.js').videoTest;
const puppeteer = require('puppeteer');

const phone = puppeteer.devices['iPhone 11'];
const tabletLandscape = puppeteer.devices['iPad landscape'];
const tablet = puppeteer.devices['iPad Pro'];
const desktop = {
  name: 'desktop',
  userAgent:
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36',
  viewport: {
    width: 1920,
    height: 1080,
    deviceScaleFactor: 1,
    isMobile: false,
    hasTouch: false,
    isLandscape: true,
  },
};

const devices = [desktop, phone, tablet, tabletLandscape];

describe('Frontend web page', () => {
  beforeAll(async () => {
    listenForErrors(browser, page);
    await page.goto('http://localhost:5309');
    await page.waitForTimeout(3000);
  });

  devices.forEach(async function (device) {
    const newName = 'frontend-browser-test-name-change-' + device.name;
    const fakeMessage =
      'this is a test chat message sent via the automated browser tests on the main web frontend from ' +
      device.name;

    interactiveChatTest(browser, page, newName, fakeMessage, device.name);
    videoTest(browser, page);

    await page.waitForTimeout(2000);
    await page.screenshot({
      path: 'screenshots/screenshot_main-' + device.name + '.png',
      fullPage: true,
    });
  });
});
