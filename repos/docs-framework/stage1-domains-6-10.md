# STAGE 1: DOCUMENTATION FRAMEWORK
# Domains 6-10

**Document ID:** STAGE1-D6D10-001
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Initial Framework
**Scope:** Data & Archives, Intelligence & Analysis, Automation & Agents, Education & Training, User Operations
**Lifespan Target:** 50+ years
**Assumptions:** Air-gapped operation, off-grid power, self-built infrastructure, future loss of original author, hardware generational change, cultural and linguistic drift over decades.

---

## META-RULES GOVERNING ALL DOMAINS

1. Every article assumes the original author is dead or unreachable.
2. Every article assumes the reader has no prior context beyond what is written in this documentation system.
3. Every article assumes hardware will be replaced multiple times across incompatible generations.
4. Every article assumes cultural norms, language usage, and institutional priorities will shift.
5. Every article must be intelligible to a competent generalist with no specialist training, or must explicitly reference prerequisite articles.
6. No article may depend on network connectivity, cloud services, or external institutions for its core content.
7. All formats referenced must be open, documented, and reproducible from specification alone.
8. Every article must declare its own obsolescence conditions.

---

## STANDARD ARTICLE TEMPLATE

All articles across all domains must conform to this structure:

```
ARTICLE: [ID]
TITLE: [Descriptive title]
VERSION: [Semantic version]
LAST REVIEWED: [Date]
REVIEW CYCLE: [Frequency]
OBSOLESCENCE CONDITIONS: [When this article should be retired or rewritten]

1. PURPOSE
   Why this article exists. What question it answers. What problem it solves.

2. SCOPE
   What this article covers and, critically, what it does not cover.

3. BACKGROUND
   Historical context, design rationale, prior decisions that led here.

4. SYSTEM MODEL
   How the relevant system works. Diagrams encouraged (in durable formats).
   Mental models. Abstractions. Component relationships.

5. RULES & CONSTRAINTS
   Hard rules that must not be violated.
   Soft rules that may be bent with documented justification.
   Constraints imposed by physics, hardware, or upstream decisions.

6. FAILURE MODES
   What can go wrong. How you know it has gone wrong.
   Severity classification. Probability estimates where possible.

7. RECOVERY PROCEDURES
   Step-by-step actions for each failure mode.
   Decision trees for ambiguous situations.
   Escalation paths.

8. EVOLUTION PATH
   How this system, process, or policy is expected to change over time.
   Known limitations of the current approach.
   Signposts that indicate when change is needed.

9. COMMENTARY SECTION
   Informal notes from authors and maintainers.
   "Why we did it this way and not that way."
   Warnings, opinions, hard-won lessons.

10. REFERENCES
    Cross-references to other articles in this system.
    External references (with full bibliographic data, not URLs).
    Standards documents cited.
```

---
---

# DOMAIN 6: DATA & ARCHIVES

**Agent Role:** Data Systems Designer
**Core Purpose:** Define data models, storage philosophy, archival strategy, format longevity, backup doctrine, and data sovereignty for a self-contained institution operating across hardware generations and decades of cultural change.

---

## 6.1 DOMAIN MAP

### Scope

Domain 6 governs everything related to the persistence, organization, integrity, and retrievability of institutional data. This includes:

- The philosophy of what data is worth keeping and why.
- The physical and logical organization of stored data.
- The selection and management of data formats for multi-decade survival.
- The backup, replication, and disaster recovery doctrine.
- Metadata standards and cataloguing systems.
- The migration strategy for moving data across hardware and software generations.
- Data sovereignty: who controls data, who accesses it, under what authority.
- The archival lifecycle from creation through active use to cold storage to eventual disposal or permanent preservation.

### Boundaries

This domain does NOT cover:

- The content of intelligence analysis (Domain 7) -- only the storage and retrieval of raw data and finished products.
- The operational procedures for daily data handling (Domain 10) -- only the rules and systems those procedures operate within.
- The automation of data workflows (Domain 8) -- only the data structures and integrity requirements that automation must respect.
- The training of personnel in data handling (Domain 9) -- only the specification of what they must learn.
- Hardware procurement or physical infrastructure (assumed covered in Domains 1-5).

### Relationships

| Related Domain | Relationship |
|---|---|
| Domain 7 (Intelligence) | D6 stores raw inputs and finished intelligence products. D7 defines what is worth collecting; D6 defines how it is stored. |
| Domain 8 (Automation) | Automated agents read from and write to D6-governed data stores. D6 defines integrity constraints; D8 defines the agents that operate within them. |
| Domain 9 (Education) | Training materials are themselves data artifacts governed by D6 archival policy. D9 defines curriculum; D6 ensures curriculum materials survive. |
| Domain 10 (Operations) | Daily operational logs, checklists, and incident records are D6 data assets. D10 defines what gets recorded; D6 defines how it is stored and for how long. |
| Domains 1-5 | Physical storage media, power systems, and network topology constrain D6 design choices. D6 must document its hardware assumptions explicitly. |

### Key Tensions

- **Comprehensiveness vs. Storage Limits:** Keeping everything is impossible off-grid. The domain must define triage.
- **Format Purity vs. Practicality:** The purest archival formats may be the hardest to work with daily.
- **Access Speed vs. Archival Durability:** Active data and archival data have conflicting optimization targets.
- **Centralization vs. Redundancy:** Single-source-of-truth vs. distributed copies for survivability.

---

## 6.2 ARTICLE LIST

| Article ID | Title | Summary |
|---|---|---|
| D6-001 | Data Philosophy: What We Keep and Why | Foundational principles for data retention. Triage framework. The institutional stance on information hoarding vs. curation. |
| D6-002 | The Canonical Data Model | Master schema for how all institutional data is structured. Entity types, relationships, namespaces. |
| D6-003 | Format Longevity Doctrine | Which file formats are approved for archival use, why, and the process for evaluating new formats. Deprecation procedures. |
| D6-004 | Storage Architecture: Logical Topology | How data is organized across storage tiers -- hot, warm, cold, archival. Directory structures. Naming conventions. |
| D6-005 | Storage Architecture: Physical Media Strategy | Which physical media types are used, their expected lifespans, rotation schedules, and failure characteristics. |
| D6-006 | Backup Doctrine | Backup frequency, methods, verification procedures, off-site (but still air-gapped) storage, and restore testing cadence. |
| D6-007 | Data Integrity & Verification | Checksumming, hash chains, bit-rot detection, integrity audit schedules, and response to detected corruption. |
| D6-008 | Metadata Standards & Cataloguing | How every data artifact is described, tagged, and made findable. The catalogue system. Controlled vocabularies. |
| D6-009 | Data Migration Doctrine | Procedures for moving data across hardware generations, software platforms, and format versions without loss. |
| D6-010 | Data Sovereignty & Access Control | Who owns data, who may access it, under what authority, and how access is granted, revoked, and audited. |
| D6-011 | Retention Schedules & Disposal Policy | How long each category of data is kept. Conditions for early disposal. Conditions for permanent preservation. The disposal process itself. |
| D6-012 | The Institutional Memory Archive | Special provisions for preserving institutional knowledge, decisions, and context beyond raw data -- the "why" archive. |
| D6-013 | Disaster Recovery for Data Systems | Comprehensive plan for reconstituting data systems after catastrophic loss, from partial to total. |
| D6-014 | Data Ingest Procedures | How new data enters the system. Validation, quarantine, classification, and integration steps. |
| D6-015 | Print & Physical Backup Doctrine | When and how to maintain paper or other non-digital backups of critical data. Selection criteria, storage, indexing. |

---

## 6.3 PRIORITY ORDER

**Tier 1 -- Existential (write first, everything else depends on these):**

1. D6-001 Data Philosophy: What We Keep and Why
2. D6-003 Format Longevity Doctrine
3. D6-006 Backup Doctrine
4. D6-007 Data Integrity & Verification

**Tier 2 -- Structural (define the architecture):**

5. D6-002 The Canonical Data Model
6. D6-004 Storage Architecture: Logical Topology
7. D6-005 Storage Architecture: Physical Media Strategy
8. D6-008 Metadata Standards & Cataloguing

**Tier 3 -- Operational (make it work day-to-day):**

9. D6-014 Data Ingest Procedures
10. D6-010 Data Sovereignty & Access Control
11. D6-011 Retention Schedules & Disposal Policy

**Tier 4 -- Resilience (long-term survival):**

12. D6-009 Data Migration Doctrine
13. D6-013 Disaster Recovery for Data Systems
14. D6-012 The Institutional Memory Archive
15. D6-015 Print & Physical Backup Doctrine

### Rationale

You cannot build storage architecture without first knowing what you are storing and why (D6-001). You cannot select formats without the longevity doctrine (D6-003). Nothing matters if backups fail (D6-006) or integrity is compromised (D6-007). The structural layer depends on these foundations. Operational articles can be drafted in parallel once structures exist. Resilience articles codify long-term survival and can be refined iteratively.

---

## 6.4 DEPENDENCIES

