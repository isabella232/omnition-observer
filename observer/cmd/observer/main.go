package main

import (
	"fmt"
	"os"

	"github.com/omnition/omnition-observer/observer/pkg/envoy"
	"github.com/omnition/omnition-observer/observer/pkg/options"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	viper.SetEnvPrefix("OBS")

	viper.SetDefault("tls_enabled", false)
	viper.BindEnv("tls_enabled")

	viper.BindEnv("tls_ca_cert")
	viper.BindEnv("tls_cert")
	viper.BindEnv("tls_key")

	viper.SetDefault("ingress_port", 15001)
	viper.BindEnv("ingress_port")

	viper.SetDefault("egress_port", 15002)
	viper.BindEnv("egress_port")

	viper.SetDefault("admin_port", 9901)
	viper.BindEnv("admin_port")
	viper.SetDefault("admin_log_path", "/dev/null")
	viper.BindEnv("admin_log_path")

	viper.SetDefault("tracing_driver", "zipkin")
	viper.BindEnv("tracing_driver")

	viper.SetDefault("tracing_host", "zipkin.default.svc.cluster.local")
	viper.BindEnv("tracing_host")

	viper.SetDefault("tracing_port", 9411)
	viper.BindEnv("tracing_port")

	viper.SetDefault("tracing_tag_headers", []string{})
	viper.BindEnv("tracing_tag_headers")

	viper.SetDefault("timeout", "15s")
	viper.BindEnv("timeout")
}

func main() {
	serialized, err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(serialized))
}

func run() ([]byte, error) {
	opts, err := buildOptions()
	if err != nil {
		return nil, err
	}

	generated, err := generateConfig(&opts)
	if err != nil {
		return nil, err
	}

	serialized, err := yaml.Marshal(&generated)
	if err != nil {
		return nil, err
	}

	return serialized, nil
}

func buildOptions() (options.Options, error) {
	return options.New(
		viper.GetInt("ingress_port"),
		viper.GetInt("egress_port"),

		viper.GetString("tracing_driver"),
		viper.GetString("tracing_host"),
		viper.GetInt("tracing_port"),
		viper.GetStringSlice("tracing_tag_headers"),

		viper.GetBool("tls_enabled"),
		viper.GetString("tls_ca_cert"),
		viper.GetString("tls_cert"),
		viper.GetString("tls_key"),

		viper.GetInt("admin_port"),
		viper.GetString("admin_log_path"),

		viper.GetDuration("timeout"),
		viper.GetInt("num_trusted_hops"),
	)
}

func generateConfig(options *options.Options) (*envoy.Config, error) {
	return envoy.New(*options)
}
