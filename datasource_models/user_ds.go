package datasource_models

type User struct {
	ID       int         `json:"id"`
	Name     string      `json:"name"`
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Address  UserAddress `json:"address"`
	Phone    string      `json:"phone"`
	Website  string      `json:"website"`
	Company  UserCompany `json:"company"`
}

type UserAddress struct {
	Street  string         `json:"street"`
	Suite   string         `json:"suite"`
	City    string         `json:"city"`
	Zipcode string         `json:"zipcode"`
	Geo     UserAddressGeo `json:"geo"`
}

type UserCompany struct {
	Name        string `json:"name"`
	CatchPhrase string `json:"catchphrase"`
	Bs          string `json:"bs"`
}

type UserAddressGeo struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}