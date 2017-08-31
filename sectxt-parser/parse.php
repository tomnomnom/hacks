#!/usr/bin/env php
<?php

class SecurityTxt {
    const FIELD_CONTACT         = 'contact';
    const FIELD_ENCRYPTION      = 'encryption';
    const FIELD_DISCLOSURE      = 'disclosure';
    const FIELD_ACKNOWLEDGEMENT = 'acknowledgement';

    private $validDisclosures = [
        'full',
        'partial',
        'none'
    ];

    private $errors = [];
    private $comments = [];

    private $fields = [
        self::FIELD_CONTACT         => [],
        self::FIELD_ENCRYPTION      => [],
        self::FIELD_DISCLOSURE      => [],
        self::FIELD_ACKNOWLEDGEMENT => [],
    ];

    public function __construct($raw = ""){
        if ($raw){
            $this->parse($raw);
        }
    }

    public function parse($raw){
        $lines = explode("\n", $raw);

        if (sizeOf($lines) < 1){
            $this->addError("empty file");
            return false;
        }

        $n = 0;
        foreach ($lines as $line){
            $n++;

            // Empty line
            $line = trim($line);
            if (!$line) continue;

            // Comment
            if ($line[0] == "#"){
                $this->comments[] = $line;
                continue;
            }

            $parts = explode(":", $line, 2);
            if (sizeOf($parts) != 2){
                $this->addError("invalid input on line {$n}: {$line}");
                continue;
            }

            $option = strToLower($parts[0]);
            $value = trim($parts[1]);

            if (!$this->validateField($option, $value, $n)){
                continue;
            }

            $this->fields[$option][] = $value;
        }

        if (sizeOf($this->fields[self::FIELD_CONTACT]) < 1){
            $this->addError("does not contain at least one contact field");
        }

    }

    private function validateField($option, $value, $lineNo = 0){
        switch ($option){
            case self::FIELD_CONTACT:
                return $this->validateContact($option, $value, $lineNo);

            case self::FIELD_DISCLOSURE:
                return $this->validateDisclosure($option, $value, $lineNo);

            case self::FIELD_ENCRYPTION:
            case self::FIELD_ACKNOWLEDGEMENT:
                return $this->validateUri($option, $value, $lineNo);

            default:
                $this->addError("invalid option '{$option}' on line {$lineNo}");
        }
        return false;
    }

    private function validateContact($option, $value, $lineNo){
        $lower = strToLower($value);
        if (!(
            filter_var($value, FILTER_VALIDATE_URL) ||
            filter_var($value, FILTER_VALIDATE_EMAIL) ||
            $this->isValidPhoneNumber($value)
        )){
            $this->addError("invalid value '{$value}' for option '{$option}' on line {$lineNo}");
            return false;
        }
        return true;
    }

    private function validateDisclosure($option, $value, $lineNo){
        if (!in_array(strToLower($value), $this->validDisclosures)){
            $this->addError("invalid value '{$value}' for option '{$option}' on line {$lineNo}; must be one of [".implode(", ", $this->validDisclosures)."]");
            return false;
        }
        return true;
    }

    private function validateUri($option, $value, $lineNo){
        if (!filter_var($value, FILTER_VALIDATE_URL)){
            $this->addError("invalid URI '{$value}' for option '{$option}' on line {$lineNo}");
            return false;
        }
        return true;
    }

    private function isValidPhoneNumber($candidate){
        return (preg_match("/^\+[0-9\(\) -]+$/", $candidate) > 0);
    }

    private function addError($msg){
        $this->errors[] = $msg;
    }

    public function hasErrors(){
        return (sizeOf($this->errors) > 0);
    }

    public function errors(){
        return $this->errors;
    }

    public function hasComments(){
        return (sizeOf($this->comments) > 0);
    }

    public function comments(){
        return $this->comments;
    }

    public function hasContact(){
        return (sizeOf($this->fields[self::FIELD_CONTACT]) > 0);
    }

    public function contact(){
        return $this->fields[self::FIELD_CONTACT];
    }

    public function hasDisclosure(){
        return (sizeOf($this->fields[self::FIELD_DISCLOSURE]) > 0);
    }

    public function disclosure(){
        return $this->fields[self::FIELD_DISCLOSURE];
    }

    public function hasEncryption(){
        return (sizeOf($this->fields[self::FIELD_ENCRYPTION]) > 0);
    }

    public function encryption(){
        return $this->fields[self::FIELD_ENCRYPTION];
    }

    public function hasAcknowledgement(){
        return (sizeOf($this->fields[self::FIELD_ACKNOWLEDGEMENT]) > 0);
    }

    public function acknowledgement(){
        return $this->fields[self::FIELD_ACKNOWLEDGEMENT];
    }

    public function fields(){
        return $this->fields;
    }

}


$infile = $argv[1]?? die("usage {$argv[0]} <inputfile>\n");

$raw = file_get_contents($infile);
if (!$raw) die("failed to read {$infile}\n");

$sectxt = new SecurityTxt();
$sectxt->parse($raw);



if ($sectxt->hasErrors()){
    echo "errors:\n";
    foreach ($sectxt->errors() as $error){
        echo "\t{$error}\n";
    }
}

echo "comments:\n";
foreach ($sectxt->comments() as $comment){
    echo "\t{$comment}\n";
}

foreach ($sectxt->fields() as $option => $field){
    if (sizeOf($field) < 1) continue;

    echo "{$option}:\n";
    foreach ($field as $value){
        echo "\t{$value}\n";
    }
}
