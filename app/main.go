package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
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

		_, err = databaseFile.Read(header)
		if err != nil {
			log.Fatal(err)
		}

		var pageSize uint16

		// The databse page size in bytes. (Offset:16 Size:2)
		if err := binary.Read(bytes.NewReader(header[16:18]), binary.BigEndian, &pageSize); err != nil {
			fmt.Println("Failed to read integer:", err)
			return
		}

		fmt.Printf("database page size: %v", pageSize)
	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}
