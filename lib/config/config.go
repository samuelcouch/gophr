package config

import (
	"bytes"
	"os"
	"strconv"

	"gopkg.in/urfave/cli.v1"
)

const (
	environmentDev  = "dev"
	environmentProd = "prod"

	envVarsPort                 = "GOPHR_PORT, PORT"
	envVarsDepotPath            = "GOPHR_DEPOT_PATH"
	envVarsDbAddress            = "GOPHR_DB_ADDR"
	envVarsEnvironment          = "GOPHR_ENV"
	envVarsSecretsPath          = "GOPHR_SECRETS_PATH"
	envVarsMigrationsPath       = "GOPHR_MIGRATIONS_PATH"
	envVarsConstructionZonePath = "GOPHR_CONSTRUCTION_ZONE_PATH"
)

// Config contains vital environment metadata used through out the backend.
type Config struct {
	IsDev                bool
	Port                 int
	DepotPath            string
	DbAddress            string
	SecretsPath          string
	MigrationsPath       string
	ConstructionZonePath string
}

func (c *Config) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("Is dev:                 ")
	buffer.WriteString(strconv.FormatBool(c.IsDev))
	buffer.WriteString("\nPort:                   ")
	buffer.WriteString(strconv.Itoa(c.Port))

	if len(c.DepotPath) > 0 {
		buffer.WriteString("\nDepot path:             ")
		buffer.WriteString(c.DepotPath)
	}

	if len(c.SecretsPath) > 0 {
		buffer.WriteString("\nSecrets path:           ")
		buffer.WriteString(c.SecretsPath)
	}

	if len(c.DbAddress) > 0 {
		buffer.WriteString("\nDatabase address:       ")
		buffer.WriteString(c.DbAddress)
	}

	if len(c.MigrationsPath) > 0 {
		buffer.WriteString("\nMigrations path:        ")
		buffer.WriteString(c.MigrationsPath)
	}

	if len(c.ConstructionZonePath) > 0 {
		buffer.WriteString("\nConstruction zone path: ")
		buffer.WriteString(c.ConstructionZonePath)
	}

	return buffer.String()
}

// GetConfig gets the configuration for the current execution environment.
func GetConfig() *Config {
	var (
		port                 int
		depotPath            string
		dbAddress            string
		secretsPath          string
		environment          string
		migrationsPath       string
		constructionZonePath string

		app            = cli.NewApp()
		actionExecuted = false
	)

	// Make the cli for config less boring.
	app.Usage = "a component of the gophr backend"

	// Map config variables 1:1 with flags.
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "environment, e",
			Value:       environmentDev,
			Usage:       "execution context of this binary",
			EnvVar:      envVarsEnvironment,
			Destination: &environment,
		},
		cli.IntFlag{
			Name:        "port, p",
			Value:       3000,
			Usage:       "http port to exposed by this binary",
			EnvVar:      envVarsPort,
			Destination: &port,
		},
		cli.StringFlag{
			Name:        "depot-path",
			Usage:       "path to the depot repos",
			EnvVar:      envVarsDepotPath,
			Destination: &depotPath,
		},
		cli.StringFlag{
			Name:        "secrets-path, s",
			Usage:       "path to the secret files",
			EnvVar:      envVarsSecretsPath,
			Destination: &secretsPath,
		},
		cli.StringFlag{
			Name:        "db-address, d",
			Value:       "127.0.0.1",
			Usage:       "address of the database",
			EnvVar:      envVarsDbAddress,
			Destination: &dbAddress,
		},
		cli.StringFlag{
			Name:        "migrations-path, m",
			Usage:       "path to the db migration files",
			EnvVar:      envVarsMigrationsPath,
			Destination: &migrationsPath,
		},
		cli.StringFlag{
			Name:        "construction-zone-path, c",
			Usage:       "path to the construction zone",
			EnvVar:      envVarsConstructionZonePath,
			Destination: &constructionZonePath,
		},
	}

	// Use the action to figure out whether the environment variables are valid.
	app.Action = func(c *cli.Context) error {
		if environment != environmentDev && environment != environmentProd {
			return cli.NewExitError("invalid environment", 1)
		}

		actionExecuted = true
		return nil
	}

	// Execute the cli; wait to see what happens afterwards.
	app.Run(os.Args)

	// If there wasn't supposed to be an action, don't continue.
	if !actionExecuted {
		os.Exit(0)
	}

	return &Config{
		IsDev:                environment == environmentDev,
		Port:                 port,
		DepotPath:            depotPath,
		DbAddress:            dbAddress,
		SecretsPath:          secretsPath,
		MigrationsPath:       migrationsPath,
		ConstructionZonePath: constructionZonePath,
	}
}
