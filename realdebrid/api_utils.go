package realdebrid

// GetTime retrieves the current server time.
// `GET /time`
func (c *RealDebridClient) GetTime() (string, error) {
	resp, err := c.client.R().Get("time")
	if err != nil {
		return "", err
	}

	return resp.String(), nil
}

// GetTimeISO retrieves the current server time in ISO format.
// `GET /time/iso`
func (c *RealDebridClient) GetTimeISO() (string, error) {
	resp, err := c.client.R().Get("time/iso")
	if err != nil {
		return "", err
	}

	return resp.String(), nil
}

// DisableAccessToken disables the current access token.
// `GET /disable_access_token`
func (c *RealDebridClient) DisableAccessToken() error {
	resp, err := c.client.R().Get("disable_access_token")
	if err != nil {
		return err
	}

	if !resp.IsSuccessState() {
		return resp.Err
	}

	return nil
}
