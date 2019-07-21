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
]
