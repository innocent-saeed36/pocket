Pocket - Tickets to Create

Persistence

- State hash follows ups: https://github.com/pokt-network/pocket/issues/361#issuecomment-1509193031
- Refactor persistence context (no heights in write, only heights in read)
  - Remove height from PersistenceWriteContext
- Dyllan’s ticket

P2P

- Actua peer discovery & churn
- Identity

Infrastructure

- Profiling node bootup

Consensus

- [Config] Block time configuration - pacemaker configs to set min block time in auto-mode
  - Hotstuff - configuration to slow down consensus
- Signature Aggregation
- State Sync optimizations
- Leader election
- Double Signing

Utility

- Generalize utility message types
  - Document message types
  - Generalize their creation & validation
  - Make it easy to add new ones
- Create actor specific protobufs
  - Share a base structure
- Validate message pool get cleared after each new block
- [Relay] Data structures
- [Servicer] Process
- Upgrades & feature flags
  - Activation height on-chain params. Similar to what we have in v0 if we need to add new flags in the future, but not affect the state hash starting from genesis.
- Utility
- [Relay] E2E Trustless Relay
- [Relay] E2E Delegated Relay
- [SPIKE] Prototype Relay Mining
  - Different difficulties for different chains
- (Jess) [SPIKE] Stake ownership
  - requirments for custodial, non custodial and revenue sharing algorithms
- [Future] Portal Validation
- [Future] [SPIKE] CRC - what do we have and need
- [Future][spike] How to get web sockets?
- Open Questions
  - How many chains can an app stake for?
  - How do we add new chains?
- Make fisherman staking permissioned

Validation

- Need to check for vulnerabilities in state storage bs state commitment

Tech Debt

- Change all addresses from `[]byte` to `string`
- Look through all TODO-like comments in the code
- Postgres - Avoid business logic at the postgres level
  - [SPIKE] Should we use an ORM?
- (state storage level) - make sure it's a depication of the tree (state commitment level)

Documentation

- Questions to answer:
  - What are all the types / protos we have?
  - What are all the message types we have?
  - How do I create a new module?
