const listenForErrors = require('./lib/errors.js').listenForErrors;

describe('Video embed page', () => {
  async function getMetaTagContent(property) {
    const selector = `meta[property="${property}"]`;

    const tag = await page.evaluate((selector) => {
      return document.head.querySelector(selector).getAttribute('content');
    }, selector);
    return tag;
  }

  beforeAll(async () => {
    await page.setViewport({ width: 1080, height: 720 });
    listenForErrors(browser, page);
    page.setUserAgent('Mastodon');
    await page.goto('http://localhost:5309');
  });

  afterAll(async () => {
    await page.waitForTimeout(5000);
    await page.screenshot({
      path: 'screenshots/screenshot_bots_share_search_scrapers.png',
      fullPage: true,
    });
  });

  it('should have rendered the simple bot accessible html page', async () => {
    await page.waitForSelector('h1');
    await page.waitForSelector('h3');

    const ogVideo = await getMetaTagContent('og:video');
    expect(ogVideo).toBe('http://localhost:5309/embed/video');

    const ogVideoType = await getMetaTagContent('og:video:type');
    expect(ogVideoType).toBe('text/html');

    // When stream is live the thumbnail is provided as the image.
    const ogImage = await getMetaTagContent('og:image');
    expect(ogImage).toBe('http://localhost:5309/thumbnail.jpg');

    const twitterUrl = await getMetaTagContent('twitter:url');
    expect(twitterUrl).toBe('http://localhost:5309/');

    const twitterImage = await getMetaTagContent('twitter:image');
    expect(twitterImage).toBe('http://localhost:5309/logo/external');
  });
});