### Internal Dependencies (within Domain 6)

```
D6-001 --> D6-002 (philosophy informs data model)
D6-001 --> D6-011 (philosophy informs retention)
D6-003 --> D6-002 (formats constrain data model)
D6-003 --> D6-009 (format changes drive migration)
D6-004 --> D6-005 (logical topology maps to physical media)
D6-006 --> D6-005 (backup methods depend on media)
D6-006 --> D6-007 (backups must be integrity-verified)
D6-007 --> D6-013 (integrity failure triggers disaster recovery)
D6-008 --> D6-002 (metadata extends the data model)
D6-008 --> D6-014 (ingest must apply metadata)
D6-009 --> D6-003 (migration depends on format doctrine)
D6-009 --> D6-005 (migration depends on media strategy)
D6-014 --> D6-008 (ingest applies metadata)
D6-014 --> D6-010 (ingest respects access control)
```

### Cross-Domain Dependencies

| D6 Article | Depends On | Provides To |
|---|---|---|
| D6-001 | D7-001 (Intelligence philosophy defines what raw data matters) | D7 (all articles rely on D6 storage guarantees) |
| D6-003 | Domains 1-5 (hardware constraints on format tooling) | D8 (automation must read/write approved formats) |
| D6-005 | Domains 1-5 (available hardware, power budget) | D10 (operations must maintain physical media) |
| D6-006 | D10-005 (backup is a daily operational procedure) | D13 (disaster recovery depends on backup integrity) |
| D6-008 | D9 (training materials need catalogue entries) | D7 (intelligence products are findable via catalogue) |
| D6-010 | D10-008 (access control is enforced operationally) | D8 (agents must respect access boundaries) |
| D6-014 | D7-003 (intelligence collection feeds ingest) | D10 (ingest is a daily operational workflow) |
| D6-015 | Domains 1-5 (physical storage space) | D13 (print backups are last-resort recovery) |

---

## 6.5 WRITING SCHEDULE

| Phase | Articles | Estimated Effort | Prerequisites |
|---|---|---|---|
| Phase 1 (Weeks 1-3) | D6-001, D6-003 | 2 articles, intensive | None (foundational) |
| Phase 2 (Weeks 3-5) | D6-006, D6-007 | 2 articles, intensive | D6-001, D6-003 |
| Phase 3 (Weeks 5-8) | D6-002, D6-004, D6-005 | 3 articles, moderate each | D6-001, D6-003 |
| Phase 4 (Weeks 8-10) | D6-008, D6-014 | 2 articles, moderate | D6-002, D6-004 |
| Phase 5 (Weeks 10-12) | D6-010, D6-011 | 2 articles, moderate | D6-001, D6-008 |
| Phase 6 (Weeks 12-16) | D6-009, D6-013, D6-012, D6-015 | 4 articles, variable | Most prior articles |

**Total estimated duration:** 16 weeks for initial drafts. Parallel work with other domains is expected.

---

## 6.6 REVIEW PROCESS

### Initial Review

Each article undergoes three review passes before acceptance:

1. **Technical Accuracy Review.** Does the article correctly describe the system, format, or process? Reviewer: a person with hands-on experience with the relevant technology. Timeline: within 2 weeks of draft.

2. **Clarity & Accessibility Review.** Can a competent generalist with no specialist background understand this article well enough to act on it? Reviewer: someone who was NOT involved in writing or designing the system. Timeline: within 1 week of technical review.

3. **Longevity & Assumptions Review.** Does this article make assumptions about hardware, software, culture, or availability that will not hold over 50 years? Does it declare its own obsolescence conditions? Reviewer: designated institutional continuity role. Timeline: within 1 week of clarity review.

### Ongoing Review

| Review Type | Frequency | Trigger |
|---|---|---|
| Scheduled review | Every 24 months | Calendar |
| Event-driven review | As needed | Hardware change, format deprecation, integrity incident, media failure |
| Comprehensive audit | Every 60 months | Calendar; all D6 articles reviewed together for coherence |

### Review Artifacts

Every review produces:
- A dated review note appended to the article's Commentary Section.
- An updated LAST REVIEWED date.
- A list of changes made or a statement of "no changes required with justification."
- An updated version number if changes were made.

---

## 6.7 MAINTENANCE PLAN

### Routine Maintenance

- **Quarterly:** Verify that all D6 articles are accessible and readable in their stored format. Spot-check cross-references for link rot within the documentation system.
- **Annually:** Run a complete format longevity assessment -- are any approved formats showing signs of ecosystem abandonment? Update D6-003 if needed.
- **Annually:** Verify that the backup doctrine (D6-006) is being followed by checking operational logs (D10 cross-reference).
- **Biannually:** Full review cycle per Section 6.6.

### Triggered Maintenance

- **New hardware generation adopted:** Immediately review D6-005 (physical media), D6-009 (migration), D6-004 (logical topology). Trigger migration planning.
- **Data integrity incident detected:** Immediately review D6-007 (integrity), D6-006 (backup), D6-013 (disaster recovery). Document the incident in Commentary.
- **New data type or source introduced:** Review D6-002 (data model), D6-014 (ingest), D6-008 (metadata). Extend as needed.
- **Personnel change in data steward role:** Review D6-010 (access control). Ensure successor has reviewed all Tier 1 articles.

### Succession Provisions

- All D6 articles must be stored in at least two independent physical locations.
- A printed summary index of all D6 articles must be maintained (cross-reference D6-015).
- The Data Systems Designer role must have a designated successor at all times, or the institution must have a documented procedure for bootstrapping someone into the role using D9 training materials.

---
---

# DOMAIN 7: INTELLIGENCE & ANALYSIS

**Agent Role:** Intelligence Analyst
**Core Purpose:** Define how the institution gathers, processes, and acts on information. Pattern recognition, signal analysis, decision support, and the epistemological discipline required to avoid self-deception in an air-gapped environment.

---

## 7.1 DOMAIN MAP

### Scope

Domain 7 governs the institution's relationship with information as a decision-making input. This includes:

- The intelligence cycle: requirements definition, collection, processing, analysis, dissemination, feedback.
- Source evaluation and reliability assessment.
- Analytical methods and frameworks for making sense of incomplete or contradictory information.
- Pattern recognition: identifying trends, anomalies, and emerging threats or opportunities.
- Decision support: how analysis products inform institutional action.
- Counter-deception and epistemic hygiene: guarding against bias, groupthink, and information corruption.
- Signal analysis: monitoring environmental, technical, and social signals relevant to institutional survival.
- The distinction between data (Domain 6) and intelligence (Domain 7): data is raw; intelligence is data that has been evaluated, analyzed, and contextualized for action.

### Boundaries

This domain does NOT cover:

- Raw data storage, format, or integrity (Domain 6).
- The automation of analytical processes (Domain 8) -- only the analytical frameworks that automation must serve.
- Training analysts (Domain 9) -- only the specification of analytical competencies required.
- The daily execution of intelligence routines (Domain 10) -- only the doctrinal frameworks those routines follow.
- Physical security or signals security (assumed in Domains 1-5).

### Relationships

| Related Domain | Relationship |
|---|---|
| Domain 6 (Data) | D6 provides the storage and retrieval infrastructure for both raw inputs and finished intelligence products. D7 defines collection requirements that inform D6 ingest procedures. |
| Domain 8 (Automation) | Automated agents can perform routine collection, pattern matching, and alerting. D7 defines the analytical standards; D8 defines the agents that support them. |
| Domain 9 (Education) | Analysts must be trained. D7 defines competency requirements; D9 designs the curriculum. |
| Domain 10 (Operations) | Intelligence products feed operational decisions. D10 defines how products are consumed; D7 defines how they are produced. |
| Domains 1-5 | Physical and network architecture constrain collection capabilities. Air-gapped operation fundamentally shapes what intelligence is possible. |

### Key Tensions

- **Comprehensiveness vs. Cognitive Overload:** Collecting everything means analyzing nothing. The domain must define focus.
- **Analytical Rigor vs. Speed:** Thorough analysis takes time; some decisions cannot wait.
- **Confidence vs. Uncertainty:** Intelligence must communicate uncertainty honestly without paralyzing decision-makers.
- **Air-Gap Constraint:** The institution cannot freely access the broader information environment. Intelligence doctrine must account for systematic blind spots.
- **Single-Analyst Risk:** In a small institution, the analyst may be one person. The domain must guard against single points of cognitive failure.

---

## 7.2 ARTICLE LIST

