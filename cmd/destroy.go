// Copyright © 2017 Stream
//

package cmd

import (
	"fmt"
	"os"

	"github.com/webmastak/vg/internal/workspace"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Recursively change permissions
func ChmodR(path string, mode os.FileMode) error {
	return filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
		if err == nil {
			err = os.Chmod(name, mode)
		}

		return err
	})
}

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy [workspaces...]",
	Short: "Removes one or multiple workspace and all their contents",
	Long: `To remove workspace 'myWorkspace' and 'someOtherWorkspace':

	vg destroy myWorkspace someOtherWorkspace
	
To remove the currently active workspace you can call the command without
any arguments:

	vg destroy
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("No workspace specified")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, wsName := range args {
			fmt.Printf("Destroying workspace %q\n", wsName)
			ws := workspace.New(wsName)
			ChmodR(ws.Path(), 0755)
			err := ws.ClearSrc()
			if err != nil {
				return err
			}

			err = os.RemoveAll(workspace.New(wsName).Path())
			if err != nil {
				return errors.Wrapf(err, "Couldn't remove workspace %q", wsName)
			}
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(destroyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// destroyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// destroyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
