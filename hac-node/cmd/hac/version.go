package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	GitCommit string
)

const (
	VersionMajor = 0
	VersionMinor = 0
	VersionPatch = 1
)

var Version = func() string {
	return fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
}()

func VersionWithCommit(gitCommit string) string {
	vsn := Version
	if len(gitCommit) >= 8 {
		vsn += "-" + gitCommit[:8]
	}
	return vsn
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "",
	Long:    ``,
	Aliases: []string{"V"},
	Run:     versionRun,
}

func versionRun(cmd *cobra.Command, args []string) {
	println(VersionWithCommit(GitCommit))
}
