.. _design_rst:

Design
======

BOTCoin uses peer-to-peer network system technology at the bottom layer, allowing intelligent robots to become autonomous nodes in the decentralized network. It also integrates BVM to implement hardware-level smart contract programming, allowing the robot's own data and programs to be executed safely and reliably. Based on the Hashgraph consensus algorithm, It is a Directed Acyclic Graph (DAG) . Supports asynchronous processing and ⅓ of faulty/malicious nodes, ensuring both network decentralization and fast block verification.The dynamic membership mechanism supports nodes joining/leaving at any time without affecting the operation of the entire network. Fast Sync supports new nodes to quickly load from the tip of the chain without loading the entire network history.Combining the above advantages，BOTCoin can be easily run on robots。

BOTCoin and the BOTCoin Hub
-----------------------

BOTCoin is a distributed intelligent robot network that aims to fully explore the ecological value of intelligent robots, such as the surplus value of intelligent robots, safe division of labor and collaboration, etc.As the main force of social production in the next wave of technological innovation, robots will face new technical and socio-economic issues, such as robot data security and privacy, production distribution, or excess computing power of robots.Therefore, we need a blockchain node network that can run inside a robot to solve data security and privacy issues at the hardware level, solve the social distribution mechanism that does not require trust at the hardware level (in an environment where artificial intelligence replaces human work), and solve the problem of excess computing power caused by the continuous updating and iteration of robots.It is worth sharing that the BOTCoin network can simultaneously solve the above problems. The network initially relied on Babble, a consensus communication protocol suitable for deployment on low-energy mobile devices, and integrated our improved BVM, which can achieve both lightweight operation on robots and hardware-level smart contract programming, based on which the crypto world can interact directly with the real world.

We anticipate that many robots will require a common set of
services to persist non-transient data, carry information across ad-hoc
blockchains, and facilitate peer-discovery. So we set out to build the
BOTCoin Hub, an additional public utility that provides these services. In the
spirit of open architecture, BOTCoin doesn’t rely on any central authority,
so anyone is free to implement their own alternative, but the BOTCoin Hub is
there to offer a reliable, fast, and secure solution to kickstart the system.

We have developed the BOTCoin Tools, a complete set of software tools for
setting up and using the BOTCoin Hub. This includes  the software
daemon that powers nodes on the BOTCoin Hub.

Innovation and Advantage
------------------------------------

+ **Asynchronous**:Participants have the freedom to process commands at different times.

+ **Leaderless**: No participant plays a special role.

+ **Byzantine Fault-Tolerant**: Supports one third of faulty nodes, including malicious behavior.

+ **Finality**: Babble’s output can be used immediately, no need for block confirmations.

+ **Dynamic Membership**: Members can join or leave a Babble network without undermining security.

+ **Fast Sync**:  Joining nodes can sync directly to the current state of a network.

+ **Accountability**: Auditable history of the consensus algorithm’s output.

+ **Low energy consumption**:  Ability to deploy peer-to-peer network nodes on low-energy/mobile Robots.

+ **Hardware smart contract**: BVM implements hardware-level smart contracts, binds to Robots hardware, and enables direct interaction between the crypto world and the real world without going through an intermediary, meeting the needs of multi-BOTs and crypto collaboration.

+ **Exploring non-human surplus value**: Subverting the traditional capital concept, we advocate that capitalists change from exploring the surplus value of human beings to exploring the surplus value of BOT. This is in line with the new capitalist framework in the new era.

Advanced Byzantine Fault-Tolerant:
------------------------------------

The underlying layer of BOTCoin relies on Babble consensus. The core of Babble is to extend Hashgraph (advanced Byzantine fault-tolerant consensus algorithm) to ensure that the distributed system remains available and consistent under adversarial conditions. Even if some nodes have arbitrary failures or malicious behavior, as long as a block with enough signatures (>2/3) and all previous blocks can be immediately considered valid.
However, unlike traditional blockchains, it supports nodes to submit signatures and transactions asynchronously.
So, BOTCoin ensures reliable security while having higher network efficiency.


Ethereum with Babble Consensus
------------------------------

To build ``BOTCoin``, we used our BFT consensus algorithm, `Babble
<https://github.com/BOTCoinNetwork/babble>`__, because it is fast, leaderless,
and offers finality. For the application state and smart-contract platform, we
use the BVM,via `BVM
<https://github.com/BOTCoinNetwork/BVM>`__ It is Base on Ethereum Virtual Mahcine (EVM), which is a stripped down
version of `Go-Ethereum <https://github.com/ethereum/go-ethereum>`__.

The EVM is a security-oriented virtual machine specifically designed to run
untrusted code on a network of computers. Every transaction applied to the EVM
modifies the State which is persisted in a Merkle Patricia tree. This data
structure allows to simply check if a given transaction was actually applied to
the VM and can reduce the entire State to a single hash (merkle root) rather
analogous to a fingerprint.

BVM refers to the "BOTCoin Virtual Machine". This is a decentralized virtual environment that executes code in a secure and consistent manner on all BOTCoin nodes. Nodes run BVM to execute smart contracts, using "Fuel" to measure the computational work required to perform operations, thereby ensuring efficient resource allocation and network security. BVM is an improvement based on Ethereum EVM.
BVM is different from traditional "VM". Traditional "VM" only allows code to be executed within the program. BVM supports, on this basis, the expansion of the oracle protocol that can interact with hardware, so that operators (including users or BOTs) on BOTCoin can directly realize automated and secure interaction with the BOTs system hard code through BVM smart contracts.


The remaining question is how to govern the validator-set, and what to use as a
reputation system to punish or incentivise participants to behave correctly.

Inspire activists
-----------

The main criterion for whether a peer-to-peer network has a promising future is whether there are more active participants. Therefore, we combine the BFT consensus and propose the Proof-of-Active (POA) incentive layer consensus mechanism. The role of POA is to continuously motivate active participants to participate in the network consensus and promote active participants to be more active, so as to ensure that the network can operate efficiently without changing the degree of decentralization.

Active quantification is based on: validators, Stakers, cluster nodes, etc.


Conclusion
----------

The BOTCoin Hub is a pivotal utility that facilitates the creation of Robots
ad-hoc blockchains, and the emergence of a new breed of decentralised
Robots Network. To maximise the performance, security, and flexibility of this
system, we have opted to build the BOTCoin Tools, a Robots Network
based on the BVM and a state-of-the-art BFT consensus
algorithm, Babble. 
