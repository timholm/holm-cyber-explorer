# STAGE 4: SPECIALIZED SECURITY SYSTEMS

## Advanced Reference Documents for Air-Gap Architecture, Cryptographic Survival, Physical Security, Supply Chain Security, and Audit Framework

**Document ID:** STAGE4-SECURITY-ADVANCED
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Classification:** Specialized Systems -- These articles provide deep technical specifications for critical security domains. They extend Stage 3 operational procedures into comprehensive architectural references designed to endure across hardware generations and threat landscape shifts.

---

## How to Read This Document

This document contains five advanced security articles that belong to Stage 4 of the holm.chat Documentation Institution. Where Stage 2 established philosophy and Stage 3 established procedures, Stage 4 provides the deep technical reference material that a competent operator needs when designing, auditing, or rebuilding a security-critical subsystem from first principles.

These articles are not manuals in the Step 1, Step 2 sense. They are architectural specifications. They describe the theory behind the practice, the constraints that govern the design, the failure modes that threaten the design, and the evolution path that keeps the design relevant across decades. They are written for the person who must understand *why* the system is built the way it is, not merely *how* to operate it.

Each article is self-contained but references the others. SEC-004 (Air-Gap Architecture) is foundational -- the air gap is the institution's primary security boundary, and every other article in this batch assumes it is intact. SEC-005 (Cryptographic Survival) addresses the long-term viability of the cryptographic systems that protect data at rest and in transit across the sneakernet. SEC-006 (Physical Security Design) specifies the physical perimeter that makes the air gap meaningful. SEC-007 (Supply Chain Security) addresses the problem of trusting hardware and software that originate outside the institution's boundary. SEC-008 (Security Audit Framework) provides the methodology for verifying that all the other security systems are actually working.

If you are reading these articles decades after they were written, the specific technologies mentioned will have changed. The principles will not have. Use the principles to evaluate new technologies. Update the specifics. Preserve the reasoning.

---

---

# SEC-004 -- Air-Gap Architecture: Theory and Implementation

**Document ID:** SEC-004
**Domain:** 3 -- Security & Integrity
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, SEC-001, SEC-002, SEC-003
**Depended Upon By:** SEC-005, SEC-006, SEC-007, SEC-008, D18-001. All articles involving data transfer, external media, or boundary enforcement.

---

## 1. Purpose

This article is the complete technical specification of the air gap that defines the holm.chat Documentation Institution's primary security boundary. It defines what the air gap is in precise terms, where the boundary lies in physical and electromagnetic space, what constitutes a breach of that boundary, what controlled exceptions exist, and how the operator verifies that the gap remains intact over time.

SEC-001 established the air gap as an architectural decision. SEC-002 defined access controls that assume the air gap is intact. This article specifies the air gap itself -- its theory, its implementation, its verification, and its maintenance. It is the document you consult when you need to answer the question: "Is this institution actually air-gapped, or have we been deceiving ourselves?"

The air gap is not a single control. It is a composite boundary comprising physical isolation, electromagnetic containment, procedural discipline, and ongoing verification. A failure in any one of these dimensions can compromise the gap while the others appear intact. This article addresses all four dimensions with the rigor that the institution's primary security boundary demands.

## 2. Scope

This article covers the complete air-gap architecture:

- The theoretical definition of an air gap and its security properties.
- Physical isolation requirements: network interface removal, cable management, room layout.
- Electromagnetic considerations: emanations security (TEMPEST), RF isolation, power-line leakage.
- The USB and removable media policy: what devices may cross the boundary, under what conditions, using what procedures.
- The sneakernet protocol: how data moves into and out of the institution in a controlled manner.
- Verification procedures: how to confirm the air gap has not been breached.
- Exception documentation: every authorized exception to the air gap, why it exists, and the compensating controls around it.

This article does not cover the quarantine procedures for incoming media (see D18-001), the cryptographic protocols used to protect data in transit across the sneakernet (see SEC-005), or the physical perimeter that prevents unauthorized physical access (see SEC-006). It defines the electromagnetic and procedural boundary; those other articles define the complementary controls.

## 3. Background

### 3.1 What an Air Gap Actually Is

An air gap, in its simplest definition, is the absence of any electronic communication path between the institution's systems and any external network. No Ethernet cable. No WiFi connection. No Bluetooth pairing. No cellular modem. No infrared port. No acoustic coupling. No power-line networking. No communication path of any kind that could be exploited by an attacker who does not have physical access to the institution's hardware.

This definition is more demanding than it first appears. A system that has WiFi hardware installed but "disabled in software" is not air-gapped. A system that has an Ethernet port with no cable attached but with the network interface driver loaded is not air-gapped. A system that has no network hardware but sits next to a networked system on the same power circuit may not be fully air-gapped if power-line communication is possible. The air gap must be physical, verified, and maintained -- not assumed, configured, and forgotten.

### 3.2 The Security Properties of an Air Gap

A properly implemented air gap provides the following security properties:

**No remote exploitation.** An attacker cannot deliver exploits over a network connection that does not exist. This eliminates the largest category of attacks against modern systems: network-based intrusion, remote code execution, phishing, drive-by downloads, and man-in-the-middle attacks.

**No remote exfiltration.** Data cannot be extracted over a network connection that does not exist. Even if an attacker somehow compromises a system within the air gap (through supply chain attack or physical access), they cannot transmit the data out without physical access to the system or to a removable media device.

**No remote command and control.** Malware that relies on communication with a command-and-control server cannot function within an air-gapped environment. It may persist on the system, but it cannot receive instructions or exfiltrate data.

**Explicit consent for all data movement.** Every piece of data that enters or leaves the institution must be physically carried across the boundary by the operator. This means the operator has the opportunity -- and the obligation -- to inspect, approve, or reject every data transfer.

### 3.3 The Limitations of an Air Gap

An air gap is not omnipotent. It does not protect against:

- An attacker with physical access to the hardware.
- Supply chain attacks that pre-install malware on hardware or software before it enters the air-gapped environment.
- Insider threats (in this institution, the operator themselves acting against their own interests due to coercion or compromise).
- Electromagnetic emanations that leak data through radio frequency emissions, acoustic signals, thermal variations, or power consumption patterns.
- Social engineering that causes the operator to voluntarily carry compromised media across the gap.
- Data destruction attacks that do not require exfiltration.

These limitations are addressed by the complementary security articles: SEC-006 for physical security, SEC-007 for supply chain security, and the emanations controls specified in this article.

## 4. System Model

### 4.1 The Boundary Definition

The air-gap boundary is defined as a physical perimeter enclosing all institutional computing hardware, storage media, and cabling. Everything inside this perimeter is the "clean side." Everything outside is the "dirty side." The boundary is enforced at four levels:

**Level 1: Network Hardware Elimination.**
All wireless network interfaces (WiFi, Bluetooth, cellular, NFC, infrared) must be physically removed from institutional hardware, not merely disabled. For laptops with soldered wireless cards, the antenna leads must be physically disconnected and the wireless firmware must be removed from the operating system. For desktop systems, wireless cards must be physically removed. For any system where physical removal is impossible, the system must not be used as institutional hardware.

All wired network interfaces that are not used for internal-only institutional networking (such as a local network between institutional machines that never touches the outside world) must have their Ethernet ports physically blocked or their interface hardware removed. If internal networking is used between institutional machines, those interfaces and cables must be clearly labeled and the network must be physically incapable of reaching any non-institutional device.

**Level 2: Electromagnetic Containment.**
Institutional hardware should be positioned to minimize electromagnetic emanations that could be intercepted from outside the physical perimeter. This does not require a full TEMPEST-certified facility for most threat models, but it does require awareness of the basic principles:

- Institutional hardware should not be placed against exterior walls where emanations could be captured from outside the building.
- CRT monitors (if used) emit significantly stronger emanations than LCD panels and should be avoided or shielded.
- USB cables, HDMI cables, and power cables act as antennas. Cable runs should be minimized and cables should not exit the secure area unnecessarily.
- Unshielded cables running parallel to each other can create crosstalk. Separate data cables from power cables where practical.
- If the threat model includes a sophisticated adversary with emanations-capture capability, consider a Faraday enclosure for the most sensitive systems. For most single-operator institutions, placing hardware in an interior room with no shared walls to public spaces provides adequate protection.

**Level 3: Power Isolation.**
Power-line networking and power-line data leakage are real attack vectors. Institutional hardware should be powered through an uninterruptible power supply (UPS) that provides electrical isolation from the building's main wiring. The UPS acts as a buffer that breaks the direct electrical connection between institutional hardware and the external power grid. A UPS with online (double-conversion) topology provides the best isolation because all power passes through a battery-inverter stage, eliminating high-frequency signals that could carry data.

If off-grid power is used (solar, generator), the isolation is inherent -- the institution's power supply is physically separate from any external grid.

**Level 4: Procedural Enforcement.**
The physical and electromagnetic controls above prevent electronic breaches. Procedural controls prevent human-assisted breaches. The following procedures constitute the air gap's procedural layer:

- No personal electronic devices (phones, smartwatches, wireless earbuds) may be brought within the physical perimeter during operations that involve sensitive data. These devices contain wireless radios that could be exploited as exfiltration channels.
- No USB device may be connected to institutional hardware without completing the sneakernet protocol defined in Section 4.3.
- No hardware may be added to the institutional environment without completing the supply chain verification defined in SEC-007.
- The operator must log every instance of a device crossing the air-gap boundary, whether inbound or outbound.

### 4.2 The USB and Removable Media Policy

USB devices are the primary authorized channel for data to cross the air-gap boundary. They are also the primary attack vector against air-gapped systems. The USB policy must balance these realities.

**Authorized USB Devices:**
The institution maintains a registry of authorized USB storage devices. Each device is assigned a unique identifier (physically labeled on the device), a classification (inbound-only, outbound-only, or bidirectional), and a custodian. No USB storage device may be used with institutional hardware unless it appears in the registry.

**Device Classes:**
- *Inbound-only devices* carry data into the institution. They are used to import software updates, reference data, and other material from external sources. After import, data on these devices passes through the quarantine process defined in D18-001 before touching any institutional system.
- *Outbound-only devices* carry data out of the institution. Their use cases are limited: transferring backup verification data to a geographically separate backup site, or exporting data that the operator has deliberately chosen to share with the external world.
- *Bidirectional devices* may carry data in both directions. Their use should be minimized because they represent the highest risk: a device that has touched the external world and then touches institutional hardware could carry contamination in either direction.

**USB Port Control:**
Institutional systems should have only the minimum number of USB ports enabled. On Linux systems, USBGuard or equivalent tooling should be configured to whitelist only authorized device identifiers and block all others. Physical USB port blockers (plastic inserts that prevent unauthorized insertion) should be installed on all unused ports.

**Write Protection:**
When transferring data out of the institution, use USB devices with physical write-protect switches. Enable write protection before connecting the device to any non-institutional system. This prevents the device from being contaminated during external use.

### 4.3 The Sneakernet Protocol

The sneakernet is the institution's controlled process for moving data across the air-gap boundary. It is called a sneakernet because data is physically carried ("by sneaker") rather than transmitted electronically.

**Inbound Transfer (External to Institutional):**

