<?php
require_once __DIR__.'/request.php';


// You can create a new request using a URL and sane defaults
$req = new Request("https://httpbin.org");

// You *could* just call $req->send() now and it would work, but you can
// override just about anything about the request and it won't stop you.

$resp = $req
    // The transport is either tls or tcp
    ->setTransport('tls')

    // The actual host and port to connect to (i.e. not the Host header)
    ->setHost('httpbin.org')
    ->setPort('443')

    // Set the method; can be any string you want
    ->setMethod('POST')

    // Set the path; same deal: can be anything you like
    ->setPath('/anything')

    // setQuery will build a query string from a map, but you
    // can give it a raw string like "param1=value1" too
    ->setQuery([
        'param1' => 'value1',
        'param2' => 'value2'
    ])

    // The bit after the '#' on the path
    ->setFragment('thefragment')

    // Set the protocol... Again no restrictions
    ->setProto('HTTP/1.1')

    // Add some headers
    ->addHeader("Host: httpbin.org")
    ->addHeader("Origin: http://evil.com")

    // Set the end-of-line sequence
    ->setEol("\r\n")

    // Send the request
    ->send();

// The request and response both have toString methods to make them easier to inspect
echo $req->toString();
echo PHP_EOL;
echo $resp->toString();

// There's one or two convenience methods on the response object to help you assert
// stuff about the response. Open to suggestions for more
if ($resp->hasHeader('access-control-allow-origin')){
    echo "Found ACAO: ".$resp->getHeader('Access-Control-Allow-Origin').PHP_EOL;
}

var_dump($resp->getBody());
if ($resp->bodyMatches("/{.*}/m")){
    echo "It's probably JSON\n";
}


/*
E.g. output:

POST /anything?param1=value1&param2=value2#thefragment HTTP/1.1
Host: httpbin.org
Origin: http://evil.com

HTTP/1.1 200 OK
Connection: keep-alive
Server: meinheld/0.6.1
Date: Fri, 25 Aug 2017 15:10:22 GMT
Content-Type: application/json
Access-Control-Allow-Origin: http://evil.com
Access-Control-Allow-Credentials: true
X-Powered-By: Flask
X-Processed-Time: 0.00148296356201
Content-Length: 388
Via: 1.1 vegur

{
  "args": {
    "param1": "value1", 
    "param2": "value2"
  }, 
  "data": "", 
  "files": {}, 
  "form": {}, 
  "headers": {
    "Connection": "close", 
    "Content-Length": "0", 
    "Host": "httpbin.org", 
    "Origin": "http://evil.com"
  }, 
  "json": null, 
  "method": "POST", 
  "origin": "123.123.123.123", 
  "url": "http://httpbin.org/anything?param1=value1&param2=value2"
}
Found ACAO: http://evil.com
*/
