# InfluxPi

This repository contains installation scripts for InfluxDB, Telegraf and Chronograf
on the Raspian distribution for the Raspberry Pi. I guess it may work for other
installations but this is the one I'm interested in. To use, the following commands
can be executed:

```
bash% script/download-install-influxdb.sh
bash% script/download-install-chronograf.sh
bash% script/download-install-telegraf.sh
```

For all the scripts, you can use the `-f` flag to remove older versions of 
software and/or the `-u <username>` flag to indicate the user under 
which the database should run as, with a default of 'influxdb'. You should
use the same username for all packages. The various defaults within the script 
can be changed. They are:

  * __INFLUXDB_URL__ where the InfluxDB binary comes from. Currently https://dl.influxdata.com/influxdb/releases/influxdb-1.4.2_linux_armhf.tar.gz
  * __CHRONOGRAF_URL__ where the InfluxDB binary comes from. Currently https://dl.influxdata.com/chronograf/releases/chronograf-1.3.10.0_linux_armhf.tar.gz
  * __TELEGRAF_URL__ where the InfluxDB binary comes from. Currently https://dl.influxdata.com/telegraf/releases/telegraf-1.4.4_linux_armhf.tar.gz
  * __TELEGRAF_VERSION__ should be changed if you change the URL, due to the way that
    distributon is packaged up.
  * __PREFIX__ where the installation is made. Currently `/opt`

If you then want to run the command-line client, then you can add the following
line to your ~/.bash_profile

```
  # influxdb
  export INFLUXDB="/opt/influxdb"
  export PATH="${PATH}:${INFLUXDB}/usr/bin"
```

License
-------

Copyright 2017 David Thorpe

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
