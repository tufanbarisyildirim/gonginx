# Gonginx
![reportcard](https://goreportcard.com/badge/github.com/tufanbarisyildirim/gonginx) [![Actions Status](https://github.com/tufanbarisyildirim/gonginx/workflows/Go/badge.svg)](https://github.com/tufanbarisyildirim/gonging/actions)


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
directive   : Keyword [parameters] Semicolon [block]
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
  Parser is the main package that analyzes and turns nginx structred files into objects. It basically has 3 libraries, `lexer` explodes it into `token`s and `parser` turns tokens into config objects which are in their own package, 
- ### [Config](/config)
  Config package is representation of any context, directive or their parameters in golang. So basically they are models and also AST
- ### [Dumper](/dumper)
  Dumper is the package that holds styling configuration only. 

#### TODO
- [ ]  associate comments with config objects to print them on config generation and make it configurable with `dumper.Style`
- [ ]  move any context wrapper into their own file (remove from parser)
- [ ]  wire config object properties to their sub object (Directives & Block)   
       e.g, S`etting UpstreamServer.Address` should update `Upstream.Directive.Parameters[0]` if that's ugly, find another way to bind data between config and AST
- [ ]  Parse included files recusively, keep relative path on load, save all in a related structure and make that optional in dumper.Style

## Limitations
There is no known limitations yet. PRs are more then welcome if you want to implement a specific directive / block, please read [Contributing](CONTRIBUTING.md) before your first PR.

## License
[MIT License](LICENSE)