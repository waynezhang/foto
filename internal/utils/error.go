package utils

import (
	"os"

	"github.com/waynezhang/foto/internal/log"
)

func CheckFatalError(err error, errMessage string) {
	if err == nil {
		return
	}

	log.Fatal("%s (%s)", errMessage, err)
	os.Exit(1)
}
