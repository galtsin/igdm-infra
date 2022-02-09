package instagram

import "fmt"

type User struct {
	ID            string
	Username      string
	FullName      string
	IsPrivate     bool
	ProfilePicURL string
	IsVerified    bool
}

func (u User) Validate() error {
	if u.ID == "" {
		return fmt.Errorf("ID should not be empty")
	}

	if u.Username == "" {
		return fmt.Errorf("Username should not be empty")
	}

	return nil
}
