# STAGE 3: OPERATIONAL DOCTRINE -- EDUCATION & TRAINING

## Domain 9 Operations Articles D9-003 through D9-007

**Document ID:** STAGE3-EDUCATION-OPS
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Classification:** Stage 3 -- Operational Doctrine. Practical manuals for assessing competence, designing curricula, capturing tacit knowledge, cross-training across functions, and writing documentation that teaches.
**Depends On:** ETH-001, CON-001, GOV-001, SEC-001, OPS-001, D9-001, D9-002
**Depended Upon By:** All subsequent Domain 9 articles. Referenced by all domains requiring competence verification, training design, or knowledge preservation.

---

## How to Read This Document

This document contains five operational articles for Domain 9: Education & Training. They are numbered D9-003 through D9-007, following the root education philosophy article D9-001 and the competency model article D9-002 established in Stage 2 and early Stage 3.

Where D9-001 established why the institution invests in education -- the reproduction of institutional competence across time and across people -- and D9-002 defined what competencies the institution requires, these five articles define how. How to assess whether someone has a skill. How to build a curriculum that teaches it. How to extract knowledge from the head of a person who cannot easily articulate it. How to ensure no critical skill lives in only one person. How to write documentation that functions as a teacher when no teacher is present.

These are operational documents. They contain procedures, templates, matrices, and rubrics. They are designed for immediate use. Each article can be consulted independently, though they reference one another extensively. If you are building a training program from scratch, read them in order. If you need a specific procedure, navigate directly.

Every procedure in this document derives its authority from D9-001 and, through it, from the Core Charter. The self-teaching requirement (R-D9-01), the documentation-first principle (R-D9-02), the tacit knowledge capture mandate (R-D9-03), the cross-training requirement (R-D9-04), and the format longevity constraint (R-D9-05) are the governing rules. No procedure here may contradict those rules. Where ambiguity arises, D9-001 prevails.

If you are a future operator reading this for the first time: these articles were written when the institution had one person. They are designed to scale from one to several, but they begin with one. Do not skip the single-operator procedures because you are also a single operator. The discipline of formal assessment, curriculum design, and knowledge capture applies even when you are assessing yourself, teaching yourself, and capturing your own knowledge. Especially then.

---

---

# D9-003 -- Competency Assessment Framework

**Document ID:** D9-003
**Domain:** 9 -- Education & Training
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, OPS-001, D9-001, D9-002
**Depended Upon By:** D9-004, D9-005, D9-006, D9-007, D9-011. Referenced by all articles requiring competence verification.

---

## 1. Purpose

This article defines how the holm.chat Documentation Institution assesses whether a person possesses the skills and knowledge required to perform institutional functions. Assessment is the feedback mechanism that closes the loop between education and competence. Without it, the institution cannot know whether its training works, whether its operators are prepared, whether a succession candidate is ready, or whether a skill gap is forming before it becomes a crisis.

D9-001 established a four-stage competence progression: Awareness, Guided Execution, Independent Execution, and Adaptive Mastery. D9-002 defined the specific competencies required for each institutional role. This article provides the operational machinery for determining where a person stands on that progression for any given competency, and for doing so with enough rigor that the result is trustworthy.

Assessment in a single-operator institution raises an immediate objection: "I am assessing myself. How can that be honest?" The answer is the same one that SEC-001 gives for access controls: self-discipline is not self-delusion. The procedures in this article are designed to make self-assessment rigorous by making it structured, documented, and verifiable against objective criteria. A self-assessment that follows the rubric in this article is more trustworthy than an unstructured guess, even though neither is as reliable as external assessment. When external assessors are available -- during succession, apprenticeship, or advisory relationships -- the same procedures apply with the addition of independent verification.

## 2. Scope

This article covers:

- The skill taxonomy: the complete classification of skills across all institutional domains.
- Assessment methods: practical tests, knowledge checks, scenario exercises, and observation protocols.
- Scoring rubrics: how to convert observed performance into competence stage ratings.
- Progression criteria: what specific evidence is required to certify a person at each competence stage.
- Self-assessment procedures: how a single operator assesses their own competence honestly.
- Assessment scheduling: when assessments occur and what triggers unscheduled assessment.
- Assessment records: how results are documented and stored.

This article does not cover:

- The definition of institutional competencies (see D9-002).
- Curriculum design for addressing identified gaps (see D9-004).
- Tacit knowledge assessment, which requires specialized capture techniques (see D9-005).
- Training material quality assessment (see D9-007).

## 3. Background

### 3.1 The Problem of Honest Self-Assessment

Research in cognitive science consistently shows that people are poor judges of their own competence. The Dunning-Kruger effect is the most widely cited finding: people with low competence in a domain tend to overestimate their ability, while experts tend to underestimate theirs. In a single-operator institution, this is not an academic concern. If the operator overestimates their ability to manage backups, the institution will not discover the gap until a restore fails. If the operator underestimates their documentation skill, they may invest training time in an area where they are already sufficient while neglecting a genuine weakness.

The mitigation is not to abandon self-assessment -- there may be no alternative -- but to structure it so heavily that subjective judgment is constrained by objective criteria. Every competence stage has specific, observable indicators. Every assessment method produces an artifact -- a completed test, a documented scenario walkthrough, a verified procedure execution -- that can be reviewed later. The operator does not ask, "Am I good at this?" The operator asks, "Can I produce the specific evidence that the rubric requires for this competence stage?" The question shifts from self-perception to self-demonstration.

### 3.2 Assessment as Institutional Memory

Assessment records serve a second function beyond competence verification: they are a form of institutional memory. A complete assessment history tells a future operator or auditor what the institution's competence profile looked like at any point in time. It reveals trends -- skills that degraded, skills that improved, skills that were never adequately developed. During a succession event, the incoming operator can review the outgoing operator's assessment history to understand what the institution's strengths and vulnerabilities are. This is why assessment records are classified as Tier 1 data per D6-001.

### 3.3 The Skill Taxonomy

Skills in this institution span ten domains, but they cluster into four categories that cut across domain boundaries:

**Category 1: Technical Operations.** Skills involving direct interaction with hardware and software. Examples: configuring backup systems, managing storage, maintaining the operating system, performing data migrations. These skills are assessed primarily through practical demonstration.

**Category 2: Information Management.** Skills involving the creation, organization, preservation, and retrieval of information. Examples: writing documentation, classifying data, applying metadata, managing the documentation corpus. These skills are assessed through a combination of practical demonstration and artifact review.

**Category 3: Governance and Process.** Skills involving decision-making, planning, and institutional management. Examples: applying the amendment process, conducting dispute resolution, performing risk assessments, managing the complexity budget. These skills are assessed through scenario exercises and decision log review.

**Category 4: Meta-Skills.** Skills that span all categories: the ability to learn from documentation, the ability to teach others, the ability to recognize when something is wrong, the ability to adapt procedures to changed circumstances. These skills are assessed through observation, scenario exercises, and longitudinal review.

## 4. System Model

### 4.1 The Assessment Methods

The institution uses four assessment methods, each suited to different skill types and competence stages:

**Method 1: Practical Test.** The person performs a real or simulated task while the outcome is documented. This is the primary method for Technical Operations skills. A practical test has four components: a task specification (what must be done), success criteria (what a correct outcome looks like), a time frame (how long the task should reasonably take), and a verification procedure (how the outcome is checked). For self-assessment, the operator performs the task and then verifies the outcome against the success criteria. For external assessment, the assessor observes or reviews the outcome independently.

Example: "Perform a full backup and restore cycle. Success: the restored data is byte-identical to the original, verified by checksum comparison. Time frame: within the time limits specified in D6-002. Verification: compare checksums of original and restored data sets."

**Method 2: Knowledge Check.** The person answers questions or explains concepts without performing a task. This is the primary method for assessing Awareness (Stage 1) and the conceptual understanding component of higher stages. Knowledge checks can be written or verbal. They must test understanding, not memorization. Questions should require the person to explain why, not just what. A question like "What is the backup rotation schedule?" tests memorization. A question like "Why does the institution use a grandfather-father-son backup rotation instead of simple chronological rotation, and under what circumstances would you deviate from it?" tests understanding.

Knowledge checks are documented as question-and-answer pairs. For self-assessment, the operator writes their answer before consulting the documentation, then compares. The gap between what the operator wrote and what the documentation says is the assessment finding.

**Method 3: Scenario Exercise.** The person is presented with a hypothetical situation and must describe or demonstrate how they would respond. This is the primary method for Governance and Process skills and for assessing Independent Execution (Stage 3) and Adaptive Mastery (Stage 4). Scenarios should be realistic and should include complications that force judgment calls.

Example: "During a routine backup, you discover that the backup media shows read errors on verification. The previous two backup sets are also on media from the same batch. Describe your response, including which articles you would consult, what immediate actions you would take, and how you would communicate the situation in the operational log."

Scenario exercises are documented as the scenario prompt and the person's complete response. For self-assessment, the operator writes their response, then reviews it against the relevant procedures, identifying any steps they missed or misjudged.

**Method 4: Observation Protocol.** An assessor watches the person perform their normal duties over a defined period, noting competence indicators. This method is only available when two or more people are present. It is the most valuable assessment method because it reveals habitual behavior rather than test-day performance. Observation protocols define what to look for (specific behaviors, decision patterns, procedural adherence), how long to observe (minimum one full operational cycle for daily tasks), and how to record observations (structured observation notes, not impressionistic summaries).

### 4.2 The Scoring Rubric

Each competency is scored against the four-stage progression defined in D9-001. The rubric defines what specific evidence places a person at each stage:

**Stage 1 -- Awareness.** Evidence required: the person can identify the function, state its purpose in the institution, name the primary documentation source for it, and recognize when it needs to be performed. Evidence is typically gathered through Knowledge Check (Method 2). Threshold: the person answers correctly for at least 80% of awareness-level questions about the competency.

**Stage 2 -- Guided Execution.** Evidence required: the person can perform the function by following documented procedures, with documentation open and available, achieving a correct outcome. Evidence is gathered through Practical Test (Method 1) with documentation access permitted. Threshold: the person achieves a correct outcome on 100% of critical steps and at least 90% of all steps. Any failed critical step is an automatic failure regardless of overall percentage.

**Stage 3 -- Independent Execution.** Evidence required: the person can perform the function without consulting documentation for routine operations, can explain the reasoning behind each major step, can recognize when conditions deviate from normal, and can handle at least one non-trivial scenario exercise correctly. Evidence is gathered through Practical Test (Method 1) without documentation access, Knowledge Check (Method 2) focused on reasoning, and Scenario Exercise (Method 3). Threshold: correct outcome on practical test, satisfactory explanation of reasoning (assessed against the Background and System Model sections of the relevant article), and a response to the scenario exercise that identifies the correct course of action even if minor details are imperfect.

**Stage 4 -- Adaptive Mastery.** Evidence required: the person can modify procedures when conditions require it, can diagnose novel problems, can teach the function to others, and can identify improvements to existing procedures. Evidence is gathered through advanced Scenario Exercise (Method 3) involving novel situations not covered by existing documentation, Observation Protocol (Method 4) over an extended period, and review of the person's contributions to documentation, Commentary Sections, and procedure improvements. Threshold: Stage 4 is assessed holistically rather than by checklist. The assessor (or the self-assessing operator) must document specific instances demonstrating adaptive capability.

### 4.3 The Assessment Schedule

**Annual Comprehensive Assessment.** Once per year, the operator conducts a full assessment of all competencies listed in D9-002. This assessment covers every competency at the currently claimed stage. Results are recorded in the assessment log. Gaps are documented and fed into the curriculum review process per D9-004.

**Triggered Assessment.** An assessment is triggered by any of the following events: a competency is required for a task and the operator has not been assessed on it within the past twelve months; a procedure fails and the root cause may be a competency gap; a new person assumes any institutional role; a succession event is initiated; or the quarterly review identifies a potential skill degradation.

