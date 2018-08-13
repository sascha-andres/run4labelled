package run4labelled

import (
	"github.com/pkg/errors"
	"os"
	p "path"
	"path/filepath"
	"strings"
)

type (
	// Configuration contains all data required to run the program
	Configuration struct {
		Label    string   `yaml:"label"`    // Label is the name of the file to look for
		Excludes []string `yaml:"excludes"` // Excludes is a list of directory names to ignore
		Run      string   `yaml:"run"`      // Run is executed in each labelled directory

		sendChannel chan Execute // sendChannel is used to tell executor to execute other process
	}

	// Execute is sent to executor with relevant information
	Execute struct {
		Directory string
		Command   string
	}
)

// SetChannel must be called to allow the logic to send commands to be executed to the executor
func (c *Configuration) SetChannel(toExecutor chan Execute) error {
	if nil == toExecutor {
		return errors.New("no channel provided")
	}
	c.sendChannel = toExecutor
	return nil
}

func (c *Configuration) Walk(baseDir string) error {
	if nil != c.sendChannel {
		defer close(c.sendChannel)
	}
	if s, err := os.Stat(baseDir); os.IsNotExist(err) || !s.IsDir() {
		return errors.New("base directory not fine")
	}
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			for _, val := range c.Excludes {
				if strings.HasSuffix(path, val) {
					return filepath.SkipDir
				}
			}
			if _, err := os.Stat(p.Join(path, c.Label)); err == nil {
				if nil != c.sendChannel {
					c.sendChannel <- Execute{
						Directory: path,
						Command:   c.Run,
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
