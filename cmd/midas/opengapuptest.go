package main

import (
	"fmt"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"context"
)

type OpenGapUpTestResults struct {
	Symbol string					`bson:"symbol"`
	Date int						`bson:"date"`
	PrePriceHigh float64			`bson:"prepricehigh"`
	PrePriceOpenToHighPct float64	`bson:"prepriceopentohighpct"`
	PrePriceHighTime int			`bson:"prepricehightime"`
	PrePriceHighToClosePct float64	`bson:"prepricehightoclosepct"`
	Vol5m int						`bson:"vol5m"`
	Vol10m int						`bson:"vol10m"`
	Vol15m int						`bson:"vol15m"`
	Vol20m int						`bson:"vol20m"`
	Vol25m int						`bson:"vol25m"`
	Vol30m int						`bson:"vol30m"`
	OverPreHigh30m int				`bson:"overprehigh30m"`
	OverPreHigh60m int				`bson:"overprehigh60m"`
}

func runopengapuptests(chart Chart, date time.Time, opengapuptestresults *OpenGapUpTestResults) {
	fmt.Println("Running tests...\n")
	opengapuptestresults.Symbol = chart.Chart.Result[0].Meta.Symbol
	opengapuptestresults.Date = chart.Chart.Result[0].Timestamp[0]
	getprepriceopentohighpct(chart, date, opengapuptestresults)
	getprepricehightoclosepct(chart, date, opengapuptestresults)
	getvolafteropen(chart, date, opengapuptestresults)
	getoverprehigh(chart, date, opengapuptestresults)
}

func getprepriceopentohighpct(chart Chart, date time.Time, opengapuptestresults *OpenGapUpTestResults) {
	to := time.Date(date.Year(), date.Month(), date.Day(), 13, 25, 0, 0, time.UTC)
	prepricehigh := 0.0
	prepricehightime := 0
	pct := 0.0
	open := chart.Chart.Result[0].Indicators.Quote[0].Open[0]

	for i, t := range chart.Chart.Result[0].Timestamp {
		high := chart.Chart.Result[0].Indicators.Quote[0].High[i]
		if high > prepricehigh {
			prepricehigh = high
			prepricehightime = t
		}
		if t == int(to.Unix()) {
			break
		}
	}
	pct = (prepricehigh - open) / open
	opengapuptestresults.PrePriceHigh = prepricehigh
	opengapuptestresults.PrePriceOpenToHighPct = pct
	opengapuptestresults.PrePriceHighTime = prepricehightime
}

func getprepricehightoclosepct(chart Chart, date time.Time, opengapuptestresults *OpenGapUpTestResults) {
	to := time.Date(date.Year(), date.Month(), date.Day(), 13, 25, 0, 0, time.UTC)
	preclose := 0.0
	pct := 0.0


	for i, t := range chart.Chart.Result[0].Timestamp {
		if t == int(to.Unix()) {
			preclose = chart.Chart.Result[0].Indicators.Quote[0].Close[i]
			break
		}
	}
	pct = (opengapuptestresults.PrePriceHigh - preclose) / preclose
	opengapuptestresults.PrePriceHighToClosePct = pct
}

func getvolafteropen(chart Chart, date time.Time, opengapuptestresults *OpenGapUpTestResults) {
	time5 := time.Date(date.Year(), date.Month(), date.Day(), 13, 30, 0, 0, time.UTC)
	time10 := time.Date(date.Year(), date.Month(), date.Day(), 13, 35, 0, 0, time.UTC)
	time15 := time.Date(date.Year(), date.Month(), date.Day(), 13, 40, 0, 0, time.UTC)
	time20 := time.Date(date.Year(), date.Month(), date.Day(), 13, 45, 0, 0, time.UTC)
	time25 := time.Date(date.Year(), date.Month(), date.Day(), 13, 50, 0, 0, time.UTC)
	time30 := time.Date(date.Year(), date.Month(), date.Day(), 13, 55, 0, 0, time.UTC)

	for i, t := range chart.Chart.Result[0].Timestamp {
		switch t {
		case int(time5.Unix()):
			opengapuptestresults.Vol5m = chart.Chart.Result[0].Indicators.Quote[0].Volume[i]
		case int(time10.Unix()):
			opengapuptestresults.Vol10m = chart.Chart.Result[0].Indicators.Quote[0].Volume[i]
		case int(time15.Unix()):
			opengapuptestresults.Vol15m = chart.Chart.Result[0].Indicators.Quote[0].Volume[i]
		case int(time20.Unix()):
			opengapuptestresults.Vol20m = chart.Chart.Result[0].Indicators.Quote[0].Volume[i]
		case int(time25.Unix()):
			opengapuptestresults.Vol25m = chart.Chart.Result[0].Indicators.Quote[0].Volume[i]
		case int(time30.Unix()):
			opengapuptestresults.Vol30m = chart.Chart.Result[0].Indicators.Quote[0].Volume[i]
			break
		}
	}
}

