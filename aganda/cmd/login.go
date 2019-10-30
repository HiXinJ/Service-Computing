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
	"os"

	"github.com/hixinj/aganda/entity"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "log in so that you can use aganda to manage meetings",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("login called")
		s := &entity.Storage{}
		userName, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		if err := s.UserLogin(userName, password); err == nil {
			fmt.Println("log in succeeded!")
		} else {
			fmt.Fprintf(os.Stderr, err.Error())
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringP("user", "u", "", "user name")
	loginCmd.MarkFlagRequired("user")
	loginCmd.Flags().StringP("password", "p", "", "your password")
	loginCmd.MarkFlagRequired("password")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
