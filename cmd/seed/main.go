package main

import (
	"database/sql"
	"fmt"
	"math/rand/v2"
	"os"
	"path"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

func mustParseDate(dateStr string) time.Time {
	t, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		panic(err)
	}
	return t
}

func projectMention(projectName string) string {
	return fmt.Sprintf(`<p class="para-node"><span class="mention" data-type="mention" data-id="%s" data-label="%s" data-mention-suggestion-char="@" data-mention-id="%s" contenteditable="false">@%s</span></p>`, projectName, projectName, projectName, projectName)
}

func main() {
	args := os.Args[1:]

	db_url := args[0]

	db_path := path.Join("./", db_url)

	if _, err := os.Stat(db_path); err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite", db_path)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		panic(err)
	}

	projectNames := []string{"OKT123", "CIP123", "ECO123", "MAX123"}

	startDate := mustParseDate("2025-01-01")

	for _, seedProject := range projectNames {
		_, err := tx.Exec(
			`INSERT into projects (name)
			VALUES (?)`,
			seedProject,
		)

		if err != nil {
			panic(err)
		}
	}

	for day := range 20 {
		projectCount := rand.IntN(4) + 1
		mentions := []string{}
		excerpts := []string{}

		for i := range projectCount {
			mentions = append(mentions, projectNames[i])
			excerpts = append(excerpts, projectMention(projectNames[i]))
		}

		var id int

		date := startDate.AddDate(0, 0, day).Format(time.DateOnly)

		content := strings.Join(excerpts, " ") + "Meeting notes lorem ipsum dolor sit amet."
		row := tx.QueryRow(`INSERT into notes (html_content, note_date) VALUES (?, ?) RETURNING id`, content, date)

		err = row.Scan(&id)

		if err != nil {
			panic(err)
		}

		for _, project := range mentions {

			_, err := tx.Exec(
				`INSERT INTO project_excerpts (project_name, note_id, excerpt, note_date)
				VALUES (?, ?, ?, ?)`,
				project, id, content,
				date,
			)

			if err != nil {
				panic(err)
			}
		}
	}

	tx.Commit()
}
