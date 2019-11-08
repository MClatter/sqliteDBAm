package main

import (
	"database/sql"
	"fmt"
	"sort"
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
	synField := strings.Split(noteSlices[36].flds, string(''))[1]
	firstLetter := noteSlices[36].flds[:1]
	synString := strings.Trim(strings.Split(synField, "Synonyms:</small><p class='ex'>")[1], "</p>")
	synSlice := strings.Split(synString, ", ")
	sort.Strings(synSlice)

	firstLetterSlice := []string{}
	for i := 0; i < len(synSlice); i++ {
		if synSlice[i][:1] == firstLetter {
			firstLetterSlice = append(firstLetterSlice, synSlice[i])
			synSlice = append(synSlice[:i], synSlice[i+1:]...)
			i--
		}
	}

	//fmt.Println(synField)
	//fmt.Println(firstLetter)
	fmt.Println(synString)
	fmt.Println(synSlice)
	fmt.Println(firstLetterSlice)
}

/*29: ONE synonym; 32 ...erial; 34 <span class="sentence">contain</span>; 36 many synonyms
 */
