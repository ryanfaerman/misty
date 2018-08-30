package cmd

import (
	"context"

	"github.com/ryanfaerman/misty/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(displayCmd)
}

var displayCmd = &cobra.Command{
	Use:   "change-led",
	Short: "change the color of the led",
	RunE: func(cmd *cobra.Command, args []string) error {

		client := client.New("http://10.0.1.5/", nil)
		err := client.Display.ChangeLED(context.TODO(), 255, 0, 0)
		log.Info(err)

		return nil

	},
}
