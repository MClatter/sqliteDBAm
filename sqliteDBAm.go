// ©, 2019, Михаил Тарабанько.

// A simple throwaway program for updating my Anki database.
// Targeted notes are queried from the db. Fields with synonyms are extracted from them.
// Synonyms are sorted and those with matching first letter are transfered to the hints area.
// Russian and polish words get special treatment.

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

func sortSynsHints(field, firstLetter string) ([]string, []string, []string) {
	fieldSplit := strings.Split(field, ":</small><p class='ex'>")
	if len(fieldSplit) == 1 {
		fieldSplit = strings.Split(field, ":</small><p class=\"ex\">")
	}
	synString := strings.Replace(fieldSplit[1], "</p>", "", 1)
	synSplit := strings.Split(synString, ", ")
	sort.Strings(synSplit)

	firstLetterSlice := []string{}
	ruDefSlice := []string{}
	for i := 0; i < len(synSplit); i++ {
		if synSplit[i][:1] == firstLetter {
			firstLetterSlice = append(firstLetterSlice, synSplit[i])
			synSplit = append(synSplit[:i], synSplit[i+1:]...)
			i--
			continue
		}
		if synSplit[i][0] > 207 {
			ruDefSlice = append(ruDefSlice, synSplit[i])
			synSplit = append(synSplit[:i], synSplit[i+1:]...)
			i--
		}
	}
	return synSplit, firstLetterSlice, ruDefSlice
}

func removePolish(field string) string {
	ruDef := strings.Trim(strings.Split(field, "</big>")[0], "<big>aąbcćdeęfghijklłmnńoóprsśtuwyzźżAĄBCĆDEĘFGHIJKLŁMNŃOÓPRSŚTUWYZ, ")
	polishFreefield := "<big>" + ruDef + "</big>" + strings.Split(field, "</big>")[1]
	return polishFreefield
}

func noteHandler(noteSlices []note) []note {
	for record := range noteSlices {
		field := strings.Split(noteSlices[record].flds, "")[1]
		firstLetter := noteSlices[record].flds[:1]
		if strings.Contains(field, " onclick") {
			noteSlices[record].flds = strings.Replace(noteSlices[record].flds, "", ""+firstLetter+"", 1)
			continue
		}

		if strings.Contains(field, "</p><span class=\"sentence\">") {
			if strings.Contains(field, "<small>Synonym") {
				field = strings.Replace(strings.Replace(field, "</p><span class=\"sentence\">", ", ", -1), "/span", "/p", -1)
			}
			field = strings.Replace(strings.Replace(field, "<span class=\"sentence\">", "<small>Synonyms:</small><p class='ex'>", -1), "/span", "/p", -1)
		}
		if strings.Contains(field, "<small>Synonym") {
			synSplit, firstLetterSlice, ruDefSlice := sortSynsHints(field, firstLetter)
			if len(firstLetterSlice) > 0 {
				firstLetter = strings.Join(firstLetterSlice, ", ")
			}
			if len(ruDefSlice) > 0 {
				field = strings.Replace(field, "<big>", "<big>"+strings.Join(ruDefSlice, ", ")+", ", 1)
			}
			if len(synSplit) > 0 {
				field = strings.SplitAfter(field, "<small>Synonym")[0] + "s:</small><p class='ex'>" + strings.Join(synSplit, ", ") + "</p>"

			} else {
				field = strings.Split(field, "<small>Synonym")[0]
			}
		}
		noteSlices[record].flds = strings.Replace(noteSlices[record].flds, "", ""+firstLetter+"", 1)
		fieldSplit := strings.Split(noteSlices[record].flds, "")
		fieldSplit[1] = removePolish(field)
		noteSlices[record].flds = strings.Join(fieldSplit, "")
	}
	return noteSlices
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

	noteSlices = noteHandler(noteSlices)

	for i := range noteSlices {
		_, err := db.Exec("UPDATE notes SET flds = $1 WHERE id = $2", noteSlices[i].flds, noteSlices[i].id)
		if err != nil {
			panic(err)
		}
	}
}
