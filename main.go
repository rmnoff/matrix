package main

import (
  "os"
  "flag"
  "log"
  "fmt"
  "time"
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


type ConstantText struct {
  Header string `json:"header"`
  Content string `json:"content"`
}

type Prediction struct {
  ImageName string `json:"ImageName"`
  Title string `json:"Title"`
  BlockType string `json:"BlockType"`
  Blocks []Block `json:"Blocks"`
}

type Block struct {
  Id int `db:"id"`
  Content string `db:"content"`
  Created sql.NullString `db:"created"`
  Edited sql.NullString `db:"edited"`
  PredType int `db:"type_id"`
  Lang int `db:"lang_id"`
  Personal bool `db:"personal"`
  Type string `json:Type`
  Title string `json:Title`
  TintColor *string `json:TintColor`
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
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('relationship') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('parents') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('parents IMPORTANT') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('kids') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('destiny') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('destiny common') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('past life') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('programms') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('life guide') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('health') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('health recommendation') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('year prediction') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('19-7') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('sexuality') ON CONFLICT DO NOTHING;`)
  tx.MustExec(`INSERT INTO predictionType(name) VALUES('lessons from children') ON CONFLICT DO NOTHING;`)
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
		return c.Write(`{"ok": true, "error": null, "data": "Future authorisation end"}`)
  })











  api.Get("/check/<input>", func(c *routing.Context) error {
    input := c.Param("input")
    timestamp := input[:len(input) - 4]
    gender := string(input[len(input) - 4])
    personal := string(input[len(input) - 3])
    languageShort := input[len(input) - 2:]
    fmt.Println(timestamp, gender, personal, languageShort)
    if timestamp == "" {
      marshalled, _ := json.Marshal(Response{false, "No timestamp provided", nil})
      return c.Write(marshalled)
    }
    combo := countBD(timestamp)
    if combo == nil {
      marshalled, _ := json.Marshal(Response{false, "Timestamp corrupted", nil})
      return c.Write(marshalled)
    }
    finalCombos := setAllCombos(combo)
    // prog1 := fmt.Sprintf("prog1: [%v %v %v]", combo[0], finalCombos[3], finalCombos[2])
    // prog2 := fmt.Sprintf("prog2: [%v %v %v]", combo[1], finalCombos[5], finalCombos[4])
    // prog3 := fmt.Sprintf("prog3: [%v %v %v]", combo[2], finalCombos[7], finalCombos[6])
    // prog4 := fmt.Sprintf("prog4: [%v %v %v]", finalCombos[13], finalCombos[15], finalCombos[19])
    // prog5 := fmt.Sprintf("prog5: [%v %v %v]", finalCombos[16], finalCombos[14], finalCombos[20])

    isPersonal := false
    if personal == "p" {
      isPersonal = true
    }

    language := Language{}
    err = db.Get(&language, "SELECT * FROM predictionlang WHERE short = $1", languageShort)
    if err != nil {
      language.Id = 1
    }

    pastLife := Prediction{}
    pastLifeBlock := Block{}
    pastLifeBlock.Type = "info"
    pastLifeBlock.Title = "Previous Life Common"
    pastLifePredictionCombo := fmt.Sprintf("%d-%d-%d", finalCombos[8], finalCombos[9], finalCombos[0])
    fmt.Println(pastLifePredictionCombo)
    err = db.Get(&pastLifeBlock, "SELECT * FROM prediction WHERE type_id=9 AND personal=$2 AND lang_id=$3 AND id IN(SELECT prediction_id FROM predictionrel WHERE combination=$1)", pastLifePredictionCombo, isPersonal, language.Id)
    if err != nil {
      log.Println(err)
      marshalled, _ := json.Marshal(Response{false, "Can't parse past life prediction", nil})
      return c.Write(marshalled)
    }
    pastLife.Title = "Previous Life"
    pastLife.Blocks = []Block{pastLifeBlock}
    pastLife.BlockType = "default"

    personalFeatures := Prediction{}
    personalFeaturesPos := Block{}
    personalFeaturesPosSecond := Block{}
    personalFeaturesPos.Type = "info"
    personalFeaturesPos.Title = "Personal Features Positive"
    personalFeaturesPosCombo := fmt.Sprintf("%d", combo[0])
    personalFeaturesPosSecondCombo := fmt.Sprintf("%d", combo[1])
    personalFeaturesNeg := Block{}
    personalFeaturesNegSecond := Block{}
    personalFeaturesNeg.Type = "info"
    personalFeaturesNeg.Title = "Personal Features Negative"
    personalFeaturesNegCombo := fmt.Sprintf("%d", combo[0])
    personalFeaturesNegSecondCombo := fmt.Sprintf("%d", combo[1])
    personalFeaturesSoc := Block{}
    personalFeaturesSoc.Type = "info"
    personalFeaturesSoc.Title = "Personal Features Social"
    personalFeaturesSocCombo := fmt.Sprintf("%d", finalCombos[1])
    err = db.Get(&personalFeaturesPos, "SELECT * FROM prediction WHERE type_id=1 AND personal=$2 AND lang_id=$3 AND id=(SELECT prediction_id FROM predictionrel WHERE combination=$1)", personalFeaturesPosCombo, isPersonal, language.Id)
    if err != nil {
      log.Println(err)
      marshalled, _ := json.Marshal(Response{false, "Can't parse positive personal features 1", nil})
      return c.Write(marshalled)
    }
    err = db.Get(&personalFeaturesPosSecond, "SELECT * FROM prediction WHERE type_id=1 AND personal=$2 AND lang_id=$3 AND id=(SELECT prediction_id FROM predictionrel WHERE combination=$1)", personalFeaturesPosSecondCombo, isPersonal, language.Id)
    if err != nil {
      log.Println(err)
      marshalled, _ := json.Marshal(Response{false, "Can't parse positive personal features 2", nil})
      return c.Write(marshalled)
    }
    personalFeaturesPos.Content = fmt.Sprintf("%s %s", personalFeaturesPos.Content, personalFeaturesPosSecond.Content)
    err = db.Get(&personalFeaturesNeg, "SELECT * FROM prediction WHERE type_id=2 AND personal=$2 AND lang_id=$3 AND id=(SELECT prediction_id FROM predictionrel WHERE combination=$1)", personalFeaturesNegCombo, isPersonal, language.Id)
    if err != nil {
      log.Println(err)
      marshalled, _ := json.Marshal(Response{false, "Can't parse negative personal features 1", nil})
      return c.Write(marshalled)
    }
    err = db.Get(&personalFeaturesNegSecond, "SELECT * FROM prediction WHERE type_id=2 AND personal=$2 AND lang_id=$3 AND id=(SELECT prediction_id FROM predictionrel WHERE combination=$1)", personalFeaturesNegSecondCombo, isPersonal, language.Id)
    if err != nil {
      log.Println(err)
      marshalled, _ := json.Marshal(Response{false, "Can't parse negative personal features 2", nil})
      return c.Write(marshalled)
    }
    personalFeaturesNeg.Content = fmt.Sprintf("%s %s", personalFeaturesNeg.Content, personalFeaturesNegSecond.Content)
    err = db.Get(&personalFeaturesSoc, "SELECT * FROM prediction WHERE type_id=3 AND personal=$2 AND lang_id=$3 AND id=(SELECT prediction_id FROM predictionrel WHERE combination=$1)", personalFeaturesSocCombo, isPersonal, language.Id)
    if err != nil {
      log.Println(err)
      marshalled, _ := json.Marshal(Response{false, "Can't parse social personal features", nil})
      return c.Write(marshalled)
    }
    personalFeatures.Title = "Personal Features"
    personalFeatures.Blocks = []Block{personalFeaturesPos,personalFeaturesNeg,personalFeaturesSoc}
    personalFeatures.BlockType = "default"

    relationship := Prediction{}
    relationshipBlock := Block{}
    relationshipBlockSecond := Block{}
    relationshipBlockThird := Block{}
    relationshipBlock.Type = "info"
    relationshipBlock.Title = "Relationship Common"
    relationshipCombo := fmt.Sprintf("%d", finalCombos[11])
    relationshipComboSecond := fmt.Sprintf("%d", finalCombos[8])
    relationshipComboThird := fmt.Sprintf("%d", finalCombos[10])
    fmt.Println(relationshipCombo)
    err = db.Get(&relationshipBlock, "SELECT * FROM prediction WHERE type_id=5 AND personal=$2 AND lang_id=$3 AND id=(SELECT prediction_id FROM predictionrel WHERE combination=$1)", relationshipCombo, isPersonal, language.Id)
    if err != nil {
      log.Println(err)
      marshalled, _ := json.Marshal(Response{false, "Can't parse relationship prediction", nil})
      return c.Write(marshalled)
    }
    err = db.Get(&relationshipBlockSecond, "SELECT * FROM prediction WHERE type_id=5 AND personal=$2 AND lang_id=$3 AND id=(SELECT prediction_id FROM predictionrel WHERE combination=$1)", relationshipComboSecond, isPersonal, language.Id)
    if err != nil {
      log.Println(err)
      marshalled, _ := json.Marshal(Response{false, "Can't parse relationship prediction", nil})
      return c.Write(marshalled)
    }
    err = db.Get(&relationshipBlockThird, "SELECT * FROM prediction WHERE type_id=5 AND personal=$2 AND lang_id=$3 AND id=(SELECT prediction_id FROM predictionrel WHERE combination=$1)", relationshipComboThird, isPersonal, language.Id)
    if err != nil {
      log.Println(err)
      marshalled, _ := json.Marshal(Response{false, "Can't parse relationship prediction", nil})
      return c.Write(marshalled)
    }
    relationshipBlock.Content = fmt.Sprintf("%s %s %s", relationshipBlock.Content, relationshipBlockSecond.Content, relationshipBlockThird.Content)
    relationship.Title = "Relationship"
    relationship.Blocks = []Block{relationshipBlock}
    relationship.BlockType = "default"

    lifeGuide := Prediction{}
    lifeGuideBlock := Block{}
    lifeGuideBlockSecond := Block{}
    lifeGuideBlockThird := Block{}
    lifeGuideBlock.Type = "info"
    lifeGuideBlock.Title = "Life Guide Common"
    lifeGuideCombo := fmt.Sprintf("%d", combo[0])
    lifeGuideSecondCombo := fmt.Sprintf("%d", combo[1])
    lifeGuideThirdCombo := fmt.Sprintf("%d", finalCombos[1])
    fmt.Println(lifeGuideCombo)
    err = db.Get(&lifeGuideBlock, "SELECT * FROM prediction WHERE type_id=11 AND personal=$2 AND lang_id=$3 AND id=(SELECT prediction_id FROM predictionrel WHERE combination=$1)", lifeGuideCombo, isPersonal, language.Id)
    if err != nil {
      log.Println(err)
      marshalled, _ := json.Marshal(Response{false, "Can't parse life guide prediction", nil})
      return c.Write(marshalled)
    }
    err = db.Get(&lifeGuideBlockSecond, "SELECT * FROM prediction WHERE type_id=11 AND personal=$2 AND lang_id=$3 AND id=(SELECT prediction_id FROM predictionrel WHERE combination=$1)", lifeGuideSecondCombo, isPersonal, language.Id)
    if err != nil {
      log.Println(err)
      marshalled, _ := json.Marshal(Response{false, "Can't parse life guide prediction", nil})
      return c.Write(marshalled)
    }
    err = db.Get(&lifeGuideBlockThird, "SELECT * FROM prediction WHERE type_id=11 AND personal=$2 AND lang_id=$3 AND id=(SELECT prediction_id FROM predictionrel WHERE combination=$1)", lifeGuideThirdCombo, isPersonal, language.Id)
    if err != nil {
      log.Println(err)
      marshalled, _ := json.Marshal(Response{false, "Can't parse life guide prediction", nil})
      return c.Write(marshalled)
    }
    lifeGuideBlock.Content = fmt.Sprintf("%s %s %s", lifeGuideBlock.Content, lifeGuideBlockSecond.Content, lifeGuideBlockThird.Content)
    lifeGuide.Title = "Life Guide"
    lifeGuide.Blocks = []Block{lifeGuideBlock}
    lifeGuide.BlockType = "default"

    sex := Prediction{}
    sexBlock := Block{}
    sexBlock.Type = "info"
    sexBlock.Title = "Sexiness Common"
    sexCombo := fmt.Sprintf("%d-%d-%d", finalCombos[1], finalCombos[17], finalCombos[18])
    fmt.Println(sexCombo)
    err = db.Get(&sexBlock, "SELECT * FROM prediction WHERE type_id=228 AND personal=$2 AND lang_id=$3 AND id=(SELECT prediction_id FROM predictionrel WHERE combination=$1)", sexCombo, isPersonal, language.Id)
    if err != nil {
      log.Println(err)
      marshalled, _ := json.Marshal(Response{false, "Can't parse life guide prediction", nil})
      return c.Write(marshalled)
    }
    sex.Title = "Sexiness"
    sex.Blocks = []Block{sexBlock}
    sex.BlockType = "default"

    destiny := Prediction{}
    destinyBlock := Block{}
    destinyBlock.Type = "info"
    destinyBlock.Title = "Destiny Common"
    destinyCombo := fmt.Sprintf("%d-%d-%d", finalCombos[21], finalCombos[22], finalCombos[23])
    fmt.Println(destinyCombo)
    err = db.Get(&destinyBlock, "SELECT * FROM prediction WHERE type_id=228 AND personal=$2 AND lang_id=$3 AND id=(SELECT prediction_id FROM predictionrel WHERE combination=$1)", destinyCombo, isPersonal, language.Id)
    if err != nil {
      log.Println(err)
      marshalled, _ := json.Marshal(Response{false, "Can't parse life guide prediction", nil})
      return c.Write(marshalled)
    }
    destinyCommonBlock := Block{}
    destinyCommonBlock.Type = "info"
    destinyCommonBlock.Title = "Destiny Common"
    destinyCommonCombo := fmt.Sprintf("%d-%d-%d", finalCombos[21], finalCombos[22], finalCombos[23])
    fmt.Println(destinyCommonCombo)
    err = db.Get(&destinyCommonBlock, "SELECT * FROM prediction WHERE type_id=228 AND personal=$2 AND lang_id=$3 AND id=(SELECT prediction_id FROM predictionrel WHERE combination=$1)", destinyCommonCombo, isPersonal, language.Id)
    if err != nil {
      log.Println(err)
      marshalled, _ := json.Marshal(Response{false, "Can't parse life guide prediction", nil})
      return c.Write(marshalled)
    }
    destiny.Title = "Destiny"
    destiny.Blocks = []Block{destinyBlock, destinyCommonBlock}
    destiny.BlockType = "default"

    predictions := []Prediction{pastLife, personalFeatures, relationship, lifeGuide, sex, destiny}

    for _, pred := range predictions {
      for _, content := range pred.Blocks {
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
    }
    marshalled, _ := json.Marshal(Response{true, "", predictions})
    return c.Write(marshalled)
  })













  api.Get("/check/template/<input>", func(c *routing.Context) error {
    c.SetContentType("application/json; charset=utf8")
    return c.Write(`{
      "ok": true,
      "data": [
        {
          "imageName": "personal_features",
          "title": "Personal feautures",
          "blockType": "default",
          "blocks": [
            {
              "type": "expandable",
              "title": "В позитивном состоянии",
              "text": "Текст про позитивное состояние человека и небольшой ответ.",
              "tintColor": null
            },
            {
              "type": "expandable",
              "title": "В негативном состоянии",
              "text": "Большой и объемный ответ, где раскрываются особенности человека в негативном состоянии. Чего стоит избегать, а что стоит в таком состоянии делать - все рассказывается в этом пункте раздела персональных особенностей.",
              "tintColor": null
            },
            {
              "type": "expandable",
              "title": "Общение",
              "text": "Каков человек в общении? Расскажет этот раздел, как только база данных будет заполнена.",
              "tintColor": null
            }
          ]
        },
        {
          "imageName": "destiny",
          "title": "Destiny",
          "blockType": "default",
          "blocks": [
            {
              "type": "info",
              "title": "Общая информация",
              "text": "Здесь рассказывается о предназначении человека непосредственно в той жизни, которой он сейчас живет",
              "tintColor": null
            },
            {
              "type": "expandable",
              "title": "Важно",
              "text": "Отдельный пункт \"важно\", где рассказывается о какой-то особенности человека",
              "tintColor": null
            }
          ]
        },
        {
          "imageName": "children",
          "title": "Children",
          "blockType": "default",
          "blocks": [
            {
              "type": "expandable",
              "title": "Важно",
              "text": "Здесь иногда бывает текст, но нечасто, потому что не у всех есть",
              "tintColor": null
            },
            {
              "type": "info",
              "title": "Общее",
              "text": "Общая информация раздела дети",
              "tintColor": null
            }
          ]
        },
        {
          "imageName": "life_guidance",
          "title": "Life Guidance",
          "blockType": "default",
          "blocks": [
            {
              "type": "info",
              "title": "Общее",
              "text": "Общая информация раздела уроки жизни",
              "tintColor": null
            }
          ]
        },
        {
          "imageName": "money",
          "title": "Money",
          "blockType": "default",
          "blocks": [
            {
              "type": "info",
              "title": "Общее",
              "text": "Общая информация раздела деньги",
              "tintColor": null
            },
            {
              "type": "info",
              "title": "Для достижения успеха",
              "text": "Общая информация раздела деньги, подраздел для достижения успеха",
              "tintColor": null
            },
            {
              "type": "expandable",
              "title": "Важно",
              "text": "Общая информация раздела деньги, подраздел важно",
              "tintColor": null
            }
          ]
        },
        {
          "imageName": "parents",
          "title": "Parents",
          "blockType": "default",
          "blocks": [
            {
              "type": "info",
              "title": "Общее",
              "text": "Общая информация раздела родители",
              "tintColor": null
            },
            {
              "type": "info",
              "title": "Обиды на родителей",
              "text": "Общая информация раздела родители, подраздел обиды на родителей",
              "tintColor": null
            },
            {
              "type": "expandable",
              "title": "Важно",
              "text": "Общая информация раздела родители, подраздел важно",
              "tintColor": null
            },
            {
              "type": "expandable",
              "title": "Важно",
              "text": "Общая информация раздела родители, подраздел программы",
              "tintColor": null
            }
          ]
        },
        {
          "imageName": "previous_life",
          "title": "Previous Life",
          "blockType": "default",
          "blocks": [
            {
              "type": "info",
              "title": "Общее",
              "text": "Общая информация раздела прошлая жизнь. Если в этот раздел нужно добавить еще категории - напишите Саше.",
              "tintColor": null

            }
          ]
        },
        {
          "imageName": "programs",
          "title": "Programs",
          "blockType": "default",
          "blocks": [
            {
              "type": "info",
              "title": "Программа 1",
              "text": "Общая информация по разделу программы",
              "tintColor": null
            },
            {
              "type": "info",
              "title": "Программа 2",
              "text": "Общая информация по разделу программы",
              "tintColor": null
            },
            {
              "type": "info",
              "title": "Программа 3",
              "text": "Общая информация по разделу программы",
              "tintColor": null
            },
            {
              "type": "info",
              "title": "Программа 4",
              "text": "Общая информация по разделу программы",
              "tintColor": null
            }
          ]
        },
        {
          "imageName": "relationships",
          "title": "Relationships",
          "blockType": "default",
          "blocks": [
            {
              "type": "info",
              "title": "Общее",
              "text": "Общая информация раздела отношения. Если в этот раздел нужно добавить еще категории - напишите Саше.",
              "tintColor": null

            }
          ]
        },
        {
          "imageName": "sexiness",
          "title": "Sexiness",
          "blockType": "default",
          "blocks": [
            {
              "type": "info",
              "title": "Общее",
              "text": "Общая информация раздела сексуальность. Если в этот раздел нужно добавить еще категории - напишите Саше.",
              "tintColor": null

            }
          ]
        },
        {
          "imageName": "health",
          "title": "Health",
          "blockType": "health",
          "blocks": [
            {
              "type": "info",
              "title": "Общее",
              "text": "Общая информация раздела здоровье. Если в этот раздел нужно добавить еще категории - напишите Саше.",
              "tintColor": null
            }
          ]
        }
      ]
    }
    `)
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