| Article ID | Title | Summary |
|---|---|---|
| D7-001 | Intelligence Philosophy: Knowing What We Don't Know | Foundational epistemological stance. What intelligence means for this institution. The role of humility, uncertainty, and structured doubt. |
| D7-002 | The Intelligence Cycle: Requirements to Action | The full lifecycle: defining what we need to know, collecting, processing, analyzing, disseminating, and incorporating feedback. |
| D7-003 | Collection Doctrine: Sources, Methods, Constraints | What sources are available to an air-gapped institution. How to evaluate source reliability. Collection priorities and methods. |
| D7-004 | Analytical Frameworks & Methods | Structured analytical techniques: competing hypotheses, red-teaming, scenario planning, trend extrapolation, pre-mortem analysis. |
| D7-005 | Cognitive Bias & Epistemic Hygiene | Catalogue of known cognitive biases with institutional countermeasures. Procedures for challenging assumptions. Devil's advocate protocols. |
| D7-006 | Threat Assessment Methodology | How to identify, classify, assess, and prioritize threats to the institution across all domains (physical, technical, social, environmental). |
| D7-007 | Opportunity Recognition Framework | The positive mirror of threat assessment -- identifying favorable conditions, emerging possibilities, and windows for action. |
| D7-008 | Intelligence Product Standards | Formats, templates, classification levels, and quality standards for intelligence products. How to communicate findings, confidence levels, and recommendations. |
| D7-009 | Environmental Monitoring: Signals & Indicators | What signals to monitor (weather, equipment degradation, social dynamics, resource levels) and how to build indicator frameworks. |
| D7-010 | Decision Support Doctrine | How intelligence products connect to institutional decision-making. Briefing protocols. The analyst's role and its limits -- informing, not deciding. |
| D7-011 | Estimative Language & Uncertainty Communication | A precise vocabulary for communicating probability, confidence, and uncertainty. Calibration exercises. |
| D7-012 | Intelligence Failure: Detection, Response, Learning | How to recognize when analysis was wrong. Post-mortem procedures. Institutional learning from intelligence failures. |
| D7-013 | Single-Analyst Resilience Protocols | Special provisions for when the intelligence function is performed by one person. Self-check procedures, structured aids, anti-bias toolkits. |

---

## 7.3 PRIORITY ORDER

**Tier 1 -- Foundational (epistemic bedrock):**

1. D7-001 Intelligence Philosophy: Knowing What We Don't Know
2. D7-005 Cognitive Bias & Epistemic Hygiene
3. D7-002 The Intelligence Cycle: Requirements to Action

**Tier 2 -- Methodological (how to do the work):**

4. D7-004 Analytical Frameworks & Methods
5. D7-003 Collection Doctrine: Sources, Methods, Constraints
6. D7-011 Estimative Language & Uncertainty Communication

**Tier 3 -- Applied (specific applications):**

7. D7-006 Threat Assessment Methodology
8. D7-009 Environmental Monitoring: Signals & Indicators
9. D7-007 Opportunity Recognition Framework
10. D7-008 Intelligence Product Standards

**Tier 4 -- Resilience (guarding against failure):**

11. D7-010 Decision Support Doctrine
12. D7-012 Intelligence Failure: Detection, Response, Learning
13. D7-013 Single-Analyst Resilience Protocols

### Rationale

Without epistemic humility and bias awareness (D7-001, D7-005), all subsequent analysis is unreliable. The cycle (D7-002) provides the operational skeleton. Methods and collection doctrine make the work possible. Applied articles address specific use cases. Resilience articles protect the function itself from degradation.

---

## 7.4 DEPENDENCIES

### Internal Dependencies (within Domain 7)

```
D7-001 --> D7-002 (philosophy shapes the cycle)
D7-001 --> D7-005 (philosophy demands epistemic hygiene)
D7-002 --> D7-003 (the cycle defines when collection happens)
D7-002 --> D7-004 (the cycle defines when analysis happens)
D7-002 --> D7-008 (the cycle defines product requirements)
D7-004 --> D7-006 (threat assessment is an applied analytical framework)
D7-004 --> D7-007 (opportunity recognition is an applied analytical framework)
D7-005 --> D7-013 (bias awareness is critical for single-analyst resilience)
D7-005 --> D7-012 (bias is a primary cause of intelligence failure)
D7-008 --> D7-010 (products must meet decision-support needs)
D7-008 --> D7-011 (products must use estimative language)
D7-009 --> D7-006 (monitoring feeds threat assessment)
D7-009 --> D7-007 (monitoring feeds opportunity recognition)
```

### Cross-Domain Dependencies

| D7 Article | Depends On | Provides To |
|---|---|---|
| D7-001 | None (foundational) | D6-001 (intelligence philosophy informs data philosophy) |
| D7-003 | D6-014 (collection feeds data ingest) | D6-001 (collection requirements inform what data to keep) |
| D7-004 | D9 (analysts need training in methods) | D8 (automation can support structured methods) |
| D7-006 | Domains 1-5 (physical threat context) | D10 (threat assessments drive operational decisions) |
| D7-008 | D6-003 (products stored in approved formats) | D10 (products consumed in daily operations) |
| D7-009 | D8 (automated monitoring agents) | D10 (monitoring is daily operational work) |
| D7-010 | D10 (decision-making is operational) | D10 (intelligence feeds operations) |
| D7-013 | D9 (training for self-check procedures) | D9 (defines competency requirements for analyst training) |

---

## 7.5 WRITING SCHEDULE

| Phase | Articles | Estimated Effort | Prerequisites |
|---|---|---|---|
| Phase 1 (Weeks 1-3) | D7-001, D7-005 | 2 articles, intensive (conceptual and research-heavy) | None |
| Phase 2 (Weeks 3-5) | D7-002, D7-011 | 2 articles, moderate | D7-001 |
| Phase 3 (Weeks 5-8) | D7-004, D7-003 | 2 articles, intensive (methodological detail) | D7-002, D7-005 |
| Phase 4 (Weeks 8-11) | D7-006, D7-009, D7-007 | 3 articles, moderate each | D7-004, D7-003 |
| Phase 5 (Weeks 11-13) | D7-008, D7-010 | 2 articles, moderate | D7-002, D7-011 |
| Phase 6 (Weeks 13-16) | D7-012, D7-013 | 2 articles, moderate | D7-005, D7-004 |

**Total estimated duration:** 16 weeks for initial drafts.

---

## 7.6 REVIEW PROCESS

### Initial Review

Each article undergoes three review passes:

1. **Analytical Soundness Review.** Are the methods well-defined? Are the frameworks internally consistent? Does the article avoid the very biases it warns against? Reviewer: a second analyst or a generalist with strong critical thinking. Timeline: within 2 weeks of draft.

2. **Practical Applicability Review.** Can this article actually be used by someone performing the intelligence function, including under stress, with limited time, and with incomplete information? Reviewer: someone who will perform or has performed the role. Timeline: within 1 week of soundness review.

3. **Adversarial Review.** What happens if this article's assumptions are wrong? What happens if the methods are applied poorly? Does the article account for its own misuse? Reviewer: designated contrarian or red-team role. Timeline: within 1 week of applicability review.

### Ongoing Review

| Review Type | Frequency | Trigger |
|---|---|---|
| Scheduled review | Every 18 months | Calendar (shorter cycle due to evolving threat landscape) |
| Post-incident review | As needed | Intelligence failure, missed threat, false alarm |
| Methodology refresh | Every 36 months | Calendar; assess whether analytical methods remain fit for purpose |
| Comprehensive audit | Every 60 months | Calendar; full domain coherence review |

### Review Artifacts

Same as Domain 6: dated commentary note, updated review date, change list or no-change justification, version increment if changed.

---

## 7.7 MAINTENANCE PLAN

### Routine Maintenance

- **Quarterly:** Review current intelligence requirements. Are we asking the right questions? Is collection aligned with actual needs?
- **Semiannually:** Calibration exercise for estimative language (D7-011). Test whether stated confidence levels matched outcomes.
- **Annually:** Full review of threat assessment (D7-006) and opportunity framework (D7-007) for currency and relevance.
- **Every 18 months:** Full article review cycle per Section 7.6.

### Triggered Maintenance

- **Intelligence failure occurs:** Immediately invoke D7-012 post-mortem procedures. Review implicated articles. Update Commentary sections.
- **New threat category identified:** Review D7-006, D7-009. Extend if needed.
- **Analyst personnel change:** Immediately review D7-013 (single-analyst resilience). Ensure successor has completed D9 training for all D7 competencies.
- **New collection source becomes available:** Review D7-003. Evaluate source reliability per established doctrine.

### Succession Provisions

- The Intelligence Analyst role must produce a standing "state of analysis" document updated quarterly, summarizing current assessments, active requirements, and key uncertainties. This document is itself a D6 archival artifact.
- All analytical frameworks (D7-004) must be documented well enough that a successor can apply them without oral instruction.
- D7-013 (single-analyst resilience) must be the first article a new analyst reads.

---
---

# DOMAIN 8: AUTOMATION & AGENTS

**Agent Role:** Automation Architect
**Core Purpose:** Define the automation philosophy, agent boundaries, task delegation rules, human-in-the-loop requirements, autonomous operation limits, and the governance of automated systems for a self-contained institution that must remain controllable and comprehensible across decades.

---

## 8.1 DOMAIN MAP

### Scope

Domain 8 governs all automated processes, scripts, agents, and decision-making systems that operate with reduced or no human intervention. This includes:

