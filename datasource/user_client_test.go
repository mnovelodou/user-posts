package datasource

import (
	"context"
	"encoding/json"
	"testing"
	"user_posts/datasource_models"
)

func TestUserServiceImpl_GetUser(t *testing.T) {
	userService := NewUserClient()
	ctx := context.Background()
	user, err := userService.GetUser(ctx, 10)
	if err != nil {
		t.Fatalf("unexpected error on GetUser: %s", err.Error())
	}
	expectedJSON := `{
  "id": 10,
  "name": "Clementina DuBuque",
  "username": "Moriah.Stanton",
  "email": "Rey.Padberg@karina.biz",
  "address": {
    "street": "Kattie Turnpike",
    "suite": "Suite 198",
    "city": "Lebsackbury",
    "zipcode": "31428-2261",
    "geo": {
      "lat": "-38.2386",
      "lng": "57.2232"
    }
  },
  "phone": "024-648-3804",
  "website": "ambrose.net",
  "company": {
    "name": "Hoeger LLC",
    "catchPhrase": "Centralized empowering task-force",
    "bs": "target end-to-end models"
  }
}`
	expectedUser := datasource_models.User{}
	err = json.Unmarshal([]byte(expectedJSON), &expectedUser)
	if err != nil {
		t.Fatalf("Unexpected error at unmarshalling expected user: %s", err.Error())
	}

	if *user != expectedUser {
		t.Errorf("Expected user:\n %v\n\n but got %v", expectedUser, *user)
	}
}
