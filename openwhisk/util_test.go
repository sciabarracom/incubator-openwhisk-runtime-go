package openwhisk

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"
)

func startTestServer(compiler string) (*httptest.Server, string, *os.File) {
	// temporary workdir
	cur, _ := os.Getwd()
	dir, _ := ioutil.TempDir("", "action")
	file, _ := filepath.Abs("_test")
	os.Symlink(file, dir+"/_test")
	os.Chdir(dir)
	log.Printf(dir)
	// setup the server
	buf, _ := ioutil.TempFile("", "log")
	ap := NewActionProxy(dir, compiler, buf, buf)
	ts := httptest.NewServer(ap)
	//log.Printf(ts.URL)
	//doPost(ts.URL+"/init", `{value: {code: ""}}`)
	return ts, cur, buf
}

func stopTestServer(ts *httptest.Server, cur string, buf *os.File) {
	runtime.Gosched()
	// wait 2 seconds before declaring a test done
	time.Sleep(2 * time.Second)
	os.Chdir(cur)
	ts.Close()
	dump(buf)
}

var testHTTPServerLastRequest string
var testHTTPServerLastBody string

func testHTTPServer(resp string) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		testHTTPServerLastRequest = r.Method + " " + r.URL.String()
		testHTTPServerLastBody = string(body)
		fmt.Fprintln(w, resp)
	}))
	return ts
}

func doGet(url string) (string, int, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", -1, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", -1, err
	}
	return string(body), res.StatusCode, nil
}

func doPost(url string, message string) (string, int, error) {
	buf := bytes.NewBufferString(message)
	res, err := http.Post(url, "application/json", buf)
	if err != nil {
		return "", -1, err
	}
	defer res.Body.Close()
	resp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", -1, err
	}
	return string(resp), res.StatusCode, nil
}

func doRun(ts *httptest.Server, message string) {
	if message == "" {
		message = `{"name":"Mike"}`
	}
	resp, status, err := doPost(ts.URL+"/run", `{ "value": `+message+`}`)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%d %s", status, resp)
	}
	if !strings.HasSuffix(resp, "\n") {
		fmt.Println()
	}
}

func doInit(ts *httptest.Server, message string) {
	resp, status, err := doPost(ts.URL+"/init", message)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%d %s", status, resp)
	}
}

func initCode(file string, main string) string {
	dat, _ := ioutil.ReadFile(file)
	body := initBodyRequest{Code: string(dat)}
	if main != "" {
		body.Main = main
	}
	j, _ := json.Marshal(initRequest{Value: body})
	return string(j)
}

func initBinaryEnv(file string, main string, env map[string]interface{}) string {
	dat, _ := ioutil.ReadFile(file)
	enc := base64.StdEncoding.EncodeToString(dat)
	body := initBodyRequest{Code: enc, Binary: true, Env: env}
	if main != "" {
		body.Main = main
	}
	j, _ := json.Marshal(initRequest{Value: body})
	return string(j)
}

func initBytes(dat []byte, main string) string {
	enc := base64.StdEncoding.EncodeToString(dat)
	body := initBodyRequest{Binary: true, Code: enc}
	if main != "" {
		body.Main = main
	}
	j, _ := json.Marshal(initRequest{Value: body})
	return string(j)
}

func initBinary(file string, main string) string {
	dat, _ := ioutil.ReadFile(file)
	return initBytes(dat, main)
}

func abs(in string) string {
	out, _ := filepath.Abs(in)
	return out
}

func dump(file *os.File) {
	//file.Read()
	buf, _ := ioutil.ReadFile(file.Name())
	fmt.Print(string(buf))
	//fmt.Print(file.ReadAll())
	os.Remove(file.Name())
}

func sys(cli string, args ...string) {
	os.Chmod(cli, 0755)
	cmd := exec.Command(cli, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(err)
	} else {
		fmt.Print(string(out))
	}
}

func exists(dir, filename string) error {
	path := fmt.Sprintf("%s/%d/%s", dir, highestDir(dir), filename)
	_, err := os.Stat(path)
	return err
}

func detectExecutable(dir, filename string) bool {
	path := fmt.Sprintf("%s/%d/%s", dir, highestDir(dir), filename)
	file, _ := ioutil.ReadFile(path)
	return IsExecutable(file, runtime.GOOS)
}

func waitabit() {
	time.Sleep(2000 * time.Millisecond)
}

func removeLineNr(out string) string {
	var re = regexp.MustCompile(`:\d+:\d+`)
	return re.ReplaceAllString(out, "::")
}

func grep(search string, data string) {
	re := regexp.MustCompile(search)
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		if re.Match([]byte(line)) {
			fmt.Println(strings.TrimSpace(line))
		}
	}
}
func replace(search string, replace string, data string) string {
	var re = regexp.MustCompile(search)
	return re.ReplaceAllString(data, replace)
}

func TestMain(m *testing.M) {
	Debugging = os.Getenv("VSCODE_PID") != "" // enable debug of tests
	sys("_test/build.sh")
	if !Debugging {
		// silence logging when not running under vscode for test
		log.SetOutput(ioutil.Discard)
	}
	// increase timeouts for init
	DefaultTimeoutStart = 1000 * time.Millisecond

	// build some test stuff
	// go ahead
	code := m.Run()
	os.Exit(code)
}
