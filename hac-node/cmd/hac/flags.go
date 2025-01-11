package main

import "github.com/spf13/cobra"

func urlFlag(cmd *cobra.Command, url *string) {
	cmd.Flags().StringVarP(url, "url", "u", "http://127.0.0.1:26657", "hac-cl service url")
}
