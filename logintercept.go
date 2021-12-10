package otplock

import (
	"io"
	"log"
)

func newLog(w io.Writer) *log.Logger {
	return log.New(w, "", 0)
}
