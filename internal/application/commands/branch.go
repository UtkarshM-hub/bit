package commands

import (
	"fmt"
	"path/filepath"

	"github.com/UtkarshM-hub/Lit/internal/application/core"
	util "github.com/UtkarshM-hub/Lit/internal/application/core/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(branchCmd)
}

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "creates branch so separate the different workflows of the project",
	Long:  `creates branch so separate the different workflows of the project`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Arguments not provided")
			return
		}

		branchName := args[0]

		dir, err := util.FindDirectory(".lit")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// indexFilePath := filepath.Join(dir, "./.lit/index")

		// check if branch already exist or not
		branch_refFile_Path := filepath.Join(dir, "/.lit/refs/heads/"+branchName)

		err = util.DoesExists(branch_refFile_Path)
		if err == nil {
			fmt.Println("The branch you want to create is already 🔥")
			return
		}

		err=core.CreateBranch(dir,branchName)	
		if err!=nil{
			fmt.Println("Error occured while creating branch")
		}
	},
}
