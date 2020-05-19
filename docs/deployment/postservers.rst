Post. Servers
===========================


.. figure:: ../_static/images/deployment/deploypost.png
    :align: center
    :figwidth: 600px
    :target: ../_static/images/deployment/deploypost.png

Droplet
--------------------------------------------

A plain ubuntu Server used to drop created Implants 

Name → Just the name of the resource
VPS/Domain → Choose one from the battery
Https Port
Path for Implant folder → ``/var/www/”path”/implant``
Endpoint: ``https://domain/implantpath/implant``

.. figure:: ../_static/images/deployment/droplet.png
    :align: right
    :figwidth: 300px
    :target: ../_static/images/deployment/droplet.png

Reverse SSH
--------------------------------------------------

The reverse SSH will create a ubuntu Server with a sshd connection and a “anonymous” user configured with its own keys.
This user is configured without a bash shell, the idea is that the implant will connect with that anonymous user and serve its own SSHD Server from Golang Code.

In this way, what it looks from the foothold as a SSH outbound connection, will be a remote bash/cmd serving.
***[More details in developer Guide]***


.. warning::
    For the moment these SSH are not fully interactive

Reverse RDESKTOP
--------------------------------------------------

TBD


Empire, Metasploit
--------------------------------------------------

.. figure:: ../_static/images/deployment/msftworking.png
    :align: right
    :figwidth: 300px
    :target: ../_static/images/deployment/msftworking.png

Similarly to the reverse SSH server, Empire and MSFT will create remote handlers to receive incoming shells.
***The handler is configured with https self signed certificates.***
You can choose the handler port

