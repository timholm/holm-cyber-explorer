# STAGE 3: OPERATIONAL DOCTRINE -- FEDERATION OPERATIONS

## Domain 17: Scaling & Federation -- Operational Manuals D17-002 through D17-006

**Document ID:** STAGE3-FEDERATION-OPS
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Classification:** Operational Doctrine -- Each article is a practical operational manual governing specific federation procedures.

---

## How to Read This Document

This document contains five operational doctrine articles for federation. These are the Stage 3 manuals that translate the philosophy of D17-001 (Federation Philosophy) into concrete, step-by-step procedures. Where D17-001 establishes what the institution believes about federation, these articles establish what the institution does.

These articles are written to be executed, not merely read. They contain procedures that must be followed exactly when performing federation operations. They contain checklists that must be completed. They contain failure modes that must be recognized. They are addressed to the operator who is standing at a workbench with hardware in front of them, preparing to commission a new node, synchronize with a partner, resolve a conflict, verify trust, or separate a departing node.

Read D17-001 first. These articles assume you have internalized the sovereignty principle, the partition-tolerance default, the trust model, and the no-central-authority constraint. Without that philosophical grounding, these procedures will seem arbitrary. With it, every step will trace back to a principle you already understand.

---

---

# D17-002 -- New Node Setup and Commissioning

**Document ID:** D17-002
**Domain:** 17 -- Scaling & Federation
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** D17-001, ETH-001, CON-001, SEC-001, OPS-001
**Depended Upon By:** D17-003, D17-004, D17-005, D17-006. All subsequent federation operations assume a node has been commissioned according to this procedure.

---

## 1. Purpose

This article defines the complete procedure for creating a new federated node from bare hardware to first successful synchronization. It is the birth certificate procedure for a node -- the process by which raw components become a sovereign participant in the institution's federation.

A node that is commissioned incorrectly carries that error for its entire operational life. A weak key generated at commissioning remains a weak key forever. A misconfigured identity propagates through every sync package the node ever sends. A hardware platform that cannot meet the institution's storage and integrity requirements will fail under load, corrupting data that may have taken years to accumulate. The commissioning procedure exists to prevent these errors by making the correct path explicit and the incorrect path recognizable.

This article is addressed to the person physically assembling hardware, installing software, and performing the commissioning ceremony. It assumes you have the technical competence to install an operating system, generate cryptographic keys, and verify checksums. It does not assume you have done this before in a federation context. Every step is explained, every decision is justified, and every verification is mandatory.

## 2. Scope

This article covers:

- Hardware requirements and selection criteria for a new node.
- Operating system and software installation procedures.
- Identity generation: the creation of the node's unique cryptographic identity.
- Initial configuration of institutional software and directory structures.
- The commissioning ceremony: the formal process that transforms configured hardware into a recognized node.
- Key exchange with existing federation nodes.
- First sync: the initial data exchange that populates the new node.
- Verification procedures that confirm the node is correctly commissioned.

This article does not cover:

- Ongoing node maintenance and operations (see OPS-001 and domain-specific operational articles).
- Hardware repair or replacement procedures (see infrastructure domain articles).
- The philosophical justification for federation (see D17-001).
- Trust establishment beyond the initial key exchange (see D17-005).
- Content creation or importation on the new node (see Domain 18).

## 3. Background

### 3.1 Why Commissioning Matters

Every federation failure mode described in D17-001 -- drift, trust erosion, sovereignty capture, sync corruption -- can be seeded at commissioning. A node that begins with a poorly documented identity cannot be reliably verified later. A node that begins with incompatible software versions will produce sync packages that other nodes cannot process. A node that begins without a proper key exchange has no cryptographic basis for trust.

Commissioning is the moment of maximum control. The operator has the hardware in front of them, the software in hand, and no legacy data to constrain their choices. Every subsequent operation is an amendment to this initial state. Getting the initial state right is therefore the highest-leverage activity in the entire lifecycle of a node.

### 3.2 The Air-Gap Constraint on Commissioning

In a connected system, commissioning a new node typically involves downloading software from the internet, pulling configuration from a central server, and registering with a coordination service. None of these are available here. Every piece of software must be carried to the new node on physical media. Every configuration parameter must be set manually or imported from a prepared configuration package. Identity registration happens through physical key exchange, not through an online directory.

This makes commissioning slower and more deliberate than in a connected system. That deliberateness is a feature. Every step requires human attention. Every step produces an artifact -- a log entry, a checksum, a signed certificate -- that can be audited later. The commissioning process is, by design, a procedure that cannot be rushed.

### 3.3 The Commissioning Ceremony

The term "ceremony" is used deliberately. This is not merely a technical installation. It is a formal event that creates a new sovereign entity within the federation. The ceremony has witnesses (at minimum, the operator of the new node and one operator of an existing node). It produces a permanent record. It is dated and signed. The formality serves two purposes: it ensures that no step is skipped, and it creates institutional memory of the moment the node came into existence.

## 4. System Model

The commissioning process operates in four phases, each of which must be completed before the next begins:

**Phase 1: Hardware Preparation.**
Raw hardware is assembled, tested, and validated against the institution's hardware requirements. The hardware becomes a platform but is not yet a node.

**Phase 2: Software Installation.**
The operating system and institutional software are installed from verified media. The software environment is configured. The platform becomes a system but is not yet a node.

**Phase 3: Identity Generation.**
The node's unique cryptographic identity is generated. Keys are created, the node identifier is assigned, and the identity certificate is produced. The system becomes an entity but is not yet a recognized node.

**Phase 4: Commissioning Ceremony.**
The new entity is formally introduced to the federation. Keys are exchanged with existing nodes. The commissioning record is created and signed. The first sync is performed and verified. The entity becomes a commissioned node.

### 4.1 Hardware Requirements

The following are minimum hardware requirements for a federation node. These requirements exist to ensure that the node can fulfill its obligations as a federation participant for a minimum of ten years without hardware-imposed limitations.

**Storage:** Minimum 2 TB of primary storage. The storage must be durable (SSD preferred over spinning disk for primary storage; spinning disk acceptable for archival storage). The storage must support the institution's integrity-checking requirements -- specifically, the filesystem must support checksums or the institution must implement application-level checksums. A second storage device of equal or greater capacity is required for backup.

**Processing:** The processor must be capable of generating cryptographic hashes for integrity verification within reasonable time. Specifically, it must be able to SHA-256 hash a 1 GB file in under 60 seconds. This is a low bar that any modern processor meets, but it establishes a floor.

**Memory:** Minimum 8 GB RAM. This is sufficient for the institutional software, the operating system, and the integrity-verification processes to operate simultaneously without swapping.

**Physical Security:** The hardware must support physical access control. At minimum, it must have a lockable enclosure or be storable in a locked space. If the hardware is a general-purpose computer, it must support full-disk encryption with a passphrase known only to the node's operator.

**Removable Media Interface:** The hardware must have at least one interface for removable storage media (USB, or equivalent). This is the sync channel -- the physical interface through which sync packages are exchanged.

**Power:** The hardware should be operable from the node's primary power source. If the node is off-grid, the hardware's power requirements must be compatible with the available generation and storage capacity. Power consumption during active operation and idle state must be documented.

### 4.2 Software Stack

The institutional software stack consists of three layers:

**Layer 1: Operating System.** A stable, well-supported, open-source operating system. The specific distribution is not mandated, but it must be one for which security patches are available for a minimum of five years and which the node operator can maintain without internet access. The operating system installation image must be verified against a known-good hash before installation.

**Layer 2: Institutional Core.** The core institutional software that manages content storage, indexing, integrity verification, and sync package creation and consumption. This software is distributed as part of the institution's software archive and is versioned according to the institution's versioning policy.

**Layer 3: Institutional Tools.** Ancillary tools for content creation, format conversion, search, and reporting. These are useful but not essential for federation. A node can participate in federation with only Layers 1 and 2.

## 5. Rules & Constraints

- **R-D17-002-01: Hardware Must Be Verified Before Software Installation.** All hardware components must pass a burn-in test (minimum 48 hours of continuous operation under load) before the operating system is installed. Hardware failures discovered after commissioning are far more costly than hardware failures discovered during testing.

- **R-D17-002-02: All Software Must Be Installed from Verified Media.** No software may be installed from unverified sources. The installation media must have a documented chain of custody and must be integrity-checked (hash-verified) before use. The hashes of all installed software must be recorded in the commissioning log.

- **R-D17-002-03: Identity Keys Must Be Generated on the Node Itself.** The node's cryptographic identity keys must be generated on the node's own hardware using the node's own entropy sources. Keys must never be generated on another machine and transferred. The key generation process must use a cryptographically secure random number generator. The private key must never leave the node.

- **R-D17-002-04: The Commissioning Ceremony Requires a Witness.** At least one operator of an existing federation node must witness the commissioning of the new node. If no existing federation node exists (i.e., this is the founding node), the commissioning witness is the operator's own future self -- the ceremony is performed, documented, and the record is preserved for future verification.

- **R-D17-002-05: The Commissioning Log Is Permanent.** Every step of the commissioning process must be recorded in a commissioning log. The log includes: date and time of each step, hardware serial numbers, software versions and hashes, key fingerprints, witness identities, and the outcome of each verification check. The commissioning log is a permanent institutional record that is never modified or deleted.

- **R-D17-002-06: First Sync Must Be Verified End-to-End.** The first synchronization between the new node and an existing node must be verified by comparing content on both nodes after the sync completes. This is the only way to confirm that the sync pipeline is functioning correctly.

- **R-D17-002-07: A Node Is Not Commissioned Until the Record Says So.** The commissioning is not complete until the commissioning record is signed by the operator and the witness, stored on the new node, and a copy is delivered to at least one existing federation node. Until this record exists, the node is hardware with software on it. It is not a node.

## 6. Failure Modes

