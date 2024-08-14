/*
 *
 * Copyright 2024 KylinSoft  Co., Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * /
 */
package common

import "universal_os_upgrader/pkg/template"

const (
	HandleBackup  = "backup"
	HandleUpdate  = "update"
	HandleRecover = "recover"
)

func GetRearShell(handleType string) (string, error) {
	shell := `
#!/bin/bash

BACKUP_DIR="/root"
REAR_CONF_FILE="/etc/rear/local.conf"
REAR_MKBACKUP_SCRIPT="/tmp/rear-mkbackup.sh"

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
echo "OUTPUT=ISO" >> $REAR_CONF_FILE
echo "OUTPUT_URL=null" >> $REAR_CONF_FILE
echo "BACKUP=NETFS" >> $REAR_CONF_FILE
echo "BACKUP_URL=iso:///backup" >> $REAR_CONF_FILE
echo "ISO_DIR=$BACKUP_DIR" >> $REAR_CONF_FILE
echo "BACKUP_PROG_EXCLUDE=(${BACKUP_PROG_EXCLUDE[@]} '/media' '/var/tmp' '/var/crash' '/tmp')" >> $REAR_CONF_FILE
echo "NETFS_KEEP_OLD_BACKUP_COPY=yes" >> $REAR_CONF_FILE
echo "MODULES=('all_modules')" >> $REAR_CONF_FILE
EOF

    chmod +x "$REAR_MKBACKUP_SCRIPT"
}

create_recovery_iso() {
    bash "$REAR_MKBACKUP_SCRIPT"

    if ! rear -d -v mkbackup; then
        echo "Failed to create recovery ISO."
        exit 1
    fi

    rm -f "$REAR_MKBACKUP_SCRIPT"

    echo "Recovery ISO created successfully."
}

handle_backup() {
    setup_environment
    configure_rear
    create_recovery_iso
}

handle_update() {
    handle_backup
}

handle_recover() {
    echo "Comming soon..."
}

case {{ .handleType }} in
    backup)
        handle_backup
        ;;
    update)
        handle_update
        ;;
    recover)
        handle_recover
        ;;
    help)
        echo "Usage: $0 {backup|update|recover|help}"
	exit 1
        ;;
    *)
        echo "Usage: $0 {backup|update|recover}"
        exit 1
        ;;
esac
`
	datastore := map[string]interface{}{}
	datastore["handleType"] = handleType
	shellfile, err := template.TemplateRender(shell, datastore)
	if err != nil {
		return "", err
	}

	return shellfile, nil
}