1. The operator identifies data that needs to be imported (software package, reference document, firmware update).
2. On an external system (not institutional hardware), the operator downloads the data and records its cryptographic hash (SHA-256 minimum).
3. The operator transfers the data to an inbound-only USB device.
4. The operator physically carries the USB device to the quarantine station (a designated institutional machine used only for incoming media processing, as defined in D18-001).
5. The quarantine station verifies the cryptographic hash against the recorded value.
6. The quarantine station performs malware scanning and integrity checks.
7. If the data passes quarantine, it is transferred to the institutional network. If it fails, it is logged and rejected.
8. The USB device is securely wiped after use.

**Outbound Transfer (Institutional to External):**

1. The operator identifies data that needs to be exported and documents the justification.
2. The data is copied to an outbound-only USB device with write protection available.
3. The operator enables write protection on the USB device before removing it from institutional hardware.
4. The transfer is logged, including: date, data description, destination, justification, device identifier.
5. The USB device is securely wiped after external use, before any future institutional use.

**Transfer Logging:**
Every sneakernet transfer, inbound or outbound, is recorded in the air-gap transfer log. This log is append-only and includes: date, direction (inbound/outbound), device identifier, data description, cryptographic hash, quarantine result (for inbound), justification (for outbound), and operator confirmation.

### 4.4 Exception Registry

No air-gap policy survives contact with reality without exceptions. The critical requirement is that every exception is documented, justified, time-bounded, and subject to compensating controls.

The exception registry is a formal document that lists every authorized deviation from the air-gap architecture. Each entry includes:

- **Exception ID:** Sequential identifier.
- **Description:** What the exception permits.
- **Justification:** Why the exception is necessary.
- **Duration:** When the exception expires (all exceptions must be time-bounded).
- **Compensating Controls:** What additional security measures are in place to mitigate the risk.
- **Review Date:** When the exception will be reviewed for continued necessity.
- **Approval:** Tier 2 decision per GOV-001 (architectural decision with 30-day waiting period).

As of the institution's founding, the exception registry should be empty. If it is not empty, every entry deserves scrutiny.

## 5. Rules & Constraints

- **R-SEC04-01:** No wireless network interface may exist in an enabled state on any institutional hardware. Physical removal is required where possible; physical disconnection of antenna leads where removal is impossible.
- **R-SEC04-02:** All USB devices used with institutional hardware must be registered in the authorized device registry. Unregistered devices must be rejected by both procedural and technical controls (USBGuard or equivalent).
- **R-SEC04-03:** Every data transfer across the air-gap boundary must be logged in the air-gap transfer log with full metadata as specified in Section 4.3.
- **R-SEC04-04:** Exceptions to the air-gap architecture require Tier 2 governance approval (30-day waiting period) and must be recorded in the exception registry.
- **R-SEC04-05:** The air-gap verification procedure (Section 6 recovery, and ongoing verification) must be executed at least quarterly and after any hardware change.
- **R-SEC04-06:** No personal wireless device may be present within the physical perimeter during sensitive operations.
- **R-SEC04-07:** Power to institutional hardware must pass through an isolating UPS or be sourced from an independent off-grid supply.

## 6. Failure Modes

- **Creeping connectivity.** A device with an overlooked wireless capability is introduced. A USB device with hidden wireless hardware (e.g., a USB device with an embedded cellular modem) is connected. A new piece of hardware includes a network interface the operator did not know about. Mitigation: hardware audit at acquisition (SEC-007), quarterly air-gap verification, USBGuard whitelist enforcement.
- **Procedural decay.** The sneakernet protocol is followed loosely. Transfers are not logged. Quarantine is skipped "because the source is trusted." The protocol simplifies itself through laziness until it provides no security at all. Mitigation: the air-gap transfer log provides an auditable record. The security audit (SEC-008) checks for procedural adherence.
- **Emanations leakage.** A change in hardware placement, cable routing, or room configuration creates an emanations path that did not previously exist. Mitigation: reassess emanations posture after any physical change to the environment. The physical security review (SEC-006) includes emanations considerations.
- **USB contamination.** A USB device used for bidirectional transfer carries malware from the external environment into the institution, bypassing quarantine due to a zero-day exploit or a scanning tool that is not up to date. Mitigation: minimize bidirectional device use, maintain quarantine scanning tools, use defense in depth (even if malware enters, least-privilege access controls limit its impact).
- **False confidence.** The operator believes the air gap is intact because it was configured correctly at installation, without recognizing that the gap requires ongoing verification. A firmware update re-enables a wireless interface. A hardware replacement introduces a new network capability. Mitigation: quarterly verification, post-change verification, the mindset that the air gap must be continuously proven, not assumed.

## 7. Recovery Procedures

1. **If the air gap has been breached:** Immediately disconnect the breaching pathway (remove the network cable, power down the wireless device, isolate the compromised system). Assess the duration and scope of the breach. If the breach was momentary and no data transfer occurred, log the incident, restore the gap, and reinforce the control that failed. If data transfer may have occurred, treat the event as a full security incident per D10-003: isolate affected systems, perform forensic analysis, assess what data was exposed, and rebuild from known-good backups if contamination is suspected.
2. **If an unauthorized USB device was connected:** Remove the device immediately. Assess what the device may have introduced or extracted. If the device was connected for more than a few seconds, assume it may have delivered a payload. Scan the system. If contamination is suspected, rebuild from backup. Log the incident. Review USBGuard configuration to determine why the device was not blocked.
3. **If emanations leakage is suspected:** Assess the environment for changes that may have created the leakage path. Reposition hardware, reroute cables, or install shielding as appropriate. If sensitive data may have been captured through emanations, assess the impact based on the sensitivity of the data and the sophistication of the likely adversary. For most personal institutions, emanations capture requires a sophisticated and specifically targeted attacker; adjust the response to the realistic threat level.
4. **If the air-gap transfer log has been neglected:** Reconstruct what transfers can be reconstructed from memory and from examining the USB device registry and quarantine logs. Mark reconstructed entries as "RECONSTRUCTED." Resume disciplined logging immediately. Schedule a full security audit (SEC-008) within 30 days.

## 8. Evolution Path

- **Years 0-5:** The air gap is being established for the first time. The sneakernet protocol will feel cumbersome. Resist the temptation to simplify it prematurely. Establish the logging discipline. Build the authorized device registry. Run quarterly verifications. The habits formed now define the air gap's long-term integrity.
- **Years 5-15:** The air gap should be routine. The primary risk shifts from implementation errors to procedural decay. Focus on audit adherence and on updating emanations posture as hardware is replaced.
- **Years 15-30:** Hardware generations will change substantially. New connectivity technologies (successors to WiFi, Bluetooth, and cellular that we cannot predict) will need to be evaluated and excluded. The principle remains the same: if it can communicate, it must be removed or disabled physically.
- **Years 30-50+:** The air gap as a concept may be challenged by new paradigms. Technologies that embed networking so deeply into hardware that removal is impractical may emerge. The institution must evaluate whether the air gap can be maintained with available hardware, or whether compensating controls (such as Faraday enclosures) must become primary rather than supplementary.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The air gap is the most conspicuous security decision this institution makes, and also the most frequently misunderstood. People hear "air-gapped" and think it means "not connected to the internet." It means far more than that. It means no electronic communication path of any kind -- not through the air, not through the wires, not through the power lines, not through the USB port without explicit, logged, quarantined consent. It is a total communication boundary, not merely an internet disconnection. The day we start treating the air gap casually is the day it stops protecting us. I have written this article to make "casual" impossible, by specifying every dimension of the gap in enough detail that half-measures are recognizable as such.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 1: Sovereignty; Principle 2: Integrity Over Convenience)
- CON-001 -- The Founding Mandate (air-gap mandate, Section 3.3)
- SEC-001 -- Threat Model and Security Philosophy (air-gap decision, Section 3.2; trust model, Section 4.3)
- SEC-002 -- Access Control Procedures (physical access as first layer)
- SEC-003 -- Key Management (cryptographic controls for data crossing the boundary)
- SEC-006 -- Physical Security Design (physical perimeter that enforces the air gap)
- SEC-007 -- Supply Chain Security (hardware verification before entry into air-gapped environment)
- D18-001 -- Import & Quarantine Procedures (quarantine process for inbound media)
- NIST SP 800-123 -- Guide to General Server Security (baseline hardening reference)
- Stage 1 Documentation Framework, Domain 3: Security & Integrity

---

---

# SEC-005 -- Long-Term Cryptographic Survival

**Document ID:** SEC-005
**Domain:** 3 -- Security & Integrity
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, SEC-001, SEC-003, SEC-004
**Depended Upon By:** SEC-008, D6-002, D12-001. All articles involving encryption, hashing, digital signatures, or long-term data integrity.

---

## 1. Purpose

This article addresses a problem that most security documents ignore because their time horizon is measured in years, not decades: how does the institution maintain cryptographic security over a fifty-year operational lifetime?

Cryptographic algorithms are not permanent. They are mathematical constructions whose security depends on computational assumptions -- assumptions about how hard certain problems are to solve. These assumptions erode over time as computational power increases, as mathematical breakthroughs occur, and as new computing paradigms (most notably quantum computing) emerge. An encryption algorithm that is unbreakable today may be trivially broken in twenty years. A hash function that guarantees integrity today may be vulnerable to collisions in fifteen.

The institution's data must remain confidential and its integrity must remain verifiable across the full fifty-year horizon. This requires not merely selecting strong algorithms today, but building a cryptographic architecture that can survive algorithm obsolescence -- a property known as crypto-agility. This article specifies that architecture: what it means, how it works, and how the operator maintains it across decades of cryptographic evolution.

## 2. Scope

This article covers:

- The theory of cryptographic longevity and why algorithms become obsolete.
- The quantum computing threat: what it is, what it affects, and what it does not.
- The crypto-agility architecture: how the institution's systems are designed to survive algorithm transitions.
- The cryptographic watchlist: how the operator monitors the health of deployed algorithms.
- Migration procedures: how to re-encrypt, re-hash, and re-sign data when an algorithm is broken or deprecated.
- Algorithm selection criteria for long-term deployments.
- The relationship between key management (SEC-003) and cryptographic survival.

This article does not cover the specific key management procedures (see SEC-003), the specific encryption configurations for particular systems (see domain-specific articles), or the operational procedures for daily cryptographic operations (see Stage 3 articles). It provides the strategic framework within which all cryptographic decisions are made.

## 3. Background

### 3.1 Why Algorithms Die

Cryptographic algorithms become obsolete through three mechanisms:

**Analytical attacks.** Mathematicians discover weaknesses in the algorithm's structure that allow it to be broken faster than brute force. This happened to MD5 (collision attacks demonstrated in 2004), SHA-1 (theoretical weaknesses published in 2005, practical collisions demonstrated in 2017), and DES (effective key length recognized as insufficient by the 1990s). Analytical attacks are unpredictable. An algorithm that appears secure today could be broken by a paper published tomorrow.

**Computational advances.** Increases in available computing power reduce the effective security margin of algorithms. A 1024-bit RSA key was considered secure in 2000; by 2020, it was considered inadequate not because the algorithm was broken but because the computational cost of factoring 1024-bit numbers had decreased to within range of well-funded adversaries. This mechanism is gradual and somewhat predictable, but the rate of advance can spike with architectural innovations (GPUs, ASICs, distributed computing).

