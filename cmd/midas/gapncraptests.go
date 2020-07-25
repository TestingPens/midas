package main

import (
	"fmt"
	"os"
	"time"
)

type DataPoints struct {
	PriceOpen      float64
	PctHOD         float64
	PctLOD         float64
	PctPMH         float64
	HODTime        int
	LODTime        int
	PctLowBef12    float64
	PctLowAft12    float64
	UnderPMHAft12  int
	UnderOpenAft12 int
	Pct5mOpen      float64
	Pct10mOpen     float64
	Pct15mOpen     float64
}

type Results struct {
	PriceOpen      float64
	PctHOD         float64
	PctLOD         float64
	PctPMH         float64
	HODTime        int
	LODTime        int
	PctLowBef12    float64
	PctLowAft12    float64
	UnderPMHAft12  int
	UnderOpenAft12 int
	Pct5mOpen      float64
	Pct10mOpen     float64
	Pct15mOpen     float64
}

type DataList struct {
	dplist []DataPoints
}

func rungapncraptests() {
	fmt.Println("Running tests...\n")
	var charts []Chart
	var datalist DataList
	var results Results

	getallcharts(&charts)
	analysegapncrap(charts, &datalist)
	getresults(&datalist, &results)
	printresults(results, len(datalist.dplist))
}

func printresults(results Results, n int) {
	fmt.Println(fmt.Sprintf("Open to HOD Pct: %f", results.PctHOD/float64(n)))
	fmt.Println(fmt.Sprintf("Open to LOD Pct: %f", results.PctLOD/float64(n)))
	fmt.Println(fmt.Sprintf("Open to PMH Pct: %f", results.PctPMH/float64(n)))
	fmt.Println(fmt.Sprintf("Open to Low Pct Before 12: %f", results.PctLowBef12/float64(n)))
	fmt.Println(fmt.Sprintf("Open to Low Pct After 12: %f", results.PctLowAft12/float64(n)))
	fmt.Println(fmt.Sprintf("5 Min Pct: %f", results.Pct5mOpen/float64(n)))
	fmt.Println(fmt.Sprintf("10 Min Pct: %f", results.Pct10mOpen/float64(n)))
	fmt.Println(fmt.Sprintf("15 Min Pct: %f", results.Pct15mOpen/float64(n)))
}

func getresults(datalist *DataList, results *Results) {

	for _, datapoints := range datalist.dplist {
		results.PctHOD += datapoints.PctHOD
		results.PctLOD += datapoints.PctLOD
		results.PctPMH += datapoints.PctPMH
		results.PctLowBef12 += datapoints.PctLowBef12
		results.PctLowAft12 += datapoints.PctLowAft12
		results.HODTime += datapoints.HODTime
		results.LODTime += datapoints.LODTime
		results.UnderPMHAft12 += datapoints.UnderPMHAft12
		results.UnderOpenAft12 += datapoints.UnderOpenAft12
		results.Pct5mOpen += datapoints.Pct5mOpen
		results.Pct10mOpen += datapoints.Pct10mOpen
		results.Pct15mOpen += datapoints.Pct15mOpen
	}
}

func analysegapncrap(charts []Chart, datalist *DataList) {
	for _, chart := range charts {
		var datapoints DataPoints
		getopen(chart, &datapoints)
		getpmh(chart, &datapoints)
		gethodlod(chart, &datapoints)
		getpricelowbef12(chart, &datapoints)
		getpricelowaft12(chart, &datapoints)
		getunderpmhaft12(chart, &datapoints)
		getunderopenaft12(chart, &datapoints)
		getkeyprices(chart, &datapoints)

		datalist.dplist = append(datalist.dplist, datapoints)
	}
}

