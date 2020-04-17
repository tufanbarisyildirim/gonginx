# Gonginx
![reportcard](https://goreportcard.com/badge/github.com/tufanbarisyildirim/gonginx) [![Actions Status](https://github.com/tufanbarisyildirim/gonginx/workflows/Go/badge.svg)](https://github.com/tufanbarisyildirim/gonging/actions)


Gonginx is an Nginx configuration parser helps you to parse, edit, regenerate your nginx config files in your go applications. It makes managing your banalcer configurations easier. We use this library in a tool that discovers microservices and updates our nginx balancer config. We will make it opensource soon.

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
  Parser is the main package that analyzes and turns nginx structred files into objects. It basically has 2 libraries, `lexer` explodes it into `token`s and `parser` turns tokens into config objects which are in their own package, 
- ### [Config](/config)
  Config package is representation of any context, directive or their parameters in golang. So basically they are models and also AST
- ### [Dumper (in progress)](/dumper)
  Dumper is the package that can print any model with some styling options. 

### Supporting Blocks/Directives - TODO
Generated a to-do/feature list from a full nginx config examle to track how is going.
Most common directives will be checked when they implemented. But blocks will be checked when we fully support their sub directives.

#### General TODO
- [ ]  associate comments with config objects to print them on config generation
- [ ]  move any context wrapper into their own file (remove from parser)
- [ ]  wire config object properties to their sub object (Directives & Block)   
       e.g, S`etting UpstreamServer.Address` should update `Upstream.Directive.Parameters[0]` if that's ugly, find another way to bind data between config and AST

#### TODO for directives, parsing


## Limitations
There is no limitation yet, because its the limt itself :) I haven't implemented all features yet. PRs are more then welcome if you want to implement a specific directive / block

# [Contributing](CONTRIBUTING.md)
Any PR is welcome!

## License
[MIT License](LICENSE)