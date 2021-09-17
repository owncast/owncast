async function interactiveChatTest(browser, page, newName, chatMessage) {
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

  it('should allow changing the username', async () => {
    await page.waitForSelector('#username-display');
    await page.click('#username-display');

    await page.waitForSelector('#username-change-input');
    await page.click('#username-change-input');

    await page.type('#username-change-input', 'a new name');

    await page.waitForSelector('#username-change-input');
    await page.click('#username-change-input');

    await page.waitForSelector('#username-change-input');
    await page.click('#username-change-input');

    await page.waitForSelector('#username-change-input');

    await page.click('#username-change-input');

    await page.keyboard.press('Control');
    await page.keyboard.press('a');

    await page.type('#username-change-input', newName);

    await page.waitForSelector('#button-update-username');
    await page.click('#button-update-username');
  });

  it('should allow typing a chat message', async () => {
    await page.waitForTimeout(2000);

    await page.waitForSelector('#message-input');
    await page.click('#message-input');

    await page.type('#message-input', 'undefined');

    await page.waitForSelector('#message-input');
    await page.click('#message-input');
    await page.keyboard.press('Control');
    await page.keyboard.press('a');
    await page.type('#message-input', chatMessage);
    await page.keyboard.press('Enter');
  });
}

module.exports.interactiveChatTest = interactiveChatTest;
