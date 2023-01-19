package helpers

import (
	"log"
	"regexp"
)

func HandleCreateError(err error) {
	if err != nil {
		r := regexp.MustCompile("already exists")
		if !r.MatchString(err.Error()) {
			log.Fatalln(err)
		}
	}
}
