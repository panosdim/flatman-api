#!/bin/sh
if [ "$(id -u)" != "0" ]; then
    exec sudo bash "$0" "$@"
fi

is_service_exists() {
    x="$1"
    if systemctl status "${x}" 2>/dev/null | grep -Fq "Active:"; then
        return 0
    else
        return 1
    fi
}

INSTALL_PATH=/opt/flatman
EXEC_NAME=flatman
SERVICE_NAME=$EXEC_NAME.service

# Build software

go build -o $EXEC_NAME
ret_code=$?
if [ $ret_code != 0 ]; then
    printf "Error: [%d] when building executable. Check that you have go tools installed." $ret_code
    exit $ret_code
fi

# Check if needed files exist
if [ -f .env ] && [ -f $EXEC_NAME ] && [ -f $SERVICE_NAME ]; then
    # Check if we upgrade or install for first time
    if is_service_exists "$SERVICE_NAME"; then
        systemctl stop $SERVICE_NAME
        cp $EXEC_NAME $INSTALL_PATH
        cp .env $INSTALL_PATH
        systemctl start $SERVICE_NAME
    else
        mkdir -p $INSTALL_PATH
        cp $EXEC_NAME $INSTALL_PATH
        cp .env $INSTALL_PATH
        cp $SERVICE_NAME /usr/lib/systemd/system
        systemctl start $SERVICE_NAME
        systemctl enable $SERVICE_NAME
	cp flatman.conf $NGINX_CONF_PATH
        nginx -s reload
    fi
else
    echo "Not all needed files found. Installation failed."
    exit 1
fi
