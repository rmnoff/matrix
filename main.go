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

  // api.Get("/check/<input>", func(c *routing.Context) error {
  //   input := c.Param("input")
  //   timestamp := input[:len(input) - 4]
  //   gender := string(input[len(input) - 4])
  //   personal := string(input[len(input) - 3])
  //   languageShort := input[len(input) - 2:]
  //   fmt.Println(timestamp, gender, personal, languageShort)
  //   if timestamp == "" {
  //     marshalled, _ := json.Marshal(Response{false, "No timestamp provided", nil})
  //     return c.Write(marshalled)
  //   }
  //   combo := countBD(timestamp)
  //   if combo == nil {
  //     marshalled, _ := json.Marshal(Response{false, "Timestamp corrupted", nil})
  //     return c.Write(marshalled)
  //   }
  //   finalCombos := setAllCombos(combo)
  //   // prog1 := fmt.Sprintf("prog1: [%v %v %v]", combo[0], finalCombos[3], finalCombos[2])
  //   // prog2 := fmt.Sprintf("prog2: [%v %v %v]", combo[1], finalCombos[5], finalCombos[4])
  //   // prog3 := fmt.Sprintf("prog3: [%v %v %v]", combo[2], finalCombos[7], finalCombos[6])
  //   // prog4 := fmt.Sprintf("prog4: [%v %v %v]", finalCombos[13], finalCombos[15], finalCombos[19])
  //   // prog5 := fmt.Sprintf("prog5: [%v %v %v]", finalCombos[16], finalCombos[14], finalCombos[20])
  //
  //   isPersonal := false
  //   if personal == "p" {
  //     isPersonal = true
  //   }
  //
  //   language := Language{}
  //   err = db.Get(&language, "SELECT * FROM predictionlang WHERE short = $1", languageShort)
  //   if err != nil {
  //     language.Id = 1
  //   }
  //
  //   pastLife := Prediction{}
  //   pastLifeBlock := Block{}
  //   pastLifeBlock.Type = "info"
  //   pastLifeBlock.Title = "Previous Life Common"
  //   pastLifePredictionCombo := fmt.Sprintf("%d-%d-%d", finalCombos[8], finalCombos[9], finalCombos[0])
  //   fmt.Println(pastLifePredictionCombo)
  //   err = db.Get(&pastLifeBlock, "SELECT * FROM prediction WHERE type_id=9 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", pastLifePredictionCombo, isPersonal, language.Id)
  //   if err != nil {
  //     log.Println(err)
  //     marshalled, _ := json.Marshal(Response{false, "Can't parse past life prediction", nil})
  //     return c.Write(marshalled)
  //   }
  //   pastLife.Title = "Previous Life"
  //   pastLife.Blocks = []Block{pastLifeBlock}
  //   pastLife.BlockType = "default"
  //
  //   personalFeatures := Prediction{}
  //   personalFeaturesPos := Block{}
  //   personalFeaturesPosSecond := Block{}
  //   personalFeaturesPos.Type = "info"
  //   personalFeaturesPos.Title = "Personal Features Positive"
  //   personalFeaturesPosCombo := fmt.Sprintf("%d", combo[0])
  //   personalFeaturesPosSecondCombo := fmt.Sprintf("%d", combo[1])
  //   personalFeaturesNeg := Block{}
  //   personalFeaturesNegSecond := Block{}
  //   personalFeaturesNeg.Type = "info"
  //   personalFeaturesNeg.Title = "Personal Features Negative"
  //   personalFeaturesNegCombo := fmt.Sprintf("%d", combo[0])
  //   personalFeaturesNegSecondCombo := fmt.Sprintf("%d", combo[1])
  //   personalFeaturesSoc := Block{}
  //   personalFeaturesSoc.Type = "info"
  //   personalFeaturesSoc.Title = "Personal Features Social"
  //   personalFeaturesSocCombo := fmt.Sprintf("%d", finalCombos[1])
  //   err = db.Get(&personalFeaturesPos, "SELECT * FROM prediction WHERE type_id=1 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", personalFeaturesPosCombo, isPersonal, language.Id)
  //   if err != nil {
  //     log.Println(err)
  //     marshalled, _ := json.Marshal(Response{false, "Can't parse positive personal features 1", nil})
  //     return c.Write(marshalled)
  //   }
  //   err = db.Get(&personalFeaturesPosSecond, "SELECT * FROM prediction WHERE type_id=1 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", personalFeaturesPosSecondCombo, isPersonal, language.Id)
  //   if err != nil {
  //     log.Println(err)
  //     marshalled, _ := json.Marshal(Response{false, "Can't parse positive personal features 2", nil})
  //     return c.Write(marshalled)
  //   }
  //   personalFeaturesPos.Content = fmt.Sprintf("%s %s", personalFeaturesPos.Content, personalFeaturesPosSecond.Content)
  //   err = db.Get(&personalFeaturesNeg, "SELECT * FROM prediction WHERE type_id=2 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", personalFeaturesNegCombo, isPersonal, language.Id)
  //   if err != nil {
  //     log.Println(err)
  //     marshalled, _ := json.Marshal(Response{false, "Can't parse negative personal features 1", nil})
  //     return c.Write(marshalled)
  //   }
  //   err = db.Get(&personalFeaturesNegSecond, "SELECT * FROM prediction WHERE type_id=2 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", personalFeaturesNegSecondCombo, isPersonal, language.Id)
  //   if err != nil {
  //     log.Println(err)
  //     marshalled, _ := json.Marshal(Response{false, "Can't parse negative personal features 2", nil})
  //     return c.Write(marshalled)
  //   }
  //   personalFeaturesNeg.Content = fmt.Sprintf("%s %s", personalFeaturesNeg.Content, personalFeaturesNegSecond.Content)
  //   err = db.Get(&personalFeaturesSoc, "SELECT * FROM prediction WHERE type_id=3 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", personalFeaturesSocCombo, isPersonal, language.Id)
  //   if err != nil {
  //     log.Println(err)
  //     marshalled, _ := json.Marshal(Response{false, "Can't parse social personal features", nil})
  //     return c.Write(marshalled)
  //   }
  //   personalFeatures.Title = "Personal Features"
  //   personalFeatures.Blocks = []Block{personalFeaturesPos,personalFeaturesNeg,personalFeaturesSoc}
  //   personalFeatures.BlockType = "default"
  //
  //   relationship := Prediction{}
  //   relationshipBlock := Block{}
  //   relationshipBlockSecond := Block{}
  //   relationshipBlockThird := Block{}
  //   relationshipBlock.Type = "info"
  //   relationshipBlock.Title = "Relationship Common"
  //   relationshipCombo := fmt.Sprintf("%d", finalCombos[11])
  //   relationshipComboSecond := fmt.Sprintf("%d", finalCombos[8])
  //   relationshipComboThird := fmt.Sprintf("%d", finalCombos[10])
  //   fmt.Println(relationshipCombo)
  //   err = db.Get(&relationshipBlock, "SELECT * FROM prediction WHERE type_id=5 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", relationshipCombo, isPersonal, language.Id)
  //   if err != nil {
  //     log.Println(err)
  //     marshalled, _ := json.Marshal(Response{false, "Can't parse relationship prediction", nil})
  //     return c.Write(marshalled)
  //   }
  //   err = db.Get(&relationshipBlockSecond, "SELECT * FROM prediction WHERE type_id=5 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", relationshipComboSecond, isPersonal, language.Id)
  //   if err != nil {
  //     log.Println(err)
  //     marshalled, _ := json.Marshal(Response{false, "Can't parse relationship prediction", nil})
  //     return c.Write(marshalled)
  //   }
  //   err = db.Get(&relationshipBlockThird, "SELECT * FROM prediction WHERE type_id=5 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", relationshipComboThird, isPersonal, language.Id)
  //   if err != nil {
  //     log.Println(err)
  //     marshalled, _ := json.Marshal(Response{false, "Can't parse relationship prediction", nil})
  //     return c.Write(marshalled)
  //   }
  //   relationshipBlock.Content = fmt.Sprintf("%s %s %s", relationshipBlock.Content, relationshipBlockSecond.Content, relationshipBlockThird.Content)
  //   relationship.Title = "Relationship"
  //   relationship.Blocks = []Block{relationshipBlock}
  //   relationship.BlockType = "default"
  //
  //   lifeGuide := Prediction{}
  //   lifeGuideBlock := Block{}
  //   lifeGuideBlockSecond := Block{}
  //   lifeGuideBlockThird := Block{}
  //   lifeGuideBlock.Type = "info"
  //   lifeGuideBlock.Title = "Life Guide Common"
  //   lifeGuideCombo := fmt.Sprintf("%d", combo[0])
  //   lifeGuideSecondCombo := fmt.Sprintf("%d", combo[1])
  //   lifeGuideThirdCombo := fmt.Sprintf("%d", finalCombos[1])
  //   err = db.Get(&lifeGuideBlock, "SELECT * FROM prediction WHERE type_id=11 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", lifeGuideCombo, isPersonal, language.Id)
  //   if err != nil {
  //     log.Println(err)
  //     marshalled, _ := json.Marshal(Response{false, "Can't parse life guide prediction 1", nil})
  //     return c.Write(marshalled)
  //   }
  //   err = db.Get(&lifeGuideBlockSecond, "SELECT * FROM prediction WHERE type_id=11 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", lifeGuideSecondCombo, isPersonal, language.Id)
  //   if err != nil {
  //     log.Println(err)
  //     marshalled, _ := json.Marshal(Response{false, "Can't parse life guide prediction 2", nil})
  //     return c.Write(marshalled)
  //   }
  //   err = db.Get(&lifeGuideBlockThird, "SELECT * FROM prediction WHERE type_id=11 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", lifeGuideThirdCombo, isPersonal, language.Id)
  //   if err != nil {
  //     log.Println(err)
  //     marshalled, _ := json.Marshal(Response{false, "Can't parse life guide prediction 3", nil})
  //     return c.Write(marshalled)
  //   }
  //   lifeGuideBlock.Content = fmt.Sprintf("%s %s %s", lifeGuideBlock.Content, lifeGuideBlockSecond.Content, lifeGuideBlockThird.Content)
  //   lifeGuide.Title = "Life Guide"
  //   lifeGuide.Blocks = []Block{lifeGuideBlock}
  //   lifeGuide.BlockType = "default"
  //
  //   sex := Prediction{}
  //   sexBlock := Block{}
  //   sexBlock.Type = "info"
  //   sexBlock.Title = "Sexiness Common"
  //   sexCombo := fmt.Sprintf("%d-%d-%d", finalCombos[1], finalCombos[17], finalCombos[18])
  //   fmt.Println(sexCombo)
  //   err = db.Get(&sexBlock, "SELECT * FROM prediction WHERE type_id=228 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", sexCombo, isPersonal, language.Id)
  //   if err != nil {
  //     log.Println(err)
  //     marshalled, _ := json.Marshal(Response{false, "Can't parse sexiness prediction", nil})
  //     return c.Write(marshalled)
  //   }
  //   sex.Title = "Sexiness"
  //   sex.Blocks = []Block{sexBlock}
  //   sex.BlockType = "default"
  //
  //   destiny := Prediction{}
  //   destinyBlock := Block{}
  //   destinyBlockSecond := Block{}
  //   destinyBlockThird := Block{}
  //   destinyBlock.Type = "info"
  //   destinyBlock.Title = "Destiny Common"
  //   destinyCombo := fmt.Sprintf("%d", finalCombos[21])
  //   destinySecondCombo := fmt.Sprintf("%d", finalCombos[22])
  //   destinyThirdCombo := fmt.Sprintf("%d", finalCombos[23])
  //   err = db.Get(&destinyBlock, "SELECT * FROM prediction WHERE type_id=8 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", destinyCombo, isPersonal, language.Id)
  //   if err != nil {
  //     log.Println(err)
  //     marshalled, _ := json.Marshal(Response{false, "Can't parse destiny prediction", nil})
  //     return c.Write(marshalled)
  //   }
  //   err = db.Get(&destinyBlockSecond, "SELECT * FROM prediction WHERE type_id=8 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", destinySecondCombo, isPersonal, language.Id)
  //   if err != nil {
  //     log.Println(err)
  //     marshalled, _ := json.Marshal(Response{false, "Can't parse destiny prediction", nil})
  //     return c.Write(marshalled)
  //   }
  //   err = db.Get(&destinyBlockThird, "SELECT * FROM prediction WHERE type_id=8 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", destinyThirdCombo, isPersonal, language.Id)
  //   if err != nil {
  //     log.Println(err)
  //     marshalled, _ := json.Marshal(Response{false, "Can't parse destiny prediction", nil})
  //     return c.Write(marshalled)
  //   }
  //   destinyBlock.Content = fmt.Sprintf("%s %s %s", destinyBlock.Content, destinyBlockSecond.Content, destinyBlockThird.Content)
  //   // destinyCommonBlock := Block{}
  //   // destinyCommonBlock.Type = "info"
  //   // destinyCommonBlock.Title = "Destiny Common"
  //   // destinyCommonCombo := fmt.Sprintf("%d-%d-%d", finalCombos[21], finalCombos[22], finalCombos[23])
  //   // fmt.Println(destinyCommonCombo)
  //   // err = db.Get(&destinyCommonBlock, "SELECT * FROM prediction WHERE type_id=220 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", destinyCommonCombo, isPersonal, language.Id)
  //   // if err != nil {
  //   //   log.Println(err)
  //   //   marshalled, _ := json.Marshal(Response{false, "Can't parse destiny common prediction", nil})
  //   //   return c.Write(marshalled)
  //   // }
  //   destiny.Title = "Destiny"
  //   // destiny.Blocks = []Block{destinyBlock, destinyCommonBlock}
  //   destiny.Blocks = []Block{destinyBlock}
  //   destiny.BlockType = "default"
  //
  //   money := Prediction{}
  //   moneyBlock := Block{}
  //   moneyBlock.Type = "info"
  //   moneyBlock.Title = "Money"
  //   moneySuccessBlock := Block{}
  //   moneySuccessBlock.Type = "info"
  //   moneySuccessBlock.Title = "To become successful"
  //   moneySuccessSecondBlock := Block{}
  //   moneySuccessThirdBlock := Block{}
  //   moneySuccessFourthBlock := Block{}
  //   moneyImportantBlock := Block{}
  //   moneyImportantBlock.Type = "expandable"
  //   moneyImportantBlock.Title = "Important!"
  //   moneyBlockCombo := fmt.Sprintf("%d", finalCombos[12])
  //   moneySuccessBlockCombo := fmt.Sprintf("%d", finalCombos[10])
  //   moneySuccessSecondBlockCombo := fmt.Sprintf("%d", combo[2])
  //   moneySuccessThirdBlockCombo := fmt.Sprintf("%d", finalCombos[6])
  //   moneySuccessFourthBlockCombo := fmt.Sprintf("%d", finalCombos[7])
  //   moneyImportantCombo := fmt.Sprintf("", )
  //   err = db.Get(&moneyBlock, "SELECT * FROM prediction WHERE type_id=4 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", moneyBlockCombo, isPersonal, language.Id)
  //   if err != nil {
  //     log.Println(err)
  //     marshalled, _ := json.Marshal(Response{false, "Can't parse money prediction", nil})
  //     return c.Write(marshalled)
  //   }
  //   err = db.Get(&moneySuccessBlock, "SELECT * FROM prediction WHERE type_id=214 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", moneySuccessBlockCombo, isPersonal, language.Id)
  //   if err != nil {
  //     log.Println(err)
  //     marshalled, _ := json.Marshal(Response{false, "Can't parse success prediction 1", nil})
  //     return c.Write(marshalled)
  //   }
  //   err = db.Get(&moneySuccessSecondBlock, "SELECT * FROM prediction WHERE type_id=214 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", moneySuccessSecondBlockCombo, isPersonal, language.Id)
  //   if err != nil {
  //     log.Println(err)
  //     marshalled, _ := json.Marshal(Response{false, "Can't parse success prediction 2", nil})
  //     return c.Write(marshalled)
  //   }
  //   err = db.Get(&moneySuccessThirdBlock, "SELECT * FROM prediction WHERE type_id=214 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", moneySuccessThirdBlockCombo, isPersonal, language.Id)
  //   if err != nil {
  //     log.Println(err)
  //     marshalled, _ := json.Marshal(Response{false, "Can't parse success prediction 3", nil})
  //     return c.Write(marshalled)
  //   }
  //   err = db.Get(&moneySuccessFourthBlock, "SELECT * FROM prediction WHERE type_id=214 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", moneySuccessFourthBlockCombo, isPersonal, language.Id)
  //   if err != nil {
  //     log.Println(err)
  //     marshalled, _ := json.Marshal(Response{false, "Can't parse success prediction 4", nil})
  //     return c.Write(marshalled)
  //   }
  //   moneySuccessBlock.Content = fmt.Sprintf("%s %s %s %s", moneySuccessBlock.Content, moneySuccessSecondBlock.Content, moneySuccessThirdBlock.Content, moneySuccessFourthBlock.Content)
  //   money.Title = "Money"
  //   money.Blocks = []Block{moneyBlock, moneySuccessBlock, moneySuccessSecondBlock, moneySuccessThirdBlock, moneySuccessFourthBlock}
  //   money.BlockType = "default"
  //
  //   kids := Prediction{}
  //   kidsBlock := Block{}
  //   kidsBlock.Type = "info"
  //   kidsBlock.Title = "Lessons from kids"
  //   kidsImportantBlock := Block{}
  //   kidsImportantBlock.Type = "expandable"
  //   kidsImportantBlock.Title = "Important!"
  //   kidsCombo := fmt.Sprintf("%d", combo[0])
  //   kidsComboSecond := fmt.Sprintf("%d", finalCombos[2])
  //   kidsComboThird := fmt.Sprintf("%d", finalCombos[3])
  //   // toopa poebat'
  //   if combo[0] == 15 || combo[0] == 8 || combo[0] == 18 {
  //     if finalCombos[2] == 7 {
  //       if finalCombos[3] == 22 {
  //
  //       }
  //     }
  //     if finalCombos[2] == 17 {
  //       if finalCombos[3] == 9 {
  //
  //       }
  //     }
  //     if finalCombos[2] == 12 {
  //       if finalCombos[3] == 6 {
  //
  //       }
  //     }
  //   }
  //
  //
  //   predictions := []Prediction{pastLife, personalFeatures, relationship, lifeGuide, sex, destiny, money}
  //
  //   for _, pred := range predictions {
  //     for _, content := range pred.Blocks {
  //       contentbygender := ContentByGender{}
  //       err := json.Unmarshal([]byte(content.Content), &contentbygender)
  //       if err == nil {
  //         if gender == "m" {
  //           content.Content = contentbygender.Male
  //         } else {
  //           content.Content = contentbygender.Female
  //         }
  //       }
  //     }
  //   }
  //   marshalled, _ := json.Marshal(Response{true, "", predictions})
  //   return c.Write(marshalled)
  // })

  api.Get("/check/new/<input>", func(c *routing.Context) error {
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
    // toCheck := [][]int{
    //   {fc.A, fc.A2, fc.A1},
    //   {fc.B, fc.B2, fc.B1},
    //   {fc.F, fc.Y, fc.O},
    //   {fc.K, fc.G, fc.U},
    //   {fc.E, fc.E1, fc.E2},
    //   {fc.D1, fc.X1, fc.X},
    //   {fc.X, fc.X2, fc.C1},
    //   {fc.H, fc.J, fc.M},
    //   {fc.N, fc.T, fc.Z},
    //   {fc.A, fc.B, fc.L},
    //   {fc.A2, fc.B2, fc.L1},
    //   {fc.A1, fc.B1, fc.L2},
    //   {fc.A3, fc.B3, fc.L3},
    //   {fc.E, fc.E, fc.L4},
    //   {fc.D1, fc.C1, fc.L5},
    //   {fc.D, fc.C, fc.L6},
    //   {fc.D3, fc.C3, fc.E3},
    // }
    // if(checkAnswers(toCheck, []int{17,5,6}, true)) {
    //   blocks = append(blocks, getAnswerFromTable(db, "'17-5-6'", 10, languageShort, gender, personal))
    // }

    // ------ PERSONAL FEATURES BEGIN --------
    personalFeaturesBlocks := []Block{}
    personalFeaturesBlocks = append(personalFeaturesBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A), 1, languageShort, gender, personal))
    personalFeaturesBlocks[len(personalFeaturesBlocks) - 1].Title = "personal features positive"
    personalFeaturesBlocks[len(personalFeaturesBlocks) - 1].Type = "info"
    if fc.A != fc.B {
      personalFeaturesBlocks[len(personalFeaturesBlocks) - 1].Content = fmt.Sprintf("%s %s", personalFeaturesBlocks[len(personalFeaturesBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.B), 1, languageShort, gender, personal).Content)
    }
    personalFeaturesBlocks = append(personalFeaturesBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A), 2, languageShort, gender, personal))
    personalFeaturesBlocks[len(personalFeaturesBlocks) - 1].Title = "personal features negative"
    personalFeaturesBlocks[len(personalFeaturesBlocks) - 1].Type = "info"
    if fc.A != fc.B {
      personalFeaturesBlocks[len(personalFeaturesBlocks) - 1].Content = fmt.Sprintf("%s %s", personalFeaturesBlocks[len(personalFeaturesBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.B), 2, languageShort, gender, personal).Content)
    }
    personalFeaturesBlocks = append(personalFeaturesBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.E), 3, languageShort, gender, personal))
    personalFeaturesBlocks[len(personalFeaturesBlocks) - 1].Title = "personal features social"
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
    destinyBlocks[len(destinyBlocks) - 1].Title = "destiny"
    destinyBlocks[len(destinyBlocks) - 1].Type = "info"
    if fc.H != fc.J {
      // destinyBlocks = append(destinyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.J), 8, languageShort, gender, personal))
      destinyBlocks[len(destinyBlocks) - 1].Content = fmt.Sprintf("%s %s", destinyBlocks[len(destinyBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.J), 8, languageShort, gender, personal).Content)
    }
    if fc.H != fc.M && fc.J != fc.M {
      // destinyBlocks = append(destinyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.M), 8, languageShort, gender, personal))
      destinyBlocks[len(destinyBlocks) - 1].Content = fmt.Sprintf("%s %s", destinyBlocks[len(destinyBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.M), 8, languageShort, gender, personal).Content)
    }
    // destinyBlocks = append(destinyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.N), 8, languageShort, gender, personal))
    destinyBlocks[len(destinyBlocks) - 1].Content = fmt.Sprintf("%s %s", destinyBlocks[len(destinyBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.N), 8, languageShort, gender, personal).Content)
    if fc.T != fc.N {
      // destinyBlocks = append(destinyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.T), 8, languageShort, gender, personal))
      destinyBlocks[len(destinyBlocks) - 1].Content = fmt.Sprintf("%s %s", destinyBlocks[len(destinyBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.T), 8, languageShort, gender, personal).Content)
    }
    if fc.Z != fc.N && fc.Z != fc.T {
      // destinyBlocks = append(destinyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Z), 8, languageShort, gender, personal))
      destinyBlocks[len(destinyBlocks) - 1].Content = fmt.Sprintf("%s %s", destinyBlocks[len(destinyBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Z), 8, languageShort, gender, personal).Content)
    }
    destinyBlocks = append(destinyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.S), 220, languageShort, gender, personal))
    destinyBlocks[len(destinyBlocks) - 1].Title = "destiny common"
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
    moneyBlocks[len(moneyBlocks) - 1].Title = "money"
    moneyBlocks[len(moneyBlocks) - 1].Type = "info"
    moneyBlocks = append(moneyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.X), 214, languageShort, gender, personal))
    moneyBlocks[len(moneyBlocks) - 1].Title = "to become successful"
    moneyBlocks[len(moneyBlocks) - 1].Type = "info"
    if fc.X != fc.C {
      // moneyBlocks = append(moneyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.C), 214, languageShort, gender, personal))
      moneyBlocks[len(moneyBlocks) - 1].Content = fmt.Sprintf("%s %s", moneyBlocks[len(moneyBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.C), 214, languageShort, gender, personal).Content)
    }
    if fc.X != fc.C1 && fc.C != fc.C1 {
      // moneyBlocks = append(moneyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.C1), 214, languageShort, gender, personal))
      moneyBlocks[len(moneyBlocks) - 1].Content = fmt.Sprintf("%s %s", moneyBlocks[len(moneyBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.C1), 214, languageShort, gender, personal).Content)
    }
    if fc.X != fc.C2 && fc.C != fc.C2 && fc.C1 != fc.C2 {
      // moneyBlocks = append(moneyBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.C2), 214, languageShort, gender, personal))
      moneyBlocks[len(moneyBlocks) - 1].Content = fmt.Sprintf("%s %s", moneyBlocks[len(moneyBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.C2), 214, languageShort, gender, personal).Content)
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
      moneyBlocks[len(moneyBlocks) - 1].Title = "important"
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
      moneyBlocks[len(moneyBlocks) - 1].Title = "important"
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
      moneyBlocks[len(moneyBlocks) - 1].Title = "important"
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
      moneyBlocks[len(moneyBlocks) - 1].Title = "important"
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
      moneyBlocks[len(moneyBlocks) - 1].Title = "important"
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
      lessonsBlocks[len(lessonsBlocks) - 1].Title = fmt.Sprintf("program %d", lessonsCount)
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
        lessonsBlocks[len(lessonsBlocks) - 1].Title = fmt.Sprintf("program %d", lessonsCount)
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
    sexinessBlocks[len(sexinessBlocks) - 1].Title = "sexiness"
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
    pastLifeBlocks[len(pastLifeBlocks) - 1].Title = "previous life"
    pastLifeBlocks[len(pastLifeBlocks) - 1].Type = "info"
    pastlife := Prediction{}
    pastlife.Title = "Past Life"
    pastlife.ImageName = "previous_life"
    pastlife.Blocks = pastLifeBlocks
    pastlife.BlockType = "default"
    // ------ PAST LIFE END --------
    // ------ PARENTS  BEGIN --------
    parentsBlocks := []Block{}
    parentsBlocks = append(parentsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.F), 234, languageShort, gender, personal))
    parentsBlocks[len(parentsBlocks) - 1].Title = "Parents"
    parentsBlocks[len(parentsBlocks) - 1].Type = "info"
    if fc.D1 == 21 {
      parentsBlocks = append(parentsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", 1), 217, languageShort, gender, personal))
      parentsBlocks[len(parentsBlocks) - 1].Title = "Important"
      parentsBlocks[len(parentsBlocks) - 1].Type = "expandable"
    }
    parentsBlocks = append(parentsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.G), 234, languageShort, gender, personal))
    parentsBlocks[len(parentsBlocks) - 1].Title = "Possible insult on parents"
    parentsBlocks[len(parentsBlocks) - 1].Type = "info"
    // if fc.F != fc.Y {
      // parentsBlocks = append(parentsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Y), 234, languageShort, gender, personal))
      parentsBlocks[len(parentsBlocks) - 1].Content = fmt.Sprintf("%s %s", parentsBlocks[len(parentsBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.Y), 234, languageShort, gender, personal).Content)
    // }
    // if fc.F != fc.O && fc.Y != fc.O {
      // parentsBlocks = append(parentsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.O), 234, languageShort, gender, personal))
      parentsBlocks[len(parentsBlocks) - 1].Content = fmt.Sprintf("%s %s", parentsBlocks[len(parentsBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.O), 234, languageShort, gender, personal).Content)
    // }
    // if fc.G != fc.K {
      // parentsBlocks = append(parentsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K), 234, languageShort, gender, personal))
      parentsBlocks[len(parentsBlocks) - 1].Content = fmt.Sprintf("%s %s", parentsBlocks[len(parentsBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.K), 234, languageShort, gender, personal).Content)
    // }
    // if fc.G != fc.U && fc.K != fc.U {
      // parentsBlocks = append(parentsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.U), 234, languageShort, gender, personal))
      parentsBlocks[len(parentsBlocks) - 1].Content = fmt.Sprintf("%s %s", parentsBlocks[len(parentsBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.U), 234, languageShort, gender, personal).Content)
    // }
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
    kidsBlocks[len(kidsBlocks) - 1].Title = "Children"
    kidsBlocks[len(kidsBlocks) - 1].Type = "info"
    if fc.A != fc.A2 {
      // kidsBlocks = append(kidsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A2), 66, languageShort, gender, personal))
      kidsBlocks[len(kidsBlocks) - 1].Content = fmt.Sprintf("%s %s", kidsBlocks[len(kidsBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A2), 66, languageShort, gender, personal).Content)
    }
    if fc.A != fc.A1 && fc.A2 != fc.A1 {
      // kidsBlocks = append(kidsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A1), 66, languageShort, gender, personal))
      kidsBlocks[len(kidsBlocks) - 1].Content = fmt.Sprintf("%s %s", kidsBlocks[len(kidsBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.A1), 66, languageShort, gender, personal).Content)
    }
    toCheck = [][]int{
      {fc.A, fc.A1, fc.A2},
      {fc.A1, fc.A, fc.A2},
    }
    if checkAnswers(toCheck, []int{6,17,5}) {
      // kidsBlocks = append(kidsBlocks, getAnswerFromTable(db, "'6-17-5'", 232, languageShort, gender, personal))
      kidsBlocks[len(kidsBlocks) - 1].Content = fmt.Sprintf("%s %s", kidsBlocks[len(kidsBlocks) - 1].Content, getAnswerFromTable(db, "'6-17-5'", 232, languageShort, gender, personal).Content)
    }
    if checkAnswers(toCheck, []int{7,15,22}) {
      // kidsBlocks = append(kidsBlocks, getAnswerFromTable(db, "'7-15-22'", 232, languageShort, gender, personal))
      kidsBlocks[len(kidsBlocks) - 1].Content = fmt.Sprintf("%s %s", kidsBlocks[len(kidsBlocks) - 1].Content, getAnswerFromTable(db, "'7-15-22'", 232, languageShort, gender, personal).Content)
    }
    if checkAnswers(toCheck, []int{8,9,17}) {
      // kidsBlocks = append(kidsBlocks, getAnswerFromTable(db, "'8-9-17'", 232, languageShort, gender, personal))
      kidsBlocks[len(kidsBlocks) - 1].Content = fmt.Sprintf("%s %s", kidsBlocks[len(kidsBlocks) - 1].Content, getAnswerFromTable(db, "'8-9-17'", 232, languageShort, gender, personal).Content)
    }
    if checkAnswers(toCheck, []int{8,13,21}) {
      // kidsBlocks = append(kidsBlocks, getAnswerFromTable(db, "'8-13-21'", 232, languageShort, gender, personal))
      kidsBlocks[len(kidsBlocks) - 1].Content = fmt.Sprintf("%s %s", kidsBlocks[len(kidsBlocks) - 1].Content, getAnswerFromTable(db, "'8-13-21'", 232, languageShort, gender, personal).Content)
    }
    if checkAnswers(toCheck, []int{6,12,18}) {
      // kidsBlocks = append(kidsBlocks, getAnswerFromTable(db, "'6-12-18'", 232, languageShort, gender, personal))
      kidsBlocks[len(kidsBlocks) - 1].Content = fmt.Sprintf("%s %s", kidsBlocks[len(kidsBlocks) - 1].Content, getAnswerFromTable(db, "'6-12-18'", 232, languageShort, gender, personal).Content)
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
    relationshipsBlocks[len(relationshipsBlocks) - 1].Title = "Relationships"
    relationshipsBlocks[len(relationshipsBlocks) - 1].Type = "info"
    if fc.X1 != fc.D1 {
      // relationshipsBlocks = append(relationshipsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.D1), 5, languageShort, gender, personal))
      relationshipsBlocks[len(relationshipsBlocks) - 1].Content = fmt.Sprintf("%s %s", relationshipsBlocks[len(relationshipsBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.D1), 5, languageShort, gender, personal).Content)
    }
    if fc.X1 != fc.X && fc.D1 != fc.X {
      // relationshipsBlocks = append(relationshipsBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.X), 5, languageShort, gender, personal))
      relationshipsBlocks[len(relationshipsBlocks) - 1].Content = fmt.Sprintf("%s %s", relationshipsBlocks[len(relationshipsBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.X), 5, languageShort, gender, personal).Content)
    }
    toCheck = [][]int{
      {fc.X, fc.X1},
      {fc.X1, fc.X},
    }
    if checkAnswers(toCheck, []int{22,7}) {
      relationshipsBlocks = append(relationshipsBlocks, getAnswerFromTable(db, "'22-7'", 231, languageShort, gender, personal))
      relationshipsBlocks[len(relationshipsBlocks) - 1].Title = "Important"
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
      healthBlocks[len(healthBlocks) - 1].Content = fmt.Sprintf("%s %s", healthBlocks[len(healthBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal).Content)
    }
    healthBlocks = append(healthBlocks, getAnswerFromTable(db, "2", 12, languageShort, gender, personal))
    healthBlocks[len(healthBlocks) - 1].Title = "Occipital and temporal lobes of the brain, eyes, ears, nose, face, upper jaw, upper jaw teeth, optic nerve, cerebral cortex."
    healthBlocks[len(healthBlocks) - 1].Type = "info"
    health2 := prepareArray([]int{fc.A2, fc.B2, fc.L1})
    for _, item := range health2 {
      // healthBlocks = append(healthBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal))
      healthBlocks[len(healthBlocks) - 1].Content = fmt.Sprintf("%s %s", healthBlocks[len(healthBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal).Content)
    }
    healthBlocks = append(healthBlocks, getAnswerFromTable(db, "3", 12, languageShort, gender, personal))
    healthBlocks[len(healthBlocks) - 1].Title = "Thyroid gland, trachea, bronchi, throat, vocal cords, shoulders, arms, seventh cervical vertebra, all cervical vertebrae, lower jaw, lower jaw teeth."
    healthBlocks[len(healthBlocks) - 1].Type = "info"
    health3 := prepareArray([]int{fc.A1, fc.B1, fc.L2})
    for _, item := range health3 {
      // healthBlocks = append(healthBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal))
      healthBlocks[len(healthBlocks) - 1].Content = fmt.Sprintf("%s %s", healthBlocks[len(healthBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal).Content)
    }
    healthBlocks = append(healthBlocks, getAnswerFromTable(db, "4", 12, languageShort, gender, personal))
    healthBlocks[len(healthBlocks) - 1].Title = "Heart, circulatory system, respiratory system, lungs, bronchi, thoracic spine, ribs, shoulder scapular area, chest."
    healthBlocks[len(healthBlocks) - 1].Type = "info"
    health4 := prepareArray([]int{fc.A3, fc.B3, fc.L3})
    for _, item := range health4 {
      // healthBlocks = append(healthBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", item), 12, languageShort, gender, personal))
      healthBlocks[len(healthBlocks) - 1].Content = fmt.Sprintf("%s %s", healthBlocks[len(healthBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal).Content)
    }
    healthBlocks = append(healthBlocks, getAnswerFromTable(db, "5", 12, languageShort, gender, personal))
    healthBlocks[len(healthBlocks) - 1].Title = "Gastrointestinal tract, abdominal organs, pancreas, spleen, liver, gallbladder, small intestine, central part of the spine."
    healthBlocks[len(healthBlocks) - 1].Type = "info"
    health5 := prepareArray([]int{fc.E, fc.L4})
    for _, item := range health5 {
      // healthBlocks = append(healthBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal))
      healthBlocks[len(healthBlocks) - 1].Content = fmt.Sprintf("%s %s", healthBlocks[len(healthBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal).Content)
    }
    healthBlocks = append(healthBlocks, getAnswerFromTable(db, "6", 12, languageShort, gender, personal))
    healthBlocks[len(healthBlocks) - 1].Title = "drenal glands, uterus and ovaries, kidneys, intestines, prostate gland in men, lumbar spinal column."
    healthBlocks[len(healthBlocks) - 1].Type = "info"
    health6 := prepareArray([]int{fc.D1, fc.C1, fc.L5})
    for _, item := range health6 {
      // healthBlocks = append(healthBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal))
      healthBlocks[len(healthBlocks) - 1].Content = fmt.Sprintf("%s %s", healthBlocks[len(healthBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal).Content)
    }
    healthBlocks = append(healthBlocks, getAnswerFromTable(db, "7", 12, languageShort, gender, personal))
    healthBlocks[len(healthBlocks) - 1].Title = "Urogenital system, lower limbs, large intestine, tailbone, sacrum, legs."
    healthBlocks[len(healthBlocks) - 1].Type = "info"
    health7 := prepareArray([]int{fc.D, fc.C, fc.L6})
    for _, item := range health7 {
      // healthBlocks = append(healthBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal))
      healthBlocks[len(healthBlocks) - 1].Content = fmt.Sprintf("%s %s", healthBlocks[len(healthBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal).Content)
    }
    healthBlocks = append(healthBlocks, getAnswerFromTable(db, "8", 12, languageShort, gender, personal))
    healthBlocks[len(healthBlocks) - 1].Title = "Circulatory system, nervous system, lymphatic system, immune system, those organs that are found throughout the body, general failure of the body."
    healthBlocks[len(healthBlocks) - 1].Type = "info"
    health8 := prepareArray([]int{fc.D3, fc.C3, fc.E3})
    for _, item := range health8 {
      // healthBlocks = append(healthBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal))
      healthBlocks[len(healthBlocks) - 1].Content = fmt.Sprintf("%s %s", healthBlocks[len(healthBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", item), 225, languageShort, gender, personal).Content)
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
    lifeGuideBlocks[len(lifeGuideBlocks) - 1].Title = "Life Guidance"
    lifeGuideBlocks[len(lifeGuideBlocks) - 1].Type = "info"
    if fc.A != fc.B {
      // lifeGuideBlocks = append(lifeGuideBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.B), 11, languageShort, gender, personal))
      lifeGuideBlocks[len(lifeGuideBlocks) - 1].Content = fmt.Sprintf("%s %s", lifeGuideBlocks[len(lifeGuideBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.B), 11, languageShort, gender, personal).Content)
    }
    if fc.A != fc.E && fc.B != fc.E {
      // lifeGuideBlocks = append(lifeGuideBlocks, getAnswerFromTable(db, fmt.Sprintf("%d", fc.E), 11, languageShort, gender, personal))
      lifeGuideBlocks[len(lifeGuideBlocks) - 1].Content = fmt.Sprintf("%s %s", lifeGuideBlocks[len(lifeGuideBlocks) - 1].Content, getAnswerFromTable(db, fmt.Sprintf("%d", fc.E), 11, languageShort, gender, personal).Content)
    }
    lifeguide := Prediction{}
    lifeguide.Title = "Life Guidance"
    lifeguide.ImageName = "life_guidance"
    lifeguide.Blocks = lifeGuideBlocks
    lifeguide.BlockType = "default"
    // ------ LIFE GUIDE END --------
    // ------ YEAR FORECAST BEGIN --------
    // layout := "01/02/2006"
    // t := time.Now()
    // dayNow := t.Year()
    // monthNow := int(t.Month())
    // yearNow := t.Day()
    // age := time.Parse(layout, fmt.Sprintf("%d/%d/%d", dayNow, monthNow, yearNow))
    // diff := age.Unix() - timestamp
    // years := diff / (60 * 60 * 24 * 365)
    // counter := 0
    // ------ YEAR FORECAST END --------
    predictions := []Prediction{personalfeatures, destiny, money, programs, sexiness, pastlife, parents, kids, relationships, health, lifeguide}
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
  fmt.Println(icombo[0])
  fmt.Println(icombo[1])
  fmt.Println(icombo[2])
  fmt.Println(d)
  fmt.Println(e)
  fmt.Println(a1)
  fmt.Println(a2)
  fmt.Println(b1)
  fmt.Println(b2)
  fmt.Println(c1)
  fmt.Println(c2)
  fmt.Println(d1)
  fmt.Println(d2)
  fmt.Println(x)
  fmt.Println(x1)
  fmt.Println(x2)
  fmt.Println(f)
  fmt.Println(g)
  fmt.Println(y)
  fmt.Println(k)
  fmt.Println(e1)
  fmt.Println(e2)
  fmt.Println(o)
  fmt.Println(u)
  fmt.Println(h)
  fmt.Println(j)
  fmt.Println(m)
  fmt.Println(n)
  fmt.Println(t)
  fmt.Println(z)
  fmt.Println(s)
  fmt.Println(b3)
  fmt.Println(a3)
  fmt.Println(l)
  fmt.Println(l1)
  fmt.Println(l2)
  fmt.Println(l3)
  fmt.Println(l4)
  fmt.Println(l5)
  fmt.Println(l6)
  fmt.Println(d3)
  fmt.Println(c3)
  fmt.Println(e3)
  return Combos{icombo[0],icombo[1],icombo[2],d,e,a1,a2,b1,b2,c1,c2,d1,d2,x,x1,x2,f,g,y,k,e1,e2,o,u,h,j,m,n,t,z,s,b3,a3,l,l1,l2,l3,l4,l5,l6,d3,c3,e3}
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
