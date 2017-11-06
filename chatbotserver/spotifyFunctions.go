package chatbot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	Favorite struct {
		ID        bson.ObjectId `bson:"_id,omitempty"`
		Uuid      string        `bson:"uuid"`
		Trackid   string        `bson:"trackid"`
		TrackName string        `bson:"trackName"`
	}
)

// AuthorizeSpotify is a function to authorizing with spotify
func AuthorizeSpotify() string {

	cacheFile, err := spotifyTokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached spotify token file. %v", err)
	}
	tok, err := spotifyTokenFromFile(cacheFile)
	if err != nil {
		tok = getSpotifyTokenFromWeb()
		saveSpotifyToken(cacheFile, tok)
	}
	return tok
}

func getNewSpotifyToken() string {

	cacheFile, err := spotifyTokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached spotify token file. %v", err)
	}
	tok := getSpotifyTokenFromWeb()
	saveSpotifyToken(cacheFile, tok)

	return tok
}

func getSpotifyTokenFromWeb() string {
	//create a headers map
	url := "https://accounts.spotify.com/api/token"
	var jsonStr = []byte(`grant_type=client_credentials`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic ZDc0NzRjMTg0OGU0NDFmM2FiOTAyMGQyNzM2OTE2ZGE6Y2M4NTU3YTE0ZmFkNDNiNTliMDI4MDc5YmE3ZTM2Yjc=")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	//converting body to json
	bodyJSON := JSON{}
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&bodyJSON)
	spotifyAuthorizationToken, _ := bodyJSON["access_token"].(string)
	fmt.Println("spotify access token:", spotifyAuthorizationToken)
	return spotifyAuthorizationToken
}

func saveSpotifyToken(file string, token string) {
	fmt.Printf("Saving spotify token file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func spotifyTokenFromFile(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	t := ""
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

func spotifyTokenCacheFile() (string, error) {
	tokenCacheDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("spotifyToken.json")), nil
}

func Get_featured_playlists() string {
	body, _ := sendGetRequest("v1/browse/featured-playlists", "")

	jsonParsed, _ := gabs.ParseJSON(body)

	featured_playlists_string := stringfyJSON(jsonParsed, "playlist")

	return featured_playlists_string
}

func Get_artist_tracks(singerName string) (string, error) {
	singerName = strings.Replace(singerName, " ", "%20", -1) // replacing spaces by %20 as required by spotify api
	artist_id, err := Get_artist_id(singerName)

	if err != nil {
		return "", err
	}

	body, _ := sendGetRequest("v1/artists/"+artist_id+"/top-tracks?country=US", "")
	jsonParsed, _ := gabs.ParseJSON(body)

	tracks := jsonParsed.Path("tracks.preview_url")

	return tracks.String(), nil
}

func Get_artist_id(singerName string) (string, error) {
	//replacing spaces with %20
	strings.Replace("", singerName, singerName, -1)
	body, _ := sendGetRequest("v1/search?type=artist&q="+singerName+"&limit=1", "")

	jsonParsed, _ := gabs.ParseJSON(body)

	ids := jsonParsed.Path("artists.items.id")
	if strings.TrimLeft(strings.TrimRight(ids.String(), "\"]"), "[\"") == "{}" {
		return "", errors.New("invalid artist")
	}

	return strings.TrimLeft(strings.TrimRight(ids.String(), "\"]"), "[\""), nil
}

func Get_artist_info(singerName string) string {
	singerName = strings.Replace(singerName, " ", "%20", -1) // replacing spaces by %20 as required by spotify api
	artist_id, err := Get_artist_id(singerName)

	if err != nil {
		return err.Error()
	}

	body, _ := sendGetRequest("v1/artists/"+artist_id, "")
	jsonParsed, _ := gabs.ParseJSON(body)

	artistInfo := stringfyArtist(jsonParsed)

	return artistInfo
}

func get_new_releases() string {
	body, _ := sendGetRequest("v1/browse/new-releases", "")
	jsonParsed, _ := gabs.ParseJSON(body)

	newReleases := stringfyJSON(jsonParsed, "album")

	return newReleases
}

func Get_mood(mood string) string {
	body, _ := sendGetRequest("v1/browse/categories/"+mood+"/playlists", "")
	jsonParsed, _ := gabs.ParseJSON(body)

	playlists := stringfyJSON(jsonParsed, "playlist")

	return playlists
}

func min(x int, y int) int {
	if x > y {
		return y
	} else {
		return x
	}
}

func search(keyword string) string {
	// artist, _ := sendGetRequest("v1/search?q="+keyword+"&type=artist", "")
	// artistJsonParsed, _ := gabs.ParseJSON(artist)
	playlist, _ := sendGetRequest("v1/search?q="+keyword+"&type=playlist", "")
	playlistJsonParsed, _ := gabs.ParseJSON(playlist)
	album, _ := sendGetRequest("v1/search?q="+keyword+"&type=album", "")
	albumJsonParsed, _ := gabs.ParseJSON(album)
	track, _ := sendGetRequest("v1/search?q="+keyword+"&type=track", "")
	trackJsonParsed, _ := gabs.ParseJSON(track)

	result := "Artists: \n\n"
	// result += stringfyArtist(artistJsonParsed)
	result += Get_artist_info(keyword)
	result += "\n\n Playlists: \n\n"
	result += stringfyJSON(playlistJsonParsed, "playlist")
	result += "\n\n Albums: \n\n"
	result += stringfyJSON(albumJsonParsed, "album")
	result += "\n\n Tracks: \n\n"
	result += stringfyJSON(trackJsonParsed, "track")

	return string(result)
}

func getTrackID(trackName string) string {
	trackName = strings.Replace(trackName, " ", "%20", -1)
	body, _ := sendGetRequest("v1/search?q="+trackName+"&type=track&limit=1", "")
	jsonParsed, _ := gabs.ParseJSON(body)
	ids := jsonParsed.Path("tracks.items.id")
	if strings.TrimLeft(strings.TrimRight(ids.String(), "\"]"), "[\"") == "{}" {
		return "nil"
	}
	return strings.TrimLeft(strings.TrimRight(ids.String(), "\"]"), "[\"")
}

func add_to_favorites(uuid string, trackid string, trackName string) (string, error) {
	db, err := mgo.Dial(db_uri)
	collection := db.DB("botify").C("Favorites")
	err = collection.Insert(&Favorite{Uuid: uuid, Trackid: trackid, TrackName: trackName})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return "This track already exists in your favorites", nil
		}
		return "", err
	} else {
		return trackName + " successfully added to your favourites", nil
	}
}

