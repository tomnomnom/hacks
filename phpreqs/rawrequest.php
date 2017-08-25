<?php
require_once __DIR__.'/response.php';

class RawRequest {
    public $transport = 'tcp';
    public $host      = 'localhost';
    public $port      = '80';
    public $request   = '';
    public $eol       = "\r\n";

    public function __construct($raw){
        $this->request = $raw;
    }

    public function send(){
        $sock = @stream_socket_client("{$this->transport}://{$this->host}:{$this->port}");
        if (!$sock) return new \Response(null);

        stream_set_write_buffer($sock, 0);
        stream_set_read_buffer($sock, 0);

        fwrite($sock, $this->toString());
        fwrite($sock, $this->eol);

        $r = new \Response($sock);
        $r->readHeaders();
        $r->readBody();
        $r->close();
        
        return $r;
    }

    public function setTransport($transport){
        $this->transport = $transport;
        return $this;
    }

    public function setHost($host){
        $this->host = $host;
        return $this;
    }

    public function setPort($port){
        $this->port = $port;
        return $this;
    }

    public function setEol($eol){
        $this->eol = $eol;
        return $this;
    }

    public function toString(){
        return $this->request;
    }
}
