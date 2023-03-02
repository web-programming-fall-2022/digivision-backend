package cmd

import "github.com/spf13/cobra"

func New() *cobra.Command {
	root := &cobra.Command{
		Use:   "dvs <subcommand>",
		Short: "dvs Daemon",
		Long:  `dvs is a gRPC microservice. More info at https://github.com/web-programming-fall-2022/digivision-backend`,
		Run:   nil,
	}
	addServeCmd(root)
	return root
}
