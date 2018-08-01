<?php

$lines = file('map');

$trans = [];
foreach ($lines as $line){
    $line = trim($line);
    $p = explode("|", $line);

    $char = $p[0];
    $subs = $p[1];

    if (!isset($trans[$char])){
        $trans[$char] = [];
    }
    $trans[$char][] = $subs;
}

foreach ($trans as $char => $subs){
    $subs = array_map(function($sub){
        $sub = str_pad($sub, 8, "0", STR_PAD_LEFT);
        return "\"\U{$sub}\"";
    }, $subs);
    $char = str_pad($char, 8, "0", STR_PAD_LEFT);
    echo "\"\U{$char}\": {".implode(', ', $subs)."},\n";
}

function codepointToEscape($cp){

}
