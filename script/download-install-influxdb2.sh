#!/bin/bash
# Script to download and install InfluxDB on a Raspberry Pi
# but might also work on other systems
# Author: David Thorpe <djt@mutablelogic.com>
#
# Usage:
#   download-install.sh [-f] [-u username]
#
# Flag -f will remove any existing installations first

#####################################################################

# This is the URL for downloading the InfluxDB dist
INFLUXDB_URL="https://dl.influxdata.com/influxdb/releases/influxdb_2.0.0-alpha.6_linux_arm64.tar.gz"
# PREFIX is the parent directory of the influxdb setup
PREFIX="/opt"
# USERNAME is the username for the influx processes
USERNAME="influxdb"
# PRODUCT is the name of the product
PRODUCT="influxdb2"
# VARPATH
VARPATH="/var/lib/${PRODUCT}"
# SERVICENAME
SERVICENAME="${PRODUCT}.service"
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
INFLUXDB_FILENAME=`basename "${INFLUXDB_URL}"`
INFLUXDB_PATH="${TEMP_DIR}/${INFLUXDB_FILENAME}"

echo "Downloading ${INFLUXDB_FILENAME}"
"${CURL_BIN}" "${INFLUXDB_URL}" -s -o "${INFLUXDB_PATH}" || exit 2

if [ ! -f "${INFLUXDB_PATH}" ] ; then
  echo "Cannot download distribution"
  rm -fr "${TEMP_DIR}"
  exit 2
fi

# Unarchive and obtain the folder name
echo "Unarchiving"
tar -C "${TEMP_DIR}" -zxf "${INFLUXDB_PATH}"
INFLUXDB_PATH=`find "${TEMP_DIR}" -maxdepth 1 -mindepth 1 -type d -print`
if [ ! -d "${INFLUXDB_PATH}" ]; then
  echo "Cannot unpack distribution"
  rm -fr "${TEMP_DIR}"
  exit 2
fi

# Move the folder into the PREFIX directory
INFLUXDB_DIST=`basename "${INFLUXDB_PATH}"`
if [ -d "${PREFIX}/${INFLUXDB_DIST}" ] ; then
  if [ "${FORCE}" = "1" ] ; then
    rm -fr "${PREFIX}/${INFLUXDB_DIST}" || exit 3
  else
      echo "Distribution already exists: ${PREFIX}/${INFLUXDB_DIST}"
      echo "(use -f flag to remove the folder first)"
      rm -fr "${TEMP_DIR}"
      exit 3
  fi
fi

if [ -e "${PREFIX}/${PRODUCT}" ] ; then
  if [ "${FORCE}" = "1" ] ; then
    rm "${PREFIX}/${PRODUCT}" || exit 3
  else 
      echo "Distribution already exists: ${PREFIX}/${PRODUCT}"
      echo "(use -f flag to remove this symbolic link first)"
      rm -fr "${TEMP_DIR}"
      exit 3
  fi
fi
