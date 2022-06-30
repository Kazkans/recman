package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type ingredient struct {
	name     string
	quantity string
}

type recipe struct {
	name         string
	time         int // in minutes
	ingredients  []ingredient
	instructions string
}

type stateFn func(*bufio.Scanner, *recipe) stateFn

func nameFn(scanner *bufio.Scanner, rec *recipe) stateFn {
	names := strings.Fields(scanner.Text())
	if len(names) >= 2 {
		rec.name = strings.ToLower(strings.Join(names[1:], " "))
	}

	return noneState
}

func timeFn(scanner *bufio.Scanner, rec *recipe) stateFn {
	times := strings.Fields(scanner.Text())
	if len(times) == 2 {
		time, err := strconv.Atoi(times[1])
		if err == nil {
			rec.time = time
		}
	}

	return noneState
}

func ingredientsFn(scanner *bufio.Scanner, rec *recipe) stateFn {
	if !scanner.Scan() {
		return nil
	}
	if scanner.Text() == "" {
		return ingredientsFn
	}

	fields := strings.Fields(scanner.Text())
	keyword := fields[0]
	if state, no_none := getState(keyword); no_none {
		return state
	}

	splits := strings.Split(scanner.Text(), ";")
	splits[0], splits[1] = strings.TrimSpace(splits[0]), strings.TrimSpace(splits[1])
	switch len(splits) {
	case 1:
		rec.ingredients = append(rec.ingredients, ingredient{strings.ToLower(splits[0]), ""})
	case 2:
		rec.ingredients = append(rec.ingredients, ingredient{strings.ToLower(splits[1]), splits[0]})
	}

	return ingredientsFn
}

func instructionsFn(scanner *bufio.Scanner, rec *recipe) stateFn {
	if !scanner.Scan() {
		return nil
	}
	if scanner.Text() == "" {
		return instructionsFn
	}

	fields := strings.Fields(scanner.Text())
	keyword := fields[0]
	if state, no_none := getState(keyword); no_none {
		return state
	}

	stripped := strings.TrimSpace(scanner.Text())
	if rec.instructions != "" {
		rec.instructions += "\n"
	}
	rec.instructions += stripped

	return instructionsFn
}

func noneState(scanner *bufio.Scanner, rec *recipe) stateFn {
	if !scanner.Scan() {
		return nil
	}
	if scanner.Text() == "" {
		return noneState
	}

	keyword := strings.ToLower(strings.Fields(scanner.Text())[0])

	state, _ := getState(keyword)
	return state
}

// returns state function and if it is not noneState
func getState(keyword string) (stateFn, bool) {
	switch keyword {
	case "ingredients:":
		return ingredientsFn, true
	case "instructions:":
		return instructionsFn, true
	case "time:":
		return timeFn, true
	case "name:":
		return nameFn, true
	default:
		return noneState, false
	}
}

func parse(input string) recipe {
	f, err := os.Open(input)

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	rec := recipe{}
	for state := noneState; state != nil; {
		state = state(scanner, &rec)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return rec
}

/*
type recipe struct {
	name string
	time int // in minutes
	ingredients []ingredient
	instructions string
}*/

func printRecipe(rec *recipe) {
	recipe := ""
	recipe += "Name: " + rec.name + "\n\n"
	if rec.time != 0 {
		recipe += "Time: " + strconv.Itoa(rec.time) + "\n\n"
	}
	recipe += "Ingredients:\n"
	for _, ing := range rec.ingredients {
		if ing.quantity != "" {
			recipe += ing.quantity + ";"
		}
		recipe += ing.name + "\n"
	}
	recipe += "\nInstructions:\n" + rec.instructions

	fmt.Println(recipe)
}
