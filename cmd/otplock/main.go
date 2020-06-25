package main

import (
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/mjwhitta/log"
	"gitlab.com/mjwhitta/otplock"
)

// Exit status
const (
	Good            int = 0
	InvalidOption   int = 1
	InvalidArgument int = 2
	ExtraArguments  int = 3
	Exception       int = 4
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			if flags.verbose {
				panic(r.(error).Error())
			}
			log.ErrX(Exception, r.(error).Error())
		}
	}()

	var e error
	var otp *otplock.OTPLock
	var sig = make(chan os.Signal)

	validate()

	if otp, e = otplock.New(flags.port); e != nil {
		panic(e)
	}

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Start OTPLock
	go func() {
		otp.Run(flags.unsafe)
	}()

	// Stop OTPLock on ^C
	<-sig
	otp.Stop()
}
