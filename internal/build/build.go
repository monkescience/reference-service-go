package build

import "runtime/debug"

// ServiceName is the canonical service identifier used for telemetry resources.
const ServiceName = "reference-service-go"

const unknownBuildValue = "(unknown)"

// Build metadata set via -ldflags in release builds.
var (
	version string
	commit  string
	date    string
)

// Version returns the injected version or falls back to build info.
func Version() string {
	if version != "" {
		return version
	}

	return buildInfoVersion()
}

// Commit returns the injected commit or falls back to build info.
func Commit() string {
	if commit != "" {
		return commit
	}

	return buildInfoSetting("vcs.revision")
}

// Date returns the injected build date or falls back to build info.
func Date() string {
	if date != "" {
		return date
	}

	return buildInfoSetting("vcs.time")
}

func buildInfoVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" {
		return unknownBuildValue
	}

	return info.Main.Version
}

func buildInfoSetting(key string) string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return unknownBuildValue
	}

	for _, setting := range info.Settings {
		if setting.Key == key && setting.Value != "" {
			return setting.Value
		}
	}

	return unknownBuildValue
}
