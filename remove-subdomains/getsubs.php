#!/usr/bin/env php
<?php
$suffixes = @file('suffixes.txt');
if (!$suffixes) die('failed to load suffixes.txt');

$suffixes = array_map(function($line){
    return trim($line, " \t\n\r\0\x0B.");
}, $suffixes);

$suffixes = array_filter($suffixes, function($line){
    return $line != "";
});

$map = [];

foreach ($suffixes as $suffix){
   $map[$suffix] = true; 
}

while (($line = fgets(STDIN)) != false){
    $candidate = trim($line);
    $parts = explode(".", $candidate);

    $tail = [];
    array_unshift($tail, array_pop($parts));

    if (isset($map[$tail[0]])){
        // We've got a tld, try for an sld
        array_unshift($tail, array_pop($parts));

        if (!isset($map[implode(".", $tail)])){
            // No match, put it back
            array_push($parts, array_shift($tail));
        }
    }

    array_pop($parts);
    //echo array_pop($parts).".".implode(".", $tail).PHP_EOL;
    echo implode(".", $parts).PHP_EOL;
}