**Post-Training Assessment.** After completing any training activity (self-study, curriculum module, cross-training rotation), the relevant competency is assessed to verify that the training achieved its objective.

### 4.4 Assessment Records

Every assessment produces a record containing: the date of the assessment; the competency assessed; the assessment method used; the evidence produced (test results, written answers, scenario responses); the competence stage assigned; the assessor's identity (including "self" for self-assessment); and any notes about gaps, near-misses, or areas for further development. Assessment records are stored in the documentation corpus, classified as Tier 1 data per D6-001, and are included in the annual operational review.

## 5. Rules & Constraints

- **R-D9-03-01:** Every institutional competency defined in D9-002 must be assessed at least annually for every person holding an operational role. There are no exemptions for experience or tenure.
- **R-D9-03-02:** Self-assessment must follow the structured methods defined in this article. Unstructured self-assessment ("I think I am good at this") does not count and must not be recorded as a valid assessment.
- **R-D9-03-03:** Stage 2 assessment for any competency must include a practical demonstration with documentation available. A knowledge check alone is insufficient for Stage 2 certification.
- **R-D9-03-04:** Stage 3 assessment must include a practical demonstration without routine documentation access and at least one scenario exercise. Being able to follow a procedure is not the same as understanding it.
- **R-D9-03-05:** Critical competencies -- those required for data preservation, security, and institutional continuity -- must be assessed at the claimed stage at least every six months. The list of critical competencies is maintained in D9-002 and reviewed annually.
- **R-D9-03-06:** Assessment results that reveal a gap between claimed competence and demonstrated competence must trigger a curriculum response within thirty days per D9-004.
- **R-D9-03-07:** Assessment records must be preserved for the lifetime of the institution. They are never deleted or overwritten. Corrections are appended, not substituted.

## 6. Failure Modes

- **Self-assessment inflation.** The operator consistently rates themselves at a higher competence stage than their actual performance warrants. Over time, the assessment records paint a picture of universal competence that does not match reality. Mitigation: the scoring rubric requires specific evidence, not self-rating. The question is not "How good am I?" but "Can I produce this evidence?" Additionally, post-incident reviews (D10-003) should cross-reference the operator's claimed competence in the relevant area with their actual performance during the incident.

- **Assessment avoidance.** The operator skips assessments because they are time-consuming and the results are uncomfortable. The annual assessment becomes biannual, then occasional, then forgotten. Mitigation: R-D9-03-01 is non-negotiable. The annual assessment is integrated into the operational calendar per OPS-001 and is tracked in the quarterly governance health check per GOV-007.

- **Teaching to the test.** The operator memorizes the specific practical tests and scenario exercises rather than developing genuine competence. They can pass the assessment but cannot handle real situations that deviate from the test script. Mitigation: scenario exercises should be varied between assessment cycles. The assessor (even if self) should create new scenarios that test the same competencies in different contexts. The assessment record should note whether the same scenario was reused.

- **Assessment without consequences.** Assessments are conducted, gaps are identified, but nothing changes. The gap is recorded and then ignored. Mitigation: R-D9-03-06 requires a curriculum response within thirty days. The quarterly review checks whether identified gaps have been addressed.

- **Over-assessment paralysis.** The operator becomes so focused on assessment that it consumes time needed for actual operations. Assessment becomes an end in itself. Mitigation: the annual comprehensive assessment should take no more than two full operational days. Triggered and post-training assessments should take no more than two hours each. If assessment is consuming more time than this, the assessment procedures need simplification, not the operator's schedule needs more time.

## 7. Recovery Procedures

1. **If assessments have not been conducted for more than twelve months:** Do not attempt to conduct a full comprehensive assessment immediately. Begin with the critical competencies identified in D9-002. Assess those within two weeks. Then schedule the remaining competencies over the following sixty days. Record the assessment gap in the operational log and the decision log as a governance lapse.

2. **If assessment records have been lost:** Conduct a fresh comprehensive assessment. Record the loss of historical records as a data incident per D6-002 procedures. Note in the new assessment records that historical data is unavailable. Do not fabricate retrospective records.

3. **If a competency gap is discovered during a live incident:** Handle the incident first, using documentation to compensate for the gap. After the incident is resolved, conduct a triggered assessment of the relevant competency. Record the incident-discovered gap in both the incident log and the assessment log. Initiate a curriculum response per D9-004.

4. **If self-assessment appears unreliable:** This is the hardest recovery. If you suspect that your self-assessments are consistently inaccurate, switch to evidence-only assessment: record yourself performing practical tests (video or detailed log), then review the recording against the rubric as though you were assessing someone else. If an external assessor becomes available, prioritize external reassessment of the competencies you are most uncertain about.

## 8. Evolution Path

- **Years 0-5:** Self-assessment is the primary mode. The founder establishes the habit of structured assessment and begins building the assessment record. The rubrics are tested against real performance and refined. Expect to discover that some rubrics are too easy (passing them does not indicate real competence) and others are too strict (failing them does not indicate real deficiency).

- **Years 5-15:** If a second person becomes involved, external assessment becomes possible. This is a critical milestone. The first external assessment will likely reveal gaps that self-assessment missed. Use this as an opportunity to calibrate the rubrics, not as an indictment of the self-assessment process.

- **Years 15-30:** The assessment record now spans a significant period. Longitudinal trends become visible. The institution can identify competencies that consistently degrade over time (candidates for more frequent assessment) and competencies that remain stable (candidates for less frequent assessment).

- **Years 30-50+:** Assessment procedures should be well-calibrated by now. The primary concern shifts from establishing the system to ensuring it remains meaningful -- that assessments are genuinely testing competence rather than merely perpetuating a bureaucratic ritual.

- **Signpost for revision:** If a competency gap causes an incident that assessment records show was assessed at Stage 3 or above, the assessment method for that competency is broken and must be redesigned.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
I am deeply skeptical of self-assessment, and I am building a system that relies on it. That tension is intentional. The alternative to structured self-assessment is not external assessment -- there is no external assessor available. The alternative is no assessment at all, which means no feedback, no gap identification, and no mechanism for improvement.

The rubric is the key. A good rubric turns self-assessment from a subjective judgment ("I think I know this") into an evidence-gathering exercise ("Can I produce this specific result?"). I have tried to write rubrics that are hard to game. The practical tests require actual outcomes, not claims. The scenario exercises require written responses that can be reviewed against the documentation. The knowledge checks test reasoning, not recall.

I expect the first annual assessment to be humbling. I expect to discover competencies I thought I had but cannot demonstrate. That is the entire point. The assessment that reveals a gap is infinitely more valuable than the assessment that confirms what I want to believe.

## 10. References

- ETH-001 -- Ethical Foundations of the Institution (Principle 2: Integrity Over Convenience; Principle 3: Transparency)
- CON-001 -- The Founding Mandate (institutional transferability, documentation completeness)
- GOV-001 -- Authority Model (succession preparedness, decision documentation)
- OPS-001 -- Operations Philosophy (operational calendar, documentation-first principle)
- D6-001 -- Data Philosophy (Tier 1 classification for assessment records)
- D9-001 -- Education Philosophy (four-stage competence progression, self-teaching requirement)
- D9-002 -- The Competency Model (competency definitions, role requirements, critical competency list)
- D9-004 -- Curriculum Design and Maintenance (curriculum response to identified gaps)
- D9-005 -- Tacit Knowledge Capture Procedures (specialized assessment for tacit knowledge)
- D9-007 -- Documentation-as-Teacher (documentation quality as enabler of self-study assessment)
- D10-003 -- Incident Response Procedures (post-incident competency review)
- GOV-007 -- Governance Health Checks (quarterly review of assessment compliance)

---

---

# D9-004 -- Curriculum Design and Maintenance

**Document ID:** D9-004
**Domain:** 9 -- Education & Training
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, OPS-001, D9-001, D9-002, D9-003
**Depended Upon By:** D9-005, D9-006, D9-007, D9-011. Referenced by all articles requiring structured learning sequences.

---

## 1. Purpose

This article defines how the holm.chat Documentation Institution designs, builds, and maintains training curricula. A curriculum is more than a list of things to learn. It is a deliberate sequence: what to learn first, what to learn next, what prerequisites must be satisfied before advancing, and how to verify that learning objectives have been met at each stage. Without a curriculum, education is ad hoc -- the learner reads whatever seems relevant, skips whatever seems boring, and develops a patchwork of knowledge with invisible gaps.

D9-001 established that documentation is the primary teacher. D9-002 defined what must be taught. D9-003 defined how to assess whether teaching worked. This article defines the structure that connects them: the curriculum that transforms a corpus of documentation into a coherent learning path.

The founding challenge is that the institution's documentation corpus was not written as a textbook. It was written as a reference library -- organized by domain and function, not by learning sequence. A new operator confronting the documentation for the first time faces hundreds of articles across ten domains with no obvious starting point and no clear order. The curriculum solves this problem by providing the map: start here, go there, read this before that, stop and practice, then continue.

## 2. Scope

This article covers:

- Learning objective design: how to define what a curriculum module must teach.
- Sequencing: how to order learning activities for effective knowledge building.
- Prerequisite mapping: how to identify and document what must be known before each module.
- Creating lessons from existing documentation articles.
- Curriculum structure: modules, units, and learning paths.
- The curriculum review schedule: when and how curricula are updated.
- Curriculum maintenance: keeping curricula aligned with evolving documentation and institutional needs.

This article does not cover:

- Competency definitions (see D9-002).
- Assessment design (see D9-003; this article references assessment but does not define assessment methods).
- Tacit knowledge curricula (see D9-005 for capture; tacit knowledge is integrated into curricula per Section 4.4 of this article).
- Training material format standards (governed by D9-001 R-D9-05 and D6-003).

## 3. Background

### 3.1 The Documentation-to-Curriculum Translation Problem

The institution's documentation is comprehensive by design. Every article follows the canonical template. Every article contains Purpose, Background, System Model, Rules, Failure Modes, Recovery, and Evolution sections. This structure makes each article self-contained and reference-quality. But self-contained reference articles are not optimized for learning. A learner who reads SEC-003 (Cryptographic Key Management) without first understanding SEC-001 (Security Philosophy) and SEC-002 (Access Control Procedures) will absorb mechanics without context. They will know what to do but not why, which D9-001 identifies as the most dangerous knowledge gap.

A curriculum solves this by imposing a pedagogical layer on top of the documentation corpus. It does not replace the documentation. It wraps it -- adding reading order, pre-reading context, post-reading exercises, and connecting narrative between articles that were written independently.

### 3.2 Prerequisites as Safety Mechanisms

In this institution, prerequisite mapping is not merely pedagogical best practice. It is a safety mechanism. Some functions, if performed by someone who lacks prerequisite knowledge, can cause institutional damage. An operator who attempts cryptographic key rotation without understanding the key inventory, or who performs a major system update without understanding rollback procedures, risks data loss. The curriculum's prerequisite chain is a guardrail: it ensures that by the time a learner reaches a dangerous procedure, they have the conceptual foundation to execute it safely.

### 3.3 Curriculum in a Single-Operator Context

When the operator is both the curriculum designer and the sole student, curriculum design serves three purposes. First, it forces the designer to articulate the learning sequence explicitly rather than relying on the organic, undocumented order in which they happened to learn things. Second, it creates the artifact that a future learner needs -- the map the original operator navigated without but the successor cannot. Third, it reveals gaps in the documentation itself: if a curriculum module requires knowledge that no article teaches, either the curriculum is wrong or the documentation is incomplete.

## 4. System Model

### 4.1 Curriculum Architecture

The curriculum is organized into three hierarchical levels:

**Level 1: Learning Paths.** A learning path is a complete curriculum for a defined role or objective. Examples: "New Operator Onboarding Path," "Security Specialist Path," "Data Management Path." Each learning path spans multiple domains and brings the learner to a defined competence stage across a set of competencies. A person may follow multiple learning paths concurrently or sequentially. The institution maintains at minimum one mandatory learning path: the New Operator Onboarding Path, which is referenced in D9-011.

