package config

//Location represents a location in nginx config
type Location struct {
	*Block
	Modifier string
	Match    string
}

//ToString serialize location
func (l *Location) ToString() string {
	return l.Block.ToString()
}
