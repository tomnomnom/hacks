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

            switch ($option){
                case self::FIELD_CONTACT:
                    $lower = strToLower($value);
                        if (!(
                            filter_var($value, FILTER_VALIDATE_URL) ||
                            filter_var($value, FILTER_VALIDATE_EMAIL) ||
                            $this->validatePhone($value)
                        )){
                            $this->addError("invalid value '{$value}' for option '{$parts[0]}' on line {$n}");
                            continue 2;
                        }

                    break;

                case self::FIELD_DISCLOSURE:
                    if (!in_array(strToLower($value), $this->validDisclosures)){
                        $this->addError("invalid value '{$value}' for option '{$parts[0]}' on line {$n}; must be one of [".implode(", ", $this->validDisclosures)."]");
                        continue 2;
                    }
                    break;

                case self::FIELD_ENCRYPTION:
                case self::FIELD_ACKNOWLEDGEMENT:
                    if (!filter_var($value, FILTER_VALIDATE_URL)){
                        $this->addError("invalid URI '{$value}' for option '{$parts[0]}' on line {$n}");
                        continue 2;
                    }
                    break;


                default:
                    $this->addError("invalid option '{$parts[0]}' on line {$n}");
                    continue 2;
                    break;
            }

            $this->fields[$option][] = $value;
        }

        if (sizeOf($this->fields[self::FIELD_CONTACT]) < 1){
            $this->addError("does not contain at least one contact field");
        }

    }

    public function comments(){
        return $this->comments;
    }

    public function fields(){
        return $this->fields;
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

    private function validatePhone($candidate){
        return (preg_match("/^\+[0-9\(\) -]+$/", $candidate) > 0);
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
