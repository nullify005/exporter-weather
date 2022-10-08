/*
Copyright © 2022 Lee Webb <nullify005@gmail.com>
*/
package cmd

import (
	"time"

	"github.com/nullify005/exporter-weather/internal/watcher"
	"github.com/spf13/cobra"
)

// watchCmd represents the watch command
var (
	listen   string
	interval *time.Duration
	watchCmd = &cobra.Command{
		Use:   "watch geohash",
		Short: "Continuously poll for observations at location",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			w := watcher.New(args[0], watcher.WithInterval(*interval), watcher.WithListen(listen))
			w.Watch()
		},
	}
)

func init() {
	rootCmd.AddCommand(watchCmd)
	watchCmd.Flags().StringVarP(&listen, "listen", "l", "127.0.0.1:2112", "Listen on the following interface & port")
	interval = watchCmd.Flags().DurationP("interval", "i", 30*time.Second, "Polling interval for observations")
}
