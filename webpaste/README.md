Installing Webpaste:

1. Clone hacks repo: https://github.com/tomnomnom/hacks
2. Go to webpaste directory, run go build, and you will get "webpaste" binary file.
3. Before starting webpaste, set environment variable: export WEBPASTE_TOKEN=iloveweb
4. By default webpaste runs on port 8080,

```
$ ./webpaste -h
Usage of ./webpaste:
  -a string
        address to listen on (default "0.0.0.0")
  -p string
        port to listen on (default "8080")
  -u    only print unique lines
```

Installing Extension:

1. Open the Extension Manager by following:

Kebab menu(three vertical dots) -> More Tools -> Extensions

2. If the developer mode is not turned on, turn it on by clicking the toggle in the top right corner.

3. Now click on Load unpacked button on the top left

4. Go the directory where you have webpaste/extension folder and select it.

5. Extension is loaded now.

6. Right click on chrome extension, go to "Options"

Put Server name:

http://localhost:8080 or http://ip:port

Same token as set above example: iloveweb

For Snippets, cloned directory has Google and JS extraction snippets.
```
cat hacks/webpaste/extension/snippets.js
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
```
Copy only the values inside " " in expension option, Save.

7. open google, search something: example: site:yahoo.com url:?

Run webpaste binary on terminal, Click on webpaste extension, click on "Google URLs", and you will see URL from the google search engine in your terminal.

