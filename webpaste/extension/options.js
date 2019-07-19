var server = document.getElementById('server')
var token = document.getElementById('token')
var save = document.getElementById('save')

chrome.storage.sync.get('config', function(data){
    if (data.config){
        server.value = data.config.server || 'localhost'
        token.value = data.config.token || 'notoken'
    }
})

save.addEventListener('click', function(){
    chrome.storage.sync.set({config: {
        "server": server.value,
        "token": token.value,
    }}, function(){
        console.log('saved', server.value, token.value)
    })
})

