package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"
	"time"
)

func ImportJSON(filename string) {
	base := path.Base(filename)
	base = base[:len(base)-len(path.Ext(base))]

	data, _ := ioutil.ReadFile(datafile)
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil || v == nil {
		fmt.Println(err)
		break
	}

	switch v.(type) {
	case map[string]interface{}:

		ScanJSON([]interface{}{v.(map[string]interface{})}, base)

	case []interface{}:

		ScanJSON(v, base)

	}

}

func ScanJSON(f interface{}, tablename string) {
	keys := make(map[string]string)
	headerString := " (_id integer not null primary key, "
	var queries []string
	primaryKeys := 1

	for _, v := range f.([]interface{}) {
		switch vv := v.(type) {
		case map[string]interface{}:
			columns := "id, "
			values := fmt.Sprintf("%d, ", primaryKeys)

			for key, val := range vv {

				switch vval := val.(type) {
				case string:
					columns += fmt.Sprintf("%s, ", strings.Replace(key, " ", "", -1))

					values += fmt.Sprintf(`"%s", `, val.(string))
					if keys[key] == "" {
						keys[key] = "text"
						headerString += strings.Replace(key, " ", "", -1) + " text, "
					}
				case int:
					columns += fmt.Sprintf("%s, ", strings.Replace(key, " ", "", -1))

					values += fmt.Sprintf("%d, ", val)
					if keys[key] == "" {
						keys[key] = "integer"
						headerString += strings.Replace(key, " ", "", -1) + " integer, "
					}

				case bool:
					columns += fmt.Sprintf("%s, ", strings.Replace(key, " ", "", -1))

					if vval {
						values += fmt.Sprintf("%d, ", 1)
					} else {
						values += fmt.Sprintf("%d, ", 0)
					}

					if keys[key] == "" {
						keys[key] = "boolean"
						headerString += strings.Replace(key, " ", "", -1) + " boolean, "
					}

				case float64:
					columns += fmt.Sprintf("%s, ", strings.Replace(key, " ", "", -1))

					values += fmt.Sprintf("%f, ", val)
					if keys[key] == "" {
						keys[key] = "float"
						headerString += strings.Replace(key, " ", "", -1) + " float, "
					}
				case []interface{}:

					_, err := db.Exec(fmt.Sprintf(`create table if not exists %s (id integer);`, strings.Title(key)))
					if err != nil {
						log.Printf("%q\n", err)
						return
					}

					_, err = db.Exec(fmt.Sprintf(`insert into %s (id) values (%d);`, strings.Title(key), primaryKeys))
					if err != nil {
						log.Printf("%q\n", err)
						return
					}

				case nil:
					columns += fmt.Sprintf("%s, ", strings.Replace(key, " ", "", -1))

					values += `"", `
				default:
					columns += fmt.Sprintf("%s, ", strings.Replace(key, " ", "", -1))

					values += `"", `
					fmt.Println(vval)
				}
			}

			queries = append(queries, fmt.Sprintf("insert into %s (%s) values (%s);", tablename, columns[:len(columns)-2], values[:len(values)-2]))
			primaryKeys += 1
		}
	}

	_, err := db.Exec(`create table ` + tablename + headerString[:len(headerString)-2] + `);`)
	if err != nil {
		log.Printf("%q\n", err)
		return
	}

	for _, s := range queries {
		_, err = db.Exec(s)
		if err != nil {
			log.Printf("%q: %s\n", err, s)
			return
		}

	}

}

func json2xml(filename string) (all string) {
	datefmt := "Jan 1, 2006 1:01pm (UTC)"

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
	return fmt.Sprintf("<xml created=\"%v\">%s</xml>", time.Now().Format(datefmt), all)

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
