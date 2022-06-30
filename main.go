package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"
)

func b2i(b bool) int8 {
	if b {
		return 1
	}
	return 0
}

func main() {
	dbFile := flag.String("d", "recipes.db", "Path to recipe database")

	recFile := flag.String("rec", "", "Path to file that contains recipe")
	ingName := flag.String("i", "", "Ingredient to output recipies that contains it")
	rmvName := flag.String("rm", "", "Name of recipe to remove")
	recName := flag.String("n", "", "Name of recipe to retrieve")

	flag.Parse()

	if b2i(*recFile != "")+b2i(*recName != "")+b2i(*rmvName != "")+b2i(*ingName != "") != 1 {
		fmt.Print("Need to provide path to recipe or name to retrieve")
		return
	}

	db, err := openDB(*dbFile)
	if err != nil {
		log.Fatal(err)
	}

	if *recFile != "" {
		r := parse(*recFile)
		addRecipe(&r, db)
	} else if *recName != "" {
		var r recipe
		r, err = getRecipe(strings.ToLower(*recName), db)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Print("Couldn't find such recipe")
			} else {
				fmt.Print(err)
			}
		} else {
			printRecipe(&r)
		}
	} else if *rmvName != "" {
		rmvRecipe(strings.ToLower(*rmvName), db)
		fmt.Printf("Succesfully removed %v\n", *rmvName)
	} else if *ingName != "" {
		findIng(strings.ToLower(*ingName), db)
	}
}
