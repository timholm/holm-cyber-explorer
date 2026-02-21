# Stage 1: Documentation Framework
# Domains 11-15: Administration, Disaster Recovery, Evolution & Adaptation, Research & Theory, Ethics & Safeguards

**Document ID:** STAGE1-D11-D15
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Initial Framework
**Scope:** Complete Stage 1 output for five domains — domain maps, article lists, priority ordering, dependency graphs, writing schedules, review processes, and maintenance plans.
**Meta-Assumptions:** 50+ year institutional lifespan. Future maintainers will not have access to the original author. Hardware platforms will change. Cultural context will shift. The institution is air-gapped and off-grid. All documentation must be self-contained and self-explanatory.

---

## Article Template Reference

Every article produced under this framework must conform to the following structure:

| Section | Purpose |
|---|---|
| **Title** | Unambiguous, searchable name |
| **Purpose** | Why this article exists; what problem it solves |
| **Scope** | What is covered; what is explicitly excluded |
| **Background** | Historical context, motivation, prior decisions |
| **System Model** | How the system or process works at a conceptual level |
| **Rules & Constraints** | Hard rules, invariants, non-negotiable boundaries |
| **Failure Modes** | What can go wrong; how you know it has gone wrong |
| **Recovery Procedures** | Step-by-step restoration from each failure mode |
| **Evolution Path** | How this article and its subject are expected to change over time |
| **Commentary Section** | Maintainer notes, dissenting opinions, lessons learned (append-only) |
| **References** | Cross-references to other articles, external sources, historical documents |

---
---

# DOMAIN 11: ADMINISTRATION

**Agent Role:** Administration Planner

---

## 11.1 Domain Map

### 11.1.1 Core Purpose

Define and maintain the administrative backbone of the institution: resource management, scheduling, budgeting, procurement, inventory, asset tracking, personnel coordination, and operational governance. Administration is the connective tissue that enables every other domain to function. Without it, the institution degrades into ad-hoc improvisation.

### 11.1.2 Scope

**In scope:**
- Budgeting and financial planning (non-monetary resource accounting included)
- Procurement of hardware, supplies, consumables
- Inventory management and asset lifecycle tracking
- Scheduling of maintenance, reviews, rotations, and ceremonies
- Personnel administration (roles, succession, onboarding, offboarding)
- Record-keeping standards and archival of administrative documents
- Resource allocation and conflict resolution
- Reporting structures and accountability chains
- Administrative tooling (ledgers, logs, forms — physical and digital)
- Operational tempo and cadence definitions

**Out of scope:**
- Technical system architecture (Domain 1-5, assumed)
- Disaster recovery procedures (Domain 12)
- Ethical governance (Domain 15)
- Research methodology (Domain 14)

### 11.1.3 Boundaries

Administration does not make policy. It executes policy set by governance structures. Administration does not determine what is ethical. It implements safeguards defined by Domain 15. Administration does not design systems. It tracks, schedules, and resources the systems designed by technical domains.

The boundary between Administration and Governance is: governance decides *what* and *why*; administration decides *how*, *when*, and *with what*.

### 11.1.4 Key Relationships

| Related Domain | Nature of Relationship |
|---|---|
| Domain 12 (Disaster Recovery) | Administration maintains the resource stockpiles, schedules, and personnel rosters that DR depends on. DR triggers administrative emergency procedures. |
| Domain 13 (Evolution & Adaptation) | Administration tracks version inventories, schedules migrations, and resources transition efforts. Evolution defines what changes; Administration executes the logistics. |
| Domain 14 (Research & Theory) | Administration allocates resources to research, schedules research reviews, and archives research outputs. |
| Domain 15 (Ethics & Safeguards) | Administration enforces ethical constraints in procurement, personnel management, and resource allocation. Ethics audits administrative processes. |
| Domains 1-10 (assumed technical/foundational) | Administration provides scheduling, inventory, and resource services to all technical domains. |

---

## 11.2 Article List

| Article ID | Title | Summary |
|---|---|---|
| ADM-001 | Principles of Institutional Administration | Foundational philosophy: why administration exists, what it optimizes for, how it relates to institutional survival. |
| ADM-002 | Budget and Resource Accounting Framework | How resources (monetary and non-monetary) are tracked, allocated, forecasted, and audited. Includes barter, labor-hours, and material goods. |
| ADM-003 | Procurement Procedures and Supply Chain Management | How goods, hardware, consumables, and services are identified, sourced, acquired, verified, and received. Includes off-grid and degraded-supply-chain scenarios. |
| ADM-004 | Inventory Management and Asset Registry | How every physical and digital asset is catalogued, tracked through its lifecycle, audited, and eventually decommissioned. |
| ADM-005 | Scheduling and Operational Cadence | Master scheduling framework: maintenance windows, review cycles, rotation schedules, seasonal adjustments, ceremony calendars. |
| ADM-006 | Personnel Administration and Role Management | Onboarding, offboarding, role assignment, succession planning, skills tracking, availability management. |
| ADM-007 | Record-Keeping Standards and Archival Procedures | What records must be kept, in what format, for how long, how they are indexed, and how they are retired. |
| ADM-008 | Reporting Structures and Accountability | Who reports to whom, what reports are required, at what frequency, and what triggers escalation. |
| ADM-009 | Administrative Forms, Ledgers, and Templates | Canonical forms for all administrative processes. Physical and digital formats. Version control for forms. |
| ADM-010 | Resource Conflict Resolution | What happens when two domains or functions compete for the same scarce resource. Prioritization frameworks. |
| ADM-011 | Administrative Continuity During Personnel Loss | How administrative functions survive the loss of key personnel, including the loss of all current administrators. |
| ADM-012 | Facilities and Physical Plant Management | Maintenance of physical spaces, utilities, environmental controls, and physical security as administrative functions. |
| ADM-013 | Administrative Self-Audit and Process Review | How administration audits its own processes for drift, waste, corruption, and obsolescence. |

---

## 11.3 Priority Order

Priority is determined by: (a) what must exist before other articles can be written, (b) what is most critical to daily operations, (c) what has the most cross-domain dependencies.

| Priority | Article ID | Rationale |
|---|---|---|
| 1 | ADM-001 | Foundational. All other articles reference the principles established here. |
| 2 | ADM-007 | Record-keeping must be defined before any other process can be documented operationally. |
| 3 | ADM-004 | Inventory is the backbone of procurement, budgeting, and asset tracking. |
| 4 | ADM-002 | Budget framework depends on knowing what assets exist (ADM-004) and how records work (ADM-007). |
| 5 | ADM-006 | Personnel administration is needed before scheduling and reporting can be defined. |
| 6 | ADM-005 | Scheduling depends on personnel (ADM-006) and inventory (ADM-004). |
| 7 | ADM-008 | Reporting depends on personnel (ADM-006) and scheduling (ADM-005). |
| 8 | ADM-003 | Procurement depends on budget (ADM-002) and inventory (ADM-004). |
| 9 | ADM-009 | Forms codify processes defined in earlier articles. |
| 10 | ADM-010 | Conflict resolution requires understanding of budget (ADM-002) and reporting (ADM-008). |
| 11 | ADM-011 | Continuity planning depends on personnel (ADM-006) and all core processes being documented. |
| 12 | ADM-012 | Facilities management depends on inventory (ADM-004), scheduling (ADM-005), and budget (ADM-002). |
| 13 | ADM-013 | Self-audit is the capstone. It requires all other processes to exist before it can audit them. |

---

## 11.4 Dependencies

### 11.4.1 Internal Dependencies (within Domain 11)

```
ADM-001 --> ADM-007 --> ADM-004 --> ADM-002
                                 \-> ADM-006 --> ADM-005 --> ADM-008
                                                         \-> ADM-012
ADM-002 + ADM-004 --> ADM-003
ADM-002 + ADM-008 --> ADM-010
ADM-006 + all core --> ADM-011
All articles --> ADM-009 (forms derived from processes)
All articles --> ADM-013 (audit requires all processes)
```

### 11.4.2 Cross-Domain Dependencies

| Article | Depends On (External) | Depended On By (External) |
|---|---|---|
| ADM-001 | Institutional charter/governance (assumed Domain 1-10) | All domains (administrative principles) |
| ADM-002 | DR-003 (emergency resource reserves) | EVL-006 (migration budgeting), RES-004 (research resource allocation) |
| ADM-003 | ETH-005 (ethical procurement constraints) | DR-005 (replacement part sourcing) |
| ADM-004 | EVL-002 (asset version tracking) | DR-002 (asset inventory for recovery), DR-004 (backup media tracking) |
| ADM-005 | DR-001 (DR drill scheduling) | RES-003 (research scheduling), EVL-004 (migration scheduling) |
| ADM-006 | ETH-008 (personnel ethics standards) | DR-008 (personnel for DR teams), all domains (role definitions) |
| ADM-011 | DR-009 (continuity planning) | DR-009 (mutual dependency — continuity) |
| ADM-012 | DR-006 (facility recovery) | All domains (physical infrastructure) |

