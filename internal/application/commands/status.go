package commands

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	core "github.com/UtkarshM-hub/Lit/internal/application/core"

	util "github.com/UtkarshM-hub/Lit/internal/application/core/util"
	color "github.com/gookit/color"
	"github.com/spf13/cobra"
)

type FileInfo struct {
	FileName       string
	FilePath       string
	FileSize       uint64
	FilePerm       uint32
	FileModifiedAt time.Time
	FileStatus     string
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "shows status of working directory and staging area",
	Long:  `shows status of working directory and staging area`,
	Run: func(cmd *cobra.Command, args []string) {
		// find .lit directory
		dir, err := util.FindDirectory(".lit")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		indexFilePath := filepath.Join(dir, "./.lit/index")

		files := core.GetFilesStatus(dir)

		tracked, untracked, modified, deleted, err := core.GetStatus(files, indexFilePath)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		statusMp := map[string]string{"M": "modified", "D": "deleted", "N": "new file"}

		// print tracked files
		fmt.Println("On branch <branch name>")
		fmt.Println("Your branch is up to date with '<origin>/<branch name>'.")
		fmt.Printf("\nChanges to be committed:\n")
		fmt.Printf("  (use 'lit restore --staged <file>...' to unstage)\n")
		for _, v := range tracked {
			fileP:=strings.Replace(v.FilePath,dir+"/","",-1)
			color.Green.Printf("\t%v:    %v\n", statusMp[v.FileStatus], fileP)
		}

		// modified and deleted
		if len(modified) > 0 {
			fmt.Printf("\nChanges not staged for commit:")
			fmt.Printf("\n   (use 'lit add/rm <file>...' to update what will be committed)")
			fmt.Printf("\n   (use 'lit restore <file>...' to discard changes in working directory)")
			for _, v := range modified {
				fileP:=strings.Replace(v.FilePath,dir+"/","",-1)
				color.Red.Printf("\n\t%v:    %v", statusMp[v.FileStatus], fileP)
			}
		}

		if len(deleted)>0{
			for _, v := range deleted {
				fileP:=strings.Replace(v.FilePath,dir+"/","",-1)
				color.Red.Printf("\n\t%v:    %v", statusMp[v.FileStatus], fileP)
			}
		}

		// Untracked files
		if len(untracked) > 0 {
			fmt.Printf("\nUntracked files:")
			fmt.Printf("  (use 'lit add <file>...' to include in what will be committed)")
			for _, v := range untracked {
				fileP:=strings.Replace(v.FilePath,dir+"/","",-1)
				color.Red.Printf("\n\t%v:    %v", statusMp[v.FileStatus], fileP)
			}
		}

		fmt.Println()
	},
}
