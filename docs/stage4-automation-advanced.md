# STAGE 4: SPECIALIZED SYSTEMS -- AUTOMATION ADVANCED REFERENCE

## Domain 8 Advanced Articles D8-007 through D8-011

**Document ID:** STAGE4-AUTOMATION-ADVANCED
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Classification:** Stage 4 -- Specialized Systems. Advanced reference documents that provide detailed technical and procedural depth for Domain 8: Automation & Agents. These articles build upon the philosophical foundations of D8-001 (Automation Restraint Doctrine) and the operational procedures of earlier Domain 8 articles. They assume mastery of D8-001 through D8-006.

---

## How to Read This Document

This document contains five advanced articles for Domain 8: Automation & Agents. Stage 2 established the philosophy -- automation exists to serve, never to govern, and the default posture is restraint. Stage 3 translated that philosophy into operational procedures. Stage 4 provides the specialized depth that the operator needs when automation has been justified, approved, and is being designed, built, monitored, tested, and maintained at a level of rigor that matches the institution's fifty-year horizon.

These are reference documents. They are not written to be read in a single sitting but to be consulted when the operator is designing a new automation, debugging a failing one, auditing the health of existing systems, or preparing for succession. Each article is self-contained, but they cross-reference heavily. When an article refers to a concept defined elsewhere, the reference is explicit.

The audience for these articles is a competent operator who has read D8-001 and accepted its premises. If you have not read D8-001, stop here and read it. These articles will seem excessively bureaucratic without the context of why restraint matters. With that context, they will seem like the minimum necessary discipline for systems that must outlast their creators.

If you are a future maintainer encountering these articles for the first time: they were written to help you. Every template, every checklist, every standard exists because the alternative -- undocumented, unmonitored, untested automation -- is the single greatest threat to institutional comprehensibility across generations.

---

---

# D8-007 -- Automation Specification Language

**Document ID:** D8-007
**Domain:** 8 -- Automation & Agents
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, SEC-001, OPS-001, D8-001, D8-002, D8-003, D8-006, D6-001, D6-008
**Depended Upon By:** D8-008, D8-009, D8-010, D8-011. All articles that create, modify, or audit automation.

---

## 1. Purpose

This article defines the Automation Specification Language: a standardized, human-readable format for describing every automation in the institution. The specification is not source code. It is the document that must exist before source code is written, that must be updated whenever source code changes, and that must remain comprehensible decades after the source code has been forgotten.

The fundamental problem this article solves is the gap between what an automation does and what a future maintainer can understand about what it does. Source code describes mechanism -- the how. The specification describes intent, context, boundaries, and consequences -- the what, the why, the when, the what-if, and the what-then. A maintainer who reads only the source code knows what the machine does. A maintainer who reads the specification knows what the institution intended, what constraints were considered, what failures were anticipated, and what the human fallback looks like.

Per D8-001 Criterion 3 (Comprehensibility), every automation must be fully understood by a competent generalist reading its documentation and source code. This article defines what "documentation" means in that context. It is the standard against which comprehensibility is measured. An automation whose specification is incomplete, outdated, or incomprehensible has failed the comprehensibility criterion regardless of how elegant its code may be.

## 2. Scope

This article covers:

- The canonical specification template that every automation must use.
- The required fields and their definitions.
- How to describe inputs, outputs, triggers, dependencies, failure modes, and human checkpoints.
- How to write a specification that a future maintainer can understand without running the code.
- The specification review process and update requirements.
- How specifications relate to the automation registry defined in D8-006.

This article does not cover:

- The source code itself or programming standards (see D5-006 for configuration management; coding standards are implementation-specific).
- The governance approval process for new automation (see D8-010).
- Monitoring and observability requirements (see D8-008 and D8-011).
- Testing requirements (see D8-010).

## 3. Background

### 3.1 Why a Specification Language Exists

Every institution that has accumulated automation over years has experienced the same failure: someone opens a script, reads the code, and cannot determine why it exists, what triggers it, what happens if it fails, or whether it is safe to modify. Comments in the code help. A README file helps more. But neither provides the structured, auditable, complete description that institutional governance requires.

The Automation Specification Language exists because unstructured documentation degrades. A free-form description encourages the author to document what they find interesting and omit what they find obvious -- and what is obvious to the creator is precisely what is opaque to the inheritor. A structured template forces completeness. Every field must be filled. Every question must be answered. The template is the institution's way of asking the automation creator: "Have you thought about this?" for every category of concern.

### 3.2 The Specification as Contract

The specification is a contract between the automation's creator and every future maintainer. It promises: this is what the automation does, these are the conditions under which it operates, these are the ways it can fail, and this is how to handle those failures. If the automation's behavior diverges from its specification, the specification is authoritative -- the automation has a bug, or the specification needs updating. Either way, the divergence must be resolved. Silent drift between specification and behavior is a failure mode addressed in Section 6.

### 3.3 Human-Readable by Design

The specification is written in plain language, not in code, not in a domain-specific language that requires a parser, and not in a format that requires specialized tools to read. It is a text document. It uses the institution's standard document format per D6-001 and D6-003. A future maintainer who has access to nothing more than a text editor can read, understand, and modify any specification. This is a deliberate constraint. Per D8-001, comprehensibility outranks efficiency.

## 4. System Model

### 4.1 The Canonical Specification Template

Every automation in the institution must have a specification document conforming to the following template. Each field is mandatory unless explicitly marked optional.

**AUTOMATION SPECIFICATION**

```
SPEC-ID:           [Unique identifier, format: AUTO-[DOMAIN]-[SEQUENCE]]
SPEC-VERSION:      [Semantic version of this specification]
AUTOMATION-NAME:   [Human-readable name]
CREATED-BY:        [Name and role of the creator]
CREATED-DATE:      [Date of initial creation]
LAST-MODIFIED-BY:  [Name and role of last modifier]
LAST-MODIFIED:     [Date of last modification]
AUTOMATION-LEVEL:  [Level 1-4 per D8-001 Section 4.2]
GOVERNANCE-APPROVAL: [Reference to approval record per GOV-001, required for Level 3-4]
STATUS:            [Draft | Active | Suspended | Retired]

1. PURPOSE
   Why this automation exists. What institutional need it serves. What manual
   process it replaces or assists. The justification against D8-001's five
   criteria, summarized in one paragraph with reference to the full
   justification record.

2. INPUTS
   Every input the automation consumes, listed individually:
   - Input name
   - Source (where the input comes from: file, system state, sensor, human entry)
   - Format (data type, encoding, structure)
   - Validation rules (what constitutes valid input, how invalid input is handled)
   - Frequency (how often the input is expected)
   - What happens if the input is missing or delayed

3. OUTPUTS
   Every output the automation produces, listed individually:
   - Output name
   - Destination (where the output goes: file, log, system state, human notification)
   - Format (data type, encoding, structure)
   - Expected values (what normal output looks like)
   - Abnormal values (what indicates a problem)
   - Retention requirements (per D6-001 data tier classification)

4. TRIGGERS
   What causes the automation to execute:
   - Trigger type (scheduled, event-driven, manual invocation, dependency chain)
   - Trigger source (cron entry, systemd timer, file watcher, human command)
   - Trigger frequency (how often the automation runs under normal conditions)
   - Trigger conditions (any prerequisites that must be true for execution)
   - What happens if the trigger fires but conditions are not met

5. DEPENDENCIES
   Everything the automation requires to function:
   - Software dependencies (packages, libraries, interpreters, with version constraints)
   - System dependencies (filesystems, services, network [none for air-gapped], devices)
   - Data dependencies (files, databases, configuration that must exist)
   - Automation dependencies (other automations that must run first or concurrently)
   - Human dependencies (approvals, confirmations, inputs per D8-003)

6. FAILURE MODES
   Every known way the automation can fail:
   - Failure mode name
   - Cause (what conditions produce this failure)
   - Detection (how the operator knows this failure has occurred)
   - Impact (what happens to the institution if this failure occurs)
   - Severity (Critical / High / Medium / Low)
   - Recovery action (what the operator does to recover)
   - Automated recovery (if any, what the automation does to self-recover)

7. HUMAN CHECKPOINTS
   Every point at which human attention is required:
   - Checkpoint name
   - When it occurs (before execution, during execution, after execution)
   - What the human must verify
   - What response options the human has (approve, modify, reject, defer)
   - What happens if the human does not respond within the expected timeframe
   - Required competency level of the human (per D9-002)

8. MANUAL FALLBACK
   The complete procedure for performing this automation's task manually:
   - Step-by-step instructions
   - Estimated time for manual execution
   - Required tools and access
   - How to verify that the manual execution produced the correct result
   (This section must be tested periodically per D8-001, R-D8-05.)

9. KILL PROCEDURE
   How to stop this automation immediately:
   - Graceful stop command (allows current operation to complete)
   - Immediate kill command (terminates without cleanup)
   - Physical kill method (power, storage disconnection per D8-009)
   - Post-kill verification (how to confirm the automation has actually stopped)
   - Post-kill cleanup (any manual steps required after emergency stop)

10. LOGGING AND OBSERVABILITY
    What the automation records about its own operation:
    - Log location (filesystem path)
    - Log format (per D8-011 standards)
    - Log rotation policy
    - Metrics exposed (per D8-011 requirements)
    - Health check endpoint or method
    - How to verify the automation is functioning correctly from outside

11. CHANGE HISTORY
    Dated record of every modification to this specification:
    - Date
    - Author
    - Summary of change
    - Reason for change
```

### 4.2 Specification Naming and Storage

Specifications are stored alongside the automation they describe, in the same directory, with the naming convention `[automation-name].spec.txt`. A copy of every active specification is maintained in the Automation Registry defined in D8-006. The registry copy is the authoritative version. If the local copy and the registry copy diverge, the registry copy prevails until the divergence is investigated and resolved.

### 4.3 Writing for the Future Maintainer

The specification must be written as if the reader has never seen the automation, has never spoken to its creator, and has only the specification and the source code. This means:

Avoid jargon unless it is defined in the institutional glossary. Avoid references to context that is not documented. Avoid phrases like "as you know" or "obviously" -- nothing is obvious to someone reading this for the first time in 2046.

The Purpose section should explain the problem, not just the solution. A specification that says "runs the backup rotation script" tells the reader what happens. A specification that says "prevents the backup directory from exceeding disk capacity by removing verified backups older than 90 days, because the primary storage volume has a fixed capacity of X and the daily backup size averages Y" tells the reader why it happens, what constraints drove the design, and what assumptions might need revisiting.

The Failure Modes section should be honest, not optimistic. Every automation can fail. The specification must admit this and enumerate the failures the creator considered. A specification with no failure modes documented is not a specification of a perfect automation -- it is a specification written by someone who did not think carefully enough.

### 4.4 The Specification Review Cycle

Every specification must be reviewed:

