package migrate

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
	"github.com/zerops-dev/warpcamp-showcase/src/database"
)

//go:embed resources/*.sql
var fs embed.FS

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate (up|down|status)",
		Short: "performs database migration",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Usage()
			}
			db, err := database.NewConnection(cmd.Context())
			if err != nil {
				return err
			}
			defer db.Close()

			var direction migrate.MigrationDirection
			switch args[0] {
			case "up":
				direction = migrate.Up
			case "down":
				direction = migrate.Down
			case "status":
				return logAppliedMigrations(db)
			default:
				return fmt.Errorf("invalid argument '%s' use 'up' or 'down'", args[0])
			}

			migrations, err := migrate.Exec(
				db,
				"postgres",
				&migrate.EmbedFileSystemMigrationSource{
					FileSystem: fs,
					Root:       "resources",
				},
				direction,
			)
			if err != nil {
				return err
			}

			slog.Info(fmt.Sprintf("executed %d migrations", migrations))
			return logAppliedMigrations(db)
		},
	}

	return cmd
}

func logAppliedMigrations(db *sql.DB) error {
	records, err := migrate.GetMigrationRecords(db, "postgres")
	if err != nil {
		return err
	}
	logger := slog.Default()
	if len(records) == 0 {
		logger.Info("no migrations applied")
		return nil
	}
	for _, r := range records {
		logger.With("appliedAt", r.AppliedAt).Info(r.Id)
	}
	return nil
}
