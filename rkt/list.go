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

//+build linux

package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/coreos/rkt/Godeps/_workspace/src/github.com/appc/spec/schema"
	"github.com/coreos/rkt/Godeps/_workspace/src/github.com/appc/spec/schema/lastditch"
	"github.com/coreos/rkt/Godeps/_workspace/src/github.com/appc/spec/schema/types"
	"github.com/coreos/rkt/Godeps/_workspace/src/github.com/dustin/go-humanize"
	"github.com/coreos/rkt/Godeps/_workspace/src/github.com/spf13/cobra"
	common "github.com/coreos/rkt/common"
	"github.com/coreos/rkt/networking/netinfo"
)

var (
	cmdList = &cobra.Command{
		Use:   "list",
		Short: "List pods",
		Run:   runWrapper(runList),
	}
	flagNoLegend   bool
	flagFullOutput bool
)

func init() {
	cmdRkt.AddCommand(cmdList)
	cmdList.Flags().BoolVar(&flagNoLegend, "no-legend", false, "suppress a legend with the list")
	cmdList.Flags().BoolVar(&flagFullOutput, "full", false, "use long output format")
}

func runList(cmd *cobra.Command, args []string) int {
	var errors []error
	tabBuffer := new(bytes.Buffer)
	tabOut := getTabOutWithWriter(tabBuffer)

	if !flagNoLegend {
		if flagFullOutput {
			fmt.Fprintf(tabOut, "UUID\tAPP\tIMAGE NAME\tIMAGE ID\tSTATE\tSINCE\tNETWORKS\n")
		} else {
			fmt.Fprintf(tabOut, "UUID\tAPP\tIMAGE NAME\tSTATE\tSINCE\tNETWORKS\n")
		}
	}

	if err := walkPods(includeMostDirs, func(p *pod) {
		pm := schema.PodManifest{}

		if !p.isPreparing && !p.isAbortedPrepare && !p.isExitedDeleting {
			// TODO(vc): we should really hold a shared lock here to prevent gc of the pod
			pmf, err := p.readFile(common.PodManifestPath(""))
			if err != nil {
				errors = append(errors, newPodListReadError(p, err))
				return
			}

			if err := pm.UnmarshalJSON(pmf); err != nil {
				errors = append(errors, newPodListLoadError(p, err, pmf))
				return
			}

			if len(pm.Apps) == 0 {
				errors = append(errors, newPodListZeroAppsError(p))
				return
			}
		}

		type printedApp struct {
			uuid    string
			appName string
			imgName string
			imgID   string
			state   string
			nets    string
			whenStr string
		}

		var appsToPrint []printedApp
		uuid := p.uuid.String()
		state := p.getState()
		nets := fmtNets(p.nets)

		created, err := p.getRootModTime()
		if err != nil {
			errors = append(errors, fmt.Errorf("unable to get modification time for pod %q root: %v", uuid, err))
		}
		exited, err := p.getExitTime()
		if err != nil {
			errors = append(errors, fmt.Errorf("unable to get exit time for pod %q: %v", uuid, err))
		}
		var whenStr string
		if flagFullOutput {
			if p.isExited || p.isExitedGarbage {
				if !exited.IsZero() {
					whenStr = exited.Format(defaultTimeLayout)
				}
			} else if !created.IsZero() {
				whenStr = created.Format(defaultTimeLayout)
			}
		} else {
			if p.isExited || p.isExitedGarbage {
				if !exited.IsZero() {
					whenStr = humanize.Time(exited)
				}
			} else if !created.IsZero() {
				whenStr = humanize.Time(created)
			}
		}

		if !flagFullOutput {
			uuid = uuid[:8]
		}
		for _, app := range pm.Apps {
			imageName, err := getImageName(p, app.Name)
			if err != nil {
				errors = append(errors, newPodListLoadImageManifestError(p, err))
				imageName = "--"
			}

			var imageID string
			if flagFullOutput {
				imageID = app.Image.ID.String()[:19]
			}

			appsToPrint = append(appsToPrint, printedApp{
				uuid:    uuid,
				appName: app.Name.String(),
				imgName: imageName,
				imgID:   imageID,
				state:   state,
				nets:    nets,
				whenStr: whenStr,
			})
			// clear those variables so they won't be
			// printed for another apps in the pod as they
			// are actually describing a pod, not an app
			uuid = ""
			state = ""
			nets = ""
			whenStr = ""
		}
		// if we reached that point, then it means that the
		// pod and all its apps are valid, so they can be
		// printed
		for _, app := range appsToPrint {
			if flagFullOutput {
				fmt.Fprintf(tabOut, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", app.uuid, app.appName, app.imgName, app.imgID, app.state, app.whenStr, app.nets)
			} else {
				fmt.Fprintf(tabOut, "%s\t%s\t%s\t%s\t%s\t%s\n", app.uuid, app.appName, app.imgName, app.state, app.whenStr, app.nets)
			}
		}

	}); err != nil {
		stderr("Failed to get pod handles: %v", err)
		return 1
	}

	if len(errors) > 0 {
		sep := "----------------------------------------"
		stderr("%d error(s) encountered when listing pods:", len(errors))
		stderr("%s", sep)
		for _, err := range errors {
			stderr("%s", err.Error())
			stderr("%s", sep)
		}
		stderr("Misc:")
		stderr("  rkt's appc version: %s", schema.AppContainerVersion)
		stderr("%s", sep)
		// make a visible break between errors and the listing
		stderr("")
	}

	tabOut.Flush()
	stdout("%s", tabBuffer.String())
	return 0
}

func newPodListReadError(p *pod, err error) error {
	lines := []string{
		fmt.Sprintf("Unable to read pod %s manifest:", p.uuid.String()),
		fmt.Sprintf("  %v", err),
	}
	return fmt.Errorf("%s", strings.Join(lines, "\n"))
}

func newPodListLoadError(p *pod, err error, pmj []byte) error {
	lines := []string{
		fmt.Sprintf("Unable to load pod %s manifest, because it is invalid:", p.uuid.String()),
		fmt.Sprintf("  %v", err),
	}
	pm := lastditch.PodManifest{}
	if err := pm.UnmarshalJSON(pmj); err != nil {
		lines = append(lines, "  Also, failed to get any information about invalid pod manifest:")
		lines = append(lines, fmt.Sprintf("    %v", err))
	} else {
		if len(pm.Apps) > 0 {
			lines = append(lines, "Objects related to this error:")
			for _, app := range pm.Apps {
				lines = append(lines, fmt.Sprintf("  %s", appLine(app)))
			}
		} else {
			lines = append(lines, "No other objects related to this error")
		}
	}
	return fmt.Errorf("%s", strings.Join(lines, "\n"))
}

func newPodListZeroAppsError(p *pod) error {
	return fmt.Errorf("Pod %s contains zero apps", p.uuid.String())
}

func newPodListLoadImageManifestError(p *pod, err error) error {
	return fmt.Errorf("pod %s ImageManifest could not be loaded: %v", p.uuid.String(), err)
}

func appLine(app lastditch.RuntimeApp) string {
	return fmt.Sprintf("App: %q from image %q (%s)",
		app.Name, app.Image.Name, app.Image.ID)
}

func fmtNets(nis []netinfo.NetInfo) string {
	var parts []string
	for _, ni := range nis {
		// there will be IPv6 support soon so distinguish between v4 and v6
		parts = append(parts, fmt.Sprintf("%v:ip4=%v", ni.NetName, ni.IP))
	}
	return strings.Join(parts, ", ")
}

func getImageName(p *pod, appName types.ACName) (string, error) {
	aim, err := p.getAppImageManifest(appName)
	if err != nil {
		return "", fmt.Errorf("problem retrieving ImageManifests from pod: %v", err)
	}

	imageName := aim.Name.String()
	if version, ok := aim.Labels.Get("version"); ok {
		imageName = fmt.Sprintf("%s:%s", imageName, version)
	}

	return imageName, nil
}
