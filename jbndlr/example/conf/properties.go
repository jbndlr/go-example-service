package conf

// Properties : Application settings as parsed from defauls and ENV.
type Properties struct {
	Service struct {
		Version string `yaml:"Version"`
		Name    string `yaml:"Name"`
	} `yaml:"Service"`
	API struct {
		ESIDPort                   int16 `yaml:"ESIDPort" envconfig:"API_ESID_PORT"`
		RESTPort                   int16 `yaml:"RESTPort" envconfig:"API_REST_PORT"`
		GRPCPort                   int16 `yaml:"GRPCPort" envconfig:"API_GRPC_PORT"`
		LimitMaxConcurrentRequests int   `yaml:"LimitMaxConcurrentRequests" envconfig:"API_LIMIT_MAX_CONCURRENT_REQUESTS"`
		LimitMaxWindowRequests     int   `yaml:"LimitMaxWindowRequests" envconfig:"API_LIMIT_MAX_WINDOW_REQUESTS"`
		LimitWindowSeconds         int   `yaml:"LimitWindowSeconds" envconfig:"API_LIMIT_WINDOW_SECONDS"`
		HTTPReadTimeout            int   `yaml:"HTTPReadTimeout" envconfig:"API_HTTP_READ_TIMEOUT"`
		HTTPWriteTimeout           int   `yaml:"HTTPWriteTimeout" envconfig:"API_HTTP_WRITE_TIMEOUT"`
		GracefulSeconds            int   `yaml:"GracefulSeconds" envconfig:"API_GRACEFUL_SECONDS"`
	} `yaml:"API"`
	Auth struct {
		JWTKey      string `yaml:"JWTKey" envconfig:"AUTH_JWT_KEY"`
		ExpireAfter int16  `yaml:"ExpireAfter" envconfig:"AUTH_EXPIRE_AFTER"`
		RenewWindow int16  `yaml:"RenewWindow" envconfig:"AUTH_RENEW_WINDOW"`
	} `yaml:"Auth"`
}
