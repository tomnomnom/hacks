[
    {
        "name": "Google URLs",
        "snippet": "[...document.querySelectorAll('div.r>a:first-child')].map(n=>n.href)",
        "postSnippet": "document.location=document.querySelectorAll('a#pnnext')[0].href;"
    },

    {
        "name": "GitHub Code Results",
        "snippet": "[...document.querySelectorAll('#code_search_results a.text-bold')].map(n=>n.href)",
        "postSnippet": "document.location=document.querySelectorAll('a.next_page')[0].href;"
    },
]
