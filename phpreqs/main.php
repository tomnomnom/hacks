<?php


class SimpleRequest extends RawRequest {
    public $method   = 'GET';
    public $scheme   = 'http';
    public $path     = '/';
    public $query    = '';
    public $fragment = '';
    public $proto    = 'HTTP/1.1';

    public $headers = [];

    public function __construct($url, array $headers = []){
        $p = parse_url($url);

        $this->scheme = $p['scheme'] ?? $this->scheme;
        if (strToLower($this->scheme) == 'https'){
            $this->port = 443;
            $this->transport = 'tls';
        }

        $this->host     = $p['host']     ?? $this->host;
        $this->port     = $p['port']     ?? $this->port;
        $this->path     = $p['path']     ?? $this->path;
        $this->query    = $p['query']    ?? $this->query;
        $this->fragment = $p['fragment'] ?? $this->fragment;

        $this->addHeader("Host: {$this->host}");

        foreach ($headers as $header){
            $this->addHeader($header);
        }

    }

    public function addHeader($header){
        $this->headers[] = $header;
    }

    public function toString(){
        $out = "{$this->method} {$this->path} {$this->proto}\r\n";
        foreach ($this->headers as $header){
            $out .= "{$header}\r\n";
        }

        return $out;
    }
}

class RawRequest {
    public $transport = 'tcp';
    public $host      = 'localhost';
    public $port      = '80';
    public $request   = '';

    public function __construct($raw){
        $this->request = $raw;
    }

    public function send(){
        $sock = @stream_socket_client("{$this->transport}://{$this->host}:{$this->port}");
        if (!$sock) return new \Response(null);

        stream_set_write_buffer($sock, 0);
        stream_set_read_buffer($sock, 0);

        fwrite($sock, $this->toString());
        fwrite($sock, "\r\n");

        return new \Response($sock);
    }

    public function toString(){
        return $this->request;
    }
}

class Response {
    private $socket = null;
    private $rawStatus = '';
    private $headers = [];
    private $body = '';

    public function __construct($socket = null){
        $this->socket = $socket;
    }

    public function readHeaders(){
        if (!$this->socket) return;

        // HTTP/1.1 200 OK
        $status = $this->readLine();
        if ($status === false){
            return;
        }
        $this->rawStatus = trim($status);

        // Content-Type: text/html etc...
        while (true){
            $line = $this->readLine();
            if ($line === false) break;

            $line = trim($line);
            if ($line == "") break;

            $this->headers[] = $line;
        }
    }

    // Read the body from the response using either the Content-Length header to
    // specify how many bytes to read, or just read until the connection closes.
    public function readBody(){
        if ($this->hasHeader('Content-Length')){
            $this->body = $this->readBytes((int) $this->getHeader('Content-Length'));
        } else {
            $this->body = $this->readUntilClosed();
        }
    }

    private function readLine(){
        if (!$this->socket) return false;
        return fgets($this->socket);
    }

    private function readBytes($count){
        if (!$this->socket) return false;
        if ($count <= 0) return false;
        return fread($this->socket, $count);
    }

    private function readUntilClosed(){
        $out = '';
        while (true){
            $line = $this->readLine();
            if ($line === false) break;

            $out .= $line;
        }

        return $out;
    }

    public function close(){
        if (!$this->socket) return false;
        fclose($this->socket);
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

    public function toString(){
        $out = "{$this->rawStatus}\r\n";

        foreach ($this->headers as $header){
            $out .= "{$header}\r\n";
        }

        $out .= "\r\n";
        $out .= $this->body;

        return $out;
    }

}


$req = new SimpleRequest("http://httpbin.org/headers", [
    'Origin: http://example.com'
]);
$resp = $req->send();
$resp->readHeaders();
$resp->readBody();
$resp->close();

echo $req->toString();
echo PHP_EOL;
echo $resp->toString();

