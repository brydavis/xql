package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
)

func main() {
	DecodeJSON("../../data/odd.json")

}

func DecodeJSON(filename string) {
	base := path.Base(filename)
	base = base[:len(base)-len(path.Ext(base))]
	raw, _ := ioutil.ReadFile(filename)

	var v interface{}
	json.Unmarshal(raw, &v)

	var data string
	switch t := v.(type) {
	case []interface{}:
		for _, row := range t {
			m := row.(map[string]interface{})
			data += ProcessMap(base, m)
		}
	case map[string]interface{}:
		data += ProcessMap(base, t)
	}

	lower := strings.ToLower(base)
	b, _ := ioutil.ReadFile("templates/base.tmpl")

	t := template.New(base)
	t, _ = t.Parse(string(b))

	filename = fmt.Sprintf("output/%s.go", lower)
	f, _ := os.Create(filename)

	t.Execute(f, struct {
		LowerName string
		Struct    string
		Name      string
	}{
		lower,
		fmt.Sprintf("type %s", data),
		base,
	})

	exec.Command("go", "fmt", filename).Run()
	exec.Command("goimports", "-w", filename).Run()
}

func ProcessMap(name string, data map[string]interface{}) string {
	s := fmt.Sprintf("%s struct {\n", strings.Title(name))
	for key, val := range data {
		switch vval := val.(type) {
		case []interface{}:
			s += ParseSlice(key, vval)
		case map[string]interface{}:
			s += ProcessMap(key, vval)
		default:
			s += fmt.Sprintf("\t%s\t%T\n", strings.Title(key), vval)
		}
	}
	return s + "\n}\n"
}

func ParseSlice(key string, vals []interface{}) string {
	s := fmt.Sprintf("\t%s\t[]", strings.Title(key))
	switch vals[0].(type) {
	case []interface{}:
		// s += fmt.Sprintf("\t%s\t%T %v\n", key, vals)
	default:
		s += fmt.Sprintf("%T\n", vals[0])
	}
	return s + "\n"
}
