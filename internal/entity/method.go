package entity

type Method struct {
	*Request
	HashTemplate string
}

type Request struct {
	Url                 string
	Headers             map[string]string
	ReqMethod           string
	ReqBody             string
	RspSuccessCode      int
	RspFormat           string
	RspResultFieldNames []string
}
