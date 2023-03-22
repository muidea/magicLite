package model

type Log struct {
	ID         int    `json:"id" orm:"id key auto" tag:"view"`
	Address    string `json:"address" orm:"address" tag:"view"`
	Memo       string `json:"memo" orm:"memo" tag:"view"`
	Creater    int    `json:"creater" orm:"creater" tag:"view"`
	CreateTime int64  `json:"createTime" orm:"createTime" tag:"view"`
	Namespace  string `json:"namespace" orm:"namespace"`
}
