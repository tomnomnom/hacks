<?php

$hosts = [
    //'https://httpbin.org',
    'http://httpbin.org'
];

foreach ($hosts as $host) {
    echo req($host, "/headers", [
        'Host: httpbin.org',
        'Origin: null'
    ]);
}

function req($host, $path = "/", array $headers = []){

    $p = parse_url($host);
    $host = $p['host'] ?? 'localhost';
    $port = $p['port'] ?? 80;
    $scheme = $p['scheme'] ?? 'http';

    if ($port == 80 && $scheme == 'https'){
        $port = 443;
    }

    $transport = "tcp";
    if ($scheme == 'https'){
        $transport = 'tls';
    }

    $fp = stream_socket_client("{$transport}://{$host}:{$port}");
    if (!$fp) return "";

    stream_set_write_buffer($fp, 0);
    stream_set_read_buffer($fp, 0);

    // This needs to support keep-alive rubbish too
    // But that means parsing the content length header
    // and... Nah.
    fwrite($fp, "GET {$path} HTTP/1.1\r\n");
    fwrite($fp, "Connection: close\r\n");

    foreach ($headers as $header){
        fwrite($fp, "{$header}\r\n");
    }
    fwrite($fp, "\r\n");

    $resp = stream_get_contents($fp);
    fclose($fp);

    return $resp;
}