---

## 11.5 Writing Schedule

Assumes one primary author with review support. Each article estimated at 3-7 working days for draft, 2-3 days for review, 1-2 days for revision.

| Phase | Articles | Duration | Notes |
|---|---|---|---|
| Phase 1: Foundations | ADM-001, ADM-007 | Weeks 1-3 | Must be complete before Phase 2. |
| Phase 2: Core Infrastructure | ADM-004, ADM-002, ADM-006 | Weeks 4-8 | Can partially overlap. ADM-004 should start first. |
| Phase 3: Operational Processes | ADM-005, ADM-008, ADM-003 | Weeks 9-13 | ADM-005 and ADM-008 can overlap. ADM-003 starts after ADM-002 is reviewed. |
| Phase 4: Codification | ADM-009, ADM-010 | Weeks 14-17 | Forms and conflict resolution codify earlier work. |
| Phase 5: Resilience & Capstone | ADM-011, ADM-012, ADM-013 | Weeks 18-23 | Continuity, facilities, and self-audit close out the domain. |

**Total estimated duration:** 23 weeks (approximately 6 months).

---

## 11.6 Review Process

### 11.6.1 Review Stages

1. **Self-Review** (Author, 1 day): Author re-reads after 48-hour cooling period. Checks against article template compliance.
2. **Peer Review** (1-2 reviewers, 3 days): At least one reviewer from within the Administration domain, at least one from a dependent domain.
3. **Cross-Domain Review** (domain leads of dependent domains, 5 days): Ensures consistency with articles in Domains 12-15 and any upstream domains.
4. **Adversarial Review** (1 reviewer, 2 days): A reviewer explicitly tries to find ambiguity, missing failure modes, and unstated assumptions.
5. **Final Acceptance** (Author + 1 senior reviewer): Sign-off. Article enters the canonical documentation set.

### 11.6.2 Review Criteria

- Template compliance (all 11 sections present and substantive)
- Internal consistency (no contradictions within the article)
- External consistency (no contradictions with other articles)
- Clarity for future maintainers (could someone in 50 years understand this without the author?)
- Failure mode coverage (are the realistic failure modes addressed?)
- Testability (can the procedures be tested or drilled?)

### 11.6.3 Dispute Resolution

If reviewers and author disagree, the disagreement is recorded in the Commentary Section of the article. If the disagreement is on a factual or safety matter, it escalates to the cross-domain review panel. Unresolved disputes are flagged and revisited at the next scheduled review cycle.

---

## 11.7 Maintenance Plan

### 11.7.1 Review Cadence

| Review Type | Frequency | Trigger |
|---|---|---|
| Scheduled Full Review | Every 3 years | Calendar |
| Incremental Review | Annually | Calendar |
| Triggered Review | As needed | Personnel change, procurement failure, audit finding, DR event, technology migration |
| Sunset Review | When article subject becomes obsolete | Evolution domain (EVL) notification |

### 11.7.2 Version Control

- Every change produces a new version number (major.minor.patch).
- Major: structural change to process or policy.
- Minor: clarification, addition of a failure mode, updated reference.
- Patch: typo, formatting, commentary addition.
- All previous versions are retained. No version is ever deleted.

### 11.7.3 Succession

- At least two people must be familiar with each article at all times.
- If only one person understands an article, that is an administrative emergency and triggers immediate cross-training.
- The Administration domain itself must have a documented succession plan (ADM-011).

### 11.7.4 Decay Detection

Signs that an article is decaying:
- Referenced procedures no longer match actual practice.
- Forms referenced in ADM-009 have been informally modified without updating the article.
- Personnel named in the article are no longer with the institution.
- Technology referenced has been replaced.
- No Commentary Section entries in more than 5 years (suggests nobody is reading it).

---
---

# DOMAIN 12: DISASTER RECOVERY

**Agent Role:** Resilience Engineer

---

## 12.1 Domain Map

### 12.1.1 Core Purpose

Define every credible failure scenario, the procedures to recover from each, the verification of backup integrity, continuity of operations during degraded states, and the complete rebuild of the institution from catastrophic loss. This domain is the institution's immune system. It must be the most rigorously tested, most frequently reviewed, and most pessimistic domain in the entire documentation set.

### 12.1.2 Scope

**In scope:**
- Failure taxonomy (hardware, software, data, personnel, facilities, supply chain, knowledge)
- Backup strategy, verification, rotation, and geographic distribution
- Recovery procedures for each failure class
- Continuity of operations during degraded modes
- Full rebuild procedures (from bare metal, from backup, from partial knowledge)
- DR testing and drill frameworks
- Recovery time objectives (RTO) and recovery point objectives (RPO)
- Cascade failure analysis
- Communication during disasters (internal coordination without external networks)
- Post-incident review and learning

**Out of scope:**
- Routine maintenance (Domain 11, Domain 5 assumed)
- Ethical decisions during crises (Domain 15 — but cross-referenced)
- Long-term evolution and planned migrations (Domain 13)
- Research into new resilience techniques (Domain 14)

### 12.1.3 Boundaries

Disaster Recovery is reactive and preparatory. It does not design systems (that is the role of technical domains). It does not decide resource allocation (that is Administration). It does not set ethical boundaries for triage decisions (that is Ethics). DR defines what to do when things break, and it verifies that the preparations are adequate.

The boundary between DR and Evolution (Domain 13) is: DR handles unplanned failures; Evolution handles planned transitions. When a planned transition goes wrong, it becomes a DR event.

### 12.1.4 Key Relationships

| Related Domain | Nature of Relationship |
|---|---|
| Domain 11 (Administration) | Administration maintains the inventories, schedules drills, and manages the personnel that DR depends on. DR triggers administrative emergency procedures. |
| Domain 13 (Evolution) | Evolution plans transitions; DR catches transition failures. DR informs Evolution about fragility points that should be addressed in future designs. |
| Domain 14 (Research) | Research may investigate new resilience techniques. DR provides real-world failure data to Research. |
| Domain 15 (Ethics) | Ethics defines triage priorities, acceptable losses, and red lines during crises. DR must operate within ethical boundaries even under extreme pressure. |
| Domains 1-10 (assumed) | DR provides recovery procedures for all technical systems defined in upstream domains. |

---

## 12.2 Article List

| Article ID | Title | Summary |
|---|---|---|
| DR-001 | Disaster Recovery Philosophy and Principles | Foundational assumptions: what we protect, in what order, and why. Defines recovery priorities, acceptable loss thresholds, and the DR mindset. |
| DR-002 | Failure Taxonomy and Threat Model | Comprehensive classification of failure types: hardware, software, data corruption, personnel loss, facility loss, supply chain collapse, knowledge loss. Probability and impact assessments. |
| DR-003 | Backup Strategy and Data Protection | What is backed up, how often, in what format, to what media, in how many copies, at what geographic separation. Includes backup-of-backups and the backup verification schedule. |
| DR-004 | Backup Verification and Integrity Testing | Procedures for proving that backups are actually restorable. Includes bit-rot detection, media degradation monitoring, and test-restore drills. |
| DR-005 | Hardware Failure Recovery | Procedures for recovering from failure of individual components: storage, compute, networking, power, cooling. Includes spare parts inventory requirements. |
| DR-006 | Facility Loss and Relocation | What happens if the physical site is destroyed or rendered uninhabitable. Relocation procedures, off-site asset caches, and minimum viable facility requirements. |
| DR-007 | Data Corruption and Integrity Failure Recovery | Procedures for detecting and recovering from silent data corruption, filesystem damage, database inconsistency, and cryptographic key loss. |
| DR-008 | Personnel Loss and Knowledge Recovery | What happens when key personnel are lost (temporarily or permanently). Knowledge recovery from documentation, cross-training verification, and emergency skill acquisition. |
| DR-009 | Continuity of Operations in Degraded Mode | How the institution continues to function when operating at reduced capacity. Defines degraded-mode tiers, minimum viable operations, and graceful degradation sequences. |
| DR-010 | Full Institutional Rebuild from Catastrophic Loss | The worst-case procedure: rebuilding the entire institution from the most minimal surviving assets. Defines the "rebuild kit" — the minimum set of artifacts needed. |
| DR-011 | DR Drill and Testing Framework | How DR procedures are tested without causing actual damage. Drill types, frequencies, scoring criteria, and lessons-learned integration. |
| DR-012 | Cascade Failure Analysis and Prevention | How to identify and prevent failures that propagate across systems. Isolation boundaries, circuit breakers, and blast radius containment. |
| DR-013 | Post-Incident Review and Learning | How incidents are analyzed after resolution. Root cause analysis methodology, blameless review process, and integration of lessons into documentation. |
| DR-014 | Communication and Coordination During Disasters | How people communicate and coordinate during a disaster when normal communication channels may be unavailable. Includes command structure, decision authority, and status reporting. |

