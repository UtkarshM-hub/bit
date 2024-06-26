package commands

import (
	"fmt"
	"path/filepath"

	"github.com/UtkarshM-hub/bit/internal/application/core"
	util "github.com/UtkarshM-hub/bit/internal/application/core/util"
	color "github.com/gookit/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkoutCmd)
}

var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "Used to switch between two branches",
	Long:  `Used to switch between two different workflows aka branches`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Arguments not provided")
			return
		}

		branchName := args[0]

		dir, err := util.FindDirectory(".bit")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// indexFilePath := filepath.Join(dir, "./.bit/index")

		// check if branch already exist or not
		branch_refFile_Path := filepath.Join(dir, "/.bit/refs/heads/"+branchName)

		err = util.DoesExists(branch_refFile_Path)
		if err != nil {
			fmt.Printf("Branch with name '%v' doesn't exist\n", branchName)
			return
		}

		// prevent user from switching branch if current branch containers changes to be stages
		indexFilePath := filepath.Join(dir, "./.bit/index")

		files := core.GetFilesStatus(dir)

		tracked, untracked, modified, deleted, err := core.GetStatus(files, indexFilePath)

		if len(tracked) != 0 || len(untracked) != 0 || len(modified) != 0 || len(deleted) != 0 {
			fmt.Println("On branch <branch_name>")
			color.Red.Println("Commit changes before switching to another branch")
			return
		}

		core.Checkout(dir, branchName)
	},
}
