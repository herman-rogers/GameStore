package client

// ServerStatus checks the current status of the server
const ServerStatus = "/api/server/status"

// FacebookLogin calls the FB api and initiates login request
const FacebookLogin = "/api/facebook/login"

// FacebookCallback is the callback url that login verification / user details
const FacebookCallback = "/api/facebook/callback"

// FacebookLoginSuccess displays facebook login success page
const FacebookLoginSuccess = "/api/facebook/success"

// FacebookLoginFailed displays facebook login failure page
const FacebookLoginFailed = "/api/facebook/failure"

// GetFacebookData is an api url to get user facebook data from cache
// Requires user login verification, device identifier, and SHA256 key id
const GetFacebookData = "/api/facebook/data"

// DeleteFacebookData is an api url to remove user from token cache
// Requires user login verification, device identifier, and SHA256 key id
const DeleteFacebookData = "/api/facebook/delete"
