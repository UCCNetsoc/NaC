#!/bin/bash
editor=$EDITOR
export EDITOR=less
ansible-vault edit vars/secrets.yml
export EDITOR=$editor
