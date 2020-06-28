package http

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

//TokenResponse is a struct to store the Access Token.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

//GetToken will get the access token needed to make the Get requests from Yahoo.
func GetToken() string {
	postURL := "https://api.login.yahoo.com/oauth2/get_token"

	refreshToken := "AGJDl16PDg6zecKSzY64quVsxOwwblSv.9VT241B0mn6n5DBvQ--"
	clientID := "dj0yJmk9VUxqWmo1ZU9JMm9LJmQ9WVdrOVFYUk1aMjAyTXpRbWNHbzlNQS0tJnM9Y29uc3VtZXJzZWNyZXQmc3Y9MCZ4PWZk"
	clientSecret := "70b9060e7c7df616f536e95286a6744f0374c3f0"
	redirectURI := "https://wolfdynastyfootball.com"

	resp, err := http.PostForm(postURL, url.Values{
		"refresh_token": {refreshToken},
		"redirect_uri":  {redirectURI},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"grant_type":    {"refresh_token"},
	})
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var results TokenResponse

	json.Unmarshal([]byte([]byte(body)), &results)
	if err != nil {
		log.Fatal(err)
	}

	return results.AccessToken
}
