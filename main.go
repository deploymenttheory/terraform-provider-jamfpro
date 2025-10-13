package main

import (
	"context"
	"flag"
	"log"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

const noLogPrefix = 0

func main() {
	ctx := context.Background()

	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	// Prevent logger from prepending date/time to logs, which breaks log-level parsing/filtering
	log.SetFlags(noLogPrefix)

	// Upgrade SDKv2 provider from protocol 5 to protocol 6
	upgradedSdkProvider, err := tf5to6server.UpgradeServer(
		ctx,
		provider.Provider().GRPCProvider,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create the Framework provider (protocol 6)
	frameworkProvider := provider.FrameworkProvider("dev")()

	// Mux both providers together using protocol 6
	muxServer, err := tf6muxserver.NewMuxServer(ctx,
		func() tfprotov6.ProviderServer { return upgradedSdkProvider },
		providerserver.NewProtocol6(frameworkProvider),
	)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf6server.ServeOpt

	if debug {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	err = tf6server.Serve(
		"registry.terraform.io/deploymenttheory/jamfpro",
		muxServer.ProviderServer,
		serveOpts...,
	)

	if err != nil {
		log.Fatal(err)
	}
}
