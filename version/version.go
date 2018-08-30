package version

import (
	"fmt"
	"runtime"
)

var (
	ApplicationName = "unknown"
	CommitHash      = "unknown"
	BuildDate       = "unknown"
	BuildNumber     = "v0-dev"
)

type VersionInfo struct {
	ApplicationName string
	CommitHash      string
	BuildDate       string
	BuildTarget     string
	BuildNumber     string
	GoVersion       string
}

func (vi VersionInfo) String() string {
	return fmt.Sprintf("%s %s (%s) %s - BuildDate: %s", vi.ApplicationName, vi.BuildNumber, vi.CommitHash, vi.BuildTarget, vi.BuildDate)
}

var Version VersionInfo

func init() {
	Version = VersionInfo{
		ApplicationName: ApplicationName,
		CommitHash:      CommitHash,
		BuildNumber:     BuildNumber,
		BuildDate:       BuildDate,
		BuildTarget:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		GoVersion:       runtime.Version(),
	}
}
