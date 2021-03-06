Implants
===========================


Generate the different executables to keep connection to the foothold and their C2 (redirectors)


.. figure:: ../_static/images/deployment/craftimplant.png
    :align: center
    :figwidth: 600px
    :target: ../_static/images/deployment/craftimplant.png



What is a module?
-------------------

Siesta Time provides you options that define essential properties of the implant.
For the moment modules can be chosen for both ``egression`` and ``persistence`` 
In the future other modules could be chosen for more specific implant properties and behaviour once executed


Basic Configurations
--------------------------------------------

Time to Live (TTL)
~~~~~~~~~~~~~~~~~~~~~~
Seconds to finish kill itself if the implant is not able to reach any of the redirectors

Response Time
~~~~~~~~~~~~~~~~~~~~~~
Seconds to response to a job arrived through any redirector


Network Modules
--------------------------------------------------

Network modules are platform-agnostic. It defines the communication channel between the implant and the redirectors. Each module has its own parameters.

.. figure:: ../_static/images/deployment/networkmodules.png
    :align: center
    :figwidth: 600px
    :target: ../_static/images/deployment/networkmodules.png


.. note::
	Every network module has its ``No Infrastructure`` version which let operators define IP/Domains that are not part of the Hive battery. This will trigger the ``Offline`` mode to generate Implants. In this mode Hive will not deploy redirectors.

Self-Signed HTTPS Go
~~~~~~~~~~~~~~~~~~~~~~

Https Client: Go
Redirector Certificate: Self Signed

TLS Port → Choose the Redirector https listening port

Paranoid Self Signed HTTPS Go
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Https Client: Go
Redirector Certificate: Self Signed
The implant will check target redirector fingerprint 

TLS Port → Choose the Redirector https listening port

.. warning::
	This module will avoid the implant to egress if there is a https TLS proxy on the middle

SaaS: Gmail API
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Https Client: Go
Redirector Certificate: Google Servers


SaaS: Gmail API,Mimic Https Client Fingerprints
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Https Client: Mimic target TLS Fingerprint
Redirector Certificate: Google Servers

UserAgent → Choose the client User Agent
TLS JA3C Fingerprint → Choose the Redirector TLS Fingerprint

.. warning::
	This module will try to evade most of NIDS non DPI (deep packet inspection) based


Persistence Modules
--------------------------------------------------

Once the red team achieves the execution of the implant through a delivery method the most important next step after checking C2 connectivity (or even before) will be to persist. This will let the implant to re-execute itself after any device shutdown/re-login

Windows - schtaks 
~~~~~~~~~~~~~~~~~~~~~~

Use windows  C++ ``comsupp.lib`` and ``taskschd.lib`` to create a task that runs on user login

.. figure:: ../_static/images/deployment/persistence.png
    :align: center
    :figwidth: 600px
    :target: ../_static/images/deployment/persistence.png

Name → task name that will appear on windows task list
Path → Windows path for target Implant Binary ( ``$(UserHOME)\“folder1\folder2\execuablenameand.extension`` )

Darwin - launchd 
~~~~~~~~~~~~~~~~~~~~~~

Using Go file libraries, it writes a new launchd service. MacOSX will fetch it on user reboot/login and load the implant.

.. figure:: ../_static/images/deployment/persistencelaunchd.png
    :align: center
    :figwidth: 600px
    :target: ../_static/images/deployment/persistencelaunchd.png

Name → Launchd file name that will appear on windows task list
Path → Darwin path for target Implant Binary  ( ``$(UserHOME)/“folder1/folder2/execuablenameand.extension`` )

Linux - XDG 
~~~~~~~~~~~~~~~~~~~~~~

User Desktop Linux devices with graphical interface respect some specification that let users to configure default tasks on user login.
Files are written using default GO file libraries.

.. figure:: ../_static/images/deployment/persistencexdg.png
    :align: center
    :figwidth: 600px
    :target: ../_static/images/deployment/persistencexdg.png


Name → Autostart file name that will appear on windows task list
Path → Linux path for target Implant Binary ( ``$(UserHOME)/“folder1/folder2/execuablenameand.extension`` )

.. note::
	Extra doc: `xdg`_ 

Redirectors
--------------------------------------------------

C2 Servers where implants try to connect to retrieve commands from Hive. Implants will try one by one to connect to them, once they find out which one the can reach by the previous network methods.

Online 
~~~~~~~~~~~~~~~~~~~~~~

.. figure:: ../_static/images/deployment/redsonline.png
    :align: center
    :figwidth: 600px
    :target: ../_static/images/deployment/redsonline.png


.. warning::
	Domains that are active already will not appear on the list

Offline
~~~~~~~~~~~~~~~~~~~~~~

.. figure:: ../_static/images/deployment/redsoffline.png
    :align: center
    :figwidth: 600px
    :target: ../_static/images/deployment/redsoffline.png

.. _xdg: https://specifications.freedesktop.org/autostart-spec/autostart-spec-latest.html