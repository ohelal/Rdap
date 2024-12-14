package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage RDAP query cache",
	Long: `Manage the RDAP query cache. Available subcommands:
  - stats: Show cache statistics
  - clear: Clear the cache`,
}

var cacheStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show cache statistics",
	Run: func(cmd *cobra.Command, args []string) {
		stats := getCacheStats()
		switch outputStyle {
		case "json":
			fmt.Println(formatJSON(stats, true))
		case "table":
			headers := []string{"Field", "Value"}
			rows := [][]string{
				{"Total Items", fmt.Sprintf("%d", stats.TotalItems)},
			}
			for key, value := range stats.Items {
				rows = append(rows, []string{key, fmt.Sprintf("%v", value)})
			}
			renderTable(headers, rows)
		default:
			fmt.Printf("Cache Statistics:\n")
			fmt.Printf("Total Items: %d\n", stats.TotalItems)
			for key, value := range stats.Items {
				fmt.Printf("%s: %v\n", key, value)
			}
		}
	},
}

var cacheClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear the cache",
	Run: func(cmd *cobra.Command, args []string) {
		clearCache()
	},
}

func init() {
	rootCmd.AddCommand(cacheCmd)
	cacheCmd.AddCommand(cacheStatsCmd)
	cacheCmd.AddCommand(cacheClearCmd)
}
