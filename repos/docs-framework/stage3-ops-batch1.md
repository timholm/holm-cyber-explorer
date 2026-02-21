# STAGE 3: OPERATIONAL DOCTRINE (BATCH 1)

## Actionable Procedures for Security, Data Integrity, and Incident Response

**Document ID:** STAGE3-OPS-BATCH1
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Classification:** Operational Procedures -- These articles translate Stage 2 philosophy into step-by-step actionable procedures. They are designed to be followed by a single operator with no one to ask.

---

## How to Read This Document

This document contains five operational articles that belong to Stage 3 of the holm.chat Documentation Institution. Stage 2 established the philosophy -- what we believe, why we believe it, and how we think about problems. Stage 3 turns that philosophy into procedures -- what to do, when to do it, and how to verify it was done correctly.

These are manuals. They are written for the person who needs to do the work, not the person who wants to think about the work. They assume you have read the relevant Stage 2 philosophy articles and understood the principles. They do not re-argue the principles. They implement them.

If you are a future maintainer following these documents: every step is written to be self-contained. Where a step requires context, the context is provided or referenced. Where a step requires judgment, the criteria for that judgment are stated. Where a step could go wrong, the recovery path is documented. You should be able to execute every procedure in this document from the documentation alone, at Stage 2 competence per D9-001.

If something in these procedures does not work -- because hardware has changed, because software has evolved, because a tool referenced here no longer exists -- do not abandon the procedure. Adapt it. The principles behind each procedure are stated in the Background section. Use those principles to find the equivalent procedure for your reality. Then update this document.

---

---

# SEC-002 -- Access Control Procedures

**Document ID:** SEC-002
**Domain:** 3 -- Security & Integrity
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, SEC-001, OPS-001
**Depended Upon By:** SEC-003, D6-002, D6-003, D10-003. All articles involving user accounts, permissions, or physical access.

---

## 1. Purpose

This article defines the practical access control procedures for the holm.chat Documentation Institution. It specifies who gets access to what, how access is granted, how access is revoked, how access decisions are recorded, and how the institution verifies that access controls remain effective over time.

Access control in an air-gapped, single-operator institution is different from access control in a corporate environment. There is no Active Directory. There is no SSO provider. There is no IT department to submit a ticket to. The operator is the identity provider, the access administrator, the auditor, and the primary user -- all simultaneously. This article acknowledges that reality while maintaining the discipline that SEC-001 demands.

The procedures in this article implement three principles from SEC-001: the trust model (Section 4.3), which defines what and whom the institution trusts; the security mindset (Section 4.1), which demands defense in depth; and the human-is-the-system pillar, which recognizes that the operator's discipline is the most critical access control.

This article is written for the operator who needs to set up access controls from scratch, maintain them over years, and hand them to a successor who has never seen the system before.

## 2. Scope

This article covers:

- User account architecture: what accounts exist, what they are for, and how they relate to each other.
- Permission models: how permissions are assigned and what the principle of least privilege means in practice.
- Physical access controls: locks, keys, placement, and physical security of hardware.
- Account creation procedures: step-by-step process for creating a new account.
- Account modification procedures: how to change permissions.
- Account revocation procedures: how to disable or remove access.
- Audit trail requirements: what gets logged, where, and how to review it.
- Periodic access review: scheduled verification that access controls match intent.

This article does not cover:

- Cryptographic key management (see SEC-003).
- Network security (not applicable; the institution is air-gapped per SEC-001).
- Encryption at rest (see SEC-003 for key management; specific encryption procedures in subsequent Domain 3 articles).
- Incident response for access breaches (see D10-003).

## 3. Background

### 3.1 Why Access Control Matters for a Single Operator

The objection is immediate: "I am one person. Why do I need access controls?" The answer has four parts.

First, separation of privilege protects against operator error. When your daily-use account cannot modify system configurations, you cannot accidentally destroy a critical service while trying to edit a document. The air gap protects against external attackers. Privilege separation protects against the most common attacker: yourself, distracted or tired or in a hurry.

Second, access controls create accountability. When every action is traceable to a specific account with specific privileges, the audit trail tells you what happened and in what context. An action performed by the root account has different implications than the same action performed by a service account.

Third, succession requires it. When a new operator assumes responsibility, they need their own accounts with appropriate access. The access control architecture must support this transition. If the only account is root and the only password is in the founder's head, succession is a crisis, not a transition.

Fourth, defense in depth demands it. SEC-001 Pillar 1 states: assume breach, prevent breach. If any single layer is compromised -- a password is guessed, a session is left open, a physical key is lost -- the damage must be contained by other layers. Access controls are one of those layers.

### 3.2 The Principle of Least Privilege

Every account, every process, and every service should have exactly the permissions it needs to perform its function and no more. This is the principle of least privilege, and it is the governing principle of all access decisions in this institution.

Least privilege is inconvenient. It means you will sometimes need to switch accounts to perform a different task. It means services will sometimes fail because they lack a permission you forgot to grant. It means access requests require thought rather than rubber-stamping. This inconvenience is the cost of defense in depth. Per ETH-001 Principle 2, integrity prevails over convenience.

### 3.3 Physical Access as the First Layer

In an air-gapped institution, physical access is the perimeter. There is no network to breach. The first and most important access control is: who can physically touch the hardware? All other access controls assume the physical perimeter is intact. If someone has unauthorized physical access, all software access controls can be bypassed given sufficient time and skill. This is why physical security is addressed first in this article, not as an afterthought.

## 4. System Model

### 4.1 The Account Architecture

The institution maintains five categories of accounts, each with a distinct purpose and privilege level:

**Category 1: The Root Account.**
The superuser account with unrestricted access to all systems. This account can modify anything, delete anything, and bypass all software restrictions. It is the most dangerous account in the institution.

Usage rules:
- The root account is used only for system-level tasks that no other account can perform: initial system setup, disk encryption management, bootloader modifications, kernel updates, and emergency recovery.
- The root account is never used for daily operations. Never.
- Root sessions are logged with enhanced detail (every command, every file access).
- Root access requires a deliberate authentication step separate from other accounts (a different password, stored separately, per SEC-003).
- After every root session, the operator documents what was done and why in the operational log.

**Category 2: The Administrative Account.**
A non-root account with elevated privileges, typically through sudo or equivalent mechanism. This account can perform system administration tasks -- installing software, modifying configurations, managing services -- with explicit authorization for each elevated action.

Usage rules:
- The administrative account is used for system maintenance, software updates, configuration changes, and account management.
- Each elevated action requires explicit authentication (password entry for each sudo invocation; do not configure passwordless sudo).
- The administrative account is not used for daily work such as document editing, data entry, or routine file management.

**Category 3: The Operator Account.**
The daily-use account for the person operating the institution. This account has access to operational data, documentation, daily tools, and the operational log. It cannot modify system configurations, install software, or access other accounts' files without explicit elevation.

Usage rules:
- This is the default account for daily operations per D10-002.
- The operator account can read and write to the operational data directories, the documentation corpus, and the operational log.
- The operator account cannot modify system configurations, install or remove software, or access the root account's files.
- The operator account is the identity recorded in the operational log for daily operations.

**Category 4: Service Accounts.**
Non-interactive accounts that run automated services -- backup processes, integrity verification, monitoring daemons. Each service account has precisely the permissions needed for its function and no more.

Usage rules:
- One service account per service. Do not share service accounts between services.
- Service accounts cannot log in interactively. They are used only to run their designated service.
- Service account permissions are documented in the service's operational documentation per D8-001.
- Service accounts are reviewed during the quarterly access audit.

**Category 5: The Succession Account.**
A dormant account prepared for the designated successor per GOV-001 Section 7.1. This account is created and configured but disabled until the succession event occurs.

Usage rules:
- The succession account is created during the succession planning process.
- It is configured with operator-level permissions (Category 3).
- It is disabled (locked) until a succession event is initiated.
- Its existence and credentials are documented in the succession packet per GOV-001.
- The succession account is verified annually to confirm it can be activated.

### 4.2 The Permission Model

Permissions are implemented through the operating system's native access control mechanisms (Unix file permissions, POSIX ACLs, or equivalent). The permission model follows three rules:

**Rule 1: Default deny.** New files and directories are created with the minimum permissions necessary. The default umask is set to 077 (owner access only) for all accounts.

**Rule 2: Group-based access for shared resources.** Where multiple accounts need access to the same resource (e.g., the operator account and a backup service account both need to read the documentation corpus), access is granted through group membership rather than world-readable permissions.

**Rule 3: No world-readable files on sensitive data.** Files classified as Tier 1 or Tier 2 per D6-001 are never world-readable. Access is always through explicit user or group permissions.

### 4.3 The Physical Access Model

Physical access is controlled in three concentric zones:

**Zone 1: The Equipment Zone.** The physical space containing the institution's servers, storage, and networking equipment (even if disabled per SEC-001). This zone is secured by a locked door with a physical key or combination lock. Only the operator and designated successors should have access. The key or combination is documented in the succession packet.

**Zone 2: The Work Zone.** The physical space where the operator interacts with the institution's systems -- the desk, monitor, keyboard, and any peripheral devices. This zone may be the same room as Zone 1 or an adjacent space. When the operator leaves the work zone, the screen must be locked and any removable media must be secured.

**Zone 3: The Storage Zone.** The physical location of backup media, paper documentation, and the succession packet. This zone may be in the same building as Zones 1 and 2, or at a geographically separate location for disaster resilience per D6-001 R-D6-02. This zone is secured independently of Zones 1 and 2.

## 5. Rules & Constraints

- **R-SEC2-01:** The root account shall not be used for any task that can be accomplished by another account. Violations must be documented as security incidents in the operational log.
- **R-SEC2-02:** Every human account must have a unique, strong password managed per SEC-003. No password may be shared between accounts.
- **R-SEC2-03:** Every service account must be documented with its purpose, its permissions, and the service it runs. Undocumented service accounts must be disabled until documented.
- **R-SEC2-04:** The default umask for all accounts shall be 077 unless a documented operational requirement necessitates a more permissive setting.
- **R-SEC2-05:** Physical access to Zone 1 must be controlled by a lock whose key or combination is known only to authorized personnel. Changes to physical access must be documented in the decision log as a Tier 3 decision per GOV-001.
- **R-SEC2-06:** When the operator leaves the work zone, the screen must be locked. This is a habit, not a judgment call.
- **R-SEC2-07:** Access controls must be audited quarterly. The audit verifies that all accounts are documented, all permissions are justified, all service accounts are still needed, and no unauthorized accounts exist.
- **R-SEC2-08:** The succession account must be verified annually by testing that it can be activated and that its permissions are correct.
- **R-SEC2-09:** All account creation, modification, and revocation must be recorded in the institutional decision log.
- **R-SEC2-10:** Authentication logs must be preserved as Tier 1 data per D6-001.

## 6. Failure Modes

- **Root account creep.** The operator starts using root for convenience. "I will just quickly..." becomes the norm. Every action runs with maximum privileges, and the audit trail loses its value because everything is done as root. Mitigation: R-SEC2-01 is absolute. Configure the root account's shell prompt to be visually distinctive (red, or with a warning label) so that operating as root is always visible and uncomfortable.

- **Password reuse.** The operator uses the same password for multiple accounts because unique passwords are annoying to remember. A single compromise exposes everything. Mitigation: SEC-003 defines the password management strategy. Each account has a unique password stored in the secure credential store.

- **Orphaned accounts.** A service is decommissioned but its service account remains. Over time, orphaned accounts accumulate, each a potential vector for confusion or misuse. Mitigation: the quarterly access audit (R-SEC2-07) specifically checks for accounts whose associated service no longer exists.

- **Physical access complacency.** The door is propped open. The screen is not locked. Backup media sits on an unlocked shelf. Over time, the physical perimeter dissolves. Mitigation: R-SEC2-05 and R-SEC2-06 are non-negotiable. The daily operational checklist per D10-002 includes a physical security item.

- **Succession account rot.** The succession account is created once and never tested. When the succession event occurs, the account does not work -- the password has been forgotten, the home directory was never created, or the permissions are wrong. Mitigation: R-SEC2-08 mandates annual verification.

- **Permission drift.** Permissions are loosened to solve a problem and never tightened again. Over years, the careful permission model erodes into effective open access. Mitigation: the quarterly audit includes a permission review against the documented permission model.

## 7. Recovery Procedures

**7.1 If root has been used for routine operations:**
1. Immediately stop using root for non-root tasks.
2. Review the root session logs for the period of misuse.
3. Identify any changes made as root that should have been made by a less-privileged account.
4. Verify that no unintended file ownership or permission changes resulted from root usage.
5. Run: `find / -user root -newer [reference_timestamp] -not -path '/proc/*' -not -path '/sys/*'` to identify files recently modified by root.
6. Correct any ownership or permission issues.
7. Document the incident in the operational log.
8. Add a root-usage check to the weekly review: examine root session logs for the week and flag any non-essential root usage.

