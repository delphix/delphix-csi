#!/usr/bin/env bash
# @author: Daniel Stolf <daniel.stolf@delphix.com> - 02/2020
# Changelog:
###

[ -f ~/.bashrc ] && source ~/.bashrc
[ -f ~/.bash_profile ] && source ~/.bash_profile

export PATH=$PATH:/usr/sbin

clustername=$($KUBECTLBIN  config get-contexts $CONTEXTNAME | grep -v CURRENT| awk '{print $3}')
$KUBECTLBIN config use-context "$CONTEXTNAME" > /dev/null 2>&1
K8SMASTER=$($KUBECTLBIN  cluster-info | grep "is running at" | grep master| rev | awk '{print $1}'| rev | sed 's/\x1b\[[0-9;]*m//g')

if [ "$clustername" != "$CLUSTER"] || [ "$K8SMASTER" != "$APIENDPOINT" ]; then
echo "ERROR: configuration missmatch in kubectl, APIENDPOINT or CLUSTER NAME changed. Run the following commands to validate: "
echo "$KUBECTLBIN  config get-contexts $CONTEXTNAME"
echo "$KUBECTLBIN  cluster-info"
exit 1
fi
#0 - kubectl path | 1- context name | 2- cluster name | 3- username | 4- API endpoint | 5- Server Version
echo "${KUBECTLBIN}|${l}|${K8SMASTER}|${VERSION}" | tr -d '\n' 
