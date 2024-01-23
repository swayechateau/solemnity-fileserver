package models

type Recovery struct {
	Email    string
	IP       string
	Domain   string
	Code     string
	AccessId string
}
