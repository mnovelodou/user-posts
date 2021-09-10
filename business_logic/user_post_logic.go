package business_logic

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"user_posts/datasource"
	"user_posts/datasource_models"
	"user_posts/service_models"
)

// UserPostLogic contains the business logic behind getting the user and its posts using the clients.
type UserPostLogic struct {
	userClient datasource.UserClient
	postClient datasource.PostClient
}

type IDNotFound struct{}

func (i *IDNotFound) Error() string {
	return ""
}

func NewUserPostLogic(userClient datasource.UserClient, postClient datasource.PostClient) *UserPostLogic {
	return &UserPostLogic{userClient, postClient}
}

// GetUserPost gets user-posts given a userID
func (upl *UserPostLogic) GetUserPost(ctx context.Context, userID int) (*service_models.User, error) {
	if userID < 1 || userID > 10 {
		// There are only userID from 1 to 10
		return nil, &IDNotFound{}
	}

	var userDS *datasource_models.User
	var postsDS []*datasource_models.Post
	wg := sync.WaitGroup{}
	// Adding 2 because we will run 2 go routines one for each client request
	wg.Add(2)
	// errCh is a buffered channel because if both clients gives errors, a non-buffered channel would
	// block when waitgroup waits as non-buffered channel cannot write both errors without reading the first
	// error in the channel
	errCh := make(chan error, 2)

	// Get user from userClient
	go func() {
		defer wg.Done()
		var err error = nil
		userDS, err = upl.userClient.GetUser(ctx, userID)
		if err != nil {
			errCh <- err
		}
	}()

	// Get posts from postClient
	go func() {
		defer wg.Done()
		var err error = nil
		postsDS, err = upl.postClient.GetUserPosts(ctx, userID)
		if err != nil {
			errCh <- err
		}
	}()

	wg.Wait()
	close(errCh)
	if len(errCh) > 0 {
		return nil, fmt.Errorf("have errors while consuming web services:\n%s", collectErrors(errCh))
	}

	// At this point, both calls are successful, so we only map datasource_models.User and Post into service_models
	user := UserDatasourceToService(userDS)
	user.Posts = make([]*service_models.Post, 0, len(postsDS))
	for _, postDS := range postsDS {
		user.Posts = append(user.Posts, PostDataSourceToService(postDS))
	}

	return user, nil
}

// collectErrors builds an string with all error messages in the error channel
func collectErrors(errCh chan error) string {
	var sb strings.Builder
	for err := range errCh {
		sb.WriteString(err.Error())
		sb.WriteString("\n")
	}
	return sb.String()
}