// RAINBOND, Application Management Platform
// Copyright (C) 2014-2017 Goodrain Co., Ltd.

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Rainbond,
// one or multiple Commercial Licenses authorized by Goodrain Co., Ltd.
// must be obtained first.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"github.com/urfave/cli"
	"github.com/Sirupsen/logrus"
	"github.com/goodrain/rainbond/pkg/grctl/clients"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/apcera/termtables"
	"fmt"
	"strings"
)


func NewCmdNode() cli.Command {
	c:=cli.Command{
		Name:  "node",
		Usage: "获取节点信息。grctl node",
		Action: func(c *cli.Context) error {
			Common(c)
			return getNode(c)
		},
	}
	return c
}


func NewCmdNodeRes() cli.Command {
	c:=cli.Command{
		Name:  "noderes",
		Usage: "获取计算节点资源信息  grctl noderes",
		Action: func(c *cli.Context) error {
			Common(c)
			return getNodeWithResource(c)
		},
	}
	return c
}

func getNodeWithResource(c *cli.Context) error {
	ns, err :=clients.K8SClient.Core().Nodes().List(metav1.ListOptions{})
	if err != nil {
		logrus.Errorf("获取节点列表失败,details: %s",err.Error())
		return err
	}
	table := termtables.CreateTable()
	table.AddHeaders("NodeName", "Version", "CapCPU(核)", "AllocatableCPU(核)","UsedCPU(核)", "CapMemory(M)","AllocatableMemory(M)","UsedMemory(M)")
	for _,v:=range ns.Items {
		capCPU:=v.Status.Capacity.Cpu().Value()
		capMem:=v.Status.Capacity.Memory().Value()
		allocCPU:=v.Status.Allocatable.Cpu().Value()
		allocMem:=v.Status.Allocatable.Memory().Value()
		table.AddRow(v.Name,v.Status.NodeInfo.KubeletVersion,capCPU,allocCPU,capCPU-allocCPU,capMem/1024/1024,allocMem/1024/1024,capMem/1024/1024-allocMem/1024/1024)
	}
	fmt.Println(table.Render())
	return nil
}

func getNode(c *cli.Context) error {
	ns, err :=clients.K8SClient.Core().Nodes().List(metav1.ListOptions{})
	if err != nil {
		logrus.Errorf("获取节点列表失败,details: %s",err.Error())
		return err
	}
	table := termtables.CreateTable()
	table.AddHeaders("Name", "Status", "Namespace","Unschedulable", "KubeletVersion","Labels")

	for _,v:=range ns.Items{
		cs:=v.Status.Conditions
		status:="unknown"
		for _,cv:=range cs{
			status=string(cv.Status)
			if strings.Contains(status,"rue"){
				status=string(cv.Type)
				break
			}
		}
		m:=v.Labels
		labels:=""
		for k:=range m {
			labels+=k
			labels+=" "
		}
		table.AddRow(v.Name, status,v.Namespace,v.Spec.Unschedulable, v.Status.NodeInfo.KubeletVersion,labels)
	}
	fmt.Println(table.Render())
	return nil
}
