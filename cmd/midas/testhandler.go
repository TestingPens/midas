package main

import (
	"time"
)

func testhandler(symbol string, date time.Time, token string) {

	getchart(symbol, date, "5", token)
	rungapncraptests()
}