**Paradigm shifts.** A fundamentally new computing paradigm renders the mathematical assumptions underlying an algorithm invalid. The most significant current example is quantum computing. Shor's algorithm, when executed on a sufficiently powerful quantum computer, can factor large integers and compute discrete logarithms in polynomial time, which breaks RSA, DSA, Diffie-Hellman, and elliptic curve cryptography. Grover's algorithm provides a quadratic speedup for brute-force search, effectively halving the key length of symmetric ciphers (AES-256 becomes equivalent to AES-128 against a quantum adversary).

### 3.2 The Fifty-Year Problem

Over fifty years, the institution can expect:

- At least two major hash function transitions (comparable to MD5 to SHA-1 to SHA-2 to SHA-3).
- At least one major asymmetric cryptography transition (comparable to the anticipated transition from classical to post-quantum algorithms).
- Multiple increases in recommended key lengths for symmetric ciphers.
- The possible emergence of computing paradigms beyond quantum that are not yet imagined.

Data encrypted today must either remain secure for fifty years under its current encryption, or the institution must be capable of re-encrypting it under a stronger algorithm before the current one is broken. Integrity hashes computed today must either remain collision-resistant for fifty years, or the institution must be capable of re-hashing data under a stronger function before the current one is compromised.

This is not a theoretical concern. It is a design requirement.

### 3.3 The Harvest-Now, Decrypt-Later Threat

Even in an air-gapped institution, the harvest-now-decrypt-later threat is relevant. If encrypted backups are stored at geographically separate locations (as they should be for disaster recovery), those backups are physically outside the air-gap perimeter. An adversary who obtains a backup today cannot decrypt it today -- but they can store it and wait for the encryption algorithm to be broken. If the institution uses RSA-2048 for backup encryption and a practical quantum computer is built in 2045, every backup encrypted before the migration to post-quantum cryptography is retroactively compromised.

This is why crypto-agility is not optional. The cost of re-encrypting data proactively is small. The cost of discovering that decades of backups are now readable by an adversary is catastrophic.

## 4. System Model

### 4.1 The Crypto-Agility Architecture

Crypto-agility means that the institution's systems are designed so that the cryptographic algorithms they use can be replaced without redesigning the systems themselves. This requires four architectural properties:

**Property 1: Algorithm Abstraction.**
No system should hardcode a specific cryptographic algorithm. All systems should reference algorithms through an abstraction layer -- a configuration file, a library interface, or a protocol parameter -- that allows the algorithm to be changed without modifying the system's code or architecture. For example, full-disk encryption should use LUKS (Linux Unified Key Setup), which supports multiple cipher suites and allows re-encryption with a different cipher without decrypting the entire volume first.

**Property 2: Metadata Completeness.**
Every encrypted file, every hashed value, and every digital signature must include metadata specifying which algorithm was used, which key (by identifier, not by value) was used, and when the operation was performed. Without this metadata, the institution cannot determine which data needs to be migrated when an algorithm is deprecated. The metadata format should be standardized across all institutional systems.

**Property 3: Layered Encryption.**
For the most sensitive data, the institution may employ layered encryption: encrypting data with one algorithm and then encrypting the result with a second, independent algorithm. If either algorithm is broken, the data remains protected by the other. This is computationally expensive and should be reserved for data whose confidentiality must be maintained for the full fifty-year horizon regardless of cost.

**Property 4: Re-encryption Capability.**
The institution must maintain the operational capability to re-encrypt, re-hash, and re-sign its entire data corpus within a reasonable timeframe. This means maintaining tools, procedures, and sufficient computational resources to perform a full cryptographic migration. The migration procedure must be tested periodically (see Section 4.4) so that it works when needed, not merely in theory.

### 4.2 Algorithm Selection Criteria

When selecting cryptographic algorithms for institutional use, the following criteria apply:

- **Standardization:** Prefer algorithms that have been standardized by recognized bodies (NIST, ISO, IETF). Standardization reflects broad expert review.
- **Maturity:** Prefer algorithms that have been publicly analyzed for at least ten years without practical attacks being found. Novelty is a risk factor, not an advantage, in cryptography.
- **Conservative parameters:** Select key lengths and parameters at or above the highest recommended level. Use AES-256 rather than AES-128. Use SHA-512 rather than SHA-256 when performance permits. The marginal cost of larger parameters is small; the security margin is substantial.
- **Post-quantum readiness:** For any new deployment, evaluate whether a post-quantum alternative exists and whether it is mature enough for institutional use. As of this writing, NIST has standardized ML-KEM (formerly CRYSTALS-Kyber) for key encapsulation and ML-DSA (formerly CRYSTALS-Dilithium) for digital signatures. These should be evaluated for institutional adoption.
- **Implementation quality:** An algorithm is only as strong as its implementation. Prefer well-audited, open-source implementations. Avoid implementing cryptographic algorithms from scratch.

### 4.3 The Cryptographic Watchlist

The cryptographic watchlist is a maintained document that tracks the health status of every cryptographic algorithm in use within the institution. For each algorithm, the watchlist records:

- **Algorithm name and version.**
- **Institutional use cases:** Where in the institution this algorithm is deployed.
- **Current status:** One of four levels:
  - *Green:* No known weaknesses. Considered secure for the foreseeable future.
  - *Yellow:* Theoretical weaknesses published or computational margins shrinking. Begin planning migration.
  - *Orange:* Practical attacks demonstrated in academic settings or algorithm formally deprecated by standards body. Execute migration within 12 months.
  - *Red:* Practical attacks demonstrated in real-world conditions. Execute emergency migration immediately.
- **Last review date:** When the algorithm's status was last assessed.
- **Migration target:** The algorithm that will replace this one if migration is needed.

The watchlist is reviewed at least annually during the annual security review (SEC-008). The review involves consulting the current state of cryptographic research, which requires bringing information across the air gap via the sneakernet protocol. This is one of the essential inbound data transfers that justifies the sneakernet's existence.

### 4.4 Migration Procedures

When an algorithm moves to Orange or Red status on the watchlist, the institution executes a cryptographic migration:

**Phase 1: Inventory.**
Identify every system, file, key, hash, and signature that uses the deprecated algorithm. The metadata completeness property (Section 4.1, Property 2) makes this possible. Without complete metadata, this phase becomes an archaeological excavation through the entire data corpus.

**Phase 2: Priority Assessment.**
Rank the inventory by sensitivity and exposure. Data stored outside the air gap (offsite backups) is the highest priority because it is most exposed to the harvest-now-decrypt-later threat. Data stored inside the air gap is lower priority because physical access is required to exploit any weakness.

**Phase 3: Migration Execution.**
Re-encrypt, re-hash, or re-sign data using the migration target algorithm. For encrypted data, this means decrypting with the old algorithm and re-encrypting with the new one. For hashed data, this means recomputing hashes with the new function. For signatures, this means re-signing with the new algorithm.

**Phase 4: Verification.**
Verify that all migrated data is accessible and intact under the new algorithm. Test decryption. Test hash verification. Test signature validation.

**Phase 5: Cleanup.**
Once migration is verified, securely delete data encrypted solely under the deprecated algorithm (if copies under the new algorithm exist and are verified). Update the watchlist. Update cryptographic metadata across all institutional systems.

**Phase 6: Documentation.**
Record the migration in the institutional decision log (Tier 2 decision). Update all articles that reference the deprecated algorithm. Add a commentary entry to this article documenting the migration and the lessons learned.

## 5. Rules & Constraints

- **R-SEC05-01:** Every encrypted file, hash value, and digital signature in the institution must include metadata identifying the algorithm, key identifier, and date of operation.
- **R-SEC05-02:** No cryptographic algorithm may be deployed in the institution without first being evaluated against the selection criteria in Section 4.2 and recorded on the cryptographic watchlist.
- **R-SEC05-03:** The cryptographic watchlist must be reviewed at least annually. The review must include consultation of current cryptographic research obtained via the sneakernet.
- **R-SEC05-04:** When an algorithm reaches Orange status, migration planning must begin within 30 days and migration execution must complete within 12 months.
- **R-SEC05-05:** When an algorithm reaches Red status, migration execution must begin immediately and complete as fast as is operationally feasible.
- **R-SEC05-06:** The migration procedure must be tested at least once every three years on a representative data sample, even when no migration is pending. An untested migration procedure is not a procedure.
- **R-SEC05-07:** Post-quantum cryptographic readiness must be evaluated during each annual review. When mature post-quantum algorithms are available, migration planning must begin.
- **R-SEC05-08:** Layered encryption must be used for any data whose confidentiality must be maintained for the full institutional lifetime and that is stored outside the air-gap perimeter.

## 6. Failure Modes

- **Watchlist neglect.** The cryptographic watchlist is not reviewed. An algorithm degrades from Green to Red without the institution noticing. Data encrypted under the broken algorithm is exposed. Mitigation: the annual security review (SEC-008) includes a mandatory watchlist review. The review is a checklist item, not an optional addendum.
- **Metadata absence.** Cryptographic operations are performed without recording metadata. When migration is needed, the institution cannot determine which data uses the deprecated algorithm. Migration becomes guesswork. Mitigation: R-SEC05-01 makes metadata mandatory. Automated tooling should enforce metadata recording.
- **Migration paralysis.** The migration procedure has never been tested. When a real migration is needed, unforeseen complications delay it beyond the window of safety. Mitigation: R-SEC05-06 requires periodic testing of the migration procedure itself.
- **Algorithm monoculture.** The institution uses a single algorithm for all purposes. When that algorithm is broken, everything is compromised simultaneously. Mitigation: use different algorithms for different purposes where practical (e.g., AES for bulk encryption, a different algorithm family for key encapsulation), so that a break in one does not compromise all.
- **Quantum surprise.** A practical quantum computer is developed sooner than anticipated, breaking all classical asymmetric cryptography before the institution has migrated to post-quantum alternatives. Mitigation: begin post-quantum migration planning now, even though the timeline is uncertain. Use layered encryption for the most sensitive data. Prioritize symmetric algorithms (which are quantum-resistant at doubled key lengths) over asymmetric algorithms for long-term protection.

## 7. Recovery Procedures

1. **If an algorithm is broken unexpectedly:** Treat it as a Red-status event. Execute the migration procedure immediately. Assess what data was encrypted solely under the broken algorithm and is stored outside the air gap. That data should be considered potentially compromised. Assess the impact. Migrate all remaining data as rapidly as possible. Document the event.
2. **If cryptographic metadata is missing:** Develop tooling to scan the data corpus and identify algorithm usage by file headers, magic numbers, and structural analysis. This is labor-intensive and error-prone. Once identified, add the missing metadata. Treat this as a lesson in why metadata discipline matters. Add a commentary entry documenting the remediation.
3. **If the migration procedure fails mid-execution:** Do not panic. The pre-migration data still exists (migrations should never destroy the source until the destination is verified). Diagnose the failure. Fix the procedure. Resume from the point of failure. This is why Phase 4 (verification) exists as a separate step: it catches problems before Phase 5 (cleanup) makes them irreversible.
4. **If post-quantum migration is needed urgently:** Prioritize data stored outside the air gap. Use hybrid encryption (classical + post-quantum) as a transitional measure if full post-quantum migration is not immediately feasible. Hybrid encryption ensures that data is protected if either algorithm is secure, which is the safest posture during a transition of uncertain duration.

