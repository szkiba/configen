# ConfiGen

ConfiGen is a template based configuration file generator tool with optional JSON Schema based validation support.

## Features:

- Generate arbitrary number of files
- Single executable binary
- Go template language
- Supports JSON, YAML, TOML data files
- JSON Schema based validation of generated files
- Supports local and remote schemas

## Usage

```
Usage:
  configen [options] [args]

Template based configuration generator.

You can specify multiple environments, input directories and values files.
Frequently used options has alternative positional argument syntax.

Options:
  -t, --template=directory    Input directory [arg: directory] (default: templates)
  -o, --output=directory      Output directory (default: dist)
  -s, --schema=directory      Schema directory (default: schemas)
  -f, --values=file           Data values file [arg: +file] (default: values.yaml)
      --set=name:value        Set value [arg: name=value]
      --loose                 Disable schema validation
      --dry-run               Skip writing output files
      --dump                  Dump intermediate files
  -q, --quiet                 Suppress console output
  -p, --package=file          Package descriptor template (default: package.json)
  -e, --env=environment       Staging environment name [arg: @environment]
      --dir=directory         Set working directory
  -V, --version               Show version information
  -w, --watch                 Watch and generate on filesystem changes
      --port=number           HTTP port for watch mode (default: random) [$PORT]

Help Options:
  -h, --help                  Show this help message
```