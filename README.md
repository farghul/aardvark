# Aardvark

Aardvark is a WordPress website (blog) archiving application. It Automates the process of archiving blog sites.

![Aardvark](aardvark.svg)

## Prerequisites

- [Dart](https://dart.dev) programming language installed to enable building executables from source code.

- An `json` file named for your enviroment containing server specific data, for example:

``` json
{
    "address": "Full website URL",
    "assets": "Path to your WordPress website assets",
    "ephemeral": "Path to your temp folder",
    "home": "Home folder of the authorized user",
    "install": "Full path to the web/wp folder",
    "lists": "Full path to the folder holding the results of the getSites() function",
    "server": "Server name hosting WordPress",
    "sites": "Title of the file created by getSites() function",
    "folder": "Root folder of your WordPress installation",
    "user": "User authorized to make changes on the server"
}
```

## Build

``` zsh
dart compile exe bin/main.dart
```

## Run

``` zsh
[program] [environment flag] [environment] [website address]
```

### Example

``` zsh
aardvark -e test https://www.example.com
```

## License

Code is distributed under [The Unlicense](https://github.com/farghul/aardvark/blob/main/LICENSE.md) and is part of the Public Domain.
