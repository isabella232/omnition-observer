package options

import (
	"strings"
	"time"

	"github.com/ansel1/merry"
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

	TracingDriver     string
	TracingHost       string
	TracingPort       int
	TracingTagHeaders []string

	TimeoutDuration time.Duration
}

func New(
	ingressPort int,
	egressPort int,
	tracingDriver string,
	tracingHost string,
	tracingPort int,
	tracingTagHeaders []string,
	tlsEnabled bool,
	tlsCACert string,
	tlsCert string,
	tlsKey string,
	adminPort int,
	adminLogPath string,
	timeoutDuration time.Duration,
) (Options, error) {
	if tlsEnabled {
		if tlsCert == "" || tlsKey == "" {
			return Options{}, merry.New("TLS cannot be enabled without certificate cert and key")
		}
	}

	// Defaulting to zipkin
	if strings.Trim(tracingDriver, " ") == "" {
		tracingDriver = "zipkin"
	}

	return Options{
		IngressPort: ingressPort,
		EgressPort:  egressPort,

		TracingDriver:     tracingDriver,
		TracingHost:       tracingHost,
		TracingPort:       tracingPort,
		TracingTagHeaders: tracingTagHeaders,

		TLSEnabled: tlsEnabled,
		TLSCert:    tlsCert,
		TLSCACert:  tlsCACert,
		TLSKey:     tlsKey,

		AdminPort:    adminPort,
		AdminLogPath: adminLogPath,

		TimeoutDuration: timeoutDuration,
	}, nil
}
