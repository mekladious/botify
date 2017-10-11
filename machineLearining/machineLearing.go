package machineLearining

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	db_uri string = "mongodb://admin:admin@ds161833.mlab.com:61833/botify"
)

type (
	entry struct {
		Question string `bson:"question"`
		Answer   string `bson:"answer"`
	}
)

func UnknownAnswer(message string) string {

	db, err := mgo.Dial(db_uri)
	defer func() {
		db.Close()
	}()
	if err != nil {
		fmt.Println(err)
	}
	collection := db.DB("botify").C("test")
	// e := entry{Question: "hi"}
	// err = collection.Insert(&e)
	var results []entry
	collection.Find(bson.M{"question": message}).Select(bson.M{"answer": 1}).All(&results)
	//collection.Find(nil).All(&results)
	if len(results) == 0 {
		e := entry{Question: message}
		err = collection.Insert(&e)
		return ""
	}
	return results[0].Answer
}
