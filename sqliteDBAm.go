package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type note struct {
	id   int
	flds string
}

func main() {
	db, err := sql.Open("sqlite3", "C:/Users/Michael/Documents/SQLite scripts/collection.anki2")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, flds FROM notes WHERE flds LIKE '%%' AND tags LIKE '%SMIntermediate%'")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	noteSlices := []note{}
	for rows.Next() {
		noteVar := note{}
		err := rows.Scan(&noteVar.id, &noteVar.flds)
		if err != nil {
			fmt.Println(err)
			continue
		}
		noteSlices = append(noteSlices, noteVar)
	}
	synonymsField := strings.Split(noteSlices[34].flds, string(''))[1]
	firstLetter := noteSlices[30].flds[:1]
	fmt.Println(synonymsField) //29: ONE synonym; 32 ...erial
	fmt.Println(firstLetter)

	fmt.Println(len(noteSlices))
}
