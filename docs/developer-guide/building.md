# Building TFLint

Go 1.24 or higher is required to build TFLint from source code. Clone the source code and run the `make` command. Built binary will be placed in `dist` directory.

```console
$ git clone https://github.com/terraform-linters/tflint.git
$ cd tflint
$ make
mkdir -p dist
go build -v -o dist/tflint
```

## Run tests

If you change code, make sure that the tests you add and existing tests will be passed:

```console
$ make test
```

Some tests depending on Git submodules. Running `make test` will update these automatically, but if you run tests directly with `go test`, you need to update submodules manually:

```sh
git submodule init
git submodule update
```

## Run E2E tests

You can check the actual CLI behavior by running the E2E tests. Since the E2E tests uses the installed `tflint` command, it is necessary to add the path into `$PATH` environment so that the binary built by `go install` can be referenced.

```console
$ make e2e
```
