package otplock

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	hl "gitlab.com/mjwhitta/hilighter"
	"gitlab.com/mjwhitta/log"
	"gitlab.com/mjwhitta/pathname"
)

// advanced will process advanced level requests.
func (otp *OTPLock) advanced(
	w http.ResponseWriter,
	r *http.Request,
	srcMeta metadata,
) {
	var bytes []byte
	var e error
	var expires time.Duration
	var guid string
	var key []byte
	var meta metadata
	var payload string
	var t int
	var tmp []byte
	var wg = sync.WaitGroup{}

	if e = r.ParseForm(); e != nil {
		hl.Fprintf(w, errPg, e.Error())
		return
	}

	// Get values from form
	t, e = strconv.Atoi(r.Form.Get("expires"))
	if e != nil {
		hl.Fprintf(w, errPg, e.Error())
		return
	}

	expires, e = time.ParseDuration(strconv.Itoa(t) + "s")
	if e != nil {
		hl.Fprintf(w, errPg, e.Error())
		return
	}

	payload = r.Form.Get("payload")
	if payload == "" {
		hl.Fprintf(w, errPg, "No payload provided")
		return
	}
	payload = strings.Join(strings.Fields(payload), "")

	// Convert hex payload to byte array
	if tmp, e = hex.DecodeString(payload); e != nil {
		hl.Fprintf(w, errPg, e.Error())
		return
	}

	// Generate a random key of the same length as payload
	key = make([]byte, len(tmp))
	if _, e = io.ReadFull(rand.Reader, key); e != nil {
		hl.Fprintf(w, errPg, e.Error())
		return
	}

	// Store metadata
	guid = uuid.New().String()
	meta = metadata{
		Expires: time.Now().Add(expires),
		GUID:    guid,
		Key:     key,
	}
	otp.Keys.Put(guid, meta)

	log.Goodf(
		"User %s created key %s valid until %s",
		r.RemoteAddr,
		guid,
		meta.Expires.Format(time.RFC3339),
	)

	// Encrypt payload with OTP key
	for i, b := range key {
		tmp[i] ^= b
	}

	log.SubInfo("Compiling payload for user " + r.RemoteAddr)

	// Compile in background thread
	wg.Add(1)
	go func() {
		var dir string
		var f *os.File
		var src string

		defer wg.Done()

		// Change to the correct directory
		dir = filepath.Join("wwwotp", srcMeta.GUID)
		if e = os.Chdir(dir); e != nil {
			return
		}

		// Read the source template
		bytes, e = os.ReadFile(srcMeta.Filename + ".template")
		if e != nil {
			return
		}

		// Replace dynamic variables
		src = strings.ReplaceAll(
			string(bytes),
			"ENCHEX",
			hex.EncodeToString(tmp),
		)
		src = strings.ReplaceAll(
			src,
			"OTPURL",
			srcMeta.Endpoint+"/"+guid,
		)

		// Write source to file
		if f, e = os.Create(srcMeta.Filename); e != nil {
			return
		}
		defer f.Close()

		f.WriteString(src)

		// Commpile source
		execute(srcMeta.Compile)

		// Read the compiled binary
		if bytes, e = os.ReadFile(srcMeta.Binary); e != nil {
			return
		}
	}()
	wg.Wait()

	if e != nil {
		hl.Fprintf(w, errPg, e.Error())
		return
	}

	w.Write(bytes)
}

