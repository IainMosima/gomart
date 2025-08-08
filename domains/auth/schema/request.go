package schema

type LoginRequest struct {
	RedirectURI string `json:"redirect_uri" validate:"required,url"`
	State       string `json:"state,omitempty"`
}

type TokenExchangeRequest struct {
	Code        string `json:"code" validate:"required"`
	State       string `json:"state,omitempty"`
	RedirectURI string `json:"redirect_uri" validate:"required,url"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type CompleteSetupRequest struct {
	Phone      *string `json:"phone,omitempty" validate:"omitempty,e164"`
	Address    *string `json:"address,omitempty"`
	City       *string `json:"city,omitempty"`
	PostalCode *string `json:"postal_code,omitempty"`
}

type LogoutRequest struct {
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