**Level 2: Modules.** A module is a self-contained unit of learning focused on a single competency or a tightly related cluster of competencies. Examples: "Module: Backup Operations," "Module: Access Control," "Module: Document Lifecycle Management." Each module has defined learning objectives, a list of documentation articles to study, practice exercises, and an assessment gate at the end. A module typically corresponds to one or two operational articles and takes between four and twenty hours to complete, depending on complexity.

**Level 3: Units.** A unit is a single learning activity within a module. A unit might be "Read SEC-002 Sections 1-4," "Complete practical exercise: create a new user account with appropriate permissions," or "Answer knowledge check questions on the principle of least privilege." Units are the atomic element of the curriculum. They are sequenced within the module so that each unit builds on the previous one.

### 4.2 Learning Objective Design

Every module begins with learning objectives. A learning objective states what the learner will be able to do after completing the module, expressed as an observable, assessable behavior. Learning objectives follow this template:

"After completing this module, the learner will be able to [action verb] [specific task or knowledge] [to what standard] [under what conditions]."

Examples:
- "After completing this module, the learner will be able to perform a complete backup cycle, including verification, to the standards defined in D6-002, using only the documentation and the institution's backup hardware."
- "After completing this module, the learner will be able to explain the institution's authority model, including all decision tiers and their requirements, without consulting documentation, to a standard sufficient for Stage 3 assessment per D9-003."

Learning objectives must be specific enough to assess. "Understand backup procedures" is not a learning objective. "Perform a verified backup and restore" is. Every learning objective must map to at least one assessment method defined in D9-003.

### 4.3 Prerequisite Mapping

Every module declares its prerequisites: the modules that must be completed (or the competencies that must be demonstrated) before the module can begin. Prerequisites are of two types:

**Hard prerequisites.** The module cannot be meaningfully attempted without this prior knowledge. Attempting it would either be futile (the learner cannot understand the material) or dangerous (the learner might misapply procedures). Hard prerequisites must be verified through assessment before the module begins.

**Soft prerequisites.** The module is enhanced by this prior knowledge but can be completed without it, possibly with reduced comprehension or efficiency. Soft prerequisites are recommended but not enforced.

The prerequisite map is maintained as a directed acyclic graph (DAG). No circular dependencies are permitted. The map is stored as a documented table in the curriculum record and is reviewed whenever a new module is added or an existing module is substantially revised.

**Prerequisite Map Maintenance Procedure:**
1. When creating a new module, identify every concept, skill, or piece of knowledge the module assumes the learner already possesses.
2. For each assumption, determine whether it is taught in an existing module. If yes, that module is a prerequisite. If no, either create a new module to teach it or add it to the new module's early units as foundational material.
3. Classify each prerequisite as hard or soft.
4. Verify that the addition does not create a circular dependency in the DAG.
5. Update the prerequisite map document.
6. Review the affected learning paths to ensure the new prerequisite does not create an unreasonable barrier.

### 4.4 Creating Lessons from Documentation Articles

The primary raw material for curriculum modules is the existing documentation corpus. The procedure for converting a documentation article into a learning sequence is as follows:

**Step 1: Identify the article's learning content.** Read the article and list every distinct concept, procedure, and principle it contains. Use the section structure as a guide: Purpose teaches "why this matters," Background teaches "what context is needed," System Model teaches "how it works," Rules teaches "what the constraints are," Failure Modes teaches "what can go wrong," and Recovery teaches "what to do about it."

**Step 2: Sequence the content.** Order the learning content from foundational to advanced. This sequence often differs from the article's section order. The article may reference concepts in Section 4 that require understanding from Section 3 of a different article. The curriculum sequence accounts for these cross-article dependencies.

**Step 3: Break the content into units.** Each unit should take between thirty minutes and two hours to complete. A unit should teach one concept or one procedure, not both simultaneously. If a procedure requires understanding a concept, the concept unit precedes the procedure unit.

**Step 4: Add practice activities.** For each procedural unit, design a practice exercise. The exercise should use real or simulated institutional systems. For each conceptual unit, design a knowledge check. The practice activities should align with the assessment methods in D9-003 so that the learner's practice prepares them for formal assessment.

**Step 5: Write the connective narrative.** Between units, add brief text that explains how the current unit connects to the next. This narrative provides the "why are we learning this next?" context that the standalone documentation article does not supply.

**Step 6: Define the assessment gate.** At the end of the module, define the assessment that the learner must pass to certify completion. The assessment gate must align with the module's learning objectives and with the appropriate competence stage per D9-003.

### 4.5 The Curriculum Review Schedule

Curricula are living documents. They must be reviewed and updated on a defined schedule:

**Triggered review:** Whenever a documentation article that serves as a curriculum source is substantially revised, every module that references that article must be reviewed within thirty days. "Substantially revised" means changes to procedures, system models, or rules -- not editorial corrections or Commentary Section additions.

**Post-assessment review:** When assessment results per D9-003 reveal a systematic pattern of learners failing to achieve the expected competence stage after completing a module, the module is reviewed for effectiveness. Systematic failure is defined as the same competency gap appearing in two or more consecutive assessments.

**Annual review:** Once per year, the complete curriculum -- all learning paths, modules, and prerequisite maps -- is reviewed for currency, completeness, and alignment with D9-002 competency requirements. The annual review is integrated into the operational calendar per OPS-001.

**Succession-triggered review:** When a succession event is initiated, the New Operator Onboarding Path is reviewed immediately and in its entirety. The incoming operator's experience with the onboarding path is documented as feedback for the next revision.

## 5. Rules & Constraints

- **R-D9-04-01:** Every competency defined in D9-002 must be addressed by at least one curriculum module. Competencies without corresponding modules represent training gaps that must be resolved.
- **R-D9-04-02:** Every curriculum module must have defined learning objectives expressed as observable, assessable behaviors. Vague objectives like "understand" or "be familiar with" are not acceptable.
- **R-D9-04-03:** Every curriculum module must declare its prerequisites. Modules with no prerequisites must explicitly state "None" rather than leaving the field blank.
- **R-D9-04-04:** The prerequisite map must be a directed acyclic graph. Circular dependencies are a design error and must be resolved before the curriculum is published.
- **R-D9-04-05:** The New Operator Onboarding Path is mandatory and must be maintained at all times. Its first module must require no prerequisites beyond general literacy and technical aptitude per D9-001 Section 3.3.
- **R-D9-04-06:** Curriculum modules must reference the specific documentation articles they are built from. If an article is updated, the referencing module must be reviewed within thirty days.
- **R-D9-04-07:** The complete curriculum must be reviewed annually. The review must verify alignment with D9-002, currency of all source article references, and the validity of the prerequisite map.

## 6. Failure Modes

- **Curriculum-documentation drift.** The documentation evolves but the curriculum is not updated. Learners study outdated procedures and develop incorrect competencies. Mitigation: R-D9-04-06 triggered review and R-D9-04-07 annual review. The documentation change log should be cross-referenced against the curriculum module list during every review.

- **Prerequisite spirals.** Over time, the prerequisite chain grows so long that reaching any advanced module requires completing an unreasonable number of prior modules. The curriculum becomes a barrier to learning rather than a facilitator. Mitigation: during the annual review, measure the longest prerequisite chain for each learning path. If any chain exceeds what a motivated learner can complete in a reasonable timeframe (defined as twice the estimated learning time), reassess which prerequisites are truly hard versus soft.

- **Curriculum without practice.** Modules consist entirely of reading assignments with no exercises. The learner completes the module having read a great deal but having done nothing. Mitigation: R-D9-04-02 requires assessable learning objectives, and assessable objectives require practice activities that prepare for assessment.

- **Orphaned modules.** A module is created for a competency that is later removed from D9-002 or assigned to a different module. The orphaned module remains in the curriculum, confusing learners who complete it and discover it leads nowhere. Mitigation: the annual review cross-references every module against D9-002. Modules with no corresponding competency are archived or reassigned.

- **The single-path trap.** Only one learning path exists (the onboarding path), and it must be completed in full before the learner can do anything. For a successor who needs to perform a specific function urgently, a single mandatory path is an obstacle. Mitigation: design learning paths with early-exit points. After the core modules, the path should branch into domain-specific tracks. A learner who needs backup competence urgently can follow the backup track without completing the security track first, as long as the hard prerequisites are satisfied.

## 7. Recovery Procedures

1. **If no curriculum exists:** Begin with the New Operator Onboarding Path. Its first module should cover the five Core Charter articles (ETH-001, CON-001, GOV-001, SEC-001, OPS-001) in sequence. Its second module should cover the domain philosophy articles for the domains most critical to immediate operations. Build additional modules as capacity permits, prioritizing competencies classified as critical in D9-002.

2. **If the curriculum has drifted significantly from current documentation:** Freeze the curriculum. Conduct a full alignment review, comparing every module against its source articles. Mark modules as "under review" until they are verified. Update modules in priority order: critical competencies first, then remaining modules by learning path order.

3. **If a learner cannot complete the onboarding path due to excessive prerequisites:** Identify the minimum set of modules required for the learner's immediate operational needs. Create a temporary abbreviated path that covers those modules. Record the abbreviation and the rationale in the decision log. Schedule the remaining modules as soon as the learner's immediate operational burden permits.

4. **If assessment consistently shows that a module is not teaching effectively:** Do not simply add more content. First, review the learning objectives -- are they correctly defined? Then review the units -- are concepts presented before procedures? Then review the practice activities -- do they genuinely prepare for the assessment? Revise the specific component that is failing rather than rewriting the entire module.

## 8. Evolution Path

- **Years 0-5:** The curriculum is being built alongside the documentation. The founder creates the first onboarding path by documenting the sequence in which they would teach someone the institution from scratch. This initial curriculum is rough. It will be heavily revised after the first real learner attempts it. The priority is existence over perfection.

- **Years 5-15:** The curriculum should be tested by at least one real learner. This test is the most valuable feedback the curriculum will ever receive. Every point of confusion, every prerequisite gap, every module that takes longer than expected -- these are the data that transform a theoretical curriculum into an effective one.

- **Years 15-30:** The curriculum should be mature enough to function without the founder. The prerequisite map is well-established. Multiple learning paths exist for different operational roles. The annual review is routine. The primary challenge is preventing ossification -- ensuring that the curriculum evolves with the institution rather than preserving a historical learning sequence that no longer reflects current operations.

- **Years 30-50+:** The curriculum has been used by multiple learners across potentially multiple succession events. It is a tested, refined educational instrument. The documentation articles it references may have been revised many times, but the curriculum's structure -- learning paths, modules, units, prerequisite maps -- has proven its value as the layer that transforms a reference library into a teaching system.

- **Signpost for revision:** If a new operator completes the onboarding path and still cannot perform critical functions at Stage 2, the curriculum has failed its primary mission and requires fundamental redesign.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
Building a curriculum for an institution where I am the only student feels like building a road in a wilderness. I know the destination. I know the terrain. But I have never had to explain the route to someone else, and the act of drawing the map is revealing how much of my own learning was accidental, undirected, and poorly sequenced.

The hardest part of curriculum design is not creating the content -- the documentation corpus is the content. The hardest part is sequencing. I learned this institution's systems in the order I built them, which is emphatically not the order someone should learn them. I built the security model after the backup system because the backup system was more urgent, but a learner should understand security philosophy before learning backup procedures because backup decisions depend on security principles. The curriculum must impose a pedagogical order that differs from the historical order.

I am also discovering prerequisite dependencies I did not realize existed. Every time I start designing a module, I find myself listing knowledge that I acquired informally and never documented as a dependency. This is exactly the kind of tacit prerequisite that will trip up a future learner. The curriculum design process is forcing me to make these invisible dependencies visible, which is valuable in its own right.

## 10. References