**7.2 If an unauthorized account is discovered:**
1. Do not delete it immediately. It may contain evidence of how it was created.
2. Disable the account: `usermod -L [account_name]`
3. Record the discovery in the operational log.
4. Examine the account's creation date, home directory, shell, and group memberships.
5. Review system logs to determine how and when the account was created.
6. Assess whether the account was used for any actions. Review auth logs for login attempts.
7. If the account was created by an authorized process (e.g., a software installation), document this and decide whether to keep it (as a documented service account) or remove it.
8. If the account cannot be explained, treat this as a security incident per D10-003.

**7.3 If physical access has been compromised:**
1. Assess the scope: what equipment was accessible, for how long, and by whom (if known).
2. If the compromise was brief and the intruder is known and trusted: document the incident, review physical security measures, and improve them.
3. If the compromise was extended or the intruder is unknown: assume all accessible systems may have been tampered with. Verify system integrity using the procedures in D6-003. Check for unauthorized boot media, hardware modifications, or new devices.
4. Change all passwords on systems that were physically accessible.
5. If tampering is suspected but cannot be confirmed, consider rebuilding affected systems from known-good backups.
6. Document everything in the operational log and the decision log.

**7.4 If the succession account has failed during a succession event:**
1. Use the administrative account to create a new account for the successor.
2. Configure the new account with operator-level permissions.
3. Set a temporary password that the successor changes immediately on first login.
4. Document the failure and the remediation in the decision log.
5. Update the succession packet with the new account details.
6. Investigate why the succession account failed and correct the root cause.

## 8. Evolution Path

- **Years 0-5:** The access control architecture is established. The operator is learning to live with the inconvenience of privilege separation. This is the period where habits form. Resist the temptation to loosen controls because you are the only user. You are building the architecture for decades, not optimizing for today's convenience.

- **Years 5-15:** The quarterly audits should have settled into routine. The account inventory should be stable. The main challenge is preventing drift -- the slow loosening of permissions that happens when the same person has been the only user for years and the controls feel unnecessary.

- **Years 15-30:** Succession becomes the dominant access control concern. The succession account must be tested and current. The physical access model must account for the successor's needs. If a succession event occurs, the access control procedures must be comprehensible to the new operator.

- **Years 30-50+:** Multiple generations of operators may have used the system. The account history in the logs becomes a record of institutional continuity. The access control model should be reviewed for whether it still reflects the institution's needs and the available technology.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
I want to be honest about the tension in this article. I am writing access control procedures for a system where I am the only user. Every fiber of my practical self says: just use one account, use sudo when needed, and stop overthinking it. And if this were a personal project with a five-year horizon, that advice would be correct.

But this is not a five-year project. It is a fifty-year institution. The access controls I establish now will be the foundation for every successor's experience of the system. If I build the habit of discipline now -- separate accounts, least privilege, quarterly audits -- the institution will be ready for succession without a frantic retrofit. If I take shortcuts now, every successor will inherit those shortcuts and the vulnerabilities they create.

The physical access section may seem excessive for a home environment. It is not. The single most likely way for this institution to be compromised is not a sophisticated attack. It is a visitor who sits down at an unlocked terminal, a child who presses buttons on an unsecured server, or a burglar who steals hardware that contains unencrypted credentials. Physical security is not paranoia. It is the recognition that the air gap only works if the physical perimeter is intact.

## 10. References

- ETH-001 -- Ethical Foundations of the Institution (Principle 1: Sovereignty; Principle 2: Integrity Over Convenience)
- CON-001 -- The Founding Mandate (self-sovereign operation, documentation completeness)
- GOV-001 -- Authority Model (succession protocol, decision tiers for access changes)
- SEC-001 -- Threat Model and Security Philosophy (trust model, defense in depth, physical threats)
- OPS-001 -- Operations Philosophy (operational tempo for access audits, documentation-first principle)
- SEC-003 -- Cryptographic Key Management (password and key storage)
- D6-001 -- Data Philosophy (data tier classification, Tier 1 preservation for logs)
- D6-002 -- Backup Doctrine (backup of access control configurations)
- D10-002 -- Daily Operations Doctrine (daily physical security checks)
- D10-003 -- Incident Response Procedures (handling access-related incidents)

---

---

# SEC-003 -- Cryptographic Key Management

**Document ID:** SEC-003
**Domain:** 3 -- Security & Integrity
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, SEC-001, SEC-002, OPS-001
**Depended Upon By:** D6-002, D6-003, D10-003. All articles involving encryption, authentication, or digital signatures.

---

## 1. Purpose

This article defines the complete cryptographic key management procedures for the holm.chat Documentation Institution. It covers the lifecycle of every key the institution uses: how keys are generated, where they are stored, when they are rotated, how they are revoked, and what happens when a key is lost or compromised.

Cryptographic keys are the institution's most sensitive assets. They protect the confidentiality of encrypted data, the integrity of signed documents, and the authenticity of the operator's identity. A lost key can render data permanently inaccessible. A compromised key can undermine every security guarantee the institution provides. A poorly managed key is a time bomb -- it functions perfectly until the moment it fails, and then the damage is often irreversible.

This article is written for a single operator managing keys without a dedicated security team, without a hardware security module (HSM) cluster, and without the luxury of a corporate PKI. The procedures are practical, achievable, and designed for the constraints of an air-gapped, off-grid institution. They are also designed for succession -- every key management decision must be comprehensible and reproducible by a future operator who has only this document and the institution's stored artifacts to guide them.

## 2. Scope

This article covers:

- The institution's key inventory: what categories of keys exist and what each is for.
- Key generation procedures: how to generate keys safely on an air-gapped system.
- Key storage architecture: where keys are kept, in what form, with what protections.
- Key rotation schedules: when keys should be replaced and how.
- Key revocation procedures: how to declare a key invalid and propagate that declaration.
- Key recovery procedures: what to do when a key is lost.
- Compromised key response: what to do when a key is suspected or confirmed compromised.
- The key ceremony: the formal process for generating, distributing, and recording high-value keys.
- Succession provisions: how key material is made available to a legitimate successor.

This article does not cover:

- Which encryption algorithms to use (see subsequent Domain 3 implementation articles; this article is algorithm-agnostic by design, because algorithms change over decades).
- Full-disk encryption setup procedures (see Domain 3 implementation articles).
- Network authentication (not applicable; the institution is air-gapped).

## 3. Background

### 3.1 Why Key Management Is the Hardest Part of Cryptography

The mathematics of modern cryptography are sound. The algorithms, when properly implemented, are essentially unbreakable within practical time horizons. The weakness is never the math. The weakness is always the management -- how keys are generated (insufficient randomness), how they are stored (written on a sticky note), how they are transported (sent in plaintext), how they are retired (never), and how they are recovered (they cannot be).

For this institution, the key management challenge is compounded by two factors. First, the air gap means there is no online key management service. Every key management action is manual, performed by the operator, with physical media. Second, the fifty-year horizon means keys must survive hardware generations, software changes, format migrations, and operator succession. A key stored on a USB drive in 2026 must still be accessible in 2056 -- or the data it protects must be re-encrypted before the old key becomes inaccessible.

### 3.2 The Single-Operator Key Management Problem

In a corporate environment, key management is distributed across multiple people with multiple roles: key officers, key custodians, recovery agents. This distribution provides checks, balances, and redundancy. A single operator has none of this. They generate their own keys, store their own keys, use their own keys, and must somehow provide for the recovery of their own keys if they become incapacitated.

This article addresses the single-operator problem through three mechanisms: rigorous procedure (compensating for the lack of peer review with documented, verifiable steps), redundant storage (compensating for single-point-of-failure risk with multiple independent copies), and succession provisioning (ensuring that keys outlive the operator who created them).

### 3.3 Key Categories in This Institution

The institution uses keys for four purposes, each with different management requirements:

- **Disk encryption keys:** Protect data at rest on the institution's storage devices. Loss means permanent data inaccessibility. Compromise means an attacker with physical access can read the data.
- **Authentication keys:** Prove the identity of the operator to the system (SSH keys, login credentials). Loss means temporary lockout. Compromise means unauthorized access.
- **Signing keys:** Prove the integrity and authorship of documents and backups. Loss means inability to sign new material. Compromise means an attacker can forge the operator's identity.
- **Backup encryption keys:** Protect backup media that may be stored at geographically remote locations. Loss means backup data is inaccessible. Compromise means backup data is exposed.

## 4. System Model

### 4.1 The Key Inventory

The institution maintains a key inventory -- a document that records every active key, its purpose, its creation date, its rotation schedule, and the location of all copies. The key inventory is itself a Tier 1 document per D6-001, but it does not contain the keys themselves. It contains metadata about the keys: identifiers, fingerprints, purposes, and locations.

The key inventory is stored in the documentation corpus and reviewed during every quarterly access audit per SEC-002.

Each entry in the key inventory contains:
- **Key ID:** A unique identifier for the key.
- **Key Type:** Disk encryption, authentication, signing, or backup encryption.
- **Algorithm and Key Length:** What algorithm and key size the key uses.
- **Creation Date:** When the key was generated.
- **Rotation Date:** When the key is next due for rotation.
- **Expiration Date:** When the key must be retired regardless of rotation status.
- **Primary Storage Location:** Where the active copy of the key resides.
- **Backup Storage Locations:** Where backup copies are stored.
- **Fingerprint/Hash:** A hash of the key for verification without exposing the key itself.
- **Associated Data:** What data or systems the key protects.
- **Succession Provisions:** Whether and how this key is included in the succession packet.

### 4.2 Key Generation Procedures

All key generation must occur on the institution's air-gapped systems. Keys are never generated on connected systems and imported.

**Step-by-step procedure for generating a new key:**

1. Log in to the administrative account per SEC-002.
2. Open a terminal session. Verify you are on the air-gapped system (no network interfaces active): `ip link show` -- all interfaces should be DOWN or absent.
3. Verify the system's random number generator is functioning: `cat /proc/sys/kernel/random/entropy_avail` -- the value should be at least 256. If it is lower, generate entropy by performing disk I/O, moving the mouse, or using a hardware random number generator if available.
4. Generate the key using the appropriate tool for the key type:
   - For disk encryption (LUKS): `cryptsetup luksFormat [device]` -- the passphrase is the key. Alternatively, use a keyfile: `dd if=/dev/urandom of=/path/to/keyfile bs=4096 count=1`
   - For SSH authentication keys: `ssh-keygen -t ed25519 -C "[key_purpose]@holm.chat [date]"`
   - For GPG signing/encryption keys: `gpg --full-generate-key` -- select appropriate algorithm and key size.
   - For backup encryption: `gpg --full-generate-key` with a dedicated key for backup purposes, or generate a symmetric key: `dd if=/dev/urandom of=/path/to/backup-key bs=32 count=1`
5. Immediately record the key in the key inventory with all required fields.
6. Create backup copies per Section 4.3.
7. Verify the key works by performing a test operation (encrypt and decrypt a test file, sign and verify a test document, authenticate a test session).
8. Record the key generation in the operational log, including the key ID and fingerprint (not the key itself).

**Critical: never record the actual key material in the operational log, the decision log, or the documentation corpus. Only record identifiers and fingerprints.**

### 4.3 Key Storage Architecture

Keys are stored in a three-tier architecture, each tier providing a different balance of accessibility and security:

**Tier A: Active Storage.**
The key as used by the running system. For disk encryption, this is the passphrase or keyfile that unlocks the volume at boot. For SSH keys, this is the private key in the user's `.ssh` directory. For GPG keys, this is the key in the user's GPG keyring.

Security: Protected by the system's access controls per SEC-002. The disk encryption key protects the other keys at rest (the SSH and GPG keys are on an encrypted volume).

**Tier B: Local Backup.**
A copy of the key stored on separate physical media in the equipment zone (Zone 1 per SEC-002). This media is encrypted with a separate passphrase known only to the operator. For physical media: a USB drive stored in a locked container within Zone 1.

Procedure for creating Tier B backups:
1. Prepare a USB drive. Format it and create a LUKS-encrypted partition with a strong, unique passphrase.
2. Mount the encrypted partition.
3. Copy all key material to the encrypted drive: keyfiles, exported GPG keys (`gpg --export-secret-keys --armor > keys.asc`), SSH private keys.
4. Include a plain-text file named `INVENTORY.txt` that lists what keys are on this drive, their key IDs, and what they are for.
5. Unmount and close the encrypted partition.
6. Label the physical drive clearly: "KEY BACKUP -- TIER B -- [date]"
7. Store in the locked container in Zone 1.

