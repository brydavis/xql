package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"time"
)

func main() {
	x := json2xml("../data/odd.json")
	fmt.Printf("dataset:\n\n%s\n", x)
}

func json2xml(filename string) (data string) {
	datefmt := "Jan 1, 2006 1:01pm (UTC)"

	base := path.Base(filename)
	base = base[:len(base)-len(path.Ext(base))]
	raw, _ := ioutil.ReadFile(filename)

	var v interface{}
	json.Unmarshal(raw, &v)

	switch t := v.(type) {
	case []interface{}:
		for _, val := range t {
			data += fmt.Sprintf("<%s>%v</%s>\n", base, Walk(val, ""), base)
		}
	case map[string]interface{}:
		for key, val := range t {
			data += fmt.Sprintf("<%s>%v</%s>\n", base, Walk(val, key), base)
		}
	}
	return fmt.Sprintf("<xml created=\"%v\">%s</xml>", time.Now().Format(datefmt), data)

}

func Walk(v interface{}, parent string) interface{} {
	switch t := v.(type) {
	case []interface{}:
		var data string
		for _, val := range t {
			data += fmt.Sprintf("%v", Walk(val, parent))
		}
		return data
	case map[string]interface{}:
		var data string
		for key, val := range t {
			data += fmt.Sprintf("%v", Walk(val, key))
		}
		return data
	case float64:
		return fmt.Sprintf("<%v>%v</%v>", parent, t, parent)
	case bool:
		return fmt.Sprintf("<%v>%v</%v>", parent, t, parent)
	case string:
		return fmt.Sprintf("<%v>%v</%v>", parent, t, parent)
	default:
		return fmt.Sprintf("<%v>%v</%v>", parent, t, parent)
	}
}

// func Indent(tabs int) (indent string) {
// 	for i := 0; i < tabs; i++ {
// 		indent += "\t"
// 	}
// 	return
// }
