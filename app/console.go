package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"text/template"
)

var FuncMap = template.FuncMap{
	"Itemize": Itemize,
}

type Results struct {
	Tables map[string]string
	Query,
	Data string
}

func Itemize(schema string) (results string) {
	re := regexp.MustCompile("[A-z0-9 _]+[(]")
	schemaSplit := strings.Split(strings.Replace(re.ReplaceAllString(schema, ""), ")", "", -1), ",")
	results = "\n"
	for _, v := range schemaSplit {
		vv := strings.Trim(v, " ")
		results += fmt.Sprintf("\t%s\n", vv)
	}
	return
}

func ListenAndServe(port int, conn *sql.DB) {
	http.HandleFunc("/static/", StaticHandler)
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		RootHandler(res, req, conn)
	})

	fmt.Printf("Listening @ %d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}

func StaticHandler(res http.ResponseWriter, req *http.Request) {
	b, _ := ioutil.ReadFile(req.URL.Path[:])

	res.Write(b)
}

func RootHandler(res http.ResponseWriter, req *http.Request, conn *sql.DB) {
	b, _ := ioutil.ReadFile("views/index.html")

	t := template.New("")
	t = t.Funcs(FuncMap)
	t, _ = t.Parse(string(b))

	var tableName string
	var tableSchema string
	schemas := make(map[string]string)

	rows, err := conn.Query("select name, sql from sqlite_master;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&tableName, &tableSchema)
		if err != nil {
			log.Fatal(err)
		}
		// log.Println(tableName, tableSchema)
		schemas[tableName] = tableSchema
	}

	if rows.Err() != nil {
		log.Fatal(err)
	}

	if req.Method == "POST" {
		req.ParseForm()
		query := req.FormValue("query")

		rows, err := conn.Query(query)
		if err != nil {
			fmt.Println(err.Error(), "\n")
		}
		defer rows.Close()

		columns, _ := rows.Columns()
		count := len(columns)
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)

		var store []map[string]interface{}

		for rows.Next() {
			for i, _ := range columns {
				valuePtrs[i] = &values[i]
			}

			rows.Scan(valuePtrs...)

			row := make(map[string]interface{})
			for i, col := range columns {
				var v interface{}
				val := values[i]
				b, ok := val.([]byte)

				if ok {
					v = string(b)
				} else {
					v = val
				}
				row[col] = v
			}

			store = append(store, row)
		}

		b, _ := json.Marshal(store)
		t.Execute(res, Results{Query: query, Data: string(b), Tables: schemas})

	} else {
		t.Execute(res, Results{Query: "", Data: "", Tables: schemas})

	}

}
