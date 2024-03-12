package commands

import (
	"fmt"
	"path/filepath"

	"github.com/UtkarshM-hub/Lit/internal/application/core"
	util "github.com/UtkarshM-hub/Lit/internal/application/core/util"
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

		dir, err := util.FindDirectory(".lit")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// indexFilePath := filepath.Join(dir, "./.lit/index")

		// check if branch already exist or not
		branch_refFile_Path := filepath.Join(dir, "/.lit/refs/heads/"+branchName)

		err = util.DoesExists(branch_refFile_Path)
		if err != nil {
			fmt.Printf("Branch with name '%v' doesn't exist\n", branchName)
			return
		}

		// change active branch
		err = core.ChangeActiveBranch(dir, branchName)
		if err != nil {
			fmt.Println(err)
		}

		core.DecompressFile(filepath.Join(dir, "/.lit/objects/3c/cd72b349cb144d17a6a66003344568ae929c6f"), filepath.Join(dir, "/name"))
	},
}