- ETH-001 -- Ethical Foundations of the Institution (Principle 4: Longevity Over Novelty)
- CON-001 -- The Founding Mandate (documentation completeness, institutional transferability)
- GOV-001 -- Authority Model (succession protocol, onboarding requirements per Section 7)
- OPS-001 -- Operations Philosophy (operational calendar, documentation-first principle)
- D6-001 -- Data Philosophy (Tier 1 classification for curriculum materials)
- D6-003 -- Format Longevity Doctrine (training material format constraints)
- D9-001 -- Education Philosophy (self-teaching requirement, four-stage competence model, documentation-as-teacher principle)
- D9-002 -- The Competency Model (competency definitions, role-based requirements, critical competency list)
- D9-003 -- Competency Assessment Framework (assessment methods, scoring rubrics, assessment gates)
- D9-005 -- Tacit Knowledge Capture Procedures (tacit knowledge integration into curriculum)
- D9-007 -- Documentation-as-Teacher (documentation quality requirements for curriculum source material)
- D9-011 -- New Member Onboarding Sequence (mandatory onboarding learning path)
- Stage 1 Meta-Framework (META-00-ART-001), Section 3: Article Template

---

---

# D9-005 -- Tacit Knowledge Capture Procedures

**Document ID:** D9-005
**Domain:** 9 -- Education & Training
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, OPS-001, D9-001, D9-002, D9-003
**Depended Upon By:** D9-004, D9-006, D9-007. Referenced by all articles involving knowledge preservation or succession preparation.

---

## 1. Purpose

This article defines how the holm.chat Documentation Institution extracts, documents, and preserves knowledge that exists only in someone's head. Tacit knowledge is the knowledge that resists documentation: the intuition that something is wrong before any diagnostic confirms it, the judgment about when to deviate from a procedure, the understanding of why a system was built a certain way that never made it into the written record, the muscle memory of a procedure performed hundreds of times.

D9-001 Section 3.1 established the critical distinction between explicit and tacit knowledge. Explicit knowledge can be written down and is -- that is what the documentation corpus is. Tacit knowledge cannot be easily written down because the person who possesses it often cannot articulate it. They know more than they can say. This article provides the procedures for helping them say it, or at least for capturing approximations of it that are useful to those who come after.

The urgency of tacit knowledge capture increases as the institution ages. In the early years, the founder's tacit knowledge is constantly exercised and therefore constantly available. But tacit knowledge is perishable. Skills not exercised degrade. Context not revisited fades. And most critically, when a person becomes unavailable -- through succession, incapacitation, or death -- every piece of uncaptured tacit knowledge is lost permanently. This article exists to wage an ongoing campaign against that loss, knowing the campaign can never fully succeed but that partial success is vastly better than no effort at all.

## 2. Scope

This article covers:

- Identifying tacit knowledge: how to recognize what has not been documented.
- The knowledge extraction session: a structured interview format for drawing out tacit knowledge.
- Think-aloud capture: the procedure for recording an expert's reasoning while performing a task.
- Observation protocols: how to watch an expert work and document what they do that the procedures do not describe.
- The knowledge gap audit: a systematic method for identifying areas where tacit knowledge likely exists but has not been captured.
- Storage and integration: how to preserve captured tacit knowledge and feed it into the documentation and curriculum systems.

This article does not cover:

- The philosophical justification for tacit knowledge capture (see D9-001 Section 3.1).
- Assessment of tacit knowledge (see D9-003, with the caveat that Stage 4 competence inherently includes tacit components).
- Curriculum design that incorporates captured tacit knowledge (see D9-004 Section 4.4).
- Format requirements for captured materials (governed by D9-001 R-D9-05 and D6-003).

## 3. Background

### 3.1 Why Tacit Knowledge Resists Capture

Tacit knowledge resists documentation for three reasons, each of which must be understood to design effective capture procedures.

First, the expert does not know they know it. After years of operating a system, the expert's responses become automatic. They check a log file "just because" before running a procedure. They listen to a disk drive's sound without consciously deciding to. They choose one approach over another based on a pattern they recognized without being able to name it. When asked "How did you know that?" the honest answer is often "I do not know. I just did." Capture procedures must be designed to surface these automatic responses.

Second, tacit knowledge is context-dependent. An expert's judgment about whether an anomaly merits investigation depends on a vast web of contextual factors -- what happened last time, what changed recently, what time of year it is, how the system has been behaving lately. Extracting the judgment without the context produces a rule that is correct in the captured context but misleading in a different one. Capture procedures must capture context alongside conclusion.

Third, tacit knowledge often resides in the body, not the mind. The physical feel of a connector that is about to fail. The visual pattern of a log file that indicates impending trouble. The rhythm of a procedure that tells the expert whether it is running normally. These embodied forms of knowledge are the hardest to capture and require observation-based methods rather than interview-based methods.

### 3.2 The Capture Paradox for a Single Operator

In a single-operator institution, the person whose tacit knowledge needs to be captured and the person who must perform the capture are the same individual. This is the capture paradox. It is difficult to observe your own automatic behaviors because the act of self-observation disrupts them. It is difficult to interview yourself because you share the same blind spots.

The procedures in this article address the capture paradox through two strategies. First, they provide structured self-capture methods: guided self-interviews with specific question protocols, think-aloud exercises performed while executing tasks, and post-task reflection journals that prompt for tacit elements. Second, they define protocols for when a second person is available, so that opportunities for external capture are not wasted.

### 3.3 Tacit Knowledge as Institutional Risk

D9-001 identifies the "undocumented expert" as a primary failure mode. Tacit knowledge capture is the preventive control for that risk. The institution's risk exposure from uncaptured tacit knowledge is proportional to two factors: the criticality of the function (how much damage occurs if the knowledge is lost) and the exclusivity of the knowledge (how many people possess it). A function that is critical and known to only one person represents the highest risk. The knowledge gap audit in Section 4.5 uses this risk framework to prioritize capture efforts.

## 4. System Model

### 4.1 Identifying Tacit Knowledge

Tacit knowledge hides in predictable places. The following indicators suggest that tacit knowledge exists and has not been captured:

**Indicator 1: The undocumented decision.** A procedure's documentation says what to do but a step exists where the operator makes a choice based on judgment rather than documented criteria. "Check the log and decide if it looks normal" is a judgment call that depends on tacit knowledge of what "normal" looks like.

**Indicator 2: The routine deviation.** The operator regularly deviates from the documented procedure in small ways -- skipping a step that is rarely needed, adding a check that is not in the procedure, performing steps in a different order. These deviations represent learned optimizations that the documentation has not absorbed.

**Indicator 3: The diagnostic shortcut.** When something goes wrong, the operator goes directly to the likely cause without following the full diagnostic procedure. This shortcut reflects pattern-matching ability developed through experience.

**Indicator 4: The contextual warning.** The operator knows that a particular operation is risky under certain conditions that the documentation does not mention. "Do not run the integrity check during backup because it causes I/O contention" is explicit knowledge only if it is documented. If it lives only in the operator's head, it is tacit.

**Indicator 5: The inherited workaround.** A procedure works around a limitation that was never formally documented. The workaround has become so routine that the operator no longer recognizes it as a workaround. When asked why they do it that way, they say "It has always been done that way."

### 4.2 The Knowledge Extraction Session

The knowledge extraction session is a structured interview designed to draw out tacit knowledge from an expert. It can be conducted by a second person interviewing the expert, or it can be conducted as a guided self-interview by the expert alone. The session produces a written or recorded artifact that is then processed into documentation.

**Preparation (30 minutes before the session):**
1. Select the function or domain to be explored. Focus on one function per session. Broad sessions produce shallow results.
2. Review the existing documentation for that function. Note where the documentation says "check," "verify," "assess," or "determine" without specifying exactly how. These are tacit knowledge hotspots.
3. Prepare the question protocol. Use the standard questions below, adapted to the specific function.

**The Standard Question Protocol:**

Phase 1 -- Process walkthrough questions:
- "Walk me through how you actually perform [function], step by step, including anything you do that is not in the documentation."
- "At which steps do you make a judgment call? What information do you use to make that judgment?"
- "Are there any steps you sometimes skip? Under what conditions?"
- "Are there any steps you add that are not documented? Why?"

Phase 2 -- Diagnostic questions:
- "When something goes wrong with [function], what is the first thing you check? Why that first?"
- "How do you tell the difference between a minor anomaly and a serious problem?"
- "Can you describe a time when you caught a problem early? What tipped you off?"
- "What are the warning signs that [function] is about to fail, even before it actually fails?"

Phase 3 -- Context and history questions:
- "Why is [function] done this way instead of some other way? What alternatives were considered?"
- "What has changed about how you perform [function] since you started? Why did it change?"
- "What is the most common mistake someone new to [function] would make? How would you prevent it?"
- "If you could only tell your successor three things about [function], what would they be?"

Phase 4 -- Edge case questions:
- "Under what conditions would you deviate from the documented procedure? Give a specific example."
- "What is the most unusual situation you have encountered with [function]? How did you handle it?"
- "What assumptions does the current procedure make that might not always be true?"

**Session execution (60-90 minutes):**
Proceed through the question protocol. Record the session if possible (audio or written notes). Do not edit the expert's responses during the session. Capture everything, even statements that seem obvious. What seems obvious to the expert is often the most valuable tacit knowledge.

**Post-session processing (60-120 minutes):**
1. Review the session record. Identify every piece of information that is not in the existing documentation.
2. Classify each piece as: (a) should be added to the procedure documentation, (b) should be added to the Commentary Section of the relevant article, (c) should be preserved as a supplementary artifact referenced by the documentation, or (d) is too context-dependent to generalize but should be preserved as a dated case study.
3. Create the appropriate documentation updates or artifacts.
4. Record the session in the tacit knowledge capture log: date, function explored, key findings, documentation changes made.

### 4.3 Think-Aloud Capture

Think-aloud capture is the procedure for recording an expert's reasoning in real time while they perform a task. It is the most effective method for capturing diagnostic reasoning, procedural judgment, and situational awareness.

**Procedure:**
1. Select a task to be performed. Ideally, this is a real operational task, not a simulation, so that the expert's natural behavior is captured.
2. Before beginning, instruct the expert (or instruct yourself, for self-capture): "As you perform this task, speak aloud everything you are thinking. Say what you are looking at, what you are checking, what you are deciding, and why. If you notice something, say what you noticed and why it matters. Do not filter for relevance -- say everything."
3. Record the think-aloud using audio recording, or have a second person take notes. If self-capturing, use audio recording -- it is not possible to write notes and perform think-aloud simultaneously.
4. Perform the task from start to finish with continuous narration.
5. After the task, review the recording. Identify statements that reveal tacit knowledge: judgment calls, pattern recognition, contextual checks, procedural deviations, and risk assessments that are not documented.
6. Process the findings per the post-session procedure in Section 4.2.

**Self-capture adaptation:** When performing think-aloud capture alone, the act of narrating disrupts the automatic behavior that the capture is trying to surface. Accept this limitation. The disruption itself is informative: when you find yourself unable to explain why you are doing something, that is a strong signal of deeply tacit knowledge. Note these moments explicitly: "I am doing X but I cannot articulate why."

### 4.4 Observation Protocol

When two people are available, observation is the most powerful capture method. The observer watches the expert perform tasks and documents behaviors, decisions, and patterns that the expert may not be aware of.

**Observer instructions:**
1. Do not interrupt the expert during the task unless safety requires it. Observation captures natural behavior; interruption disrupts it.
2. Record the following: every step the expert performs, including steps not in the documentation; every check or glance the expert makes, including visual checks of indicators, sounds, or physical states; every pause where the expert appears to be making a decision; every deviation from the documented procedure; and every verbal comment the expert makes (even muttering).
3. After the task, conduct a brief debriefing: show the expert your observation notes and ask them to explain any behavior you noted that you did not understand.
4. Process the findings per the post-session procedure in Section 4.2.

### 4.5 The Knowledge Gap Audit

