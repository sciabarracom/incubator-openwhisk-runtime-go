package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

// LastError seen
var LastError error

// FatalIf panics if error is not nil
func FatalIf(err error) {
	LastError = err
	if err != nil {
		log.Error(err)
		panic(err)
	}
}

// LogIf logs a warning if the error is not nil
// returns true if the err is not nil
func LogIf(err error) bool {
	LastError = err
	if err != nil {
		log.Warn(err)
		return true
	}
	return false
}

// Recover recovers from a panic returning an error
func Recover(perr *error) {
	if r := recover(); r != nil {
		switch v := r.(type) {
		case error:
			*perr = v
		case string:
			*perr = fmt.Errorf("%s", v)
		default:
			*perr = fmt.Errorf("%v", v)
		}
	}
}

// DryRunFlag allows for dry execution
var DryRunFlag = kingpin.Flag("dry-run", "show commands instead of executing them").
	Short('n').Default("false").Bool()

// buffer for dry run rusults
var dryRunBuf = []string{}

// DryRunPush saves dummy results for dry run execution
func DryRunPush(buf ...string) {
	dryRunBuf = buf
}

// DryRunPop returns a value from the buffer of dry run results
// returns an empty string if the  buffer is empty
func DryRunPop(buf ...string) string {
	res := ""
	if len(dryRunBuf) > 0 {
		res = dryRunBuf[0]
		dryRunBuf = dryRunBuf[1:]
	}
	return res
}

// Sys executes a command in a convenient way:
// it splits the paramenter in arguments if separated by spaces,
// then accepts multiple arguments;
// logs errors in stderr and prints output in stdout;
// also returns output as a string, or empty if errors.
// If the command starts with "@" do not print the output.
// It also honor the DryRunFlag, in this case it always print the command
func Sys(cli string, args ...string) string {
	a := strings.Split(strings.ReplaceAll(cli, "\n", " "), " ")
	params := args
	if len(a) > 1 {
		params = append(a[1:], args...)
	}
	exe := strings.TrimPrefix(a[0], "@")
	silent := strings.HasPrefix(a[0], "@")

	if *DryRunFlag {
		if len(params) > 0 {
			fmt.Printf("%s %s\n", exe, strings.Join(params, " "))
		} else {
			fmt.Println(exe)
		}
		return DryRunPop()
	}

	log.Tracef("< %s %v\n", exe, params)
	cmd := exec.Command(exe, params...)
	out, err := cmd.CombinedOutput()
	res := string(out)
	log.Tracef("> %s", res)
	if !LogIf(err) && !silent {
		fmt.Printf(res)
	}
	return res
}

// SysCd works as Sys,
// but cd into the directory first,
// and restore the original directory after.
func SysCd(cd string, cli string, args ...string) string {
	orig, err := os.Getwd()
	FatalIf(err)
	FatalIf(os.Chdir(cd))
	res := Sys(cli, args...)
	FatalIf(os.Chdir(orig))
	return res
}

// SysSh execute a command as a shell script;
// if it starts with "@" it does not print the output;
// returns the output as a string.
func SysSh(cmd string) string {
	if strings.HasPrefix(cmd, "@") {
		return Sys("@sh -c", strings.TrimPrefix(cmd, "@"))
	}
	return Sys("sh -c", cmd)
}
