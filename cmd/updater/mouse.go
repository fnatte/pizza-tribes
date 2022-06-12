package main

import (
	"math/rand"
	"strings"

	"github.com/fnatte/pizza-tribes/internal/models"

	_ "embed"
)

//go:embed names.txt
var namesTxt string

var names []string

func init() {
	names = strings.Split(namesTxt, "\n")
}

func GetNewMouseName(existingMice map[string]*models.Mouse) string {
	usedNames := []string{}
	for _, m := range existingMice {
		usedNames = append(usedNames, m.Name)
	}

	for {
		name := names[rand.Intn(len(names))]

		for _, usedName := range usedNames {
			if name == usedName {
				continue
			}
		}

		return name
	}
}
