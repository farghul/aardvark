# Aardvark

Aardvark is a WordPress website (blog) archiving tool. It Automates the process of archiving, restoring, and deleting blog sites.

![Aardvark](aardvark.svg)

## Prerequisites

- Googles' [Go language](https://go.dev) installed to enable building executables from source code.

- An `env.json` file containing enviromental data, for example:

``` json
{
    "address": "example.com",
    "assets": "Path to your website assets",
    "folder": "Root folder of your WordPress installation",
    "server": "Server name hosting WordPress",
    "wordpress": "Full path to the web/wp folder",
    "user": "User authorized to make changes on the server"
}
```

## Build

From the root folder containing *main.go*, use the command that matches your environment:

### Windows & Mac:

``` zsh
go build -o [name] .
```

### Linux:

``` zsh
GOOS=linux GOARCH=amd64 go build -o [name] .
```

## Run

``` zsh
[program] [operation flag] [website address]
```

## Operations

Current operational flags are:

- -a (Archive)

- -r (Restore)

- -d (Delete)

### Example Deployment

``` zsh
aardvark -a www.example.com
```

**NOTE**: Do not include the `http://` or `https://` prefixes.

## License

Code is distributed under [The Unlicense](https://github.com/farghul/aardvark/blob/main/LICENSE.md) and is part of the Public Domain.
