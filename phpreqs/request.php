<?php
require_once __DIR__.'/rawrequest.php';

class Request extends RawRequest {
    public $method   = 'GET';
    public $scheme   = 'http';
    public $path     = '/';
    public $query    = '';
    public $fragment = '';
    public $proto    = 'HTTP/1.1';

    public $headers = [];

    public function __construct($url = "", array $headers = []){
        if (!$url) return;

        $p = parse_url($url);
        if (!$p) return;

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

        foreach ($headers as $header){
            $this->addHeader($header);
        }
    }

    public function setMethod($method){
        $this->method = $method;
        return $this;
    }

    public function setPath($path){
        $this->path = $path;
        return $this;
    }

    public function setProto($proto){
        $this->proto = $proto;
        return $this;
    }

    public function addHeader($header){
        $this->headers[] = $header;
        return $this;
    }

    public function setQuery($q){
        if (is_array($q)){
            $q = http_build_query($q);
        }
        $this->query = $q;
        return $this;
    }

    public function setFragment($fragment){
        $this->fragment = $fragment;
        return $this;
    }

    public function toString(){
        $path = $this->path;

        if ($this->query){
            $path .= "?{$this->query}";
        }

        if ($this->fragment){
            $path .= "#{$this->fragment}";
        }

        $out = "{$this->method} {$path} {$this->proto}{$this->eol}";
        foreach ($this->headers as $header){
            $out .= "{$header}{$this->eol}";
        }

        return $out;
    }
}
