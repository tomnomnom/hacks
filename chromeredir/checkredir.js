const puppeteer = require('puppeteer');

if (!process.argv[2]) {
    console.log("usage: node checkredir.js <URL>")
    process.exit(1);
}

url = process.argv[2];

(async () => {
  const browser = await puppeteer.launch({ignoreHTTPSErrors: true});
  const page = await browser.newPage();
  await page.goto(url);

  const finalLocation = await page.evaluate(() => {
    return document.location.href;
  });

  console.log('Final URL:', finalLocation);

  await browser.close();
})();