The knowledge gap audit is a systematic review that identifies functions where tacit knowledge likely exists but has not been captured. It should be conducted annually, aligned with the curriculum review per D9-004.

**Procedure:**
1. List every function in the institution's competency model per D9-002.
2. For each function, answer two questions: (a) Is there a person who performs this function significantly better or faster than the documentation alone would enable? If yes, tacit knowledge exists. (b) Has a tacit knowledge capture session been conducted for this function within the past twenty-four months? If no, capture may be overdue.
3. Score each function on a risk matrix: criticality (how important is this function?) multiplied by exclusivity (how many people can perform it?). High criticality and high exclusivity equals high risk.
4. Rank the functions by risk score. The highest-risk functions with uncaptured tacit knowledge are the priority for the next capture cycle.
5. Schedule capture sessions for the top-priority functions per the operational calendar.
6. Record the audit results in the tacit knowledge capture log.

## 5. Rules & Constraints

- **R-D9-05-01:** Every Tier 1 institutional function (as classified in D9-002) must undergo at least one tacit knowledge capture session within twenty-four months of being identified. This requirement is ongoing: when functions are reclassified or new Tier 1 functions are identified, the clock resets.
- **R-D9-05-02:** Tacit knowledge capture sessions must produce a documented artifact: a written record, an audio recording, or a video recording. Sessions without artifacts do not count toward the requirement.
- **R-D9-05-03:** Findings from capture sessions must be integrated into the documentation or curriculum within sixty days. Capture without integration is collection without purpose.
- **R-D9-05-04:** The knowledge gap audit must be conducted annually. Results must be recorded and must inform the next year's capture schedule.
- **R-D9-05-05:** All tacit knowledge artifacts must comply with the format longevity requirements of D6-003. Audio and video must be in open formats. Written records must be in plain text or Markdown.
- **R-D9-05-06:** When a succession event is imminent or in progress, tacit knowledge capture becomes the highest-priority education activity. All other education activities may be deferred to maximize capture from the outgoing operator.

## 6. Failure Modes

- **Capture procrastination.** Tacit knowledge capture is always less urgent than operational tasks. It is easy to defer indefinitely because the consequences of not doing it are invisible until someone is gone. Mitigation: R-D9-05-01 establishes a hard deadline. The annual knowledge gap audit creates a scheduled trigger. The quarterly governance review checks compliance.

- **Shallow capture.** The capture session is conducted but remains at the surface level. The expert describes what they do but not why, or provides general principles without specific examples. Mitigation: the question protocol in Section 4.2 is designed to push past surface answers. Phase 2 (diagnostic questions) and Phase 4 (edge case questions) specifically target the deeper layers of tacit knowledge.

- **Capture without integration.** Sessions are recorded, artifacts are filed, but the findings never make it into the documentation or curriculum. The artifacts exist but nobody consults them. Mitigation: R-D9-05-03 requires integration within sixty days. The post-session processing procedure explicitly classifies each finding into a documentation or curriculum action.

- **Self-capture blind spots.** The single operator cannot surface their own blind spots because, by definition, they cannot see them. Mitigation: accept this limitation. Use the indicators in Section 4.1 to find blind spots indirectly. Use think-aloud capture to reveal moments where you act without explanation. When a second person becomes available, prioritize having them observe you rather than interview you -- observation reveals what self-reporting cannot.

- **Over-capture.** The operator becomes so thorough in capturing tacit knowledge that the resulting artifacts are voluminous, unorganized, and impossible for a future learner to navigate. Mitigation: the post-session processing procedure includes a classification step that routes findings to specific documentation locations. Raw artifacts should be preserved but clearly labeled as raw. The processed, integrated version in the documentation is what learners should consult.

## 7. Recovery Procedures

1. **If tacit knowledge capture has never been conducted:** Begin with the knowledge gap audit (Section 4.5). Identify the three highest-risk functions. Schedule one capture session per week until those three functions are captured. Then continue at a sustainable pace for remaining functions.

2. **If capture artifacts exist but have not been integrated:** Conduct a triage of all unprocessed artifacts. For each, apply the post-session processing procedure retroactively. Prioritize integration for Tier 1 functions. Establish a regular processing schedule to prevent future backlog.

3. **If the sole expert is imminently unavailable (emergency capture):** Abandon the standard session format. Use the Phase 4 question: "If you could only tell your successor three things about each function, what would they be?" Record everything. Do not process during the session. Capture breadth over depth -- some information about every function is more valuable than deep information about one function when time is critically limited.

4. **If captured tacit knowledge is in obsolete formats:** Migrate to current open formats per D6-003. If the content cannot be migrated (e.g., the playback technology no longer exists), treat the knowledge as lost and conduct new capture sessions for the affected functions. Document the format failure as a lessons-learned entry in the Commentary Section of D6-003.

## 8. Evolution Path

- **Years 0-5:** The founder captures their own tacit knowledge, which is the most awkward and least effective mode of capture. The primary value in this period is building the habit and the infrastructure: establishing the capture log, conducting the first knowledge gap audits, learning what the question protocol reveals and does not reveal. Expect the first attempts to feel forced and to produce less value than expected. Persist. The skill of self-capture improves with practice.

- **Years 5-15:** If a second person becomes available, even temporarily, this is the golden opportunity for external capture. A fresh observer will notice things about the founder's behavior that the founder has never been able to articulate. Prioritize observation-based capture during these periods.

- **Years 15-30:** The capture log now spans years. It becomes a historical record of institutional knowledge evolution. Early captures can be compared with current practices to identify how tacit knowledge has changed as the institution matured. This longitudinal view is itself a form of institutional knowledge.

- **Years 30-50+:** Tacit knowledge capture should be a routine institutional practice. Each succession event generates a burst of capture activity. The procedures should be well-calibrated from years of use. The primary challenge is ensuring that captures remain genuine attempts to surface real knowledge rather than pro-forma exercises conducted to satisfy the requirement.

- **Signpost for revision:** If a succession event results in the loss of critical tacit knowledge despite capture efforts, the capture procedures are insufficient and must be redesigned.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The idea of interviewing myself about my own tacit knowledge is, frankly, absurd on its face. The whole point of tacit knowledge is that I do not know I have it. How am I supposed to ask myself about things I do not know I know?

And yet, the first time I tried the think-aloud capture while performing a backup -- narrating every thought, every glance at a status indicator, every decision -- I discovered three things I do that are not in my backup documentation. I check the disk temperature before starting a large write. I glance at the system load to estimate how long the backup will take. I listen for a change in the disk access pattern that indicates the backup has moved from small files to large files. None of these are in the procedure. All of them inform how I manage the backup process. All of them would be invisible to a successor following my documentation.

The think-aloud method is uncomfortable. Narrating your own thoughts feels self-conscious, and the self-consciousness disrupts the natural flow. But the disruption is the point. When I stumble and say, "I am checking this because... I do not actually know why," that is the moment of discovery. That is the tacit knowledge surfacing, even if it surfaces as a question rather than an answer.

I expect that most of my tacit knowledge will take years to fully capture, and some will never be captured at all. The goal is not perfection. The goal is to reduce the gap between what I know and what the documentation preserves. Every session narrows that gap, even if it can never close it.

## 10. References

- ETH-001 -- Ethical Foundations of the Institution (Principle 3: Transparency of Operation; Principle 4: Longevity Over Novelty)
- CON-001 -- The Founding Mandate (institutional transferability, knowledge preservation)
- GOV-001 -- Authority Model (succession protocol, knowledge transfer during succession)
- OPS-001 -- Operations Philosophy (documentation-first principle, operational tempo)
- D6-001 -- Data Philosophy (Tier 1 classification for knowledge artifacts)
- D6-003 -- Format Longevity Doctrine (format requirements for capture artifacts)
- D9-001 -- Education Philosophy (tacit knowledge philosophy, Section 3.1; self-teaching requirement; undocumented expert failure mode)
- D9-002 -- The Competency Model (function classification, Tier 1 function identification)
- D9-003 -- Competency Assessment Framework (Stage 4 assessment includes tacit components)
- D9-004 -- Curriculum Design and Maintenance (integration of captured tacit knowledge into curricula)
- D9-007 -- Documentation-as-Teacher (writing documentation that encodes tacit knowledge)

---

---

# D9-006 -- Cross-Training Procedures

**Document ID:** D9-006
**Domain:** 9 -- Education & Training
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, OPS-001, D9-001, D9-002, D9-003, D9-004
**Depended Upon By:** D9-007, D9-011. Referenced by all articles involving succession, resilience, or single-point-of-failure analysis.

---

## 1. Purpose

This article defines how the holm.chat Documentation Institution ensures that no critical function depends on a single person's knowledge. Cross-training is the practice of deliberately ensuring that multiple people can perform each essential function, so that the loss or unavailability of any one person does not cripple the institution. In a multi-person institution, cross-training is achieved through structured rotation, paired operations, and skill shadowing. In a single-operator institution, cross-training takes a different form: it means ensuring that the documentation, curricula, and knowledge artifacts are sufficient for a second person to achieve competence, and it means the operator actively practices functions they do not routinely perform so that skills do not atrophy through disuse.

D9-001 established the cross-training imperative in R-D9-04: "Cross-training must ensure that no institutional function is dependent on a single person's knowledge." This article provides the operational machinery for fulfilling that imperative. It defines the cross-training matrix, the minimum coverage requirements, the procedures for cross-training when multiple people are available, and the self-sustaining practices for when only one person is present.

The stakes are simple. A single point of knowledge failure is a single point of institutional failure. If only one person can perform a function and that person is unavailable, the function does not get performed. If the function is critical -- backup verification, security audit, data recovery -- the institution is degraded and potentially at risk. Cross-training is the insurance policy against this risk.

## 2. Scope

This article covers:

- The cross-training matrix: a comprehensive map of who can perform what, at what competence stage.
- Minimum coverage requirements: the minimum number of people (or the minimum documentation quality) required for each function.
- Cross-training schedules: how often cross-training activities occur.
- Paired operations: the practice of performing tasks with a partner for knowledge transfer.
- Skill shadowing: the practice of observing someone perform a task you will need to learn.
- Skill maintenance: how to prevent cross-trained skills from atrophying.
- Single-operator cross-training: what cross-training means when there is only one person.

This article does not cover:

- The definition of institutional functions (see D9-002).
- Assessment of cross-trained competencies (see D9-003).
- Curriculum design for cross-training modules (see D9-004).
- Tacit knowledge transfer during cross-training (see D9-005).

## 3. Background

### 3.1 The Bus Factor

The "bus factor" of a project is the number of people who would need to be hit by a bus before the project is unable to continue. A bus factor of one means the loss of a single person is fatal. For a single-operator institution, the bus factor is inherently one for all functions. Cross-training cannot change this arithmetic when there is only one person. But it can change two things: the speed at which a successor achieves competence (which is the difference between a brief disruption and an institutional crisis), and the depth of the documentation and knowledge artifacts available to that successor.

When multiple people are present -- during apprenticeship, advisory relationships, or post-succession periods -- cross-training becomes the direct mechanism for raising the bus factor above one. This article defines procedures for both realities.

### 3.2 Skill Atrophy

A skill that is not exercised degrades. Research on skill retention shows that procedural skills (how to perform a task) degrade faster than conceptual understanding (why the task is done), and that complex skills degrade faster than simple ones. An operator who was cross-trained on a function three years ago and has not performed it since may retain awareness (Stage 1) and perhaps the ability to follow documentation (Stage 2), but their independent execution capability (Stage 3) has likely degraded.

Cross-training is not a one-time activity. It requires ongoing maintenance -- periodic practice, reassessment, and refresher training. The cross-training schedule in this article accounts for skill atrophy by defining maintenance intervals for cross-trained competencies.

### 3.3 Cross-Training in a Single-Operator Institution

When there is only one person, "cross-training" means ensuring transferability rather than actual multi-person coverage. The single operator's cross-training obligations are:

First, ensuring that all critical functions have documentation of sufficient quality to bring a new person to Stage 2 competence without live instruction. This is the self-teaching requirement of D9-001 applied specifically to cross-training.

Second, periodically performing functions they do not routinely perform, to verify that the documentation works and to maintain their own breadth of competence. An operator who normally manages backups but never performs security audits should periodically perform a security audit to verify both their competence and the documentation's completeness.

Third, maintaining the cross-training matrix as a planning document that identifies what a second person would need to learn, in what order, to achieve minimum coverage.

## 4. System Model

### 4.1 The Cross-Training Matrix

The cross-training matrix is a two-dimensional table mapping every institutional function to every person who can perform it, along with their competence stage. For a single-operator institution, the matrix has one person column and serves as a transferability readiness assessment. When additional people are present, the matrix expands to track multi-person coverage.

**Matrix structure:**

| Function | Person A (Stage) | Person B (Stage) | Coverage Status | Documentation Quality |
|----------|-----------------|-----------------|----------------|----------------------|
| Full backup cycle | 4 | 2 | Covered | Verified self-teachable |
| Key rotation | 3 | 1 | At risk | Needs review |
| OS update | 4 | -- | Single point | Verified self-teachable |

**Column definitions:**
- **Function:** The institutional function, drawn from the D9-002 competency model.
- **Person columns:** The competence stage of each person for that function, per D9-003 assessment. A dash indicates no competence.
- **Coverage Status:** "Covered" means at least two people at Stage 2 or above. "At risk" means only one person at Stage 2 or above, but documentation is self-teachable. "Single point" means only one person can perform the function and documentation has not been verified as self-teachable. "Critical gap" means the function cannot be performed by any currently available person.
- **Documentation Quality:** Whether the documentation for this function has been verified as self-teachable per the self-teaching test in D9-007.

**Matrix maintenance:** The matrix is updated whenever an assessment is completed (D9-003), whenever a person joins or leaves the institution, and during the annual comprehensive review. The matrix is stored in the documentation corpus as a Tier 1 document.

### 4.2 Minimum Coverage Requirements

The institution defines three tiers of coverage requirements based on function criticality:

**Tier 1 -- Critical Functions.** Functions without which the institution cannot operate or whose failure would cause irreversible harm: data backup and recovery, security key management, system integrity verification, and any function classified as critical in D9-002. Minimum coverage: at least two people at Stage 2 or above, OR one person at Stage 3 or above with documentation verified as self-teachable to Stage 2. In a single-operator institution, the standard is: one person at Stage 3 or above with documentation verified as self-teachable to Stage 2.

**Tier 2 -- Important Functions.** Functions that support institutional operations but whose temporary absence does not cause irreversible harm: documentation maintenance, curriculum updates, routine system administration, operational logging. Minimum coverage: at least one person at Stage 2 or above with documentation verified as self-teachable.

**Tier 3 -- Supporting Functions.** Functions that enhance institutional capability but whose absence can be tolerated for extended periods: advanced diagnostics, optimization, non-critical automation, educational material creation. Minimum coverage: at least one person at Stage 1 (Awareness) with documentation available.

Any function that falls below its tier's minimum coverage requirement is flagged in the cross-training matrix as a gap requiring remediation.

### 4.3 Cross-Training Procedures for Multiple People

When two or more people are available, the following cross-training methods are used:

**Paired Operations.** Two people perform a task together. The primary operator performs the task while the secondary operator observes, asks questions, and takes notes. After the task, roles are discussed: the primary explains decisions made, the secondary identifies steps they did not understand. On the next occurrence of the task, roles swap: the secondary performs while the primary supervises.

Paired operations schedule: every Tier 1 function should be performed as a paired operation at least once per quarter when two people are available. Tier 2 functions should be paired at least once per six months.

**Skill Shadowing.** The learner observes the expert performing their normal duties over an extended period (minimum one full operational day). Unlike paired operations, shadowing is observational only -- the learner does not participate but watches, listens, and notes. Shadowing is particularly effective for capturing the routine decision-making, prioritization, and situational awareness that formal training does not teach.

Shadowing schedule: when a new person joins the institution, they should shadow the primary operator for at least five full operational days before beginning independent work. During ongoing operations, each person should shadow the other's domain-specific work for at least one day per quarter.

**Rotation.** People periodically swap primary responsibility for functions. Person A takes over Person B's functions for a defined period (one week minimum) while Person B is available for consultation but does not perform the function. Rotation tests whether the non-primary operator can genuinely perform the function independently.

Rotation schedule: Tier 1 functions should be rotated at least once per year. Each rotation period should be long enough to encounter at least one full cycle of the function (e.g., a full backup rotation cycle for backup functions).

### 4.4 Single-Operator Cross-Training Practices

When only one person is present, cross-training becomes self-maintenance and transferability verification:

**Practice rotations.** The operator identifies functions they do not routinely perform (e.g., disaster recovery procedures, key rotation, full system rebuild). At least quarterly, the operator performs one of these functions -- either in production (if the function is due) or in a test environment. The operator documents any difficulties encountered and updates the procedure documentation accordingly.

**Self-teaching verification.** The operator selects one critical function per quarter and attempts to perform it using only the written documentation, as if they were a new learner. They deliberately do not rely on memory. Every point where the documentation is insufficient is noted and addressed.

**Documentation-as-cross-training.** For every function they perform, the operator maintains documentation at a quality level sufficient for a Stage 2 learner. This is not merely a documentation requirement -- it is the single operator's only viable cross-training strategy. The successor who arrives after the operator is gone will have only the documentation. Its quality is the cross-training.

### 4.5 Skill Maintenance Schedule

Cross-trained skills require periodic exercise to prevent atrophy. The maintenance schedule defines how often each competence stage must be exercised to remain valid:

- **Stage 2 (Guided Execution):** Must be exercised at least every twelve months. If not exercised, competence stage drops to Stage 1 (Awareness) in the cross-training matrix.
- **Stage 3 (Independent Execution):** Must be exercised at least every six months. If not exercised, competence stage drops to Stage 2 in the cross-training matrix.
- **Stage 4 (Adaptive Mastery):** Must be exercised at least every three months for the skill to remain at Stage 4. Without exercise, Stage 4 degrades to Stage 3 within six months.

"Exercised" means performing the function in a real or realistic simulated context. Reading the documentation alone does not count as exercise. Completing a refresher assessment per D9-003 does count.

## 5. Rules & Constraints

- **R-D9-06-01:** The cross-training matrix must be maintained and current at all times. It must be updated within seven days of any assessment, any change in personnel, or any change in function classification.
- **R-D9-06-02:** No Tier 1 function may remain at "Single point" coverage status for more than six months without an active remediation plan. The remediation plan must include either training a second person or verifying that the documentation is self-teachable to Stage 2.
- **R-D9-06-03:** When two or more people are available, paired operations must be conducted for every Tier 1 function at least quarterly.
- **R-D9-06-04:** The single operator must perform at least one practice rotation per quarter, targeting functions not routinely performed.
- **R-D9-06-05:** Cross-trained competencies that have not been exercised within the maintenance schedule defined in Section 4.5 must be downgraded in the matrix. The operator does not get to claim a competence stage that has not been recently demonstrated.
- **R-D9-06-06:** The cross-training matrix must be reviewed as part of every succession planning activity. The matrix directly informs the incoming operator's training priority.
- **R-D9-06-07:** Self-teaching verification (Section 4.4) must be performed for at least one critical function per quarter.

## 6. Failure Modes

- **Matrix neglect.** The cross-training matrix is created once and never updated. Over time, it reflects a historical competence profile that does not match current reality. Decisions based on the stale matrix are unreliable. Mitigation: R-D9-06-01 requires updates within seven days of relevant changes. The quarterly review includes a matrix currency check.

- **Cross-training as checkbox.** Cross-training activities are performed at the minimum required level without genuine engagement. The secondary operator during paired operations is physically present but mentally disengaged. The practice rotation is performed hastily with shortcuts. Mitigation: cross-training activities should be followed by assessment per D9-003. If the post-activity assessment shows no improvement or maintained competence, the activity was ineffective and should be redesigned.

- **Coverage complacency in single-operator mode.** The single operator recognizes that they cannot actually achieve multi-person coverage and stops trying. The cross-training matrix becomes a record of one person's competencies rather than a transferability plan. Mitigation: reframe the single-operator cross-training obligation explicitly around documentation quality and self-teaching verification. The metric is not "Can two people do this?" but "Could a new person learn this from what exists?"

- **Skill atrophy denial.** The operator believes they retain a skill they have not practiced in years. The cross-training matrix shows Stage 3, but actual competence has degraded to Stage 2 or below. Mitigation: R-D9-06-05 automatically downgrades competencies that have not been exercised. The operator must demonstrate, not claim.

- **Rotation avoidance for uncomfortable functions.** The operator avoids practice rotations for functions they dislike or fear, such as disaster recovery or key ceremony procedures. These are often the most critical functions. Mitigation: practice rotation targets should be selected based on the cross-training matrix gaps, not the operator's preferences. The quarterly selection process should explicitly identify the function the operator is least comfortable with as a strong candidate.

## 7. Recovery Procedures

1. **If the cross-training matrix does not exist:** Create it immediately. List every function from D9-002. For each function, record the current operator's assessed competence stage (or estimated stage if no formal assessment exists; schedule formal assessment within thirty days). Assess documentation quality for each function. Identify all gaps against the coverage requirements in Section 4.2.

2. **If Tier 1 functions are at "Single point" status beyond the six-month limit:** This is an urgent institutional risk. For each affected function, initiate documentation quality verification within two weeks. If the documentation is not self-teachable, prioritize making it self-teachable over all non-critical operational work. This is a triage situation: the institution's ability to survive the loss of one person depends on the quality of this documentation.

3. **If cross-training has lapsed for multiple quarters:** Do not attempt to resume all cross-training simultaneously. Prioritize: Tier 1 functions first, then Tier 2. For each function, begin with a self-teaching verification (Section 4.4) to assess current documentation quality. Then schedule practice rotations for functions where personal competence may have degraded.

4. **If a cross-trained person's competence has degraded significantly:** Reclassify them at their current actual stage in the matrix. Design a refresher curriculum module per D9-004 that targets the specific competencies that degraded. Refresher training is typically faster than initial training because the conceptual foundation usually persists even when procedural skill has atrophied.

## 8. Evolution Path

- **Years 0-5:** The cross-training matrix is a single-column document. The operator's primary obligation is documentation quality and self-teaching verification. The practice rotations build the habit of touching unfamiliar functions regularly. This period establishes the infrastructure that cross-training will use when a second person arrives.

- **Years 5-15:** If a second person becomes available, the cross-training matrix gains its second column and the paired operations, shadowing, and rotation procedures activate. This is the period where the matrix becomes a true operational tool rather than a planning document.

- **Years 15-30:** Multiple succession events may have occurred. The cross-training matrix has historical depth -- it can show how coverage has evolved over time. The institution should have refined its sense of which functions are most vulnerable to single-point failures and which cross-training methods are most effective.

- **Years 30-50+:** Cross-training should be deeply embedded in institutional culture. New members expect cross-training as part of their onboarding. The matrix is a routine management tool. The primary challenge is preventing the system from becoming bureaucratic -- ensuring that cross-training activities remain genuine skill-building exercises rather than compliance theater.

- **Signpost for revision:** If the loss of any single person causes an operational disruption lasting more than two weeks for a Tier 1 function, the cross-training system has failed for that function and the procedures must be revised.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
Cross-training procedures for an institution with one person is an exercise in uncomfortable honesty. The matrix has one column. Every function is a single point of failure. The bus factor is one for everything. Writing this article forces me to stare at that reality rather than politely ignore it.

The honest answer to "How do you cross-train when there is only one person?" is: you cannot. Not really. What you can do is prepare for the day when a second person arrives. You can make the documentation self-teachable. You can capture your tacit knowledge. You can maintain the cross-training matrix as a transferability readiness assessment. And you can practice functions you do not routinely perform, so that at least the one person you have is as broadly competent as possible.

