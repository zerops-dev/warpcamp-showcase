package control

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/zerops-dev/warpcamp-showcase/src/database"
)

func EmptyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "empty",
		Short: "empties database",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := database.NewConnection(cmd.Context())
			if err != nil {
				return err
			}

			result, err := db.ExecContext(cmd.Context(), "DELETE FROM messages;")
			if err != nil {
				return err
			}

			affected, err := result.RowsAffected()
			if err != nil {
				return err
			}

			slog.With("rows", affected).Info("database cleared")

			return nil
		},
	}

	return cmd
}
