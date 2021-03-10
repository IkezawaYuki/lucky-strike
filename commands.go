package main

import (
	"github.com/IkezawaYuki/lucky-strike/internal/terminal"
	"github.com/hashicorp/terraform-svchost/disco"
	"github.com/hashicorp/terraform/command/cliconfig"
	"github.com/mitchellh/cli"
)

const runningInAutomationEnvMode = "TF_IN_AUTOMATION"

var Commands map[string]cli.CommandFactory

var PrimaryCommands []string

var HiddenCommands map[string]struct{}

var Ui cli.Ui

func initCommands(
	originalWorkingDir string,
	streams *terminal.Streams,
	config *cliconfig.Config,
	services *disco.Disco,
	providerSrc getproviders.Source,

)