---

## 12.3 Priority Order

| Priority | Article ID | Rationale |
|---|---|---|
| 1 | DR-001 | Foundational philosophy. Every other DR article references this. |
| 2 | DR-002 | Threat model must exist before recovery procedures can be designed. |
| 3 | DR-003 | Backup strategy is the single most critical protective measure. |
| 4 | DR-004 | Backups are worthless if not verified. This must follow immediately after DR-003. |
| 5 | DR-009 | Degraded-mode operations are the most common DR scenario. Most failures are partial, not total. |
| 6 | DR-005 | Hardware failure is the most frequent failure type. |
| 7 | DR-007 | Data corruption is the most insidious failure type (can go undetected). |
| 8 | DR-008 | Personnel loss is the hardest to recover from in a small institution. |
| 9 | DR-014 | Communication must be defined before complex multi-person recovery scenarios. |
| 10 | DR-012 | Cascade analysis depends on understanding individual failure modes (DR-005, DR-007). |
| 11 | DR-006 | Facility loss is low probability but requires extensive planning. |
| 12 | DR-011 | Drill framework requires all procedures to exist before it can test them. |
| 13 | DR-010 | Full rebuild is the capstone — the ultimate DR scenario. Requires all other procedures as inputs. |
| 14 | DR-013 | Post-incident review requires incident experience or drill results. Can be written in parallel with DR-011. |

---

## 12.4 Dependencies

### 12.4.1 Internal Dependencies (within Domain 12)

```
DR-001 --> DR-002 --> DR-003 --> DR-004
                  \-> DR-009
                  \-> DR-005 --> DR-012
                  \-> DR-007 --> DR-012
                  \-> DR-008
                  \-> DR-006
DR-005 + DR-007 + DR-008 + DR-009 --> DR-014
DR-all procedures --> DR-011 (drill framework tests everything)
DR-all procedures --> DR-010 (full rebuild references everything)
DR-011 + DR-010 --> DR-013 (post-incident review)
```

### 12.4.2 Cross-Domain Dependencies

| Article | Depends On (External) | Depended On By (External) |
|---|---|---|
| DR-001 | ETH-001 (ethical principles for triage), ADM-001 (administrative principles) | All domains (DR philosophy) |
| DR-002 | EVL-001 (evolution risks as threat inputs) | EVL-003 (fragility data for migration planning) |
| DR-003 | ADM-004 (asset inventory for backup media tracking) | ADM-002 (backup costs in budget) |
| DR-005 | ADM-003 (procurement of spare parts), ADM-004 (spare parts inventory) | EVL-005 (hardware sunset triggers) |
| DR-006 | ADM-012 (facilities management) | ADM-012 (facility recovery informs facility management) |
| DR-008 | ADM-006 (personnel records), ADM-011 (admin continuity) | ADM-011 (mutual dependency — personnel continuity) |
| DR-009 | ETH-006 (ethical triage during degraded operations) | ADM-005 (degraded-mode scheduling) |
| DR-010 | All domains (rebuild requires all documentation) | EVL-010 (institutional architecture for rebuild design) |
| DR-011 | ADM-005 (drill scheduling) | RES-008 (drill data as research input) |
| DR-013 | RES-002 (analysis methodology) | RES-008 (incident data as research input), ETH-010 (ethical review of incidents) |

---

## 12.5 Writing Schedule

| Phase | Articles | Duration | Notes |
|---|---|---|---|
| Phase 1: Foundations | DR-001, DR-002 | Weeks 1-4 | DR-002 is the largest article in the domain — threat modeling is extensive. |
| Phase 2: Backup System | DR-003, DR-004 | Weeks 5-9 | These two articles are tightly coupled and should be written by the same author. |
| Phase 3: Core Recovery | DR-009, DR-005, DR-007 | Weeks 10-16 | Can partially overlap. DR-009 should start first (most common scenario). |
| Phase 4: Human & Communication | DR-008, DR-014 | Weeks 17-21 | People-focused recovery. |
| Phase 5: Advanced Analysis | DR-012, DR-006 | Weeks 22-26 | Cascade analysis and facility loss are complex but lower priority. |
| Phase 6: Capstone | DR-011, DR-010, DR-013 | Weeks 27-34 | Drill framework, full rebuild, and post-incident review. DR-010 is the largest single article — it references everything. |

**Total estimated duration:** 34 weeks (approximately 8.5 months).

**Note:** DR is intentionally given more time than other domains. It is the most safety-critical documentation set and must be exhaustively detailed.

---

## 12.6 Review Process

### 12.6.1 Review Stages

1. **Self-Review** (Author, 1 day): Template compliance, completeness check.
2. **Tabletop Walkthrough** (Author + 2 reviewers, 1 day): Walk through each procedure step by step on paper. Challenge every assumption.
3. **Cross-Domain Review** (Administration, Ethics, Evolution leads, 5 days): Ensure DR articles are consistent with administrative realities, ethical boundaries, and evolution plans.
4. **Adversarial Red Team Review** (1-2 reviewers, 3 days): Reviewers actively try to find scenarios where the procedures would fail. They invent new failure modes and check if the article handles them.
5. **Practical Test** (where possible, variable duration): For articles like DR-003, DR-004, DR-005 — actually test the procedure on non-production systems before accepting the article.
6. **Final Acceptance** (Author + Resilience Engineer lead + 1 senior reviewer): Sign-off.

### 12.6.2 Special Review Requirements

- DR-010 (Full Rebuild) must be reviewed by someone who was NOT involved in building the original system. This tests whether the rebuild instructions are truly self-contained.
- DR-003 and DR-004 (Backup) must include a verified test-restore as part of the acceptance criteria.
- All DR articles must include explicit "you will know recovery has succeeded when..." criteria.

### 12.6.3 Dispute Resolution

Disputes in DR articles default to the more conservative/pessimistic position. If there is disagreement about whether a failure mode is credible, it is included. Exclusion requires unanimous agreement that the scenario is not credible. All disputes are recorded in the Commentary Section.

---

## 12.7 Maintenance Plan

### 12.7.1 Review Cadence

| Review Type | Frequency | Trigger |
|---|---|---|
| Scheduled Full Review | Every 2 years | Calendar. DR has a shorter review cycle than other domains. |
| Drill-Triggered Review | After every DR drill | Drill results may reveal procedure gaps. |
| Incident-Triggered Review | After every real incident | Mandatory, even if the incident was minor. |
| Technology-Triggered Review | When any system component changes | Hardware swap, software update, configuration change. |
| Annual Backup Verification | Annually | Calendar. Includes test-restore of at least one full backup. |

### 12.7.2 Version Control

Same as Domain 11 (major.minor.patch, all versions retained).

### 12.7.3 Succession

- DR documentation must be understandable by someone with no prior institutional knowledge. This is the ultimate succession requirement — DR-010 (Full Rebuild) is, in essence, the succession plan for the entire institution.
- At least three people must be familiar with DR-003 (Backup) and DR-004 (Verification) at all times.
- Annual "cold read" test: someone who has never read a specific DR article is asked to execute its procedures. Their success or failure drives article revision.

### 12.7.4 Decay Detection

- DR procedures not drilled in more than 2 years are considered potentially decayed.
- Backup media not verified in more than 1 year is considered untrustworthy.
- Spare parts inventory not audited in more than 1 year is considered unreliable.
- DR articles referencing personnel who have departed are immediately flagged for revision.

---
---

# DOMAIN 13: EVOLUTION & ADAPTATION

**Agent Role:** Evolution Strategist

---

## 13.1 Domain Map

### 13.1.1 Core Purpose

Define how the institution changes deliberately over time. Technology migrations, version transitions, paradigm shifts, sunset procedures, upgrade paths, and the philosophy of managed change. While DR (Domain 12) handles unplanned change, Evolution handles planned change. The goal is to ensure the institution can adapt to a fundamentally different technological and cultural landscape over a 50+ year horizon without losing continuity of purpose or integrity of data.

### 13.1.2 Scope

**In scope:**
- Philosophy of institutional evolution (when to change, when to resist change)
- Technology migration planning and execution
- Version transition procedures (hardware, software, data formats, documentation formats)
- Paradigm shift management (when foundational assumptions change)
- Sunset and decommissioning procedures
- Upgrade path design and validation
- Compatibility management (forward and backward)
- Legacy system maintenance during transition periods
- Change impact assessment methodology
- Institutional memory preservation during transitions
- Technology horizon scanning (evaluating emerging technologies without internet access)

**Out of scope:**
- Emergency/unplanned changes (Domain 12)
- Day-to-day operational scheduling (Domain 11)
- Ethical evaluation of new technologies (Domain 15, but cross-referenced)
- Research into new technologies (Domain 14, but tightly coupled)

