package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"regexp"

	_ "github.com/mattn/go-sqlite3"

	"bytes"
	"encoding/hex"
	"io"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
)

var ( // Command Line Options
	commands     = flag.String("sql", "", "SQL Command(s) to run on the data")
	source_text  = flag.String("source", "stdin", "Source file to load, or defaults to stdin")
	delimiter    = flag.String("dlm", ",", "Delimiter between fields -dlm=tab for tab, -dlm=0x## to specify a character code in hex")
	lazyQuotes   = flag.Bool("lazy-quotes", false, "Enable LazyQuotes in the csv parser")
	header       = flag.Bool("header", false, "Treat file as having the first row as a header row")
	outputHeader = flag.Bool("output-header", false, "Display column names in output")
	// tableName    = flag.String("table-name", "tbl", "Override the default table name (tbl)")
	save_to = flag.String("save-to", "", "If set, sqlite3 db is left on disk at this path")
	console = flag.Bool("console", false, "Start web server / terminal-like app with connection to database")
	verbose = flag.Bool("verbose", false, "Enable verbose logging")
)

func FileExt(sourceFileName string) string {
	return path.Ext(sourceFileName)[1:]
}

func GenerateTableName(pathway string) string {
	baseFile := path.Base(pathway)
	flen := len(baseFile)
	extlen := len(path.Ext(baseFile))
	return baseFile[:flen-extlen]
}

func ParseCSV() {

}

func main() {
	// Parse Command Line Flags
	flag.Parse()

	if *console && (*source_text == "stdin") {
		log.Fatalln("Can not open console with pipe input, read a file instead")
	}

	separator := determineSeparator(delimiter)

	// Open Database
	db := openDB(save_to, console)

	// Open the input source
	// var fp *os.File
	fp := openFileOrStdin(source_text)
	defer fp.Close()

	// Init a structured text reader
	reader := csv.NewReader(fp)
	reader.FieldsPerRecord = 0
	reader.Comma = separator
	reader.LazyQuotes = *lazyQuotes

	// Read the first row
	first_row, read_err := reader.Read()

	if read_err != nil {
		log.Fatalln(read_err)
	}

	// type Frame map[string][]string

	// 	file, err := os.Open(filename)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	defer file.Close()

	// 	reader := csv.NewReader(file)
	// 	reader.FieldsPerRecord = -1

	// 	raw, err := reader.ReadAll()
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	var rawface [][]interface{}
	// 	for _, v := range raw {
	// 		// fmt.Println(k, v)
	// 		i := make([]interface{}, len(v))
	// 		for kk, vv := range v {
	// 			i[kk] = vv
	// 		}
	// 		rawface = append(rawface, i)
	// 	}

	// 	df := Frame{}
	// 	for k, head := range raw[0] {
	// 		df[head] = []string{}
	// 		for _, row := range raw[1:] {
	// 			df[head] = append(df[head], row[k])
	// 		}
	// 	}

	var headerRow []string

	if *header {
		headerRow = first_row
		first_row, read_err = reader.Read()

		if read_err != nil {
			log.Fatalln(read_err)
		}
	} else {
		headerRow = make([]string, len(first_row))

		// Name each field after the column
		reStartDigit := regexp.MustCompile("^[0-9]")
		for i := 0; i < len(first_row); i++ {
			if reStartDigit.MatchString(first_row[i]) {
				headerRow[i] = "c" + first_row[i]
			} else {
				headerRow[i] = first_row[i]
			}
		}
	}

	tableName := GenerateTableName(fp.Name())
	// Create the table to load to
	createTable(tableName, &headerRow, db, verbose)

	// Start the clock for importing
	t0 := time.Now()

	// Create transaction
	tx, tx_err := db.Begin()

	if tx_err != nil {
		log.Fatalln(tx_err)
	}

	// Load first row
	stmt := createLoadStmt(tableName, &headerRow, tx)
	loadRow(&first_row, tx, stmt, verbose)

	// Read the data
	for {
		row, file_err := reader.Read()
		if file_err == io.EOF {
			break
		} else if file_err != nil {
			log.Println(file_err)
		} else {
			loadRow(&row, tx, stmt, verbose)
		}
	}
	stmt.Close()
	tx.Commit()

	t1 := time.Now()

	if *verbose {
		fmt.Fprintf(os.Stderr, "Data loaded in: %v\n", t1.Sub(t0))
	}

	// Determine what sql to execute
	sqls_to_execute := strings.Split(*commands, ";")

	t0 = time.Now()

	// Execute given SQL
	for _, sql_cmd := range sqls_to_execute {
		if strings.Trim(sql_cmd, " ") != "" {
			result, err := db.Query(sql_cmd)
			if err != nil {
				log.Fatalln(err)
			}
			displayResult(result, outputHeader, separator)
		}
	}

	t1 = time.Now()
	if *verbose {
		fmt.Fprintf(os.Stderr, "Queries run in: %v\n", t1.Sub(t0))
	}

	if *console {
		ListenAndServe(8080, db)
	}

	// Open console
	// if *console {
	// 	db.Close()
	// 	args := []string{*openPath}
	// 	if *outputHeader {
	// 		args = append(args, "-header")
	// 	}
	// 	cmd := exec.Command("sqlite3", args...)

	// 	cmd.Stdin = os.Stdin
	// 	cmd.Stdout = os.Stdout
	// 	cmd.Stderr = os.Stderr
	// 	cmd_err := cmd.Run()
	// 	if cmd.Process != nil {
	// 		cmd.Process.Release()
	// 	}

	// 	if len(*save_to) == 0 {
	// 		os.RemoveAll(filepath.Dir(*openPath))
	// 	}

	// 	if cmd_err != nil {
	// 		log.Fatalln(cmd_err)
	// 	}
	// } else if len(*save_to) == 0 {
	// 	db.Close()
	// 	os.Remove(*openPath)
	// } else {
	// 	db.Close()
	// }
}

