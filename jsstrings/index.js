let fs = require("fs");
let acorn = require("acorn");

let fn = process.argv[2];

if (fn == ""){
   console.log("usage: jsstrings <file>");
   process.exit() ;
}

fs.readFile(fn, "utf8", function(err, data) {
    // TODO: Trim the strings
    for (let token of acorn.tokenizer(data,{
        onComment: function(block, text, start, end){
            console.log(text);
        }
    })) {
        if (token.type == acorn.tokTypes.string){
            console.log(token.value);
        }
    }
});
