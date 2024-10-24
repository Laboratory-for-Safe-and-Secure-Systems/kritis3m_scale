package cli

import (
	"github.com/spf13/cobra"
	"log"
)

func init() {
	rootCmd.AddCommand(startServer)
}

var startServer = &cobra.Command{
	Use:   "start",
	Short: "starts the node server",
	Long:  "starts the node server",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("starting kritis3m_scale app")

		//first: start headscale app:
		app, _ := getKritis3mScaleApp()
		app.Serve()

	},
}
