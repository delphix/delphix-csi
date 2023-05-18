#! /bin/sh
# @author: Gaurav Manhas <gaurav.manhas@delphix.com> - 05/2023
###
shopt -s expand_aliases
[ -f ~/.bashrc ] && source ~/.bashrc
[ -f ~/.bash_profile ] && source ~/.bash_profile

export PATH=$PATH:/usr/sbin
# kubectl path

sudo umount "${MOUNTPATH}"

if [ -d "${MOUNTPATH} ]; then
    echo "Mount path removed Successfuly"
else
    exit 1
fi
