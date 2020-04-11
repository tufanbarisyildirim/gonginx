package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

//Config  represents a whole config file.
type Config struct {
	Block
	FilePath string
}

//ToString return config as string
func (c *Config) ToString() string {
	return c.Block.ToString()
}

//ToByteArray return config as byte array
func (c *Config) ToByteArray() []byte {
	return c.Block.ToByteArray()
}

//SaveToFile save config to a file
//TODO: add custom file / folder path support
func (c *Config) SaveToFile() error {
	//wrilte file
	dirPath := filepath.Dir(c.FilePath)
	if _, err := os.Stat(dirPath); err != nil {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	//write main file
	err := ioutil.WriteFile(c.FilePath, c.ToByteArray(), 0644)
	if err != nil {
		return err
	}

	//write sub files (inlude /file/path)
	for _, statement := range c.Block.Statements {
		if fs, ok := statement.(FileStatement); ok {
			err := fs.SaveToFile()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
