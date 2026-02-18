# STAGE 3: OPERATIONAL DOCTRINE -- PLATFORM & CORE SYSTEMS

## Domain 5 Operations Articles D5-002 through D5-006

**Document ID:** STAGE3-PLATFORM-OPS
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Classification:** Stage 3 -- Operational Doctrine. Practical manuals for the daily operation of the institution's platform layer.
**Depends On:** ETH-001, CON-001, GOV-001, SEC-001, OPS-001, Stage 1 Domain 5 Framework
**Depended Upon By:** All articles involving platform maintenance, system administration, network operations, or configuration management.

---

## How to Read This Document

This document contains five operational articles for Domain 5: Platform & Core Systems. These are not philosophical documents. They are manuals. They tell you what to do, when to do it, and what to do when something goes wrong.

Each article follows the canonical template defined in Stage 1. Each article traces its authority back to the root documents in Stage 2. Where a procedure in this document conflicts with a principle in Stage 2, the principle prevails and this document must be amended.

These articles assume you have read OPS-001 (Operations Philosophy) and understand the operational tempo. They assume you have read SEC-001 (Threat Model and Security Philosophy) and understand the air-gap constraint. If you have not read those documents, stop here and read them first. The procedures in this document will not make sense without that context.

If you are performing a specific task and need a specific procedure, use the table of contents below to navigate directly. If you are reading for the first time, read sequentially. The articles build on one another.

---
---

# D5-002 -- Operating System Maintenance Procedures

**Document ID:** D5-002
**Domain:** 5 -- Platform & Core Systems
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, SEC-001, OPS-001, P-001, P-002, P-003, P-012
**Depended Upon By:** D5-003, D5-004, D5-005, D5-006, all Domain 12 (Disaster Recovery) articles involving system rebuilds.

---

## 1. Purpose

This article defines how the operating system is maintained on every machine within the institution. Maintenance means patching, updating, upgrading, and occasionally choosing not to do any of those things. In a conventional, internet-connected environment, operating system maintenance is largely automated: package managers pull updates from remote repositories, security patches arrive through automatic update services, and the operator's primary decision is whether to reboot now or later. None of that applies here.

This institution is air-gapped. There are no remote repositories. There are no automatic updates. There is no mechanism by which a security advisory will appear on your screen unbidden. Every update that enters this institution does so because a human being physically transported it across the air gap, verified its integrity, tested it in isolation, and deliberately applied it. This article defines the procedures that make that process reliable, repeatable, and safe.

The purpose is not merely technical. Operating system maintenance in an air-gapped environment is a governance act. Every patch applied is a decision to trust external code. Every patch withheld is a decision to accept a known vulnerability. Both decisions carry risk. Both decisions must be documented. This article provides the framework for making those decisions well.

## 2. Scope

This article covers:

- The process for identifying which OS updates are available and relevant.
- The process for acquiring updates across the air gap.
- The process for testing updates before production deployment.
- The process for applying updates to production systems.
- Rollback procedures when updates cause problems.
- Kernel management: when to update the kernel, when to hold.
- The decision framework for major version upgrades versus staying on a current release.
- Security patch triage for offline systems.

This article does NOT cover:

- The initial operating system installation (see P-003).
- Application-level software updates (see D5-006 for configuration management, P-002 for software acquisition).
- Hardware-level firmware updates (see Domain 4 articles).
- The quarantine process for incoming media (see Domain 18). This article assumes media has already cleared quarantine.

## 3. Background

Operating system maintenance is the single most routine and most consequential platform task. It is routine because it happens regularly -- monthly at minimum, often more frequently. It is consequential because a failed update can render a system unbootable, and a missed security patch can leave a known vulnerability open indefinitely.

In internet-connected environments, the community handles much of this burden. Security advisories propagate automatically. Package managers resolve dependencies. Rollback mechanisms are built into update tools. In an air-gapped environment, the operator assumes every one of these responsibilities personally. You are the distribution mirror. You are the advisory aggregation service. You are the dependency resolver. This is the cost of sovereignty, and it must be paid consistently.

The historical record of institutions that deferred maintenance -- whether physical or digital -- is uniformly grim. Deferred maintenance compounds. A missed patch becomes two missed patches becomes a system so far behind that updating it requires a full reinstallation. This article exists to prevent that cascade.

### 3.1 The Air-Gap Maintenance Paradox

There is a fundamental tension in maintaining an air-gapped system. The air gap exists to protect the institution from external threats. But operating system updates come from the outside world. Every update is, by definition, a controlled breach of the air gap. The institution must manage this tension explicitly: allowing necessary updates through while maintaining the integrity that the air gap provides.

The resolution is procedural, not technical. The air gap is not breached by the update itself but would be breached by an uncontrolled update process. A controlled, documented, verified update process preserves the spirit of the air gap even as it necessarily crosses the physical boundary.

## 4. System Model

### 4.1 The Update Pipeline

Operating system updates flow through a five-stage pipeline. No stage may be skipped.

**Stage 1: Intelligence Gathering.**
The operator monitors external sources for information about available updates. This happens outside the air gap, on a separate, non-institutional device. Sources include: the distribution's security advisory mailing list, the distribution's release notes, and relevant security databases (CVE databases, distribution-specific trackers). The operator records which updates are available and which are security-relevant.

**Stage 2: Triage and Selection.**
Not every available update needs to be applied. The operator evaluates each update against three criteria: Does it fix a security vulnerability that is relevant to our deployment? Does it fix a bug we have encountered or are likely to encounter? Does it provide functionality we need? Updates that meet none of these criteria are deferred. Updates that meet at least one are candidates for the next stages.

**Stage 3: Acquisition and Verification.**
Selected updates are downloaded on the external device. Their cryptographic signatures are verified against known-good keys. The packages are transferred to physical media -- USB drives or optical media -- following the quarantine procedures defined in Domain 18. Once quarantined media is cleared, the packages are available for testing.

**Stage 4: Testing.**
Updates are applied to a non-production test system first. This test system must mirror the production configuration as closely as feasible. After applying updates, the operator runs the test suite defined in Section 4.4 of this article. The test system remains in the updated state for a minimum observation period of 48 hours before production deployment is authorized.

**Stage 5: Production Deployment.**
Updates are applied to production systems one at a time, starting with the least critical system. After each system is updated, the operator runs the production verification checklist defined in Section 4.5. If any verification step fails, the operator halts deployment and initiates rollback procedures.

### 4.2 The Security Patch Triage Matrix

Security patches receive special handling. They are classified by two dimensions: severity and relevance.

**Severity** follows the standard CVSS framework, simplified for institutional use:
- **Critical:** Remote code execution, privilege escalation, or data exposure without authentication.
- **High:** Vulnerabilities that require some precondition (authenticated access, specific configuration) but have severe impact.
- **Medium:** Vulnerabilities with limited impact or that require unlikely preconditions.
- **Low:** Theoretical vulnerabilities, information disclosure of non-sensitive data, or denial-of-service against non-critical services.

**Relevance** is assessed against the institution's actual deployment:
- **Direct:** The vulnerable component is installed and actively used.
- **Indirect:** The vulnerable component is installed but not actively used, or is a dependency of an active component.
- **None:** The vulnerable component is not installed.

The triage matrix:

| | Critical | High | Medium | Low |
|---|---|---|---|---|
| **Direct** | Apply within 7 days | Apply within 14 days | Apply at next maintenance cycle | Apply at convenience |
| **Indirect** | Apply within 14 days | Apply at next maintenance cycle | Apply at convenience | Defer |
| **None** | No action | No action | No action | No action |

These timelines are targets, not absolutes. The air-gap acquisition process adds latency that internet-connected systems do not face. The timelines begin when the operator becomes aware of the vulnerability, not when the patch is released.

### 4.3 Kernel Management

The kernel is the most sensitive component of the operating system. A failed kernel update can render a system unbootable. Kernel updates therefore receive additional scrutiny.

**Rules for kernel updates:**
- Always keep at least two known-good kernels installed and available at boot time.
- Before updating the kernel, verify that the bootloader is configured to fall back to the previous kernel if the new one fails to boot.
- Kernel updates are never applied to production systems on the same day they are tested. A minimum 72-hour observation period on the test system is required for kernel updates, compared to 48 hours for other updates.
- Kernel major version changes (e.g., 5.x to 6.x) are treated as Tier 2 architectural decisions per GOV-001. They require a written proposal, a 30-day waiting period, and a documented rollback plan.
- Kernel configuration changes (recompilation, module changes) follow the same testing pipeline as kernel version updates.

### 4.4 The Test Suite

After applying updates to the test system, the following checks must pass:

1. The system boots successfully to the expected runlevel or target.
2. All expected filesystems mount correctly.
3. All expected network interfaces come up with correct configuration.
4. All core services start and pass their health checks (see D5-005).
5. Disk I/O performance is within 10% of baseline measurements.
6. Memory utilization at idle is within 15% of baseline.
7. All scheduled cron jobs or systemd timers are intact and correctly configured.
8. The system log shows no errors or unexpected warnings during boot and initial operation.
9. A test backup completes successfully and can be restored.
10. Authentication works correctly for all configured accounts.

If any check fails, the update is not cleared for production deployment. The failure is documented, and the operator investigates before proceeding.

### 4.5 Production Verification Checklist

After applying updates to each production system:

1. All items from the test suite (Section 4.4) are verified.
2. All production workloads resume correctly.
3. All inter-system connections are functional (see D5-004 for network verification).
4. Monitoring dashboards show normal values (see D5-005).
5. A snapshot or backup of the updated system state is created before proceeding to the next system.

### 4.6 The Upgrade Decision Framework

Periodically, the institution must decide whether to upgrade to a new major OS release or remain on the current one. This decision is governed by the following framework:

**Stay on the current release when:**
- The current release is still receiving security patches.
- No hardware support gap has emerged (new hardware requires a newer kernel or driver).
- No critical software dependency requires the newer release.
- The effort of upgrading exceeds the benefit.

**Upgrade to the new release when:**
- The current release is approaching end of security support.
- Hardware support requires it.
- Critical software functionality requires it.
- The new release has been available for at least six months and the institution has observed its stability in community reports.

Major OS upgrades are Tier 2 architectural decisions per GOV-001. They require planning, documentation, and a full rollback strategy. The rollback strategy for a major upgrade must include the ability to restore from a pre-upgrade backup to a fully functional state.

## 5. Rules & Constraints