### 13.1.3 Boundaries

Evolution is about deliberate, planned transformation. It is not about reacting to crises (that is DR). It is not about maintaining the status quo (that is Administration). It is not about researching whether a new approach is sound (that is Research). Evolution takes inputs from Research ("this new approach is viable"), Administration ("we have resources for this transition"), and Ethics ("this change is permissible") and produces migration plans, sunset schedules, and upgrade procedures.

The critical boundary: Evolution must never be used to justify change for its own sake. Every proposed evolution must demonstrate that it serves the institution's long-term survival and mission.

### 13.1.4 Key Relationships

| Related Domain | Nature of Relationship |
|---|---|
| Domain 11 (Administration) | Administration resources and schedules evolution activities. Evolution generates new administrative requirements. |
| Domain 12 (Disaster Recovery) | DR handles evolution failures. Evolution reduces DR risk by replacing fragile systems. DR informs Evolution about fragility points. |
| Domain 14 (Research) | Research validates new technologies. Evolution consumes research outputs to plan migrations. |
| Domain 15 (Ethics) | Ethics evaluates whether proposed changes are permissible. Evolution must obtain ethical clearance for significant changes. |
| Domains 1-10 (assumed) | Evolution modifies, replaces, or retires systems defined in technical domains. |

---

## 13.2 Article List

| Article ID | Title | Summary |
|---|---|---|
| EVL-001 | Philosophy of Institutional Evolution | When and why to change. The tension between stability and adaptation. Criteria for initiating change. The precautionary principle applied to institutional transformation. |
| EVL-002 | Version Management and Compatibility Policy | How versions are numbered, tracked, and managed across all systems. Forward and backward compatibility requirements. Version lifecycle (active, deprecated, sunset, archived). |
| EVL-003 | Change Impact Assessment Methodology | How to evaluate the full impact of a proposed change before committing to it. Risk assessment, dependency mapping, rollback feasibility, and blast radius estimation. |
| EVL-004 | Technology Migration Planning | How to plan a migration from one technology to another. Phased rollout, parallel running, cutover procedures, and success criteria. |
| EVL-005 | Hardware Lifecycle and Replacement Strategy | How to plan for hardware aging, obsolescence, and replacement. Includes strategies for stockpiling, interoperability, and graceful degradation as hardware ages. |
| EVL-006 | Data Format Migration and Preservation | How to migrate data from one format to another without loss. Includes format obsolescence risk assessment, canonical format selection, and migration verification. |
| EVL-007 | Documentation Evolution and Format Migration | How the documentation itself evolves. Format changes, restructuring, terminology updates, and maintaining referential integrity across the documentation set. |
| EVL-008 | Sunset and Decommissioning Procedures | How to retire a system, technology, or process. Data extraction, dependency severing, archive creation, and post-sunset verification. |
| EVL-009 | Paradigm Shift Management | What happens when a foundational assumption changes (e.g., a core technology becomes untenable, a fundamental approach is invalidated). How to manage deep structural change without losing institutional coherence. |
| EVL-010 | Institutional Architecture and Design Principles for Longevity | Meta-article on designing all systems for evolvability. Modularity, loose coupling, interface stability, and the "replaceable parts" philosophy. |
| EVL-011 | Technology Horizon Scanning Without Internet Access | How to evaluate emerging technologies when the institution is air-gapped and off-grid. Information sourcing, evaluation frameworks, and controlled experimentation. |
| EVL-012 | Legacy System Coexistence and Bridge Strategies | How old and new systems coexist during transition periods. Bridge layers, adapters, dual-write strategies, and coexistence duration limits. |
| EVL-013 | Post-Migration Verification and Stabilization | How to verify that a migration succeeded and the new system is stable. Includes burn-in periods, regression testing, and rollback decision points. |

---

## 13.3 Priority Order

| Priority | Article ID | Rationale |
|---|---|---|
| 1 | EVL-001 | Foundational philosophy. Defines when change is justified. |
| 2 | EVL-010 | Design principles for longevity shape all future evolution decisions. |
| 3 | EVL-002 | Version management is infrastructure for all other evolution activities. |
| 4 | EVL-003 | Change impact assessment must exist before any migration is planned. |
| 5 | EVL-005 | Hardware lifecycle is the most immediate and concrete evolution concern for an off-grid institution. |
| 6 | EVL-006 | Data format preservation is the second most urgent concern — data outlives hardware. |
| 7 | EVL-004 | Migration planning depends on impact assessment (EVL-003) and hardware/data context (EVL-005, EVL-006). |
| 8 | EVL-008 | Sunset procedures are needed as soon as migrations begin. |
| 9 | EVL-012 | Legacy coexistence is needed during any transition. |
| 10 | EVL-013 | Post-migration verification closes the loop on EVL-004. |
| 11 | EVL-007 | Documentation evolution is important but the documentation set must first exist before it can evolve. |
| 12 | EVL-009 | Paradigm shift management is high-impact but low-frequency. |
| 13 | EVL-011 | Horizon scanning is ongoing but not urgent for initial framework. |

---

## 13.4 Dependencies

### 13.4.1 Internal Dependencies (within Domain 13)

```
EVL-001 --> EVL-010 --> EVL-002
EVL-002 --> EVL-003 --> EVL-004 --> EVL-013
                    \-> EVL-005
                    \-> EVL-006
EVL-004 --> EVL-008 (sunset is part of migration)
EVL-004 --> EVL-012 (coexistence during migration)
EVL-001 --> EVL-009 (paradigm shifts reference philosophy)
EVL-010 --> EVL-007 (documentation evolution follows design principles)
EVL-001 --> EVL-011 (horizon scanning follows evolution philosophy)
```

### 13.4.2 Cross-Domain Dependencies

| Article | Depends On (External) | Depended On By (External) |
|---|---|---|
| EVL-001 | ETH-001 (ethical principles constrain evolution), DR-002 (threat model informs evolution priorities) | DR-002 (evolution risks feed threat model), ADM-001 (evolution philosophy informs administration) |
| EVL-002 | ADM-004 (asset registry tracks versions) | DR-005 (version info needed for hardware recovery), ADM-004 (version tracking in inventory) |
| EVL-003 | DR-012 (cascade failure analysis informs impact assessment) | ADM-002 (change costs feed budget), DR-002 (change risks feed threat model) |
| EVL-004 | ADM-005 (scheduling), ADM-002 (budget) | DR-009 (degraded mode during migration) |
| EVL-005 | ADM-003 (procurement), ADM-004 (inventory) | DR-005 (hardware recovery), ADM-002 (hardware budget) |
| EVL-006 | RES-006 (format research), RES-005 (data integrity research) | DR-007 (data recovery format awareness) |
| EVL-008 | ADM-004 (asset registry for decommissioning) | ADM-004 (asset removal from registry) |
| EVL-009 | RES-001 (research methodology for paradigm evaluation) | ETH-009 (ethical implications of paradigm shifts) |
| EVL-011 | RES-001 (research methodology) | RES-003 (research priorities from horizon scanning) |

---

## 13.5 Writing Schedule

| Phase | Articles | Duration | Notes |
|---|---|---|---|
| Phase 1: Philosophy & Architecture | EVL-001, EVL-010 | Weeks 1-5 | These are deep, philosophical articles. Allow extra time. |
| Phase 2: Infrastructure | EVL-002, EVL-003 | Weeks 6-10 | Version management and impact assessment are foundational infrastructure. |
| Phase 3: Concrete Concerns | EVL-005, EVL-006 | Weeks 11-16 | Hardware and data format migration are the most tangible, immediate concerns. |
| Phase 4: Migration Operations | EVL-004, EVL-008, EVL-012 | Weeks 17-23 | The operational heart of evolution: how to actually do migrations. |
| Phase 5: Verification & Documentation | EVL-013, EVL-007 | Weeks 24-28 | Post-migration verification and documentation evolution. |
| Phase 6: Long-Horizon | EVL-009, EVL-011 | Weeks 29-34 | Paradigm shifts and horizon scanning are long-term concerns. |

**Total estimated duration:** 34 weeks (approximately 8.5 months).

---

## 13.6 Review Process

### 13.6.1 Review Stages

1. **Self-Review** (Author, 1 day): Template compliance, internal consistency.
2. **Longevity Review** (1 reviewer, 3 days): Specifically asks: "Will this article still make sense in 50 years? Does it depend on assumptions that may not hold?" This is unique to the Evolution domain.
3. **Peer Review** (1-2 reviewers, 3 days): Technical accuracy and procedural completeness.
4. **Cross-Domain Review** (DR, Administration, Ethics leads, 5 days): Ensures evolution articles are consistent with recovery, operational, and ethical constraints.
5. **Historical Perspective Review** (1 reviewer, 2 days): Reviewer examines the article through the lens of past technology transitions (e.g., "Would this approach have worked during the transition from X to Y?"). This grounds the article in historical precedent.
6. **Final Acceptance** (Author + Evolution Strategist lead + 1 senior reviewer): Sign-off.

