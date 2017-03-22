package api

// ServerStatus checks the current status of the server
const ServerStatus = "/api/status"

// FacebookLogin calls the FB api and initiates login request
const FacebookLogin = "/api/facebooklogin"

// FacebookCallback is the callback url that login verification / user details
const FacebookCallback = "/api/facebookcallback"

// FacebookLoginSuccess displays facebook login success page
const FacebookLoginSuccess = "/api/facebooksuccess"

// FacebookLoginFailed displays facebook login failure page
const FacebookLoginFailed = "/api/facebookfailure"

// GetFacebookData is an api url to get user facebook data from cache
// Requires user login verification, device identifier, and SHA256 key id
const GetFacebookData = "/api/facebookdata"

// DeleteFacebookData is an api url to remove user from token cache
// Requires user login verification, device identifier, and SHA256 key id
const DeleteFacebookData = "/api/facebookdelete"