- **R-D5-002-01:** No operating system update shall be applied to a production system without first being tested on a non-production system.
- **R-D5-002-02:** All updates must be cryptographically verified before application. Updates that fail signature verification are rejected and reported as a security incident.
- **R-D5-002-03:** At least two bootable kernels must be available on every production system at all times.
- **R-D5-002-04:** The operator must conduct intelligence gathering for available updates at least monthly. Security-critical sources must be checked at least bi-weekly.
- **R-D5-002-05:** Every update applied to a production system must be recorded in the operational log with: date, system affected, packages updated, reason for update, test results, and post-deployment verification results.
- **R-D5-002-06:** Major OS version upgrades require Tier 2 governance per GOV-001.
- **R-D5-002-07:** No update shall be applied during a period when rollback capabilities are compromised (e.g., when backups are in progress, when the test system is unavailable, or when the operator will be unavailable for the observation period).

## 6. Failure Modes

- **Update neglect.** The operator stops checking for updates. The system falls months or years behind. Eventually, updating requires a chain of intermediate updates or a full reinstallation. Mitigation: the monthly intelligence gathering requirement (R-D5-002-04) and the operational tempo in OPS-001.
- **Untested deployment.** Under time pressure, the operator skips the test system and applies updates directly to production. The update breaks something. Rollback is required but may be incomplete. Mitigation: R-D5-002-01 is a hard constraint. No exceptions.
- **Failed rollback.** An update causes problems, but the rollback mechanism fails -- the previous kernel was not preserved, the backup is corrupt, or the rollback procedure was never tested. Mitigation: R-D5-002-03 (multiple kernels) and the requirement to verify rollback capability before deploying.
- **Signature verification bypass.** The operator cannot verify a package signature and applies it anyway because it seems urgent. A compromised package enters the institution. Mitigation: R-D5-002-02 is absolute. If a signature cannot be verified, the package is not applied. Period.
- **Test-production divergence.** The test system drifts from production configuration. Updates pass testing but fail in production because the environments differ. Mitigation: quarterly configuration audits comparing test and production systems (see D5-006).

## 7. Recovery Procedures

**If an update has rendered a system unbootable:**
1. Boot from the previous kernel using the bootloader fallback mechanism.
2. If the previous kernel boots successfully, remove the problematic update and document the failure.
3. If no kernel boots, boot from the institution's emergency recovery media (a bootable USB or optical disk maintained per OPS-001, R-OPS-07).
4. Mount the system's root filesystem from recovery media.
5. Restore the system from the most recent known-good backup.
6. Document the incident, including the specific update that caused the failure, the symptoms, and the recovery steps taken.

**If an update has caused service failures without rendering the system unbootable:**
1. Identify the specific update that caused the failure using the operational log and package manager history.
2. Roll back the specific package to its previous version.
3. Verify that all services are restored.
4. Document the incident.
5. Investigate the root cause before re-attempting the update.

**If the system has fallen significantly behind on updates:**
1. Do not attempt to apply all missed updates at once.
2. Begin with security-critical patches, prioritized by the triage matrix.
3. Apply updates in small batches, testing each batch before proceeding.
4. If the system is more than one major release behind, evaluate whether a clean reinstallation from current media is more practical than incremental updating. This decision is a Tier 3 governance decision.
5. Document the catch-up plan and its execution.

## 8. Evolution Path

- **Years 0-2:** The update pipeline is being established. Expect to refine the test suite and verification checklists as you discover what matters and what does not. The intelligence gathering process will stabilize as you learn which sources are reliable.
- **Years 2-5:** The pipeline should be routine. Focus on ensuring the test system stays synchronized with production. Begin accumulating historical data on which types of updates cause problems.
- **Years 5-15:** You will face your first major OS upgrade decision. The framework in Section 4.6 will be tested. Document the experience thoroughly for future reference.
- **Years 15-30:** Multiple hardware generations will have passed. The OS may have been replaced entirely. The principles in this article -- test before deploy, verify signatures, maintain rollback capability, document everything -- should survive even if every specific procedure changes.
- **Years 30-50+:** A successor may be maintaining the system. The update pipeline must be documented well enough that they can follow it without oral instruction. The Commentary Section should contain decades of operational wisdom about what types of updates tend to cause problems and how to handle them.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The most dangerous moment in OS maintenance is not the update itself. It is the moment you convince yourself that "this one is safe" and skip the test system. Every failed production deployment in the history of system administration was preceded by that thought. The pipeline exists because human judgment about update safety is unreliable. Follow the pipeline.

The timelines in the triage matrix will feel slow to anyone accustomed to internet-connected systems where patches are applied within hours of release. Remember that the air gap itself is a significant mitigation for many vulnerability classes. A remote code execution vulnerability is substantially less threatening to a system with no network-reachable attack surface. The timelines reflect this reality. Urgency should be proportional to actual risk, not to the severity score alone.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 2: Integrity Over Convenience)
- CON-001 -- The Founding Mandate (air-gap mandate)
- SEC-001 -- Threat Model and Security Philosophy (supply chain threats, Category 3)
- OPS-001 -- Operations Philosophy (operational tempo, documentation-first principle)
- P-001 -- Platform Philosophy and Architectural Principles
- P-002 -- Software Acquisition and Verification
- P-003 -- Operating System Selection and Configuration
- P-012 -- Platform Evolution and Migration Planning
- Domain 18 -- Import & Quarantine (quarantine procedures for incoming media)
- Domain 12 -- Disaster Recovery (system restore procedures)

---
---

# D5-003 -- Core Service Configuration and Hardening

**Document ID:** D5-003
**Domain:** 5 -- Platform & Core Systems
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, SEC-001, OPS-001, P-001, P-003, P-004, P-005, D5-002
**Depended Upon By:** D5-004, D5-005, D5-006, all Domain 3 implementation articles, all higher-layer service configurations.

---

## 1. Purpose

This article defines the configuration and hardening requirements for every core service the institution operates. A core service is a service that other services depend on, that the operator interacts with during normal operation, or whose failure would degrade the institution's ability to function. Core services are the nervous system of the platform. If they are misconfigured, everything built on top of them is unreliable. If they are unhardened, every vulnerability in the stack above them is amplified.

This article serves three purposes simultaneously. First, it establishes configuration baselines -- the known-good state that each service should be in. Second, it provides hardening checklists -- the specific steps taken to reduce each service's attack surface and increase its resilience. Third, it defines verification procedures -- how to confirm that configurations have not drifted from their baselines over time.

Configuration drift is the silent killer of long-lived systems. A setting is changed during troubleshooting and never reverted. A default value is assumed to be correct and never verified. A hardening step is applied to one system but not another. Over months and years, the actual configuration diverges from the documented configuration until the documentation is fiction. This article exists to prevent that divergence.

## 2. Scope

This article covers the following core services:

- DNS (internal domain resolution)
- NTP (time synchronization)
- Syslog / logging infrastructure
- Authentication and authorization services
- SSH (secure shell for local administration)
- Firewall (host-level packet filtering)
- Cron / scheduled task management

This article does NOT cover:

- Network-level architecture and segmentation (see D5-004).
- Monitoring and alerting systems (see D5-005).
- Application-layer services built on top of these core services.
- The initial installation and selection of these services (see P-003 and P-004).

## 3. Background

The concept of "hardening" comes from military fortification: reducing the surface area that an attacker can target. In system administration, hardening means disabling unnecessary features, restricting access to necessary features, and configuring services to fail safely rather than fail open.

Hardening is often treated as a one-time activity: you harden a system when you deploy it and then forget about it. This is insufficient for a fifty-year institution. Hardening must be verified continuously because configurations drift, updates reset settings, and new features enable new attack surfaces. A hardened system that is not regularly verified is a system that was hardened once and may not be hardened now.

For an air-gapped institution, hardening serves a dual purpose. The obvious purpose is security against external threats -- which, given the air gap, are primarily supply-chain attacks and insider threats. The less obvious but equally important purpose is operational reliability. A hardened service is a service with fewer moving parts, fewer enabled features to malfunction, and fewer interactions to produce unexpected behavior. Hardening is simplification, and simplification is survival.

### 3.1 The Service Dependency Map

Understanding service dependencies is critical because a failure in one core service can cascade. The dependency map for the institution's core services is:

```
Authentication <-- SSH <-- Remote Administration
      |
      v
NTP (time sync) <-- Logging (timestamps) <-- All services (audit trail)
      |
      v
DNS (name resolution) <-- All services that resolve names
      |
      v
Firewall (host-level) <-- All network-facing services
      |
      v
Cron (scheduling) <-- Maintenance tasks, log rotation, health checks
```

A failure in NTP, for example, does not merely mean the clock is wrong. It means log timestamps are unreliable, certificate validation may fail (if timestamps are sufficiently wrong), and any time-based authentication tokens become suspect. Dependencies are multiplicative, not additive.

## 4. System Model

### 4.1 Configuration Baselines

A configuration baseline is a documented, version-controlled record of every configuration file and every significant setting for a service. The baseline is the "known-good state." Deviations from the baseline must be deliberate, documented, and reviewed.

For each core service, the baseline document must contain:

1. **Service identity:** Name, version, purpose within the institution.
2. **Configuration files:** Full path to every configuration file, with a cryptographic hash of the expected contents.
3. **Enabled features:** An explicit list of features that are enabled and why.
4. **Disabled features:** An explicit list of features that were deliberately disabled and why. This is as important as the enabled list, because it records decisions.
5. **Listening interfaces and ports:** Every network port the service listens on, and on which interfaces. If a service should not listen on a network at all, this must be stated.
6. **Permissions:** File permissions on configuration files, binary files, and data directories. User and group ownership.
7. **Resource limits:** Any CPU, memory, or file descriptor limits applied to the service.
8. **Logging configuration:** What the service logs, where it logs, and at what verbosity level.
9. **Dependencies:** What other services this service requires to function.
10. **Startup configuration:** How the service starts (systemd unit, init script, etc.), what happens if it fails to start, and whether it restarts automatically.

### 4.2 Hardening Checklists by Service

#### 4.2.1 DNS (Internal Resolution)

The institution runs an internal DNS server for name resolution within the air-gapped network. There is no external DNS. There are no upstream forwarders. The DNS server is authoritative for all internal zones and resolves nothing else.

Hardening checklist:
- Recursion is disabled or restricted to the local network only.
- Zone transfers are disabled or restricted to authorized secondary servers only.
- The DNS server does not listen on any external-facing interface (there should be none, but verify).
- DNSSEC is configured for internal zones if the DNS implementation supports it.
- The DNS cache is sized appropriately and cache poisoning mitigations are enabled.
- Version information is not disclosed in responses.
- The DNS process runs as a dedicated, unprivileged user.
- Configuration files are owned by root and readable only by the service user.
- Rate limiting is enabled to prevent resource exhaustion from misconfigured clients.

#### 4.2.2 NTP (Time Synchronization)

