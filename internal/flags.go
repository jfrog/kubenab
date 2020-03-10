package internal

// this variables will be set on compilation-time.
var (
	DebugAvail bool   // DebugAvail describes if the binary was compiled with debug features enabled.
	Version    string // Version contains the current version string and will be set on compile time
	Commit     string // Commit contains the latest git commit hash
	BuildDate  string // BuildDate contains the date (UTC) when the binary was built
)

// this variables are used for saving flag arguments
var (
	Debug      bool   // Debug describes if the debug lvl is enabled or not
	ListenAddr string // ListenAddr contains the address where it should listen
	ListenPort uint   // ListenPort contains the port where it should listen
)