- When the automation is modified in any way.
- When the automation's dependencies change (software update, hardware change, data format change).
- During the annual automation inventory review per D8-001, R-D8-07.
- When a failure occurs that was not documented in the Failure Modes section.
- When a new operator assumes responsibility for the automation.

The review verifies that the specification accurately describes the automation's current behavior. Divergences are documented and resolved. If the automation's behavior has changed without a corresponding specification update, this is treated as a documentation debt and prioritized for resolution.

## 5. Rules & Constraints

- **R-D8-07-01:** Every automation in the institution must have a specification conforming to the template defined in Section 4.1. No automation may operate in production without a complete, current specification.
- **R-D8-07-02:** The specification must be written before the automation is deployed. Draft specifications are permitted during development, but all fields must be complete before production deployment.
- **R-D8-07-03:** The specification must be updated within 72 hours of any change to the automation's behavior, inputs, outputs, dependencies, or failure characteristics.
- **R-D8-07-04:** Specifications must be stored in plain text format per D6-003. No proprietary formats. No formats requiring specialized rendering tools.
- **R-D8-07-05:** The Manual Fallback section (Section 8 of the template) must be tested at least annually. Test results are recorded in the specification's Change History.
- **R-D8-07-06:** The specification is classified as Tier 2 data (Operational Data) per D6-001. It follows the retention and backup requirements for that tier.
- **R-D8-07-07:** When an automation is retired per D8-006, its specification is preserved in the archive with a status change to "Retired" and a dated note explaining why the automation was retired. Retired specifications are never deleted; they are reclassified as Tier 3 (Reference Data).

## 6. Failure Modes

- **Specification drift.** The automation changes but the specification does not. Over time, the specification describes something that no longer exists, and the actual automation is undocumented. Detection: annual review finds discrepancies. Impact: future maintainer relies on inaccurate documentation. Mitigation: R-D8-07-03 requires updates within 72 hours. The annual review catches what the 72-hour rule misses.

- **Specification theater.** The operator fills in the template mechanically, using vague or boilerplate language, without genuine thought about each field. The specification exists but does not actually help a future reader. Detection: review by a second party (or self-review after a cooling-off period) reveals fields that do not contain actionable information. Mitigation: the specification review process must evaluate quality, not just completeness.

- **Template rigidity.** The template does not accommodate a new type of automation that does not fit the existing fields. The operator either forces the automation into the template or abandons the template entirely. Mitigation: the template may be extended with additional fields as needed, provided all canonical fields are still present. Extensions are documented in the Change History. If the template itself needs revision, this is a Tier 3 governance action per GOV-001.

- **Orphaned specification.** The automation is deleted or replaced, but the specification remains active in the registry. The registry lists automation that no longer exists. Detection: annual inventory reconciliation between the registry and actual running systems. Mitigation: the retirement procedure in D8-006 requires specification status update as part of the retirement checklist.

- **Incomprehensible specification.** The specification was written by someone with deep context and is impenetrable to someone without it. Detection: succession testing -- a new operator attempts to understand the automation from its specification alone. Mitigation: specifications should be written for the least-experienced plausible reader, and reviewed for comprehensibility as part of the annual review.

## 7. Recovery Procedures

1. **If specifications are missing for existing automations:** Declare a specification sprint. Inventory all active automations from the registry and from system inspection (running processes, cron entries, systemd timers). For each automation without a specification, create one. If the automation's behavior cannot be determined through documentation review and code reading, use D8-011 (Legacy Automation) procedures to analyze it. Prioritize by automation level: Level 4 and Level 3 automations first, then Level 2, then Level 1.

2. **If specification drift is widespread:** Freeze non-critical automation changes. Audit each specification against its automation's actual behavior. Update specifications to match reality. If reality has diverged in ways that violate institutional constraints (e.g., the automation now touches Tier 1 data without documented human checkpoints), the automation must be corrected to match the specification, not the other way around.

3. **If the template is found inadequate:** Document the specific inadequacy. Propose template extensions through the governance process. In the interim, use the existing template with an additional "EXTENSIONS" section at the end, clearly marked as non-standard. When the template is formally updated, migrate all specifications to the new version.

4. **If a specification is found to be incomprehensible by a new operator:** The new operator documents their confusion -- which fields were unclear, what context was missing, what assumptions were unstated. The specification is rewritten collaboratively, with the original author (if available) providing context and the new operator verifying that the rewrite is genuinely comprehensible. If the original author is unavailable, the new operator rewrites based on their investigation and marks the specification as "reconstructed" with the date and their confidence level.

## 8. Evolution Path

- **Years 0-5:** The specification template is new. Early specifications will be imperfect as the operator learns what level of detail is genuinely useful versus what is bureaucratic overhead. Expect to refine the template based on experience. The first time a specification saves the operator from a mistake or enables quick recovery from a failure, the value of the discipline will be clear.

- **Years 5-15:** The specification corpus should be stable. The annual review cycle should be routine. The primary challenge is maintaining specification accuracy as automations evolve. This is also the period to test whether specifications are truly comprehensible by having someone unfamiliar with each automation attempt to understand it from the specification alone.

- **Years 15-30:** Specifications written in years 0-5 will be tested by time. Do they still make sense? Are the assumptions they document still valid? Are the failure modes they describe still relevant, or have new ones emerged? This is the first real test of the specification as a multi-generational document.

- **Years 30-50+:** A successor inherits the specification corpus. The quality of these specifications determines whether they can understand and maintain the institution's automation or must rebuild from scratch. Every shortcut taken in writing specifications is a debt imposed on this future operator.

- **Signpost for revision:** If the specification template consistently requires the same extensions, incorporate those extensions into the canonical template. If operators consistently find certain fields unhelpful, evaluate whether those fields should be simplified or restructured -- but resist the urge to remove fields merely because they are inconvenient. The inconvenience is the point.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
I have worked with systems where the documentation was the code and the code was the documentation. The argument is seductive: code is precise, unambiguous, and always up-to-date because it is the thing that actually runs. The problem is that code tells you what the machine does. It does not tell you what the institution intended. It does not tell you what was considered and rejected. It does not tell you what the operator should do when the code fails. Code is a monologue from the machine. A specification is a conversation between the creator and the future maintainer.

The template may seem excessive for a simple cron job that rotates log files. It is not. That simple cron job will run unsupervised for years. It will be inherited by someone who did not create it. It will eventually fail in a way that its creator did not anticipate. When that day comes, the specification is the difference between a fifteen-minute fix and a three-hour investigation.

Write the specification as if you are writing a letter to someone who will need your help and will not be able to reach you. Because that is exactly what it is.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 3: Transparency of Operation)
- CON-001 -- The Founding Mandate (comprehensibility requirement)
- OPS-001 -- Operations Philosophy (documentation-first principle)
- GOV-001 -- Authority Model (governance tiers for template changes)
- D6-001 -- Data Philosophy (data tier classification for specifications)
- D6-003 -- Format Longevity Doctrine (plain text format requirement)
- D6-008 -- Metadata Standards (metadata requirements for specification documents)
- D8-001 -- Automation Restraint Doctrine (justification criteria, automation levels, comprehensibility criterion)
- D8-002 -- Agent Design Principles (design constraints that specifications must reflect)
- D8-003 -- Human-in-the-Loop Doctrine (human checkpoint requirements)
- D8-006 -- Agent Lifecycle (registry, retirement procedures)
- D8-009 -- Emergency Override & Kill Procedures (kill procedure section of specification)
- D8-011 -- Automation Observability Standards (logging and metrics sections of specification)

---

---

# D8-008 -- Monitoring Automation: Watching the Watchers

**Document ID:** D8-008
**Domain:** 8 -- Automation & Agents
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, SEC-001, OPS-001, D8-001, D8-002, D8-003, D8-007, D6-001
**Depended Upon By:** D8-009, D8-010, D8-011, D10-009. All operational articles involving automation health assessment.

---

## 1. Purpose

This article defines how the institution monitors its automated systems without creating infinite regress -- the problem where monitoring systems must themselves be monitored, and the monitors of those monitors must also be monitored, ad infinitum. This is not a theoretical concern. It is the central architectural challenge of automation monitoring in a small, single-operator institution where the complexity budget defined in OPS-001 is finite and every layer of monitoring consumes resources that cannot be spent elsewhere.

The solution is architectural, not technical. It rests on three principles: independent verification through dissimilar mechanisms, a defined terminus where monitoring ends and human attention begins, and a discipline of alert management that prevents the operator from drowning in noise and missing the signals that matter.

Monitoring automation is itself automation. It is subject to every constraint in D8-001, including the justification framework, the comprehensibility requirement, and the manual fallback obligation. A monitoring system that the operator cannot understand is worse than no monitoring at all, because it provides false confidence -- the operator believes someone is watching when in fact nothing effective is watching.

This article provides the architecture for monitoring that is honest about its limitations, effective within those limitations, and designed to age gracefully across hardware generations and operator transitions.

## 2. Scope

This article covers:

- The infinite regress problem and the institution's architectural solution.
- Independent verification strategies that avoid common-mode failures.
- Health metrics for automation: what to measure and what the measurements mean.
- Alert fatigue prevention: how to design alerts that the operator will actually read.
- The automation health dashboard: what it shows, how it is maintained, and when it is consulted.
- The monitoring terminus: where automated monitoring ends and human ritual begins.

This article does not cover:

- The internal observability of individual automations (see D8-011).
- The specification of what each automation must log (see D8-007 and D8-011).
- Emergency override procedures triggered by monitoring alerts (see D8-009).
- General system monitoring unrelated to automation (see D5-005).
- Incident response procedures (see D10-005).

## 3. Background

### 3.1 The Infinite Regress Problem

If automated systems must be monitored, and the monitor is an automated system, then the monitor must also be monitored. This creates a logical chain that has no natural terminus. In large organizations, this problem is managed through organizational layering: different teams monitor different layers, and the probability that all layers fail simultaneously is acceptably low. A single-operator institution has no such luxury. One person cannot operate, monitor, and meta-monitor simultaneously.

The institution resolves this problem not by extending the chain to infinity but by terminating it deliberately. The monitoring architecture has exactly two layers of automated monitoring, and the third layer is the human operator performing a structured ritual. This is not a compromise. It is the correct design for this context. Automated monitoring catches the failures that happen between human inspections. The human inspections catch the failures in the automated monitoring.

### 3.2 The Dissimilar Redundancy Principle

A monitoring system that shares failure modes with the system it monitors is useless at precisely the moment it is most needed. If the monitoring system runs on the same disk that the automation writes to, a disk failure takes out both the automation and its monitor simultaneously. If the monitoring system uses the same logging infrastructure, a logging failure blinds both.

