var server = document.getElementById('server')
var token = document.getElementById('token')
var add = document.getElementById('add')
var save = document.getElementById('save')
let snippets = document.getElementById('snippets')

chrome.storage.sync.get('config', function(data){
    if (!data.config){
        return
    }

    if (!data.config.snippets){
        data.config.snippets = [
            {"name": "", "code": "", "onsuccess": ""}
        ] 
    }

    server.value = data.config.server || 'localhost'
    token.value = data.config.token || 'notoken'

    data.config.snippets.map(s => {
        snippets.appendChild(snippetTemplate(s))
    })

})

add.addEventListener('click', () => {
    snippets.appendChild(snippetTemplate({}))
})

save.addEventListener('click', function(){
    let snipData = [...snippets.querySelectorAll('div.snippet')].map(el => {
        return {
            "name": el.querySelector('.name').value,
            "code": el.querySelector('.code').value,
            "onsuccess": el.querySelector('.onsuccess').value,
        } 
    }).filter(s => s.name && s.code)
    
    chrome.storage.sync.set({config: {
        "server": server.value,
        "token": token.value,
        "snippets": snipData,
    }}, function(){
        console.log('saved:', server.value, token.value)
        save.innerText = 'Saved!'
        setTimeout(() => {
            save.innerText = 'Save'
        }, 1000)
    })
})


function snippetTemplate(data){
    let dp = new DOMParser()
    let snippet = dp.parseFromString(`
        <div class=snippet>
            <label>Name:</label>
            <input type=text class=name>

            <label>Code (should return an array of strings):</label>
            <textarea class=code></textarea>

            <label>On Success (code to run after data has been sent):</label>
            <textarea class=onsuccess></textarea>

            <div> <button class=delete>Delete</button> </div>
        </div>
    `, 'text/html')

    snippet.querySelector('.name').value = data.name || ""
    snippet.querySelector('.code').value = data.code || ""
    snippet.querySelector('.onsuccess').value = data.onsuccess || ""
    snippet.querySelector('.delete').addEventListener('click', (e) => {
        let snip = e.target.parentNode.parentNode
        snip.parentNode.removeChild(snip)
    })

    return snippet.querySelector('.snippet')
}
