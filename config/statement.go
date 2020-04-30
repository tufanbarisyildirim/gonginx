package config

import "github.com/tufanbarisyildirim/gonginx/dumper"

//Statement represents any statement
type Statement interface {
	ToString(style *dumper.Style) string
	GetName() string //the directive name.
}

//FileStatement a statement that saves its own file
type FileStatement interface {
	Statement
	SaveToFile(style *dumper.Style) error
}

//IncludeStatement represents include statement in nginx
type IncludeStatement interface {
	FileStatement
}
