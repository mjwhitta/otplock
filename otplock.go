package otplock

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mjwhitta/errors"
	"github.com/mjwhitta/log"
	"github.com/mjwhitta/safety"
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
func New(port uint16) *OTPLock {
	// Ensure port is valid
	if port == 0 {
		port = 8080
	}

	// Create OTPLock instance
	return &OTPLock{
		Addr:        fmt.Sprintf("0.0.0.0:%d", port),
		Keys:        safety.NewMap(),
		Root:        uuid.New().String(),
		serverMutex: &sync.Mutex{},
		stopped:     make(chan bool, 1),
	}
}

// Run will listen for incoming connections and return the requested
// OTP if still valid.
func (otp *OTPLock) Run(allowUnsafe ...bool) error {
	var e error

	// Use a mutex to ensure this function only runs one at a time
	otp.serverMutex.Lock()

	if otp.server != nil {
		// Make sure to unlock before returning
		otp.serverMutex.Unlock()
		return errors.New("already running")
	}

	// Create HTTP server
	otp.server = &http.Server{
		Addr:              otp.Addr,
		ErrorLog:          newLog(otp),
		ReadHeaderTimeout: 10 * time.Second, //nolint:mnd // 10 secs
	}

	// Unlock now that the server is created
	otp.serverMutex.Unlock()

	// Store unsafe
	otp.AllowUnsafe = (len(allowUnsafe) > 0) && allowUnsafe[0]

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
	default:
		e = errors.Newf("unexpected server error: %w", e)
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
			4*time.Second, //nolint:mnd // 4 secs
		)
		defer cancel()

		// Shutdown the server
		_ = otp.server.Shutdown(ctx)
		otp.server = nil

		// Wait for shutdown complete
		<-otp.stopped
	}

	otp.serverMutex.Unlock()
}

// Write is used by the http.Server to log errors.
func (otp *OTPLock) Write(b []byte) (int, error) {
	log.Err(string(b))
	return len(b), nil
}
