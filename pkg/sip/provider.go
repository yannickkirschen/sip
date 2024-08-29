package sip

import (
	"os"
)

type Provider func(key Key) (string, bool)

var Providers = []Provider{LookupEnv}

func RegisterProvider(f Provider) {
	Providers = append(Providers, f)
}

func LookupEnv(key Key) (string, bool) {
	return os.LookupEnv(key.ToEnvVar())
}
