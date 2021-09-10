package service_models

type User struct {
	ID       int      `json:"id"`
	UserInfo UserInfo `json:"userInfo"`
	Posts    []*Post  `json:"posts"`
}

type UserInfo struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}