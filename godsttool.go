package main

import (
	"errors"
	"fmt"
	"log"
	"sort"
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

// locTransition encapsulates a single DST transition in a location.
type locTransition struct {
	t        time.Time
	l        *time.Location
	zone     string
	timezone string
}

// newLocTransitions returns a slice of DST transitions for a location or an error.
func newLocTransitions(location string) ([]*locTransition, error) {
	loc, err := time.LoadLocation(location)
	if err != nil {
		return nil, err
	}

	dayOne := time.Date(year, 1, 1, 0, 0, 0, 0, loc)
	t1, err := findTransition(dayOne, loc)
	if err != nil {
		return nil, err
	}
	zone1, _ := t1.Zone()

	t2, err := findTransition(t1, loc)
	if err != nil {
		return nil, err
	}
	zone2, _ := t2.Zone()

	return []*locTransition{
		{
			t:        t1,
			zone:     zone1,
			timezone: location,
			l:        loc,
		},
		{
			t:        t2,
			zone:     zone2,
			timezone: location,
			l:        loc,
		}}, nil
}

func main() {
	locations := []string{
		"Australia/Sydney",
		"America/Los_Angeles",
		"Europe/Zurich",
	}

	var ts []*locTransition

	for _, l := range locations {
		lt, err := newLocTransitions(l)
		if err != nil {
			log.Fatal(err)
		}

		ts = append(ts, lt...)
	}

	// Sort the timezone transitions in time order.
	sort.Slice(ts, func(i, j int) bool {
		return ts[i].t.Before(ts[j].t)
	})

	for _, s := range ts {
		fmt.Printf("Transition: %s, Zone: %s\n", s.t, s.zone)
	}
}
