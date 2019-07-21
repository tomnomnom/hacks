var server = document.getElementById('server')
var token = document.getElementById('token')
var snippet = document.getElementById('snippet')
var postSnippet = document.getElementById('postSnippet')
var save = document.getElementById('save')

chrome.storage.sync.get('config', function(data){
    if (data.config){
        server.value = data.config.server || 'localhost'
        token.value = data.config.token || 'notoken'
        snippet.value = data.config.snippet || '// Snippet'
        postSnippet.value = data.config.postSnippet || ''
    }
})

save.addEventListener('click', function(){
    chrome.storage.sync.set({config: {
        "server": server.value,
        "token": token.value,
        "snippet": snippet.value,
        "postSnippet": postSnippet.value,
    }}, function(){
        console.log('saved', server.value, token.value)
    })
})

