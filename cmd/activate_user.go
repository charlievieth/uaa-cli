package cmd

import (
	"errors"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"code.cloudfoundry.org/uaa-cli/utils"
	"github.com/spf13/cobra"
)

func ActivateUserCmd(um uaa.UserManager, printer cli.Printer, username, origin, attributes string) error {
	user, err := um.GetByUsername(username, origin, attributes)
	if err != nil {
		return err
	}
	if user.Meta == nil {
		return errors.New("The user did not have expected metadata version.")
	}
	err = um.Activate(user.ID, user.Meta.Version)
	if err != nil {
		return err
	}
	log.Infof("Account for user %v successfully activated.", utils.Emphasize(user.Username))

	return nil
}

func ActivateUserValidations(cfg uaa.Config, args []string) error {
	if err := EnsureContextInConfig(cfg); err != nil {
		return err
	}

	if len(args) == 0 {
		return errors.New("The positional argument USERNAME must be specified.")
	}
	return nil
}

var activateUserCmd = &cobra.Command{
	Use:   "activate-user USERNAME",
	Short: "Activate a user by username",
	PreRun: func(cmd *cobra.Command, args []string) {
		err := ActivateUserValidations(GetSavedConfig(), args)
		NotifyValidationErrors(err, cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		um := uaa.UserManager{GetHttpClient(), cfg}
		err := ActivateUserCmd(um, cli.NewJsonPrinter(log), args[0], origin, attributes)
		NotifyErrorsWithRetry(err, cfg, log)
	},
}

func init() {
	RootCmd.AddCommand(activateUserCmd)
	activateUserCmd.Annotations = make(map[string]string)
	activateUserCmd.Annotations[USER_CRUD_CATEGORY] = "true"
}
