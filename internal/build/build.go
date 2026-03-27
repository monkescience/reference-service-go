package build

// Version is set at build time via -ldflags "-X reference-service-go/internal/build.Version=...".
var Version string

//nolint:gochecknoinits // init is used to enforce version is set at build time.
func init() {
	if Version == "" {
		panic("version must be set at build time using -ldflags")
	}
}
