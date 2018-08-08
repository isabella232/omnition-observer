package options

import (
	"github.com/ansel1/merry"
)

// Tracing system identifiers
const (
	_ = iota
	ZIPKIN
	JAEGER
)

// Config ..
type Options struct {
	TLSEnabled bool
	TLSCACert  string
	TLSCert    string
	TLSKey     string

	AdminPort    int
	AdminLogPath string
	IngressPort  int
	EgressPort   int

	TracingDriver string
	TracingHost   string
	TracingPort   int
}

func New(
	IngressPort int,
	EgressPort int,
	TracingDriver string,
	TracingHost string,
	TracingPort int,
	TLSEnabled bool,
	TLSCACert string,
	TLSCert string,
	TLSKey string,
	AdminPort int,
	AdminLogPath string,
) (Options, error) {
	if TLSEnabled {
		if TLSCert == "" || TLSKey == "" {
			return Options{}, merry.New("TLS cannot be enabled without certificate cert and key")
		}
	}
	return Options{
		IngressPort: IngressPort,
		EgressPort:  EgressPort,

		TracingDriver: TracingDriver,
		TracingHost:   TracingHost,
		TracingPort:   TracingPort,

		TLSEnabled: TLSEnabled,
		TLSCert:    TLSCert,
		TLSCACert:  TLSCACert,
		TLSKey:     TLSKey,

		AdminPort:    AdminPort,
		AdminLogPath: AdminLogPath,
	}, nil
}
