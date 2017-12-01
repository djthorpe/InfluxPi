#!/bin/bash
# Script to download and install Chronograf on a Raspberry Pi
# but might also work on other systems
# Author: David Thorpe <djt@mutablelogic.com>
#
# Usage:
#   download-install-chronograf.sh [-f] [-u username]
#
# Flag -f will remove any existing installations first

#####################################################################

# This is the URL for downloading the Chronograf dist
CHRONOGRAF_URL="https://dl.influxdata.com/chronograf/releases/chronograf-1.3.10.0_linux_armhf.tar.gz"
# PREFIX is the parent directory of the influxdb setup
PREFIX="/opt"
# USERNAME is the username for the influx processes
USERNAME="influxdb"
# VARPATH
VARPATH="/var/lib/influxdb"
# SERVICENAME
SERVICENAME="chronograf.service"
# FORCE set to 1 will result in any existing installation being
# removed first
FORCE=0

#####################################################################
# PROCESS FLAGS

while getopts 'fu:' FLAG ; do
  case ${FLAG} in
    f)
	  FORCE=1
      ;;
    u)
	  USERNAME=${OPTARG}
      ;;      
    \?)
      echo "Invalid option: -${OPTARG}"
	  exit 1
      ;;
  esac
done

#####################################################################
# CHECKS

# Temporary location
TEMP_DIR=`mktemp -d`
if [ ! -d "${TEMP_DIR}" ]; then
  echo "Missing temporary directory: ${TEMP_DIR}"
  exit 1
fi

# Ensure script is run as root user
if [ "${USER}" != "root" ]; then
  echo "This script must be run as root user"
  exit 1
fi

# Ensure we have curl
CURL_BIN=`which curl`
if [ ! -x "${CURL_BIN}" ] ; then
  echo "Missing curl"
  exit 1
fi

#####################################################################
# DOWNLOAD AND INSTALL

# Create the prefix directory if necessary
install -d "${PREFIX}" || exit -1

# Download the code
CHRONOGRAF_FILENAME=`basename "${CHRONOGRAF_URL}"`
CHRONOGRAF_PATH="${TEMP_DIR}/${CHRONOGRAF_FILENAME}"

echo "Downloading ${CHRONOGRAF_FILENAME}"
"${CURL_BIN}" "${CHRONOGRAF_URL}" -s -o "${CHRONOGRAF_PATH}" || exit 2

if [ ! -f "${CHRONOGRAF_PATH}" ] ; then
  echo "Cannot download distribution"
  rm -fr "${TEMP_DIR}"
  exit 2
fi

# Unarchive and obtain the folder name
echo "Unarchiving"
tar -C "${TEMP_DIR}" -zxf "${CHRONOGRAF_PATH}"
CHRONOGRAF_PATH=`find "${TEMP_DIR}" -maxdepth 1 -mindepth 1 -type d -print`
if [ ! -d "${CHRONOGRAF_PATH}" ]; then
  echo "Cannot unpack distribution"
  rm -fr "${TEMP_DIR}"
  exit 2
fi

# Move the folder into the PREFIX directory
CHRONOGRAF_DIST=`basename "${CHRONOGRAF_PATH}"`
if [ -d "${PREFIX}/${CHRONOGRAF_DIST}" ] ; then
  if [ "${FORCE}" = "1" ] ; then
    rm -fr "${PREFIX}/${CHRONOGRAF_DIST}" || exit 3
  else
      echo "Distribution already exists: ${PREFIX}/${CHRONOGRAF_DIST}"
      echo "(use -f flag to remove the folder first)"
      rm -fr "${TEMP_DIR}"
      exit 3
  fi
fi

if [ -e "${PREFIX}/chronograf" ] ; then
  if [ "${FORCE}" = "1" ] ; then
    rm "${PREFIX}/chronograf" || exit 3
  else 
      echo "Distribution already exists: ${PREFIX}/chronograf"
      echo "(use -f flag to remove this symbolic link first)"
      rm -fr "${TEMP_DIR}"
      exit 3
  fi
fi

echo "Making link: ${PREFIX}/${CHRONOGRAF_DIST} -> ${PREFIX}/chronograf"
cd "${PREFIX}" || exit 3
mv "${CHRONOGRAF_PATH}" "." || exit 3
rm -fr "${TEMP_DIR}" || exit 3
rm -f "chronograf" || exit 3
ln -s "${CHRONOGRAF_DIST}" chronograf || exit 3

#####################################################################
# MAKE USERS AND GROUPS

USERID=`id -u ${USERNAME} 2> /dev/null`
if [ "${USERID}" = "" ]; then
  echo "Making users and groups for ${USERNAME}"
  useradd  -M -U -s /usr/sbin/nologin -d "${VARPATH}" -r "${USERNAME}" || exit 4
fi

#####################################################################
# MAKE VAR DIRECTORY

echo "Making ${VARPATH}"
install -d "${VARPATH}" -o "${USERNAME}" -g "${USERNAME}" || exit 5
chown -R "${USERNAME}:${USERNAME}" "${VARPATH}" || exit 5


#####################################################################
# UNLOAD SERVICE

SERVICE_LOADED=`systemctl list-units | grep ${SERVICENAME}`
if [ ! "${SERVICE_LOADED}" = "" ] ; then
  echo "Existing ${SERVICENAME} loaded, removing"
  systemctl stop ${SERVICENAME} || exit 6
  systemctl disable ${SERVICENAME} || exit 6
  rm -f "/etc/systemd/system/${SERVICENAME}" || exit 6
  systemctl daemon-reload
  systemctl reset-failed
fi

#####################################################################
# CREATE THE CONFIGURATION FILE

SYSTEMCTL_FILE="${PREFIX}/chronograf/usr/lib/chronograf/scripts/chronograf.service"
BIN_FILE="${PREFIX}/chronograf/usr/bin/chronograf"
BOLT_PATH="${VARPATH}/chronograf-v1.db"
CANNED_PATH="${PREFIX}/chronograf/usr/share/chronograf/canned"

if [ ! -f "${SYSTEMCTL_FILE}" ] ; then
  echo "Missing systemctl service file: $SYSTEMCTL_FILE"
  exit 7
fi

if [ ! -x "${BIN_FILE}" ] ; then
  echo "Missing daemon executable file: $BIN_FILE"
  exit 7
fi

if [ ! -d "${CANNED_PATH}" ] ; then
  echo "Missing canned path: $CANNED_PATH"
  exit 7
fi

echo "Creating service file /etc/systemd/system/${SERVICENAME}"
cat "${SYSTEMCTL_FILE}" \
  | sed "s/User=.*/User=${USERNAME}/g" \
  | sed "s/Group=.*/Group=${USERNAME}/g" \
  | sed "s/EnvironmentFile=.*/Environment=CHRONOGRAF_OPTS=/g" \
  | sed "s/ExecStart=.*/ExecStart=${BIN_FILE//\//\\/} -b ${BOLT_PATH//\//\\/} -c ${CANNED_PATH//\//\\/} \$CHRONOGRAF_OPTS/g" \
  > /etc/systemd/system/${SERVICENAME} || exit 7

#####################################################################
# LOAD THE SERVICE

echo "Starting ${SERVICENAME}"
systemctl daemon-reload || exit 8
systemctl start ${SERVICENAME} || exit 8
sleep 1
systemctl status ${SERVICENAME} -l || exit 8


