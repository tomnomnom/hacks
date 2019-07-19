chrome.runtime.onInstalled.addListener(function(){
    // Needs to check if config exists or not first
    //chrome.storage.sync.set({config: {
        //server: "localhost",
        //token: "notoken"
    //}}, function(){})

    chrome.declarativeContent.onPageChanged.removeRules(undefined, function(){
        chrome.declarativeContent.onPageChanged.addRules([{
            conditions: [new chrome.declarativeContent.PageStateMatcher({
                pageUrl: {hostEquals: 'github.com'}
            })],
            actions: [new chrome.declarativeContent.ShowPageAction()]
        }])
    })
})
