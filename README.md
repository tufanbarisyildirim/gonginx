<p align="center"><img src="./gopher.png" width="360"></p>
<p align="center">
<a href="https://goreportcard.com/badge/github.com/tufanbarisyildirim/gonginx"><img src="https://goreportcard.com/badge/github.com/tufanbarisyildirim/gonginx" alt="Report Card" /></a>
<a href="https://github.com/tufanbarisyildirim/gonging/actions"><img src="https://github.com/tufanbarisyildirim/gonginx/workflows/Go/badge.svg" alt="Actions Status" /></a>
</p>

# Gonginx
Gonginx is an Nginx configuration parser helps you to parse, edit, regenerate your nginx config files in your go applications. It makes managing your balancer configurations easier. 

## Basic grammar of an nginx config file
```yacc

%token Keyword Variable BlockStart BlockEnd Semicolon Regex

%%

config      :  /* empty */ 
            | config directives
            ;
block       : BlockStart directives BlockEnd
            ;
directives  : directives directive
            ;
directive   : Keyword [parameters] (semicolon|block)
            ;
parameters  : parameters keyword
            ;
keyword     : Keyword 
            | Variable 
            | Regex
            ;
```

## Core Components
- ### [Parser](/parser) 
  Parser is the main package gonginx
- ### [Config](/config)
  Config package gonginx
- ### [Dumper](dumper.go)
  Dumper is the package gonginx

#### TODO
- [ ]  associate comments with config objects to print them on config generation and make it configurable with `dumper.Style`
- [ ]  move any context wrapper into their own file (remove from parser)
- [ ]  Parse included files recusively, keep relative path on load, save all in a related structure and make that optional in dumper.Style
- [ ]  Implement specific searches, like finding servers by server_name (domain) or any upstream by target etc.

## Limitations
There is no known limitations yet. PRs are more then welcome if you want to implement a specific directive / block, please read [Contributing](CONTRIBUTING.md) before your first PR.

## License
[MIT License](LICENSE)