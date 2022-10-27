#!/usr/bin/env bash
# @author: Daniel Stolf <daniel.stolf@delphix.com> - 10/2022
# Changelog:
# 10/2022 - change environment discovery to allow compatibility with microk8s and kubectl aliases
###

[ -f ~/.bashrc ] && source ~/.bashrc
[ -f ~/.bash_profile ] && source ~/.bash_profile

export PATH=$PATH:/usr/sbin
# kubectl path
KUBECTLBIN=$(which kubectl 2> /dev/null | gawk '{ match($0, /(\(|'\`'|'\'')(.*)(\)|'\`'|'\'')/, arr); if(arr[2] != "") print arr[2] }')
if [ ! "${KUBECTLBIN}." == "." ]; then
    for l in $(${KUBECTLBIN} config get-contexts | grep -v CURRENT| awk '{print $2"|"$3"|"$4}'); do
        $KUBECTLBIN config use-context "$CONTEXTNAME" > /dev/null 2>&1
        VERSION=$($KUBECTLBIN version --output=json | sed -ne '/serverVersion/,$ p' | gawk '{ match($0, /(\"gitVersion\": \")(.*)(\")/, arr); if(arr[2] != "") print arr[2] }' )
        K8SMASTER=$($KUBECTLBIN  cluster-info | grep "is running at" | grep master| rev | awk '{print $1}'| rev | sed 's/\x1b\[[0-9;]*m//g')
        #0 - kubectl path | 1- context name | 2- cluster name | 3- username | 4- API endpoint | 5- Server Version
        echo "${KUBECTLBIN}|${l}|${K8SMASTER}|${VERSION}" | tr -d '\n'
    done
else
  exit 0
fi