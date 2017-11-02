package chatbot

import (
	"strings"
)

var (
	// WelcomeMessage A constant to hold the welcome message
	WelcomeMessage = "Welcome, what do you want to order?"

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

func sampleProcessor(session Session, message string) (string, error) {
	if strings.Contains(strings.ToLower(message), "featured playlists") {
		featuredPlaylists := Get_featured_playlists()
		return featuredPlaylists, nil
	} else {
		result := checkForSymbols(UnknownAnswer(message))
		if result != "" {
			return result, nil
		}
	}
	if strings.Contains(strings.ToLower(message), "new") && strings.Contains(strings.ToLower(message), "release") {
		newReleases := get_new_releases()
		return newReleases, nil
	}
	return "Sorry I didn't understand you .. For now you can get featured playlists and new releases.. more features coming soon", nil

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
