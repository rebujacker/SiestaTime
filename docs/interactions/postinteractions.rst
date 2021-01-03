Post. Servers
=====================================================

.. figure:: ../_static/images/interactions/handlers.png
    :align: right
    :figwidth: 300px
    :target: ../_static/images/interactions/handlers.png

Most of the previously created Handlers can be interacted with. This menu will be found on “Handlers” tab


Interact with SSH - Full Interactive Shell
--------------------------------------------

Interact with a SSH Rev Handler will be possible once a bichito has been attached to it. 

.. figure:: ../_static/images/interactions/sshinteract.png
    :align: center
    :figwidth: 600px
    :target: ../_static/images/interactions/sshinteract.png

.. note::
    The process and connection will finish once you close the Interactive shell

.. warning::
    For the moment these ssh are full Interactive for Linux and OSX (no windows) ``TBD``

Interact with SSH opening a SOCKS5
--------------------------------------------

Interact with a SSH Rev Handler will be possible once a bichito has been attached to it.
This will open a SOCKS5 in the Operator's machine ready to proxy all traffict through implanted device. 

.. warning::
    You will need to kill SSH sessions in the target staging server to avoid extra process/threads to be running in the implanted device


Kill SSH Sessions
--------------------------------------------

Interact with a targer Post/Staging Server and kill any connected POST SESSIONS


Empire and Metasploit
--------------------------------------------------

Interacting with both MSFT and Empire is very straightforward. Once created just click “Interact” and you will be provided by a console directly bonded to the target's MSFT/Empire Server.

.. figure:: ../_static/images/interactions/empire.png
    :align: center
    :figwidth: 600px
    :target: ../_static/images/interactions/empire.png