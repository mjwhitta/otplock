package otplock

import "time"

type metadata struct {
	Binary   string
	Compile  string
	Endpoint string
	Expires  time.Time
	Filename string
	GUID     string
	Key      []byte
}
