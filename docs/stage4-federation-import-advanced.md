# STAGE 4: SPECIALIZED SYSTEMS -- FEDERATION & IMPORT ADVANCED REFERENCE

## Domains 17-18: Scaling & Federation, Import & Quarantine -- Advanced Technical Specifications

**Document ID:** STAGE4-FEDERATION-IMPORT-ADVANCED
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Classification:** Advanced Reference Documents -- Each article is a detailed technical specification for specialized systems within its domain.

---

## How to Read This Document

This document contains five advanced reference articles spanning Domains 17 and 18. These are Stage 4 documents -- specialized systems articles that provide the deep technical detail required to implement, operate, and maintain the federation and import systems defined in earlier stages.

Stage 1 defined what these domains contain. Stage 2 established why they exist and what principles govern them. Stage 3 provided operational procedures for performing the work. This stage -- Stage 4 -- provides the engineering-level specifications that make the operational procedures concrete and unambiguous. Where Stage 3 says "build a sync package," this stage defines exactly what a sync package contains, byte by byte. Where Stage 3 says "scan for malware," this stage defines how scanning works without internet connectivity.

These articles assume you have read and internalized D17-001 (Federation Philosophy), D18-001 (Import & Quarantine Philosophy), and the Stage 3 operational manuals D17-002 through D17-006. Without that foundation, these specifications will seem arbitrary. With it, every field in every data structure will trace back to a principle you already understand.

---

---

# D17-007 -- Sync Package Specification

**Document ID:** D17-007
**Domain:** 17 -- Scaling & Federation
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** D17-001, D17-003, D17-005, SEC-001, OPS-001
**Depended Upon By:** D17-008 (governance sync uses sync packages), D18-003 (malware scanning applies to incoming sync packages), D18-005 (provenance tracking references sync package metadata). All federation operations that create or consume sync packages depend on this specification.

---

## 1. Purpose

This article defines the complete technical specification of a sync package -- the unit of data exchange between federated nodes. D17-001 introduced the sync package as a concept. D17-003 defined the operational procedures for creating, transporting, and receiving sync packages. This article specifies the exact format: the byte-level structure, the fields, the cryptographic requirements, the manifest schema, and the validation rules that every sync package must satisfy.

A sync package is a diplomatic pouch in digital form. It must be self-contained, self-verifying, and self-describing. The receiving node must be able to determine, from the package alone, who sent it, what it contains, whether it has been tampered with, and whether the sender had the authority to send it. No external lookup, no network call, no reference to any resource outside the package and the receiving node's own trust store may be required to validate a sync package.

This specification exists because ambiguity in a data exchange format is a vulnerability. If two nodes disagree about what a field means, data will be misinterpreted. If the integrity verification procedure is imprecise, corrupted data will pass validation on one node and fail on another. If the encryption scheme is underspecified, a future implementation may choose a weak algorithm. Every field, every algorithm, and every validation rule is specified here so that two independent implementations, written decades apart, can produce and consume interoperable sync packages.

## 2. Scope

This article covers:

- The sync package container format and its layered structure.
- The plaintext header specification: fields, encoding, and validation rules.
- The signed manifest specification: schema, content hash requirements, and signature format.
- The encrypted payload specification: encryption algorithm, key derivation, and payload structure.
- The integrity seal specification: hash algorithm, signature format, and verification procedure.
- Version negotiation between nodes running different software versions.
- The sync package lifecycle: creation, validation, application, and archival.
- Extensibility: how the format evolves without breaking backward compatibility.

This article does not cover:

- The operational procedures for sync (see D17-003).
- The physical transport of sync media (see D17-003, Section 4.3).
- The trust establishment that produces encryption keys (see D17-005).
- The conflict resolution that occurs after a sync package is applied (see D17-004).
- Content-level validation of items within the payload (see Domain 18 for imported content validation).

## 3. Background

### 3.1 Why a Formal Specification Matters

D17-003 describes the sync package at a structural level: header, manifest, payload, integrity seal. That description is sufficient for understanding the sync protocol. It is not sufficient for implementing it. Implementation requires precision that operational documentation does not provide. What character encoding does the header use? What hash algorithm produces the content hashes in the manifest? What happens when the manifest references a file whose path contains characters that are invalid on the receiving node's filesystem? These questions cannot be answered by philosophy or procedure. They require a specification.

The specification also serves as a contract between implementations. If two nodes are running different versions of the institutional software -- or, in the distant future, entirely different software that implements the same specification -- the sync package format is their shared language. The specification defines that language with enough precision that a conforming implementation can be written from this document alone.

### 3.2 Design Principles for the Format

The sync package format is governed by four design principles, each derived from the institution's philosophical foundations:

**Self-containment.** A sync package includes everything needed to validate and process it. This derives from R-D17-02 (Partition Is Normal) and R-D17-05 (Sync Packages Are Self-Contained). The receiving node may have been disconnected from all other nodes for years. It cannot consult any external resource. The package must be complete.

**Transparency.** The package format is fully documented. There are no proprietary components, no patent-encumbered algorithms, no fields whose meaning is known only to the software authors. This derives from ETH-001, Principle 3 (Transparency of Operation). A future operator must be able to read a sync package with nothing more than a hex editor and this specification.

**Integrity.** Every byte of the package is covered by the integrity seal. Every item in the payload is individually hashed in the manifest. Corruption at any level is detectable. This derives from SEC-001's integrity requirements.

**Longevity.** The format uses stable, well-understood algorithms and encodings. No algorithm is chosen for novelty. Every algorithm is chosen because it has a long track record of security analysis and because open-source implementations exist in multiple languages. This derives from ETH-001, Principle 4 (Longevity Over Novelty).

### 3.3 Relationship to Existing Descriptions

D17-003, Section 4.1 describes the sync package as consisting of four components: Header, Manifest, Payload, and Integrity Seal. This specification retains that four-component structure and elaborates each component to implementation-ready detail. Nothing in this specification contradicts D17-003. Everything in this specification extends it.

## 4. System Model

### 4.1 Package Container Format

A sync package is a single file with the extension `.syncpkg`. The file is a concatenation of four sections in fixed order, each preceded by a section delimiter that identifies the section type and its length. The overall structure is:

```
[Package Preamble]
[Section 1: Header]
[Section 2: Manifest]
[Section 3: Payload]
[Section 4: Integrity Seal]
```

**Package Preamble.** The first 64 bytes of the file. Contains:

- Bytes 0-7: Magic number `SYNCPKG\0` (ASCII, null-terminated). Identifies the file as a sync package.
- Bytes 8-11: Format version (uint32, big-endian). Current version: `0x00000001` (version 1).
- Bytes 12-15: Header section offset (uint32, big-endian). Byte offset from the start of the file to the beginning of the header section.
- Bytes 16-19: Header section length (uint32, big-endian).
- Bytes 20-23: Manifest section offset (uint32, big-endian).
- Bytes 24-27: Manifest section length (uint32, big-endian).
- Bytes 28-31: Payload section offset (uint32, big-endian).
- Bytes 32-35: Payload section length (uint32, big-endian).
- Bytes 36-39: Seal section offset (uint32, big-endian).
- Bytes 40-43: Seal section length (uint32, big-endian).
- Bytes 44-63: Reserved (zero-filled). Reserved for future format versions.

The preamble is always plaintext and unencrypted. It allows the receiving node to locate each section without parsing the entire file.

### 4.2 Header Specification

The header is a UTF-8 encoded, newline-delimited set of key-value pairs. Each line consists of a key, a colon, a space, and a value. Keys are ASCII alphanumeric with hyphens. Values are UTF-8 strings. The header is not encrypted, allowing the receiving node to identify the package before any cryptographic operations.

Required header fields:

