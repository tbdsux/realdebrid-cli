package realdebrid

import (
	"bytes"
	"errors"
	"os"
	"strings"
)

type Torrent struct {
	ID       string   `json:"id"`
	Filename string   `json:"filename"`
	Hash     string   `json:"hash"`
	Bytes    int64    `json:"bytes"`
	Host     string   `json:"host"`
	Split    int64    `json:"split"`
	Progress int64    `json:"progress"`
	Status   string   `json:"status"`
	Added    string   `json:"added"`
	Links    []string `json:"links"`
	Ended    string   `json:"ended,omitempty"`
	Speed    int64    `json:"speed,omitempty"`
	Seeders  int64    `json:"seeders,omitempty"`
}

type TorrentInfo struct {
	*Torrent
	OriginalFilename string `json:"original_filename"`
	OriginalBytes    int64  `json:"original_bytes"`
	Files            []struct {
		ID       int64  `json:"id"`
		Path     string `json:"path"`
		Bytes    int64  `json:"bytes"`
		Selected int    `json:"selected"`
	}
}

// GetTorrents retrieves the list of torrents of the user.
// `GET /torrents`
func (c *RealDebridClient) GetTorrents() ([]Torrent, error) {
	var output []Torrent

	resp, err := c.client.R().SetSuccessResult(&output).Get("torrents")
	if err != nil || !resp.IsSuccessState() {
		return nil, resp.Err
	}

	return output, nil
}

// GetTorrentsInfo gets all informations on the specified torrent id.
// `GET /torrents/info/{id}`
func (c *RealDebridClient) GetTorrentsInfo(id string) (*TorrentInfo, error) {
	var output TorrentInfo

	resp, err := c.client.R().SetSuccessResult(&output).Get("torrents/info/" + id)
	if err != nil || !resp.IsSuccessState() {
		return nil, resp.Err
	}

	return &output, nil
}

type TorrentAvailableHost struct {
	Host        string `json:"host"`
	MaxFileSize int64  `json:"max_file_size"`
}

// GetTorrentsAvailableHosts retrieves the list of available hosts to upload the torrent to.
// `GET /torrents/availableHosts`
func (c *RealDebridClient) GetTorrentsAvailableHosts() ([]TorrentAvailableHost, error) {
	var output []TorrentAvailableHost

	resp, err := c.client.R().SetSuccessResult(&output).Get("torrents/availableHosts")
	if err != nil || !resp.IsSuccessState() {
		return nil, resp.Err
	}

	return output, nil
}

type AddTorrent struct {
	ID  string `json:"id"`
	URI string `json:"uri"`
}

// AddTorrent uploads a torrent file to RealDebrid.
// The file should be a valid torrent file and exists at the specified filepath.
// `PUT /torrents/addTorrent`
func (c *RealDebridClient) AddTorrent(filepath string) (*AddTorrent, error) {
	var output AddTorrent

	if _, err := os.Stat(filepath); err != nil {
		return nil, errors.New("file does not exist: " + filepath)
	}

	// Read the file content
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, errors.New("failed to read file: " + err.Error())
	}

	resp, err := c.client.R().
		SetSuccessResult(&output).
		SetHeader("Content-Type", "application/octet-stream").
		SetBody(bytes.NewBuffer(data)).
		Put("torrents/addTorrent")
	if err != nil || !resp.IsSuccessState() {
		return nil, resp.Err
	}

	// TODO: add upload callback?

	return &output, nil
}

// AddMagnet adds a magnet link to RealDebrid.
// The magnet link should be a valid magnet URI.
// `POST /torrents/addMagnet`
func (c *RealDebridClient) AddMagnet(magnet string) (*AddTorrent, error) {
	var output AddTorrent

	if magnet == "" {
		return nil, errors.New("magnet link cannot be empty")
	}

	resp, err := c.client.R().
		SetSuccessResult(&output).
		SetFormData(map[string]string{
			"magnet": magnet,
		}).
		Post("torrents/addMagnet")
	if err != nil || !resp.IsSuccessState() {
		return nil, resp.Err
	}

	return &output, nil
}

// SelectTorrentFiles selects files in a torrent to start it.
// If fileIds is empty, all files will be selected.
// `POST /torrents/selectFiles/{id}`
func (c *RealDebridClient) SelectTorrentFiles(id string, fileIds []string) error {
	files := ""

	if len(fileIds) == 0 {
		files = "all"
	} else {
		files = strings.Join(fileIds, ",")
	}

	resp, err := c.client.R().
		SetFormData(map[string]string{
			"files": files,
		}).Post(
		"torrents/selectFiles/" + id,
	)
	if err != nil || !resp.IsSuccessState() {
		return resp.Err
	}

	return nil
}

// DeleteTorrent deletes a torrent from your torrent list.
// `DELETE /torrents/delete/{id}`
func (c *RealDebridClient) DeleteTorrent(id string) error {
	resp, err := c.client.R().Delete(
		"torrents/delete/" + id,
	)
	if err != nil || !resp.IsSuccessState() {
		return resp.Err
	}

	return nil
}
