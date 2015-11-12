package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

func ImportCSV(datafile, tablename string, db *sql.DB) {

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
		headers += heads + ` text, `
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

}
