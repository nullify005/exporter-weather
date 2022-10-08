/*
Copyright © 2022 Lee Webb <nullify005@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/nullify005/exporter-weather/internal/search"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var (
	searchCmd = &cobra.Command{
		Use:                   "search location ...",
		Short:                 "Return a Geohash for a named location",
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			for _, v := range args {
				err = search.Search(v)
			}
			if err != nil {
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(searchCmd)
}
