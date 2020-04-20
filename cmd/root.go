package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/victornm/es-backend/api"
	"log"
	"net/http"
	"os"
	"strconv"
)

func newRootCommand() *cobra.Command {
	defaultConfig := struct {
		secret          string
		jwtExpiredHours int
		frontendBaseURL string
		apiBaseURL      string
		httpPort        int
		sqlUrl          string

		oauth2GoogleClientID     string
		oauth2GoogleClientSecret string
	}{
		secret:                   envString("SECRET", "z91NRBxicpx2qjvO"),
		jwtExpiredHours:          envInt("TOKEN_EXPIRED_HOURS", 24),
		frontendBaseURL:          envString("FRONTEND_BASE_URL", "http://localhost:3000"),
		apiBaseURL:               envString("API_BASE_URL", "http://localhost:8080"),
		httpPort:                 envInt("HTTP_PORT", 8080),
		sqlUrl:                   envString("SQL_URL", "postgres://postgres:admin@localhost:5432/postgres?sslmode=disable&search_path=public"),
		oauth2GoogleClientID:     envString("OAUTH2_GOOGLE_CLIENT_ID", ""),
		oauth2GoogleClientSecret: envString("OAUTH2_GOOGLE_CLIENT_SECRET", ""),
	}

	var (
		config   = &api.ServerConfig{}
		httpPort int
	)

	cmd := &cobra.Command{
		Use: "app",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Flags().
				StringVar(&config.JWTSecret, "secret", defaultConfig.secret, "secret key for JWT")
			cmd.Flags().
				IntVar(&config.JWTExpiredHours, "token-expired-hours", defaultConfig.jwtExpiredHours, "expired duration in hour for JWT token")
			cmd.Flags().
				StringVar(&config.FrontendBaseURL, "frontend-base-url", defaultConfig.frontendBaseURL, "")
			cmd.Flags().
				StringVar(&config.APIBaseURL, "api-base-url", defaultConfig.apiBaseURL, "")
			cmd.Flags().
				IntVar(&httpPort, "http-port", defaultConfig.httpPort, "port listening")
			cmd.Flags().
				StringVar(&config.SqlConnString, "sql-url", defaultConfig.sqlUrl, "connection string to database")
			cmd.Flags().
				StringVar(&config.OAuth2GoogleClientID, "oauth2-google-client-id", defaultConfig.oauth2GoogleClientID, "connection string to database")
			cmd.Flags().
				StringVar(&config.OAuth2GoogleClientSecret, "oauth2-google-client-secret", defaultConfig.oauth2GoogleClientSecret, "connection string to database")

			log.Printf("Init server with config: %+v\n", config)

			s := api.NewServer(config)
			s.Init()
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpPort), s))
		},
	}

	// add sub commands
	cmd.AddCommand(newMigrateCommand())

	return cmd
}

func Execute() {
	err := newRootCommand().Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func envString(key string, value string) string {
	env := os.Getenv(key)
	if len(env) == 0 {
		return value
	}

	return env
}

func envInt(key string, value int) int {
	env, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return value
	}

	return env
}
