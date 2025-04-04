.. _botcoin_configuration_rst:

Botcoin Configuration
====================

All the configuration required to run a node is stored under a directory with a
specific structure. By default, ``botcoin`` will look for this directory in
``$HOME/.botcoin/botcoin-dir/`` [1]_ (on Linux), but it is possible to override 
this with the ``--config`` flag.

The directory must respect the following stucture:

::
  
  .botcoin/botcoin-config/
  ├── babble
  │   ├── peers.genesis.json
  │   ├── peers.json
  │   └── priv_key
  ├── eth
  │   ├── genesis.json
  └── monetd.toml

You would not normally need to access these configuration files directly. The
``botcoin config`` tool provides a CLI interfaces to set up a network. The
command ``botcoin config location`` provides further details of the default 
filepaths used for your instance.

Eth
---

The ``eth/genesis.json`` file defines prefunded accounts in the state, as well
as the POA smart-contract. This file is useful to predefine a set of accounts
that own all the initial tokens at the inception of the network. In addition,
the ``poa`` section contains information about the POA smart-contract.

Example ``genesis.json`` defining one prefunded account (the ABI and bytecode
of the smart-contract have been truncated):

.. code:: json

    {
        "alloc": {
                "a10aae5609643848ff1bceb76172652261db1d6c": {
                        "balance": "1234567890000000000000",
                        "moniker": "node0"
                }
        },
        "poa": {
                "address": "0xaBBAABbaaBbAABbaABbAABbAABbaAbbaaBbaaBBa",
                "abi": "[\n\t{\n\t\t\"constant\": true, ... ]",
                "code": "6080604052600436106101095760003560e01c8063..."
                }
    }

Babble
------

-  **babble/genesis.peers.json**: defines Babble's initial peer-set.

-  **babble/peers.json**: defines Babble's current peer-set

-  **babble/priv\_key**: contains the validator's private key for Babble.

Run Options
-----------

Options pertaining to the operation of the node are read from the 
<config>/monetd.toml file, or overwritten by the following flags:

::

    Flags:
          --api-listen string           IP:PORT of HTTP API service (default ":8080")
          --babble.advertise string     advertise IP:PORT of Babble node
          --babble.bootstrap            bootstrap Babble from database (default true)
          --babble.cache-size int       number of items in LRU caches (default 50000)
          --babble.heartbeat duration   heartbeat timer milliseconds (time between gossips) (default 200ms)
          --babble.listen string        bind IP:PORT of Babble node (default "192.168.0.13:1337")
          --babble.maintenance-mode     start babble in suspended (non-gossipping) state
          --babble.max-pool int         max number of pool connections (default 2)
          --babble.moniker string       friendly name
          --babble.suspend-limit int    number of undetermined-events since last run that will trigger automatic suspension (default 300)
          --babble.sync-limit int       max number of Events per sync (default 1000)
          --babble.timeout duration     TCP timeout milliseconds (default 1s)
      -c, --config string               configuration directory (default "/home/martin/.monet/monetd-config")
      -d, --data string                 data directory (default "/home/martin/.monet/monetd-data")
          --eth.cache int               megabytes of memory allocated to internal caching (min 16MB / database forced) (default 128)
          --eth.min-gas-price string    minimum gasprice of transactions submitted through this node (ex 1K, 1M, 1G, etc.) (default "0")
      -h, --help                        help for run
    
    Global Flags:
      -v, --verbose   verbose output

Example of a monet.toml file:

::

  verbose = "false"
  api-listen = ":8080"

  [babble]
    listen = "192.168.1.3:1337"
    advertise = "mypublicaddress:1337"
    heartbeat = "500ms"
    timeout = "1s"
    cache-size = 50000
    sync-limit = 1000
    max-pool = 2
    bootstrap = false

  [eth]
    cache = 128


.. [1] This location is for Linux instances. Mac and Windows uses a different
       path. The path for your instance can be ascertained with this command:
       ``botcoin config location``
