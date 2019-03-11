#!/bin/bash

FEATURE_GATES="KubeletPluginsWatcher=true,BlockVolume=true,CSIDriverRegistry=true,MountPropagation=true"

setenforce 0
sed -i 's/SELINUX=enforcing/SELINUX=disabled/g' /etc/sysconfig/selinux

swapoff -a
sed -i '/swap/s/^/#/g' /etc/fstab

modprobe br_netfilter
cat >> /etc/sysctl.conf <<EOF
net.bridge.bridge-nf-call-ip6tables=1
net.bridge.bridge-nf-call-iptables=1
EOF
sysctl -p

yum install -y yum-utils
yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
yum install -y docker-ce device-mapper-persistent-data lvm2
systemctl enable --now docker

cat > /etc/yum.repos.d/kubernetes.repo <<EOF
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg
       https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
EOF

setup_kubectl() {
 mkdir -p /home/vagrant/.kube
 cp -i /etc/kubernetes/admin.conf /home/vagrant/.kube/config
 chown -R vagrant:vagrant /home/vagrant/.kube/config
}

yum install tc
yum install -y kubelet kubeadm kubectl

if [ -z "$2" ]; then
  setup_kubectl

  /vagrant/script/kubernetes-join.sh

  echo KUBELET_NODE_IP_ARGS=\""--node-ip=$1"\" >> /var/lib/kubelet/kubeadm-flags.env
  cp /vagrant/script/10-kubeadm.conf /etc/systemd/system/kubelet.service.d/10-kubeadm.conf

  systemctl daemon-reload
else
  systemctl enable kubelet.service

  kubeadm init --apiserver-advertise-address="$1" --pod-network-cidr="$2" --feature-gates="$FEATURE_GATES"

  kubeadm token create --print-join-command > /vagrant/script/kubernetes-join.sh
  chmod +x /vagrant/script/kubernetes-join.sh

  setup_kubectl

  kubectl apply -f /vagrant/script/kubernetes-flannel.yml
fi
