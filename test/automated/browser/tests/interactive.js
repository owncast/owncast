const puppeteer = require('puppeteer');
const browser = await puppeteer.launch()
const page = await browser.newPage()
await page.goto('http://localhost:8080/index-standalone-chat-readwrite.html')

await page.setViewport({ width: 1526, height: 754 })

await page.waitForSelector('#message-input')
await page.click('#message-input')

await page.type('#message-input', 'undefined')

await page.waitForSelector('#message-input')
await page.click('#message-input')

await page.type('#message-input', 'Test chat message 123.')

await page.waitForSelector('#username-display')
await page.click('#username-display')

await page.waitForSelector('#username-change-input')
await page.click('#username-change-input')

await page.type('#username-change-input', 'a new name')

await page.waitForSelector('#username-change-input')
await page.click('#username-change-input')

await page.waitForSelector('#username-change-input')
await page.click('#username-change-input')

await page.waitForSelector('#username-change-input')
await page.click('#username-change-input')

await page.type('#username-change-input', 'new name')

await page.waitForSelector('#button-update-username')
await page.click('#button-update-username')

await page.screenshot({ path: 'screenshot_1.png', fullPage: true })

await browser.close()