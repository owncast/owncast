async function interactiveChatTest(
  browser,
  page,
  newName,
  chatMessage,
  device
) {
  it('should have the chat input', async () => {
    await page.waitForSelector('#message-input');
  });

  it('should have the chat input enabled', async () => {
    const isDisabled = await page.evaluate(
      'document.querySelector("#message-input").getAttribute("disabled")'
    );
    expect(isDisabled).not.toBe('true');
  });

  it('should have the username label', async () => {
    await page.waitForSelector('#username-display');
  });

  it('should allow changing the username on ' + device, async () => {
    await page.waitForSelector('#username-display');
    await page.evaluate(() =>
      document.querySelector('#username-display').click()
    );
    await page.waitForSelector('#username-change-input');
    await page.evaluate(() =>
      document.querySelector('#username-change-input').click()
    );
    await page.evaluate(() =>
      document.querySelector('#username-change-input').click()
    );

    await page.type('#username-change-input', newName);

    await page.waitForSelector('#button-update-username');

    await page.evaluate(() =>
      document.querySelector('#button-update-username').click()
    );

    // page.keyboard.press('Enter');
    await page.waitForTimeout(3000);
  });

  it('should allow typing a chat message', async () => {
    await page.waitForSelector('#message-input');
    await page.evaluate(() => document.querySelector('#message-input').click());
    await page.waitForTimeout(1000);
    await page.focus('#message-input');
    await page.keyboard.type(chatMessage);
    page.keyboard.press('Enter');
  });
}

module.exports.interactiveChatTest = interactiveChatTest;