### 13.6.2 Special Review Requirements

- EVL-009 (Paradigm Shifts) must be reviewed by someone outside the technical domains to check for blind spots.
- EVL-006 (Data Format Migration) must include a proof-of-concept migration of a small dataset as part of acceptance.
- EVL-010 (Design Principles) must be reviewed against every other domain's articles for consistency.

### 13.6.3 Dispute Resolution

Disputes in Evolution articles default to the more conservative position (slower change, more verification). The bias is toward stability unless the case for change is overwhelming. All disputes are recorded in the Commentary Section with full reasoning from all parties.

---

## 13.7 Maintenance Plan

### 13.7.1 Review Cadence

| Review Type | Frequency | Trigger |
|---|---|---|
| Scheduled Full Review | Every 5 years | Calendar. Evolution articles change slowly — they describe long-term processes. |
| Migration-Triggered Review | After every completed migration | Lessons learned must be integrated. |
| Technology Assessment | Every 2 years | Horizon scanning update. |
| Paradigm Review | Every 10 years | Fundamental reassessment of assumptions. |

### 13.7.2 Version Control

Same as Domain 11 (major.minor.patch, all versions retained).

### 13.7.3 Succession

- Evolution planning is the most "strategic" domain. It requires institutional context that takes years to develop.
- At least two people must be deeply familiar with EVL-001, EVL-009, and EVL-010 at all times.
- New Evolution Strategists must serve a minimum 2-year apprenticeship before being the sole author of evolution articles.

### 13.7.4 Decay Detection

- Evolution articles that have never been used to plan an actual migration are suspect — they may be theoretical rather than practical.
- If a technology migration occurs without referencing the Evolution domain, the domain has failed its purpose. This is a critical decay signal.
- If EVL-011 (Horizon Scanning) has not produced any outputs in 3 years, the institution may be stagnating.

---
---

# DOMAIN 14: RESEARCH & THEORY

**Agent Role:** Research Director

---

## 14.1 Domain Map

### 14.1.1 Core Purpose

Define how the institution creates new knowledge, tests hypotheses, evaluates technologies, maintains intellectual rigor, and distinguishes verified knowledge from opinion, assumption, and folklore. In a long-lived institution, the temptation to rely on "the way we've always done it" is enormous. Research & Theory provides the antidote: systematic methods for questioning assumptions, testing beliefs, and generating new understanding.

### 14.1.2 Scope

**In scope:**
- Research methodology (experimental design, hypothesis testing, controlled observation)
- Intellectual rigor standards (evidence requirements, burden of proof, peer review)
- Knowledge creation and validation workflows
- Experimental frameworks for testing new technologies, processes, and approaches
- Documentation of research outputs (papers, findings, negative results)
- Research ethics (separate from but informed by Domain 15)
- Knowledge classification (verified, provisional, speculative, deprecated)
- Research resource management (time, equipment, personnel)
- Replication and reproducibility requirements
- Integration of research findings into operational practice
- Preservation of negative results and failed experiments

**Out of scope:**
- Operational procedures (Domains 1-11)
- Ethical policy-making (Domain 15)
- Migration execution (Domain 13)
- Disaster response (Domain 12)

### 14.1.3 Boundaries

Research & Theory provides validated knowledge to other domains. It does not make operational decisions. It does not set ethical policy. It does not execute migrations. It investigates, validates, and reports. Other domains consume research outputs and translate them into action.

The boundary between Research and Evolution is: Research asks "Is this viable?" Evolution asks "How do we transition to it?" Research must complete before Evolution begins.

The boundary between Research and Ethics is: Research defines methodology; Ethics defines what research is permissible.

### 14.1.4 Key Relationships

| Related Domain | Nature of Relationship |
|---|---|
| Domain 11 (Administration) | Administration allocates resources to research. Research outputs may change administrative procedures. |
| Domain 12 (Disaster Recovery) | Research analyzes incident data. Research may develop new resilience techniques. DR provides real-world failure data. |
| Domain 13 (Evolution) | Research validates new technologies. Evolution consumes research findings. Research informs evolution priorities. |
| Domain 15 (Ethics) | Ethics constrains research methods. Research may investigate ethical questions empirically. |
| Domains 1-10 (assumed) | Research may investigate improvements to any technical system. |

---

## 14.2 Article List

| Article ID | Title | Summary |
|---|---|---|
| RES-001 | Research Philosophy and Intellectual Standards | Why the institution conducts research. Epistemological foundations. The relationship between knowledge and action. Standards for what counts as "knowing" something. |
| RES-002 | Research Methodology Framework | How to design and conduct research. Hypothesis formation, experimental design, controls, variables, data collection, analysis, and conclusions. Adapted for a small, resource-constrained, air-gapped institution. |
| RES-003 | Research Prioritization and Resource Allocation | How research projects are proposed, evaluated, prioritized, and resourced. Balancing curiosity-driven and mission-critical research. |
| RES-004 | Experimental Framework for Technology Evaluation | Specific methodology for evaluating new technologies, tools, and approaches. Includes isolation requirements, benchmark design, and acceptance criteria. |
| RES-005 | Data Integrity in Research | How research data is collected, stored, validated, and protected from corruption. Chain of custody for experimental data. |
| RES-006 | Knowledge Classification and Confidence Levels | Taxonomy of knowledge: verified (replicated, high confidence), provisional (single study, moderate confidence), speculative (theoretical, low confidence), deprecated (superseded or disproven). |
| RES-007 | Peer Review and Validation Process for Research Outputs | How research findings are reviewed, challenged, and validated before being accepted into the institutional knowledge base. |
| RES-008 | Negative Results and Failed Experiments Registry | Why and how to document experiments that failed, hypotheses that were disproven, and approaches that did not work. Preventing future researchers from repeating known failures. |
| RES-009 | Integration of Research into Operational Practice | How validated research findings are translated into changes in operational procedures, system designs, or institutional policies. The gap between "we know this" and "we do this." |
| RES-010 | Research Ethics and Permissible Experimentation | Ethical boundaries for research activities. What experiments are permissible, what safeguards are required, and what oversight mechanisms exist. (Cross-references Domain 15 extensively.) |
| RES-011 | Long-Term Research Programs and Multi-Generational Studies | How to design and maintain research programs that span decades or generations of researchers. Continuity, handoff, and living research documents. |
| RES-012 | Replication and Reproducibility Standards | Requirements for replicating research results. What constitutes a successful replication. How to handle replication failures. |
| RES-013 | Research Archive and Knowledge Preservation | How research outputs (positive and negative) are archived, indexed, and made accessible to future researchers. Includes format longevity considerations. |

---

## 14.3 Priority Order

| Priority | Article ID | Rationale |
|---|---|---|
| 1 | RES-001 | Foundational philosophy. Defines what research means in this institution. |
| 2 | RES-002 | Methodology is the core of the domain. Everything else depends on it. |
| 3 | RES-006 | Knowledge classification must exist before research outputs can be properly categorized. |
| 4 | RES-005 | Data integrity is infrastructure for all research activities. |
| 5 | RES-007 | Peer review validates research outputs. Must exist before any findings are accepted. |
| 6 | RES-010 | Research ethics must be established early to prevent permissibility disputes. |
| 7 | RES-003 | Prioritization requires understanding of methodology (RES-002) and ethics (RES-010). |
| 8 | RES-004 | Technology evaluation is the most common research type for this institution. |
| 9 | RES-012 | Replication standards depend on methodology (RES-002) and data integrity (RES-005). |
| 10 | RES-008 | Negative results registry can begin as soon as any research is conducted. |
| 11 | RES-009 | Integration into practice requires research outputs to exist. |
| 12 | RES-013 | Archive design depends on understanding what research outputs look like. |
| 13 | RES-011 | Long-term research programs are important but not urgent for initial framework. |

---

## 14.4 Dependencies

### 14.4.1 Internal Dependencies (within Domain 14)

```
RES-001 --> RES-002 --> RES-005 --> RES-012
                    \-> RES-006 --> RES-007
                    \-> RES-010
RES-002 + RES-010 --> RES-003 --> RES-004
RES-007 --> RES-009 (integration requires validated findings)
RES-002 + RES-005 --> RES-008 (negative results follow same rigor)
RES-006 + RES-008 --> RES-013 (archive stores all classified knowledge)
RES-001 + RES-002 --> RES-011 (long-term programs require deep methodology)
```

### 14.4.2 Cross-Domain Dependencies

