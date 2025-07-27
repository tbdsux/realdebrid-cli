package realdebrid

type Hosts = map[string]struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

// GetHosts retrieves supported hosts.
// This function does not require authentication.
// `GET /hosts`
func (c *RealDebridClient) GetHosts() (*Hosts, error) {
	var output Hosts

	resp, err := c.client.R().SetSuccessResult(&output).Get("hosts")
	if err != nil || !resp.IsSuccessState() {
		return nil, resp.Err
	}

	return &output, nil
}

type HostsStatus = map[string]struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Image             string `json:"image"`
	Supported         int    `json:"supported"`
	Status            string `json:"status"`
	CheckTime         string `json:"check_time"`
	CompetitorsStatus map[string]struct {
		Status    string `json:"status"`
		CheckTime string `json:"check_time"`
	} `json:"competitors_status"`
}

// GetHostsStatus gets the status of supported hosters or not and their status on competitors.
// `GET /hosts/status`
func (c *RealDebridClient) GetHostsStatus() (*HostsStatus, error) {
	var output HostsStatus

	resp, err := c.client.R().SetSuccessResult(&output).Get("hosts/status")
	if err != nil || !resp.IsSuccessState() {
		return nil, resp.Err
	}

	return &output, nil
}

// GetHostsRegex gets all supported folder Regex, useful to find supported links inside a document.
// This function does not require authentication.
// `GET /hosts/regex`
func (c *RealDebridClient) GetHostsRegex() ([]string, error) {
	var output []string

	resp, err := c.client.R().SetSuccessResult(&output).Get("hosts/regex")
	if err != nil || !resp.IsSuccessState() {
		return nil, resp.Err
	}

	return output, nil
}

// GetHostsRegexFolder gets all supported folder Regex, useful to find supported links inside a document.
// This function does not require authentication.
// `GET /hosts/regexFolder`
func (c *RealDebridClient) GetHostsRegexFolder() ([]string, error) {
	var output []string

	resp, err := c.client.R().SetSuccessResult(&output).Get("hosts/regexFolder")
	if err != nil || !resp.IsSuccessState() {
		return nil, resp.Err
	}

	return output, nil
}

// GetHostsDomains gets all hoster domains supported on the service.
// This function does not require authentication.
// `GET /hosts/domains`
func (c *RealDebridClient) GetHostsDomains() ([]string, error) {
	var output []string

	resp, err := c.client.R().SetSuccessResult(&output).Get("hosts/domains")
	if err != nil || !resp.IsSuccessState() {
		return nil, resp.Err
	}

	return output, nil
}
