package realdebrid

import "fmt"

type StreamingMediaInfos struct {
	Filename           string            `json:"filename"`
	Hoster             string            `json:"hoster"`
	Link               string            `json:"link"`
	Type               string            `json:"type"`
	Season             string            `json:"season,omitempty"`
	Episode            string            `json:"episode,omitempty"`
	Year               string            `json:"year,omitempty"`
	Duration           float64           `json:"duration"`
	Bitrate            int               `json:"bitrate"`
	Size               int64             `json:"size"`
	Details            MediaInfoDetails  `json:"details"`
	PosterPath         string            `json:"poster_path"`
	AudioImage         string            `json:"audio_image"`
	BackdropPath       string            `json:"backdrop_path"`
	BaseURL            string            `json:"base_url"`
	AvailableFormats   map[string]string `json:"available_formats"`
	AvailableQualities map[string]string `json:"available_qualities"`
	ModelURL           string            `json:"model_url"`
	Host               string            `json:"host"`
}

type MediaInfoDetailsVideo struct {
	Stream     string  `json:"stream"`
	Lang       string  `json:"lang"`                 // Language in plain text (ex "English", "French")
	LangISO    string  `json:"lang_iso"`             // Language in iso_639 (ex fre, eng)
	Codec      string  `json:"codec"`                // Codec of the video (ex "h264", "divx")
	Colorspace string  `json:"colorspace,omitempty"` // Colorspace of the video (ex "yuv420p")
	Width      int     `json:"width,omitempty"`      // Width and height of the video
	Height     int     `json:"height,omitempty"`     // Width and height of the video
	Sampling   int     `json:"sampling,omitempty"`   // Audio sampling rate
	Channels   float64 `json:"channels,omitempty"`   // Number of channels (ex 2, 5.1, 7.1)
	Type       string  `json:"type,omitempty"`       // Format of subtitles
}

type MediaInfoDetails struct {
	Video     map[string]MediaInfoDetailsVideo `json:"video"`
	Audio     map[string]MediaInfoDetailsVideo `json:"audio"`
	Subtitles map[string]MediaInfoDetailsVideo `json:"subtitles"`
}

// GetStreamingMediaInfos retrieves detailed information about a streaming media by its ID.
// ID is from `GetDownloads` or `UnrestricLink` response.
// `GET /streaming/mediaInfos/{id}`
func (c *RealDebridClient) GetStreamingMediaInfos(id string) (*StreamingMediaInfos, error) {
	if id == "" {
		return nil, fmt.Errorf("media ID cannot be empty")
	}

	var mediaInfo StreamingMediaInfos
	resp, err := c.client.R().SetPathParam("id", id).SetSuccessResult(&mediaInfo).Get("streaming/mediaInfos/{id}")
	if err != nil || !resp.IsSuccessState() {
		return nil, fmt.Errorf("failed to get streaming media infos for ID %s: %w", id, err)
	}

	return &mediaInfo, nil
}

type StreamingTranscode struct {
	Apple    TranscodeInfo `json:"apple"`
	Dash     TranscodeInfo `json:"dash"`
	LiveMP4  TranscodeInfo `json:"liveMP4"`
	H264WebM TranscodeInfo `json:"h264WebM"`
}

type TranscodeInfo struct {
	Full string `json:"full"`
}

// GetStreamingTranscode retrieves the transcoding links for given ID.
// ID is from `GetDownloads` or `UnrestricLink` response.
// `GET /streaming/transcode/{id}`
func (c *RealDebridClient) GetStreamingTranscode(id string) (*StreamingTranscode, error) {
	if id == "" {
		return nil, fmt.Errorf("transcode ID cannot be empty")
	}

	var transcode StreamingTranscode
	resp, err := c.client.R().SetPathParam("id", id).SetSuccessResult(&transcode).Get("streaming/transcode/{id}")
	if err != nil || !resp.IsSuccessState() {
		return nil, fmt.Errorf("failed to get streaming transcode for ID %s: %w", id, err)
	}

	return &transcode, nil
}
