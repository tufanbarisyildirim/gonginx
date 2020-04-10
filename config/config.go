package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
)

//Config  represents a whole config file.
type Config struct {
	FilePath   string
	Statements []Statement
}

//ToString return config as string
func (c *Config) ToString() string {
	return string(c.ToByteArray())
}

//ToByteArray return config as byte array
func (c *Config) ToByteArray() []byte {
	var buf bytes.Buffer

	for _, statement := range c.Statements {
		buf.WriteString(statement.ToString())
		buf.WriteString("\n")
	}

	return buf.Bytes()
}

//TokenLiteral returns the first token of config
func (c *Config) TokenLiteral() string {
	return ""
}

//SaveToFile save config to a file
func (c *Config) SaveToFile() error {
	//wrilte file
	dirPath := filepath.Dir(c.FilePath)
	if _, err := os.Stat(dirPath); err != nil {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(c.FilePath, c.ToByteArray(), 0644)
}