## 8. Evolution Path

- **Years 0-5:** Establish the cryptographic watchlist. Deploy algorithms with conservative parameters. Implement metadata recording across all systems. Test the migration procedure on a representative data sample. Evaluate post-quantum algorithms for maturity and institutional suitability.
- **Years 5-15:** Expect at least one watchlist status change requiring attention. Execute the first real migration if needed. By this period, post-quantum algorithms should be mature enough for institutional deployment. Plan and execute the post-quantum migration if not already done.
- **Years 15-30:** The post-quantum migration should be complete. The watchlist continues to evolve. New algorithm families may emerge. The crypto-agility architecture proves its value as the institution navigates its first full algorithm lifecycle.
- **Years 30-50+:** The cryptographic landscape will have changed in ways we cannot predict. The crypto-agility architecture and the watchlist discipline are the institution's defenses against that unpredictability. They do not require predicting the future. They require responding to the future as it arrives.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
Cryptography is the only security domain where the threat travels backward in time. Data encrypted today can be captured today and decrypted in the future when the algorithm breaks. This means that every cryptographic decision we make today has consequences that extend to the end of the institution's lifetime. It also means that procrastination is uniquely dangerous in this domain: by the time you realize you should have migrated, the window for safe migration may have already closed. I have designed this article to make procrastination structurally difficult, through the watchlist, the mandatory reviews, and the periodic migration testing. Future maintainers: take the watchlist seriously. The algorithms will change. The only question is whether you change with them proactively or reactively.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 2: Integrity Over Convenience)
- CON-001 -- The Founding Mandate (long-term data preservation mandate)
- SEC-001 -- Threat Model and Security Philosophy (threat categories 2 and 3)
- SEC-003 -- Key Management (key lifecycle and storage)
- SEC-004 -- Air-Gap Architecture (sneakernet protocol for watchlist updates)
- SEC-008 -- Security Audit Framework (annual cryptographic review)
- NIST SP 800-57 -- Recommendation for Key Management (algorithm lifecycle guidance)
- NIST Post-Quantum Cryptography Standardization (ML-KEM, ML-DSA standards)
- Domain 6 -- Data Architecture (data format and metadata standards)
- Domain 12 -- Disaster Recovery (backup encryption requirements)
- Stage 1 Documentation Framework, Domain 3: Security & Integrity

---

---

# SEC-006 -- Physical Security Design

**Document ID:** SEC-006
**Domain:** 3 -- Security & Integrity
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, SEC-001, SEC-002, SEC-004
**Depended Upon By:** SEC-007, SEC-008, D4-001. All articles involving site design, hardware placement, or environmental controls.

---

## 1. Purpose

This article specifies the physical security design for the site that houses the holm.chat Documentation Institution. In an air-gapped institution, physical security is the perimeter. There is no network firewall because there is no network. The walls, the locks, the doors, and the environmental controls are the firewall. If physical security fails, every other security control can be bypassed.

This article addresses the full spectrum of physical security: perimeter design, access control hardware, environmental threat mitigation (fire, flood, intrusion, and temperature), surveillance, tamper detection, and the audit procedures that verify these controls remain effective over time. It is written for the operator who must design, install, maintain, and periodically verify the physical security of a site that must endure for decades.

Physical security is the domain where the institution most directly confronts the physical world. Algorithms do not rust. Documentation does not flood. But the hardware that stores the institution's data exists in a physical environment that can burn, flood, freeze, shake, or be broken into. This article is the institution's defense against those realities.

## 2. Scope

This article covers:

- Site selection criteria and layout principles.
- Perimeter design: walls, doors, windows, and barriers.
- Access control hardware: locks, keys, and authentication mechanisms.
- Environmental threat mitigation: fire suppression, water protection, temperature and humidity control.
- Tamper detection: seals, sensors, and indicators.
- Surveillance: monitoring capabilities compatible with a single-operator model.
- Physical security audit procedures.

This article does not cover logical access controls (see SEC-002), cryptographic protections (see SEC-003 and SEC-005), or the air-gap architecture beyond its physical enforcement dimension (see SEC-004). It defines the physical environment; those other articles define the controls that operate within it.

## 3. Background

### 3.1 The Physical Security Imperative

The founding mandate (CON-001) declares that the institution is self-sovereign: it depends on no external entity. The air-gap architecture (SEC-004) eliminates electronic attack vectors. These decisions concentrate the institution's vulnerability surface onto the physical domain. An attacker who wishes to compromise an air-gapped institution must achieve physical access. A natural disaster that threatens an air-gapped institution threatens its hardware directly, with no possibility of cloud failover or remote backup retrieval.

This concentration of risk onto the physical domain is a deliberate design choice. It trades distributed, hard-to-manage electronic risks for concentrated, manageable physical risks. But the word "manageable" carries obligation: the physical risks must actually be managed. A well-designed physical security posture converts the air gap from a limitation into a strength. A poorly designed one makes the air gap irrelevant, because the physical boundary it depends on does not hold.

### 3.2 Threat Landscape for a Personal Institution

The physical threat landscape for a personal institution differs from an enterprise data center. Enterprises face sophisticated adversaries with significant resources. A personal institution is more likely to face:

- **Opportunistic theft:** Burglars who target electronics for resale value, not for data.
- **Environmental events:** Fire, flooding, severe weather, and power surges that are local, not targeted.
- **Accidental damage:** The operator or a household member inadvertently damaging hardware through spills, drops, or mishandling.
- **Gradual degradation:** Temperature extremes, humidity, dust, and pests that degrade hardware over years.
- **Targeted intrusion:** Less likely but nonzero, especially over a fifty-year horizon as the institution accumulates value.

The physical security design must be proportionate to these threats. It does not require a bank vault. It requires deliberate, layered protection that addresses each threat category with appropriate controls.

## 4. System Model

### 4.1 The Concentric Perimeter Model

Physical security is organized as concentric perimeters, each providing an additional layer of protection:

**Perimeter 1: The Property Boundary.**
The outermost layer. For a residential site, this is the property line. Controls at this layer include fencing, exterior lighting, visibility from public areas (which deters some attackers and attracts others; assess based on local threat model), and landscaping that does not provide concealment near the building. This layer provides deterrence and detection, not prevention.

**Perimeter 2: The Building Envelope.**
The building itself. Controls at this layer include reinforced exterior doors, window locks or bars (appropriate to local building codes and aesthetics), deadbolt locks, and door frame reinforcement. This layer prevents casual intrusion and delays determined intrusion.

**Perimeter 3: The Secure Room.**
The room or area within the building that houses institutional hardware. This is the critical perimeter. Controls at this layer include a solid-core or metal door with a high-security lock, no windows (or windows that are permanently sealed and reinforced), walls that extend fully to the structural ceiling (not merely to a drop ceiling that can be bypassed), and environmental controls. This room is where the air-gap perimeter (SEC-004) is physically enforced.

**Perimeter 4: The Hardware Enclosures.**
The racks, cabinets, or enclosures that contain institutional hardware. Controls at this layer include lockable server racks or cabinets, tamper-evident seals on hardware cases, and cable management that prevents accidental disconnection.

Each perimeter should provide independent protection. An attacker who bypasses one perimeter should face a distinct challenge at the next. This is defense in depth applied to physical space.

### 4.2 Access Control Hardware

**Lock Selection:**
The secure room should use a high-security lock with the following properties:
- Pick resistance: rated for a minimum of 10 minutes of skilled attack (look for UL 437 or equivalent certification).
- Bump resistance: the lock should be resistant to lock bumping attacks.
- Key control: the key blanks should not be freely available. The lock manufacturer should restrict key blank distribution to authorized dealers.
- Rekey capability: the lock should be rekeyable without replacement, to support key changes during succession or after a suspected compromise.

**Key Management (Physical):**
Physical keys are managed with the same discipline as cryptographic keys:
- A master key is stored in a secure, fire-rated safe at the site.
- A duplicate key is stored at a geographically separate location, accessible to the designated successor (per the succession protocol in GOV-001).
- A key log records every key issued, to whom, and when. Key issuance is a logged event.
- If a key is lost or suspected compromised, the lock is rekeyed within 24 hours.

**Electronic Access Control (Optional):**
Electronic access control (keypads, card readers, biometric scanners) may be used as a supplementary layer. However, electronic access control must not be the sole access mechanism. Electronic systems fail -- batteries die, circuits corrode, software corrupts. A mechanical lock must always provide a backup access path. Electronic access control's primary value is auditability: it can log entry times and identities automatically, which a mechanical lock cannot.

If electronic access control is deployed, it must be powered by the institution's UPS, not by external power alone. It must fail secure (locked) rather than fail open. Its audit logs must be integrated into the institution's operational logging.

### 4.3 Environmental Threat Mitigation

**Fire Protection:**
- The secure room should contain a fire suppression system appropriate for electronics: clean-agent (gas-based) suppression rather than water sprinklers. If clean-agent suppression is not feasible, a handheld clean-agent fire extinguisher should be immediately accessible within the room.
- Smoke detection should be installed inside the secure room and in adjacent areas. Smoke detectors should be tested monthly as part of the operational tempo.
- Fire-resistant storage (a fire-rated safe or fire-resistant cabinet) should house the most critical backup media and physical key duplicates. A one-hour fire rating minimum is recommended.
- Hardware should be positioned away from electrical panels, water heaters, and other ignition sources.

**Water Protection:**
- The secure room should not be located below grade (basements flood) if any alternative exists. If below-grade placement is unavoidable, install a sump pump with battery backup and a water sensor that provides audible alarm.
- Hardware should be elevated at least 15 centimeters (6 inches) above the floor on racks or platforms.
- No plumbing should run through or above the secure room if the building's layout permits this.
- A water sensor should be placed at floor level in the secure room. The sensor should produce an audible alarm.

**Temperature and Humidity Control:**
- The secure room should maintain a temperature between 18-27 degrees Celsius (64-80 degrees Fahrenheit) and relative humidity between 40-60 percent.
- A temperature and humidity logger should run continuously, with data reviewed during monthly operations.
- Hardware generates heat. Ensure adequate ventilation or cooling for the thermal load of the installed equipment.
- Storage media (especially magnetic media and optical discs) are sensitive to humidity extremes. Long-term media storage should be in the most climate-controlled area available.

**Power Protection:**
- All institutional hardware must be connected through an uninterruptible power supply (UPS) as specified in SEC-004.
- The UPS should provide sufficient runtime for a graceful shutdown of all systems (minimum 15 minutes recommended).
- Surge protection must be present on all power connections.
- If off-grid power is available (solar, generator), automatic or manual transfer capability should exist to maintain power during grid outages.

### 4.4 Tamper Detection

Tamper detection provides evidence that an unauthorized person has accessed or modified institutional hardware.

**Tamper-Evident Seals:**
- All hardware cases should be sealed with serialized tamper-evident seals (tape or stickers that show visible evidence of removal).
- The serial numbers of all seals should be recorded in the tamper-detection log.
- Seals should be checked during monthly operations. Any broken or missing seal triggers an investigation.
- When hardware is legitimately opened for maintenance, the old seal is destroyed, the opening is logged, and a new seal is applied.

