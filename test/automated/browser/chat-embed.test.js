const listenForErrors = require('./lib/errors.js').listenForErrors;
const interactiveChatTest = require('./tests/chat.js').interactiveChatTest;

describe('Chat read-write embed page', () => {
  beforeAll(async () => {
    await page.setViewport({ width: 600, height: 700 });
    listenForErrors(browser, page);
    await page.goto('http://localhost:5309/embed/chat/readwrite');
  });

  afterAll(async () => {
    await page.waitForTimeout(3000);
    await page.screenshot({ path: 'screenshots/screenshot_chat_embed.png', fullPage: true });
  });

  const newName = 'frontend-browser-embed-test-name-change';
  const fakeMessage = 'this is a test chat message sent via the automated browser tests on the read/write chat embed page.'

  interactiveChatTest(browser, page, newName, fakeMessage);
});

describe('Chat read-only embed page', () => {
  beforeAll(async () => {
    await page.setViewport({ width: 500, height: 700 });
    listenForErrors(browser, page);
    await page.goto('http://localhost:5309/embed/chat/readonly');
  });

  it('should have the messages container', async () => {
    await page.waitForSelector('#messages-container');
  });

});
