package commands

import (
	"fmt"

	"github.com/UtkarshM-hub/Lit/internal/application/core"
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
		if commitMessage==""{
			fmt.Println("Commit message not found")
			return
		}
		err:=core.Commit(commitMessage)
		if err!=nil{
			fmt.Println(err.Error())
			return
		}
	},
}
