package envoy

import "time"

type Protocol int

const (
	_ = iota
	HTTP1
	HTTP2
	TCP
)

type TrafficDirection int

const (
	_ TrafficDirection = iota
	EGRESS
	INGRESS
)

type SocketAddress struct {
	Address   string
	PortValue int `yaml:"port_value"`
}

type Address struct {
	SocketAddress SocketAddress `yaml:"socket_address"`
}

type Admin struct {
	AccessLogPath string `yaml:"access_log_path"`
	Address       Address
}

type FilterChainMatch struct {
	ApplicationProtocols string `yaml:"application_protocols,omitempty"`
	TransportProtocol    string `yaml:"transport_protocol,omitempty"`
}

type RetryPriority struct{}
type RetryHostPredicate struct{}

type RetryPolicy struct {
	RetryOn                       string             `yaml:"retry_on"`
	NumRetries                    uint32             `yaml:"num_retries,omitempty"`
	PerRetryTimeout               *time.Duration     `yaml:"per_try_timeout,omitempty"`
	RetryPriority                 RetryPriority      `yaml:"retry_priority,omitempty"`
	RetryHostPredicate            RetryHostPredicate `yaml:"retry_host_predicate,omitempty"`
	HostSelectionRetryMaxAttempts int64              `yaml:"host_selection_retry_max_attempts,omitempty"`
	RetriableStatusCodes          []uint32           `yaml:"retriable_status_code,omitempty"`
}

type VirtualHostRouteMatch struct {
	Prefix string
}

type VirtualHostRouteCluster struct {
	Cluster     string
	RetryPolicy RetryPolicy
}
type VirtualHostRouteRedirect struct {
	PathRedirect  string `yaml:"path_redirect"`
	HTTPSRedirect bool   `yaml:"https_redirect"`
}

type VirtualHostRoute struct {
	Match    VirtualHostRouteMatch
	Route    VirtualHostRouteCluster  `yaml:",omitempty"`
	Redirect VirtualHostRouteRedirect `yaml:",omitempty"`
}

type VirtualHost struct {
	Name    string
	Domains []string
	Routes  []VirtualHostRoute
}

type RouteConfig struct {
	Name         string
	VirtualHosts []VirtualHost `yaml:"virtual_hosts"`
}

type HTTPFilter struct {
	Name   string
	Config struct{}
}

type FilterConfigTracing struct {
	OperationName string `yaml:"operation_name"`
}

type FilterConfig struct {
	StatPrefix        string              `yaml:"stat_prefix"`
	CodecType         string              `yaml:"codec_type,omitempty"`
	GenerateRequestID bool                `yaml:"generate_request_id,omitempty"`
	UseRemoteAddress  bool                `yaml:"use_remote_address,omitempty"`
	Tracing           FilterConfigTracing `yaml:",omitempty"`
	RouteConfig       RouteConfig         `yaml:"route_config,omitempty"`
	HTTPFilters       []HTTPFilter        `yaml:"http_filters,omitempty"`
	Cluster           string              `yaml:"cluster,omitempty"`
}

type Filter struct {
	Name   string
	Config FilterConfig
}

type DataSource struct {
	InlineString string `yaml:"inline_string,omitempty"`
	FileName     string `yaml:"filename,omitempty"`
}

type TLSCertificate struct {
	CertificateChain DataSource `yaml:"certificate_chain"`
	PrivateKey       DataSource `yaml:"private_key"`
}

type ValidationContext struct {
	TrustedCA DataSource `yaml:"trusted_ca"`
}

type CommonTLSContext struct {
	ALPNProtocols     string            `yaml:"alpn_protocols,omitempty"`
	TLSCertificates   []TLSCertificate  `yaml:"tls_certificates,omitempty"`
	ValidationContext ValidationContext `yaml:"validation_context,omitempty"`
}

type TLSContext struct {
	CommonTLSContext CommonTLSContext `yaml:"common_tls_context"`
}

type FilterChain struct {
	FilterChainMatch FilterChainMatch `yaml:"filter_chain_match,omitempty"`
	Filters          []Filter
	TLSContext       TLSContext `yaml:"tls_context,omitempty"`
}

type ListenerFilter struct {
	Name string
}

type Listener struct {
	Name            string
	Address         Address
	Transparent     bool
	ListenerFilters []ListenerFilter `yaml:"listener_filters"`
	FilterChains    []FilterChain    `yaml:"filter_chains"`
}

type HTTP2ProtocolOptions struct {
	MaxConcurrentStreams int `yaml:"max_concurrent_streams"`
}

type Cluster struct {
	Name                 string
	ConnectTimeout       string `yaml:"connect_timeout"`
	Type                 string
	LBPolicy             string               `yaml:"lb_policy"`
	HTTP2ProtocolOptions HTTP2ProtocolOptions `yaml:"http2_protocol_options,omitempty"`
	TLSContext           TLSContext           `yaml:"tls_context,omitempty"`
	Hosts                []ClusterHost        `yaml:"hosts,omitempty"`
}

type ClusterHost struct {
	SocketAddress SocketAddress `yaml:"socket_address"`
}

type StaticResources struct {
	Listeners []Listener
	Clusters  []Cluster
}

type TracingHTTPConfig struct {
	CollectorCluster  string `yaml:"collector_cluster"`
	CollectorEndpoint string `yaml:"collector_endpoint"`
}

type TracingHTTP struct {
	Name   string
	Config TracingHTTPConfig
}
type Tracing struct {
	HTTP TracingHTTP
}

type Config struct {
	Admin           Admin
	StaticResources StaticResources `yaml:"static_resources"`
	Tracing         Tracing
}
