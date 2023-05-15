package http

func GetBasicAuth(username string, password string) map[string]string {
	return map[string]string{"username": username, "password": password}
}
