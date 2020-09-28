package config

// The HTTPConfig config fort http transport.Server.
type HTTPConfig struct {
	Addr string `default:":8080"`
}

// The GRPCConfig config fort grpc transport.Server.
type GRPCConfig struct {
	Addr string `default:":8081"`
}
