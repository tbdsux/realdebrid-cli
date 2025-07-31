package realdebrid

import (
	"bytes"
	"errors"
	"fmt"
	"os"
)

type UnrestrictProps struct {
	Link     string
	Password string // Password to unlock the file access hoster side
	Remote   int    // 0 or 1, use Remote traffic, dedicated servers and account sharing protections lifted

}

type UnrestrictedLink struct {
	ID          string                        `json:"id"`
	Filename    string                        `json:"filename"`
	MimeType    string                        `json:"mimeType"`
	FileSize    int64                         `json:"filesize"`
	Link        string                        `json:"link"`
	Host        string                        `json:"host"`
	Chunks      int64                         `json:"chunks"`
	Download    string                        `json:"download"`
	Streamable  int                           `json:"streamable"`
	Type        string                        `json:"type,omitempty"`
	Alternative []UnrestrictedLinkAlternative `json:"alternative,omitempty"` // Alternative links if available
}

type UnrestrictedLinkAlternative struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Download string `json:"download"`
	Type     string `json:"type"`
}

// UnrestricLink retrieves an unrestricted link from a given URL.
// The link must be a valid URL that RealDebrid can process.
// `POST /unrestrict/link`
func (c *RealDebridClient) UnrestricLink(props *UnrestrictProps) (*UnrestrictedLink, error) {
	if props.Link == "" {
		return nil, fmt.Errorf("link cannot be empty")
	}

	var formData = map[string]string{
		"link": props.Link,
	}
	if props.Password != "" {
		formData["password"] = props.Password
	}
	if props.Remote != 0 {
		formData["remote"] = fmt.Sprint(props.Remote)
	}

	var link UnrestrictedLink
	resp, err := c.client.R().
		SetFormData(formData).
		SetSuccessResult(&link).
		Post("unrestrict/link")
	if err != nil || !resp.IsSuccessState() {
		return nil, fmt.Errorf("failed to unrestrict link: %w", err)
	}

	return &link, nil
}

type UnrestrictCheckProps struct {
	Link     string
	Password string
}

type UnrestricCheck struct {
	Host        string `json:"host"`
	HostIcon    string `json:"host_icon"`
	HostIconBig string `json:"host_icon_big"`
	Link        string `json:"link"`
	Filename    string `json:"filename"`
	Filesize    int64  `json:"filesize"`
	Supported   int    `json:"supported"`
}

// UnrestrictCheck checks if a link can be unrestricted.
// Does not require authentication.
// `POST /unrestrict/check`
func (c *RealDebridClient) UnrestrictCheck(props *UnrestrictCheckProps) (*UnrestricCheck, error) {
	if props.Link == "" {
		return nil, fmt.Errorf("link cannot be empty")
	}

	var formData = map[string]string{
		"link": props.Link,
	}
	if props.Password != "" {
		formData["password"] = props.Password
	}

	var check UnrestricCheck
	resp, err := c.client.R().
		SetFormData(formData).
		SetSuccessResult(&check).
		Post("unrestrict/check")
	if err != nil || !resp.IsSuccessState() {
		return nil, fmt.Errorf("failed to check unrestrict link: %w", err)
	}

	return &check, nil
}

// UnrestrictFolder retrieves unrestricted links from a folder URL.
// `POST /unrestrict/folder`
func (c *RealDebridClient) UnrestrictFolder(link string) ([]string, error) {
	if link == "" {
		return nil, fmt.Errorf("link cannot be empty")
	}

	var formData = map[string]string{
		"link": link,
	}

	var result []string
	resp, err := c.client.R().
		SetFormData(formData).
		SetSuccessResult(&result).
		Post("unrestrict/folder")
	if err != nil || !resp.IsSuccessState() {
		return nil, fmt.Errorf("failed to unrestrict folder: %w", err)
	}

	return result, nil
}

// UnrestricContainerFile retrieves unrestricted links from a container file.
// Filetypes: RSDF, CCF, CCF3, DLC
// `POST /unrestrict/containerFile`
func (c *RealDebridClient) UnrestricContainerFile(filepath string) ([]string, error) {
	var result []string

	if _, err := os.Stat(filepath); err != nil {
		return nil, errors.New("file does not exist: " + filepath)
	}

	// Read the file content
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, errors.New("failed to read file: " + err.Error())
	}

	resp, err := c.client.R().
		SetSuccessResult(&result).
		SetHeader("Content-Type", "application/octet-stream").
		SetBody(bytes.NewBuffer(data)).
		Post("unrestrict/containerFile")
	if err != nil || !resp.IsSuccessState() {
		return nil, fmt.Errorf("failed to unrestrict container file: %w", err)
	}

	return result, nil
}

// UnrestrictContainerLink retrieves unrestricted links from a container link.
// The link should point to a container file (RSDF, CCF, CCF3, DLC).
// `POST /unrestrict/containerLink`
func (c *RealDebridClient) UnrestrictContainerLink(link string) ([]string, error) {
	if link == "" {
		return nil, fmt.Errorf("link cannot be empty")
	}

	var formData = map[string]string{
		"link": link,
	}

	var result []string
	resp, err := c.client.R().
		SetFormData(formData).
		SetSuccessResult(&result).
		Post("unrestrict/containerLink")
	if err != nil || !resp.IsSuccessState() {
		return nil, fmt.Errorf("failed to unrestrict container link: %w", err)
	}

	return result, nil
}
