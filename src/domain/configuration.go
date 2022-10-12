package domain

// Configuration /*
type Configuration struct {
	ProxyServer    ProxyServer `yaml:"proxy_server"`
	ChromeExecPath string      `yaml:"chrome_exec_path"`
	UserAgent      string      `yaml:"user_agent"`
	Headless       bool        `yaml:"headless"`
	Timeout        int         `yaml:"timeout"`
	Cookies        string      `yaml:"cookies"`
	XrayExt        XrayExt     `yaml:"xray_ext"`
	Worker         int         `yaml:"worker"`
}
type ProxyServer struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
type XrayExt struct {
	Enabled  bool   `yaml:"enabled"`
	XrayPath string `yaml:"xray_path"`
}
