package was

import (
	"log"
	"strconv"
	"strings"

	"github.com/davionchai/snowflake-was/internal/arguments"
	"github.com/davionchai/snowflake-was/internal/operator"
	"github.com/davionchai/snowflake-was/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	sf "github.com/snowflakedb/gosnowflake"
)

var (
	snowflakeUsername      string
	snowflakePassword      string
	snowflakeAccount       string
	snowflakeRole          string
	snowflakeWarehouseRun  string
	snowflakeAuthenticator string

	snowflakeWarehouseAutoscale string
	minSize                     string
	maxSize                     string

	queuedThreshold        int
	queuedBaseNumber       int
	defaultQueueCheckpoint int
	cycleSeconds           int
)

var wasFlags = arguments.NewWasFlags()
var wasConfigKeys = arguments.NewWasConfigKeys()

var WasCmd = &cobra.Command{
	Use:   "was",
	Short: "Auto scaling snowflake warehouse",
	Long: `snowflake-was is a Go written utility that helps to auto scale 
Snowflake warehouse size. It does so by monitoring warehouse activity
and responding to different escalation points to upsize or downsize the warehouse size`,
	Run: func(cmd *cobra.Command, args []string) {
		sfConfig, appConfig, defaultCheck := argumentResolver()
		if viper.ConfigFileUsed() == "" && defaultCheck == true {
			cmd.Help()
		} else {
			operator.OperateWAS(sfConfig, appConfig)
		}
	},
}

func init() {
	// read in environment variables that match
	viper.SetEnvPrefix("was")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	WasCmd.Flags().StringVarP(&snowflakeUsername, wasFlags.SnowflakeUsername, "u", "", "snowflake username")
	WasCmd.Flags().StringVarP(&snowflakePassword, wasFlags.SnowflakePassword, "p", "", "snowflake password, only allowed when authenticator is `snowflake`")
	WasCmd.Flags().StringVarP(&snowflakeAccount, wasFlags.SnowflakeAccount, "a", "", "snowflake account, ie xyz.us-east-1")
	WasCmd.Flags().StringVarP(&snowflakeRole, wasFlags.SnowflakeRole, "r", "", "snowflake role")
	WasCmd.Flags().StringVarP(&snowflakeWarehouseRun, wasFlags.SnowflakeWarehouseRun, "w", "", "snowflake warehouse to run queries")
	WasCmd.Flags().StringVarP(&snowflakeAuthenticator, wasFlags.SnowflakeAuthenticator, "e", "", "snowflake authenticator type")

	WasCmd.Flags().StringVarP(&snowflakeWarehouseAutoscale, wasFlags.SnowflakeWarehouseAutoscale, "s", "", "snowflake warehouse to be autoscaled")
	WasCmd.Flags().StringVarP(&minSize, wasFlags.MinSize, "i", "", "minimum warehouse size to be downsized")
	WasCmd.Flags().StringVarP(&maxSize, wasFlags.MaxSize, "x", "", "maximum warehouse size to be upsized")

	WasCmd.Flags().IntVarP(&queuedThreshold, wasFlags.QueuedThreshold, "t", -1, "min amount of queue to trigger upsize event")
	WasCmd.Flags().IntVarP(&queuedBaseNumber, wasFlags.QueuedBaseNumber, "b", -1, "base number for escalation point algo, affects how queudPowerOfNumber is calculated")
	WasCmd.Flags().IntVarP(&defaultQueueCheckpoint, wasFlags.DefaultQueueCheckpoint, "q", -1, "affects how fast a warehouse is upsized and downsized, maximum 15 and min 0, suggested value is 5")
	WasCmd.Flags().IntVarP(&cycleSeconds, wasFlags.CycleSeconds, "c", -1, "determine how often the app perform warehouse autoscaling activity")

	viper.BindPFlags(WasCmd.Flags())
	// viper.BindPFlag("sf_username", WasCmd.Flags().Lookup("sf-username"))

}

// order of argument resolver
//  1. read from yaml file (if exists)
//  2. read from env (overwrite yaml)
//  3. read from cli (overwrite yaml & env var)
//  4. all default are assumed empty for err catching through flags
func argumentResolver() (*sf.Config, *arguments.WasArguments, bool) {
	sfConfig, checkSfConfig := buildSfConfig()
	appConfig, checkAppConfig := buildAppConfig()

	if checkSfConfig == true || checkAppConfig == true {
		return nil, nil, true
	}

	return sfConfig, appConfig, false
}

