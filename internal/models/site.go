package models

type Site struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Interval int    `json:"interval"`
}