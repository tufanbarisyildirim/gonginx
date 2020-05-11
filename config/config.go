package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/tufanbarisyildirim/gonginx/dumper"
)

//Config  represents a whole config file.
type Config struct {
	*Block
	FilePath string
}

//ToString return config as string
func (c *Config) ToString(style *dumper.Style) string {
	return c.Block.ToString(style)
}

//ToByteArray return config as byte array
func (c *Config) ToByteArray(style *dumper.Style) []byte {
	return c.Block.ToByteArray(style)
}

//SaveToFile save config to a file
//TODO: add custom file / folder path support
func (c *Config) SaveToFile(style *dumper.Style) error {
	//wrilte file
	dirPath := filepath.Dir(c.FilePath)
	if _, err := os.Stat(dirPath); err != nil {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return err //TODO: do we reallt need to find a way to test dir creating error?
		}
	}

	//write main file
	err := ioutil.WriteFile(c.FilePath, c.ToByteArray(style), 0644)
	if err != nil {
		return err //TODO: do we need to find a way to test writing error?
	}

	//write sub files (inlude /file/path)
	for _, directive := range c.Block.Directives {
		if fs, ok := (interface{}(directive)).(FileDirective); ok {
			err := fs.SaveToFile(style)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//FindDirectives find directives from whole config block
func (c *Config) FindDirectives(directiveName string) []IDirective {
	return c.Block.FindDirectives(directiveName)
}

//Servers find directives from whole config block
func (c *Config) Servers() []*Server {
	var servers []*Server
	directives := c.Block.FindDirectives("upstream")
	for _, directive := range directives {
		s, _ := NewServer(directive)
		servers = append(servers, s)
	}
	return servers
}

//FindUpstreams find directives from whole config block
func (c *Config) FindUpstreams() []*Upstream {
	var upstreams []*Upstream
	directives := c.Block.FindDirectives("upstream")
	for _, directive := range directives {
		up, _ := NewUpstream(directive)
		upstreams = append(upstreams, up)
	}
	return upstreams
}
