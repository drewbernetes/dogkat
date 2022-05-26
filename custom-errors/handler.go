package custom_errors

import "log"

func LogFatalError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