In an air-gapped environment, there are no external NTP servers. Time synchronization operates within the internal network. One system is designated as the authoritative time source (the stratum 1 server, synchronized to a local hardware clock or GPS receiver if available), and all other systems synchronize to it.

Hardening checklist:
- The NTP daemon listens only on internal interfaces.
- Authentication is enabled between NTP peers (symmetric key or autokey).
- The NTP daemon does not allow monlist or other amplification-prone queries from unauthorized sources.
- The stratum 1 server's local clock source is documented, and its drift characteristics are recorded.
- NTP clients are configured to reject time steps larger than a defined threshold (e.g., 1000 seconds) without manual confirmation, to prevent accidental or malicious clock jumps.
- The NTP process runs as a dedicated, unprivileged user.
- If no hardware clock source is available, the institution must document its time accuracy expectations and the implications for time-sensitive operations (log correlation, certificate validity).

#### 4.2.3 Syslog / Logging Infrastructure

Centralized logging is essential for security, troubleshooting, and auditing. All systems forward their logs to a central log server. The log server is a critical service whose integrity must be protected.

Hardening checklist:
- Log transport between systems is encrypted (TLS-protected syslog).
- The central log server accepts logs only from known, authorized systems.
- Log files are stored on a dedicated partition or volume to prevent log exhaustion from filling the root filesystem.
- Log rotation is configured with retention periods defined by institutional policy.
- Logs are append-only at the application level. The log server does not allow remote deletion or modification of logs.
- The log server process runs as a dedicated, unprivileged user.
- Log files are permission-restricted: writable only by the logging service, readable by authorized administrators.
- Log integrity verification is performed regularly (checksums or signed log chains).
- The logging infrastructure has its own logging -- a separate mechanism to detect if the primary logging system itself fails.

#### 4.2.4 Authentication and Authorization

The institution's authentication service provides identity verification for all human and service accounts. This may be implemented as a local directory service (e.g., LDAP), a centralized account database, or a carefully managed set of local accounts with synchronized configuration.

Hardening checklist:
- Accounts follow the principle of least privilege. No user has more access than required for their role.
- Service accounts are separate from human accounts and have only the permissions required for their specific function.
- Password policies enforce minimum complexity and length. Passphrases are preferred over passwords.
- Account lockout or rate-limiting is enabled for failed authentication attempts.
- All authentication events (successful and failed) are logged.
- Root or superuser access is controlled through sudo or equivalent, with full command logging.
- Direct root login is disabled on all services except the local console.
- Unused default accounts are disabled or removed.
- Authentication credentials are stored using strong, salted hashing (bcrypt, scrypt, or equivalent).
- Session timeouts are configured for all interactive services.

#### 4.2.5 SSH (Secure Shell)

SSH is the primary mechanism for remote administration within the internal network.

Hardening checklist:
- Password authentication is disabled. Key-based authentication only.
- Root login via SSH is disabled.
- SSH listens only on internal interfaces.
- Cipher suites are restricted to strong, current algorithms. Weak ciphers are explicitly disabled.
- Idle session timeout is configured. Maximum authentication attempts per connection is limited.
- SSH key sizes meet the institutional minimum (RSA 4096-bit or Ed25519).
- SSH host keys are inventoried and their fingerprints are documented for verification.
- Port forwarding and tunneling are disabled unless explicitly required and documented.
- The SSH daemon logs all connection events at a useful verbosity level.

#### 4.2.6 Firewall (Host-Level)

Every system runs a host-level firewall with a default deny-all policy for both inbound and outbound traffic.

Hardening checklist:
- Default inbound policy: DROP. Default outbound policy: DROP.
- Rules are defined by service, not by port number alone. Each rule is documented with its purpose.
- ICMP is permitted for diagnostics but rate-limited. Logging is enabled for dropped packets.
- Firewall rules are loaded at boot before any network services start.
- The firewall configuration is version-controlled with the configuration baselines.
- IPv6 is explicitly disabled if not used, or explicitly firewalled if used.
- Anti-spoofing rules are in place. The rule set is reviewed quarterly.

#### 4.2.7 Cron / Scheduled Tasks

Cron manages scheduled maintenance tasks, log rotation, health checks, and other periodic operations. Its reliability directly affects institutional maintenance.

Hardening checklist:
- Access to cron is restricted to authorized users only (cron.allow/cron.deny).
- Every cron job is documented in the configuration baseline with its purpose, schedule, and expected runtime.
- Cron job output is captured and directed to the logging infrastructure, not to local mailboxes or /dev/null.
- Cron jobs run with the minimum necessary privileges.
- Resource-intensive cron jobs are scheduled to avoid conflicts with each other and with peak usage periods.
- The cron daemon's own health is monitored (see D5-005).
- Failed cron jobs generate alerts.
- Cron job configurations are version-controlled.

### 4.3 Configuration Drift Detection

Configuration drift is detected through regular comparison of actual system state against documented baselines. The drift detection procedure runs as follows:

1. **Generate current state:** For each core service, extract the current configuration into a comparable format (configuration file hashes, running parameter values, listening ports).
2. **Compare against baseline:** Diff the current state against the documented baseline.
3. **Classify deviations:** Each deviation is classified as:
   - **Authorized:** A documented, approved change that needs to be incorporated into the baseline.
   - **Unauthorized but benign:** An undocumented change that does not affect security or functionality (e.g., a log file that has grown). The baseline should be updated.
   - **Unauthorized and concerning:** An undocumented change that could affect security or functionality. Investigation is required.
4. **Remediate:** Authorized deviations update the baseline. Benign deviations update the baseline with a note. Concerning deviations trigger an incident investigation.
5. **Record:** All drift detection results are recorded in the operational log, regardless of whether deviations were found.

Drift detection must be performed at least quarterly. It should be performed after any maintenance activity, OS update, or configuration change.

## 5. Rules & Constraints

- **R-D5-003-01:** Every core service must have a documented configuration baseline before it is placed into production use.
- **R-D5-003-02:** The hardening checklists in Section 4.2 are mandatory minimums. Additional hardening may be applied but the listed items must not be omitted.
- **R-D5-003-03:** Configuration drift detection must be performed at least quarterly and after any system change.
- **R-D5-003-04:** All configuration changes to core services must be documented in the operational log before the change is made (intent) and after the change is verified (result).
- **R-D5-003-05:** No core service may run as root unless technically required and documented. When root execution is necessary, the justification must be recorded in the configuration baseline.
- **R-D5-003-06:** The service dependency map must be maintained and updated whenever a new service is added or a dependency changes.
- **R-D5-003-07:** Configuration files for core services must be stored in version control (see D5-006).

## 6. Failure Modes

- **Configuration drift undetected.** Quarterly audits are skipped. The actual configuration diverges from documentation. When a system needs to be rebuilt, the documentation produces a system different from the one that was running. Mitigation: R-D5-003-03 and integration with the operational tempo in OPS-001.
- **Hardening regression.** A software update resets a hardening setting to its default. The regression is not detected because drift detection did not run after the update. Mitigation: drift detection after every update (D5-002 integration) and explicit post-update verification.
- **Dependency cascade failure.** A core service fails, and the operator does not understand the dependency chain, causing them to troubleshoot the wrong service. Mitigation: the dependency map in Section 3.1 and its maintenance requirement (R-D5-003-06).
- **Over-hardening.** A service is hardened so aggressively that it cannot function. Features required by legitimate operations are disabled. Mitigation: testing every hardening change before production deployment, and documenting the rationale for every disabled feature so it can be re-evaluated.
- **Documentation fiction.** The baselines are documented but never verified. They describe a system that existed once but no longer does. The documentation provides false confidence. Mitigation: drift detection is the cure. Baselines without verification are theories, not facts.

## 7. Recovery Procedures

**If a core service has been misconfigured:**
1. Compare the current configuration against the version-controlled baseline.
2. Identify the specific deviation.
3. If the baseline is known-good, restore the baseline configuration.
4. Restart the service.
5. Verify functionality using the service's health check.
6. Document the misconfiguration, its cause, and the restoration.

**If configuration drift has accumulated over a long period:**
1. Perform a full drift detection against all core service baselines.
2. Classify each deviation per Section 4.3.
3. For concerning deviations, investigate before remediating. The deviation may be hiding a deeper problem.
4. Remediate in small batches, verifying after each batch.
5. Update baselines to reflect the verified-good state.
6. Schedule a follow-up audit within 30 days to confirm stability.

**If a core service fails and the cause is unknown:**
1. Consult the service dependency map. Verify that all dependencies are healthy before investigating the failed service itself.
2. Check the service's logs for error messages.
3. Compare the current configuration against the baseline.
4. If the configuration matches the baseline and dependencies are healthy, the problem is likely a software bug or resource issue. Check resource utilization (disk, memory, CPU).
5. If the service cannot be restored, restore from the most recent known-good backup and apply changes incrementally until the problem recurs.

## 8. Evolution Path

- **Years 0-2:** Establish configuration baselines for all core services. The initial baselines will be refined frequently as operational experience reveals necessary adjustments. Expect the hardening checklists to grow as edge cases are discovered.
- **Years 2-5:** Baselines should be stable. Drift detection should be routine. The focus shifts to ensuring that updates (D5-002) do not regress hardening settings.
- **Years 5-15:** Core services may be replaced as better alternatives emerge. When a service is replaced, the new service receives a full baseline and hardening treatment before entering production. The old service's baseline is archived, not deleted.
- **Years 15-30:** The hardening landscape will have changed. Algorithms considered strong today may be weak. New attack vectors may have emerged. The hardening checklists must be reviewed against current best practices during annual reviews.
- **Years 30-50+:** The specific services listed in this article may no longer exist. But the principles -- document everything, harden everything, detect drift, maintain the dependency map -- must persist regardless of which specific software implements them.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The hardening checklists in this article will feel exhaustive to someone who has never had a system compromised and paranoid to someone who has. They are neither. They are the minimum. Every item in every checklist addresses a real attack vector or a real operational failure that has been documented in the broader system administration community. If an item seems unnecessary, you are fortunate enough not to have encountered the scenario it prevents. Do not confuse that fortune with evidence that the scenario is impossible.

The service dependency map may seem obvious when you first draw it. It will not seem obvious at 3 AM when a service is down and you cannot remember whether NTP failing could possibly affect authentication. Maintain the map. Consult the map. Update the map. It is a tool for your worst moments, not your best ones.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 3: Transparency of Operation)
- SEC-001 -- Threat Model and Security Philosophy (Pillar 2: Security Is a Property)
- OPS-001 -- Operations Philosophy (documentation-first principle)
- P-001 -- Platform Philosophy and Architectural Principles
- P-003 -- Operating System Selection and Configuration
- P-004 -- Network Services Architecture
- P-005 -- Identity and Access Management
- D5-002 -- Operating System Maintenance Procedures (update-triggered drift detection)
- D5-004 -- Network Architecture and Segmentation (network-level hardening context)
- D5-005 -- System Monitoring and Alerting (health check integration)
- D5-006 -- Configuration Management Without Internet (version control for configs)
- Domain 3 -- Security & Integrity (security requirements these configurations implement)

