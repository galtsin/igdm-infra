package types

import (
	"channels-instagram-dm/domain/model/instagram"
)

type User struct {
	ID            interface{} `json:"pk"`
	Username      string      `json:"username"`
	FullName      string      `json:"full_name"`
	IsPrivate     bool        `json:"is_private"`
	ProfilePicURL string      `json:"profile_pic_url"`
	IsVerified    bool        `json:"is_verified"`
	// HasAnonymousProfilePicture   bool   `json:"has_anonymous_profile_picture"`
	// IsUsingUnifiedInboxForDirect int `json:"is_using_unified_inbox_for_direct"`
}

func (u User) ToModel() (instagram.User, error) {
	model := instagram.User{
		ID:            ValueToString(u.ID),
		Username:      u.Username,
		FullName:      u.FullName,
		IsPrivate:     u.IsPrivate,
		ProfilePicURL: u.ProfilePicURL,
		IsVerified:    u.IsVerified,
	}

	if err := model.Validate(); err != nil {
		return instagram.User{}, err
	}

	return model, nil
}