// config will handle incoming connections to the configuration
// dashboard.
func (otp *OTPLock) config(w http.ResponseWriter, r *http.Request) {
	var e error

	if e = r.ParseForm(); e != nil {
		hl.Fprintf(w, errPg, e.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		// If GET, return config page
		switch r.Form.Get("level") {
		case "advanced":
			if otp.AllowUnsafe {
				hl.Fprint(w, advancedNew)
			} else {
				hl.Fprintf(w, errPg, "Advanced config disabled")
			}
		default:
			hl.Fprint(w, simpleDashboard)
		}
	case http.MethodPost:
		// If POST, process config
		switch r.Form.Get("level") {
		case "advanced":
			otp.newAdv(w, r)
		default:
			otp.simple(w, r)
		}
	}
}

// dynamic will handle all incoming connections.
func (otp *OTPLock) dynamic(w http.ResponseWriter, r *http.Request) {
	var data interface{}
	var guid string = pathname.Basename(r.URL.Path)
	var meta metadata
	var ok bool

	// Main config dashboard
	if guid == otp.Root {
		otp.config(w, r)
		return
	}

	// Check if GUID exists
	if data, ok = otp.Keys.Get(guid); ok {
		meta = data.(metadata)

		if len(meta.Key) != 0 {
			// GUID exists so check if key expired
			if time.Now().Before(meta.Expires) {
				// Key not expired, return it
				log.Goodf("User %s got key %s", r.RemoteAddr, guid)
				w.Write(meta.Key)
			} else {
				// Key expired
				otp.Keys.Delete(guid)
				log.Warnf(
					"User %s attempted to get key %s",
					r.RemoteAddr,
					guid,
				)
				hl.Fprint(w, notFound)
			}
		} else {
			// GUID exists but it wasn't a key
			switch r.Method {
			case http.MethodGet:
				// Return advanced dashboard
				hl.Fprint(w, advancedDashboard)
			case http.MethodPost:
				otp.advanced(w, r, meta)
				return
			}
		}
	} else {
		// GUID not found
		hl.Fprint(w, notFound)
	}
}

// newAdv will process advanced level requests.
func (otp *OTPLock) newAdv(w http.ResponseWriter, r *http.Request) {
	var binary string
	var compile string
	var e error
	var endpoint string
	var f *os.File
	var filename string
	var guid string
	var meta metadata
	var source string

	// Get values from form
	binary = r.Form.Get("binary")
	if binary == "" {
		hl.Fprintf(w, errPg, "No binary name provided")
		return
	}

	endpoint = r.Form.Get("endpoint")
	if endpoint == "" {
		hl.Fprintf(w, errPg, "No endpoint provided")
		return
	}

	filename = pathname.Basename(r.Form.Get("filename"))
	if filename == "" {
		hl.Fprintf(w, errPg, "No filename provided")
		return
	}

	compile = r.Form.Get("compile")
	if compile == "" {
		hl.Fprintf(w, errPg, "No compile command provided")
		return
	}

	source = r.Form.Get("source")
	if source == "" {
		hl.Fprintf(w, errPg, "No template source provided")
		return
	}

	// Store metadata
	guid = uuid.New().String()
	meta = metadata{
		Binary:   binary,
		Compile:  compile,
		Endpoint: endpoint,
		Filename: filename,
		GUID:     guid,
	}
	otp.Keys.Put(guid, meta)

	// Make directory
	os.MkdirAll(filepath.Join("wwwotp", guid), os.ModePerm)

	// Create file
	f, e = os.Create(
		filepath.Join("wwwotp", guid, filename) + ".template",
	)
	if e != nil {
		hl.Fprintf(w, errPg, e.Error())
		return
	}
	defer f.Close()

	// Write file
	f.WriteString(source)

	// Redirect user to advanced page
	hl.Fprintf(w, advancedResp, endpoint+"/"+guid, "/"+guid)
}

// simple will process simple level requests.
func (otp *OTPLock) simple(w http.ResponseWriter, r *http.Request) {
	var e error
	var endpoint string
	var expires time.Duration
	var guid string
	var key []byte
	var meta metadata
	var payload string
	var t int
	var tmp []byte

	// Get values from form
	endpoint = r.Form.Get("endpoint")
	if endpoint == "" {
		hl.Fprintf(w, errPg, "No endpoint provided")
		return
	}

	t, e = strconv.Atoi(r.Form.Get("expires"))
	if e != nil {
		hl.Fprintf(w, errPg, e.Error())
		return
	}

	expires, e = time.ParseDuration(strconv.Itoa(t) + "s")
	if e != nil {
		hl.Fprintf(w, errPg, e.Error())
		return
	}

	payload = r.Form.Get("payload")
	if payload == "" {
		hl.Fprintf(w, errPg, "No payload provided")
		return
	}
	payload = strings.Join(strings.Fields(payload), "")

	// Convert hex payload to byte array
	if tmp, e = hex.DecodeString(payload); e != nil {
		hl.Fprintf(w, errPg, e.Error())
		return
	}

	// Generate a random key of the same length as payload
	key = make([]byte, len(tmp))
	if _, e = io.ReadFull(rand.Reader, key); e != nil {
		hl.Fprintf(w, errPg, e.Error())
		return
	}

	// Store metadata
	guid = uuid.New().String()
	meta = metadata{
		Expires: time.Now().Add(expires),
		Key:     key,
	}
	otp.Keys.Put(guid, meta)

	log.Goodf(
		"User %s created key %s valid until %s",
		r.RemoteAddr,
		guid,
		meta.Expires.Format(time.RFC3339),
	)

	// Encrypt payload with OTP key
	for i, b := range key {
		tmp[i] ^= b
	}

	// Return URL and encrypted payload to user
	hl.Fprintf(
		w,
		simpleResp,
		endpoint+"/"+guid,
		hex.EncodeToString(tmp),
	)
}
