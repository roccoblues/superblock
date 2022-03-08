package twitter

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// User data returned from the API.
type User struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	UserName        string    `json:"username"`
	CreatedAt       time.Time `json:"created_at"`
	Description     string    `json:"description"`
	Location        string    `json:"location"`
	PinnedTweetID   string    `json:"pinned_tweet_id"`
	ProfileImageUrl string    `json:"profile_image_url"`
	Protected       bool      `json:"protected"`
	URL             string    `json:"url"`
	Verified        bool      `json:"verified"`
	Withheld        bool      `json:"withheld"`
}

type meResponse struct {
	Data *User `json:"data"`
	Meta *meta `json:"meta,omitempty"`
}

type blockingRequest struct {
	TargetUserID string `json:"target_user_id"`
}

type blockingResponse struct {
	Data struct {
		Blocking bool `json:"blocking"`
	}
}

// DefaultUserFields is the default list of user fields requested.
const DefaultUserFields = "id,name,username,created_at,description,location,pinned_tweet_id,profile_image_url,protected,url,verified,withheld"

// Me returns information about an authorized user.
// API: https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users-me
// Rate limit: 75 requests per 15-minute window per each authenticated user.
func (c *Client) Me(params *url.Values) (*User, error) {
	if params == nil {
		params = &url.Values{}
		params.Set("user.fields", DefaultUserFields)
	}
	req, err := c.NewRequest(http.MethodGet, fmt.Sprintf("/users/me?%s", params.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result meResponse
	err = c.DoRequest(req, &result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

// User returns information about a single user.
// API: https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users-id
// Rate limit: 900 requests per 15-minute window per each authenticated user.
func (c *Client) User(id string, params *url.Values) (*User, error) {
	if params == nil {
		params = &url.Values{}
		params.Set("user.fields", DefaultUserFields)
	}
	req, err := c.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s?%s", id, params.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result meResponse
	err = c.DoRequest(req, &result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

// Block causes the user to block the target user.
// API: https://developer.twitter.com/en/docs/twitter-api/users/blocks/api-reference/post-users-user_id-blocking
// Rate limit: 50 requests per 15-minute window per each authenticated user.
func (c *Client) Block(user string, targetUser string) (bool, error) {
	body := blockingRequest{TargetUserID: targetUser}
	req, err := c.NewRequest(http.MethodPost, fmt.Sprintf("/users/%s/blocking", user), body)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	var result blockingResponse
	err = c.DoRequest(req, &result)
	if err != nil {
		return false, err
	}

	return result.Data.Blocking, nil
}
