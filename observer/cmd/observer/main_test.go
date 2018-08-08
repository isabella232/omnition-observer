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
	TracingDriver  string
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

	assert.Equal(t, ingress.Name, "ingress_listener")
	assert.Equal(t, ingress.Address.SocketAddress.PortValue, 15001)
	assert.Equal(t, len(ingress.FilterChains), 3)
	assert.Equal(t, ingress.FilterChains[0].Filters[0].Config.Tracing.OperationName, "ingress")
	assert.Equal(t, ingress.FilterChains[0].Filters[0].Config.RouteConfig.VirtualHosts[0].Routes[0].Route.Cluster, "h1_ingress_cluster")
	assert.Equal(t, ingress.FilterChains[1].Filters[0].Config.RouteConfig.VirtualHosts[0].Routes[0].Route.Cluster, "h2_ingress_cluster")

	assert.Equal(t, egress.Name, "egress_listener")
	assert.Equal(t, egress.Address.SocketAddress.PortValue, 15002)
	assert.Equal(t, len(egress.FilterChains), 3)
	assert.Equal(t, egress.FilterChains[0].Filters[0].Config.Tracing.OperationName, "egress")
	assert.Equal(t, egress.FilterChains[0].Filters[0].Config.RouteConfig.VirtualHosts[0].Routes[0].Route.Cluster, "h1_egress_cluster")
	assert.Equal(t, egress.FilterChains[1].Filters[0].Config.RouteConfig.VirtualHosts[0].Routes[0].Route.Cluster, "h2_egress_cluster")

	assert.Equal(t, c.StaticResources.Clusters[0].Name, "h1_ingress_cluster")
	assert.Equal(t, c.StaticResources.Clusters[1].Name, "h1_egress_cluster")
	assert.Equal(t, c.StaticResources.Clusters[2].Name, "h2_ingress_cluster")
	assert.Equal(t, c.StaticResources.Clusters[3].Name, "h2_egress_cluster")
	assert.Equal(t, c.StaticResources.Clusters[4].Name, "tcp_ingress_cluster")
	assert.Equal(t, c.StaticResources.Clusters[5].Name, "tcp_egress_cluster")
	assert.Equal(t, c.StaticResources.Clusters[6].Name, "tracing_zipkin_cluster")

}

func TestCMDWithOptions(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// TODO(owais): Improve tabel tests by storing expected results in the table
	envSets := []map[string]string{
		map[string]string{
			"OBS_ADMIN_PORT":     "2020",
			"OBS_ADMIN_LOG_PATH": "dasdkjasd",
			"OBS_INGRESS_PORT":   "12345",
			"OBS_TRACING_DRIVER": "gibberish",
			"OBS_TRACING_HOST":   "my-tracing-address",
			"OBS_TRACING_PORT":   "6543",
			"OBS_EGRESS_PORT":    "54321",
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
				assert.Equal(t, "ingress_listener", c.StaticResources.Listeners[0].Name)

			case "OBS_EGRESS_PORT":
				assert.Equal(t, v, strconv.Itoa(c.StaticResources.Listeners[1].Address.SocketAddress.PortValue))
				assert.Equal(t, "egress_listener", c.StaticResources.Listeners[1].Name)

			case "OBS_TRACING_DRIVER":
				// hardcoded to zipkin right now
				assert.Equal(t, "envoy.zipkin", c.Tracing.HTTP.Name)

			case "OBS_TRACING_HOST":
				assert.Equal(t, v, c.StaticResources.Clusters[6].Hosts[0].SocketAddress.Address)

			case "OBS_TRACING_PORT":
				assert.Equal(t, v, strconv.Itoa(c.StaticResources.Clusters[6].Hosts[0].SocketAddress.PortValue))

			case "OBS_TLS_ENABLED", "OBS_TLS_CERT", "OBS_TLS_KEY":
				if envSet["OBS_TLS_ENABLED"] == "true" {
					if envSet["OBS_TLS_CERT"] == "" || envSet["OBS_TLS_KEY"] == "" {
						assert.NotNil(t, cmdErr)
						assert.Contains(t, string(output), "TLS cannot be enabled without certificate cert and key")
					} else {

						httpChain := c.StaticResources.Listeners[0].FilterChains[0]
						http2Chain := c.StaticResources.Listeners[0].FilterChains[1]

						certs1 := httpChain.TLSContext.CommonTLSContext.TLSCertificates[0]
						assert.Equal(t, envSet["OBS_TLS_CERT"], certs1.CertificateChain.InlineString)
						assert.Equal(t, envSet["OBS_TLS_KEY"], certs1.PrivateKey.InlineString)
						route1 := httpChain.Filters[0].Config.RouteConfig.VirtualHosts[0].Routes[0]
						assert.False(t, route1.Redirect.HTTPSRedirect)
						assert.Empty(t, route1.Redirect.PathRedirect)
						assert.Equal(t, "h1_ingress_cluster", route1.Route.Cluster)

						certs2 := http2Chain.TLSContext.CommonTLSContext.TLSCertificates[0]
						assert.Equal(t, envSet["OBS_TLS_CERT"], certs2.CertificateChain.InlineString)
						assert.Equal(t, envSet["OBS_TLS_KEY"], certs2.PrivateKey.InlineString)
						route2 := http2Chain.Filters[0].Config.RouteConfig.VirtualHosts[0].Routes[0]
						assert.False(t, route2.Redirect.HTTPSRedirect)
						assert.Empty(t, route2.Redirect.PathRedirect)
						assert.Equal(t, "h2_ingress_cluster", route2.Route.Cluster)

						route3 := c.StaticResources.Listeners[0].FilterChains[2].Filters[0].Config.RouteConfig.VirtualHosts[0].Routes[0]
						assert.Empty(t, route3.Route.Cluster)
						assert.Equal(t, true, route3.Redirect.HTTPSRedirect)
						assert.Equal(t, "/", route3.Redirect.PathRedirect)

						route4 := c.StaticResources.Listeners[0].FilterChains[3].Filters[0].Config.RouteConfig.VirtualHosts[0].Routes[0]
						assert.Empty(t, route4.Route.Cluster)
						assert.Equal(t, true, route4.Redirect.HTTPSRedirect)
						assert.Equal(t, "/", route4.Redirect.PathRedirect)
					}
				} else {
					assert.Nil(t, c.StaticResources.Listeners[0].FilterChains[0].TLSContext)
				}

			default:
				assert.Fail(t, "Don't know how to test env var: "+k)
			}
		}
	}
}
