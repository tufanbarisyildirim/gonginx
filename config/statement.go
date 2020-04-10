package config

//Statement represents any statement
type Statement interface {
	ToString() string
}

//DirectiveStatement represents directives
type DirectiveStatement interface {
	Statement
	directiveStatement()
}

//ContextStatement represents blocks
type ContextStatement interface {
	Statement
	contextStatement()
}

//FileStatement a statement that saves its own file
type FileStatement interface {
	Statement
	SaveToFile() error
}

//IncludeStatement represents include statement in nginx
type IncludeStatement interface {
	FileStatement
	includeStatement()
}
