package webserver

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	instance        *WebServer
	instancePort    int
	instanceFileDir string
)

func Start() (string, error) {
	go func() {
		err := StartWebServer(0, "")
		if err != nil {
			fmt.Println("server failed to start", err)
		}
	}()

	n := 20
	for {
		time.Sleep(100 * time.Millisecond)
		if instance != nil && instance.IsRunning() {
			return instance.ServerURL(), nil
		}
		n -= 1
		if n <= 0 {
			return "", fmt.Errorf("server failed to start")
		}
	}
}

func Stop() {
	if instance != nil && instance.IsRunning() {
		instance.Stop()
	}
}

// NOTE: make sure to call this method before starting the server
func SetPort(port int) {
	instancePort = port
}

// NOTE: make sure to call this method before starting the server
// TODO: rename this to SetFileDir
func SetConfig(fileDir string) {
	instanceFileDir = fileDir
}

func Url() string {
	if instance != nil {
		return instance.ServerURL()
	}
	return ""
}

func StartWebServer(port int, fileDir string) error {
	ws, err := NewWebServer(port, fileDir)
	if err != nil {
		return err
	}
	return ws.Run()
}

func GetWebServerPort() (int, error) {
	if instance == nil {
		return 0, fmt.Errorf("server is not initialized")
	}
	if !instance.IsRunning() {
		return 0, fmt.Errorf("server is not running")
	}
	return instance.Port(), nil
}

type WebServer struct {
	server    *http.Server
	listener  net.Listener
	serverURL string
	port      int
	fileDir   string
	running   bool
}

func NewWebServer(port int, fileDir string) (*WebServer, error) {
	if instancePort > 0 && port == 0 {
		port = instancePort
	}
	if instanceFileDir != "" && fileDir == "" {
		fileDir = instanceFileDir
	}
	ws := &WebServer{
		port:    port,
		fileDir: fileDir,
	}
	instance = ws
	return ws, nil
}

func (s *WebServer) Run() error {
	var err error

	// NOTE: if port is 0, then it will use a randomly available port
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	if s.port == 0 {
		// determine randomly assigned port
		tcpAddr, ok := s.listener.Addr().(*net.TCPAddr)
		if !ok {
			return fmt.Errorf("unable to determine randomly assigned tcp port")
		}
		port := tcpAddr.Port
		s.port = int(port)
	}

	fmt.Println("listening on port", s.port)

	s.running = true
	defer func() {
		s.running = false
	}()

	// TODO: use s.server.ServeTLS, etc..
	addr := fmt.Sprintf("127.0.0.1:%d", s.port)

	s.serverURL = fmt.Sprintf("http://%s", addr)

	s.server = &http.Server{
		Addr:              addr,
		Handler:           s.Routes(),
		ReadTimeout:       45 * time.Second,
		WriteTimeout:      45 * time.Second,
		IdleTimeout:       45 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	err = s.server.Serve(s.listener)
	if err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *WebServer) Stop() error {
	if s.server != nil {
		err := s.server.Close()
		if err != nil {
			return err
		}
		err = s.listener.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *WebServer) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.NoCache)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("."))
	})

	FileServer(r, "/", http.Dir(s.fileDir))

	return r
}

func (s *WebServer) IsRunning() bool {
	return s.running
}

func (s *WebServer) ServerURL() string {
	return s.serverURL
}

func (s *WebServer) Port() int {
	return s.port
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Head(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
