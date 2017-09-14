<!doctype html>
<html>
    <head>
        <title>unimap</title>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
        <style>
            body {
                font-family: sans-serif;
            }
            #charRange {
                width: 100%;
            }
            .cell {
                font-size: 24px;
                background: #EFEFEF;
                padding: 1px;
                padding-top: 5px;
                margin: 2px;
                width: 1.5em;
                height: 1.5em;
                text-align: center;
                vertical-align: middle;
                display: inline-block;
            }
            .container {
                width: 100%;
                max-width: 1000px;
                margin: 0 auto;
                text-align: center;
            }
            
        </style>
    </head>
    <?php
        const CHAR_COUNT = 256;
    ?>
    <body>

        <div class=container>
            <div>
                <input type=range id=charRange min=0 max=65535 step=1 value=0>
            </div>

            <?php
                for ($i = 0; $i < CHAR_COUNT; $i++){
                    echo "<span class=cell id=cell-{$i}></span>";
                }
            ?>
        </div>
    </body>
    <script>
        (function(){
            var charCount = <?= CHAR_COUNT; ?>;
            var cells = [];
            for (var i = 0; i < charCount; i++){
                cells[i] = document.getElementById('cell-'+i);
            }
            var range = document.getElementById('charRange');

            var prev = null;

            function update(){
                lower = parseInt(range.value, 10);
                if (lower != prev){
                    for (var i = lower; i < lower + charCount; i++){
                        var c = String.fromCodePoint(i);
                        cells[i-lower].innerHTML = c;
                    }
                    prev = lower;
                }
            
                window.requestAnimationFrame(update);
            }

            window.requestAnimationFrame(update);


        })();
    </script>
</html>