- **Weak key generation.** The node's cryptographic keys are generated with insufficient entropy, making them guessable or derivable. This compromises every sync package the node will ever produce and every trust relationship it will ever establish. Mitigation: R-D17-002-03 requires on-node generation with a cryptographically secure random number generator. The key generation step must include an entropy verification (checking that the system's entropy pool is adequately seeded before key generation begins). If in doubt, perform physical entropy generation (keyboard input, disk timing jitter) before generating keys.

- **Unverified software installation.** The software installed on the node has been tampered with, contains malware, or is an incorrect version. The node operates on a corrupted foundation. Mitigation: R-D17-002-02 requires hash verification of all installation media. The hash values must be obtained from a separate, trusted channel -- not from the same media being verified.

- **Hardware failure after commissioning.** A hardware component that passed burn-in testing fails after the node is commissioned and populated with data. If the failure affects storage, data may be lost. Mitigation: R-D17-002-01 reduces this risk but cannot eliminate it. The backup storage requirement provides the recovery path. The first backup should be taken immediately after first sync.

- **Commissioning without witness.** The operator commissions the node alone, skipping the witness requirement. This removes the independent verification that the process was followed correctly. Mitigation: R-D17-002-04 makes the witness mandatory. A commissioning log that lacks a witness signature is flagged as incomplete, and the node's trust level with other nodes is reduced until the witness requirement is satisfied.

- **Identity collision.** Two nodes are inadvertently assigned the same identifier. This is extraordinarily unlikely with properly generated cryptographic identifiers but would be catastrophic if it occurred, as sync packages from either node would be indistinguishable. Mitigation: Node identifiers must be generated using a cryptographic hash of the node's public key, making collisions computationally infeasible.

- **First sync failure masking ongoing problems.** The first sync appears to succeed but produces subtle data discrepancies that are not detected because the verification was superficial. Mitigation: R-D17-002-06 requires end-to-end content comparison after first sync. This comparison must be exhaustive, not a spot check. Every file transferred must have its integrity verified on the receiving node.

## 7. Recovery Procedures

1. **If hardware fails burn-in testing:** Replace the failed component. Restart the burn-in test from the beginning with the new component. Do not attempt to salvage partial burn-in results.

2. **If software hash verification fails:** Do not install the software. Obtain a fresh copy from a trusted source. If no trusted source is available, do not proceed with commissioning. A node built on unverified software cannot be trusted.

3. **If key generation produces a weak key (detected by key quality checks):** Delete the key material. Investigate the entropy source. Reseed the entropy pool. Generate new keys. Do not reuse or "strengthen" a weak key.

