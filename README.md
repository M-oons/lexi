# lexi

Small CLI for word utilities.

## Build

```bash
make build
```

## Usage
```
Usage:
  lexi [flags]
  lexi [command]

Available Commands:
  help        Help about any command
  random      Generate random words from a list

Flags:
  -h, --help   help for lexi
```

## Commands

### random

Generate random words from a list.

```
Usage:
  lexi random [flags]

Flags:
  -n, --count int          maximum number of words to generate (default 1)
  -f, --format string      output format for each line; use {}/{l}/{u} placeholders for random words in default case, lower case or upper case respectively (default "{}")
  -h, --help               help for random
      --max int            maximum length of words to generate
      --min int            minimum length of words to generate
  -o, --order string       order of output words
  -r, --regex string       regular expression to filter words
      --seed int           seed for random number generator (default -1)
      --separator string   separator to use between words (defaults to newline for files, comma for inline words)
  -w, --words string       words to use for randomization
```