**Tier C: Off-site Backup.**
A copy of the key stored at a geographically separate location (Zone 3 per SEC-002). This copy is encrypted and stored with the succession packet or in a safe deposit box or other secure off-site location.

Procedure for creating Tier C backups:
1. Prepare a USB drive. Encrypt it as with Tier B, but use a different passphrase from Tier B.
2. Copy all key material as with Tier B.
3. Include `INVENTORY.txt` as with Tier B.
4. Also include a printed sheet in a sealed envelope alongside the USB drive that contains: the drive's encryption passphrase, the list of key IDs and their purposes, and basic instructions for accessing the keys. This printed sheet is the succession document for key recovery.
5. Transport the drive and envelope to the off-site location.
6. Record the off-site location in the key inventory.

### 4.4 Key Rotation Schedule

Key rotation limits the exposure window if a key is compromised without the operator's knowledge. Different key types have different rotation schedules:

- **Disk encryption keys:** Rotate every 5 years or when the encryption algorithm is superseded. Rotation requires re-encrypting the volume with a new key, which is a major operation. Schedule it during the annual operations cycle.
- **Authentication keys (SSH):** Rotate every 2 years. Generate a new key, deploy it, verify access, then remove the old key.
- **Signing keys (GPG):** Rotate every 3 years. Generate a new key, publish the new fingerprint in the documentation, sign new material with the new key, retain the old key (in revoked state) for verification of previously signed material.
- **Backup encryption keys:** Rotate every 2 years. After rotation, do not destroy the old key immediately -- existing backups encrypted with the old key must remain decryptable. Retain old backup keys in the key inventory with status "RETIRED -- RETAINED FOR DECRYPTION" until all backups encrypted with that key have been either re-encrypted or aged out of the retention schedule.

**Step-by-step rotation procedure:**
1. Generate the new key per Section 4.2.
2. Create Tier B and Tier C backups of the new key per Section 4.3.
3. Deploy the new key: install it in the active location, configure the relevant system to use it.
4. Test the new key: verify it works for its intended purpose.
5. If the old key protected existing data (disk encryption, backup encryption): plan and execute the data re-encryption or establish a documented retention of the old key.
6. Update the key inventory: mark the old key as ROTATED with the date and the new key's ID. Add the new key with all required fields.
7. Update the Tier B and Tier C backups to include the new key.
8. Record the rotation in the operational log and the decision log.

### 4.5 Key Revocation

A key is revoked when it is known or suspected to be compromised, or when the person it authenticates is no longer authorized.

**Immediate revocation procedure:**
1. Generate a replacement key per Section 4.2 (do this first, so you have a working key before revoking the old one).
2. Revoke the old key:
   - For SSH keys: remove the public key from all `authorized_keys` files. Remove the private key from the user's `.ssh` directory.
   - For GPG keys: generate a revocation certificate (`gpg --gen-revoke [key_id]`), import it to the keyring, and export the revoked public key for distribution.
   - For disk encryption keys: change the LUKS passphrase (`cryptsetup luksChangeKey [device]`) or remove the compromised keyslot.
   - For backup encryption keys: this is complex. Existing backups remain encrypted with the old key. They cannot be retroactively re-encrypted. Flag all backups encrypted with the revoked key as "ENCRYPTED WITH REVOKED KEY" in the backup inventory. Prioritize re-creating critical backups with the new key.
3. Deploy the replacement key per Section 4.4 step 3.
4. Update the key inventory: mark the old key as REVOKED with the date and reason.
5. Update all backup copies (Tier B and Tier C) to remove the revoked key and include the replacement.
6. Document the revocation in the operational log, the decision log, and the key inventory.

### 4.6 The Key Ceremony

A key ceremony is a formal procedure for generating high-value keys -- specifically, the master disk encryption key and the primary signing key. The ceremony provides a structured, documented process that ensures the key is generated securely and that all backup and succession provisions are established from the beginning.

**The key ceremony procedure:**

1. **Preparation (the day before):**
   - Clear the schedule. No other operational work on ceremony day.
   - Prepare all media: USB drives for Tier B and Tier C backups, formatted and ready for encryption.
   - Prepare printed forms: key inventory entries, succession documents.
   - Verify the system's random number generator.

2. **Generation (ceremony day, step 1):**
   - Log in to the administrative account.
   - Verify air-gap integrity: no network interfaces active.
   - Generate the key per Section 4.2.
   - Record the key fingerprint on the printed form.

3. **Backup (ceremony day, step 2):**
   - Create Tier B backup per Section 4.3.
   - Create Tier C backup per Section 4.3.
   - Verify both backups by mounting and confirming the key material is present and readable.

4. **Verification (ceremony day, step 3):**
   - Use the key for its intended purpose on a test target.
   - Use the Tier B backup copy to perform the same test (proving the backup is functional).

5. **Documentation (ceremony day, step 4):**
   - Complete the key inventory entry.
   - Record the ceremony in the operational log with: date, key ID, fingerprint, key type, all storage locations, and a statement that verification succeeded.
   - Record in the decision log as a Tier 2 decision.

6. **Succession (ceremony day, step 5):**
   - Update the succession packet with the new key's existence and the Tier C backup location.
   - Seal the Tier C backup and succession documents in the envelope.

7. **Transport (within 7 days):**
   - Deliver the Tier C backup to the off-site location.
   - Record the delivery date in the key inventory.

## 5. Rules & Constraints

- **R-SEC3-01:** All key generation must occur on the institution's air-gapped systems. No exceptions.
- **R-SEC3-02:** Every key must exist in at least three copies: active, Tier B backup, and Tier C backup. No key may exist in only one copy.
- **R-SEC3-03:** The key inventory must be maintained as a Tier 1 document per D6-001. It must be accurate at all times. Discrepancies between the inventory and reality are security incidents.
- **R-SEC3-04:** Key material must never appear in the operational log, the decision log, the documentation corpus, or any other document. Only key identifiers and fingerprints may be recorded in documents.
- **R-SEC3-05:** Key rotation schedules must be followed. Overdue rotations must be flagged during the quarterly access audit per SEC-002 and remediated within 30 days.
- **R-SEC3-06:** Old keys protecting existing data must be retained in RETIRED status until all data they protect has been re-encrypted or aged out of the retention schedule. Retired keys are stored in the key archive (a designated section of the Tier B and Tier C backups).
- **R-SEC3-07:** The succession packet must contain sufficient information for a successor to access all Tier C key backups and all encrypted institutional data. This packet must be tested annually.
- **R-SEC3-08:** Compromised keys must be revoked within 24 hours of the compromise being detected. The revocation procedure in Section 4.5 is mandatory.
- **R-SEC3-09:** The key ceremony procedure (Section 4.6) must be used for all master disk encryption keys and primary signing keys. Lesser keys may use the simplified generation procedure in Section 4.2.

## 6. Failure Modes

- **Key loss.** The active key and all backups are lost or destroyed. Encrypted data protected by the lost key is permanently inaccessible. This is the most catastrophic key management failure. Mitigation: the three-copy rule (R-SEC3-02), geographic separation of Tier C backups, and regular backup verification.

- **Key compromise without detection.** An attacker obtains a copy of a key without the operator's knowledge. The attacker can access encrypted data or forge signatures indefinitely. Mitigation: physical security per SEC-002, key rotation (which limits the window of undetected compromise), and integrity verification per D6-003 (which may detect unauthorized access).

- **Rotation procrastination.** Key rotation is always scheduled for next month. Years pass. Keys become ancient, algorithms become deprecated, and the institution's cryptographic foundation grows brittle. Mitigation: the quarterly audit per SEC-002 flags overdue rotations. The rotation schedule is in the key inventory, which is reviewed quarterly.

- **Succession packet decay.** The succession packet is prepared once and never updated. Keys are rotated but the succession packet still contains the old keys. When the succession event occurs, the successor cannot decrypt the data. Mitigation: every key rotation must include updating the succession packet. R-SEC3-07 requires annual testing.

- **Passphrase amnesia.** The operator forgets the passphrase for the Tier B or Tier C backup media. The backup keys are inaccessible. Mitigation: passphrases for backup media are recorded in the sealed succession envelope at the off-site location. The operator should also consider a password manager database (itself encrypted) on the air-gapped system for day-to-day passphrase management.

- **Algorithm obsolescence.** The cryptographic algorithm used to generate keys becomes deprecated or broken. Keys generated with the old algorithm are no longer trustworthy. Mitigation: the rotation schedule forces periodic key regeneration, which naturally migrates to current algorithms. The annual security review per OPS-001 should assess algorithm health.

## 7. Recovery Procedures

**7.1 If an active key is lost but backups exist:**
1. Retrieve the Tier B backup from Zone 1.
2. Decrypt the backup media using the Tier B passphrase.
3. Copy the lost key to the appropriate active location.
4. Test the recovered key.
5. Generate a new Tier B backup to replace the one you just opened (the passphrase is now "used" and should be considered less secure).
6. Document the recovery in the operational log.
7. Investigate how the active key was lost and implement preventive measures.

**7.2 If active and Tier B keys are lost but Tier C exists:**
1. Retrieve the Tier C backup from the off-site location.
2. Open the sealed envelope to obtain the decryption passphrase.
3. Decrypt the backup media.
4. Copy the keys to the appropriate active locations.
5. Test all recovered keys.
6. Immediately create new Tier B and Tier C backups with new passphrases.
7. Investigate the loss. This is a serious incident per D10-003.
8. Document everything in the operational log and the decision log.

**7.3 If all copies of a key are lost:**
1. Accept that data encrypted with the lost key is permanently inaccessible. This cannot be reversed.
2. Document the loss: what key was lost, what data it protected, and the scope of the impact.
3. Assess the damage: list all data that is now inaccessible.
4. Generate new keys for all future use.
5. Create new backups at all tiers.
6. Review and revise key management procedures to prevent recurrence.
7. Record this as a catastrophic incident in the decision log (Tier 2 decision to acknowledge and respond to the loss).

**7.4 If a key is suspected compromised:**
1. Immediately generate a replacement key per Section 4.2.
2. Revoke the compromised key per Section 4.5.
3. Assess what data or systems the compromised key protected.
4. For disk encryption: assume data may have been copied. Assess the physical access history. If physical access was not compromised, the encrypted data at rest is still protected (the attacker needs the key AND physical access to the media).
5. For signing keys: assess what was signed with the compromised key. Notify (through documentation) that signatures from the compromised key should not be trusted after the suspected compromise date.
6. For backup encryption: assess which backups are encrypted with the compromised key. Prioritize re-creating critical backups with the new key.
7. Document the incident fully in the operational log and follow D10-003 incident response procedures.

## 8. Evolution Path

- **Years 0-5:** The key management infrastructure is being established. Perform the key ceremony for the initial master keys. Establish the rotation schedule. Build the habit of quarterly key inventory review. This is the foundation period -- get it right now and the next decades are manageable.

- **Years 5-15:** The first key rotations occur. This is the first real test of the rotation procedures. The Tier C backups should have been verified multiple times. The succession packet should have been tested. Cryptographic algorithm assessments during annual reviews may begin to flag the first algorithm changes.

- **Years 15-30:** Hardware generations will have changed. The USB drives used for Tier B and Tier C backups may need to be migrated to new media types. The algorithms used at founding may be approaching end of life. Plan and execute cryptographic migrations proactively, before the old algorithms are considered broken.

- **Years 30-50+:** Multiple key generations will have occurred. The key inventory's history section becomes a record of the institution's cryptographic evolution. A successor inheriting the institution should be able to trace the full key lineage from founding to present.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
Key management is the one area where I am most tempted to cut corners and most afraid of the consequences. The three-copy rule feels excessive for a single person. The key ceremony feels theatrical. The rotation schedules feel premature when the keys are brand new.

But I have seen what happens when key management is casual. I have seen organizations locked out of their own data because the one person who knew the passphrase left and nobody thought to ask. I have seen backup tapes sitting in a vault for a decade, perfectly preserved, perfectly encrypted, perfectly inaccessible because the decryption key was on a server that was decommissioned five years ago.

The key ceremony is not theater. It is a forcing function. It forces me to create backups before I need them, to document before I forget, and to plan for succession before succession is urgent. The day I generate a master key is the day I am most focused on that key's lifecycle. The ceremony captures that focus and turns it into procedure.

One note for future operators: the specific algorithms and tools referenced in this article will almost certainly be outdated by the time you read it. `ed25519` may have been superseded. `gpg` may have been replaced. The procedures -- generate securely, store redundantly, rotate regularly, revoke promptly, document everything -- will not be outdated. Adapt the tools. Keep the procedures.

## 10. References

