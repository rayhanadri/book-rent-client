package variables

import (
	"library-client/model"
)

// var IsLoggedIn = false
var CurrentUser model.User

var ApiBaseUrl string = "https://stark-citadel-27445-77ccd9936ec4.herokuapp.com/api/v1/"
var AccessToken string = ""
var RefreshToken string = ""
var IsLoggedIn bool = false
