# Coding Agent Guidelines for Gonginx

## Explanation

Gonginx is a Go library for parsing, editing, and regenerating nginx
configuration files. The main packages are:

- **parser** – turns nginx config files into structured objects.
- **config** – models directives and blocks.
- **dumper** – renders configuration objects back into text.

Examples demonstrating typical usage can be found in the `examples/`
folder and in `GUIDE.md`.

## Examples

Parse a configuration file and print the listen ports of all servers:

```go
p, err := parser.NewParser("nginx.conf")
if err != nil {
    log.Fatal(err)
}
conf, err := p.Parse()
if err != nil {
    log.Fatal(err)
}
servers := conf.FindDirectives("server")
for _, srv := range servers {
    for _, listen := range srv.GetBlock().FindDirectives("listen") {
        fmt.Println(listen.GetParameters()[0].GetValue())
    }
}
```

Add a server to an upstream block and dump the result:

```go
p := parser.NewStringParser(`http{ upstream backend{} }`)
conf, _ := p.Parse()
up := conf.FindUpstreams()[0]
up.AddServer(&config.UpstreamServer{Address: "127.0.0.1:443"})
fmt.Println(dumper.DumpConfig(conf, dumper.IndentedStyle))
```

More detailed examples, including writing configs to disk, are in
`GUIDE.md`.

## Contribution Guide

1. Follow the [Code of Conduct](CODE_OF_CONDUCT.md).
2. Before sending a pull request, run `make test` from the repository
   root to format the code and execute all tests.
3. Commit messages should be clear and contain logical units, e.g.

   ```
   fix parser include handling

   * handle recursive includes
   * add regression test
   ```

4. If you add new functionality, include tests that fail without your
   change and pass with it.

## Using the Library

Install dependencies via Go modules and import packages directly from
`github.com/tufanbarisyildirim/gonginx`:

```go
import (
    "github.com/tufanbarisyildirim/gonginx/config"
    "github.com/tufanbarisyildirim/gonginx/dumper"
    "github.com/tufanbarisyildirim/gonginx/parser"
)
```

Create a parser with `parser.NewParser` (for files) or
`parser.NewStringParser` (for strings). After modifying the returned
`config.Config` object, convert it back to text with
`dumper.DumpConfig` or persist it using `dumper.WriteConfig`.