**Position Indicators:**
- Mark the exact position of hardware on racks or shelves with tape or markers. Displacement from the marked position indicates that hardware has been moved.
- Record the positions photographically during quarterly security reviews.

**Intrusion Indicators:**
- The secure room door should have a tamper indicator (a simple seal on the door frame, or a more sophisticated sensor) that reveals whether the door has been opened.
- If the institution is left unattended for extended periods (travel, absence), the tamper indicators become the primary means of detecting unauthorized access during the absence.

### 4.5 Surveillance

Surveillance in a single-operator institution must be practical and privacy-respecting:

- If cameras are used, they should monitor the exterior of the building and the approach to the secure room, not the interior of the secure room itself (to prevent surveillance of sensitive data on screens).
- Camera storage should be local (within the air-gapped environment or on dedicated, non-networked storage), not cloud-based.
- Camera systems should not have WiFi or network capability. If network-capable cameras are the only option available, their network interfaces must be physically disabled per SEC-004.
- Motion-activated lighting on the building exterior is a low-cost, high-value deterrent that does not create privacy concerns.
- An audible alarm system (door contacts, glass-break sensors, motion detectors) provides intrusion detection during periods when the operator is absent.

## 5. Rules & Constraints

- **R-SEC06-01:** Institutional hardware must be housed in a dedicated secure room that implements at minimum Perimeters 3 and 4 as described in Section 4.1.
- **R-SEC06-02:** The secure room must use a high-security lock with pick, bump, and key-control resistance. Mechanical backup access must exist independent of any electronic access control.
- **R-SEC06-03:** Fire detection must be present in the secure room. Fire suppression appropriate for electronics must be accessible within the room.
- **R-SEC06-04:** Hardware must be elevated above floor level. Water detection must be present at floor level in the secure room.
- **R-SEC06-05:** All hardware cases must be sealed with serialized tamper-evident seals. Seal integrity must be checked monthly.
- **R-SEC06-06:** Temperature and humidity must be monitored continuously in the secure room. Readings outside the acceptable range trigger an investigation and corrective action within 24 hours.
- **R-SEC06-07:** Physical keys must be managed per the key management procedure in Section 4.2, including off-site backup and a key issuance log.
- **R-SEC06-08:** Any surveillance system must be non-networked and must store recordings locally. Cameras must not surveil the interiors of rooms where sensitive data is displayed.

## 6. Failure Modes

- **Single-perimeter dependence.** The operator relies entirely on one perimeter (typically the building envelope) and neglects the others. When that perimeter fails (a door is left unlocked, a window is broken), there is no backup. Mitigation: the concentric perimeter model requires multiple independent layers.
- **Environmental creep.** Conditions in the secure room gradually degrade -- temperature rises because of added hardware, humidity increases because a ventilation path is blocked, dust accumulates because cleaning is deferred. Hardware degrades silently until it fails. Mitigation: continuous monitoring and monthly review of environmental data.
- **Tamper-detection theater.** Seals are applied but never checked. Seals degrade from heat or age and can no longer show evidence of tampering. The tamper-detection system exists on paper but detects nothing. Mitigation: monthly seal verification is a mandatory operational task. Replace degraded seals proactively.
- **Key compromise.** A physical key is lost, copied, or stolen. The operator does not rekey the lock because it is inconvenient or expensive. An unauthorized person retains access indefinitely. Mitigation: immediate rekeying policy (Section 4.2). The cost of rekeying is small compared to the cost of undetected unauthorized access.
- **Complacent surveillance.** Alarm batteries die. Camera storage fills up and stops recording. Motion lights burn out. The security system degrades through neglect. Mitigation: monthly testing of all surveillance and alarm components as part of the operational tempo.

## 7. Recovery Procedures

1. **If unauthorized physical access is suspected:** Inspect all tamper-evident seals. Check position indicators. Review any surveillance footage. If evidence of intrusion is found, assume all systems within the breached perimeter may be compromised. Perform integrity checks on all institutional data (hash verification against known-good records). If tampering with hardware is evident, rebuild affected systems from known-good backups stored outside the breached area. Rekey all locks. Replace all tamper seals. Document the incident thoroughly.
2. **If environmental damage has occurred (fire, flood):** Assess the physical damage before powering on any equipment. For water damage, do not power on wet equipment; allow thorough drying and inspect for corrosion before attempting recovery. For fire damage, assess both direct damage and smoke/soot contamination. Restore from off-site backups if on-site hardware is unrecoverable. Document the event, including what was lost and what was recovered.
3. **If a physical key has been lost:** Rekey the affected lock within 24 hours. Issue new keys. Update the key log. Inspect tamper-detection indicators for any evidence that the lost key was used before the rekey. If the key loss occurred during a period when the site was unattended, treat as a potential unauthorized access event.
4. **If environmental monitoring reveals chronic out-of-range conditions:** Identify the cause. Upgrade cooling, ventilation, or dehumidification as needed. If conditions cannot be brought within range, consider relocating hardware to a more suitable area. Chronic out-of-range conditions are not acceptable; they are a slow-motion data loss event.

## 8. Evolution Path

- **Years 0-5:** Establish the physical security baseline. Install the secure room, locks, environmental monitoring, tamper seals, and basic surveillance. The initial implementation need not be perfect; it must be functional and improvable. Document the baseline thoroughly so that improvements can be measured against it.
- **Years 5-15:** Refine based on experience. Upgrade components that have proven inadequate. Replace aging environmental sensors. The physical security audit (SEC-008) will identify gaps.
- **Years 15-30:** The building itself may change. Renovations, moves, or succession may require the entire physical security design to be re-implemented at a new site. This article's principles remain valid regardless of location. The specific implementation must be adapted.
- **Years 30-50+:** Physical security hardware (locks, sensors, cameras) will have been replaced multiple times. The concentric perimeter model and the layered approach should endure. New physical security technologies will emerge; evaluate them against the institution's threat model and adopt what is justified.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
Physical security is the most tangible and the most neglected dimension of personal infrastructure security. People will spend hours configuring encryption and minutes choosing a door lock. This is backward. In an air-gapped institution, the door lock is more important than the encryption algorithm, because the door lock is the first line of defense and the encryption is the last. I have written this article to correct that imbalance. The physical security design does not need to be expensive or elaborate. It needs to be deliberate, layered, and maintained. A cheap lock on a solid door, checked monthly, is better than an expensive lock on a hollow door, never checked.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 1: Sovereignty; Principle 2: Integrity Over Convenience)
- CON-001 -- The Founding Mandate (self-sovereign infrastructure)
- SEC-001 -- Threat Model and Security Philosophy (physical threats, Category 1)
- SEC-002 -- Access Control Procedures (physical access as first layer)
- SEC-004 -- Air-Gap Architecture (physical enforcement of the air gap)
- SEC-008 -- Security Audit Framework (physical security audit procedures)
- GOV-001 -- Authority Model (succession protocol and key handover)
- Domain 4 -- Infrastructure & Power (hardware placement and power requirements)
- NIST SP 800-116 -- Guidelines for the Use of PIV Credentials in Facility Access
- Stage 1 Documentation Framework, Domain 3: Security & Integrity

---

---

# SEC-007 -- Supply Chain Security for Air-Gapped Systems

**Document ID:** SEC-007
**Domain:** 3 -- Security & Integrity
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, SEC-001, SEC-004, SEC-005, SEC-006
**Depended Upon By:** SEC-008, D4-001, D5-001, D18-001. All articles involving hardware procurement, software acquisition, or system deployment.

---

## 1. Purpose

This article addresses one of the hardest problems in air-gapped security: how do you trust hardware and software that originated outside your security boundary? Every piece of hardware in the institution was manufactured by someone else. Every piece of software was written by someone else. Every firmware image, every BIOS update, every USB controller was designed, fabricated, and packaged by entities whose trustworthiness the institution cannot verify from first principles.

This is the supply chain problem. It is not unique to air-gapped systems, but it is uniquely consequential for them. A networked system can be updated remotely if a supply chain compromise is discovered. An air-gapped system that has been compromised at the supply chain level may not exhibit any detectable symptoms -- the compromise may lie dormant in firmware, invisible to the operating system, waiting for a trigger that may never come or that the operator will never recognize.

This article specifies the procedures for procuring hardware and software as safely as possible in a world where the supply chain cannot be fully trusted. It defines verification procedures, clean room setup protocols, vendor diversification strategies, and the ongoing monitoring discipline that partially compensates for the irreducible uncertainty of trusting things you did not build yourself.

## 2. Scope

This article covers:

- The theory of supply chain risk and why it matters for air-gapped systems.
- Hardware procurement procedures: where to buy, what to verify, how to inspect.
- Firmware integrity checking: what can be verified and what cannot.
- Software verification procedures: source authentication, hash verification, reproducible builds.
- The clean room setup procedure: how to bring new hardware into the institutional environment safely.
- Vendor diversification: why single-vendor dependence is a supply chain risk and how to mitigate it.
- Ongoing monitoring for supply chain compromises.

This article does not cover the quarantine procedures for data on incoming media (see D18-001), the cryptographic tools used for verification (see SEC-003 and SEC-005), or the physical inspection of hardware for tampering during regular operations (see SEC-006). It defines the initial trust-establishment process that hardware and software must pass before they are admitted to the institutional environment.

## 3. Background

### 3.1 The Trust Hierarchy Problem

SEC-001 Section 4.3 defines the institution's trust model: the institution trusts hardware that has been "physically inspected and tested" and software that has been "verified from source and audited." These statements are aspirational. In practice, the depth of inspection and verification that a single operator can perform is limited.

You cannot inspect a CPU at the transistor level. You cannot verify that a hard drive's firmware does what its source code (which you do not have) says it does. You cannot audit every line of a Linux kernel. The trust model is not binary -- fully trusted or fully untrusted. It is a gradient, and the institution's supply chain procedures are designed to push hardware and software as far toward the "trusted" end of that gradient as is practically achievable.

The honest accounting required by ETH-001 Principle 6 demands that we acknowledge what supply chain security cannot do: it cannot guarantee that hardware is uncompromised. It can only raise the cost and difficulty of a supply chain attack to the point where it exceeds the likely adversary's resources and motivation. For a personal institution, this is a meaningful defense. For a nation-state target, it would not be. The procedures in this article are calibrated for the former.

### 3.2 Categories of Supply Chain Risk

**Hardware Supply Chain Risks:**
- Firmware implants: malicious code inserted into the firmware of storage controllers, network interfaces, management engines (Intel ME, AMD PSP), or peripheral controllers. These operate below the operating system and are invisible to standard security tools.
- Hardware trojans: modifications to the physical silicon that alter the chip's behavior in ways not specified by its design. Extremely difficult to detect. Currently considered a nation-state-level attack.
- Counterfeit components: components that claim to be one thing but are another. May have different performance characteristics, different failure modes, or different (potentially malicious) firmware.
- Interdiction: a targeted attack where hardware in transit is intercepted and modified before delivery. Documented in the Snowden disclosures for targeted nation-state operations.

