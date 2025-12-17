package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	v1 "github.com/premkit/premkit/handlers/v1"
	"github.com/premkit/premkit/models"
	"github.com/premkit/premkit/server"

	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Initialize a temp boltdb folder
	err := setup()
	if err != nil {
		os.Exit(1)
	}
	defer teardown()

	// Run the tests
	result := m.Run()

	os.Exit(result)
}

func waitForServer(url string, maxWait time.Duration) error {
	deadline := time.Now().Add(maxWait)
	delay := 10 * time.Millisecond

	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			return nil
		}

		time.Sleep(delay)
		delay *= 2
		if delay > 100*time.Millisecond {
			delay = 100 * time.Millisecond
		}
	}
	return fmt.Errorf("server did not become ready within %v", maxWait)
}

func setup() error {
	originalDataFile := viper.GetString("data_file")
	defer func() {
		viper.Set("data_file", originalDataFile)
		os.RemoveAll("/tmp/integration.db")
	}()
	viper.Set("data_file", "/tmp/integration.db")

	// Start a premkit server
	config := server.Config{
		HTTPPort:  9141,
		HTTPSPort: 0,

		TLSKeyFile:  "",
		TLSCertFile: "",
	}
	go server.Run(&config)

	// Wait for the server to be ready
	if err := waitForServer("http://localhost:9141/", time.Second); err != nil {
		return fmt.Errorf("premkit server failed to start: %w", err)
	}

	// Start a custom upstream echo server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s := make(map[string]interface{})
		s["query"] = r.URL.Query().Encode()
		b, err := json.Marshal(s)
		if err != nil {
			fmt.Printf("Error encoding response: %#v\n", err)
			os.Exit(1)
		}
		io.WriteString(w, string(b))
	})
	go func() {
		http.ListenAndServe(":9142", nil)
	}()

	// Wait for the echo server to be ready
	if err := waitForServer("http://localhost:9142/", time.Second); err != nil {
		return fmt.Errorf("echo server failed to start: %w", err)
	}

	// Register that upstream
	service := models.Service{
		Name: "echo",
		Path: "/path/echo",
		Upstreams: []*models.Upstream{
			&models.Upstream{
				URL:                "http://localhost:9142/path/echo/",
				IncludeServicePath: false,
				InsecureSkipVerify: false,
			},
		},
	}
	registerParams := v1.RegisterServiceParams{
		Service:         &service,
		ReplaceExisting: true,
	}

	request := gorequest.New()
	resp, _, errs := request.Post("http://localhost:9141/premkit/v1/service").
		Send(registerParams).
		End()
	if len(errs) != 0 {
		fmt.Printf("Unexpected errs from registering service: %q", errs)
		os.Exit(1)
	}
	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Unxpected non-201 from registering service: %d", resp.StatusCode)
		os.Exit(1)
	}

	// fmt.Printf("Integration web server is up for the next 10 minutes...\n")
	// time.Sleep(time.Minute * 10)

	return nil
}

func teardown() {
	// When the tests finish, the web servers will stop...  this is hacky
}

// Test an upstream with a querystring
func TestQuerystring(t *testing.T) {
	// Register a dynamic upstream
	resp, body, errs := gorequest.New().Get("http://localhost:9141/path/echo?a=b").End()

	require.Equal(t, 0, len(errs), "there should be no errors")
	require.Equal(t, http.StatusOK, resp.StatusCode)

	r := make(map[string]interface{})
	err := json.Unmarshal([]byte(body), &r)
	require.NoError(t, err)

	fmt.Printf("r = %#v\n", r)
	assert.Equal(t, r["query"], "a=b")
}