---
---

# D5-004 -- Network Architecture and Segmentation

**Document ID:** D5-004
**Domain:** 5 -- Platform & Core Systems
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, SEC-001, OPS-001, P-001, P-004, D5-003
**Depended Upon By:** D5-005, D5-006, all Domain 3 network security articles, all higher-layer service deployments.

---

## 1. Purpose

This article defines the internal network architecture of the institution. It specifies how systems are connected, how network traffic is segmented and controlled, what "air-gapped" means at every level of the network, and how the network is monitored for anomalies.

The word "air-gapped" is used frequently in this institution's documentation. Most people understand it as "not connected to the internet." That understanding is incomplete. A properly air-gapped network is not merely disconnected from the internet -- it is designed so that reconnection is physically difficult, so that the absence of external connectivity is a structural property rather than a configuration option. This article defines what that structural property looks like in practice, including every exception, every boundary, and every physical enforcement mechanism.

Additionally, this article defines the internal segmentation of the network. Even though all systems are behind the air gap, not all systems should be able to communicate with all other systems freely. Internal segmentation limits the blast radius of any compromise, contains misconfigurations, and enforces the principle of least privilege at the network level. A flat internal network is a convenience that becomes a liability the moment anything goes wrong.

## 2. Scope

This article covers:

- The physical network topology of the institution.
- The logical network topology, including VLANs and subnets.
- Firewall rules between network segments.
- The definition and enforcement of the air gap.
- Every authorized exception to the air gap and how those exceptions are controlled.
- Network monitoring procedures.
- Physical versus logical segmentation and when each is used.
- Inter-system communication policies.

This article does NOT cover:

- The configuration of individual services on the network (see D5-003).
- The physical cabling and hardware specifications (see Domain 4).
- Wireless networking (there is none; see Section 4.2).
- External communication of any kind (there is none; the absence is the point).

## 3. Background

Network architecture is the skeleton of any digital institution. It determines which systems can talk to which other systems, what traffic is visible to whom, and how failures in one part of the network affect other parts. In an internet-connected environment, network architecture must defend against external attackers, internal compromises, lateral movement, data exfiltration, and dozens of other threats. In an air-gapped environment, many of these threats are eliminated, but they are replaced by a different set of requirements.

The primary network threat in an air-gapped institution is not external penetration. It is internal misconfiguration. A wrong firewall rule that blocks a critical service. A VLAN misconfiguration that allows a test system to access production data. A network change that subtly breaks monitoring. These are the failures that air-gapped institutions face, and they are addressed through careful design, documentation, and verification rather than through complex intrusion detection systems.

### 3.1 Why Segment an Air-Gapped Network?

A common question: if the network is air-gapped, why bother with internal segmentation? The answer has three parts.

First, defense in depth. The air gap is one layer of protection. Internal segmentation is another. If the air gap is somehow compromised -- through a supply-chain attack on imported media, through a malicious device introduced during maintenance, or through a mistake -- internal segmentation limits the attacker's ability to move laterally.

Second, operational hygiene. Segmentation prevents operational accidents. A misconfigured service in the test segment cannot accidentally corrupt production data if the segments are properly isolated. A script running in the development environment cannot inadvertently access the backup server. Segmentation is as much about preventing mistakes as preventing attacks.

Third, monitoring clarity. In a flat network, all traffic looks the same. In a segmented network, traffic between segments is visible and auditable. A packet crossing from the test segment to the production segment is an event worth investigating. In a flat network, you would not even notice it.

## 4. System Model

### 4.1 Physical Network Topology

The institution's physical network uses wired Ethernet exclusively. Every connection between systems is a physical cable. There are no wireless access points, no Bluetooth adapters, no infrared ports. Wireless capability is a physical air-gap violation because it creates a communication channel that cannot be controlled by cable management alone.

The physical topology is a modified star: all systems connect to one or more managed switches, which connect to a central router or firewall appliance that enforces segmentation rules. The specific hardware is documented in Domain 4. This article addresses the logical design that runs on that hardware.

**Physical enforcement of the air gap:**
- No system in the institution has a network interface connected to any network outside the institution. This is not a software configuration. It is a physical absence: there is no cable, no port, no antenna.
- Network interfaces that are not in use are physically disabled. On systems where this is possible at the hardware level (BIOS/UEFI disabling, physical disconnection), hardware-level disabling is used. Software-level disabling is a secondary measure, not a primary one.
- The institution's network switches do not have uplink ports connected to any external network.
- There is no modem, no cellular adapter, no satellite link. The absence of these devices is verified during the annual security audit.

### 4.2 Wireless Policy

There is no wireless networking. WiFi capabilities are disabled at the BIOS/UEFI level where possible and via kernel module blacklisting as a secondary measure. Bluetooth is disabled at the hardware and software level on all systems. Any device entering the institution that has wireless capability must have that capability physically disabled or used under Domain 18 quarantine procedures.

Wireless signals cross air gaps. A WiFi-enabled device inside the air-gapped network is a potential bridge to any wireless network within range. The air gap requires the absence of wireless, not merely the disabling of it at the software level.

### 4.3 Logical Network Segmentation

The internal network is divided into segments using VLANs (Virtual Local Area Networks). Each segment serves a specific function and has defined rules for what traffic may enter and leave it.

**Segment 1: Core Services (VLAN 10)**
Purpose: Hosts the institution's core infrastructure services -- DNS, NTP, authentication, central logging.
Access policy: All other segments may send requests to Core Services on defined ports. Core Services may respond. Core Services may not initiate connections to other segments except for specific, documented purposes (e.g., the log server initiating a pull from a system that supports pull-based logging).
Systems in this segment: DNS server, NTP server, authentication server, central log server.

**Segment 2: Production Systems (VLAN 20)**
Purpose: Hosts the institution's primary operational systems -- the systems that store and process institutional data.
Access policy: Production systems may access Core Services. Production systems may access the Backup segment for backup operations. Production systems may not access the Test segment or the Administration segment directly.
Systems in this segment: Primary storage servers, database servers, application servers.

**Segment 3: Backup and Archive (VLAN 30)**
Purpose: Hosts backup servers and archive storage.
Access policy: Production systems may push data to the Backup segment. The Backup segment may pull data from Production. The Backup segment may not access any other segment. The Backup segment has no path to the Test segment or the Administration segment.
Systems in this segment: Backup servers, archive storage, tape libraries (if applicable).

**Segment 4: Test and Development (VLAN 40)**
Purpose: Hosts test systems, development environments, and the OS update test system (per D5-002).
Access policy: The Test segment may access Core Services for basic infrastructure (DNS, NTP, authentication). The Test segment may not access Production or Backup. Traffic from the Test segment to Production is explicitly blocked.
Systems in this segment: Test mirrors of production systems, development workstations, update testing systems.

**Segment 5: Administration (VLAN 50)**
Purpose: The operator's primary workstation and administrative tools. This is the segment from which the operator manages all other segments.
Access policy: The Administration segment may access all other segments via SSH and defined management ports. No other segment may initiate connections to the Administration segment. This is the most privileged segment.
Systems in this segment: The operator's administrative workstation, monitoring dashboards, configuration management tools.

**Segment 6: Quarantine (VLAN 60)**
Purpose: Isolated segment for processing incoming media and data per Domain 18.
Access policy: The Quarantine segment has no access to any other segment. Data leaves Quarantine only through a defined, manual transfer process (sneakernet) after clearing quarantine procedures. The Quarantine segment may have its own isolated Core Services (local DNS, NTP) or may function without them entirely.
Systems in this segment: Quarantine workstation, malware scanning systems.

### 4.4 Firewall Rules Architecture

Firewall rules implement the access policies defined in Section 4.3. Rules are maintained at two levels:

**Network-level (perimeter) firewall:** Implemented on the central router or firewall appliance. Controls inter-VLAN traffic. Default policy: deny all inter-VLAN traffic. Explicit rules permit only the traffic described in the segment access policies above.

**Host-level firewall:** Implemented on every individual system per D5-003, Section 4.2.6. Provides defense in depth. Even if the network-level firewall is misconfigured, the host-level firewall provides an additional layer of protection.

The firewall rule set is documented in a dedicated firewall rules document that is version-controlled alongside configuration baselines (D5-006). Every rule has:
- A unique identifier.
- Source segment and/or IP range.
- Destination segment and/or IP range.
- Protocol and port.
- Direction (inbound/outbound relative to the segment).
- Purpose: a human-readable explanation of why this rule exists.
- Date added and authorization reference (operational log entry or governance decision).

Rules without documented purposes are treated as unauthorized configurations and are removed.

### 4.5 The Air-Gap Exception Register

The air gap is absolute during normal operations. However, there are specific, controlled scenarios where the boundary between the institution and the outside world is crossed. These are not breaches of the air gap. They are controlled crossings, and they are enumerated here because an unenumerated crossing is, by definition, unauthorized.

**Exception 1: Software and Update Import.**
Media enters the institution carrying software updates, operating system patches, or new software packages. This crossing is governed by Domain 18 (Import & Quarantine) and D5-002 (OS Maintenance). The media never connects to an institutional network. It connects only to the Quarantine segment, which has no connectivity to other segments.

**Exception 2: Data Import.**
External data (documents, reference materials, archives) enters the institution through the same quarantine process. Again, media connects only to the Quarantine segment.

**Exception 3: Data Export.**
Institutional data may occasionally need to leave the institution (e.g., off-site backup rotation). Data export follows a defined procedure: data is written to media on a system in the Backup segment, the media is physically removed, and its removal is logged. The system that wrote the media has no external connectivity at any point during this process.

**Exception 4: Hardware Maintenance.**
Hardware leaving the institution undergoes data destruction per Domain 3. Hardware entering undergoes quarantine per Domain 18.

There are no other exceptions. Any new crossing type must be added to this register through a Tier 2 governance decision per GOV-001 before the crossing occurs.

### 4.6 Network Monitoring

Network monitoring in an air-gapped institution serves different purposes than in an internet-connected one. There is no intrusion detection for external attacks because there is no external attack surface. Instead, network monitoring focuses on:

- **Traffic pattern verification:** Is traffic flowing between authorized segments only? Is any traffic appearing on segments where it should not?
- **Bandwidth baseline monitoring:** Is traffic volume consistent with expected patterns? A sudden spike could indicate a misconfiguration or a compromised system generating unusual traffic.
- **Service availability:** Are the expected services on each segment reachable and responding?
- **ARP table monitoring:** Are the MAC addresses on each segment consistent with the known hardware inventory? An unexpected MAC address is a physical security concern.
- **VLAN integrity:** Are VLAN tags being respected? Is any traffic leaking between VLANs outside the firewall?

Network monitoring data is forwarded to the central log server in the Core Services segment and reviewed as part of the daily health check (D5-005).

## 5. Rules & Constraints

- **R-D5-004-01:** The air gap is enforced physically, not merely logically. No institutional system may possess an active network interface connected to any non-institutional network.
- **R-D5-004-02:** All wireless networking capabilities must be disabled at the hardware level on all institutional systems.
- **R-D5-004-03:** Inter-VLAN traffic is denied by default. Only explicitly documented and authorized rules permit traffic between segments.
- **R-D5-004-04:** Every firewall rule must have a documented purpose. Rules without documented purposes are unauthorized.
- **R-D5-004-05:** The Air-Gap Exception Register (Section 4.5) is the complete and authoritative list of authorized boundary crossings. Any crossing not on the register is a security incident.
- **R-D5-004-06:** The Quarantine segment must have no network connectivity to any other institutional segment.
- **R-D5-004-07:** Network architecture changes are Tier 2 governance decisions per GOV-001.
- **R-D5-004-08:** The network topology, VLAN assignments, and firewall rules must be reviewed at least annually and after any network change.

## 6. Failure Modes

- **VLAN misconfiguration.** A VLAN assignment error allows a system in one segment to access another segment's traffic. This can result in data leakage or unauthorized access. Mitigation: verify VLAN assignments after every network change. Test by attempting cross-segment access from each segment and confirming it is blocked where it should be.
- **Firewall rule creep.** Over time, temporary rules are added and never removed. The firewall becomes permissive by accumulation. Mitigation: annual firewall rule review. Every rule must be re-justified or removed.
- **Air-gap drift.** A wireless adapter is installed during a hardware replacement and not disabled. A modem is connected for "temporary" external access and never removed. Mitigation: the annual physical audit verifies the absence of unauthorized communication paths.
- **Switch failure.** A managed switch fails, disrupting multiple segments. Mitigation: spare switches in inventory. Documented switch configuration that can be applied to a replacement within hours.
- **Monitoring blind spot.** A segment is not monitored, allowing anomalies to go undetected. Mitigation: the monitoring configuration must cover all segments. The monitoring system's own coverage is verified during quarterly reviews (D5-005).
- **Quarantine bypass.** Media is connected directly to a non-quarantine system, bypassing the quarantine process. Mitigation: physical and procedural controls. USB ports on non-quarantine systems may be physically disabled or access-controlled.

## 7. Recovery Procedures

**If a VLAN misconfiguration is discovered:**
1. Immediately correct the VLAN assignment on the affected switch port.
2. Review traffic logs to determine whether unauthorized traffic crossed segment boundaries during the misconfiguration.
3. If unauthorized traffic occurred, treat this as a security incident per SEC-001.
4. Document the misconfiguration, its duration, its impact assessment, and the correction.

**If a firewall rule is found to be overly permissive:**
1. Tighten the rule immediately to the minimum required permission.
2. Review logs for any traffic that exploited the permissive rule.
3. Update the firewall rule documentation.
4. Conduct a full firewall rule review to check for similar issues.

**If an unauthorized communication path is discovered (air-gap violation):**
1. Physically disconnect the unauthorized path immediately.
2. Treat this as a security incident. Assume all systems reachable through the unauthorized path may be compromised.
3. Conduct a full forensic review of affected systems per SEC-001.
4. Document the violation, the response, and the remediation.
5. Implement physical controls to prevent recurrence.

**If network connectivity is lost between segments:**
1. Check physical connectivity: cables, switch ports, switch power.
2. Check switch configuration: VLAN assignments, port status.
3. Check the central firewall: rule set, interface status.
4. If hardware has failed, replace with spare hardware and apply the documented configuration from version control.
5. Verify full connectivity using the network verification checklist.
6. Document the outage, its cause, its duration, and the resolution.

## 8. Evolution Path

- **Years 0-2:** The network is being built for the first time. Expect VLAN assignments and firewall rules to be revised frequently as services are deployed and operational requirements become clear. Document every change. The initial rule set will be the foundation for all future changes.
- **Years 2-5:** The network architecture should be stable. Changes should be infrequent and governed by Tier 2 decisions. The focus shifts to monitoring and verification: confirming that the network still matches its documentation.
- **Years 5-15:** Network hardware will need replacement. When replacing switches or the firewall appliance, the documented configuration is applied to the new hardware and verified. This is the first real test of whether the documentation is sufficient.
- **Years 15-30:** Networking technology may have evolved significantly. New standards, new hardware capabilities, new segmentation mechanisms. The principles in this article -- segment by function, deny by default, enforce physically, document everything -- must guide the adoption of new technology.
- **Years 30-50+:** The physical network may look nothing like what exists today. But the logical architecture -- separate segments for separate functions, controlled boundaries between them, physical enforcement of the air gap -- should be recognizable. Technology changes. Architecture principles endure.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The Quarantine segment deserves special attention. It is tempting to give it "read-only" access to some internal service -- maybe DNS, for convenience during quarantine operations. Resist this temptation. The Quarantine segment is where untrusted data enters the institution. Any network path from Quarantine to other segments is a path that compromised data could travel. If the Quarantine workstation needs DNS, give it a local hosts file. If it needs time synchronization, set the clock manually. Inconvenience in the Quarantine segment is a feature, not a bug.

The six-segment architecture described here may seem elaborate for a small institution. It is. But segments are cheap. They cost nothing once the managed switch is in place. What they provide -- isolation, monitoring clarity, blast radius containment -- is worth far more than the modest effort of maintaining VLAN assignments. Do not be tempted to collapse segments "for simplicity." Flat networks are simple in the way that unlocked doors are simple: the simplicity saves effort until the moment it costs everything.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 1: Sovereignty of the Individual)
- CON-001 -- The Founding Mandate (air-gap mandate, Section 3.3)
- SEC-001 -- Threat Model and Security Philosophy (trust model, physical threats)
- OPS-001 -- Operations Philosophy (operational tempo, maintenance routines)
- P-001 -- Platform Philosophy and Architectural Principles
- P-004 -- Network Services Architecture
- D5-003 -- Core Service Configuration and Hardening (host-level firewall, service hardening)
- D5-005 -- System Monitoring and Alerting (network monitoring integration)
- Domain 4 -- Infrastructure & Power (physical network hardware)
- Domain 18 -- Import & Quarantine (quarantine procedures)

---
---

# D5-005 -- System Monitoring and Alerting

**Document ID:** D5-005
**Domain:** 5 -- Platform & Core Systems
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, OPS-001, P-001, P-008, D5-003, D5-004
**Depended Upon By:** D5-006, all Domain 10 daily operations articles, all Domain 12 disaster recovery trigger articles.

---

## 1. Purpose

This article defines what the institution monitors, how it monitors, what thresholds trigger alerts, and what the operator does with monitoring data. Monitoring is how the institution knows it is alive. Without monitoring, the operator is blind -- unable to distinguish a healthy system from one that is silently failing.

The distinction between monitoring and alerting is important. Monitoring is the continuous observation of system state. Alerting is the notification that a monitored value has crossed a threshold that requires attention. Monitoring without alerting produces data that nobody looks at. Alerting without monitoring produces noise that cannot be investigated. Both must exist and both must be calibrated.

In an air-gapped institution, monitoring has a particular character. There are no external monitoring services to call home to. There is no cloud dashboard. There are no pager integrations. The monitoring system runs entirely within the institution, and alerts are delivered to the operator through mechanisms that exist within the institution. This means the operator must be physically present or following a defined schedule to receive alerts. Monitoring is useful only if someone is looking.

This article also defines the daily health check -- the specific, step-by-step procedure the operator follows every day to verify that the institution is functioning normally. The daily health check is not optional. It is the first and most important line of defense against silent failures.

## 2. Scope

This article covers:

- What metrics and log data are monitored across all institutional systems.
- How monitoring data is collected, stored, and retained.
- Alert thresholds for each monitored metric.
- Alert classification (critical, warning, informational).
- Dashboard design for the operator.
- The daily health check procedure.
- What "normal" looks like -- baseline definitions for all monitored values.
- Log aggregation for offline systems.

This article does NOT cover:

- The configuration of individual services that produce log data (see D5-003).
- Network-level monitoring architecture (see D5-004, Section 4.6; this article covers the consumption of that data).
- The response to specific types of incidents (see SEC-001 and Domain 12 for incident response and disaster recovery).
- Physical environment monitoring -- temperature, humidity, power (see Domain 4).

## 3. Background

Monitoring is the institutional sense of sight. A system without monitoring is a system operating in the dark. The operator may believe everything is fine because nothing has obviously failed, but "nothing has obviously failed" is not the same as "everything is working correctly." Disks fill silently. Memory leaks develop over weeks. Certificate expiration dates approach. Backup jobs fail and nobody notices because the failure generates a log entry that nobody reads.

The history of system failures is largely a history of monitoring failures. Not failures of the monitoring system itself, but failures of attention -- monitoring data that was collected but not reviewed, alerts that were generated but not investigated, dashboards that existed but were not consulted. This article addresses all three: what to collect, what to alert on, and what to look at.

### 3.1 Monitoring in an Air-Gapped Environment

An air-gapped institution has no cloud dashboards, no mobile alerts, no automated remediation. This is a simplification. The monitoring system needs to do exactly four things: collect data from all systems, store data reliably, present data clearly, and alert when something is wrong. It serves one person -- the operator -- and must be simple enough to understand, maintain, and trust.

## 4. System Model

### 4.1 What to Monitor

Every institutional system is monitored across four categories:

**Category 1: Resource Utilization**
- CPU usage: overall percentage, per-core breakdown, load average (1, 5, 15 minute).
- Memory: total, used, available, swap usage.
- Disk: usage percentage per mount point, I/O rates, I/O latency.
- Network: bandwidth in/out per interface, packet error rates, packet drop rates.

**Category 2: Service Health**
- Service status: running, stopped, failed, for every core service defined in D5-003.
- Service responsiveness: time to respond to a health check probe.
- Service-specific metrics: DNS query rate and resolution time, NTP offset, authentication success/failure rates, log ingestion rate.

