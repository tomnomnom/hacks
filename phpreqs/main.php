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

