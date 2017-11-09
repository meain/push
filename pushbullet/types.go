package pushbullet

type header struct {
	name  string
	value string
}

type deviceResponse struct {
	Devices []Device `json:"devices"`
}