- ETH-001 -- Ethical Foundations of the Institution (Principle 2: Integrity Over Convenience)
- CON-001 -- The Founding Mandate (data sovereignty, self-contained operation)
- GOV-001 -- Authority Model (succession protocol, decision tiers)
- SEC-001 -- Threat Model and Security Philosophy (encryption mandate R-SEC-02, key storage R-SEC-04)
- SEC-002 -- Access Control Procedures (account architecture, physical zones, quarterly audits)
- OPS-001 -- Operations Philosophy (operational tempo for key rotation, annual security review)
- D6-001 -- Data Philosophy (key inventory as Tier 1 data, data tier system)
- D6-002 -- Backup Doctrine (backup encryption key management)
- D6-003 -- Data Integrity Verification (detecting unauthorized access via integrity checks)
- D10-003 -- Incident Response Procedures (compromised key response)

---

---

# D6-002 -- Backup Doctrine

**Document ID:** D6-002
**Domain:** 6 -- Data & Archives
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, SEC-001, SEC-002, SEC-003, OPS-001, D6-001
**Depended Upon By:** D6-003, D10-003. All articles involving data recovery, disaster recovery, or data durability.

---

## 1. Purpose

This article defines the complete backup strategy for the holm.chat Documentation Institution. It specifies what gets backed up, how often, where backups are stored, in what format, how backups are verified, how restorations are tested, and what to do when a backup fails.

A backup that has not been tested is not a backup. It is a hope. This article is written to eliminate hope from the backup process and replace it with verified, documented, reproducible procedures. Every backup created under these procedures will have been verified at creation and tested through periodic restoration. Every backup failure will be detected, documented, and remediated.

This article implements the data redundancy requirements of D6-001 (particularly R-D6-02: Tier 1 data in at least three copies on two media types with geographic separation) and the backup verification requirements of OPS-001 (monthly test restorations). It translates those requirements into specific, step-by-step procedures that a single operator can execute.

## 2. Scope

This article covers:

- The backup strategy: what is backed up and at what frequency.
- The 3-2-1 rule as adapted for air-gapped, off-grid operation.
- Backup media selection and management.
- Step-by-step backup procedures for each data tier.
- Backup verification: how to confirm a backup is complete and correct.
- Test restoration procedures: how to verify a backup can actually be restored.
- Backup rotation and retention: how long backups are kept and when they are recycled.
- Backup encryption: how backup media is protected.
- The backup log: what is recorded and why.

This article does not cover:

- Data integrity verification beyond backup verification (see D6-003).
- Full disaster recovery procedures (see subsequent Domain 12 articles).
- Storage hardware selection and maintenance (see subsequent Domain 6 articles).
- Key management for backup encryption keys (see SEC-003).

## 3. Background

### 3.1 The 3-2-1 Rule

The 3-2-1 rule is the industry standard for backup resilience: maintain at least 3 copies of data, on at least 2 different types of media, with at least 1 copy off-site. This rule was developed for networked environments where "off-site" might mean a cloud backup or a remote data center. For an air-gapped institution, the rule must be adapted.

**The adapted 3-2-1 rule for this institution:**
- **3 copies:** The primary data on the institution's active storage, a local backup on separate physical media, and an off-site backup at a geographically separate location.
- **2 media types:** The active storage and at least one backup should be on different types of media (e.g., if active storage is on SSDs, at least one backup should be on spinning disks, or on tape, or on optical media). This protects against media-class failures -- a manufacturing defect that affects a batch of SSDs will not affect your spinning disk backup.
- **1 off-site:** At least one backup copy must be stored at a location far enough from the primary site that a single disaster (fire, flood, earthquake) cannot destroy both. "Far enough" depends on the threat model. For most scenarios, a separate building at least 30 kilometers away is sufficient. For maximum resilience, a separate geographic region.

### 3.2 The Air-Gap Backup Challenge

In a networked environment, backups can be automated end-to-end: a backup agent runs, compresses and encrypts the data, transmits it to a remote server, and verifies the transfer. In an air-gapped environment, the off-site copy requires physical transport. Someone must physically carry a backup medium from the primary site to the off-site location. This introduces delay, effort, and the risk of loss or damage during transport.

This article addresses the air-gap challenge by defining clear procedures for backup creation, encryption, physical labeling, transport, and verification at both origin and destination. The process is manual by necessity but structured to minimize the chance of error.

### 3.3 What Backups Protect Against

Backups protect against: hardware failure (disk death), software corruption (filesystem errors, application bugs), operator error (accidental deletion, misconfiguration), and site-level disasters (fire, flood, theft). They do not protect against: threats that the operator is unaware of (if all copies are corrupted silently, the backups are also corrupted), or threats that affect all copies simultaneously (if the backup encryption key is lost, all encrypted backups are inaccessible).

Understanding what backups do and do not protect against is essential. Backups are one layer of defense. They complement, but do not replace, data integrity verification (D6-003), access controls (SEC-002), and key management (SEC-003).

## 4. System Model

### 4.1 What Gets Backed Up

All institutional data is backed up according to its tier classification per D6-001:

**Tier 1: Institutional Memory (backed up daily).**
- The documentation corpus (all articles, all versions).
- The decision log and operational log.
- The key inventory (not the keys themselves -- those are managed per SEC-003).
- Commentary sections and all governance records.
- Backup frequency: daily incremental, weekly full.

**Tier 2: Operational Data (backed up weekly).**
- System configurations and service definitions.
- Active project data and working documents.
- Automation scripts and service account configurations per D8-001.
- Backup frequency: weekly incremental, monthly full.

**Tier 3: Reference Data (backed up monthly).**
- Research materials, imported knowledge bases, archived media.
- Backup frequency: monthly full.

**Tier 4: Transient Data (not backed up).**
- Temporary files, scratch data, working copies.
- Per D6-001, Tier 4 data has no preservation requirement.

### 4.2 The Backup Architecture

The institution maintains three backup locations, implementing the adapted 3-2-1 rule:

**Location A: On-site local backup.**
- Physical media: a dedicated storage device (external drive, NAS, or secondary internal drive) separate from the active storage.
- Media type: should differ from active storage if possible (if active is SSD, backup to HDD or vice versa).
- Contents: all Tier 1, Tier 2, and Tier 3 backups.
- Encryption: the backup device is encrypted with the backup encryption key per SEC-003.
- Stored in Zone 1 per SEC-002.

**Location B: On-site archival backup.**
- Physical media: a set of removable media (external HDDs, optical discs, or tape) stored in Zone 1.
- This is the monthly archival snapshot -- a complete image of all backed-up data.
- Encryption: each medium is encrypted.
- These are rotated on a schedule defined in Section 4.5.

**Location C: Off-site backup.**
- Physical media: removable media transported to Zone 3 per SEC-002.
- Contents: all Tier 1 data (always), Tier 2 data (at least quarterly), Tier 3 data (at least annually).
- Encryption: each medium is encrypted with the backup encryption key. The decryption key is available at the off-site location per SEC-003 Tier C provisions.
- Transport schedule: monthly for Tier 1 updates, quarterly for comprehensive off-site backup.

### 4.3 Step-by-Step Backup Procedures

**Daily Incremental Backup (Tier 1 data to Location A):**

1. Log in to the administrative account.
2. Verify the backup destination device is connected and mounted.
3. Run the backup tool. The recommended approach is `rsync` with checksums:
   ```
   rsync -avz --checksum --delete \
     --backup --backup-dir=/backup/incremental/$(date +%Y-%m-%d) \
     /path/to/tier1-data/ \
     /path/to/backup-location-a/tier1/current/
   ```
   This creates a current mirror of Tier 1 data with changed files preserved in dated incremental directories.
4. Verify the backup completed without errors. Check the rsync exit code: `echo $?` should be 0.
5. Verify the backup size: `du -sh /path/to/backup-location-a/tier1/current/` -- compare against the source size. They should be within 1% of each other.
6. Generate a checksum manifest of the backup:
   ```
   find /path/to/backup-location-a/tier1/current/ -type f -exec sha256sum {} \; > \
     /path/to/backup-location-a/tier1/manifests/$(date +%Y-%m-%d).sha256
   ```
7. Record the backup in the backup log: date, time, type (daily incremental), source, destination, size, exit code, any errors.
8. Unmount the backup device if it is removable.

**Weekly Full Backup (Tier 1 data to Location A):**

1. Follow steps 1-2 of the daily procedure.
2. Create a full archive:
   ```
   tar -czf /path/to/backup-location-a/tier1/weekly/$(date +%Y-%m-%d)-full.tar.gz \
     /path/to/tier1-data/
   ```
3. Generate a checksum of the archive:
   ```
   sha256sum /path/to/backup-location-a/tier1/weekly/$(date +%Y-%m-%d)-full.tar.gz > \
     /path/to/backup-location-a/tier1/weekly/$(date +%Y-%m-%d)-full.tar.gz.sha256
   ```
4. Verify the archive is valid:
   ```
   tar -tzf /path/to/backup-location-a/tier1/weekly/$(date +%Y-%m-%d)-full.tar.gz > /dev/null
   ```
   This lists the archive contents without extracting, confirming the archive is readable. Exit code 0 means success.
5. Record in the backup log.

**Monthly Full Backup (All tiers to Location B):**

1. Prepare the Location B medium: connect the external drive or prepare optical media.
2. If the medium is new, format and encrypt it per SEC-003.
3. Create full archives of each tier:
   ```
   tar -czf /path/to/location-b/tier1-$(date +%Y-%m-%d).tar.gz /path/to/tier1-data/
   tar -czf /path/to/location-b/tier2-$(date +%Y-%m-%d).tar.gz /path/to/tier2-data/
   tar -czf /path/to/location-b/tier3-$(date +%Y-%m-%d).tar.gz /path/to/tier3-data/
   ```
4. Generate checksums for each archive.
5. Create a manifest file listing all archives, their sizes, and their checksums.
6. Verify each archive is valid (tar -tzf test).
7. Label the physical medium: "MONTHLY BACKUP -- [date] -- TIERS 1/2/3"
8. Store in Zone 1 per SEC-002.
9. Record in the backup log.

**Off-site Backup (to Location C):**

1. Prepare the removable medium. Format and encrypt per SEC-003.
2. Copy the most recent Tier 1 full backup and all incremental backups since the last off-site update.
3. On quarterly intervals, include Tier 2 data. On annual intervals, include Tier 3 data.
4. Generate checksums and a manifest.
5. Label the medium: "OFF-SITE BACKUP -- [date] -- CONTENTS: [tier list]"
6. Transport to Zone 3. Handle with care. Do not leave unattended during transport.
7. At the off-site location, verify the medium is readable: mount, check the manifest, spot-check one or two files.
8. Record the transport in the backup log: date, contents, destination, verification result.

### 4.4 Backup Verification Procedures

**Verification at creation (every backup):**
- Check tool exit codes (0 = success).
- Compare source and destination sizes.
- Generate and store checksum manifests.
- Verify archive readability (tar -tzf or equivalent).

**Monthly test restoration:**
Per OPS-001, a randomly selected backup set must be test-restored monthly. The procedure:

1. Select a backup to test. Use a random selection method: roll a die to choose the tier, then select a date.
2. Prepare a test restoration target: a scratch partition or temporary directory that will be erased after the test.
3. Copy the backup to the test target.
4. Extract the backup:
   ```
   tar -xzf /path/to/test-target/backup-archive.tar.gz -C /path/to/test-target/restored/
   ```
5. Verify the restored data:
   - Compare file counts: `find /path/to/test-target/restored/ -type f | wc -l` vs. expected.
   - Spot-check individual files: open several, verify content is correct and readable.
   - If a checksum manifest exists for the backup, verify checksums:
     ```
     cd /path/to/test-target/restored/ && sha256sum -c /path/to/manifest.sha256
     ```
6. Record the test results in the backup log: date, backup tested (date and tier), restoration success (yes/no), any issues found, time taken.
7. Clean up the test target: securely erase the restored data.

If the test restoration fails:
- Identify the failure point: was the archive corrupt, was the checksum wrong, or was the restoration procedure incorrect?
- Attempt restoration from a different backup of the same data.
- If multiple backups fail, escalate per D10-003 (incident response).
- Do not wait until the next monthly test. Fix the backup process immediately.

### 4.5 Backup Rotation and Retention

Backups are retained according to this schedule:

- **Daily incrementals:** Retained for 30 days. After 30 days, deleted to reclaim space.
- **Weekly full backups:** Retained for 12 weeks (approximately 3 months).
- **Monthly full backups (Location B):** Retained for 12 months. After 12 months, the medium is recycled (securely erased and reformatted for reuse).
- **Off-site backups (Location C):** The most recent 4 off-site backups are retained (approximately 4 months of monthly transports). Older media is retrieved, securely erased, and recycled.
- **Annual archival snapshot:** One full backup from each year is retained permanently as a historical archive. This is a Tier 1 data artifact. It is stored at both Location B and Location C.