func getoverprehigh(chart Chart, date time.Time, opengapuptestresults *OpenGapUpTestResults) {
	time30 := time.Date(date.Year(), date.Month(), date.Day(), 13, 55, 0, 0, time.UTC)
	time60 := time.Date(date.Year(), date.Month(), date.Day(), 14, 25, 0, 0, time.UTC)

	for i, t := range chart.Chart.Result[0].Timestamp {
		switch t {
		case int(time30.Unix()):
			close30m := chart.Chart.Result[0].Indicators.Quote[0].Close[i]
			if close30m > opengapuptestresults.PrePriceHigh {
				opengapuptestresults.OverPreHigh30m = 1
			}
		case int(time60.Unix()):
			close60m := chart.Chart.Result[0].Indicators.Quote[0].Close[i]
			if close60m > opengapuptestresults.PrePriceHigh {
				opengapuptestresults.OverPreHigh60m = 0
			}
			break
		}
	}
}

func writeresults(collection string, opengapuptestresults *OpenGapUpTestResults) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("midas")
	col := db.Collection(collection)

	_, err = col.InsertOne(ctx, opengapuptestresults)
	if err != nil {
		panic(err)
	}
}

func getaverages(collection string) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)
	db := client.Database("midas")
	col := db.Collection("opengapuptests")
	var opengapuptestresults []OpenGapUpTestResults

	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}

	if err = cursor.All(ctx, &opengapuptestresults); err != nil {
		panic(err)
	}

	var long30mresults []OpenGapUpTestResults
	var long60mresults []OpenGapUpTestResults
	var short30mresults []OpenGapUpTestResults
	var short60mresults []OpenGapUpTestResults

	for _, test := range opengapuptestresults {
		if test.OverPreHigh30m == 1 {
			long30mresults = append(long30mresults, test)
		} else {
			short30mresults = append(short30mresults, test)
		}

		if test.OverPreHigh60m == 1 {
			long60mresults = append(long60mresults, test)
		} else {
			short60mresults = append(short60mresults, test)
		}
	}

	fmt.Println("Long after 30m...")
	l30mpphigh := 0.0
	l30mppopentohigh := 0.0
	l30mpphighttoclose := 0.0
	l30mvol5m := 0
	l30mvol10m := 0
	l30mvol15m := 0
	l30mvol20m := 0
	l30mvol25m := 0
	l30mvol30m := 0


	for _, test := range long30mresults {
		l30mpphigh += test.PrePriceHigh
		l30mppopentohigh += test.PrePriceOpenToHighPct
		l30mpphighttoclose += test.PrePriceHighToClosePct
		l30mvol5m += test.Vol5m
		l30mvol10m += test.Vol10m
		l30mvol15m += test.Vol15m
		l30mvol20m += test.Vol20m
		l30mvol25m += test.Vol25m
		l30mvol30m += test.Vol30m
	}
	fmt.Println(fmt.Sprintf("Pre Price High: %f | Pre Price Open to High: %f | Pre Price High to Open: %f | 5m Vol: %d | 10m Vol: %d | 15m Vol: %d | 20m Vol: %d | 25m Vol: %d | 30m Vol: %d \n", l30mpphigh/float64(len(long30mresults)), l30mppopentohigh/float64(len(long30mresults)), l30mpphighttoclose/float64(len(long30mresults)), l30mvol5m/len(long30mresults), l30mvol10m/len(long30mresults), l30mvol15m/len(long30mresults), l30mvol20m/len(long30mresults), l30mvol25m/len(long30mresults), l30mvol30m/len(long30mresults)))

	fmt.Println("Short after 30m...")
	s30mpphigh := 0.0
	s30mppopentohigh := 0.0
	s30mpphighttoclose := 0.0
	s30mvol5m := 0
	s30mvol10m := 0
	s30mvol15m := 0
	s30mvol20m := 0
	s30mvol25m := 0
	s30mvol30m := 0


	for _, test := range short30mresults {
		s30mpphigh += test.PrePriceHigh
		s30mppopentohigh += test.PrePriceOpenToHighPct
		s30mpphighttoclose += test.PrePriceHighToClosePct
		s30mvol5m += test.Vol5m
		s30mvol10m += test.Vol10m
		s30mvol15m += test.Vol15m
		s30mvol20m += test.Vol20m
		s30mvol25m += test.Vol25m
		s30mvol30m += test.Vol30m
		// fmt.Println(test.Symbol, test.PrePriceHighToClosePct)
	}
	fmt.Println(fmt.Sprintf("Pre Price High: %f | Pre Price Open to High: %f | Pre Price High to Open: %f | 5m Vol: %d | 10m Vol: %d | 15m Vol: %d | 20m Vol: %d | 25m Vol: %d | 30m Vol: %d \n", s30mpphigh/float64(len(short30mresults)), s30mppopentohigh/float64(len(short30mresults)), s30mpphighttoclose/float64(len(short30mresults)), s30mvol5m/len(short30mresults), s30mvol10m/len(short30mresults), s30mvol15m/len(short30mresults), s30mvol20m/len(short30mresults), s30mvol25m/len(short30mresults), s30mvol30m/len(short30mresults)))
}