func createTable(tableName string, columnNames *[]string, db *sql.DB, verbose *bool) error {
	var buffer bytes.Buffer

	buffer.WriteString("CREATE TABLE IF NOT EXISTS " + (tableName) + " (")

	for i, col := range *columnNames {
		var col_name string

		reg := regexp.MustCompile(`[^a-zA-Z0-9]`)

		col_name = reg.ReplaceAllString(col, "_")
		if *verbose && col_name != col {
			fmt.Fprintf(os.Stderr, "Column %x renamed to %s\n", col, col_name)
		}

		buffer.WriteString(col_name + " TEXT")

		if i != len(*columnNames)-1 {
			buffer.WriteString(", ")
		}
	}

	buffer.WriteString(");")

	_, err := db.Exec(buffer.String())

	if err != nil {
		log.Fatalln(err)
	}

	return err
}

func createLoadStmt(tableName string, values *[]string, db *sql.Tx) *sql.Stmt {
	if len(*values) == 0 {
		log.Fatalln("Nothing to build insert with!")
	}
	var buffer bytes.Buffer

	buffer.WriteString("INSERT INTO " + (tableName) + " VALUES (")
	for i := range *values {
		buffer.WriteString("?")
		if i != len(*values)-1 {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString(");")
	stmt, err := db.Prepare(buffer.String())
	if err != nil {
		log.Fatalln(err)
	}
	return stmt
}

func loadRow(values *[]string, db *sql.Tx, stmt *sql.Stmt, verbose *bool) error {
	if len(*values) == 0 {
		return nil
	}
	vals := make([]interface{}, 0)
	for _, val := range *values {
		vals = append(vals, val)
	}
	_, err := stmt.Exec(vals...)
	if err != nil && *verbose {
		fmt.Fprintln(os.Stderr, "Bad row: ", err)
	}
	return err
}

type csvWriter struct {
	*csv.Writer
}

func (w csvWriter) put(record []string) {
	if err := w.Write(record); err != nil {
		log.Fatalln(err)
	}
}

func displayResult(rows *sql.Rows, outputHeader *bool, sep rune) {
	cols, cols_err := rows.Columns()

	if cols_err != nil {
		log.Fatalln(cols_err)
	}

	out := csvWriter{csv.NewWriter(os.Stdout)}

	out.Comma = sep

	if *outputHeader {
		out.put(cols)
	}

	rawResult := make([][]byte, len(cols))
	result := make([]string, len(cols))

	dest := make([]interface{}, len(cols))
	for i := range cols {
		dest[i] = &rawResult[i]
	}

	for rows.Next() {
		rows.Scan(dest...)

		for i, raw := range rawResult {
			result[i] = string(raw)
		}

		out.put(result)
	}

	out.Flush()
}

func openFileOrStdin(path *string) *os.File {
	var fp *os.File
	var err error
	if (*path) == "stdin" {
		fp = os.Stdin
		err = nil
	} else {
		fp, err = os.Open(*cleanPath(path))
	}

	if err != nil {
		log.Fatalln(err)
	}

	return fp
}

// func OpenFiles(path *string) []*os.File {
// 	var fps []*os.File
// 	var err error
// 	for _, v := range strings.SplitN(*path, " ", -1) {
// 		if (*path) == "stdin" {
// 			fps = append(fps, os.Stdin)
// 			err = nil
// 		} else {
// 			fps, err = os.Open(*cleanPath(v))
// 		}

// 		if err != nil {
// 			log.Fatalln(err)
// 		}

// 	}

// 	return fps
// }

func cleanPath(path *string) *string {
	var result string
	usr, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}

	if (*path)[:2] == "~/" {
		dir := usr.HomeDir + "/"
		result = strings.Replace(*path, "~/", dir, 1)
	} else {
		result = (*path)
	}

	abs_result, abs_err := filepath.Abs(result)
	if abs_err != nil {
		log.Fatalln(err)
	}

	clean_result := filepath.Clean(abs_result)

	return &clean_result
}

func openDB(path *string, no_memory *bool) *sql.DB {
	currentPath, _ := os.Getwd()
	currentDir := filepath.Base(currentPath)
	dbFile := filepath.Join(currentPath, fmt.Sprintf("%s.db", currentDir))

	// fmt.Printf("currentPath: %s\ncurrentDir: %s\ndbFile: %s\n", currentPath, currentDir, dbFile)

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalln(err)
	}

	return db
}

func determineSeparator(delimiter *string) rune {
	var separator rune

	if (*delimiter) == "tab" {
		separator = '\t'
	} else if strings.Index((*delimiter), "0x") == 0 {
		dlm, hex_err := hex.DecodeString((*delimiter)[2:])

		if hex_err != nil {
			log.Fatalln(hex_err)
		}

		separator, _ = utf8.DecodeRuneInString(string(dlm))
	} else {
		separator, _ = utf8.DecodeRuneInString(*delimiter)
	}
	return separator
}
