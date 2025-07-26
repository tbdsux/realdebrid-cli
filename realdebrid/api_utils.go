package realdebrid

func (c *RealDebridClient) GetTime() (string, error) {
	resp, err := c.R().Get("time")
	if err != nil {
		return "", err
	}

	return resp.String(), nil
}

func (c *RealDebridClient) GetTimeISO() (string, error) {
	resp, err := c.R().Get("time/iso")
	if err != nil {
		return "", err
	}

	return resp.String(), nil
}
