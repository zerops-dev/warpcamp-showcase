package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/l-hellmann/sigterm"
	"github.com/spf13/cobra"
	"github.com/zerops-dev/warpcamp-showcase/src/http"
	"github.com/zerops-dev/warpcamp-showcase/src/migrate"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}),
		),
	)

	cmd := &cobra.Command{
		Use:           "app",
		SilenceErrors: true,
	}

	cmd.AddCommand(
		migrate.Command(),
		http.Command(),
	)

	if err := cmd.ExecuteContext(sigterm.Context()); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