**Secure erasure procedure for recycled backup media:**
1. Overwrite the entire medium: `dd if=/dev/urandom of=/dev/[device] bs=4M status=progress`
2. Verify the overwrite by attempting to mount the medium (it should fail or show no recognizable data).
3. Reformat and re-encrypt the medium.
4. Record the recycling in the backup log: date, medium identifier, confirmation of erasure.

### 4.6 The Backup Log

The backup log is a Tier 1 data artifact per D6-001. It records every backup operation, every verification, every test restoration, and every failure. It is the audit trail that proves the backup process is functioning.

Each entry contains:
- Date and time.
- Operation type: daily backup, weekly backup, monthly backup, off-site transport, test restoration, rotation/recycling.
- Source and destination.
- Data tiers included.
- Size of backup.
- Exit codes and verification results.
- Any errors or anomalies.
- Operator identifier.

The backup log is itself backed up as part of the Tier 1 data.

## 5. Rules & Constraints

- **R-D6B-01:** Tier 1 data must be backed up daily. No exceptions. If the daily backup is missed, it must be performed the next day and the gap documented in the backup log.
- **R-D6B-02:** All backups must be verified at creation per Section 4.4. An unverified backup is not counted toward the 3-2-1 requirement.
- **R-D6B-03:** Test restorations must be performed monthly per Section 4.4. Test results must be recorded in the backup log. A failed test restoration triggers immediate investigation.
- **R-D6B-04:** All backup media must be encrypted per SEC-003. Unencrypted backup media is a security violation.
- **R-D6B-05:** Off-site backups must be transported at least monthly for Tier 1 data. The transport must be documented in the backup log.
- **R-D6B-06:** The backup log is Tier 1 data and is subject to all Tier 1 preservation requirements per D6-001.
- **R-D6B-07:** Backup media that reaches the end of its retention period must be securely erased before recycling per Section 4.5.
- **R-D6B-08:** The annual archival snapshot must be created and stored at both Location B and Location C. It is retained permanently.
- **R-D6B-09:** Backup procedures must be documented in sufficient detail that a successor can execute them from the documentation alone, per D9-001's self-teaching requirement.

## 6. Failure Modes

- **Backup neglect.** Backups are not performed on schedule. The gap between the last backup and current data grows. When a failure occurs, the data loss is measured in weeks or months instead of hours. Mitigation: R-D6B-01 makes daily Tier 1 backup mandatory. The daily checklist per D10-002 includes a backup verification item.

- **Silent backup failure.** The backup tool runs but produces corrupt or incomplete output. The operator sees "backup complete" and moves on. The backup is useless. Mitigation: the verification procedures in Section 4.4, which check exit codes, compare sizes, generate checksums, and verify archive readability.

- **Untested backups.** Backups are created and verified but never test-restored. When a real restoration is needed, it fails due to a problem that verification alone could not detect. Mitigation: R-D6B-03 mandates monthly test restorations.

- **Backup media degradation.** Physical media degrades over time. A backup that was valid when created becomes unreadable years later. Mitigation: the rotation schedule limits how long any single medium is relied upon. The annual archival snapshot is stored on fresh media each year.

- **All-eggs-in-one-basket.** All backups are stored in the same physical location. A fire, flood, or theft destroys everything -- primary data and all backups. Mitigation: the 3-2-1 rule with mandatory off-site backup (Location C).

- **Encryption key loss for backup media.** The backup encryption key is lost, rendering all backups inaccessible. Mitigation: SEC-003 requires the backup encryption key to exist in three copies with geographic separation. The succession packet includes the key.

## 7. Recovery Procedures

**7.1 If a daily backup was missed:**
1. Perform the backup immediately when the gap is discovered.
2. Record the gap in the backup log: dates missed, reason (if known).
3. If more than 3 consecutive daily backups were missed, also perform a weekly full backup to ensure a complete baseline exists.
4. Review why the backup was missed. If the procedure is too burdensome, consider automation at Level 2 or 3 per D8-001 (supervised or monitored autonomous).

**7.2 If a test restoration fails:**
1. Do not panic. The fact that you discovered the failure during a test is the system working correctly.
2. Identify the failure point: archive corruption, checksum mismatch, extraction error, or content error.
3. Test a different backup of the same data (different date, different location).
4. If multiple backups fail, assess the scope: is the backup tool producing bad output, or has the source data been corrupted?
5. If the backup tool is at fault, fix or replace the tool, then create new backups immediately.
6. If the source data is corrupted, this is a data integrity incident per D6-003.
7. Document everything in the backup log and the operational log.

**7.3 If backup media has failed:**
1. Replace the failed media.
2. Create fresh backups on the new media.
3. Verify the new backups per Section 4.4.
4. If the failed media was the only copy at its location (e.g., the only Location B medium), the institution has dropped below the 3-2-1 threshold. Prioritize creating the replacement immediately.
5. Document the failure in the backup log. Note the media type, age, and any signs of degradation to inform future media selection.

**7.4 If the off-site backup is inaccessible:**
1. Assess why it is inaccessible: location unavailable, media damaged, encryption key missing.
2. If the location is temporarily unavailable, create an additional on-site backup to maintain redundancy until off-site access is restored.
3. If the off-site media is damaged, transport a fresh off-site backup as soon as possible.
4. If the encryption key is missing, follow SEC-003 recovery procedures.
5. The institution should never operate with zero off-site backups for more than 30 days.

## 8. Evolution Path

- **Years 0-5:** The backup procedures are being established. The first year is about building the habit: daily backups, weekly fulls, monthly archives, monthly test restorations. Expect to refine the procedures as you discover what works and what does not. The backup log will be sparse at first and should grow into a rich record.

- **Years 5-15:** The backup procedures should be stable and habitual. The first media replacements will occur as backup drives age out. The first annual archival snapshots will accumulate. This is when the long-term storage strategy begins to be tested by reality.

- **Years 15-30:** Storage technology will have changed. The media types used for backups may need to migrate to new technology. Plan media migrations proactively -- do not wait until the old media can no longer be read. The backup log's history becomes a valuable record of the institution's data resilience over time.

- **Years 30-50+:** The annual archival snapshots span decades. The backup architecture has been through multiple media technology transitions. A successor inheriting the institution should find a clear, documented backup system that they can continue operating from the documentation alone.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The most important sentence in this article is: "A backup that has not been tested is not a backup." I have been burned by this exact failure in previous systems. Backups ran faithfully for months. The logs said "success." The checksums were generated. Everything looked perfect. Then the disk died, the restoration was attempted, and it failed -- because the backup tool had been silently misconfigured and was backing up an empty directory. The verification procedures in this article exist because of that experience.

The air-gap constraint makes backups harder than they need to be in a connected world. I cannot use rsync over SSH to a remote server. I cannot spin up a cloud backup. Every off-site backup requires me to physically carry a drive to another location. This is annoying. It is also the price of sovereignty. The off-site backup that I physically transport is the off-site backup that no one can remotely delete, corrupt, or hold hostage.

The monthly test restoration is the procedure I most dread and most value. It takes time. It feels wasteful when everything is working. And it is the only way to know that the backups are real. I have committed to it in this document, and I am holding myself to it. Future operators: hold yourselves to it too.

## 10. References

- ETH-001 -- Ethical Foundations of the Institution (Principle 2: Integrity Over Convenience)
- CON-001 -- The Founding Mandate (data sovereignty, self-contained operation)
- SEC-001 -- Threat Model and Security Philosophy (Category 1: Physical Threats, Category 2: Data Integrity Threats)
- SEC-002 -- Access Control Procedures (physical zones for backup storage)
- SEC-003 -- Cryptographic Key Management (backup encryption keys)
- OPS-001 -- Operations Philosophy (monthly test restorations, operational tempo)
- D6-001 -- Data Philosophy (data tier system, R-D6-02 redundancy requirements, backup log as Tier 1 data)
- D6-003 -- Data Integrity Verification (complementary integrity procedures)
- D8-001 -- Automation Philosophy (automation levels for backup processes)
- D10-002 -- Daily Operations Doctrine (daily backup verification in checklist)
- D10-003 -- Incident Response Procedures (escalation for backup failures)

---

---

# D6-003 -- Data Integrity Verification

**Document ID:** D6-003
**Domain:** 6 -- Data & Archives
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, SEC-001, OPS-001, D6-001, D6-002
**Depended Upon By:** D10-003. All articles involving data storage, archival, or verification.

---

## 1. Purpose

This article defines the procedures for verifying that the data stored by the holm.chat Documentation Institution has not been corrupted, tampered with, or silently degraded. It answers the question that every backup strategy leaves open: how do you know that the data you have right now is still the data you think it is?

Backups protect against data loss. Integrity verification protects against data corruption -- the silent, insidious degradation that can occur without any visible sign of failure. A file can be corrupted and still be present, still be the right size, still have the right name and modification date. Only a content-level check -- comparing the file's actual content against a known-good reference -- can detect silent corruption.

This article provides specific tools, commands, and procedures for generating integrity baselines, performing periodic verification scans, detecting corruption, and responding when corruption is found. It is a manual for the operator who needs to trust their data over decades.

## 2. Scope

This article covers:

- The theory and practice of data integrity verification.
- Checksum generation: creating integrity baselines for all data.
- Periodic integrity scanning: scheduled verification against baselines.
- Real-time integrity monitoring: filesystem-level protections.
- Corruption detection: what to do when a checksum mismatch is found.
- Integrity verification for backup media.
- Tools and commands for each procedure.

This article does not cover:

- Backup creation procedures (see D6-002).
- Cryptographic signing for authorship verification (see SEC-003 for key management).
- Filesystem selection and configuration (see subsequent Domain 6 implementation articles).
- Data format verification (see D6-001 for format philosophy).

## 3. Background

### 3.1 The Silent Corruption Problem

Data corruption comes in two forms. Loud corruption is obvious: a disk fails, a file becomes unreadable, an application crashes. The operator knows immediately that something is wrong. Silent corruption is invisible: a bit flips in a stored file due to cosmic radiation, media degradation, or a controller error. The file is still readable. It opens without error. But its content has changed. A digit in a spreadsheet is wrong. A word in a document has been altered. A byte in an archive is different, and the archive may or may not decompress correctly depending on where the corruption occurred.

Silent corruption is particularly dangerous for a long-lived institution because it compounds over time. A single bit flip today may not be noticed for years. If backups are made after the corruption occurs, the corruption propagates to the backups. By the time the corruption is detected, there may be no uncorrupted copy left anywhere.

### 3.2 The Checksum Solution

A checksum is a mathematical fingerprint of a file's content. If even a single bit of the file changes, the checksum changes. By computing checksums when data is known to be good and then periodically recomputing and comparing, the institution can detect corruption that would otherwise be invisible.

The choice of checksum algorithm matters. Simple checksums (CRC32) can miss certain corruption patterns. Cryptographic hashes (SHA-256) are computationally more expensive but detect any change with overwhelming probability. For an institution that prioritizes integrity over convenience (ETH-001 Principle 2), SHA-256 or equivalent is the standard.

### 3.3 Integrity Verification as Continuous Practice

Integrity verification is not something you do once. It is a continuous practice -- a regular, scheduled activity that builds confidence in the institution's data over time. Each successful verification is evidence that the data remains intact. Each failed verification is an early warning that allows recovery before the corruption becomes irrecoverable.

The procedures in this article integrate with the operational tempo defined in OPS-001: daily spot-checks during the daily health check, weekly targeted scans during weekly operations, monthly comprehensive scans during the monthly backup verification, and annual full-corpus verification during the annual operations cycle.

## 4. System Model

### 4.1 The Integrity Baseline

The integrity baseline is a complete manifest of checksums for every file in the institution's data store. It is the reference against which all future verification compares. If the baseline is wrong -- if it is computed from already-corrupted data -- the verification will miss the corruption. Therefore, the baseline must be established with care.

**Step-by-step procedure for creating the initial baseline:**

1. Ensure the data is in a known-good state. For the initial baseline, this means the data was recently created, recently imported through the quarantine process per SEC-001, or recently restored from a verified backup.
2. Generate checksums for all files, organized by data tier:
   ```
   find /path/to/tier1-data/ -type f -exec sha256sum {} \; > \
     /path/to/integrity/baselines/tier1-baseline-$(date +%Y-%m-%d).sha256
   ```
   Repeat for Tier 2 and Tier 3 data.
