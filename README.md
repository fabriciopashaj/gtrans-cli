# gtrans-cli
A simple CLI app in Go for interfacing with Google Translate from the terminal.

## Installation
```bash
go install github.com/fabriciopashaj/gtrans-cli
```

## Usage
When run the program will read commands from `Stdin`. Running the `help` command gives the following output.
```
> help
| Get internal environment variable
| > get envVar
| Set internal environment variables in the form key=value
| > set source=es target=en
| Translate text after command from Env["source"] to Env["target"]
| By default Env["source"] = "auto"
| > trt Krankenhaus
| Translate text in file in path {source} and output to file in path {dest}
| If {dest} is not an absolute path, it will be interpreted relative to the directory of {source}
| > trf {source} {dest}
| If {dest} is not provided translated text will be outputed to Stdout instead.
| > trf {source}
| Print this help message
| > help
```

`gtrans-cli` maintains an internal environment dictionary `map[string]string` that you can modify. Currently only the keys `source`, `target` and their values are of importance.

## Dependecies
- [](github.com/Conight/go-googletrans) v0.2.3
- [](https://github.com/chzyer/readline) v1.5.1
