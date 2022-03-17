async function listenForErrors(browser, page) {
  const ignoredErrors = [
    'ERR_ABORTED',
    'MEDIA_ERR_SRC_NOT_SUPPORTED',
    '404',
    'JSHandle@error',
  ];

  // Emitted when the page emits an error event (for example, the page crashes)
  page.on('error', (error) => {
    throw new Error(`❌ ${error}`);
  });

  browser.on('error', (error) => {
    throw new Error(`❌ ${error}`);
  });

  // Emitted when a script within the page has uncaught exception
  page.on('pageerror', (error) => {
    throw new Error(`❌ ${error}`);
  });

  // Catch all failed requests like 4xx..5xx status codes
  page.on('requestfailed', (request) => {
    const ignoreError = ignoredErrors.some((e) =>
      request.failure().errorText.includes(e)
    );
    if (!ignoreError) {
      throw new Error(
        `❌ url: ${request.url()}, errText: ${
          request.failure().errorText
        }, method: ${request.method()}`
      );
    }
  });

  // Listen for console errors in the browser.
  page.on('console', (msg) => {
    const type = msg._type;
    if (type !== 'error') {
      return;
    }

    const ignoreError = ignoredErrors.some((e) => msg._text.includes(e));
    if (!ignoreError) {
      throw new Error(`❌ ${msg._text}`);
    }
  });
}

module.exports.listenForErrors = listenForErrors;
