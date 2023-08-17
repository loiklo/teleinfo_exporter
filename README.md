# teleinfo_exporter
Prometheus exporter for french Linky teleinfo.
Export data read on the tic serial port and expose metrics, ready to be scrapped by Prometheus.

## Build
The build for the ARM arch force the arm6 compatibility to be able to run the exporter on **RaspberryPi Zero W** or **RaspberryPi 1**.
The build for the ARM version can't be done on a RaspberryPi 1 with the precompiled Golang compiler which requires arm7. The tips is to build the arm6 binary on an amd64 to arm7 architecture.

```bash
alex$ make
```

## Installation
```bash
root# mkdir -p /opt/teleinfo_exporter/homedir
#cp teleinfo_exporter_armv6l /opt/teleinfo_exporter/.
root# useradd -d /opt/teleinfo_exporter/homedir -G dialout -s /bin/false teleinfo-exporter
```

```bash
root# cat /etc/systemd/system/teleinfo-exporter.service 
[Unit]
  Description=teleinfo_exporter
  After=time-sync.target
[Service]
  User=teleinfo-exporter
  Group=teleinfo-exporter
  ExecStart=/opt/teleinfo_exporter/teleinfo_exporter_armv6l
  WorkingDirectory=/opt/teleinfo_exporter/homedir
  Restart=on-failure
  RestartSec=10
[Install]
  WantedBy=multi-user.target
```

```bash
root# systemctl daemon-reload
root# systemctl enable teleinfo-exporter.service
root# systemctl start teleinfo-exporter.service
```

## Example of data output
```bash
alex$ curl http://localhost:9150/metrics
# HELP tic_iinst Intensite instantanee
# TYPE tic_iinst gauge
tic_iinst{adco="012345678901"} 15
# HELP tic_index Index
# TYPE tic_index counter
tic_index{adco="012345678901",color="blue",option="tempo",phase="hc"} 285447
tic_index{adco="012345678901",color="blue",option="tempo",phase="hp"} 237281
tic_index{adco="012345678901",color="red",option="tempo",phase="hc"} 12327
tic_index{adco="012345678901",color="red",option="tempo",phase="hp"} 5585
tic_index{adco="012345678901",color="white",option="tempo",phase="hc"} 52366
tic_index{adco="012345678901",color="white",option="tempo",phase="hp"} 79858
# HELP tic_papp Puissance apparente
# TYPE tic_papp gauge
tic_papp{adco="012345678901"} 3510
```

## ToDo
- Allow to set the serial port in parameters
- Allow to change the http port
