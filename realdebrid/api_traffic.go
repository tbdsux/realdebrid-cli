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

func (c *RealDebridClient) GetTraffic() (*Traffic, error) {
	var tr Traffic
	resp, err := c.R().SetSuccessResult(&tr).Get("traffic")
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

func (c *RealDebridClient) GetTrafficDetails(start string, end string) (*TrafficDetails, error) {
	var trd TrafficDetails
	resp, err := c.R().SetSuccessResult(&trd).SetQueryParam("start", start).SetQueryParam("end", end).Get("traffic/details")
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccessState() {
		return nil, resp.Err
	}

	return &trd, nil
}
