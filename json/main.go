package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

func main() {
	name := "MyJSON"
	raw := []byte(`[
		  {
		    "num": 6.13,
		    "strs": [
		      "a",
		      "b"
		    ],
		    "bln": true,
		    "flts": [
		      7,
		      9.2
		    ],
		    "nested": {
		      "n1": 74,
		      "n2": [
		        true,
		        false
		      ]
		    }
		  }
		]`)

	j := DecodeJSON(name, raw)
	// WriteTempFile(name, []byte(j))
	CreateTemplate(name, []byte(j))

}

func DecodeJSON(name string, raw []byte) string {
	St := "type "
	var i interface{}
	json.Unmarshal(raw, &i)
	switch data := i.(type) {
	case []interface{}:
		for _, row := range data {
			m := row.(map[string]interface{})
			St += ProcessMap(name, m)
		}
	case map[string]interface{}:
		St += ProcessMap(name, data)
	}
	return St
}

// func ProcessMap(data map[string]interface{}) map[string]interface{} {
// 	// m:= make(map[string]interface{})
// 	s := `
// type A struct {
// 	`

// 	for key, val := range data {
// 		switch vval := val.(type) {
// 		case int, int16, int32, int64:
// 			s += fmt.Sprintf("\t%s\tint %d\n", key, vval)
// 		case float32, float64:
// 			s += fmt.Sprintf("\t%s\tfloat64 %f\n", key, vval)
// 		case bool:
// 			s += fmt.Sprintf("\t%s\tbool %t\n", key, vval)
// 		case string:
// 			s += fmt.Sprintf("\t%s\tstring %s\n", key, vval)
// 		case []interface{}:
// 			s += ParseSlice(key, vval)
// 		case map[string]interface{}:
// 			ProcessMap(vval)
// 		}

// 	}
// 	fmt.Println(s + "\n}")
// 	return data
// }

// func ParseSlice(key string, vals []interface{}) string {
// 	s := fmt.Sprintf(`
// %s struct {
// 	`, key)

// 	switch vals[0].(type) {
// 	case int, int16, int32, int64:
// 		s += fmt.Sprintf("\t%s\t[]int %d\n", key, vals)
// 	case float32, float64:
// 		s += fmt.Sprintf("\t%s\t[]float64 %f\n", key, vals)
// 	case bool:
// 		s += fmt.Sprintf("\t%s\t[]bool %t\n", key, vals)
// 	case string:
// 		s += fmt.Sprintf("\t%s\t[]string %s\n", key, vals)
// 	case []interface{}:
// 		s += fmt.Sprintf("\t%s\t[][] %v\n", key, vals)
// 	}

// 	return s + "\n}"

// }

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

func WriteTempFile(name string, data []byte) {
	lower := strings.ToLower(name)
	data = []byte(fmt.Sprintf("package %s\n\n%s\nfunc NewStruct() %s {\n\tns := make(%s)\n\treturn ns\n}", lower, string(data), name, name))

	os.Mkdir(lower, 0700)
	filename := fmt.Sprintf("%s", lower)
	filepath := fmt.Sprintf("%s/%s.go", lower, filename)
	ioutil.WriteFile(filepath, data, 0700)

	exec.Command("go", "fmt", filepath).Run()
	exec.Command("goimports", "-w", filepath).Run()
	// err := exec.Command("go", "build", filepath).Run()
	// if err != nil {
	// fmt.Println(err)
	// }

	// out, err := exec.Command("./" + filename).Output()
	// if err != nil {
	// 	out = []byte(err.Error())
	// }

	// exec.Command("rm", filename).Run()
	// exec.Command("rm", filepath).Run()
}

func CreateTemplate(name string, data []byte) {

	lower := strings.ToLower(name)
	b, _ := ioutil.ReadFile(fmt.Sprintf("templates/%s.tmpl", lower))

	t := template.New(name)
	t, _ = t.Parse(string(b))
	// t.Execute(os.Stdout, struct {
	// 	LowerName string
	// 	Struct    string
	// 	Name      string
	// }{
	// 	lower,
	// 	string(data),
	// 	name,
	// })

	os.Mkdir(lower, 0700)
	// filename := "main" //fmt.Sprintf("%s", lower)
	filepath := fmt.Sprintf("%s/main.go", lower)
	f, _ := os.Create(filepath)
	os.Chmod(filepath, 0700)

	t.Execute(f, struct {
		LowerName string
		Struct    string
		Name      string
	}{
		lower,
		string(data),
		name,
	})

	exec.Command("go", "fmt", filepath).Run()
	exec.Command("goimports", "-w", filepath).Run()
	if err := exec.Command("go", "build", "-o", fmt.Sprintf("%s/%s", lower, lower), filepath).Run(); err != nil {
		fmt.Println(err)
	}

	// exec.Command("rm", filename).Run()
	// exec.Command("rm", filepath).Run()
}