**Category 3: Security Indicators**
- Authentication failures: count per time period, source, account.
- Firewall drops: count per time period, source/destination.
- File integrity: changes to critical system files (see D5-003 configuration baselines).
- Privileged operations: sudo usage, root logins, service account activity.

**Category 4: System Integrity**
- SMART data from all disks: reallocated sectors, pending sectors, temperature, power-on hours.
- Filesystem integrity: results of periodic fsck or equivalent.
- Backup status: last successful backup time, backup size, verification result.
- Time synchronization: NTP offset and jitter from the authoritative source.
- Certificate expiration dates for any internal TLS certificates.

### 4.2 How to Monitor

Monitoring data is collected through a three-layer system:

**Layer 1: Agent-based collection.** A lightweight monitoring agent runs on every system. The agent collects resource utilization metrics, service health status, and security indicators at configurable intervals (default: 60 seconds for resource metrics, 300 seconds for service health). The agent pushes data to the central monitoring server in the Core Services segment (VLAN 10).

**Layer 2: Log aggregation.** System logs and application logs are forwarded to the central log server via the logging infrastructure described in D5-003. The monitoring system reads from the log server to identify log-based alerts (error messages, authentication failures, etc.).

**Layer 3: Active probing.** The monitoring server actively probes services at defined intervals to verify they are responsive. This is distinct from agent-reported status: the agent says "I think the service is running," while the probe verifies "I can actually reach the service and get a correct response."

Monitoring data is stored on the central monitoring server. Retention: high-resolution (60-second intervals) for 7 days, medium-resolution (5-minute averages) for 90 days, low-resolution (1-hour averages) for 2 years, alert history indefinitely. Data older than its retention period is archived (not deleted) during the first five years; the policy may be revised thereafter based on storage constraints.

### 4.3 Alert Thresholds and Classification

Alerts are classified into three levels:

**Critical:** Requires immediate attention. The system is degraded or at risk of imminent failure. The operator should investigate within one hour of the next scheduled check, or immediately if present.

**Warning:** Requires attention at the next maintenance window. The system is functional but trending toward a problem. The operator should investigate within 24 hours.

**Informational:** No action required. Recorded for trend analysis and audit purposes. Reviewed during weekly operations.

Default alert thresholds (these are starting values; operational experience will refine them):

| Metric | Warning | Critical |
|---|---|---|
| CPU load average (15 min) | > 80% of core count | > 95% of core count |
| Memory usage | > 80% | > 95% |
| Swap usage | > 20% | > 50% |
| Disk usage (any mount) | > 75% | > 90% |
| Disk I/O latency | > 50ms average | > 200ms average |
| NTP offset | > 100ms | > 1 second |
| Service health check failure | 1 consecutive failure | 3 consecutive failures |
| Authentication failures | > 5 in 10 minutes | > 20 in 10 minutes |
| SMART reallocated sectors | Any non-zero increase | > 10 increase in 30 days |
| Backup age | > 36 hours since last success | > 72 hours since last success |
| Certificate expiration | < 30 days | < 7 days |
| Disk error rate (read/write) | Any non-zero | > 10 per hour |
| Firewall drops (unexpected) | > 100 per hour from single source | > 1000 per hour from single source |

These thresholds must be calibrated against actual baselines. A system that normally runs at 70% CPU should have different thresholds than one that normally runs at 10%. Baseline calibration is part of the initial deployment process and is recalibrated quarterly.

### 4.4 Dashboard Design

The monitoring dashboard is the operator's primary interface to the institution's health. It must be designed for rapid comprehension -- the operator should be able to assess overall health in under 30 seconds.

The dashboard is organized in three tiers:

**Tier 1: Overview (the "traffic light" view).**
One line per system. Each system shows: name, segment, overall status (green/yellow/red). Green means all monitored values are within normal range. Yellow means at least one warning-level alert is active. Red means at least one critical alert is active. This view answers the question: "Is everything OK?" in under ten seconds.

**Tier 2: System detail.**
Selected by clicking or navigating from the overview. Shows all monitored metrics for a single system, with current values, 24-hour trend graphs, and active alerts. This view answers the question: "What is wrong with this system?"

**Tier 3: Metric detail.**
Selected from the system detail view. Shows the full history of a single metric with configurable time range, threshold lines, and annotations for maintenance events. This view answers the question: "How has this metric behaved over time?"

The dashboard must be accessible from the Administration segment (VLAN 50). It should be usable from a text-based terminal if the graphical interface is unavailable. The monitoring system must not require an external browser or framework that complicates the institutional stack.

### 4.5 The Daily Health Check Procedure

This is the most important procedure in this article. It is performed every day, takes 15-30 minutes, and is the operator's primary mechanism for detecting problems early. The procedure is performed from the operator's workstation in the Administration segment.

**Step 1: Dashboard Review (5 minutes).**
Open the monitoring dashboard. Review the Tier 1 overview. Confirm that all systems show green. If any system shows yellow or red, note the system and proceed to Step 2 with that system as a priority.

**Step 2: Alert Review (5 minutes).**
Review all alerts generated since the last daily check. For each alert:
- If it is a critical alert that has not been resolved, investigate immediately after completing the daily check.
- If it is a warning alert, note it for investigation during this check or the next maintenance window.
- If it is informational, note any unusual patterns.

**Step 3: Log Scan (5 minutes).**
Review the last 24 hours of centralized logs for error-level messages. You are not reading every log line. You are scanning for patterns: repeated errors, new error types, errors from systems that do not normally produce errors. The log aggregation system should provide a summary view of error counts by system and by message type.

**Step 4: Resource Trend Review (3 minutes).**
Check 7-day trends for disk usage across all systems. Disks are the most common resource to silently exhaust. Confirm that no disk is on a trajectory to hit 90% within the next 30 days at the current growth rate. Check memory trends for any system showing a consistent upward trend that could indicate a memory leak.

**Step 5: Backup Verification (2 minutes).**
Confirm that the most recent backup job for every production system completed successfully. Check the backup age metric. If any backup is older than expected, note it for investigation.

**Step 6: Time Synchronization (1 minute).**
Verify that all systems show NTP offsets within the acceptable range. Time drift affects log correlation, certificate validation, and scheduled task accuracy.

**Step 7: Record (2 minutes).**
Record the daily check results in the operational log. If everything is normal, the entry may be brief: "Daily check completed. All systems nominal. No active alerts." If anomalies were found, record them with the planned response.

### 4.6 What "Normal" Looks Like

"Normal" is defined by baselines. A baseline is the range of values that a metric exhibits during typical, healthy operation. Baselines must be established during initial deployment by observing each system under typical workload for at least two weeks.

Baselines are documented for every monitored metric on every system. The baseline document for each system contains:
- The metric name.
- The typical range (e.g., CPU load average: 0.2-0.8).
- The typical pattern (e.g., CPU spikes to 2.0 during the daily backup window from 02:00-03:00).
- Known exceptions (e.g., disk usage increases by approximately 500 MB per week due to log accumulation).

Baselines are recalibrated:
- After any significant configuration change.
- After any hardware change.
- After adding or removing a major workload.
- At least annually, during the annual operations review.

A value outside its baseline is not necessarily a problem. But a value outside its baseline that is not explained by a known exception is an anomaly worth investigating.

## 5. Rules & Constraints

- **R-D5-005-01:** Every institutional system must be monitored. No system may operate unmonitored.
- **R-D5-005-02:** The daily health check procedure (Section 4.5) is mandatory. It is performed every day the operator is present. If the operator is absent for more than 48 hours, the first task upon return is to perform the health check and review all alerts generated during the absence.
- **R-D5-005-03:** Critical alerts must be investigated within one hour of detection. Warning alerts must be investigated within 24 hours. Informational alerts are reviewed during weekly operations.
- **R-D5-005-04:** Alert thresholds must be calibrated against actual baselines, not default values alone. Default values are starting points, not final values.
- **R-D5-005-05:** Monitoring data retention must follow the policy in Section 4.2. Alert history is retained indefinitely.
- **R-D5-005-06:** The monitoring system itself must be monitored. A meta-monitoring check verifies that the monitoring agent is reporting from every system and that the monitoring server is storing and processing data correctly.
- **R-D5-005-07:** Baselines must be documented for every monitored metric and recalibrated at least annually.

## 6. Failure Modes

- **Alert fatigue.** Too many alerts are generated, most of them unactionable. The operator begins ignoring alerts. A genuine critical alert is missed because it is buried in noise. Mitigation: threshold calibration based on actual baselines. If an alert fires frequently without requiring action, the threshold needs adjustment, not the operator's patience.
- **Silent monitoring failure.** The monitoring agent on a system crashes or stops reporting. If the monitoring system does not detect this absence, the system operates unmonitored without anyone knowing. Mitigation: the meta-monitoring check (R-D5-005-06) detects agent failures. A system that has not reported in more than twice its expected interval generates a critical alert.
- **Dashboard neglect.** The dashboard exists but the operator stops looking at it. The daily health check is replaced by a quick assumption that "everything is probably fine." Mitigation: the daily check is a procedural requirement, not a suggestion. The operational log records that it was performed.
- **Baseline staleness.** Baselines are established once and never updated. As the system's workload changes, the baselines become irrelevant. Alerts fire for new-normal values. The operator adjusts by loosening thresholds until they are meaningless. Mitigation: annual baseline recalibration and recalibration after significant changes.
- **Storage exhaustion of monitoring data.** The monitoring system's own storage fills up, and it begins losing data or failing entirely. Mitigation: monitoring data retention is managed by the retention policy. The monitoring server's own disk usage is a monitored metric with appropriate thresholds.
- **Cargo cult checking.** The operator performs the daily check but does not actually engage with the data. They go through the motions, record "all nominal" in the log, and leave. Mitigation: the daily check procedure includes specific items to look at (Step 4: disk trends, Step 6: NTP offsets). These items require reading actual values, not merely confirming the dashboard is green.

## 7. Recovery Procedures

**If the monitoring system itself has failed:**
1. The monitoring server failure should be detected by the meta-monitoring mechanism (if the monitoring server is also the meta-monitor, this creates a gap -- see Commentary Section).
2. Verify the monitoring server is operational. Check its system logs, disk space, and service status.
3. If the monitoring server has crashed, restart it. Verify data collection resumes.
4. If monitoring data has been lost, document the gap. Manual health checks of all systems should be performed to cover the unmonitored period.
5. Investigate the cause of the monitoring failure to prevent recurrence.

**If alert fatigue has developed:**
1. Review all active alerts. Classify each as: genuine alert requiring action, recurring nuisance alert, or miscalibrated threshold.
2. For nuisance alerts, adjust the threshold or investigate the underlying cause. A threshold that fires constantly is either set wrong or revealing a chronic problem.
3. For miscalibrated thresholds, recalibrate against current baselines.
4. Document the recalibration in the operational log.

