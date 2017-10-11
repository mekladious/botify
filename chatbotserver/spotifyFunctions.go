package chatbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	return string(body)
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
