package main

import (
	"log"
	"net/http"
	"os"
	"path"

	"golang.org/x/net/webdav"
	"gopkg.in/yaml.v2"
)

// WebDAVConfig 定义单个 WebDAV 服务的配置
type WebDAVConfig struct {
	RootDir  string `yaml:"root_dir"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Mode     string `yaml:"mode"` // 支持 "readonly" 或 "readwrite"
}

// Config 定义整个配置文件结构
type Config struct {
	Port     string                  `yaml:"port"`
	Services map[string]WebDAVConfig `yaml:"services"`
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func readOnlyMiddleware(next http.Handler, isReadOnly bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isReadOnly {
			switch r.Method {
			case "PUT", "POST", "DELETE", "MKCOL", "PROPPATCH", "MOVE", "COPY":
				http.Error(w, "Forbidden: Read-only mode", http.StatusForbidden)
				log.Printf("Blocked %s request in read-only mode: %v", r.Method, r.URL)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// withBasicAuth 实现基本认证中间件
func withBasicAuth(next http.Handler, username, password string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != username || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// 加载配置文件
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 设置 HTTP 服务器
	mux := http.NewServeMux()

	// 动态注册每个 WebDAV 服务
	for name, cfg := range config.Services {
		// 检查 root_dir 是否存在
		if _, err := os.Stat(cfg.RootDir); os.IsNotExist(err) {
			log.Printf("Warning: Root directory %s for %s does not exist, creating it...", cfg.RootDir, name)
			if err := os.MkdirAll(cfg.RootDir, 0755); err != nil {
				log.Fatalf("Failed to create directory %s: %v", cfg.RootDir, err)
			}
		}

		// 创建 WebDAV 处理程序
		prefix := path.Join("/", name) + "/"
		handler := &webdav.Handler{
			FileSystem: webdav.Dir(cfg.RootDir),
			Prefix:     prefix,
			LockSystem: webdav.NewMemLS(),
			Logger: func(r *http.Request, err error) {
				if err != nil {
					log.Printf("WEBDAV [%s %s]: %v, ERROR: %v\n", prefix, r.Method, r.URL, err)
				} else {
					log.Printf("WEBDAV [%s %s]: %v\n", prefix, r.Method, r.URL)
				}
			},
		}

		// 判断是否为只读模式
		isReadOnly := cfg.Mode == "readonly"

		// 应用中间件：先认证，再检查只读模式
		authHandler := withBasicAuth(handler, cfg.Username, cfg.Password)
		finalHandler := readOnlyMiddleware(authHandler, isReadOnly)

		// 注册到 mux
		mux.Handle(prefix, finalHandler)
		log.Printf("Registered WebDAV service at %s with root %s (mode: %s)", prefix, cfg.RootDir, cfg.Mode)
	}

	// 启动服务器
	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: mux,
	}

	log.Printf("WebDAV server starting on :%s...", config.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
