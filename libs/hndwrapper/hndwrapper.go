package hndwrapper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type BaseResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type ResponseWithData struct {
	BaseResponse
	Data any `json:"data"`
}

type Wrapper[Req any, Rep any] struct {
	fn func(ctx context.Context, req *Req) (*Rep, error)
}

func New[Req any, Rep any](fn func(ctx context.Context, req *Req) (*Rep, error)) *Wrapper[Req, Rep] {
	return &Wrapper[Req, Rep]{
		fn: fn,
	}
}

func (o *Wrapper[Req, Rep]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bodyRaw, err := io.ReadAll(r.Body)
	if handleError(w, err, "read body") {
		return
	}

	log.Println("Incoming request", r.Method, r.URL.Path+":\n\t\t", string(bodyRaw))

	var req *Req

	if len(bodyRaw) > 0 {
		req = new(Req)
		err = json.Unmarshal(bodyRaw, req)
		if handleError(w, err, "parse json") {
			return
		}
	}

	result, err := o.fn(r.Context(), req)
	if handleError(w, err, "domain") {
		return
	}

	if result != nil {
		sendJson(w, ResponseWithData{
			BaseResponse: BaseResponse{
				Success: true,
			},
			Data: result,
		})
	} else {
		sendJson(w, BaseResponse{
			Success: true,
		})
	}
}

func handleError(w http.ResponseWriter, err error, msg string) bool {
	if err == nil {
		return false
	}

	sendJson(w, BaseResponse{
		Success: false,
		Error:   fmt.Errorf(msg+": %w", err).Error(),
	})

	return true
}

func sendJson(w http.ResponseWriter, obj any) {
	rawData, err := json.Marshal(&obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal error"))
	}

	_, _ = w.Write(rawData)
}
