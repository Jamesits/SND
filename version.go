package snd

import "fmt"

var versionMajor = 0
var versionMinor = 2
var versionRevision = 0
var versionGitCommitHash string
var versionCompileTime string
var versionCompileHost string
var versionGitStatus string

func GetVersionNumberString() string {
	return fmt.Sprintf("%d.%d.%d", versionMajor, versionMinor, versionRevision)
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

	if len(versionGitStatus) == 0 {
		versionGitStatus = "dirty"
	}

	return fmt.Sprintf("SND/%s (+https://github.com/Jamesits/SND; Compiled on %s for commit %s (%s) at %s)", GetVersionNumberString(), versionCompileHost, versionGitCommitHash, versionGitStatus, versionCompileTime)
}
