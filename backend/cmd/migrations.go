package cmd

import (
	"fmt"

	"github.com/svartlfheim/gomigrator"
	"github.com/svartlfheim/mimisbrunnr/internal/config"
)

func handleMigrationsUp(m *gomigrator.Migrator, cfg *config.AppConfig, args []string) error {
	return m.Up(gomigrator.MigrateToLatest)
}

func handleMigrationsList(m *gomigrator.Migrator, cfg *config.AppConfig, args []string) error {
	ml, err := m.ListMigrations()

	if err != nil {
		return err
	}

	for _, m := range ml {
		fmt.Printf("Migration: %s (%s)\n", m.Id, m.Status)
	}

	return nil
}

func handleMigrationsDown(m *gomigrator.Migrator, cfg *config.AppConfig, args []string) error {
	return m.Down(gomigrator.MigrateToNothing)
}
