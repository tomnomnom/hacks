const puppeteer = require('puppeteer');

puppeteer.launch({ignoreHTTPSErrors: true}).then(async browser => {
    const page = await browser.newPage();
    await page.setRequestInterception(true);
    await page.setUserAgent('Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.81 Safari/537.36');

    page.on('request', interceptedRequest => {
        interceptedRequest.continue();
    });

    page.on('response', resp => {
        if (resp._headers['content-type'] == undefined){
            return;
        }

        if (
            resp._headers['content-type'].match(/(javascript|json)/i) 
        ){
            console.log(resp.url());
        }
    });

    let url = process.argv[2];
    await page.goto(url);
    await browser.close();
}).catch(err => {
    process.exit();
});
