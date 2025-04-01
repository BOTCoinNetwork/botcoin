.. _join_rst:

Joining a Network
=================

This section describes how to join an existing network that is already running,
such as the one created in :ref:`quickstart_rst`.

Here's a summary of the steps required to join an existing network built with
the botcoin:

.. code::

    $ botcoin keys new node1
    $ botcoin config pull [address]:[port] --key node1
    $ botcoincli poa nominate -h [address] -p [port] --from [node1 address] --pwd [password file for node1 key] --moniker node1 [node1 address]

    # wait to be accepted in the whitelist, which can be checked with
    $ monetcli poa whitelist
    # or
    $ botcoincli poa nomineelist

    $ botcoin run

Where [address] and [port] correspond to the endpoint of an existing peer in
the network.

**This scenario is designed to be run on a machine other than the one that is
running the existing node.**

Create An Account
-----------------

We need to generate a new key-pair for our account:

.. code:: bash

    $ botcoin keys new node1
    Passphrase:
    Repeat passphrase:
    Address: 0x5a735fC1235ce1E60eb5f9B9BCacb643a9Da27F4

Pull the Configuration From an Existing Node
---------------------------------------------

We now pull the ``botcoin`` configuration files from an existing peer. The
syntax for this command is:

.. code:: bash

    $ botcoin config pull [peer] [--key] [--address]

The peer parameter is the address/IP of an existing node on the network. The
network's configuration is requested from this peer. If the address does not
specify a port, the default API port (8080) is assumed.

We also need to specify the IP address of our own node. For a live network that
would clearly be a public IP address, but for an exploratory testnet, we would
recommend using an internal IP address. On Linux ``ifconfig`` will give you IP
address information. This can be set by using the --address flag. If not
specified ``botcoin`` will pick the first non-loopback address.

The ``--key`` parameter specifies the keyfile to use by moniker.

Thus we need to run the following command, but replace ``192.168.1.5:8080``
with the endpoint of the existing peer.

.. code:: bash

    $ botcoin config pull 192.168.1.5:8080 --key node1

Apply to Join the Network
-------------------------

If we tried to run ``botcoin`` at this stage, it would not be allowed to join
the other node because it isn't whitelisted yet. So we need to apply to the
whitelist first.

We do so with the ``botcoincli poa nominate`` command. The syntax is:

.. code:: bash

    $ botcoincli poa nominate -h <existing node> --from <moniker> --moniker <nominee moniker> --pwd <passphrase file> <nominee address>

But we can also do it interactively. **On the existing instance (node0), run
the following interactive ``botcoincli`` session**:

.. code:: bash

    botcoincli i
    __  __                          _        ____   _       ___
   |  \/  |   ___    _ __     ___  | |_     / ___| | |     |_ _|
   | |\/| |  / _ \  | '_ \   / _ \ | __|   | |     | |      | |
   | |  | | | (_) | | | | | |  __/ | |_    | |___  | |___   | |
   |_|  |_|  \___/  |_| |_|  \___|  \__|    \____| |_____| |___|

   Mode:        Interactive
   Data Dir:    /home/user/.monet
   Config File: /home/user/.monet/botcoincli.toml
   Keystore:    /home/user/.monet/keystore

    Commands:
     [...]


    botcoincli$ poa nominate
    ? From:  node0
    ? Passphrase:  [hidden]
    ? Nominee:  0x960c13654c477ac1d2d7f8fc7ae84d93a2225257
    ? Moniker:  node1

    You (0xa10aae5609643848ff1bceb76172652261db1d6c) nominated 'node1' (0x960c13654c477ac1d2d7f8fc7ae84d93a2225257)

    botcoincli$ poa nomineelist
    .------------------------------------------------------------------------------.
    | Moniker |                  Address                   | Up Votes | Down Votes |
    |---------|--------------------------------------------|----------|------------|
    | Node1   | 0x960c13654c477ac1d2d7f8fc7ae84d93a2225257 |        0 |          0 |
    '------------------------------------------------------------------------------'

Now that, we have applied to the whitelist (via node0), we need all the
entities in the current whitelist to vote for us. At the moment, only node0 is
in the whitelist, so let's cast a vote.

.. code:: bash

    botcoincli$ poa whitelist
    .------------------------------------------------------.
    | Moniker |                  Address                   |
    |---------|--------------------------------------------|
    | Node0   | 0xa10aae5609643848ff1bceb76172652261db1d6c |
    '------------------------------------------------------'

    botcoincli$ poa vote
    ? From:  node0
    ? Passphrase:  [hidden]
    ? Nominee:  0x960c13654c477ac1d2d7f8fc7ae84d93a2225257
    ? Verdict:  Yes
    You (0xa10aae5609643848ff1bceb76172652261db1d6c) voted 'Yes' for '0x960c13654c477ac1d2d7f8fc7ae84d93a2225257'.
    Election completed with the nominee being 'Accepted'.

    botcoin$ poa whitelist
    .------------------------------------------------------.
    | Moniker |                  Address                   |
    |---------|--------------------------------------------|
    | Node0   | 0xa10aae5609643848ff1bceb76172652261db1d6c |
    | Node1   | 0x960c13654c477ac1d2d7f8fc7ae84d93a2225257 |
    '------------------------------------------------------'

Finaly node1 made it into the whitelist.

Starting the Node
-----------------

To start node1, run the simple ``botcoin run`` command. You should be able see
the JoinRequest going through consensus, and being accepted by the PoA
contract.

.. code:: bash

    $ botcoin run
