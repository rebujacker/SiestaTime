Hive
===========================


Hive is the main Operations Server of Siesta Time. Is the first element that needs to be deployed.
The user has the option to do it ``Offline``, or using terraform to deploy Hive with target Virtual Private Cloud resources.
Configurations of target VPC to use will be saved in a txt file by the syntax of: ``./SiestaTime/installconfig/config<VPC>.txt``
There is the option to deploy without using a VPC, or ``Offline``, this will install Hive in the current host.


What this does? What is Hive?
------------------------------

Hive is the main “Operation Server”. Hive will receive every command/job from authenticated clients and process or redirect them to a target foothold.
The main DB of siesta time will be in this server, with all information/configurations from the red team operators’ actions.

* Install in the installer's device host system dependencies (gcc,apache utils,...)
* Parse config<VPC>.txt file, use the parameters for creating a hive.tf (terraform plan)
* Download Go and the required dependencies to compile Hive.
* Create Hive sqlite DB
* Download Terraform and apply plan

In the same way, the hive.tf plan will (and this will be performed in the same host if Offline):

* Install Hive OS/Distro. dependencies
* Download go and their dependencies, to be able to compile Implants
* Download Terraform and terraform plugins
* Create ``/usr/local/STHive`` folder structure
* Upload OS dependencies, keys, sqlite DB, compiled Hive binary ...
* Configure Hive as a service





Online - AWS
----------------

**Steps to Prepare AWS Servers**

* Find EC2 Information
    * Prepare AWS key and credentials for target VPC
    * AccessKey/SecretKey
    * EC2 → “My Security Credentials” → “Access Keys”
    * `ami`_ 
    * `region`_ 
    * Create key pair on target region and Download “.pem” key


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

* Run 

.. prompt:: bash $

    ./hive.sh installaws

Offline
-----------
Offline option let operators to install hive in a target host without the use of terraform or any kind of VPC credentials.

``./hive.sh installOffline <IP> <Port> <targetFolder> <adminUsername> <adminPassword>``

.. prompt:: bash $

    ./hive.sh installOffline 0.0.0.0 6232 /usr/local/ admin admin



.. note::
    Every installing option comes with a "No Darwin" version of it. This will let hive to work without the need of Darwin dependencies (but loosing MacOSX implant abilities)

.. prompt:: bash $

    ./hive.sh installawsNoDarwin
    ./hive.sh installOfflineNoDarwin
    [...]

Uninstall
-----------

.. prompt:: bash $

    ./hive.sh remove

.. warning:: When installed Offline remove will not erase created/configured host data and packages


.. _ami: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/finding-an-ami.html
.. _region: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RegionsAndAvailabilityZones.html