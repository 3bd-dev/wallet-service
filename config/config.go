package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

const (
	envPrefix              = ""
	DefaultContextDeadline = 15 * time.Second
)

// HTTPServer contains the configuration for the HTTP server.
type HTTPServer struct {
	// Port is the HTTP server port.
	Port int `envconfig:"HTTP_SERVER_PORT" default:"8080"`

	// IdleTimeout is the maximum amount of time to wait for an open connection
	// when processing no requests and keep-alives are enabled. If this value is
	// 0, ReadTimeout value be used.
	IdleTimeout time.Duration `envconfig:"HTTP_SERVER_IDLE_TIMEOUT" default:"60s"`

	// ReadTimeout is the maximum duration for reading the entire request,
	// including the body.
	ReadTimeout time.Duration `envconfig:"HTTP_SERVER_READ_TIMEOUT" default:"1s"`

	// WriteTimeout is the maximum duration before timing out
	// writes of the response.
	WriteTimeout time.Duration `envconfig:"HTTP_SERVER_WRITE_TIMEOUT" default:"2s"`

	// ShutdownTimeout is the maximum duration before timing out
	ShutdownTimeout time.Duration `envconfig:"HTTP_SERVER_SHUTDOWN_TIMEOUT" default:"5s"`
}

// Database contains configuration for the Postgres Database.
type Database struct {
	URL          string `envconfig:"DATABASE_URL" required:"true"` // DO NOT ADD  prefix for this env. this dbmate is depend on this env without prefix
	MaxOpenConns int    `envconfig:"DATABASE_MAX_OPEN_CONNECTIONS" default:"10"`
	MaxIdleConns int    `envconfig:"DATABASE_MAX_IDLE_CONNECTIONS" default:"5"`
}

type Service struct {
	Version     string `envconfig:"SERVICE_VERSION" default:""`
	Environment string `envconfig:"SERVICE_ENVIRONMENT" required:"true"`
}

type PaymentGatewayA struct {
	BaseURL                  string        `envconfig:"GATEWAY_A_API_BASE_URL"`
	RetryAttempt             int           `envconfig:"GATEWAY_A_RETRY_ATTEMPT" default:"3"`   // Number of retry attempts for failed requests
	RetryDelay               time.Duration `envconfig:"GATEWAY_A_RETRY_DELAY" default:"1s"`    // Delay between retries
	CBMaxRequests            uint32        `envconfig:"GATEWAY_A_CB_MAX_REQUESTS" default:"5"` // Max requests allowed in half-open state
	CBInterval               time.Duration `envconfig:"GATEWAY_A_CB_INTERVAL" default:"60s"`   // Interval to reset the failure counter
	CBTimeout                time.Duration `envconfig:"GATEWAY_A_CB_TIMEOUT" default:"30s"`    // Time to stay open before testing recovery
	CBMaxConsecutiveFailures uint32        `envconfig:"GATEWAY_A_CB_MAX_CONSECUTIVE_FAILURES" default:"3"`
	CBMaxTotalFailures       uint32        `envconfig:"GATEWAY_A_CB_MAX_TOTAL_FAILURES" default:"5"`
}

type PaymentGatewayB struct {
	BaseURL                  string        `envconfig:"GATEWAY_B_API_BASE_URL"`
	RetryAttempt             int           `envconfig:"GATEWAY_B_RETRY_ATTEMPT" default:"3"`
	RetryDelay               time.Duration `envconfig:"GATEWAY_B_RETRY_DELAY" default:"1s"`
	MaxRequests              uint32        `envconfig:"GATEWAY_B_CB_MAX_REQUESTS" default:"5"` // Max requests allowed in half-open state
	Interval                 time.Duration `envconfig:"GATEWAY_B_CB_INTERVAL" default:"60s"`   // Interval to reset the failure counter
	Timeout                  time.Duration `envconfig:"GATEWAY_B_CB_TIMEOUT" default:"30s"`    // Time to stay open before testing recovery
	CBMaxConsecutiveFailures uint32        `envconfig:"GATEWAY_B_CB_MAX_CONSECUTIVE_FAILURES" default:"3"`
	CBMaxTotalFailures       uint32        `envconfig:"GATEWAY_B_CB_MAX_TOTAL_FAILURES" default:"5"`
}
type PaymentGatewayConfig struct {
	// GatewayA configuration.
	GatewayA PaymentGatewayA

	// GatewayB configuration.
	GatewayB        PaymentGatewayB
	CallbackPattern string `envconfig:"PAYMENT_CALLBACK_PATTERN" required:"true"`
}

// Config holds all configuration in a struct to make the transition to the
// new config structs easier.
type Config struct {
	Service              Service
	Server               HTTPServer
	Database             Database
	PaymentGatewayConfig PaymentGatewayConfig
}

// Load parses configuration from environment.
func Load() (Config, error) {
	cfg := Config{}
	err := envconfig.Process(envPrefix, &cfg)
	return cfg, err
}