- The philosophy of automation: what should be automated, what must not be, and why.
- Agent design principles: how automated agents are scoped, constrained, and monitored.
- Human-in-the-loop requirements: which decisions require human confirmation and how that confirmation is obtained.
- Autonomous operation limits: hard ceilings on what automation may do without human awareness.
- Task delegation: how work is divided between humans and automated systems.
- Agent lifecycle: creation, testing, deployment, monitoring, retirement of automated processes.
- Failure modes of automation: cascading failures, runaway processes, silent errors, automation surprise.
- The governance of automation: who may create, modify, or retire automated processes and under what authority.

### Boundaries

This domain does NOT cover:

- The data that automated systems operate on (Domain 6) -- only the rules for how agents interact with data.
- The analytical methods automated systems might support (Domain 7) -- only the boundaries within which automation operates.
- The training of personnel to manage automation (Domain 9) -- only the specification of competencies required.
- The daily operation and monitoring of automated systems (Domain 10) -- only the design and governance frameworks.
- The specific hardware or software platforms (assumed in Domains 1-5) -- only the automation principles that are platform-independent.

### Relationships

| Related Domain | Relationship |
|---|---|
| Domain 6 (Data) | Automated agents are major consumers and producers of data. D6 defines integrity constraints; D8 defines agents that must respect them. |
| Domain 7 (Intelligence) | Automation can support collection, processing, and pattern recognition. D7 defines analytical standards; D8 defines the automation that serves those standards. |
| Domain 9 (Education) | Personnel must be trained to create, manage, and override automated systems. D8 defines competency requirements; D9 designs training. |
| Domain 10 (Operations) | Automated systems are part of daily operations. D10 defines monitoring and intervention procedures; D8 defines the systems being monitored. |
| Domains 1-5 | Hardware and power constraints fundamentally limit what automation is possible. Agent design must account for power budgets, hardware reliability, and air-gapped constraints. |

### Key Tensions

- **Automation Efficiency vs. Human Comprehension:** The more automated a system is, the harder it is for humans to understand what it is doing and why. Over decades, this becomes an existential risk.
- **Autonomy vs. Control:** Greater automation requires less human effort but introduces risk of unintended action. The tension intensifies when the humans who built the system are gone.
- **Complexity vs. Maintainability:** Sophisticated automation requires sophisticated maintenance. In a small institution, maintainability must win.
- **Brittleness vs. Adaptability:** Highly automated systems can be brittle -- small changes in conditions cause failures. Manual processes are slower but more adaptable.
- **Trust Decay:** Future maintainers may not understand or trust inherited automation. The domain must account for the "black box" problem across generations.

---

## 8.2 ARTICLE LIST

| Article ID | Title | Summary |
|---|---|---|
| D8-001 | Automation Philosophy: The Case for Restraint | Foundational principles for when and why to automate. The default stance is human execution; automation must justify itself against specific criteria. |
| D8-002 | Agent Design Principles | How automated agents are designed: single-purpose, minimal authority, observable, interruptible, explainable. The "smallest useful agent" principle. |
| D8-003 | Human-in-the-Loop Doctrine | Which actions always require human confirmation. How confirmation is obtained. Timeout and default behaviors when a human is unavailable. |
| D8-004 | Autonomous Operation Limits | Hard ceilings on autonomous action. What automation may never do without human presence. Escalation triggers. Emergency halt procedures. |
| D8-005 | Task Delegation Framework | How to decide whether a task should be performed by a human, a supervised agent, or an autonomous agent. Decision matrix and criteria. |
| D8-006 | Agent Lifecycle: Build, Test, Deploy, Monitor, Retire | The full lifecycle of an automated process from conception through retirement. Stage gates, testing requirements, deployment approval, retirement criteria. |
| D8-007 | Agent Monitoring & Observability | How automated systems report their status, actions, and errors. Logging requirements. Dashboard design principles. Alerting thresholds. |
| D8-008 | Automation Failure Modes & Cascading Risk | Catalogue of how automated systems fail: silent errors, runaway processes, cascading failures, stale data, infinite loops, resource exhaustion. |
| D8-009 | Emergency Override & Kill Procedures | How to immediately halt any automated process. Physical and logical kill switches. Recovery after emergency halt. Post-incident review requirements. |
| D8-010 | Automation Governance: Authority & Accountability | Who may create, modify, deploy, or retire automated processes. Approval requirements. Change control. Accountability chain. |
| D8-011 | Legacy Automation: Understanding Inherited Systems | How future maintainers should approach automated systems they did not build. Comprehension procedures, testing strategies, and the decision to keep, modify, or retire. |
| D8-012 | Automation & Data Integrity Contracts | The interface between automated agents and data systems. What guarantees agents must provide regarding data they read or write. Transaction boundaries. |
| D8-013 | Scheduled vs. Event-Driven Automation | Design patterns for time-based and event-based automation. Trade-offs, failure characteristics, and selection criteria. |

---

## 8.3 PRIORITY ORDER

**Tier 1 -- Foundational (the doctrine of restraint):**

1. D8-001 Automation Philosophy: The Case for Restraint
2. D8-003 Human-in-the-Loop Doctrine
3. D8-004 Autonomous Operation Limits

**Tier 2 -- Design (how to build correctly):**

4. D8-002 Agent Design Principles
5. D8-005 Task Delegation Framework
6. D8-012 Automation & Data Integrity Contracts

**Tier 3 -- Operations (how to run and watch):**

7. D8-006 Agent Lifecycle: Build, Test, Deploy, Monitor, Retire
8. D8-007 Agent Monitoring & Observability
9. D8-009 Emergency Override & Kill Procedures

**Tier 4 -- Resilience (failure and legacy):**

10. D8-008 Automation Failure Modes & Cascading Risk
11. D8-010 Automation Governance: Authority & Accountability
12. D8-011 Legacy Automation: Understanding Inherited Systems
13. D8-013 Scheduled vs. Event-Driven Automation

### Rationale

The philosophy of restraint (D8-001) must come first because the greatest risk in automation is doing too much of it. Human-in-the-loop requirements (D8-003) and hard limits (D8-004) establish non-negotiable boundaries before any agent is designed. Design principles and delegation frameworks provide the "how." Lifecycle and observability enable ongoing operation. Failure modes and governance protect the institution long-term. Legacy automation (D8-011) addresses the multi-generational challenge directly.

---

## 8.4 DEPENDENCIES

### Internal Dependencies (within Domain 8)

```
D8-001 --> D8-002 (philosophy constrains design)
D8-001 --> D8-003 (philosophy defines human role)
D8-001 --> D8-005 (philosophy informs delegation criteria)
D8-003 --> D8-004 (human-in-the-loop defines autonomy ceiling)
D8-002 --> D8-006 (design principles are applied throughout lifecycle)
D8-002 --> D8-007 (observability is a design principle)
D8-002 --> D8-012 (data contracts are a design requirement)
D8-004 --> D8-009 (autonomy limits define when override is needed)
D8-006 --> D8-010 (lifecycle requires governance at each stage)
D8-006 --> D8-011 (retirement criteria inform legacy assessment)
D8-007 --> D8-008 (monitoring detects failure modes)
D8-008 --> D8-009 (failure detection triggers override)
D8-005 --> D8-013 (delegation informs scheduling pattern choice)
```

### Cross-Domain Dependencies

| D8 Article | Depends On | Provides To |
|---|---|---|
| D8-001 | D7-001 (intelligence philosophy -- automation must not replace judgment) | D10 (operational procedures reference automation philosophy) |
| D8-002 | Domains 1-5 (hardware platform capabilities) | D6 (agents must conform to data system interfaces) |
| D8-003 | D10 (human availability is an operational reality) | D7 (analytical automation requires human judgment confirmation) |
| D8-007 | D6-008 (monitoring logs are data artifacts with metadata) | D10 (operations consumes monitoring dashboards) |
| D8-009 | Domains 1-5 (physical kill switches) | D10 (override is an operational emergency procedure) |
| D8-010 | D10-008 (governance integrates with access control) | D9 (governance roles require training) |
| D8-011 | D9 (training for comprehension of inherited systems) | All domains (all domains may contain inherited automation) |
| D8-012 | D6-002 (data model defines what agents interact with) | D6-007 (data integrity depends on agent compliance) |

---

## 8.5 WRITING SCHEDULE

| Phase | Articles | Estimated Effort | Prerequisites |
|---|---|---|---|
| Phase 1 (Weeks 1-3) | D8-001, D8-003 | 2 articles, intensive (doctrinal, philosophical) | None |
| Phase 2 (Weeks 3-5) | D8-004, D8-002 | 2 articles, intensive | D8-001, D8-003 |
| Phase 3 (Weeks 5-8) | D8-005, D8-012, D8-006 | 3 articles, moderate each | D8-001, D8-002 |
| Phase 4 (Weeks 8-11) | D8-007, D8-009 | 2 articles, moderate | D8-002, D8-004 |
| Phase 5 (Weeks 11-14) | D8-008, D8-010, D8-013 | 3 articles, moderate | D8-006, D8-007 |
| Phase 6 (Weeks 14-16) | D8-011 | 1 article, intensive (conceptually challenging) | Most prior articles |

