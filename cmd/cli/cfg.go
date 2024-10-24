package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(import_cfg_cmd)

}

var import_cfg_cmd = &cobra.Command{
	Use:   "import",
	Short: "imports acl.json into db",
	Long:  "imports acl into db.",
	Run: func(cmd *cobra.Command, args []string) {
		app, _ := getKritis3mScaleApp()

		app.Import()
		// app.Import_policies()
	},
}
