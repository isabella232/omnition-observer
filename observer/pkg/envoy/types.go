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

func (direction TrafficDirection) String() string {
	switch direction {
	case INGRESS:
		return "INBOUND"
	case EGRESS:
		return "OUTBOUND"
	default:
		return ""
	}
}

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

type VirtualHostRouteMatch struct {
	Prefix string
}

type VirtualHostRouteCluster struct {
	Cluster string
	Timeout *time.Duration `yaml:"timeout,omitempty"`
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
	Config struct{} `yaml:"typed_config"`
}

type FilterConfigTracing struct {
	OverallSampling Value    `yaml:"overall_sampling,omitempty"`
	CustomTags      []string `yaml:"custom_tags,omitempty"`
}

type FilterConfig struct {
	ConfigType        string              `yaml:"@type"`
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
	Name        string
	TypedConfig FilterConfig `yaml:"typed_config,omitempty"`
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
	TLSContext       *TLSContext `yaml:"tls_context,omitempty"`
}

type ListenerFilter struct {
	Name string
}

type Listener struct {
	Name            string
	Direction       string `yaml:"traffic_direction"`
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
	DnsLookupFamily      string               `yaml:"dns_lookup_family,omitempty"`
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

type Config struct {
	Admin           Admin
	StaticResources StaticResources `yaml:"static_resources"`
	Tracing         Tracing
}

type Value struct {
	Value float32 `yaml:"value"`
}
