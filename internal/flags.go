package internal

var (
	Debug      bool   // Debug describes if the debug lvl is enabled or not
	DebugAvail bool   // DebugAvail describes if the binary was compiled with debug features enabled.
	Version    string // Version contains the current version string and will be set on compile time
	Commit     string // Commit contains the latest git commit hash
	BuildDate  string // BuildDate contains the date (UTC) when the binary was built
)
