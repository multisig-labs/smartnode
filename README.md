# Rocket Pool - Smart Node Package

<p align="center">
  <img src="https://raw.githubusercontent.com/rocket-pool/rocketpool/master/images/logo.png?raw=true" alt="Rocket Pool - Next Generation Decentralised Ethereum Proof-of-Stake (PoS) Infrastructure Service and Pool" width="500" />
</p>

---

`Rocket Pool 2.5` is a next generation Ethereum proof of stake (PoS) infrastructure service designed to be highly decentralised, distributed and compatible with Ethereum 2.0, the new consensus protocol that Ethereum will transition to in 2020.

Running a Rocket Pool smart node allows you to stake on Ethereum 2.0 with only 16 ETH, and earn a higher return than you would outside the network.

This repository contains the source code for:

* The Rocket Pool smart node client (CLI), which is used to manage a smart node either locally or remotely (over SSH)
* The Rocket Pool smart node service, which provides an API for client communication and performs background node tasks
* The Rocket Pool smart node PoW proxy, a simple Eth 1.0 client proxy server for communicating with remote Eth 1.0 services such as Infura

The smart node service and PoW proxy are designed to be run as part of a docker stack and generally do not need to be installed manually.
See the [Rocket Pool dockerhub](https://hub.docker.com/u/rocketpool) page for a complete list of docker images.


## Installation

See the [Smart Node Installer](https://github.com/rocket-pool/smartnode-install) repository for supported platforms and installation instructions.


## CLI Commands

The following commands are available via the smart node client:

- `rocketpool service install` - Install the Rocket Pool service either locally or to a remote server
- `rocketpool service config` - Configure the Rocket Pool service for use
- `rocketpool service status` - Display the current status of the Rocket Pool service
- `rocketpool service start` - Start the Rocket Pool service to begin running a smart node
- `rocketpool service pause` - Pause the Rocket Pool service temporarily
- `rocketpool service stop` - Pause the Rocket Pool service temporarily
- `rocketpool service terminate` - Terminate the Rocket Pool service and remove all associated docker containers & volumes
- `rocketpool service logs [services...]` - View the logs for one or more services running as part of the docker stack
- `rocketpool service stats` - Display resource usage statistics for the Rocket Pool service
- `rocketpool service version` - Display version information for the Rocket Pool client & service

- `rocketpool wallet status` - Display the current status of the node's wallet
- `rocketpool wallet init` - Initialize the node's password and wallet
- `rocketpool wallet recover` - Recover a node wallet from a mnemonic phrase
- `rocketpool wallet rebuild` - Rebuild validator keystores from derived keys
- `rocketpool wallet export` - Export the node's wallet information

- `rocketpool node status` - Display the current status of the node
- `rocketpool node register` - Register the node with the Rocket Pool network
- `rocketpool node set-withdrawal-address [address]` - Set the address which node rewards & refunds are sent to
- `rocketpool node set-timezone` - Update the node's timezone location
- `rocketpool node swap-ggp` - Swap old GGP tokens for new GGP
- `rocketpool node stake-ggp` - Stake GGP against the node to collateralize minipools
- `rocketpool node withdraw-ggp` - Withdraw GGP staked against the node
- `rocketpool node deposit` - Make a deposit to create a minipool and begin staking
- `rocketpool node send [amount] [token] [to]` - Send an amount of ETH or tokens to an address
- `rocketpool node burn [amount] [token]` - Burn reward tokens for ETH

- `rocketpool minipool status` - Display the current status of all minipools run by the node
- `rocketpool minipool refund` - Refund ETH from minipools which have had user-deposited ETH assigned to them
- `rocketpool minipool dissolve` - Dissolve initialized minipools and recover deposited ETH from them
- `rocketpool minipool exit` - Exit active minipool validators from the beacon chainand close them
- `rocketpool minipool close` - Close minipools which have timed out and been dissolved

- `rocketpool auction status` - Display the current status of the GGP auction contract and lots
- `rocketpool auction lots` - Display the details of all GGP lots
- `rocketpool auction create-lot` - Create a new GGP lot from GGP in the auction contract
- `rocketpool auction bid-lot` - Bid ETH on an active GGP lot
- `rocketpool auction claim-lot` - Clean GGP from a cleared lot you bid on
- `rocketpool auction recover-lot` - Recover unclaimed GGP from a cleared lot back to the auction contract

- `rocketpool odao status` - Display the current status of the oracle DAO
- `rocketpool odao members` - Display the details of all oracle DAO members
- `rocketpool odao proposals` - Display the details of all oracle DAO proposals
- `rocketpool odao propose-invite [address] [id] [url]` - Invite a member to join the oracle DAO
- `rocketpool odao propose-leave` - Propose leaving the oracle DAO
- `rocketpool odao propose-replace [address] [id] [url]` - Propose replacing your position in the oracle DAO with a new member
- `rocketpool odao propose-kick` - Propose kicking a member from the oracle DAO
- `rocketpool odao cancel-proposal` - Cancel a proposal you created
- `rocketpool odao vote-proposal` - Vote on a proposal
- `rocketpool odao execute-proposal` - Execute a passed proposal
- `rocketpool odao join` - Join the oracle DAO (requires an executed invite proposal)
- `rocketpool odao leave` - Leave the oracle DAO (requires an executed leave proposal)
- `rocketpool odao replace` - Replace your position in the oracle DAO (requires an executed replace proposal)

- `rocketpool network node-fee` - Display the current network node commission rate for new minipools
- `rocketpool network ggp-price` - Display the current network GGP price information

- `rocketpool queue status` - Display the current status of the deposit pool
- `rocketpool queue process` - Process the deposit pool by assigning user-deposited ETH to available minipools