**If a system has been operating unmonitored:**
1. Install or restore the monitoring agent.
2. Perform a manual health check of the system: check resource utilization, service status, log files, and disk health.
3. Verify the system's configuration against its baseline (D5-003, D5-006).
4. Establish or re-establish baselines for the system.
5. Document the period of unmonitored operation and the results of the manual health check.

## 8. Evolution Path

- **Years 0-2:** Baselines are being established. Alert thresholds are being calibrated. The daily check procedure will be refined as the operator learns what matters and what is noise. Expect significant threshold adjustment during this period.
- **Years 2-5:** The monitoring system should be stable and well-calibrated. The daily check should feel routine. Trending data begins to reveal long-term patterns: seasonal disk growth, hardware degradation curves, workload shifts.
- **Years 5-15:** SMART data will begin to predict disk failures. Memory and CPU trends will inform capacity planning. The monitoring data archive becomes a valuable historical record of institutional health.
- **Years 15-30:** The monitoring tools themselves may need replacement. When replacing the monitoring system, preserve the alert history and as much historical data as feasible. The principles -- know what normal looks like, alert on deviations, check every day -- remain constant.
- **Years 30-50+:** A successor will inherit the monitoring system. The daily check procedure and the dashboard must be comprehensible without oral instruction. The baselines document and the Commentary Section provide the context the successor needs.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
There is a meta-monitoring problem that this article must be honest about: if the monitoring server is the only system that performs monitoring, then the monitoring server itself is a single point of failure for the entire monitoring function. The meta-monitoring check (R-D5-005-06) helps, but if the monitoring server is completely offline, who detects that?

The practical answer for a single-operator institution is the daily health check. The operator physically looks at the dashboard every day. If the dashboard is down, the operator knows immediately. This is not a sophisticated solution. It is a human solution. And for a single-operator institution, the human is the ultimate monitoring system. The daily check is not just about looking at data. It is about confirming that the systems that collect and present data are themselves alive.

The alert thresholds in Section 4.3 are deliberately conservative. It is better to start with thresholds that alert too often and tighten them over time than to start with thresholds that never alert and discover their inadequacy during a failure. Expect the first few months to require frequent threshold adjustment. This is normal. It is the calibration process, not a flaw.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 6: Honest Accounting of Limitations)
- OPS-001 -- Operations Philosophy (daily operations, Section 4.1)
- P-001 -- Platform Philosophy and Architectural Principles
- P-008 -- System Monitoring and Health
- D5-003 -- Core Service Configuration and Hardening (service health checks, logging infrastructure)
- D5-004 -- Network Architecture and Segmentation (network monitoring, Section 4.6)
- Domain 4 -- Infrastructure & Power (physical environment monitoring)
- Domain 12 -- Disaster Recovery (alert-triggered recovery procedures)

---
---

# D5-006 -- Configuration Management Without Internet

**Document ID:** D5-006
**Domain:** 5 -- Platform & Core Systems
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, OPS-001, P-001, P-009, D5-002, D5-003, D5-004, D5-005
**Depended Upon By:** All Domain 5 operational articles, all Domain 10 change management articles, all Domain 12 system rebuild articles.

---

## 1. Purpose

This article defines how the institution manages system configurations without access to internet-connected tools, cloud-based configuration management services, remote package repositories, or any of the other infrastructure that modern configuration management takes for granted. Configuration management is the discipline of ensuring that every system in the institution is in a known, documented, reproducible state at all times.

The phrase "configuration management" is heavily loaded in modern IT. It typically refers to tools like Ansible, Puppet, Chef, or SaltStack -- tools that pull configurations from central servers, install packages from remote repositories, and enforce desired state through automated agents. These tools assume internet connectivity, or at minimum, connectivity to centralized infrastructure that itself has internet access. An air-gapped institution has neither.

This does not mean configuration management is impossible without internet. It means it looks different. The principles are the same: document the desired state, detect deviations from the desired state, correct deviations, and record changes. The mechanisms are different: version control instead of cloud repositories, manual application instead of automated agents (or locally-hosted automation), and rigorous change procedures instead of "push a commit and let the tool handle it."

This article exists because unmanaged configurations are the primary source of operational unpredictability in long-lived systems. A system whose configuration is not documented, not version-controlled, and not audited is a system that will, over time, become incomprehensible to its operator. And an incomprehensible system is a system that cannot be maintained, cannot be recovered from failure, and cannot be transferred to a successor.

## 2. Scope

This article covers:

- The version control system for configurations.
- Change management procedures for all configuration changes.
- The configuration audit process.
- How configurations are documented for reproducibility.
- How configuration management interacts with the other operational articles (D5-002 through D5-005).
- How to manage configurations across multiple systems without centralized automation.

This article does NOT cover:

- The specific configuration baselines for core services (see D5-003).
- The process for acquiring new software (see P-002 and D5-002).
- Data management and backup procedures (see P-007 and Domain 12).
- The selection of version control software (see P-003 and P-009).

## 3. Background

Configuration management is the art of knowing what your systems look like. Not what you think they look like. Not what they looked like when you set them up. What they look like right now, this moment, and why they look that way.

In a small, air-gapped institution, the temptation is to manage configurations informally. There are only a handful of systems. The operator built them all. Surely they remember how everything is configured? This temptation is the single greatest threat to long-term operational integrity. The operator does not remember. They remember the broad strokes. They forget the detail that took two hours to troubleshoot. They forget the setting that was changed six months ago to work around a bug that has since been fixed. They forget the configuration file that was copied from a forum post and never fully understood.

The broader system administration community learned this lesson decades ago, which is why configuration management tools exist. But those tools were built for internet-connected environments with teams of administrators. This article adapts the principles for a single operator in an air-gapped environment.

### 3.1 The Configuration Lifecycle

Every configuration in the institution goes through a lifecycle:

1. **Creation:** A configuration is written for the first time, usually during initial system deployment.
2. **Documentation:** The configuration is recorded in version control with comments explaining its purpose and rationale.
3. **Deployment:** The configuration is applied to the target system.
4. **Verification:** The applied configuration is verified to match the documented version.
5. **Operation:** The system operates with this configuration.
6. **Change:** A change is needed. The change goes through the change management process (Section 4.2).
7. **Audit:** The actual configuration is compared against the documented version to detect drift (see D5-003, Section 4.3).
8. **Retirement:** The configuration is no longer needed (system decommissioned, service replaced). The configuration is archived, not deleted.

Every configuration, at every point in this lifecycle, has a documented state in version control. If a configuration exists on a system but not in version control, that configuration is unauthorized -- even if it works correctly and was applied by the operator.

## 4. System Model

### 4.1 The Version Control System

All configurations are stored in a local version control system. The specific tool is selected per P-009, but the requirements are:

- **Local operation:** The version control system must function entirely within the air-gapped network. No external server, no cloud service, no remote repository.
- **Full history:** Every change to every configuration must be preserved with the date, the author, and the reason for the change.
- **Branching and merging:** The ability to maintain separate branches for testing changes before committing them to the production branch.
- **Diffing:** The ability to compare any two versions of a configuration and see exactly what changed.
- **Integrity:** The repository must be checksummed or cryptographically verifiable to detect corruption.

The version control repository is hosted on a system in the Core Services segment (VLAN 10) or the Administration segment (VLAN 50). It is backed up as part of the institutional backup strategy. The repository is itself a critical asset -- losing version control history is a significant operational loss.

**Repository structure:**

```
/configs/
  /systems/
    /[hostname]/
      /etc/                  -- System configuration files
      /services/             -- Service-specific configuration
      /firewall/             -- Host firewall rules
      /cron/                 -- Scheduled tasks
      /baseline.sha256       -- Hash file for baseline verification
      /README                -- System description, role, dependencies
  /network/
    /vlans/                  -- VLAN definitions
    /firewall-rules/         -- Network-level firewall rules
    /switch-configs/         -- Managed switch configurations
    /topology.txt            -- Current network topology description
  /shared/
    /templates/              -- Configuration templates used across systems
    /scripts/                -- Deployment and verification scripts
    /policies/               -- Configuration policies (password policy, etc.)
  /archive/
    /decommissioned/         -- Configurations of retired systems
    /superseded/             -- Previous versions of replaced configurations
```

### 4.2 The Change Management Process

Every configuration change -- every change, no matter how small -- follows this process:

**Step 1: Propose.**
The operator writes a brief change proposal in the operational log. The proposal includes:
- What is being changed.
- Why it is being changed.
- What systems are affected.
- What the expected impact is.
- What the rollback plan is if the change causes problems.

For Tier 4 operational decisions (per GOV-001), the proposal may be a single paragraph. For Tier 3 policy decisions, the proposal must include a waiting period per GOV-001. For Tier 2 architectural decisions, the full governance process applies.

**Step 2: Implement in version control.**
The change is made to the configuration file in version control first, not on the target system. The commit message includes the operational log reference and the rationale for the change.

**Step 3: Test.**
If a test system is available (and for OS-level and core service changes, one must be per D5-002), the change is applied to the test system first and verified.

**Step 4: Deploy.**
The changed configuration is applied to the production system. Application may be manual (copying the file and restarting the service) or via a locally-hosted deployment script. The deployment method is documented.

**Step 5: Verify.**
After deployment, the operator verifies that:
- The configuration on the target system matches the version in version control (hash comparison).
- The affected service is functioning correctly.
- No unintended side effects are observed.

**Step 6: Record.**
The operator records the change completion in the operational log, including the verification results. If the change required a rollback, that is recorded as well.

### 4.3 Manual vs. Automated Deployment

In an internet-connected environment, configuration management tools automate deployment: the operator commits a change, and the tool propagates it to all affected systems. In an air-gapped institution, automation is limited but not impossible.

**Option A: Manual deployment.**
The operator copies the configuration file to the target system and applies it. This is the simplest approach and the least error-prone for a small number of systems. The risk is inconsistency: when the same change must be applied to multiple systems, manual deployment may result in slightly different results.

**Option B: Local automation.**
A locally-hosted automation tool (a simple shell script, a locally-running instance of a configuration management tool, or a custom deployment script) reads the desired configuration from version control and applies it to target systems over the internal network. This is more complex to set up but reduces the risk of inconsistency across systems.

The institution's preference is: use manual deployment until the number of systems makes it impractical, then introduce local automation. Local automation is subject to the same version control and change management requirements as any other configuration. The automation tool's own configuration is version-controlled.

When using manual deployment, the operator must use a checklist for multi-system changes:

