##  db-cli
 [![Build Status](https://travis-ci.org/moshloop/db-cli.svg?branch=master)](https://travis-ci.org/moshloop/db-cli)

A CLI tool for executing SQL or importing JSON files into databases

```
Usage: db-cli

Utility to execute SQL or import a JSON file against a DB
Options:
  -driver string
      Database driver name, e.g. mssql, mysql, postgres  (default "mssql")
  -file string
      JSON array to insert into a table
  -pass string
      Database password
  -sql string
      SQL Command to run
  -table string
      Table name to insert into
  -url string
      Database connection URL
  -user string
      Database username
  -version
      show version
```