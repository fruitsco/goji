package goji

// Environment represents the environment of the application.
type Environment string

const (
	// EnvironmentDevelopment represents the development environment.
	EnvironmentDevelopment Environment = "development"

	// EnvironmentProduction represents the production environment.
	EnvironmentProduction Environment = "production"

	// EnvironmentTest represents the test environment.
	EnvironmentTest Environment = "test"
)
