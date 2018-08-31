package models

// SystemVersion contains information about the system version.
// See https://www.zevenet.com/zapidoc_ce_v3.1/#show-version
type SystemVersion struct {
	ApplianceVersion string `json:"appliance_version"`
	Hostname         string `json:"hostname"`
	KernelVersion    string `json:"kernel_version"`
	SystemDate       string `json:"system_date"`
	ZevenetVersion   string `json:"zevenet_version"`
	RestApiVersion   string `json:"rest_api_version"`
}