**Total estimated duration:** 16 weeks for initial drafts.

---

## 8.6 REVIEW PROCESS

### Initial Review

Each article undergoes three review passes:

1. **Safety & Constraint Review.** Does this article adequately protect against automation overreach? Are the limits clear, enforceable, and testable? Does it account for what happens when the automation's creator is no longer available? Reviewer: the most cautious person available -- someone who distrusts automation. Timeline: within 2 weeks of draft.

2. **Implementability Review.** Can the principles in this article actually be implemented with the institution's available technology and personnel? Are the requirements realistic for a small, off-grid, air-gapped operation? Reviewer: someone with hands-on systems experience. Timeline: within 1 week of safety review.

3. **Generational Transfer Review.** Will this article make sense to someone 20 years from now who did not build the systems it describes? Does D8-011 (legacy automation) adequately prepare them? Reviewer: a non-specialist. Timeline: within 1 week of implementability review.

### Ongoing Review

| Review Type | Frequency | Trigger |
|---|---|---|
| Scheduled review | Every 24 months | Calendar |
| Post-incident review | As needed | Automation failure, near-miss, override event |
| Technology change review | As needed | New hardware, new software platform, new automation tool adopted |
| Comprehensive audit | Every 60 months | Calendar; full domain coherence review with emphasis on legacy assessment |

---

## 8.7 MAINTENANCE PLAN

### Routine Maintenance

- **Monthly:** Review automation monitoring logs (cross-reference D10). Are all automated processes behaving as expected? Any anomalies?
- **Quarterly:** Verify that all emergency override procedures (D8-009) still function. Test kill switches.
- **Annually:** Review the inventory of all active automated processes against D8-006 lifecycle records. Are any past their intended retirement date? Are any undocumented?
- **Biannually:** Full article review cycle per Section 8.6.

### Triggered Maintenance

- **Automation failure:** Invoke D8-008 failure analysis. Review D8-009 override procedures. Update implicated articles.
- **New automation deployed:** Ensure it conforms to D8-002 (design principles), D8-003 (human-in-the-loop), D8-010 (governance). Add to lifecycle registry.
- **Personnel change in automation role:** Review D8-010 (governance), D8-011 (legacy systems). Ensure successor has been trained per D9 requirements.
- **Hardware platform change:** Review all active automation for compatibility. D8-006 lifecycle may trigger retirements.

### Succession Provisions

- Every automated process must have a human-readable description stored independently of the automation itself (not just code comments, but a standalone document per D8-006 requirements).
- D8-011 (Legacy Automation) is the succession document of last resort. It must be comprehensive enough that a new maintainer can approach any inherited automation safely.
- A complete inventory of all active automated processes must be maintained in the D6 data archive, with each entry linking to its D8-006 lifecycle record.

---
---

# DOMAIN 9: EDUCATION & TRAINING

**Agent Role:** Education Designer
**Core Purpose:** Define how knowledge is transferred across generations, how new maintainers learn the institution's systems and culture, how competency is assessed, and how the institution survives the loss of any individual's expertise.

---

## 9.1 DOMAIN MAP

### Scope

Domain 9 governs the institution's capacity to reproduce its own competence. This includes:

- The philosophy of institutional education: why formal knowledge transfer matters, especially in a self-contained system.
- Curriculum design for each institutional role and domain.
- Skill assessment and certification: how the institution knows someone is ready to perform a role.
- Mentorship and apprenticeship frameworks: structured relationships for knowledge transfer.
- Documentation as pedagogy: how the documentation system itself serves an educational function.
- Self-directed learning support: resources and methods for learning without a teacher.
- Tacit knowledge capture: strategies for preserving knowledge that is difficult to write down.
- Cross-training: ensuring multiple people can perform critical functions.
- Degraded-capacity training: how to maintain institutional knowledge when personnel are few.

### Boundaries

This domain does NOT cover:

- The content of what is taught (that lives in Domains 1-8 and 10) -- only the frameworks for how it is taught and learned.
- The daily operational execution of training sessions (Domain 10) -- only the design and assessment standards.
- The data storage of training materials (Domain 6) -- only the specification of what materials must exist.
- The automation of training delivery or assessment (Domain 8) -- only the educational requirements that automation must serve.

### Relationships

| Related Domain | Relationship |
|---|---|
| Domain 6 (Data) | Training materials are data artifacts stored per D6 policy. D6 ensures materials survive; D9 defines what materials must exist. |
| Domain 7 (Intelligence) | Intelligence analysis requires trained analysts. D7 defines competency requirements; D9 designs the training. |
| Domain 8 (Automation) | Automation management requires trained operators. D8 defines competency requirements; D9 designs the training. |
| Domain 10 (Operations) | Training delivery is an operational activity. D10 schedules and executes training events; D9 designs them. |
| All Domains (1-10) | Every domain has associated training requirements. D9 is a cross-cutting concern that serves all other domains. |

### Key Tensions

- **Comprehensive Training vs. Operational Time:** Every hour spent training is an hour not spent operating. The institution must balance investment in future capacity against present needs.
- **Standardized vs. Adaptive Curriculum:** A fixed curriculum is reproducible but may not adapt to individual learners or changing conditions.
- **Explicit Knowledge vs. Tacit Knowledge:** Not everything can be written down. The domain must find strategies for preserving experiential, intuitive, and embodied knowledge.
- **Depth vs. Breadth:** Training everyone in everything is impossible. The domain must define minimum cross-training and where specialization is acceptable.
- **Assessment Validity:** How do you know someone actually understands something vs. merely passing a test? In a small institution, the consequences of false certification are severe.

---

## 9.2 ARTICLE LIST

| Article ID | Title | Summary |
|---|---|---|
| D9-001 | Education Philosophy: Reproducing Institutional Competence | Why formal education matters for institutional survival. The distinction between information transfer and competence development. The institution's stance on learning. |
| D9-002 | The Competency Model: What Everyone Must Know | Baseline competencies required of all institutional members regardless of role. Shared knowledge, shared skills, shared understanding. |
| D9-003 | Role-Specific Curriculum Design | Framework for designing training curricula tied to specific institutional roles. How to derive a curriculum from domain documentation. |
| D9-004 | Skill Assessment & Certification Standards | How competency is evaluated. Types of assessment (demonstration, oral examination, practical exercise, portfolio). Certification criteria. Recertification requirements. |
| D9-005 | Mentorship & Apprenticeship Framework | Structured one-on-one knowledge transfer. Mentor responsibilities, apprentice responsibilities, milestone definitions, transition-to-independence criteria. |
| D9-006 | Self-Directed Learning Support | Resources, methods, and structures for learning without a teacher. Study guides, practice exercises, self-assessment tools, learning journals. |
| D9-007 | Tacit Knowledge Capture Strategies | Methods for preserving knowledge that resists documentation: video/audio recording, structured interviews, worked examples, narrated demonstrations, decision logs. |
| D9-008 | Cross-Training Requirements | Minimum cross-training matrix: which roles must be coverable by whom. Single-point-of-failure identification in human capital. |
| D9-009 | Training Material Standards | Format, structure, and quality requirements for all training materials. Consistency with D6 archival standards. Template for instructional documents. |
| D9-010 | Degraded-Capacity Education Plan | How to maintain knowledge transfer when the institution has minimal personnel. Emergency training priorities. "One person left" scenario. |
| D9-011 | New Member Onboarding Sequence | The complete sequence for bringing a new person into the institution, from orientation through role-specific training to independent operation. |
| D9-012 | Training Effectiveness Evaluation | How to assess whether training actually produces competence. Feedback loops, outcome tracking, curriculum revision triggers. |
| D9-013 | The Documentation System as Teacher | How the documentation system itself is designed to be educational. Navigation aids, progressive disclosure, worked examples within articles. |

---

## 9.3 PRIORITY ORDER

**Tier 1 -- Existential (without these, knowledge dies with people):**

1. D9-001 Education Philosophy: Reproducing Institutional Competence
2. D9-002 The Competency Model: What Everyone Must Know
3. D9-008 Cross-Training Requirements
4. D9-010 Degraded-Capacity Education Plan

**Tier 2 -- Structural (the frameworks for formal training):**

5. D9-003 Role-Specific Curriculum Design
6. D9-004 Skill Assessment & Certification Standards
7. D9-005 Mentorship & Apprenticeship Framework
8. D9-009 Training Material Standards

**Tier 3 -- Applied (specific mechanisms):**

9. D9-011 New Member Onboarding Sequence
10. D9-006 Self-Directed Learning Support
11. D9-007 Tacit Knowledge Capture Strategies

**Tier 4 -- Refinement (continuous improvement):**

12. D9-012 Training Effectiveness Evaluation
13. D9-013 The Documentation System as Teacher

### Rationale

The philosophy (D9-001) and baseline competency model (D9-002) must exist before any curriculum can be designed. Cross-training (D9-008) and degraded-capacity planning (D9-010) address the most acute risk: loss of the only person who knows something. Structural articles enable systematic training. Applied articles address specific scenarios. Refinement articles enable continuous improvement but are not critical for initial operation.

