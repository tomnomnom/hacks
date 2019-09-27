const puppeteer = require('puppeteer')
const readline = require('readline');

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout,
    terminal: false
});

var urls = []
rl.on('line', async (url) => {
    urls.push(url)
})

rl.on('close', async () => {
    const browser = await puppeteer.launch({ignoreHTTPSErrors: true})

    Promise.all(urls.map(url => {
        return new Promise(async (resolve) => {
            var page = await browser.newPage()
            await page.goto(url)
            var destination = await page.evaluate(() => {
                return {"domain": document.domain, "href": document.location.href}
            })

            var u = new URL(url)

            if (u.host != destination.domain){
                resolve(`${url} redirects to ${destination.href}`)
            } else {
                resolve(null)
            }
        })
    })).then((values) => {
        console.log(values.filter(v => v != null).join('\n')) 
        browser.close()
    })

})
