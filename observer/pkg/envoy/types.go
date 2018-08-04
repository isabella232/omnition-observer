package envoy

const (
	_ = iota
	HTTP1
	HTTP2
	TCP
)

const (
	_ = iota
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
	ApplicationProtocols string `yaml:"application_protocols"`
}

type VirtualHostRouteMatch struct {
	Prefix string
}

type VirtualHostRouteCluster struct {
	Cluster string
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
	Name       string
	RequireTLS bool `yaml:"required_tls,omitempty"`
	Domains    []string
	Routes     []VirtualHostRoute
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
	Cluster           string
}

type Filter struct {
	Name   string
	Config FilterConfig
}

type DataSource struct {
	// add other ways to specify data like path to file
	InlineString string `yaml:"inline_string"`
}

type TLSCertificate struct {
	CertificateChain DataSource `yaml:"certificate_chain"`
	PrivateKey       DataSource `yaml:"private_key"`
}

type CommonTLSContext struct {
	TLSCertificates []TLSCertificate `yaml:"tls_certificates"`
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

type TLSValidationContext struct {
	TrustedCA DataSource `yaml:"trusted_ca"`
}

type UpstreamCommonTLSContext struct {
	ValidationContext TLSValidationContext `yaml:"validation_context"`
}

type UpstreamTLSContext struct {
	CommonTLSContext UpstreamCommonTLSContext `yaml:"common_tls_context"`
}

type Cluster struct {
	Name                 string
	ConnectTimeout       string `yaml:"connect_timeout"`
	Type                 string
	LBPolicy             string               `yaml:"lb_policy"`
	Features             string               `yaml:",omitempty"`
	HTTP2ProtocolOptions HTTP2ProtocolOptions `yaml:"http2_protocol_options,omitempty"`
	Hosts                []ClusterHost
	TLSContext           UpstreamTLSContext `yaml:"tls_context,omitempty"`
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
