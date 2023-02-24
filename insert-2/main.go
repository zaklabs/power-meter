package main

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

const (
	myToken  = "ot3rcVMa4Y1V8S3otRMFItQiX0vChFbgeSOe2gOlv7RjYk5NpMnUaJBOND9U9JvILr50tCus6v36tL4pJ2SBQQ=="
	myOrg    = "5261a190ceb20da3"
	myBucket = "as"
)

func main() {
	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClient("http://10.81.81.36:8099", myToken)

	// Use blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking(myOrg, myBucket)

	//  1.  Create point using full params constructor
	p := influxdb2.NewPoint("stat",
		map[string]string{"unit": "temperature"},
		map[string]interface{}{"avg": 24.5, "max": 45.0},
		time.Now())
	// write point immediately
	writeAPI.WritePoint(context.Background(), p)

	// 2.  Create point using fluent style
	p = influxdb2.NewPointWithMeasurement("stat").
		AddTag("unit", "temperature").
		AddField("avg", 23.2).
		AddField("max", 45.0).
		SetTime(time.Now())
	writeAPI.WritePoint(context.Background(), p)

	// 3.  Or write directly line protocol
	line := fmt.Sprintf("stat,unit=temperature avg=%f,max=%f", 23.5, 45.0)
	writeAPI.WriteRecord(context.Background(), line)

	// Get query client
	queryAPI := client.QueryAPI(myOrg)
	// Get parser flux query result
	result, err := queryAPI.Query(context.Background(), `from(bucket:"as")|> range(start: -1h) |> filter(fn: (r) => r._measurement == "stat")`)
	if err == nil {
		// Use Next() to iterate over query result lines
		for result.Next() {
			// read result
			fmt.Printf("row: %s\n", result.Record().String())
		}
		if result.Err() != nil {
			fmt.Printf("Query error: %s\n", result.Err().Error())
		}
	}
	// Ensures background processes finishes
	client.Close()
}
