/*
Copyright © 2023 Ryder Retzlaff <ryder@retzlaff.family>
*/
package cmd

import (
	"os"

	"github.com/pyscripter99/tget/internals"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tget -u url",
	Short: "A golang implementation of wget",
	Long:  `A golang implementation of wget.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get url & download file
		url := cmd.Flag("url").Value.String()
		if url == "" {
			cmd.Help()
			return
		}
		output := cmd.Flag("output").Value.String()
		err := internals.Download(url, output, cmd.Flag("user-agent").Value.String())
		if err != nil {
			internals.Fatal(err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("user-agent", "a", "", "set user agent")
	rootCmd.Flags().StringP("url", "u", "", "url to download")
	rootCmd.Flags().StringP("output", "o", "", "output file path")
}
