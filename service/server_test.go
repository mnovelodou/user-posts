package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"user_posts/business_logic"
	"user_posts/datasource"
	"user_posts/datasource_models"
)

// Test successful call
func TestServer(t *testing.T) {
	expBody := `{"id":1,"userInfo":{"name":"Leanne Graham","username":"Bret","email":"Sincere@april.biz"},"posts":[{"id":1,"title":"sunt aut facere repellat provident occaecati excepturi optio reprehenderit","body":"quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto"},{"id":2,"title":"qui est esse","body":"est rerum tempore vitae\nsequi sint nihil reprehenderit dolor beatae ea dolores neque\nfugiat blanditiis voluptate porro vel nihil molestiae ut reiciendis\nqui aperiam non debitis possimus qui neque nisi nulla"},{"id":3,"title":"ea molestias quasi exercitationem repellat qui ipsa sit aut","body":"et iusto sed quo iure\nvoluptatem occaecati omnis eligendi aut ad\nvoluptatem doloribus vel accusantium quis pariatur\nmolestiae porro eius odio et labore et velit aut"},{"id":4,"title":"eum et est occaecati","body":"ullam et saepe reiciendis voluptatem adipisci\nsit amet autem assumenda provident rerum culpa\nquis hic commodi nesciunt rem tenetur doloremque ipsam iure\nquis sunt voluptatem rerum illo velit"},{"id":5,"title":"nesciunt quas odio","body":"repudiandae veniam quaerat sunt sed\nalias aut fugiat sit autem sed est\nvoluptatem omnis possimus esse voluptatibus quis\nest aut tenetur dolor neque"},{"id":6,"title":"dolorem eum magni eos aperiam quia","body":"ut aspernatur corporis harum nihil quis provident sequi\nmollitia nobis aliquid molestiae\nperspiciatis et ea nemo ab reprehenderit accusantium quas\nvoluptate dolores velit et doloremque molestiae"},{"id":7,"title":"magnam facilis autem","body":"dolore placeat quibusdam ea quo vitae\nmagni quis enim qui quis quo nemo aut saepe\nquidem repellat excepturi ut quia\nsunt ut sequi eos ea sed quas"},{"id":8,"title":"dolorem dolore est ipsam","body":"dignissimos aperiam dolorem qui eum\nfacilis quibusdam animi sint suscipit qui sint possimus cum\nquaerat magni maiores excepturi\nipsam ut commodi dolor voluptatum modi aut vitae"},{"id":9,"title":"nesciunt iure omnis dolorem tempora et accusantium","body":"consectetur animi nesciunt iure dolore\nenim quia ad\nveniam autem ut quam aut nobis\net est aut quod aut provident voluptas autem voluptas"},{"id":10,"title":"optio molestias id quia eum","body":"quo et expedita modi cum officia vel magni\ndoloribus qui repudiandae\nvero nisi sit\nquos veniam quod sed accusamus veritatis error"}]}`
	server := NewServer(business_logic.NewUserPostLogic(datasource.NewUserClient(), datasource.NewPostClient()))
	server.Start(":8181")
	defer server.Shutdown(context.Background())

	response, err := http.Get("http://localhost:8181/v1/user-posts/1")
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		t.Fatalf("Expected successful response code but got %d", response.StatusCode)
	}

	gotBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	gotBody := strings.TrimSpace(string(gotBodyBytes))
	if gotBody != expBody {
		t.Errorf("expected body:\n%s\n\nGot:\n%s", expBody, gotBody)
	}
}

// This have multiple tests with bad requests.
func TestBadCalls(t *testing.T) {
	server := NewServer(business_logic.NewUserPostLogic(datasource.NewUserClient(), datasource.NewPostClient()))
	server.Start(":8181")
	defer server.Shutdown(context.Background())

	tests := []struct{
		name string
		request string
		statusCode int
	}{
		{"noID", "http://localhost:8181/v1/user-posts/", 404},
		{"noNumberID", "http://localhost:8181/v1/user-posts/sdf", 404},
		{"IDNotFoundUpper", "http://localhost:8181/v1/user-posts/11", 404},
		{"IDNotFoundLower", "http://localhost:8181/v1/user-posts/0", 404},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := http.Get(test.request)
			if err != nil {
				t.Error(err)
				return
			}

			if response.StatusCode != test.statusCode {
				t.Errorf("expected status code 400 instead of %d", response.StatusCode)
			}
		})
	}
}

type RandomErrorPost struct{}

func (rep *RandomErrorPost) GetUserPosts(ctx context.Context, _ int) ([]*datasource_models.Post, error) {
	return nil, fmt.Errorf("random error at post service")
}

// Test internal errors
func TestInternalErrors(t *testing.T) {
	server := NewServer(business_logic.NewUserPostLogic(datasource.NewUserClient(), &RandomErrorPost{}))
	server.Start(":8181")
	defer server.Shutdown(context.Background())

	response, err := http.Get("http://localhost:8181/v1/user-posts/1")
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != 500 {
		t.Errorf("expected status code of 500 instead of %d", response.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	body := string(bodyBytes)
	if !strings.Contains(body, "random error at post service") {
		t.Errorf("Expected error to be \"random error at post service\", but is: %s", body)
	}
}