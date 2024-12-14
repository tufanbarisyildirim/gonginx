package config

import "errors"

// Location represents a location in nginx config
type Location struct {
	*Directive
	Modifier string
	Match    string
	Parent   IBlock
}

// SetParent change the parent block
func (l *Location) SetParent(parent IBlock) {
	l.Parent = parent
}

// GetParent the parent block
func (l *Location) GetParent() IBlock {
	return l.Parent
}

// NewLocation initialize a location block from a directive
func NewLocation(directive IDirective) (*Location, error) {
	dir, ok := directive.(*Directive)
	if !ok {
		return nil, errors.New("no ")
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
