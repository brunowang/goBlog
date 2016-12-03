package engine

var config map[string]string

func init() {
	config = make(map[string]string)
	config["adminName"] = "admin"
	config["adminPass"] = "admin"
}
func Config(key string) string {
	return config[key]
}
