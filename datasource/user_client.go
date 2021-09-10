package datasource

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"user_posts/datasource_models"
)

const userEndpoint = "https://jsonplaceholder.typicode.com/users/%d"

type ServerError struct {
	code int
	msg  string
}

func (e *ServerError) Error() string {
	return fmt.Sprintf("Server returned status code %d:\n%s", e.code, e.msg)
}

type UserClient interface {
	GetUser(ctx context.Context, id int) (*datasource_models.User, error)
}

type UserClientImpl struct {
	*http.Client
}

func NewUserClient() UserClient {
	return &UserClientImpl{http.DefaultClient}
}

func (u *UserClientImpl) GetUser(ctx context.Context, ID int) (*datasource_models.User, error) {
	url := fmt.Sprintf(userEndpoint, ID)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error while creating GET request on UserClient: %s", err.Error())
	}
	request = request.WithContext(ctx)
	response, err := u.Do(request)
	if err != nil {
		return nil, fmt.Errorf("problem while doing request to %s: %s", url, err.Error())
	}

	if response.StatusCode < 200 && response.StatusCode >= 300 {
		body, _ := ioutil.ReadAll(response.Body)
		return nil, &ServerError{ response.StatusCode, string(body) }
	}

	decoder := json.NewDecoder(response.Body)
	user := &datasource_models.User{}
	err = decoder.Decode(user)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall user: %s", err.Error())
	}

	return user, nil
}