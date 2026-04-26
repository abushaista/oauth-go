package command

import "context"

// Command is a marker interface for all commands
type Command interface{}

// CommandHandler is the interface all command handlers must implement
type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) (interface{}, error)
}

// AuthorizeCommand initiates the authorization code flow
type AuthorizeCommand struct {
	ClientID            string
	RedirectURI         string
	ResponseType        string
	Scope               string
	State               string
	CodeChallenge       string
	CodeChallengeMethod string
	UserID              string
}

// TokenCommand exchanges an authorization code for tokens
type TokenCommand struct {
	GrantType    string
	Code         string
	CodeVerifier string
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scope        string
}

// RefreshCommand refreshes an access token using a refresh token
type RefreshCommand struct {
	GrantType    string
	RefreshToken string
	ClientID     string
	ClientSecret string
	Scope        string
}

// LoginCommand authenticates a user with credentials
type LoginCommand struct {
	Username string
	Password string
}

// ConsentCommand records user consent for a client
type ConsentCommand struct {
	UserID   string
	ClientID string
	Scopes   []string
}
