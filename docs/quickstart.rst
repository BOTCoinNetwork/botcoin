.. _quickstart_rst:

Getting Started
===============

In this document we explain how to run a single node and how to use
``botcoincli`` to interact with it. In another section, we will explain how to
join an existing network. For details about any command, please refer to the
:ref:`specification<monetd_commands_rst>`.

Creating A Single Node Network
------------------------------

In short, run the following three commands to start a standalone node:

.. code::

    $ botcoin keys new node0
    $ botcoin config build node0
    $ botcoin run

The ``keys new`` command will prompt us for a password, and generate a new
encrypted keyfile in the default keystore ``~/.monet/keystore``. We identified
our key with the ``node0`` moniker.

The ``config build`` command takes our key, and generates a minimal network
configuration with a single validator node, and a prefunded account. Again, the
configuration is written to ``~/.monet``. [1]_

Finally, the ``run`` command starts a monetd node, which will default to using
the configuration files in ``~/.monet``. [1]_

Using monetcli
--------------

Let's use ``monetcli`` to query the newly created node. First of all, install
``monetcli`` with ``npm install -g monetcli``.

While ``monetd`` is still running, open another terminal and start ``monetcli``
in interactive mode:

.. code:: bash

    $monetcli i
     __  __                          _        ____   _       ___ 
    |  \/  |   ___    _ __     ___  | |_     / ___| | |     |_ _|
    | |\/| |  / _ \  | '_ \   / _ \ | __|   | |     | |      | | 
    | |  | | | (_) | | | | | |  __/ | |_    | |___  | |___   | | 
    |_|  |_|  \___/  |_| |_|  \___|  \__|    \____| |_____| |___|
                                                                 
    Mode:        Interactive
    Data Dir:    /home/user/.monet
    Config File: /home/user/.monet/monetcli.toml
    Keystore:    /home/user/.monet/keystore
   
     Change datadir by: $ monetcli --datadir [path] [command]
   
     Commands:
   
       exit                                 Exit Monet CLI
       help [command...]                    Provides help for a given command.
       accounts create [options] [moniker]  Creates an encrypted keypair locally
       accounts get [options] [address]     Fetches account details from a connected node
       accounts list [options]              List all accounts in the local keystore directory
       accounts update [options] [moniker]  Update passphrase for a local account
       accounts import [options] [moniker]  Import an encrypted keyfile to the keystore
       config set [options]                 Set values of the configuration inside the data directory
       config view [options]                Output current configuration file
       poa check [options] [address]        Check whether an address is on the whitelist
       poa nominate [options] [address]     Nominate an address to proceed to election
       poa nomineelist [options]            List nominees for a connected node
       poa vote [options] [address]         Vote for an nominee currently in election
       poa whitelist [options]              List whitelist entries for a connected node
       transfer [options]                   Initiate a transfer of token(s) to an address
       info [options]                       Display information about node
       version [options]                    Display current version of cli
       block [options] [block]              Display details of a block by index
       validators [options] [round]         Get validators by consensus round
       history [options]                    Show validator history
       clear [options]                      Clear output on screen

Type ``info`` to check the status of the node:

.. code::

    monetcli$ info
    monetcli http GET localhost:8080/info
    ┌────────────────────────┬────────────┐
    │ consensus_events       │ 0          │
    │ consensus_transactions │ 0          │
    │ events_per_second      │ 0.00       │
    │ id                     │ 2002112846 │
    │ last_block_index       │ -1         │
    │ last_consensus_round   │ nil        │
    │ last_peer_change       │ -1         │
    │ min_gas_price          │ 0          │
    │ moniker                │ node0      │
    │ num_peers              │ 1          │
    │ round_events           │ 0          │
    │ rounds_per_second      │ 0.00       │
    │ state                  │ Babbling   │
    │ sync_rate              │ 1.00       │
    │ transaction_pool       │ 0          │
    │ type                   │ babble     │
    │ undetermined_events    │ 0          │
    └────────────────────────┴────────────┘

Type ``accounts list`` to get a list of accounts in the keystore, and the
balance associated with them.

.. code::

    monetcli$ accounts list
    monetcli info keystore /home/user/.monet/keystore
    monetcli info node localhost:8080
    ┌─────────┬────────────────────────────────────────────┬─────────────┬───────┐
    │ Moniker │ Address                                    │ Balance     │ Nonce │
    ├─────────┼────────────────────────────────────────────┼─────────────┼───────┤
    │ node0   │ 0xa10aae5609643848fF1Bceb76172652261dB1d6c │ ~1234.5678T │ 0     │
    └─────────┴────────────────────────────────────────────┴─────────────┴───────┘

