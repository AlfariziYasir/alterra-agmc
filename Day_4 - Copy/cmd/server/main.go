package main

import (
	"api-mvc/database/migration"
	"api-mvc/internal/http"
	"api-mvc/pkg/logger"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Go Blog API"
	app.Description = "Implementing back-end services for blog application"

	app.Commands = []cli.Command{
		{
			Name:        "migrations",
			Description: "migrations looks at the currently active migration version and will migrate all the way up (applying all up migrations)",
			Action: func(c *cli.Context) error {
				return migration.Up()
			},
		},
		{
			Name:        "drop",
			Description: "drop deletes everything in the database",
			Action: func(c *cli.Context) error {
				return migration.Drop()
			},
		},
		{
			Name:        "status",
			Description: "cek status",
			Action: func(c *cli.Context) error {
				return migration.Drop()
			},
		},
		{
			Name:        "start",
			Description: "start the server",
			Action: func(c *cli.Context) error {
				return http.Start()
			},
		},
		{
			Name:        "launch",
			Description: "launch migrate all the way up (applying all up migrations) and start the server",
			Action: func(c *cli.Context) error {
				err := migration.Up()
				if err != nil {
					return err
				}
				return http.Start()
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Log().Fatal().Err(err).Msg("failed to run server")
	}
}
