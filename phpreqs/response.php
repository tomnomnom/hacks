<?php

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

        return $this;
    }

    // Read the body from the response using either the Content-Length header to
    // specify how many bytes to read, or just read until the connection closes.
    public function readBody(){
        if ($this->hasHeader('Content-Length')){
            $this->body = $this->readBytes((int) $this->getHeader('Content-Length'));
        } else {
            $this->body = $this->readUntilClosed();
        }

        return $this;
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
        return $this;
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

    public function getBody(){
        return $this->body;
    }

    public function bodyMatches($pattern){
        return (preg_match($pattern, $this->body) === 1);
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
