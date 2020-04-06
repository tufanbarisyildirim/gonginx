package ast

type Node interface {
	TokenLiteral() string
	String() string
}

type Block interface {
	Node
}


type Directive interface {
	Node
	Block
}

type Config struct {
	Directives []Directive
}

//  nginx blocks

//	events {
//  	worker_connections  4096;  ## Default: 1024
//	}
//
//	http {
//  	server { # php/fastcgi
//    	location {
//    	}
//  }
//
//  upstream {
//  }
//}
