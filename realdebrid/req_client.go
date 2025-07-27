package realdebrid

import (
	"fmt"

	"github.com/imroc/req/v3"
)

type ApiError struct {
	ErrorMessage string `json:"error"`
	ErrorCode    int64  `json:"error_code"`
}

func (msg *ApiError) Error() string {
	return fmt.Sprintf("API Error: %s", msg.ErrorMessage)
}

type RealDebridClient struct {
	client *req.Client
}

func NewClient(apiKey string) *RealDebridClient {
	return &RealDebridClient{
		client: req.C().
			SetCommonBearerAuthToken(apiKey).
			SetBaseURL("https://api.real-debrid.com/rest/1.0/").
			SetCommonErrorResult(&ApiError{}).
			EnableDumpEachRequest().
			OnAfterResponse(func(client *req.Client, resp *req.Response) error {
				if resp.Err != nil {
					return nil
				}

				if apiErr, ok := resp.ErrorResult().(*ApiError); ok {
					resp.Err = apiErr
					return nil
				}

				if !resp.IsSuccessState() {
					return fmt.Errorf("bad response, raw dyno:\n%s", resp.Dump())
				}

				return nil
			}),
	}
}
