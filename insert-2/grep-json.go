package main

import (
    "net/http"
    "io"
    "encoding/json"
    "fmt"
    "time"
    "log"
    influxdb2 "github.com/influxdata/influxdb-client-go/v2"
    "context"
)

const (
	myToken  = "ot3rcVMa4Y1V8S3otRMFItQiX0vChFbgeSOe2gOlv7RjYk5NpMnUaJBOND9U9JvILr50tCus6v36tL4pJ2SBQQ=="
	myOrg    = "5261a190ceb20da3"
	myBucket = "power_meter"
)

func main() {
  client := influxdb2.NewClient("http://10.81.81.36:8099", myToken)

	// Use blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking(myOrg, myBucket)

  voltage1 := float64(0.0)
  current1 := float64(0.0)
  power1 := float64(0.0)
  energy1 := float64(0.0)
  frequency1 := float64(0.0)
  pf1 := float64(0.0)

  voltage2 := float64(0.0)
  current2 := float64(0.0)
  power2 := float64(0.0)
  energy2 := float64(0.0)
  frequency2 := float64(0.0)
  pf2 := float64(0.0)

  voltage3 := float64(0.0)
  current3 := float64(0.0)
  power3 := float64(0.0)
  energy3 := float64(0.0)
  frequency3 := float64(0.0)
  pf3 := float64(0.0)

  p := influxdb2.NewPoint("stat",
		map[string]string{"unit": "power_data"},
		map[string]interface{}{
      "voltage1": 0.0, 
      "current1": 0.0,
      "power1": 0.0,
      "energy1": 0.0,
      "frequency1": 0.0,
      "pf1": 0.0,
      "voltage2": 0.0, 
      "current2": 0.0,
      "power2": 0.0,
      "energy2": 0.0,
      "frequency2": 0.0,
      "pf2": 0.0,
      "voltage3": 0.0, 
      "current3": 0.0,
      "power3": 0.0,
      "energy3": 0.0,
      "frequency3": 0.0,
      "pf3": 0.0,
    },
		time.Now())
	// write point immediately
	writeAPI.WritePoint(context.Background(), p)
  
  fmt.Printf("Running...\n")

  for true {
    r, _ := http.Get("http://10.81.84.26:8390/data")
  //   http://10.81.84.26:8390/data
    defer r.Body.Close()
  
    b, _ := io.ReadAll(r.Body)
    
    var res map[string]interface{}
    jsonErr := json.Unmarshal(b, &res)
    if jsonErr != nil {
      log.Fatal(jsonErr)
    }else{
      // fmt.Printf("Voltage 1 : %0.1f\n", res["voltage1"].(float64))
      // fmt.Printf("Voltage 2 : %0.1f\n", res["voltage2"].(float64))

      voltage1 = res["voltage1"].(float64)
      current1 = res["current1"].(float64)
      power1 = res["power1"].(float64)
      energy1 = res["energy1"].(float64)
      frequency1 = res["frequency1"].(float64)
      pf1 = res["pf1"].(float64)

      voltage2 = res["voltage2"].(float64)
      current2 = res["current2"].(float64)
      power2 = res["power2"].(float64)
      energy2 = res["energy2"].(float64)
      frequency2 = res["frequency2"].(float64)
      pf2 = res["pf2"].(float64)

      voltage3 = res["voltage3"].(float64)
      current3 = res["current3"].(float64)
      power3 = res["power3"].(float64)
      energy3 = res["energy3"].(float64)
      frequency3 = res["frequency3"].(float64)
      pf3 = res["pf3"].(float64)
    }
    
    p = influxdb2.NewPointWithMeasurement("stat").
        AddTag("unit", "power_data").
        AddField("voltage1", voltage1).
        AddField("current1", current1).
        AddField("power1", power1).
        AddField("energy1", energy1).
        AddField("frequency1", frequency1).
        AddField("pf1", pf1).
        AddField("voltage2", voltage2).
        AddField("current2", current2).
        AddField("power2", power2).
        AddField("energy2", energy2).
        AddField("frequency2", frequency2).
        AddField("pf2", pf2).
        AddField("voltage3", voltage3).
        AddField("current3", current3).
        AddField("power3", power3).
        AddField("energy3", energy3).
        AddField("frequency3", frequency3).
        AddField("pf3", pf3).
        SetTime(time.Now())
    writeAPI.WritePoint(context.Background(), p)
  //   print(res["voltage1"].(string))
      
    time.Sleep(1 * time.Second)
  }
  client.Close()
}


// {
//     "voltage1": 229.5,
//     "current1": 21.75,
//     "power1": 4599.1,
//     "energy1": 3714.43,
//     "frequency1": 50.1,
//     "pf1": 0.92,
//     "voltage2": 235.9,
//     "current2": 15.88,
//     "power2": 3231.3,
//     "energy2": 2519.45,
//     "frequency2": 50.1,
//     "pf2": 0.86,
//     "voltage3": 0,
//     "current3": 0,
//     "power3": 0,
//     "energy3": 0,
//     "frequency3": 0,
//     "pf3": 0
//     }