func getopen(chart Chart, datapoints *DataPoints) {
	date := time.Unix(int64(chart.Chart.Result[0].Timestamp[0]), 0)
	opentime := time.Date(date.Year(), date.Month(), date.Day(), 13, 30, 0, 0, time.UTC)
	for i, t := range chart.Chart.Result[0].Timestamp {
		if t == int(opentime.Unix()) {
			datapoints.PriceOpen = chart.Chart.Result[0].Indicators.Quote[0].Open[i]
			return
		}
	}
	fmt.Println("Couldnt not find open price for " + chart.Chart.Result[0].Meta.Symbol)
	os.Exit(1)
}

func getpmh(chart Chart, datapoints *DataPoints) {
	date := time.Unix(int64(chart.Chart.Result[0].Timestamp[0]), 0)
	opentime := time.Date(date.Year(), date.Month(), date.Day(), 13, 30, 0, 0, time.UTC)
	temphigh := 0.0

	for i, t := range chart.Chart.Result[0].Timestamp {
		if t < int(opentime.Unix()) {
			high := chart.Chart.Result[0].Indicators.Quote[0].High[i]
			if high > temphigh {
				temphigh = high
			}
		} else {
			break
		}
	}

	if temphigh != 0.0 {
		datapoints.PctPMH = (temphigh - datapoints.PriceOpen) / datapoints.PriceOpen
	} else {
		fmt.Println("Could not find PMH for: " + chart.Chart.Result[0].Meta.Symbol)
		os.Exit(1)
	}
}

func gethodlod(chart Chart, datapoints *DataPoints) {
	date := time.Unix(int64(chart.Chart.Result[0].Timestamp[0]), 0)
	opentime := time.Date(date.Year(), date.Month(), date.Day(), 13, 30, 0, 0, time.UTC)
	closetime := time.Date(date.Year(), date.Month(), date.Day(), 20, 0, 0, 0, time.UTC)
	temphod := datapoints.PriceOpen
	temphodtime := 0
	templod := datapoints.PriceOpen
	templodtime := 0

	for i, t := range chart.Chart.Result[0].Timestamp {
		if t > int(opentime.Unix()) && t < int(closetime.Unix()) {
			high := chart.Chart.Result[0].Indicators.Quote[0].High[i]
			low := chart.Chart.Result[0].Indicators.Quote[0].Low[i]
			if high > temphod {
				temphod = high
				temphodtime = t
			}

			if low < templod {
				templod = low
				templodtime = t
			}
		}
	}

	if temphod != datapoints.PriceOpen && templod != datapoints.PriceOpen {
		datapoints.PctHOD = (temphod - datapoints.PriceOpen) / datapoints.PriceOpen
		datapoints.HODTime = temphodtime
		datapoints.PctLOD = (templod - datapoints.PriceOpen) / datapoints.PriceOpen
		datapoints.LODTime = templodtime
	} else {
		fmt.Println("Could not find HOD/LOD for: " + chart.Chart.Result[0].Meta.Symbol)
		os.Exit(1)
	}

}

func getpricelowbef12(chart Chart, datapoints *DataPoints) {
	date := time.Unix(int64(chart.Chart.Result[0].Timestamp[0]), 0)
	opentime := time.Date(date.Year(), date.Month(), date.Day(), 13, 30, 0, 0, time.UTC)
	midday := time.Date(date.Year(), date.Month(), date.Day(), 16, 0, 0, 0, time.UTC)
	templow := datapoints.PriceOpen

	for i, t := range chart.Chart.Result[0].Timestamp {
		if t > int(opentime.Unix()) && t < int(midday.Unix()) {
			low := chart.Chart.Result[0].Indicators.Quote[0].Low[i]
			if low < templow {
				templow = low
			}
		}
	}

	if templow != datapoints.PriceOpen {
		datapoints.PctLowBef12 = (templow - datapoints.PriceOpen) / datapoints.PriceOpen
	} else {
		fmt.Println("Could not find low before 12 for: " + chart.Chart.Result[0].Meta.Symbol)
		os.Exit(1)
	}
}