So we have a prefunded account. The same account is used as a validator in
Babble, and as a Tenom-holding account in the ledger. This is the same account,
node0, that we created in the previous steps, with the encrypted private key
residing in ``~/.monet/keystore``.

Now, let's create a new key using ``monetcli``, and transfer some tokens to it.

.. code:: bash

    monetcli$ accounts create
    ? Moniker:  node1
    ? Output Path:  /home/user/.monet/keystore
    ? Passphrase:  [hidden]
    ? Re-enter passphrase:  [hidden]
    monetcli info keystore /home/user/.monet/keystore
    {
        "version":3,
        "id":"89970faf-8754-468e-903c-c9d3248a08cc",
        "address":"960c13654c477ac1d2d7f8fc7ae84d93a2225257",
        "crypto":{
            "ciphertext":"7aac819c1bed442d77897b690e5c2f14416589c7bdd6bdd2b5df5d03584ce0ec",
            "cipherparams":{
                "iv":"3d15a67d76293c3b7123f2bde76ba120"
            },
            "cipher":"aes-128-ctr",
            "kdf":"scrypt",
            "kdfparams":{
                "dklen":32,
                "salt":"730dd67f175a77c9833a230e334719292cbb735607795b1b84484e3d04783510",
                "n":8192,
                "r":8,
                "p":1
            },
            "mac":"7535c31c277a698207d278cd1f1df90747463e390b822cfef7d2faf8f1fa1809"
        }
    }

Like ``monetd keys new`` this command created a new key and wrote the encrypted
keyfile in ~/.monet/keystore. Let's double check that the key was created:

.. code:: bash

    monetcli$ accounts list
    monetcli info keystore /home/user/.monet/keystore
    monetcli info node localhost:8080
    ┌─────────┬────────────────────────────────────────────┬─────────────┬───────┐
    │ Moniker │ Address                                    │ Balance     │ Nonce │
    ├─────────┼────────────────────────────────────────────┼─────────────┼───────┤
    │ node0   │ 0xa10aae5609643848fF1Bceb76172652261dB1d6c │ ~1234.5678T │ 0     │
    ├─────────┼────────────────────────────────────────────┼─────────────┼───────┤
    │ node1   │ 0x960c13654c477ac1d2d7f8fc7ae84d93a2225257 │ 0T          │ 0     │
    └─────────┴────────────────────────────────────────────┴─────────────┴───────┘

Now, let's transfer 100 tokens to it.

.. code:: bash

    monetcli$ transfer
    ? From:  node0 (1,234,567,890,000,000,000,000)
    ? Enter password:  [hidden]
    ? To 0x960c13654c477ac1d2d7f8fc7ae84d93a2225257
    ? Value:  100
    {
      "from": "0xa10aae5609643848fF1Bceb76172652261dB1d6c",
      "to": "0x960c13654c477ac1d2d7f8fc7ae84d93a2225257",
      "value": "100T",
      "gas": 1000000,
      "gasPrice": "0T"
    }
    Transaction fee: 0T
    ? Submit transaction Yes
    Transaction submitted successfully.


Finally, we can check the account balances again to verify the outcome of the
transfer:

.. code:: bash

    monetcli$ accounts list --exact
    monetcli info keystore /home/user/.monet/keystore
    monetcli info node localhost:8080
    ┌─────────┬────────────────────────────────────────────┬─────────────┬───────┐
    │ Moniker │ Address                                    │ Balance     │ Nonce │
    ├─────────┼────────────────────────────────────────────┼─────────────┼───────┤
    │ node0   │ 0xa10aae5609643848fF1Bceb76172652261dB1d6c │ 1134.56789T │ 1     │
    ├─────────┼────────────────────────────────────────────┼─────────────┼───────┤
    │ node1   │ 0x960c13654c477ac1d2d7f8fc7ae84d93a2225257 │ 100T        │ 0     │
    └─────────┴────────────────────────────────────────────┴─────────────┴───────┘

.. [1] This location is for Linux instances. Mac and Windows uses a different
       path. The path for your instance can be ascertain with this command:
       ``monetd config location``
