<?php
$url = "https://hackerone.com/graphql";
$authtoken = $argv[1]?? die('needs auth token');

$query = <<<QUERY
query Settings { 
    query{ 
        id,
        teams(first: 50 after: "%s") {
            pageInfo {
                hasNextPage,
                hasPreviousPage
            },
            edges{
                cursor,
                node{
                    _id,
                    handle,
                    structured_scopes {
                        edges {
                            node {
                                id,
                                asset_type,
                                asset_identifier,
                                eligible_for_submission,
                                eligible_for_bounty,
                                max_severity,
                                archived_at,
                                instruction
                            }
                        }
                    }
                }
            }
        }
    }   
}
QUERY;

$gen = function($cursor = "") use($query){
	return json_encode([
		'query' => sprintf($query, $cursor),
        'variables' => (object) []
	]);
};


$cursor = "";
do {
    $params = [
        'http' => [
            'method' => 'POST',
            'header' => "Content-Type: application/json\r\n".
                        "Origin: https://hackerone.com\r\n".
                        "Referer: https://hackerone.com/programs\r\n".
                        "X-Auth-Token: {$authtoken}",
            'content' => $gen($cursor)
        ]
    ];
    $context = stream_context_create($params);
    $fp = fopen($url, 'rb', false, $context);
    $result = $fp ? stream_get_contents($fp) : null;
    $result = json_decode($result);
    if (!$result) die('response error');
    
    $hasNextPage = $result->data->query->teams->pageInfo->hasNextPage;

    foreach ($result->data->query->teams->edges as $edge){
        $cursor = $edge->cursor;
        foreach ($edge->node->structured_scopes->edges as $scope){
            $scope = $scope->node;
            if (!$scope->eligible_for_submission){
                continue;
            }
            if (strToLower($scope->asset_type) != "url"){
                continue;
            }

            echo $scope->asset_identifier.PHP_EOL;
        } 
    }

} while($hasNextPage);

