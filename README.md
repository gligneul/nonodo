# NoNodo

![Build](https://github.com/gligneul/nonodo/actions/workflows/build.yaml/badge.svg)

NoNodo is a development node for the [Cartesi Rollups](https://docs.cartesi.io/cartesi-rollups/) that was designed to work with applications running in the host machine instead of the Cartesi machine.
So, the application developer doesn't need to be concerned with compiling their application to RISC-V.
The application back-end should run in the developer's machine and call the Rollup HTTP API to process advance and inspect inputs.

NoNodo is a valuable development workflow help, but there are some caveats the developer must be aware of:

- The application will eventually need to be compiled to RISC-V or use a RISC-V runtime in case of interpreted languages;
- In this mode, the application will not be running inside the sandbox of the Cartesi machine and will not block operations that won't be allowed when running inside a Cartesi machine, like accessing remote resources;
- This mode only works for applications that use the Cartesi Rollups HTTP API and doesn't work with applications using the low-level API;
- Performance inside a Cartesi machine will be much lower than running on the host.

## Installation

### Pre-requisites

NoNodo uses the Anvil as the underlying Ethereum node.
To install Anvil, read the instructions on the [Foundry book](https://book.getfoundry.sh/getting-started/installation).

### Binaries

Download the binaries from the [relase page](https://github.com/gligneul/nonodo/releases).

### Installing from source

Installing NoNodo using the [`go install`](https://go.dev/ref/mod#go-install) command is possible.
To use this command, install Go following the instructions on the [Go website](https://go.dev/doc/install).
Then, run the command below to install NoNodo.

```sh
go install github.com/gligneul/nonodo@latest
```

The command above installs NoNodo into the `bin` directory inside the directory defined by the `GOPATH` environment variable.
If you don't set the `GOPATH` variable, the default install location is `$HOME/go/bin`.
So, to call the `nonodo` command directly, you should add it to the `PATH` variable.
The command below exemplifies that.

```sh
export PATH="$HOME/go/bin:$PATH"
```

## Usage

To start NoNodo with the default configuration, run the command below.

```sh
nonodo
```

With the default configuration, NoNodo starts an Anvil node with the Cartesi Rollups contracts deployed.
NoNodo uses the same deployment used by [sunodo](https://docs.sunodo.io/), so the contract addresses are the same.
NoNodo offers some flags to configure Anvil, starting with `--anvil-*`.

NoNodo exposes the Cartesi Rollups GraphQL (`/graphql`) and Inspect (`/inspect`) APIs for the application front-end, and the Rollup (`/rollup`) API for the application back-end.
NoNodo uses the HTTP address and port set by the `--http-address` and `--http-port` flags.
By default, NoNodo binds to the `http://127.0.0.1:8080/` address.

### Running the Application

NoNodo can run the application back-end as a sub-process.
If this process exits, NoNodo will also stop its execution.
This option is helpful to keep the whole development context in a single terminal.
To use this option, pass the command to run the application after `--`.

```sh
nonodo -- ./my-app
```

#### Built-in Echo Application

NoNodo has a built-in echo application that generates a voucher, a notice, and a report for each advance input.
The echo also generates a report for each inspect input.
This option is useful when testing the application front-end without a working back-end.
To start NoNodo with the built-in echo application, use the `--built-in-echo` flag.

```sh
nonodo --built-in-echo
```

### Sending inputs

To send an input to the Cartesi application, you may use cast, a command-line tool from the foundry
package. For instance, the invocation below sends an input with contents `0xdeadbeef` to the running
application.

```sh
INPUT=0xdeadbeef; \
INPUT_BOX_ADDRESS=0x59b22D57D4f067708AB0c00552767405926dc768; \
APPLICATION_ADDRESS=0x70ac08179605AF2D9e75782b8DEcDD3c22aA4D0C; \
cast send \
    --mnemonic "test test test test test test test test test test test junk" \
    --rpc-url "http://localhost:8545" \
    $INPUT_BOX_ADDRESS "addInput(address,bytes)(bytes32)" $APPLICATION_ADDRESS $INPUT
```

### GraphQL API

NoNodo exposes the GraphQL reader API in the endpoint `http://127.0.0.1:8080/graphql`.
You may access this address to use the GraphQL interactive playground in your web browser.
You can also make POST requests directly to the GraphQL API.
For instance, the command below gets the number of inputs.

```sh
QUERY='query { inputs { totalCount } }'; \
curl \
    -X POST \
    -H 'Content-Type: application/json' \
    -d "{\"query\": \"$QUERY\"}" \
    http://127.0.0.1:8080/graphql
```

### Inspect API

Nonodo exposes the Inspect API in the endpoint `http://127.0.0.1:8080/inspect`.
Like the Rollups Node, you may send inspect requests with GET or POST methods.
For instance, the command below sends an inspect input with payload `hi` to NoNodo.

```
curl -X POST -d "hi" http://127.0.0.1:8080/inspect/
```
