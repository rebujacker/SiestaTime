Theory: Why Siesta Time ?
===========================

Concepts
--------------------------------------------------

Hive
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Is the main “Operation Server”. Hive will receive every command/job from authenticated clients and process or redirect them to a target foothold.
The main DB of siesta time will be in this server, with all information/configurations from the red team operators’ actions.

Operator
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Is the equivalent of the user in Siesta Time. The creator of hive or “Admin” (first user) will be the only one able to add new Operators.
Operators added will be able to compile their own client and connect to Hive.

VPS/VPC
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
“Virtual Private Service/Cloud” are sets of credentials saved in Hive that can be used to deploy redirectors that backbone implants’ connection and Post./Staging Servers to interact with them later on.

Domain
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Hive will be able to store a set of credentials to manipulate a target domain at will. Used to map its resolution to the generated Server’s infrastructure selected VPCs.  Once an element is requested to be created (implant,post. server...)

SaaS
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
“Software as a Service” are sets of credentials from an internet service that let implants to egress using string data

Implant
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Implants are composed by a number of redirectors and the compiled executables for different platforms (linux,darwin and windows).
The implants will have an array of redirectors to connect to, that will be in the shape of a target IP,domain or SaaS account.
These redirectors will be deployed in the creation of the implant

Redirector
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Host  with a server software running as a service. Its purpose is to redirect jobs from footholds (bichitos) to the Hive, and vice versa.They are automatically deployed on implant creation.
In offline mode, the redirector executable can be downloaded to be installed in any desired host

Bichito
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
The main implant of Siesta Time. They are generated in the shape of an executable for a target platform. Once executed, they will appear as an online process within the created implant and attached to a redirector. 
On the implant creation it is possible to choose the capabilities of
the Bichitos. How will egress through the network, his time to day, persistence… these are the modules

Staging/Post
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Servers whose objective is performing delivery, staging and post-exploitation tasks.
Operators can directly connect/interact with them.

Report
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
These elements are a plain text file that holds every Job processed by Hive, and every command typed on Staging/Posts servers

Client
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
GUI on electron and go application running in operators devices.
Clients will authenticate against hive and send jobs to it. They can also interact with Post. Servers thanks to a Hive tunnelization.