**Software Supply Chain Risks:**
- Compromised source repositories: an attacker gains access to a software project's source code and inserts malicious code. The resulting builds are then distributed to all users.
- Compromised build infrastructure: the source code is clean, but the build process inserts malicious code into the binary. This is undetectable by source code review.
- Compromised distribution channels: the build is clean, but the download mirror has been tampered with. The binary the user receives is not the binary the project built.
- Dependency confusion: a malicious package with a similar name to a legitimate dependency is published, and build systems inadvertently download the malicious version.

### 3.3 The Air-Gapped Amplification Effect

Supply chain risks are amplified in an air-gapped environment for two reasons. First, the institution cannot easily receive security updates that patch a discovered supply chain compromise. The process of importing patches across the air gap is slower and more deliberate than a network update, which means the window of vulnerability is longer. Second, the institution's primary defense -- the air gap -- may be bypassed entirely by a supply chain attack, because the compromised component enters the environment through the legitimate hardware acquisition process, not through a network connection.

This amplification effect is why supply chain security is a dedicated article rather than a subsection of another article. It deserves concentrated attention.

## 4. System Model

### 4.1 Hardware Procurement Procedures

**Vendor Selection:**
- Prefer vendors with a track record of transparency and security responsiveness.
- Prefer hardware with open-source firmware options (e.g., coreboot-compatible motherboards, drives with open firmware projects).
- Avoid hardware with mandatory management engines or baseboard management controllers that cannot be disabled or neutralized (recognizing that complete avoidance may be impractical for modern x86 hardware; see Section 4.3 for mitigation).
- Purchase from authorized dealers or directly from manufacturers. Avoid gray-market hardware, open-box returns, or used hardware of unknown provenance.

**Procurement Diversification:**
- Do not source all hardware from a single vendor or a single retail channel. If a supply chain attack targets one vendor or one distribution path, diversification limits the institution's exposure.
- For critical components (storage drives, motherboards), purchase from at least two different manufacturers. This way, a firmware compromise in one manufacturer's product does not affect the entire institution.
- When replacing hardware, consider switching vendors rather than purchasing the same brand again. This limits the window during which a single vendor compromise could affect the institution.

**Purchase Practices:**
- Purchase hardware in person from retail locations when practical. This makes targeted interdiction more difficult (the attacker would need to compromise the retailer's entire stock rather than a single shipment).
- If ordering online, use unpredictable shipping addresses and timing. Do not announce hardware purchases publicly.
- Upon receipt, inspect packaging for signs of tampering: broken seals, re-taped boxes, inconsistent labels, unexpected weight.

### 4.2 Software Verification Procedures

**Source Authentication:**
Every piece of software that enters the institution must be verified as originating from its claimed source. This means:

- Download software from the project's official distribution channel, not from third-party mirrors when possible.
- Verify GPG signatures on software packages. The project's signing key must be obtained through an independent channel (for example, from the project's official website, cross-referenced with a key server) and verified against the key fingerprint published in the project's documentation.
- Verify SHA-256 (or stronger) checksums against checksums published by the project on a channel independent of the download channel.
- For operating system images, verify the chain: download the ISO, verify its signature against the distribution's signing key, and verify that the signing key is the expected key.

**Reproducible Builds:**
When available, prefer software that supports reproducible builds -- builds where anyone can compile the source code and produce a binary identical to the official release. This allows the operator to verify (or to trust that others have verified) that the binary matches the source. As of this writing, Debian, Arch Linux, and several other distributions have made significant progress on reproducible builds.

**Minimal Installation:**
Install only what is needed. Every additional software package increases the supply chain attack surface. The institution's systems should run a minimal operating system with only the packages required for their function. Unused packages should be removed, not merely ignored.

### 4.3 Firmware Integrity

Firmware is the most difficult component to verify because it operates below the operating system and is often distributed as opaque binary blobs without source code.

**What Can Be Verified:**
- Firmware updates from the manufacturer can be verified by checking their digital signatures (if the manufacturer provides signed updates).
- The firmware version running on a device can be compared against the expected version using operating system tools (e.g., dmidecode for BIOS, smartctl for drive firmware, lspci/lsusb for device firmware).
- Some firmware can be dumped and compared hash-for-hash against a known-good reference image.

**What Cannot Be Verified:**
- The behavior of firmware at the code level, without source code access, is opaque.
- Hardware management engines (Intel ME, AMD PSP) run their own operating system below the main CPU. They cannot be fully disabled on most modern hardware. Their code is proprietary and signed by the manufacturer.
- Storage controller firmware, USB controller firmware, and peripheral firmware are typically proprietary and unsigned (or signed by the manufacturer with keys the operator cannot verify).

**Mitigation for Unverifiable Firmware:**
- Use coreboot or Heads as a BIOS replacement on supported hardware. These open-source BIOS projects replace the manufacturer's proprietary firmware with auditable code.
- For Intel ME, use me_cleaner or the HAP (High Assurance Platform) bit to disable most ME functionality where supported.
- Prefer hardware platforms that minimize the amount of unverifiable firmware running at any time.
- Accept the residual risk honestly. Per ETH-001 Principle 6, acknowledge what cannot be verified rather than pretending the risk does not exist.

### 4.4 The Clean Room Setup Procedure

When new hardware enters the institution, it passes through a clean room setup process before being connected to any institutional data or network.

**Step 1: Physical Inspection.**
Open the hardware case. Inspect the motherboard, drive bays, and expansion slots for unexpected components. Compare against the manufacturer's documentation. Look for additional chips, modified traces, or components that do not appear in the schematic. Photograph the hardware internals for the institutional record.

