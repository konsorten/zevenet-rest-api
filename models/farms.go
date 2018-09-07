package models

// FarmProfile is an enumeration of possible farm profiles.
type FarmProfile string

const (
	FarmProfile_HTTP      FarmProfile = "http"
	FarmProfile_HTTPS     FarmProfile = "https"
	FarmProfile_Level4NAT FarmProfile = "l4xnat"
	FarmProfile_DataLink  FarmProfile = "datalink"
)

// FarmInfo contains the list of all available farms.
// See https://www.zevenet.com/zapidoc_ce_v3.1/#list-all-farms
type FarmInfo struct {
	FarmName  string `json:"farmname" example:"test"`
	Profile   string/*FarmProfile*/ `json:"profile" example:"http"`
	VirtualIP string `json:"vip" example:"1.1.1.1"`
	Links     struct {
		Self   string `json:"self" example:"."`
		Detail string `json:"detail" example:"./http/test"`
	} `json:"links"`
}
