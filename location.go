package gonginx

// Location represents a location in nginx config
type Location struct {
	*Directive
	Modifier string
	Match    string
	Parent   IBlock
}

func (l *Location) SetParent(parent IBlock) {
	l.Parent = parent
}

func (l *Location) GetParent() IBlock {
	return l.Parent
}
