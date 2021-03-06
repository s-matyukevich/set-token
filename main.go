package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cloudfoundry/cli/cf/commandregistry"
	"github.com/cloudfoundry/cli/cf/trace"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/cloudfoundry/go-ccapi/v3/client"
	"github.com/fatih/color"
)

type SetTokenPlugin struct {
	accessToken  string
	refreshToken string
	client       string
	clientSecret string
}

func (p *SetTokenPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	p.parseParameters(args)
	traceLogger := trace.NewLogger(color.Output, false, "true", "")
	deps := commandregistry.NewDependency(os.Stdout, traceLogger, "")
	p.getAccessToken(deps)
}

func (p *SetTokenPlugin) parseParameters(args []string) {
	if args[0] != "set-token" {
		p.exitWithUsage(fmt.Errorf("Wrong command: %s", args[0]))
	}
	flagSet := flag.NewFlagSet("set-token", flag.ContinueOnError)
	flagSet.StringVar(&p.refreshToken, "r", "", "Refresh token")
	flagSet.StringVar(&p.client, "c", "", "CF OAuth client")
	flagSet.StringVar(&p.clientSecret, "s", "", "CF OAuth client secret")
	err := flagSet.Parse(args[1:])
	if err != nil || len(flagSet.Args()) > 1 {
		p.exitWithUsage(err)
	}
	if p.refreshToken == "" {
		p.exitWithUsage(fmt.Errorf("Refresh token is not set"))
	}
	if p.client == "" {
		p.client = "cf"
	}
}

func (p *SetTokenPlugin) getAccessToken(deps commandregistry.Dependency) {
	config := deps.Config
	refresher := client.NewTokenRefresher(config.AuthenticationEndpoint(), p.client, p.clientSecret)
	var err error
	p.accessToken, p.refreshToken, err = refresher.Refresh(p.refreshToken)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	config.SetAccessToken("bearer " + p.accessToken)
	config.SetRefreshToken(p.refreshToken)
	config.SetCFOAuthClient(p.client)
	config.SetCFOAuthClientSecret(p.clientSecret)
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
					Usage: "cf set-token [-r REFRESH_TOKEN] [-c OAUTH_CLIENT] [-p OAUTH_CLIENT_PASSWORD]",
					Options: map[string]string{
						"r": "Refresh token",
						"c": "CF OAuth client",
						"s": "CF OAuth client secret",
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
