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
	"os"
	"os/exec"

)

func NewCmdExec() cli.Command {
	c:=cli.Command{
		Name:  "exec",
		Usage: "进入容器方法。grctl exec NAMESPACE POD_NAME COMMAND ",
		Action: func(c *cli.Context) error {
			Common(c)
			return execContainer(c)
		},
	}
	return c
}

// grctl exec NAMESPACE POD_ID COMMAND
func execContainer(c *cli.Context) error {
	//podID := c.Args().Get(1)
	args := c.Args().Tail()
	tenantID:=c.Args().First()


	//podID := c.Args().First()
	//args := c.Args().Tail()
	//tenantID, err := clients.FindNamespaceByPod(podID)

	//clients.K8SClient.Core().Namespaces().Get("",metav1.GetOptions{}).


	kubeCtrl, err := exec.LookPath("kubectl")
	if err != nil {
		logrus.Error("Don't fnd the kubectl")
		return err
	}
	if len(args) == 0 {
		args = []string{"bash"}
	}
	//logrus.Infof("using namespace %s,podid %s",tenantID)
	defaultArgs := []string{kubeCtrl, "exec", "-it", "--namespace=" + tenantID}
	args = append(defaultArgs, args...)
	//logrus.Info(args)
	cmd := exec.Cmd{
		Env:    os.Environ(),
		Path:   kubeCtrl,
		Args:   args,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := cmd.Run(); err != nil {
		logrus.Error("Exec error.", err.Error())
		return err
	}
	return nil
}

