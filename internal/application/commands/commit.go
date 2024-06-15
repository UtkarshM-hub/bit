package commands

import (
	"fmt"
	"path/filepath"

	"github.com/UtkarshM-hub/Lit/internal/application/core"
	"github.com/UtkarshM-hub/Lit/internal/application/core/util"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var commitMessage string

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "Used to assign message to commit")
}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Creates a commit out of current staging area",
	Long:  `Creates a commit out of current staging area`,
	Run: func(cmd *cobra.Command, args []string) {
		if commitMessage == "" {
			fmt.Println("Commit message not found")
			return
		}

		dir, err := util.FindDirectory(".bit")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		indexFilePath := filepath.Join(dir, "./.bit/index")

		files := core.GetFilesStatus(dir)

		tracked, _, modified, deleted, err := core.GetStatus(files, indexFilePath)

		if len(tracked) == 0 && len(modified) == 0 && len(deleted) == 0 {
			fmt.Println("On branch <branch_name>")
			color.Red.Println("Nothing to commit, working tree clean")
			return
		}

		err = core.Commit(commitMessage)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	},
}
