package list

import (
	"fmt"
	"github.com/continuul/go-archetype/lib/archetype"
	"github.com/logrusorgru/aurora"

	"github.com/spf13/cobra"
)

// Cmd represents the list command
var Cmd = &cobra.Command{
	Use:   "list",
	Short: "List the available archetypes",
	Run: func(cmd *cobra.Command, args []string) {
		archetypes, err := archetype.AssetDir("")
		if err != nil {
			fmt.Errorf("ERROR: failed to list available archetypes: %s", err)
		}
		for _, archetype := range archetypes {
			fmt.Println(aurora.Green(fmt.Sprintf(" - %s", archetype)))
		}
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