---

## 9.4 DEPENDENCIES

### Internal Dependencies (within Domain 9)

```
D9-001 --> D9-002 (philosophy defines what baseline competence means)
D9-001 --> D9-003 (philosophy guides curriculum design approach)
D9-002 --> D9-008 (baseline competencies inform cross-training matrix)
D9-002 --> D9-011 (baseline competencies define onboarding scope)
D9-003 --> D9-004 (curriculum requires assessment)
D9-003 --> D9-005 (curriculum delivered through mentorship)
D9-003 --> D9-009 (curriculum requires training materials)
D9-004 --> D9-012 (assessment feeds effectiveness evaluation)
D9-005 --> D9-007 (mentorship is key channel for tacit knowledge)
D9-006 --> D9-009 (self-directed learning requires well-structured materials)
D9-008 --> D9-010 (cross-training informs degraded-capacity priorities)
D9-009 --> D9-013 (material standards extend to documentation design)
D9-012 --> D9-003 (effectiveness feedback revises curriculum)
```

### Cross-Domain Dependencies

| D9 Article | Depends On | Provides To |
|---|---|---|
| D9-001 | None (foundational) | All domains (every domain has training requirements) |
| D9-002 | All domains (competency model derived from domain requirements) | D10 (operations requires trained personnel) |
| D9-003 | D7 (intelligence competencies), D8 (automation competencies), D10 (operations competencies) | Each source domain (trained people to fill roles) |
| D9-004 | D10 (assessment is an operational activity) | D7-013 (analyst self-assessment), D8-010 (automation governance requires certified personnel) |
| D9-007 | D6-003 (capture formats must be archival-approved) | D6-012 (institutional memory archive includes tacit knowledge captures) |
| D9-009 | D6-003 (training material formats per archival standards) | D6 (training materials stored per D6 policy) |
| D9-010 | D8-013 (degraded capacity may increase automation reliance) | D10 (degraded operations depend on degraded-capacity training) |
| D9-011 | D10 (onboarding is an operational process) | All domains (new members enter through onboarding) |

---

## 9.5 WRITING SCHEDULE

| Phase | Articles | Estimated Effort | Prerequisites |
|---|---|---|---|
| Phase 1 (Weeks 1-3) | D9-001, D9-002 | 2 articles, intensive | None |
| Phase 2 (Weeks 3-6) | D9-008, D9-010 | 2 articles, intensive (existential planning) | D9-001, D9-002 |
| Phase 3 (Weeks 6-9) | D9-003, D9-004, D9-009 | 3 articles, moderate each | D9-001, D9-002 |
| Phase 4 (Weeks 9-12) | D9-005, D9-011 | 2 articles, moderate | D9-003, D9-004 |
| Phase 5 (Weeks 12-14) | D9-006, D9-007 | 2 articles, moderate | D9-003, D9-009 |
| Phase 6 (Weeks 14-16) | D9-012, D9-013 | 2 articles, moderate | D9-004, D9-009 |

**Total estimated duration:** 16 weeks for initial drafts.

**Critical note:** Domain 9 cannot be fully completed until Domains 6, 7, 8, and 10 have at least their foundational articles (Tier 1) drafted, because D9-003 (role-specific curriculum) derives its content from those domains. Initial drafts of D9-003 should use placeholder references to be filled in as source domain articles are completed.

---

## 9.6 REVIEW PROCESS

### Initial Review

Each article undergoes three review passes:

1. **Pedagogical Review.** Does this article reflect sound principles of adult education and competence development? Are the assessment methods valid? Would the training actually produce the claimed competence? Reviewer: someone with teaching or training experience, or failing that, someone who has recently been in the learner role. Timeline: within 2 weeks of draft.

2. **Completeness Review.** Does the training framework cover all critical competencies? Are there gaps that would leave a future maintainer unable to perform a required function? This review requires cross-referencing against all other domains. Reviewer: a generalist with access to all domain documentation. Timeline: within 2 weeks of pedagogical review.

3. **Solo-Operator Review.** Does this article still work if the institution has only one person? Can the frameworks be applied by someone training themselves from documentation alone? Reviewer: someone attempting to learn a skill solely from the documentation. Timeline: within 2 weeks of completeness review.

### Ongoing Review

| Review Type | Frequency | Trigger |
|---|---|---|
| Scheduled review | Every 24 months | Calendar |
| Post-training review | After each training cycle | Training event completion |
| Curriculum update | As needed | Source domain article changes that affect competency requirements |
| Comprehensive audit | Every 60 months | Calendar; full domain coherence review |

---

## 9.7 MAINTENANCE PLAN

### Routine Maintenance

- **After each training event:** Collect learner feedback. Record in Commentary section of relevant articles. Identify curriculum gaps.
- **Annually:** Review cross-training matrix (D9-008) against current personnel. Update gap analysis. Prioritize cross-training for the coming year.
- **Annually:** Review degraded-capacity plan (D9-010) for realism. Is it still feasible given current personnel and systems?
- **Biannually:** Full article review cycle per Section 9.6.

### Triggered Maintenance

- **Source domain article updated:** Review derived curriculum (D9-003) for alignment. Update training materials (D9-009) as needed.
- **Personnel departure:** Immediately assess cross-training coverage (D9-008). Identify exposed single points of failure. Prioritize knowledge capture (D9-007) from departing person if possible.
- **New member joins:** Execute onboarding sequence (D9-011). Use the experience to evaluate and refine the sequence.
- **Assessment failure or competency gap discovered:** Review relevant curriculum, materials, and assessment methods. Determine whether the gap is in training design or delivery.

### Succession Provisions

- The Education Designer role is itself subject to the cross-training matrix (D9-008). At least one other person must be capable of maintaining the education system.
- All training materials must be self-contained enough to function without a live instructor (D9-006 requirement), even if live instruction is preferred.
- D9-010 (Degraded-Capacity Education Plan) is the ultimate fallback. It must be the most accessible, most clearly written article in the entire domain.

---
---

# DOMAIN 10: USER OPERATIONS

**Agent Role:** Operations Lead
**Core Purpose:** Define daily workflows, operational procedures, routine maintenance, incident response, operational tempo, and the human rhythms that keep the institution functioning day after day for decades.

---

## 10.1 DOMAIN MAP

### Scope

Domain 10 governs the moment-to-moment, day-to-day, week-to-week operation of the institution. This is where doctrine meets reality. This includes:

- Daily operational routines and checklists.
- Operational tempo: the rhythm and pace of institutional activity.
- Routine maintenance schedules for all systems.
- Incident detection, classification, response, and post-incident review.
- Shift structures, duty rosters, and workload management (even for one person).
- Operational logs and records.
- Health and sustainability of the operator: fatigue management, morale, work-life boundaries.
- Seasonal and cyclical operational changes.
- Degraded operations: maintaining function with reduced personnel or impaired systems.
- The interface between human activity and automated systems during daily operations.

### Boundaries

This domain does NOT cover:

- The design of systems being operated (Domains 1-8) -- only the procedures for operating and maintaining them.
- The training of operators (Domain 9) -- only the operational context that training must prepare people for.
- Strategic decision-making or long-term planning -- only the execution of established procedures and the escalation of issues that exceed operational authority.
- The design of data systems (Domain 6) -- only the operational procedures for daily data handling.

### Relationships

| Related Domain | Relationship |
|---|---|
| Domain 6 (Data) | Operations generates data (logs, records) governed by D6. Operations executes D6 procedures (backups, integrity checks) on schedule. |
| Domain 7 (Intelligence) | Operations consumes intelligence products for decision-making. Operations provides raw observational data to intelligence. |
| Domain 8 (Automation) | Operations monitors and manages automated systems. Operations intervenes when automation fails or reaches its limits. |
| Domain 9 (Education) | Operations provides the context for training. Operations delivers training events. Operations identifies competency gaps through daily experience. |
| Domains 1-5 | Operations physically maintains and operates all hardware, power, and infrastructure systems. D10 is the "hands on equipment" domain. |

### Key Tensions

- **Thoroughness vs. Sustainability:** Comprehensive daily checklists produce burnout. The operational tempo must be survivable indefinitely.
- **Routine vs. Alertness:** The more routine operations become, the greater the risk of complacency and missed anomalies.
- **Documentation vs. Doing:** Time spent documenting operations is time not spent operating. But undocumented operations cannot be transferred or improved.
- **Single-Operator Reality:** Many self-built institutions are operated by one person. D10 must work for a team of one, scaling up gracefully if that changes.
- **Rigidity vs. Judgment:** Procedures provide consistency, but operators must exercise judgment when conditions deviate from procedure. The domain must define when to follow the book and when to think.

---

## 10.2 ARTICLE LIST

