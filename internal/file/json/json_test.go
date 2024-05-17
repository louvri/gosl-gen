package json

import (
	"encoding/json"
	"os"
	"testing"
)

const path1 string = "./build/test_object.json"
const path2 string = "./build/test_array.json"

func init() {
	create := func(path, data string) {
		err := os.MkdirAll("./build", 0777)
		if err != nil {
			panic(err)
		}
		os.Remove(path)
		file, err := os.Create(path)
		if err != nil {
			panic(err.Error())
		}
		defer file.Close()
		if _, err := file.WriteString(data); err != nil {
			panic(err)
		}
	}
	create(path1, `{
		"hello":"world",
		"color":"beautiful"
	}`)
	create(path2, `[{
		"hello":"world",
		"color":"beautiful"
	}]`)
}

func TestObject(t *testing.T) {
	jsn := New(path1)
	result, err := jsn.Object()
	if err != nil {
		t.Fatal(err.Error())
	}
	if result["hello"] != "world" {
		t.Fail()
	}
	out, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(string(out))
}

func TestArray(t *testing.T) {
	jsn := New(path2)
	result, err := jsn.Array()
	if err != nil {
		t.Fatal(err.Error())
	}
	if result[0]["hello"] != "world" {
		t.Fail()
	}
	out, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(string(out))
}
