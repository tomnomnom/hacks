#!/usr/bin/env php
<?php

const FIELD_CONTACT         = 'contact';
const FIELD_ENCRYPTION      = 'encryption';
const FIELD_DISCLOSURE      = 'disclosure';
const FIELD_ACKNOWLEDGEMENT = 'acknowledgement';

$infile = $argv[1]?? die("usage {$argv[0]} <inputfile>\n");

$raw = file_get_contents($infile);
if (!$raw) die("failed to read {$infile}\n");

$lines = explode("\n", $raw);
if (sizeOf($lines) < 1) die("empty file\n");

$errors = [];
$comments = [];
$fields = [
    FIELD_CONTACT         => [],
    FIELD_ENCRYPTION      => [],
    FIELD_DISCLOSURE      => [],
    FIELD_ACKNOWLEDGEMENT => [],
];

$validDisclosures = [
    'full',
    'partial',
    'none'
];

$n = 0;
foreach ($lines as $line){
    $n++;

    // Empty line
    $line = trim($line);
    if (!$line) continue;

    // Comment
    if ($line[0] == "#"){
        $comments[] = $line;
        continue;
    }

    $parts = explode(":", $line, 2);
    if (sizeOf($parts) != 2){
        $errors[] = "invalid input on line {$n}: {$line}";
        continue;
    }

    $option = strToLower($parts[0]);
    $value = trim($parts[1]);

    switch ($option){
        case FIELD_CONTACT:
            // TODO: actual validation
            break;

        case FIELD_DISCLOSURE:
            if (!in_array(strToLower($value), $validDisclosures)){
                $errors[] = "invalid value '{$value}' for option '{$parts[0]}' on line {$n}; must be one of [".implode(", ", $validDisclosures)."]";
                continue 2;
            }
            break;

        case FIELD_ENCRYPTION:
        case FIELD_ACKNOWLEDGEMENT:
            if (!filter_var($value, FILTER_VALIDATE_URL)){
                $errors[] = "invalid URI '{$value}' for option '{$parts[0]}' on line {$n}"; 
                continue 2;
            }
            break;

        
        default:
            $errors[] = "invalid option '{$parts[0]}' on line {$n}";
            continue 2;
            break;
    }

    $fields[$option][] = $value;
}

if (sizeOf($fields[FIELD_CONTACT]) < 1){
    $errors[] = "does not contain at least one contact field";
}

if (sizeOf($errors) > 0){
    echo "errors:\n";
    foreach ($errors as $error){
        echo "\t{$error}\n";
    }
}

echo "comments:\n";
foreach ($comments as $comment){
    echo "\t{$comment}\n";
}

foreach ($fields as $option => $field){
    if (sizeOf($field) < 1) continue;

    echo "{$option}:\n";
    foreach ($field as $value){
        echo "\t{$value}\n";
    }
}