| Article | Depends On (External) | Depended On By (External) |
|---|---|---|
| RES-001 | ETH-001 (ethical foundations inform research philosophy) | EVL-001 (research philosophy informs evolution philosophy) |
| RES-002 | None (foundational) | DR-013 (post-incident analysis uses research methodology), EVL-003 (impact assessment uses analytical methods) |
| RES-003 | ADM-002 (budget for research), ADM-005 (scheduling) | ADM-002 (research costs in budget) |
| RES-004 | EVL-011 (horizon scanning identifies technologies to evaluate) | EVL-004 (migration planning uses evaluation results), EVL-005 (hardware decisions use evaluation) |
| RES-005 | ADM-004 (data storage inventory) | DR-007 (data integrity methods) |
| RES-006 | None (foundational) | All domains (knowledge confidence levels) |
| RES-008 | None (foundational) | All domains (avoiding known failures) |
| RES-009 | ADM-001 (administrative processes that may change) | EVL-004 (research-driven migrations), ADM processes (operational changes) |
| RES-010 | ETH-001, ETH-002, ETH-005 (ethical framework) | ETH-010 (ethical review of research activities) |
| RES-013 | EVL-006 (data format preservation), ADM-007 (archival standards) | EVL-006 (archive format decisions), DR-010 (rebuild includes research archive) |

---

## 14.5 Writing Schedule

| Phase | Articles | Duration | Notes |
|---|---|---|---|
| Phase 1: Foundations | RES-001, RES-002 | Weeks 1-5 | RES-002 (Methodology) is the most complex article. Allow extra time. |
| Phase 2: Infrastructure | RES-006, RES-005, RES-010 | Weeks 6-12 | Classification, data integrity, and ethics form the infrastructure for all research. |
| Phase 3: Validation | RES-007, RES-012 | Weeks 13-17 | Peer review and replication standards. |
| Phase 4: Operations | RES-003, RES-004 | Weeks 18-23 | Prioritization and technology evaluation. |
| Phase 5: Integration & Preservation | RES-008, RES-009, RES-013 | Weeks 24-30 | Negative results, practice integration, and archival. |
| Phase 6: Long-Horizon | RES-011 | Weeks 31-34 | Multi-generational research programs. |

**Total estimated duration:** 34 weeks (approximately 8.5 months).

---

## 14.6 Review Process

### 14.6.1 Review Stages

1. **Self-Review** (Author, 1 day): Template compliance, internal consistency.
2. **Methodological Review** (1-2 reviewers with research experience, 5 days): Is the methodology sound? Are the standards achievable? Are there logical gaps?
3. **Practicality Review** (1 reviewer from Administration, 3 days): Are the resource requirements realistic? Can this actually be done in a small, off-grid institution?
4. **Cross-Domain Review** (Evolution, Ethics leads, 5 days): Ensures research articles are consistent with evolution plans and ethical constraints.
5. **Epistemological Review** (1 reviewer, 2 days): Challenges the foundational assumptions. Asks: "What if this way of knowing is wrong? What alternatives exist?"
6. **Final Acceptance** (Author + Research Director lead + 1 senior reviewer): Sign-off.

### 14.6.2 Special Review Requirements

- RES-001 (Philosophy) must be reviewed by someone from outside the Research domain to avoid insularity.
- RES-002 (Methodology) must include at least one worked example demonstrating the methodology.
- RES-010 (Research Ethics) must be jointly reviewed with the Ethics domain lead (Domain 15).

### 14.6.3 Dispute Resolution

Disputes in Research articles are resolved by evidence and argument, not by authority. If a dispute cannot be resolved, both positions are documented in the Commentary Section with supporting arguments. The article is published with the dispute noted. Resolution is deferred to future evidence.

---

## 14.7 Maintenance Plan

### 14.7.1 Review Cadence

| Review Type | Frequency | Trigger |
|---|---|---|
| Scheduled Full Review | Every 5 years | Calendar. Research methodology evolves slowly. |
| Research-Output-Triggered Review | When a significant finding changes understanding | Any finding classified as "verified" that contradicts existing articles. |
| Methodology Review | Every 3 years | Reassess whether the methodology is still appropriate. |
| Archive Review | Every 2 years | Ensure the research archive is accessible and indexed. |

### 14.7.2 Version Control

Same as Domain 11 (major.minor.patch, all versions retained).

### 14.7.3 Succession

- Research methodology is the hardest competency to transfer. It requires both theoretical understanding and practical experience.
- At least two people must be trained in RES-002 (Methodology) at all times.
- New Research Directors must have conducted at least three complete research projects under supervision before assuming the role.

### 14.7.4 Decay Detection

- If no new research has been conducted in 2 years, the Research domain is decaying.
- If the Negative Results Registry (RES-008) has no entries, either no research is being done or negative results are being discarded (both are problems).
- If research findings are being integrated into practice (RES-009) without going through peer review (RES-007), the validation process has been bypassed.
- If RES-006 (Knowledge Classification) has not been updated to reflect new findings, classifications may be stale.

---
---

# DOMAIN 15: ETHICS & SAFEGUARDS

**Agent Role:** Ethics Officer

---

## 15.1 Domain Map

### 15.1.1 Core Purpose

Define the ethical boundaries of the institution, the mechanisms that enforce those boundaries, the oversight structures that detect violations, and the protections for those who report violations. This domain is the conscience of the institution. It exists to prevent the institution from becoming something its founders would not recognize or endorse, while also allowing for the possibility that the founders' ethical understanding was incomplete.

This is the only domain that has authority to halt activities in other domains.

### 15.1.2 Scope

**In scope:**
- Ethical principles and foundational values
- Red lines (actions that are never permissible regardless of circumstances)
- Safeguard mechanisms (structural protections against ethical violations)
- Oversight frameworks (who watches the watchers)
- Whistleblower protections (how individuals can safely report concerns)
- Harm prevention (anticipating and preventing foreseeable harm)
- Ethical review of new technologies, processes, and institutional changes
- Conflict of interest management
- Power concentration prevention
- Ethical decision-making frameworks for novel situations
- Rights and protections for institutional members
- Ethical evolution (how ethical standards themselves change over time)
- Transparency and accountability requirements

**Out of scope:**
- Operational procedures (Domains 1-11)
- Technical system design (Domains 1-10 assumed)
- Research methodology (Domain 14, but ethics constrains research)
- Disaster response execution (Domain 12, but ethics constrains triage)

### 15.1.3 Boundaries

Ethics & Safeguards has a unique position: it has veto power over activities in all other domains but does not have directive power (it cannot tell other domains what to do, only what they must not do). It is a check, not a command structure.

The boundary between Ethics and Governance (assumed upstream) is: Governance sets institutional direction; Ethics ensures that direction does not violate fundamental principles.

The boundary between Ethics and Administration is: Administration executes processes; Ethics ensures those processes are fair, transparent, and non-harmful.

The boundary between Ethics and Research is: Research defines methodology; Ethics defines what research is permissible.

### 15.1.4 Key Relationships

| Related Domain | Nature of Relationship |
|---|---|
| Domain 11 (Administration) | Ethics constrains administrative processes (e.g., fair procurement, non-discriminatory personnel management). Administration implements ethical requirements. |
| Domain 12 (Disaster Recovery) | Ethics defines triage priorities and acceptable losses during crises. DR operates within ethical boundaries even under extreme pressure. |
| Domain 13 (Evolution) | Ethics evaluates whether proposed changes are permissible. Evolution must obtain ethical clearance for significant changes. |
| Domain 14 (Research) | Ethics constrains research methods. Ethics defines what experiments are permissible. Research provides empirical input to ethical deliberation. |
| All Domains | Ethics has oversight authority over all domains. All domains must comply with ethical red lines. |

---

## 15.2 Article List

