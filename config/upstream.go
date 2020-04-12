package config

//Upstream represents `upstream{}` block
type Upstream struct {
	*Directive
	Name string
}

//ToString convert it to a string
func (us *Upstream) ToString() string {
	return us.Directive.ToString()
}
