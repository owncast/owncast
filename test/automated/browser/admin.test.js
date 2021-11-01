const listenForErrors = require('./lib/errors.js').listenForErrors;
const ADMIN_USERNAME = 'admin';
const ADMIN_PASSWORD = 'abc123';

describe('Admin page', () => {
  beforeAll(async () => {
    await page.setViewport({ width: 1080, height: 720 });
    listenForErrors(browser, page);

    // set HTTP Basic auth
    await page.authenticate({
      username: ADMIN_USERNAME,
      password: ADMIN_PASSWORD,
    });

    await page.goto('http://localhost:5309/admin');
  });

  afterAll(async () => {
    await page.waitForTimeout(3000);
    await page.screenshot({ path: 'screenshots/admin.png', fullPage: true });
  });

  it('should have rendered the admin home page', async () => {
    await page.waitForSelector('.home-container');
  });
});
