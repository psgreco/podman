package main

import (
	"os"

	"github.com/containers/common/pkg/auth"
	"github.com/containers/image/v5/types"
	"github.com/containers/podman/v2/cmd/podman/registry"
	"github.com/containers/podman/v2/pkg/domain/entities"
	"github.com/containers/podman/v2/pkg/registries"
	"github.com/spf13/cobra"
)

var (
	logoutOptions = auth.LogoutOptions{}
	logoutCommand = &cobra.Command{
		Use:   "logout [flags] [REGISTRY]",
		Short: "Logout of a container registry",
		Long:  "Remove the cached username and password for the registry.",
		RunE:  logout,
		Args:  cobra.MaximumNArgs(1),
		Example: `podman logout quay.io
  podman logout --authfile dir/auth.json quay.io
  podman logout --all`,
	}
)

func init() {
	// Note that the local and the remote client behave the same: both
	// store credentials locally while the remote client will pass them
	// over the wire to the endpoint.
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Mode:    []entities.EngineMode{entities.ABIMode, entities.TunnelMode},
		Command: logoutCommand,
	})
	flags := logoutCommand.Flags()

	// Flags from the auth package.
	flags.AddFlagSet(auth.GetLogoutFlags(&logoutOptions))
	logoutOptions.Stdout = os.Stdout
	logoutOptions.AcceptUnspecifiedRegistry = true
}

// Implementation of podman-logout.
func logout(cmd *cobra.Command, args []string) error {
	sysCtx := types.SystemContext{
		AuthFilePath:             logoutOptions.AuthFile,
		SystemRegistriesConfPath: registries.SystemRegistriesConfPath(),
	}
	return auth.Logout(&sysCtx, &logoutOptions, args)
}
