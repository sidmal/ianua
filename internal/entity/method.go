package entity

type Method struct {
	Url                  string
	RequestMethod        string
	RequestBody          string
	RequestHeaders       []map[string]string
	SecurityHashTemplate string
}
