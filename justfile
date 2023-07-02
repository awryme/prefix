binpath := "./.bin/prefix"
mainpath := "./cmd/prefix/"

build:
    go build -v -o {{binpath}} {{mainpath}}

run +args: build
    {{binpath}} {{args}}