package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

var datefmt = "Jan 1, 2006 1:01pm (UTC)"

func main() {
	x := structify("../../data/odd.json")
	fmt.Printf("dataset:\n\n%s\n", x)
}

func structify(filename) (all string) {
	base := path.Base(filename)
	base = base[:len(base)-len(path.Ext(base))]
	raw, _ := ioutil.ReadFile(filename)

	var v interface{}
	if err := json.Unmarshal(raw, &v); err != nil {
		fmt.Println(err)
	}

	switch t := v.(type) {
	case []interface{}:
		for _, val := range t {
			all += fmt.Sprintf("<%s>%v</%s>\n", base, Walk(val, ""), base)
		}
	case map[string]interface{}:
		for key, val := range t {
			all += fmt.Sprintf("<%s>%v</%s>\n", base, Walk(val, key), base)
		}
	}
	return fmt.Sprintf(`
		type %s struct {
			
		}`, all)

}

func Walk(v interface{}, parent string) interface{} {
	switch t := v.(type) {
	case []interface{}:
		var all string
		for _, val := range t {
			all += fmt.Sprintf("%v", Walk(val, parent))
		}
		return all
	case map[string]interface{}:
		var all string
		for key, val := range t {
			all += fmt.Sprintf("%v", Walk(val, key))
		}
		return all
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