The practice rotation requirement is the part of this article I expect to resist the most. It is easier to keep doing the things I know well than to practice the things I rarely do. But the things I rarely do are precisely the things that will be hardest for a successor to learn, because they have the least documentation exercise, the least tacit knowledge capture, and the most accumulated assumptions. Forcing myself to perform them periodically is the only mechanism I have for keeping them honest.

## 10. References

- ETH-001 -- Ethical Foundations of the Institution (Principle 2: Integrity Over Convenience; Principle 4: Longevity Over Novelty)
- CON-001 -- The Founding Mandate (institutional transferability, no single point of failure)
- GOV-001 -- Authority Model (succession planning, knowledge transfer requirements)
- OPS-001 -- Operations Philosophy (operational calendar, sustainability)
- D6-001 -- Data Philosophy (Tier 1 classification for cross-training matrix)
- D9-001 -- Education Philosophy (cross-training imperative R-D9-04, self-teaching requirement R-D9-01, bus factor analysis)
- D9-002 -- The Competency Model (function classification, criticality tiers)
- D9-003 -- Competency Assessment Framework (assessment methods for cross-trained competencies, maintenance assessment)
- D9-004 -- Curriculum Design and Maintenance (refresher curriculum modules, learning path design)
- D9-005 -- Tacit Knowledge Capture Procedures (knowledge transfer during paired operations and shadowing)
- D9-007 -- Documentation-as-Teacher (self-teaching verification test, documentation quality requirements)
- D9-011 -- New Member Onboarding Sequence (initial cross-training for new members)
- GOV-007 -- Governance Health Checks (quarterly matrix review)

---

---

# D9-007 -- Documentation-as-Teacher: Writing for Learning

**Document ID:** D9-007
**Domain:** 9 -- Education & Training
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, OPS-001, D9-001, D9-002, D9-003, D9-004
**Depended Upon By:** All Domain 9 articles. Referenced by every domain where documentation serves an educational function, which is every domain.

---

## 1. Purpose

This article defines how to write documentation that teaches. In this institution, documentation is the primary teacher. D9-001 established this principle. D9-004 defined how to organize documentation into curricula. This article addresses the foundational challenge: how to write individual articles, procedures, and reference documents so that a person reading them can learn -- not merely look up information, but develop genuine understanding.

The distinction matters. A reference document answers the question, "What is the procedure?" A teaching document answers, "What is the procedure, why does it work, what does it look like when done correctly, what does it look like when done incorrectly, and how would you adapt it if circumstances change?" The first serves someone who already understands and needs a reminder. The second serves someone who does not yet understand and needs to learn. This institution's documentation must serve both functions simultaneously, because the same document that serves today's experienced operator must serve tomorrow's brand-new successor.

Writing documentation that teaches is harder than writing documentation that records. It requires the author to model the reader's ignorance -- to remember what it was like to not know what they now know, and to provide the scaffolding that bridges the gap. This article provides the principles, methods, and tests for achieving that standard.

## 2. Scope

This article covers:

- Instructional design principles for documentation: how cognitive science informs documentation writing.
- Worked examples: how to write step-by-step demonstrations that teach through concrete cases.
- Progressive disclosure: how to layer information so that the reader is not overwhelmed.
- Concept scaffolding: how to build understanding from simple to complex.
- The self-study test: how to verify that documentation can teach without a live teacher.
- Writing for unknown readers: how to avoid assumptions about the reader's background.
- The relationship between documentation structure and learning effectiveness.

This article does not cover:

- The canonical article template (defined in the Stage 1 Meta-Framework).
- Format and storage requirements (governed by D6-003 and D9-001 R-D9-05).
- Curriculum sequencing across multiple articles (see D9-004).
- Tacit knowledge capture for integration into documentation (see D9-005).
- Assessment design (see D9-003).

## 3. Background

### 3.1 The Curse of Knowledge

The most pervasive obstacle to writing documentation that teaches is the curse of knowledge: once you know something, you cannot easily imagine not knowing it. The expert who writes a procedure unconsciously skips steps that seem obvious. They use terminology they have internalized without defining it. They describe a system state as "normal" without specifying what normal looks like. Every one of these gaps is a trap for the learner.

The curse of knowledge cannot be eliminated. It is a fundamental cognitive limitation. But it can be managed through discipline, technique, and testing. The techniques in this article -- worked examples, progressive disclosure, concept scaffolding -- are established instructional design methods for compensating for the curse of knowledge. The self-study test in Section 4.5 is the ultimate check: if someone unfamiliar with the content can learn from the documentation, the curse has been adequately managed.

### 3.2 Two Audiences, One Document

This institution's documentation serves two audiences simultaneously. The experienced operator needs a reference: quick to consult, organized by function, focused on the "what" and "how." The learner needs a teacher: patient in explanation, organized by concept, focused on the "why" and "how do I know it is working." These needs conflict. Reference documentation is terse. Teaching documentation is expansive. Reference documentation assumes context. Teaching documentation provides it.

The article template resolves this tension through section design. The Purpose and Background sections serve the learner by providing context and motivation. The System Model section serves both audiences by describing how things work. The Rules section serves the reference user by providing clear constraints. The Failure Modes and Recovery sections serve the learner by anticipating what goes wrong. This design is deliberate. It allows the experienced operator to skip to the sections they need while ensuring the learner encounters a complete educational sequence when reading front to back.

This article enhances that design by defining how to write within each section so that the educational function is maximized without compromising the reference function.

### 3.3 Why "Self-Teachable" Is the Standard

D9-001 established the self-teaching requirement (R-D9-01): every function must be learnable from documentation alone, to at least Stage 2 competence. This article defines how to write documentation that meets that standard. The standard is not "documentation exists." The standard is "documentation teaches." The difference is the difference between a parts list and an assembly manual. Both contain the same information. Only one enables the reader to build something.

## 4. System Model

### 4.1 Instructional Design Principles for Documentation

Six principles from cognitive science and instructional design govern how this institution writes documentation:

**Principle 1: Activate prior knowledge before introducing new content.** Begin each article, section, or procedure by connecting to something the reader already knows. The Purpose section does this at the article level by explaining why the topic matters in terms the reader can relate to. Within procedural sections, begin each step with a brief statement of what the step accomplishes before describing how to perform it.

Application: Instead of "Run `zpool scrub tank1`," write "To verify the integrity of all data on the storage pool, run `zpool scrub tank1`. This reads every block and compares it against its checksum, which will reveal any silent data corruption."

**Principle 2: Present concepts before procedures.** Understanding why a procedure works makes it easier to follow, easier to remember, and easier to troubleshoot when it fails. The Background and System Model sections accomplish this at the article level. Within procedural sections, each major procedure should be preceded by a brief explanation of what it accomplishes and why it works.

Application: Before documenting a backup verification procedure, explain what verification checks (checksum comparison, file count, restore test) and why each check matters (checksums detect corruption, file counts detect missing files, restore tests detect format or tooling problems).

**Principle 3: Manage cognitive load.** The reader's working memory is limited. Presenting too much information at once overwhelms the reader and impedes learning. Manage cognitive load through chunking (breaking information into manageable pieces), sequencing (presenting one concept at a time), and progressive disclosure (providing detail only when the reader is ready for it). See Section 4.3.

Application: A complex procedure should be broken into numbered phases, each with a clear purpose. The reader should be able to complete one phase, pause, and resume without losing context.

**Principle 4: Provide worked examples.** Abstract descriptions become concrete through examples. A worked example walks the reader through a specific case, showing every step with real (or realistic) data, expected outputs, and decision points. Worked examples are the single most effective instructional technique for procedural knowledge. See Section 4.2.

**Principle 5: Make the invisible visible.** Many institutional processes have states that are invisible to the reader: what a healthy system looks like versus an unhealthy one, what a correct output looks like versus an incorrect one, what "normal" means in a specific context. Documentation must make these states explicit. Use comparison (show both the correct and incorrect output), annotation (explain what each element of an output means), and visualization (use diagrams or structured layouts when text alone is insufficient).

Application: When a procedure produces output, show the expected output. When the procedure could produce error output, show that too, with an explanation of what each error means and what to do about it.

**Principle 6: Teach the failure modes alongside the success path.** The standard documentation pattern is: here is how to do the thing, and here (in a separate section far away) is what can go wrong. For learning purposes, failure information should be integrated with the success path. When a step could fail, say so immediately: "If this command returns an error code, see Section 7.2 before proceeding." The learner should never be surprised by a failure during execution. They should encounter the possibility of failure in the documentation before they encounter it in reality.

### 4.2 Worked Examples

A worked example is a complete, step-by-step demonstration of a procedure using specific, concrete data. It is the documentation equivalent of a live demonstration. The reader follows along, seeing exactly what happens at each step, what the expected outputs are, and where decision points occur.

**Structure of a worked example:**

1. **Setup:** State the initial conditions. What system state exists before the example begins? What has the reader already done? What tools are open? What data exists?

2. **Execution:** Walk through every step with specific values, not placeholders. Instead of "Enter the backup destination path," write "Enter `/mnt/backup/2026-02-16-full` as the backup destination path." Show the exact command or action, the exact output, and the exact next step.

3. **Checkpoints:** At key points, state what the reader should see and what it means. "The output should show `scrub repaired 0B in 0h42m.` The `0B` indicates no corruption was found. If you see a non-zero value here, stop and consult Section 7.3."

4. **Decision points:** Where the procedure branches, show both branches. "If the verification shows all checksums match, proceed to Step 7. If any checksums do not match, proceed to Section 7.4: Recovery from Verification Failure."

5. **Completion:** State the expected final state. What has changed? How can the reader verify that the procedure succeeded?

**When to include worked examples:**
- Every Tier 1 critical procedure must have at least one worked example.
- Any procedure with more than ten steps should have a worked example.
- Any procedure that involves judgment calls should have a worked example that demonstrates the judgment in context.
- Any procedure that is commonly performed incorrectly (based on assessment data or incident history) should have a worked example that explicitly demonstrates the correct approach and contrasts it with the common error.

### 4.3 Progressive Disclosure

Progressive disclosure is the practice of presenting information in layers, from general to specific, from simple to complex. The reader encounters the essential information first and can choose to go deeper as needed. This prevents information overload while ensuring completeness.

Progressive disclosure operates at three levels in this institution's documentation:

**Article level:** The section structure itself is a progressive disclosure mechanism. Purpose tells you why (broadest context). Background tells you what you need to know (contextual detail). System Model tells you how it works (technical detail). Rules tell you what constrains it (specific requirements). Failure Modes tell you what can go wrong (edge cases). Recovery tells you how to fix it (contingency detail). A reader who needs only the overview reads Purpose and Background. A reader who needs implementation reads System Model and Rules. A reader troubleshooting a problem reads Failure Modes and Recovery.

**Section level:** Within each section, begin with the most important statement and elaborate progressively. The first paragraph of each section should be a summary. Subsequent paragraphs add detail, context, and edge cases. A reader who reads only the first paragraph of each section should get a coherent, if abbreviated, understanding.

**Procedure level:** Complex procedures use a layered approach. The main procedure provides the standard path with sufficient detail for execution. Footnotes, side notes, or clearly labeled sub-sections provide additional detail for unusual circumstances, edge cases, or deeper explanation. The reader performing the procedure for the first time follows the main path. The reader encountering an anomaly consults the sub-sections.

### 4.4 Concept Scaffolding

Concept scaffolding is the practice of building complex understanding from simpler foundations, explicitly connecting each new concept to previously established ones. Scaffolding ensures that the reader never encounters a concept for which they lack the prerequisites.

**Scaffolding techniques:**

