package commands

import (
	"fmt"
	"log"
	"path"

	"github.com/UtkarshM-hub/bit/internal/application/core"
	util "github.com/UtkarshM-hub/bit/internal/application/core/util"
	"github.com/spf13/cobra"
)

var staged bool

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().BoolVar(&staged, "staged", false, "Used to specify type of restoration")
}

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "restore the file according to previous commit",
	Long:  `restore the file to its previous state according to the previous commit object`,
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

		var Restorable []string
		if args[0] == "." {
			// removing files from staging area
			fmt.Println("Please specify the filename")
			return
		} else {
			// removing specific files from staging area
			for _, v := range args {
				newPath := path.Join(dir, v)
				Restorable = append(Restorable, newPath)
			}
		}

		// restore the files either from staged or from previous commit
		err = core.Restore(Restorable, dir, staged)
		if err != nil {
			log.Fatal(err)
		}
	},
}
