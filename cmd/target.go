// Copyright © 2017 Jennifer Hamon <jhamon@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"github.com/jhamon/uaa-cli/config"
	"github.com/jhamon/uaa-cli/uaa"
	"github.com/spf13/cobra"
	"os"
)

var skipSSLValidation bool

func printTarget(target uaa.Target, status string, version string) {
	fmt.Println("Target: " + target.BaseUrl)
	fmt.Println("Status: " + status)
	fmt.Println("UAA Version: " + version)
	fmt.Printf("SkipSSLValidation: %v\n", target.SkipSSLValidation)
}

func TraceRetryMsg(c uaa.Config) {
	if !c.Trace {
		fmt.Println("Retry with --trace for more information.")
	}
}

func showTarget() {
	c := GetSavedConfig()
	target := c.GetActiveTarget()

	if target.BaseUrl == "" {
		printTarget(target, "", "")
		return
	}

	info, err := uaa.Info(GetHttpClient(), c)
	if err != nil {
		printTarget(target, "ERROR", "unknown")
		os.Exit(1)
	}

	printTarget(target, "OK", info.App.Version)
}

func updateTarget(newTarget string) {
	c := GetSavedConfig()

	target := uaa.Target{
		SkipSSLValidation: skipSSLValidation,
		BaseUrl:           newTarget,
	}

	c.AddTarget(target)
	_, err := uaa.Info(GetHttpClientWithConfig(c), c)
	if err != nil {
		fmt.Println(fmt.Sprintf("The target %s could not be set.", newTarget))
		TraceRetryMsg(c)
		os.Exit(1)
	}

	config.WriteConfig(c)
	fmt.Println("Target set to " + newTarget)
}

var targetCmd = &cobra.Command{
	Use:   "target UAA_URL",
	Short: "Set the url of the UAA you'd like to target",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			showTarget()
		} else {
			updateTarget(args[0])
		}
	},
}

func init() {
	RootCmd.AddCommand(targetCmd)
	targetCmd.Flags().BoolVarP(&skipSSLValidation, "skip-ssl-validation", "k", false, "Disable security validation on requests to this target")
}
