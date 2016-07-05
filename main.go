package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cloudfoundry/cli/cf/commandregistry"
	"github.com/cloudfoundry/cli/cf/trace"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/fatih/color"
)

type SetTokenPlugin struct {
	accessToken  string
	refreshToken string
}

func (p *SetTokenPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	p.parseParameters(args)
	traceLogger := trace.NewLogger(color.Output, false, "true", "")
	deps := commandregistry.NewDependency(os.Stdout, traceLogger, "")
	p.loadAndModifyConfig(deps)
}

func (p *SetTokenPlugin) parseParameters(args []string) {
	if args[0] != "set-token" {
		p.exitWithUsage(fmt.Errorf("Wrong command: %s", args[0]))
	}
	flagSet := flag.NewFlagSet("set-token", flag.ContinueOnError)
	flagSet.StringVar(&p.accessToken, "a", "", "Authccess token")
	flagSet.StringVar(&p.refreshToken, "r", "", "Refresh token")
	err := flagSet.Parse(args[1:])
	if err != nil || len(flagSet.Args()) > 1 {
		p.exitWithUsage(err)
	}
	if p.accessToken == "" || p.refreshToken == "" {
		p.exitWithUsage(fmt.Errorf("You must set either access or refresh token, or both of them"))
	}
}

func (p *SetTokenPlugin) loadAndModifyConfig(deps commandregistry.Dependency) {
	config := deps.Config
	config.SetAccessToken(p.accessToken)
	config.SetRefreshToken(p.refreshToken)
}

func (p *SetTokenPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "SetTokenPlaugin",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "set-token",
				HelpText: "Allows you to manually set authentication and refresh tokens.",

				UsageDetails: plugin.Usage{
					Usage: "cf set-token [-a ACCESS_TOKEN] [-r REFRESH_TOKEN]",
					Options: map[string]string{
						"-a": "Access token",
						"-r": "Refresh token",
					},
				},
			},
		},
	}
}

func (p *SetTokenPlugin) exitWithUsage(err error) {
	metadata := p.GetMetadata()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Usage: " + metadata.Commands[0].UsageDetails.Usage)
	os.Exit(1)
}

func main() {
	plugin.Start(new(SetTokenPlugin))
}
