package main

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/omnition/omnition-observer/observer/pkg/envoy"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestCMDBasic(t *testing.T) {
	config, err := run()
	assert.Nil(t, err)
	c, err := unmarshalConfig(config)
	assert.Nil(t, err)

	assert.Equal(t, len(c.StaticResources.Listeners), 2)

	ingress := c.StaticResources.Listeners[0]
	egress := c.StaticResources.Listeners[1]

	assert.Equal(t, ingress.Name, "ingress_listener")
	assert.Equal(t, ingress.Address.SocketAddress.PortValue, 15001)
	assert.Equal(t, len(ingress.FilterChains), 3)
	assert.Equal(t, ingress.FilterChains[0].Filters[0].TypedConfig.RouteConfig.VirtualHosts[0].Routes[0].Route.Cluster, "h1_ingress_cluster")
	assert.Equal(t, ingress.FilterChains[1].Filters[0].TypedConfig.RouteConfig.VirtualHosts[0].Routes[0].Route.Cluster, "h2_ingress_cluster")

	assert.Equal(t, egress.Name, "egress_listener")
	assert.Equal(t, egress.Address.SocketAddress.PortValue, 15002)
	assert.Equal(t, len(egress.FilterChains), 3)
	var nilSlice []string
	assert.Equal(t, egress.FilterChains[0].Filters[0].TypedConfig.Tracing.CustomTags, nilSlice)
	assert.Equal(t, egress.FilterChains[0].Filters[0].TypedConfig.RouteConfig.VirtualHosts[0].Routes[0].Route.Cluster, "h1_egress_cluster")
	assert.Equal(t, egress.FilterChains[1].Filters[0].TypedConfig.RouteConfig.VirtualHosts[0].Routes[0].Route.Cluster, "h2_egress_cluster")

	assert.Equal(t, c.StaticResources.Clusters[0].Name, "h1_ingress_cluster")
	assert.Equal(t, c.StaticResources.Clusters[1].Name, "h1_egress_cluster")
	assert.Equal(t, c.StaticResources.Clusters[2].Name, "h2_ingress_cluster")
	assert.Equal(t, c.StaticResources.Clusters[3].Name, "h2_egress_cluster")
	assert.Equal(t, c.StaticResources.Clusters[4].Name, "tcp_ingress_cluster")
	assert.Equal(t, c.StaticResources.Clusters[5].Name, "tcp_egress_cluster")
	assert.Equal(t, c.StaticResources.Clusters[6].Name, "tracing_zipkin_cluster")
}

