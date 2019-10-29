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

    let values = await Promise.all(urls.map(async url => {
            var page = await browser.newPage()
            await page.goto(url)
            var destination = await page.evaluate(() => {
                return {"domain": document.domain, "href": document.location.href}
            })

            var u = new URL(url)

            if (u.host != destination.domain){
                return `${url} redirects to ${destination.href}`
            } else {
                return null
            }
		})
    )
	console.log(values.filter(v => v != null).join('\n'))
	browser.close()
})

