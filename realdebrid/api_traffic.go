package realdebrid

type Traffic = map[string]struct {
	Left  int64  `json:"left"`
	Bytes int64  `json:"bytes"`
	Links int64  `json:"links"`
	Limit int64  `json:"limit"`
	Type  string `json:"type"`
	Extra int64  `json:"extra"`
	Reset string `json:"reset"`
}

// GetTraffic retrieves the traffic limitations for limited hosters (limits, current usage, extra packages).
// `GET /traffic`
func (c *RealDebridClient) GetTraffic() (*Traffic, error) {
	var tr Traffic
	resp, err := c.client.R().SetSuccessResult(&tr).Get("traffic")
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccessState() {
		return nil, resp.Err
	}

	return &tr, nil
}

type TrafficDetails = map[string]struct {
	Host  map[string]int64 `json:"host"`
	Bytes int64            `json:"bytes"`
}

// GetTrafficDetails retrieves detailed traffic information on each hoster during a specified period.
// The `start` and `end` parameters should be in ISO 8601 format (e.g., "2023-01-01T00:00:00Z").
// `GET /traffic/details`
func (c *RealDebridClient) GetTrafficDetails(start string, end string) (*TrafficDetails, error) {
	var trd TrafficDetails
	resp, err := c.client.R().SetSuccessResult(&trd).SetQueryParam("start", start).SetQueryParam("end", end).Get("traffic/details")
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccessState() {
		return nil, resp.Err
	}

	return &trd, nil
}
