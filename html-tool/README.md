# html-tool

Take URLs or filenames for HTML documents on stdin and extract tag contents, attribute values, or comments.

## Install

```
▶ go get -u github.com/tomnomnom/hacks/html-tool
```

## Usage

```
▶ html-tool 
Accept URLs or filenames for HTML documents on stdin and extract parts of them.

Usage: html-tool <mode> [<args>]

Modes:
	tags <tag-names>        Extract text contained in tags
	attribs <attrib-names>  Extract attribute values
	comments                Extract comments

Examples:
	cat urls.txt | html-tool tags title a strong
	find . -type f -name "*.html" | html-tool attribs src href
	cat urls.txt | html-tool comments
```

## TODO
* Support selectors with https://github.com/ericchiang/css
* Option to ignore certificate errors
* Option to output filename/URL with output
