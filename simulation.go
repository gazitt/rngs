package main

import (
	"fmt"
)

type Simulation struct {
	allnames map[string]struct{}
	oldNames map[string]struct{}
	newNames map[string]struct{}
	exists   func(string) bool
	force    bool
}

func GetSimulationFunc(force bool, a []string, existsfunc func(string) bool) func(string, string) error {
	s := new(Simulation)
	s.oldNames = make(map[string]struct{}, len(a))
	for _, v := range a {
		s.oldNames[v] = struct{}{}
	}
	s.allnames = s.oldNames
	s.newNames = make(map[string]struct{}, len(a))
	s.exists = existsfunc
	return s.simulate
}

func (s *Simulation) simulate(oldname, newname string) error {
	if s.force {
		return nil
	}

	if _, alreadyThere := s.oldNames[newname]; alreadyThere {
		return fmt.Errorf("Conflict: %s", newname)
	}
	if _, alreadyThere := s.newNames[newname]; alreadyThere {
		return fmt.Errorf("Conflict: %s", newname)
	}
	if _, alreadyThere := s.allnames[newname]; !alreadyThere && s.exists(newname) {
		return fmt.Errorf("Conflict: %s", newname)
	}
	delete(s.oldNames, oldname)
	s.newNames[newname] = struct{}{}
	return nil
}
