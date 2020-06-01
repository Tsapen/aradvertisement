package auth

// TokenDetails contains token details.
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

// AccessDetails struct contains access details.
type AccessDetails struct {
	AccessUUID string
	Username   string
}

// User struct contains user information.
type User struct {
	Username string
	Password string
}

// DB is a database interface.
type DB interface {
	CreateAuth(string, *TokenDetails) error
	FetchAuth(authD *AccessDetails) (string, error)
	DeleteAuth(string) error
}
