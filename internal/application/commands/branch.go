package commands

import (
	"fmt"
	"path/filepath"

	"github.com/UtkarshM-hub/bit/internal/application/core"
	util "github.com/UtkarshM-hub/bit/internal/application/core/util"
	"github.com/spf13/cobra"
)

var list_branches bool

func init() {
	rootCmd.AddCommand(branchCmd)
	branchCmd.Flags().BoolVarP(&list_branches, "list", "a", false, "Used to list branches with current active branch")
}

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "creates branch so separate the different workflows of the project",
	Long:  `creates branch so separate the different workflows of the project`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && !list_branches {
			fmt.Println("Arguments not provided")
			return
		}

		dir, err := util.FindDirectory(".bit")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if list_branches {
			core.ListBranches(dir)
			return
		}

		branchName := args[0]

		// indexFilePath := filepath.Join(dir, "./.bit/index")

		// check if branch already exist or not
		branch_refFile_Path := filepath.Join(dir, "/.bit/refs/heads/"+branchName)

		err = util.DoesExists(branch_refFile_Path)
		if err == nil {
			fmt.Println("Branch already exists")
			return
		}

		err = core.CreateBranch(dir, branchName)
		if err != nil {
			fmt.Println("Error occured while creating branch")
		}
	},
}
