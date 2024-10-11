package models

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type UserResponse struct {
	Message string              `json:"message,omitempty"`
	User    UserWithoutPassword `json:"user,omitempty"`
}

type TokenResponse struct {
	Message string `json:"message,omitempty"`
	Token string `json:"token"`
}

type CollectionResponse struct {
	Message    string `json:"message,omitempty"`
	Collection Collection
}

type WordResponse struct {
	Message      string `json:"message,omitempty"`
	Word         Word
}
