package chatbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Jeffail/gabs"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	Favorite struct {
		ID			bson.ObjectId `bson:"_id,omitempty"`
		Uuid		string `bson:"uuid"`
		Trackid		string `bson:"trackid"`
	}
)

// AuthorizeSpotify is a function to authorizing with spotify
func AuthorizeSpotify() string {
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

func Get_featured_playlists() string {
	body, _ := sendGetRequest("v1/browse/featured-playlists", "")
	// bodyJSON := JSON{}
	// err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&bodyJSON)
	jsonParsed, _ := gabs.ParseJSON(body)

	names := jsonParsed.Path("playlists.items.name")
	hrefs := jsonParsed.Path("playlists.items.external_urls.spotify")

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
	x := len(s2)

	//x := len(s4)
	sFinal := ""
	for i := 0; i < x-1; i++ {
		sFinal = sFinal + s2[i] + " : " + s4[i] + " \n"
	}

	fmt.Println(sFinal)
	//fmt.Println(hrefs)
	return sFinal
}

func Get_artist_tracks(singerName string) string {
	singerName = strings.Replace(singerName, " ", "%20", -1) // replacing spaces by %20 as required by spotify api
	artist_id := Get_artist_id(singerName)

	body, _ := sendGetRequest("v1/artists/"+artist_id+"/top-tracks?country=US", "")
	jsonParsed, _ := gabs.ParseJSON(body)

	tracks := jsonParsed.Path("tracks.preview_url")

	return tracks.String()
}

func Get_artist_info(singerName string) string {
	singerName = strings.Replace(singerName, " ", "%20", -1) // replacing spaces by %20 as required by spotify api
	artist_id := Get_artist_id(singerName)

	body, _ := sendGetRequest("v1/artists/"+artist_id, "")
	jsonParsed, _ := gabs.ParseJSON(body)
	//name := ""
	//followers := 0
	//test := jsonParsed.Path(name)

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
	fmt.Println(sFinal)

	return sFinal
	//return string(body)
}

func Get_artist_id(singerName string) string {
	//replacing spaces with %20
	strings.Replace("", singerName, singerName, -1)
	body, _ := sendGetRequest("v1/search?type=artist&q="+singerName+"&limit=1", "")

	jsonParsed, _ := gabs.ParseJSON(body)

	ids := jsonParsed.Path("artists.items.id")

	return strings.TrimLeft(strings.TrimRight(ids.String(), "\"]"), "[\"")
}

func get_new_releases() string {
	body, _ := sendGetRequest("v1/browse/new-releases", "")
	jsonParsed, _ := gabs.ParseJSON(body)

	names := jsonParsed.Path("albums.items.name")
	hrefs := jsonParsed.Path("albums.items.external_urls.spotify")

	s1 := ""
	children, _ := names.Children()
	for _, child := range children {
		s1 = s1 + child.Data().(string)
		s1 = s1 + "$"

	}

	s2 := strings.Split(s1, "$")

	s3 := ""
	children1, _ := hrefs.Children()
	for _, child := range children1 {
		s3 = s3 + child.Data().(string)
		s3 = s3 + "$"

	}

	s4 := strings.Split(s3, "$")
	x := len(s4)

	//x := len(s4)
	sFinal := ""
	for i := 0; i < x-1; i++ {
		sFinal = sFinal + s2[i] + " : " + s4[i] + " \n"
	}

	fmt.Println(sFinal)
	//fmt.Println(hrefs)
	return sFinal
}

func Get_mood(mood string) string {
	body, _ := sendGetRequest("v1/browse/categories/"+mood+"/playlists", "")
	jsonParsed, _ := gabs.ParseJSON(body)

	names := jsonParsed.Path("playlists.items.name")
	hrefs := jsonParsed.Path("playlists.items.external_urls.spotify")

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
	x := len(s2)

	//x := len(s4)
	sFinal := ""
	for i := 0; i < x-1; i++ {
		sFinal = sFinal + s2[i] + " : " + s4[i] + " \n"
	}

	fmt.Println(sFinal)
	//fmt.Println(hrefs)
	return sFinal
	//playlist := Get_playlist_tracks()
}

func min(x int, y int) int {
	if x > y {
		return y
	} else {
		return x
	}
}

func search(keyword string) string {
	artist, _ := sendGetRequest("v1/search?q="+keyword+"&type=artist", "")
	playlist, _ := sendGetRequest("v1/search?q="+keyword+"&type=playlist", "")
	album, _ := sendGetRequest("v1/search?q="+keyword+"&type=album", "")
	track, _ := sendGetRequest("v1/search?q="+keyword+"&type=track", "")

	result := append(artist, playlist...)
	result = append(result, album...)
	result = append(result, track...)

	return string(result)
}

func add_to_favorites(uuid string, trackid string) (string, error){
	db, err := mgo.Dial(db_uri)
	collection := db.DB("botify").C("Favorites")
	err = collection.Insert(&Favorite{Uuid:uuid, Trackid:trackid})
	if err!= nil{
		return "", err
	} else{
		return "success", nil
	}
}

func get_favorites(uuid string) (JSON, error){
	db, err := mgo.Dial(db_uri)
	collection := db.DB("botify").C("Favorites")
	
	var results []Favorite
	collection.Find(bson.M{"uuid": uuid}).All(&results)
	// collection.Find(nil).All(&results)

	res := JSON{"Favorites":results}
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
