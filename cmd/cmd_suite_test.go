package cmd_test

import (
	"fmt"
	"os"

	"github.com/koding/vagrantutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/spf13/afero"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}

var vagrant *vagrantutil.Vagrant
var projectPath = fmt.Sprintf("%s/src/github.com/lgug2z/sgrab", os.Getenv("GOPATH"))
var err error

var _ = BeforeSuite(func() {
	vagrantFileContents := `# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.ssh.insert_key = false
  config.vm.box = "ubuntu/trusty64"
  config.vm.hostname = "vagrant"

  config.vm.provider "virtualbox" do |vb|
    # Use VBoxManage to customize the VM. For example to change memory:
    vb.customize ["modifyvm", :id, "--memory", "2048", "--cpus", "2"]
  end

  config.vm.provision "shell",
    inline: "touch /westworld-s01e01.mkv"
end
`
	vagrant, err = vagrantutil.NewVagrant(projectPath)
	Expect(err).ToNot(HaveOccurred())

	Expect(vagrant.Create(vagrantFileContents)).To(Succeed())

	up, err := vagrant.Up()
	for _ = range up {
	}

	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	destroy, err := vagrant.Destroy()
	for _ = range destroy {
	}

	Expect(err).ToNot(HaveOccurred())

	fs := afero.NewOsFs()
	Expect(fs.Remove(fmt.Sprintf("%s/%s", projectPath, "Vagrantfile"))).To(Succeed())
	Expect(fs.Remove(fmt.Sprintf("%s/%s", projectPath, ".vagrant"))).To(Succeed())
})
