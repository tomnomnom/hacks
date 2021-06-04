# kxss

I don't know what this name is. There might even be something that has this name already,
but I'm on a plane with no internet connection so I can't check that right now. Also:
who said you needed StackOverflow to be able to write code? ¯\\\_(ツ)\_/¯

## Idea

So the general idea is:

* Take URLs with params on stdin. These might have come from waybackurls or maybe a Burp session
* Request the URLs, check the response body for any reflected parameters. There will be many false positives here.
* For any reflected parameters, re-request with some random alphanumeric value appended to the param
	* Only one param is appended to at a time. This is to avoid breaking the request when a different param is required

And the bit that's not done yet:
* For any params that passed the appended-to check (i.e. ones we're *really* sure are reflected), start trying special characters too (`<"'>` etc)

The "best" thing to do here would probably be to test one special character at a time, but it's really a
trade-off between the number of requests we're issuing and accuracy.

## Things

* It'd be nice to support POST params at some point too
* We're going to be generating a *lot* of requests, often to the same hosts
	* This needs some kind of rate-limiting with re-queing so that we don't overwhelm hosts, but can still cover lots of ground quickly

## Usage

At the moment there's a test server in `cmd/testserver`. Run that with `go run cmd/testserver/main.go` and
then do something like this:

```
▶ go build
▶ echo 'http://localhost:5566/?name=Tom&age=33&fake=Buck' | ./kxss
got reflection of appended param age on http://localhost:5566/?name=Tom&age=33&fake=Buck
got reflection of appended param name on http://localhost:5566/?name=Tom&age=33&fake=Buck
```

The source of the handler in the test server looks like this:

```
qs := r.URL.Query()
log.Printf("req: %#v", qs)
w.Header().Add("Content-Type", "text/html")
fmt.Fprintf(w, "Hello, %s!\n", qs.Get("name"))
fmt.Fprintf(w, "I hear you're %s years old!\n", qs.Get("age"))
fmt.Fprint(w, "My name is Buck and I'm here to greet you.\n")
```

Note that "Buck" appears in the query string and in the page output, but is not reported
as reflected. That's because we appended a random value to the end of it to double check.

## Payloads and context

As well as figuring out what special characters are allowed in reflected params, we need
to come up with a way to figure out what context (or contexts!) the reflected param appears in.
Really we need to figure out the context first and use that to prioritise the special character
checking part.

E.g. if we have something like this:

```
<h1>$reflectedParam</h1>
```

...then we need to check for `<` first as we can rule out XSS the fastest that way. We still
probably want to try multiple tricks for that char (e.g. double urlencoding, that funky %EF%BC%9C thing etc)

But if we have, say, this:

```
<iframe src=$reflectedParam></iframe>
```

...then we might not want to mess about with `<` and other special chars initially, but instead
look at `javascript:` type payloads because they're more likely to work.

A given context will have a variety of characters / payloads that we'll want to try.

Honestly: I think this bit is going to be Hard with a capital aitch; especially when the context
is inbetween `<script>` tags and HTML parsers can't save us.

Another option is to just blindly try lots of payloads. That feels a bit more barbaric, but
my goodness does it make for more simple code. Those funky payloads that work in lots of contexts
are probably worth playing with too, although they have an increased chance of being blocked by
WAFs or other filters too.

My suggestion is that for an MVP we settle for checking individual characters and having an output
that's something like:

```
o hai, i found reflected param X on url Y that allows these chars: <>"'
```

...and then letting the human go and investigate. We're basically never going to be able to be
as good as a human at figuring out the contexts etc, but it would be nice for the tool to at
least try some full payloads at some point in the future.

## Validation

If we do end up trying full payloads at any point I think basically the only way to reliably
validate that XSS fires is to use headless chrome (probably with `chromedp`) to load a URL
and then hook the alert boxes. This is resource intensive though so it should only be done
when we've exhausted what we can do with Go's HTTP client.
