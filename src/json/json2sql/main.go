package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

import _ "github.com/mattn/go-sqlite3"

// type Base struct {
// 	Name    string
// 	File    os.File
// 	Query   string
// 	Results [][][]string // []map[string]interface{}
// 	DB      *sql.DB
// }

func ImportJSON(filename string) {
	base := path.Base(filename)
	base = base[:len(base)-len(path.Ext(base))]

	data, _ := ioutil.ReadFile(filename)
	var f interface{}
	if err := json.Unmarshal(data, &f); err != nil || f == nil {
		fmt.Println(err)
		break
	}

	switch f.(type) {
	case map[string]interface{}:
		m := f.(map[string]interface{})
		var n []interface{}
		n = append(n, m)
		var o interface{}
		o = n
		ImportJSON(base)
	case []interface{}:
		ImportJSON(base)
	}

}

// func (base *Base) Exec() error {
// 	db, err := sql.Open("sqlite3", "./"+base.Name)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	rows, err := db.Query(base.Query + `;`)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer rows.Close()

// 	columns, _ := rows.Columns()
// 	count := len(columns)
// 	values := make([]interface{}, count)
// 	valuePtrs := make([]interface{}, count)

// 	for rows.Next() {
// 		for i, _ := range columns {
// 			valuePtrs[i] = &values[i]
// 		}

// 		rows.Scan(valuePtrs...)
// 		// store := make(map[string]interface{})
// 		var store [][]string // [][]string{{"item1", "value1"}, {"item2", "value2"}, {"item3", "value3"}}
// 		for i, col := range columns {
// 			var v interface{}
// 			val := values[i]
// 			b, ok := val.([]byte)

// 			if ok {
// 				v = string(b)
// 			} else {
// 				v = val
// 			}
// 			switch v.(type) {
// 			case string:
// 				store = append(store, []string{col, v.(string)})
// 			case int, int32, int64:
// 				store = append(store, []string{col, strconv.Itoa(int(v.(int64)))})
// 			}
// 		}
// 		// base.Results = append(base.Results, store)
// 		base.Results = append(base.Results, store)
// 	}
// 	return nil
// }

func ScanJSON(f interface{}, db *sql.DB, tablename string) {
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

func main() {
	ImportJSON("foo.db")
	// base.CreateTable("complex.json", "Accounts")
	// base.Select("ProductTitle", "Price", "Color").From("products")
	// base.Query = `select "Price", "Color", 50 as Fifty from products`
	// if err := base.Exec(); err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(base.Results, "\n")

	// base.ExportTable("products", "products.csv")

	// base.ListenAndServe(8099)

	// base.CreateTable("data/data-example-1.xml", "Products")
}
