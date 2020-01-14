package main

import (
  "flag"
  "log"
  "fmt"
  "time"
  "strconv"
  "strings"

  "database/sql"

  "github.com/valyala/fasthttp"
  _ "github.com/lib/pq"
  "github.com/jmoiron/sqlx"
)

var (
  addr = flag.String("addr", ":8080", "TCP address to listen to")
)

var schema := `
CREATE TABLE IF NOT EXISTS text (
  id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  content TEXT,
  edited DATE,
  created DATE
);

CREATE TABLE IF NOT EXISTS user (
  id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  email VARCHAR(255),
  firstname VARCHAR(255),
  lastname VARCHAR(255),
  password TEXT,
  birthdate TEXT,
  gender BOOLEAN
);`

func countBD(bd string) {
  i, err := strconv.ParseInt(bd, 10, 64)
  if err != nil {
    log.Fprintf(err)
  }
  tm := time.Unix(i, 0)
  UNIXsplit := strings.Split(tm, " ")
  date := UNIXsplit[0]
  dateSplit := strings.Split(date, "-")
  daySplit, monthSplit, yearSplit := strings.Split(s[0],""), strings.Split(s[1],""), strings.Split(s[2],"")
  parsedDayOne, err := strconv.ParseInt(day, 10, 64)
  if err != nil {
    log.Fprintf(err)
  }
  parsedDayTwo, err := strconv.ParseInt(day, 10, 64)
  if err != nil {
    log.Fprintf(err)
  }
  parsedMonthOne, err := strconv.ParseInt(month, 10, 64)
  if err != nil {
    log.Fprintf(err)
  }
  parsedMonthTwo, err := strconv.ParseInt(month, 10, 64)
  if err != nil {
    log.Fprintf(err)
  }
  parsedYearOne, err := strconv.ParseInt(year, 10, 64)
  if err != nil {
    log.Fprintf(err)
  }
  parsedYearTwo, err := strconv.ParseInt(year, 10, 64)
  if err != nil {
    log.Fprintf(err)
  }
  parsedYearThree, err := strconv.ParseInt(year, 10, 64)
  if err != nil {
    log.Fprintf(err)
  }
  parsedYearFour, err := strconv.ParseInt(year, 10, 64)
  if err != nil {
    log.Fprintf(err)
  }
  daySum := parsedDayOne + parsedDayTwo
  monthSum := parsedMonthOne + parsedMonthTwo
  yearSum := parsedYearOne + parsedYearTwo + parsedYearThree + parsedYearFour
}

func main() {
  flag.Parse()

  handler := requestHandler

  if err := fasthttp.ListenAndServe(*addr, handler); err != nil {
    log.Fatalf("Error in ListenAndServe: %s", err)
  }
}

func requestHandler(ctx *fasthttp.RequestCtx) {
  fmt.Fprintf(ctx, "Hello, world\n\n")
}
