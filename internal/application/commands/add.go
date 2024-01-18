package commands

import (
	"fmt"
	"path/filepath"

	"github.com/UtkarshM-hub/Lit/internal/application/core"
	util "github.com/UtkarshM-hub/Lit/internal/application/core/util"
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
		if len(args) == 0 {
			fmt.Println("Arguments not provided")
			return
		}

		dir, err := util.FindDirectory(".lit")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		indexFilePath := filepath.Join(dir, "./.lit/index")

		if args[0] == "." {
			AllFiles := core.GetFilesStatus(dir)
			_,untracked,modified,deleted,err := core.GetStatus(AllFiles,indexFilePath)
			if err!=nil{
				fmt.Println(err.Error())
				return
			}
			core.CoreAdd(dir,untracked,modified,deleted)
		} else {
			fmt.Println(args)
		}
	},
}
