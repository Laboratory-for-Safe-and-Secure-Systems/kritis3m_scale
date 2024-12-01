package cli

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(nodeCmd)
	// Add listNodesCmd as a subcommand of nodeCmd
	listNodesCmd.Flags().IntP("id", "i", -1, "search node with id")
	listNodesCmd.Flags().Bool("with_cfg", false, "include configs")
	nodeCmd.AddCommand(listNodesCmd) // Attach listNodesCmd to nodeCmd

	nodeCmd.AddCommand(activateConfigCmd) // Attach listNodesCmd to nodeCmd
	activateConfigCmd.Flags().Int("node_id", -1, "activate node with id <node_id>")
	activateConfigCmd.Flags().Int("cfg_id", -1, "activate configuration with config id <cfg_id>")
	activateConfigCmd.MarkFlagRequired("node_id")
	activateConfigCmd.MarkFlagRequired("cfg_id")

	nodeCmd.AddCommand(listActiveCmd) // Attach listNodesCmd to nodeCmd

	nodeCmd.AddCommand(configCmd)        // Attach listNodesCmd to nodeCmd
	configCmd.AddCommand(listconfigsCmd) // Attach listNodesCmd to nodeCmd

	listconfigsCmd.Flags().IntP("id", "i", -1, "search config with id")
	listconfigsCmd.Flags().Bool("with_appls", false, "include application ids")

}

var nodeCmd = &cobra.Command{
	Use:     "nodes",
	Short:   "Manage the nodes of Kritis3m-Scale",
	Aliases: []string{"node", "machine", "machines"},
}

var configCmd = &cobra.Command{
	Use:     "configs",
	Short:   "Manage the application configurations of Kritis3m-Scale",
	Aliases: []string{"config", "configurations"},
}

var listActiveCmd = &cobra.Command{
	Use:     "lsa",
	Short:   "lists active configurations ",
	Aliases: []string{"lsa", "list_active", "show_active"},
	Run: func(cmd *cobra.Command, args []string) {

		app, err := getKritis3mScaleApp()
		if err != nil {
			log.Err(err).Msg("error occured in list nodes")
			return
		}
		app.ListActive()
	},
}

var listNodesCmd = &cobra.Command{
	Use:     "list",
	Short:   "list nodes",
	Aliases: []string{"ls", "list", "show"},
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("activate configuration")

		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			log.Err(err).Msg("err getting id")
			return
		}
		includeConfig, _ := cmd.Flags().GetBool("with_cfg")
		if err != nil {

			log.Err(err).Msg("err getting cfg flag")
			return
		}

		app, err := getKritis3mScaleApp()
		if err != nil {
			log.Err(err).Msg("error occured in list nodes")
			return
		}
		app.ListNodes(id, includeConfig)
	},
}

var activateConfigCmd = &cobra.Command{
	Use:     "activate",
	Short:   "activate configuration",
	Aliases: []string{"activate", "select"},
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("in list Nodes")

		node_id, err := cmd.Flags().GetInt("node_id")
		if err != nil {
			log.Err(err).Msg("err getting id")
			return
		}
		cfg_id, err := cmd.Flags().GetInt("cfg_id")
		if err != nil {

			log.Err(err).Msg("err getting cfg flag")
			return
		}
		if (cfg_id < 0) ||
			(node_id < 0) {
			log.Error().Msg("cfg or node id must be a positive number")
			return
		}

		app, err := getKritis3mScaleApp()
		if err != nil {
			log.Err(err).Msg("error occured in list nodes")
			return
		}
		app.ActivateConfig(uint(node_id), uint(cfg_id))
	},
}

var listconfigsCmd = &cobra.Command{
	Use:     "list",
	Short:   "list configs",
	Aliases: []string{"ls", "list", "show"},
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("in list Nodes")

		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			log.Err(err).Msg("err getting id")
			return
		}

		includeAppls, err := cmd.Flags().GetBool("with_appls")
		if err != nil {

			log.Err(err).Msg("err getting  appl flag")
			return
		}

		app, err := getKritis3mScaleApp()
		if err != nil {
			log.Err(err).Msg("error occured in list nodes")
			return
		}
		app.Listconfigs(id, includeAppls)
	},
}

// var listConfigsCmd = &cobra.Command{
// 	Use:     "config",
// 	Short:   "Show configurations",
// 	Aliases: []string{"config", "configuration"},
// }