- `Package-ID`: A unique identifier for this package. Format: `<sender-node-id>-<timestamp>-<sequence>`, where timestamp is ISO 8601 UTC and sequence is a zero-padded four-digit integer to handle multiple packages created in the same second.
- `Format-Version`: The sync package format version. Must match the version in the preamble.
- `Sender-Node-ID`: The unique identifier of the sending node (the SHA-256 hash of the sender's public key, hex-encoded, first 32 characters).
- `Recipient-Node-ID`: The unique identifier of the intended recipient node (same format as Sender-Node-ID).
- `Created-Timestamp`: The creation time of the package, ISO 8601 UTC format.
- `Software-Version`: The version of the institutional software that created the package, in semantic versioning format (MAJOR.MINOR.PATCH).
- `Standards-Version`: The version of the institutional standards this package conforms to.
- `Sync-Type`: One of `FULL`, `INCREMENTAL`, or `GOVERNANCE`. FULL contains all shared data. INCREMENTAL contains only changes since the last sync with this recipient. GOVERNANCE contains only governance-related content (used by D17-008).
- `Description`: A human-readable description of the package contents. Maximum 500 characters.
- `Previous-Sync-ID`: The Package-ID of the most recent sync package successfully exchanged with this recipient, or `NONE` if this is the first sync. This field enables the receiving node to detect gaps in the sync sequence.
- `Item-Count`: The number of items in the payload (uint32, decimal string).
- `Payload-Compression`: The compression algorithm applied to the payload before encryption. One of `NONE`, `ZSTD`, or `GZIP`. Zstandard is preferred for its balance of compression ratio and speed.
- `Encryption-Algorithm`: The algorithm used to encrypt the payload. Currently: `AES-256-GCM`. This field enables future algorithm migration.
- `Hash-Algorithm`: The hash algorithm used for content hashes in the manifest and for the integrity seal. Currently: `SHA-256`. This field enables future algorithm migration.

Optional header fields:

- `Urgency`: One of `ROUTINE`, `PRIORITY`, `CRITICAL`. Defaults to `ROUTINE`. CRITICAL indicates the package contains urgent content such as trust revocation notices.
- `Expiry-Date`: ISO 8601 UTC date after which the package should be considered stale and not applied. Used for time-sensitive governance actions.
- `Human-Note`: An additional free-text note from the sender to the recipient. Maximum 2000 characters.

### 4.3 Manifest Specification

The manifest is a structured document listing every item in the payload. It is encoded as UTF-8 text in a defined schema and is cryptographically signed by the sender.

The manifest consists of a manifest header and a list of item entries.

**Manifest header fields:**

- `Manifest-Version`: Currently `1`.
- `Total-Items`: The count of item entries.
- `Total-Size`: The total uncompressed size of all payload items in bytes.
- `Manifest-Created`: ISO 8601 UTC timestamp.
- `Sender-Fingerprint`: The full SHA-256 fingerprint of the sender's public key (64 hex characters).

**Item entry fields (one set per item):**

- `Item-Path`: The logical path of the item within the institutional content structure. Uses forward-slash separators regardless of the sender's operating system. Maximum 1024 characters. Characters restricted to: UTF-8 alphanumeric, hyphen, underscore, period, forward-slash, and space.
- `Item-Type`: One of `DOCUMENT`, `METADATA`, `GOVERNANCE`, `MEMORY`, `BINARY`, `SYSTEM`.
- `Item-Hash`: The SHA-256 hash of the item's content (64 hex characters).
- `Item-Size`: The uncompressed size of the item in bytes.
- `Item-Modified`: The last modification timestamp of the item, ISO 8601 UTC.
- `Item-Action`: One of `CREATE`, `UPDATE`, `DELETE`. CREATE indicates a new item. UPDATE indicates a modification to an existing item. DELETE indicates the item should be removed on the receiving node (the payload contains the deletion record, not the deleted content).
- `Item-Previous-Hash`: For UPDATE actions, the SHA-256 hash of the version this update is based on. This enables the receiving node to detect conflicts (if the local version's hash does not match Item-Previous-Hash, the item has been modified locally since the last sync). For CREATE and DELETE actions, this field is `NONE`.

The manifest is signed using the sender's private key with the Ed25519 signature algorithm. The signature covers the entire manifest text (from the first byte of the manifest header through the last byte of the last item entry). The signature is appended to the manifest as a final field:

- `Manifest-Signature`: The Ed25519 signature, base64-encoded.
- `Signing-Key-Fingerprint`: The SHA-256 fingerprint of the public key corresponding to the signing private key.

### 4.4 Payload Specification

The payload contains the actual data items listed in the manifest. The payload is first assembled as an uncompressed archive, then compressed (if Payload-Compression is not NONE), then encrypted.

**Payload archive format.** The payload is a TAR archive (POSIX.1-2001 / pax format) containing one file per manifest item entry. Each file in the archive is named with the Item-Path from the manifest. The archive contains no directories as separate entries -- directories are implied by the file paths. The archive contains no symlinks, no hard links, no device nodes, and no special files. Only regular files are permitted.

For DELETE actions, the archive contains a deletion record file at the Item-Path with the extension `.deleted` appended. The deletion record is a UTF-8 text file containing: the original Item-Hash (the hash of the content being deleted), the deletion timestamp, and the reason for deletion.

**Compression.** If Payload-Compression specifies an algorithm, the TAR archive is compressed as a single unit before encryption. The compressed data is what gets encrypted.

**Encryption.** The compressed (or uncompressed) archive is encrypted using AES-256-GCM. The encryption key is derived from the bilateral trust relationship between the sender and recipient, using HKDF-SHA256 with the following inputs:

- Input key material: The shared secret established during the key exchange ceremony (D17-005).
- Salt: The Package-ID, UTF-8 encoded.
- Info: The string `syncpkg-payload-encryption`, UTF-8 encoded.

The AES-256-GCM nonce (96 bits) is randomly generated and prepended to the ciphertext. The GCM authentication tag (128 bits) is appended to the ciphertext.

The encrypted payload section therefore contains: `[nonce (12 bytes)][ciphertext][auth-tag (16 bytes)]`.

### 4.5 Integrity Seal Specification

The integrity seal covers the entire package: header, manifest, and encrypted payload. It provides tamper detection for the package as a whole, complementing the manifest's per-item hashes and the payload's GCM authentication.

The integrity seal is computed as follows:

1. Concatenate the raw bytes of the header section, manifest section, and payload section (including encryption overhead).
2. Compute the SHA-256 hash of the concatenation.
3. Sign the hash using the sender's Ed25519 private key.

The seal section contains:

- `Seal-Hash`: The SHA-256 hash (64 hex characters).
- `Seal-Signature`: The Ed25519 signature of the hash, base64-encoded.
- `Seal-Key-Fingerprint`: The SHA-256 fingerprint of the signing key.
- `Seal-Timestamp`: The time the seal was created, ISO 8601 UTC. This should be within seconds of Created-Timestamp in the header. A significant discrepancy indicates the package was modified after initial creation.

### 4.6 Version Negotiation

When two nodes run different software versions, the sync package format version determines compatibility. The negotiation rules are:

**Backward compatibility.** A node running software version N must be able to read sync packages created by software version N or any earlier version. This means every format change must be backward-compatible: new fields may be added to the header and manifest, but existing fields may not be removed or have their semantics changed.

**Forward tolerance.** A node receiving a sync package with a format version higher than its own should attempt to process it. Unknown header fields are ignored. Unknown manifest item types are logged and the items are placed in a quarantine directory for manual review. If the format version is two or more major versions ahead, the node should reject the package and log a warning that a software upgrade is needed.

**Version advertisement.** Each sync package's header includes both the format version and the software version. A node that consistently receives packages with a higher software version knows that it should consider upgrading. This information is logged but does not trigger any automatic action.

**Migration packages.** When a format version change occurs that requires data migration, the upgrading node may include a `Format-Migration-Guide` item in the payload (Item-Type: SYSTEM). This is a human-readable document explaining the changes and any actions the recipient should take.

## 5. Rules & Constraints

- **R-D17-007-01: Strict Format Compliance.** Every sync package must comply with this specification in its entirety. A package that omits a required field, uses an unsupported encoding, or violates any structural requirement defined in Section 4 is non-conforming and must be rejected by the receiving node. Partial compliance is not acceptable.

- **R-D17-007-02: Cryptographic Algorithm Constraints.** The encryption algorithm (AES-256-GCM), hash algorithm (SHA-256), signature algorithm (Ed25519), and key derivation function (HKDF-SHA256) are the mandatory algorithms for format version 1. Alternative algorithms may be introduced in future format versions but must not replace these algorithms in version 1 packages. Algorithm agility is achieved through format versioning, not through runtime negotiation.

- **R-D17-007-03: No Executable Content in Payload.** The payload must not contain executable files (binaries, scripts, bytecode) unless the Item-Type is explicitly SYSTEM and the item is a format migration guide or a similar institutional tool. Executable content in sync packages is a supply chain attack vector and is forbidden for all other item types.

- **R-D17-007-04: Path Normalization.** All Item-Path values must be normalized: no `.` or `..` components, no double slashes, no trailing slashes, no absolute paths (all paths are relative to the institutional content root). A manifest containing non-normalized paths is non-conforming.

- **R-D17-007-05: Maximum Package Size.** A single sync package must not exceed 128 GB (after compression, before encryption). This limit exists to ensure that packages can be transported on commonly available removable media and can be processed within reasonable time and memory constraints. If the sync scope exceeds this limit, it must be split into multiple sequential packages, linked by a `Continuation-Of` header field.

- **R-D17-007-06: Hash Verification Is Mandatory Before Application.** After decryption and decompression, every item in the payload must be hash-verified against its manifest entry before any item is applied to the production system. If any single item fails hash verification, the entire package is rejected. This is the atomic integrity rule: the package succeeds as a whole or fails as a whole.

- **R-D17-007-07: Manifest Must Be Verified Before Decryption.** The manifest signature must be verified before the payload is decrypted. If the manifest signature is invalid, the payload is never decrypted. This prevents the processing of payloads whose contents cannot be trusted, even if the outer seal verifies.

## 6. Failure Modes

- **Format version mismatch beyond tolerance.** The receiving node encounters a package with a format version it cannot process. The package is unreadable. Mitigation: the version negotiation rules in Section 4.6 define the tolerance range. When a package is rejected due to version mismatch, the rejection is logged with the specific version numbers, providing the operator with clear information about what upgrade is needed. The sending node should include a Format-Migration-Guide when creating packages in a new format version.

- **Algorithm deprecation.** A cryptographic algorithm used in the format is found to be insecure. All existing packages signed or encrypted with that algorithm are potentially compromised. Mitigation: algorithm changes require a new format version. The old format version remains supported for reading (backward compatibility) but is no longer used for creating new packages. A migration period allows nodes to upgrade. The Commentary Section of this article should record any algorithm deprecation decisions.

- **Path encoding conflicts.** The sender's institutional content structure uses characters in Item-Path that the receiver's filesystem cannot represent. Mitigation: the restricted character set in Section 4.3 is deliberately conservative. If a path contains characters outside this set, the manifest is non-conforming and the package is rejected. The sender must rename the item before including it in a sync package.

- **Nonce reuse in encryption.** The same AES-256-GCM nonce is used for two different packages encrypted with the same key. This catastrophically compromises the encryption of both packages. Mitigation: the nonce is randomly generated (96 bits, providing a collision probability of approximately 2^-48 per pair of packages). Additionally, the key derivation uses the Package-ID as salt, meaning that even if the same nonce were generated twice, the derived keys would differ because the Package-IDs differ. This provides defense in depth against nonce reuse.

- **Manifest-payload inconsistency.** The manifest lists items that are not in the payload, or the payload contains items not listed in the manifest. Mitigation: the hash verification procedure (R-D17-007-06) catches missing or corrupted items. An item in the payload but not in the manifest is detected during extraction (unknown files in the archive are flagged and quarantined). Both conditions result in package rejection.

- **Package truncation.** The package file is incomplete -- truncated during write to media or during read from media. Mitigation: the preamble contains section offsets and lengths. If the file is shorter than the offset plus length of any section, the truncation is detected immediately. The integrity seal also fails on a truncated package.

## 7. Recovery Procedures

1. **If a sync package fails format validation:** Log the specific validation failure (which field, what value, what was expected). Do not attempt to repair the package. Notify the sender with the specific failure details. The sender should investigate whether the failure is a software bug (the package was created incorrectly), a transport corruption (the package was damaged in transit), or a tampering attempt (the package was deliberately modified). The sender creates and transmits a replacement package.

2. **If the receiving node cannot process the package's format version:** Log the version mismatch. Check whether a software upgrade is available on local media. If an upgrade is available and has been validated per Domain 18 procedures, consider upgrading before reprocessing the package. If no upgrade is available, notify the sender that a lower format version is needed. The sender should create a new package using a format version within the receiver's tolerance range.

3. **If cryptographic algorithm deprecation is announced:** Do not immediately stop processing existing packages using the deprecated algorithm. Establish a migration timeline. During the migration period, accept packages using either the old or new algorithm. After the migration deadline, reject packages using the deprecated algorithm. Document the migration in the governance log and in this article's Commentary Section.

4. **If nonce reuse is detected or suspected:** Identify all packages that may be affected. If the affected packages used the same derived key (same sender-recipient pair), consider the encrypted payloads of those packages as potentially compromised. Re-request the data via new sync packages. Rotate the bilateral encryption keys by performing a new key exchange. Document the incident.

5. **If a package exceeds the maximum size limit:** The sender must split the sync scope into multiple packages. Each subsequent package references the previous one via the `Continuation-Of` header field. The receiving node processes the packages in order, treating them as a single logical sync. If any package in the continuation chain fails validation, the entire chain is rejected.

## 8. Evolution Path

- **Years 0-5:** Format version 1 is established and tested. Expect to discover edge cases not covered by this specification -- unusual characters in paths, unexpected file sizes, compression ratios that challenge memory constraints. Document these discoveries in the Commentary Section. Resist the urge to increment the format version for minor issues; instead, clarify the specification through addenda.

- **Years 5-15:** The first format version increment may be needed. New item types, new metadata fields, or new compression algorithms may justify a version 2. The version 2 specification must be written as a companion to this document, not a replacement. Version 1 packages must remain readable indefinitely.

- **Years 15-30:** Cryptographic algorithm review becomes critical. SHA-256 and AES-256 are considered strong as of 2026, but cryptographic landscapes change. The institution should monitor (through imported security literature) whether these algorithms remain appropriate. If migration is needed, it should be planned over a multi-year timeline, not rushed.

- **Years 30-50+:** The sync package format may have gone through multiple version increments. The preamble structure, with its magic number and version field, should remain stable across all versions. A node running in 2076 should be able to read the first eight bytes of any sync package ever created and determine whether it can process it. That is the fifty-year promise of this specification.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The choice of TAR as the payload archive format deserves explanation. TAR is old, simple, well-understood, and implemented in every operating system. It has no compression built in (compression is handled separately), no encryption built in (encryption is handled separately), and no complex metadata that could introduce compatibility issues. It is, in the best sense, boring. Boring is what we want in a format that must be readable in fifty years. ZIP was considered and rejected because its integrated compression and its handling of character encodings have caused decades of interoperability problems. TAR with external compression and external encryption gives us full control over each layer.

The choice of Ed25519 for signatures and AES-256-GCM for encryption reflects a preference for algorithms that are fast, well-analyzed, and implemented in every major cryptographic library. Ed25519 in particular was chosen because it produces short signatures, has no known weaknesses as of 2026, and is deterministic (the same message and key always produce the same signature, eliminating a class of implementation bugs).

The 128 GB package size limit is generous by current standards but may seem small in the future. It can be raised in a future format version. The continuation mechanism ensures that even version 1 can handle arbitrarily large syncs through multiple packages.

## 10. References

- D17-001 -- Federation Philosophy (sync package concept; R-D17-05 self-containment requirement)
- D17-003 -- Air-Gapped Synchronization Procedures (sync protocol; package creation and reception steps)
- D17-005 -- Trust Establishment and Verification (key exchange; shared secrets for encryption)
- D17-004 -- Conflict Resolution Between Nodes (Item-Previous-Hash enables conflict detection)
- SEC-001 -- Threat Model and Security Philosophy (cryptographic requirements; integrity model)
- ETH-001 -- Ethical Foundations (Principle 3: Transparency; Principle 4: Longevity Over Novelty)
- OPS-001 -- Operations Philosophy (complexity budget; operational sustainability)
- Domain 18 -- Import & Quarantine (content validation after payload extraction)
- Domain 20 -- Institutional Memory (sync event logging)

---

---

# D17-008 -- Distributed Governance Protocol

**Document ID:** D17-008
**Domain:** 17 -- Scaling & Federation
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** D17-001, D17-003, D17-007, D17-004, D17-005, GOV-001, SEC-001
**Depended Upon By:** All federation governance activities. Referenced by D17-004 (Class 4 governance conflicts), D17-006 (withdrawal governance), and Domain 20 (governance event recording).

---

## 1. Purpose

This article defines the protocol by which governance operates across multiple sovereign nodes in the federation. Governance within a single node is straightforward: the operator makes decisions according to GOV-001. Governance across a federation of sovereign nodes is fundamentally different. No node has authority over any other. No majority can compel a minority. No central body can issue decrees. Yet the federation must maintain sufficient constitutional consistency to function as a coherent association rather than a collection of strangers who happen to share a data format.

The distributed governance protocol solves this problem not by creating central authority but by creating a structured process through which sovereign nodes can propose changes, deliberate on them, vote, exercise veto rights, ratify decisions, and record the results -- all through the air-gapped, high-latency communication channel of sync packages. Every governance action is a document. Every deliberation is an exchange of documents. Every decision is a signed record. The protocol transforms the federation's limitation -- asynchronous, batch-oriented communication -- into a feature: deliberation is slow, which means it is deliberate.

This article is addressed to every operator who participates in federation governance: proposing changes to shared standards, voting on federation-wide policies, exercising veto rights, or ratifying decisions that affect the constitutional consistency of the federation.

## 2. Scope

This article covers:

- The governance sync: a specialized sync package type for governance actions.
- The proposal format: how governance proposals are structured and submitted.
- The deliberation process: how nodes exchange positions on pending proposals.
- Voting mechanisms: how votes are cast, counted, and recorded.
- Veto rights: when and how a node may veto a proposal.
- Ratification procedures: how approved proposals become binding.
- Constitutional consistency: how to maintain a coherent shared constitutional baseline without central authority.
- The governance timeline: how the protocol handles the latency inherent in air-gapped communication.

This article does not cover:

- Single-node governance (see GOV-001).
- The sync protocol mechanics (see D17-003 and D17-007).
- Trust management (see D17-005), though trust level affects governance participation rights.
- Conflict resolution for non-governance content (see D17-004).
- The philosophical justification for distributed governance (see D17-001, Section 4.3).

## 3. Background

### 3.1 The Governance Problem in Sovereign Federation

GOV-001 defines the authority model for a single node. The operator is the ultimate authority. Decisions are made, recorded, and implemented. The governance system is clean because there is exactly one decision-maker.

Federation breaks this cleanliness. When multiple sovereign nodes agree to share standards, maintain compatible data formats, and synchronize knowledge, they create a shared space that no single node governs. Changes to this shared space -- a new data format, a revised naming convention, an updated security standard -- affect all nodes. No node has the right to impose such changes unilaterally. But waiting for perfect consensus on every change would paralyze the federation.

The distributed governance protocol navigates between these extremes. It provides a structured process that respects sovereignty (no node is compelled), enables progress (proposals can move forward with sufficient support), and preserves coherence (the shared constitutional baseline is maintained through explicit versioning and ratification).

### 3.2 Governance Through Sync Packages

All governance communication between nodes occurs through sync packages of type GOVERNANCE (as defined in D17-007, Section 4.2). These packages contain governance documents -- proposals, deliberations, votes, ratifications -- rather than content data. They follow the same cryptographic and integrity requirements as any sync package. They are subject to the same transport, quarantine, and validation procedures.

The use of sync packages for governance has three important consequences. First, governance actions are authenticated: every document is signed by its author's node key, providing non-repudiation. Second, governance actions are permanent: every document is recorded in institutional memory on the receiving node. Third, governance actions are asynchronous: there is no real-time negotiation. Each node reads proposals, deliberates internally, and responds in its next governance sync. This asynchrony makes governance slow. It also makes it thoughtful.

### 3.3 The Constitutional Baseline

The constitutional baseline is the set of shared agreements that define the federation: the shared data formats, naming conventions, quality standards, governance rules, and ethical principles that all federated nodes agree to follow. The baseline is versioned. Each version is a complete, self-contained document that specifies all shared agreements in effect at that version.

The baseline is not a constitution in the legal sense -- it has no enforcement mechanism beyond trust. A node that violates the baseline faces trust consequences (D17-005), not legal ones. But the baseline provides the coherence that makes federation meaningful. Without it, federation is just data exchange. With it, federation is a community of practice.

## 4. System Model

### 4.1 Governance Document Types

The distributed governance protocol operates through six document types, each with a defined format and purpose:

**Proposal (GOV-PROP).** A formal request to change the constitutional baseline. Contains: the proposed change, the rationale, the impact assessment, the proposed implementation timeline, and the author's node identifier and signature. A proposal is identified by a unique Proposal-ID: `GOV-PROP-<author-node-id>-<timestamp>-<sequence>`.

**Deliberation (GOV-DELIB).** A formal response to a proposal. Contains: the Proposal-ID being addressed, the responding node's position (SUPPORT, OPPOSE, AMEND, ABSTAIN), the rationale for that position, and any proposed amendments. Multiple deliberation documents may be issued by the same node for the same proposal as the discussion evolves.

**Amendment (GOV-AMEND).** A formal modification to a pending proposal. Must reference the original Proposal-ID. May be proposed by any node, not just the original author. Amendments are themselves subject to deliberation and voting.

**Vote (GOV-VOTE).** A formal vote on a proposal (or amended proposal). Contains: the Proposal-ID, the voting node's decision (APPROVE, REJECT, ABSTAIN), the rationale, and a commitment to implement if approved. Votes are final once cast -- a node may not change its vote, but it may issue a subsequent vote on an amended version of the proposal.

**Veto (GOV-VETO).** A formal invocation of veto rights. Contains: the Proposal-ID, the vetoing node's objection, the specific harm the proposal would cause to the vetoing node's sovereignty or institutional integrity, and whether the veto is absolute (the node will withdraw from the federation if the proposal is adopted) or conditional (the node requests specific modifications to address its concerns).

**Ratification (GOV-RATIFY).** The record that a proposal has been approved and becomes part of the constitutional baseline. Contains: the Proposal-ID, the vote tally, the ratification timestamp, the new baseline version number, and the full text of the ratified change. The ratification document is signed by all approving nodes.

### 4.2 The Governance Lifecycle

A governance action proceeds through five phases:

**Phase 1: Proposal.** A node identifies a need for change to the constitutional baseline. The node drafts a GOV-PROP document and includes it in its next governance sync to all federated nodes. The proposal enters a deliberation period.

**Phase 2: Deliberation.** Each receiving node reviews the proposal and may issue GOV-DELIB documents expressing its position. Amendments (GOV-AMEND) may be proposed. The deliberation period has a minimum duration of 90 days from the date the last federation node receives the proposal. This minimum exists because air-gapped communication is slow; a node that receives the proposal late must have adequate time to deliberate. The proposing node may extend but never shorten the deliberation period.

**Phase 3: Voting.** After the deliberation period closes, each node casts a GOV-VOTE. The voting period is 60 days from the close of deliberation. Nodes that do not cast a vote within the voting period are recorded as ABSTAIN. Votes are transported via governance sync packages.

**Phase 4: Veto Window.** After all votes are received (or the voting period closes), there is a 30-day veto window. Any node may issue a GOV-VETO during this window. A veto halts the ratification process and triggers a resolution phase (see Section 4.4).

**Phase 5: Ratification.** If the proposal receives sufficient approval (see Section 4.3) and no veto is issued, the proposing node prepares a GOV-RATIFY document and distributes it via governance sync. Each receiving node applies the ratified change to its local copy of the constitutional baseline and increments the baseline version number.

### 4.3 Voting Thresholds

The voting threshold depends on the category of the proposed change:

**Category A -- Foundational Changes.** Changes to the ethical principles (ETH-001), the constitutional mandate (CON-001), or the governance model itself (GOV-001 or this protocol). Threshold: unanimous approval of all federated nodes. Abstentions are not permitted for Category A changes -- each node must explicitly approve or reject. This is the highest bar because these changes affect the identity of the institution itself.

**Category B -- Standard Changes.** Changes to shared data formats, naming conventions, quality standards, or operational procedures. Threshold: approval by at least two-thirds of federated nodes. Abstentions are counted as neither approval nor rejection but are recorded.

**Category C -- Administrative Changes.** Changes to federation operational parameters (sync frequency recommendations, format version adoption timelines, trust audit schedules). Threshold: approval by a simple majority of federated nodes. Abstentions count as neither approval nor rejection.

In all categories, the total number of federated nodes is the count of nodes with active trust relationships (Trust Level 1 or above) at the time the voting period opens. Nodes that withdrew during the voting period are excluded from the count.

### 4.4 Veto Rights and Resolution

Any federated node may exercise a veto on any proposal, regardless of category. The veto is the ultimate expression of sovereignty in the governance protocol -- a node's declaration that a proposed change would harm its institutional integrity to a degree that cannot be accepted.

A veto triggers the following resolution process:

1. The proposing node receives the veto and reviews the objection.
2. The proposing node may withdraw the proposal, modify the proposal to address the veto, or maintain the proposal unchanged.
3. If the proposal is modified, it re-enters deliberation as an amended proposal (Phase 2) with a new deliberation period.
4. If the proposal is maintained unchanged, the vetoing node must decide: accept the proposal despite its objection (withdraw the veto), or maintain the veto.
5. If the veto is maintained and the proposal has otherwise met its voting threshold, the federation faces a constitutional crisis. The options are: the proposal is abandoned, the federation formally acknowledges a divergence on this point (with the vetoing node following different rules), or the vetoing node withdraws from the federation (D17-006).

The veto is deliberately powerful. It can block any change, regardless of how many other nodes support it. This power exists because the federation is voluntary. A node that is compelled to accept changes it considers harmful is no longer sovereign. The price of this power is that it can be used to obstruct progress. The Commentary Section of this article should record every veto and its outcome, building a body of precedent for future governance.

### 4.5 Constitutional Baseline Management

The constitutional baseline is a versioned document maintained independently by each node. Consistency between nodes is achieved through the ratification process, not through central distribution.

**Baseline version numbering.** Versions are sequential integers. Each ratification increments the version by one. There is no branching -- the baseline is a linear sequence of versions.

**Baseline divergence detection.** Each governance sync package includes the sender's current baseline version in the header (Standards-Version field). If a receiving node observes a different baseline version, it knows that the two nodes are not in agreement. This may be temporary (one node has not yet applied a recent ratification) or it may indicate a genuine divergence (one node rejected or vetoed a change that the other accepted).

**Baseline reconciliation.** When two nodes have different baseline versions, they exchange their baseline documents. Each node reviews the differences and identifies whether the divergence is a timing issue (the lagging node can simply apply pending ratifications) or a substantive disagreement (the nodes have made different governance decisions). Timing issues are resolved by applying pending ratifications in order. Substantive disagreements are escalated as Class 4 governance conflicts (D17-004).

### 4.6 The Governance Sync Schedule

Governance syncs follow the same physical transport mechanisms as regular syncs (D17-003) but may be scheduled independently. The recommended governance sync frequency is:

- During active deliberation: every sync opportunity (every time physical transport occurs between nodes).
- During quiet periods: at least quarterly, coinciding with the OPS-001 quarterly review cycle.
- For urgent governance actions (trust revocations, security incidents): immediately, using whatever transport is available. The Urgency field in the sync package header should be set to CRITICAL.

## 5. Rules & Constraints

- **R-D17-008-01: No Unilateral Changes to Shared Standards.** No node may unilaterally change any element of the constitutional baseline. All changes must go through the governance lifecycle defined in Section 4.2. A node may adopt stricter local standards that exceed the baseline, but it may not relax baseline requirements unilaterally.

- **R-D17-008-02: Deliberation Periods Are Minimum, Not Maximum.** The 90-day deliberation period and 60-day voting period are minimums. They may be extended by the proposing node or by mutual agreement. They may never be shortened. Governance decisions made under time pressure are governance decisions made badly.

- **R-D17-008-03: All Governance Documents Are Permanent.** Every GOV-PROP, GOV-DELIB, GOV-AMEND, GOV-VOTE, GOV-VETO, and GOV-RATIFY document is a permanent institutional record. These documents are stored in institutional memory on all nodes that receive them and are never modified or deleted. They form the legislative history of the federation.

- **R-D17-008-04: Votes Are Final.** A node may not change its vote on a specific proposal after the vote has been cast and distributed. If the proposal is amended, the amended version is a new proposal that requires new votes. This prevents vote manipulation through iterative pressure.

- **R-D17-008-05: Veto Must Specify Harm.** A veto document must specify the concrete harm the proposal would cause to the vetoing node. "I disagree" is not a valid veto. "This change would require me to convert 10,000 documents to a new format within 90 days, which exceeds my operational capacity" is a valid veto. The requirement to specify harm discourages frivolous vetoes while preserving the right for legitimate use.

- **R-D17-008-06: Trust Level Affects Governance Rights.** Only nodes at Trust Level 1 or above may participate in governance. A node at Trust Level 0 (Untrusted) has no standing to propose, deliberate, vote, or veto. Within the eligible trust levels, all nodes have equal governance rights regardless of trust level. A Level 1 node's vote counts the same as a Level 3 node's vote.

- **R-D17-008-07: Ratification Is Not Implementation.** Ratification means the change has been approved. Implementation -- the actual adoption of the change on each node -- proceeds according to the timeline specified in the proposal. Ratification without implementation is a governance failure that must be tracked and addressed during the quarterly review.

## 6. Failure Modes

- **Governance paralysis.** The deliberation and voting process is so slow that the federation cannot adapt to changing circumstances. Proposals take years to ratify. The baseline becomes outdated. Mitigation: the minimum timelines (90 days deliberation, 60 days voting, 30 days veto) are designed to be adequate, not excessive. If paralysis develops, the cause is usually not the timeline but the scope of proposals. Large, sweeping proposals that change many things at once are harder to deliberate than focused proposals that change one thing. The remedy is smaller, more frequent proposals.

- **Governance fatigue.** Operators stop participating in governance because the volume of proposals exceeds their capacity to review them. Proposals pass with only a few votes because most nodes abstain. The governance process degenerates into rule by the active minority. Mitigation: the proposal process should be self-limiting. Each proposal requires effort to create, deliberate on, and vote on. If governance fatigue develops, the federation should declare a governance moratorium -- a period during which no new proposals are submitted and the backlog is cleared.

- **Veto abuse.** A single node uses its veto power to block all changes, effectively granting itself control over the federation's evolution. The veto was designed to protect sovereignty; it is being used to enforce stagnation. Mitigation: R-D17-008-05 requires vetoes to specify concrete harm. A pattern of vetoes that cannot articulate concrete harm is a pattern of abuse. The other nodes may choose to proceed despite the veto, understanding that this may result in the vetoing node's withdrawal. This is a drastic outcome, but it is preferable to allowing one node to hold the federation hostage.

- **Baseline fragmentation.** Different nodes operate under different baseline versions because ratifications are not applied consistently. The federation's shared standards become fiction -- each node follows a different version. Mitigation: baseline version is included in every sync package header. Nodes can detect baseline divergence during every sync. The quarterly governance audit should include a baseline version check across all trusted nodes.

- **Constitutional crisis.** A veto on a Category A proposal creates an irreconcilable disagreement about the federation's foundational principles. The federation cannot continue as a single coherent entity. Mitigation: this is not a failure of the protocol. It is a genuine disagreement that the protocol has surfaced and formalized. The resolution options -- withdrawal, divergence, or compromise -- are documented in Section 4.4. The protocol does not prevent constitutional crises. It ensures they are handled transparently.

- **Governance sync starvation.** Physical transport between nodes becomes infrequent, and governance documents are delayed for months or years. Proposals expire. Voting periods close before all nodes have received the proposal. Mitigation: governance syncs should be prioritized in transport scheduling. If physical transport is severely constrained, the minimum deliberation period should be extended to accommodate the reality. Governance that ignores communication constraints is governance in name only.

## 7. Recovery Procedures

1. **If a governance document is lost in transit:** The sending node retransmits the document in its next governance sync. Governance documents are idempotent -- receiving the same document twice causes no harm (the second copy is logged as a duplicate and discarded). If the lost document was a vote and the voting period has closed, the sending node's position is recorded as ABSTAIN for that round. The proposal may be resubmitted if the lost votes would have changed the outcome.

2. **If baseline versions have diverged across the federation:** Each node declares its current baseline version in its next governance sync. The node with the most recent version distributes a baseline reconciliation package containing all ratification documents between the lowest and highest version in the federation. Each lagging node reviews and applies the pending ratifications. If any ratification is rejected by a lagging node (because it was adopted during a period when that node was disconnected and unable to participate), the rejection is treated as a retroactive veto and enters the resolution process of Section 4.4.

3. **If a veto has caused a constitutional crisis:** The federation must convene a formal resolution process. Each node prepares a position paper explaining its stance. Position papers are exchanged via governance sync. The federation has three paths: compromise (the proposal is modified to address the veto), divergence (the vetoing node follows different rules on this specific point, with the divergence formally documented), or separation (the vetoing node withdraws, following D17-006). The choice of path is itself a governance decision, though one that may proceed outside the normal proposal process due to its urgency.

4. **If governance fatigue has degraded participation:** Declare a governance moratorium of at least 90 days. During the moratorium, no new proposals are submitted. Existing proposals remain in their current phase. Use the moratorium to clear the deliberation backlog. After the moratorium, resume governance with a commitment to focused, limited-scope proposals.

5. **If a node has been making unilateral changes to the baseline:** Identify the specific changes. Determine whether they are beneficial (improvements that should be proposed through the governance process), harmful (changes that violate the baseline), or neutral (local customizations that do not affect interoperability). For beneficial changes, the node should submit them as formal proposals retroactively. For harmful changes, the node must revert them. For neutral changes, document them as local variations.

## 8. Evolution Path

- **Years 0-5:** The distributed governance protocol is largely theoretical. With one or two nodes, governance decisions are made through direct conversation, not through formal proposals. Use this period to establish the baseline document and the habit of recording governance decisions, even when the formal protocol is not yet needed.

- **Years 5-15:** The first formal governance cycle may occur. The protocol will be tested against real disagreements and real communication constraints. Expect the deliberation periods to feel too long when the issue is simple and too short when the issue is complex. The Commentary Section should record these experiences.

- **Years 15-30:** The governance protocol should be mature. The baseline should have gone through multiple version increments. The governance archive should contain enough proposals, deliberations, votes, and ratifications to constitute a legislative history. Future operators should be able to read this history and understand why the federation's shared standards are what they are.

- **Years 30-50+:** The governance protocol may be operated by people who did not write the constitution they are amending. The protocol must be comprehensible to these future operators. The formal structure -- proposals, deliberations, votes, vetoes, ratifications -- should be intuitive enough that a new operator can participate after reading this document. If the protocol requires oral tradition to understand, it has failed.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The 90-day minimum deliberation period will seem absurdly long to anyone accustomed to internet-speed governance. It is not long enough. Consider the scenario: Node A creates a proposal and includes it in a sync package. The sync package is transported physically to Node B, which may take days or weeks. Node B reads the proposal, deliberates internally, and prepares a response. The response is included in a sync package back to Node A, which takes more days or weeks. Node A reads the response and may need to respond in turn. Each round trip takes weeks. Ninety days may accommodate only two or three rounds of deliberation. This is not a bug. Governance decisions that affect the constitutional baseline of a voluntary federation of sovereign nodes deserve weeks of deliberation, not hours.

The veto power is the most controversial element of this protocol. It gives any single node the power to block any change. This is an enormous power, and it will be abused at some point. The requirement to specify concrete harm is a guardrail, not a barrier. A determined obstructionist can always articulate some harm. The true check on veto abuse is social, not procedural: a node that vetoes everything will find itself increasingly isolated, its trust level declining, its influence waning. The protocol provides the tools. The community provides the judgment.

I considered and rejected a consensus model (all nodes must agree on everything). Consensus works for small groups with continuous communication. It does not work for sovereign nodes communicating asynchronously via sneakernet. I also considered and rejected a pure majority-rule model. Majority rule in a small federation (three to five nodes) means a single defection changes the outcome. The tiered threshold system -- unanimity for foundational changes, supermajority for standards, simple majority for administration -- balances these concerns.

## 10. References

- D17-001 -- Federation Philosophy (sovereignty principle; no-central-authority constraint; distributed governance model)
- D17-003 -- Air-Gapped Synchronization Procedures (governance sync transport)
- D17-004 -- Conflict Resolution Between Nodes (Class 4 governance conflicts)
- D17-005 -- Trust Establishment and Verification (trust levels affect governance rights)
- D17-006 -- Node Withdrawal and Data Separation (withdrawal as consequence of governance failure)
- D17-007 -- Sync Package Specification (GOVERNANCE sync type; format for governance documents)
- GOV-001 -- Authority Model (single-node governance; decision tiers; governance log)
- SEC-001 -- Threat Model and Security Philosophy (authentication of governance documents)
- ETH-001 -- Ethical Foundations (Principle 1: Sovereignty; Principle 3: Transparency)
- OPS-001 -- Operations Philosophy (quarterly review; governance in operational tempo)
- Domain 20 -- Institutional Memory (governance archive as permanent record)

---

---

# D18-003 -- Malware Scanning for Air-Gapped Systems

**Document ID:** D18-003
**Domain:** 18 -- Import & Quarantine
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** D18-001, SEC-001, OPS-001, D17-003
**Depended Upon By:** D18-004 (format conversion follows scanning), D18-005 (scan results are part of provenance), and all import operations that process external content.

---

## 1. Purpose

This article defines how the institution scans imported content for malware in an air-gapped environment where internet-connected signature updates are impossible. Malware scanning is one of the most critical functions of the import pipeline (D18-001, Section 4.1, Stage 2: Quarantine), yet it is also one of the most challenging to perform offline. Connected systems rely on continuously updated signature databases, cloud-based reputation services, and real-time threat intelligence feeds. None of these are available here. The institution must detect malware using only the tools and knowledge it carries locally.

This article specifies the complete malware scanning approach: offline signature management, heuristic analysis, behavioral sandboxing, and the manual review procedures that complement automated tools. It defines the scan procedure step by step, from initial media insertion through final disposition. It addresses the hardest question in air-gapped security: how do you protect against threats you have never seen, using tools that cannot phone home?

The answer is defense in depth. No single scanning technique is sufficient. Signature scanning catches known threats. Heuristic analysis catches threats that resemble known patterns. Behavioral sandboxing catches threats that execute malicious actions regardless of their signature. Manual review catches threats that fool all automated tools. Together, these layers provide a scanning capability that is imperfect -- all scanning is imperfect -- but substantially reduces the risk of malware entering the institution.

This article is addressed to the operator performing the scan, to the person maintaining the scanning tools and signature databases, and to the future operator who inherits a scanning infrastructure they did not build.

## 2. Scope

This article covers:

- The scanning architecture: how scanning tools are organized within the quarantine environment.
- Offline signature management: how to maintain and update malware signatures without internet access.
- Heuristic scanning: how to detect unknown malware through structural and behavioral patterns.
- Behavioral analysis in sandbox environments: how to execute suspicious content safely and observe its behavior.
- The scan procedure: the complete step-by-step process for scanning imported content.
- Scan result interpretation: how to read and act on scan findings.
- What to do when malware is detected: containment, documentation, and recovery.
- Maintaining scanning effectiveness over time in an air-gapped context.

This article does not cover:

- The broader quarantine architecture (see D18-001, Section 4.3).
- Format validation and conversion (see D18-004).
- Provenance tracking (see D18-005).
- The import pipeline as a whole (see D18-001).
- Network-level security (not applicable in an air-gapped system).
- Scanning of content that originates within the institution (covered by internal integrity checks in SEC-001).

## 3. Background

### 3.1 The Air-Gap Scanning Problem

Connected malware scanners update their signature databases daily, sometimes hourly. They submit unknown samples to cloud analysis services. They query reputation databases. They receive real-time intelligence about emerging threats. Strip all of this away, and what remains?

What remains is a scanner that knows about yesterday's threats -- or last month's, or last year's, depending on when its signatures were last updated. The scanner cannot learn about new threats until the operator physically carries updated signatures across the air gap. Between updates, the scanner is blind to any malware created since its last signature update.

This blindness is the central challenge. It cannot be eliminated. It can be mitigated through three strategies: maximizing the frequency and quality of signature updates (Section 4.2), supplementing signatures with heuristic detection that catches threat patterns rather than specific threats (Section 4.3), and adding behavioral analysis that detects what malware does rather than what it looks like (Section 4.4). Together, these strategies provide reasonable protection. Separately, each is insufficient.

### 3.2 Threat Categories for Imported Content

Imported content can carry several categories of threats:

**Embedded malware.** Malicious code hidden within otherwise legitimate files. Examples include macro viruses in documents, exploit payloads in image files, and steganographic command channels in media files. This is the most common threat category for an air-gapped institution because it rides within the content the institution deliberately imports.

**Trojan applications.** Software that appears legitimate but contains malicious functionality. This primarily applies to software tools imported for institutional use (D18-001 acknowledges that software tools may be imported). Trojans are detected through signature matching, heuristic analysis, and behavioral observation.

**Archive bombs.** Files designed to exhaust system resources when processed -- zip bombs that expand to enormous sizes, recursive archives that trigger infinite loops, or files that consume all available memory during scanning. These may not be "malware" in the traditional sense but can disable the quarantine system.

**Format exploits.** Files that exploit vulnerabilities in the software used to view or process them. A malicious PDF that exploits a vulnerability in the PDF viewer, for instance. These are particularly dangerous because they attack the scanning infrastructure itself.

### 3.3 The Acceptable Risk Decision

No scanning system provides absolute protection. The institution must accept residual risk. The acceptable risk level is determined by the operator's judgment, informed by the scanning results and the content's trust score (D18-001, Section 4.2). High-trust content from a verified source with clean scan results carries low residual risk. Low-trust content from an unknown source with ambiguous scan results carries high residual risk. The operator's decision to admit or reject is informed by the scan but not dictated by it (R-D18-07: The Operator Decides).

## 4. System Model

### 4.1 Scanning Architecture

The scanning system operates entirely within the quarantine environment (D18-001, Section 4.3). It consists of four components:

**The Quarantine Workstation.** A dedicated machine, physically or logically isolated from all production systems. This machine is considered expendable -- if scanning triggers malware execution, the workstation can be rebuilt from a known-good image without affecting the institution. The workstation runs the scanning tools and provides the sandbox environment.

**The Signature Store.** A local database of malware signatures, stored on the quarantine workstation. The signature store is updated through a deliberate, controlled process (Section 4.2). The store includes signatures from one or more open-source antivirus engines (ClamAV is the reference implementation due to its open-source nature, cross-platform availability, and offline capability).

**The Heuristic Engine.** A set of rules and patterns that detect suspicious characteristics in files without matching specific signatures. The heuristic engine is configured locally and can be extended by the operator. It operates on structural analysis -- examining file headers, embedded objects, script content, and format anomalies.

**The Sandbox.** An isolated execution environment where suspicious files can be opened or executed and their behavior observed. The sandbox is a virtual machine or container that is destroyed and recreated from a clean image after each use. The sandbox has no network access (by definition, given the air gap), no access to institutional data, and no persistent storage that survives reset.

### 4.2 Offline Signature Management

Signature updates are the most critical maintenance task for the scanning system. The procedure for updating signatures is:

**Acquisition.** On an internet-connected machine that is not part of the institution (a separate device used solely for this purpose, or a public terminal), download the latest signature databases from the antivirus project's official distribution channel. For ClamAV, this is the official mirror network. Verify the download against the project's published checksums and GPG signatures.

**Transport preparation.** Copy the verified signature files to clean transport media (prepared according to D17-003, Section 4.3 media preparation procedures). Create a manifest of the signature files: filenames, sizes, SHA-256 hashes, and the date of the signature database.

**Quarantine intake.** Transport the media to the institution. Connect it to the quarantine workstation. Do not connect it to any production system. Verify the manifest against the files on the media. Compare the GPG signatures against the known signing key for the antivirus project (which should already be stored in the quarantine workstation's keyring from a previous, verified import).

**Installation.** Replace the existing signature database on the quarantine workstation with the new one. Run the scanner's self-test procedure to verify the new signatures load correctly. Scan a set of known-clean test files and a set of known-malware test files (the EICAR test file and any locally maintained test samples) to verify both clean-pass and malware-detection functionality.

**Documentation.** Record the update in the scanning log: date of update, signature database version, number of signatures, source, and verification results.

**Recommended update frequency.** Signature updates should be performed at least quarterly, coinciding with the OPS-001 quarterly review cycle. Monthly updates are preferred when operationally feasible. The gap between signature updates represents the window of vulnerability to new malware. This gap should be as short as practical.

### 4.3 Heuristic Scanning

Heuristic scanning detects threats by analyzing the structural and behavioral characteristics of files rather than matching specific signatures. The heuristic engine applies the following checks:

**File type verification.** The file's actual content is compared against its declared type (file extension, MIME type). A file with a `.jpg` extension that contains executable code is flagged. A file with a `.pdf` extension that has an invalid PDF header is flagged. This catches the simplest form of disguised malware.

**Embedded executable detection.** All files are scanned for embedded executable content: PE headers (Windows executables), ELF headers (Linux executables), Mach-O headers (macOS executables), and common script shebangs (`#!/bin/sh`, `#!/usr/bin/python`, etc.). Executable content embedded in non-executable files is flagged.

**Macro and active content detection.** Documents in formats that support macros or active content (Microsoft Office formats, PDF with JavaScript, HTML with embedded scripts) are scanned for the presence of such content. Active content is flagged regardless of whether it is malicious -- the operator must review and decide whether the active content is acceptable.

**Obfuscation detection.** Files are analyzed for signs of deliberate obfuscation: base64-encoded executable payloads, character encoding tricks that disguise code, encrypted or compressed sections within otherwise plaintext files, and unusual entropy patterns that suggest encrypted or compressed embedded content.

**Format anomaly detection.** Each file is parsed according to its declared format specification. Anomalies -- headers that do not match the specification, field values that exceed expected ranges, structures that are valid but unusual -- are flagged. Format anomalies are not necessarily malicious, but they correlate with exploit attempts that rely on parser vulnerabilities.

**Archive analysis.** Archive files (ZIP, TAR, RAR, 7z, and others) are extracted and their contents are recursively scanned. Nested archives are extracted to a configurable maximum depth (default: 5 levels). Archive bombs are detected by monitoring the decompression ratio -- if the extracted size exceeds 100 times the compressed size, extraction is halted and the archive is flagged.

### 4.4 Behavioral Analysis in Sandbox

Behavioral analysis is the third scanning layer. It applies to files that pass signature and heuristic scanning but remain suspicious, or to file types that are inherently risky (executables, documents with active content, interactive media).

**Sandbox preparation.** The sandbox environment is a virtual machine running a clean operating system image with standard document viewers, media players, and a file manager. The virtual machine is configured with:

- No network interfaces (redundant with the air gap, but enforced at the VM level as defense in depth).
- No shared folders with the host.
- Snapshot capability so the VM can be restored to its clean state after each analysis session.
- Process monitoring tools that log all process creation, file system changes, registry modifications (on Windows), and system call activity.
- Memory monitoring to detect unusual memory allocation patterns.

**Execution procedure.** The suspicious file is copied into the sandbox. The operator opens or executes the file using the appropriate application. The process monitor records all activity for a defined observation period (minimum 5 minutes for documents, minimum 15 minutes for executables). The operator observes the application's behavior for obvious anomalies: unexpected process spawning, attempts to modify system files, creation of new files in unexpected locations, or high CPU/memory usage that suggests cryptographic or data exfiltration activity.

**Result analysis.** After the observation period, the process monitor log is reviewed. Behavioral indicators of malware include:

- Process spawning: the file opens an application, which spawns unexpected child processes.
- File system modification: the application writes to locations outside its expected working directory.
- Persistence attempts: the application creates startup entries, scheduled tasks, or service registrations.
- Data collection: the application reads files it has no reason to access (address books, credential stores, document directories).
- Resource exhaustion: the application consumes excessive CPU, memory, or disk space.

**Post-analysis cleanup.** After analysis, the sandbox is reverted to its clean snapshot. No state from the analysis session persists. The process monitor log is exported to the quarantine workstation for documentation.

### 4.5 The Complete Scan Procedure

The following procedure is executed for every piece of imported content, without exception:

**Step 1: Media intake.** The transport media is connected to the quarantine workstation. The media is never connected to any production system. Record the intake in the scanning log: date, media type, source, carrier identity.

**Step 2: Signature scan.** Run the signature-based scanner against all files on the media. Record the results: files scanned, threats detected (if any), scanner version, signature database version.

**Step 3: Heuristic scan.** Run the heuristic engine against all files. Record all flags and anomalies.

**Step 4: Triage.** Review the results from Steps 2 and 3. Classify each file into one of four categories:
- CLEAN: No signatures matched, no heuristic flags. Proceed to format validation (D18-004).
- FLAGGED: Heuristic flags but no signature match. Requires manual review (Step 5).
- DETECTED: Signature match. Malware detected. Proceed to containment (Step 7).
- SUSPICIOUS: Heuristic flags suggest behavioral analysis is warranted. Proceed to sandbox (Step 6).

**Step 5: Manual review.** For FLAGGED files, the operator examines the heuristic flags and makes a judgment call. If the flags are explainable and benign (for example, a document contains a macro that is expected and documented by the source), the file may be reclassified as CLEAN. If the flags are unexplainable or concerning, the file is reclassified as SUSPICIOUS.

**Step 6: Behavioral sandbox.** For SUSPICIOUS files, execute the behavioral analysis procedure (Section 4.4). If the behavioral analysis reveals malicious activity, reclassify as DETECTED. If behavior is benign, reclassify as CLEAN (with a note that behavioral analysis was performed).

**Step 7: Containment.** For DETECTED files, follow the malware detection response procedure (Section 4.6).

**Step 8: Documentation.** Record the complete scan results for all files in the import log. The scan record becomes part of the content's provenance chain (D18-005).

### 4.6 Malware Detection Response

When malware is detected, the following response procedure is executed:

1. **Isolate.** The infected file is moved to a designated quarantine vault on the quarantine workstation. The vault is an encrypted directory with restricted access. The file is not deleted -- it is preserved for analysis and for updating the institution's local threat knowledge.

2. **Assess scope.** Determine whether the malware may have affected the quarantine workstation. If the malware was detected by signature scan before any file was opened, the workstation is likely unaffected. If the malware was detected during behavioral analysis (meaning the file was executed in the sandbox), the sandbox was designed to contain it, but verify the sandbox isolation was not breached.

3. **Document.** Record the detection in the security incident log: date, file name, file hash, detection method (signature, heuristic, or behavioral), malware classification (if known), source of the content, and the carrier who transported the media.

4. **Notify.** If the content came from a federated node, notify the sending node via the next sync package that malware was detected in their transmission. Include the file hash and the detection details. The sending node must investigate whether their system is compromised.

5. **Decide on remaining content.** If the media contained multiple files and only some are infected, the clean files may proceed through the import pipeline, but with heightened scrutiny. The presence of malware on the media raises the trust concern for all content from that source (D18-001, Trust Scoring).

6. **Rebuild if necessary.** If there is any suspicion that the quarantine workstation has been compromised, rebuild it from the known-good image. This is a precautionary measure. The cost of rebuilding the quarantine workstation is low. The cost of operating a compromised quarantine system is catastrophic.

7. **Update local threat knowledge.** If the detected malware is not in the current signature database (detected by heuristic or behavioral analysis), create a local signature or detection rule. This local addition supplements the official signature database and protects against re-encounter of the same threat.

## 5. Rules & Constraints

- **R-D18-003-01: All Imported Content Must Be Scanned.** No content enters the institution without passing through the scan procedure defined in Section 4.5. There are no exceptions, no fast paths, and no trusted-source exemptions. Content from federated nodes at Trust Level 3 is still scanned. Content from the operator's own external devices is still scanned. The air gap is the institution's wall. The scanner is the gate in the wall.

- **R-D18-003-02: Signature Updates Must Be Verified.** Signature database updates must be integrity-verified (checksum and GPG signature) before installation. An unverified signature update is itself a potential attack vector -- a tampered signature database could suppress detection of specific threats.

- **R-D18-003-03: The Quarantine Workstation Is Expendable.** The quarantine workstation must be rebuildable from a known-good image at any time. The workstation must not contain any data that does not also exist elsewhere (on the production system or on backup media). Rebuilding the quarantine workstation is an inconvenience, not a loss. Design accordingly.

- **R-D18-003-04: Behavioral Analysis Requires Sandbox Isolation.** No suspicious file may be opened or executed outside the sandbox environment. Opening a suspicious file on the quarantine workstation's host operating system (outside the sandbox VM) defeats the isolation model and may compromise the workstation.

- **R-D18-003-05: Detection Triggers Documentation, Not Deletion.** Detected malware is quarantined and documented, not immediately deleted. The malware sample has forensic and intelligence value. It may inform future detection rules. It may reveal information about the source's security posture. Deletion without documentation destroys this value.

- **R-D18-003-06: Scan Results Are Part of Provenance.** The scan results for every piece of imported content must be recorded as part of that content's provenance record (D18-005). Future operators must be able to determine what scanning was performed, what tools were used, what signature database version was current, and what the results were.

## 6. Failure Modes

- **Signature staleness.** The signature database has not been updated in months or years. The scanner misses threats that have emerged since the last update. The operator develops false confidence in clean scan results. Mitigation: the scanning log records the signature database version and its date. The quarterly review must include a signature freshness check. If signatures are more than six months old, the scan results should be considered incomplete, and additional heuristic and behavioral scrutiny is warranted.

- **Heuristic false positives.** The heuristic engine flags legitimate content as suspicious. If false positives are frequent, the operator develops "flag fatigue" and begins ignoring heuristic results. The heuristic layer becomes useless. Mitigation: false positives must be documented and the heuristic rules adjusted. The goal is a false positive rate low enough that flags are taken seriously. A heuristic engine that flags everything is as useless as one that flags nothing.

- **Sandbox escape.** Malware executing in the sandbox exploits a vulnerability in the virtualization software to escape the sandbox and affect the host quarantine workstation. Mitigation: the sandbox VM software must be kept current (updated through the same air-gapped signature update process). The quarantine workstation itself is expendable (R-D18-003-03). The workstation has no access to production systems. A sandbox escape compromises the quarantine workstation, not the institution.

- **Scanner compromise.** The scanning tools themselves are compromised -- either through a tampered installation or through a vulnerability exploited by a malicious file being scanned. A compromised scanner may report clean results for infected files. Mitigation: the scanner is installed from verified media (hash-checked). The scanner's own integrity is verified periodically by comparing its binaries against known-good hashes. Behavioral analysis provides an independent detection layer that does not depend on the signature scanner.

- **Unknown file format.** A file arrives in a format that the scanner cannot analyze. The scanner reports "unable to scan" rather than "clean." The operator may mistake "unable to scan" for "safe." Mitigation: unscannable files must be clearly labeled as such. They should be treated with maximum suspicion -- subjected to behavioral analysis in the sandbox and admitted only if the operator has specific reason to trust the source.

- **Archive bomb resource exhaustion.** A deliberately crafted archive consumes all available disk space or memory during decompression, rendering the quarantine workstation unusable. Mitigation: the heuristic engine monitors decompression ratios (Section 4.3). Archive extraction operates under resource limits: maximum extraction size (configurable, default 10 GB), maximum nesting depth (default 5 levels), and maximum file count (default 100,000 files). Extraction that exceeds any limit is halted and the archive is flagged.

## 7. Recovery Procedures

1. **If malware has entered the production system (scanner failure):** This is a security incident per SEC-001. Identify the infected content through its import provenance record. Isolate the infected files. Determine whether the malware has executed or spread within the production system. If the malware has not executed (it was a dormant payload in a data file), remove the infected files and restore from backup if needed. If the malware has executed, assume the production system is compromised and follow SEC-001 incident response: isolate, assess, rebuild if necessary. Document the incident. Investigate how the scanner missed the threat and update scanning procedures.

2. **If the quarantine workstation is compromised:** Disconnect the workstation. Do not attempt to clean it. Rebuild from the known-good image. Reinstall the signature database from verified media. Verify the rebuild by running the self-test suite. Resume scanning operations. Document the compromise and its cause.

3. **If signatures cannot be updated (no access to internet-connected machine):** Continue scanning with stale signatures. Increase reliance on heuristic scanning and behavioral analysis. Document the signature staleness in the scanning log. Reduce the trust level of all scan results until signatures can be updated. Consider whether any trusted federation partner can provide a verified signature update via sync package.

4. **If heuristic false positives have reached intolerable levels:** Review and tune the heuristic rules. Identify which rules are generating the most false positives. Adjust thresholds or disable overly aggressive rules. Document each adjustment with the rationale. Re-scan recently flagged content to determine whether previously flagged items are actually clean.

5. **If a detected malware sample needs deeper analysis than the institution can perform:** Preserve the sample in the quarantine vault. Document everything known about it: detection method, behavioral observations, file metadata, source provenance. If the institution has a trusted external contact with security expertise, consider providing the sample (on dedicated media, properly labeled and handled) for external analysis. The external analysis results are imported back through the standard import pipeline.

## 8. Evolution Path

- **Years 0-5:** The scanning infrastructure is being established. The first signature database is imported. The first heuristic rules are configured. Expect the scanning process to be slow and manual. Document every scan, every detection, and every false positive. These early records will calibrate the scanning system for years to come.

- **Years 5-15:** The scanning system should be mature. Signature updates should be routine. Heuristic rules should be tuned based on years of operational experience. The behavioral sandbox should have been rebuilt multiple times as virtualization technology evolves. The scanning log should be rich enough to reveal trends: what kinds of threats are most common? What sources produce the most detections? What file formats are most risky?

- **Years 15-30:** The threat landscape will have changed significantly. Malware techniques that are common in 2026 may be obsolete. New attack vectors will have emerged. The scanning system must evolve to address new threats while retaining the ability to detect old ones (which may be re-encountered on archival media or from long-separated federation partners). The heuristic engine should be the most actively maintained component, as it is the most adaptable.

- **Years 30-50+:** The specific scanning tools may have been replaced multiple times. ClamAV may no longer exist. The principles -- signature matching, heuristic analysis, behavioral sandboxing, defense in depth -- should remain valid. Future operators should treat this article as a specification of what scanning must accomplish, not a prescription of which tools to use. The tools are the means. The principles are the end.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The honest truth about air-gapped malware scanning is that it is inferior to connected scanning. There is no way around this. A connected scanner with hourly updates and cloud-based analysis will catch threats that our offline scanner misses. We accept this inferiority because the air gap provides security benefits that outweigh the scanning gap. The air gap prevents exfiltration, prevents remote exploitation, prevents unauthorized access. The scanning gap is the price we pay for those protections.

The mitigation strategy -- defense in depth, three scanning layers, manual operator judgment -- is the best available response. It is not perfect. It will miss things. The question is not "will malware ever enter the institution?" It is "when malware enters, will we detect it quickly and contain it effectively?" The scanning system is the first line of defense. The institution's internal integrity checks (checksums, audit logs, behavioral monitoring of running systems) are the second line. Neither line is impenetrable. Together, they provide reasonable protection.

I chose ClamAV as the reference implementation because it is open source, it has been maintained for over twenty years, its signature database is freely available, and it runs on every major operating system. It is not the best scanner in terms of detection rates. Commercial scanners consistently outperform it in comparative testing. But commercial scanners require licenses, depend on proprietary infrastructure, and may cease to exist in ten years. ClamAV's open-source nature means the institution can maintain it indefinitely, fork it if necessary, and understand its operation completely. Longevity over novelty, as ETH-001 instructs.

## 10. References

- D18-001 -- Import & Quarantine Philosophy (quarantine architecture; trust scoring; import pipeline stages)
- SEC-001 -- Threat Model and Security Philosophy (three-pillar security mindset; incident response; supply chain threats)
- OPS-001 -- Operations Philosophy (quarterly review; operational sustainability; complexity budget)
- D17-003 -- Air-Gapped Synchronization Procedures (media preparation; transport media handling)
- D17-007 -- Sync Package Specification (scanning applies to incoming sync packages)
- D18-004 -- Format Conversion Pipeline (format conversion follows scanning)
- D18-005 -- Provenance Tracking and Chain of Custody (scan results as provenance data)
- ETH-001 -- Ethical Foundations (Principle 4: Longevity Over Novelty)
- Domain 19 -- Quality Assurance (scanning system audit procedures)
- Domain 20 -- Institutional Memory (malware incident records; scanning log as permanent record)

---

---

# D18-004 -- Format Conversion Pipeline

**Document ID:** D18-004
**Domain:** 18 -- Import & Quarantine
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** D18-001, D18-003, SEC-001, OPS-001, ETH-001
**Depended Upon By:** D18-005 (format conversion is a provenance event), D17-007 (sync packages carry content in institutional formats). All content that enters the institution in a non-standard format passes through this pipeline.

---

## 1. Purpose

This article defines how the institution converts imported content from external formats to institutional standard formats. D18-001 established that open formats are preferred (R-D18-06) and that the institution must be able to process content across decades of format evolution. This article provides the operational specification for achieving both goals: a format conversion pipeline that transforms incoming content into formats the institution can store, read, verify, and maintain for fifty years.

Format conversion is not merely a technical convenience. It is an act of institutional self-preservation. A document stored in a proprietary format depends on the continued existence of the proprietary software that reads it. When that software disappears -- and over fifty years, all software disappears -- the document becomes unreadable. A document converted to an open, well-documented standard format depends only on the continued existence of the standard's specification, which is a document the institution can store alongside the content itself. The conversion pipeline transforms fragile content into durable content.

This article specifies the conversion procedures for common format families, the quality verification that ensures conversion fidelity, and the format conversion registry that documents the institution's conversion capabilities and their limitations. It is addressed to the operator performing conversions, to the person maintaining the conversion tools, and to the future operator who must decide how to handle a format that did not exist when this article was written.

## 2. Scope

This article covers:

- The format conversion pipeline: its position within the import process and its operational flow.
- Institutional standard formats: the target formats for each content category.
- Conversion procedures for common format families: documents, images, audio, video, databases, and structured data.
- Quality verification after conversion: how to confirm that the converted content faithfully represents the original.
- The format conversion registry: a living document that catalogs the institution's conversion capabilities.
- Handling unconvertible formats: what to do when conversion is not possible or would result in unacceptable quality loss.
- The relationship between format conversion and long-term preservation.

This article does not cover:

- Malware scanning (see D18-003; scanning precedes conversion).
- Provenance tracking (see D18-005; conversion events are recorded in provenance).
- The philosophical justification for format preferences (see D18-001, R-D18-06 and ETH-001, Principle 4).
- Content creation within the institution (institutional content is created in standard formats from the beginning).
- Format specifications themselves (the institution should maintain copies of the relevant format specifications as reference documents).

## 3. Background

### 3.1 The Format Durability Problem

The history of digital formats is a history of obsolescence. WordPerfect files from the 1980s are largely unreadable without specialized recovery tools. Lotus 1-2-3 spreadsheets require emulators. RealMedia video files depend on a player that no longer exists. PageMaker layout files, HyperCard stacks, Flash animations -- each was once ubiquitous, and each is now an archaeological challenge.

The institution cannot afford this cycle. Content imported today must be readable in 2076. The only way to achieve this is to store content in formats whose specifications are publicly available, whose implementations exist in multiple independent software projects, and whose design prioritizes longevity over feature richness. These are the institutional standard formats.

### 3.2 The Conversion Fidelity Challenge

Format conversion is inherently lossy in some cases. A Microsoft Word document contains formatting, comments, revision tracking, embedded fonts, and macros that may not have exact equivalents in the target format. A JPEG image's lossy compression cannot be reversed to recover the original pixel data. A proprietary database format may use data types that have no standard equivalent.

The conversion pipeline must handle this reality honestly. Some conversions are lossless (or close enough to lossless that the differences are negligible). Some conversions involve measurable quality loss. Some conversions are impossible without unacceptable degradation. For each conversion, the pipeline must document what was preserved, what was lost, and whether the loss is acceptable. The operator makes the final judgment, informed by the pipeline's quality report.

### 3.3 The Dual-Preservation Strategy

For content where conversion involves any loss, the institution practices dual preservation: the original file in its original format is preserved alongside the converted file in the standard format. This ensures that if a future operator gains access to better conversion tools (or the original format's reader software), the original content is available for re-conversion. The original file is stored in a clearly marked archive area, tagged as "original format, preserved for potential future conversion." The converted file is the operational copy -- the one the institution uses, indexes, and serves.

Dual preservation consumes additional storage. This cost is accepted as insurance against conversion loss. The only exception is when the original format has been positively identified as containing executable or active content that poses a security risk (D18-003). In that case, the original may be stripped of active content before archival, with the stripping documented in provenance.

## 4. System Model

### 4.1 Pipeline Position and Flow

The format conversion pipeline operates after malware scanning (D18-003) and before admission to the production system (D18-001, Section 4.1, Stage 3: Validation). The flow is:

1. Content arrives from the quarantine scanning stage with a CLEAN designation.
2. The content's format is identified (file extension, magic bytes, MIME type analysis).
3. If the content is already in an institutional standard format, it bypasses conversion and proceeds to validation.
4. If the content is in a non-standard format, the conversion registry is consulted to determine whether a conversion path exists.
5. If a conversion path exists, the content is converted. Quality verification is performed.
6. If quality verification passes, the converted content proceeds to validation. The original is preserved per the dual-preservation strategy.
7. If no conversion path exists, or if quality verification fails, the content is flagged for operator review.

### 4.2 Institutional Standard Formats

The institution designates the following standard formats by content category. These choices are governed by four criteria: open specification (publicly documented), multiple implementations (not dependent on a single vendor), proven longevity (has existed for at least ten years), and fitness for purpose (adequate for the content type's requirements).

**Documents -- Text-Primary:**
- Primary standard: UTF-8 plain text (.txt). For content where formatting is not essential.
- Secondary standard: PDF/A-2 (.pdf). For content where visual layout, figures, and typographic formatting must be preserved. PDF/A is the archival subset of PDF, designed explicitly for long-term preservation. It prohibits features that impair archivability (embedded JavaScript, external resource dependencies, encryption).
- Tertiary standard: HTML5 with embedded CSS (.html). For content where hyperlinks, semantic structure, and moderate formatting are needed. All resources must be embedded or co-located; no external references.

**Documents -- Structured:**
- Primary standard: Markdown with YAML front matter (.md). For structured text documents, technical documentation, and institutional articles.
- Secondary standard: XML with a documented schema (.xml). For content that requires strict structural validation.

**Images:**
- Primary standard: PNG (.png). For images where lossless quality is required (diagrams, screenshots, technical illustrations).
- Secondary standard: JPEG (.jpg) with quality factor documented. For photographic images where some compression loss is acceptable.
- Tertiary standard: SVG (.svg). For vector graphics.
- Preservation standard: TIFF (.tiff). For archival-quality images, especially scanned documents.

**Audio:**
- Primary standard: FLAC (.flac). Lossless audio compression. For audio content where fidelity is important.
- Secondary standard: Opus (.opus) in Ogg container. For audio content where file size must be minimized with acceptable lossy compression. Opus provides excellent quality-to-size ratio.
- Preservation standard: WAV (.wav), 16-bit or 24-bit PCM. For archival-quality audio.

**Video:**
- Primary standard: AV1 in Matroska container (.mkv). Open, royalty-free video codec with excellent compression. For general video content.
- Secondary standard: H.264/AVC in Matroska container (.mkv). Widely supported, well-understood codec. For video content where AV1 encoding is not feasible.
- Preservation standard: FFV1 in Matroska container (.mkv). Lossless video compression designed for archival use.

**Databases and Structured Data:**
- Primary standard: SQLite (.sqlite, .db). Self-contained, serverless, public-domain database format. For structured data that requires query capability.
- Secondary standard: CSV with a companion schema document (.csv + .schema.json). For tabular data where a full database is unnecessary.
- Tertiary standard: JSON (.json). For hierarchical structured data.

**Archives:**
- Standard: TAR with Zstandard compression (.tar.zst). For bundled content. TAR for the archive structure, Zstandard for compression.

### 4.3 Conversion Procedures by Format Family

**Document Conversion:**

Microsoft Office formats (.docx, .xlsx, .pptx) are converted using LibreOffice in headless command-line mode. LibreOffice is the reference tool because it is open source, cross-platform, and produces acceptable PDF/A output for most documents. The conversion command is run in the quarantine environment. After conversion, the operator performs a visual comparison between the original (opened in the sandbox or viewed as a screenshot) and the converted PDF/A to verify layout fidelity.

Legacy document formats (.doc, .xls, .ppt, .rtf, .wpd) follow the same process. Conversion quality is generally lower for older formats. The operator must accept higher quality loss or maintain the original for future re-conversion.

E-book formats (.epub, .mobi) are converted to HTML5 (for reflowable text content) or PDF/A (for fixed-layout content). Conversion tools include Calibre (open source, well-maintained).

**Image Conversion:**

Proprietary image formats (.psd, .ai, .raw camera formats) are converted to TIFF (for preservation) and PNG or JPEG (for operational use). GIMP and ImageMagick are the reference tools. RAW camera formats are converted using dcraw or LibRaw, which are open-source RAW processing libraries. Conversion quality is verified by comparing pixel dimensions, color depth, and a visual spot-check of the converted image against the original.

WebP and HEIF/HEIC images are converted to PNG (lossless) or JPEG (lossy), depending on the source image characteristics. Lossless sources are converted to PNG. Lossy sources are converted to JPEG with a quality factor at least as high as the estimated source quality.

**Audio Conversion:**

Proprietary or uncommon audio formats (.wma, .aac in proprietary containers, .m4a) are converted to FLAC (lossless) or Opus (lossy). FFmpeg is the reference tool. For lossy-to-lossy conversion (for example, MP3 to Opus), the operator must acknowledge that the conversion introduces additional generational quality loss. The original should be preserved. Quality verification includes: comparing duration (must be identical to the millisecond), checking for encoding artifacts (pops, clicks, dropouts), and verifying that metadata (title, artist, album) is preserved or manually transferred.

**Video Conversion:**

Proprietary or uncommon video formats (.wmv, .avi with proprietary codecs, .flv, .mov with ProRes) are converted to AV1 or H.264 in Matroska containers. FFmpeg is the reference tool. Video conversion is computationally intensive; the quarantine workstation must have adequate processing capacity. Quality verification includes: comparing duration and frame count, checking for visual artifacts (blocking, color shifts, frame drops), verifying audio-video synchronization, and confirming that subtitle tracks (if present) are preserved.

For archival video, the source is converted to FFV1 lossless in Matroska container. This produces large files but preserves every frame exactly.

**Database and Structured Data Conversion:**

Proprietary database formats (.mdb, .accdb for Microsoft Access; exports from MySQL, PostgreSQL, or other server databases) are imported into SQLite. The conversion verifies: row counts match between source and target, column types are appropriate, and queries against the converted database produce the same results as queries against the source.

Spreadsheets with data-only content (no formulas, no macros) may be converted to CSV with a companion schema. The schema document records column names, data types, and any encoding notes. Quality verification compares row and column counts, and spot-checks values.

### 4.4 Quality Verification After Conversion

Quality verification is mandatory for every conversion. The verification procedure has three levels:

**Level 1: Automated verification.** The conversion tool's exit code is checked (non-zero indicates failure). File size is checked (a converted file that is dramatically smaller or larger than expected may indicate a problem). For text-based formats, character count and word count are compared. For media formats, duration, dimensions, and codec parameters are verified.

**Level 2: Structural verification.** The converted file is opened in the appropriate viewer and its structure is inspected. For documents: are all pages present? Are headings, tables, and images intact? For images: are dimensions correct? Is the color space appropriate? For audio and video: does playback start and end correctly? Is the quality subjectively acceptable?

**Level 3: Content verification.** A sample of the content is compared in detail between the original and the converted version. For a 100-page document, spot-check at least 5 pages (the first, the last, and three randomly selected). For a database, run the same query against both versions and compare results. For media, compare representative segments.

The quality verification produces a conversion quality report, which is stored alongside the converted content and becomes part of the provenance record (D18-005). The report includes: the original format, the target format, the conversion tool and version, the verification level achieved, any quality losses noted, and the operator's judgment on whether the conversion is acceptable.

### 4.5 The Format Conversion Registry

The format conversion registry is a living document maintained by the institution. It catalogs:

- Every source format the institution has encountered.
- The conversion path for each source format (which tool, which settings, which target format).
- The expected quality level for each conversion (lossless, near-lossless, acceptable lossy, significant lossy).
- Known limitations of each conversion (features that are lost, edge cases that fail).
- The date each conversion path was last tested.

The registry is updated whenever a new format is encountered, whenever a conversion tool is updated, and whenever a conversion procedure is revised. It is stored as an institutional document within Domain 18 and is included in federation sync packages so that all nodes share the same conversion knowledge.

The registry also maintains a "cannot convert" list: formats for which no acceptable conversion path exists. For these formats, the institution's options are: admit the content in its original format with a note that it may become unreadable (acceptable for low-importance content), seek or develop a conversion tool (for high-importance content), or reject the content.

### 4.6 Handling Unconvertible Formats

When content arrives in a format for which no conversion path exists, the operator follows this decision procedure:

1. **Assess importance.** Is this content critical to the institution's mission? If not, consider rejecting it and documenting the rejection reason (R-D18-04: Rejection Must Be Documented).

2. **Research alternatives.** Is there an alternative version of the same content in a convertible format? A research paper in a proprietary format may also be available as a PDF. A dataset in a proprietary database may also be available as a CSV export.

3. **Investigate conversion tools.** Is there an open-source tool that can read this format? The operator may need to import and install new software (following the standard import pipeline, including malware scanning). Document any new tools in the format conversion registry.

4. **Accept with risk notation.** If the content is important enough to admit despite the format risk, store it in its original format with a provenance record that clearly states: the format is not an institutional standard, no conversion was performed, and the content may become unreadable if the format's reader software becomes unavailable. Schedule a periodic review (at least annually) to check whether a conversion path has become available.

5. **Extract what can be extracted.** Even when full conversion is impossible, partial extraction may be feasible. Text can often be extracted from proprietary document formats, even when layout is lost. Metadata can often be extracted from media formats. Record what was extracted and what was not.

## 5. Rules & Constraints

- **R-D18-004-01: Standard Formats Are Required for Operational Content.** All content that enters the institution's operational content store must be in an institutional standard format (Section 4.2) or must have a documented exception approved by the operator with a risk notation in the provenance record.

- **R-D18-004-02: Conversion Quality Must Be Verified.** No converted content may be admitted to the production system without completing at least Level 1 and Level 2 quality verification (Section 4.4). Level 3 verification is required for content classified as high-importance.

- **R-D18-004-03: Originals Are Preserved.** When conversion involves any quality loss, the original file in its original format is preserved alongside the converted version. The original is stored in the designated archival area and tagged with its format, source, and the date of preservation.

- **R-D18-004-04: Conversion Tools Must Be Open Source.** The institution's conversion pipeline must use open-source tools whose source code is available for inspection, modification, and long-term maintenance. Proprietary conversion tools may be used only when no open-source alternative exists, and their use must be documented as a temporary measure in the format conversion registry.

- **R-D18-004-05: The Format Conversion Registry Must Be Maintained.** The registry must be reviewed and updated at least annually during the OPS-001 quarterly review cycle. Conversion paths that have not been tested within two years must be re-tested. New formats encountered since the last review must be added.

- **R-D18-004-06: Lossy Conversion Requires Operator Approval.** Automated conversion may proceed for lossless or near-lossless conversions. Any conversion that involves significant quality loss requires explicit operator approval, documented in the provenance record. The operator must understand what is lost and judge whether the loss is acceptable for this specific content.

## 6. Failure Modes

- **Conversion tool obsolescence.** The conversion tool itself becomes unmaintainable or unavailable. LibreOffice, FFmpeg, or GIMP cease development. The institution can no longer perform conversions that it previously handled routinely. Mitigation: the institution should maintain local copies of the conversion tools' source code and build environments. The format conversion registry should note dependencies on specific tools so that obsolescence risks are visible.

- **Silent quality degradation.** A conversion produces output that appears correct on casual inspection but has subtle errors: font substitution that changes character spacing, color profile conversion that shifts hues, audio resampling that introduces artifacts. Mitigation: the three-level quality verification procedure is designed to catch these issues. Level 3 (detailed content comparison) is the defense against silent degradation. It should be performed for all important content, not just for content flagged as risky.

- **Format standard evolution.** The institutional standard format itself evolves. A new version of PDF/A is released. A new lossless audio codec supersedes FLAC. The institution must decide whether to adopt the new standard, which requires re-converting existing content, or maintain the old standard, which may become less well-supported. Mitigation: standard format changes are governance decisions (D17-008 for federated institutions). They should be made deliberately, with a migration plan. Dual preservation means the originals are available for re-conversion to the new standard.

- **Conversion pipeline bottleneck.** A large import -- thousands of files in non-standard formats -- overwhelms the conversion pipeline. The operator faces a backlog that stretches for weeks. Mitigation: the pipeline should support batch conversion for common format families. The operator should triage the import: convert high-priority content first, queue lower-priority content. The OPS-001 complexity budget should account for conversion workload.

- **Format identification failure.** A file's format cannot be identified. The extension is misleading, the magic bytes are ambiguous, and the content resists analysis. Mitigation: files with unidentifiable formats should be treated as unconvertible (Section 4.6) and handled through the unconvertible format decision procedure. They should also be flagged for additional malware scrutiny (D18-003), as format obfuscation is a common malware technique.

- **Cascading conversion loss.** Content is converted through multiple intermediate formats before reaching the institutional standard, with each step introducing additional quality loss. Mitigation: the format conversion registry must specify direct conversion paths from source to target wherever possible. Multi-step conversions should be flagged and avoided unless no single-step path exists. When multi-step conversion is necessary, the original is preserved (dual preservation) and the number of steps is documented.

## 7. Recovery Procedures

1. **If a conversion is found to have introduced unacceptable quality loss after admission:** Retrieve the original file from the archival area (dual preservation). Re-perform the conversion with different settings or a different tool. If no acceptable conversion is achievable, maintain the original file as the operational copy with a format risk notation, and schedule periodic re-conversion attempts as tools improve.

2. **If a conversion tool is no longer available:** Consult the format conversion registry for alternative tools. If no alternative exists, the affected conversion paths are added to the "cannot convert" list. Content in affected formats is admitted in original format with risk notations. The institution should actively seek replacement tools through its import process.

3. **If the institutional standard format itself becomes problematic:** Convene a format review through the governance process (D17-008 for federated institutions, GOV-001 for single-node). Evaluate alternative standard formats. If a migration is approved, create a migration plan that includes: re-converting all content from the old standard to the new one, updating the format conversion registry, updating all references in institutional documentation, and preserving backward-reading capability for the old standard during a transition period.

4. **If a batch conversion produces inconsistent results:** Halt the batch. Identify which conversions succeeded and which failed or produced quality issues. Process the failures individually, adjusting settings or tools as needed. Do not resume batch processing until the cause of inconsistency is identified.

5. **If format identification is consistently failing for a new format family:** Research the format family. Obtain or develop identification rules (magic bytes, header structure). Update the format identification tools. Add the format family to the conversion registry, either with a conversion path or on the "cannot convert" list.

## 8. Evolution Path

- **Years 0-5:** The format conversion pipeline handles a limited set of well-known formats. The conversion registry is small. Most imported content is likely in common formats (PDF, JPEG, MP3, DOCX) for which established conversion paths exist. Focus on establishing the pipeline, calibrating quality verification, and building the registry.

- **Years 5-15:** The registry should be comprehensive for the formats the institution has encountered. Conversion tools may need updating. New formats may have emerged that are not yet in the registry. The dual-preservation archive may be growing significantly; assess its storage impact and ensure it is managed sustainably.

- **Years 15-30:** Some institutional standard formats may themselves face obsolescence pressure. PDF/A is well-established and likely to endure, but audio and video codecs evolve more rapidly. This is the period when the institution may face its first standard format migration. The governance process for such a migration should be initiated well before the old format becomes problematic.

- **Years 30-50+:** The format landscape will be unrecognizable compared to 2026. Formats that are ubiquitous today will be historical curiosities. The conversion pipeline's value is not in the specific tools it uses (those will have been replaced many times) but in the principles it embodies: convert to open standards, verify quality, preserve originals, document everything. These principles should guide format decisions long after the specific tools named in this article are forgotten.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The choice of institutional standard formats is the most opinionated decision in this article, and the one most likely to be revisited. PDF/A for documents, PNG and TIFF for images, FLAC for audio, AV1 for video, SQLite for databases -- these are defensible choices in 2026, but technology changes. I have tried to choose formats that optimize for longevity rather than features. PDF/A is not the most capable document format. FLAC is not the most efficient audio codec. SQLite is not the most powerful database. But each is well-documented, widely implemented, and designed (or at least well-suited) for long-term preservation.

The dual-preservation strategy is insurance against my own fallibility. If I have chosen the wrong standard format, or if the conversion tools have introduced errors I did not detect, the original is still there. A future operator with better tools or better judgment can re-convert from the original. This means the institution stores some content twice, which costs storage space. Storage is cheap. Lost content is priceless.

FFmpeg deserves special mention. It is the Swiss Army knife of media conversion -- capable of reading and writing virtually every audio and video format ever created. It is open source, it is actively maintained, and it runs everywhere. If a single tool were chosen for all media conversion, FFmpeg would be that tool. Its command-line interface is complex and poorly documented, which is its primary weakness. The format conversion registry should include specific FFmpeg commands for each conversion path, tested and verified, so that the operator does not need to become an FFmpeg expert to perform routine conversions.

## 10. References

- D18-001 -- Import & Quarantine Philosophy (R-D18-06: open format preference; import pipeline stages; trust scoring)
- D18-003 -- Malware Scanning for Air-Gapped Systems (scanning precedes conversion; active content handling)
- D18-005 -- Provenance Tracking and Chain of Custody (conversion events as provenance records)
- ETH-001 -- Ethical Foundations (Principle 4: Longevity Over Novelty; Principle 2: Integrity Over Convenience)
- SEC-001 -- Threat Model and Security Philosophy (format-based threats; active content restrictions)
- OPS-001 -- Operations Philosophy (complexity budget; quarterly review; operational sustainability)
- D17-007 -- Sync Package Specification (sync packages carry content in institutional formats)
- D17-008 -- Distributed Governance Protocol (format standard changes are governance decisions)
- Domain 19 -- Quality Assurance (conversion quality verification standards)
- Domain 20 -- Institutional Memory (format conversion registry as institutional knowledge)

---

---

# D18-005 -- Provenance Tracking and Chain of Custody

**Document ID:** D18-005
**Domain:** 18 -- Import & Quarantine
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** D18-001, D18-003, D18-004, SEC-001, OPS-001, D20-001
**Depended Upon By:** All institutional content that originated outside the institution. Referenced by D17-007 (sync package provenance), D17-005 (trust assessment relies on provenance), and Domain 19 (provenance completeness as quality metric).

---

## 1. Purpose

This article defines how the institution documents the complete history of every piece of imported content -- where it came from, how it arrived, who handled it, what transformations were applied, and what decisions were made about it at every stage of the import pipeline. This history is the content's provenance. Provenance is not optional metadata. It is the foundation on which all trust assessments rest.

D18-001 established the provenance imperative: "Provenance answers the question: 'Why should we trust this content?'" This article operationalizes that imperative by specifying the provenance record format, the chain of custody protocol, the verification procedures, and the provenance audit that ensures the institution's provenance system remains complete and honest.

A piece of content without provenance is an orphan. The institution does not know where it came from, how it was obtained, who vetted it, or what happened to it between creation and admission. An orphan may be perfectly legitimate. It may also be fabricated, corrupted, or malicious. Without provenance, the institution cannot distinguish between these possibilities. With provenance, it can make informed trust decisions.

This article is addressed to every person who imports content into the institution, every person who transforms or processes imported content, and every future operator who inherits institutional content and needs to understand its history.

## 2. Scope

This article covers:

- The provenance record format: the complete specification of a provenance record.
- Chain of custody tracking: how to document every handler and every transfer.
- Provenance events: the taxonomy of events that generate provenance entries.
- Provenance verification: how to confirm that a provenance record is complete and consistent.
- The provenance audit: the periodic review of provenance completeness across the institution.
- Provenance for federated content: how provenance works for content received from other nodes.
- Provenance integrity: how to protect provenance records from tampering.

This article does not cover:

- The import pipeline procedures (see D18-001).
- Malware scanning details (see D18-003; scan results feed into provenance).
- Format conversion details (see D18-004; conversion events feed into provenance).
- The philosophical justification for provenance (see D18-001, Section 3.3).
- Provenance of internally created content (covered by content-creation domain articles; this article covers only imported content).
- Institutional memory architecture (see D20-001; provenance records are stored within the institutional memory framework).

## 3. Background

### 3.1 What Provenance Means for Trust

D18-001 defined four dimensions of trust: Source Trust, Carrier Trust, Format Trust, and Content Trust. Provenance is the evidentiary basis for the first two dimensions. Source Trust answers "who created this content?" -- and the only way to answer that question with any confidence is to have a record that traces the content back to its origin. Carrier Trust answers "how did this content get here?" -- and the only way to answer that is to document the chain of custody from origin to admission.

Without provenance, trust assessment is guesswork. With provenance, trust assessment is evidence-based. The evidence may be incomplete (provenance rarely captures every moment of a file's existence), but incomplete evidence is infinitely better than no evidence. Provenance does not guarantee trustworthiness. It provides the information needed to make a reasoned judgment about trustworthiness.

### 3.2 The Chain of Custody Concept

Chain of custody is a concept borrowed from legal evidence handling. In a courtroom, evidence is admissible only if the prosecution can demonstrate an unbroken chain of custody: who collected the evidence, how it was stored, who handled it at each stage, and how it was protected from tampering. Any break in the chain introduces reasonable doubt about the evidence's integrity.

The institution applies the same concept to imported content. The chain of custody documents every handler (human or automated), every location (physical or digital), and every transformation (copy, conversion, scan) between the content's origin and its admission to the institutional production system. A complete chain of custody means the institution can account for the content's entire journey. A break in the chain means there is a period during which the institution cannot verify what happened to the content. Breaks are not fatal -- they are common, especially for content obtained from the internet -- but they reduce trust.

### 3.3 Provenance Is Append-Only

Provenance records, like decision logs (D20-001), are append-only. An event that has been recorded cannot be unrecorded. If an error is discovered in a provenance entry, a correction entry is appended that references and corrects the original. The original entry remains. This append-only constraint protects against revisionism: no one can alter a content's history after the fact. The provenance record is what it is. Errors are corrected by addition, not by modification.

This constraint also means that provenance records grow monotonically. Every event adds to the record; nothing removes from it. Over decades, the provenance records for long-lived content can become substantial. This growth is acceptable and expected. The cost of storing provenance is trivial compared to the cost of losing it.

## 4. System Model

### 4.1 The Provenance Record Format

Every piece of imported content has exactly one provenance record. The provenance record is a structured document containing a header and a chronological sequence of provenance events. The record is stored as a UTF-8 text file in a defined schema.

**Provenance Record Header:**

- `Provenance-ID`: A unique identifier for this provenance record. Format: `PROV-<content-hash>-<import-date>`, where content-hash is the first 16 characters of the SHA-256 hash of the content as first received by the institution.
- `Content-ID`: The institutional content identifier assigned to the content after admission.
- `Content-Hash-Original`: The SHA-256 hash of the content as first received, before any transformation.
- `Content-Hash-Current`: The SHA-256 hash of the current version of the content (updated after each transformation).
- `Record-Created`: ISO 8601 UTC timestamp of record creation.
- `Record-Version`: The provenance record format version. Currently: `1`.

**Provenance Event Format:**

Each event in the record contains:

- `Event-Sequence`: An integer, starting at 1, incrementing by 1 for each event. Provides unambiguous ordering.
- `Event-Type`: One of the event types defined in Section 4.2.
- `Event-Timestamp`: ISO 8601 UTC timestamp of the event.
- `Event-Actor`: The identifier of the person or system that performed the action. For human actors: name or operator identifier. For automated systems: tool name and version.
- `Event-Location`: Where the event occurred. For physical events: site identifier. For digital events: system identifier (quarantine workstation, production system, etc.).
- `Event-Description`: A human-readable description of what happened.
- `Event-Input-Hash`: The SHA-256 hash of the content before this event (if the event modifies the content).
- `Event-Output-Hash`: The SHA-256 hash of the content after this event (if the event modifies the content).
- `Event-Evidence`: References to supporting documentation. For scan events: the scan report identifier. For conversion events: the quality verification report identifier. For decisions: the decision log entry identifier.
- `Event-Notes`: Optional free-text notes by the actor.

### 4.2 Provenance Event Taxonomy

The following event types are defined. Each represents a distinct action in the content's lifecycle:

**ORIGIN.** The content was created at its source. This event is recorded based on the best information available to the importer. It may be precise ("downloaded from [URL] on [date] by [person]") or approximate ("obtained from a USB drive labeled [label] given to us by [person] who stated it originated from [source]"). The origin event is often the least certain event in the chain, and that uncertainty should be honestly documented.

**ACQUISITION.** The content was obtained by the person who will carry it to the institution. This event marks the transition from external custody to the importer's custody. It records: who acquired the content, from where, on what date, and on what medium.

**TRANSPORT.** The content was physically moved from one location to another. Each physical movement is a separate transport event. If the content was carried on a USB drive from a home office to the institution's site, that is one transport event. If it was shipped via mail, that is another. The transport event records: departure location, departure date, arrival location, arrival date, transport method, and the identity of the carrier.

**INTAKE.** The content arrived at the institution's boundary and was received into the quarantine system. This event marks the transition from external handling to institutional handling. It records: the date, the receiving operator, the media type, and the initial file inventory.

**SCAN.** The content was scanned for malware (D18-003). The scan event records: the scanning tool and version, the signature database version, the heuristic engine configuration, the scan result (CLEAN, FLAGGED, DETECTED, SUSPICIOUS), and a reference to the full scan report.

**CONVERSION.** The content was converted from one format to another (D18-004). The conversion event records: the source format, the target format, the conversion tool and version, the conversion settings, the quality verification level achieved, and references to the conversion quality report. The Event-Input-Hash and Event-Output-Hash fields are critical for this event type, as they link the pre-conversion and post-conversion versions.

**REVIEW.** The content was manually reviewed by a human operator. The review event records: who performed the review, what aspects were reviewed (content accuracy, format suitability, relevance to institutional mission), and the review outcome (approved, rejected, approved with conditions).

**DECISION.** A decision was made about the content's disposition. The decision event records: the decision (admit, reject, defer, quarantine further), the rationale, and the decision-maker. This event links to the decision log (D20-001, GOV-001).

**ADMISSION.** The content was admitted to the institutional production system. This event marks the transition from quarantine to production. It records: the date, the admitting operator, the final content hash, and the institutional content identifier assigned.

**MODIFICATION.** The content was modified after admission. This event is used for post-admission changes: corrections, updates, annotations. Each modification is a separate event with its own input and output hashes.

**FEDERATION-RECEIVE.** The content was received from a federated node via sync package. This event records: the sending node identifier, the sync package identifier (D17-007), the trust level of the sending node at the time of receipt, and any provenance information included by the sending node.

**FEDERATION-SEND.** The content was sent to a federated node via sync package. This event records: the receiving node identifier, the sync package identifier, and the date of transmission.

### 4.3 Chain of Custody Tracking

The chain of custody is the subset of provenance events that track the physical and digital custody of the content. It is the sequence of ORIGIN, ACQUISITION, TRANSPORT, INTAKE, and ADMISSION events -- the events that describe where the content was and who was responsible for it at each moment.

A complete chain of custody has no temporal gaps. For every moment between ORIGIN and ADMISSION, the provenance record identifies who had custody. In practice, complete chains are rare for externally sourced content. The content's life before ACQUISITION is often undocumented. The origin may be approximate. The chain between creation and acquisition may include multiple unknown handlers.

The institution records the chain as completely as possible and explicitly notes gaps. A gap entry has the following format:

- `Event-Type`: `CUSTODY-GAP`
- `Event-Description`: "Chain of custody is unknown between [start condition] and [end condition]."
- `Event-Notes`: Any information that mitigates the gap (e.g., "the content's cryptographic signature was verified, providing confidence that it was not modified during the gap despite unknown handling").

Custody gaps reduce trust. A content item with a complete chain of custody receives higher Carrier Trust than one with gaps. The size and nature of the gaps factor into the trust assessment.

### 4.4 Provenance for Federated Content

Content received from a federated node presents a special provenance challenge. The content may have its own provenance record on the sending node, which includes events that occurred before the content arrived at the receiving node. The receiving node should obtain and incorporate this upstream provenance.

The process is:

1. When a federated node sends content in a sync package, the content's provenance record is included alongside the content (as an Item-Type: METADATA entry in the sync package manifest, per D17-007).

2. The receiving node creates a new provenance record for the content, beginning with a FEDERATION-RECEIVE event.

3. The upstream provenance (from the sending node) is appended to the receiving node's record as a clearly demarcated section: "Upstream Provenance from Node [identifier]." The upstream provenance is treated as a claim by the sending node, not as verified fact. Its trustworthiness depends on the trust level of the sending node.

4. Subsequent events on the receiving node (local scanning, conversion, review, admission) are added to the local provenance record in the normal way.

This produces a composite provenance record: the receiving node's own observations plus the sending node's claimed history. The distinction between observed and claimed provenance must be maintained. The receiving node can verify its own entries. It cannot independently verify the sending node's entries (though it can assess their plausibility).

### 4.5 Provenance Integrity Protection

Provenance records must be protected from tampering with the same rigor as institutional memory records (D20-001). The protection mechanisms are:

**Append-only storage.** Provenance records are stored in an append-only format. The storage system must prevent modification of existing entries. On filesystems that do not natively support append-only files, application-level controls enforce the constraint, backed by periodic integrity verification.

**Hash chaining.** Each provenance event includes a hash of the previous event (similar to a blockchain). The first event's "previous hash" field is the content's original hash (Content-Hash-Original). Each subsequent event's previous hash is the SHA-256 hash of the previous event's complete text. This creates a tamper-evident chain: modifying any event invalidates the hash chain from that point forward.

**Periodic integrity verification.** During the quarterly review (OPS-001), the integrity of all provenance records is verified by re-computing the hash chain for each record and comparing against the stored hashes. Any discrepancy triggers a security incident investigation per SEC-001.

**Backup replication.** Provenance records are replicated to backup media with the same priority as the content they describe. A backup that preserves content but not its provenance is an incomplete backup.

### 4.6 The Provenance Audit

The provenance audit is a periodic review of provenance completeness and quality across all institutional content. It is performed at least annually, aligned with the OPS-001 quarterly review cycle (the annual audit is more comprehensive than the quarterly integrity check).

The audit examines:

**Completeness.** Does every piece of imported content have a provenance record? Are there content items in the production system without corresponding provenance records (orphans)? Are there provenance records without corresponding content items (ghosts)?

**Coverage.** For each provenance record, are all required event types present? Every imported item should have, at minimum: INTAKE, SCAN, DECISION, and ADMISSION events. Missing events indicate that the import pipeline was not followed correctly.

**Chain integrity.** Are all hash chains valid? Do the Event-Input-Hash and Event-Output-Hash fields correctly link transformation events? Can the content's current hash be traced back through its transformation chain to its original hash?

**Gap assessment.** How many provenance records contain CUSTODY-GAP events? What proportion of the institution's content has incomplete chains of custody? Is the proportion increasing or decreasing over time?

**Trust correlation.** Do the trust scores assigned to content during import correlate with the provenance quality? Content with poor provenance should have lower trust scores. If high-trust content has poor provenance, the trust scoring model may need recalibration.

The audit produces a provenance health report, stored in institutional memory (Domain 20). The report includes summary statistics, identified issues, and recommended actions.

## 5. Rules & Constraints

- **R-D18-005-01: Every Import Requires a Provenance Record.** No content may be admitted to the institutional production system without a provenance record. The record must be created at or before the INTAKE event and must be maintained through admission. Content without provenance is not admitted.

- **R-D18-005-02: Provenance Records Are Append-Only.** Provenance events may not be modified or deleted after creation. Corrections are made by appending correction events that reference the original entry. The hash chain (Section 4.5) enforces this constraint cryptographically.

- **R-D18-005-03: Gaps Must Be Documented.** When the chain of custody cannot be determined for a period of the content's history, a CUSTODY-GAP event must be explicitly recorded. Gaps may not be papered over with fabricated events. Honest gaps are preferable to dishonest completeness.

- **R-D18-005-04: Transformation Events Must Include Hashes.** Every event that modifies the content (CONVERSION, MODIFICATION) must include both the Event-Input-Hash and Event-Output-Hash. These hashes create the verifiable transformation chain that links the current content to its original form.

- **R-D18-005-05: Provenance Travels with Content.** When content is included in a sync package for federation (D17-007), its provenance record must be included alongside it. A receiving node must be able to review the content's history before deciding whether to admit it.

- **R-D18-005-06: Upstream Provenance Is Claimed, Not Verified.** Provenance information received from a federated node is recorded as the sending node's claim. The receiving node notes the trust level of the sending node alongside the upstream provenance. Upstream provenance is not independently verifiable and must not be presented as though it were.

- **R-D18-005-07: Provenance Integrity Is Verified Quarterly.** The hash chain integrity of all provenance records must be verified at least quarterly as part of the OPS-001 review cycle. Any integrity failure is treated as a potential security incident.

## 6. Failure Modes

- **Provenance neglect.** The operator skips provenance recording because it is time-consuming and the content seems obviously trustworthy. Over time, the institution accumulates content without provenance. Trust assessments become unreliable. Mitigation: R-D18-005-01 makes provenance mandatory. The provenance audit (Section 4.6) detects orphaned content. The quarterly review includes a provenance completeness check. Automation should handle as much of the provenance recording as possible (the scanning system automatically generates SCAN events, the conversion pipeline automatically generates CONVERSION events, the admission system automatically generates ADMISSION events). The operator's burden should be limited to recording the events that require human knowledge: ORIGIN, ACQUISITION, TRANSPORT, and REVIEW.

- **Provenance fabrication.** The provenance record is inaccurate -- either through carelessness (recording the wrong source, the wrong date) or deliberate deception (fabricating a prestigious origin for content of unknown provenance). Mitigation: provenance records are append-only and hash-chained, making after-the-fact fabrication detectable. Real-time fabrication (recording false information at the time of import) is harder to detect. The defense is institutional culture: ETH-001's emphasis on honesty and transparency, combined with the understanding that fabricated provenance corrodes the trust system that protects the institution. The provenance audit may detect implausible records (a content item whose provenance claims a source that did not exist at the claimed date, for example).

- **Provenance overload.** The provenance system generates so much data that it becomes unmanageable. Every minor event creates a provenance entry. The provenance records are longer than the content they describe. Mitigation: the event taxonomy (Section 4.2) is deliberately limited to significant events. Routine system actions (file copies within the production system, index rebuilds, cache refreshes) are not provenance events. Only events that affect the content's identity, integrity, custody, or trust status generate provenance entries.

- **Hash chain corruption.** A storage error corrupts a provenance record, breaking the hash chain. The integrity verification detects the break but cannot determine whether the corruption was accidental or deliberate. Mitigation: provenance records are replicated to backup media. If a hash chain break is detected, the backup copy is compared to the damaged copy. If the backup is intact, the damaged copy is replaced. If both copies are damaged at the same point, a deeper investigation is warranted (coincident corruption at the same event is suspicious).

- **Upstream provenance manipulation.** A federated node sends content with fabricated upstream provenance -- claiming a source, a chain of custody, or a scan history that is false. The receiving node has no way to independently verify the upstream provenance and may make trust decisions based on false information. Mitigation: R-D18-005-06 requires upstream provenance to be labeled as claimed, not verified. The trust level of the sending node (D17-005) is the primary defense: a Level 3 Trusted Partner's provenance claims carry more weight than a Level 1 node's claims. Persistent upstream provenance inconsistencies should affect the sending node's trust level.

- **Provenance-content desynchronization.** The provenance record and the content it describes become disconnected -- the content is moved, renamed, or modified without a corresponding provenance event. The record describes a content item that no longer exists at the recorded location, or a content item exists without a linked provenance record. Mitigation: the Content-ID and Content-Hash fields in the provenance header provide the link between provenance and content. The provenance audit detects desynchronization (orphans and ghosts). The institution's content management procedures must ensure that provenance and content are moved, backed up, and deleted together.

## 7. Recovery Procedures

1. **If content is discovered without provenance (orphan content):** Create a provenance record retroactively. The first event should be a RECONSTRUCTION event (a special-purpose event type used only for retroactive documentation) that documents: when the content was discovered to be without provenance, what is known about its origin (if anything), and who performed the reconstruction. The record should be marked as reconstructed, not original. Reconstructed provenance has lower evidentiary value than original provenance, and this should be reflected in the content's trust score.

2. **If a provenance hash chain is broken:** Retrieve the backup copy of the provenance record. If the backup is intact, replace the damaged record. If the backup is also damaged, reconstruct the chain to the extent possible from other sources: scan logs, conversion logs, decision logs, and institutional memory entries. Mark the reconstructed portion of the chain as reconstructed.

3. **If provenance fabrication is discovered:** Treat it as a security incident per SEC-001. Identify all content affected by the fabricated provenance. Re-assess the trust scores for all affected content. If the fabrication is limited to one source or one importer, quarantine all content from that source or importer and re-evaluate. Document the fabrication in the decision log and in this article's Commentary Section. Update import procedures to prevent recurrence (additional verification steps, second-person review for provenance from the affected source).

4. **If provenance records have become desynchronized from content:** Run the provenance audit to identify all orphans (content without provenance) and ghosts (provenance without content). For orphans, create reconstructed provenance records. For ghosts, investigate: was the content deleted? Was it moved? Was the provenance record created for content that was never admitted? Update the ghost records with the findings. Implement tighter coupling between content operations and provenance operations to prevent future desynchronization.

5. **If the provenance audit reveals systemic incompleteness:** Identify the pipeline stage where provenance recording is failing. Common causes: a tool in the pipeline does not automatically generate provenance events, an operator is skipping a step, or a recently introduced process does not have provenance integration. Fix the root cause. Retroactively reconstruct provenance for content that passed through the broken stage, to the extent possible.

## 8. Evolution Path

- **Years 0-5:** Provenance tracking is being established. The early provenance records will be crude -- manually created, possibly incomplete, with many custody gaps for content obtained from the internet. This is acceptable. The discipline of creating provenance records is more important than the records' completeness. As the import pipeline matures, automation will handle more of the provenance recording.

- **Years 5-15:** Provenance should be well-integrated into the import pipeline. Most events should be generated automatically. The provenance audit should be routine. The institution should have enough provenance data to correlate provenance quality with content reliability: does content with better provenance tend to be more trustworthy? This correlation validates the provenance system.

- **Years 15-30:** The provenance records from the founding era are now historical documents in their own right. They document not just the content's history but the institution's early practices. Future operators will read these records and learn how the institution's import procedures evolved. The records should be treated as institutional memory, not just content metadata.

- **Years 30-50+:** The provenance system should be one of the institution's most mature and reliable systems. It has been running for decades, generating records for every import, and surviving multiple operator successions. The hash chain provides a tamper-evident record stretching back to the institution's founding. This is one of the institution's most valuable assets -- an unbroken, verifiable record of how every piece of external content entered the institution and what happened to it. Protect it accordingly.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
Provenance tracking is the least glamorous part of the import process and the most important. Nobody enjoys filling out provenance records. The temptation to skip the origin documentation, to leave the chain of custody vague, to not bother with the transport event because "it was just a USB drive from my desk to the quarantine workstation" -- these temptations will be constant. Resist them. Every gap in provenance is a gap in trust. Every skipped event is a question that a future operator cannot answer.

The hash chaining mechanism deserves explanation. I considered and rejected a simpler approach: just store the provenance events without cryptographic linking. The argument for simplicity was that provenance records are stored on a system that is already access-controlled, so tampering is unlikely. The argument for hash chaining was that "unlikely" is not "impossible," and the cost of hash chaining is trivial (one SHA-256 computation per event), while the cost of undetectable provenance tampering is potentially catastrophic. A provenance record that can be silently altered is a provenance record that cannot be trusted. Hash chaining makes silent alteration detectable. The marginal cost is near zero. The marginal benefit is the integrity of the entire trust system.

I also want to note the philosophical tension between provenance tracking and privacy. If content was obtained from a person, the provenance record documents that person's involvement. In some contexts, this documentation could be sensitive. The institution must balance the need for complete provenance against any privacy obligations it has to its sources. This balance is context-dependent and must be navigated by the operator's judgment, informed by ETH-001's ethical principles.

## 10. References

- D18-001 -- Import & Quarantine Philosophy (provenance imperative; trust scoring model; import pipeline stages)
- D18-003 -- Malware Scanning for Air-Gapped Systems (scan events as provenance entries)
- D18-004 -- Format Conversion Pipeline (conversion events as provenance entries; dual preservation)
- D20-001 -- Institutional Memory Philosophy (append-only records; anti-revisionism; permanent record tier)
- SEC-001 -- Threat Model and Security Philosophy (integrity verification; incident response for tampering)
- ETH-001 -- Ethical Foundations (Principle 2: Integrity Over Convenience; Principle 3: Transparency)
- OPS-001 -- Operations Philosophy (quarterly review; provenance verification in operational tempo)
- D17-005 -- Trust Establishment and Verification (trust levels affect provenance credibility)
- D17-007 -- Sync Package Specification (provenance travels with content in sync packages)
- GOV-001 -- Authority Model (decision records linked from provenance DECISION events)
- Domain 19 -- Quality Assurance (provenance completeness as quality metric; provenance audit standards)

---

---

*End of Stage 4 Specialized Systems -- Federation & Import Advanced Reference*

**Document Total:** 5 articles
**Domains Covered:** 17 (Scaling & Federation): D17-007, D17-008. 18 (Import & Quarantine): D18-003, D18-004, D18-005.
**Combined Estimated Word Count:** ~14,500 words
**Status:** All five articles ratified as of 2026-02-16.
**Next Steps:** Domain-specific implementation guides derived from these specifications. Cross-domain integration testing between federation sync packages (D17-007), governance protocol (D17-008), malware scanning (D18-003), format conversion (D18-004), and provenance tracking (D18-005).
