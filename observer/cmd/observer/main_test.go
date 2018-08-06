package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"

	"github.com/omnition/omnition-observer/observer/pkg/envoy"
)

type testCase struct {
	env map[string]string

	adminPort      int
	adminLogPath   string
	tlsEnabled     bool
	tlsCACert      string
	tlsCert        string
	tlsKey         string
	ingressPort    int
	egressPort     int
	tracingDriver  string
	tracingAddress string
	tracingPort    int
}

func TestMain(m *testing.M) {
	err := os.Chdir("../../")
	make := exec.Command("make", "go_build")
	err = make.Run()
	if err != nil {
		fmt.Printf("could not build binary: %v", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestCMDBasic(t *testing.T) {
	dir, err := os.Getwd()

	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(path.Join(dir, "build/observer"))

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	c := envoy.Config{}
	err = yaml.Unmarshal(output, &c)

	assert.Nil(t, err)
	assert.Equal(t, len(c.StaticResources.Listeners), 2)

	ingress := c.StaticResources.Listeners[0]
	egress := c.StaticResources.Listeners[1]

	assert.Equal(t, ingress.Name, "omnition_ingress_listener")
	assert.Equal(t, ingress.Address.SocketAddress.PortValue, 15001)
	assert.Equal(t, len(ingress.FilterChains), 3)
	assert.Equal(t, ingress.FilterChains[0].Filters[0].Config.Tracing.OperationName, "ingress")
	assert.Equal(t, ingress.FilterChains[0].Filters[0].Config.RouteConfig.VirtualHosts[0].Routes[0].Route.Cluster, "http1_ingress_cluster")
	assert.Equal(t, ingress.FilterChains[1].Filters[0].Config.RouteConfig.VirtualHosts[0].Routes[0].Route.Cluster, "http2_ingress_cluster")

	assert.Equal(t, egress.Name, "omnition_egress_listener")
	assert.Equal(t, egress.Address.SocketAddress.PortValue, 15002)
	assert.Equal(t, len(egress.FilterChains), 3)
	assert.Equal(t, egress.FilterChains[0].Filters[0].Config.Tracing.OperationName, "egress")
	assert.Equal(t, egress.FilterChains[0].Filters[0].Config.RouteConfig.VirtualHosts[0].Routes[0].Route.Cluster, "http1_egress_cluster")
	assert.Equal(t, egress.FilterChains[1].Filters[0].Config.RouteConfig.VirtualHosts[0].Routes[0].Route.Cluster, "http2_egress_cluster")
}

func TestCMDWithOptions(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// TODO(owais): Improve tabel tests by storing expected results in the table
	envSets := []map[string]string{
		map[string]string{
			"OBS_ADMIN_PORT":      "2020",
			"OBS_ADMIN_LOG_PATH":  "dasdkjasd",
			"OBS_INGRESS_PORT":    "12345",
			"OBS_TRACING_DRIVER":  "gibberish",
			"OBS_TRACING_ADDRESS": "my-tracing-address",
			"OBS_TRACING_PORT":    "6543",
			"OBS_EGRESS_PORT":     "54321",
		},
		map[string]string{
			"OBS_TLS_ENABLED": "true",
		},
		map[string]string{
			"OBS_TLS_ENABLED": "true",
			"OBS_TLS_CERT":    "some cert",
			"OBS_TLS_KEY":     "some key",
		},
	}

	for _, envSet := range envSets {
		cmd := exec.Command(path.Join(dir, "build/observer"))

		envPairs := []string{}
		for k, v := range envSet {
			envPairs = append(envPairs, k+"="+v)
		}

		cmd.Env = append(os.Environ(), envPairs...)
		output, cmdErr := cmd.CombinedOutput()

		c := envoy.Config{}
		if err == nil {
			err = yaml.Unmarshal(output, &c)
			assert.Nil(t, err)
		}

		for k, v := range envSet {
			switch k {

			case "OBS_ADMIN_PORT":
				assert.Equal(t, v, strconv.Itoa(c.Admin.Address.SocketAddress.PortValue))

			case "OBS_ADMIN_LOG_PATH":
				assert.Equal(t, v, c.Admin.AccessLogPath)

			case "OBS_INGRESS_PORT":
				assert.Equal(t, v, strconv.Itoa(c.StaticResources.Listeners[0].Address.SocketAddress.PortValue))
				assert.Equal(t, "omnition_ingress_listener", c.StaticResources.Listeners[0].Name)

			case "OBS_EGRESS_PORT":
				assert.Equal(t, v, strconv.Itoa(c.StaticResources.Listeners[1].Address.SocketAddress.PortValue))
				assert.Equal(t, "omnition_egress_listener", c.StaticResources.Listeners[1].Name)

			case "OBS_TLS_ENABLED", "OBS_TLS_CERT", "OBS_TLS_KEY":
				if envSet["OBS_TLS_ENABLED"] == "true" {
					if envSet["OBS_TLS_CERT"] == "" || envSet["OBS_TLS_KEY"] == "" {
						assert.NotNil(t, cmdErr)
						assert.Contains(t, string(output), "TLS cannot be enabled without certificate cert and key")
					} else {
						certs := c.StaticResources.Listeners[0].FilterChains[0].TLSContext.CommonTLSContext.TLSCertificates[0]
						assert.Equal(t, envSet["OBS_TLS_CERT"], certs.CertificateChain.InlineString)
						assert.Equal(t, envSet["OBS_TLS_KEY"], certs.PrivateKey.InlineString)
					}
				} else {
					assert.Nil(t, c.StaticResources.Listeners[0].FilterChains[0].TLSContext)
				}
			case "OBS_TRACING_DRIVER":
				// hardcoded to zipkin right now
				assert.Equal(t, "envoy.zipkin", c.Tracing.HTTP.Name)

			case "OBS_TRACING_ADDRESS":
				assert.Equal(t, v, c.StaticResources.Clusters[6].Hosts[0].SocketAddress.Address)

			case "OBS_TRACING_PORT":
				assert.Equal(t, v, strconv.Itoa(c.StaticResources.Clusters[6].Hosts[0].SocketAddress.PortValue))

			default:
				assert.Fail(t, "Don't know how to test env var: "+k)
			}
		}
	}
}
