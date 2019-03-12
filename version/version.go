package version

import (
	"runtime"
	"time"
)

var (
	build Build
)

type Build struct {
	Version   string
	GitSHA    string
	BuildTime time.Time
	GoVersion string
}

func init() {
	build.Version = version
	build.GitSHA = gitSHA
	build.BuildTime, _ = time.Parse(time.UnixDate, buildTime)
	build.GoVersion = runtime.Version()
}

func GetBuild() Build {
	return build
}

func Version() string {
	return build.Version
}

func GitSHA() string {
	return build.GitSHA
}

func BuildTime() time.Time {
	return build.BuildTime
}

func GoVersion() string {
	return build.GoVersion
}
