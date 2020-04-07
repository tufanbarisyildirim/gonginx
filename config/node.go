package config

//Node represents any meaningful node of nginx config
type Node interface {
	TokenLiteral() string
	String() string
}
