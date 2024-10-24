package cli

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import_policies",
	Short: "imports the polices of acl.json and writes those into the database",
	Long:  "imports the polices of acl.json and writes those into the database",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("starting kritis3m_scale app")

		//first: start headscale app:
		app, _ := getKritis3mScaleApp()
		app.Import()

	},
}
