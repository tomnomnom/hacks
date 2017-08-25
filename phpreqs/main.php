<?php
require_once __DIR__.'/request.php';


$req = new Request("http://httpbin.org");

$resp = $req->setMethod('POST')
    ->setPath('/anything')
    ->setQuery([
        'param1' => 'value1',
        'param2' => 'value2'
    ])
    ->setFragment('thefragment')
    ->setProto('HTTP/1.1')
    ->addHeader("Origin: http://evil.com")
    ->setEol("\r\n")
    ->send();

echo $req->toString();
echo PHP_EOL;
echo $resp->toString();

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

*/
