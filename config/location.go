package config

import "github.com/tufanbarisyildirim/gonginx/dumper"

//Location represents a location in nginx config
type Location struct {
	*Directive
	Modifier string
	Match    string
}

//ToString serialize location
func (l *Location) ToString(style *dumper.Style) string {
	return l.Directive.ToString(style)
}
