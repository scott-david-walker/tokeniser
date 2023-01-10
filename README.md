# Tokeniser

Simple GitHub Action written in go to replace strings in files

## Inputs

- `files` - Glob expression of files to search through. E.G. `**/*.yaml`
- `prefix` - Prefix to use when matching tokens, defaults to `#{`
- `suffix` - Suffix to use when matching tokens, defaults to `}#`
- `fail-if-no-provided-replacement` - Will fail if a token is found but there was no variable provided to replace it. Defaults to `true` 

## Example

If you wanted to replace `#{version}#` in all of your yaml files

```yml
- uses: scott-david-walker/tokeniser@v1.0.5
  with:
    files: '**/*.yaml'
  env:
    version: "1.0.5"  
```

If you want to use a different format, you can use a different prefix and suffix.
Same example as before, but this time with `{{ }}`

```yml
- uses: scott-david-walker/tokeniser@v1.0.5
  with:
    prefix: '{{'
    suffix: '}}'
    files: '**/*.yaml'
  env:
    version: "1.0.5"
```