func buildSfConfig() (*sf.Config, bool) {
	username := parseCommand(wasFlags.SnowflakeUsername, wasConfigKeys.SnowflakeUsername)
	password := parseCommand(wasFlags.SnowflakePassword, wasConfigKeys.SnowflakePassword)
	account := parseCommand(wasFlags.SnowflakeAccount, wasConfigKeys.SnowflakeAccount)
	role := parseCommand(wasFlags.SnowflakeRole, wasConfigKeys.SnowflakeRole)
	warehouse := parseCommand(wasFlags.SnowflakeWarehouseRun, wasConfigKeys.SnowflakeWarehouseRun)
	_authenticator := parseCommand(wasFlags.SnowflakeAuthenticator, wasConfigKeys.SnowflakeAuthenticator)

	// check if is default
	checkList := map[string]interface{}{
		"username":      username,
		"password":      password,
		"account":       account,
		"role":          role,
		"warehouse":     warehouse,
		"authenticator": _authenticator,
	}
	if checkDefaultValue(checkList) == true {
		return nil, true
	}

	authenticator, err := utils.GetSnowflakeAuthType(_authenticator)
	if err != nil {
		log.Fatal(err)
	}
	// design pattern of gosnowflake using pointer for session_parameters assignment
	query_tag := "snowflake-was"
	session_parameters := map[string]*string{"QUERY_TAG": &query_tag}

	// only AuthTypeSnowflake allow for password
	if authenticator == sf.AuthTypeSnowflake {
		if password == "" {
			log.Fatal("Password cannot be empty if authenticator is [snowflake]")
		}
		return &sf.Config{
			User:          username,
			Password:      password,
			Account:       account,
			Role:          role,
			Warehouse:     warehouse,
			Authenticator: authenticator,
			Params:        session_parameters,
		}, false
	} else {
		return &sf.Config{
			User:          username,
			Account:       account,
			Role:          role,
			Warehouse:     warehouse,
			Authenticator: authenticator,
			Params:        session_parameters,
		}, false
	}
}

func buildAppConfig() (*arguments.WasArguments, bool) {
	warehouseAutoscale := parseCommand(wasFlags.SnowflakeWarehouseAutoscale, wasConfigKeys.SnowflakeWarehouseAutoscale)
	minSize := parseCommand(wasFlags.MinSize, wasConfigKeys.MinSize)
	maxSize := parseCommand(wasFlags.MaxSize, wasConfigKeys.MaxSize)

	queuedThreshold, err := strconv.Atoi(
		parseCommand(wasFlags.QueuedThreshold, wasConfigKeys.QueuedThreshold),
	)
	if err != nil {
		log.Fatal(err)
	}
	queuedBaseNumber, err := strconv.Atoi(
		parseCommand(wasFlags.QueuedBaseNumber, wasConfigKeys.QueuedBaseNumber),
	)
	if err != nil && queuedBaseNumber != 0 {
		log.Fatal(err)
	}
	defaultQueueCheckpoint, err := strconv.Atoi(
		parseCommand(wasFlags.DefaultQueueCheckpoint, wasConfigKeys.DefaultQueueCheckpoint),
	)
	if err != nil {
		log.Fatal(err)
	}
	cycleSeconds, err := strconv.Atoi(
		parseCommand(wasFlags.CycleSeconds, wasConfigKeys.CycleSeconds),
	)
	if err != nil {
		log.Fatal(err)
	}

	// check if is default
	checkList := map[string]interface{}{
		"snowflakeWarehouseAutoscale": warehouseAutoscale,
		"minSize":                     minSize,
		"maxSize":                     maxSize,
		"queuedThreshold":             queuedThreshold,
		"queuedBaseNumber":            queuedBaseNumber,
		"defaultQueueCheckpoint":      defaultQueueCheckpoint,
		"cycleSeconds":                cycleSeconds,
	}
	if checkDefaultValue(checkList) == true {
		return nil, true
	}

	appConfig := &arguments.WasArguments{
		SnowflakeWarehouseAutoscale: warehouseAutoscale,
		MinSize:                     minSize,
		MaxSize:                     maxSize,

		QueuedThreshold:        queuedThreshold,
		QueuedBaseNumber:       queuedBaseNumber,
		DefaultQueueCheckpoint: defaultQueueCheckpoint,
		CycleSeconds:           cycleSeconds,
	}

	return appConfig, false
}

func parseCommand(wasFlag, wasConfigKey string) string {
	v := viper.GetViper()
	targetFlag := v.GetString(wasFlag)
	// if cli and env is not default (meaning assigned)
	if targetFlag != "" && targetFlag != "-1" {
		return targetFlag
	}
	// read from yaml
	taregtConfig := v.GetString(wasConfigKey)
	if taregtConfig != "" && taregtConfig != "-1" {
		return taregtConfig
	}
	// return default
	return targetFlag
}

// "false" means there is at least 1 provided arguments;
// "true" means given struct is at default (empty)
func checkDefaultValue(targetMap map[string]interface{}) bool {
	for _, value := range targetMap {
		switch val := value.(type) {
		case string:
			if val != "" {
				return false
			}
		case int:
			if val != -1 {
				return false
			}
		}
	}
	return true
}
