package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mjwhitta/log"
	"github.com/mjwhitta/otplock"
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

	var otp *otplock.OTPLock
	var sig chan os.Signal = make(chan os.Signal, 1)

	validate()

	// Create OTPLock service
	otp = otplock.New(flags.port)

	// Catch ^C
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Start OTPLock
	go func() {
		_ = otp.Run(flags.unsafe)
	}()

	// Stop OTPLock on ^C
	<-sig
	otp.Stop()
}
