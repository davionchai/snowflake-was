package operator

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/davionchai/snowflake-was/internal/arguments"
	"github.com/davionchai/snowflake-was/pkg/databases"
	"github.com/davionchai/snowflake-was/pkg/utils"
	"golang.org/x/exp/constraints"

	sf "github.com/snowflakedb/gosnowflake"
)

type ShowWarehouseResult struct {
	Name            string
	State           string
	Type            string
	Size            string
	Running         int
	Queued          int
	IsDefault       string
	IsCurrent       string
	AutoSuspend     int
	AutoResume      string
	Available       string
	Provisioning    string
	Quiescing       string
	Other           string
	CreatedOn       time.Time
	ResumedOn       time.Time
	UpdatedOn       time.Time
	Owner           string
	Comment         string
	ResourceMonitor string
	Actives         int
	Pendings        int
	Failed          int
	Suspended       int
	Uuid            string
}

func OperateWAS(sfConfig *sf.Config, appConfig *arguments.WasArguments) {
	queueCheckpoint := appConfig.DefaultQueueCheckpoint
	for {
		conn, err := databases.NewSnowflakeConn(sfConfig)
		if err != nil {
			log.Fatalf("failed to create DSN from Config: %v, err: %v", sfConfig, err)
		}
		defer conn.DB.Close()

		// show warehouses
		showQuery := fmt.Sprintf("show warehouses like '%s';", strings.ToLower(appConfig.SnowflakeWarehouseAutoscale))
		showRows, err := conn.RunQuery(showQuery)
		if err != nil {
			log.Fatal(err)
		}
		defer showRows.Close()

		var showRowResults []ShowWarehouseResult

		for showRows.Next() {
			var rowResult ShowWarehouseResult
			if err := showRows.Scan(
				&rowResult.Name,
				&rowResult.State,
				&rowResult.Type,
				&rowResult.Size,
				&rowResult.Running,
				&rowResult.Queued,
				&rowResult.IsDefault,
				&rowResult.IsCurrent,
				&rowResult.AutoSuspend,
				&rowResult.AutoResume,
				&rowResult.Available,
				&rowResult.Provisioning,
				&rowResult.Quiescing,
				&rowResult.Other,
				&rowResult.CreatedOn,
				&rowResult.ResumedOn,
				&rowResult.UpdatedOn,
				&rowResult.Owner,
				&rowResult.Comment,
				&rowResult.ResourceMonitor,
				&rowResult.Actives,
				&rowResult.Pendings,
				&rowResult.Failed,
				&rowResult.Suspended,
				&rowResult.Uuid,
			); err != nil {
				log.Fatalf("Error at scanning result %v", err)
			}
			showRowResults = append(showRowResults, rowResult)
		}
		if err := showRows.Err(); err != nil {
			log.Fatalf("Error at scanning result %v", err)
		}
		if len(showRowResults) > 1 {
			log.Fatalf(
				"Only 1 warehouse is allowed. Found [%v] warehouses.", len(showRowResults),
			)
		}
		showRowResult := &showRowResults[0]
		warehouseCenter, err := utils.NewWarehouseCenter(showRowResult.Size, appConfig.MinSize, appConfig.MaxSize)
		if err != nil {
			log.Fatal(err)
		}

		// escalation point algo
		// can be enhanced to progressively add n to queueCheckpoint
		// 	given queuedThreshold grows greater than queuedMultiplier
		//  can use queuedThreshold/queuedMultiplier and get multiplyFactor
		if showRowResult.Queued >= appConfig.QueuedThreshold {
			// example queued = 15, threshold = 5, then exponent = (3-1) = 2
			// 	if baseNumber = 2, then will be 2^2 = 4
			// smallest exponent is always 0, so min increment is always 1
			// this will ensure upsize fast if queue is high
			queuedExponent := (showRowResult.Queued / appConfig.QueuedThreshold) - 1
			queuedPowerOfNumber := intPow(appConfig.QueuedBaseNumber, queuedExponent)

			if queueCheckpoint > appConfig.DefaultQueueCheckpoint {
				queueCheckpoint = min(queueCheckpoint+queuedPowerOfNumber, 15)
			} else {
				queueCheckpoint = appConfig.DefaultQueueCheckpoint + queuedPowerOfNumber
			}
		} else if showRowResult.Queued < appConfig.QueuedThreshold {
			if queueCheckpoint < appConfig.DefaultQueueCheckpoint {
				queueCheckpoint = max(queueCheckpoint-1, 0)
			} else {
				queueCheckpoint = appConfig.DefaultQueueCheckpoint - 1
			}
		}
		log.Printf("checkpoint hit %v", queueCheckpoint)

		sizingEventTriggered := false
		var sizingEvent string
		if queueCheckpoint == 15 {
			sizingEventTriggered = warehouseCenter.Upsize()
			sizingEvent = "upsizing"
			queueCheckpoint = appConfig.DefaultQueueCheckpoint
		} else if queueCheckpoint == 0 {
			sizingEventTriggered = warehouseCenter.Downsize()
			sizingEvent = "downsizing"
			queueCheckpoint = appConfig.DefaultQueueCheckpoint
		}

		if sizingEventTriggered {
			log.Printf("%s warehouse to [%s]", sizingEvent, warehouseCenter.Size)
			alterQuery := fmt.Sprintf(
				"alter warehouse %s set warehouse_size = %s;",
				appConfig.SnowflakeWarehouseAutoscale,
				warehouseCenter.Size,
			)
			_, err := conn.DB.Exec(alterQuery)
			if err != nil {
				log.Fatal(err)
			}
		}

		time.Sleep(time.Duration(appConfig.CycleSeconds) * time.Second)
	}
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func intPow(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}
