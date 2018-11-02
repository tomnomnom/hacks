# tojson

Turn lines of stdin into JSON.

## Examples

An example bit of output:

```
▶ ls -l | tail -n1
-rw-rw-r-- 1 tom tom 1365 Nov  2 14:20 main.go
```

By default each line becomes an item in an array:

```
▶ ls -l | tail -n1 | tojson
[
    "-rw-rw-r-- 1 tom tom 1365 Nov  2 14:20 main.go"
]
```

You can use `--format=2d-array` to split the fields on each line:

```
▶ ls -l | tail -n1 | tojson --format=2d-array
[
    [
        "-rw-rw-r--",
        "1",
        "tom",
        "tom",
        "1365",
        "Nov",
        "2",
        "14:20",
        "main.go"
    ]
]
```

Or you can use `--format=map` and specify the name of each field:

```
▶ ls -l | tail -n1 | tojson --format=map mode links users group size month day time name
[
    {
        "day": "2",
        "group": "tom",
        "links": "1",
        "mode": "-rw-rw-r--",
        "month": "Nov",
        "name": "main.go",
        "size": "1365",
        "time": "14:20",
        "users": "tom"
    }
]
```

## Install

```
▶ go get -u github.com/tomnomnom/hacks/tojson
```
