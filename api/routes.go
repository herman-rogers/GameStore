
package api

// ServerStatus checks the current status of the server
const ServerStatus = "/api/status"

// FacebookLogin calls the FB api and initiates login request
const FacebookLogin = "/api/facebooklogin"

// FacebookCallback is the callback url that login verification / user details
const FacebookCallback = "/api/facebookcallback"

// FacebookLoginPage displays facebook login status to user
const FacebookLoginPage = "/api/facebookpage"

// GetFacebookData is an api url to get user facebook data from cache
// Requires user login verification, device identifier, and SHA256 key id
const GetFacebookData = "/api/facebookdata"

// DeleteFacebookData is an api url to remove user from token cache
// Requires user login verification, device identifier, and SHA256 key id
const DeleteFacebookData = "/api/facebookdelete"