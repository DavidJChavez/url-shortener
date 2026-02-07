package model

import "time"

type URL struct {
	ID          int64     `json:"id"`
	Code        string    `json:"code"`
	OriginalURL string    `json:"originalURL"`
	Clicks      int64     `json:"clicks"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CreateURLRequest struct {
	URL string `json:"url"`
}

type CreateURLResponse struct {
	ShortURL    string `json:"shortURL"`
	OriginalURL string `json:"originalURL"`
	Code        string `json:"code"`
}

type StatsResponse struct {
	Clicks    int64     `json:"clicks"`
	CreatedAt time.Time `json:"createdAt"`
}