func getpricelowaft12(chart Chart, datapoints *DataPoints) {
	date := time.Unix(int64(chart.Chart.Result[0].Timestamp[0]), 0)
	midday := time.Date(date.Year(), date.Month(), date.Day(), 16, 0, 0, 0, time.UTC)
	closetime := time.Date(date.Year(), date.Month(), date.Day(), 20, 0, 0, 0, time.UTC)
	templow := datapoints.PriceOpen

	for i, t := range chart.Chart.Result[0].Timestamp {
		if t > int(midday.Unix()) && t < int(closetime.Unix()) {
			low := chart.Chart.Result[0].Indicators.Quote[0].Low[i]
			if low < templow {
				templow = low
			}
		}
	}

	if templow != datapoints.PriceOpen {
		datapoints.PctLowAft12 = (templow - datapoints.PriceOpen) / datapoints.PriceOpen
	} else {
		fmt.Println("Could not find low after 12 for: " + chart.Chart.Result[0].Meta.Symbol)
		os.Exit(1)
	}
}

func getunderpmhaft12(chart Chart, datapoints *DataPoints) {
	date := time.Unix(int64(chart.Chart.Result[0].Timestamp[0]), 0)
	midday := time.Date(date.Year(), date.Month(), date.Day(), 16, 0, 0, 0, time.UTC)
	closetime := time.Date(date.Year(), date.Month(), date.Day(), 20, 0, 0, 0, time.UTC)
	underpmh := 1

	for i, t := range chart.Chart.Result[0].Timestamp {
		if t > int(midday.Unix()) && t < int(closetime.Unix()) {
			high := chart.Chart.Result[0].Indicators.Quote[0].High[i]
			highpct := (high - datapoints.PriceOpen) / datapoints.PriceOpen
			if highpct > datapoints.PctPMH {
				underpmh = 0
				break
			}
		}
	}

	datapoints.UnderPMHAft12 = underpmh
}

func getunderopenaft12(chart Chart, datapoints *DataPoints) {
	date := time.Unix(int64(chart.Chart.Result[0].Timestamp[0]), 0)
	midday := time.Date(date.Year(), date.Month(), date.Day(), 16, 0, 0, 0, time.UTC)
	closetime := time.Date(date.Year(), date.Month(), date.Day(), 20, 0, 0, 0, time.UTC)
	underopen := 1

	for i, t := range chart.Chart.Result[0].Timestamp {
		if t > int(midday.Unix()) && t < int(closetime.Unix()) {
			high := chart.Chart.Result[0].Indicators.Quote[0].High[i]
			if high > datapoints.PriceOpen {
				underopen = 0
				break
			}
		}
	}

	datapoints.UnderOpenAft12 = underopen
}

func getkeyprices(chart Chart, datapoints *DataPoints) {
	date := time.Unix(int64(chart.Chart.Result[0].Timestamp[0]), 0)
	time5m := time.Date(date.Year(), date.Month(), date.Day(), 13, 35, 0, 0, time.UTC)
	time10m := time.Date(date.Year(), date.Month(), date.Day(), 13, 40, 0, 0, time.UTC)
	time15m := time.Date(date.Year(), date.Month(), date.Day(), 13, 45, 0, 0, time.UTC)

	for i, t := range chart.Chart.Result[0].Timestamp {
		if t == int(time5m.Unix()) {
			datapoints.Pct5mOpen = (chart.Chart.Result[0].Indicators.Quote[0].Close[i] - datapoints.PriceOpen) / datapoints.PriceOpen
		}

		if t == int(time10m.Unix()) {
			datapoints.Pct10mOpen = (chart.Chart.Result[0].Indicators.Quote[0].Close[i] - datapoints.PriceOpen) / datapoints.PriceOpen
		}

		if t == int(time15m.Unix()) {
			datapoints.Pct15mOpen = (chart.Chart.Result[0].Indicators.Quote[0].Close[i] - datapoints.PriceOpen) / datapoints.PriceOpen
		}

		if t > int(time15m.Unix()) {
			break
		}
	}

}