Dissimilar redundancy means: the monitoring mechanism must be architecturally different from the mechanism it monitors. A cron job is not effectively monitored by another cron job on the same system using the same scheduler. It is effectively monitored by a systemd timer on a different system checking for the expected output, or by a human checking a physical indicator, or by an independent process that verifies outcomes rather than observing processes.

### 3.3 Alert Fatigue: The Silent Killer

Alert fatigue is the condition where the operator receives so many alerts that they stop reading them carefully, respond to them mechanically, or ignore them entirely. It is the single most common failure mode in monitoring systems, and it is more dangerous than having no monitoring at all, because it creates the illusion of oversight.

The institution's approach to alert fatigue is simple and absolute: every alert that reaches the operator must be actionable. If an alert does not require the operator to do something specific, it is not an alert -- it is noise, and it must be removed from the alert channel. The operator must never be trained, by repeated exposure to false or irrelevant alerts, to ignore the alert channel.

## 4. System Model

### 4.1 The Three-Layer Monitoring Architecture

**Layer 1: Self-Monitoring (Automated).**
Each automation monitors its own health and reports its status. This is the observability requirement from D8-001 Criterion 4, implemented per D8-011 standards. Each automation writes logs, exposes metrics, and reports success or failure at the conclusion of each execution. Self-monitoring is the first line of detection. Its limitation is that it cannot detect its own failure to report -- a crashed automation does not log that it has crashed.

**Layer 2: Independent Monitoring (Automated).**
A dedicated monitoring system, architecturally separate from the automations it watches, verifies that Layer 1 is functioning. The independent monitor does not observe the automation's internal state. It observes outcomes: Did the expected output appear? Did the log file update? Did the status flag change? Did the automation complete within its expected timeframe? The independent monitor checks for the presence of expected evidence, not for the process itself. This is the watchdog layer. Its limitation is that it is itself an automation and can itself fail.

**Layer 3: Human Verification (Ritual).**
The operator performs a structured daily review of both the automation outputs and the monitoring system's reports. This is not an ad-hoc glance at a dashboard. It is a checklist-driven inspection per OPS-001 operational tempo. The operator verifies: Are all automations reporting? Is the monitoring system itself reporting? Do the reports match expected patterns? This human layer is the terminus of the monitoring chain. It catches failures in both Layer 1 and Layer 2.

### 4.2 Independent Verification Strategies

The independent monitor (Layer 2) uses the following strategies to verify automation health:

**Heartbeat checking.** Each automation writes a timestamp to a known location upon successful execution. The independent monitor checks these timestamps. If a timestamp is older than the expected interval plus a defined tolerance, the automation is flagged as potentially failed. The heartbeat file must reside on a different filesystem than the automation's primary workspace, to survive the storage failures most likely to kill the automation.

**Output verification.** For automations that produce tangible outputs, the independent monitor verifies that expected output exists, was created within the expected timeframe, and meets minimum validity criteria (non-zero size, expected format).

**Log analysis.** The independent monitor scans automation logs for error patterns, unexpected warnings, or the absence of expected success messages. This is a lightweight scan for known-bad patterns, not comprehensive log analysis.

**Resource monitoring.** The independent monitor tracks resource consumption: CPU time, memory usage, disk space, execution duration. Significant deviations from baseline indicate a problem or an environmental change warranting investigation.

### 4.3 Health Metrics for Automation

Every automation must expose, through its logging and status reporting per D8-011, the following health metrics:

- **Last successful execution time.** When did the automation last complete without errors?
- **Last execution duration.** How long did the most recent execution take, compared to the historical average?
- **Execution success rate.** Over the last N executions, what percentage completed without errors?
- **Error count and classification.** How many errors in the current period, and what types?
- **Resource high-water marks.** Peak CPU, memory, and disk usage during the current period.
- **Dependency status.** Are all required inputs, services, and upstream automations available?
- **Output freshness.** When was the most recent output produced, and is it within expected bounds?

These metrics are collected by the independent monitor and presented on the automation health dashboard.

### 4.4 Alert Design and Fatigue Prevention

Alerts are classified into three tiers:

**Tier 1: Action Required.** Something has failed or is actively failing. The operator must investigate and take action. Examples: automation has not reported in more than twice its expected interval, output verification has failed, resource consumption has exceeded safety thresholds. Tier 1 alerts are rare. If Tier 1 alerts occur more than twice per week on average across all automations, the alerting thresholds are miscalibrated and must be adjusted.

**Tier 2: Attention Warranted.** Something is anomalous but not yet failed. The operator should investigate at the next scheduled review. Examples: execution duration has increased by more than 50% from baseline, success rate has dropped below 95% over the last 10 executions, a non-critical dependency is unavailable. Tier 2 alerts are informational. They appear on the dashboard but do not interrupt the operator.

**Tier 3: Logged Only.** Something has been noted for trend analysis. No immediate action is needed. Examples: minor resource fluctuations, single-instance errors that recovered automatically, expected maintenance-window anomalies. Tier 3 items appear only in logs and periodic reports.

**The Alert Hygiene Discipline:** Every Tier 1 alert that fires must be investigated. Every Tier 1 alert that proves to be a false positive must result in a threshold adjustment or alert removal. The operator must never learn to ignore Tier 1 alerts. If the operator finds themselves dismissing Tier 1 alerts without investigation, the alert system has failed and must be redesigned immediately.

### 4.5 The Automation Health Dashboard

The dashboard is a single, consolidated view of all automation health. It is not a web application. It is a text-based report generated at defined intervals (default: hourly during operational hours, once overnight), stored as a file, and consulted by the operator during daily review.

The dashboard displays:

- A list of all active automations with their current status (Healthy, Warning, Failed, Unknown).
- The timestamp of each automation's last successful execution.
- Any active Tier 1 or Tier 2 alerts.
- The monitoring system's own health status (self-check).
- A summary of the last 24 hours: total executions, total failures, any anomalies.

The dashboard generation itself is an automation. It has its own specification per D8-007. Its failure is detected by the human verification layer: if the dashboard file has not been updated, the operator knows the monitoring system has a problem.

## 5. Rules & Constraints

- **R-D8-08-01:** The institution shall maintain exactly three layers of monitoring as defined in Section 4.1. Adding a fourth automated layer is prohibited without Tier 2 governance approval per GOV-001, because it adds complexity without proportionate benefit.
- **R-D8-08-02:** The independent monitoring system (Layer 2) must not share critical dependencies with the automations it monitors. At minimum, it must use a different execution mechanism (e.g., if automations use cron, the monitor uses a systemd timer, or vice versa) and write to a different storage location.
- **R-D8-08-03:** Tier 1 alerts must be investigated within 4 hours of detection. If the operator is unavailable (sleep, absence), investigation must occur at the next operational period. The alert must persist until acknowledged.
- **R-D8-08-04:** The operator must perform the Layer 3 human verification ritual daily during normal operations. The verification checklist is documented in the operational runbook per D10-002.
- **R-D8-08-05:** Every false positive Tier 1 alert must result in a threshold adjustment or alert removal within 7 days of occurrence. Recurring false positives are an operational failure.
- **R-D8-08-06:** The automation health dashboard must be generated and accessible at all times. Its generation mechanism is a critical automation and is subject to the highest monitoring rigor: its heartbeat is checked by the independent monitor, and its output is verified by the human layer.
- **R-D8-08-07:** Alert thresholds must be reviewed and recalibrated quarterly as part of the operational review cycle. Thresholds that have not been adjusted in over a year must be explicitly re-justified.

## 6. Failure Modes

- **Infinite regress creep.** The operator, anxious about monitoring reliability, adds additional monitoring layers. Each layer adds complexity, consumes resources, and requires maintenance. The monitoring system becomes the most complex and fragile component in the institution. Mitigation: R-D8-08-01 caps automated monitoring at two layers. The human layer is the terminus.

- **Alert flood.** A systemic event (power fluctuation, storage issue) causes many automations to fail simultaneously. The monitoring system generates dozens of alerts, overwhelming the operator. Mitigation: alert correlation -- the independent monitor groups related alerts when they share a common cause (e.g., all failures on a specific host or storage volume) and presents a single root-cause alert.

- **Monitor-target coupling.** The monitoring system shares a failure mode with the system it monitors. Both fail together. The operator sees no alerts and believes everything is working. Mitigation: R-D8-08-02 requires architectural separation. The human verification layer catches the case where both fail.

- **Dashboard staleness.** The dashboard generation automation fails. The operator consults a stale dashboard and believes everything is healthy. Mitigation: the dashboard must display its own generation timestamp prominently. The operator's verification checklist includes confirming the dashboard timestamp is current.

- **Alert fatigue onset.** False positives accumulate. The operator begins dismissing alerts without investigation. A real failure is missed. Mitigation: R-D8-08-05 requires false positive resolution. The quarterly threshold review (R-D8-08-07) catches systemic miscalibration.

- **Monitoring theater.** Monitoring exists and reports green, but the metrics it checks are superficial and would not detect the failures that actually threaten the institution. Mitigation: the annual automation review (D8-001, R-D8-07) must evaluate whether monitoring is checking the right things, not just whether it is checking things.

## 7. Recovery Procedures

1. **If the independent monitor (Layer 2) has failed:** The human verification layer should detect this within one operational cycle (24 hours). The operator notices that the dashboard is stale or that heartbeat checks are not being performed. Recovery: restart the monitoring system. Review logs to determine why it failed. If the failure mode is shared with monitored automations, redesign to eliminate the coupling. During the period when Layer 2 is down, increase the frequency of Layer 3 human checks.

2. **If alert fatigue has set in:** Acknowledge the problem honestly. Review all active alerts. Disable or reclassify every alert that is not genuinely actionable. Reset thresholds to reduce false positives, accepting a temporary increase in the risk of missed detections. Re-establish the discipline that every Tier 1 alert is investigated. This is a process reset, and it should be documented in the Commentary Section.

3. **If monitoring and monitored systems have failed simultaneously:** This is a significant incident. The human verification layer is the detection mechanism. Recovery: bring systems back online one at a time, starting with the monitoring system. Verify each automation manually before relying on automated monitoring. Document the common failure mode and redesign to eliminate it.

4. **If the dashboard is discovered to be stale:** Immediately perform a manual check of all automation health indicators. Restart the dashboard generation automation. Determine how long the dashboard has been stale and audit the period of staleness for any failures that may have been missed.

## 8. Evolution Path

- **Years 0-5:** The monitoring architecture is being established. Start simple. Monitor heartbeats and output presence. Add more sophisticated checks only as you learn which failures actually occur. Resist the urge to build comprehensive monitoring on day one -- you do not yet know what needs monitoring most.

- **Years 5-15:** The monitoring system should be mature. Alert thresholds should be well-calibrated based on years of operational data. The primary challenge is maintaining monitoring accuracy as automations evolve. Every change to an automation must be reflected in its monitoring configuration.

