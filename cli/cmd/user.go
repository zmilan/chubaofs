// Copyright 2018 The Chubao Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/chubaofs/chubaofs/proto"
	"github.com/chubaofs/chubaofs/sdk/master"
	"github.com/spf13/cobra"
)

const (
	cmdUserUse   = "user [COMMAND]"
	cmdUserShort = "Manage cluster users"
)

func newUserCmd(client *master.MasterClient) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   cmdUserUse,
		Short: cmdUserShort,
		Args:  cobra.MinimumNArgs(0),
	}
	cmd.AddCommand(
		newUserCreateCmd(client),
		newUserInfoCmd(client),
		newUserListCmd(client),
		newUserPermCmd(client),
		newUserUpdateCmd(client),
		newUserDeleteCmd(client),
	)
	return cmd
}

const (
	cmdUserCreateUse   = "create [USER ID]"
	cmdUserCreateShort = "Create a new user"
)

func newUserCreateCmd(client *master.MasterClient) *cobra.Command {
	var optPassword string
	var optAccessKey string
	var optSecretKey string
	var optUserType string
	var optYes bool
	var cmd = &cobra.Command{
		Use:   cmdUserCreateUse,
		Short: cmdUserCreateShort,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var userID = args[0]
			var password = optPassword
			var accessKey = optAccessKey
			var secretKey = optSecretKey
			var userType = proto.UserTypeFromString(optUserType)

			if !userType.Valid() {
				errout("Invalid user type.")
				os.Exit(1)
			}

			// ask user for confirm
			if !optYes {
				// display information before create
				var displayPassword = "[default]"
				if optPassword != "" {
					displayPassword = optPassword
				}
				var displayAccessKey = "[auto generate]"
				var displaySecretKey = "[auto generate]"
				if optAccessKey != "" && optSecretKey != "" {
					displayAccessKey = optAccessKey
					displaySecretKey = optSecretKey
				}
				var displayUserType = userType.String()
				fmt.Printf("Create a new ChubaoFS cluster user\n")
				stdout("  User ID   : %v\n", userID)
				stdout("  Password  : %v\n", displayPassword)
				stdout("  Access Key: %v\n", displayAccessKey)
				stdout("  Secret Key: %v\n", displaySecretKey)
				stdout("  Type      : %v\n", displayUserType)
				stdout("\nConfirm (yes/no)[yes]: ")
				var userConfirm string
				_, _ = fmt.Scanln(&userConfirm)
				if userConfirm != "yes" && len(userConfirm) != 0 {
					stdout("Abort by user.\n")
					return
				}
			}

			var param = proto.UserCreateParam{
				ID:        userID,
				Password:  password,
				AccessKey: accessKey,
				SecretKey: secretKey,
				Type:      userType,
			}
			var userInfo *proto.UserInfo
			if userInfo, err = client.UserAPI().CreateUser(&param); err != nil {
				errout("Create user failed: %v\n", err)
				os.Exit(1)
			}

			// display operation result
			stdout("Create user success:\n")
			printUserInfo(userInfo)
			return
		},
	}
	cmd.Flags().StringVar(&optPassword, "password", "", "Specify user password")
	cmd.Flags().StringVar(&optAccessKey, "access-key", "", "Specify user access key for object storage interface authentication")
	cmd.Flags().StringVar(&optSecretKey, "secret-key", "", "Specify user secret key for object storage interface authentication")
	cmd.Flags().StringVar(&optUserType, "user-type", "normal", "Specify user type [normal | admin]")
	cmd.Flags().BoolVarP(&optYes, "yes", "y", false, "Answer yes for all questions")
	return cmd
}

const (
	cmdUserUpdateUse   = "update [USER ID]"
	cmdUserUpdateShort = "Update information about specified user"
)

