package realdebrid

import (
	"bytes"
	"errors"
	"fmt"
	"os"
)

type Settings struct {
	DownloadPorts                []string          `json:"download_ports"`
	DownloadPort                 string            `json:"download_port"`
	DownloadProtocols            []string          `json:"download_protocols"`
	DownloadProtocol             string            `json:"download_protocol"`
	Locales                      map[string]string `json:"locales"`
	Locale                       string            `json:"locale"`
	StreamingQualities           []string          `json:"streaming_qualities"`
	StreamingQuality             string            `json:"streaming_quality"`
	MobileStreamingQuality       string            `json:"mobile_streaming_quality"`
	StreamingLanguages           map[string]string `json:"streaming_languages"`
	StreamingLanguagePreference  string            `json:"streaming_language_preference"`
	StreamingCastAudio           []string          `json:"streaming_cast_audio"`
	StreamingCastAudioPreference string            `json:"streaming_cast_audio_preference"`
}

// GetSettings retrieves the user's settings.
// `GET /settings`
func (c *RealDebridClient) GetSettings() (*Settings, error) {
	var settings Settings

	resp, err := c.client.R().SetSuccessResult(&settings).Get("settings")
	if err != nil || !resp.IsSuccessState() {
		return nil, resp.Err
	}

	return &settings, nil
}

type SettingName string

var SettingDownloadPort SettingName = "download_port"
var SettingLocale SettingName = "locale"
var SettingStreamingLanguagePreference SettingName = "streaming_language_preference"
var SettingStreamingQuality SettingName = "streaming_quality"
var SettingMobileStreamingQuality SettingName = "mobile_streaming_quality"
var SettingStreamingCastAudioPreference SettingName = "streaming_cast_audio_preference"

// UpdateSettings updates a specific setting.
// `POST /settings/update`
func (c *RealDebridClient) UpdateSettings(name SettingName, value string) error {
	if name == "" || value == "" {
		return fmt.Errorf("setting name and value cannot be empty")
	}

	var formData = map[string]string{
		"setting_name":  string(name),
		"setting_value": value,
	}

	resp, err := c.client.R().SetFormData(formData).Post("settings/update")
	if err != nil || !resp.IsSuccessState() {
		return fmt.Errorf("failed to update setting %s: %w", name, err)
	}

	return nil
}

// ConvertPoints converts the user's fidelity points.
// `POST /settings/convertPoints`
func (c *RealDebridClient) ConvertPoints() error {
	resp, err := c.client.R().Post("settings/convertPoints")
	if err != nil || !resp.IsSuccessState() {
		return fmt.Errorf("failed to convert points: %w", err)
	}
	return nil
}

// ChangePassword sends a verification email to change the user's password.
// `POST /settings/changePassword`
func (c *RealDebridClient) ChangePassword() error {
	resp, err := c.client.R().Post("settings/changePassword")
	if err != nil || !resp.IsSuccessState() {
		return fmt.Errorf("failed to change password: %w", err)
	}
	return nil
}

// PutAvatarFile uploads a new avatar image for the user.
// The file should be a valid image file.
// Max image size: 1MB, supported formats: JPG, PNG, GIF, max dimensions: 150x150.
// `PUT /settings/avatarFile`
func (c *RealDebridClient) PutAvatarFile(filepath string) error {
	if filepath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	if _, err := os.Stat(filepath); err != nil {
		return fmt.Errorf("file does not exist: %s", filepath)
	}

	// Read the file content
	data, err := os.ReadFile(filepath)
	if err != nil {
		return errors.New("failed to read file: " + err.Error())
	}

	mime := ""
	if ext := filepath[len(filepath)-4:]; ext == ".jpg" || ext == ".jpeg" {
		mime = "image/jpeg"
	} else if ext == ".png" {
		mime = "image/png"
	} else if ext == ".gif" {
		mime = "image/gif"
	} else {
		return fmt.Errorf("unsupported file type: %s", ext)
	}

	resp, err := c.client.R().
		SetHeader("Content-Type", mime).
		SetBody(bytes.NewBuffer(data)).
		Put("settings/avatarFile")
	if err != nil || !resp.IsSuccessState() {
		return fmt.Errorf("failed to upload avatar: %w", err)
	}

	return nil
}

// DeleteAvatar removes the user's avatar image to set it back to default.
// `DELETE /settings/avatarDelete`
func (c *RealDebridClient) DeleteAvatar() error {
	resp, err := c.client.R().Delete("settings/avatarDelete")
	if err != nil || !resp.IsSuccessState() {
		return fmt.Errorf("failed to delete avatar: %w", err)
	}
	return nil
}