- **Years 15-30:** Hardware transitions will test whether the monitoring architecture is sufficiently abstract. If monitoring is tightly coupled to specific tools or platforms, it will need to be rebuilt. The three-layer architecture should survive; the implementation of each layer will evolve.

- **Years 30-50+:** A successor must be able to understand and maintain the monitoring system. The monitoring system's own specification (per D8-007) and this article are their guide. If the monitoring system has become the most complex component in the institution, something has gone wrong.

- **Signpost for revision:** If the three-layer architecture consistently proves insufficient -- if failures regularly escape all three layers -- the architecture needs revisiting. If the three layers catch failures reliably but the institution finds the daily human verification too burdensome, consider whether automation maturity justifies extending the human review interval, but do not extend it without evidence that the automated layers are trustworthy.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The temptation with monitoring is to build something elaborate. Dashboards with graphs. Real-time metrics. Color-coded status panels. These are satisfying to build and impressive to look at. They are also maintenance burdens that serve the operator's desire for control more than the institution's need for oversight.

The most effective monitoring I have ever used was a text file that said "BACKUP OK -- 2026-02-16 03:15" and a checklist that said "verify backup timestamp is from today." When the timestamp was not from today, I knew something was wrong. That is the level of simplicity this institution should aspire to. Build complexity only when simplicity has proven insufficient, and not one moment before.

The human verification layer is not a weakness in the architecture. It is the architecture's greatest strength. A human who looks at a system every day develops an intuition for what normal looks like that no automated system can replicate. The daily check is not overhead -- it is the most sophisticated monitoring the institution has.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 2: Integrity Over Convenience; Principle 3: Transparency)
- CON-001 -- The Founding Mandate (single-operator constraint)
- SEC-001 -- Threat Model and Security Philosophy (defense in depth, Pillar 1)
- OPS-001 -- Operations Philosophy (operational tempo, complexity budget)
- D8-001 -- Automation Restraint Doctrine (Criterion 4: Observability; manual fallback requirement)
- D8-002 -- Agent Design Principles (minimal inter-automation dependencies)
- D8-007 -- Automation Specification Language (monitoring system must have its own specification)
- D8-009 -- Emergency Override & Kill Procedures (response to critical monitoring alerts)
- D8-011 -- Automation Observability Standards (what each automation must expose)
- D5-005 -- Service Health and System Monitoring (general system monitoring context)
- D6-001 -- Data Philosophy (data tier classification for monitoring logs)
- D10-002 -- Daily Operations Runbook (human verification checklist integration)
- D10-005 -- Incident Response (escalation from monitoring alerts to incidents)

---

---

# D8-009 -- Scheduled Task Architecture

**Document ID:** D8-009
**Domain:** 8 -- Automation & Agents
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, SEC-001, OPS-001, D8-001, D8-002, D8-003, D8-007, D8-008, D6-001
**Depended Upon By:** D8-010, D8-011, D10-002, D10-009. All articles involving scheduled or timed operations.

---

## 1. Purpose

This article defines how the institution designs, implements, maintains, and audits scheduled tasks -- the cron jobs, systemd timers, and any other time-triggered automations that run without human initiation. Scheduled tasks are the backbone of institutional automation. They are also the most common source of automation debt, because they are easy to create, easy to forget, and easy to accumulate until the system is running dozens of scheduled operations that no one fully understands.

The institution's approach to scheduled tasks is defined by D8-001's restraint doctrine: every scheduled task must be justified, specified, monitored, and auditable. But scheduled tasks present unique challenges that this article addresses specifically: temporal conflicts between tasks, resource contention when multiple tasks run simultaneously, the difficulty of testing time-dependent behavior, the long-term maintenance of tasks that run for years without human attention, and the registry problem of knowing, at any moment, exactly what scheduled tasks exist and what they do.

This article provides the architecture that transforms scheduled tasks from a scattered collection of cron entries into a coherent, documented, maintainable system. The Scheduled Task Registry defined here is the authoritative record of every time-triggered automation in the institution.

## 2. Scope

This article covers:

- Design principles for scheduled tasks in an air-gapped, single-operator environment.
- The Scheduled Task Registry: what it contains and how it is maintained.
- Implementation standards for cron jobs and systemd timers.
- The scheduling audit: how to verify that scheduled tasks are correct, conflict-free, and current.
- Conflict detection between scheduled tasks.
- Logging and verification requirements specific to scheduled tasks.
- Resource management for concurrent scheduled operations.

This article does not cover:

