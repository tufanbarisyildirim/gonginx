package config

//Location represents a location in nginx config
type Location struct {
	*Directive
	Modifier string
	Match    string
}

//ToString serialize location
func (l *Location) ToString() string {
	return l.Directive.ToString()
}