func newUserUpdateCmd(client *master.MasterClient) *cobra.Command {
	var optAccessKey string
	var optSecretKey string
	var optUserType string
	var optYes bool
	var cmd = &cobra.Command{
		Use:   cmdUserUpdateUse,
		Short: cmdUserUpdateShort,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var userID = args[0]
			var accessKey = optAccessKey
			var secretKey = optSecretKey
			var userType proto.UserType
			if optUserType != "" {
				userType = proto.UserTypeFromString(optUserType)
				if !userType.Valid() {
					errout("Invalid user type.\n")
					os.Exit(1)
				}
			}

			if !optYes {
				var displayAccessKey = "[no change]"
				if optAccessKey != "" {
					displayAccessKey = optAccessKey
				}
				var displaySecretKey = "[no change]"
				if optSecretKey != "" {
					displaySecretKey = optSecretKey
				}
				var displayUserType = "[no change]"
				if optUserType != "" {
					displayUserType = optUserType
				}
				fmt.Printf("Update ChubaoFS cluster user\n")
				stdout("  User ID   : %v\n", userID)
				stdout("  Access Key: %v\n", displayAccessKey)
				stdout("  Secret Key: %v\n", displaySecretKey)
				stdout("  Type      : %v\n", displayUserType)
				stdout("\nConfirm (yes/no)[yes]: ")
				var userConfirm string
				_, _ = fmt.Scanln(&userConfirm)
				if userConfirm != "yes" && len(userConfirm) != 0 {
					stdout("Abort by user.\n")
					os.Exit(1)
					return
				}
			}
			if accessKey == "" && secretKey == "" && optUserType == "" {
				stdout("No update.\n")
				os.Exit(1)
				return
			}
			var param = proto.UserUpdateParam{
				UserID:    userID,
				AccessKey: accessKey,
				SecretKey: secretKey,
				Type:      userType,
			}
			var userInfo *proto.UserInfo
			if userInfo, err = client.UserAPI().UpdateUser(&param); err != nil {
				errout("Update user failed: %v\n", err)
				os.Exit(1)
			}

			stdout("Update user success:\n")
			printUserInfo(userInfo)
			return
		},
	}
	cmd.Flags().StringVar(&optAccessKey, "access-key", "", "Update user access key")
	cmd.Flags().StringVar(&optSecretKey, "secret-key", "", "Update user secret key")
	cmd.Flags().StringVar(&optUserType, "user-type", "", "Update user type [normal | admin]")
	cmd.Flags().BoolVarP(&optYes, "yes", "y", false, "Answer yes for all questions")
	return cmd
}

const (
	cmdUserDeleteUse   = "delete [USER ID]"
	cmdUserDeleteShort = "Delete specified user"
)

func newUserDeleteCmd(client *master.MasterClient) *cobra.Command {
	var optYes bool
	//var optForce bool
	var cmd = &cobra.Command{
		Use:   cmdUserDeleteUse,
		Short: cmdUserDeleteShort,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var userID = args[0]

			if !optYes {
				stdout("Delete user [%v] (yes/no)[no]:", userID)
				var userConfirm string
				_, _ = fmt.Scanln(&userConfirm)
				if userConfirm != "yes" {
					stdout("Abort by user.\n")
					os.Exit(1)
					return
				}
			}

			if err = client.UserAPI().DeleteUser(userID); err != nil {
				errout("Delete user failed:\n%v\n", err)
				os.Exit(1)
			}
			stdout("Delete user success.\n")
			return
		},
	}
	cmd.Flags().BoolVarP(&optYes, "yes", "y", false, "Answer yes for all questions")
	//cmd.Flags().BoolVarP(&optForce, "force", "f", false, "Force to delete user")
	return cmd
}

const (
	cmdUserInfoUse   = "info [USER ID]"
	cmdUserInfoShort = "Show detail information about specified user"
)

func newUserInfoCmd(client *master.MasterClient) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   cmdUserInfoUse,
		Short: cmdUserInfoShort,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var userID = args[0]
			var userInfo *proto.UserInfo
			if userInfo, err = client.UserAPI().GetUserInfo(userID); err != nil {
				errout("Get user info failed: %v\n", err)
				os.Exit(1)
			}
			printUserInfo(userInfo)
		},
	}

	return cmd
}

