package models

import (
	"time"
)

// User struct represents the structure of a user document in MongoDB
type User struct {
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	TimeCreate  time.Time `json:"timecreate"`
	TimeLogin   time.Time `json:"timelogin"`
	Permissions string    `json:"permissions"`
	TokenKey    string    `json:"tokenkey"`
}
