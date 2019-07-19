let copy = document.getElementById('copy')

chrome.storage.sync.get('config', function(data){
    if (!data.config){
        return
    }

    let server = data.config.server || 'localhost'
    let token = data.config.token || 'notoken'

    copy.addEventListener('click', function(e){
        chrome.tabs.query({active: true, currentWindow: true}, function(tabs){
            chrome.tabs.executeScript(
                tabs[0].id,
                {code: '[...document.querySelectorAll("#code_search_results a.text-bold")].map(n=>n.href)'},

                function(results){
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
                            {code: 'document.location=document.querySelectorAll("a.next_page")[0].href;'}
                        )
                    }).catch((err) => {
                        alert(err)
                    })
                }
            )
        })
    })

})
