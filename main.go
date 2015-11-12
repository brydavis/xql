package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

import _ "github.com/mattn/go-sqlite3"

// type Base struct {
// 	Name     string
// 	File     os.File
// 	Database string
// 	Query    string
// 	Results  [][][]string // []map[string]interface{}
// 	DB       *sql.DB
// }

func CreateDatabase(filename string) error {
	os.Remove(filename)
	_, err := os.Create(filename)
	return err
}

func CreateTable(base, datafile, tablename string) {
	db, err := sql.Open("sqlite3", "./"+base)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	switch filepath.Ext(datafile) {
	case ".xml":
		// ImportXML(datafile)

	case ".json":
		// b, _ := ioutil.ReadFile(datafile)
		// var f interface{}
		// if err := json.Unmarshal(b, &f); err != nil || f == nil {
		// 	fmt.Println(err)
		// 	break
		// }

		// switch f.(type) {
		// case map[string]interface{}:
		// 	m := f.(map[string]interface{})
		// 	var n []interface{}
		// 	n = append(n, m)
		// 	var o interface{}
		// 	o = n
		// 	ImportJSON(o, db, tablename)
		// case []interface{}:
		// 	ImportJSON(f, db, tablename)
		// }

	case ".csv":
		data, _ := os.Open(datafile)

		reader := csv.NewReader(data)
		reader.FieldsPerRecord = -1

		raw, err := reader.ReadAll()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var headers string
		for _, heads := range raw[0] {
			headers += strings.Replace(strings.Title(heads), " ", "", -1) + ` text, `
		}

		stmt := `
		create table ` + tablename + ` (id integer not null primary key, ` + headers[:len(headers)-2] + `);
		delete from ` + tablename + `;
		`
		_, err = db.Exec(stmt)
		if err != nil {
			log.Printf("%q: %s\n", err, stmt)
			return
		}

		var primaryKey int
		for _, row := range raw[1:] {
			primaryKey++
			values := strconv.Itoa(primaryKey) + `, `
			for _, element := range row {
				values += `"` + element + `", `
			}

			stmt := `insert into ` + tablename + ` values (` + values[:len(values)-2] + `);`
			_, err = db.Exec(stmt)
			if err != nil {
				log.Printf("%q: %s\n", err, stmt)
				return
			}
		}
	default:
		fmt.Println("handle SQL file")

	}

	db.Close()

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

// func (base *Base) ExportTable(table, filename string) error {
// 	switch filepath.Ext(filename) {
// 	case ".csv":
// 		fmt.Println("writing CSV file for table ", table)
// 		base.Query = `select * from ` + table
// 		base.Exec()

// 		file, err := os.Create(filename)
// 		if err != nil {
// 			fmt.Println(err)
// 			return err
// 		}
// 		defer file.Close()

// 		writer := csv.NewWriter(file)
// 		for _, row := range base.Results {
// 			for _, e := range row {
// 				if err := writer.Write(e); err != nil {
// 					fmt.Println(err)
// 					break
// 				}
// 			}
// 		}

// 		writer.Flush()

// 	case ".xml":
// 		fmt.Println("writing XML file for table ", table)
// 	case ".json":
// 		fmt.Println("writing JSON file for table ", table)
// 	default:
// 		fmt.Println("writing SQL file for table ", table)
// 	}

// 	return nil
// }

// func ImportJSON(f interface{}, db *sql.DB, tablename string) {
// 	keys := make(map[string]string)
// 	headerString := " (id integer not null primary key, "
// 	var queries []string
// 	primaryKeys := 1

// 	for _, v := range f.([]interface{}) {
// 		switch vv := v.(type) {
// 		case map[string]interface{}:
// 			columns := "id, "
// 			values := fmt.Sprintf("%d, ", primaryKeys)

// 			for key, val := range vv {

// 				switch vval := val.(type) {
// 				case string:
// 					columns += fmt.Sprintf("%s, ", strings.Replace(key, " ", "", -1))

// 					values += fmt.Sprintf(`"%s", `, val.(string))
// 					if keys[key] == "" {
// 						keys[key] = "text"
// 						headerString += strings.Replace(key, " ", "", -1) + " text, "
// 					}
// 				case int:
// 					columns += fmt.Sprintf("%s, ", strings.Replace(key, " ", "", -1))

// 					values += fmt.Sprintf("%d, ", val)
// 					if keys[key] == "" {
// 						keys[key] = "integer"
// 						headerString += strings.Replace(key, " ", "", -1) + " integer, "
// 					}

// 				case bool:
// 					columns += fmt.Sprintf("%s, ", strings.Replace(key, " ", "", -1))

// 					if vval {
// 						values += fmt.Sprintf("%d, ", 1)
// 					} else {
// 						values += fmt.Sprintf("%d, ", 0)
// 					}

// 					if keys[key] == "" {
// 						keys[key] = "boolean"
// 						headerString += strings.Replace(key, " ", "", -1) + " boolean, "
// 					}

// 				case float64:
// 					columns += fmt.Sprintf("%s, ", strings.Replace(key, " ", "", -1))

// 					values += fmt.Sprintf("%f, ", val)
// 					if keys[key] == "" {
// 						keys[key] = "float"
// 						headerString += strings.Replace(key, " ", "", -1) + " float, "
// 					}
// 				case []interface{}:

// 					_, err := db.Exec(fmt.Sprintf(`create table if not exists %s (id integer);`, strings.Title(key)))
// 					if err != nil {
// 						log.Printf("%q\n", err)
// 						return
// 					}

// 					_, err = db.Exec(fmt.Sprintf(`insert into %s (id) values (%d);`, strings.Title(key), primaryKeys))
// 					if err != nil {
// 						log.Printf("%q\n", err)
// 						return
// 					}

// 				case nil:
// 					columns += fmt.Sprintf("%s, ", strings.Replace(key, " ", "", -1))

// 					values += `"", `
// 				default:
// 					columns += fmt.Sprintf("%s, ", strings.Replace(key, " ", "", -1))

// 					values += `"", `
// 					fmt.Println(vval)
// 				}
// 			}

// 			queries = append(queries, fmt.Sprintf("insert into %s (%s) values (%s);", tablename, columns[:len(columns)-2], values[:len(values)-2]))
// 			primaryKeys += 1
// 		}
// 	}

// 	_, err := db.Exec(`create table ` + tablename + headerString[:len(headerString)-2] + `);`)
// 	if err != nil {
// 		log.Printf("%q\n", err)
// 		return
// 	}

// 	for _, s := range queries {
// 		_, err = db.Exec(s)
// 		if err != nil {
// 			log.Printf("%q: %s\n", err, s)
// 			return
// 		}

// 	}

// }

func ListenAndServe(port int) {
	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/static/", StaticHandler)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}

func StaticHandler(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, req.URL.Path[1:])
}

func RootHandler(res http.ResponseWriter, req *http.Request) {
	db, err := sql.Open("sqlite3", "./foo.db") //+base.Name)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	file, _ := ioutil.ReadFile("views/base.html")
	t := template.New("")
	t, _ = t.Parse(string(file))

	// r, err := db.Query("select sql from sqlite_master")
	// defer r.Close()
	// fmt.Println(r)

	// for r.Next() {
	// 	var schema string
	// 	r.Scan(&schema)

	// 	// fmt.Println(strings.Trim(schema, "("))

	// 	rgx := regexp.MustCompile("([A-z| _-]+[(])|[)]")
	// 	schema = rgx.ReplaceAllString(schema, "")
	// 	s := strings.Split(schema, ",")

	// 	for _, v := range s {
	// 		v = strings.TrimSpace(v)
	// 		ss := strings.Split(v, " ")
	// 		fmt.Println(ss[0])
	// 	}

	// 	fmt.Println(schema)
	// }

	r, _ := db.Query("select name from sqlite_master")
	// defer r.Close()
	tbl_names := make(map[string]map[string]string)

	for r.Next() {
		var n string
		r.Scan(&n)

		tbl_names[n] = make(map[string]string)

	}

	for tbl, _ := range tbl_names {
		r, _ = db.Query(fmt.Sprintf("pragma table_info(%s)", tbl))
		for r.Next() {
			var cid int
			var dflt_value interface{}
			var name string
			var notnull int
			var pk int
			var data_type string

			r.Scan(&cid, &name, &data_type, &dflt_value, &notnull, &pk)

			tbl_names[tbl][name] = data_type

		}
	}


	if req.Method == "GET" {
		t.Execute(res, nil)
	} else {
		req.ParseForm()
		query := req.FormValue("query")
		var output bool
		if req.FormValue("output") == "json" {
			output = true
		}

		fmt.Println(query)

		rows, err := db.Query(fmt.Sprintf("%s;", query))
		defer rows.Close()

		if err != nil {
			fmt.Println(err)
		} else {
			columns, _ := rows.Columns()
			count := len(columns)
			values := make([]interface{}, count)
			valuePtrs := make([]interface{}, count)

			var results []interface{}

			for rows.Next() {
				for i, _ := range columns {
					valuePtrs[i] = &values[i]
				}

				rows.Scan(valuePtrs...)
				store := make(map[string]interface{})
				for i, col := range columns {
					var v interface{}
					val := values[i]
					b, ok := val.([]byte)

					if ok {
						v = string(b)
					} else {
						v = val
					}
					store[col] = v
				}
				// fmt.Println(store)
				results = append(results, store)

			}

			dump, _ := json.Marshal(results)
			schemas, _ := json.Marshal(tbl_names)

			t.Execute(res, struct {
				Query      string
				Results    interface{}
				Schemas    interface{}
				JsonOutput bool
			}{
				query,
				string(dump),
				string(schemas),
				output,
			})
		}
	}

}

func main() {
	port := 8099
	CreateDatabase("foo.db")
	CreateTable("foo.db", "data/crime.csv", "Data")
	fmt.Printf("Ready @ %d\n", port)
	ListenAndServe(port)

}
