let copy = document.getElementById('copy')

copy.addEventListener('click', function(e){

    chrome.storage.sync.get('config', function(data){
        if (!data.config){
            return
        }

        console.log('config', data.config)

        let server = data.config.server || 'localhost'
        let token = data.config.token || 'notoken'
        let snippet = data.config.snippet || '// Snippet'
        let postSnippet = data.config.postSnippet || ''

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
