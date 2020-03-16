package main

import (
  "os"
  "flag"
  "log"
  "fmt"
  "time"
  "sort"
  "strconv"
  "strings"
  "net/http"
  "encoding/json"

  "database/sql"
  _ "github.com/lib/pq"
  "github.com/jmoiron/sqlx"

  "github.com/jackwhelpton/fasthttp-routing"
  // "github.com/jackwhelpton/fasthttp-routing/content"
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
  birthdate = "834883200"
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

type User struct {
  Id int `db:"id"`
  Email sql.NullString `db:"email"`
  Firstname sql.NullString `db:"firstname"`
  Lastname sql.NullString `db:"lastname"`
  Password sql.NullString `db:"password"`
  Birthdate sql.NullString `db:"birthdate"`
  Gender bool `db:"gender"`
}


type ConstantText struct {
  Header string `json:"header"`
  Content string `json:"content"`
}

type Prediction struct {
  ImageName string `json:"imageName"`
  Title string `json:"title"`
  BlockType string `json:"blockType"`
  Blocks []Block `json:"blocksOrig"`
  BlocksJSON []BlockJSON `json:"blocks"`
}

type Block struct {
  Id int `db:"id",json:"id"`
  Content string `db:"content",json:"content"`
  Created sql.NullString `db:"created",json:"created"`
  Edited sql.NullString `db:"edited",json:"edited"`
  PredType int `db:"type_id",json:"type_id"`
  Lang int `db:"lang_id",json:"lang_id"`
  Personal bool `db:"personal",json:"personal"`
  Type string `json:type`
  Title string `json:title`
  TintColor *string `json:tintColor`
}

type BlockJSON struct {
  Id int `json:"id"`
  Content string `json:"text"`
  Created sql.NullString `json:"created"`
  Edited sql.NullString `json:"edited"`
  PredType int `json:"type_id"`
  Lang int `json:"lang_id"`
  Personal bool `json:"personal"`
  Type string `json:"type"`
  Title string `json:"title"`
  TintColor *string `json:"tintColor"`
}

type PredictionType struct {
  Id int `db:"id"`
  Name string `db:"name"`
  Short sql.NullString `db:"short"`
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

type ResponseNew struct {
  Ok bool `json:"ok"`
  Error string `json:"error"`
  Data []Block `json:"data"`
}

type ResponseTemplate struct {
  Ok bool `json:"ok"`
  Error string `json:"error"`
  Data string `json:"data"`
}

type ResponseType struct {
  Ok bool `json:"ok"`
  Error string `json:"error"`
  Data [][]PredictionType `json:"data"`
}

type ResponseAdd struct {
  Id int `json:"id"`
}

type Language struct {
  Id int `db:"id"`
  Name string `db:"name"`
  Short string `db:"short"`
}

type Combos struct {
  A  int
  B  int
  C  int
  D  int
  E  int
  A1 int
  A2 int
  B1 int
  B2 int
  C1 int
  C2 int
  D1 int
  D2 int
  X  int
  X1 int
  X2 int
  F  int
  G  int
  Y  int
  K  int
  E1 int
  E2 int
  O  int
  U  int
  H  int
  J  int
  M  int
  N  int
  T  int
  Z  int
  S  int
  B3 int
  A3 int
  L  int
  L1 int
  L2 int
  L3 int
  L4 int
  L5 int
  L6 int
  D3 int
  C3 int
  E3 int
  F1 int
  A4 int
  A5 int
  A6 int
  F2 int
  F3 int
  F4 int
  F5 int
  F6 int
  F7 int
  F8 int
  B4 int
  B5 int 
  B6 int 
  K1 int 
  B7 int 
  B8 int 
  K2 int 
  K3 int 
  K4 int 
  K5 int 
  C4 int 
  K7 int 
  K8 int 
  K6 int 
  C5 int 
  C6 int 
  C7 int 
  Y1 int 
  Y2 int 
  Y3 int 
  Y4 int 
  Y5 int 
  Y6 int 
  Y7 int 
  D4 int 
  Y8 int 
  D6 int 
  D5 int 
  D7 int 
  D8 int 
  D9 int 
  G1 int 
  G2 int 
  G3 int 
  G4 int 
  G5 int 
  G6 int 
  G7 int 
  T1 int 
  T2 int 
  T3 int 
  T4 int 
  T5 int 
  T6 int 
  T7 int 
  O1 int 
  O2 int 
  O3 int 
  X3 int 
  X4 int 
  X5 int 
  X6 int 
  X8 int 
  E4 int 
  E5 int 
  E6 int 
  E7 int 
  E8 int 
  Z1 int 
  Z2 int 
  Z4 int 
  Z5 int 
  Z6 int 
  Z7 int 
  Z8 int 
  H1 int
  H3 int 
  H4 int 
  H5 int 
  H6 int 
  H7 int 
  H8 int 
  N1 int
}

// personalfeatures, destiny, money, programs, sexiness, pastlife, parents, kids, relationships, health, lifeguide

type LocaleBlock struct {
  En string
  Ru string
}

type Locale struct {
  PersonalFeaturesPos LocaleBlock
  PersonalFeaturesNeg LocaleBlock
  PersonalFeaturesSoc LocaleBlock
  Destiny20           LocaleBlock
  Destiny40           LocaleBlock
  DestinyCommon       LocaleBlock
  Money               LocaleBlock
  Programs            LocaleBlock
  Program             LocaleBlock
  Sexiness            LocaleBlock
  PastLife            LocaleBlock
  Parents             LocaleBlock
  Kids                LocaleBlock
  Relationship        LocaleBlock
  Health              LocaleBlock
  LifeGuide           LocaleBlock
  Important           LocaleBlock
  InsultParentsMen    LocaleBlock
  InsultParentsWomen  LocaleBlock
  Resentment          LocaleBlock
  ToBecomeSuccessful  LocaleBlock

}

func newLocale() Locale {
  return Locale{
    LocaleBlock{"Личные Качества позитив", "Personal Features Positive"},
    LocaleBlock{"Личные Качества негатив", "Personal Features Negative"},
    LocaleBlock{"Личные Качества общение", "Personal Features Social"},
    LocaleBlock{"Предназначение 20-40", "Destiny 20-40"},
    LocaleBlock{"Предназначение 40-60", "Destiny 40-60"},
    LocaleBlock{"Предназначение общее", "Destiny Common"},
    LocaleBlock{"Деньги", "Money"},
    LocaleBlock{"Программы", "Programs"},
    LocaleBlock{"Программа ", "Program "},
    LocaleBlock{"Сексуальность", "Sexiness"},
    LocaleBlock{"Прошлая жизнь", "Previous life"},
    LocaleBlock{"Родители", "Parents"},
    LocaleBlock{"Дети", "Kids"},
    LocaleBlock{"Отношения", "Relationship"},
    LocaleBlock{"Здоровье", "Health"},
    LocaleBlock{"Руководство по жизни", "Life Guidance"},
    LocaleBlock{"Важно", "Important"},
    LocaleBlock{"Ссора с родителями (муж.)", "Possible insult against parents (men)"},
    LocaleBlock{"Ссора с родителями (жен.)", "Possible insult against parents (women)"},
    LocaleBlock{"Обида на родителей", "Resentment against parents"},
    LocaleBlock{"Для достижения успеха", "To become successful"},
  }
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
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('money important') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('to become successful') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('relationships') ON CONFLICT DO NOTHING;`)
  // tx.MustExec(`INSERT INTO predictionType(name) VALUES('parents') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('parents IMPORTANT') ON CONFLICT DO NOTHING;`)
  // tx.MustExec(`INSERT INTO predictionType(name) VALUES('children') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('destiny') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('destiny common') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('past life') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('programms') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('life guide') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('health') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('health recommendation') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('year forecast') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('19-7') ON CONFLICT DO NOTHING;`)
  // tx.MustExec(`INSERT INTO predictionType(name) VALUES('sexiness') ON CONFLICT DO NOTHING;`)
  // tx.MustExec(`INSERT INTO predictionType(name) VALUES('lessons from children') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('children IMPORTANT') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('exception 22-7') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('resentment against parents') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('divorce') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('parent programms') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('information about destiny userfulinformation') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('information in personal features userfulinformation') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionLang(name) VALUES('russian') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionLang(name) VALUES('english') ON CONFLICT DO NOTHING;`)
  tx.Commit()

  router := routing.New()
  router.Use(
		access.Logger(log.Printf),
		slash.Remover(fasthttp.StatusMovedPermanently),
		fault.Recovery(log.Printf),
	)
  router.Use(func(req *routing.Context) error {
		// origin := string(req.Request.Header.Peek("Origin"))
		req.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
    req.Response.Header.Set("Access-Control-Allow-Credentials", "true")
    req.Response.Header.SetBytesV("Access-Control-Allow-Origin", req.Request.Header.Peek("Origin"))
		if err := req.Next(); err != nil {
			if httpError, ok := err.(routing.HTTPError); ok {
				req.Response.SetStatusCode(httpError.StatusCode())
			} else {
				req.Response.SetStatusCode(http.StatusInternalServerError)
			}
			req.SetContentType("application/json; charset=UTF-8")
			req.SetBody([]byte("lol"))
		}
		return nil
	})

  api := router.Group("/api/v1")
  // api.Use(content.TypeNegotiator(content.JSON))
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
    email := c.PostForm("email")
    password := c.PostForm("password")
    user := User{}
    err := db.Get(&user, "SELECT * FROM userProfile WHERE email = $1", email)
    if err != nil {
      return c.Write(`{"ok": false, "error": "Can't parse users table", "data": null}`)
    }
    if user.Id <= 0 {
      return c.Write(`{"ok": false, "error": "User not found", "data": null}`)
    }
    if user.Password.String != password {
      return c.Write(`{"ok": false, "error": "E-mail or password incorrect", "data": null}`)
    }
    return c.Write(fmt.Sprintf(`{"ok": true, "error": null, "data": {"email": "%s", "firstname": "%s", "lastname": "%s", "birthdate": "%s"}}`, user.Email.String, user.Firstname.String, user.Lastname.String, user.Birthdate.String))
  })
  api.Post("/register", func(c *routing.Context) error {
    email := c.PostForm("email")
    password := c.PostForm("password")
    firstname := c.PostForm("firstname")
    lastname := c.PostForm("lastname")
    birthdate := c.PostForm("birthdate")
    gender := c.PostForm("gender")
    userExists := User{}
    err := db.Get(&userExists, "SELECT * FROM userProfile WHERE email = $1", email)
    if err == nil || userExists.Id > 0 {
      return c.Write(`{"ok": false, "error": "User already exists", "data": null}`)
    }
    tx := db.MustBegin()
    tx.MustExec(`INSERT INTO userProfile(email,password,firstname,lastname,birthdate,gender) VALUES($1,$2,$3,$4,$5,$6)`, email, password, firstname, lastname, birthdate, gender)
    tx.Commit()
    return c.Write(`{"ok": true, "error": null, "data": "User registered"}`)
  })

  api.Get("/check/new/<input>", func(c *routing.Context) error {
    locale := newLocale()
    // blocks := []Block{}
    input := c.Param("input")
    timestamp := input[:len(input) - 4]
    gender := string(input[len(input) - 4])
    personal := string(input[len(input) - 3])
    languageShort := input[len(input) - 2:]
    if timestamp == "" {
      marshalled, _ := json.Marshal(Response{false, "No timestamp provided", nil})
      return c.Write(marshalled)
    }
    combo := countBD(timestamp)
    if combo == nil {
      marshalled, _ := json.Marshal(Response{false, "Timestamp corrupted", nil})
      return c.Write(marshalled)
    }
    fc := setAllCombosNew(combo)
    // ------ PERSONAL FEATURES BEGIN --------
    personalFeaturesBlocks := []Block{}
    personalFeaturesBlocks = append(personalFeaturesBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A), 1, languageShort, gender, personal))
    if languageShort == "ru" { personalFeaturesBlocks[len(personalFeaturesBlocks) - 1].Title = locale.PersonalFeaturesPos.Ru } else { personalFeaturesBlocks[len(personalFeaturesBlocks) - 1].Title = locale.PersonalFeaturesPos.En }
    personalFeaturesBlocks[len(personalFeaturesBlocks) - 1].Type = "info"
    if fc.A != fc.B {
      personalFeaturesBlocks[len(personalFeaturesBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", personalFeaturesBlocks[len(personalFeaturesBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.B), 1, languageShort, gender, personal).Content)
    }
    personalFeaturesBlocks = append(personalFeaturesBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A), 2, languageShort, gender, personal))
    if languageShort == "ru" { personalFeaturesBlocks[len(personalFeaturesBlocks) - 1].Title = locale.PersonalFeaturesNeg.Ru } else { personalFeaturesBlocks[len(personalFeaturesBlocks) - 1].Title = locale.PersonalFeaturesNeg.En }
    personalFeaturesBlocks[len(personalFeaturesBlocks) - 1].Type = "info"
    if fc.A != fc.B {
      personalFeaturesBlocks[len(personalFeaturesBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", personalFeaturesBlocks[len(personalFeaturesBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.B), 2, languageShort, gender, personal).Content)
    }
    personalFeaturesBlocks = append(personalFeaturesBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.E), 3, languageShort, gender, personal))
    if languageShort == "ru" { personalFeaturesBlocks[len(personalFeaturesBlocks) - 1].Title = locale.PersonalFeaturesSoc.Ru } else { personalFeaturesBlocks[len(personalFeaturesBlocks) - 1].Title = locale.PersonalFeaturesSoc.En }
    personalFeaturesBlocks[len(personalFeaturesBlocks) - 1].Type = "info"
    personalfeatures := Prediction{}
    personalfeatures.Title = "Personal Features"
    personalfeatures.ImageName = "personal_features"
    personalfeatures.Blocks = personalFeaturesBlocks
    personalfeatures.BlockType = "default"
    // ------ PERSONAL FEATURES END   --------
    // ------ DESTINY BEGIN           --------
    destinyBlocks := []Block{}
    destinyBlocks = append(destinyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.H), 8, languageShort, gender, personal))
    if languageShort == "ru" { destinyBlocks[len(destinyBlocks) - 1].Title = locale.Destiny20.Ru } else { destinyBlocks[len(destinyBlocks) - 1].Title = locale.Destiny20.En }
    // destinyBlocks[len(destinyBlocks) - 1].Title = "destiny 20-40"
    destinyBlocks[len(destinyBlocks) - 1].Type = "info"
    if fc.H != fc.J {
      // destinyBlocks = append(destinyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.J), 8, languageShort, gender, personal))
      destinyBlocks[len(destinyBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", destinyBlocks[len(destinyBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.J), 8, languageShort, gender, personal).Content)
    }
    if fc.H != fc.M && fc.J != fc.M {
      // destinyBlocks = append(destinyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.M), 8, languageShort, gender, personal))
      destinyBlocks[len(destinyBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", destinyBlocks[len(destinyBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.M), 8, languageShort, gender, personal).Content)
    }
    // destinyBlocks = append(destinyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.N), 8, languageShort, gender, personal))
    destinyBlocks = append(destinyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.N), 8, languageShort, gender, personal))
    if languageShort == "ru" { destinyBlocks[len(destinyBlocks) - 1].Title = locale.Destiny40.Ru } else { destinyBlocks[len(destinyBlocks) - 1].Title = locale.Destiny40.En }
    // destinyBlocks[len(destinyBlocks) - 1].Title = "destiny 40-60"
    destinyBlocks[len(destinyBlocks) - 1].Type = "info"
    if fc.T != fc.N {
      // destinyBlocks = append(destinyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.T), 8, languageShort, gender, personal))
      destinyBlocks[len(destinyBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", destinyBlocks[len(destinyBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.T), 8, languageShort, gender, personal).Content)
    }
    if fc.Z != fc.N && fc.Z != fc.T {
      // destinyBlocks = append(destinyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Z), 8, languageShort, gender, personal))
      destinyBlocks[len(destinyBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", destinyBlocks[len(destinyBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Z), 8, languageShort, gender, personal).Content)
    }
    destinyBlocks = append(destinyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.S), 220, languageShort, gender, personal))
    if languageShort == "ru" { destinyBlocks[len(destinyBlocks) - 1].Title = locale.DestinyCommon.Ru } else { destinyBlocks[len(destinyBlocks) - 1].Title = locale.DestinyCommon.En }
    // destinyBlocks[len(destinyBlocks) - 1].Title = "destiny common"
    destinyBlocks[len(destinyBlocks) - 1].Type = "info"
    destiny := Prediction{}
    destiny.Title = "Destiny"
    destiny.ImageName = "destiny"
    destiny.Blocks = destinyBlocks
    destiny.BlockType = "default"
    // ------ DESTINY END            --------
    // ------ MONEY BEGIN            --------
    moneyBlocks := []Block{}
    moneyBlocks = append(moneyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.X2), 4, languageShort, gender, personal))
    if languageShort == "ru" { moneyBlocks[len(moneyBlocks) - 1].Title = locale.Money.Ru } else { moneyBlocks[len(moneyBlocks) - 1].Title = locale.Money.En }
    // moneyBlocks[len(moneyBlocks) - 1].Title = "money"
    moneyBlocks[len(moneyBlocks) - 1].Type = "info"
    moneyBlocks = append(moneyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.X), 214, languageShort, gender, personal))
    if languageShort == "ru" { moneyBlocks[len(moneyBlocks) - 1].Title = locale.ToBecomeSuccessful.Ru } else { moneyBlocks[len(moneyBlocks) - 1].Title = locale.ToBecomeSuccessful.En }
    // moneyBlocks[len(moneyBlocks) - 1].Title = "to become successful"
    moneyBlocks[len(moneyBlocks) - 1].Type = "info"
    if fc.X != fc.C {
      // moneyBlocks = append(moneyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.C), 214, languageShort, gender, personal))
      moneyBlocks[len(moneyBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", moneyBlocks[len(moneyBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.C), 214, languageShort, gender, personal).Content)
    }
    if fc.X != fc.C1 && fc.C != fc.C1 {
      // moneyBlocks = append(moneyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.C1), 214, languageShort, gender, personal))
      moneyBlocks[len(moneyBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", moneyBlocks[len(moneyBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.C1), 214, languageShort, gender, personal).Content)
    }
    if fc.X != fc.C2 && fc.C != fc.C2 && fc.C1 != fc.C2 {
      // moneyBlocks = append(moneyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.C2), 214, languageShort, gender, personal))
      moneyBlocks[len(moneyBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", moneyBlocks[len(moneyBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.C2), 214, languageShort, gender, personal).Content)
    }
    toCheck := [][]int{
      {fc.C, fc.C1},
      {fc.C1, fc.X},
      {fc.X, fc.X2},
      {fc.X2, fc.X},
      {fc.X2, fc.C1},
      {fc.C1, fc.X2},
    }
    fmt.Println(toCheck)
    important := false
    if checkAnswers(toCheck, []int{19,7}) {
      if !important {
        important = true
      }
      moneyBlocks = append(moneyBlocks, getAnswerFromTable(db, "'1'", 227, languageShort, gender, personal))
      if languageShort == "ru" { moneyBlocks[len(moneyBlocks) - 1].Title = locale.Important.Ru } else { moneyBlocks[len(moneyBlocks) - 1].Title = locale.Important.En }
      // moneyBlocks[len(moneyBlocks) - 1].Title = "important"
      moneyBlocks[len(moneyBlocks) - 1].Type = "expandable"
    }
    toCheck = [][]int{
      {fc.C, fc.C1, fc.C2},
      {fc.C1, fc.C, fc.C2},
    }
    if checkAnswers(toCheck, []int{7,15,22}) {
      if !important {
        important = true
      }
      moneyBlocks = append(moneyBlocks, getAnswerFromTable(db, "'7-15-22'", 213, languageShort, gender, personal))
      if languageShort == "ru" { moneyBlocks[len(moneyBlocks) - 1].Title = locale.Important.Ru } else { moneyBlocks[len(moneyBlocks) - 1].Title = locale.Important.En }
      // moneyBlocks[len(moneyBlocks) - 1].Title = "important"
      moneyBlocks[len(moneyBlocks) - 1].Type = "expandable"
    }
    toCheck = [][]int{
      {fc.C, fc.C1, fc.C2},
    }
    if checkAnswers(toCheck, []int{8,9,17}, true) {
      if !important {
        important = true;
      }
      moneyBlocks = append(moneyBlocks, getAnswerFromTable(db, "'8-9-17'", 213, languageShort, gender, personal))
      if languageShort == "ru" { moneyBlocks[len(moneyBlocks) - 1].Title = locale.Important.Ru } else { moneyBlocks[len(moneyBlocks) - 1].Title = locale.Important.En }
      // moneyBlocks[len(moneyBlocks) - 1].Title = "important"
      moneyBlocks[len(moneyBlocks) - 1].Type = "expandable"
    }
    toCheck = [][]int{
      {fc.C, fc.C1, fc.C2},
      {fc.C1, fc.C, fc.C2},
    }
    if checkAnswers(toCheck, []int{8,14,22}) {
      if !important {
        important = true
      }
      moneyBlocks = append(moneyBlocks, getAnswerFromTable(db, "'8-14-22'", 213, languageShort, gender, personal))
      if languageShort == "ru" { moneyBlocks[len(moneyBlocks) - 1].Title = locale.Important.Ru } else { moneyBlocks[len(moneyBlocks) - 1].Title = locale.Important.En }
      // moneyBlocks[len(moneyBlocks) - 1].Title = "important"
      moneyBlocks[len(moneyBlocks) - 1].Type = "expandable"
    }
    toCheck = [][]int{
      {fc.A, fc.A2, fc.A1},
      {fc.B, fc.B2, fc.B1},
      {fc.F, fc.Y, fc.O},
      {fc.K, fc.G, fc.U},
      {fc.E, fc.E1, fc.E2},
      {fc.D1, fc.X1, fc.X},
      {fc.X, fc.X2, fc.C1},
      {fc.H, fc.J, fc.M},
      {fc.N, fc.T, fc.Z},
      {fc.A, fc.B, fc.L},
      {fc.A2, fc.B2, fc.L1},
      {fc.A1, fc.B1, fc.L2},
      {fc.A3, fc.B3, fc.L3},
      {fc.E, fc.E, fc.L4},
      {fc.D1, fc.C1, fc.L5},
      {fc.D, fc.C, fc.L6},
      {fc.D3, fc.C3, fc.E3},
    }
    fmt.Println(toCheck)
    if checkAnswers(toCheck, []int{5,14,19}, true) {
      moneyBlocks = append(moneyBlocks, getAnswerFromTable(db, "'5-14-19'", 213, languageShort, gender, personal))
      if languageShort == "ru" { moneyBlocks[len(moneyBlocks) - 1].Title = locale.Important.Ru } else { moneyBlocks[len(moneyBlocks) - 1].Title = locale.Important.En }
      // moneyBlocks[len(moneyBlocks) - 1].Title = "important"
      moneyBlocks[len(moneyBlocks) - 1].Type = "expandable"
    }
    toCheck = [][]int{
      {fc.A, fc.A2, fc.A1},
    }
    if checkAnswers(toCheck, []int{5,14,19}, true) {
      moneyBlocks = append(moneyBlocks, getAnswerFromTable(db, "'5-14-19'", 217, languageShort, gender, personal))
      if languageShort == "ru" { moneyBlocks[len(moneyBlocks) - 1].Title = locale.Important.Ru } else { moneyBlocks[len(moneyBlocks) - 1].Title = locale.Important.En }
      // moneyBlocks[len(moneyBlocks) - 1].Title = "important"
      moneyBlocks[len(moneyBlocks) - 1].Type = "expandable"
    }
    money := Prediction{}
    money.Title = "Money"
    money.ImageName = "money"
    money.Blocks = moneyBlocks
    money.BlockType = "default"
    // ------ MONEY END   --------
    // ------ PROGRAMS BEGIN ---------
    programsBlocks := []Block{}
    // toCheck = [][]int{{fc.A, fc.A1, fc.A2}}
    // if checkAnswers(toCheck, []int{5,14,19}, true) {
    //   programsBlocks = append(programsBlocks, getAnswerFromTable(db, "'5-14-19'", 101, languageShort, gender, personal))
    // }
    lessonsCount := 0
    lessonsBlocks := []Block{}
    toCheck = [][]int{{fc.C, fc.C2, fc.C1}}
    if checkAnswers(toCheck, []int{17,5,6}, true) {
      lessonsBlocks = append(lessonsBlocks, getAnswerFromTable(db, "'5-6-17'", 10, languageShort, gender, personal))
      if languageShort == "ru" { lessonsBlocks[len(lessonsBlocks) - 1].Title = fmt.Sprintf("%s %d", locale.Program.Ru, lessonsCount) } else { lessonsBlocks[len(lessonsBlocks) - 1].Title = fmt.Sprintf("%s %d", locale.Program.En, lessonsCount) }
      // lessonsBlocks[len(lessonsBlocks) - 1].Title = fmt.Sprintf("program %d", lessonsCount)
      lessonsBlocks[len(lessonsBlocks) - 1].Type = "expandable"
      lessonsCount = lessonsCount + 1
    }
    toCheck = [][]int{
      {fc.A, fc.A2, fc.A1},
      {fc.B, fc.B2, fc.B1},
      {fc.C, fc.C2, fc.C1},
      {fc.F, fc.Y, fc.O},
      {fc.K, fc.G, fc.U},
      {fc.E, fc.E1, fc.E2},
      {fc.D1, fc.X1, fc.X},
      {fc.X, fc.X2, fc.C1},
      {fc.H, fc.J, fc.M},
      {fc.N, fc.T, fc.Z},
      {fc.A, fc.B, fc.L},
      {fc.A2, fc.B2, fc.L1},
      {fc.A1, fc.B1, fc.L2},
      {fc.A3, fc.B3, fc.L3},
      {fc.E, fc.E, fc.L4},
      {fc.D1, fc.C1, fc.L5},
      {fc.D, fc.C, fc.L6},
      {fc.D3, fc.C3, fc.E3},
    }
    toCompare := [][]int{
      {13,6,19},{9,20,11},{5,17,12},{10,11,19},
      {9,14,5},{10,8,16},{15,22,7},{11,16,22},
      {13,7,20},{11,3,19},{9,6,15},{9,6,17},
      {11,9,16},{21,7,14},{13,18,5},{18,11,11},
      {17,3,20},{20,3,10},{17,22,5},{9,8,17},
      {13,21,8},{9,10,19},{20,11,4},{18,10,10},
      {18,6,12},{18,7,7},{22,4,9},{22,4,8},
      {22,22,8},{18,8,8},{15,21,6},{19,8,7},
      {14,22,8},{10,15,5},{10,5,22},{18,5,5},
    }
    for _, lesson := range toCompare {
      if checkAnswers(toCheck, lesson, true) {
        fmt.Println(lesson)
        sort.Ints(lesson)
        answer := fmt.Sprintf("'%d-%d-%d'", lesson[0], lesson[1], lesson[2])
        fmt.Println(answer)
        if lesson[0] == 8 && lesson[1] == 10 && lesson[2] == 16 {
          fmt.Println(toCheck[0])
          fmt.Println(toCheck[1])
          fmt.Println(toCheck[2])
          fmt.Println(toCheck[3])
          fmt.Println(toCheck[4])
          fmt.Println(toCheck[5])
          fmt.Println(toCheck[6])
          fmt.Println(toCheck[7])
          fmt.Println(toCheck[8])
          fmt.Println(toCheck[9])
          fmt.Println(toCheck[10])
          fmt.Println(toCheck[11])
          fmt.Println(toCheck[12])
          fmt.Println(toCheck[13])
          fmt.Println(toCheck[14])
          fmt.Println(toCheck[15])
          fmt.Println(toCheck[16])
          fmt.Println(toCheck[17])
        }
        lessonsBlocks = append(lessonsBlocks, getAnswerFromTable(db, answer, 10, languageShort, gender, personal))
        if languageShort == "ru" { lessonsBlocks[len(lessonsBlocks) - 1].Title = fmt.Sprintf("%s %d", locale.Program.Ru, lessonsCount) } else { lessonsBlocks[len(lessonsBlocks) - 1].Title = fmt.Sprintf("%s %d", locale.Program.En, lessonsCount) }
        // lessonsBlocks[len(lessonsBlocks) - 1].Title = fmt.Sprintf("program %d", lessonsCount)
        lessonsBlocks[len(lessonsBlocks) - 1].Type = "expandable"
        lessonsCount = lessonsCount + 1
      }
    }
    for _, lblock := range lessonsBlocks {
      programsBlocks = append(programsBlocks, lblock)
    }
    programs := Prediction{}
    programs.Title = "Programs"
    programs.ImageName = "programs"
    programs.Blocks = programsBlocks
    programs.BlockType = "default"
    // ------ PROGRAMS END --------
    // ------ SEXINESS BEGIN --------
    sexinessBlocks := []Block{}
    sexual := fmt.Sprintf("'%d-%d-%d'", fc.E, fc.E1, fc.E2)
    sexinessBlocks = append(sexinessBlocks, getAnswerFromTable(db, sexual, 76, languageShort, gender, personal))
    if languageShort == "ru" { sexinessBlocks[len(sexinessBlocks) - 1].Title = locale.Sexiness.Ru } else { sexinessBlocks[len(sexinessBlocks) - 1].Title = locale.Sexiness.En }
    // sexinessBlocks[len(sexinessBlocks) - 1].Title = "sexiness"
    sexinessBlocks[len(sexinessBlocks) - 1].Type = "info"
    sexiness := Prediction{}
    sexiness.Title = "Sexiness"
    sexiness.ImageName = "sexiness"
    sexiness.Blocks = sexinessBlocks
    sexiness.BlockType = "default"
    // ------ SEXINESS END --------
    // ------ PAST LIFE BEGIN --------
    pastLifeBlocks := []Block{}
    mainLesson := fmt.Sprintf("'%d-%d-%d'", fc.D1, fc.D2, fc.D)
    pastLifeBlocks = append(pastLifeBlocks, getAnswerFromTable(db, mainLesson, 9, languageShort, gender, personal))
    if languageShort == "ru" { pastLifeBlocks[len(pastLifeBlocks) - 1].Title = locale.PastLife.Ru } else { pastLifeBlocks[len(pastLifeBlocks) - 1].Title = locale.PastLife.En }
    // pastLifeBlocks[len(pastLifeBlocks) - 1].Title = "previous life"
    pastLifeBlocks[len(pastLifeBlocks) - 1].Type = "info"
    pastlife := Prediction{}
    pastlife.Title = "Past Life"
    pastlife.ImageName = "previous_life"
    pastlife.Blocks = pastLifeBlocks
    pastlife.BlockType = "default"
    // ------ PAST LIFE END --------
    // ------ PARENTS  BEGIN --------
    parentsBlocks := []Block{}
    // if fc.D1 == 21 {
    //   parentsBlocks = append(parentsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", 1), 217, languageShort, gender, personal))
    //   parentsBlocks[len(parentsBlocks) - 1].Title = "Important"
    //   parentsBlocks[len(parentsBlocks) - 1].Type = "expandable"
    // }
    toCheck = [][]int{
      {fc.A, fc.A2, fc.A1},
    }
    if checkAnswers(toCheck, []int{17,5,6}, true) {
      parentsBlocks = append(parentsBlocks, getAnswerFromTable(db, "'5-6-17'", 217, languageShort, gender, personal))
      if languageShort == "ru" { parentsBlocks[len(parentsBlocks) - 1].Title = locale.Important.Ru } else { parentsBlocks[len(parentsBlocks) - 1].Title = locale.Important.En }
      // parentsBlocks[len(parentsBlocks) - 1].Title = "Important"
      parentsBlocks[len(parentsBlocks) - 1].Type = "expandable"
    }
    if checkAnswers(toCheck, []int{13,8,21}, true) {
      parentsBlocks = append(parentsBlocks, getAnswerFromTable(db, "'8-13-21'", 217, languageShort, gender, personal))
      if languageShort == "ru" { parentsBlocks[len(parentsBlocks) - 1].Title = locale.Important.Ru } else { parentsBlocks[len(parentsBlocks) - 1].Title = locale.Important.En }
      // parentsBlocks[len(parentsBlocks) - 1].Title = "Important"
      parentsBlocks[len(parentsBlocks) - 1].Type = "expandable"
    }
    parentsBlocks = append(parentsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.F), 234, languageShort, gender, personal))
    if languageShort == "ru" { parentsBlocks[len(parentsBlocks) - 1].Title = locale.InsultParentsMen.Ru } else { parentsBlocks[len(parentsBlocks) - 1].Title = locale.InsultParentsMen.En }
      // parentsBlocks[len(parentsBlocks) - 1].Title = "Possible insult on parents (men)"
    parentsBlocks[len(parentsBlocks) - 1].Type = "info"
    // if fc.F != fc.Y {
      // parentsBlocks = append(parentsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Y), 234, languageShort, gender, personal))
      parentsBlocks[len(parentsBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", parentsBlocks[len(parentsBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Y), 234, languageShort, gender, personal).Content)
    // }
    // if fc.F != fc.O && fc.Y != fc.O {
      // parentsBlocks = append(parentsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.O), 234, languageShort, gender, personal))
      parentsBlocks[len(parentsBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", parentsBlocks[len(parentsBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.O), 234, languageShort, gender, personal).Content)
    // }
    parentsBlocks = append(parentsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.G), 234, languageShort, gender, personal))
    if languageShort == "ru" { parentsBlocks[len(parentsBlocks) - 1].Title = locale.InsultParentsWomen.Ru } else { parentsBlocks[len(parentsBlocks) - 1].Title = locale.InsultParentsWomen.En }
    // parentsBlocks[len(parentsBlocks) - 1].Title = "Possible insult on parents (women)"
    parentsBlocks[len(parentsBlocks) - 1].Type = "info"
    // if fc.G != fc.K {
      // parentsBlocks = append(parentsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K), 234, languageShort, gender, personal))
      parentsBlocks[len(parentsBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", parentsBlocks[len(parentsBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K), 234, languageShort, gender, personal).Content)
    // }
    // if fc.G != fc.U && fc.K != fc.U {
      // parentsBlocks = append(parentsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.U), 234, languageShort, gender, personal))
      parentsBlocks[len(parentsBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", parentsBlocks[len(parentsBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.U), 234, languageShort, gender, personal).Content)
    // }
    parentsBlocks = append(parentsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A), 232, languageShort, gender, personal))
    if languageShort == "ru" { parentsBlocks[len(parentsBlocks) - 1].Title = locale.Resentment.Ru } else { parentsBlocks[len(parentsBlocks) - 1].Title = locale.Resentment.En }
    // parentsBlocks[len(parentsBlocks) - 1].Title = "Resentment against parents"
    parentsBlocks[len(parentsBlocks) - 1].Type = "info"
    parentsBlocks[len(parentsBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", parentsBlocks[len(parentsBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A2), 232, languageShort, gender, personal).Content)
    parentsBlocks[len(parentsBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", parentsBlocks[len(parentsBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A1), 232, languageShort, gender, personal).Content)
    parents := Prediction{}
    parents.Title = "Parents"
    parents.ImageName = "parents"
    parents.Blocks = parentsBlocks
    parents.BlockType = "default"
    // ------ PARENTS END --------
    // ------ KIDS BEGIN --------
    // if gender == "m" {
    //   blocks = append(blocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.F), 234, languageShort, gender, personal))
    //   if fc.F != fc.Y {
    //     blocks = append(blocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Y), 234, languageShort, gender, personal))
    //   }
    //   if fc.F != fc.O && fc.Y != fc.O {
    //     blocks = append(blocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.O), 234, languageShort, gender, personal))
    //   }
    // } else {
    //   blocks = append(blocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.G), 234, languageShort, gender, personal))
    //   if fc.G != fc.K {
    //     blocks = append(blocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K), 234, languageShort, gender, personal))
    //   }
    //   if fc.G != fc.U && fc.K != fc.U {
    //     blocks = append(blocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.U), 234, languageShort, gender, personal))
    //   }
    // }
    kidsBlocks := []Block{}
    kidsBlocks = append(kidsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A), 66, languageShort, gender, personal))
    if languageShort == "ru" { kidsBlocks[len(kidsBlocks) - 1].Title = locale.Kids.Ru } else { kidsBlocks[len(kidsBlocks) - 1].Title = locale.Kids.En }
    // kidsBlocks[len(kidsBlocks) - 1].Title = "Children"
    kidsBlocks[len(kidsBlocks) - 1].Type = "info"
    if fc.A != fc.A2 {
      // kidsBlocks = append(kidsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A2), 66, languageShort, gender, personal))
      kidsBlocks[len(kidsBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", kidsBlocks[len(kidsBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A2), 66, languageShort, gender, personal).Content)
    }
    if fc.A != fc.A1 && fc.A2 != fc.A1 {
      // kidsBlocks = append(kidsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A1), 66, languageShort, gender, personal))
      kidsBlocks[len(kidsBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", kidsBlocks[len(kidsBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A1), 66, languageShort, gender, personal).Content)
    }
    toCheck = [][]int{
      {fc.A, fc.A1, fc.A2},
      {fc.A1, fc.A, fc.A2},
    }
    // if checkAnswers(toCheck, []int{6,17,5}) {
    //   // kidsBlocks = append(kidsBlocks, getAnswerFromTable(db, "'6-17-5'", 232, languageShort, gender, personal))
    //   kidsBlocks[len(kidsBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", kidsBlocks[len(kidsBlocks) - 1].Content, getAnswerFromTable(db, "'5-6-17'", 230, languageShort, gender, personal).Content)
    // }
    if checkAnswers(toCheck, []int{7,15,22}) {
      // kidsBlocks = append(kidsBlocks, getAnswerFromTable(db, "'7-15-22'", 232, languageShort, gender, personal))
      kidsBlocks = append(kidsBlocks, getAnswerFromTable(db, "'7-15-22'", 230, languageShort, gender, personal))
      if languageShort == "ru" { kidsBlocks[len(kidsBlocks) - 1].Title = locale.Important.Ru } else { kidsBlocks[len(kidsBlocks) - 1].Title = locale.Important.En }
      // kidsBlocks[len(kidsBlocks) - 1].Title = "Important"
      kidsBlocks[len(kidsBlocks) - 1].Type = "expandable"
    }
    if checkAnswers(toCheck, []int{8,9,17}) {
      // kidsBlocks = append(kidsBlocks, getAnswerFromTable(db, "'8-9-17'", 232, languageShort, gender, personal))
      kidsBlocks = append(kidsBlocks, getAnswerFromTable(db, "'8-9-17'", 230, languageShort, gender, personal))
      if languageShort == "ru" { kidsBlocks[len(kidsBlocks) - 1].Title = locale.Important.Ru } else { kidsBlocks[len(kidsBlocks) - 1].Title = locale.Important.En }
      // kidsBlocks[len(kidsBlocks) - 1].Title = "Important"
      kidsBlocks[len(kidsBlocks) - 1].Type = "expandable"
    }
    // if checkAnswers(toCheck, []int{8,13,21}) {
    //   // kidsBlocks = append(kidsBlocks, getAnswerFromTable(db, "'8-13-21'", 232, languageShort, gender, personal))
    //   kidsBlocks[len(kidsBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", kidsBlocks[len(kidsBlocks) - 1].Content, getAnswerFromTable(db, "'8-13-21'", 230, languageShort, gender, personal).Content)
    // }
    if checkAnswers(toCheck, []int{6,12,18}) {
      // kidsBlocks = append(kidsBlocks, getAnswerFromTable(db, "'6-12-18'", 232, languageShort, gender, personal))
      kidsBlocks = append(kidsBlocks, getAnswerFromTable(db, "'6-12-18'", 230, languageShort, gender, personal))
      if languageShort == "ru" { kidsBlocks[len(kidsBlocks) - 1].Title = locale.Important.Ru } else { kidsBlocks[len(kidsBlocks) - 1].Title = locale.Important.En }
      // kidsBlocks[len(kidsBlocks) - 1].Title = "Important"
      kidsBlocks[len(kidsBlocks) - 1].Type = "expandable"
    }
    kids := Prediction{}
    kids.Title = "Children"
    kids.ImageName = "children"
    kids.Blocks = kidsBlocks
    kids.BlockType = "default"
    // ------ KIDS END --------
    // ------ RELATIONSHIPS BEGIN --------
    relationshipsBlocks := []Block{}
    relationshipsBlocks = append(relationshipsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.X1), 5, languageShort, gender, personal))
    if languageShort == "ru" { relationshipsBlocks[len(relationshipsBlocks) - 1].Title = locale.Relationship.Ru } else { relationshipsBlocks[len(relationshipsBlocks) - 1].Title = locale.Relationship.En }
    // relationshipsBlocks[len(relationshipsBlocks) - 1].Title = "Relationships"
    relationshipsBlocks[len(relationshipsBlocks) - 1].Type = "info"
    if fc.X1 != fc.D1 {
      // relationshipsBlocks = append(relationshipsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.D1), 5, languageShort, gender, personal))
      relationshipsBlocks[len(relationshipsBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", relationshipsBlocks[len(relationshipsBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.D1), 5, languageShort, gender, personal).Content)
    }
    if fc.X1 != fc.X && fc.D1 != fc.X {
      // relationshipsBlocks = append(relationshipsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.X), 5, languageShort, gender, personal))
      relationshipsBlocks[len(relationshipsBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", relationshipsBlocks[len(relationshipsBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.X), 5, languageShort, gender, personal).Content)
    }
    toCheck = [][]int{
      {fc.X, fc.X1},
      {fc.X1, fc.X},
    }
    if checkAnswers(toCheck, []int{22,7}) {
      relationshipsBlocks = append(relationshipsBlocks, getAnswerFromTable(db, "'1'", 231, languageShort, gender, personal))
      if languageShort == "ru" { relationshipsBlocks[len(relationshipsBlocks) - 1].Title = locale.Important.Ru } else { relationshipsBlocks[len(relationshipsBlocks) - 1].Title = locale.Important.En }
      // relationshipsBlocks[len(relationshipsBlocks) - 1].Title = "Important"
      relationshipsBlocks[len(relationshipsBlocks) - 1].Type = "expandable"
    }
    relationships := Prediction{}
    relationships.Title = "Relationships"
    relationships.ImageName = "relationships"
    relationships.Blocks = relationshipsBlocks
    relationships.BlockType = "default"
    // ------ RELATIONSHIPS END --------
    // ------ HEALTH BEGIN --------
    healthBlocks := []Block{}
    healthBlocks = append(healthBlocks, getAnswerFromTable(db, "1", 12, languageShort, gender, personal))
    healthBlocks[len(healthBlocks) - 1].Title = "The brain, hair, upper part of the skull."
    healthBlocks[len(healthBlocks) - 1].Type = "info"
    // healthBlocks = append(healthBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A), 225, languageShort, gender, personal))
    // healthBlocks = append(healthBlocks, getAnswerFromTable(db, "1", 12, languageShort, gender, personal))
    health1 := prepareArray([]int{fc.A, fc.B, fc.L})
    for _, item := range health1 {
      // healthBlocks = append(healthBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal))
      healthBlocks[len(healthBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", healthBlocks[len(healthBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal).Content)
    }
    healthBlocks = append(healthBlocks, getAnswerFromTable(db, "2", 12, languageShort, gender, personal))
    healthBlocks[len(healthBlocks) - 1].Title = "Occipital and temporal lobes of the brain, eyes, ears, nose, face, upper jaw, upper jaw teeth, optic nerve, cerebral cortex."
    healthBlocks[len(healthBlocks) - 1].Type = "info"
    health2 := prepareArray([]int{fc.A2, fc.B2, fc.L1})
    for _, item := range health2 {
      // healthBlocks = append(healthBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal))
      healthBlocks[len(healthBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", healthBlocks[len(healthBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal).Content)
    }
    healthBlocks = append(healthBlocks, getAnswerFromTable(db, "3", 12, languageShort, gender, personal))
    healthBlocks[len(healthBlocks) - 1].Title = "Thyroid gland, trachea, bronchi, throat, vocal cords, shoulders, arms, seventh cervical vertebra, all cervical vertebrae, lower jaw, lower jaw teeth."
    healthBlocks[len(healthBlocks) - 1].Type = "info"
    health3 := prepareArray([]int{fc.A1, fc.B1, fc.L2})
    for _, item := range health3 {
      // healthBlocks = append(healthBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal))
      healthBlocks[len(healthBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", healthBlocks[len(healthBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal).Content)
    }
    healthBlocks = append(healthBlocks, getAnswerFromTable(db, "4", 12, languageShort, gender, personal))
    healthBlocks[len(healthBlocks) - 1].Title = "Heart, circulatory system, respiratory system, lungs, bronchi, thoracic spine, ribs, shoulder scapular area, chest."
    healthBlocks[len(healthBlocks) - 1].Type = "info"
    health4 := prepareArray([]int{fc.A3, fc.B3, fc.L3})
    for _, item := range health4 {
      // healthBlocks = append(healthBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", item), 12, languageShort, gender, personal))
      healthBlocks[len(healthBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", healthBlocks[len(healthBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal).Content)
    }
    healthBlocks = append(healthBlocks, getAnswerFromTable(db, "5", 12, languageShort, gender, personal))
    healthBlocks[len(healthBlocks) - 1].Title = "Gastrointestinal tract, abdominal organs, pancreas, spleen, liver, gallbladder, small intestine, central part of the spine."
    healthBlocks[len(healthBlocks) - 1].Type = "info"
    health5 := prepareArray([]int{fc.E, fc.L4})
    for _, item := range health5 {
      // healthBlocks = append(healthBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal))
      healthBlocks[len(healthBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", healthBlocks[len(healthBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal).Content)
    }
    healthBlocks = append(healthBlocks, getAnswerFromTable(db, "6", 12, languageShort, gender, personal))
    healthBlocks[len(healthBlocks) - 1].Title = "drenal glands, uterus and ovaries, kidneys, intestines, prostate gland in men, lumbar spinal column."
    healthBlocks[len(healthBlocks) - 1].Type = "info"
    health6 := prepareArray([]int{fc.D1, fc.C1, fc.L5})
    for _, item := range health6 {
      // healthBlocks = append(healthBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal))
      healthBlocks[len(healthBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", healthBlocks[len(healthBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal).Content)
    }
    healthBlocks = append(healthBlocks, getAnswerFromTable(db, "7", 12, languageShort, gender, personal))
    healthBlocks[len(healthBlocks) - 1].Title = "Urogenital system, lower limbs, large intestine, tailbone, sacrum, legs."
    healthBlocks[len(healthBlocks) - 1].Type = "info"
    health7 := prepareArray([]int{fc.D, fc.C, fc.L6})
    for _, item := range health7 {
      // healthBlocks = append(healthBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal))
      healthBlocks[len(healthBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", healthBlocks[len(healthBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal).Content)
    }
    healthBlocks = append(healthBlocks, getAnswerFromTable(db, "8", 12, languageShort, gender, personal))
    healthBlocks[len(healthBlocks) - 1].Title = "Circulatory system, nervous system, lymphatic system, immune system, those organs that are found throughout the body, general failure of the body."
    healthBlocks[len(healthBlocks) - 1].Type = "info"
    health8 := prepareArray([]int{fc.D3, fc.C3, fc.E3})
    for _, item := range health8 {
      // healthBlocks = append(healthBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal))
      healthBlocks[len(healthBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", healthBlocks[len(healthBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal).Content)
    }
    health := Prediction{}
    health.Title = "Health"
    health.ImageName = "health"
    health.Blocks = healthBlocks
    health.BlockType = "health"
    // ------ HEALTH END --------
    // ------ LIFE GUIDE BEGIN --------
    lifeGuideBlocks := []Block{}
    lifeGuideBlocks = append(lifeGuideBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A), 11, languageShort, gender, personal))
    if languageShort == "ru" { lifeGuideBlocks[len(lifeGuideBlocks) - 1].Title = locale.LifeGuide.Ru } else { lifeGuideBlocks[len(lifeGuideBlocks) - 1].Title = locale.LifeGuide.En }
    // lifeGuideBlocks[len(lifeGuideBlocks) - 1].Title = "Life Guidance"
    lifeGuideBlocks[len(lifeGuideBlocks) - 1].Type = "info"
    if fc.A != fc.B {
      // lifeGuideBlocks = append(lifeGuideBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.B), 11, languageShort, gender, personal))
      lifeGuideBlocks[len(lifeGuideBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", lifeGuideBlocks[len(lifeGuideBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.B), 11, languageShort, gender, personal).Content)
    }
    if fc.A != fc.E && fc.B != fc.E {
      // lifeGuideBlocks = append(lifeGuideBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.E), 11, languageShort, gender, personal))
      lifeGuideBlocks[len(lifeGuideBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", lifeGuideBlocks[len(lifeGuideBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.E), 11, languageShort, gender, personal).Content)
    }
    lifeguide := Prediction{}
    lifeguide.Title = "Life Guidance"
    lifeguide.ImageName = "life_guidance"
    lifeguide.Blocks = lifeGuideBlocks
    lifeguide.BlockType = "default"
    // ------ LIFE GUIDE END --------
    // ------ YEAR FORECAST BEGIN --------
    forecastBlocks := []Block{}
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.B), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "20"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.D), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.H), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.D), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "60"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.B), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.H), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.B8), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "21-22"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.G3), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.O1), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.G3), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "61-62"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.B8), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.O1), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.B7), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "22.5-23"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.G2), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.O2), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.G2), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "62.5-63"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.B7), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.O2), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K2), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "23.5-24"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.G4), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.O3), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.G4), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "63.5-64"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K2), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.O3), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K1), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "25"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.G1), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.X3), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.G1), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "65"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K1), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.X3), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K4), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "26-27"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.G6), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.X4), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.G6), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "66-67"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K4), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.X4), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K3), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "27.5-28"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.G5), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.X5), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.G5), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "67.5-68"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K3), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.X5), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K5), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "28.5-29"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.G7), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.X6), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.G7), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "68.5-69"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K5), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.X6), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "30"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.G), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.U), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.G), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "70"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.U), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K8), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "31-32"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.T3), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.X8), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.T3), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "71-72"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K8), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.X8), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K7), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "32.5-33"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.T2), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.E4), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.T2), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "72.5-73"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K7), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.E4), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K6), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "33.5-34"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.T4), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.E5), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.T4), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "73.5-74"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K6), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.E5), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.C4), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "35"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.T1), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.E6), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.T1), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "75"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.C4), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.E6), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.C6), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "36-37"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.T5), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.E7), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.T5), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "76-77"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.C6), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.E7), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.C5), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "37.5-38"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.T5), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.E8), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.T5), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "77.5-78"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.C5), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.E8), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.C7), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "38.5-39"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.T7), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Z1), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.T7), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "78.5-79"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.C7), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Z1), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.C), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "40"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.J), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "80"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.C), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.J), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Y3), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "41-42"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A5), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Z2), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Y2), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "42.5-43"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A4), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Z4), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Y4), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "43.5-44"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A6), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Z5), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Y1), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "45"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.F1), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Z6), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Y6), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "46-47"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.F3), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Z7), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Y5), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "47.5-48"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.F2), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Z8), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Y7), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "48.5-49"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.F4), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.H1), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Y), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "50"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.F), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.N), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.D6), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "51-52"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.F7), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.H3), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Y8), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "52.5-53"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.F6), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.H4), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.D5), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "53.5-54"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.F8), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.H5), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.D4), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "55"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.F5), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.H6), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.D8), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "56-57"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.B5), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.H7), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.D7), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "57.5-58"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.B4), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.H8), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.D9), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "58.5-59"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.B6), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.N1), 13, languageShort, gender, personal).Content)
    forecastBlocks = append(forecastBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.D), 13, languageShort, gender, personal))
    forecastBlocks[len(forecastBlocks) - 1].Title = "60"
    forecastBlocks[len(forecastBlocks) - 1].Type = "info"
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.B), 13, languageShort, gender, personal).Content)
    forecastBlocks[len(forecastBlocks) - 1].Content = fmt.Sprintf("%s\n\n%s", forecastBlocks[len(forecastBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.H), 13, languageShort, gender, personal).Content)
    forecast := Prediction{}
    forecast.Title = "Forecast"
    forecast.ImageName = "forecast"
    forecast.Blocks = forecastBlocks
    forecast.BlockType = "forecast"
    // ------ YEAR FORECAST END --------
    predictions := []Prediction{personalfeatures, destiny, money, programs, sexiness, pastlife, parents, kids, relationships, health, lifeguide, forecast}
    for i, prediction := range predictions {
      for _, block := range prediction.Blocks {
        nblock := BlockJSON{ block.Id, block.Content, block.Created, block.Edited, block.PredType, block.Lang, block.Personal, block.Type, block.Title, block.TintColor }
        prediction.BlocksJSON = append(prediction.BlocksJSON, nblock)
      }
      prediction.Blocks = []Block{}
      predictions[i] = prediction
    }
    marshalled, _ := json.Marshal(Response{true, "", predictions})
    return c.Write(marshalled)
    // return c.Write(ResponseNew{true, "", blocks})
  })

  api.Get("/show/types", func(c *routing.Context) error {
    types := []PredictionType{}
    langs := []PredictionType{}
    err = db.Select(&types, "SELECT * FROM predictionType")
    if err != nil {
      marshalled, _ := json.Marshal(Response{false, "Can't parse types", nil})
      return c.Write(marshalled)
    }
    err = db.Select(&langs, "SELECT * FROM predictionLang")
    if err != nil {
      fmt.Println(err)
      marshalled, _ := json.Marshal(Response{false, "Can't parse languages", nil})
      return c.Write(marshalled)
    }
    marshalled, _ := json.Marshal(ResponseType{true, "", [][]PredictionType{types, langs}})
    return c.Write(marshalled)
  })

  api.Get("/show/predictions", func(c *routing.Context) error {
    predictions := []Prediction{}
    err = db.Select(&predictions, "SELECT * FROM prediction")
    if err != nil {
      marshalled, _ := json.Marshal(Response{false, "Can't parse predictions", nil})
      return c.Write(marshalled)
    }
    marshalled, _ := json.Marshal(Response{true, "", predictions})
    return c.Write(marshalled)
  })

  api.Post("/add", func(c *routing.Context) error {
    ptypeid := c.PostForm("ptypeid")
    combo := c.PostForm("combo")
    prediction := c.PostForm("prediction")
    personal := c.PostForm("personal")
    language := c.PostForm("language")
    fmt.Println(ptypeid, combo, prediction, personal, language)
    currPrediction := Block{}
    db.Get(&currPrediction, "SELECT * FROM prediction ORDER BY id DESC LIMIT 1")
    tx := db.MustBegin()
    tx.MustExec(`INSERT INTO prediction(content,type_id,personal,lang_id) VALUES($1,$2,$3,$4)`, prediction, ptypeid, personal, language)
    tx.MustExec(`INSERT INTO predictionRel(prediction_id,combination) VALUES($1,$2)`, currPrediction.Id + 1, combo)
    tx.Commit()
    marshalled, _ := json.Marshal(ResponseAdd{currPrediction.Id + 1})
    return c.Write(marshalled)
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

func setAllCombosNew(icombo []int) Combos {
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
  l4 := checkGreater(e+e)
  l5 := checkGreater(d1+c1)
  l6 := checkGreater(icombo[2]+d)
  d3 := checkGreater(icombo[0]+a2+a1+a3+e+c1+icombo[2])
  c3 := checkGreater(icombo[1]+b2+b1+b3+e+d1+d)
  e3 := checkGreater(l+l1+l2+l3+l4+l5+l6)
  f1 := checkGreater(icombo[0]+f)
  a4 := checkGreater(icombo[0]+f1)
  a5 := checkGreater(icombo[0]+a4)
  a6 := checkGreater(a4+f1)
  f2 := checkGreater(f1+f)
  f3 := checkGreater(f1+f2)
  f4 := checkGreater(f2+f)
  f5 := checkGreater(f+icombo[1])
  f6 := checkGreater(f+f5)
  f7 := checkGreater(f+f6)
  f8 := checkGreater(f6+f5)
  b4 := checkGreater(f5+icombo[1])
  b5 := checkGreater(f5+b4)
  b6 := checkGreater(b4+icombo[1])
  k1 := checkGreater(icombo[1]+g)
  b7 := checkGreater(icombo[1]+k1)
  b8 := checkGreater(icombo[1]+b7)
  k2 := checkGreater(b7+k1)
  k3 := checkGreater(k1+g)
  k4 := checkGreater(k1+k3)
  k5 := checkGreater(k3+g)
  c4 := checkGreater(g+icombo[2])
  k7 := checkGreater(g+c4)
  k8 := checkGreater(g+k7)
  k6 := checkGreater(k7+c4)
  c5 := checkGreater(c4+icombo[2])
  c6 := checkGreater(c4+c5)
  c7 := checkGreater(c5+icombo[2])
  y1 := checkGreater(icombo[2]+y)
  y2 := checkGreater(icombo[2]+y1)
  y3 := checkGreater(icombo[2]+y2)
  y4 := checkGreater(y2+y1)
  y5 := checkGreater(y1+y)
  y6 := checkGreater(y1+y5)
  y7 := checkGreater(y5+y)
  d4 := checkGreater(d+y)
  y8 := checkGreater(d4+y)
  d6 := checkGreater(y8+y)
  d5 := checkGreater(d4+y8)
  d7 := checkGreater(d+d4)
  d8 := checkGreater(d7+d4)
  d9 := checkGreater(d+d7)
  g1 := checkGreater(k+d)
  g2 := checkGreater(d+g1)
  g3 := checkGreater(g2+d)
  g4 := checkGreater(g1+g2)
  g5 := checkGreater(k+g1)
  g6 := checkGreater(g5+k)
  g7 := checkGreater(g5+g)
  t1 := checkGreater(icombo[0]+k)
  t2 := checkGreater(t1+k)
  t3 := checkGreater(t2+k)
  t4 := checkGreater(t1+t2)
  t5 := checkGreater(icombo[0]+t1)
  t6 := checkGreater(t1+t5)
  t7 := checkGreater(icombo[0]+t5)
  o1 := checkGreater(b8+g3)
  o2 := checkGreater(b7+g2)
  o3 := checkGreater(k2+g4)
  x3 := checkGreater(k1+g1)
  x4 := checkGreater(k4+g6)
  x5 := checkGreater(k3+g5)
  x6 := checkGreater(k5+g7)
  x8 := checkGreater(k8+t3)
  e4 := checkGreater(k7+t2)
  e5 := checkGreater(k6+t4)
  e6 := checkGreater(c4+t1)
  e7 := checkGreater(c6+t5)
  e8 := checkGreater(c5+t5)
  z1 := checkGreater(c7+t7)
  z2 := checkGreater(y3+a5)
  z4 := checkGreater(y2+a4)
  z5 := checkGreater(y4+a6)
  z6 := checkGreater(y1+f1)
  z7 := checkGreater(y6+f3)
  z8 := checkGreater(y5+f2)
  h1 := checkGreater(y7+f4)
  h3 := checkGreater(d6+f7)
  h4 := checkGreater(y8+f6)
  h5 := checkGreater(d5+f8)
  h6 := checkGreater(d4+f5)
  h7 := checkGreater(d8+b5)
  h8 := checkGreater(d7+b4)
  n1 := checkGreater(d9+b6)
  return Combos{icombo[0],icombo[1],icombo[2],d,e,a1,a2,b1,b2,c1,c2,d1,d2,x,x1,x2,f,g,y,k,e1,e2,o,u,h,j,m,n,t,z,s,b3,a3,l,l1,l2,l3,l4,l5,l6,d3,c3,e3,f1,a4,a5,a6,f2,f3,f4,f5,f6,f7,f8,b4,b5,b6,k1,b7,b8,k2,k3,k4,k5,c4,k7,k8,k6,c5,c6,c7,y1,y2,y3,y4,y5,y6,y7,d4,y8,d6,d5,d7,d8,d9,g1,g2,g3,g4,g5,g6,g7,t1,t2,t3,t4,t5,t6,t7,o1,o2,o3,x3,x4,x5,x6,x8,e4,e5,e6,e7,e8,z1,z2,z4,z5,z6,z7,z8,h1,h3,h4,h5,h6,h7,h8,n1}
}

func testEq(a, b []int) bool {

    if (a == nil) != (b == nil) {
        return false;
    }

    if len(a) != len(b) {
        return false
    }

    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }

    return true
}

func checkAnswers(input [][]int, answer []int, random_opt ...bool) bool {
  random := false
  if len(random_opt) > 0 {
    random = random_opt[0]
  }
  for _, item := range input {
    if random {
      sort.Ints(item)
      sort.Ints(answer)
    }
    if testEq(item, answer) {
      return true
    }
  }
  return false
}

func getAnswerFromTable(db *sqlx.DB, id string, tableNumber int, lang string, sex string, personal_short string, short_opt ...bool) Block {
  if (id != "") {
    personal_bool := false
    // personal := "common"
    if personal_short == "p" {
      // personal = "personal"
      personal_bool = true
    }
    language := Language{}
    err := db.Get(&language, "SELECT * FROM predictionlang WHERE short = $1", lang)
    if err != nil {
      language.Id = 1
    }
    // field := fmt.Sprintf("%s_%s", personal, lang)
    // table := fmt.Sprintf("answers%d", tableNumber)
    _, err = strconv.Atoi(id)
    if err == nil {
      id = fmt.Sprintf("'%s'", id)
    }
    query := fmt.Sprintf("SELECT * FROM prediction WHERE type_id=%d AND personal=%v AND lang_id=%d AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=%s)", tableNumber, personal_bool, language.Id, id)
    // query := fmt.Sprintf("SELECT %s FROM %s WHERE id=%s", table, field, id)
    fmt.Println(query)
    block := Block{}
    err = db.Get(&block, query)
    if err != nil {
      fmt.Println(err)
      return Block{}
    }
    contentbygender := ContentByGender{}
    err = json.Unmarshal([]byte(block.Content), &contentbygender)
    fmt.Println(err)
    if err == nil {
      if sex == "m" {
        block.Content = contentbygender.Male
      } else {
        block.Content = contentbygender.Female
      }
    }
    return block
  } else {
    return Block{}
  }
}

func contains(arr []int, item int) bool {
   for _, a := range arr {
      if a == item {
         return true
      }
   }
   return false
}

func prepareArray(inputArray []int) []int {
  out := []int{}
  for _, item := range inputArray {
    if !contains(out, item) {
      out = append(out, item)
    }
  }
  return out;
}
