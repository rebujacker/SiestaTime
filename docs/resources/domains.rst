Domains
===========================

Domains let the operator to add domain resources that can be used in the C2 deployment. Operators can list them to see if they are in use by any infraestructure already with the attribute ``active``


List Domain Info/Remove: ``Domains/SaaS --> example.xyz``

Add New domain: ``Domains/SaaS --> Add Domain``


Domain Status/Remove
-----------------------

.. figure:: ../_static/images/addresources/domainstatus.png
    :align: center
    :figwidth: 600px
    :target: ../_static/images/addresources/domainstatus.png


Godaddy
----------------

Keys --> `godaddykeys`_


.. figure:: ../_static/images/addresources/godaddy.png
    :align: center
    :figwidth: 600px
    :target: ../_static/images/addresources/godaddy.png


Gmail
-----------

* How to get my gmail connected App Credentials?
    * Create gmail Account
    * Follow instructions --> `goquickstart`_
    * Put both ``credentials.json`` and ``token.json`` Strings



.. figure:: ../_static/images/addresources/gmail.png
    :align: center
    :figwidth: 600px
    :target: ../_static/images/addresources/gmail.png


.. note::
	You need to specify the following gmail app access ``google.ConfigFromJSON(b, gmail.GmailModifyScope)``

.. _godaddykeys: https://developer.godaddy.com/keys/
.. _goquickstart: https://developers.google.com/gmail/api/quickstart/go