| Article ID | Title | Summary |
|---|---|---|
| D10-001 | Operations Philosophy: The Art of Sustained Function | Foundational principles for daily operations. Why routine matters. Why sustainability matters more than intensity. The operator's mindset. |
| D10-002 | Daily Operations Checklist Framework | How daily checklists are structured, prioritized, and maintained. The master checklist template. Criticality classification. |
| D10-003 | Operational Tempo & Rhythm Design | How to design a sustainable operational tempo: daily, weekly, monthly, quarterly, annual rhythms. Seasonal adjustments. Rest periods. |
| D10-004 | Routine Maintenance Schedules | Master maintenance schedule for all systems across all domains. Frequency, procedures, tools required, acceptance criteria. |
| D10-005 | Incident Classification & Response Framework | How to detect, classify, and respond to incidents. Severity levels, response procedures, escalation paths, post-incident review requirements. |
| D10-006 | Operational Logging Standards | What gets recorded, when, in what format, and where. The operational log as institutional memory. Minimum viable logging for a solo operator. |
| D10-007 | Duty Structures & Workload Management | Shift design for teams, solo-operator protocols, workload limits, mandatory rest, and task prioritization under overload. |
| D10-008 | Access Control: Operational Procedures | Day-to-day procedures for granting, using, and revoking system access. Key management. Authentication procedures. Physical security routines. |
| D10-009 | Automation Monitoring: Operational Procedures | How operators monitor automated systems during daily operations. What to watch, how often, what anomalies look like, intervention thresholds. |
| D10-010 | Degraded Operations Doctrine | How to maintain institutional function when systems are impaired, personnel are reduced, or resources are constrained. Triage framework. Minimum viable operations. |
| D10-011 | Seasonal & Cyclical Operations Guide | How operations change with seasons (power availability, weather, agricultural cycles if applicable). Preparation checklists. Transition procedures. |
| D10-012 | Operator Health & Sustainability | Fatigue management, stress indicators, morale maintenance, work-life boundaries, and the institutional duty of care to its operators. |
| D10-013 | Post-Incident Review Process | Detailed procedures for learning from incidents. Blame-free analysis. Root cause identification. Corrective action tracking. |
| D10-014 | Operational Handoff Procedures | How to transfer operational responsibility between people (shift change, role transition, extended absence). State transfer, briefing requirements. |

---

## 10.3 PRIORITY ORDER

**Tier 1 -- Existential (the institution cannot function without these):**

1. D10-001 Operations Philosophy: The Art of Sustained Function
2. D10-002 Daily Operations Checklist Framework
3. D10-005 Incident Classification & Response Framework
4. D10-010 Degraded Operations Doctrine

**Tier 2 -- Structural (the backbone of routine):**

5. D10-003 Operational Tempo & Rhythm Design
6. D10-004 Routine Maintenance Schedules
7. D10-006 Operational Logging Standards
8. D10-007 Duty Structures & Workload Management

**Tier 3 -- Integration (connecting to other domains):**

9. D10-009 Automation Monitoring: Operational Procedures
10. D10-008 Access Control: Operational Procedures
11. D10-014 Operational Handoff Procedures

**Tier 4 -- Sustainability (long-term human factors):**

12. D10-012 Operator Health & Sustainability
13. D10-011 Seasonal & Cyclical Operations Guide
14. D10-013 Post-Incident Review Process

### Rationale

The philosophy (D10-001) and daily checklist (D10-002) are what keep the institution alive each day. Incident response (D10-005) handles the unexpected. Degraded operations (D10-010) addresses the most likely stressed scenarios. Tempo, maintenance, and logging create sustainable routine. Integration articles connect operations to other domains. Sustainability articles protect the human beings doing the work.

---

## 10.4 DEPENDENCIES

### Internal Dependencies (within Domain 10)

```
D10-001 --> D10-002 (philosophy shapes daily routine design)
D10-001 --> D10-003 (philosophy defines sustainable tempo)
D10-001 --> D10-012 (philosophy demands operator care)
D10-002 --> D10-004 (daily checklists include maintenance tasks)
D10-002 --> D10-006 (daily checklists include logging requirements)
D10-003 --> D10-007 (tempo design informs duty structures)
D10-003 --> D10-011 (tempo changes seasonally)
D10-005 --> D10-013 (incidents trigger post-incident review)
D10-005 --> D10-010 (some incidents cause degraded operations)
D10-007 --> D10-014 (duty structures require handoff procedures)
D10-007 --> D10-012 (workload management relates to health)
D10-009 --> D10-005 (automation anomalies may be incidents)
D10-010 --> D10-007 (degraded ops redefine workload)
```

### Cross-Domain Dependencies

