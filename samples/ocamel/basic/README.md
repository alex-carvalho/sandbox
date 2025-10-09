# Basic OCaml Project with Dune

## Building

```shell
dune build
```

## Running tests

```shell
dune test
```

## Running individual modules

```shell
# Build and run the library
dune exec -- ocaml -I _build/lib lib/variables.ml
dune exec -- ocaml -I _build/lib lib/functions.ml
```