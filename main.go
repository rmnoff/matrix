package main

import (
  "os"
  "flag"
  "log"
  "fmt"
  "time"
  "strconv"
  "strings"
  "encoding/json"

  "database/sql"
  _ "github.com/lib/pq"
  "github.com/jmoiron/sqlx"

  "github.com/jackwhelpton/fasthttp-routing"
  "github.com/jackwhelpton/fasthttp-routing/content"
  "github.com/jackwhelpton/fasthttp-routing/fault"
  "github.com/jackwhelpton/fasthttp-routing/slash"
  "github.com/jackwhelpton/fasthttp-routing/access"
  "github.com/jackwhelpton/fasthttp-routing/file"
  "github.com/erikdubbelboer/fasthttp"
)

func parsePsqlElements(url string) (string, string, string, string, string) {
  split := strings.Split(url, "@")
  unamepwdsplit := strings.Split(split[0], "//")
  unamepwd := strings.Split(unamepwdsplit[1], ":")
  uname := unamepwd[0]
  pwd := unamepwd[1]
  urlportdbname := strings.Split(split[1], ":")
  link := urlportdbname[0]
  portdbname := strings.Split(urlportdbname[1], "/")
  port := portdbname[0]
  dbname := portdbname[1]
  return uname, pwd, link, port, dbname
}

var (
  port      = os.Getenv("PORT")
  // port      = "8080"
  addr      = flag.String("addr", fmt.Sprintf(":%s", port), "TCP address to listen to")
  psqlURL   = os.Getenv("DATABASE_URL")
  dbuname, dbpwd, dblink, dbport, dbname = parsePsqlElements(psqlURL)
  // dblink   = "manny.db.elephantsql.com"
  // dbuname = "fzspbstv"
  // dbname = "fzspbstv"
  // dbpwd   = "ImSLvDaU_NNF1IvdEViKTqezbPwmnXMx"
  // dbport  = "5432"
  psqlInfo  = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s" +
    " sslmode=disable", dblink, dbport, dbuname, dbpwd, dbname)
)

var schema = `
CREATE TABLE IF NOT EXISTS prediction (
  id SERIAL PRIMARY KEY,
  content TEXT,
  edited TEXT,
  created TEXT,
  type_id INT,
  lang_id INT,
  personal BOOLEAN
);

CREATE TABLE IF NOT EXISTS additionalDate (
  id SERIAL PRIMARY KEY,
  firstname VARCHAR(255),
  lastname VARCHAR(255),
  birthdate TEXT,
  owner_id INT
);

CREATE TABLE IF NOT EXISTS predictionRel (
  id SERIAL PRIMARY KEY,
  prediction_id INT,
  combination VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS predictionType (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) UNIQUE
);

CREATE TABLE IF NOT EXISTS predictionLang (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) UNIQUE
);

CREATE TABLE IF NOT EXISTS constantText (
  id SERIAL PRIMARY KEY,
  content TEXT,
  edited TEXT,
  created TEXT,
  type_id INT
);

CREATE TABLE IF NOT EXISTS constantType (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS userProfile (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255),
  firstname VARCHAR(255),
  lastname VARCHAR(255),
  password TEXT,
  birthdate TEXT,
  gender BOOLEAN
);`


type ConstantText struct {
  Header string `json:"header"`
  Content string `json:"content"`
}

type Prediction struct {
  Id int `db:"id"`
  Content string `db:"content"`
  Created sql.NullString `db:"created"`
  Edited sql.NullString `db:"edited"`
  Type int `db:"type_id"`
  Lang int `db:"lang_id"`
  Personal bool `db:"personal"`
  // Foreword []ConstantText `db:"foreword"`
  // ImageName sql.NullString `json:"imageName"`
}

type PredictionType struct {
  Id int `db:"id"`
  Name string `db:"name"`
}

type ContentByGender struct {
  Female string `json:"f"`
  Male string `json:"m"`
}

type Response struct {
  Ok bool `json:"ok"`
  Error string `json:"error"`
  Data []Prediction `json:"data"`
}


