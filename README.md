# Stats network manager

## Description

This tool was created as an utility tool for other legacy scripts that collect statistics from network elements.

Functions of current app:

1. collect information about network elements managed by NMS via SOAP requests
2. set trusted managers for a specific network element type (9500 MPR)
3. stores results in files as legacy scripts uses files, not db

## TODO

-   [ ] init templating
-   [ ] change logging to zap or zerolog
-   [ ] use SQLite to store info about network elements and policy used
-   [ ] restructure this app to become a micro-service, use REST API for frontend
-   [ ] add another micro-service to replace legacy scripts to collect staticstics via snmp
-   [ ] move functional of snmp (trusted managers) to new micro-service
-   [ ] use InfluxDB or VictoriaMetrics to store statics
-   [ ] use protobuf + grpc to connect services
-   [ ] use one more micro-service to collect and convert statistics to format accepted by client
