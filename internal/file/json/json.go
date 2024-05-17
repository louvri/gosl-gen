package json

import (
	"encoding/json"
	"io"
	"os"
)

type File struct {
	path string
}

func New(path string) File {
	return File{
		path: path,
	}
}

func (f *File) Object() (map[string]interface{}, error) {
	response := make(map[string]interface{})
	file, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(byteValue), &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (f *File) Array() ([]map[string]interface{}, error) {
	response := make([]map[string]interface{}, 0)
	file, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(byteValue), &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