3. Generate a checksum of each baseline file itself (a meta-checksum):
   ```
   sha256sum /path/to/integrity/baselines/tier1-baseline-$(date +%Y-%m-%d).sha256 > \
     /path/to/integrity/baselines/tier1-baseline-$(date +%Y-%m-%d).sha256.meta
   ```
4. Store the baseline files in the integrity baselines directory, which is itself Tier 1 data per D6-001.
5. Create a backup copy of the baseline on Location A per D6-002.
6. Record the baseline creation in the operational log: date, file count, total data size, time taken.

**Updating the baseline:**
The baseline must be updated whenever data legitimately changes -- when files are added, modified, or deleted as part of normal operations. There are two approaches:

- **Incremental update:** After making changes, regenerate checksums only for the changed files and merge them into the baseline.
   ```
   sha256sum /path/to/changed-file > /tmp/new-checksum.txt
   # Manually replace the old checksum line in the baseline file
   ```
- **Periodic full regeneration:** On a scheduled basis (monthly is recommended for active data), regenerate the entire baseline. This catches any changes the incremental process may have missed.
   ```
   find /path/to/tier1-data/ -type f -exec sha256sum {} \; > \
     /path/to/integrity/baselines/tier1-baseline-$(date +%Y-%m-%d).sha256
   ```

The recommended practice is incremental updates during daily operations and full regeneration monthly.

### 4.2 Periodic Integrity Scanning

Integrity scanning compares current file checksums against the baseline. Mismatches indicate either legitimate changes that were not recorded in the baseline (a process error) or corruption (a data integrity event).

**Daily spot-check (during daily operations, 2-3 minutes):**
Select a small random sample of files from Tier 1 data and verify their checksums:
```
# Select 10 random files from the Tier 1 baseline
shuf -n 10 /path/to/integrity/baselines/tier1-baseline-current.sha256 > /tmp/daily-check.sha256
sha256sum -c /tmp/daily-check.sha256
```
If all checks pass, record "Daily integrity spot-check: PASS (10 files)" in the operational log.
If any check fails, investigate immediately per Section 4.4.

**Weekly targeted scan (during weekly operations, 15-30 minutes):**
Verify all Tier 1 data:
```
sha256sum -c /path/to/integrity/baselines/tier1-baseline-current.sha256
```
Record the results in the operational log. Note any mismatches and their resolution.

**Monthly comprehensive scan (during monthly operations, 1-2 hours):**
Verify all Tier 1 and Tier 2 data:
```
sha256sum -c /path/to/integrity/baselines/tier1-baseline-current.sha256
sha256sum -c /path/to/integrity/baselines/tier2-baseline-current.sha256
```
Also verify the integrity of the most recent backup per D6-002 Section 4.4.

**Annual full-corpus verification (during annual operations, half day or more):**
Verify all data across all tiers:
```
sha256sum -c /path/to/integrity/baselines/tier1-baseline-current.sha256
sha256sum -c /path/to/integrity/baselines/tier2-baseline-current.sha256
sha256sum -c /path/to/integrity/baselines/tier3-baseline-current.sha256
```
Regenerate the baselines after successful verification to create a fresh reference point.

### 4.3 Filesystem-Level Integrity

Beyond manual checksum verification, the institution should use a filesystem that provides built-in integrity checking. The recommended approach is to use a copy-on-write filesystem with built-in checksumming (such as ZFS or Btrfs). These filesystems compute checksums for every data block and verify them on every read, detecting silent corruption in real time rather than waiting for a scheduled scan.

**If using ZFS:**
- Enable checksums (enabled by default): `zfs get checksum [pool/dataset]`
- Perform a scrub monthly (this verifies all data on the pool against stored checksums):
   ```
   zpool scrub [pool_name]
   ```
- Check scrub results:
   ```
   zpool status [pool_name]
   ```
   Look for "errors: No known data errors" on a healthy pool, or specific error counts if corruption was found.
- ZFS scrubs should be part of the monthly operations cycle per OPS-001.

**If using Btrfs:**
- Perform a scrub monthly:
   ```
   btrfs scrub start /path/to/mount
   ```
- Check scrub results:
   ```
   btrfs scrub status /path/to/mount
   ```

**If using a filesystem without built-in integrity (ext4, XFS):**
Manual checksum verification per Section 4.2 is the only defense. Consider migrating to a checksumming filesystem when a suitable opportunity arises (Tier 2 decision per GOV-001).

### 4.4 Corruption Detection and Response

When a checksum mismatch is detected, follow this procedure:

**Step 1: Confirm the mismatch.**
Re-run the checksum on the flagged file:
```
sha256sum /path/to/flagged-file
```
Compare against the baseline. If the checksum now matches, the first result was a transient error (rare but possible). Record the incident and monitor for recurrence.

**Step 2: Determine if the change is legitimate.**
Check whether the file was intentionally modified since the baseline was created. Consult the operational log: was this file edited as part of a documented change? Check the file's modification timestamp against the baseline creation date.

If the change is legitimate: update the baseline with the new checksum. Record the update in the operational log.

**Step 3: If the change is not legitimate (corruption confirmed):**
1. Do not modify the corrupted file. Leave it in place for analysis.
2. Record the corruption in the operational log: file path, expected checksum, actual checksum, file size, modification date.
3. Classify the corruption severity:
   - **Single file, Tier 3 or Tier 4:** Low severity. Restore from the most recent good backup.
   - **Single file, Tier 1 or Tier 2:** Medium severity. Restore from backup and investigate the cause.
   - **Multiple files:** High severity. This may indicate a media failure, filesystem corruption, or something more serious. Escalate per D10-003.
   - **Corruption in both primary and backup copies:** Critical severity. The corruption predates the backups. Search for an older uncorrupted backup. If found, restore from it. If not found, document the permanent data loss.
4. Restore the uncorrupted version from the most recent backup that predates the corruption:
   ```
   cp /path/to/backup/copy/of/file /path/to/original/location/
   ```
5. Verify the restored file:
   ```
   sha256sum /path/to/original/location/file
   ```
   Compare against the baseline.
6. Investigate the cause: check filesystem health, storage device SMART data, system logs for I/O errors.
7. If the cause is media degradation, plan media replacement per the evolution path in this article.

### 4.5 Integrity Verification for Backup Media

Backup media is subject to the same corruption risks as primary storage. The procedures in D6-002 Section 4.4 verify backups at creation. This section addresses ongoing integrity verification of stored backups.

**Monthly: verify most recent Location A backup:**
```
sha256sum -c /path/to/backup-location-a/tier1/manifests/most-recent.sha256
```

**Quarterly: verify a randomly selected Location B backup:**
Select one of the Location B archival media. Mount it. Verify the archives:
```
sha256sum -c /path/to/location-b/manifest.sha256
```

**Annually: verify the off-site backup (Location C):**
During the annual operations cycle, retrieve or visit the off-site location. Mount the off-site backup media. Verify the checksums. Return the media.

If backup media fails verification, the backup is unreliable and must be replaced per D6-002 recovery procedures.

## 5. Rules & Constraints

- **R-D6I-01:** An integrity baseline must exist for all Tier 1 and Tier 2 data. The baseline must be stored as Tier 1 data per D6-001.
- **R-D6I-02:** Daily spot-checks of Tier 1 data integrity must be performed as part of the daily operations cycle per D10-002.
- **R-D6I-03:** Weekly full verification of Tier 1 data must be performed during the weekly operations cycle.
- **R-D6I-04:** Monthly comprehensive verification of Tier 1 and Tier 2 data must be performed during the monthly operations cycle.
- **R-D6I-05:** Annual full-corpus verification of all tiers must be performed during the annual operations cycle.
- **R-D6I-06:** All checksum computations must use SHA-256 or a stronger algorithm. Weaker algorithms (MD5, CRC32) are not sufficient for integrity verification.
- **R-D6I-07:** When corruption is detected, the response procedure in Section 4.4 must be followed. Corruption must not be ignored or dismissed without investigation.
- **R-D6I-08:** If a copy-on-write filesystem with built-in checksumming is used, filesystem scrubs must be performed at least monthly.
- **R-D6I-09:** Integrity verification results must be recorded in the operational log. Both successes and failures are recorded.

## 6. Failure Modes

- **Baseline drift.** The baseline is not updated when files legitimately change. Over time, every file mismatches because the baseline is outdated. The operator starts ignoring mismatches because "they are all expected." Real corruption hides in the noise of expected mismatches. Mitigation: the incremental update procedure and monthly full regeneration. Mismatches due to baseline drift are a process failure, not a data failure, but they undermine the verification system just as effectively.

- **Verification theater.** The scans run but the results are not reviewed. The operator sees a wall of "OK" lines and does not notice the three "FAILED" lines buried among them. Mitigation: pipe verification output through a filter that highlights only failures:
   ```
   sha256sum -c baseline.sha256 2>&1 | grep -v ': OK$'
   ```
   If this command produces no output, all checks passed. Any output means a failure that requires attention.

- **Corrupted baseline.** The baseline file itself is corrupted, leading to false positives or false negatives. Mitigation: the meta-checksum (a checksum of the baseline file) in Section 4.1 step 3. Verify the baseline's own integrity before using it for verification.

- **Media-class failure.** An entire class of storage media fails simultaneously (a manufacturing defect, a bad batch of drives). If all copies of the data are on the same media class, all copies may be corrupted. Mitigation: the two-media-type requirement in D6-002 (the 3-2-1 rule). The annual verification of all tiers catches cross-media issues.

- **Verification too infrequent.** Corruption occurs but is not detected until long after the event. By the time it is detected, the corruption has propagated to all backups. Mitigation: the tiered verification schedule (daily spot-checks, weekly full Tier 1, monthly comprehensive) is designed to catch corruption before it propagates to all backup generations.

- **False sense of security from filesystem scrubs.** The operator relies on ZFS or Btrfs scrubs and stops performing manual checksum verification. The filesystem catches block-level corruption but may not catch application-level corruption (e.g., a bug that writes incorrect data to a file). Mitigation: filesystem scrubs and manual checksum verification serve different purposes. Both are necessary. Neither replaces the other.

## 7. Recovery Procedures

**7.1 If baseline drift has made integrity scanning useless:**
1. Stop relying on the outdated baseline.
2. Assess the current state of the data: is there any reason to believe the data is corrupted?
3. If the data is believed to be good (no errors reported, no anomalies observed): regenerate the baseline from the current data. This is a new starting point. Accept that you cannot detect corruption that occurred between the old baseline and now.
4. If the data's integrity is uncertain: restore from the most recent verified backup, then regenerate the baseline from the restored data.
5. Implement the incremental update procedure to prevent future drift.
6. Document the incident in the operational log.

**7.2 If corruption is found in both primary data and backups:**
1. This is a critical event. Record it in the operational log and the decision log.
2. Determine when the corruption occurred by checking historical baselines and backup manifests. Find the oldest backup that still has the corrupted checksum and the newest backup that has the correct checksum.
3. Restore from the newest backup with the correct checksum.
4. If no uncorrupted backup exists: the data is permanently corrupted. Document the loss per D6-001 Recovery Procedure 5.
5. Investigate the cause thoroughly. If the cause is a software bug, assess what other data may be affected.
6. Escalate per D10-003.

**7.3 If the integrity verification system itself has been neglected:**
1. Start with the daily spot-check. Run it today.
2. Within the current week, run a full Tier 1 verification.
3. Within the current month, run a comprehensive Tier 1 and Tier 2 verification.
4. If significant corruption is found, follow Section 4.4.
5. If no corruption is found, consider yourself fortunate and resume the scheduled verification cycle.
6. Document the gap in the operational log.

## 8. Evolution Path

- **Years 0-5:** The integrity verification system is being established. Create the initial baselines. Build the habit of daily spot-checks and weekly scans. Learn the tools. Discover the operational rhythm. The first years will also reveal whether the chosen filesystem's built-in integrity features are sufficient or whether additional measures are needed.

- **Years 5-15:** The verification system should be routine. The baselines have been regenerated many times. The monthly and annual scans should have a history of results in the operational log. This history itself is valuable: it reveals patterns of media degradation, filesystem behavior, and data growth.

- **Years 15-30:** Storage technology transitions will require rebuilding baselines on new media. The transition is an opportunity: when data is migrated to new storage, verify every file during the migration. The migration checksum verification serves double duty as an integrity check and a migration validation.

- **Years 30-50+:** Decades of integrity verification data in the operational log. The institution should be able to point to a continuous chain of verification results that demonstrate data integrity across its entire lifetime. This chain is one of the institution's most valuable artifacts.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The daily spot-check -- ten random files, two minutes -- is the procedure I am most proud of in this article. Not because it is technically sophisticated (it is trivially simple) but because it is sustainable. I can do it when I am tired. I can do it when I am busy. I can do it when I have spent the day dealing with a crisis and have nothing left. Ten files. Two minutes. And those two minutes, repeated every day for decades, build a continuous chain of evidence that the institution's most critical data is intact.

