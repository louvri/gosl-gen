package config

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	//go:embed template
	templates embed.FS
)

type Option int

const (
	Read Option = iota
	Write
)

type File struct {
	path   string
	option Option
	wb     *os.File
	rb     fs.File
}

func New(path string, option Option) File {
	return File{
		path:   path,
		option: option,
	}
}

func (f *File) Read(process func(data string) error) error {
	if f.rb == nil {
		file, err := templates.Open(f.path)
		if err != nil {
			return err
		}
		f.rb = file
	}
	scanner := bufio.NewScanner(f.rb)
	for scanner.Scan() {
		err := process(fmt.Sprintf("%s\n", scanner.Text()))
		if err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (f *File) Write(data string) error {
	if f.wb == nil {
		os.Remove(f.path)
		dir := filepath.Dir(f.path)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0777)
			if err != nil {
				return err
			}
		}
		file, err := os.Create(f.path)
		if err != nil {
			return err
		}
		f.wb = file
	}
	if _, err := f.wb.WriteString(data); err != nil {
		return err
	}
	return nil
}

func (f *File) Close() error {
	if f.wb != nil {
		return f.wb.Close()
	}
	if f.rb != nil {
		return f.rb.Close()
	}
	return nil
}