1. List all systems that require the change.
2. Apply the change to each system, in order, checking off each one.
3. Verify the change on each system using hash comparison.
4. Record completion for each system in the operational log.

### 4.4 The Configuration Audit

The configuration audit is a comprehensive comparison of every production system's actual configuration against its documented, version-controlled configuration. It is the institutional defense against configuration drift and undocumented changes.

**Audit frequency:** Quarterly (integrated with the quarterly operations cycle in OPS-001).

**Audit procedure:**

1. **Generate inventory.** List every system in the institution and its role.
2. **For each system:**
   a. Generate cryptographic hashes of all configuration files on the live system.
   b. Compare these hashes against the baseline hashes in version control.
   c. For any mismatched hash, generate a diff showing the exact changes.
   d. Classify each deviation per D5-003, Section 4.3 (authorized, benign, concerning).
3. **For the network:**
   a. Dump the current configuration from all managed switches and the firewall appliance.
   b. Compare against the version-controlled network configurations.
   c. Classify deviations.
4. **For version control itself:**
   a. Verify the integrity of the version control repository (run an integrity check).
   b. Review the commit log since the last audit. Do all commits have adequate descriptions and operational log references?
5. **Compile the audit report.** The report lists every deviation found, its classification, and the remediation action taken. The report is stored in the operational log and the version control repository.
6. **Remediate.** Authorized deviations update the baseline. Benign deviations update the baseline with documentation. Concerning deviations trigger investigation and, if necessary, a security incident per SEC-001.

### 4.5 Configuration Reproducibility

The ultimate test of configuration management is reproducibility: can a system be rebuilt from documentation alone, without access to the original system?

Every system's configuration must be documented well enough to pass this test. The documentation includes:

1. **The base platform:** Which OS version, which installation profile, which packages. This information is version-controlled and maintained even as the OS is updated (D5-002).
2. **The configuration files:** Every non-default configuration file, in version control.
3. **The deployment order:** In what order must configurations be applied? Are there dependencies (e.g., "create the user before setting the file ownership")?
4. **The verification steps:** How does the operator confirm that the rebuilt system matches the original?
5. **The service startup order:** In what order must services be started, and what must be verified between each?

Reproducibility is tested during the annual operations cycle. The test: select one non-critical system, rebuild it from documentation alone, and verify it matches the original. This test reveals documentation gaps that no amount of reading can find.

### 4.6 Interaction with Other Operational Articles

Configuration management interacts with every other operational article in this document:

- **D5-002 (OS Maintenance):** OS updates may change configuration files. After any OS update, drift detection must be run to identify configuration changes introduced by the update. Updated configurations are committed to version control with the update as the rationale.
- **D5-003 (Core Service Hardening):** The hardening baselines defined in D5-003 are stored in the version control repository. Drift detection (D5-003, Section 4.3) is a subset of the full configuration audit defined here.
- **D5-004 (Network Architecture):** Network configurations (switch configs, firewall rules, VLAN assignments) are version-controlled with the same rigor as system configurations.
- **D5-005 (Monitoring):** The monitoring system's own configuration is version-controlled. Changes to alert thresholds, dashboard configurations, and monitoring agent settings follow the change management process.

## 5. Rules & Constraints

- **R-D5-006-01:** Every configuration file on every production system must exist in version control. Configurations that exist on a system but not in version control are unauthorized.
- **R-D5-006-02:** Every configuration change must follow the change management process in Section 4.2. Emergency changes may abbreviate the proposal step, but must complete all other steps and must document the emergency justification retroactively.
- **R-D5-006-03:** The configuration audit (Section 4.4) must be performed quarterly.
- **R-D5-006-04:** Configuration reproducibility must be tested annually by rebuilding a system from documentation alone.
- **R-D5-006-05:** The version control repository must be backed up as a critical institutional asset. Loss of version control history is a significant operational incident.
- **R-D5-006-06:** Every commit to the configuration repository must include a description of the change and a reference to the operational log entry that authorized it.
- **R-D5-006-07:** The configuration repository must maintain full history. History-rewriting operations (force push, rebase that discards commits) are prohibited.

## 6. Failure Modes

- **Configuration drift accumulation.** Configurations are changed on live systems without being committed to version control. Over time, the version control repository describes a system that no longer exists. Mitigation: quarterly audits (R-D5-006-03) and the discipline of always changing version control first, then deploying (Section 4.2, Step 2).
- **Version control neglect.** The repository exists but is not used. Configurations are managed informally, in the operator's memory. Mitigation: R-D5-006-01 makes version-controlled configuration a hard requirement, not a best practice. The quarterly audit will reveal neglect immediately.
- **Repository corruption.** The version control repository becomes corrupted (disk failure, software bug). Configuration history is lost. Mitigation: the repository is backed up (R-D5-006-05). The backup includes the full repository with complete history. Backup integrity is verified quarterly (OPS-001).
- **Change management bypass.** Under time pressure, the operator skips the change management process and modifies a production system directly. The change works, so it is never documented. Mitigation: the configuration audit will detect the undocumented change. Emergency changes are permitted but must be documented retroactively (R-D5-006-02).
- **Documentation-reality divorce.** The documentation is thorough but does not match reality because it was written aspirationally ("this is how we want it to be") rather than descriptively ("this is how it actually is"). Mitigation: the reproducibility test (R-D5-006-04) is the cure. If the system rebuilt from documentation does not match the production system, the documentation is wrong.
- **Over-complexity in automation.** The operator builds an elaborate local automation system that becomes itself a source of complexity and fragility. The automation tool breaks, and the operator cannot manage configurations manually because they have forgotten how. Mitigation: automation is introduced only when manual deployment becomes impractical (Section 4.3). Manual deployment skills must be maintained regardless of automation status.

## 7. Recovery Procedures

**If the version control repository has been corrupted or lost:**
1. Restore the repository from the most recent backup.
2. Identify any changes made between the backup and the corruption event by comparing the restored repository against live system configurations.
3. Commit any identified changes to the restored repository with notes indicating they were recovered after a repository incident.
4. Verify repository integrity after restoration.
5. Document the incident, the data loss (if any), and the restoration in the operational log.
6. Review the backup strategy for the repository -- was the backup frequency sufficient?

**If configuration drift has been discovered during an audit:**
1. For each drifted configuration, determine whether the live system or the repository contains the correct version.
2. If the live system is correct (an undocumented but intentional change), commit the live configuration to the repository with full documentation of the change and its rationale. Add a note that the change was documented retroactively.
3. If the repository is correct (an unintentional change on the live system), restore the repository version to the live system and verify.
4. If the correct version is uncertain, investigate. Check the operational log, check system logs, check if the change correlates with any maintenance activity. Do not commit or restore until the correct version is determined.
5. Document every resolution in the audit report.

**If a system must be rebuilt and the documentation is incomplete:**
1. If other systems with similar roles exist, use their configurations as a reference.
2. Rebuild to the closest approximation possible from documentation.
3. Test thoroughly before placing the rebuilt system in production.
4. Document everything that was missing from the original documentation and add it to the repository.
5. Treat this as a documentation failure and review the reproducibility testing process to ensure it catches such gaps in the future.

**If the change management process has been consistently bypassed:**
1. Acknowledge the failure in the operational log.
2. Perform a full configuration audit immediately.
3. Document all undocumented changes.
4. Investigate why the process was bypassed. Was it too burdensome? Too slow? Too complex? If so, revise the process (Tier 3 governance decision) to make it practical while maintaining its essential requirements.
5. Resume the process. Do not declare the problem solved until at least two quarterly audits show compliance.

## 8. Evolution Path

- **Years 0-2:** The version control repository is being populated for the first time. Every system deployment generates configuration artifacts that must be committed. The change management process will feel cumbersome. Follow it anyway. The habits formed now define the institution's configuration culture.
- **Years 2-5:** The repository is comprehensive. The quarterly audits should show minimal drift. The reproducibility test should succeed for most systems. The focus shifts from establishing the process to refining it.
- **Years 5-15:** The repository contains a rich history of how the institution's configuration has evolved. This history becomes valuable during troubleshooting ("When did this setting change? What was it before?") and during planning ("How has this system's resource configuration grown over time?"). The repository is no longer just a management tool. It is institutional memory.
- **Years 15-30:** Version control tools may have changed. If the tool is replaced, the full history must be migrated to the new tool. Configuration history is a critical institutional asset that spans the lifetime of the institution, not the lifetime of any particular tool.
- **Years 30-50+:** A successor will inherit the repository. The repository, its structure, and the change management process must be comprehensible without oral instruction. The README files in the repository and the Commentary Section of this article provide the context the successor needs. The repository is, in a real sense, a letter to the future.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The most important sentence in this article is in Section 4.2, Step 2: "The change is made to the configuration file in version control first, not on the target system." This is counterintuitive. The natural workflow is to change the system, verify it works, and then document it. But that workflow creates a window where the live system and the documented configuration diverge. That window is where configuration drift begins. By changing version control first, you ensure that the repository is always at least as current as the live systems. The reverse -- live systems leading the repository -- is the path to eventual documentation fiction.

The quarterly audit will, in the early years, consistently find drift. This is not failure. It is the audit working. Every audit that finds drift and corrects it is the process functioning as designed. Worry when the audits stop finding anything -- that either means the process has been perfected (unlikely in the first few years) or the audits are not being performed thoroughly enough.

On the question of local automation: resist the urge to build elaborate tooling. For an institution with a handful of systems, a well-organized repository and disciplined manual deployment are more reliable than a fragile automation system that the operator does not fully understand. Automation should be introduced when the pain of manual deployment exceeds the pain of maintaining automation -- and not a moment before.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 3: Transparency of Operation)
- CON-001 -- The Founding Mandate (self-built, self-understood infrastructure)
- OPS-001 -- Operations Philosophy (documentation-first principle, operational tempo)
- P-001 -- Platform Philosophy and Architectural Principles
- P-009 -- Configuration Management and Reproducibility
- D5-002 -- Operating System Maintenance Procedures (update-triggered configuration changes)
- D5-003 -- Core Service Configuration and Hardening (configuration baselines and drift detection)
- D5-004 -- Network Architecture and Segmentation (network configuration management)
- D5-005 -- System Monitoring and Alerting (monitoring configuration management)
- GOV-001 -- Authority Model (governance tiers for configuration changes)
- Domain 12 -- Disaster Recovery (system rebuild from configuration documentation)

---

*End of Stage 3 Platform Operations -- Five Operational Doctrine Articles*

**Document Total:** 5 articles (D5-002 through D5-006)
**Status:** All five articles ratified as of 2026-02-16.
**Next:** Additional Stage 3 articles for remaining Domain 5 topics and other operational domains.