The danger I fear most is not catastrophic disk failure. That is dramatic and obvious and triggers immediate response. The danger I fear most is the single bit flip in a file I do not check for three years, which propagates through five backup generations, and which I discover only when I try to open a document that has become subtly, silently wrong. The verification schedule in this article is designed to catch that bit flip before it can hide.

I want to be honest: I will sometimes find the weekly Tier 1 scan tedious. It takes time. It produces the same "all OK" result week after week. And then one week it will not. That is the week that justifies every scan that preceded it.

## 10. References

- ETH-001 -- Ethical Foundations of the Institution (Principle 2: Integrity Over Convenience; Principle 6: Honest Accounting of Limitations)
- CON-001 -- The Founding Mandate (data sovereignty, preservation mandate)
- SEC-001 -- Threat Model and Security Philosophy (Category 2: Data Integrity Threats, bit rot)
- OPS-001 -- Operations Philosophy (operational tempo for verification cycles)
- D6-001 -- Data Philosophy (data tier system, integrity as institutional priority)
- D6-002 -- Backup Doctrine (backup verification, complementary protection)
- SEC-003 -- Cryptographic Key Management (hash algorithms for checksums)
- D10-002 -- Daily Operations Doctrine (daily spot-check in operational cycle)
- D10-003 -- Incident Response Procedures (escalation for corruption events)

---

---

# D10-003 -- Incident Response Procedures

**Document ID:** D10-003
**Domain:** 10 -- User Operations
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, SEC-001, SEC-002, SEC-003, OPS-001, D6-001, D6-002, D6-003, D10-002
**Depended Upon By:** All articles that reference incident handling or escalation.

---

## 1. Purpose

This article defines what to do when something goes wrong. It provides classification criteria for incidents, response procedures for each severity level, decision trees for common scenarios, escalation paths, communication protocols, and the post-incident review process.

In a single-operator, air-gapped institution, incident response does not look like the corporate process of war rooms, conference calls, and incident commanders. The operator is alone. There is no one to call for help. The operator is simultaneously the detector, the investigator, the responder, the communicator, and the post-incident reviewer. This article is written for that reality.

The core principle of this article is: when something goes wrong, stop, think, document, then act. The natural human response to an incident is urgency -- to fix it immediately, to make it stop, to restore normality. That urgency is dangerous. Hasty action without understanding can turn a recoverable incident into a catastrophic one. The procedures in this article are designed to impose structure on the operator's response, ensuring that investigation precedes remediation, that documentation happens in real time rather than after the fact, and that the root cause is addressed rather than merely the symptoms.

## 2. Scope

This article covers:

- The definition of an incident and how it differs from a deviation.
- The severity classification system: how to assess an incident's severity.
- Response procedures for each severity level.
- Decision trees for common incident types.
- The post-incident review process.
- The incident log and its relationship to the operational log and the decision log.

This article does not cover:

- Specific backup restoration procedures (see D6-002).
- Specific data integrity verification procedures (see D6-003).
- Specific access control breach procedures (see SEC-002).
- Specific key compromise procedures (see SEC-003).
- Full disaster recovery (see subsequent Domain 12 articles).

Those articles provide the detailed technical procedures for specific types of recovery. This article provides the incident management framework that determines when and how those procedures are invoked.

## 3. Background

### 3.1 The Difference Between a Deviation and an Incident

Per D10-002 Section 4.4, deviations from the daily checklist are classified as minor, moderate, or critical. An incident is a specific type of deviation: an event that has caused or may cause harm to the institution's data, systems, security, or operational capability. Not every deviation is an incident. A slightly slow backup is a deviation. A backup that produces corrupt output is an incident. A high CPU temperature reading is a deviation. A disk reporting SMART errors is an incident.

The distinction matters because incidents require a structured response process -- investigation, classification, remediation, documentation, and review. Deviations may require only a note in the operational log and monitoring for recurrence.

When in doubt about whether an event is a deviation or an incident, classify it as an incident. It is better to over-respond and discover the situation is benign than to under-respond and discover it is serious.

### 3.2 Why Single-Operator Incident Response Is Hard

In a team environment, incident response benefits from parallel processing: one person investigates while another documents while another prepares remediation. A single operator must do all of these sequentially. The temptation is to skip the documentation and go straight to fixing the problem. This temptation must be resisted. Documentation during an incident is not overhead. It is the mechanism by which the operator maintains clarity, avoids circular investigation, and creates the record that enables the post-incident review.

The procedures in this article are designed for sequential execution by a single person. They are deliberately structured to prevent the operator from becoming overwhelmed by requiring one step at a time, with explicit pause points for assessment.

### 3.3 The Post-Incident Review Is Not Optional

Every incident, regardless of severity, receives a post-incident review. The review is not a punishment. It is not a search for blame. In a single-operator institution, blame is meaningless -- there is only one person. The review is a learning opportunity. It asks: what happened, why did it happen, how was it detected, how was it resolved, and what changes should be made to prevent recurrence?

Incidents that are resolved but not reviewed will recur. The institution does not learn from incidents it does not study.

## 4. System Model

### 4.1 Severity Classification

Every incident is classified into one of four severity levels. The classification determines the response procedure, the documentation requirements, and the urgency of remediation.

**Severity 1 (S1): Critical.**
Definition: An event that has caused or is actively causing irreversible harm to the institution's data, security, or operational capability. The institution's ability to fulfill its mission per CON-001 is in immediate jeopardy.

Examples:
- Loss of data with no known backup.
- Confirmed compromise of encryption keys.
- Physical breach of the equipment zone with evidence of tampering.
- Multiple simultaneous disk failures affecting primary data and backups.
- Evidence of unauthorized modification of institutional memory (Tier 1 data).

Response timeframe: Immediate. All other operations are suspended.

**Severity 2 (S2): High.**
Definition: An event that has caused significant harm or poses a significant risk of harm if not addressed promptly. The institution is still operational but its resilience is degraded.

Examples:
- Failure of a single disk with data recoverable from backup.
- Corruption detected in Tier 1 or Tier 2 data with known-good backups available.
- Failure of a backup device (one copy of the 3-2-1 is missing).
- Suspected but unconfirmed security breach.
- Failure of a service account or automated system that handles critical data.

Response timeframe: Within 24 hours. Daily operations may continue but the incident takes priority over non-essential tasks.

**Severity 3 (S3): Medium.**
Definition: An event that has caused limited harm or inconvenience. The institution is fully operational and resilient, but something needs attention.

Examples:
- Corruption detected in Tier 3 data with backup available.
- A non-critical service has failed.
- A backup completed with warnings but verification indicates the backup is usable.
- SMART data indicates a disk is approaching end of life (not yet failing).
- An access control anomaly that was benign on investigation.

Response timeframe: Within 7 days. Normal operations continue.

**Severity 4 (S4): Low.**
Definition: An event that was detected by monitoring, investigated, and found to be benign or trivially resolvable. Documented for the record and for pattern analysis.

Examples:
- A transient checksum mismatch that resolved on re-check.
- A service that restarted automatically after a brief failure.
- An unexpected log entry that was traced to normal operation.
- A minor configuration inconsistency corrected in minutes.

Response timeframe: At the operator's convenience, within the current operational cycle.

### 4.2 The Incident Response Procedure

Every incident, regardless of severity, follows this seven-step procedure. The steps are the same; the depth and urgency differ by severity level.

**Step 1: Detect and Record.**
The incident is detected -- through the daily checklist, through an alert, through observation, or through integrity verification. The moment it is detected, open the incident log and record:
- Date and time of detection.
- How the incident was detected.
- Initial description of the event (what is happening or what was observed).
- Preliminary severity assessment (this may change as investigation proceeds).

Do not skip this step. Do not say "I will document it after I fix it." The documentation starts now.

**Step 2: Assess and Classify.**
Determine the severity level using the criteria in Section 4.1. Ask:
- Has irreversible harm occurred or is it imminent? (If yes: S1.)
- Has significant harm occurred or is the institution's resilience degraded? (If yes: S2.)
- Is the harm limited and the institution fully operational? (If yes: S3.)
- Was the event benign or trivially resolved? (If yes: S4.)

Record the classification in the incident log with your reasoning.

**Step 3: Contain.**
Prevent the incident from causing further harm. The containment action depends on the incident type:
- For data corruption: stop writes to the affected area. Isolate the affected storage.
- For security breach: disable the compromised account or service. Do not destroy evidence.
- For hardware failure: disconnect the failing device if possible. Switch to backup or redundant hardware.
- For software failure: stop the failing service. Fall back to manual procedures per D8-001.

Record the containment actions in the incident log.

**Step 4: Investigate.**
Determine what happened, why it happened, and what the full scope of the impact is. This is the most time-consuming step and the most important. Do not rush it.

Investigation checklist:
- What changed recently? (Consult the operational log and decision log.)
- What are the system logs showing? (Check syslog, application logs, authentication logs.)
- Is the problem isolated or widespread? (Check other systems, other storage devices, other services.)
- When did the problem start? (Trace backward through logs to find the earliest indication.)
- Is the problem ongoing or has it stopped? (Monitor for continued activity.)

Record all findings in the incident log as they are discovered. This creates a real-time narrative that will be invaluable during the post-incident review.

**Step 5: Remediate.**
Fix the problem. The specific remediation depends on the incident type:
- For data corruption: restore from the most recent verified backup per D6-002.
- For security breach: revoke compromised credentials per SEC-002/SEC-003. Rebuild affected systems if necessary.
- For hardware failure: replace the failing component. Restore data if needed.
- For software failure: fix the configuration, update the software, or replace the service.

Record all remediation actions in the incident log.

**Step 6: Verify.**
Confirm that the remediation was successful and that the incident is truly resolved.
- Verify data integrity after restoration per D6-003.
- Verify access controls after a security incident per SEC-002.
- Verify the affected service is operational.
- Verify that the root cause has been addressed, not just the symptoms.

Record the verification results in the incident log.

**Step 7: Close and Review.**
Close the incident in the incident log. Record:
- Total duration of the incident (from detection to closure).
- Final severity classification (may differ from initial assessment).
- Summary of root cause.
- Summary of remediation actions.
- Impact assessment: what was affected, what was lost (if anything), what was the operational downtime.

Then conduct the post-incident review per Section 4.5.

### 4.3 Decision Trees for Common Incidents

**Decision Tree 1: Checksum Mismatch Detected**
```
Checksum mismatch detected
|
+-- Re-check: does mismatch persist?
    |
    +-- No: Transient error. Record as S4. Monitor for recurrence.
    |
    +-- Yes: Was the file legitimately modified?
        |
        +-- Yes: Update baseline. Record as process deviation, not incident.
        |
        +-- No: CORRUPTION CONFIRMED.
            |
            +-- How many files affected?
                |
                +-- Single file, Tier 3/4: S3. Restore from backup.
                +-- Single file, Tier 1/2: S2. Restore from backup. Investigate cause.
                +-- Multiple files: S2 minimum (S1 if widespread).
                    Isolate storage. Check backups for corruption.
                    |
                    +-- Backups clean? Restore from backup. Investigate media health.
                    +-- Backups also corrupt? S1. Find oldest clean backup. Accept data loss if none exists.
```

**Decision Tree 2: Disk Reports SMART Errors**
```
SMART errors detected
|
+-- Error type?
    |
    +-- Reallocated sectors (increasing): Disk is degrading.
    |   +-- Fewer than 10 sectors: S3. Monitor weekly. Plan replacement.
    |   +-- More than 10 sectors: S2. Verify all data on disk. Plan immediate replacement.
    |   +-- Rapidly increasing: S2. Immediate backup. Replace as soon as possible.
    |
    +-- Pending sectors: S3. Monitor. May resolve or escalate.
    |
    +-- Temperature warning: S4. Check cooling. Record in log.
    |
    +-- Self-test failure: S2. Full data verification. Plan replacement.
```

**Decision Tree 3: Unauthorized Account or Process Detected**
```
Unrecognized account or process found
|
+-- Is it documented in the service account registry (SEC-002)?
    |
    +-- Yes: Verify it matches documentation. If match: S4 false alarm.
    |         If mismatch: S3. Investigate discrepancy.
    |
    +-- No: Undocumented account/process.
        |
        +-- Can it be traced to a legitimate action? (Software install, system update)
            |
            +-- Yes: S3. Document the account. Update registry.
            |
            +-- No: S2 minimum.
                Disable the account/process immediately.
                Check access logs. Check for other anomalies.
                |
                +-- Evidence of unauthorized activity? S1. Full security audit.
                +-- No evidence beyond the account itself? S2. Investigate origin. Monitor.
```

