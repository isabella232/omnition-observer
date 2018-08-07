package envoy

import (
	"github.com/omnition/omnition-observer/observer/pkg/options"
)

func newFilterChain(
	direction TrafficDirection,
	protocol Protocol,
	httpsRedirect bool,
	opts options.Options,
) FilterChain {

	drName := "ingress"
	if direction == EGRESS {
		drName = "egress"
	}

	protoLabel := ""
	protoCode := ""
	switch protocol {
	case TCP:
		return FilterChain{
			Filters: []Filter{
				Filter{
					Name: "envoy.tcp_proxy",
					Config: FilterConfig{
						StatPrefix: drName + "_tcp",
						Cluster:    "tcp_" + drName + "_cluster",
					},
				},
			},
		}
	case HTTP1:
		protoLabel = "h1"
		protoCode = "http/1.1"
	case HTTP2:
		protoLabel = "h2"
		protoCode = "http/2.0"
	}

	label := protoLabel + "_" + drName

	chain := FilterChain{
		FilterChainMatch: FilterChainMatch{
			ApplicationProtocols: protoCode,
		},

		Filters: []Filter{
			Filter{
				Name: "envoy.http_connection_manager",
				Config: FilterConfig{
					StatPrefix:        label,
					CodecType:         "AUTO",
					GenerateRequestID: true,
					UseRemoteAddress:  true,
					Tracing: FilterConfigTracing{
						OperationName: drName,
					},
					RouteConfig: RouteConfig{
						Name: label + "_route",
						VirtualHosts: []VirtualHost{
							VirtualHost{
								Name:    label + "_vhost",
								Domains: []string{"*"},
								Routes: []VirtualHostRoute{
									VirtualHostRoute{
										Match: VirtualHostRouteMatch{"/"},
									},
								},
							},
						},
					},
					HTTPFilters: []HTTPFilter{
						HTTPFilter{Name: "envoy.grpc_http1_bridge"},
						HTTPFilter{Name: "envoy.router"},
					},
				},
			},
		},
	}

	if opts.TLSEnabled {
		// Setup TLS certificates
		if direction == INGRESS && !httpsRedirect {
			chain.TLSContext = TLSContext{
				CommonTLSContext{
					ALPNProtocols: protoCode,
					TLSCertificates: []TLSCertificate{
						TLSCertificate{
							CertificateChain: DataSource{
								InlineString: opts.TLSCert,
							},
							PrivateKey: DataSource{
								InlineString: opts.TLSKey,
							},
						},
					},
				},
			}
		}

		if httpsRedirect {
			// Setup HTTP > HTTPS redirect
			chain.Filters[0].Config.RouteConfig.VirtualHosts[0].Routes[0].Redirect = VirtualHostRouteRedirect{
				PathRedirect:  "/",
				HTTPSRedirect: true,
			}
		} else {
			// Setup actual route
			if direction == INGRESS {
				chain.FilterChainMatch.TransportProtocol = "tls"
			}
			chain.Filters[0].Config.RouteConfig.VirtualHosts[0].Routes[0].Route = VirtualHostRouteCluster{
				Cluster: protoLabel + "_" + drName + "_cluster",
			}
		}
	} else {
		// TLS not configured. Always handle route as is
		chain.Filters[0].Config.RouteConfig.VirtualHosts[0].Routes[0].Route = VirtualHostRouteCluster{
			Cluster: protoLabel + "_" + drName + "_cluster",
		}
	}

	return chain
}

func newListener(direction TrafficDirection, opts options.Options) Listener {
	port := opts.IngressPort
	name := "ingress_listener"
	if direction == EGRESS {
		port = opts.EgressPort
		name = "egress_listener"
	}

	listener := Listener{
		Name: name,
		Address: Address{
			SocketAddress{
				Address:   "0.0.0.0",
				PortValue: port,
			},
		},
		Transparent: true,
		ListenerFilters: []ListenerFilter{
			ListenerFilter{"envoy.listener.original_dst"},
			ListenerFilter{"envoy.listener.tls_inspector"},
			ListenerFilter{"envoy.listener.text_protocol_inspector"},
		},
	}

	chains := []FilterChain{}

	// HTTP1 Chain
	chains = append(chains, newFilterChain(direction, HTTP1, false, opts))
	// HTTP1 Chain
	chains = append(chains, newFilterChain(direction, HTTP2, false, opts))

	if direction == INGRESS && opts.TLSEnabled {
		// HTTP > HTTPS redirect for incoming traffic
		chains = append(chains, newFilterChain(direction, HTTP1, true, opts))
		chains = append(chains, newFilterChain(direction, HTTP2, true, opts))
	}

	listener.FilterChains = append(chains, newFilterChain(direction, TCP, false, opts))

	return listener
}

func newCluster(direction TrafficDirection, protocol Protocol, opts options.Options) Cluster {
	drName := "ingress"
	if direction == EGRESS {
		drName = "egress"
	}

	protoCode := ""
	protoLabel := ""
	features := ""
	switch protocol {
	case HTTP1:
		protoCode = "http/1.1"
		protoLabel = "h1"
	case HTTP2:
		protoCode = "http/2.0"
		protoLabel = "h2"
		features = "http2"
	case TCP:
		protoLabel = "tcp"
	}

	c := Cluster{
		Name:           protoLabel + "_" + drName + "_cluster",
		ConnectTimeout: "0.5s",
		Type:           "ORIGINAL_DST",
		LBPolicy:       "ORIGINAL_DST_LB",
		Features:       features,
	}
	if protocol == HTTP2 {
		c.HTTP2ProtocolOptions = HTTP2ProtocolOptions{
			MaxConcurrentStreams: 2147483647,
		}
	}

	if direction == EGRESS && opts.TLSEnabled && opts.TLSCACert != "" {
		c.TLSContext = TLSContext{
			CommonTLSContext{
				ALPNProtocols: protoCode,
				ValidationContext: ValidationContext{
					DataSource{
						InlineString: opts.TLSCACert,
						//FileName: "/etc/ssl/certs/ca-certificates.crt",
					},
				},
			},
		}
	}
	return c
}

func newTracingCluster(opts options.Options) Cluster {
	// TODO(owais): Add support for jaeger native tracing
	return Cluster{
		Name:           "tracing_zipkin_cluster",
		ConnectTimeout: "1s",
		Type:           "strict_dns",
		LBPolicy:       "round_robin",
		Hosts: []ClusterHost{
			ClusterHost{
				SocketAddress{
					Address:   opts.TracingHost,
					PortValue: opts.TracingPort,
				},
			},
		},
	}
}

func New(opts options.Options) (Config, error) {
	cfg := Config{
		Admin{
			opts.AdminLogPath,
			Address{
				SocketAddress{
					"0.0.0.0",
					opts.AdminPort,
				},
			},
		},
		StaticResources{
			Listeners: []Listener{
				newListener(INGRESS, opts),
				newListener(EGRESS, opts),
			},
			Clusters: []Cluster{
				newCluster(INGRESS, HTTP1, opts),
				newCluster(EGRESS, HTTP1, opts),
				newCluster(INGRESS, HTTP2, opts),
				newCluster(EGRESS, HTTP2, opts),
				newCluster(INGRESS, TCP, opts),
				newCluster(EGRESS, TCP, opts),
				newTracingCluster(opts),
			},
		},
		Tracing{
			TracingHTTP{
				"envoy.zipkin",
				TracingHTTPConfig{
					"tracing_zipkin_cluster",
					"/api/v1/spans",
				},
			},
		},
	}
	return cfg, nil
}
