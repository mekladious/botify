package chatbot

import (
	"regexp"
	"strings"
)

var (
	// WelcomeMessage A constant to hold the welcome message
	WelcomeMessage = "Hello, Botify is ready to inspire you ;)"

	// sessions = {
	//   "uuid1" = Session{...},
	//   ...
	// }
	sessions = map[string]Session{}

	processor = sampleProcessor

	SpotifyAuthorizationToken string
)

type (
	// Session Holds info about a session
	Session map[string]interface{}

	// JSON Holds a JSON object
	JSON map[string]interface{}

	// Processor Alias for Process func
	Processor func(session Session, message string) (string, error)
)

func sampleProcessor(session Session, message string, uuid string) (string, string, error) {
	message = strings.ToLower(message)
	if strings.Contains(message, "featured playlists") {
		featuredPlaylists := Get_featured_playlists()
		return featuredPlaylists, "", nil
	} else if strings.Contains(message, "alarm") {
		if strings.Contains(message, "show") {
			alarms := GetAlarms(uuid)
			return alarms, "", nil
		}
		if strings.Contains(message, "delete") {
			deleteAlarm := DeleteAlarm(uuid, after(message, ": "))
			return deleteAlarm, "", nil
		}
		singerName := between(message, "want", "to")
		alarmTime := after(message, "at")
		if singerName == "" || alarmTime == "" {
			return "please use the format 'i want (artist name) to alarm me at (time hh:mm)'", "", nil
		}
		tracks, er := Get_artist_tracks(singerName)
		if er != nil {
			return er.Error(), "", nil
		}
		err := InsertAlarmInGoogleCalendar(alarmTime, uuid, tracks, singerName)
		reply := "Done, Alarm is set. " + singerName + " will wake you up at " + alarmTime + "."
		if err != "" {
			reply = err
		}

		return reply, "", nil
	} else if strings.Contains(message, "i am") {
		mood := after(message, "i am")
		if strings.Contains(mood, "happy") || strings.Contains(mood, "excite") || strings.Contains(mood, "cheerful") {
			Moody := Get_mood("party")
			return Moody, "", nil
		} else if strings.Contains(mood, "tired") || strings.Contains(mood, "chilling") || strings.Contains(mood, "bored") || strings.Contains(mood, "stressed") {
			Moody := Get_mood("chill")
			return Moody, "", nil
		} else if strings.Contains(mood, "moody") || strings.Contains(mood, "unstable") {
			Moody := Get_mood("mood")
			return Moody, "", nil
		} else if strings.Contains(mood, "angry") || strings.Contains(mood, "furious") || strings.Contains(mood, "annoyed") || strings.Contains(mood, "sad") {
			Moody := Get_mood("rock")
			return Moody, "", nil
		} else if strings.Contains(mood, "excercising") || strings.Contains(mood, "working out") || strings.Contains(mood, "training") {
			Moody := Get_mood("workout")
			return Moody, "", nil
		} else if strings.Contains(mood, "studying") || strings.Contains(mood, "thinking") || strings.Contains(mood, "wrapped up") {
			Moody := Get_mood("focus")
			return Moody, "", nil
		} else if strings.Contains(mood, "sleep") || strings.Contains(mood, "drowsy") || strings.Contains(mood, "exhausted") {
			Moody := Get_mood("sleep")
			return Moody, "", nil
		} else if strings.Contains(mood, "love") || strings.Contains(mood, "affectionate") {
			Moody := Get_mood("romance")
			return Moody, "", nil
		} else if strings.Contains(mood, "travelling") || strings.Contains(mood, "on the road") {
			Moody := Get_mood("travel")
			return Moody, "", nil
		} else if strings.Contains(mood, "play") || strings.Contains(mood, "fun") {
			Moody := Get_mood("gaming")
			return Moody, "", nil
		} else if strings.Contains(mood, "going out") || strings.Contains(mood, "laugh") {
			Moody := Get_mood("comedy")
			return Moody, "", nil
		} else {
			return "unknown mood try happy,bored, love,tired, moody, sad, excercising, studying, sleepy, travelling,playing,laugh,  ...", "", nil
		}

	} else if strings.Contains(message, "info of") {
		artist := after(message, "info of")
		info, images := Get_artist_info(artist)
		return info, images, nil
	} else if strings.Contains(strings.ToLower(message), "search") || strings.Contains(strings.ToLower(message), "play") {
		message = strings.Replace(message, "search", "", -1)
		message = strings.Replace(message, "play", "", -1)
		message = strings.Replace(message, " ", "+", -1)
		reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
		message = reg.ReplaceAllString(message, "")

		results, images := search(message)

		if results != "null" {
			return results, images, nil
		} else {
			return "No results were found for your search, please try again", "", nil
		}

	} else if strings.Contains(strings.ToLower(message), "new") && strings.Contains(strings.ToLower(message), "release") {
		newReleases := get_new_releases()
		return newReleases, "", nil
	} else if strings.Contains(strings.ToLower(message), "favorite") {
		if strings.Contains(strings.ToLower(message), "add") {
			//trackId := "7c0XG5cIJTrrAgEC3ULPiq" //dummy data
			trackName := after(message, ":")
			trackid := getTrackID(trackName)
			if trackid == "nil" {
				return "track not found", "", nil
			}
			res, err := add_to_favorites(uuid, trackid, trackName)
			return res, "", err
		} else if strings.Contains(strings.ToLower(message), "show") {
			res, err := get_favorites(uuid)
			return res, "", err
		} else if strings.Contains(strings.ToLower(message), "delete") {
			res, err := delete_favorite(uuid, after(message, ":"))
			return res, "", err
		} else {
			return "", "supported functions: add, show, delete", nil
		}
	} else {
		result := checkForSymbols(UnknownAnswer(message))
		if result != "" {
			return result, "", nil
		}
	}

	return "Sorry I didn't understand you .. For now you can get featured playlists and new releases.. more features coming soon", "", nil

	// // Make sure a history key is defined in the session which points to a slice of strings
	// _, historyFound := session["history"]
	// if !historyFound {
	// 	session["history"] = []string{}
	// }

	// // Fetch the history from session and cast it to an array of strings
	// history, _ := session["history"].([]string)

	// // Make sure the message is unique in history
	// for _, m := range history {
	// 	if strings.EqualFold(m, message) {
	// 		return "", fmt.Errorf("You've already ordered %s before!", message)
	// 	}
	// }

	// // Add the message in the parsed body to the messages in the session
	// history = append(history, message)

	// // Form a sentence out of the history in the form Message 1, Message 2, and Message 3
	// l := len(history)
	// wordsForSentence := make([]string, l)
	// copy(wordsForSentence, history)
	// if l > 1 {
	// 	wordsForSentence[l-1] = "and " + wordsForSentence[l-1]
	// }
	// sentence := strings.Join(wordsForSentence, ", ")

	// //Save the updated history to the session
	// session["history"] = history

	// return fmt.Sprintf("So, you want %s! What else?", strings.ToLower(sentence)), nil
}
