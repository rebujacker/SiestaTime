Siesta Time
===========================

*siestatime: Red Team Automation Tool*

.. image::  https://github.com/rebujacker/SiestaTime/blob/master/src/client/electronGUI/static/icons/png/STicon.png
    :height: 64px
    :width: 64px
    :alt: STime logo

+----------------+--------------------------------------------------+
|Project site    | https://github.com/rebujacker/SiestaTime         |
+----------------+--------------------------------------------------+
|Issues          | https://github.com/rebujacker/SiestaTime/issues/ |
+----------------+--------------------------------------------------+
|Author          | Alvaro Folgado (rebujacker)                      |
+----------------+--------------------------------------------------+
|Documentation   | https://siestatime.readthedocs.io                |
+----------------+--------------------------------------------------+
|Last Version    | Beta 1.0 (1.0b1)                                 |
+----------------+--------------------------------------------------+
|Golang versions | 1.10.3 or above                                  |
+----------------+--------------------------------------------------+

What's Siesta Time, What is its Purpose?
===================================================

Red Team Automation tool powered by go and terraform.

Red Team operations require substantial efforts to both create implants and a resilient C2 infrastructure. SiestaTime aims to merge these ideas into a tool with an easy-to-use GUI, which facilitates implant and infrastructure automation alongside its actors reporting.
SiestaTime allows operators to provide registrar, SaaS and VPS credentials in order to deploy a resilient and ready to use Red Team infrastructure. The generated implants will blend-in as legitimate traffic by communicating to the infrastructure using SaaS channels and/or common network methods.

Use your VPS/Domains battery to deploy staging servers and inject your favorite shellcode for interactive sessions, clone sites and hide your implants ready to be downloaded, deploy more redirectors if needed. All these jobs/interactions will be saved and reported to help the team members with the documentation process.

SiestaTime is built entirely in Golang, with the ability to generate Implants for multiple platforms, interact with different OS resources, and perform efficient C2 communications. Terraform used to deploy/destroy different Infrastructure.

This will help increase companies red teams efficiency, improving industry security standards and make the defenders to catch-up , being ready for real threats.


This tool has both **Educational Purposes** and aims to help **security industry** and **defenders**.


Short/Flash install
===================================

Fast Examples of Siesta Time Installation (look at the doc for dependencies and detailed process)

Offline
----------------------


.. code-block:: bash

    #Ubuntu
    ./hive.sh installOffline 0.0.0.0 6232 /usr/local/ admin admin
    ./stime.sh install admin admin 0.0.0.0 6232 8000 $(openssl x509 -fingerprint -sha256 -noout -in ./installConfig/hive.pem | cut -d '=' -f2)
    ./stime.sh



Online
----------------------

* Complete ``SiestaTime/installConfig/configAWS.txt``

::

    USERNAME : Admin Username
    PASSWORD : Admin password
    port : HTTPS Hive port listener
    accesskey: AWS accesskey
    secretkey: AWS secretkey
    Region: AWS region
    Keyname: AWS keyname (without .pem)
    ami: aws ami 
    itype: AWS ec2 itype

* Copy AWS key to ``SiestaTime/installConfig/<keyname>.pem``


.. code-block:: bash

    # Ubuntu
    ./hive.sh installaws
    ./stime.sh install admin admin 0.0.0.0 6232 8000 $(openssl x509 -fingerprint -sha256 -noout -in ./installConfig/hive.pem | cut -d '=' -f2)
    ./stime.sh



Available features
==================================

**This is a Beta**

The tool miss yet a lot of work and, most importantly, **bug fixing**

**Current Modules/Abilities**

Hive:
    - VPS 
        - AWS
    - Domain
        - GO Daddy
    - SaaS
        - Gmail API

Stagings/POST.:
    - Droplet
    - Reverse SSH
    - Reverse SOCKS5 SSH
    - MSF Handler: HTTPS Let's Encrypt
    - Empire Handler: HTTPS Let's Encrypt

Reporting:
    - Basic Text Reports (Jobs,Staging/Post. Sessions)

Bichito:

- Network Egression:
    - HTTPS Paranoid GO
    - Self-Signed HTTPS GO
    - Gmail API
    - Gmail API - Mimic TLS

- Persistence:
    - Windows - schtasks
    - Linux - XDG
    - Darwin - launchd

- Interaction:
    - Bichiterpreter (Job Based): exec (using os.exec),sysinfo,upload, [...]
    - Inject Launchers (using os.exec)
    - Reverse SSH
    - Reverse SOCKS5 SSH


- <Future Abilities>


Documentation
==============================

User Guide: https://siestatime.readthedocs.io

Developer Guide: **TBD**

Presented at Defcon 27 Red Team Village.
Slides from: https://redteamvillage.io/ --> https://www.slideshare.net/AlvaroFolgadoRueda1/siestatime-defcon27-red-team-village

Presented at Defcon 28 Red Team Village.
Slides from: https://redteamvillage.io/ --> https://www.slideshare.net/AlvaroFolgadoRueda1/total-eagression-defcon


Contributing
=============================

Any collaboration is welcome! The Bigger the tool modules set is, the better testing options could be addressed in future Assessments.

Red Teamers and Offensive Security Engineers call for code/modules! :)

There are many tasks to do. You can check the `Issues <https://github.com/rebujacker/SiestaTime/issues/>`_ and send us a Pull Request.


Disclaimer
===================================

Author/Contributors will not be responsible for the malfunctioning or weaponization of this code

License
========================

This project is distributed under `GPL V3 license <https://github.com/rebujacker/SiestaTime/LICENSE>`_
