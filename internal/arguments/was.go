package arguments

// warehouse auto scaling app arguments
type WasArguments struct {
	SnowflakeWarehouseAutoscale string
	MinSize                     string
	MaxSize                     string

	QueuedThreshold        int
	QueuedBaseNumber       int
	DefaultQueueCheckpoint int
	CycleSeconds           int
}

// warehouse auto scaling flags string
type wasCLI struct {
	SnowflakeUsername      string
	SnowflakePassword      string
	SnowflakeAccount       string
	SnowflakeRole          string
	SnowflakeWarehouseRun  string
	SnowflakeAuthenticator string

	SnowflakeWarehouseAutoscale string
	MinSize                     string
	MaxSize                     string

	QueuedThreshold        string
	QueuedBaseNumber       string
	DefaultQueueCheckpoint string
	CycleSeconds           string
}

func NewWasFlags() *wasCLI {
	wasFlags := &wasCLI{
		SnowflakeUsername:      "sf-username",
		SnowflakePassword:      "sf-password",
		SnowflakeAccount:       "sf-account",
		SnowflakeRole:          "sf-role",
		SnowflakeWarehouseRun:  "sf-warehouse-run",
		SnowflakeAuthenticator: "sf-authenticator",

		SnowflakeWarehouseAutoscale: "sf-warehouse-autoscale",
		MinSize:                     "min-size",
		MaxSize:                     "max-size",

		QueuedThreshold:        "queued-threshold",
		QueuedBaseNumber:       "queued-base-number",
		DefaultQueueCheckpoint: "default-queue-checkpoint",
		CycleSeconds:           "cycle-seconds",
	}

	return wasFlags
}

func NewWasConfigKeys() *wasCLI {
	wasConfigKeys := &wasCLI{
		SnowflakeUsername:      "was.snowflakeUsername",
		SnowflakePassword:      "was.snowflakePassword",
		SnowflakeAccount:       "was.snowflakeAccount",
		SnowflakeRole:          "was.snowflakeRole",
		SnowflakeWarehouseRun:  "was.snowflakeWarehouseRun",
		SnowflakeAuthenticator: "was.snowflakeAuthenticator",

		SnowflakeWarehouseAutoscale: "was.snowflakeWarehouseAutoscale",
		MinSize:                     "was.minSize",
		MaxSize:                     "was.maxSize",

		QueuedThreshold:        "was.queuedThreshold",
		QueuedBaseNumber:       "was.queuedBaseNumber",
		DefaultQueueCheckpoint: "was.defaultQueueCheckpoint",
		CycleSeconds:           "was.cycleSeconds",
	}

	return wasConfigKeys
}
