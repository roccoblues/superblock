package twitter

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Tweet data returned from the API.
type Tweet struct {
	ID        string    `json:"id"`
	AuthorID  string    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	Text      string    `json:"text"`
}

type tweetResponse struct {
	Data *Tweet `json:"data"`
	Meta *meta  `json:"meta,omitempty"`
}

type likingUsersResponse struct {
	Data []*User `json:"data"`
	Meta *meta   `json:"meta,omitempty"`
}

// LikingUsersInterator is used to iterate over liking users results.
type LikingUsersInterator struct {
	client *Client
	meta   *meta
	path   string
	params *url.Values
}

// DefaultTweetFields is the default list of tweet fields requested.
const DefaultTweetFields = "id,author_id,created_at,text"

// Tweet returns information about a single Tweet.
// API: https://developer.twitter.com/en/docs/twitter-api/tweets/lookup/api-reference/get-tweets-id
// Rate limit: 900 requests per 15-minute window per each authenticated user.
func (c *Client) Tweet(tweetID string, params *url.Values) (*Tweet, error) {
	if params == nil {
		params = &url.Values{}
		params.Set("tweet.fields", DefaultTweetFields)
	}
	req, err := c.NewRequest(http.MethodGet, fmt.Sprintf("/tweets/%s?%s", tweetID, params.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result tweetResponse
	err = c.DoRequest(req, &result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

// LikingUsers returns an iterator to get information about a Tweet’s liking users.
// API: https://developer.twitter.com/en/docs/twitter-api/tweets/likes/api-reference/get-tweets-id-liking_users
// Rate limit: 75 requests per 15-minute window per each authenticated user.
func (c *Client) LikingUsers(id string, params *url.Values) *LikingUsersInterator {
	path := fmt.Sprintf("/tweets/%s/liking_users", id)
	if params == nil {
		params = &url.Values{}
		params.Set("user.fields", "id,name,created_at")
	}
	return &LikingUsersInterator{
		client: c,
		path:   path,
		params: params,
	}
}

// NextPage advances the iterator and returns true if there are more results.
func (i *LikingUsersInterator) NextPage() bool {
	if i.meta == nil {
		return true
	}
	if i.meta.NextToken == "" {
		return false
	}
	i.params.Set("pagination_token", i.meta.NextToken)
	return true
}

// Users returns the liking users result for the current pagination token.
func (i *LikingUsersInterator) User() ([]*User, error) {
	req, err := i.client.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", i.path, i.params.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result likingUsersResponse
	err = i.client.DoRequest(req, &result)
	if err != nil {
		return nil, err
	}

	i.meta = result.Meta
	return result.Data, nil
}
