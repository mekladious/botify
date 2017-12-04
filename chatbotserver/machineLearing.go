package chatbot

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/mgo.v2"
)

const (
	db_uri string = "mongodb://admin:admin@ds161833.mlab.com:61833/botify"
)

type (
	entry struct {
		Question string `bson:"question"`
		Answer   string `bson:"answer"`
	}
	acceptedAnswer struct {
		Answer     string
		Percentage float32
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
	collection := db.DB("botify").C("MachineLearning")
	// e := entry{Question: "hi"}
	// err = collection.Insert(&e)
	var results []entry
	//collection.Find(bson.M{"question": message}).Select(bson.M{"answer": 1}).All(&results)
	collection.Find(nil).All(&results)
	if len(results) == 0 {
		e := entry{Question: message}
		err = collection.Insert(&e)
		return ""
	}

	//var acceptedAnswers []acceptedAnswer
	for _, result := range results {
		result.Question = strings.ToLower(result.Question)
		answer := result.Answer
		var replaced_var_index []int
		contain_var := false
		var_name := ""
		if strings.Contains(result.Question, "_") {
			//replacing any variable _v_ in db question with regex to be compared in the next step with the input question
			reg := regexp.MustCompile(`_[a-z0-9]+_`)
			replaced_var_index = reg.FindIndex([]byte(result.Question))
			var_name = reg.FindString(result.Question)
			regex_db_question := string(reg.ReplaceAll([]byte(result.Question), []byte("[a-z]{3,}")))
			result.Question = regex_db_question
			contain_var = true
		}

		//cheking message with DB question
		r := regexp.MustCompile(result.Question)
		if r.MatchString(message) {
			if contain_var { // get variable from question
				//answer = result.Question[0:replaced_var_index[0]]
				variable_value := ""
				for i := replaced_var_index[0]; i <= replaced_var_index[0]+10 && i < len(message); i++ {
					if string(message[i]) != " " {
						variable_value += string(message[i])
					} else {
						break
					}
				}
				answer = strings.Replace(answer, var_name, variable_value, -1)
			} else {
				answer = result.Answer
			}
			return answer
		}

	}
	// for _, aa := range acceptedAnswers {
	// 	fmt.Println(aa.Answer)
	// }
	return ""
}

func checkForSymbols(answer string, session Session, uuid string) (string, string, string, string, error) {
	if strings.Contains(answer, "#Get_featured_playlists") {
		answer = strings.Replace(answer, "#Get_featured_playlists", Get_featured_playlists(), -1)
	} else if strings.Contains(answer, "=") { // = means that the message was equivalent to some message go process the equivalent message and return the result

		message, images, tracks, alarmTime, err := sampleProcessor(session, after(answer, "="), uuid)
		return message, images, tracks, alarmTime, err
	}

	return answer, "", "", "", nil
}
