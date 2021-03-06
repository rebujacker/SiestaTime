Client
===========================


Once Hive is Online, operators can connect to it, but they need to install their client first


What this does? What is the Client
------------------------------------

GUI electron powered and Go application running in Operators' devices. The go client executable run a localhost server that feeds data to the electron app. In the same time, the go client executable authenticate against Hive with credentials passed in compiled time.
Clients will authenticate against hive and send jobs to it. They can also interact with Post. Servers thanks to a Hive tunnelization system

* Install Client OS dependencies
* Use config<VPC>.txt to get user credentials
* Download Go and dependencies to compile client
* Configure and install electron app



Install
-----------

Run ``stime.sh`` for client installation. This will compile stclient with provided credentials and configuration, and will generate the GUI folder.
``stime.sh install <username> <password> <Hive IP/Domain> <Hive Port> <Client Port> <Hive VPC certificate Fingerprint>``

.. prompt:: bash $

    ./stime.sh install admin admin 13.57.31.79 6232 8000 $(openssl x509 -fingerprint -sha256 -noout -in ./installConfig/hive.pem | cut -d '=' -f2)


Run the Client
---------------


.. prompt:: bash $

    ./stime.sh

Uninstall
-----------

Will remove go dependencies 
``./stime.sh remove``

