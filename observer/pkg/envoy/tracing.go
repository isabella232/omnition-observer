package envoy

// Tracing system identifiers
const (
	ZIPKIN = "zipkin"
	JEAGER = "jeager"
)

type TracingZipkinConfig struct {
	ConfigType               string `yaml:"@type"`
	CollectorCluster         string `yaml:"collector_cluster"`
	CollectorEndpoint        string `yaml:"collector_endpoint"`
	CollectorEndpointVersion string `yaml:"collector_endpoint_version,omitempty"`
}

// ref: https://github.com/jaegertracing/jaeger-client-cpp
type TracingJeagerConfig struct {
	ConfigType   string       `yaml:"@type"`
	Library      string       `yaml:"library"`
	JeagerConfig JeagerConfig `yaml:"config"`
}

type JeagerConfig struct {
	ServiceName string               `yaml:"service_name"`
	Sampler     JeagerConfigSampler  `yaml:"sampler"`
	Reporter    JeagerConfigReporter `yaml:"reporter"`
	Tags        string               `yaml:"tags,omitempty"`
}

type JeagerConfigSampler struct {
	SamplerType string  `yaml:"type"`
	Param       float32 `yaml:"param"`
}

type JeagerConfigReporter struct {
	CollectorEndpoint string `yaml:"endpoint"`
}

type TracingHTTP struct {
	Name   string
	Config interface{} `yaml:"typed_config"`
}

type Tracing struct {
	Http TracingHTTP
}
