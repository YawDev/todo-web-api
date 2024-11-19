package messages

var AccountNotFound string = "existing account not found"
var AlreadyLoggedIn string = "User is already logged in"
var UserNotFound string = "user not found"
var SomethingWentWrong string = "Internal error. Something went wrong"
var NoTokenProvided string = "No token provided."
var InvalidToken string = "invalid token."
var InvalidRefreshToken string = "invalid refresh token."
var InvalidPassword string = "Invalid Password Credentials"
var AccessTokenError string = "Error while generating access token."
var UnauthorizedRefreshToken = "refresh token unauthorized"
var NotFound string = "Not found"

var SuccessLogout = "User logged out successfully"
var SuccessLogin = "User logged in successfully"
var SuccessUserCreate = "User created successfully"
var SuccessListCreate = "List created successfully"
var SuccessTaskCreate = "Task created successfully"

var UserNotFoundInDb = "User record not found in db"
var TaskNotFoundInDb = "Task record not found in db"
var ListNotFoundInDb = "List record not found in db"

var FailedTaskDelete = "Task delete failed"
var FailedListDelete = "List delete failed"
var FailedUserDelete = "User delete failed"

var TaskQueryInternalError string = "something went wrong while fetching task"
var ListQueryInternalError string = "something went wrong while fetching list"
var UserQueryInternalError string = "something went wrong while fetching user"
