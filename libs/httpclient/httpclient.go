package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type BaseResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type ResponseWithData struct {
	BaseResponse
	Data any `json:"data"`
}

func Send[Req any, Rep any](ctx context.Context, method, uri string, timeout time.Duration, req *Req, rep *Rep, wrapResponse bool) error {
	var reqStream io.Reader

	if req != nil {
		rawData, err := json.Marshal(req)
		if err != nil {
			return fmt.Errorf("fail to encode to json (%s %s): %w", method, uri, err)
		}
		reqStream = bytes.NewBuffer(rawData)
	}

	ctx, fnCancel := context.WithTimeout(ctx, timeout)
	defer fnCancel()

	httpRequest, err := http.NewRequestWithContext(ctx, method, uri, reqStream)
	if err != nil {
		return fmt.Errorf("create request (%s %s): %w", method, uri, err)
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("do request (%s %s): %w", method, uri, err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode < 200 || httpResponse.StatusCode > 299 {
		return fmt.Errorf("wrong status code (%s %s): %d", method, uri, httpResponse.StatusCode)
	}

	repData, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return fmt.Errorf("fail to read response body (%s %s): %w", method, uri, err)
	}

	if rep != nil && len(repData) > 0 {
		if wrapResponse {
			repObj := &ResponseWithData{
				Data: rep,
			}
			err = json.Unmarshal(repData, repObj)
			if err != nil {
				return fmt.Errorf("fail to decode json (%s %s): %w, %s", method, uri, err, string(repData))
			}
			if !repObj.Success {
				return fmt.Errorf("success: false (%s %s), %s", method, uri, string(repData))
			}
		} else {
			err = json.Unmarshal(repData, rep)
			if err != nil {
				return fmt.Errorf("fail to decode json (%s %s): %w, %s", method, uri, err, string(repData))
			}
		}
	} else {
		repObj := &BaseResponse{}
		err = json.Unmarshal(repData, repObj)
		if err != nil {
			return fmt.Errorf("fail to decode json (%s %s): %w, %s", method, uri, err, string(repData))
		}
		if !repObj.Success {
			return fmt.Errorf("success: false (%s %s), %s", method, uri, string(repData))
		}
	}

	return nil
}
