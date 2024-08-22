/*
 * Copyright 2024 KylinSoft  Co., Ltd.
 * KubeMate is licensed under the Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
 * PURPOSE.
 * See the Mulan PSL v2 for more details.
 */

package common

import (
	"fmt"
	"universal_os_upgrader/pkg/template"
)

func GetRearShell(handleType string, nfsServer string, nfsPath string) (string, error) {
	shell := fmt.Sprintf(`#!/bin/bash

REAR_CONF_FILE="/etc/rear/local.conf"
REAR_MKBACKUP_SCRIPT="/tmp/rear-mkbackup.sh"
BACKUP_URL="nfs://%s%s"

setup_environment() {
    sudo yum install -y rear genisoimage syslinux psmisc

    ARCH=$(uname -m)

    if [ "$ARCH" = "x86_64" ]; then
        sudo yum install -y grub2-efi-x64-modules
    elif [ "$ARCH" = "aarch64" ]; then
        sudo yum install -y grub2-efi-aa64-modules
    else
        echo "Unsupported architecture: $ARCH"
        exit 1
    fi

    OS_ID=$(grep '^ID=' /etc/os-release | cut -d '=' -f 2 | tr -d '"')

    if [ "$OS_ID" = "kylin" ]; then
        sudo yum install -y kylin-lsb
    elif [ "$OS_ID" = "openEuler" ]; then
        sudo yum install -y openeuler-lsb
    else
        echo "Unsupported OS: $OS_ID"
        exit 1
    fi
}

configure_rear() {
    if [ -f "$REAR_CONF_FILE" ]; then
        sudo rm -f "$REAR_CONF_FILE"
    fi

    cat > "$REAR_MKBACKUP_SCRIPT" << EOF
#!/bin/bash
echo "OUTPUT=PXE" >> $REAR_CONF_FILE
echo "BACKUP=NETFS" >> $REAR_CONF_FILE
echo "BACKUP_URL=$BACKUP_URL" >> $REAR_CONF_FILE
echo "BACKUP_PROG_EXCLUDE=(${BACKUP_PROG_EXCLUDE[@]} '/media' '/var/tmp' '/var/crash' '/tmp')" >> $REAR_CONF_FILE
echo "NETFS_KEEP_OLD_BACKUP_COPY=yes" >> $REAR_CONF_FILE
echo "MODULES=('all_modules')" >> $REAR_CONF_FILE
echo "COPY_AS_IS+=('/usr/lib64/libsepol.so.2')" >> $REAR_CONF_FILE
EOF

    chmod +x "$REAR_MKBACKUP_SCRIPT"
}

execute_backup() {
    bash "$REAR_MKBACKUP_SCRIPT"

    if ! rear -d -v mkbackup; then
        echo "Failed to execute backup."
        exit 1
    fi

    rm -f "$REAR_MKBACKUP_SCRIPT"

    echo "Execute backup created successfully."
}

handle_backup() {
    setup_environment
    configure_rear
    execute_backup
}

case {{ .handleType }} in
    backup)
        handle_backup
        ;;
    *)
        echo "Usage: $0 {backup|update|recover}"
        exit 1
        ;;
esac
`, nfsServer, nfsPath)

	datastore := map[string]interface{}{}
	datastore["handleType"] = handleType
	shellfile, err := template.TemplateRender(shell, datastore)
	if err != nil {
		return "", err
	}

	return shellfile, nil
}
