// Copyright 2016 The rkt Authors
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

package hvlkvm

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rkt/rkt/networking/netinfo"
	"github.com/rkt/rkt/stage1/init/kvm/hypervisor"
)

// StartCmd takes path to stage1, UUID of the pod, path to kernel, network
// describers, memory in megabytes and quantity of cpus and prepares command
// line to run LKVM process
func StartCmd(wdPath, uuid, kernelPath string, nets []*netinfo.NetInfo, cpu, mem int64, debug bool) []string {
	machineID := strings.Replace(uuid, "-", "", -1)
	driverConfiguration := hypervisor.KvmHypervisor{
		Bin: "./lkvm",
		KernelParams: []string{
			"systemd.default_standard_error=journal+console",
			"systemd.default_standard_output=journal+console",
			"systemd.machine_id=" + machineID,
		},
	}

	driverConfiguration.InitKernelParams(debug)

	startCmd := []string{
		filepath.Join(wdPath, driverConfiguration.Bin),
		"run",
		"--name", "rkt-" + uuid,
		"--no-dhcp",
		"--cpu", strconv.Itoa(int(cpu)),
		"--mem", strconv.Itoa(int(mem)),
		"--console=virtio",
		"--kernel", kernelPath,
		"--disk", "stage1/rootfs", // relative to run/pods/uuid dir this is a place where systemd resides
		"--params", strings.Join(driverConfiguration.KernelParams, " "),
	}
	return append(startCmd, kvmNetArgs(nets)...)
}

// kvmNetArgs returns additional arguments that need to be passed
// to lkvm tool to configure networks properly. Logic is based on
// network configuration extracted from Networking struct
// and essentially from activeNets that expose NetDescriber behavior
func kvmNetArgs(nets []*netinfo.NetInfo) []string {
	var lkvmArgs []string

	for _, net := range nets {
		ipConf := net.FirstIPConfig()
		if ipConf == nil {
			continue
		}

		lkvmArgs = append(lkvmArgs, "--network")
		lkvmArgs = append(
			lkvmArgs,
			fmt.Sprintf("mode=tap,tapif=%s,host_ip=%s,guest_ip=%s", net.IfName, ipConf.Gateway, ipConf.Address.IP.String()),
		)
	}

	return lkvmArgs
}
