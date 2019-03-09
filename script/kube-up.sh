#!/bin/bash
# minikube start --feature-gates=KubeletPluginsWatcher=true,BlockVolume=true,CSIBlockVolume=true,CSIPersistentVolume=true,CSIDriverRegistry=true,CSINodeInfo=true
minikube start --feature-gates=KubeletPluginsWatcher=true,BlockVolume=true,CSIDriverRegistry=true

# replace /var sys links with abs links
VAR_CMD="sudo find /var -type l -execdir bash -c 'ln -sfn \"\$(readlink -f \"\$0\")\" \"\$0\"' {} \\;"
minikube ssh "$VAR_CMD"

# replace etc sys links with abs links
ETC_CMD="sudo find /etc -type l -execdir bash -c 'ln -sfn \"\$(readlink -f \"\$0\")\" \"\$0\"' {} \\;"
minikube ssh "$ETC_CMD"
