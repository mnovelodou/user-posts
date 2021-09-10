package business_logic

import (
	"user_posts/datasource_models"
	"user_posts/service_models"
)

func UserDatasourceToService(userDS *datasource_models.User) *service_models.User {
	user := &service_models.User{}
	user.ID = userDS.ID
	user.UserInfo.Name = userDS.Name
	user.UserInfo.Username = userDS.Username
	user.UserInfo.Email = userDS.Email
	return user
}

func PostDataSourceToService(postDS *datasource_models.Post) *service_models.Post {
	post := &service_models.Post{}
	post.ID = postDS.ID
	post.Title = postDS.Title
	post.Body = postDS.Body
	return post
}