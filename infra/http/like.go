package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type LikeData struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type LikingUsersResponse struct {
	Data []LikeData `json:"data"`
}

func LikingUsers(twitterID string, at string) []LikeData {
	url := fmt.Sprintf("https://api.twitter.com/2/tweets/%s/liking_users", twitterID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+at)
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	var decoded LikingUsersResponse
	if err := json.Unmarshal(byteArray, &decoded); err != nil {
		log.Fatal(err)
	}
	return decoded.Data
}