4. **If first sync fails:** Diagnose the failure. Common causes include: incompatible software versions (check version metadata in the sync package header), corrupted transport media (re-verify the media's integrity), or misconfigured sync parameters (review the node's sync configuration against the federation's shared standards). After fixing the cause, retry the sync from the beginning. Do not attempt to resume a partial sync.

5. **If the commissioning ceremony is interrupted:** Do not resume from where you stopped. Return to the beginning of Phase 4. The ceremony is atomic -- it either completes fully or it has not happened. Partial commissioning is not a valid state.

6. **If the commissioning log is lost:** If the log is lost before the ceremony is complete, restart the entire commissioning from Phase 1. If the log is lost after the ceremony is complete, reconstruct it from the artifacts: key fingerprints, software hashes, and witness testimony. Document the reconstruction and note that the original log was lost. This reconstruction has reduced evidentiary value compared to the original log.

## 8. Evolution Path

- **Years 0-5:** The commissioning procedure will be performed rarely -- perhaps only once for the founding node and once or twice more for early federation partners. Each commissioning is a learning event. Document what worked, what was confusing, and what took longer than expected. Expect to revise this procedure after each commissioning.

- **Years 5-15:** The procedure should be stable. Hardware requirements may need updating as technology evolves (storage expectations increase, processor architectures change). The software stack will have been revised at least once. Ensure that commissioning documentation stays current with the actual software being deployed.

- **Years 15-30:** Hardware that was current at founding may be obsolete. The hardware requirements section will need revision. The commissioning procedure itself should be technology-neutral enough to survive hardware generations, but the specific thresholds (storage minimums, processing benchmarks) will need updating.

- **Years 30-50+:** The commissioning procedure may be performed by people who never met the founding operator. The procedure must be self-sufficient -- comprehensible and executable without oral tradition. If a future operator cannot commission a new node using only this document and the referenced materials, this document has failed.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
I commissioned the founding node without this document, because this document did not yet exist. The founding node's commissioning was therefore informal and undocumented. One of the first tasks I should undertake is to retroactively document the founding node's configuration to the extent possible -- its hardware specifications, software versions, key fingerprints, and the date of its creation. This retroactive documentation will be incomplete compared to what this procedure produces for future nodes, but it is better than nothing. The founding node's commissioning gap is a debt that should be acknowledged in the institutional memory.

The 48-hour burn-in requirement may feel excessive for consumer hardware that has already been factory-tested. It is not. Factory testing catches manufacturing defects. Burn-in testing catches infant mortality failures -- components that passed manufacturing QA but fail under sustained load. Forty-eight hours is a compromise between thoroughness and practicality. If your schedule allows, extend it to 72 hours.

## 10. References

- D17-001 -- Federation Philosophy (sovereignty principle; no-central-authority constraint; trust model)
- ETH-001 -- Ethical Foundations (Principle 1: Sovereignty; Principle 3: Transparency)
- CON-001 -- The Founding Mandate (air-gap mandate; hardware sovereignty)
- SEC-001 -- Threat Model and Security Philosophy (key management; cryptographic standards)
- OPS-001 -- Operations Philosophy (complexity budget; operational sustainability)
- D17-005 -- Trust Establishment and Verification (key exchange procedures referenced during commissioning)
- Domain 18 -- Import & Quarantine (quarantine procedures for software installation media)
- Domain 20 -- Institutional Memory (commissioning log as permanent record)

---

---

# D17-003 -- Air-Gapped Synchronization Procedures

**Document ID:** D17-003
**Domain:** 17 -- Scaling & Federation
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** D17-001, D17-002, SEC-001, OPS-001
**Depended Upon By:** D17-004 (conflict resolution depends on sync), D17-005 (trust verification occurs during sync), D17-006 (withdrawal requires final sync).

---

## 1. Purpose

This article defines the complete procedure for synchronizing data between federated nodes without any network connection. It is the circulatory system of the federation -- the process by which knowledge, decisions, and institutional memory flow between sovereign nodes via physical media.

Synchronization is the federation's most frequent and most vulnerable operation. Every sync is a border crossing. Every sync package is a diplomatic pouch: sealed by the sender, transported through uncontrolled space, and opened by the receiver under controlled conditions. The integrity of the federation depends on the integrity of this process. A single corrupted sync that goes undetected can propagate bad data across the entire federation. A single intercepted sync package can expose institutional knowledge to unauthorized parties.

This article provides the step-by-step protocol for creating, transporting, receiving, validating, and applying sync packages. It is written for the operator who is about to plug a USB drive into their node and either export data for another node or import data from one. Every step matters. No step may be skipped.

## 2. Scope

This article covers:

- Physical media preparation and selection for sync transport.
- Sync package creation: what is included, how it is structured, how it is sealed.
- Transport security: protecting the sync package during physical transit.
- Sync package reception and initial validation.
- Integrity verification: how to confirm the package has not been tampered with or corrupted.
- Applying a validated sync package to the receiving node.
- The complete sync protocol, step by step, for both the sending and receiving node.
- Failure handling: what to do when a sync package fails validation at any stage.

This article does not cover:

- The initial commissioning sync (see D17-002, Phase 4).
- Conflict resolution when a sync reveals divergent changes (see D17-004).
- Trust verification beyond what is needed for sync validation (see D17-005).
- The philosophical justification for air-gapped sync (see D17-001).
- Content quarantine procedures for non-federation imports (see Domain 18).

## 3. Background

### 3.1 The Sneakernet Reality

This institution synchronizes via sneakernet -- the physical transport of storage media between locations. This is not a temporary compromise until network connectivity becomes available. It is the permanent, designed, and preferred synchronization method. The air gap is a security feature, not a limitation. Every design decision in this protocol preserves the air gap.

Sneakernet synchronization has properties that network synchronization does not. It is high-bandwidth (a USB drive can carry terabytes), high-latency (transport takes hours, days, or weeks), and unidirectional per trip (the media travels from A to B; the return trip is a separate event). It is also physically observable -- the operator can see the media, hold it, lock it in a case, and hand it directly to the recipient. There is no packet-sniffing on a sneakernet. There is no man-in-the-middle attack that does not involve a physical man in the middle.

### 3.2 The Sync Package as Diplomatic Pouch

D17-001 introduced the sync package as the unit of communication between nodes. This article defines its concrete structure. A sync package is a self-contained, integrity-protected, sender-authenticated bundle that includes everything the receiving node needs to process it. The receiving node should not need to contact any other node, consult any external resource, or make any assumption about the package's contents that is not verifiable from the package itself.

The diplomatic pouch analogy is precise. A diplomatic pouch is sealed by the sender. It is transported through territory the sender does not control. It is opened by the recipient, who verifies the seal before examining the contents. If the seal is broken, the contents are suspect. If the seal is intact, the contents are treated as authentic communications from the sender. The sync package works the same way.

### 3.3 The Risks of Physical Transport

Physical media can be lost, stolen, damaged, or tampered with during transport. The sync protocol must account for all of these:

- **Loss:** The media never arrives. The sending node still has the data. The sync simply does not happen. Loss is an inconvenience, not a catastrophe, because no data exists exclusively on the transport media.
- **Theft:** An unauthorized party obtains the media. If the sync package is encrypted, the thief gains nothing. If it is not encrypted, the thief gains read access to whatever was in the package. Encryption is therefore mandatory for any sync package that leaves secure premises.
- **Damage:** The media is physically damaged during transport. The data may be partially or fully corrupted. The integrity verification process must detect this.
- **Tampering:** An adversary intercepts the media, modifies the contents, and allows it to continue to its destination. The cryptographic seal must detect this.

## 4. System Model

### 4.1 Sync Package Structure

A sync package consists of the following components, in this order:

**Header (plaintext).** Contains: the package format version, the sender's node identifier, the intended recipient's node identifier, a timestamp, and a human-readable description of the package contents. The header is plaintext so that the recipient can identify the package before decryption.

**Manifest (signed).** A complete list of every file in the package, with its path, size, content hash (SHA-256), and modification timestamp. The manifest is signed with the sender's private key.

**Payload (encrypted).** The actual data being synchronized: new or modified documents, metadata updates, institutional memory entries, governance decisions, and any other content included in the sync scope. The payload is encrypted with a key derived from the trust relationship between the sender and recipient (see D17-005 for key exchange procedures).

**Integrity Seal (signed hash).** A SHA-256 hash of the entire package (header + manifest + payload), signed with the sender's private key. This is the cryptographic seal that detects tampering.

### 4.2 The Sync Protocol

The sync protocol has two sides: the sending procedure and the receiving procedure. They are described separately because they happen at different times and places.

**Sending Procedure (7 steps):**

1. Prepare the transport media (see Section 4.3).
2. Determine the sync scope: what data has changed since the last sync with this recipient.
3. Build the manifest: catalog every item to be included.
4. Assemble the payload: collect all items into a single archive.
5. Encrypt the payload with the recipient-specific key.
6. Generate the integrity seal over the complete package.
7. Write the package to the transport media. Verify the write by reading it back and checking the seal.

**Receiving Procedure (9 steps):**

1. Receive the transport media. Do not connect it to the production system.
2. Connect the media to the quarantine system.
3. Read the header. Verify that this package is intended for this node.
4. Verify the integrity seal. If the seal check fails, stop. Do not proceed. (See Section 7 for failure handling.)
5. Verify the manifest signature. Confirm the sender's identity.
6. Decrypt the payload.
7. Verify each item in the payload against the manifest (compare hashes).
8. Apply the validated payload to the production system, following conflict detection procedures (see D17-004 if conflicts are detected).
9. Record the sync event in the node's institutional memory: date, sender, package identifier, number of items received, number of conflicts detected, and outcome.

### 4.3 Transport Media Preparation

Transport media must be prepared before use. The preparation procedure is:

1. Use only media that the node operator has purchased new and opened from sealed packaging, or media that has been securely erased (full overwrite, not quick format) since its last use.
2. Format the media with a simple, universally readable filesystem (FAT32 for packages under 4 GB, exFAT for larger packages).
3. Write-test the media: write a known pattern to the full capacity of the media, read it back, and compare. This detects bad sectors and media failures before they corrupt a sync package.
4. Label the media physically (a paper label or engraving, not an electronic label that can be spoofed) with: the preparing node's identifier, the preparation date, and a unique media serial number.

Media that has been used for sync transport may be reused after secure erasure and re-preparation. Media that has shown any sign of failure (bad sectors, read errors, physical damage) must be destroyed, not reused.

## 5. Rules & Constraints

- **R-D17-003-01: Sync Packages Must Be Encrypted.** Every sync package that will leave the physical security perimeter of the sending node must be encrypted. The encryption key must be specific to the sender-recipient pair. Unencrypted sync packages are permitted only when the transport media never leaves the sending node's secure premises (e.g., a local backup that happens to use the sync package format).

- **R-D17-003-02: Integrity Must Be Verified Before Processing.** The receiving node must verify the integrity seal and the manifest signature before decrypting or examining the payload. If integrity verification fails, the package is rejected in its entirety. No partial processing. No "salvaging what we can."

- **R-D17-003-03: Every Sync Is Logged.** Both the sending node and the receiving node must log every sync event. The sending log records: date, recipient, package identifier, scope (what was included), and transport media identifier. The receiving log records: date, sender, package identifier, integrity check result, items received, conflicts detected, and resolution status.

- **R-D17-003-04: Transport Media Is Never Trusted.** Transport media is treated as potentially hostile, even if it was prepared by the receiving node. The quarantine process applies to all incoming media without exception. This protects against both tampering during transport and the possibility that the media was swapped entirely.

- **R-D17-003-05: Sync Scope Is Explicit.** The sending node must explicitly define what is included in each sync package. There is no "send everything" default. The operator reviews and approves the sync scope before the package is built. This prevents accidental inclusion of sensitive data and ensures the operator understands what they are sharing.

- **R-D17-003-06: No Bidirectional Sync on a Single Media.** A transport medium carries a sync package in one direction only. The sending node writes; the receiving node reads. If the receiving node wishes to send data back, it prepares a separate sync package on separate media (or on the same media after secure erasure and re-preparation). This prevents cross-contamination and simplifies the integrity model.

## 6. Failure Modes

- **Integrity seal failure.** The hash or signature check fails on the receiving end. Causes include: media corruption during transport, tampering, bit rot, or a software bug in the package creation process. Response: reject the entire package. Notify the sender via the next available communication channel. Request retransmission. Log the failure with full details.

- **Manifest mismatch.** The manifest signature verifies, but individual items in the payload do not match their manifest hashes. This indicates partial corruption -- the package structure is intact but some content was damaged. Response: reject the entire package. Do not attempt to salvage uncorrupted items. The package is atomic; it either arrives whole or it is rejected whole.

- **Decryption failure.** The payload cannot be decrypted. Causes include: wrong recipient key (the package was encrypted for a different node), key compromise and rotation (the sender used an old key), or corruption of the encrypted payload. Response: verify the header to confirm the package was intended for this node. If it was, contact the sender to diagnose the key mismatch. Do not attempt repeated decryption with different keys -- this is not a password-guessing exercise.

- **Stale sync package.** A sync package arrives that is older than a package already received from the same sender. The data in it is a subset of data already applied. Response: log the receipt but do not apply the package. Stale packages are not harmful, but applying them may revert newer data. The manifest timestamp and the sync log together provide the information needed to detect staleness.

- **Media failure during read.** The transport media develops read errors while the receiving node is attempting to read the sync package. Response: abort the read. Do not retry from the failed media -- repeated reads of failing media can cause further degradation. Request retransmission from the sender on fresh media.

- **Scope leakage.** The sending operator inadvertently includes sensitive or private data in the sync scope. This is detected after the package has left the node. Response: notify the recipient immediately (via whatever communication channel is available) to destroy the media without processing the package. If the package has already been processed, the recipient must flag and quarantine the inadvertently shared data. Log the incident. Review and tighten the scope review procedure.

## 7. Recovery Procedures

1. **If a sync package fails integrity verification:** Do not attempt to process any part of the package. Record the failure in the sync log with: the package identifier, the type of failure (seal failure, manifest mismatch, etc.), and the date. Securely erase the transport media. Contact the sending node to request retransmission. The sending node should verify its copy of the package against its own records before retransmitting.

2. **If the transport media is lost in transit:** The sending node still has all the data. Prepare a new sync package on new media and send it. Update the sync log to record the loss. If the lost media contained an encrypted package, the risk is limited to the metadata in the plaintext header. If for any reason the package was unencrypted, treat the contents as potentially compromised and take appropriate action based on the sensitivity of the data.

3. **If the transport media is suspected of tampering:** Do not connect it to any system. If it has already been connected to the quarantine system, isolate the quarantine system and rebuild it from known-good media. Log the incident with full details. Request retransmission from the sender on new media via a different transport path if possible.

4. **If a sync package is accidentally applied twice:** The second application should be idempotent if the sync software is functioning correctly (duplicate content with identical hashes should be detected and skipped). If duplicate application has caused data corruption, restore from the most recent pre-sync backup and re-apply the package once.

5. **If the sending node discovers an error in a sync package after sending:** The sending node prepares a correction package. The correction package contains a retraction notice identifying the erroneous package by its identifier, a description of the error, and corrected data if applicable. The correction follows the same sync protocol as any other package.

## 8. Evolution Path

- **Years 0-5:** Sync will be infrequent and small. The protocol may feel over-engineered for transferring a few hundred files on a USB drive. It is not. The protocol's value is in the discipline it establishes, not in the volume it handles. Perform every step even when the package is trivially small.

- **Years 5-15:** Sync packages may grow significantly as nodes accumulate content. The protocol should be evaluated for efficiency -- whether the packaging format compresses well, whether incremental syncs are working correctly, and whether the integrity verification time is acceptable for larger packages.

- **Years 15-30:** Transport media technology will change. USB may be obsolete. The protocol is transport-agnostic by design -- it defines the package structure and the verification procedures, not the physical medium. When the medium changes, update Section 4.3. The rest of the protocol should not need modification.

- **Years 30-50+:** The protocol may be executed by operators who have never used the specific media technologies of the founding era. The protocol must describe what to verify (integrity, authenticity, completeness) rather than how to verify it with a specific tool. Tools change. Verification principles do not.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The temptation with sync is to make it easy. Easy is dangerous. Easy means skipping the integrity check because "it is just a USB drive and I prepared it myself." Easy means not encrypting because "the drive is just going from my desk to the other room." Easy means not logging because "I will remember what I synced." Every shortcut in sync is a debt that compounds. One skipped integrity check is the one that would have caught the corrupted file. One unencrypted transport is the one that gets lost. One unlogged sync is the one that causes an unresolvable conflict six months later.

The protocol is long because the stakes are high. Every step exists because its absence would create a specific vulnerability. If a future operator finds a step that serves no purpose, they should document why they believe it is unnecessary in this Commentary Section -- and then keep performing it until the community has reviewed and agreed to its removal.

## 10. References

- D17-001 -- Federation Philosophy (sync package concept; partition tolerance; trust model)
- D17-002 -- New Node Setup and Commissioning (first sync procedure; initial media preparation)
- D17-004 -- Conflict Resolution Between Nodes (invoked when sync reveals conflicts)
- D17-005 -- Trust Establishment and Verification (key exchange for sync encryption)
- SEC-001 -- Threat Model and Security Philosophy (encryption requirements; integrity verification)
- OPS-001 -- Operations Philosophy (logging requirements; operational tempo)
- Domain 18 -- Import & Quarantine (quarantine system used for incoming sync media)
- Domain 20 -- Institutional Memory (sync event logging; transport records)

---

---

# D17-004 -- Conflict Resolution Between Nodes

**Document ID:** D17-004
**Domain:** 17 -- Scaling & Federation
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** D17-001, D17-002, D17-003, SEC-001, OPS-001
**Depended Upon By:** D17-005 (trust is affected by conflict behavior), D17-006 (withdrawal may be triggered by irreconcilable conflict).

---

## 1. Purpose

This article defines how the institution detects, classifies, and resolves conflicting changes when federated nodes synchronize after a period of partition. Conflict is not a bug. It is a natural and expected consequence of sovereign nodes operating independently. Every conflict represents two nodes that each made a reasonable decision independently. The institution's task is not to prevent conflict -- which would require constraining sovereignty -- but to resolve it with transparency, fairness, and a complete record.

An unresolved conflict is institutional poison. It grows. Subsequent changes on both sides build upon the conflicting versions, deepening the divergence. What was a simple disagreement about one document becomes a structural incompatibility between two branches of institutional knowledge. Early detection and prompt resolution are therefore essential. But prompt does not mean hasty. Resolution must be deliberate, documented, and agreed upon by the affected nodes.

This article is addressed to the operator who has just run a sync and received a conflict report. It tells you what to do, step by step, from detection through resolution through documentation.

## 2. Scope

This article covers:

- Conflict detection: how conflicts are identified during synchronization.
- Conflict classification: the taxonomy of conflict types and their severity levels.
- Resolution strategies: the available approaches to resolving each class of conflict.
- Manual merge procedures: how to perform a manual merge when automated resolution is not possible.
- The conflict log: what is recorded and why.
- Prevention: how to reduce the frequency and severity of future conflicts.

This article does not cover:

- The sync process itself (see D17-003).
- Conflicts between a node's data and externally imported content (see Domain 18).
- Trust implications of repeated conflicts (see D17-005).
- The philosophical framework for inter-node disagreement (see D17-001, Section 4.3).

## 3. Background

### 3.1 Why Conflicts Are Inevitable

In a connected system with real-time synchronization, conflicts can be minimized through locking (preventing two users from editing the same resource simultaneously) or through operational transformation (reconciling changes in real time as they are made). Neither technique is available in an air-gapped federation. Nodes operate independently for days, weeks, or months between syncs. Each node's operator makes changes based on the best information available to them, which may not include changes made on other nodes.

The result is that two nodes may independently modify the same document, reassign the same identifier, update the same metadata field, or make governance decisions that contradict each other. These are not errors. They are the predictable consequence of sovereignty and partition tolerance. D17-001 established that partition is the normal state and sovereignty is non-negotiable. Conflicts are the price of those principles. This article is about paying that price efficiently and transparently.

### 3.2 The Conflict Spectrum

Not all conflicts are equal. Some are trivial (two nodes independently corrected the same typo in the same way -- the "conflict" is actually agreement). Some are substantive but resolvable (two nodes edited different sections of the same document -- a straightforward merge). Some are fundamental (two nodes made contradictory governance decisions that cannot both stand). The resolution strategy must match the severity.

### 3.3 The Partition Duration Problem

The longer two nodes are partitioned, the more conflicts accumulate and the deeper those conflicts become. A partition of days produces surface-level conflicts. A partition of months produces structural conflicts. A partition of years may produce nodes that have diverged so fundamentally that full reunification is impractical. The conflict resolution procedure must scale from the trivial case to the existential case.

## 4. System Model

### 4.1 Conflict Detection

Conflict detection occurs during Step 8 of the receiving procedure (D17-003, Section 4.2). When the validated payload is being applied to the receiving node, each incoming item is compared against the corresponding local item. A conflict exists when:

- The incoming item and the local item have both been modified since the last successful sync between these two nodes.
- The modifications are different (if both nodes made identical changes, this is a convergence, not a conflict, and is resolved automatically by accepting either version).

Conflict detection operates on three levels:

**File-level conflict.** The same file has been modified on both nodes. This is the most common conflict type.

**Metadata-level conflict.** The same metadata field (classification, tags, authorship, status) has been changed on both nodes, even if the file content is unchanged.

**Structural conflict.** A file has been moved, renamed, or deleted on one node and modified on the other. Or two nodes have created different files with the same identifier.

### 4.2 Conflict Classification

Each detected conflict is classified by severity:

**Class 1 -- Trivial.** The conflict can be resolved automatically with high confidence. Examples: identical changes on both sides (convergence), whitespace-only differences, metadata timestamp differences. Resolution: automatic, with logging.

**Class 2 -- Routine.** The conflict requires human review but follows a standard resolution pattern. Examples: different edits to different sections of the same document, metadata disagreements on non-critical fields. Resolution: manual merge following standard procedures, with logging.

**Class 3 -- Substantive.** The conflict involves contradictory content changes that require judgment to resolve. Examples: both nodes rewrote the same section differently, conflicting factual claims, incompatible structural reorganizations. Resolution: deliberate review by the operator, with documented rationale for the chosen resolution.

**Class 4 -- Governance.** The conflict involves contradictory governance decisions, policy changes, or institutional-level disagreements. Examples: one node ratified a policy change that the other node rejected, conflicting amendments to the same governance document. Resolution: inter-node negotiation. Cannot be resolved unilaterally.

### 4.3 Resolution Strategies

For each conflict class, the following resolution strategies are available:

**Accept Local.** Keep the local version, discard the incoming version. Appropriate when the local version is known to be more correct or more current based on information not available to the other node.

**Accept Remote.** Keep the incoming version, discard the local version. Appropriate when the remote version is known to be more correct or more current.

**Merge.** Combine changes from both versions into a single resolved version. Appropriate when both versions contain valid changes that can coexist.

**Fork.** Maintain both versions as separate items, each clearly labeled with its origin. Appropriate when the two versions represent legitimate differences in perspective or purpose that should be preserved rather than reconciled.

**Escalate.** Defer resolution and flag the conflict for inter-node negotiation. Appropriate for Class 4 governance conflicts and any conflict where unilateral resolution would be inappropriate.

**Defer.** Temporarily set the conflict aside to be resolved later. The local node continues with its local version, but the conflict remains flagged. Appropriate only when resolution requires information or resources not currently available. Deferred conflicts must have a resolution deadline.

### 4.4 The Manual Merge Procedure

When a conflict requires manual merging (Class 2 or Class 3), the procedure is:

1. Extract both versions: the local version and the incoming version.
2. Generate a diff between the two versions to identify all points of divergence.
3. For each point of divergence, determine which version is correct or whether both contain valid changes that should be combined.
4. Create the merged version, incorporating the chosen resolution for each divergence point.
5. Verify the merged version for internal consistency (does it make sense as a whole document, not just at the individual merge points?).
6. Record the merge decision for each divergence point in the conflict log, including the rationale.
7. The merged version replaces both the local and incoming versions. It becomes the canonical version on this node and will be included in the next sync package to the other node.

### 4.5 The Conflict Log

Every conflict, regardless of class or resolution strategy, is recorded in the conflict log. The conflict log is a permanent, append-only institutional record. Each entry contains:

- Conflict identifier (unique).
- Date of detection.
- Sender node identifier.
- Affected item (file path, document identifier, or metadata field).
- Conflict class (1-4).
- Description of the conflict.
- Resolution strategy chosen.
- Rationale for the chosen strategy.
- The resolved version's content hash.
- Operator who performed the resolution.
- Date of resolution.

The conflict log serves three purposes: it provides an audit trail (anyone can review past resolutions), it enables pattern detection (recurring conflicts indicate a systemic issue), and it creates institutional memory (future operators can learn from past resolution decisions).

## 5. Rules & Constraints

- **R-D17-004-01: No Silent Resolution.** Every conflict must produce a conflict log entry, even if the resolution is automatic. Silent conflict resolution -- where conflicts are detected and resolved without any record -- is forbidden. The institution cannot learn from conflicts it does not know about.

- **R-D17-004-02: Resolution Must Precede Propagation.** A conflict must be resolved before the resolved version is included in any outgoing sync package. Unresolved conflicts must not propagate through the federation. Sending an unresolved conflict to a third node transforms a bilateral disagreement into a multilateral one.

- **R-D17-004-03: Class 4 Conflicts Cannot Be Resolved Unilaterally.** Governance conflicts require the participation of all affected nodes. A single node may not unilaterally resolve a conflict that involves contradictory governance decisions. If the affected nodes cannot agree, the conflict remains open and each node operates under its own governance decision until agreement is reached or the federation formally recognizes the divergence.

- **R-D17-004-04: Deferred Conflicts Must Have Deadlines.** A conflict may be deferred, but not indefinitely. Every deferred conflict must have a resolution deadline recorded in the conflict log. If the deadline passes without resolution, the conflict is escalated to the next higher class.

- **R-D17-004-05: The Conflict Log Is Permanent.** Conflict log entries may not be modified or deleted. If a resolution is later found to be incorrect, a new entry is added documenting the correction. The original entry remains.

- **R-D17-004-06: Pattern Analysis Is Mandatory.** The operator must review the conflict log at least quarterly to identify recurring patterns. Recurring conflicts between the same pair of nodes, or recurring conflicts on the same document or topic, indicate a systemic issue that must be addressed at its source, not merely resolved one instance at a time.

## 6. Failure Modes

- **Conflict fatigue.** The volume of conflicts is so high that the operator begins resolving them carelessly, choosing "Accept Local" for everything just to clear the queue. Quality of resolution degrades. Divergence between nodes increases rather than decreasing. Mitigation: if conflict volume exceeds the operator's capacity, reduce sync frequency to allow more changes to stabilize before synchronization. Address the root cause (often a lack of coordination on which node is responsible for which content).

- **Governance deadlock.** Two nodes disagree on a governance question and neither will yield. The conflict remains open indefinitely, blocking related decisions. Mitigation: R-D17-004-04 requires deadlines. If the deadline passes without resolution, the nodes must either negotiate, fork the federation on the disputed point, or one node must withdraw from the federation on the specific issue (see D17-006).

- **Merge corruption.** A manual merge introduces an error -- a paragraph deleted accidentally, a metadata field set incorrectly, an inconsistency between sections. The merged version is worse than either of the originals. Mitigation: the verification step in the manual merge procedure (Section 4.4, Step 5) is designed to catch this. Both original versions must be preserved until the merged version has been verified and the conflict log entry has been written.

- **Cascading conflicts.** A conflict in a foundational document (e.g., a classification scheme or a naming convention) causes secondary conflicts in every document that depends on it. Resolving the secondary conflicts without first resolving the foundational conflict wastes effort. Mitigation: during conflict triage, identify foundational conflicts first and resolve them before attempting to resolve dependent conflicts.

- **Historical revisionism.** An operator resolves a conflict by rewriting history -- changing timestamps, altering authorship records, or modifying institutional memory to make it appear that the conflict never occurred. Mitigation: R-D17-004-05 makes the conflict log permanent and append-only. Institutional memory entries are similarly immutable. Revisionism requires falsifying multiple independent records and will be detected during audit.

## 7. Recovery Procedures

1. **If a merge was performed incorrectly:** Retrieve the original local version and the original incoming version from the pre-merge archive (both versions must be preserved until the next successful sync confirms the merge was accepted by the other node). Redo the merge. Write a new conflict log entry documenting the error and the correction. Include the corrected version in the next sync package.

2. **If a conflict was resolved as the wrong class:** Reclassify the conflict with a new conflict log entry explaining the reclassification. If the resolution strategy was inappropriate for the correct class (e.g., a Class 4 conflict was resolved unilaterally as if it were Class 2), reopen the conflict and apply the correct resolution process.

3. **If conflict fatigue has degraded resolution quality:** Halt synchronization temporarily. Review the conflict log for the most recent sync batch. For each conflict resolved as "Accept Local," verify that this was the correct resolution. Reopen any that were resolved incorrectly. Before resuming synchronization, address the root cause of the conflict volume.

4. **If a governance deadlock persists past its deadline:** Convene a formal inter-node discussion (via exchange of signed position documents transported on sync media). Each node states its position, its rationale, and its proposed resolution. If agreement is still not possible, the options are: accept the fork (each node follows its own governance decision on this point, with the divergence documented), or one node withdraws from the federation on this specific matter (following D17-006 procedures for partial scope reduction).

5. **If cascading conflicts are overwhelming:** Stop resolving secondary conflicts. Identify the foundational conflict from which they cascade. Resolve the foundation first. Then re-run conflict detection on the secondary items -- many of them may resolve automatically once the foundation is settled.

## 8. Evolution Path

- **Years 0-5:** Conflicts will be rare because the federation is small and the nodes are likely operated by people who communicate regularly. Use this period to establish the discipline of the conflict log and the habit of documenting resolution rationale, even for trivial conflicts.

- **Years 5-15:** As the federation grows and nodes become more independent, conflicts will increase in frequency and complexity. The Class 4 governance conflict will likely be encountered for the first time. The conflict resolution procedure will be tested by real disagreement, not just technical divergence.

- **Years 15-30:** The conflict log should be rich enough to reveal patterns. Analysis of recurring conflicts should drive changes in federation coordination practices -- clearer delineation of which node is authoritative for which content areas, better communication of planned changes before they are made, and more disciplined use of the fork strategy when genuine differences of perspective are valuable.

- **Years 30-50+:** The conflict resolution procedures may need to handle conflicts between nodes operated by people who have never met and who inherited their nodes from previous operators. The procedures must be self-sufficient -- executable without personal relationships between operators. The formality that may seem excessive in the founding years will prove essential when the personal connections of the founding era are gone.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The hardest part of conflict resolution is not the technical merge. It is the ego management. When your version of a document conflicts with another node's version, the instinct is to believe your version is correct. Sometimes it is. Sometimes it is not. The classification system exists to remove ego from the equation. A Class 2 conflict is resolved by examining the changes, not by asserting authority. A Class 4 conflict is resolved by negotiation, not by force.

I have deliberately excluded any notion of a "master node" that wins conflicts by default. D17-001 forbids central authority. This means every conflict is a negotiation between equals. This is slower than having a master, and it is right. Sovereignty means you do not always get your way. It also means no one else always gets theirs.

The fork strategy deserves special emphasis. Not every disagreement needs to be resolved in favor of one side. Sometimes two nodes have genuinely different needs or perspectives, and both versions should exist. The fork is not a failure of conflict resolution. It is a legitimate outcome that preserves the sovereignty of both nodes.

## 10. References

- D17-001 -- Federation Philosophy (sovereignty principle; partition tolerance; governance model)
- D17-003 -- Air-Gapped Synchronization Procedures (sync protocol; conflict detection during sync)
- D17-005 -- Trust Establishment and Verification (trust implications of conflict behavior)
- D17-006 -- Node Withdrawal and Data Separation (withdrawal as response to irreconcilable conflict)
- SEC-001 -- Threat Model and Security Philosophy (integrity of conflict log)
- OPS-001 -- Operations Philosophy (operational tempo; quarterly review cycle)
- GOV-001 -- Authority Model (governance decision records; inter-node authority)
- Domain 20 -- Institutional Memory (conflict log as permanent institutional record)

---

---

# D17-005 -- Trust Establishment and Verification

**Document ID:** D17-005
**Domain:** 17 -- Scaling & Federation
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** D17-001, D17-002, SEC-001, GOV-001
**Depended Upon By:** D17-003 (sync encryption depends on trust keys), D17-004 (trust context affects conflict resolution), D17-006 (trust revocation is part of withdrawal).

---

## 1. Purpose

This article defines how federated nodes establish, maintain, verify, and when necessary revoke trust. Trust is the foundation of every federation operation. Without trust, sync packages cannot be encrypted or authenticated. Without trust, conflict resolution has no basis for good faith. Without trust, the federation is not a community of cooperating sovereigns but a collection of isolated systems that happen to exchange data.

Trust in this institution is not an abstraction. It is a concrete, measurable, documented relationship between two specific nodes, governed by explicit terms, verified through specific procedures, and revocable through a defined process. It is not a feeling. It is not an assumption. It is a protocol.

This article is addressed to every operator who will participate in establishing trust with another node, who will perform trust verification during routine operations, or who will face the difficult decision of revoking trust in a compromised or misbehaving node.

## 2. Scope

This article covers:

- The key exchange ceremony: how two nodes establish cryptographic trust.
- Trust level definitions: the specific levels of trust and what each permits.
- Trust verification procedures: how trust is confirmed during routine operations.
- Trust renewal: how trust relationships are maintained over time.
- Trust revocation: how to revoke trust when a node is compromised or misbehaving.
- The trust audit: how to systematically review all trust relationships.
- Responding to suspected compromise of a trusted node.

This article does not cover:

- The initial key generation during node commissioning (see D17-002).
- The sync protocol that uses trust-derived keys (see D17-003).
- The philosophical framework for trust (see D17-001, Section 4.1).
- Import trust scoring for non-federation content (see Domain 18).

## 3. Background

### 3.1 Trust in an Air-Gapped Context

In a connected system, trust can be verified continuously through certificate authorities, online directories, and real-time revocation. None of these mechanisms are available here. Trust in this institution is established through physical ceremony, maintained through periodic verification, and revoked through explicit notification. When you trust another node, you trust it based on a key exchange you witnessed, a fingerprint you verified in person, and a record you signed with your own hand. This is slower than automated trust. It is also stronger -- it cannot be undermined by compromising a certificate authority or poisoning a DNS cache.

### 3.2 The Cost of Misplaced Trust

A trusted node that is compromised becomes a Trojan horse. Its sync packages pass integrity verification because they are signed with the correct key -- the key that the compromised node still possesses. Its encrypted payloads are decrypted without suspicion because the decryption key is derived from the trust relationship. Everything the compromised node sends is accepted as authentic, because from a cryptographic standpoint, it is authentic. The compromise is not in the cryptography. It is in the assumption that the key holder is still trustworthy.

This is why trust must be actively maintained, not passively assumed. A key exchange ceremony establishes trust at a point in time. Verification procedures confirm that the trust remains warranted. Revocation procedures provide the mechanism to withdraw trust when it is no longer warranted. All three are essential. Establishment without verification is naive. Verification without revocation capability is powerless.

### 3.3 Trust Is Not Transitive by Default

D17-001 established that trust is not transitive. If Node A trusts Node B, and Node B trusts Node C, Node A does not automatically trust Node C. This is a deliberate design decision. Transitive trust means that the most trusting node in the federation determines the trust boundary for all nodes. If Node B trusts everyone, and trust is transitive, then Node A effectively trusts everyone too -- regardless of Node A's own judgment. This violates sovereignty.

Trust transitivity can be explicitly established. If Node A reviews Node C's credentials and decides to trust it, that trust is direct and deliberate, not inherited from Node B's decision. The process for establishing direct trust is the same regardless of whether the introduction came through another node.

## 4. System Model

### 4.1 The Key Exchange Ceremony

The key exchange ceremony is the foundational act of trust establishment. It requires the physical co-presence of operators from both nodes (or their designated agents) and produces a documented, signed record. The ceremony proceeds as follows:

**Preparation (each node, independently):**

1. Generate a ceremony-specific exchange package containing: the node's public key, the node's identifier, the node's current software version, and a timestamp.
2. Compute the fingerprint of the public key (a human-readable representation of the key, typically the first 40 characters of the SHA-256 hash of the public key, formatted in groups of four characters).
3. Print the fingerprint on paper. This physical artifact is the verification medium.

**The Ceremony (both operators, co-present):**

4. Each operator presents their printed fingerprint to the other operator.
5. Each operator reads their own fingerprint aloud while the other operator verifies it against the printed copy.
6. Each operator connects the other's exchange media (USB drive containing the exchange package) to their quarantine system and extracts the public key.
7. Each operator computes the fingerprint of the extracted public key and verifies it against the printed fingerprint and the oral recitation from Step 5.
8. If all three representations match (printed, oral, computed), the key is accepted. If any discrepancy exists, the ceremony halts and the key is rejected.
9. Both operators sign a ceremony record documenting: date, location, the identities of both operators, the fingerprints of both keys, and the outcome (accepted or rejected).
10. Each operator stores the other's public key in their node's trust store, associated with the ceremony record.

**Post-Ceremony:**

11. Each node creates a test sync package encrypted for the other node.
12. The test packages are exchanged (on fresh media) and each node verifies it can decrypt the other's test package.
13. Successful decryption confirms that the key exchange is functionally complete.

### 4.2 Trust Level Definitions

Trust is not binary. The institution defines four trust levels, each permitting different federation operations:

**Level 0 -- Untrusted.** No trust relationship exists. The node's identity may be known, but no key exchange has occurred. No sync packages are exchanged. This is the default state for all nodes.

**Level 1 -- Recognized.** A key exchange ceremony has been completed, but the trust relationship is new and has not been verified through operational experience. Sync packages may be exchanged, but the receiving node applies heightened scrutiny to incoming data (longer quarantine, more detailed content review). Conflicts are escalated rather than auto-resolved.

**Level 2 -- Established.** The trust relationship has been verified through multiple successful sync cycles without integrity failures or trust violations. The receiving node applies standard scrutiny to incoming data. Automatic resolution of Class 1 conflicts is permitted.

**Level 3 -- Trusted Partner.** The highest level of operational trust. The trust relationship has a long track record of reliable, high-quality exchanges. The receiving node may apply expedited processing to incoming data (shorter quarantine, reduced manual review for routine content). The nodes may engage in coordinated governance activities.

Trust levels are assigned per-node and per-direction. Node A may be Level 3 for Node B while Node B is Level 2 for Node A. Trust is earned through behavior, not negotiated through agreement.

### 4.3 Trust Verification Procedures

Trust is verified through three mechanisms:

**Sync-Time Verification.** Every incoming sync package implicitly verifies trust: the integrity seal is checked against the sender's stored public key. If the seal verifies, the sender's identity is confirmed for this package. If the seal fails, a trust investigation is triggered.

**Periodic Key Confirmation.** At least annually, each trust relationship must be reconfirmed through a key verification. This does not require a full ceremony -- it requires the operators to verify their stored copy of the other node's public key against the key the other node is currently using. This can be done through sync package analysis (the signing key in the most recent sync package should match the stored key) or through direct fingerprint exchange (at a meeting, or via a signed physical letter).

**Trust Audit.** A comprehensive review of all trust relationships, performed at least annually. The trust audit examines: the number and trust level of all trusted nodes, the last verification date for each, the conflict history with each, any anomalies in sync behavior, and whether any trust level should be adjusted (up or down). The audit produces a dated record stored in institutional memory.

### 4.4 Trust Revocation

Trust revocation is the most consequential trust operation. It severs a relationship and must be executed carefully. The revocation procedure is:

1. **Decision.** The operator decides to revoke trust. The reason must be documented. Reasons include: confirmed or suspected key compromise, repeated trust violations (corrupted packages, bad-faith conflict resolution, governance violations), or the other node's withdrawal from the federation.

2. **Key Removal.** The revoked node's public key is removed from the active trust store and moved to the revoked key archive. The key is not deleted -- it is preserved for historical reference and for verifying old sync packages.

3. **Notification.** A revocation notice is prepared and sent to the revoked node via the most expedient available channel (a sync package to the revoked node if one is scheduled, a physical letter, or notification through a mutually trusted third node). The notice includes: the date of revocation, the reason, and whether re-establishment is possible.

4. **Federation Notification.** A revocation notice is also sent to all other trusted nodes, so they can make their own trust decisions based on this information. This notice is informational, not directive -- each node decides for itself whether to also revoke trust.

5. **Post-Revocation Cleanup.** Review all data received from the revoked node. If the revocation was due to suspected compromise, flag all recent data from that node for re-verification. If the revocation was due to behavioral issues (not compromise), the data already received may be retained.

6. **Record.** The complete revocation, including decision rationale, notification records, and any re-verification results, is recorded in the trust log.

## 5. Rules & Constraints

- **R-D17-005-01: Trust Requires Physical Verification.** Trust may only be established through a key exchange ceremony that includes physical verification of key fingerprints. Remote-only key exchange is not permitted for the initial trust establishment. This requirement may be relaxed for trust renewal (not establishment) if the trust relationship is Level 2 or above.

- **R-D17-005-02: Trust Levels Must Be Earned.** A new trust relationship always begins at Level 1 (Recognized). Elevation to higher levels requires demonstrated reliability over multiple sync cycles. A node may not be granted Level 2 or Level 3 at the time of initial key exchange, regardless of the operators' personal relationship.

- **R-D17-005-03: Annual Verification Is Mandatory.** Every trust relationship must be verified at least annually through key confirmation or trust audit. A trust relationship that has not been verified within 14 months is automatically downgraded by one level. A trust relationship that has not been verified within 24 months is automatically revoked.

- **R-D17-005-04: Revocation Is Unilateral.** Any node may revoke trust in any other node at any time, for any reason. Revocation does not require the consent of the revoked node or of any other node. The revoking node must document its reasons, but its decision is sovereign.

- **R-D17-005-05: Revocation Notification Is Mandatory.** When trust is revoked, the revoking node must notify both the revoked node and all other trusted nodes. Failure to notify is itself a trust violation -- other nodes deserve to know that a federation relationship has changed.

- **R-D17-005-06: Trust Is Not Transferable.** Trust in one node cannot be transferred or delegated to another node. If Node A trusts Node B, Node A cannot grant Node C trust "on behalf of" Node B. Node C must establish its own trust relationship with Node A through the standard ceremony.

- **R-D17-005-07: The Trust Log Is Permanent.** All trust events (establishment, level changes, verifications, revocations) must be recorded in the trust log. The trust log is append-only and permanent.

## 6. Failure Modes

- **Ceremony theater.** The key exchange ceremony is performed perfunctorily -- fingerprints are not actually verified, the oral recitation is mumbled through, the ceremony record is signed without reading. The ceremony provides the form of trust without the substance. Mitigation: the ceremony requires three independent verifications (printed, oral, computed). All three must match. The ceremony record must include the specific fingerprints, not just a statement that verification occurred.

- **Trust level inflation.** Over time, all nodes are elevated to Level 3 because the operators are friendly and see no reason to maintain distinctions. The trust levels become meaningless. Every node receives expedited processing. The heightened scrutiny of Level 1 is never applied. Mitigation: trust level elevation must be earned through operational history and documented in the trust log with specific evidence. The annual trust audit must justify each node's current level.

- **Revocation paralysis.** An operator suspects a trust problem but delays revocation because of the social consequences -- the other operator is a friend, a family member, or a long-standing partner. The compromised trust relationship continues to be used, potentially propagating bad data. Mitigation: revocation is a technical decision, not a social one. The criteria are documented in the trust log. If the criteria for revocation are met, revocation must proceed regardless of the relationship.

- **Key compromise without detection.** A node's private key is compromised (stolen, copied) without the node operator's knowledge. The compromised key is used by an adversary to create fraudulent sync packages that pass all verification checks. Mitigation: this is the hardest failure mode to detect. Periodic key confirmation helps (the adversary would need to sustain the deception across multiple verification events). Anomaly detection in sync content helps (unexpected changes, content that contradicts the operator's known positions). But detection is not guaranteed. This is why physical security of the node (D17-002) is the primary defense.

- **Stale revocation notification.** Trust is revoked, but the notification reaches other nodes weeks or months later due to the latency of air-gapped communication. During the gap, other nodes continue to trust the revoked node. Mitigation: revocation notifications should be flagged as urgent and transported with the highest available priority. Receiving nodes, upon learning of a revocation, should review any sync packages received from the revoked node since the revocation date.

## 7. Recovery Procedures

1. **If a key exchange ceremony fails (fingerprints do not match):** Do not accept the key. Do not attempt to "fix" the mismatch. The mismatch indicates either a transmission error (the wrong key was loaded onto the exchange media) or a more serious problem (the media was tampered with). Both operators should regenerate their exchange packages independently and attempt the ceremony again with fresh media.

2. **If a trust relationship has lapsed (not verified within 24 months):** The trust is automatically revoked per R-D17-005-03. To re-establish, perform a full key exchange ceremony as though the relationship were new. The re-established relationship begins at Level 1 regardless of its previous level.

3. **If key compromise is suspected:** Immediately revoke trust in the potentially compromised node. Generate a new key pair on your own node (because if you suspect the other node's key, you should also consider whether your key has been exposed through the compromised relationship). Notify all trusted nodes. After the compromised node's operator has secured their system and generated new keys, trust may be re-established through a full ceremony.

4. **If trust level inflation is detected during audit:** Downgrade trust levels to the level supported by actual operational evidence. Document the downgrade in the trust log. Communicate the downgrade to affected nodes. This is corrective, not punitive -- it restores the trust levels to where they should have been.

5. **If a revocation was made in error:** Trust cannot be "un-revoked." If the revoking operator determines that the revocation was a mistake, the path forward is to re-establish trust through a full ceremony. The revocation record remains in the trust log. A new ceremony record documents the re-establishment and notes that it follows an erroneous revocation.

## 8. Evolution Path

- **Years 0-5:** The first key exchange ceremonies will be between people who know each other well and trust each other personally. The formal ceremony may feel unnecessary. Perform it anyway. The ceremony's value is not in its immediate effect but in the precedent it establishes and the record it creates.

- **Years 5-15:** Trust relationships will be tested by real disagreements and real operational incidents. The trust levels will begin to diverge as some nodes prove more reliable than others. The trust audit will become a meaningful exercise rather than a formality.

- **Years 15-30:** Key exchange ceremonies may occur between operators who are strangers, introduced through existing federation members. The ceremony's formality will prove its value -- it provides a structured way to establish trust without relying on personal acquaintance.

- **Years 30-50+:** The trust infrastructure must survive operator succession. When a node's operator changes, the new operator inherits the node's keys but not its trust relationships. New trust relationships must be established through ceremonies with the new operator. The institution's trust model is between nodes-as-operated-by-specific-people, not between nodes in the abstract.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The key exchange ceremony is borrowed from the PGP key-signing tradition, adapted for an air-gapped context. The core insight of the PGP key-signing party -- that trust should be established through physical verification, not through institutional intermediaries -- is exactly right for this institution. We have no certificate authority. We have no central directory. We have people who can look each other in the eye, read a fingerprint aloud, and verify it against a computed value. This is primitive by the standards of modern cryptographic infrastructure. It is also robust against every attack that does not involve physically compromising the ceremony.

The trust level system is designed to counter the natural tendency to trust people we like. Personal affinity is not a security credential. A node operated by a close friend starts at Level 1, just like a node operated by a stranger. Trust is earned through operational behavior: clean sync packages, honest conflict resolution, reliable key maintenance. This may feel cold. It is prudent.

I am particularly concerned about revocation paralysis. Revoking trust in someone you know personally is one of the hardest things an operator can be asked to do. The procedure is deliberately mechanical -- criteria, documentation, notification -- to take the personal weight off the operator's shoulders. You are not rejecting a person. You are following a protocol.

## 10. References

- D17-001 -- Federation Philosophy (trust model; non-transitivity; sovereignty)
- D17-002 -- New Node Setup and Commissioning (initial key generation)
- D17-003 -- Air-Gapped Synchronization Procedures (sync-time trust verification)
- D17-004 -- Conflict Resolution Between Nodes (trust implications of conflict behavior)
- D17-006 -- Node Withdrawal and Data Separation (trust revocation during withdrawal)
- SEC-001 -- Threat Model and Security Philosophy (key management; cryptographic standards; threat model)
- GOV-001 -- Authority Model (governance authority in trust decisions)
- OPS-001 -- Operations Philosophy (annual review cycle; operational sustainability)
- Domain 20 -- Institutional Memory (trust log as permanent record; ceremony records)

---

---

# D17-006 -- Node Withdrawal and Data Separation

**Document ID:** D17-006
**Domain:** 17 -- Scaling & Federation
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** D17-001, D17-002, D17-003, D17-004, D17-005, SEC-001, GOV-001
**Depended Upon By:** Referenced by all Domain 17 articles as the terminal procedure for federation relationships.

---

## 1. Purpose

This article defines the complete procedure for a node's departure from the federation. Withdrawal is the federation's most emotionally charged operation but it must be the most procedurally disciplined. A withdrawal performed badly -- data left entangled, keys left active, trust relationships left unresolved -- creates a wound that festers in the federation's institutional memory and a liability that the departing node carries into its post-federation existence.

D17-001 established that withdrawal is always permitted (R-D17-07). This article defines how withdrawal is performed. The right to leave is unconditional. The obligation to leave cleanly is the price of that right.

This article covers voluntary withdrawal (a node choosing to leave), involuntary withdrawal (a node being ejected following trust revocation by all federation members), and the special case of re-admission (a previously withdrawn node returning to the federation). It is addressed to the operator preparing to withdraw, the operators of nodes remaining in the federation, and the future operator who may consider re-admission after a withdrawal event.

## 2. Scope

This article covers:

- The withdrawal decision: what triggers it and how it is communicated.
- Data separation: how shared and entangled data is divided.
- Key revocation: how the withdrawing node's keys are handled.
- Shared resource division: how jointly developed content, standards, and tools are allocated.
- The withdrawal protocol: the step-by-step procedure.
- Post-withdrawal verification: how both sides confirm the separation is complete.
- Re-admission procedures: how a withdrawn node may rejoin the federation.

This article does not cover:

- The philosophical right to withdraw (established in D17-001, R-D17-07).
- Ongoing single-node operations after withdrawal (see OPS-001).
- The emotional or interpersonal dimensions of withdrawal (these are real but outside the scope of a technical procedure).
- Dissolution of the entire federation (a special case where all nodes withdraw simultaneously, governed by the same procedures applied concurrently by all nodes).

## 3. Background

### 3.1 Why Nodes Withdraw

Nodes withdraw for many reasons. Some are positive: the node's purpose has been fulfilled, the operator is consolidating to fewer nodes, or the node is being replaced by a more capable successor. Some are negative: irreconcilable conflicts with other nodes, trust failures, governance disagreements, or the operator's loss of confidence in the federation's direction. Some are neutral: the operator is moving, changing circumstances, or simply choosing to operate independently.

The withdrawal procedure does not distinguish between these reasons. The procedure is the same whether the withdrawal is amicable or acrimonious, planned or sudden, permanent or potentially temporary. The reason for withdrawal is documented in the withdrawal record, but it does not change the procedure. This is deliberate. A procedure that varies based on the emotional temperature of the situation is a procedure that will be executed inconsistently.

### 3.2 The Data Entanglement Problem

During federation, data becomes entangled. Node A receives content from Node B. Node A modifies that content. Node B receives the modification and modifies it further. The resulting document is a collaborative product with contributions from both nodes. If Node A withdraws, who owns the document? What version does each node keep?

The answer depends on the terms of the trust relationship established in D17-005 and any additional data-sharing agreements made during the federation. The withdrawal procedure enforces these terms. Where the terms are silent or ambiguous, the default is generous: both nodes retain all data they possess at the time of withdrawal. The withdrawing node does not lose data it received from the federation. The remaining nodes do not lose data they received from the withdrawing node. Collaborative works exist in full on all nodes that participated in their creation.

### 3.3 The Key Lifecycle

When a node withdraws, its cryptographic relationship with the federation changes fundamentally. Its keys, which were the basis of trust, become the basis of separation. The withdrawing node's public key, previously used to verify its sync packages, must be reclassified from "active trusted key" to "historical key -- withdrawn node." The key is not destroyed -- it is needed to verify historical sync packages and institutional memory entries. But it must no longer be used for current operations.

The remaining nodes must also manage their own key exposure. Every key that was shared with the withdrawing node (encryption keys derived from bilateral trust relationships) must be rotated. The withdrawing node possessed these keys and continues to possess them after withdrawal. If the withdrawal is adversarial, those keys must be considered compromised from the moment of withdrawal.

## 4. System Model

### 4.1 The Withdrawal Protocol

The withdrawal protocol has five phases:

**Phase 1: Declaration.**
The withdrawing node communicates its intention to withdraw. The declaration includes: the date of the intended withdrawal, the reason (which may be as brief or as detailed as the operator chooses), and whether the operator is open to discussion that might resolve the underlying issue. The declaration is delivered as a signed document via sync package or physical letter.

Declaration does not commit the node to withdraw. It opens a window for discussion. If the underlying issue is resolvable, the node may choose to remain. If the declaration period passes (minimum 30 days for voluntary withdrawal; immediate for involuntary withdrawal following trust revocation) without resolution, the withdrawal proceeds.

**Phase 2: Data Audit.**
Both the withdrawing node and the remaining nodes audit their data holdings to identify:

- Data originated by the withdrawing node that exists on other nodes.
- Data originated by other nodes that exists on the withdrawing node.
- Collaborative data with contributions from both the withdrawing node and remaining nodes.
- Metadata, institutional memory entries, and governance records related to the withdrawing node.

The audit produces a data separation plan: a document that specifies, for each category of data, what happens to it during withdrawal.

**Phase 3: Final Synchronization.**
A final sync is performed between the withdrawing node and each remaining node with which it has a trust relationship. This final sync serves two purposes: it ensures all data is as current as possible before separation (no pending changes are lost), and it provides the vehicle for delivering the withdrawal record and the data separation plan.

**Phase 4: Separation.**
The actual separation is executed according to the data separation plan:

- Trust is revoked bilaterally. The withdrawing node revokes trust in all remaining nodes. Each remaining node revokes trust in the withdrawing node. The revocation procedure follows D17-005, Section 4.4.
- Encryption keys derived from now-revoked trust relationships are rotated on all remaining nodes. The withdrawing node may also rotate its keys, but this is at its discretion since it is no longer part of the federation.
- Any data that the data separation plan specifies should be deleted is deleted. Deletion is verified using the same integrity mechanisms used for sync verification.
- The withdrawal record is finalized: a document signed by the withdrawing node and acknowledged by each remaining node, recording the date of withdrawal, the data separation plan, the key revocation records, and any post-withdrawal agreements.

**Phase 5: Verification.**
After separation, both the withdrawing node and the remaining nodes verify:

- All trust revocations have been processed.
- All key rotations have been completed.
- All data deletions specified in the separation plan have been executed.
- The withdrawal record is complete and stored in institutional memory.
- No active references to the withdrawn node's keys remain in any trust store (the keys exist in the revoked key archive, not in the active store).

### 4.2 Data Separation Defaults

When the trust relationship or data-sharing agreement does not specify how data is handled on withdrawal, the following defaults apply:

- **Data originated by the withdrawing node:** The withdrawing node retains it. Remaining nodes also retain their copies. The originator cannot compel other nodes to delete copies they legally possess.
- **Data originated by remaining nodes:** The withdrawing node retains its copies. Remaining nodes retain their originals. The same principle applies in reverse.
- **Collaborative data:** All nodes that participated in the collaboration retain the current version. No node is required to delete collaborative works.
- **Institutional memory entries:** All nodes retain all institutional memory entries. These are historical records and are never deleted during withdrawal.
- **Governance records:** All nodes retain all governance records for the same reason.
- **Federation-specific metadata:** Metadata that is only meaningful in a federation context (sync logs, trust records, conflict logs referencing the withdrawing node) is retained for historical reference but flagged as pertaining to a withdrawn node.

### 4.3 Involuntary Withdrawal

Involuntary withdrawal occurs when all remaining nodes in the federation revoke trust in a specific node. This is not an expulsion -- no central authority has the power to expel. It is the cumulative effect of individual trust decisions that result in a node having no trusted relationships remaining.

The involuntary withdrawal procedure differs from voluntary withdrawal in two ways:

- There is no declaration period. Trust revocation is effective immediately.
- The final synchronization may not occur if the trust revocation was triggered by suspected compromise. In this case, the data audit and separation plan are executed unilaterally by the remaining nodes, without the participation of the withdrawn node.

The withdrawn node retains all data it possessed at the time of withdrawal. The remaining nodes retain all data they possessed. The data separation defaults apply.

### 4.4 Re-Admission

A previously withdrawn node may seek to rejoin the federation. Re-admission is not automatic. It is a new trust establishment process, treated as though the node were being commissioned for the first time from a federation perspective.

The re-admission procedure is:

1. The returning node declares its intention to rejoin, including an explanation of what has changed since withdrawal.
2. Each remaining node independently decides whether to engage in a key exchange ceremony with the returning node. No node is obligated to re-admit.
3. If any remaining node agrees, a full key exchange ceremony is performed (D17-005, Section 4.1). The returning node begins at Trust Level 1 regardless of its previous level.
4. A re-admission record is created documenting: the original withdrawal record, the reason for return, the new key exchange ceremony records, and any conditions attached to re-admission.
5. Synchronization resumes following the standard sync protocol. The first sync after re-admission is treated with the same care as a first sync after commissioning -- full end-to-end verification.

The returning node's historical records (from its previous federation membership) are not erased. Its conflict log, trust log, and institutional memory entries from the prior period remain intact. The re-admission creates a new chapter, not a blank slate.

## 5. Rules & Constraints

- **R-D17-006-01: Withdrawal Cannot Be Blocked.** No node or group of nodes may prevent a voluntary withdrawal. The withdrawing node's right to leave is unconditional (D17-001, R-D17-07). Other nodes may attempt to resolve the underlying issue during the declaration period, but they may not delay or obstruct the withdrawal process.

- **R-D17-006-02: Data Separation Must Be Agreed Before Separation.** The data separation plan must be completed and acknowledged by both sides before the separation phase begins. If agreement cannot be reached, the defaults in Section 4.2 apply.

- **R-D17-006-03: Final Sync Is Mandatory for Voluntary Withdrawal.** A withdrawing node must perform a final sync with each remaining trusted node. This ensures no pending changes are lost and the withdrawal record is distributed. The only exception is if a trusted node refuses to participate in the final sync, in which case the withdrawing node sends the withdrawal record via physical letter.

- **R-D17-006-04: Key Rotation After Withdrawal Is Mandatory.** All remaining nodes must rotate any encryption keys derived from the trust relationship with the withdrawing node. This rotation must be completed within 30 days of the withdrawal. Until rotation is complete, sync packages between remaining nodes must use temporary keys or must not contain data that would be sensitive if the withdrawn node still possessed the decryption key.

- **R-D17-006-05: The Withdrawal Record Is Permanent.** The withdrawal record is a permanent institutional memory entry on all nodes. It is never deleted, even if the node later returns through re-admission. The record includes: the withdrawal date, the reason, the data separation plan, the key revocation records, and the verification results.

- **R-D17-006-06: Re-Admission Starts at Level 1.** A returning node always begins at Trust Level 1, regardless of the trust level it held before withdrawal. Previous operational history is relevant context but does not substitute for earning trust anew.

- **R-D17-006-07: Involuntary Withdrawal Requires Universal Revocation.** A node is involuntarily withdrawn only when all remaining nodes have independently revoked trust. A node that retains even one trust relationship is still a federation member, even if other nodes have revoked trust.

## 6. Failure Modes

- **Acrimonious withdrawal without procedure.** The withdrawing operator, angry or hurt, simply stops communicating. No declaration, no data audit, no final sync, no key revocation. The node vanishes. The remaining federation is left with active trust records for a node that no longer participates, unrotated keys, and no withdrawal record. Mitigation: the remaining nodes execute the involuntary withdrawal procedure unilaterally. They revoke trust, rotate keys, and create a withdrawal record documenting the circumstances. The withdrawn node's data remains on any node that possesses it, governed by the separation defaults.

- **Data separation disagreement.** The withdrawing node and the remaining nodes cannot agree on the data separation plan. The withdrawing node wants remaining nodes to delete content; the remaining nodes refuse (or vice versa). Mitigation: the defaults in Section 4.2 are designed for exactly this case. When agreement fails, the defaults apply: everyone keeps what they have. This may not satisfy either party, but it is fair and unambiguous.

- **Delayed key rotation.** Remaining nodes fail to rotate keys within the 30-day window, leaving them vulnerable if the withdrawn node's operator is adversarial. Mitigation: R-D17-006-04 mandates rotation within 30 days. The trust audit (D17-005) must verify that rotation has occurred. If rotation has not occurred within 30 days, the lagging node must cease sync operations with other nodes until rotation is complete.

- **Re-admission without accountability.** A node withdraws due to trust violations, then seeks re-admission without addressing the violations. The returning operator expects to resume the previous trust level. Mitigation: R-D17-006-06 requires Level 1 restart. The re-admission record must include an explanation of what has changed. Each remaining node decides independently whether the explanation is satisfactory.

- **Emotional contagion.** One node's withdrawal triggers anxiety in other nodes, leading to a cascade of withdrawals. The federation unravels not because of systemic failure but because of panic. Mitigation: this is fundamentally a human problem, not a technical one. The procedure can help by ensuring that each withdrawal is documented with its specific reason, making it clear that one node's withdrawal does not imply a systemic problem. The declaration period also helps -- 30 days provides time for the remaining federation to stabilize and for potentially anxious operators to assess the situation rationally.

- **Ghost node.** A node withdraws but the remaining federation's systems continue to reference it -- in sync schedules, in trust stores (incorrectly left in active status), in documentation. The ghost node creates confusion: new operators encounter references to a node that does not respond. Mitigation: the verification phase (Phase 5) is designed to catch this. All references to the withdrawn node must be flagged or updated. The quarterly review should include a check for ghost node references.

## 7. Recovery Procedures

1. **If a withdrawal was performed without following the protocol:** Reconstruct what should have happened. Perform the data audit retroactively. Execute key revocation and rotation. Create a withdrawal record that documents both what happened and what should have happened. Flag the withdrawal as non-standard in institutional memory.

2. **If key rotation was not completed within 30 days:** Treat all sync packages sent during the un-rotated period as potentially compromised (they were encrypted with keys known to the withdrawn node). Perform key rotation immediately. Review the content of any sync packages sent during the gap for sensitivity. If sensitive content was transmitted, notify affected nodes.

3. **If a withdrawn node's data is discovered to have been corrupted before withdrawal:** The withdrawing node is no longer reachable through federation channels. If the data corruption is limited, the remaining nodes correct it locally and document the correction. If the corruption is extensive, the remaining nodes must assess whether data received from the withdrawn node during the suspected corruption period should be quarantined and re-verified.

4. **If re-admission is denied by all remaining nodes:** The returning node remains a single-node institution. It retains all its data. It may seek to establish trust with individual nodes on a bilateral basis, which is functionally distinct from "re-admission to the federation" but achieves a similar result.

5. **If the entire federation dissolves (all nodes withdraw simultaneously):** Each node executes the voluntary withdrawal procedure with respect to every other node. All key rotations are performed. All withdrawal records are created. Each node becomes a standalone institution. The federation's institutional memory is preserved on all former member nodes. The dissolution is documented as a specific event type in institutional memory.

## 8. Evolution Path

- **Years 0-5:** Withdrawal is likely a theoretical concern. The federation is new, nodes are few, and relationships are strong. Use this period to establish the data-sharing agreements and trust relationship terms that will govern eventual withdrawals. The agreements that seem unnecessary now will be essential later.

- **Years 5-15:** The first withdrawal may occur. It will test the procedure and the operators' ability to execute it with discipline. Expect the emotional dimension to be harder than the technical dimension. Document the experience thoroughly -- future withdrawals will benefit from the precedent.

- **Years 15-30:** Withdrawals and re-admissions may become periodic events as the federation matures. The procedure should be routine enough to execute without drama. The withdrawal record archive should be rich enough to inform future decisions about federation structure and trust terms.

- **Years 30-50+:** The federation may have seen multiple withdrawals, re-admissions, and even a complete dissolution and re-formation. The withdrawal records from earlier decades will serve as institutional memory -- examples of how the federation handled separation and what worked and what did not. Future operators will learn from these records. That is why the records must be thorough.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
I wrote the withdrawal procedure before the federation exists. This is like writing a prenuptial agreement before the first date. It feels premature and slightly morbid. It is neither. The time to define how separation works is before there is any acrimony, before there is entangled data, before there are relationships that make it hard to think clearly about procedure.

The most important design decision in this article is the data separation defaults. The default is generous: everyone keeps everything. This avoids the most acrimonious disputes (demands that the other side delete content) and preserves the most institutional knowledge (nothing is lost in the separation). It means a withdrawn node walks away whole, and the remaining federation continues whole. The total amount of knowledge in the world does not decrease when a node withdraws. That feels right.

The re-admission procedure is deliberately rigorous. A node that returns to the federation starts at Level 1, full stop. This is not punitive. It is realistic. Trust was broken or abandoned. Rebuilding it requires the same process as building it the first time. The returning operator may find this frustrating, especially if they left on good terms. The procedure does not distinguish. Trust Level 1 is not a judgment. It is a starting point.

## 10. References

- D17-001 -- Federation Philosophy (R-D17-07: withdrawal right; sovereignty principle; governance model)
- D17-002 -- New Node Setup and Commissioning (re-admission parallels commissioning procedures)
- D17-003 -- Air-Gapped Synchronization Procedures (final sync protocol)
- D17-004 -- Conflict Resolution Between Nodes (unresolved conflicts at withdrawal)
- D17-005 -- Trust Establishment and Verification (trust revocation; key management; re-admission key exchange)
- SEC-001 -- Threat Model and Security Philosophy (key rotation; post-compromise procedures)
- GOV-001 -- Authority Model (governance record preservation; decision authority during withdrawal)
- OPS-001 -- Operations Philosophy (operational continuity after withdrawal; complexity budget adjustment)
- Domain 20 -- Institutional Memory (withdrawal record as permanent institutional record)

---