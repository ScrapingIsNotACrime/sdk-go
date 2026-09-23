package sinac

import (
	"runtime/debug"
	"strings"
)

const modulePath = "github.com/ScrapingIsNotACrime/sdk-go"

// Version is this module's version as recorded in the importing program's
// build info (for example "0.1.0"), or "dev" when it is unknown.
var Version = moduleVersion()

func moduleVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}
	for _, dep := range info.Deps {
		if dep.Path != modulePath {
			continue
		}
		version := dep.Version
		if dep.Replace != nil && dep.Replace.Version != "" {
			version = dep.Replace.Version
		}
		if version != "" && version != "(devel)" {
			return strings.TrimPrefix(version, "v")
		}
	}
	return "dev"
}

func userAgent() string {
	return "scrapingisnotacrime-go/" + Version
}
