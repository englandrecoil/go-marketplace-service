package dto

type CredentialsRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
