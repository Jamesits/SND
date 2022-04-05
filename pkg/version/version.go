package version

import "fmt"

var version string
var versionGitCommitHash string
var versionCompileTime string
var versionCompileHost string

func GetVersionNumberString() string {
	return version
}

func GetVersionFullString() string {
	if len(versionCompileHost) == 0 {
		versionCompileHost = "localhost"
	}

	if len(versionGitCommitHash) == 0 {
		versionGitCommitHash = "UNKNOWN"
	}

	if len(versionCompileTime) == 0 {
		versionCompileTime = "UNKNOWN TIME"
	}

	return fmt.Sprintf("SND/%s (+https://github.com/Jamesits/SND; Compiled on %s for commit %s at %s)", GetVersionNumberString(), versionCompileHost, versionGitCommitHash, versionCompileTime)
}