**Decision Tree 4: Backup Failure**
```
Backup operation failed
|
+-- What type of failure?
    |
    +-- Exit code error (non-zero):
    |   +-- Known error (disk full, permissions): S3. Fix cause. Retry backup.
    |   +-- Unknown error: S2. Investigate. Do not proceed until cause is understood.
    |
    +-- Verification failure (checksum mismatch, size mismatch):
    |   +-- Retry backup. If retry succeeds: S4. Record the transient failure.
    |   +-- Retry fails: S2. Investigate backup tool, source data, destination media.
    |
    +-- Media failure (cannot write to backup device):
        +-- S2. Replace media. Backup to alternate location immediately.
        +-- If no alternate backup exists: S1 (single point of failure).
```

**Decision Tree 5: Power Anomaly (Off-Grid System)**
```
Power anomaly detected (voltage fluctuation, battery alert, generator issue)
|
+-- Is the system on battery/UPS backup?
    |
    +-- Yes: How much runtime remains?
    |   +-- More than 30 minutes: S3. Investigate power source. Monitor battery.
    |   +-- Less than 30 minutes: S2. Initiate graceful shutdown procedures if power source cannot be restored.
    |   +-- Less than 5 minutes: S1. Immediate graceful shutdown. Data integrity at risk.
    |
    +-- No (direct power, no UPS):
        +-- S1 if power is failing. Shut down immediately to prevent corruption.
        +-- Record the gap in UPS coverage as a separate S3 issue.
```

### 4.4 The Incident Log

The incident log is separate from the operational log. It is dedicated to incident records and post-incident reviews. It is a Tier 1 data artifact per D6-001.

Each incident log entry contains:
- **Incident ID:** Sequential identifier (INC-YYYY-NNN format).
- **Date/time of detection.**
- **Severity level (initial and final if changed).**
- **Detection method:** How the incident was discovered.
- **Description:** What happened, in clear prose.
- **Containment actions:** What was done to stop the bleeding.
- **Investigation notes:** What was found, in chronological order.
- **Root cause:** Why it happened (if determinable).
- **Remediation actions:** What was done to fix it.
- **Verification results:** How it was confirmed fixed.
- **Impact assessment:** What was affected, what was lost, how long the incident lasted.
- **Post-incident review date:** When the PIR was completed.
- **Lessons learned and action items:** From the PIR.

### 4.5 The Post-Incident Review (PIR)

Every incident receives a PIR. The timing depends on severity:
- S1: PIR within 48 hours of incident closure.
- S2: PIR within 7 days of incident closure.
- S3: PIR within 30 days, or during the next monthly operations cycle.
- S4: PIR may be consolidated -- review all S4 incidents in a batch during the monthly or quarterly cycle.

**PIR procedure:**

1. **Re-read the incident log entry.** Start by re-reading the entire chronological record of the incident.

2. **Answer the five questions:**
   - **What happened?** A clear, factual summary of the event.
   - **Why did it happen?** The root cause. Use "five whys" analysis: ask "why?" at least five times until you reach a systemic cause rather than a proximate one.
     - Example: The backup failed. Why? The disk was full. Why? Old backups were not rotated. Why? The rotation script had a bug. Why? The script was not tested after the last modification. Why? The testing procedure is not triggered by script modifications. Root cause: the testing procedure gap.
   - **How was it detected?** Which monitoring, checklist, or verification process caught it? If it was caught by accident, that is itself a finding.
   - **How was it resolved?** Summary of the remediation.
   - **What should change?** Specific, actionable changes to prevent recurrence.

3. **Generate action items.** Each action item has:
   - A description of the change.
   - A responsible party (in a single-operator institution, this is always the operator -- but the item is still recorded formally).
   - A deadline.
   - A GOV-001 tier classification if the change requires a governance decision.

4. **Record the PIR in the incident log**, appended to the original incident entry.

5. **Track action items** in the decision log or operational log until completed.

## 5. Rules & Constraints

- **R-D10I-01:** All incidents must be classified using the four-level severity system defined in Section 4.1.
- **R-D10I-02:** Documentation begins at the moment of detection, not after resolution. Incident log entries must be written in real time during the incident.
- **R-D10I-03:** Every incident, regardless of severity, receives a post-incident review per Section 4.5.
- **R-D10I-04:** The incident log is a Tier 1 data artifact per D6-001 and is subject to all Tier 1 preservation requirements.
- **R-D10I-05:** S1 incidents require immediate response. All other operations are suspended until the S1 incident is contained.
- **R-D10I-06:** When in doubt about severity, classify higher rather than lower. An S2 response to an S3 event wastes some time. An S3 response to an S2 event risks harm.
- **R-D10I-07:** Post-incident review action items must be tracked to completion. Open action items are reviewed during the quarterly operations cycle per OPS-001.
- **R-D10I-08:** The incident log must be reviewed during the annual operations cycle for patterns: recurring incidents, increasing frequency, or categories of incident that the current procedures do not adequately address.
- **R-D10I-09:** Incident response procedures must be printed on physical media and stored in Zone 1 per SEC-002, per OPS-001 R-OPS-07. When the systems are down, the procedures for recovering them must be accessible without those systems.

## 6. Failure Modes

- **Panic response.** The operator detects an incident and immediately starts trying to fix it without stopping to assess, classify, or document. Actions are taken out of order. Evidence is destroyed. The remediation creates new problems. Mitigation: the seven-step procedure is designed to impose structure. Step 1 is always "record." Step 2 is always "assess." The fix does not come until Step 5.

- **Severity underestimation.** The operator classifies an incident at a lower severity than it deserves, either through optimism ("it is probably fine") or through fatigue ("I do not want to deal with this right now"). The response is inadequate, and the incident escalates. Mitigation: R-D10I-06 -- when in doubt, classify higher. The decision trees in Section 4.3 provide objective criteria.

- **Documentation debt after incidents.** The incident is resolved, but the documentation is incomplete or never written. The post-incident review cannot be conducted effectively. The same incident recurs because the lessons were not captured. Mitigation: R-D10I-02 requires real-time documentation. R-D10I-03 requires PIR for every incident.

- **PIR action items that never happen.** The post-incident review generates action items. The action items are recorded. Then nothing changes. The same vulnerability persists, and the same incident recurs. Mitigation: R-D10I-07 requires tracking to completion. The quarterly review checks for open action items.

- **Incident fatigue.** In a period of many incidents, the operator becomes desensitized. Severity creep occurs: what would have been S2 last month is now S3 because the operator is tired of high-severity responses. Mitigation: the severity criteria in Section 4.1 are objective, not relative. They do not change based on the operator's current stress level. If many incidents are occurring, the pattern itself is a finding that requires investigation during the PIR process.

- **Inaccessible procedures during a crisis.** The system is down. The incident response procedures are on the system. The operator cannot access them. Mitigation: R-D10I-09 requires physical copies of incident response procedures stored in Zone 1. Print this article. Keep it where you can reach it when the screens are dark.

## 7. Recovery Procedures

**7.1 If the incident response process itself has broken down:**
1. Stop. Take a breath.
2. Open the incident log. If the incident log is inaccessible (because it is on the system that is down), use paper. Literal paper.
3. Write down what you know: what happened, when, what you have done so far.
4. Classify the incident using the criteria in Section 4.1.
5. Follow the seven-step procedure from whatever step you are at. If you are not sure what step you are at, start at Step 4 (Investigate).
6. After the incident is resolved, reconstruct the timeline and fill in the incident log.

**7.2 If too many incidents are occurring simultaneously:**
1. Triage by severity. S1 incidents first, then S2, then S3. S4 incidents wait.
2. If multiple S1 incidents are simultaneous, prioritize by irreversibility: address the incident that will cause the most permanent harm first.
3. Contain all incidents before attempting to remediate any of them. Containment prevents escalation while you work through the queue.
4. After all incidents are contained, remediate in severity order.
5. The simultaneous incident pattern is itself a finding: something systemic may be wrong. Record this for the PIR.

**7.3 If the operator is not capable of responding (injury, illness, absence):**
1. This is a succession scenario. The institution's defenses are the automated systems that continue to run and the documentation that enables a future response.
2. When the operator returns (or a successor assumes responsibility), review all logs from the period of absence.
3. Run a full integrity verification per D6-003.
4. Address any incidents discovered retroactively, classifying them at the time of discovery rather than the time of occurrence.
5. If the period of absence was extended (weeks or more), conduct a comprehensive system audit before resuming normal operations.

**7.4 If the post-incident review reveals a systemic failure:**
1. A systemic failure is a pattern -- not a single incident but a class of incidents that share a root cause.
2. Treat the systemic failure as its own Tier 2 or Tier 3 governance decision per GOV-001.
3. Propose and implement structural changes (new procedures, new monitoring, revised architecture) rather than patching individual incidents.
4. Track the structural change through the decision log.
5. Monitor for recurrence over the following 90 days.

## 8. Evolution Path

- **Years 0-5:** The incident response process is being tested for the first time. Expect the severity criteria to need refinement as real incidents reveal edge cases the current classification does not cleanly address. The decision trees will grow as new incident types are encountered. Every incident in this period is an opportunity to improve the process.

- **Years 5-15:** The incident log should contain enough history to reveal patterns. The annual pattern review should be a substantive exercise. The PIR process should be a natural part of operations rather than an imposition.

- **Years 15-30:** A successor may need to use these procedures. The decision trees and the incident log history are their most valuable resources. A successor who can read the institution's incident history understands the institution's vulnerabilities better than any architecture document can convey.

- **Years 30-50+:** Decades of incident history. The institution has survived dozens or hundreds of incidents. The incident response process has been refined through experience. The PIR archive is an institutional resource that informs every aspect of the institution's design and operation.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
I wrote the decision trees before I wrote the rest of this article. They are the section I expect to use most often, because when something goes wrong, I will not have the patience to read paragraphs of prose. I will need a branching set of questions that leads me to the correct action. The trees are deliberately simple. They do not cover every scenario. They cover the common ones -- the ones that account for the vast majority of incidents in a system like this. Rare and exotic incidents will require the full seven-step procedure and genuine investigative thought.

The hardest discipline in incident response is Step 1: stop and document before you act. I know this because I have violated it many times in my career. Something breaks, and every instinct says "fix it now." The instinct is wrong. Fixing a problem you do not understand is how you turn a small problem into a big one. I have written this article partly as a contract with my future self: when something goes wrong, you will open the log before you open the terminal. You will write down what you see before you type a command. You will think before you act.

The post-incident review is the learning engine of this institution. Without it, incidents are isolated events. With it, incidents become lessons that accumulate into wisdom. I have committed to reviewing every incident, even the trivial ones. The trivial ones sometimes reveal the most interesting patterns.

One final note: print this article. Put it in a binder. Put the binder next to the equipment. When the system is down and the screens are dark, that binder is the only incident response resource you have. It must be there.

## 10. References

- ETH-001 -- Ethical Foundations of the Institution (Principle 2: Integrity Over Convenience; Principle 6: Honest Accounting of Limitations)
- CON-001 -- The Founding Mandate (institutional mission, resilience mandate)
- GOV-001 -- Authority Model (decision tiers for structural changes from PIR findings)
- SEC-001 -- Threat Model and Security Philosophy (threat categories, defense in depth)
- SEC-002 -- Access Control Procedures (access breach response, unauthorized account handling)
- SEC-003 -- Cryptographic Key Management (compromised key response)
- OPS-001 -- Operations Philosophy (operational tempo, quarterly review, annual audit, R-OPS-07 printed procedures)
- D6-001 -- Data Philosophy (data tier system for impact assessment, incident log as Tier 1 data)
- D6-002 -- Backup Doctrine (backup failure handling, restoration procedures)
- D6-003 -- Data Integrity Verification (corruption detection, integrity scanning)
- D8-001 -- Automation Philosophy (automated system failure, manual fallback)
- D10-002 -- Daily Operations Doctrine (deviation handling, daily checklist as detection mechanism)

---

---

*End of Stage 3 Operational Doctrine -- Batch 1*

**Document Total:** 5 articles
**Domains Covered:** Domain 3 -- Security & Integrity (SEC-002, SEC-003), Domain 6 -- Data & Archives (D6-002, D6-003), Domain 10 -- User Operations (D10-003)
**Combined Estimated Word Count:** ~15,000 words
**Status:** All five articles ratified as of 2026-02-16.
**Relationship to Stage 2:** These articles implement the philosophies established in SEC-001, OPS-001, D6-001, and D10-002. They translate principles into procedures. Where Stage 2 says "what we believe," Stage 3 says "what we do."
**Next Stage:** Additional Stage 3 operational articles for remaining domains, followed by Stage 4 implementation-specific articles.
