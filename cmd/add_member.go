package cmd

import (
	"errors"
	"net/http"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"code.cloudfoundry.org/uaa-cli/utils"
	"github.com/spf13/cobra"
)

func AddMemberPreRunValidations(config uaa.Config, args []string) error {
	if err := EnsureContextInConfig(config); err != nil {
		return err
	}

	if len(args) != 2 {
		return errors.New("The positional arguments GROUPNAME and USERNAME must be specified.")
	}

	return nil
}

func AddMemberCmd(httpClient *http.Client, config uaa.Config, groupName, username string, log cli.Logger) error {
	gm := uaa.GroupManager{httpClient, config}
	group, err := gm.GetByName(groupName, "")
	if err != nil {
		return err
	}

	um := uaa.UserManager{httpClient, config}
	user, err := um.GetByUsername(username, "", "")
	if err != nil {
		return err
	}

	err = gm.AddMember(group.ID, user.ID)
	if err != nil {
		return err
	}

	log.Infof("User %v successfully added to group %v", utils.Emphasize(username), utils.Emphasize(groupName))

	return nil
}

var addMemberCmd = &cobra.Command{
	Use:   "add-member GROUPNAME USERNAME",
	Short: "Add a user to a group",
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyValidationErrors(AddMemberPreRunValidations(cfg, args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyErrorsWithRetry(AddMemberCmd(GetHttpClient(), cfg, args[0], args[1], log), cfg, log)
	},
}

func init() {
	RootCmd.AddCommand(addMemberCmd)
	addMemberCmd.Annotations = make(map[string]string)
	addMemberCmd.Annotations[GROUP_CRUD_CATEGORY] = "true"
}
