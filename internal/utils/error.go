package utils

import (
	"github.com/rs/zerolog/log"
)

func CheckFatalError(err error, errMessage string) {
	if err == nil {
		return
	}

	log.Fatal().Msgf("%s (%s)", errMessage, err)
}
