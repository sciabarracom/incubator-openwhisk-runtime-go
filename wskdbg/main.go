package main

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

// VerboseFlag is flag for verbose
var VerboseFlag = kingpin.Flag("verbose", "Verbose").
	Short('v').Default("false").Bool()

var (
	debugCmd          = kingpin.Command("debug", "debug")
	debugInitCmd      = debugCmd.Command("init", "init")
	debugMyIPCmd      = debugCmd.Command("myip", "myip")
	debugInitIPArg    = debugInitCmd.Arg("ip", "ip").Required().String()
	debugInitFileArg  = debugInitCmd.Arg("code", "code").Required().String()
	debugInitMainArg  = debugInitCmd.Arg("main", "main").Default("main").String()
	debugInitEnvArg   = debugInitCmd.Arg("env", "env").Default("{}").String()
	debugRunCmd       = debugCmd.Command("run", "run")
	debugRunIPArg     = debugRunCmd.Arg("ip", "ip").Required().String()
	debugRunJSONArg   = debugRunCmd.Arg("json", "json").Default("{}").String()
	debugStartCmd     = debugCmd.Command("start", "start")
	debugStartIPArg   = debugStartCmd.Arg("ip", "ip").Required().String()
	debugStartJSONArg = debugStartCmd.Arg("json", "json").Default("{}").String()
	debugFwdCmd       = debugCmd.Command("fwd", "fwd")
	debugFwdIPCmd     = debugFwdCmd.Arg("ip", "ip").Required().String()
)

var (
	initCmd    = kingpin.Command("init", "init")
	runCmd     = kingpin.Command("run", "run")
	serverFlag = kingpin.Flag("server")
)

// Main entrypoint for wskide
func Main() {
	cmd := kingpin.Parse()
	if *VerboseFlag {
		log.SetLevel(log.TraceLevel)
	}
	//if !invokeParse(cmd) {
	//	kingpin.Usage()
	//}
}
