package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	// Available if you need it!
	// "github.com/pingcap/parser"
	// "github.com/pingcap/parser/ast"
)

// Usage: your_sqlite3.sh sample.db .dbinfo
func main() {
	databaseFilePath := os.Args[1]
	command := os.Args[2]

	databaseFile, err := os.Open(databaseFilePath)
	if err != nil {
		log.Fatal(err)
	}

	header := extractHeader(databaseFile)
	pageSize := extractPageSize(header)
	sqliteSchemaPage := extractSQLiteSchemaPage(databaseFile, pageSize)
	tables := extractTables(sqliteSchemaPage)

	switch command {
	case ".dbinfo":
		fmt.Printf("database page size: %v", pageSize)
		fmt.Printf("number of tables: %v\n", len(tables))
	case ".tables":
		fmt.Printf("%v\n", strings.Join(tables, " "))
	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}

func extractHeader(databaseFile *os.File) []byte {
	// The first 100 bytes of the databae file comprise the database file header.
	// https://www.sqlite.org/fileformat.html#the_database_header
	header := make([]byte, 100)

	if _, err := databaseFile.Read(header); err != nil {
		log.Fatal(err)
	}

	return header
}

func extractPageSize(header []byte) uint16 {
	var pageSize uint16

	// The databse page size in bytes. (Offset:16 Size:2)
	if err := binary.Read(bytes.NewReader(header[16:18]), binary.BigEndian, &pageSize); err != nil {
		fmt.Println("Failed to read integer:", err)
		log.Fatal()
	}

	return pageSize
}

func extractSQLiteSchemaPage(databaseFile *os.File, pageSize uint16) []byte {
	// Page 1 of a database file is the root page of a table b-tree that holds a special table named "sqlite_schema".
	// https://www.sqlite.org/fileformat.html#storage_of_the_sql_database_schema
	buffer := make([]byte, pageSize)
	if _, err := databaseFile.Read(buffer); err != nil {
		log.Fatal(err)
	}

	return buffer
}

func extractTables(buffer []byte) []string {
	var tables []string

	bufferString := string(buffer)
	pattern := regexp.MustCompile(`CREATE TABLE (.+?) \(.+?\)`)
	matches := pattern.FindAllStringSubmatch(bufferString, -1)
	for _, match := range matches {
		tables = append(tables, match[1])
	}
	sort.Strings(tables)

	return tables
}
