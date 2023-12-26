package commands

import (
	"fmt"
	"time"

	util "github.com/UtkarshM-hub/Lit/internal/application/core/util"
	core "github.com/UtkarshM-hub/Lit/internal/application/core"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes current directory as lit directory",
	Long:  `Initializes current directory as lit directory`,
	Run: func(cmd *cobra.Command, args []string) {

		t := time.Now()
		var dir string
		var err error

		// Get director path
		if len(args) == 0 || args[0] == "." {

			// get path of the current directory
			dir, err = util.GetCurrentDirectory()
			if err != nil {
				fmt.Println("Error finding path for current directory")
				return
			}
		} else {

			// join the directory path with the provided directory path
			dir, err = util.GetJoinedPaths(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		// check if it exists
		err = util.DirectoryExists(dir)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		done:=core.Init(dir)
		fmt.Println(dir,done, time.Since(t))
	},
}
