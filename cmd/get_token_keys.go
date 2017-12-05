package cmd

import (
	"net/http"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"github.com/spf13/cobra"
)

func GetTokenKeysCmd(client *http.Client, config uaa.Config) error {
	key, err := uaa.TokenKeys(client, config)

	if err != nil {
		return err
	}

	return cli.NewJsonPrinter(log).Print(key)
}

var getTokenKeysCmd = &cobra.Command{
	Use:     "get-token-keys",
	Short:   "View all keys the UAA has used to sign JWT tokens",
	Aliases: []string{"token-keys"},
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyValidationErrors(EnsureTargetInConfig(cfg), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		NotifyErrorsWithRetry(GetTokenKeysCmd(GetHttpClient(), GetSavedConfig()), GetSavedConfig(), log)
	},
}

func init() {
	RootCmd.AddCommand(getTokenKeysCmd)
	getTokenKeysCmd.Annotations = make(map[string]string)
	getTokenKeysCmd.Annotations[TOKEN_CATEGORY] = "true"
}
