package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// cfgFile config file variable
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bigbang",
	Short: "A brief description of your application",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", ".configs/config-prod.yaml", "config file")
}
