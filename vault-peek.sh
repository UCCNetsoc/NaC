#!/bin/bash
editor=$EDITOR
export EDITOR=cat
ansible-vault edit vars/secrets.yml
export EDITOR=$editor
