package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

var Tables []string

func main() {
	EncodeJSON("complex.json")

	fmt.Printf("\n%+v\n", Tables)
}

func EncodeJSON(filename string) {
	base := path.Base(filename)
	table := base[:len(base)-len(path.Ext(base))]
	raw, _ := ioutil.ReadFile(filename)

	var v interface{}
	if err := json.Unmarshal(raw, &v); err != nil {
		fmt.Println(err)
	}

	CreateTable(table, v)
}

func CreateTable(table string, v interface{}) {
	all := "_id int,\n"

	switch t := v.(type) {
	case []interface{}:
		for _, val := range t {
			all += Walk(val, table)
		}
	case map[string]interface{}:
		m := map[string]bool{}

		for key, val := range t {
			if !m[key] {
				all += Walk(val, key)
			}
		}
	}
	Tables = append(Tables, fmt.Sprintf("create table %s (\n%s\n);\n", table, all[:len(all)-2]))
}

func Walk(v interface{}, parent string) string {
	switch t := v.(type) {
	case []interface{}:

		CreateTable(parent, v)
		return fmt.Sprintf("%s_id int,\n", parent)

		// var all string
		// for _, val := range t {
		// 	all += fmt.Sprintf("%v\n", Walk(val, parent))
		// }
		// return all
	case map[string]interface{}:

		CreateTable(parent, v)
		return fmt.Sprintf("%s_id int,\n", parent)

		// var all string
		// for key, val := range t {
		// 	all += Walk(val, key)
		// }
		// return all
	case float64:
		return fmt.Sprintf("%s float,\n", parent)
	case bool:
		return fmt.Sprintf("%s bool,\n", parent)
	case string:
		return fmt.Sprintf("%s varchar(%d),\n", parent, len(t))
	default:
		return fmt.Sprintf("%s text,\n", parent)
	}
}
