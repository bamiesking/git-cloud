package utils

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
)

// https://stackoverflow.com/questions/39320371/how-start-web-server-to-open-page-in-browser-in-golang
// open opens the specified URL in the default browser of the user.
func Open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func GitRepoPath() (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", "git rev-parse --show-toplevel")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	if stderr.Len() > 0 {
		err := errors.New(stderr.String())
		return "", err
	}
	path := bytes.TrimRight(stdout.Bytes(), "\n")
	return string(path), nil
}

type (
	generator func(http.ResponseWriter, *http.Request, chan<- string)
	handler   func(http.ResponseWriter, *http.Request)
)

func OAuthHandlerGenerator(w http.ResponseWriter, r *http.Request, t chan<- string) {
	code := r.URL.Query().Get("code")
	fmt.Println(code)
	t <- code
	w.Write([]byte("You can now close this tab"))
}

func GenerateHandler(g generator, t chan<- string) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		g(w, r, t)
	}
}

func AwaitOAuthRedirect(g generator, port chan<- int, t chan string) {
	listener, err := net.Listen("tcp", ":0")
	Handle(err)
	port <- listener.Addr().(*net.TCPAddr).Port
	h := GenerateHandler(g, t)
	http.HandleFunc("/", h)
	http.Serve(listener, nil)
}

func Handle(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
