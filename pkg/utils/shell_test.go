/*
Copyright 2023 cuisongliu@qq.com.

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

package utils

import (
	"testing"
)

func TestExecShellForAny(t *testing.T) {
	if err := ExecShellForAny()([]any{RetryShell("rm -rf /tmp/fork-sealos-repo && gh repo clone labring/sealos1 /tmp/fork-sealos-repo")}); err != nil {
		t.Errorf("ExecShellForAny() error = %v", err)
	}
}

func TestRunCommandWithOutput(t *testing.T) {
	out, err := RunCommandWithOutput("cd /tmp/fork-sealos-repo && gh run list --workflow Release --json name,number,status,url,event,conclusion -q '.[0]'", false)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", out)
}
