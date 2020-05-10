package config

import "github.com/tufanbarisyildirim/gonginx/dumper"

//IDirective represents any directive
type IDirective interface {
	ToString(style *dumper.Style) string
	GetName() string //the directive name.
	GetParameters() []string
	GetBlock() *Block
}

//FileDirective a statement that saves its own file
type FileDirective interface {
	SaveToFile(style *dumper.Style) error
}

//IncludeDirective represents include statement in nginx
type IncludeDirective interface {
	FileDirective
}
