package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"time"
)

func main() {
	x := J2X("../../../data/complex.json")
	fmt.Printf("dataset:\n\n%s\n", x)
}

func J2X(filename string) (data string) {
	base := path.Base(filename)
	base = base[:len(base)-len(path.Ext(base))]
	raw, _ := ioutil.ReadFile(filename)

	var v interface{}
	json.Unmarshal(raw, &v)

	switch t := v.(type) {
	case []interface{}:
		for _, val := range t {
			data += fmt.Sprintf("<%s>%v</%s>\n", base, walk(val, ""), base)
		}
	case map[string]interface{}:
		for key, val := range t {
			// data += fmt.Sprintf("<%s>%v</%s>\n", base, walk(val, key), base)
			data += fmt.Sprintf("%v", walk(val, key))

		}
	}
	return fmt.Sprintf("<xml created=\"%v\">%s</xml>", time.Now().Format("Jan 1, 2006 1:01pm (UTC)"), data)
}

func walk(v interface{}, parent string) interface{} {
	switch t := v.(type) {
	case []interface{}:
		var data string
		for _, val := range t {
			data += fmt.Sprintf("%v", walk(val, parent))
		}
		return data
	case map[string]interface{}:
		var data string

		data += fmt.Sprintf("<%s>", parent)

		for key, val := range t {
			data += fmt.Sprintf("%v", walk(val, key))
		}

		data += fmt.Sprintf("</%s>", parent)

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
