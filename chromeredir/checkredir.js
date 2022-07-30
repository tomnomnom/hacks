const puppeteer = require('puppeteer')
const readline = require('readline');

const MAX_WORKERS = 20
const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout,
    terminal: false
});

var urls = []
var reading = true
rl.on('line', async (url) => {
    urls.push(url)                      // start queuing the read urls right away
});

(async ()=> {                           // kick off an async "thread" to read from the queue
const browser = await puppeteer.launch({ignoreHTTPSErrors: true})   // build the browser once
let working = new Set()                 // maybe not the most memory efficient to make two datastructures
for (let i = 0; i < MAX_WORKERS; i++) { // but the list as a queue is helpful and the set is helpful for different reasons
    (async ()=>{                        // only make MAX_WORKERS tasks ("threads") so we do not crash chrome
        const page = await browser.newPage()
        while (urls.length) {
            let url = urls.shift()      // grab the first URL
            working.add(url)            // mark that we are working on that URL
            try {
                await page.goto(url)
                var destination = await page.evaluate(() => {
                    return {"domain": document.domain, "href": document.location.href}
                })

                var u = new URL(url)

                if (u.host != destination.domain){
                    console.log(`${url} redirects to ${destination.href}`)
                } else {
                    console.log(`${url} does not redirect`)
                }
            } catch {
                // should an error just pass?
                console.log(`error checking ${url}`)
            } finally {
                working.delete(url)     // we are no longer working on that URL
            }
        }
        await page.close()              // clean up the page object, potentially an issue if page crashes in loop
        if (!reading && !working.size) {// I think this will prevent premature browser closure and issues with list/set desync
            browser.close()
        }
    })()
}
})()

rl.on('close', async () => {
    reading = false                     // make sure that our queue and set do not get desynced
})