**Build-up definitions.** When introducing a complex concept, begin with a simple version and refine it progressively. Do not begin with the complete, precise definition. Begin with the intuitive version. Then add nuance. Then add edge cases. Example: "A backup is a copy of your data stored separately from the original. More precisely, a backup is a verified copy -- a copy whose integrity has been confirmed through checksum comparison. More precisely still, a full backup in this institution is a verified copy of all institutional data, stored on media in a separate physical location, following the rotation schedule defined in D6-002."

**Analogy before abstraction.** When explaining an abstract concept, provide an analogy to something familiar before presenting the formal definition. Then explicitly state where the analogy breaks down. Example: "The key inventory is like a library catalog for your keys. It tells you what keys exist, where each one is stored, and what each one opens. Unlike a library catalog, the key inventory does not contain the keys themselves -- it contains only information about the keys."

**Explicit connections.** When introducing a concept that builds on a previously taught concept, say so explicitly. "This procedure uses the backup verification steps you learned in Module 3. The difference is that here you are verifying a restored backup rather than a fresh one. The checksum comparison step is identical. The file count step adds a comparison against the source system's current file count."

**Prerequisite statements.** Every section or procedure that requires prior knowledge begins with a statement of what that knowledge is and where to find it. "Before performing this procedure, you must understand the key categories defined in SEC-003, Section 3.3. If you have not read that section, do so now."

### 4.5 The Self-Study Test

The self-study test is the definitive verification that a document meets the self-teaching requirement. It is simple in concept, demanding in execution, and irreplaceable as a quality check.

**Procedure:**

**Step 1: Select the document to test.** This can be a single article, a module from D9-004, or a procedure within an article.

**Step 2: Identify a tester.** The ideal tester is someone unfamiliar with the document's subject matter but possessing the general literacy and technical aptitude that D9-001 Section 3.3 assumes. If no such person is available, the author can self-test by setting the document aside for at least two weeks (to allow short-term memory of the writing process to decay) and then attempting to follow it as a naive reader.

**Step 3: The tester attempts to learn or perform the function using only the documentation.** No verbal assistance, no hints, no clarifications. The tester records every point of confusion, every ambiguity, every step where they did not know what to do or did the wrong thing.

**Step 4: Compile the confusion log.** Every confusion point is a documentation failure. Classify each as: (a) missing information (a step or concept not covered), (b) unclear writing (the information is present but confusing), (c) missing prerequisite (the document assumes knowledge the reader lacks), or (d) incorrect information (the document is wrong).

**Step 5: Revise the documentation.** Address every item in the confusion log. Missing information is added. Unclear writing is revised. Missing prerequisites are either taught in the document or added to the prerequisite declaration. Incorrect information is corrected.

**Step 6: Retest.** After revision, repeat the test. Continue the cycle until the tester can learn or perform the function from the documentation alone, achieving at least Stage 2 competence.

**Step 7: Record the test.** Document the test in the article's Commentary Section or in the curriculum record: date, tester identity, number of revision cycles, final result. A document that has passed the self-study test is marked as "Verified self-teachable" in the cross-training matrix per D9-006.

**Self-testing when no external tester is available:** Set the document aside for a minimum of two weeks. Then approach it as a reader, not an author. Follow every step literally. Where the document says "check the log," check the log -- do not rely on memory of what the log contains. Where the document says "verify the output," verify it -- do not assume you know the output because you wrote the document. Record every moment where you relied on knowledge not in the document. Each such moment is a documentation gap.

## 5. Rules & Constraints

- **R-D9-07-01:** Every article in the documentation corpus must be written to teach, not merely to record. This means every procedure must include the "why" alongside the "what," every system description must include worked examples, and every failure mode must include recovery guidance.
- **R-D9-07-02:** Every Tier 1 critical procedure must include at least one worked example with concrete data, expected outputs, checkpoints, and decision points.
- **R-D9-07-03:** The self-study test must be performed on every article that supports a Tier 1 function before the article is considered complete. Articles that have not passed the self-study test must be marked as "Not verified self-teachable" in the cross-training matrix.
- **R-D9-07-04:** Documentation must not assume knowledge that is not either taught in the document or explicitly listed as a prerequisite. Every assumption is either satisfied or declared.
- **R-D9-07-05:** When a document is revised, the revision must maintain or improve the document's educational quality. Revisions that add procedural detail must also add the explanatory context that makes the new detail understandable to a learner.
- **R-D9-07-06:** The six instructional design principles in Section 4.1 are mandatory standards for all documentation, not optional guidelines. Reviewers must check for compliance with these principles during the documentation review process.
- **R-D9-07-07:** Documentation must use progressive disclosure to serve both experienced operators and learners. Reference-mode access (skipping to the needed section) and learning-mode access (reading front to back) must both be supported by the document's structure.

## 6. Failure Modes

- **Expert-only documentation.** The author writes for themselves: dense, allusion-heavy, assuming context that only the author possesses. The documentation is an excellent reference for the author and an impenetrable wall for everyone else. Mitigation: the self-study test (Section 4.5) is the primary defense. R-D9-07-03 makes the test mandatory for critical procedures. Additionally, the instructional design principles provide writing-time guidance that reduces the severity of the problem before testing catches it.

- **Example-free abstraction.** The documentation explains concepts and procedures in abstract terms without ever showing a concrete instance. The reader understands the theory but cannot connect it to practice. Mitigation: R-D9-07-02 requires worked examples for critical procedures. The principle of "make the invisible visible" (Section 4.1, Principle 5) pushes authors toward concreteness.

- **Information avalanche.** The author, aware that documentation must be complete, includes every possible detail with no layering or organization. The reader is overwhelmed. Critical information is buried in a sea of edge cases. Mitigation: progressive disclosure (Section 4.3) structures information in layers. The first paragraph of each section summarizes. Detail comes after context. Edge cases come after the main path.

- **Stale worked examples.** The documentation is updated when procedures change, but the worked examples are not. The example shows the old procedure, confusing learners who see a discrepancy between the example and the current text. Mitigation: when a procedure is revised, every worked example that demonstrates that procedure must be revised simultaneously. The review checklist should include "Are all worked examples current?"

- **Self-study test avoidance.** The test is time-consuming and its results are uncomfortable (they reveal how bad the documentation is). The operator skips it, marking documents as "self-teachable" without verification. Mitigation: R-D9-07-03 makes the test mandatory for Tier 1 functions. The cross-training matrix (D9-006) tracks which documents have been verified. Documents without verification are visible gaps.

- **Scaffolding gaps.** The documentation uses terminology or concepts from other articles without teaching them or declaring them as prerequisites. The reader encounters unfamiliar terms and must hunt through other documents to understand the current one. Mitigation: R-D9-07-04 requires all assumptions to be satisfied or declared. The prerequisite mapping from D9-004 provides the framework for identifying cross-article dependencies.

## 7. Recovery Procedures

1. **If documentation has been written as reference-only (not self-teachable):** Do not rewrite everything at once. Prioritize by function criticality. For each Tier 1 function, add the following to the existing documentation: a "Why this matters" paragraph at the beginning of each section, at least one worked example per major procedure, expected outputs for every command or action, and explicit prerequisite statements. Then conduct the self-study test. Revise based on results.

2. **If worked examples are outdated:** Identify all worked examples in the documentation corpus. Cross-reference each with the current version of the procedure it demonstrates. Update any example that no longer matches the current procedure. If the procedure has changed significantly, write a new example from scratch rather than patching the old one.

3. **If the self-study test reveals pervasive documentation failures:** This is not a documentation problem; it is a strategic problem. The institution's primary educational vehicle is not functioning. Treat this as a Tier 1 operational priority. Allocate dedicated time -- not margins of other work, but primary operational time -- to documentation revision. Begin with the critical path: the documents a new operator would need to learn the most essential functions. Defer non-critical documentation improvement until the critical path is self-teachable.

4. **If no tester is available for the self-study test:** Use the self-testing procedure (Section 4.5). Accept its limitations. Supplement with the think-aloud capture from D9-005: perform the procedure while narrating, and note every moment where your narration includes information not in the documentation. Each such moment is a candidate for documentation improvement.

## 8. Evolution Path

- **Years 0-5:** The documentation is being written with teaching in mind from the start. The founder applies the instructional design principles as they write. Self-study tests are conducted using the self-testing method. The documentation improves iteratively. Expect to discover that early documentation is significantly less teachable than later documentation, as the author's skill in writing for learning develops through practice.

- **Years 5-15:** The documentation should be tested by its first real learner. This is the critical test. Every self-study gap that slipped through self-testing will be revealed. The revision effort following the first real learner test will likely be the largest single improvement in documentation quality the institution ever experiences.

- **Years 15-30:** The documentation has been through multiple revision cycles and real learner tests. The instructional design principles are well-embedded. The primary challenge is maintaining educational quality during routine updates -- ensuring that procedural revisions do not degrade the teaching elements that surround them.

- **Years 30-50+:** The documentation is a mature educational instrument. It has taught multiple learners. The worked examples reflect real institutional experience across decades. The progressive disclosure structures have been refined through use. The self-study test has been validated by its results. The documentation is no longer just a reference library with teaching elements bolted on. It is a teaching system that also functions as a reference library.

- **Signpost for revision:** If a motivated learner with general technical aptitude cannot achieve Stage 2 competence for a Tier 1 function from the documentation alone, this article's principles have not been adequately implemented and the documentation for that function must be revised.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
Writing documentation that teaches is a skill I am learning as I go. I have been a consumer of documentation for years, and I know the difference between documentation that helps me learn and documentation that merely records what the author already knows. The first kind is rare. I am trying to write the first kind, and I am discovering how difficult it is.

The hardest principle to follow is "make the invisible visible." When I write a procedure, the expected outputs are obvious to me. I wrote the system. I know what every command returns. But a reader who has never seen the system has no idea what to expect, and if the actual output differs from what the documentation describes -- even in trivial ways like version numbers or timestamps -- the reader does not know whether the difference matters. Showing the expected output, and explaining which parts are variable and which parts must match exactly, is tedious to write and essential to the reader.

The self-study test is the part of this article I believe in most strongly and dread most intensely. Setting my documentation aside for two weeks and then trying to follow it as a naive reader is humbling. I find gaps everywhere. Steps I thought were clear are ambiguous. Terms I thought were obvious are undefined. Procedures I thought were complete skip steps I perform unconsciously. Every gap found is a gap fixed, but the volume of gaps is a reminder of how wide the curse of knowledge really is.

I will say this: writing documentation that teaches has made me a better operator. The act of explaining why a procedure works forces me to understand it more deeply than merely knowing how to perform it. The act of writing worked examples forces me to actually execute the procedure and verify that it works as described. The act of anticipating a learner's confusion forces me to think about edge cases I might otherwise ignore. Documentation-as-teacher teaches the teacher first.

## 10. References

- ETH-001 -- Ethical Foundations of the Institution (Principle 3: Transparency of Operation; Principle 4: Longevity Over Novelty)
- CON-001 -- The Founding Mandate (documentation completeness, institutional transferability, fifty-year horizon)
- GOV-001 -- Authority Model (succession preparedness, documentation for future operators)
- OPS-001 -- Operations Philosophy (documentation-first principle, operational tempo)
- D6-001 -- Data Philosophy (Tier 1 classification for educational documentation)
- D6-003 -- Format Longevity Doctrine (format requirements for documentation and training materials)
- D9-001 -- Education Philosophy (self-teaching requirement R-D9-01, documentation-as-teacher principle R-D9-02, four-stage competence model, curse of knowledge acknowledgment)
- D9-002 -- The Competency Model (function classification, Tier 1 function identification)
- D9-003 -- Competency Assessment Framework (Stage 2 assessment as self-teachable standard)
- D9-004 -- Curriculum Design and Maintenance (sequencing across articles, prerequisite mapping)
- D9-005 -- Tacit Knowledge Capture Procedures (integrating tacit knowledge into documentation)
- D9-006 -- Cross-Training Procedures (self-teaching verification, documentation quality in cross-training matrix)
- Stage 1 Meta-Framework (META-00-ART-001), Section 3: Article Template (canonical section structure as educational design)
