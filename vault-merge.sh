#!/bin/bash

# Credit to original script: https://victorkoronen.se/2017/07/07/merging-ansible-vaults-in-git/

set -e

ancestor_version=$1
current_version=$2
other_version=$3
conflict_marker_size=$4
merged_result_pathname=$5

ancestor_tempfile=$(mktemp tmp.XXXXXXXXXX)
current_tempfile=$(mktemp tmp.XXXXXXXXXX)
other_tempfile=$(mktemp tmp.XXXXXXXXXX)

delete_tempfiles() {
    rm -f "$ancestor_tempfile" "$current_tempfile" "$other_tempfile"
}
trap delete_tempfiles EXIT

if [[ ! -z "${VAULT_PASS}" ]]; then
    echo "$VAULT_PASS" > ./_vault_pass
	ansible-vault decrypt --output  "$ancestor_tempfile" "$ancestor_version" --vault-password-file ./_vault_pass
	ansible-vault decrypt --output  "$current_tempfile" "$current_version"   --vault-password-file ./_vault_pass
	ansible-vault decrypt --output  "$other_tempfile" "$other_version"       --vault-password-file ./_vault_pass
    rm _vault_pass
else
	ansible-vault decrypt --output  "$ancestor_tempfile" "$ancestor_version" --ask-vault-pass
	ansible-vault decrypt --output  "$current_tempfile" "$current_version"   --ask-vault-pass
	ansible-vault decrypt --output  "$other_tempfile" "$other_version"       --ask-vault-pass
fi

git merge-file "$current_tempfile" "$ancestor_tempfile" "$other_tempfile"

if [[ ! -z "${VAULT_PASS}" ]]; then
    echo "$VAULT_PASS" > ./_vault_pass
	ansible-vault encrypt --output "$current_version" "$current_tempfile" --vault-password-file ./_vault_pass
    rm _vault_pass
else
	ansible-vault encrypt --output "$current_version" "$current_tempfile" --ask-vault-pass
fi