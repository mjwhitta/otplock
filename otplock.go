package otplock

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	hl "gitlab.com/mjwhitta/hilighter"
	"gitlab.com/mjwhitta/log"
	"gitlab.com/mjwhitta/safety"
)

// OTPLock is a struct containing all metadata required to host an
// HTTP server.
type OTPLock struct {
	Addr        string
	AllowUnsafe bool
	Keys        *safety.Map
	Root        string
	server      *http.Server
	serverMutex *sync.Mutex
	stopped     chan bool
}

// New will return a pointer to a new OTPLock instance.
func New(port int) (otp *OTPLock, e error) {
	// Ensure port is valid
	if (port == 0) || (port > 65535) {
		port = 8080
	}

	// Create OTPLock instance
	otp = &OTPLock{
		Addr:        "0.0.0.0:" + strconv.Itoa(port),
		Keys:        safety.NewMap(),
		Root:        uuid.New().String(),
		serverMutex: &sync.Mutex{},
		stopped:     make(chan bool, 1),
	}

	return
}

// Run will listen for incoming connections and return the requested
// OTP if still valid.
func (otp *OTPLock) Run(allowUnsafe bool) error {
	var e error

	// Use a mutex to ensure this function only runs one at a time
	otp.serverMutex.Lock()

	if otp.server != nil {
		// Make sure to unlock before returning
		otp.serverMutex.Unlock()
		return hl.Errorf("OTPLock is already running")
	}

	// Create HTTP server
	otp.server = &http.Server{Addr: otp.Addr}

	// Store unsafe
	otp.AllowUnsafe = allowUnsafe

	// Unlock now that the server is created
	otp.serverMutex.Unlock()

	otp.server.RegisterOnShutdown(
		func() {
			log.Info("OTPLock is shutting down")

			// Signal that shutdown is complete
			otp.stopped <- true
			close(otp.stopped)
		},
	)

	// Handle all routes
	http.HandleFunc("/", otp.dynamic)

	log.Infof(
		"OTPLock can be configured at http://%s/%s",
		otp.Addr,
		otp.Root,
	)
	e = otp.server.ListenAndServe()

	switch e {
	case http.ErrServerClosed:
		// Expected
		e = nil
	}

	return e
}

// Stop will shutdown the OTPLock instance.
func (otp *OTPLock) Stop() {
	var ctx context.Context
	var cancel context.CancelFunc

	otp.serverMutex.Lock()

	if otp.server != nil {
		ctx, cancel = context.WithTimeout(
			context.Background(),
			4*time.Second,
		)
		defer cancel()

		// Shutdown the server
		otp.server.Shutdown(ctx)
		otp.server = nil

		// Wait for shutdown complete
		<-otp.stopped
	}

	otp.serverMutex.Unlock()
}
