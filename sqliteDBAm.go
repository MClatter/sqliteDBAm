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

func sortOutSynonymsAndHints(separationString, synField, firstLetter string) {
	synString := strings.Trim(strings.Split(synField, separationString)[1], "</p>")
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
	fmt.Println(synString)
	fmt.Println(synSlice)
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
	record := 4
	synField := strings.Split(noteSlices[record].flds, string(''))[1]
	fmt.Println(synField)

	firstLetter := noteSlices[record].flds[:1]
	fmt.Println(firstLetter)

	if strings.Contains(synField, "</p><span class=\"sentence\">") && strings.Contains(synField, "<small>Synonym") {
		synField = strings.Replace(strings.Replace(synField, "</p><span class=\"sentence\">", ", ", -1), "/span", "/p", -1)
	}
	if strings.Contains(synField, "<small>Synonym") {
		sortOutSynonymsAndHints(":</small><p class='ex'>", synField, firstLetter)
	}
}

// TODO:try/catch errors
/*29: ONE synonym; 32 ...erial; 34 <span class="sentence">contain</span>; 36,37 many synonyms
2,3: no syns; 4: one syn
*/
