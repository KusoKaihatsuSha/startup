package order

type Stages int

const (
	// FLAG Get data from flags
	FLAG Stages = iota + 1

	// FILE Get data from json file
	FILE

	// ENV Get data from environments
	ENV

	// PreloadConfigEnvThenFlag - Get filepath config from environments and then flags
	// Default
	PreloadConfigEnvThenFlag

	// PreloadConfigFlagThenEnv - Get filepath config from flags and then environments
	PreloadConfigFlagThenEnv

	// PreloadConfigFlag - Get filepath config from flags
	PreloadConfigFlag

	// PreloadConfigEnv - Get filepath config from environments
	PreloadConfigEnv

	// NoPreloadConfig - Get filepath config file only from ordered stages
	NoPreloadConfig
)
