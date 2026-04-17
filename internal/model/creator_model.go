package model

type UpdateProfileRequest struct {
	DisplayName     *string `json:"display_name"      validate:"omitempty,min=2,max=100"`
	Username        *string `json:"username"          validate:"omitempty,min=3,max=30,alphanum"`
	Bio             *string `json:"bio"               validate:"omitempty,max=500"`
	ThemeColor      *string `json:"theme_color"       validate:"omitempty,len=7"`
	MinDonation     *int64  `json:"min_donation"      validate:"omitempty,min=1000"`
	MaxDonation     *int64  `json:"max_donation"      validate:"omitempty"`
	ThankYouMessage *string `json:"thank_you_message" validate:"omitempty,max=300"`
}

type PublicProfileResponse struct {
	Username        string  `json:"username"`
	DisplayName     string  `json:"display_name"`
	AvatarURL       *string `json:"avatar_url"`
	Bio             *string `json:"bio"`
	ThemeColor      string  `json:"theme_color"`
	MinDonation     int64   `json:"min_donation"`
	MaxDonation     int64   `json:"max_donation"`
	ThankYouMessage string  `json:"thank_you_message"`
}

type AlertSettingsRequest struct {
	Duration  int    `json:"duration"  validate:"min=1,max=30"`
	Template  string `json:"template"  validate:"required"`
	SoundURL  string `json:"sound_url"`
	Animation string `json:"animation"`
}
