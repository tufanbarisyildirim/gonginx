package config

import "errors"

// Location represents a location block in an nginx configuration.
type Location struct {
	*Directive
	Modifier string
	Match    string
	Parent   IDirective
	Line     int
}

// SetLine sets the line number.
func (l *Location) SetLine(line int) {
	l.Line = line
}

// GetLine returns the line number.
func (l *Location) GetLine() int {
	return l.Line
}

// SetParent sets the parent directive.
func (l *Location) SetParent(parent IDirective) {
	l.Parent = parent
}

// GetParent returns the parent directive.
func (l *Location) GetParent() IDirective {
	return l.Parent
}

// NewLocation initializes a Location from a directive.
func NewLocation(directive IDirective) (*Location, error) {
	dir, ok := directive.(*Directive)
	if !ok {
		return nil, errors.New("location directive type error")
	}

	if len(dir.Parameters) == 0 {
		return nil, errors.New("no enough parameter for location")
	}
	location := &Location{
		Modifier:  "",
		Match:     "",
		Directive: dir,
	}
	if len(dir.Parameters) == 1 {
		location.Match = dir.Parameters[0].GetValue()
		return location, nil
	} else if len(dir.Parameters) == 2 {
		location.Modifier = dir.Parameters[0].GetValue()
		location.Match = dir.Parameters[1].GetValue()
		return location, nil
	}
	return nil, errors.New("too many arguments for location directive")
}

// FindDirectives finds directives by name.
func (l *Location) FindDirectives(directiveName string) []IDirective {
	block := l.GetBlock()
	if block == nil {
		return []IDirective{}
	}
	return block.FindDirectives(directiveName)
}

// GetDirectives returns all directives in the location block.
func (l *Location) GetDirectives() []IDirective {
	block := l.GetBlock()
	if block == nil {
		return []IDirective{}
	}
	return block.GetDirectives()
}
