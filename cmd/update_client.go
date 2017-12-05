package cmd

import (
	"errors"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"code.cloudfoundry.org/uaa-cli/utils"
	"github.com/spf13/cobra"
)

func UpdateClientValidations(cfg uaa.Config, args []string, clientSecret string) error {
	if err := EnsureContextInConfig(cfg); err != nil {
		return err
	}
	if len(args) < 1 {
		return MissingArgumentError("client_id")
	}
	if clientSecret != "" {
		return errors.New(`Client not updated. Please see "uaa set-client-secret -h" to learn more about changing client secrets.`)
	}
	return nil
}

func UpdateClientCmd(cm *uaa.ClientManager, clientId, displayName, authorizedGrantTypes, authorities, redirectUri, scope string, accessTokenValidity, refreshTokenValidity int64) error {
	toUpdate := uaa.UaaClient{
		ClientId:             clientId,
		DisplayName:          displayName,
		AuthorizedGrantTypes: arrayify(authorizedGrantTypes),
		Authorities:          arrayify(authorities),
		RedirectUri:          arrayify(redirectUri),
		Scope:                arrayify(scope),
		AccessTokenValidity:  accessTokenValidity,
		RefreshTokenValidity: refreshTokenValidity,
	}

	updated, err := cm.Update(toUpdate)
	if err != nil {
		return errors.New("An error occurred while updating the client.")
	}

	log.Infof("The client %v has been successfully updated.", utils.Emphasize(clientId))
	return cli.NewJsonPrinter(log).Print(updated)

}

var updateClientCmd = &cobra.Command{
	Use:   "update-client CLIENT_ID",
	Short: "Update an OAuth client registration in the UAA",
	PreRun: func(cmd *cobra.Command, args []string) {
		NotifyValidationErrors(UpdateClientValidations(GetSavedConfig(), args, clientSecret), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		cm := &uaa.ClientManager{GetHttpClient(), cfg}
		NotifyErrorsWithRetry(UpdateClientCmd(cm, args[0], displayName, authorizedGrantTypes, authorities, redirectUri, scope, accessTokenValidity, refreshTokenValidity), cfg, log)
	},
}

func init() {
	RootCmd.AddCommand(updateClientCmd)
	updateClientCmd.Annotations = make(map[string]string)
	updateClientCmd.Annotations[CLIENT_CRUD_CATEGORY] = "true"
	updateClientCmd.Flags().StringVarP(&clientSecret, "client_secret", "s", "", "client secret")
	updateClientCmd.Flag("client_secret").Hidden = true

	updateClientCmd.Flags().StringVarP(&authorizedGrantTypes, "authorized_grant_types", "", "", "list of grant types allowed with this client.")
	updateClientCmd.Flags().StringVarP(&authorities, "authorities", "", "", "scopes requested by client during client_credentials grant")
	updateClientCmd.Flags().StringVarP(&scope, "scope", "", "", "scopes requested by client during authorization_code, implicit, or password grants")
	updateClientCmd.Flags().Int64VarP(&accessTokenValidity, "access_token_validity", "", 0, "the time in seconds before issued access tokens expire")
	updateClientCmd.Flags().Int64VarP(&refreshTokenValidity, "refresh_token_validity", "", 0, "the time in seconds before issued refrsh tokens expire")
	updateClientCmd.Flags().StringVarP(&displayName, "display_name", "", "", "a friendly human-readable name for this client")
	updateClientCmd.Flags().StringVarP(&redirectUri, "redirect_uri", "", "", "callback urls allowed for use in authorization_code and implicit grants")
	updateClientCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain in which to update the client")
}
