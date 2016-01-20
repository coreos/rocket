// Copyright 2015 The rkt Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/coreos/rkt/tests/testutils"
)

const delta = 3 * time.Second

// compareTime checks if a and b are roughly equal (1s precision)
func compareTime(a time.Time, b time.Time) bool {
	diff := a.Sub(b)
	if diff < 0 {
		diff = -diff
	}
	return diff < time.Second
}

func TestRktList(t *testing.T) {
	const imgName = "rkt-list-test"

	image := patchTestACI(fmt.Sprintf("%s.aci", imgName), fmt.Sprintf("--name=%s", imgName))
	defer os.Remove(image)

	imageHash := getHashOrPanic(image)
	imgID := ImageID{image, imageHash}

	ctx := testutils.NewRktRunCtx()
	defer ctx.Cleanup()

	// Prepare image
	cmd := fmt.Sprintf("%s --insecure-options=image prepare %s", ctx.Cmd(), imgID.path)
	podUuid := runRktAndGetUUID(t, cmd)

	// Get hash
	imageID := fmt.Sprintf("sha512-%s", imgID.hash[:12])

	tmpDir := createTempDirOrPanic(imgName)
	defer os.RemoveAll(tmpDir)

	// Define tests
	tests := []struct {
		cmd           string
		shouldSucceed bool
		expect        string
	}{
		// Test that pod UUID is in output
		{
			"list --full",
			true,
			podUuid,
		},
		// Test that image name is in output
		{
			"list",
			true,
			imgName,
		},
		// Test that imageID is in output
		{
			"list --full",
			true,
			imageID,
		},
		// Remove the image
		{
			fmt.Sprintf("image rm %s", imageID),
			true,
			"successfully removed",
		},
		// Name should still show up in rkt list
		{
			"list",
			true,
			imgName,
		},
		// Test that imageID is still in output
		{
			"list --full",
			true,
			imageID,
		},
	}

	// Run tests
	for i, tt := range tests {
		runCmd := fmt.Sprintf("%s %s", ctx.Cmd(), tt.cmd)
		t.Logf("Running test #%d, %s", i, runCmd)
		runRktAndCheckOutput(t, runCmd, tt.expect, !tt.shouldSucceed)
	}
}

func getSinceTime(t *testing.T, ctx *testutils.RktRunCtx, imageID string, state string) time.Time {
	// Run rkt list --full
	rktCmd := fmt.Sprintf("%s list --full", ctx.Cmd())
	child := spawnOrFail(t, rktCmd)
	child.Wait()

	// Get prepared time
	match := fmt.Sprintf(".*%s\t%s\t(.*)\t", imageID, state)
	result, out, err := expectRegexWithOutput(child, match)
	if err != nil {
		t.Fatalf("%q regex not found, Error: %v\nOutput: %v", match, err, out)
	}
	tmStr := strings.TrimSpace(result[1])
	tm, err := time.Parse(defaultTimeLayout, tmStr)
	if err != nil {
		t.Fatalf("Error parsing %s time: %q", state, err)
	}

	return tm
}

func TestRktListWhen(t *testing.T) {
	const imgName = "rkt-list-creation-time-test"

	image := patchTestACI(fmt.Sprintf("%s.aci", imgName), fmt.Sprintf("--exec=/inspect --read-stdin"))
	defer os.Remove(image)

	imageHash := getHashOrPanic(image)
	imgID := ImageId{image, imageHash}

	ctx := testutils.NewRktRunCtx()
	defer ctx.Cleanup()

	// Prepare image
	cmd := fmt.Sprintf("%s --insecure-skip-verify prepare %s", ctx.Cmd(), imgID.path)
	podUuid := runRktAndGetUUID(t, cmd)

	// t0: prepare
	expectPrepare := time.Now()

	// Get hash
	imageID := fmt.Sprintf("sha512-%s", imgID.hash[:12])

	tmpDir := createTempDirOrPanic(imgName)
	defer os.RemoveAll(tmpDir)

	prepared := getSinceTime(t, ctx, imageID, "prepared")
	if !compareTime(expectPrepare, prepared) {
		t.Fatalf("incorrect preparation time returned. Got: %q Expect: %q (1s precision)", prepared, expectPrepare)
	}

	time.Sleep(delta)

	// t1: run
	expectRun := time.Now()

	// Run image
	cmd = fmt.Sprintf("%s run-prepared --interactive %s", ctx.Cmd(), podUuid)
	rktChild := spawnOrFail(t, cmd)

	time.Sleep(delta)

	running := getSinceTime(t, ctx, imageID, "running")
	if !compareTime(expectRun, running) {
		t.Fatalf("incorrect running time returned. Got: %q Expect: %q (1s precision)", running, expectRun)
	}

	say := "Hello"
	if err := rktChild.SendLine(say); err != nil {
		t.Fatalf("Failed to send %q to rkt: %v", say, err)
	}
	rktChild.Wait()

	// t2: exit
	expectExit := expectRun.Add(delta)

	exited := getSinceTime(t, ctx, imageID, "exited")
	if !compareTime(expectExit, exited) {
		t.Fatalf("incorrect exit time returned. Got: %q Expect: %q (1s precision)", exited, expectExit)
	}
}
