package vk

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

const clientID = "3140623"                        //VK for iPhone app client_id
const clientSecret = "VeWdmVclDCtn6ihuP1nt"       //VK for iPhone app client_secret
const authURL = "https://oauth.vk.com/token?"     //Direct Authorization URL
const apiMethodURL = "https://api.vk.com/method/" //Method request URL

// AuthResponse structure contains all parameters of response of authorization request
type AuthResponse struct {
	UserID           int    `json:"user_id"`
	ExpiresIn        int    `json:"expires_in"`
	AccessToken      string `json:"access_token"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// WallGetResponse - wall.get response
type WallGetResponse struct {
	Response struct {
		Count int `json:"count"`
		Items []struct {
			ID          int    `json:"id"`
			FromID      int    `json:"from_id"`
			OwnerID     int    `json:"owner_id"`
			Date        int    `json:"date"`
			Type        string `json:"type"`
			MarkedAsAds int    `json:"marked_as_ads"`
			PostType    string `json:"post_type"`
			Text        string `json:"text"`
			SignerID    int    `json:"signer_id"`
			IsPinned    int    `json:"is_pinned,omitempty"`
			Attachments []struct {
				Type  string `json:"type"`
				Photo struct {
					ID      int `json:"id"`
					AlbumID int `json:"album_id"`
					OwnerID int `json:"owner_id"`
					UserID  int `json:"user_id"`
					Sizes   []struct {
						Type   string `json:"type"`
						URL    string `json:"url"`
						Width  int    `json:"width"`
						Height int    `json:"height"`
					} `json:"sizes"`
					Text      string `json:"text"`
					Date      int    `json:"date"`
					PostID    int    `json:"post_id"`
					AccessKey string `json:"access_key"`
				} `json:"photo"`
			} `json:"attachments"`
			PostSource struct {
				Type string `json:"type"`
			} `json:"post_source"`
			Comments struct {
				Count         int  `json:"count"`
				CanPost       int  `json:"can_post"`
				GroupsCanPost bool `json:"groups_can_post"`
			} `json:"comments"`
			Likes struct {
				Count      int `json:"count"`
				UserLikes  int `json:"user_likes"`
				CanLike    int `json:"can_like"`
				CanPublish int `json:"can_publish"`
			} `json:"likes"`
			Reposts struct {
				Count        int `json:"count"`
				UserReposted int `json:"user_reposted"`
			} `json:"reposts"`
			Views struct {
				Count int `json:"count"`
			} `json:"views"`
			IsFavorite     bool `json:"is_favorite"`
			AdsEasyPromote struct {
				Type int `json:"type"`
			} `json:"ads_easy_promote"`
		} `json:"items"`
		NextFrom string `json:"next_from"`
	} `json:"response"`
}

// GetCommentsResponse - wall.getComments response
type GetCommentsResponse struct {
	Response struct {
		Count int `json:"count"`
		Items []struct {
			ID           int           `json:"id"`
			FromID       int           `json:"from_id"`
			PostID       int           `json:"post_id"`
			OwnerID      int           `json:"owner_id"`
			ParentsStack []interface{} `json:"parents_stack"`
			Date         int           `json:"date"`
			Text         string        `json:"text"`
			Likes        struct {
				Count      int  `json:"count"`
				UserLikes  int  `json:"user_likes"`
				CanLike    int  `json:"can_like"`
				CanPublish bool `json:"can_publish"`
			} `json:"likes"`
			Thread struct {
				Count           int           `json:"count"`
				Items           []interface{} `json:"items"`
				CanPost         bool          `json:"can_post"`
				ShowReplyButton bool          `json:"show_reply_button"`
			} `json:"thread"`
		} `json:"items"`
		Profiles []struct {
			ID              int    `json:"id"`
			FirstName       string `json:"first_name"`
			LastName        string `json:"last_name"`
			IsClosed        bool   `json:"is_closed"`
			CanAccessClosed bool   `json:"can_access_closed"`
			Sex             int    `json:"sex"`
			ScreenName      string `json:"screen_name"`
			Photo50         string `json:"photo_50"`
			Photo100        string `json:"photo_100"`
			Online          int    `json:"online"`
		} `json:"profiles"`
		Groups            []interface{} `json:"groups"`
		CurrentLevelCount int           `json:"current_level_count"`
		CanPost           bool          `json:"can_post"`
		ShowReplyButton   bool          `json:"show_reply_button"`
	} `json:"response"`
}

// Auth function makes authorization request and returns *AuthResponse structure
func Auth(login string, password string) (*AuthResponse, error) {
	var jsonResponse *AuthResponse
	var requestURL = url.Values{
		"grant_type":    {"password"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"username":      {login},
		"password":      {password},
	}
	response, err := http.Get(authURL + requestURL.Encode())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(content, &jsonResponse); err != nil {
		return nil, err
	}
	return jsonResponse, nil
}

// VkRequest function makes api method request and returns []byte JSON response
func Request(methodName string, parameters map[string]string) ([]byte, error) {
	requestURL, err := url.Parse(apiMethodURL + methodName)
	if err != nil {
		return nil, err
	}
	requestQuery := requestURL.Query()
	for key, value := range parameters {
		requestQuery.Set(key, value)
	}
	if authResponse != nil {
		requestQuery.Set("access_token", authResponse.AccessToken)
	}
	requestQuery.Set("v", "5.101")
	requestURL.RawQuery = requestQuery.Encode()

	response, err := http.Get(requestURL.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

//RandomShuffle - returns randomized arr of ints
func RandomShuffle(vals []int) []int {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]int, len(vals))
	perm := rnd.Perm(len(vals))
	for i, randIndex := range perm {
		ret[i] = vals[randIndex]
	}
	return ret
}

// InitArrayOfIndexes - generate sequence
func InitArrayOfIndexes(size int) []int {
	var res []int
	for i := 0; i < size; i++ {
		res = append(res, i)
	}

	res = RandomShuffle(res)

	return res
}

var authResponse *AuthResponse

//Init - log in to vk server and obtain access token
func Init(vkLogin string, vkPwd string) error {
	var err error
	authResponse, err = Auth(vkLogin, vkPwd)
	if err != nil {
		return err
	}
	return nil
}
