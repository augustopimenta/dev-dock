package docker

func StartProxyContainer() {
	api := Api{}
	api.Init()
	api.Run(&Container{
		"",
		"reverse-proxy",
		"jwilder/nginx-proxy",
		map[string]string{"80/tcp":"80"},
		[]string{"/var/run/docker.sock:/tmp/docker.sock:ro"},
		[]string{},
	})
}
