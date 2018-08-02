package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"

	"github.com/omnition/omnition-observer/observer/pkg/envoy"
	"github.com/omnition/omnition-observer/observer/pkg/options"
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

	viper.SetDefault("tracing_handler", "zipkin")
	viper.BindEnv("tracing_handler")

	viper.SetDefault("tracing_address", "zipkin.opsoss.svc.cluster.local")
	viper.BindEnv("tracing_address")

	viper.SetDefault("tracing_port", 9411)
	viper.BindEnv("tracing_port")
}

func main() {

	opts, err := options.New(
		viper.GetInt("ingress_port"),
		viper.GetInt("egress_port"),

		viper.GetString("tracing_handler"),
		viper.GetString("tracing_address"),
		viper.GetInt("tracing_port"),

		viper.GetBool("tls_enabled"),
		viper.GetString("tls_ca_cert"),
		viper.GetString("tls_cert"),
		viper.GetString("tls_key"),

		viper.GetInt("admin_port"),
		viper.GetString("admin_log_path"),
	)
	if err != nil {
		log.Fatal(err)
	}

	generated, err := envoy.New(opts)
	if err != nil {
		log.Fatal(err)
	}

	serialized, err := yaml.Marshal(generated)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(serialized))
}
