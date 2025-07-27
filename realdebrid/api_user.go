package realdebrid

type User struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Points     int64  `json:"points"`
	Locale     string `json:"locale"`
	Avatar     string `json:"avatar"`
	Type       string `json:"type"`
	Premium    int64  `json:"premium"`
	Expiration string `json:"expiration"`
}

func (c *RealDebridClient) GetUser() (*User, error) {
	var user User

	resp, err := c.client.R().SetSuccessResult(&user).Get("user")
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccessState() {
		return nil, resp.Err
	}

	return &user, nil
}
