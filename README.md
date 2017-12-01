# InfluxPi

This repository contains installation scripts for InfluxDB on Raspberry Pi.
To use, the following command can be executed:

```
bash% script/download-install.sh
```

You can use the `-f` flag to remove older versions of InfluxDB (force) and the
`-u <username>` flag to indicate the user under which the database should run
as, with a default of 'influxdb'.

The various defaults within the script can be changed. They are:

  * __INFLUXDB_URL__ where the InfluxDB binary comes from. Currently https://dl.influxdata.com/influxdb/releases/influxdb-1.4.2_linux_armhf.tar.gz
  * __PREFIX__ where the installation is made. Currently `/opt`

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
