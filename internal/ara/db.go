package ara

// ObjectCreationInfo contains object information.
type ObjectCreationInfo struct {
	Username  string
	Latitude  float64
	Longitude float64
	Comment   string
	Type      string
}

// ObjectSelectInfo is struct for doing select.
type ObjectSelectInfo struct {
	Latitude  float64
	Longitude float64
}

// ObjectSelectByID is struct for doing select.
type ObjectSelectByID struct {
	Username string
	Type     rune
}

// ObjectAroundResp contains object information.
type ObjectAroundResp struct {
	ID int
	ObjectCreationInfo
}

// ObjectUpdateInfo contains inforamtion for object updating.
type ObjectUpdateInfo struct {
	ID      int
	Comment string
}

// UserCreationInfo contains user information.
type UserCreationInfo struct {
	Username string
	Email    string
	Password string
}

// UserObjectSelectResp contains useful object information for user.
type UserObjectSelectResp struct {
	ID        int
	Comment   string
	Latitude  float64
	Longitude float64
}

// UserLogin struct contains info for login.
type UserLogin struct {
	Username string
	Password string
}

// DB is a database interface.
type DB interface {
	CreateObject(ObjectCreationInfo) (int, error)
	SelectObjectsAround(ObjectSelectInfo) ([]ObjectAroundResp, error)
	SelectObjectByID(int) (ObjectSelectByID, error)
	SelectUsersObjects(string) ([]UserObjectSelectResp, error)
	UpdateObject(ObjectUpdateInfo) error
	DeleteObject(int) error

	CreateUser(UserCreationInfo) error
	DeleteUser(string) error

	CheckLogin(UserLogin) (bool, error)
}