- Event-driven automations (triggered by system events rather than time).
- Human-initiated automations (manually invoked scripts or tools).
- The automation specification template (see D8-007, which provides the specification; this article provides scheduled-task-specific architectural requirements).
- Emergency override procedures (see D8-009 from the original Stage 1 numbering -- refer to the Emergency Override article in this institution's documentation corpus).
- General monitoring architecture (see D8-008).

## 3. Background

### 3.1 The Scheduled Task Accumulation Problem

Scheduled tasks accumulate. A new cron job is added to rotate logs. Another is added to check disk space. Another to verify backup integrity. Another to clean temporary files. Another to generate reports. Each one is individually justified. Each one is individually simple. But collectively, they form a system -- and that system was never designed as a whole. Tasks conflict. Tasks compete for resources. Tasks depend on other tasks' outputs but that dependency is implicit, encoded only in the scheduling times chosen by the operator.

In a conventional environment, this accumulation is managed (or tolerated) because the systems are rebuilt periodically, because multiple operators review each other's work, or because the impact of a failed scheduled task is low. In an institution designed to last fifty years, with a single operator, and where scheduled tasks may run for decades, accumulation without architecture is a path to incomprehensibility.

### 3.2 The Cron Problem

Cron is the most common scheduling mechanism on Unix-like systems and one of the least observable. Determining whether a cron job succeeded requires checking the job's own output -- the cron daemon records only invocation. Cron has no built-in dependency management, conflict detection, resource control, or health reporting.

Systemd timers provide more structure: logging integration, dependency declarations, and status reporting through standard commands. The institution should prefer systemd timers where supported. Regardless of mechanism, the architectural requirements in this article apply.

### 3.3 Temporal Reasoning is Difficult

Humans are poor at reasoning about time-based systems. "Run at 3 AM daily" seems simple until the operator realizes that three other tasks also run at 3 AM, that the filesystem the task writes to is being backed up at 3:15 AM, and that the system's maintenance window starts at 3:30 AM. Temporal conflicts are invisible until they cause failures, and they are difficult to diagnose because the failure appears to be in one task when the root cause is the interaction between tasks.

This article provides tools for making temporal relationships explicit and visible: the Scheduled Task Registry, the conflict detection procedure, and the scheduling audit.

## 4. System Model

### 4.1 The Scheduled Task Registry

The Scheduled Task Registry is a single document that lists every scheduled task in the institution. It is maintained as a plain text file in the institutional documentation, classified as Tier 2 data per D6-001. The registry is the authoritative source of truth about what scheduled tasks exist. If a scheduled task is running but not in the registry, it is unauthorized. If a scheduled task is in the registry but not running, it is an anomaly requiring investigation.

Each entry in the registry contains:

```
TASK-ID:           [Unique identifier, format: SCHED-[SEQUENCE]]
TASK-NAME:         [Human-readable name]
SPEC-REFERENCE:    [Reference to full specification per D8-007]
SCHEDULE:          [Cron expression or timer definition, in plain language AND machine format]
MECHANISM:         [cron | systemd-timer | other (specify)]
HOST:              [Which system the task runs on]
USER:              [Which user account the task runs under]
AUTOMATION-LEVEL:  [Per D8-001 Section 4.2]
RESOURCE-PROFILE:  [Low | Medium | High -- expected resource consumption]
DURATION-ESTIMATE: [Expected execution time: typical and maximum]
DEPENDENCIES:      [Other scheduled tasks that must complete before this one]
CONFLICTS:         [Other scheduled tasks that must not run concurrently]
LAST-VERIFIED:     [Date of last audit verification]
STATUS:            [Active | Suspended | Retired]
```

### 4.2 Scheduling Design Principles

**Principle 1: Spread the Load.** Scheduled tasks should be distributed across the available time, not clustered. The operator should maintain a visual timeline (a simple text-based schedule) showing when each task runs and how long it is expected to take. Gaps between tasks provide buffer for overruns and reduce the likelihood of resource contention.

**Principle 2: Declare Dependencies Explicitly.** If Task B requires the output of Task A, this dependency must be declared in both tasks' specifications and in the registry. The scheduling must ensure Task A completes before Task B starts, with sufficient margin for variation. Dependencies must never be implicit -- encoded only in the relative timing of two tasks.

**Principle 3: Guard Against Overlap.** If a task has not completed by the time its next scheduled execution arrives, the second execution must be prevented. Overlapping instances of the same task create race conditions, duplicate processing, and resource exhaustion. Every scheduled task must implement a locking mechanism that prevents concurrent execution of the same task.

**Principle 4: Fail Visibly.** A scheduled task that fails must produce evidence of its failure. Silent failure -- where a task fails but produces no output, no log entry, and no alert -- is the most dangerous failure mode. Every scheduled task must log its start, its completion (with success/failure status), and sufficient detail to diagnose failures.

**Principle 5: Design for Missed Executions.** If a system was powered off during a task's scheduled time, or if a task was suspended for maintenance, the system must handle the missed execution gracefully. Some tasks should run as soon as possible after a missed window (e.g., backups). Others should simply skip to the next scheduled time (e.g., log rotation). The behavior for missed executions must be documented in the specification.

### 4.3 Implementation Standards

**For cron-based tasks:**
- Use a dedicated cron directory (/etc/cron.d/ or equivalent) with one file per task, named after the TASK-ID.
- Each cron file must contain a header comment with the TASK-ID, task name, and a reference to the full specification.
- All cron jobs must redirect stdout and stderr to a log file specific to that task.
- All cron jobs must implement file-based locking to prevent overlapping execution.
- All cron jobs must write a heartbeat timestamp on successful completion per D8-008.

**For systemd timer-based tasks:**
- Each task has a paired .timer and .service unit file.
- Unit files must contain a Description that references the TASK-ID and task name.
- Timer accuracy settings must be appropriate for the task's requirements.
- Service units must specify resource limits (MemoryMax, CPUQuota) appropriate to the task's resource profile.
- OnFailure directives must point to an alert mechanism.

**For all scheduled tasks:**
- A wrapper script or function must handle: lock acquisition, start logging, execution, completion logging, heartbeat writing, lock release, and error handling.
- The wrapper is a standardized component maintained as institutional infrastructure, not reimplemented for each task.

### 4.4 The Scheduling Audit

The scheduling audit is performed quarterly and during the annual automation review per D8-001, R-D8-07. The audit verifies:

1. **Registry accuracy.** Every running scheduled task matches a registry entry. Every registry entry corresponds to a running scheduled task (or a deliberately suspended one). Discrepancies are investigated and resolved.

2. **Conflict analysis.** The temporal schedule is reviewed for conflicts: tasks that run at the same time and compete for the same resources, tasks whose execution times have drifted to overlap, dependencies that are no longer met by the current schedule.

3. **Resource accounting.** The combined resource consumption of all scheduled tasks is reviewed against the system's capacity. If scheduled tasks collectively consume more than 60% of any resource (CPU, memory, disk I/O) during any time window, the schedule must be adjusted to reduce peak load.

4. **Specification currency.** Each scheduled task's specification (per D8-007) is verified to match its current behavior. Schedules, dependencies, resource profiles, and failure modes are confirmed accurate.

5. **Logging verification.** Each scheduled task's logs are reviewed to confirm that logging is functioning, that logs contain sufficient detail, and that log rotation is preventing unbounded growth.

6. **Lock mechanism verification.** Each task's locking mechanism is tested to confirm it prevents overlapping execution.

### 4.5 Conflict Detection

Conflicts between scheduled tasks fall into three categories:

**Temporal conflicts.** Two or more tasks scheduled to run at overlapping times when they should not (because they share a resource, because one depends on the other, or because their combined resource consumption exceeds available capacity).

**Resource conflicts.** Two or more tasks that, when running concurrently, exhaust a shared resource (disk I/O, CPU, memory, a specific filesystem or device).

**Dependency conflicts.** A task depends on another task's output, but the scheduling does not guarantee that the upstream task completes before the downstream task starts.

Detection method: the scheduling audit (Section 4.4) compares each pair of scheduled tasks for potential conflicts. For N tasks, this requires examining N*(N-1)/2 pairs. In a well-governed institution following D8-001's restraint doctrine, N should be small enough that this analysis is tractable by hand. If N grows large enough that pairwise analysis is impractical, the institution has too many scheduled tasks and must consolidate.

## 5. Rules & Constraints

- **R-D8-09-01:** Every scheduled task must be registered in the Scheduled Task Registry. Unregistered scheduled tasks are unauthorized and must be either registered or removed upon discovery.
- **R-D8-09-02:** Every scheduled task must have a complete specification per D8-007.
- **R-D8-09-03:** Every scheduled task must implement a locking mechanism that prevents concurrent execution of the same task.
- **R-D8-09-04:** Every scheduled task must log its start time, completion time, success/failure status, and sufficient detail to diagnose failures.
- **R-D8-09-05:** Every scheduled task must write a heartbeat timestamp on successful completion, in the location monitored by the independent monitoring system per D8-008.
- **R-D8-09-06:** The Scheduled Task Registry must be audited quarterly. Discrepancies between the registry and actual system state must be resolved within 7 days.
- **R-D8-09-07:** Scheduled tasks must not be clustered. No more than two tasks may be scheduled to begin within the same 15-minute window unless they are explicitly documented as non-conflicting in the registry.
- **R-D8-09-08:** New scheduled tasks require a conflict analysis against all existing scheduled tasks before deployment. The analysis is documented in the task's specification.
- **R-D8-09-09:** The combined peak resource consumption of all concurrent scheduled tasks must not exceed 70% of any system resource. If this threshold is exceeded, the schedule must be revised.

## 6. Failure Modes

- **Registry drift.** Scheduled tasks are added, modified, or removed without updating the registry. The registry becomes inaccurate. The operator does not know what is actually running. Mitigation: the quarterly audit (R-D8-09-06) catches drift. The discipline of registry-first operation -- updating the registry before modifying the actual scheduled task -- prevents it.

- **Temporal collision.** Multiple tasks run simultaneously and compete for resources. Tasks that normally complete in minutes take hours. Downstream dependencies are violated. Mitigation: the conflict analysis required by R-D8-09-08 and the clustering prohibition in R-D8-09-07.

- **Lock failure.** A task's locking mechanism fails (stale lock file from a crashed process). The task does not run because it believes another instance is running. Or two instances run simultaneously because the locking mechanism has a race condition. Mitigation: lock mechanisms must handle stale locks (using process ID verification or timestamp-based expiration). Lock behavior is tested during the quarterly audit.

- **Silent task disappearance.** A task stops running -- its cron entry is accidentally removed, its timer unit is disabled, or a system update overwrites its configuration. The task no longer executes, but no alert fires because the monitoring system only checks for failure, not for absence. Mitigation: the heartbeat mechanism (R-D8-09-05) and independent monitoring (D8-008) detect the absence of expected heartbeats.

- **Cascading dependency failure.** Task A fails. Task B runs on Task A's missing output and produces incorrect results. Task C compounds the error. Mitigation: downstream tasks must verify inputs before executing, failing explicitly on missing or invalid data.

- **Missed execution accumulation.** After maintenance downtime, multiple missed tasks attempt simultaneous execution, overwhelming the system. Mitigation: each task's specification documents missed-execution behavior. Catch-up execution should be staggered.

## 7. Recovery Procedures

1. **If the registry is discovered to be inaccurate:** Perform an immediate reconciliation. List all actual scheduled tasks from system inspection (crontab -l, systemctl list-timers, /etc/cron.d/ contents). Compare against the registry. For each discrepancy: if the task is legitimate but unregistered, register it. If the task is registered but not running, determine why and restore or retire it. If the task is unauthorized, investigate its origin and remove or legitimize it.

2. **If a temporal collision has caused failures:** Identify all tasks involved. Determine which task failed and which completed. Reschedule conflicting tasks with sufficient separation. Rerun any failed tasks manually after confirming the collision is resolved. Update the registry with the new schedule. Document the collision and the resolution.

3. **If a stale lock is preventing task execution:** Verify that no instance of the task is actually running (check process list). If the lock is confirmed stale, remove it manually. Investigate why the lock was not cleaned up (process crash? system reboot during execution?). If the lock mechanism does not handle staleness gracefully, improve it.

4. **If multiple missed tasks are attempting to run simultaneously after a maintenance window:** Prevent the catch-up flood by temporarily suspending all scheduled tasks. Review which tasks need to catch up and which should skip to their next scheduled time. Run catch-up tasks one at a time in priority order. Verify each before proceeding to the next. Re-enable scheduled tasks once catch-up is complete.

5. **If a cascading dependency failure has occurred:** Stop all downstream tasks. Identify the root cause (the first task in the chain that failed). Fix the root cause. Rerun the chain from the failed task forward, verifying each step. Review dependency declarations to confirm they are explicit and enforced.

## 8. Evolution Path

- **Years 0-5:** Start with very few scheduled tasks. Establish the registry and the audit discipline before the task inventory grows. Build the standard wrapper script/function and use it from the beginning. The investment in infrastructure now prevents chaos later.

- **Years 5-15:** The scheduled task inventory should be stable and well-understood. The quarterly audit should be routine. New tasks should be rare and carefully justified. The primary challenge is maintaining accuracy as the underlying platform evolves (cron to systemd transitions, OS upgrades, hardware changes).

- **Years 15-30:** Review whether the scheduling mechanisms still serve the institution. Platform changes may offer better alternatives. The registry and the architectural principles should survive; the implementation may change entirely.

- **Years 30-50+:** The scheduled task corpus is inherited by a successor. The registry is their map. The specifications are their guide. The quality of the registry and specifications written in years 0-15 determines whether the successor can maintain the system or must rebuild it.

- **Signpost for revision:** If the number of scheduled tasks exceeds what the operator can audit by hand in a single quarterly session (roughly 20-30 tasks), the institution either has too many scheduled tasks or needs tooling to assist with the audit. The first response should be consolidation -- can tasks be combined? Can some tasks be retired? Tooling is the second response, and it is itself automation that must be justified per D8-001.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
I have administered systems where the crontab grew to hundreds of entries over years, each one added by a different person for a different reason, with no registry, no conflict analysis, and no audit discipline. The result was a system where nobody knew what ran when, where failures were intermittent and mysterious, and where the only safe approach to changing anything was to change nothing.

The registry may feel like overhead for an institution with five scheduled tasks. It is not overhead. It is the mechanism that ensures the institution still has five well-understood scheduled tasks in year fifteen, rather than thirty poorly understood ones.

One specific warning: resist the urge to schedule tasks at "nice" times like midnight, 3 AM, or the top of the hour. Every system administrator in history has chosen these times, and they are the most congested scheduling slots on any system that has accumulated tasks over time. Schedule at odd times -- 2:37 AM, 4:13 AM -- and you will have fewer collisions.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 2: Integrity Over Convenience; Principle 3: Transparency)
- CON-001 -- The Founding Mandate (comprehensibility, single-operator constraint)
- OPS-001 -- Operations Philosophy (complexity budget, operational tempo)
- D8-001 -- Automation Restraint Doctrine (justification framework, annual review, manual fallback)
- D8-002 -- Agent Design Principles (minimal dependencies)
- D8-007 -- Automation Specification Language (specification template for each scheduled task)
- D8-008 -- Monitoring Automation (heartbeat checking, independent verification)
- D8-011 -- Automation Observability Standards (logging requirements for scheduled tasks)
- D6-001 -- Data Philosophy (data tier classification for registry and logs)
- D10-002 -- Daily Operations Runbook (integration of scheduled task verification)

---

---

# D8-010 -- Automation Testing and Validation

**Document ID:** D8-010
**Domain:** 8 -- Automation & Agents
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, SEC-001, OPS-001, D8-001, D8-002, D8-003, D8-006, D8-007, D8-008, D8-009, D6-001
**Depended Upon By:** D8-011. All articles that create, deploy, or modify automation.

---

## 1. Purpose

This article defines how automation is tested before deployment and validated after. In conventional software engineering, testing is supported by extensive infrastructure: continuous integration servers, staging environments that mirror production, automated test suites that run on every commit, and teams of people whose sole job is quality assurance. An air-gapped, single-operator institution has none of these. Testing must be achieved with the resources available: a single operator, limited hardware, no network connectivity, and no external services.

This constraint does not reduce the need for testing. It intensifies it. In a conventional environment, a failed deployment can be rolled back quickly because the team is large and the infrastructure is redundant. In this institution, a failed automation can corrupt data, disrupt operations, and consume the operator's time for hours or days -- time that cannot be recovered. Per D8-001, R-D8-06, untested automation must not be deployed in a production capacity. This article defines what "tested" means, how testing is performed, and how validation continues after deployment.

The article also addresses the unique challenges of testing in an air-gapped environment: the inability to pull test dependencies from the internet, the difficulty of simulating production conditions on limited hardware, and the need for testing approaches that a single operator can perform without assistance.

## 2. Scope

This article covers:

- The testing philosophy for institutional automation.
- Test environments: how to create and maintain them with limited resources.
- Pre-deployment testing: unit tests, integration tests, and acceptance tests.
- Post-deployment validation: verifying that automation works correctly in production.
- Regression testing: ensuring changes do not break existing functionality.
- The Automation Acceptance Test Procedure: the formal gate between development and production.
- Canary deployment strategies adapted for single-node environments.
- Testing scheduled tasks and time-dependent behavior.