func TestCMDWithOptions(t *testing.T) {
	t.Run("Succeed with some configs", func(t *testing.T) {
		envVariables := map[string]string{
			"OBS_ADMIN_PORT":          "2020",
			"OBS_ADMIN_LOG_PATH":      "dasdkjasd",
			"OBS_INGRESS_PORT":        "12345",
			"OBS_TRACING_DRIVER":      "zipkin",
			"OBS_TRACING_HOST":        "my-tracing-address",
			"OBS_TRACING_PORT":        "6543",
			"OBS_TRACING_TAG_HEADERS": "header1 header2 header3",
			"OBS_EGRESS_PORT":         "54321",
			"OBS_NUM_TRUSTED_HOPS":    "2",
		}
		setEnvironmentVariables(t, envVariables)

		// When
		config, err := run()
		assert.Nil(t, err)
		c, err := unmarshalConfig(config)
		assert.Nil(t, err)

		// Then
		assert.Equal(t, envVariables["OBS_ADMIN_PORT"], strconv.Itoa(c.Admin.Address.SocketAddress.PortValue))
		assert.Equal(t, envVariables["OBS_ADMIN_LOG_PATH"], c.Admin.AccessLogPath)
		assert.Equal(t, envVariables["OBS_INGRESS_PORT"], strconv.Itoa(c.StaticResources.Listeners[0].Address.SocketAddress.PortValue))
		assert.Equal(t, "ingress_listener", c.StaticResources.Listeners[0].Name)
		assert.Equal(t, envVariables["OBS_EGRESS_PORT"], strconv.Itoa(c.StaticResources.Listeners[1].Address.SocketAddress.PortValue))
		assert.Equal(t, "egress_listener", c.StaticResources.Listeners[1].Name)
		// hardcoded to zipkin right now
		assert.Equal(t, "envoy.zipkin", c.Tracing.Http.Name)
		assert.Equal(t, envVariables["OBS_TRACING_HOST"], c.StaticResources.Clusters[6].Hosts[0].SocketAddress.Address)
		assert.Equal(t, envVariables["OBS_TRACING_PORT"], strconv.Itoa(c.StaticResources.Clusters[6].Hosts[0].SocketAddress.PortValue))
		h1Chain := c.StaticResources.Listeners[0].FilterChains[0]
		h2Chain := c.StaticResources.Listeners[0].FilterChains[1]
		headers := strings.Split(envVariables["OBS_TRACING_TAG_HEADERS"], " ")
		assert.Equal(t, headers, h1Chain.Filters[0].TypedConfig.Tracing.CustomTags)
		assert.Equal(t, headers, h2Chain.Filters[0].TypedConfig.Tracing.CustomTags)
		assert.Equal(t, envVariables["OBS_NUM_TRUSTED_HOPS"], strconv.Itoa(h1Chain.Filters[0].TypedConfig.TrustedHopsCount))

		assert.Nil(t, c.StaticResources.Listeners[0].FilterChains[0].TLSContext)
	})

	t.Run("Failing: Missing: OBS_TLS_CERT", func(t *testing.T) {
		// Given
		envVariables := map[string]string{
			"OBS_TLS_ENABLED": "true",
			"OBS_TLS_KEY":     "some key",
			"OBS_TLS_CERT":    "",
		}
		setEnvironmentVariables(t, envVariables)

		// When
		_, err := buildOptions()

		// Then
		assert.NotNil(t, err, "Options instantiation should fail")
	})

	t.Run("Failing: Missing: OBS_TLS_KEY", func(t *testing.T) {
		envVariables := map[string]string{
			"OBS_TLS_ENABLED": "true",
			"OBS_TLS_KEY":     "",
			"OBS_TLS_CERT":    "some cert",
		}
		setEnvironmentVariables(t, envVariables)

		// When
		_, err := buildOptions()

		// Then
		assert.NotNil(t, err, "Options instantiation should fail")
	})

	t.Run("Succeed TLS config", func(t *testing.T) {
		envVariables := map[string]string{
			"OBS_TLS_ENABLED": "true",
			"OBS_TLS_CERT":    "some cert",
			"OBS_TLS_KEY":     "some key",
		}
		setEnvironmentVariables(t, envVariables)

		// When
		config, err := run()
		assert.Nil(t, err)
		c, err := unmarshalConfig(config)
		assert.Nil(t, err)

		// Then
		httpChain := c.StaticResources.Listeners[0].FilterChains[0]
		http2Chain := c.StaticResources.Listeners[0].FilterChains[1]

		certs1 := httpChain.TLSContext.CommonTLSContext.TLSCertificates[0]
		assert.Equal(t, envVariables["OBS_TLS_CERT"], certs1.CertificateChain.InlineString)
		assert.Equal(t, envVariables["OBS_TLS_KEY"], certs1.PrivateKey.InlineString)
		route1 := httpChain.Filters[0].TypedConfig.RouteConfig.VirtualHosts[0].Routes[0]
		assert.False(t, route1.Redirect.HTTPSRedirect)
		assert.Empty(t, route1.Redirect.PathRedirect)
		assert.Equal(t, "h1_ingress_cluster", route1.Route.Cluster)

		certs2 := http2Chain.TLSContext.CommonTLSContext.TLSCertificates[0]
		assert.Equal(t, envVariables["OBS_TLS_CERT"], certs2.CertificateChain.InlineString)
		assert.Equal(t, envVariables["OBS_TLS_KEY"], certs2.PrivateKey.InlineString)
		route2 := http2Chain.Filters[0].TypedConfig.RouteConfig.VirtualHosts[0].Routes[0]
		assert.False(t, route2.Redirect.HTTPSRedirect)
		assert.Empty(t, route2.Redirect.PathRedirect)
		assert.Equal(t, "h2_ingress_cluster", route2.Route.Cluster)

		route3 := c.StaticResources.Listeners[0].FilterChains[2].Filters[0].TypedConfig.RouteConfig.VirtualHosts[0].Routes[0]
		assert.Empty(t, route3.Route.Cluster)
		assert.Equal(t, true, route3.Redirect.HTTPSRedirect)
		assert.Equal(t, "/", route3.Redirect.PathRedirect)

		route4 := c.StaticResources.Listeners[0].FilterChains[3].Filters[0].TypedConfig.RouteConfig.VirtualHosts[0].Routes[0]
		assert.Empty(t, route4.Route.Cluster)
		assert.Equal(t, true, route4.Redirect.HTTPSRedirect)
		assert.Equal(t, "/", route4.Redirect.PathRedirect)
	})
}

func setEnvironmentVariables(t *testing.T, envVars map[string]string) {
	for k, v := range envVars {
		err := os.Setenv(k, v)
		assert.Nil(t, err)
	}
}

func unmarshalConfig(serializedConfig []byte) (envoy.Config, error) {
	c := envoy.Config{}
	err := yaml.Unmarshal(serializedConfig, &c)
	return c, err
}
