package envoy

import "github.com/omnition/omnition-observer/observer/pkg/options"

func newFilterChain(direction int, protocol int, opts options.Options) FilterChain {
	drName := "ingress"
	if direction == EGRESS {
		drName = "egress"
	}

	protoName := ""
	switch protocol {
	case HTTP1:
		protoName = "http1"
	case HTTP2:
		protoName = "http2"
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
	}

	chain := FilterChain{
		Filters: []Filter{
			Filter{
				Name: "envoy.http_connection_manager",
				Config: FilterConfig{
					StatPrefix:        "omnition",
					CodecType:         "AUTO",
					GenerateRequestID: true,
					UseRemoteAddress:  true,
					Tracing: FilterConfigTracing{
						OperationName: drName,
					},
					RouteConfig: RouteConfig{
						Name: protoName + "_route",
						VirtualHosts: []VirtualHost{
							VirtualHost{
								Name:       protoName + "_vhost",
								RequireTLS: true,
								Domains:    []string{"*"},
								Routes: []VirtualHostRoute{
									VirtualHostRoute{
										Match: VirtualHostRouteMatch{"/"},
										// redirect if TLS
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

	if protocol == HTTP1 {
		chain.FilterChainMatch = FilterChainMatch{"http/1.1"}
	} else if protocol == HTTP2 {
		chain.FilterChainMatch = FilterChainMatch{"http/2"}
	}

	if opts.TLSEnabled {
		//chain.Filters[0].Config.RouteConfig.VirtualHosts[0].Routes[0].Redirect = VirtualHostRouteRedirect{
		//	PathRedirect:  "/",
		//	HTTPSRedirect: true,
		//}
		chain.TLSContext = TLSContext{
			CommonTLSContext{
				[]TLSCertificate{
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
		//} else {
		//	chain.Filters[0].Config.RouteConfig.VirtualHosts[0].Routes[0].Route = VirtualHostRouteCluster{
		//		Cluster: protoName + "_" + drName + "_cluster",
		//	}
	}

	chain.Filters[0].Config.RouteConfig.VirtualHosts[0].Routes[0].Route = VirtualHostRouteCluster{
		Cluster: protoName + "_" + drName + "_cluster",
	}
	return chain
}

func newListener(direction int, opts options.Options) Listener {
	port := opts.IngressPort
	name := "omnition_ingress_listener"
	if direction == EGRESS {
		port = opts.EgressPort
		name = "omnition_egress_listener"
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
			ListenerFilter{"envoy.listener.text_protocol_inspector"},
			ListenerFilter{"envoy.listener.tls_inspector"},
		},
		FilterChains: []FilterChain{
			newFilterChain(direction, HTTP1, opts),
			newFilterChain(direction, HTTP2, opts),
			newFilterChain(direction, TCP, opts),
		},
	}
	return listener
}

func newCluster(direction int, protocol int) Cluster {
	drName := "ingress"
	if direction == EGRESS {
		drName = "egress"
	}

	protoName := ""
	features := ""
	switch protocol {
	case HTTP1:
		protoName = "http1"
	case HTTP2:
		protoName = "http2"
		features = "http2"
	case TCP:
		protoName = "tcp"
	}

	c := Cluster{
		Name:           protoName + "_" + drName + "_cluster",
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
	return c
}

func newTracingCluster(opts options.Options) Cluster {
	// TODO(owais): Add support for jaeger native tracing
	c := Cluster{
		Name:           "zipkin_cluster",
		ConnectTimeout: "1s",
		Type:           "strict_dns",
		LBPolicy:       "round_robin",
		Hosts: []ClusterHost{
			ClusterHost{
				SocketAddress{
					Address:   opts.TracingAddress,
					PortValue: opts.TracingPort,
				},
			},
		},
	}
	if opts.TLSCACert != "" {
		c.TLSContext = UpstreamTLSContext{
			UpstreamCommonTLSContext{
				TLSValidationContext{
					DataSource{
						InlineString: opts.TLSCACert,
					},
				},
			},
		}
	}
	return c
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
				newCluster(INGRESS, HTTP1),
				newCluster(EGRESS, HTTP1),

				newCluster(INGRESS, HTTP2),
				newCluster(EGRESS, HTTP2),

				newCluster(INGRESS, TCP),
				newCluster(EGRESS, TCP),

				newTracingCluster(opts),
			},
		},
		Tracing{
			TracingHTTP{
				"envoy.zipkin",
				TracingHTTPConfig{
					"zipkin_cluster",
					"/api/v1/spans",
				},
			},
		},
	}
	return cfg, nil
}