This article does not cover:

- The automation specification (see D8-007 for the specification that defines expected behavior against which testing verifies).
- Monitoring deployed automation (see D8-008 for ongoing monitoring after deployment).
- The governance approval for deployment (see D8-010 from the original Stage 1 numbering -- refer to the Automation Governance article in this institution's documentation corpus).
- General system testing unrelated to automation (see D5-002 for OS verification).

## 3. Background

### 3.1 Why Testing Is Hard in This Environment

Testing automation requires comparing actual behavior against expected behavior. In a networked environment, expected behavior is often defined by interaction with external services, APIs, and data sources. In an air-gapped environment, all of these must be simulated or substituted. This means the test environment is always an approximation of production, and the operator must understand exactly what the approximation covers and what it does not.

Additionally, the operator is the developer, the tester, the reviewer, and the deployer. The psychological challenge is significant: the same person who wrote the automation is evaluating whether it works. Every cognitive bias that makes people poor at evaluating their own work applies. The testing procedures in this article are designed to mitigate this bias through structure: checklists that force the operator to verify specific behaviors, acceptance criteria defined before testing begins, and a cooling-off period between writing and testing.

### 3.2 The Cost of Not Testing

The cost of deploying untested automation is paid in one of three currencies: corrupted data, lost time, or lost trust. Corrupted data may not be recoverable. Lost time is irreplaceable for a single operator. Lost trust -- the operator's confidence that their systems work as intended -- degrades the operator's effectiveness and morale.

No amount of testing proves the absence of bugs. Testing demonstrates the presence of correct behavior under tested conditions. The gap between tested conditions and all possible conditions is where failures live. This article does not promise that tested automation will never fail. It promises that tested automation will have been verified under a defined set of conditions, that the conditions are documented, and that the gaps are acknowledged.

## 4. System Model

### 4.1 The Test Environment

The institution maintains a test environment for automation validation. The test environment is architecturally separate from production but mirrors its essential characteristics:

**Minimum test environment requirements:**
- A separate user account or container on the same hardware, or a separate machine if available.
- A copy of the production directory structure relevant to the automation being tested.
- Test data that is structurally identical to production data but clearly marked as test data.
- The same software versions (OS, interpreter, libraries) as production.
- Access to the same configuration files (or test copies thereof).

**Test environment limitations (documented, not hidden):**
- The test environment may have less storage, less memory, or less CPU than production. Performance testing may not be fully representative.
- The test environment cannot replicate production timing precisely. Scheduled tasks tested at different times may behave differently in production.
- The test environment does not contain production data (to prevent accidental modification). Testing with representative synthetic data introduces the risk that edge cases in real data are missed.

The test environment's configuration is documented and maintained alongside the automation specifications. Drift between the test and production environments is checked during the quarterly scheduling audit per D8-009.

### 4.2 Pre-Deployment Testing

Before an automation is deployed to production, it must pass three levels of testing:

**Level 1: Component Testing.**
Each logical component is tested in isolation with known inputs and expected outputs. Requirements: at least one test case per input, per output, and for invalid or missing input. Test cases are documented and repeatable.

**Level 2: Integration Testing.**
The automation is tested as a whole, interacting with its dependencies. Requirements: the automation executes in the test environment with all dependencies present, produces expected outputs, logs correctly, locks correctly (for scheduled tasks), and handles simulated failures (removed files, read-only filesystems, exhausted memory).

**Level 3: Acceptance Testing.**
The automation is executed simulating production operation as closely as possible. The Automation Acceptance Test Procedure (Section 4.4) governs this level.

### 4.3 Post-Deployment Validation

Passing pre-deployment testing does not end the testing obligation. After deployment to production, the automation enters a validation period during which it is monitored with heightened scrutiny:

**The validation period:**
- Duration: a minimum of three full execution cycles, or 7 days, whichever is longer.
- During validation, the automation runs at its intended schedule and in its intended mode.
- The operator reviews every execution's logs (not just the monitoring dashboard) during the validation period.
- Any failure, anomaly, or unexpected behavior during validation triggers immediate investigation.
- At the end of the validation period, the operator formally accepts the automation into production status or suspends it for further investigation.

**Validation checklist:** All outputs match expected values. Execution duration and resource consumption are within documented bounds. No unanticipated errors or warnings. Logging and heartbeat mechanisms functioning. Monitoring system (D8-008) correctly reports status. No interference with other automations or system functions.

### 4.4 The Automation Acceptance Test Procedure

The AATP is the formal gate between development and production. No automation may enter production without passing the AATP. The procedure is:

1. **Specification review.** Confirm that the automation's specification (per D8-007) is complete and current.

2. **Test evidence review.** Confirm that Level 1 and Level 2 testing have been completed and that test results are documented.

3. **Acceptance criteria definition.** Before the acceptance test begins, the operator writes down the specific criteria that constitute success. These criteria are derived from the specification but stated in concrete, verifiable terms.

4. **Acceptance test execution.** The automation is run in the test environment under conditions that simulate production as closely as possible. The operator records all observations.

5. **Result evaluation.** The operator compares observations against the acceptance criteria. Every criterion must be met.

6. **Manual fallback verification.** The manual fallback procedure (specification Section 8) is executed at least once to verify it works and to confirm the operator can perform the task without the automation.

7. **Decision.** The operator formally records one of three decisions:
   - **Accept:** The automation is cleared for production deployment. Proceed to validation period.
   - **Accept with conditions:** The automation is cleared for deployment but with specific monitoring or restriction requirements beyond the standard. Conditions are documented.
   - **Reject:** The automation is not cleared for deployment. The reasons are documented. The automation returns to development.

8. **Record.** The AATP results are recorded in the automation's specification Change History and in the institutional decision log.

### 4.5 Regression Testing

When an existing automation is modified, the modification must not break existing functionality. Regression testing verifies this:

- All existing component tests are re-executed after the modification.
- All existing integration tests are re-executed.
- The AATP is re-executed with acceptance criteria that include both the new behavior and the preservation of existing behavior.
- If the modification changes the automation's interaction with other automations, those other automations are also tested for impact.

Regression test results are documented in the specification Change History.

### 4.6 Canary Deployment in a Single-Node Environment

Traditional canary deployment routes a small percentage of traffic to the new version while the old version handles the majority. In a single-node institution, this is not directly applicable. The adapted approach:

**Parallel execution canary:** The new version of the automation is deployed alongside the old version, configured to run at a different time (offset by a small interval). Both versions execute independently. The operator compares their outputs. If the new version produces identical results, the old version is retired and the new version assumes the production schedule. If results differ, the operator investigates before switching.

**Shadow execution canary:** The new version runs in the test environment on a copy of production data (or a recent backup). Its outputs are compared against the production automation's outputs. Discrepancies are investigated before the new version replaces the old.

**Time-limited canary:** The new version is deployed to production but with a defined evaluation period (typically three execution cycles). If any anomaly occurs during the evaluation period, the old version is restored immediately using the rollback procedure documented in the specification.

### 4.7 Testing Time-Dependent Behavior

Scheduled tasks and time-dependent automations present specific testing challenges:

- **Clock manipulation.** Where possible, the test environment uses a configurable clock (faketime or similar) to test behavior at specific times, including edge cases like midnight, end of month, leap years, and daylight saving transitions.
- **Missed execution testing.** The automation is tested for correct behavior after a simulated missed execution: the task's scheduled time passes without execution, and the operator verifies that the catch-up behavior matches the specification.
- **Duration boundary testing.** The automation is tested under conditions that cause it to exceed its expected duration, verifying that the locking mechanism prevents overlapping execution and that the monitoring system detects the overrun.

## 5. Rules & Constraints

- **R-D8-10-01:** No automation may be deployed to production without passing the Automation Acceptance Test Procedure defined in Section 4.4.
- **R-D8-10-02:** All test results must be documented and stored alongside the automation specification. Test documentation is Tier 2 data per D6-001.
- **R-D8-10-03:** The test environment must be maintained and verified quarterly to ensure it accurately reflects the production environment's relevant characteristics.
- **R-D8-10-04:** Modifications to existing automation require regression testing per Section 4.5 before production deployment.
- **R-D8-10-05:** The validation period (Section 4.3) is mandatory and may not be shortened unless the operator documents a compelling justification and accepts the risk in writing.
- **R-D8-10-06:** Manual fallback procedures must be tested as part of every AATP execution (Section 4.4, Step 6).
- **R-D8-10-07:** Testing artifacts (test scripts, test data, test results) must be maintained in version control or equivalent archive. They are part of the automation's documentation.
- **R-D8-10-08:** The operator must wait a minimum of 24 hours between completing development and beginning acceptance testing, to mitigate the bias of testing what you just wrote. This cooling-off period allows the operator to approach the test with fresher eyes.

## 6. Failure Modes

- **Test environment divergence.** The test environment no longer matches production. Automation passes testing but fails in production because of differences the test did not cover. Mitigation: quarterly test environment verification (R-D8-10-03) and documentation of known test environment limitations.

- **Testing complacency.** The operator treats testing as a formality -- running the AATP quickly, not examining outputs carefully, accepting results without genuine scrutiny. Mitigation: the cooling-off period (R-D8-10-08), the requirement to document acceptance criteria before testing, and the annual review process.

- **Edge case blindness.** Testing covers normal conditions thoroughly but misses edge cases that occur rarely in production (disk full, permission changes, unexpected data formats). Mitigation: the specification's Failure Modes section (D8-007) guides test case design. Integration testing specifically includes simulated failures.

- **Regression test erosion.** Over time, existing test cases are not maintained. They fail for reasons unrelated to the automation's correctness (test data is outdated, test paths have changed). The operator stops running them. Mitigation: test maintenance is part of the quarterly audit. Test cases that fail for infrastructure reasons are fixed, not abandoned.

- **Bias in self-testing.** The same person who wrote the automation designs the tests. The tests verify what the developer expected, not what the institution needs. Mitigation: acceptance criteria derived from the specification (which was written before the code) rather than from the code itself. The cooling-off period helps.

## 7. Recovery Procedures

1. **If untested automation has been deployed to production:** Do not panic, but do not ignore it. Immediately begin a post-hoc validation. Execute the AATP against the running automation. If it passes, document the gap in process and reinforce the discipline. If it fails, suspend the automation and activate the manual fallback. Document the incident.

2. **If a tested automation fails in production:** Suspend the automation. Activate the manual fallback. Compare the production failure against the test results. Determine what condition exists in production that was not present in testing. Add this condition to future test cases. Fix the automation. Re-execute the full AATP. Redeploy with the standard validation period.

3. **If the test environment has diverged significantly from production:** Halt automation development until the test environment is reconciled. Compare the test and production environments systematically. Document all differences. Correct the test environment to match production. Re-test any automation that was tested during the period of divergence.

4. **If regression testing reveals a break caused by a recent modification:** Roll back the modification. Restore the previous version of the automation. Investigate the cause of the regression. Fix the modification to preserve existing behavior. Re-test. Document the regression and the fix.

5. **If the cooling-off period is consistently being skipped:** This is a discipline failure. Acknowledge it. Consider whether the operator is under excessive time pressure (an operational concern per OPS-001) and address the root cause. The cooling-off period exists because testing your own work immediately is unreliable. Skipping it trades a small time savings for a significant increase in deployment risk.

## 8. Evolution Path

- **Years 0-5:** Testing discipline is being established. The first automations will be simple, and the temptation to skip formal testing will be strong. Resist it. The habit formed now persists for decades. Build the test environment early. Write test cases even for trivial automations. The discipline matters more than the test coverage.

- **Years 5-15:** The test suite should be growing alongside the automation inventory. Component tests from early automations should be stable and running reliably. The primary challenge is maintaining the test environment as the production environment evolves. This is also the period when the first significant regressions will be caught by the test suite, validating the investment.

- **Years 15-30:** The test corpus is substantial. Maintaining it is a real burden. Consider whether test automation (a tool that runs all tests in sequence and reports results) is justified per D8-001's criteria. The irony of automating the testing of automation is acknowledged -- but if the testing process itself is well-specified and well-understood, it is a reasonable candidate for automation.

- **Years 30-50+:** A successor inherits both the automation and its test suite. The test suite is their safety net -- it tells them whether their modifications break anything. Poorly maintained tests are a false safety net. Well-maintained tests are one of the most valuable inheritances the institution can provide.

- **Signpost for revision:** If the testing process consistently takes longer than the development process, the testing procedures may be disproportionate to the complexity of the automations being tested. Simplify without abandoning rigor. If automations consistently pass all tests and then fail in production, the test coverage is inadequate and must be expanded.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The 24-hour cooling-off period will be the most frequently resisted requirement in this article. The operator will finish writing an automation at 11 PM, want to deploy it, and consider the cooling-off period an unnecessary delay. I know this because I am that operator, and I have felt that impulse many times.

The impulse is wrong. Every significant bug I have deployed to production was deployed in the same session I wrote the code. Every time I waited a day, I found something I had missed. Not always a bug -- sometimes a clearer way to write the code, sometimes a failure mode I had not considered, sometimes a test case I had forgotten. The 24-hour delay is not idle time. It is the time your subconscious spends reviewing your work.

One more thing: the manual fallback test (Step 6 of the AATP) will feel redundant. You wrote the fallback. You know it works. But you do not know it works in the current state of the system, with the current data, under the current conditions. Test it. The five minutes you spend confirming the fallback works is the cheapest insurance in the institution.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 2: Integrity Over Convenience)
- CON-001 -- The Founding Mandate (single-operator constraint, comprehensibility)
- SEC-001 -- Threat Model and Security Philosophy (defense in depth)
- OPS-001 -- Operations Philosophy (complexity budget, operational tempo)
- D8-001 -- Automation Restraint Doctrine (R-D8-06: untested automation must not be deployed; comprehensibility criterion)
- D8-002 -- Agent Design Principles (testability as a design principle)
- D8-006 -- Agent Lifecycle (testing requirements at each lifecycle stage)
- D8-007 -- Automation Specification Language (specification defines expected behavior for testing)
- D8-008 -- Monitoring Automation (post-deployment monitoring during validation)
- D8-009 -- Scheduled Task Architecture (testing time-dependent and scheduled behavior)
- D8-011 -- Automation Observability Standards (verifying observability during testing)
- D5-002 -- OS Maintenance Procedures (test environment maintenance)
- D6-001 -- Data Philosophy (data tier classification for test documentation)