const (
	cmdUserPermUse   = "perm [USER ID] [VOLUME] [PERM (READONLY,RO,READWRITE,RW,NONE)]"
	cmdUserPermShort = "Setup volume permission for a user"
)

func newUserPermCmd(client *master.MasterClient) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   cmdUserPermUse,
		Short: cmdUserPermShort,
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			var userID = args[0]
			var volume = args[1]
			var perm proto.Permission
			switch strings.ToLower(args[2]) {
			case "ro", "readonly":
				perm = proto.BuiltinPermissionReadOnly
			case "rw", "readwrite":
				perm = proto.BuiltinPermissionWritable
			case "none":
				perm = proto.NonePermission
			default:
				stdout("Permission must be on of ro, rw, none")
				return
			}
			stdout("Setup volume permission\n")
			stdout("  User ID   : %v\n", userID)
			stdout("  Volume    : %v\n", volume)
			stdout("  Permission: %v\n", perm.ReadableString())

			// ask user for confirm
			stdout("\nConfirm (yes/no)[yes]: ")
			var userConfirm string
			_, _ = fmt.Scanln(&userConfirm)
			if userConfirm != "yes" && len(userConfirm) != 0 {
				stdout("Abort by user.\n")
				return
			}
			var err error
			defer func() {
				if err != nil {
					errout("Setup permission failed:\n%v\n", err)
					os.Exit(1)
				}
			}()
			var userInfo *proto.UserInfo
			if userInfo, err = client.UserAPI().GetUserInfo(userID); err != nil {
				return
			}
			if _, err = client.AdminAPI().GetVolumeSimpleInfo(volume); err != nil {
				return
			}
			if perm.IsNone() {
				param := proto.NewUserPermRemoveParam(userID, volume)
				userInfo, err = client.UserAPI().RemovePolicy(param)
			} else {
				param := proto.NewUserPermUpdateParam(userID, volume)
				param.SetPolicy(perm.String())
				userInfo, err = client.UserAPI().UpdatePolicy(param)
			}
			if err != nil {
				return
			}
			printUserInfo(userInfo)
		},
	}
	return cmd
}

const (
	cmdUserListUse   = "list"
	cmdUserListShort = "List cluster users"
)

func newUserListCmd(client *master.MasterClient) *cobra.Command {
	var optKeyword string
	var cmd = &cobra.Command{
		Use:     cmdUserListUse,
		Short:   cmdUserListShort,
		Aliases: []string{"ls"},
		Run: func(cmd *cobra.Command, args []string) {
			var users []*proto.UserInfo
			var err error
			defer func() {
				if err != nil {
					errout("List cluster user failed: %v\n", err)
					os.Exit(1)
				}
			}()
			if users, err = client.UserAPI().ListUsers(optKeyword); err != nil {
				return
			}
			stdout("%v\n", userInfoTableHeader)
			for _, user := range users {
				stdout("%v\n", formatUserInfoTableRow(user))
			}
		},
	}
	cmd.Flags().StringVar(&optKeyword, "keyword", "", "Specify keyword of user name to filter")
	return cmd
}

func printUserInfo(userInfo *proto.UserInfo) {
	stdout("[Summary]\n")
	stdout("  User ID    : %v\n", userInfo.UserID)
	stdout("  Access Key : %v\n", userInfo.AccessKey)
	stdout("  Secret Key : %v\n", userInfo.SecretKey)
	stdout("  Type       : %v\n", userInfo.UserType)
	stdout("  Create Time: %v\n", userInfo.CreateTime)
	if userInfo.Policy == nil {
		return
	}
	stdout("[Volumes]\n")
	stdout("%-20v    %-12v\n", "VOLUME", "PERMISSION")
	for _, vol := range userInfo.Policy.OwnVols {
		stdout("%-20v    %-12v\n", vol, "Owner")
	}
	for vol, perms := range userInfo.Policy.AuthorizedVols {
		stdout("%-20v    %-12v\n", vol, strings.Join(perms, ","))
	}
}
