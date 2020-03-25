/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

package localpath

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"k8s.io/client-go/util/homedir"
)

func TestReplaceWinDriveLetterToVolumeName(t *testing.T) {
	path, err := ioutil.TempDir("", "repwindl2vn")
	if err != nil {
		t.Fatalf("Error make tmp directory: %v", err)
	}
	defer os.RemoveAll(path)

	if runtime.GOOS != "windows" {
		// Replace to fake func.
		getWindowsVolumeName = func(d string) (string, error) {
			return `/`, nil
		}
		// Add dummy Windows drive letter.
		path = `C:` + path
	}

	if _, err := replaceWinDriveLetterToVolumeName(path); err != nil {
		t.Errorf("Error replace a Windows drive letter to a volume name: %v", err)
	}
}

func TestHasWindowsDriveLetter(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{
		{`C:\Users\Foo\.minikube`, true},
		{`D:\minikube\.minikube`, true},
		{`C\Foo\Bar\.minikube`, false},
		{`/home/foo/.minikube`, false},
	}

	for _, tc := range cases {
		if hasWindowsDriveLetter(tc.path) != tc.want {
			t.Errorf("%s have a Windows drive letter: %t", tc.path, tc.want)
		}
	}
}

func TestConfigFile(t *testing.T) {
	configFile := ConfigFile()
	if !strings.Contains(configFile, "config.json") {
		t.Errorf("ConfigFile returned path without 'config.json': %s", configFile)
	}
}

func TestMiniPath(t *testing.T) {
	var testCases = []struct {
		env, basePath string
	}{
		{"/tmp/.minikube", "/tmp/"},
		{"/tmp/", "/tmp"},
		{"", homedir.HomeDir()},
	}
	for _, tc := range testCases {
		originalEnv := os.Getenv(MinikubeHome)
		defer func() { // revert to pre-test env var
			err := os.Setenv(MinikubeHome, originalEnv)
			if err != nil {
				t.Fatalf("Error reverting env %s to its original value (%s) var after test ", MinikubeHome, originalEnv)
			}
		}()
		t.Run(fmt.Sprintf("%s", tc.env), func(t *testing.T) {
			expectedPath := filepath.Join(tc.basePath, ".minikube")
			os.Setenv(MinikubeHome, tc.env)
			path := MiniPath()
			if path != expectedPath {
				t.Errorf("MiniPath expected to return '%s', but got '%s'", expectedPath, path)
			}
		})
	}
}
