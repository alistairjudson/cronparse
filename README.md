cronparse
=========


cronparse is a tool that can parse and expand cron expressions into the times
that said expression will run at.

e.g. 
```console
$ cronparse */15 0 1,15 * 1-5 /usr/bin/find
minute         0 15 30 45
hour           0
day of month   1 15
month          1 2 3 4 5 6 7 8 9 10 11 12
day of week    1 2 3 4 5
command        /usr/bin/find
```
## Contents
<!-- vim-markdown-toc GFM -->

- [Building](#building)
  - [Tests](#tests)
- [Architecture](#architecture)
  - [Lexical Analysis](#lexical-analysis)
  - [Parsing](#parsing)

<!-- vim-markdown-toc -->

### Building
cronparse is built using [Go][go]. To build cronparse you require the Go tool,
you can find how to do that for your specific system [here][installing-go].

Once you have installed Go, you can get up and running by cloning this
repository to a location of your choice and running:

`go run cmd/cronparse/main.go */15 0 1,15 * 1-5 /usr/bin/find`

Alternatively if you have your `$GOPATH/bin` added to your `$PATH` you can run:

`go get -u github.com/alistairjudson/cronparse/cmd/cronparse`

#### Tests
To run the tests you can simply run:

`go test ./...` 

### Linting
This project uses [golangci-lint][golangci-lint] in order to lint the project
it is configured by the `.golangci.yml` file.

### Architecture
cronparse has a couple of stages when it comes to parsing interpreting the
expression, these stages are:

#### Lexical Analysis
cronparse uses a [finite state machine][fsm] to tokenise the individual
components of the input into tokens, the way that this is implemented is heavily
inspired by the approach taken by Rob Pike in this [talk][rob-pike-talk]. If you
want to find out more about how this, this [blog post][blog-post] is pretty
cool.

In the tokenisation stage, the values themselves are not validated, only that
the syntax of the expression matches that of a component of a cron expression.

#### Parsing
The parsing of the cron expression is pretty lazy, it just uses some simple
pattern matching in order to validate the expressions, and expand the values
that they represent.

[go]: https://golang.org/
[installing-go]: https://golang.org/doc/install
[fsm]: https://en.wikipedia.org/wiki/Finite-state_machine
[rob-pike-talk]: https://www.youtube.com/watch?v=HxaD_trXwRE
[blog-post]: https://hackernoon.com/lexical-analysis-861b8bfe4cb0
[golangci-lint]: github.com/golangci/golangci-lint

