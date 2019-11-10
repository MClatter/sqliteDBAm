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

func sortOutSynonymsAndHints(separationString, field, firstLetter string) {
	fieldSplit := strings.Split(field, separationString)
	if len(fieldSplit) == 1 {
		fieldSplit = strings.Split(field, ":</small><p class=\"ex\">")
	}
	synString := strings.Trim(fieldSplit[1], "</p>")
	synSplit := strings.Split(synString, ", ")
	sort.Strings(synSplit)

	firstLetterSlice := []string{}
	for i := 0; i < len(synSplit); i++ {
		if synSplit[i][:1] == firstLetter {
			firstLetterSlice = append(firstLetterSlice, synSplit[i])
			synSplit = append(synSplit[:i], synSplit[i+1:]...)
			i--
		}
	}
	fmt.Println(synString)
	fmt.Println(synSplit)
	fmt.Println(firstLetterSlice)
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
	record := 29
	field := strings.Split(noteSlices[record].flds, string(''))[1]
	fmt.Println(field)

	firstLetter := noteSlices[record].flds[:1]
	fmt.Println(firstLetter)

	if strings.Contains(field, "</p><span class=\"sentence\">") && strings.Contains(field, "<small>Synonym") {
		field = strings.Replace(strings.Replace(field, "</p><span class=\"sentence\">", ", ", -1), "/span", "/p", -1)
	}
	if strings.Contains(field, "<small>Synonym") {
		sortOutSynonymsAndHints(":</small><p class='ex'>", field, firstLetter)
	}
}

// TODO:try/catch errors
/*29: ONE synonym; 32 ...erial; 34 <span class="sentence">contain</span>; 36,37 many synonyms
2,3: no syns; 4: one syn
*/
