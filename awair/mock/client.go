package mock

import (
  "context"
	"crypto/tls"
  "fmt"
  "io/ioutil"
  "net"
  "net/http"
  "net/http/httptest"
  "os"
  "path/filepath"
  "strings"
)

const ACCESS_TOKEN = "abc123"

type MockServer struct {
  Server *httptest.Server
  mux *http.ServeMux
  testdataPath string
}

func checkAuthorization(r *http.Request) bool {
  return r.Header.Get("Authorization") == fmt.Sprintf("Bearer %s", ACCESS_TOKEN)
}

func getTestDataPath() string {
  dir, _ := os.Getwd()
  parts := strings.Split(dir, string(filepath.Separator))
  i := 0
  for i = 0; i < len(parts); i++ {
    if parts[i] == "awair-exporter" {
      break
    }
  }

  return fmt.Sprintf("%c%s", filepath.Separator, filepath.Join(parts[0:i+1]...))
}

func NewMockServer() *MockServer {
  m := &MockServer{}
  m.mux = http.NewServeMux()
  m.mux.Handle("/v1/users/self/devices", m.checkAuthMiddleware(m.serveFile("Devices.json")))
  m.mux.Handle("/v1/users/self", m.checkAuthMiddleware(m.serveFile("UserInfo.json")))
  m.mux.Handle("/v1/users/self/devices/", m.checkAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if strings.HasSuffix(r.URL.Path, "/awair/0/air-data/latest") {
      m.serveFile("UserLatestAirData-awair-0.json").ServeHTTP(w, r)
    } else if strings.HasSuffix(r.URL.Path, "/awair-r2/0/air-data/latest") {
      m.serveFile("UserLatestAirData-awair-r2-0.json").ServeHTTP(w, r)
    } else if strings.HasSuffix(r.URL.Path, "/api-usages") {
      m.serveFile("DeviceAPIUsage.json").ServeHTTP(w, r)
    } else {
      w.WriteHeader(404)
    }
  })))

  m.testdataPath = filepath.Join(getTestDataPath(), "testdata")

  m.Server = httptest.NewTLSServer(m.mux)
  return m
}

func (m *MockServer) Close() {
  m.Server.Close()
}

func (m *MockServer) Client() *http.Client {
  return &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, m.Server.Listener.Addr().String())
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

func (m *MockServer) checkAuthMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if !checkAuthorization(r) {
      w.WriteHeader(403)
      return
    }

    next.ServeHTTP(w, r)
  })
}

func (m *MockServer) serveFile(path string) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    data, _ := ioutil.ReadFile(filepath.Join(m.testdataPath, path))
    w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    w.Write([]byte(data))
  })
}
