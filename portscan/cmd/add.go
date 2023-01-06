/*
Copyright © 2023 youngzy
Copyrights apply to this source code.
Check LICENSE for details.

*/
package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"io"
	"os"
	"portscan/scan"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:          "add <host 1>...<host n>",
	Short:        "Add new host(s) to list",
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		//hostFile, err := cmd.Flags().GetString("hosts-file")
		//if err != nil {
		//	return err
		//}
		hostFile := viper.GetString("hosts-file")
		return addAction(os.Stdout, hostFile, args)
	},
}

func addAction(out io.Writer, hostFile string, args []string) error {
	hl := &scan.HostsList{}
	if err := hl.Load(hostFile); err != nil {
		return err
	}
	for _, h := range args {
		if err := hl.Add(h); err != nil {
			return err
		}
		fmt.Fprintln(out, "Added host:", h)
	}
	return hl.Save(hostFile)
}

func init() {
	hostsCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