func main() {
  flag.Parse()

  db, err := sqlx.Connect("postgres", psqlInfo)
  if err != nil {
    log.Panic(err)
  }
  defer db.Close()
  db.MustExec(schema)
  tx := db.MustBegin()
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('personal features positive') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('personal features negative') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('personal features social') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('money') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('relationship') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('parents') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('kids') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('destiny') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('past life') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('programms') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('life guide') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('health') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('year prediction') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionLang(name) VALUES('russian') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionLang(name) VALUES('english') ON CONFLICT DO NOTHING;`)
  tx.Commit()

  router := routing.New()
  router.Use(
		access.Logger(log.Printf),
		slash.Remover(fasthttp.StatusMovedPermanently),
		fault.Recovery(log.Printf),
	)
  api := router.Group("/api/v1")
  api.Use(content.TypeNegotiator(content.JSON))
	api.Get("/healthcheck", func(c *routing.Context) error {
		return c.Write(`{"ok": true, "error": null}`)
	})
  api.Get("/doc", func(c *routing.Context) error {
		return c.Write(`{"ok": true, "error": null}`)
  })
  api.Get("/auth", func(c *routing.Context) error {
		return c.Write(`{"ok": true, "error": null, "data": "Future authorisation start"}`)
  })
  api.Post("/auth", func(c *routing.Context) error {
		return c.Write(`{"ok": true, "error": null, "data": "Future authorisation end"}`)
  })
  api.Get("/check/<input>", func(c *routing.Context) error {
    input := c.Param("input")
    timestamp := input[:len(input) - 1]
    gender := input[len(input) - 1:]
    fmt.Println(timestamp)
    fmt.Println(gender)
    if timestamp == "" {
      return c.Write(Response{false, "No timestamp provided", nil})
    }
    combo := countBD(timestamp)
    if combo == nil {
      return c.Write(Response{false, "Timestamp corrupted", nil})
    }
    finalCombos := setAllCombos(combo)
    // prog1 := fmt.Sprintf("prog1: [%v %v %v]", combo[0], finalCombos[3], finalCombos[2])
    // prog2 := fmt.Sprintf("prog2: [%v %v %v]", combo[1], finalCombos[5], finalCombos[4])
    // prog3 := fmt.Sprintf("prog3: [%v %v %v]", combo[2], finalCombos[7], finalCombos[6])
    // prog4 := fmt.Sprintf("prog4: [%v %v %v]", finalCombos[13], finalCombos[15], finalCombos[19])
    // prog5 := fmt.Sprintf("prog5: [%v %v %v]", finalCombos[16], finalCombos[14], finalCombos[20])

    pastLifePrediction := Prediction{}
    pastLifePredictionCombo := fmt.Sprintf("%d-%d-%d", finalCombos[8], finalCombos[9], finalCombos[0])
    err = db.Get(&pastLifePrediction, "SELECT * FROM prediction WHERE type_id=1 AND id=(SELECT prediction_id FROM predictionrel WHERE combination=$1)", pastLifePredictionCombo)
    if err != nil {
      log.Println(err)
      return c.Write(Response{false, "Can't parse past life prediction", nil})
    }

    personalFeaturesPos := Prediction{}
    personalFeaturesPosCombo := fmt.Sprintf("%d-%d", combo[0], combo[1])
    personalFeaturesNeg := Prediction{}
    personalFeaturesNegCombo := fmt.Sprintf("%d-%d", combo[0], combo[1])
    personalFeaturesSoc := Prediction{}
    personalFeaturesSocCombo := fmt.Sprintf("%d", finalCombos[1])
    fmt.Println(personalFeaturesPosCombo)
    fmt.Println(personalFeaturesNegCombo)
    fmt.Println(personalFeaturesSocCombo)
    err = db.Get(&personalFeaturesPos, "SELECT * FROM prediction WHERE type_id=2 AND id=(SELECT prediction_id FROM predictionrel WHERE combination=$1)", personalFeaturesPosCombo)
    if err != nil {
      log.Println(err)
      return c.Write(Response{false, "Can't parse positive personal features", nil})
    }
    err = db.Get(&personalFeaturesNeg, "SELECT * FROM prediction WHERE type_id=3 AND id=(SELECT prediction_id FROM predictionrel WHERE combination=$1)", personalFeaturesNegCombo)
    if err != nil {
      log.Println(err)
      return c.Write(Response{false, "Can't parse negative personal features", nil})
    }
    err = db.Get(&personalFeaturesSoc, "SELECT * FROM prediction WHERE type_id=4 AND id=(SELECT prediction_id FROM predictionrel WHERE combination=$1)", personalFeaturesSocCombo)
    if err != nil {
      log.Println(err)
      return c.Write(Response{false, "Can't parse social personal features", nil})
    }

    relationship := Prediction{}
    relationshipCombo := fmt.Sprintf("%d-%d-%d", finalCombos[11], finalCombos[8], finalCombos[10])
    fmt.Println(relationshipCombo)
    err = db.Get(&relationship, "SELECT * FROM prediction WHERE type_id=5 AND id=(SELECT prediction_id FROM predictionrel WHERE combination=$1)", relationshipCombo)
    if err != nil {
      log.Println(err)
      return c.Write(Response{false, "Can't parse relationship prediction", nil})
    }

    lifeGuide := Prediction{}
    lifeGuideCombo := fmt.Sprintf("%d-%d-%d", combo[0], combo[1], finalCombos[1])
    fmt.Println(lifeGuideCombo)
    err = db.Get(&lifeGuide, "SELECT * FROM prediction WHERE type_id=5 AND id=(SELECT prediction_id FROM predictionrel WHERE combination=$1)", lifeGuideCombo)
    if err != nil {
      log.Println(err)
      return c.Write(Response{false, "Can't parse life guide prediction", nil})
    }

    sex := Prediction{}
    sexCombo := fmt.Sprintf("%d-%d-%d", finalCombos[1], finalCombos[17], finalCombos[18])
    fmt.Println(sexCombo)
    err = db.Get(&sex, "SELECT * FROM prediction WHERE type_id=5 AND id=(SELECT prediction_id FROM predictionrel WHERE combination=$1)", sexCombo)
    if err != nil {
      log.Println(err)
      return c.Write(Response{false, "Can't parse life guide prediction", nil})
    }

    data := []Prediction{pastLifePrediction, personalFeaturesPos, personalFeaturesNeg, personalFeaturesSoc, relationship, lifeGuide, sex}
    for _, content := range data {
      contentbygender := ContentByGender{}
      err := json.Unmarshal([]byte(content.Content), &contentbygender)
      if err == nil {
        if gender == "m" {
          content.Content = contentbygender.Male
        } else {
          content.Content = contentbygender.Female
        }
      }
    }

    return c.Write(Response{true, "", data})
  })
  api.Get("/show/types", func(c *routing.Context) error {
    types := []PredictionType{}
    langs := []PredictionType{}
    err = db.Select(&types, "SELECT * FROM predictionType")
    if err != nil {
      return c.Write(Response{false, "Can't parse types", nil})
    }
    err = db.Select(&langs, "SELECT * FROM predictionLang")
    if err != nil {
      return c.Write(Response{false, "Can't parse languages", nil})
    }
    return c.Write([][]PredictionType{types, langs})
  })
  api.Get("/show/predictions", func(c *routing.Context) error {
    predictions := []Prediction{}
    err = db.Select(&predictions, "SELECT * FROM prediction")
    if err != nil {
      return c.Write(Response{false, "Can't parse predictions", nil})
    }
    return c.Write(predictions)
  })
  api.Post("/add", func(c *routing.Context) error {
    ptypeid := c.PostForm("ptypeid")
    combo := c.PostForm("combo")
    prediction := c.PostForm("prediction")
    personal := c.PostForm("personal")
    language := c.PostForm("language")
    currPrediction := Prediction{}
    db.Get(&currPrediction, "SELECT * FROM prediction ORDER BY id DESC LIMIT 1")
    tx := db.MustBegin()
    tx.MustExec(`INSERT INTO prediction(content,type_id,personal,lang_id) VALUES($1,$2,$3,$4)`, prediction, ptypeid, personal, language)
    tx.MustExec(`INSERT INTO predictionRel(prediction_id,combination) VALUES($1,$2)`, currPrediction.Id + 1, combo)
    tx.Commit()
    return c.Write(`Запись добавлена, нажмите назад, чтобы добавить следующую или закройте страничку.`)
  })

  router.Get("/", file.Content("ui/index.html"))

  fasthttp.ListenAndServe(*addr, router.HandleRequest)
}

func countBD(bd string) []int {
  i, err := strconv.ParseInt(bd, 10, 64)
  if err != nil {
    log.Print(err)
    return nil
  }
  tm := time.Unix(i, 0).String()
  UNIXsplit := strings.Split(tm, " ")
  date := UNIXsplit[0]
  dateSplit := strings.Split(date, "-")
  yearSplit := strings.Split(dateSplit[0],"")
  parsedYearOne, _ := strconv.ParseInt(yearSplit[0], 10, 64)
  parsedYearTwo, _ := strconv.ParseInt(yearSplit[1], 10, 64)
  parsedYearThree, _ := strconv.ParseInt(yearSplit[2], 10, 64)
  parsedYearFour, _ := strconv.ParseInt(yearSplit[3], 10, 64)
  probDay, _ := strconv.ParseInt(dateSplit[2], 10, 64)
  probMonth, _ := strconv.ParseInt(dateSplit[1], 10, 64)
  daySum := checkGreater(int(probDay))
  monthSum := probMonth
  yearSum := checkGreater(int(parsedYearOne + parsedYearTwo + parsedYearThree + parsedYearFour))
  return []int{int(daySum), int(monthSum), int(yearSum)}
}

func checkGreater(inum int) int {
  if inum > 22 {
    oinum := inum
    inum = 0
    for _, r := range strings.Split(strconv.Itoa(oinum),"") {
      num, err := strconv.Atoi(r)
      if err != nil {
        log.Print(err)
      }
      inum = inum + num
    }
  }
  return inum
}

func setAllCombos(icombo []int) []int {
  d  := checkGreater(icombo[0]+icombo[1]+icombo[2])
  e  := checkGreater(icombo[0]+icombo[1]+icombo[2]+d)
  a1 := checkGreater(icombo[0]+e)
  a2 := checkGreater(icombo[0]+a1)
  b1 := checkGreater(icombo[1]+e)
  b2 := checkGreater(icombo[1]+b1)
  c1 := checkGreater(icombo[2]+e)
  c2 := checkGreater(icombo[2]+c1)
  d1 := checkGreater(d+e)
  d2 := checkGreater(d+d1)
  x  := checkGreater(c1+d1)
  x1 := checkGreater(d1+x)
  x2 := checkGreater(c1+x)
  f  := checkGreater(icombo[0]+icombo[1])
  g  := checkGreater(icombo[1]+icombo[2])
  y  := checkGreater(icombo[2]+d)
  k  := checkGreater(icombo[0]+d)
  e1 := checkGreater(f+g+y+k)
  e2 := checkGreater(e+e1)
  o  := checkGreater(f+y)
  u  := checkGreater(k+g)
  h  := checkGreater(icombo[1]+d)
  j  := checkGreater(icombo[0]+icombo[2])
  m  := checkGreater(h+j)
  n  := checkGreater(f+y)
  t  := checkGreater(g+k)
  z  := checkGreater(n+t)
  s  := checkGreater(m+z)
  b3 := checkGreater(b1+e)
  a3 := checkGreater(a1+e)
  l  := checkGreater(icombo[0]+icombo[1])
  l1 := checkGreater(a2+b2)
  l2 := checkGreater(a1+b1)
  l3 := checkGreater(a3+b3)
  l4 := checkGreater(icombo[0]+icombo[1]+icombo[2]+d+e)
  l5 := checkGreater(d1+c1)
  l6 := checkGreater(icombo[2]+d)
  d3 := checkGreater(icombo[0]+a2+a1+a3+e+d1+d)
  c3 := checkGreater(icombo[1]+b2+b1+b3+e+c1+icombo[2])
  e3 := checkGreater(l+l1+l2+l3+l4+l5+l6)
  return []int{d,e,a1,a2,b1,b2,c1,c2,d1,d2,x,x1,x2,f,g,y,k,e1,e2,o,u,h,j,m,n,
    t,z,s,b3,a3,l,l1,l2,l3,l4,l5,l6,d3,c3,e3}
}
