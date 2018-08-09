package main

import (
	"flag"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	//_ "gopkg.in/goracle.v2"
	"net/url"
	"os"
	"strings"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"github.com/Jeffail/gabs"

	"encoding/json"
)

var (
	version      bool
	conn         string
	user         string
	file         string
	table        string
	sql          string
	pass         string
	driver       string
	buildVersion = "1.0"
)

func usage() {
	println(`Usage: db-cli [options] [command]

Utility to execute SQL or import a JSON file against a DB
Options:`)
	flag.PrintDefaults()

}

func log(s string) {
	os.Stderr.WriteString(s + "\n")

}

func main() {

	flag.BoolVar(&version, "version", false, "show version")
	flag.StringVar(&conn, "url", "", "Database connection URL")
	flag.StringVar(&user, "user", "", "Database username")
	flag.StringVar(&pass, "pass", "", "Database password")
	flag.StringVar(&driver, "driver", "mssql", "Database driver name, e.g. mssql, mysql, postgres ")
	flag.StringVar(&sql, "sql", "", "SQL Command to run")
	flag.StringVar(&file, "file", "", "JSON array to insert into a table")
	flag.StringVar(&table, "table", "", "Table name to insert into")

	flag.Usage = usage
	flag.Parse()

	if version {
		fmt.Println(buildVersion)
		return
	}

	if flag.NArg() == 0 && flag.NFlag() == 0 {
		usage()
		os.Exit(1)
	}

	u := &url.URL{
		Scheme: strings.Split(conn, "://")[0],
		User:   url.UserPassword(user, pass),
		Host:   strings.Split(strings.Split(conn, "/")[2], "?")[0],
	}

	log("Connecting to " + u.String())
	db, err := sqlx.Connect(driver, u.String())
	if err != nil {
		panic(err)
	}
	db.Ping()
	defer db.Close()

	if sql != "" {
		rows, err := db.Query(sql)
		if err != nil {
			panic(err)
		}

		cols, _ := rows.Columns()
		var results []interface{}
		for rows.Next() {
			columns := make([]interface{}, len(cols))
			columnPointers := make([]interface{}, len(cols))
			for i, _ := range columns {
				columnPointers[i] = &columns[i]
			}

			if err := rows.Scan(columnPointers...); err != nil {
				panic(err)
			}

			m := make(map[string]interface{})
			for i, colName := range cols {
				val := columnPointers[i].(*interface{})
				m[colName] = *val
			}
			println(fmt.Sprintf("%#v", m))
			results = append(results, m)

		}
		data, err := json.Marshal(results)
		if err != nil {
			panic(err)
		}
		os.Stdout.Write(data)

	}

	if file != "" {
		if table == "" {
			panic("Must specify a -table to load into")
		}
		data, _ := ioutil.ReadFile(file)
		json, _ := gabs.ParseJSON(data)
		children, _ := json.Children()
		for _, child := range (children) {
			row, _ := child.ChildrenMap()

			columns := ""
			values := ""
			var arr []interface{}

			for k := range (row) {
				if len(columns) > 0 {
					columns += ","
					values += ","
				}
				columns += k
				values += "?"
				arr = append(arr, fmt.Sprintf("%v", row[k].Data()))
			}
			sql = fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", table, columns, values)
			log(fmt.Sprintf("%#v -> %#v", sql, arr))
			stmt, err := db.Prepare(sql)
			if err != nil {
				panic(err)
			}
			_, err = stmt.Exec(arr...)
			if err != nil {
				panic(err)
			}
		}
	}
}
