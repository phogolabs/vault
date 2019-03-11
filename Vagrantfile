# -*- mode: ruby -*-
# vi: set ft=ruby :

def default_s(key, default)
  ENV[key] && ! ENV[key].empty? ? ENV[key] : default
end

def default_i(key, default)
  default_s(key, default).to_i
end

MEMORY = default_i('MEMORY', 2048)

Vagrant.configure("2") do |config|
  config.vm.box = "ol7-latest"
  config.vm.box_url = "https://yum.oracle.com/boxes/oraclelinux/latest/ol7-latest.box"

  config.vm.provider "virtualbox" do |vb|
    vb.memory = MEMORY
  end

  config.vm.define "master", primary: true do |master|
    master.vm.hostname = "master.vagrant.vm"
    master.vm.network "private_network", ip: "192.168.99.100"

    config.vm.provision "shell", 
      path: "script/kubernetes.sh",
      args: [
        "192.168.99.100",
        "192.168.0.0/16",
      ]
  end

  config.vm.define "worker1" do |worker|
    worker.vm.hostname = "worker1.vagrant.vm"
    worker.vm.network "private_network", ip: "192.168.99.101"

    config.vm.provision "shell", 
      path: "script/kubernetes.sh",
      args: ["192.168.99.101"]
  end

  config.vm.define "worker2" do |worker|
    worker.vm.hostname = "worker2.vagrant.vm"
    worker.vm.network "private_network", ip: "192.168.99.102"

    config.vm.provision "shell", 
      path: "script/kubernetes.sh",
      args: ["192.168.99.102"]
  end
end

