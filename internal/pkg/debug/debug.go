package debug

import (
	"fmt"
	"runtime/debug"
)

var info, ok = debug.ReadBuildInfo()

func init() {
	if ok {
		fmt.Println(info)
	}
}

func GetBuildInfo() (goVersion, commitHash string) {
	if !ok {
		return "", ""
	}

	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			commitHash = setting.Value
		}
	}

	return info.GoVersion, commitHash
}
