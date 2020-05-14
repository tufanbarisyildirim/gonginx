package gonginx

//Include include structure
type Include struct {
	*Directive
	IncludePath string
	Configs     []*Config
}

//TODO(tufan): move that part into dumper package
//SaveToFile saves include to its own file
//func (i *Include) SaveToFile(style *dumper.Style) error {
//	if len(i.Configs) == 0 {
//		return fmt.Errorf("included empty file %s", i.IncludePath)
//	}
//	for _, c := range i.Configs {
//		err := c.SaveToFile(style)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

//GetName get directive name
func (i *Include) GetName() string {
	return i.Directive.Name
}
