<?php

$urls = [
    'https://httpbin.org/headers'
];

$params = [];

foreach ($urls as $url) {

    $host = parse_url($url, PHP_URL_HOST);
    $path = parse_url($url, PHP_URL_PATH);

    $context = stream_context_create($params);
    $fp = stream_socket_client(
        "tls://{$host}:443",
        $errno,
        $errstr,
        30,
        STREAM_CLIENT_CONNECT,
        $context
    );
    if (!$fp) continue;

    stream_set_write_buffer($fp, 0);
    stream_set_read_buffer($fp, 0);

    fwrite($fp, "GET {$path} HTTP/1.1\r\n");
    fwrite($fp, "Host: {$host}\r\n");
    fwrite($fp, "Connection: close\r\n");
    fwrite($fp, "Origin: {$host}\r\n");
    fwrite($fp, "\r\n");

    $resp = stream_get_contents($fp);

    echo $resp;
    fclose($fp);
}
