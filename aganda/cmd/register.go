/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/hixinj/aganda/entity"
	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "注册用户",
	Long: `
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("register called")
		// username, _ := cmd.Flags().GetString("user")
		// fmt.Println("register called by " + username)
		s := &entity.Storage{}
		s.ReadUsers()

		newUser := entity.User{}
		newUser.Name, _ = cmd.Flags().GetString("user")
		newUser.Password, _ = cmd.Flags().GetString("password")
		newUser.Email, _ = cmd.Flags().GetString("email")
		newUser.Phone, _ = cmd.Flags().GetString("phone")

		if s.QueryUser(func(u entity.User) bool {
			return u.Name == newUser.Name
		}) != nil {
			fmt.Println("error: same user name exists.")
			return
		}

		s.UserList = append(s.UserList, newUser)

		s.WriteUsers()
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
	registerCmd.Flags().StringP("user", "u", "Anonymous", "user name")
	registerCmd.Flags().StringP("password", "", "", "password of the account")
	registerCmd.Flags().StringP("email", "e", "", "your email")
	registerCmd.Flags().StringP("phone", "", "", "your phone")
	registerCmd.MarkFlagRequired("user")
	registerCmd.MarkFlagRequired("password")
	registerCmd.MarkFlagRequired("email")
	registerCmd.MarkFlagRequired("phone")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
