package entities

type Device struct {
	DBModel
	Hostname         string `json:"hostname" binding:"required"`
	Username         string `json:"username"`
	UserID           string `json:"user_id" binding:"required"`
	OSName           string `json:"os_name" binding:"required"`
	OSArch           string `json:"os_arch" binding:"required"`
	MacAddress       string `json:"mac_address" binding:"required"`
	MacAddressBase64 string `json:"mac_address_base_64" gorm:"-"`
	LocalIPAddress   string `json:"local_ip_address"`
	Port             string `json:"port"`
	FetchedUnix      int64  `json:"fetched_unix" binding:"required"`
}
