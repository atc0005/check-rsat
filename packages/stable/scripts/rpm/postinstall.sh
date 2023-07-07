#!/bin/bash

# Copyright 2023 Adam Chalkley
#
# https://github.com/atc0005/check-rsat
#
# Licensed under the MIT License. See LICENSE file in the project root for
# full license information.

project_org="atc0005"
project_shortname="check-rsat"

project_fq_name="${project_org}/${project_shortname}"
project_url_base="https://github.com/${project_org}"
project_repo="${project_url_base}/${project_shortname}"
project_releases="${project_repo}/releases"
project_issues="${project_repo}/issues"
project_discussions="${project_repo}/discussions"

plugin_name_suffix=""
plugin_path="/usr/lib64/nagios/plugins"

#
# Set required SELinux context to allow plugin use when SELinux is enabled.
#

# Make sure we can locate the selinuxenabled binary.
if [ -x "$(command -v selinuxenabled)" ]; then
    selinuxenabled

    if [ $? -ne 0 ]; then
        echo -e "\n[--] SELinux is not enabled, skipping application of contexts."
    else
        # SELinux is enabled. Set context.
        echo -e "\nApplying SELinux contexts on plugins ..."

        for plugin_name in \
            check_rsat_sync_plans

        do

            echo -e "\nApplying SELinux contexts on ${plugin_path}/${plugin_name}${plugin_name_suffix}"
            restorecon -v ${plugin_path}/${plugin_name}

            if [ $? -eq 0 ]; then
                echo -e "[OK] Successfully applied SELinux contexts on ${plugin_path}/${plugin_name}${plugin_name_suffix}"
            else
                echo -e "[!!] Failed to set SELinux contexts on ${plugin_path}/${plugin_name}${plugin_name_suffix}"
            fi

        done
    fi

else
    echo "[!!] Error: Failed to locate selinuxenabled command." >&2
fi

echo
echo "Thank you for installing packages provided by the ${project_fq_name} project!"
echo
echo "Project resources:"
echo
echo "- Obtain latest release: ${project_releases}"
echo "- View/Ask questions: ${project_discussions}"
echo "- View/Open issues: ${project_issues}"
echo