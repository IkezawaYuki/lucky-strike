package main

import "github.com/mitchellh/cli"

const runningInAutomationEnvMode = "TF_IN_AUTOMATION"

var Commands map[string]cli.CommandFactory
