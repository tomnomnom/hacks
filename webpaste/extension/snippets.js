[
    {
        "name": "Google URLs",
        "code": "[...document.querySelectorAll('div.r>a:first-child')].map(n=>n.href)",
        "onsuccess": "document.location=document.querySelectorAll('a#pnnext')[0].href;"
    },
    {
        "name": "GitHub Code Results",
        "code": "[...document.querySelectorAll('#code_search_results a.text-bold')].map(n=>n.href)",
        "onsuccess": "document.location=document.querySelectorAll('a.next_page')[0].href;"
    },
    {
        "name": "Shodan Results",
        "code": "[...document.querySelectorAll('div.ip>a.fa')].map(n=>n.href)",
        "onsuccess": "try {document.location = document.querySelectorAll('a.btn')[1].href;}catch(e){if (e instanceof TypeError){document.location=document.querySelectorAll('a.btn')[0].href;} }"
    },
    {
        "name": "Censys Results",
        "code": "[...document.querySelectorAll('span.ip>a')].map(n=>n.href)",
        "onsuccess": "document.location=document.querySelectorAll('li.hover>a:first-child')[0].href;"
    },
    {
        "name": "Gist Results",
        "code": "[...document.querySelectorAll('#code_search_results a.text-bold')].map(n=>n.href)",
        "onsuccess": "document.location=document.querySelectorAll('a.next_page')[0].href;"
    },
    {
        "name": "Yandex",
        "code": "[...document.querySelectorAll('div.organic__path a.link[href]')].map(el => el.href)",
        "onsuccess": "document.location=document.querySelector('a.pager__item_kind_next').href;"
    }
]
