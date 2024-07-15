package commands

import (
	"fmt"
	"path"

	"github.com/UtkarshM-hub/bit/internal/application/core"
	util "github.com/UtkarshM-hub/bit/internal/application/core/util"
	"github.com/spf13/cobra"
)

 var removeCached bool=false

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.Flags().BoolVar(&removeCached, "cached", false, "Used to assign message to commit")
}

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "removes files from staging area",
	Long:  `removes files from staging area`,
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

		if removeCached {
			var Removables []string
			if args[0] == "." {
				// removing files from staging area

			} else {
				// removing specific files from staging area
				for _, v := range args {
					newPath := path.Join(dir, v)
					Removables = append(Removables, newPath)
				}
			}

			core.RemoveFilesFromStagingArea(dir, Removables)
		} else {
			fmt.Println("Functionlity is not added yet")
			fmt.Println("you can use --cached flag to remove cached")
		}

		
	},
}
