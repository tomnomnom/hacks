let buttons = document.getElementById('buttons')

chrome.storage.sync.get('config', function(data){
    if (!data.config || !data.config.snippets || data.config.snippets.length < 1){
        buttons.innerText = "No snippets set"
        return
    }

    data.config.snippets.map(s => {
        console.log(s)
        buttons.appendChild(buttonTemplate(s))
    })
})


buttons.addEventListener('click', function(e){

    chrome.storage.sync.get('config', function(data){
        if (!data.config){
            return
        }

        console.log('config', data.config)

        let server = data.config.server || 'localhost'
        let token = data.config.token || 'notoken'
        let snippet = e.target.dataset.code
        let postSnippet = e.target.dataset.onsuccess

        chrome.tabs.query({active: true, currentWindow: true}, function(tabs){
            chrome.tabs.executeScript(
                tabs[0].id,
                {code: snippet},

                function(results){
                    console.log('results', results)

                    fetch(server, {
                        method: 'POST',
                        mode: 'cors',
                        headers: {'Content-Type': 'application/json'},
                        body: JSON.stringify({
                            "token": token,
                            "lines": results[0]
                        })
                    }).then(() => {
                        chrome.tabs.executeScript(
                            tabs[0].id,
                            {code: postSnippet}
                        )
                    }).catch((err) => {
                        alert(err)
                    })
                }
            )
        })
    })

})

function buttonTemplate(data){
    let dp = new DOMParser()
    let button = dp.parseFromString(`
        <button></button>
    `, 'text/html').querySelector('button')

    button.innerText = data.name
    button.value = data.name
    button.dataset.code = data.code
    button.dataset.onsuccess = data.onsuccess

    return button
}
