package commands

import (
	"fmt"
	"path/filepath"

	"github.com/UtkarshM-hub/bit/internal/application/core"
	util "github.com/UtkarshM-hub/bit/internal/application/core/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "adds files to staging area",
	Long:  `adds files from working directory to staging area`,
	Run: func(cmd *cobra.Command, args []string) {

		// Currently expecting "." with the command
		// The feature of adding single file into staging area will be added soon
		if len(args) == 0 {
			fmt.Println("Arguments not provided")
			return
		}

		// find path of directory with .bit file
		dir, err := util.FindDirectory(".bit")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		indexFilePath := filepath.Join(dir, "./.bit/index")

		if args[0] == "." {
			AllFiles := core.GetFilesStatus(dir)
			_, untracked, modified, deleted, err := core.GetStatus(AllFiles, indexFilePath)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			// adding files to staging area
			core.CoreAdd(dir, untracked, modified, deleted)
		} else {
			// logic of adding single files should go here
			fmt.Println(args)
		}
	},
}