---

---

# D8-011 -- Automation Observability Standards

**Document ID:** D8-011
**Domain:** 8 -- Automation & Agents
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, SEC-001, OPS-001, D8-001, D8-002, D8-003, D8-007, D8-008, D8-009, D8-010, D6-001, D6-008
**Depended Upon By:** All articles in Domain 8 that reference logging, metrics, or status reporting. D10-009. All operational articles that consume automation state information.

---

## 1. Purpose

This article defines what every automation in the institution must expose about its internal state. Observability is the quality of a system that allows an external observer to understand what the system is doing, what it has done, and whether it is functioning correctly -- without modifying the system or accessing its internal code. It is D8-001 Criterion 4 (Observability) made concrete: specific standards for logging, metrics, status reporting, and self-description that every automation must implement.

Observability is not monitoring. Monitoring (D8-008) is what the institution does to watch automation from the outside. Observability is what the automation provides to make monitoring possible. A well-monitored but poorly observable automation is a black box with sensors taped to its exterior -- the sensors can tell you that something went wrong but not what or why. A well-observable automation is a glass box: its internal state is visible, its decisions are traceable, and its failures are diagnosable.

This article establishes the standards that make every automation in the institution a glass box. The standards are deliberately opinionated: they specify formats, required fields, naming conventions, and retention requirements. This rigidity serves comprehensibility. When every automation logs in the same format with the same fields, an operator reading logs from any automation in the institution encounters a familiar structure. When every automation exposes the same health metrics, the monitoring system (D8-008) can aggregate and compare them uniformly. Standardization is the foundation of institutional observability.

## 2. Scope

This article covers:

- Logging standards: what must be logged, in what format, with what fields.
- Metrics requirements: what quantitative data every automation must expose.
- Status reporting: how every automation communicates its current state.
- The observability audit checklist: how to verify that an automation meets observability standards.
- How to build observability into automation from the start, rather than retrofitting it.
- Log retention, rotation, and archival requirements.

This article does not cover:

- The monitoring architecture that consumes observability data (see D8-008).
- The specification template that documents observability requirements (see D8-007, Section 10).
- Testing of observability (see D8-010, which includes observability verification in its test procedures).
- General system logging unrelated to automation (see D5-005).
- Data retention policies beyond automation-specific requirements (see D6-001, D6-011).

## 3. Background

### 3.1 The Black Box Problem

D8-001 identifies the black box problem as existential for a multi-generational institution: automation that cannot be understood by its current maintainer is worse than automation that does not exist, because it creates dependency without comprehension. Observability is the primary defense against the black box problem. An automation that logs every decision, reports every outcome, and exposes every metric may still be difficult to understand -- but it is not opaque. The operator has evidence to reason from, even if they do not fully understand the mechanism that produced the evidence.

### 3.2 Why Standards, Not Guidelines

Guidelines are optional by nature. Standards are mandatory. The distinction matters because observability is most needed precisely when the operator is least inclined to invest in it: during rapid development, during emergency fixes, during "temporary" automations that become permanent. If observability is a guideline, it will be the first thing sacrificed to urgency. If it is a standard -- verifiable, auditable, enforced -- it will be present when it is needed.

### 3.3 The Log as Historical Record

In a fifty-year institution, automation logs are not just diagnostic tools. They are historical records of what the institution's automated systems did, when, and with what results. They enable future operators to understand not just the current state but the history of automation behavior. Per D6-001, automation logs are classified as Tier 2 data (Operational Data) with a defined retention period. Logs are not ephemeral. They are institutional memory.

## 4. System Model

### 4.1 The Logging Standard

Every automation must produce structured logs conforming to the following standard:

**Log format:** Each log entry is a single line of text (no multi-line entries except for stack traces or error details, which are indented under their parent entry). Each line contains the following fields, separated by a consistent delimiter (pipe character "|" is the institutional standard):

```
TIMESTAMP | AUTOMATION-ID | EXECUTION-ID | LEVEL | COMPONENT | MESSAGE
```

**Field definitions:**

- **TIMESTAMP:** ISO 8601 format with timezone, to second precision minimum. Example: 2026-02-16T03:15:42+00:00. Consistent timestamps enable chronological interleaving of logs from different automations.

- **AUTOMATION-ID:** The SPEC-ID from the automation's specification (per D8-007). Enables filtering logs by automation.

- **EXECUTION-ID:** A unique identifier for this execution run. Format: YYYYMMDD-HHMMSS-[SHORT-RANDOM]. Enables grouping all log entries from a single execution.

- **LEVEL:** One of: TRACE (detailed internal state, may be disabled in production), DEBUG (useful for diagnosis, not needed normally), INFO (normal operational events), WARN (unexpected but recovered), ERROR (something went wrong, operator should investigate), FATAL (cannot continue, terminating, immediate attention required).

- **COMPONENT:** The logical component within the automation that produced the entry. For simple automations, this is the automation name. For complex automations, this identifies the stage (e.g., "input-validation", "processing", "output-write").

- **MESSAGE:** Human-readable description of what happened. Messages must be self-contained -- understandable without referring to source code. Bad: "Condition met, proceeding." Good: "Backup file /data/backups/2026-02-16.tar.gz verified, size 4.2GB matches expected range 3.5-5.0GB."

### 4.2 Required Log Events

Every automation must log, at minimum, the following events at the specified levels:

**At execution start (INFO):** Timestamp, automation ID and version, execution ID, input summary (sufficient to reproduce the execution), and environment summary (disk space, dependency availability).

**At each significant processing milestone (INFO):** Phase or step completed, intermediate results or metrics.

**At execution completion (INFO):** Success or failure status, output summary (what was produced, where), execution duration, resource consumption summary.

**At any error or unexpected condition (WARN or ERROR):**
- Description of the condition and what the automation was doing when it occurred.
- Impact on the automation's output or continuation.
- Action taken (retried, skipped, failed, continued with degraded output).

**At termination due to unrecoverable error (FATAL):**
- Description of the fatal condition and state at time of failure.
- What outputs may be in an inconsistent state.
- Recommended recovery action for the operator.

### 4.3 Metrics Requirements

Every automation must expose the following metrics, updated at each execution:

**Execution metrics:** Last execution start time, end time, duration (seconds), and status (success/failure/partial). Cumulative execution count, success count, failure count, and success rate (as percentage).

**Resource metrics (per execution):** Peak memory usage, CPU time consumed, disk I/O (bytes read and written), and execution duration compared to historical average (percentage deviation).

**Output metrics (per execution):** Output size or count, output location, and output validation result (pass/fail, if applicable).

Metrics are written to a standardized metrics file alongside the automation's logs. The metrics file is a simple key-value text file, overwritten at each execution with the most recent values and appended to a historical metrics log:

```
# Metrics for AUTO-D6-003 (Backup Rotation)
# Generated: 2026-02-16T03:20:15+00:00
last_start=2026-02-16T03:15:42+00:00
last_end=2026-02-16T03:20:14+00:00
last_duration_seconds=272
last_status=success
total_executions=847
total_successes=845
total_failures=2
success_rate=99.76
peak_memory_mb=128
cpu_seconds=45.3
bytes_read=4521789440
bytes_written=4521789440
output_size_bytes=4521789440
output_path=/data/backups/2026-02-16.tar.gz
output_valid=true
```