| D10 Article | Depends On | Provides To |
|---|---|---|
| D10-001 | All domain philosophies (operations must serve institutional purpose) | All domains (operations executes what other domains design) |
| D10-002 | All domains (checklists derived from each domain's maintenance and monitoring requirements) | D6 (daily data handling), D7 (daily intelligence routines) |
| D10-004 | Domains 1-5 (hardware maintenance specs), D6-005 (media maintenance), D8-006 (agent lifecycle maintenance) | All domains (maintenance is the operational realization of each domain's requirements) |
| D10-005 | D8-009 (automation override is an incident response), D6-013 (data disaster recovery) | D7-012 (intelligence failures are a class of incident) |
| D10-006 | D6-008 (logs are data artifacts with metadata standards) | D6 (logs stored per D6 policy), D7 (logs are raw intelligence data) |
| D10-007 | D9-008 (cross-training determines who can cover which duties) | D9 (operational experience identifies training gaps) |
| D10-009 | D8-007 (automation monitoring design) | D8 (operational feedback on automation behavior) |
| D10-010 | D9-010 (degraded-capacity training) | D9 (degraded ops identifies critical training gaps) |
| D10-012 | D9-001 (education philosophy includes wellness) | D9 (operator health training requirements) |
| D10-014 | D9-011 (onboarding is a form of handoff) | D9 (handoff procedures require trained personnel) |

---

## 10.5 WRITING SCHEDULE

| Phase | Articles | Estimated Effort | Prerequisites |
|---|---|---|---|
| Phase 1 (Weeks 1-3) | D10-001, D10-005 | 2 articles, intensive | None (foundational) |
| Phase 2 (Weeks 3-5) | D10-002, D10-010 | 2 articles, intensive | D10-001 |
| Phase 3 (Weeks 5-8) | D10-003, D10-004, D10-006 | 3 articles, moderate each | D10-001, D10-002 |
| Phase 4 (Weeks 8-10) | D10-007, D10-009 | 2 articles, moderate | D10-003, D8-007 |
| Phase 5 (Weeks 10-13) | D10-008, D10-014, D10-013 | 3 articles, moderate | D10-005, D10-007 |
| Phase 6 (Weeks 13-16) | D10-012, D10-011 | 2 articles, moderate | D10-003, D10-007 |

**Total estimated duration:** 16 weeks for initial drafts.

**Critical note:** D10-002 (Daily Operations Checklist) and D10-004 (Routine Maintenance Schedules) are living documents that will grow substantially as other domains produce their articles. Initial drafts should establish the framework; content will be populated incrementally as source domains mature.

---

## 10.6 REVIEW PROCESS

### Initial Review

Each article undergoes three review passes:

1. **Operational Reality Review.** Has this article been tested against actual daily operations? Do the procedures work in practice, not just in theory? Are the time estimates realistic? Reviewer: the person who will actually perform the procedures, ideally after a trial period. Timeline: within 2 weeks of draft, after at least 1 week of operational testing.

2. **Sustainability Review.** Can these procedures be sustained indefinitely without burnout? Is the operational tempo survivable for a solo operator? For a small team? Reviewer: someone focused on human factors and long-term viability. Timeline: within 1 week of operational review.

3. **Degradation Review.** What happens to these procedures when conditions deteriorate? Does D10-010 adequately cover the failure modes of normal operations? Reviewer: someone with experience in crisis management or austere conditions. Timeline: within 1 week of sustainability review.

### Ongoing Review

| Review Type | Frequency | Trigger |
|---|---|---|
| Scheduled review | Every 12 months (shortest cycle -- operations evolves fastest) | Calendar |
| Post-incident review | After every significant incident | Incident closure |
| Seasonal review | Before each major seasonal transition | Calendar |
| Comprehensive audit | Every 36 months | Calendar; full domain coherence review |

### Review Artifacts

Same as other domains: dated commentary note, updated review date, change list or no-change justification, version increment.

**Additional D10-specific artifact:** Operational Review Log, a running record of all changes to operational procedures with brief rationale, maintained as a D6 data artifact.

---

## 10.7 MAINTENANCE PLAN

### Routine Maintenance

- **Weekly:** Quick assessment: are daily procedures still current? Any procedures that no longer apply or need adjustment?
- **Monthly:** Review operational logs for patterns. Are the same issues recurring? Do procedures need updating?
- **Quarterly:** Review automation monitoring procedures (D10-009) against actual automation behavior. Update thresholds and anomaly definitions.
- **Semiannually:** Review operational tempo (D10-003) against operator feedback. Is the pace sustainable? Adjust.
- **Annually:** Full article review cycle per Section 10.6. Comprehensive maintenance schedule review (D10-004).

### Triggered Maintenance

- **Incident occurs:** Follow D10-005 response framework. After resolution, execute D10-013 post-incident review. Update implicated procedures.
- **System added or removed:** Update D10-002 (daily checklist), D10-004 (maintenance schedule), D10-009 (automation monitoring).
- **Personnel change:** Update D10-007 (duty structures), execute D10-014 (handoff). If personnel decrease, invoke D10-010 (degraded operations assessment).
- **Seasonal transition:** Execute D10-011 preparation checklists. Adjust tempo per D10-003.

### Succession Provisions

- D10-002 (Daily Operations Checklist) is the single most critical operational document. It must be understandable by a newcomer on their first day, even if they cannot execute everything on it immediately.
- D10-014 (Operational Handoff) must be comprehensive enough to transfer operational state from one person to another with minimal loss, including during unplanned transitions (incapacitation of current operator).
- All operational procedures must be written in clear, imperative language suitable for someone under stress with limited context.
- A "cold start" procedure must be derivable from D10 articles alone: if the institution has been non-operational and a new person arrives, D10-002 and D10-010 together should enable them to begin operating, however minimally.

---
---

# CROSS-DOMAIN SYNTHESIS: DOMAINS 6-10

## Integration Points

The five domains in this document form a tightly coupled operational cluster. Their integration points are:

### The Data-Operations Axis (D6 <-> D10)
Operations generates data; data systems require operational maintenance. The daily checklist (D10-002) includes D6 backup and integrity procedures. Operational logs (D10-006) are D6 data artifacts. This is the most active integration point -- it fires every day.

### The Intelligence-Operations Axis (D7 <-> D10)
Intelligence products inform operational decisions. Operational observations feed intelligence collection. The intelligence cycle (D7-002) and operational tempo (D10-003) must be synchronized so that intelligence products are available when operational decisions are made.

### The Automation-Operations Axis (D8 <-> D10)
Automated systems run within the operational environment. Operations monitors them (D10-009) per automation monitoring design (D8-007). When automation fails, operations invokes override procedures (D8-009) as an incident response (D10-005). This axis is the primary human-machine interface.

### The Education-Everything Axis (D9 <-> All)
Education is the cross-cutting concern. Every domain depends on trained people. Every domain's evolution depends on the next generation understanding the current design. D9 is the domain that ensures all other domains survive personnel changes.

### The Degradation Cascade (D10-010, D9-010, D8-004)
When operations degrade (D10-010), the institution may need to rely more heavily on automation (D8-004 limits) and may have reduced capacity for training (D9-010). These three articles must be read together and must be internally consistent.

## Master Dependency Graph (Simplified)

```
D6-001 (Data Philosophy)
  |
  +--> D6-002 (Data Model) --> D8-012 (Automation Data Contracts)
  |                         --> D6-008 (Metadata) --> D10-006 (Logging)
  |
  +--> D6-003 (Format Longevity) --> D9-009 (Training Material Standards)
  |                               --> D6-009 (Migration)
  |
  +--> D6-006 (Backup) --> D10-002 (Daily Checklist, includes backup tasks)
  |
  +--> D6-007 (Integrity) --> D8-012 (Agents must preserve integrity)

D7-001 (Intelligence Philosophy)
  |
  +--> D7-002 (Intel Cycle) --> D10-003 (Tempo, synchronized with intel)
  |
  +--> D7-005 (Bias) --> D7-013 (Single-Analyst) --> D9-003 (Analyst Training)
  |
  +--> D7-003 (Collection) --> D6-014 (Data Ingest)
  |                         --> D8 (Automated collection agents)

D8-001 (Automation Philosophy)
  |
  +--> D8-003 (Human-in-Loop) --> D10-009 (Ops monitors automation)
  |
  +--> D8-004 (Autonomy Limits) --> D8-009 (Override) --> D10-005 (Incident Response)
  |
  +--> D8-002 (Agent Design) --> D8-007 (Monitoring) --> D10-009 (Ops monitoring)

D9-001 (Education Philosophy)
  |
  +--> D9-002 (Competency Model) --> D9-008 (Cross-Training) --> D10-007 (Duty Structures)
  |
  +--> D9-003 (Curriculum) --> [All domains' competency requirements]
  |
  +--> D9-010 (Degraded Capacity) --> D10-010 (Degraded Operations)

D10-001 (Operations Philosophy)
  |
  +--> D10-002 (Daily Checklist) --> [All domains' maintenance requirements]
  |
  +--> D10-005 (Incident Response) --> D10-013 (Post-Incident Review)
  |
  +--> D10-003 (Tempo) --> D10-007 (Duty) --> D10-012 (Operator Health)
```

## Recommended Cross-Domain Writing Coordination

Because these five domains are tightly coupled, the following coordination rules apply:

1. **All five domain philosophies (D6-001, D7-001, D8-001, D9-001, D10-001) should be drafted in the same phase.** They must be mutually consistent. Cross-review all five before any is finalized.

2. **D10-002 (Daily Checklist) cannot be completed until at least the Tier 1 articles of all other domains exist.** It is a derivative document. Plan for iterative expansion.

3. **D9-003 (Role-Specific Curriculum) requires content from all other domains.** It should be started early as a framework but filled in progressively.

4. **The degradation triad (D10-010, D9-010, D8-004) must be reviewed as a unit.** A change to one may require changes to the others.

5. **All cross-domain dependency links should be implemented as explicit, typed references** using the article ID system, not vague statements like "see the data section."

---

## GLOBAL WRITING SCHEDULE (DOMAINS 6-10 COMBINED)

| Combined Phase | Weeks | Focus | Articles |
|---|---|---|---|
| 1. Philosophies | 1-3 | All five domain philosophy articles drafted and cross-reviewed | D6-001, D7-001, D8-001, D9-001, D10-001 |
| 2. Existential Safeguards | 3-6 | Articles that prevent institutional death | D6-003, D6-006, D7-005, D8-003, D9-002, D10-005 |
| 3. Structural Frameworks | 6-10 | The architecture of each domain | D6-002, D6-004, D7-002, D7-004, D8-002, D8-004, D9-008, D10-002, D10-010 |
| 4. Operational Backbone | 10-14 | What makes daily operations possible | D6-007, D6-008, D7-003, D7-011, D8-005, D8-006, D9-003, D9-004, D10-003, D10-004, D10-006 |
| 5. Integration & Monitoring | 14-18 | Connecting domains, building monitoring | D6-014, D7-008, D7-009, D8-007, D8-012, D9-005, D9-009, D10-007, D10-008, D10-009 |
| 6. Resilience & Sustainability | 18-22 | Long-term survival, human factors, edge cases | D6-009, D6-013, D7-006, D7-007, D8-008, D8-009, D9-010, D9-011, D10-012, D10-014 |
| 7. Refinement & Legacy | 22-26 | Continuous improvement, legacy, remaining articles | D6-005, D6-010, D6-011, D6-012, D6-015, D7-010, D7-012, D7-013, D8-010, D8-011, D8-013, D9-006, D9-007, D9-012, D9-013, D10-011, D10-013 |

**Total estimated duration for Domains 6-10 combined:** 26 weeks (6 months) for initial drafts of all 68 articles.

---

## GLOBAL REVIEW CADENCE

| Review Event | Frequency | Scope |
|---|---|---|
| Individual article review | Per domain schedule (12-24 months) | Single article, per domain review process |
| Cross-domain consistency review | Every 18 months | All five domain philosophy articles reviewed together |
| Degradation triad review | Every 12 months | D10-010 + D9-010 + D8-004 reviewed as unit |
| Full documentation audit | Every 60 months | All 68 articles reviewed for coherence, currency, and completeness |
| Succession readiness assessment | Every 12 months | Can a new person bootstrap from documentation alone? Test this. |

---

## FINAL NOTES

### On the Nature of This Document

This Stage 1 framework is itself an artifact that must be maintained. As articles are written, this framework should be updated to reflect:
- Articles that were merged, split, renamed, or abandoned during writing.
- Dependencies that were discovered or dissolved during writing.
- Schedule adjustments driven by experience.

### On Getting Started

The single most important action is writing the five philosophy articles (D6-001, D7-001, D8-001, D9-001, D10-001). Everything else flows from these. If time and energy are limited, write the philosophies. If only one philosophy can be written, write D10-001 (Operations Philosophy), because operations is what keeps the institution alive today while the rest of the documentation system is built.

### On Imperfection

No documentation framework survives contact with reality unchanged. This framework is a starting point, not a contract. The meta-rules at the top of this document and the standard article template are the durable structures. The specific article list, priorities, and schedules are provisional. Adjust them. Document why you adjusted them. That documentation is itself an institutional memory artifact.

---

**END OF STAGE 1 FRAMEWORK: DOMAINS 6-10**
**Document Version:** 1.0.0
**Total Proposed Articles:** 68
**Total Estimated Initial Writing Duration:** 26 weeks
**Next Action:** Begin Phase 1 -- draft five domain philosophy articles in parallel.
