package common

import cd "github.com/muidea/magicCommon/def"

const GPTModule = "/module/gpt"

const (
	QueryMessage = "/message/query/"
)

type QueryParam struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type QueryResult struct {
	cd.Result
	Data string `json:"data"`
}
