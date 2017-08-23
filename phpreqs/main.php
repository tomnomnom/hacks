<?php

$hosts = [
    'https://httpbin.org',
    'http://httpbin.org'
];

foreach ($hosts as $host) {
    $resp = req($host, "/", [
        "Host: httpbin.org",
        "Origin: {$host}"
    ]);

    if ($resp->hasHeader('Access-Control-Allow-Origin')){
        echo "{$host}: ".$resp->getRawHeader('access-control-allow-origin').PHP_EOL;
    }

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

    $fp = @stream_socket_client("{$transport}://{$host}:{$port}");
    if (!$fp) return new \Response("[error]");

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

    $raw = stream_get_contents($fp);
    fclose($fp);

    return new \Response($raw);
}

class Response {
    public $rawStatus = '';
    public $headers = [];
    public $body = '';

    public function __construct($raw){
        $this->parseRaw($raw);
    }

    public function parseRaw($raw){

        $lines = explode("\n", $raw);
        $this->rawStatus = array_shift($lines);

        while ($line = trim(array_shift($lines))){

            if ($line == ''){
                // Hit the body
                break;
            }

            $this->headers[] = $line;
        }

        $this->body = implode("\n", $lines);
    }

    public function getRawHeader($search){
        $search = $this->normaliseString($search);

        foreach ($this->headers as $header){
            $parts = explode(":", $header, 2);
            if (sizeOf($parts) != 2){
                continue;
            }

            $key = $this->normaliseString($parts[0]);

            if ($search == $key){
                return $header;
            }
        }
        return "";
    }

    public function hasHeader($header){
        return ($this->getRawHeader($header) != "");
    }

    public function getHeader($header){
        $raw = $this->getRawHeader($header);
        $parts = explode(":", $raw, 2);

        if (sizeOf($parts) != 2){
            return "";
        }

        return trim($parts[1]);
    }

    private function normaliseString($str){
        return strToLower(trim($str));
    }

}











