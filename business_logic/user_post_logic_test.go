package business_logic

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"user_posts/datasource"
	"user_posts/datasource_models"
	"user_posts/service_models"
)

type RandomErrorPost struct{}

type RandomErrorUser struct{}

func (rep *RandomErrorPost) GetUserPosts(ctx context.Context, _ int) ([]*datasource_models.Post, error) {
	return nil, fmt.Errorf("random error at post service")
}

func (reu *RandomErrorUser) GetUser(ctx context.Context, _ int) (*datasource_models.User, error) {
	return nil, fmt.Errorf("random error at user service")
}

// Tests a valid requests with real clients
func TestUserPostLogic_GetUserPost(t *testing.T) {
	logic := NewUserPostLogic(datasource.NewUserClient(), datasource.NewPostClient())
	user, err := logic.GetUserPost(context.Background(), 1)
	if err != nil {
		t.Fatalf("Get user post should not have failed: %s", err.Error())
	}

	user.Posts = nil
	expectedUser := service_models.User{
		ID: 1,
		UserInfo: service_models.UserInfo{
			Name: "Leanne Graham",
			Username: "Bret",
			Email: "Sincere@april.biz",
		},
		Posts: []*service_models.Post{
			{
				ID: 1,
				Title: "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
				Body: "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto",
			},
			{
				ID: 2,
				Title: "qui est esse",
				Body: "est rerum tempore vitae\nsequi sint nihil reprehenderit dolor beatae ea dolores neque\nfugiat blanditiis voluptate porro vel nihil molestiae ut reiciendis\nqui aperiam non debitis possimus qui neque nisi nulla",
			},
			{
				ID: 3,
				Title: "ea molestias quasi exercitationem repellat qui ipsa sit aut",
				Body: "et iusto sed quo iure\nvoluptatem occaecati omnis eligendi aut ad\nvoluptatem doloribus vel accusantium quis pariatur\nmolestiae porro eius odio et labore et velit aut",
			},
			{
				ID: 4,
				Title: "eum et est occaecati",
				Body: "ullam et saepe reiciendis voluptatem adipisci\nsit amet autem assumenda provident rerum culpa\nquis hic commodi nesciunt rem tenetur doloremque ipsam iure\nquis sunt voluptatem rerum illo velit",
			},
			{
				ID: 5,
				Title: "nesciunt quas odio",
				Body: "repudiandae veniam quaerat sunt sed\nalias aut fugiat sit autem sed est\nvoluptatem omnis possimus esse voluptatibus quis\nest aut tenetur dolor neque",
			},
			{
				ID: 6,
				Title: "dolorem eum magni eos aperiam quia",
				Body: "ut aspernatur corporis harum nihil quis provident sequi\nmollitia nobis aliquid molestiae\nperspiciatis et ea nemo ab reprehenderit accusantium quas\nvoluptate dolores velit et doloremque molestiae",
			},
			{
				ID: 7,
				Title: "magnam facilis autem",
				Body: "dolore placeat quibusdam ea quo vitae\nmagni quis enim qui quis quo nemo aut saepe\nquidem repellat excepturi ut quia\nsunt ut sequi eos ea sed quas",
			},
			{
				ID: 8,
				Title: "dolorem dolore est ipsam",
				Body: "dignissimos aperiam dolorem qui eum\nfacilis quibusdam animi sint suscipit qui sint possimus cum\nquaerat magni maiores excepturi\nipsam ut commodi dolor voluptatum modi aut vitae",
			},
			{
				ID: 9,
				Title: "nesciunt iure omnis dolorem tempora et accusantium",
				Body: "consectetur animi nesciunt iure dolore\nenim quia ad\nveniam autem ut quam aut nobis\net est aut quod aut provident voluptas autem voluptas",
			},
			{
				ID: 10,
				Title: "optio molestias id quia eum",
				Body: "quo et expedita modi cum officia vel magni\ndoloribus qui repudiandae\nvero nisi sit\nquos veniam quod sed accusamus veritatis error",
			},
		},
	}

	if user.ID != expectedUser.ID {
		t.Errorf("expected user.ID %d, got %d", expectedUser.ID, user.ID)
	}

	if user.UserInfo.Name != expectedUser.UserInfo.Name {
		t.Errorf("expected user.UserInfo.Name %s, got %s", expectedUser.UserInfo.Name, user.UserInfo.Name)
	}

	if user.UserInfo.Username != expectedUser.UserInfo.Username {
		t.Errorf("expected user.UserInfo.Username %s, got %s", expectedUser.UserInfo.Username, user.UserInfo.Username)
	}

	if user.UserInfo.Email != expectedUser.UserInfo.Email {
		t.Errorf("expected user.UserInfo.Email %s, got %s", expectedUser.UserInfo.Email, user.UserInfo.Email)
	}

	for i, post := range user.Posts {
		if *post != *expectedUser.Posts[i] {
			t.Errorf("expected post to be %v but got %v", *expectedUser.Posts[i], *post)
		}
	}
}

// Tests an error returning from UserClient
func TestRandomErrorAtUserService(t *testing.T) {
	logic := NewUserPostLogic(&RandomErrorUser{}, datasource.NewPostClient())
	_, err := logic.GetUserPost(context.Background(), 1)
	if err == nil {
		t.Fatalf("expected to have an error")
	}

	if !strings.Contains(err.Error(), "random error at user service") {
		t.Errorf("error expected to have 'random error at user service' but got: %s", err.Error())
	}
}

// Tests an error returning from PostClient
func TestRandomErrorAtPostService(t *testing.T) {
	logic := NewUserPostLogic(datasource.NewUserClient(), &RandomErrorPost{})
	_, err := logic.GetUserPost(context.Background(), 1)
	if err == nil {
		t.Fatalf("expected to have an error")
	}

	if !strings.Contains(err.Error(), "random error at post service") {
		t.Errorf("error expected to have 'random error at post service' but got: %s", err.Error())
	}
}

// Tests an error from both clients
func TestRandomErrorBothServices(t *testing.T) {
	logic := NewUserPostLogic(&RandomErrorUser{}, &RandomErrorPost{})
	_, err := logic.GetUserPost(context.Background(), 1)
	if err == nil {
		t.Fatalf("expected to have an error")
	}

	if !strings.Contains(err.Error(), "random error at post service") {
		t.Errorf("error expected to have 'random error at post service' but got: %s", err.Error())
	}

	if !strings.Contains(err.Error(), "random error at user service") {
		t.Errorf("error expected to have 'random error at user service' but got: %s", err.Error())
	}
}