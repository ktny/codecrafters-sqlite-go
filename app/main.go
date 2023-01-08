package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	// Available if you need it!
	// "github.com/pingcap/parser"
	// "github.com/pingcap/parser/ast"
)

// Usage: your_sqlite3.sh sample.db .dbinfo
func main() {
	databaseFilePath := os.Args[1]
	command := os.Args[2]

	switch command {
	case ".dbinfo":
		databaseFile, err := os.Open(databaseFilePath)
		if err != nil {
			log.Fatal(err)
		}

		// The first 100 bytes of the databae file comprise the database file header.
		// https://www.sqlite.org/fileformat.html#the_database_header
		header := make([]byte, 100)

		if _, err := databaseFile.Read(header); err != nil {
			log.Fatal(err)
		}

		var pageSize uint16

		// The databse page size in bytes. (Offset:16 Size:2)
		if err := binary.Read(bytes.NewReader(header[16:18]), binary.BigEndian, &pageSize); err != nil {
			fmt.Println("Failed to read integer:", err)
			return
		}

		// Page 1 of a database file is the root page of a table b-tree that holds a special table named "sqlite_schema".
		// https://www.sqlite.org/fileformat.html#storage_of_the_sql_database_schema
		buffer := make([]byte, pageSize)
		if _, err := databaseFile.Read(buffer); err != nil {
			log.Fatal(err)
		}

		numberOfTables := strings.Count(string(buffer), "CREATE TABLE")

		fmt.Printf("database page size: %v", pageSize)
		fmt.Printf("number of tables: %v\n", numberOfTables)
	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}
