#!/usr/bin/env bash
# @author: Daniel Stolf <daniel.stolf@delphix.com> - 02/2020
# Changelog:
###

[ -f ~/.bashrc ] && source ~/.bashrc
[ -f ~/.bash_profile ] && source ~/.bash_profile

export PATH=$PATH:/usr/sbin
# kubectl path

EXPORTPATH=$(df -kh | grep "${MOUNTPATH}" | awk '{print $1}')

if [ ! "${EXPORTPATH}." == "." ]; then
    echo ${EXPORTPATH} | tr -d '\n'
else
    exit 1
fi