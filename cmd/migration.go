package cmd

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"log"
)

func newMigrateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate database",
		Long:  `migrate database`,
	}

	var sourceURL, databaseURL string

	cmd.PersistentFlags().
		StringVar(&sourceURL, "source-url", "file://database/migrations", "location of migration files")
	cmd.PersistentFlags().
		StringVar(&databaseURL, "database-url", "postgres://postgres:admin@localhost:5432/postgres?sslmode=disable&search_path=public", "connection string to database")

	cmd.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "Migrate database up",
		Long:  `run migrate database up using *.up.sql files`,
		Run: func(cmd *cobra.Command, args []string) {
			migrateUp(sourceURL, databaseURL)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "down",
		Short: "Migrate database down",
		Long:  `run migrate database up using *.down.sql files`,
		Run: func(cmd *cobra.Command, args []string) {
			migrateDown(sourceURL, databaseURL)
		},
	})

	return cmd
}

func migrateUp(sourceURL, databaseURL string) {
	m := newMigrate(sourceURL, databaseURL)

	log.Println("Migrating up...")
	err := m.Up()
	if err != nil {
		m.GracefulStop <- true
		log.Fatal(err)
	}
}

func migrateDown(sourceURL, databaseURL string) {
	m := newMigrate(sourceURL, databaseURL)

	log.Println("Migrating down...")
	err := m.Down()
	if err != nil {
		m.GracefulStop <- true
		log.Fatal(err)
	}
}

func newMigrate(sourceURL, databaseURL string) *migrate.Migrate {
	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	m.Log = &migrationLogger{}

	return m
}

type migrationLogger struct {

}

// Printf is like fmt.Printf
func (l *migrationLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Verbose should return true when verbose logging output is wanted
func (l *migrationLogger) Verbose() bool {
	return true
}