**Step 2: Firmware Baseline.**
Boot the hardware from a known-good, verified live USB (prepared from a verified operating system image). Record all firmware versions: BIOS, drive firmware, USB controller firmware, any management engine version. Compare against known-good reference values (obtained from the manufacturer's published specifications and cross-referenced when possible).

**Step 3: BIOS/Firmware Configuration.**
Configure the BIOS for security: disable network boot, disable unused interfaces, set a BIOS password, enable secure boot if the institution's Linux distribution supports it with operator-controlled keys (not vendor keys). If using coreboot, flash the verified coreboot image.

**Step 4: Operating System Installation.**
Install the institution's verified operating system image. Perform the installation from a verified live USB. Do not use any network connection during installation. After installation, verify the installation integrity (package checksums, file system integrity).

**Step 5: Hardening.**
Apply the institution's system hardening procedures: disable unnecessary services, configure USBGuard, configure firewall rules (for internal-only networking if applicable), apply file system permissions per SEC-002.

**Step 6: Integration Testing.**
Test the new hardware with institutional workloads before connecting it to institutional data. Run storage benchmarks and compare against expected values. Run memory tests (memtest86+). Burn in the hardware for a minimum of 48 hours under load to detect early failures.

**Step 7: Documentation.**
Record the entire clean room process in the institutional hardware log: hardware model, serial numbers, firmware versions, configuration choices, test results, and the date of admission to the institutional environment. Apply tamper-evident seals (SEC-006) to the hardware case.

### 4.5 Vendor Diversification Strategy

Vendor diversification is the strategy of sourcing critical components from multiple independent manufacturers so that a compromise in one manufacturer's supply chain does not compromise the entire institution.

The diversification strategy applies to:
- **Storage:** Use drives from at least two different manufacturers. If one manufacturer's firmware is compromised, the other manufacturer's drives still hold clean backups.
- **Motherboards/CPUs:** If the institution has multiple systems, use different motherboard vendors or CPU architectures where practical.
- **USB devices:** Source authorized USB devices from different manufacturers.
- **Software:** Use different implementations for critical security functions where possible (e.g., different encryption tools for primary and backup encryption, though the algorithms should both be from the approved watchlist).

The goal is not paranoia but resilience. Monoculture is fragile. Diversity is robust.

## 5. Rules & Constraints

- **R-SEC07-01:** All new hardware must pass through the clean room setup procedure (Section 4.4) before being connected to institutional data or networks.
- **R-SEC07-02:** All software must be verified (source authentication, hash verification, and signature verification) before installation on institutional systems.
- **R-SEC07-03:** Critical components (storage, motherboards) must be sourced from at least two different manufacturers.
- **R-SEC07-04:** Hardware must be purchased from authorized channels. Gray-market hardware, open-box returns, and hardware of unknown provenance are prohibited without Tier 2 governance approval and additional verification steps.
- **R-SEC07-05:** Firmware versions on all institutional hardware must be recorded in the hardware log and compared against known-good references.
- **R-SEC07-06:** The clean room setup procedure must be documented for each hardware acquisition, including photographs of hardware internals.
- **R-SEC07-07:** Vendor diversification must be reviewed during the annual security audit (SEC-008). Single-vendor dependencies must be identified and a diversification plan created.

## 6. Failure Modes

- **Verification fatigue.** The operator grows tired of verifying signatures and checking hashes. Steps are skipped. "It is from the official website, it is probably fine." Probably is not verified. Mitigation: automate verification where possible (scripted hash checks, automated GPG verification). Make verification the default, not an extra step.
- **False trust in open source.** The operator assumes that because software is open source, it has been audited and is safe. Open-source software can be compromised just as proprietary software can. The xz backdoor incident (2024) demonstrated that even well-known, widely-used open-source projects can harbor supply chain attacks. Mitigation: verify regardless of source. Trust but verify is the model; trust alone is not.
- **Clean room shortcuts.** The clean room procedure takes hours. The operator needs the hardware now. Steps are abbreviated. Integration testing is skipped. The hardware enters the institutional environment without proper vetting. Mitigation: the clean room procedure is a Tier 3 procedure that cannot be modified without governance approval. Shortcuts are policy violations, not efficiencies.
- **Vendor monoculture.** The operator purchases all hardware from one vendor because it is familiar, convenient, or because they got a good deal. A compromise in that vendor's supply chain affects every system. Mitigation: R-SEC07-03 mandates multi-vendor sourcing for critical components.
- **Firmware blindness.** The operator cannot verify firmware because the source code is not available. The risk is acknowledged but never actively managed. Over time, the residual risk accumulates as more unverifiable firmware enters the environment. Mitigation: actively seek hardware with open firmware options. Use coreboot and me_cleaner where supported. Review firmware risk during the annual security audit.

## 7. Recovery Procedures

1. **If a supply chain compromise is discovered in deployed hardware:** Assess the scope. Identify all institutional systems that contain the affected component. Determine whether the compromise allows data exfiltration (unlikely in an air-gapped environment but possible through an insider or during sneakernet transfers), data destruction, or silent data corruption. Isolate affected systems. Verify data integrity against known-good checksums. If integrity is compromised, restore from backups stored on unaffected hardware. Replace the compromised component. Document the incident.
2. **If a supply chain compromise is discovered in deployed software:** Assess the scope. Determine what the compromised software had access to. Verify data integrity. Apply a clean version of the software (verified through the procedures in Section 4.2). If the compromise allowed privilege escalation, rebuild the affected system from a verified base. Document the incident.
3. **If the clean room procedure was not followed for hardware already in production:** Retroactively perform the verification steps that can be performed on deployed hardware (firmware version check, physical inspection if hardware can be taken offline). Accept that some verification steps (like integration testing on a clean system) cannot be performed retroactively. Document the gap. Ensure the procedure is followed for all future acquisitions.
4. **If vendor diversification has been neglected:** Conduct a vendor audit. Identify single-vendor dependencies. Create a diversification plan with timelines. When hardware is next replaced, source the replacement from a different vendor. Do not force premature replacement solely for diversification; address it at natural replacement points.

## 8. Evolution Path

- **Years 0-5:** Establish the clean room procedure and vendor diversification strategy. The initial hardware procurement is the most critical: these are the systems that will handle all institutional data during the founding period. Be thorough. Document everything.
- **Years 5-15:** Hardware will be replaced as components reach end-of-life. Each replacement is an opportunity to diversify vendors and to upgrade to hardware with better open-firmware support. The clean room procedure should be routine by now.
- **Years 15-30:** The supply chain landscape will have changed. New manufacturers will exist. New firmware transparency initiatives may make verification easier. New attack techniques may make verification harder. Adapt the procedures to the current landscape while maintaining the principles.
- **Years 30-50+:** The supply chain risks will be different but no less real. The fundamental approach -- verify what you can, diversify to limit exposure, acknowledge what you cannot verify, and maintain the discipline to do this for every acquisition -- will remain valid regardless of how the technology evolves.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
Supply chain security is the domain where honesty is most uncomfortable. We depend on hardware and software we did not build, cannot fully inspect, and must trust to some degree. The procedures in this article do not eliminate supply chain risk. They manage it. They raise the bar an attacker must clear. They limit the blast radius of a compromise. And they force the operator to make deliberate, documented decisions about what to trust and why, rather than trusting by default and hoping for the best. The xz backdoor of 2024 demonstrated that even the most widely-reviewed open-source software is not immune. If anything, that incident reinforced the value of these procedures: verify everything you can, diversify to limit what a single compromise can reach, and never assume that "widely used" means "trustworthy."

## 10. References

- ETH-001 -- Ethical Foundations (Principle 6: Honest Accounting of Limitations)
- CON-001 -- The Founding Mandate (self-built mandate, Section 3.3)
- SEC-001 -- Threat Model and Security Philosophy (supply chain threats, Category 3)
- SEC-004 -- Air-Gap Architecture (sneakernet protocol for software import)
- SEC-005 -- Long-Term Cryptographic Survival (algorithm verification for cryptographic software)
- SEC-006 -- Physical Security Design (tamper-evident seals for hardware)
- SEC-008 -- Security Audit Framework (supply chain audit procedures)
- D18-001 -- Import & Quarantine Procedures (quarantine for incoming media and data)
- Stage 1 Documentation Framework, Domain 3: Security & Integrity

---

---

# SEC-008 -- Security Audit Framework

**Document ID:** SEC-008
**Domain:** 3 -- Security & Integrity
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, SEC-001, SEC-002, SEC-003, SEC-004, SEC-005, SEC-006, SEC-007, OPS-001
**Depended Upon By:** All articles in Domain 3. All articles involving security verification, compliance, or assurance.

---

## 1. Purpose

This article defines the complete security audit methodology for the holm.chat Documentation Institution. A security architecture is only as strong as its weakest implemented component, and the only way to know whether every component is actually implemented and functioning is to audit it. Regularly, systematically, and honestly.

The security articles that precede this one -- SEC-001 through SEC-007 -- define what the institution's security should look like. This article defines how the institution verifies that it actually looks that way. It is the difference between claiming to be secure and demonstrating that you are secure, to the extent that demonstration is possible.

Security auditing in a single-operator institution is unusual. There is no independent auditor. The operator audits themselves. This creates an obvious conflict of interest: the person responsible for implementing security is also responsible for evaluating it. This article addresses that conflict directly, through structured checklists, documented evidence requirements, and a commitment to honesty that must override the natural human tendency to find what you want to find rather than what is actually there.

## 2. Scope

This article covers:

- The audit philosophy: why self-audit is possible and how to make it honest.
- The audit cycle: how often different domains are audited.
- The audit checklist: what is checked in each security domain.
- Evidence requirements: what constitutes sufficient proof that a control is functioning.
- Findings classification: how to categorize and prioritize what the audit discovers.
- Remediation tracking: how to ensure findings are addressed, not merely recorded.
- The annual security review: the comprehensive audit that ties everything together.

This article does not perform audits -- it defines how audits are performed. The actual audit results are recorded in the security audit log, a separate institutional record.

## 3. Background

### 3.1 The Self-Audit Problem

The greatest challenge of auditing a single-operator institution is that the auditor and the auditee are the same person. This means:

- You know where the weaknesses are, which means you can unconsciously avoid looking at them.
- You know what the "right" answer should be for each audit check, which means you can unconsciously confirm rather than investigate.
- You have a personal interest in a clean audit result, which means you can unconsciously minimize findings.

These are not accusations. They are descriptions of well-documented cognitive biases -- confirmation bias, motivated reasoning, and the ostrich effect. They affect everyone. The mitigation is not willpower (willpower loses to bias over time). The mitigation is structure: checklists that force you to look where you might not want to, evidence requirements that prevent you from answering "yes" without proof, and a findings classification system that makes the cost of hiding a finding greater than the cost of addressing it.

### 3.2 The Value of Honest Findings

ETH-001 Principle 6 requires honest accounting of limitations. In the audit context, this means: a finding is not a failure. A finding is information. An audit that discovers problems is a successful audit. An audit that discovers nothing is either conducted against a genuinely perfect security posture (unlikely) or conducted without sufficient rigor (likely).

The purpose of the audit is not to produce a clean report. The purpose is to produce an accurate report. The operator must internalize this before conducting the first audit. An audit culture that values clean reports over accurate reports will produce clean reports that hide real problems until those problems cause real damage.

### 3.3 Audit Scope and Frequency

Not everything needs to be audited at the same frequency. The audit cycle is risk-based: controls that protect against high-likelihood, high-impact threats are audited more frequently than controls that protect against low-likelihood threats. The operational tempo (OPS-001) already includes security-related checks at daily, weekly, monthly, quarterly, and annual intervals. The audit framework builds on that tempo, adding formal verification to the routine checks.

## 4. System Model

### 4.1 The Audit Cycle

The audit cycle defines how often each security domain is formally audited:

**Monthly Audits (focused, 1-2 hours):**
- Physical security spot check: tamper seals, door locks, environmental readings.
- Air-gap verification: confirm no unauthorized network interfaces, check USBGuard logs.
- Backup integrity: test-restore a randomly selected backup set.

**Quarterly Audits (comprehensive per domain, half day):**
- One full domain audit per quarter, rotating through the four primary domains:
  - Q1: Physical security and access control (SEC-002, SEC-006).
  - Q2: Air-gap integrity and data transfer discipline (SEC-004, D18-001).
  - Q3: Cryptographic health and key management (SEC-003, SEC-005).
  - Q4: Supply chain posture and firmware integrity (SEC-007).

**Annual Audit (full security review, 2-3 days):**
- All security domains audited comprehensively.
- Threat model review and update.
- Cryptographic watchlist review.
- Security drill execution (simulated incident response).
- Findings review from all audits during the year.
- Remediation verification for all open findings.
- Security posture assessment and planning for the coming year.

### 4.2 The Audit Checklist

The following checklist is structured by security domain. Each item includes the control being verified, the evidence required to confirm the control is functioning, and the reference to the article that specifies the control.

**Domain A: Air-Gap Integrity (SEC-004)**

- A1: No wireless network interfaces enabled on any institutional system. Evidence: lspci and lsusb output showing no wireless hardware; physical inspection confirming antenna disconnection. Frequency: quarterly and after any hardware change.
- A2: All USB ports controlled by USBGuard or equivalent. Evidence: USBGuard policy review; test with unauthorized USB device (should be rejected). Frequency: quarterly.
- A3: Air-gap transfer log is current and complete. Evidence: review transfer log entries against USB device registry; cross-reference with quarantine station logs. Frequency: quarterly.
- A4: Exception registry is current. Evidence: review all entries; verify expiration dates; confirm compensating controls are in place. Frequency: quarterly.
- A5: Power isolation is intact. Evidence: UPS is functional; verify double-conversion operation; test battery backup. Frequency: quarterly.
- A6: No personal wireless devices present during sensitive operations. Evidence: procedural review; operator self-attestation. Frequency: monthly.

**Domain B: Physical Security (SEC-006)**

- B1: All perimeters intact. Evidence: physical walkthrough of all four perimeters; door and window inspection; lock function test. Frequency: quarterly.
- B2: Tamper-evident seals intact on all hardware. Evidence: visual inspection of all seals; serial number comparison against tamper log. Frequency: monthly.
- B3: Environmental conditions within range. Evidence: temperature and humidity log review for the audit period; identification of any out-of-range events. Frequency: quarterly.
- B4: Fire detection functional. Evidence: smoke detector test (test button). Frequency: monthly.
- B5: Water detection functional. Evidence: water sensor test. Frequency: quarterly.
- B6: Physical key inventory accurate. Evidence: count and verify all keys against key log; confirm off-site backup key exists. Frequency: annually.
- B7: Surveillance systems functional. Evidence: camera footage review; alarm system test; motion light test. Frequency: quarterly.

**Domain C: Access Control (SEC-002)**

- C1: Account inventory accurate. Evidence: list all system accounts; compare against authorized account list; identify and investigate any discrepancies. Frequency: quarterly.
- C2: Least privilege enforced. Evidence: for each account, list its permissions and justify each permission against the account's function. Frequency: annually.
- C3: Root account usage appropriate. Evidence: review authentication logs for root account usage; verify each instance was justified. Frequency: quarterly.
- C4: Password/passphrase strength compliant. Evidence: verify password policy enforcement; review password age. Frequency: annually.
- C5: Audit trail integrity. Evidence: verify that authentication and authorization logs are being written and are not tampered with. Frequency: quarterly.

**Domain D: Cryptographic Health (SEC-003, SEC-005)**

- D1: Cryptographic watchlist current. Evidence: review watchlist; verify each algorithm's status against current research (obtained via sneakernet). Frequency: annually.
- D2: Full-disk encryption active on all institutional systems. Evidence: LUKS status check; verify encryption is active and key slot configuration is correct. Frequency: quarterly.
- D3: Cryptographic metadata present on all encrypted/hashed data. Evidence: sample check of encrypted files and hash databases; verify metadata completeness. Frequency: annually.
- D4: Key management procedures followed. Evidence: review key inventory; verify backup key existence; verify succession key provisions. Frequency: annually.
- D5: No deprecated algorithms in active use. Evidence: cross-reference institutional systems against watchlist; identify any Orange or Red status algorithms. Frequency: annually.

**Domain E: Supply Chain Posture (SEC-007)**

- E1: All hardware passed clean room procedure. Evidence: review hardware log; verify clean room documentation exists for each device. Frequency: annually.
- E2: Software verification records complete. Evidence: review software installation log; verify GPG signature and hash verification records. Frequency: annually.
- E3: Vendor diversification maintained. Evidence: list all critical components by vendor; identify single-vendor dependencies. Frequency: annually.
- E4: Firmware versions current and documented. Evidence: compare running firmware against hardware log; identify discrepancies. Frequency: annually.

**Domain F: Operational Discipline (OPS-001)**

- F1: Operational tempo adherence. Evidence: review operational logs for the audit period; identify missed daily, weekly, or monthly operations. Frequency: quarterly.
- F2: Documentation currency. Evidence: select five random systems; compare documentation against actual configuration; identify discrepancies. Frequency: annually.
- F3: Decision log currency. Evidence: review decision log; identify any decisions that were not logged at the appropriate tier. Frequency: quarterly.
- F4: Backup test success. Evidence: review backup test-restoration logs for the audit period; verify at least one test per month was performed. Frequency: quarterly.

### 4.3 Evidence Requirements

For each audit item, the evidence must meet three criteria:

**Contemporaneous:** The evidence must be produced during or immediately before the audit. A screenshot from six months ago does not prove that the control is functioning today.

**Verifiable:** The evidence must be reproducible. Another person (such as a successor) should be able to follow the same steps and produce similar evidence. "I checked and it looked fine" is not evidence. A command output, a photograph, a log entry -- these are evidence.

**Recorded:** The evidence must be stored in the security audit log. The audit log is an append-only institutional record that preserves the history of all audit activities and findings. It is the institution's proof -- to itself and to future operators -- that security was actually maintained, not merely claimed.

### 4.4 Findings Classification

Audit findings are classified by severity:

**Critical:** A control that is completely absent or non-functional, exposing the institution to a high-likelihood, high-impact threat. Example: full-disk encryption is not active on a system containing institutional data. Remediation deadline: 7 days. If remediation cannot be completed in 7 days, the affected system must be taken offline until the control is restored.

**Major:** A control that is partially functional or significantly degraded, or a complete absence of a control for a moderate-impact threat. Example: tamper-evident seals have not been checked in three months; several seals have degraded. Remediation deadline: 30 days.

**Minor:** A control that is functional but has a deficiency that does not immediately expose the institution to significant risk. Example: the air-gap transfer log is missing entries for two transfers that can be reconstructed from other records. Remediation deadline: 90 days.

**Observation:** A condition that is not a deficiency but represents an improvement opportunity or a trend that could become a deficiency if not addressed. Example: environmental temperatures are within range but trending upward over the past quarter. Remediation deadline: next scheduled review.

### 4.5 Remediation Tracking

Every finding at Critical, Major, or Minor severity is tracked in the remediation log until resolved. The remediation log records:

- **Finding ID:** Links to the audit log entry.
- **Severity:** Critical, Major, or Minor.
- **Description:** What was found.
- **Root Cause:** Why the deficiency exists (determined during investigation).
- **Remediation Plan:** What will be done to fix it.
- **Deadline:** Based on severity classification.
- **Status:** Open, In Progress, Remediated, Verified.
- **Verification:** Evidence that the remediation was effective, collected during the next audit.

A finding is not closed when the fix is applied. It is closed when the next audit verifies that the fix is effective. This prevents the common pattern of "fixing" something in a way that does not actually work, then forgetting about it until the next audit discovers the same problem again.

### 4.6 The Annual Security Review

The annual security review is the capstone of the audit cycle. It is the most comprehensive assessment of the institution's security posture and is conducted during the annual operations cycle defined in OPS-001.

The annual review comprises:

**Day 1: Comprehensive Audit.**
Execute the full audit checklist (all domains, all items). Record evidence. Classify findings.

**Day 2: Analysis and Drill.**
- Review all findings from the past year's audits. Identify patterns. Are the same types of findings recurring? If so, the underlying cause has not been addressed.
- Update the threat model. Has anything changed in the threat landscape that affects the institution?
- Review the cryptographic watchlist. Update algorithm statuses based on the past year's research.
- Execute a security drill: simulate a realistic security incident (e.g., "your primary storage drive has failed," "you suspect unauthorized physical access occurred while you were traveling," "a critical algorithm on the watchlist has been broken"). Walk through the response procedure from start to finish. Document what worked and what did not.

**Day 3: Planning and Documentation.**
- Review the remediation log. Verify that all findings from the previous year have been addressed.
- Assess the overall security posture. Is it improving, stable, or degrading?
- Plan security improvements for the coming year. What controls need to be upgraded? What procedures need to be revised? What new threats need to be addressed?
- Update the Commentary Section of relevant security articles.
- Write the annual security review summary: a one-page document that captures the state of security at this point in time, for the institutional record and for future operators.

## 5. Rules & Constraints

- **R-SEC08-01:** The audit cycle defined in Section 4.1 is mandatory. Monthly, quarterly, and annual audits must be performed on schedule. Deviations must be documented and the missed audit must be performed within 14 days of the scheduled date.
- **R-SEC08-02:** All audit findings must be classified according to the severity system in Section 4.4 and tracked in the remediation log until verified as resolved.
- **R-SEC08-03:** Critical findings must be remediated within 7 days. If remediation is not possible within 7 days, the affected system must be isolated or taken offline.
- **R-SEC08-04:** The annual security review must include a security drill testing at least one incident response scenario.
- **R-SEC08-05:** Audit evidence must be contemporaneous, verifiable, and recorded in the security audit log.
- **R-SEC08-06:** The security audit log is append-only. Entries may not be modified or deleted. Corrections are recorded as new entries referencing the original.
- **R-SEC08-07:** Findings that recur in consecutive audits must be escalated: their severity is increased by one level (Minor becomes Major, Major becomes Critical) until the root cause is addressed.
- **R-SEC08-08:** The audit checklist must be reviewed and updated during each annual security review to reflect changes in the institution's systems and threat landscape.

## 6. Failure Modes

- **Audit avoidance.** The operator skips audits because they are time-consuming, anxiety-inducing, or apparently unnecessary. Security degrades undetected. Mitigation: the audit cycle is integrated into the operational tempo (OPS-001). Audits are scheduled, not optional. Chronic audit avoidance triggers a governance review.
- **Confirmation bias in self-audit.** The operator finds what they expect to find rather than what is actually there. Controls are marked as "pass" without rigorous verification. Mitigation: evidence requirements force the operator to produce proof, not merely assert compliance. Checklists force the operator to look at every control, not just the ones they are confident about.
- **Remediation drift.** Findings are recorded but not remediated. The remediation log accumulates open items. The operator becomes desensitized to the backlog. Mitigation: R-SEC08-07 escalates recurring findings. The annual review explicitly checks the remediation log. Open findings older than their deadline are flagged prominently.
- **Checklist ossification.** The checklist becomes a rote exercise that does not evolve with the institution's changing systems and threat landscape. New systems are deployed without being added to the checklist. The audit checks things that no longer matter and ignores things that do. Mitigation: R-SEC08-08 requires annual checklist review and update. The checklist is a living document.
- **Audit theater.** The audits are performed, the logs are filled, the findings are classified and remediated -- but none of it is honest. The operator goes through the motions without genuine investigation. This is the most dangerous failure mode because it produces a perfect-looking record of a declining security posture. Mitigation: the security drill is the antidote. You cannot fake an incident response. Either the recovery procedure works or it does not. The drill provides ground truth that the checklist alone cannot.

## 7. Recovery Procedures

1. **If audits have been skipped:** Perform a full annual-level audit immediately, regardless of where the institution is in the audit cycle. Classify all findings. Prioritize Critical and Major findings for immediate remediation. Resume the normal audit cycle from the current date.
2. **If the remediation log has accumulated a backlog:** Triage the backlog. Address Critical findings first, Major second, Minor third. For findings that have persisted beyond their deadline, escalate per R-SEC08-07. If the backlog is too large to address in a reasonable timeframe, consider which findings represent actual risk and which are paperwork. Address the risks. Accept and document the rest with a realistic remediation timeline.
3. **If the audit checklist is outdated:** Conduct a full inventory of institutional systems and controls. Compare against the current checklist. Add missing items. Remove items that reference decommissioned systems. Update evidence requirements to reflect current tooling. This is a Tier 3 governance activity.
4. **If audit honesty has been compromised:** This is the hardest recovery because it requires the operator to acknowledge their own self-deception. Begin by re-reading ETH-001 Principle 6 (Honest Accounting of Limitations) and SEC-001 Section 4.1 Pillar 3 (The Human Is the System). Then perform a full audit with the explicit goal of finding problems. Set a minimum: identify at least three genuine findings. If you cannot find three, you are not looking hard enough. Record the honesty reset in the Commentary Section.

## 8. Evolution Path

- **Years 0-5:** The audit framework is being established alongside the systems it audits. Expect the checklist to change frequently as new systems are deployed and as the operator learns what is worth checking and what is not. The first few annual reviews will be learning experiences. Document the lessons.
- **Years 5-15:** The audit framework should be stable. The checklist will evolve but not dramatically. The remediation log should show a trend of fewer recurring findings as systemic issues are addressed. The annual security drill should test increasingly realistic scenarios.
- **Years 15-30:** Succession becomes relevant. The audit framework must be documented thoroughly enough that a new operator can conduct audits without oral instruction. The security audit log becomes a historical record of the institution's security posture over time, valuable for understanding trends and for giving the new operator confidence that the systems they are inheriting have been maintained.
- **Years 30-50+:** The audit framework will have been revised many times. The specific checklist items will bear little resemblance to the originals. But the methodology -- regular, structured, evidence-based, honestly reported, and followed by verified remediation -- should be unchanged. The methodology is the framework. The checklist is merely its current expression.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
I have written this article last among the five security articles in this batch, and that is intentional. The audit framework is the verification layer that sits on top of everything else. It exists to answer the question that every other security article avoids: "Is any of this actually working?" It is the most important security article in the institution, not because it defines controls (the other articles do that) but because it forces the operator to prove, regularly and honestly, that those controls are real and not aspirational. The greatest risk to this institution's security is not a sophisticated attacker or a novel vulnerability. It is the operator's natural tendency to believe that things are working because they were set up correctly once. They were not. They are working because someone checks them, maintains them, and fixes them when they break. The audit framework is what makes that someone -- you, or whoever comes after you -- actually do the checking.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 2: Integrity Over Convenience; Principle 6: Honest Accounting)
- CON-001 -- The Founding Mandate (institutional mission and boundaries)
- GOV-001 -- Authority Model (decision tiers for audit-triggered changes)
- SEC-001 -- Threat Model and Security Philosophy (threat categories and security mindset)
- SEC-002 -- Access Control Procedures (Domain C audit items)
- SEC-003 -- Key Management (Domain D audit items)
- SEC-004 -- Air-Gap Architecture (Domain A audit items)
- SEC-005 -- Long-Term Cryptographic Survival (cryptographic watchlist review)
- SEC-006 -- Physical Security Design (Domain B audit items)
- SEC-007 -- Supply Chain Security (Domain E audit items)
- OPS-001 -- Operations Philosophy (operational tempo integration)
- D10-003 -- Incident Response Procedures (security drill reference)
- D18-001 -- Import & Quarantine Procedures (quarantine audit items)
- NIST SP 800-53 -- Security and Privacy Controls (control framework reference)
- ISO 27001 -- Information Security Management Systems (audit methodology reference)
- Stage 1 Documentation Framework, Domain 3: Security & Integrity

---

*End of Stage 4 Specialized Security Systems -- Five Advanced Articles*

**Document Total:** 5 articles (SEC-004 through SEC-008)
**Combined Estimated Word Count:** ~15,000 words
**Status:** All five articles ratified as of 2026-02-16.
**Next Stage:** Stage 4 continues with specialized articles in other domains. These five articles provide the security reference foundation for all subsequent specialized work.
