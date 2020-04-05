package ast

type Node interface {
	TokenLiteral() string
	String() string
}

type Command interface {
	Node
	commandNode()
}

type Block interface {
	Node
	GetBlocks() []Block
	GetCommands() []Command
	blockNode()
}

type Config struct {
	Block
	Blocks   []Block
	Commands []Command
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