### 4.4 Status Reporting

Every automation must provide a mechanism for the operator to query its current status at any time. The status report answers the question: "Is this automation healthy right now?"

The status mechanism is a status file written by the automation and readable by the operator and by the monitoring system (D8-008). The status file contains:

```
# Status for AUTO-D6-003 (Backup Rotation)
# Updated: 2026-02-16T03:20:15+00:00
automation_id=AUTO-D6-003
automation_name=Backup Rotation
status=healthy
last_execution=2026-02-16T03:20:14+00:00
last_result=success
next_scheduled=2026-02-17T03:15:00+00:00
active_alerts=none
dependencies_status=all_available
notes=
```

Status values are: **healthy** (last execution succeeded, no active alerts), **degraded** (last execution succeeded with warnings, or a non-critical dependency is unavailable), **failed** (last execution failed), **unknown** (no execution has occurred since deployment or since a status reset).

The status file is updated at the conclusion of each execution and whenever the automation detects a change in its dependencies. The monitoring system (D8-008) consumes status files as part of its health assessment.

### 4.5 Building Observability In

Observability cannot be retrofitted effectively. It must be designed into automation from the start.

**Before writing any automation logic, write the logging framework first.** The first lines of any automation establish the log file path, format, execution ID, and start entry. The last lines write the completion entry, metrics file, and status file. Logic goes between these bookends. This forces the developer to think in terms of observable events before implementation.

**The observability template:** The institution maintains a shell function library (or equivalent) that provides:
- `obs_init` -- initializes logging, sets execution ID, writes start entry.
- `obs_log` -- writes a log entry at a specified level.
- `obs_metric` -- records a metric value.
- `obs_complete` -- writes completion entry, metrics file, status file, heartbeat.
- `obs_fail` -- writes failure entry, updates status to failed, writes heartbeat with failure flag.

Every automation uses this library. The library is maintained as institutional infrastructure per D5-006.

### 4.6 The Observability Audit Checklist

The following checklist is used during the quarterly automation audit (D8-009, Section 4.4) and the annual automation review (D8-001, R-D8-07) to verify that each automation meets observability standards:

```
OBSERVABILITY AUDIT CHECKLIST
Automation: [SPEC-ID]  Auditor: [Name]  Date: [Date]

LOGGING:
[ ] Log file exists at documented location
[ ] Entries conform to standard format (Section 4.1)
[ ] All required log events present (Section 4.2)
[ ] Messages are human-readable and self-contained
[ ] Log rotation configured and functioning
[ ] Retention meets D6-001 requirements

METRICS:
[ ] Metrics file exists at documented location
[ ] All required metrics present (Section 4.3)
[ ] Values plausible and consistent with log entries
[ ] Historical metrics log maintained

STATUS:
[ ] Status file exists, is current, and accurately reflects actual state

MONITORING INTEGRATION:
[ ] Heartbeat written per D8-008; independent monitor checking it
[ ] Alert thresholds configured

INFRASTRUCTURE:
[ ] Uses standard observability library, or non-standard implementation
    is documented and justified

NOTES: [Free-form observations on observability quality]
```

## 5. Rules & Constraints

- **R-D8-11-01:** Every automation must produce structured logs conforming to the format defined in Section 4.1. Non-conforming log formats are a compliance failure.
- **R-D8-11-02:** Every automation must log all required events defined in Section 4.2. Missing required log events are a compliance failure.
- **R-D8-11-03:** Every automation must expose all required metrics defined in Section 4.3.
- **R-D8-11-04:** Every automation must maintain a status file conforming to Section 4.4, updated at the conclusion of each execution.
- **R-D8-11-05:** Log entries must be human-readable without reference to source code. Codes, abbreviations, and internal identifiers must be accompanied by human-readable explanations.
- **R-D8-11-06:** Automation logs are classified as Tier 2 data per D6-001. They must be retained for a minimum of one year. Historical metrics logs must be retained for the lifetime of the automation.
- **R-D8-11-07:** Log rotation must be configured for every automation. No log file may grow without bound. The rotation policy must be documented in the automation's specification (D8-007, Section 10).
- **R-D8-11-08:** New automation must use the standard observability library (Section 4.5). Exceptions require documented justification and must implement equivalent functionality.
- **R-D8-11-09:** The observability audit checklist (Section 4.6) must be completed for every automation during the quarterly audit. Automations that fail the audit must be remediated within 14 days.
- **R-D8-11-10:** Log messages must never contain sensitive data (passwords, cryptographic keys, personal information beyond what is necessary for operational diagnosis). This constraint is absolute and overrides completeness.

## 6. Failure Modes

- **Log format divergence.** Different automations use different log formats. Aggregating or comparing logs across automations becomes impossible. The operator must learn each automation's idiosyncratic format. Mitigation: R-D8-11-01 mandates the standard format. The quarterly audit verifies compliance.

- **Log verbosity extremes.** Either too verbose (generating gigabytes of unreadable logs) or too terse (logging only errors, making pre-error context invisible). Mitigation: the required log events (Section 4.2) define the minimum. TRACE and DEBUG may be disabled in production, but INFO, WARN, ERROR, and FATAL must always be active.

- **Metrics staleness.** The metrics file is not updated because the automation failed before reaching the metrics-writing step. The monitoring system reads stale metrics and reports health when the automation is actually broken. Mitigation: the monitoring system (D8-008) checks metrics timestamps, not just values. A stale metrics file triggers a warning.

- **Status file corruption.** The status file is partially written (due to a crash during write). The monitoring system reads garbled data. Mitigation: status files should be written atomically (write to a temporary file, then rename). The monitoring system validates status file format before consuming it.

- **Observability as afterthought.** Observability is added as a compliance task after development. Logs are generic, metrics superficial, status reporting perfunctory. Mitigation: Section 4.5 mandates building observability first. The AATP (D8-010) includes observability verification.

- **Sensitive data in logs.** A log message inadvertently includes passwords, keys, or personal information. Mitigation: R-D8-11-10 prohibits sensitive data. Developers must sanitize variable content before logging.

- **Log storage exhaustion.** Log rotation fails. Logs fill the filesystem. Other automations fail. Mitigation: R-D8-11-07 mandates rotation. The independent monitor (D8-008) tracks disk space. The audit verifies rotation.

## 7. Recovery Procedures

1. **If log format non-compliance is discovered:** Identify all non-compliant automations. Prioritize by criticality. Modify each to use the standard format. Verify through the audit checklist. For historical non-standard logs, document the format and maintain a conversion guide.

2. **If logs are missing for a period:** Determine cause. Fix the logging mechanism. Reconstruct the gap period from other evidence (system logs, outputs, monitoring records). Document the gap.

3. **If sensitive data is found in logs:** Redact from all copies. Rotate any exposed credentials. Fix the automation. Review all automations for similar vulnerabilities. Document per SEC-001.

4. **If log storage is exhausted:** Rotate or compress logs immediately. Verify all automations have functioning rotation. Review retention policies. Expand storage capacity if retention requirements demand it.

5. **If the observability library is inadequate:** Document the inadequacy. Extend through the standard development process (D8-010). Update all consuming automations. If log format changes, update monitoring system parsing and document the version change.

## 8. Evolution Path

- **Years 0-5:** The observability standards are new. The observability library will be built and tested. Expect to iterate on log format and metrics definitions as you learn what is genuinely useful for diagnosis versus merely comprehensive.

- **Years 5-15:** Standards should be stable, the library mature. The primary challenge is maintaining consistency as automations evolve. Historical metrics logs become valuable for trend analysis.

- **Years 15-30:** The log corpus spans a decade or more. The log format should remain stable to enable historical analysis. If the format must change, maintain a format version field and a conversion tool.

- **Years 30-50+:** A successor inherits logs, metrics, and status reports spanning decades. Because every automation used the same format, the successor can read any log from any era. The observability standards are the successor's window into automated history.

- **Signpost for revision:** If the log format consistently lacks diagnostic information, extend it. If metrics include values never consulted, simplify. But resist removing fields merely because they seem unused -- their value may appear only during rare failure investigations.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The most contentious standard in this article will be the requirement for human-readable log messages (R-D8-11-05). It is faster to write `log("ERR:BKP:CHK:FAIL:3")` than `log("ERROR: Backup checksum verification failed for /data/backups/2026-02-16.tar.gz -- computed SHA-256 does not match stored value. This is attempt 3 of 3. Backup file may be corrupt.")`. The first takes less disk space, less development time, and less thought. The second can be understood by someone who has never seen the source code, twenty years from now, in a crisis, at 3 AM.

Write the second kind. Always.

I also want to note that the pipe-delimited format was chosen over JSON deliberately. JSON is more structured and more easily parsed by machines. But it is less readable by humans when viewed in a terminal with `cat` or `less`. This institution prioritizes human readability over machine parseability. If a future operator needs machine parsing, a pipe-delimited format is trivially convertible to JSON, CSV, or any other structured format. The reverse is not true -- JSON logs are cumbersome to read raw.

One final observation: the requirement to write the observability framework before the automation logic (Section 4.5) will feel backward to experienced developers. It is not. It is the automation equivalent of test-driven development: define how you will observe the system before you build the system. What you choose to observe shapes how you think about what the system does. Build the glass box first, then put the mechanism inside it.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 3: Transparency of Operation)
- CON-001 -- The Founding Mandate (comprehensibility, longevity)
- SEC-001 -- Threat Model and Security Philosophy (sensitive data handling, incident reporting)
- OPS-001 -- Operations Philosophy (documentation-first principle, operational tempo)
- GOV-001 -- Authority Model (governance for standards changes)
- D6-001 -- Data Philosophy (data tier classification for logs; Tier 2: Operational Data)
- D6-003 -- Format Longevity Doctrine (plain text format for logs and metrics)
- D6-008 -- Metadata Standards (metadata requirements for log files)
- D6-011 -- Data Retention Schedules (retention requirements for automation logs)
- D8-001 -- Automation Restraint Doctrine (Criterion 4: Observability; R-D8-04: all automation must be observable)
- D8-002 -- Agent Design Principles (observability as a design principle)
- D8-007 -- Automation Specification Language (specification Section 10: logging and observability)
- D8-008 -- Monitoring Automation (consumes observability data; heartbeat, dashboard, alert integration)
- D8-009 -- Scheduled Task Architecture (logging requirements for scheduled tasks)
- D8-010 -- Automation Testing and Validation (observability verification during testing)
- D5-005 -- Service Health and System Monitoring (general system logging context)
- D10-009 -- Operational Monitoring Procedures (operational consumption of observability data)

---

*End of Stage 4: Specialized Systems -- Automation Advanced Reference*
