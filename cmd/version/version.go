package version

import (
	"fmt"
	"github.com/continuul/go-archetype/lib/version"
	"runtime"

	"github.com/spf13/cobra"
)

// Cmd represents the version command
var Cmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("go-archetype version %s %s/%s", version.GetHumanVersion(), runtime.GOOS, runtime.GOARCH))
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
