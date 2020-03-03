package main

import (
	"errors"
	"fmt"
	"log"
	"time"
)

const year int = 2020

// findTransition attempts to find what date a timezone will transition
// from standard to daylight savings time.
func findTransition(s time.Time, loc *time.Location) (time.Time, error) {
	start, _ := s.Zone()

	for {
		if s.After(time.Date(year, 12, 30, 23, 59, 00, 00, loc)) {
			return time.Time{}, errors.New("didnt find the transition")
		}
		s = s.Add(time.Hour * 24)

		cursor, _ := s.Zone()
		if cursor != start {
			break
		}
	}

	return s, nil
}

func main() {
	locations := []string{
		"Australia/Sydney",
		"America/Los_Angeles",
		"Europe/Zurich",
	}

	for _, l := range locations {
		fmt.Printf("Transitions for timezone: %s\n", l)
		loc, err := time.LoadLocation(l)
		if err != nil {
			log.Fatal(err)
		}

		dayOne := time.Date(year, 1, 1, 0, 0, 0, 0, loc)
		trOne, err := findTransition(dayOne, loc)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("First transition happens on: %s\n", trOne)

		trTwo, err := findTransition(trOne, loc)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Second transition is on: %s\n", trTwo)
	}
}
