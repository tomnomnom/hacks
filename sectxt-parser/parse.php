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

$comments = [];
$fields = [
    FIELD_CONTACT         => [],
    FIELD_ENCRYPTION      => [],
    FIELD_DISCLOSURE      => [],
    FIELD_ACKNOWLEDGEMENT => [],
];

foreach ($lines as $line){
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
        die("invalid line: {$line}\n");
    }

    $option = strToLower($parts[0]);
    $value = $parts[1];

    if (!isset($fields[$option])){
        die("invalid option: {$parts[0]}\n");
    }

    $fields[$option][] = $value;
}

if (sizeOf($fields[FIELD_CONTACT]) < 1){
    die("does not contain at least one contact field\n");
}

echo "comments:\n";
foreach ($comments as $comment){
    echo "\t{$comment}\n";
}

foreach ($fields as $option => $field){
    echo "{$option}:\n";

    foreach ($field as $value){
        echo "\t{$value}\n";
    }
}
