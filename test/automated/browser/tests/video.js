async function videoTest(browser, page) {
  it('should have the video container element', async () => {
    await page.waitForSelector('#video-container');
  });

  it('should have the stream info status bar', async () => {
    await page.waitForSelector('#stream-info');
  });
}

module.exports.videoTest = videoTest;
