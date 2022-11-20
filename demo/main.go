package main

import (
	"IM/demo/client"
	"IM/demo/server"
	"context"
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	flag.Parse()
	root := &cobra.Command{
		Use:     "chat",
		Version: "v1",
		Short:   "chat demo",
	}
	ctx := context.Background()

	root.AddCommand(server.NewServerStartCmd(ctx, "v1"))
	root.AddCommand(client.NewCmd(ctx))

	if err := root.Execute(); err != nil {
		logrus.WithError(err).Fatal("Could not run command")
	}
}
