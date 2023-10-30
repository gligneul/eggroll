# EggRoll 🐣🛼

![Tests](https://github.com/gligneul/eggroll/actions/workflows/test.yml/badge.svg)
[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/gligneul/eggroll)
[![Documentation](https://img.shields.io/static/v1?label=guide&message=documentation&color=blue)](https://gligneul.github.io/eggroll)

A high-level framework for Cartesi Rollups in Go.

## Requirements

EggRoll is built on top of the [Cartesi Rollups](https://docs.cartesi.io/cartesi-rollups/) infrastructure version 1.0.
To use EggRoll, you also need [sunodo](https://github.com/sunodo/sunodo/) version 0.9.

## Documentation

Check the [documentation site](https://gligneul.github.io/eggroll) for more info.

## Development

The commands below are for EggRoll developers.
EggRoll users should check the documentation site above.

### Building

```
go build
```

### Unit testing

```
go test
```

### Integration testing

```
EGGTEST_RUN_INTEGRATION=true EGGTEST_VERBOSE=true go test -p 1 -v ./examples/...
```

### Running documentation server

```
doctave serve
```
