package realdebrid

import "fmt"

type Download struct {
	ID        string `json:"id"`
	Filename  string `json:"filename"`
	MimeType  string `json:"mimeType"`
	FileSize  int64  `json:"fileSize"`
	Link      string `json:"link"`
	Host      string `json:"host"`
	Chunks    int64  `json:"chunks"`
	Download  string `json:"download"`
	Generated string `json:"generated"`
}

type GetDownloadRequest struct {
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`
	Page   int `json:"page,omitempty"`
}

// GetDownloads gets the user's downloads list.
// Defaults: offset=0, limit=10, page=1
// `GET /downloads`
func (c *RealDebridClient) GetDownloads(req *GetDownloadRequest) ([]Download, error) {
	params := map[string]string{
		"offset": "0",
		"limit":  "10",
		"page":   "1",
	}

	if req != nil && req.Offset != 0 {
		params["offset"] = fmt.Sprint(req.Offset)
	}
	if req != nil && req.Limit != 0 {
		params["limit"] = fmt.Sprint(req.Limit)
	}
	if req != nil && req.Page != 0 {
		params["page"] = fmt.Sprint(req.Page)
	}

	var downloads []Download

	resp, err := c.client.R().SetSuccessResult(&downloads).SetQueryParams(params).Get("downloads")
	if err != nil || !resp.IsSuccessState() {
		return nil, err
	}

	return downloads, nil
}

// DeleteDownload deletes a a link from the user's downloads list.
// `DELETE /downloads/delete/{id}`
func (c *RealDebridClient) DeleteDownload(id string) error {
	if id == "" {
		return fmt.Errorf("download ID cannot be empty")
	}

	resp, err := c.client.R().SetPathParam("id", id).Delete("downloads/delete/{id}")
	if err != nil || !resp.IsSuccessState() {
		return fmt.Errorf("failed to delete download with ID %s: %w", id, err)
	}

	return nil
}