func get_favorites(uuid string) (string, error) {
	db, err := mgo.Dial(db_uri)
	collection := db.DB("botify").C("Favorites")

	var results []Favorite
	collection.Find(bson.M{"uuid": uuid}).All(&results)
	if len(results) == 0 {
		return "You don't have any favorite tracks yet, use favorite add: to add new favorites", nil
	}
	// collection.Find(nil).All(&results)
	res := ""
	for i := 0; i < len(results); i++ {
		r := results[i]
		index := i + 1
		res = res + strconv.Itoa(index) + ") " + r.TrackName + ": https://open.spotify.com/track/" + r.Trackid + " \n"
	}
	// res := JSON{"Favorites": results}
	return res, err
}

func sendGetRequest(url string, body string) ([]byte, string) {
	defer func() {
		fmt.Println("Recovered from get request error: ", recover())
	}()
	if SpotifyAuthorizationToken == "" {
		log.Print("No spotify authorization token.. started to get one")
		SpotifyAuthorizationToken = AuthorizeSpotify()
	}

	fmt.Println("Sending get request to", "https://api.spotify.com/"+url)

	// var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("GET", "https://api.spotify.com/"+url, bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf8")
	req.Header.Set("Authorization", "Bearer "+SpotifyAuthorizationToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	//	fmt.Println("response Headers:", resp.Header)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(respBody))

	return respBody, string(resp.Status)
}

func stringfyArtist(jsonParsed *gabs.Container) string {
	name := jsonParsed.Path("name").Data().(string)
	followers := jsonParsed.Path("followers.total").Data().(float64)

	z := int(followers)
	f := fmt.Sprint(z)
	genres := jsonParsed.Path("genres")
	s1 := ""
	children, _ := genres.Children()
	for _, child := range children {
		s1 = s1 + child.Data().(string)
		s1 = s1 + ","

	}
	s2 := strings.Split(s1, ",")
	y := len(s2)
	s5 := ""
	for i := 0; i < y-1; i++ {
		s5 = s5 + s2[i] + " , "
	}

	images := jsonParsed.Path("images.url")
	s3 := ""
	children1, _ := images.Children()
	for _, child := range children1 {
		s3 = s3 + child.Data().(string)
		s3 = s3 + ","

	}
	s4 := strings.Split(s3, ",")
	x := len(s4)
	s6 := ""
	for i := 0; i < x-1; i++ {
		s6 = s6 + s4[i] + " , "
	}
	sFinal := ""
	sFinal += name + "\n" +
		"has " + f + " followers\n" +
		"specializes in " + s5 + ".\n" +
		"Pictures: " + s6
	return sFinal
}

func stringfyJSON(jsonParsed *gabs.Container, jsontype string) string {
	var names *gabs.Container
	var hrefs *gabs.Container
	if jsontype == "playlist" {
		names = jsonParsed.Path("playlists.items.name")
		hrefs = jsonParsed.Path("playlists.items.external_urls.spotify")
	} else if jsontype == "album" {
		names = jsonParsed.Path("albums.items.name")
		hrefs = jsonParsed.Path("albums.items.external_urls.spotify")
	} else if jsontype == "track" {
		names = jsonParsed.Path("tracks.items.name")
		hrefs = jsonParsed.Path("tracks.items.external_urls.spotify")
	}

	s1 := ""
	children, _ := names.Children()
	for _, child := range children {
		s1 = s1 + child.Data().(string)
		s1 = s1 + ","

	}

	s2 := strings.Split(s1, ",")

	s3 := ""
	children1, _ := hrefs.Children()
	for _, child := range children1 {
		s3 = s3 + child.Data().(string)
		s3 = s3 + "$"

	}

	s4 := strings.Split(s3, "$")
	x := min(len(s2), len(s4))

	sFinal := ""
	for i := 0; i < x-1; i++ {
		sFinal = sFinal + s2[i] + " : " + s4[i] + " \n"
	}
	return sFinal
}
