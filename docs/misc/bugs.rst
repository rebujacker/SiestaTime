Known Bugs
===========================


Client
--------------------------

Job Creation - Stuck
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
The client software implements a lock to avoid the triggering of too many jobs against Hive. This still can be faulty and the lock stuck to 0.
If this happens a restart on the GUI/Client is recommended.

Hive
--------------------------------------------

Jobs Queue - GUI
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Hive will execute jobs 1 by 1. Once an output for each Job is reached, the Jobs will update.
This means that while Hive is busy operators will not see any updates on the client for sent jobs.

DB Locked
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Multiple writes on Hive DB have been shown to block DB and miss writes.
This normally happens when multiple “Hive Jobs” are sequenced too simultaneously.


Redirector
--------------------------------------------------

Not common known problems are still known in redirector software
Deadlock is known to have happened in the past, but redirectors auto-restart if this happens.

.. note::
	Having a bunch of redirectors if recommended if some of them fails

Bichito
--------------------------------------------------

Not common known problems are still known in implant/Bichito software

.. note::
	Persistence is recommended in the scenario where the bichito/implant software could block its functionality on unknown bugs


Post. Servers
--------------------------------------------------

.. note::
	Faulty Staging/Post. Can be easily trashed and re-deployed. If the user feels any of them is becoming faulty, use the rapid deployment properties of STime to create it again. 