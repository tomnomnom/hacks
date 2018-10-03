# inscope

Prototype tool for filtering URLs and domains supplied on stdin to make sure they meet one of a set of regular expressions.

In theory this should be useful for filtering the output of other tools to only include items that are in scope for a bug bounty program.

## Install

```
▶ go get -u github.com/tomnomnom/hacks/inscope
```

## Usage

Pipe URLs and/or domains into it on stdin:

```
▶ cat testinput
https://example.com/footle
https://inscope.example.com/some/path?foo=bar
https://outofscope.example.net/bar
example.com
example.net

▶ cat testinput | inscope
https://example.com/footle
https://inscope.example.com/some/path?foo=bar
example.com
http://sub.example.com
```

## Scope Files

The tool reads regexes from a file called `.scope` in the current working directory.
If it doen't find one it recursively checks the parent directory until it hits the root.

Here's an example `.scope` file:

```
.*\.example\.com$
^example\.com$
.*\.example\.net$
!.*outofscope\.example\.net$
```

Each line is a regular expression to match domain names. When URLs are provided as input they
are parsed and only the hostname/domain portion is checked against the regex.

Line starting with `!` are treated as negative matches - i.e. any domain matching that regex will
be considered out of scope even if it matches one of the other regexes.
