/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"net"
	"strconv"

	"github.com/ernestio/definition-mapper/libmapper/providers/vcloud/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/vcloud/definition"
	binaryprefix "github.com/r3labs/binary-prefix"
)

// MapInstances : Maps the instances for the input payload on a ernest internal format
func MapInstances(d *definition.Definition) []*components.Instance {
	var instances []*components.Instance

	for _, instance := range d.Instances {
		var commands []string

		ip := net.ParseIP(instance.StartIP).To4()
		memory, _ := binaryprefix.GetMB(instance.Memory)

		for _, prov := range instance.Provisioner {
			if len(prov.Shell) > 0 {
				commands = prov.Shell
			}
		}

		for i := 0; i < instance.Count; i++ {
			var disks []components.InstanceDisk

			if instance.RootDisk != "" {
				size, _ := binaryprefix.GetMB(instance.RootDisk)
				disks = append(disks, components.InstanceDisk{
					ID:   0,
					Size: size,
				})
			}

			disks = append(disks, MapInstanceDisks(instance.Disks)...)

			newInstance := &components.Instance{
				Name:          instance.Name + "-" + strconv.Itoa(i+1),
				Hostname:      instance.Name + "-" + strconv.Itoa(i+1),
				Catalog:       instance.Catalog(),
				Image:         instance.Template(),
				Cpus:          instance.Cpus,
				Memory:        memory,
				Disks:         disks,
				Network:       instance.Network,
				IP:            ip.String(),
				ShellCommands: commands,
			}

			newInstance.SetDefaultVariables()

			instances = append(instances, newInstance)

			// Increment IP address
			ip[3]++
		}
	}
	return instances
}

// MapInstanceDisks : Maps the instances disks
func MapInstanceDisks(d []string) []components.InstanceDisk {
	var disks []components.InstanceDisk

	for x, disk := range d {
		size, _ := binaryprefix.GetMB(disk)
		disks = append(disks, components.InstanceDisk{
			ID:   (x + 1),
			Size: size,
		})
	}

	return disks
}
