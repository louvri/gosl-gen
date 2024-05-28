package config

import (
	"os"
	"strings"
	"testing"
)

const output string = "./build/out.toml"

func init() {
	err := os.MkdirAll("./build", 0777)
	if err != nil {
		panic(err)
	}
}

func TestReadAndWriteConfig(t *testing.T) {
	toBeRead := New("template/config.toml", Read)
	toBeWritten := New(output, Write)
	defer func() {
		err := toBeRead.Close()
		if err != nil {
			t.Fatal(err.Error())
		}
		err = toBeWritten.Close()
		if err != nil {
			t.Fatal(err.Error())
		}
	}()
	err := toBeRead.Read(func(data string) error {
		new := strings.ReplaceAll(data, "$DB_CONNECTION_STRING", "YOU KNOW")
		return toBeWritten.Write(new)
	})
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestReadAndWriteHelper(t *testing.T) {
	toBeRead := New("template/helper.gotmpl", Read)
	toBeWritten := New(output, Write)
	defer func() {
		err := toBeRead.Close()
		if err != nil {
			t.Fatal(err.Error())
		}
		err = toBeWritten.Close()
		if err != nil {
			t.Fatal(err.Error())
		}
	}()
	err := toBeRead.Read(func(data string) error {
		new := strings.ReplaceAll(data, "$MODEL_PATH", "YOU")
		return toBeWritten.Write(new)
	})
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestReadAndWriteSearchBodyRequest(t *testing.T) {
	toBeRead := New("template/search_body.gotmpl", Read)
	toBeWritten := New(output, Write)
	defer func() {
		err := toBeRead.Close()
		if err != nil {
			t.Fatal(err.Error())
		}
		err = toBeWritten.Close()
		if err != nil {
			t.Fatal(err.Error())
		}
	}()
	err := toBeRead.Read(func(data string) error {
		new := strings.ReplaceAll(data, "$MODEL_PATH", "YOU")
		return toBeWritten.Write(new)
	})
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestReadAndWriteSearchRequest(t *testing.T) {
	toBeRead := New("template/search_request.gotmpl", Read)
	toBeWritten := New(output, Write)
	defer func() {
		err := toBeRead.Close()
		if err != nil {
			t.Fatal(err.Error())
		}
		err = toBeWritten.Close()
		if err != nil {
			t.Fatal(err.Error())
		}
	}()
	err := toBeRead.Read(func(data string) error {
		new := strings.ReplaceAll(data, "$MODEL_PATH", "YOU")
		return toBeWritten.Write(new)
	})
	if err != nil {
		t.Fatal(err.Error())
	}
}
