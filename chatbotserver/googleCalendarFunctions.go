package chatbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

var srv *calendar.Service

//InsertAlarmInGoogleCalendar takes time in format hh:mm AM and user uuid and tracks to be played when alarm is on
func InsertAlarmInGoogleCalendar(alarmTime string, uuid string, tracks string) string {
	connectGoogleCalendar()
	now := time.Now()

	hour := now.Hour()
	minute := now.Minute()

	InputHour, err := strconv.Atoi(between(alarmTime, " ", ":"))
	if err != nil {
		log.Fatalf("Unable to create alarm. %v \n", err)
		return "Unable to create alarm."
	}
	InputMinute, err := strconv.Atoi(after(alarmTime, ":"))
	if err != nil {
		log.Fatalf("Unable to create alarm. %v \n", err)
		return "Unable to create alarm."
	}
	// fmt.Println("input hour", before(alarmTime, ":"))
	// fmt.Println("input hour", InputHour)
	// fmt.Println("now hour", hour)
	nextDay := false

	//if alarm time is less than now then it means that the alarm is on the next day
	if InputHour < hour {
		nextDay = true
	} else if (InputHour == hour) && (InputMinute < minute) {
		nextDay = true
	}

	if nextDay {
		now = now.AddDate(0, 0, 1)
	}

	alarm := &calendar.Event{
		Summary:     uuid,
		Description: tracks,
		Start: &calendar.EventDateTime{
			DateTime: time.Date(now.Year(), now.Month(), now.Day(), InputHour-2, InputMinute, 0, 0, time.UTC).Format(time.RFC3339),
			TimeZone: "Africa/Cairo",
		},
		End: &calendar.EventDateTime{
			DateTime: time.Date(now.Year(), now.Month(), now.Day(), InputHour-2, InputMinute, 20, 0, time.UTC).Format(time.RFC3339),
			TimeZone: "Africa/Cairo",
		},
	}

	//inserting alarm
	_, err = srv.Events.Insert("primary", alarm).Do()
	if err != nil {
		log.Fatalf("Unable to create alarm. %v \n", err)
		return "Unable to create alarm."
	}

	return ""
}

// GetAlarms gets user next alarms from calendar
func GetAlarms(uuid string) string {
	connectGoogleCalendar()
	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events. %v", err)
	}
	alarms := ""
	alarms += "Upcoming events: \n\n"
	if len(events.Items) > 0 {
		for _, i := range events.Items {
			var when string
			// If the DateTime is an empty string the Event is an all-day Event.
			// So only Date is available.
			if i.Start.DateTime != "" {
				when = i.Start.DateTime
			} else {
				when = i.Start.Date
			}
			if i.Summary == uuid {
				alarms += when + "\n"
			}
		}
	} else {
		alarms = "No upcoming events found.\n"
	}
	return alarms
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	// usr, err := user.Current()
	// if err != nil {
	// 	return "", err
	// }
	tokenCacheDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	//tokenCacheDir := filepath.Join("/")
	fmt.Println("tokenCacheDir: ", tokenCacheDir)
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("calendar-go-quickstart.json")), nil
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func connectGoogleCalendar() {
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/calendar-go-quickstart.json
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

	service, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
	} else {
		fmt.Println("google calendar connected successfully")
	}
	srv = service
}
