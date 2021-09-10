package datasource

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"user_posts/datasource_models"
)

const postEndpoint = "https://jsonplaceholder.typicode.com/posts?userId=%d"

// PostClient interface creates an abstraction for
type PostClient interface {
	GetUserPosts(ctx context.Context, userID int) ([]*datasource_models.Post, error)
}

type PostClientImpl struct {
	*http.Client
}

func NewPostClient() PostClient {
	return &PostClientImpl{http.DefaultClient}
}

func (p *PostClientImpl) GetUserPosts(ctx context.Context, userID int) ([]*datasource_models.Post, error) {
	url := fmt.Sprintf(postEndpoint, userID)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error while creating post request: %s", err.Error())
	}

	request = request.WithContext(ctx)
	response, err := p.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error while doing post request: %s", err.Error())
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, _ := ioutil.ReadAll(response.Body)
		return nil, &ServerError{response.StatusCode, string(body)}
	}

	decoder := json.NewDecoder(response.Body)
	var posts []*datasource_models.Post
	err = decoder.Decode(&posts)
	if err != nil {
		return nil, fmt.Errorf("problem while decoding json with posts: %s", err.Error())
	}

	return posts, nil
}
