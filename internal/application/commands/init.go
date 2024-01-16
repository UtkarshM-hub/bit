package commands

import (
	"fmt"
	"path/filepath"
	"time"

	core "github.com/UtkarshM-hub/Lit/internal/application/core"
	util "github.com/UtkarshM-hub/Lit/internal/application/core/util"
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
		err = util.DoesExists(dir)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// check if .lit directory already exists
		litpath := filepath.Join(dir, ".lit")
		err = util.DoesExists(litpath)
		if err == nil {
			fmt.Println("The directory is already lit 🔥")
			return
		}

		// initialize that directory as lit directory
		done := core.Init(dir)
		if done {
			fmt.Println(args[0], "initialized as 🔥 directory")
		}
		fmt.Println(dir, done, time.Since(t))
	},
}
