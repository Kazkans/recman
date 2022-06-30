package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func findIng(name string, db *sql.DB) {
	rows, _ := db.Query("SELECT recipeID FROM ingredients WHERE name = ?", name) // \""+name+"\""
	defer rows.Close()

	for rows.Next() {
		var recipeID string
		err := rows.Scan(&recipeID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(recipeID)
	}
}

func rmvRecipe(name string, db *sql.DB) {
	stmt, err := db.Prepare("DELETE FROM recipes where name = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(name)
	if err != nil {
		log.Fatal(err)
	}

	stmt, err = db.Prepare("DELETE FROM ingredients where recipeID = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(name)
	if err != nil {
		log.Fatal(err)
	}
}

func addRecipe(rec *recipe, db *sql.DB) {
	stmt, _ := db.Prepare("INSERT INTO recipes (name, time, instructions) VALUES (?, ?, ?)")
	stmt.Exec(rec.name, rec.time, rec.instructions)

	for _, ing := range rec.ingredients {
		stmt, _ := db.Prepare("INSERT INTO ingredients (recipeID, name, quantity) VALUES (?, ?, ?)")
		stmt.Exec(rec.name, ing.name, ing.quantity)
	}
	defer stmt.Close()
}

func getRecipe(name string, db *sql.DB) (rec recipe, err error) {
	row := db.QueryRow("SELECT name, time, instructions FROM recipes WHERE name=?", name)
	err = row.Scan(&rec.name, &rec.time, &rec.instructions)

	if err != nil {
		return
	}

	rows, _ := db.Query("SELECT name, quantity FROM ingredients WHERE recipeID = ?", name) // \""+name+"\""
	defer rows.Close()

	for rows.Next() {
		var ing ingredient
		err = rows.Scan(&ing.name, &ing.quantity)
		if err != nil {
			return
		}

		rec.ingredients = append(rec.ingredients, ing)
	}

	return rec, nil
}

func createTable(db *sql.DB) {
	const create string = `
  CREATE TABLE IF NOT EXISTS ingredients (
  id INTEGER NOT NULL,
  recipeID INTEGER NOT NULL,
  name TEXT NOT NULL,
  quantity TEXT,
  PRIMARY KEY(id AUTOINCREMENT)
  );
  CREATE TABLE IF NOT EXISTS recipes (
	name	TEXT NOT NULL PRIMARY KEY,
	time	INTEGER,
	instructions	TEXT
  );`

	if _, err := db.Exec(create); err != nil {
		log.Fatal(err)
	}
}

// opens database if none then errors
func readDB(filename string) (*sql.DB, error) {
	return sql.Open("sqlite3", filename)
}

// opens database if none then creates
func openDB(filename string) (*sql.DB, error) {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(filename)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	db, err := readDB(filename)
	return db, err
}
