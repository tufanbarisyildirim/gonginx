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
	for _, statement := range c.Block.Statements {
		if fs, ok := statement.(FileStatement); ok {
			err := fs.SaveToFile(style)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
