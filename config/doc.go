// Package config provide base config struct like ProjectConfig, DatabaseConfig, HTTPConfig, GRPCConfig.
//
// LoadConfigs called at the very beginning gave the definition of the project configuration.
// You can use env in your structures.
// To get started, you need to create your configuration structure that will contain the necessary structures
// from this package ProjectConfig must be defined at least.
//
// DatabaseConfig powerful struct that can create ds strings for your database.
//
// HTTPConfig and GRPCConfig define settings for transport.
package config
