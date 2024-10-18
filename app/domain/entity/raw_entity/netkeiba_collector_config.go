package raw_entity

type NetKeibaCollectorConfigs struct {
	Cookies []Cookie `json:"cookies"`
}

type Cookie struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Path     string `json:"path"`
	Domain   string `json:"domain"`
	Secure   bool   `json:"secure"`
	HttpOnly bool   `json:"http_only"`
}
