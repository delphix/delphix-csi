#!/usr/bin/env bash
# @author: Daniel Stolf <daniel.stolf@delphix.com> - 02/2020
# Changelog:
###

[ -f ~/.bashrc ] && source ~/.bashrc
[ -f ~/.bash_profile ] && source ~/.bash_profile

export PATH=$PATH:/usr/sbin
# kubectl path
KUBECTLBIN=$(which kubectl 2> /dev/null)

if [ ! "${KUBECTLBIN}." == "." ]; then
    for l in $(kubectl config get-contexts | grep -v CURRENT| awk '{print $2"|"$3"|"$4}'); do
        $KUBECTLBIN config use-context "$CONTEXTNAME" > /dev/null 2>&1
        VERSION=$($KUBECTLBIN  version | grep "Server Version" | grep -oP "(?<=GitVersion:\")[^\"]+(?=\")")
        K8SMASTER=$($KUBECTLBIN  cluster-info | grep "is running at" | grep master| rev | awk '{print $1}'| rev | sed 's/\x1b\[[0-9;]*m//g')
        #0 - kubectl path | 1- context name | 2- cluster name | 3- username | 4- API endpoint | 5- Server Version
        echo "${KUBECTLBIN}|${l}|${K8SMASTER}|${VERSION}" | tr -d '\n' 
    done
else
  exit 0
fi