| Article ID | Title | Summary |
|---|---|---|
| ETH-001 | Ethical Foundations and Core Values | The non-negotiable ethical principles of the institution. Why they exist. How they were chosen. How they relate to each other when they conflict. |
| ETH-002 | Red Lines: Absolute Prohibitions | Actions that are never permissible regardless of circumstances, justification, or pressure. The "no matter what" list. Includes rationale for each red line. |
| ETH-003 | Safeguard Mechanisms and Structural Protections | Structural features of the institution designed to prevent ethical violations. Separation of duties, mandatory review points, automatic triggers, and failsafe defaults. |
| ETH-004 | Oversight Framework: Who Watches the Watchers | How ethical oversight is structured. Layers of oversight, independence requirements, rotation, and the prevention of captured oversight. |
| ETH-005 | Ethical Procurement and Resource Use | Ethical constraints on how resources are acquired and used. Fair dealing, sustainability, and the prevention of exploitative practices. |
| ETH-006 | Ethical Triage and Crisis Decision-Making | How ethical principles apply during crises when resources are scarce and decisions are urgent. Pre-committed triage frameworks to reduce ad-hoc decision-making under pressure. |
| ETH-007 | Whistleblower Protections and Dissent Channels | How individuals can safely report ethical concerns, violations, or disagreements. Protection from retaliation. Anonymous reporting mechanisms. |
| ETH-008 | Personnel Ethics and Rights | Ethical obligations toward institutional members. Fair treatment, privacy, autonomy, consent, and the right to refuse participation in activities that violate personal ethics. |
| ETH-009 | Ethical Review of Institutional Changes | Process for ethically evaluating proposed changes to systems, processes, or policies. Triggered by Evolution (Domain 13) and by any domain proposing significant changes. |
| ETH-010 | Power Concentration Prevention | Mechanisms to prevent any individual, role, or faction from accumulating excessive power. Includes term limits, separation of authority, mandatory rotation, and structural checks. |
| ETH-011 | Ethical Evolution: How Ethics Itself Changes | How the institution's ethical standards can change over time. The meta-ethical framework. Requirements for modifying red lines. The tension between moral consistency and moral growth. |
| ETH-012 | Transparency and Accountability Standards | What information must be available to all institutional members. What decisions must be explained. What records cannot be kept secret. |
| ETH-013 | Harm Prevention and Impact Assessment | How to anticipate and prevent foreseeable harm from institutional activities. Includes a harm taxonomy and pre-mortem methodology. |
| ETH-014 | Conflict of Interest Management | How to identify, disclose, and manage conflicts of interest in decision-making, resource allocation, and oversight. |
| ETH-015 | Ethical Incident Response and Remediation | What happens when an ethical violation is discovered. Investigation, remediation, accountability, and institutional learning. |

---

## 15.3 Priority Order

| Priority | Article ID | Rationale |
|---|---|---|
| 1 | ETH-001 | Foundational. Every other Ethics article references these core values. |
| 2 | ETH-002 | Red lines must be established immediately. They are the highest-priority ethical output. |
| 3 | ETH-003 | Safeguard mechanisms enforce the red lines. They must exist before the institution begins operations. |
| 4 | ETH-007 | Whistleblower protections must exist before any ethical violation can occur, so that reporting is safe from day one. |
| 5 | ETH-004 | Oversight framework must be established early to prevent early power consolidation. |
| 6 | ETH-010 | Power concentration prevention is urgent — power consolidates quickly if unchecked. |
| 7 | ETH-008 | Personnel rights must be defined before personnel are subject to institutional processes. |
| 8 | ETH-012 | Transparency standards support oversight (ETH-004) and whistleblowing (ETH-007). |
| 9 | ETH-006 | Crisis ethics must be pre-committed before a crisis occurs. |
| 10 | ETH-005 | Ethical procurement constrains administration (Domain 11). |
| 11 | ETH-014 | Conflict of interest management depends on oversight (ETH-004) and transparency (ETH-012). |
| 12 | ETH-009 | Ethical review of changes depends on having the ethical framework in place. |
| 13 | ETH-013 | Harm prevention is ongoing but depends on the full ethical framework. |
| 14 | ETH-015 | Incident response requires the full ethical framework and oversight structure. |
| 15 | ETH-011 | Ethical evolution is the most sensitive article. It must be written last, with the full weight of the established framework behind it. |

---

## 15.4 Dependencies

### 15.4.1 Internal Dependencies (within Domain 15)

```
ETH-001 --> ETH-002 --> ETH-003
ETH-001 --> ETH-007 (whistleblower protections from core values)
ETH-001 --> ETH-004 --> ETH-010
ETH-001 --> ETH-008
ETH-004 + ETH-007 --> ETH-012 (transparency supports oversight and whistleblowing)
ETH-001 --> ETH-006 (crisis ethics from core values)
ETH-001 --> ETH-005
ETH-004 + ETH-012 --> ETH-014 (conflict of interest requires oversight and transparency)
ETH-002 + ETH-003 + ETH-004 --> ETH-009 (ethical review requires red lines, safeguards, oversight)
ETH-001 + ETH-013 --> ETH-015 (incident response requires harm framework)
All articles --> ETH-011 (ethical evolution references everything)
```

### 15.4.2 Cross-Domain Dependencies

| Article | Depends On (External) | Depended On By (External) |
|---|---|---|
| ETH-001 | Institutional charter (assumed upstream) | All domains (ethical foundation), DR-001 (DR philosophy), RES-001 (research philosophy), EVL-001 (evolution philosophy) |
| ETH-002 | None (foundational, self-contained) | All domains (red lines are universal constraints) |
| ETH-003 | ADM-008 (reporting structures for safeguard enforcement) | ADM-008 (accountability structures), DR-009 (degraded mode ethical safeguards) |
| ETH-005 | ADM-003 (procurement procedures to constrain) | ADM-003 (ethical procurement requirements) |
| ETH-006 | DR-001 (disaster context), DR-009 (degraded mode definitions) | DR-009 (ethical triage during degraded ops), DR-014 (ethical communication during disasters) |
| ETH-007 | ADM-006 (personnel structures for protection mechanisms) | ADM-006 (whistleblower integration in personnel management) |
| ETH-008 | ADM-006 (personnel administration) | ADM-006 (personnel rights in administration) |
| ETH-009 | EVL-003 (change impact assessment) | EVL-004 (ethical clearance for migrations), EVL-009 (ethical review of paradigm shifts) |
| ETH-010 | ADM-008 (reporting structures) | ADM-008 (power distribution in reporting), all domains (power checks) |
| ETH-011 | RES-001 (epistemological foundations for ethical reasoning) | EVL-001 (ethical evolution informs institutional evolution philosophy) |
| ETH-015 | DR-013 (post-incident review methodology) | DR-013 (ethical incidents feed post-incident learning) |

---

## 15.5 Writing Schedule

| Phase | Articles | Duration | Notes |
|---|---|---|---|
| Phase 1: Foundations | ETH-001, ETH-002 | Weeks 1-5 | These are the most important articles in the entire documentation set. They define what the institution will not do. Allow extensive deliberation time. |
| Phase 2: Structural Protections | ETH-003, ETH-007, ETH-004 | Weeks 6-12 | Safeguards, whistleblower protections, and oversight. These must be established before the institution begins significant operations. |
| Phase 3: Power & Rights | ETH-010, ETH-008, ETH-012 | Weeks 13-19 | Power concentration prevention, personnel rights, and transparency. |
| Phase 4: Applied Ethics | ETH-006, ETH-005, ETH-014 | Weeks 20-27 | Crisis ethics, procurement ethics, and conflict of interest. |
| Phase 5: Review & Prevention | ETH-009, ETH-013, ETH-015 | Weeks 28-35 | Ethical review of changes, harm prevention, and incident response. |
| Phase 6: Meta-Ethics | ETH-011 | Weeks 36-40 | Ethical evolution is the capstone. It must be written with extreme care. Allow extra time for deliberation and review. |

**Total estimated duration:** 40 weeks (approximately 10 months).

**Note:** Ethics is intentionally given the longest writing schedule of all domains. These articles have the most profound consequences and require the most careful deliberation. Rushing ethical frameworks is itself an ethical risk.

---

## 15.6 Review Process

### 15.6.1 Review Stages

1. **Self-Review** (Author, 2 days): Extended self-review. Author must specifically check for unstated assumptions, cultural biases, and power dynamics embedded in the text.
2. **Diverse Perspective Review** (2-3 reviewers with different backgrounds, 5 days): Ethics articles must be reviewed by people with different life experiences, values, and perspectives. Homogeneous review is a failure mode.
3. **Adversarial Ethics Review** (1-2 reviewers, 3 days): Reviewers specifically try to find ways the ethical framework could be abused, circumvented, or weaponized. They ask: "How could a bad actor use these rules to cause harm?"
4. **Cross-Domain Review** (All domain leads, 7 days): Every domain lead reviews Ethics articles because they constrain all domains. Extended review period.
5. **"Future Person" Review** (1 reviewer, 3 days): Reviewer asks: "Would a person in 50 years, with different cultural values, find these ethics reasonable? What might they find abhorrent that we find acceptable? What might they find acceptable that we prohibit?"
6. **Final Acceptance** (Author + Ethics Officer + all domain leads): Acceptance of Ethics articles requires broader sign-off than any other domain.

### 15.6.2 Special Review Requirements

- ETH-002 (Red Lines) requires unanimous acceptance from all domain leads. A single objection sends it back to revision.
- ETH-011 (Ethical Evolution) requires a supermajority (not unanimity) for acceptance, because unanimity on meta-ethics may be unachievable.
- ETH-007 (Whistleblower Protections) must be reviewed by someone who has no institutional authority (to check that the protections are adequate from the perspective of the least powerful person).
- No Ethics article may be accepted if fewer than five distinct individuals have reviewed it.

### 15.6.3 Dispute Resolution

Disputes in Ethics articles are the most consequential in the entire documentation set. The dispute resolution process is:

1. Dispute is documented in full in the Commentary Section.
2. All parties present their reasoning in writing.
3. A facilitated discussion is held (not a debate — a genuine attempt to understand all perspectives).
4. If agreement is reached, the article is revised.
5. If agreement is not reached, both positions are included in the article with clear labeling. The article is published with the dispute visible. Future maintainers can resolve it with the benefit of additional experience and perspective.
6. Red line disputes (ETH-002) that cannot be resolved default to the more restrictive position.

---

## 15.7 Maintenance Plan

### 15.7.1 Review Cadence

| Review Type | Frequency | Trigger |
|---|---|---|
| Scheduled Full Review | Every 3 years | Calendar. Ethics requires regular re-examination. |
| Incident-Triggered Review | After every ethical incident | Mandatory, regardless of severity. |
| Cultural Shift Review | Every 5 years | Deliberate reassessment of whether cultural changes require ethical updates. |
| New Member Review | When new members join | New members review ethical articles and provide fresh-eyes feedback. |
| Red Line Review | Every 10 years | Formal reassessment of absolute prohibitions with full deliberative process. |

### 15.7.2 Version Control

Same as Domain 11 (major.minor.patch, all versions retained), with one addition:
- **Red Line changes** (ETH-002) require a special version designation: the version number is prefixed with "RL-" to make red line changes visually distinct in the version history.
- All previous versions of ETH-002 are retained and prominently accessible, so that the evolution of red lines can be traced over the institution's entire history.

### 15.7.3 Succession

- The Ethics Officer role must never be filled by someone who also holds operational authority in another domain. Separation of ethical oversight from operational power is a structural requirement.
- At least three people must be deeply familiar with ETH-001 and ETH-002 at all times.
- New Ethics Officers must serve a minimum 3-year apprenticeship under an existing Ethics Officer before assuming the role.
- If no qualified Ethics Officer is available, the institution enters a conservatorship mode where no significant decisions are made until the role is filled. This is defined in ETH-003 (Safeguard Mechanisms).

### 15.7.4 Decay Detection

- If the whistleblower channel (ETH-007) has never been used, either the institution is ethically perfect (unlikely) or the channel is not trusted. Investigate.
- If ethical reviews (ETH-009) are being rubber-stamped without substantive engagement, the oversight process is decaying.
- If the Commentary Sections of Ethics articles have no entries from people outside the Ethics domain, the domain is becoming insular.
- If power concentration checks (ETH-010) have not flagged any concerns in 5 years, the checks may not be functioning.
- If the "Future Person" review (15.6.1, stage 5) consistently produces no concerns, the reviewers may not be engaging seriously with the exercise.

---
---

# CROSS-DOMAIN INTEGRATION

## Global Dependency Summary

The five domains form an interconnected system. The following diagram summarizes the primary directional dependencies:

```
                    ETH-001 (Ethical Foundations)
                         |
          +--------------+--------------+
          |              |              |
          v              v              v
     DR-001          EVL-001        RES-001
  (DR Philosophy)  (Evolution     (Research
                    Philosophy)    Philosophy)
          |              |              |
          v              v              v
     DR-002          EVL-010        RES-002
  (Threat Model)   (Design for    (Methodology)
                    Longevity)
          |              |              |
          +------+-------+------+------+
                 |              |
                 v              v
            ADM-001         ADM-004
         (Admin          (Inventory)
          Principles)
```

**Key integration points:**
1. ETH-001 is the root document. All domains ultimately trace their authority back to the ethical foundations.
2. ADM-004 (Inventory) and ADM-006 (Personnel) are the most depended-upon administrative articles.
3. DR-002 (Threat Model) and EVL-003 (Change Impact Assessment) share analytical methods with RES-002 (Methodology).
4. The DR-Administration relationship is the tightest operational coupling.
5. The Ethics-Evolution relationship is the most important governance coupling.

## Global Writing Sequence Recommendation

Given the cross-domain dependencies, the recommended global writing sequence is:

| Global Phase | Domain | Articles | Rationale |
|---|---|---|---|
| 1 | Ethics | ETH-001, ETH-002 | Ethical foundations must exist first. They constrain everything. |
| 2 | All domains | *-001 (all foundational articles) | Each domain's philosophy article, written in light of ETH-001. |
| 3 | Administration | ADM-007, ADM-004, ADM-006 | Administrative infrastructure needed by all other domains. |
| 4 | DR | DR-002, DR-003, DR-004 | Threat model and backup strategy are urgent safety requirements. |
| 5 | Ethics | ETH-003, ETH-007, ETH-004 | Structural protections before significant operations begin. |
| 6 | All domains | Remaining articles per domain schedule | Each domain follows its internal priority order. |
| 7 | All domains | Capstone articles (ADM-013, DR-010, EVL-009, RES-011, ETH-011) | Final, integrative articles written last. |

## Global Review Coordination

- When an article in one domain is revised, all dependent articles in other domains must be notified and checked for consistency.
- A cross-domain review meeting should occur quarterly to identify integration issues.
- A master dependency register must be maintained (by Administration) that tracks all cross-domain references and flags stale links.

## Global Maintenance Coordination

- Annual cross-domain consistency audit: verify that no two domains contain contradictory guidance.
- Five-year comprehensive review: all domains reviewed together to identify systemic drift.
- Ten-year foundational review: ETH-001, all *-001 articles, and EVL-009 reviewed together to assess whether foundational assumptions still hold.

---
---

# APPENDIX A: ARTICLE TEMPLATE (FULL)

For reference, the complete article template that all articles must follow:

```
ARTICLE ID: [DOMAIN-NNN]
TITLE: [Clear, unambiguous, searchable title]
VERSION: [major.minor.patch]
DATE: [YYYY-MM-DD]
AUTHOR: [Name]
REVIEWED BY: [Names, dates]
STATUS: [Draft | Under Review | Accepted | Deprecated]

1. PURPOSE
   Why this article exists. What problem it solves. What would go wrong
   if this article did not exist.

2. SCOPE
   What this article covers. What it explicitly does NOT cover.
   Boundary conditions.

3. BACKGROUND
   Historical context. Why this approach was chosen. What alternatives
   were considered and rejected. Prior decisions that inform this article.

4. SYSTEM MODEL
   How the system or process works at a conceptual level. Diagrams
   encouraged. Mental models for the reader. Assumptions stated explicitly.

5. RULES & CONSTRAINTS
   Hard rules. Invariants. Non-negotiable boundaries. Stated as clear,
   testable propositions. Each rule includes its rationale.

6. FAILURE MODES
   What can go wrong. How you know it has gone wrong. Early warning signs.
   Probability and impact assessment where possible.

7. RECOVERY PROCEDURES
   For each failure mode: step-by-step restoration procedure.
   Prerequisites, tools needed, expected duration, success criteria.

8. EVOLUTION PATH
   How this article and its subject are expected to change over time.
   Known limitations of the current approach. Conditions that would
   trigger a revision. Predicted obsolescence timeline.

9. COMMENTARY SECTION
   Append-only. Maintainer notes, dissenting opinions, lessons learned,
   observations from practice. Each entry dated and attributed.
   This section is never edited, only appended.

10. REFERENCES
    Cross-references to other articles (by Article ID).
    External sources (with full citation).
    Historical documents.
    Related articles in other domains.
```

---

# APPENDIX B: META-RULES FOR ALL DOCUMENTATION

These rules apply to every article across all five domains:

1. **Self-Containment:** Every article must be understandable without access to the original author. Assume the reader has no one to ask.

2. **Explicit Assumptions:** Every assumption must be stated. "Everyone knows" is not an assumption — it is a failure of documentation.

3. **Cultural Neutrality:** Avoid idioms, metaphors, and references that depend on a specific cultural context. Write for a reader whose cultural background is unknown.

4. **Temporal Robustness:** Do not reference dates, events, or technologies as though they are self-evidently understood. Explain context. A reader in 2076 should be able to understand an article written in 2026.

5. **Graceful Degradation:** If part of the documentation set is lost, the remaining articles should still be useful. Minimize single points of failure in the documentation structure.

6. **Honesty About Uncertainty:** If something is uncertain, say so. Do not present guesses as facts. Use the knowledge classification system (RES-006) when referencing knowledge claims.

7. **No Orphan References:** Every cross-reference must point to an article that exists (or is explicitly marked as "planned"). Dead references are documentation rot.

8. **Commentary Preservation:** The Commentary Section of every article is append-only and permanent. It is the institutional memory for that article. Deleting commentary entries is a red line violation (ETH-002).

9. **Plain Language:** Write at a level that a competent, educated person outside your specialty can understand. Jargon must be defined on first use.

10. **Versioning Discipline:** Every change, no matter how small, produces a new version number. The version history is the article's autobiography.

---

**END OF STAGE 1: DOCUMENTATION FRAMEWORK — DOMAINS 11-15**
