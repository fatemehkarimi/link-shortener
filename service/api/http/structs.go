package http

import "time"

type RequestCreateLink struct {
	URL string `json:"URL"`
}

type ResponseCreateLink struct {
	Hash       string    `json:"hash"`
	Create_at  time.Time `json:"create_at"`
	Expires_at time.Time `json:"expires_at"`
}

type RequestGetURLByHash struct {
	Hash string `json:"hash"`
}

type ResponseGetURLByHash struct {
	URL *string `json:"URL"`
}
