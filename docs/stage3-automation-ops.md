# STAGE 3: AUTOMATION OPERATIONS

## Operational Doctrine for Domain 8 -- Automation & Agents

**Document ID:** STAGE3-AUTOMATION-OPS
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Classification:** Operational Procedures -- These articles translate the D8-001 Automation Restraint Doctrine into step-by-step actionable procedures. They are designed to be followed by a single operator with no one to ask.

---

## How to Read This Document

This document contains five operational articles for Domain 8 (Automation & Agents) of the holm.chat Documentation Institution. They implement the philosophy established in D8-001 -- the Automation Restraint Doctrine -- by providing the concrete procedures, templates, checklists, and workflows that make that philosophy operational.

D8-001 tells you why the institution approaches automation with skepticism, why every automation must justify its existence, and why comprehensibility outweighs convenience. These articles tell you how to design automation that meets those standards, how to oversee it during operation, how to audit it annually, how to shut it down in an emergency, and how to evaluate automation you did not build.

Read D8-001 first. These articles assume you have internalized its principles. They do not re-argue the case for restraint. They implement it.

These articles are written for the person doing the work. Every procedure is designed to be executable by a single operator following the documentation alone, at Stage 2 competence per D9-001. Where judgment is required, the criteria for that judgment are stated. Where a procedure can fail, the recovery path is documented.

If you are a future maintainer inheriting these procedures: the specific tools and commands may have changed by the time you read this. The principles behind each procedure are stated in the Background sections. Use those principles to adapt. Then update this document.

---

---

# D8-002 -- Agent Design Principles and Standards

**Document ID:** D8-002
**Domain:** 8 -- Automation & Agents
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, SEC-001, OPS-001, D8-001
**Depended Upon By:** D8-003, D8-004, D8-005, D8-006, D10-002. All articles involving the creation or modification of automation.

---

## 1. Purpose

This article defines the practical standards for designing automation that meets the institutional requirements established in D8-001. It translates the five-criterion justification test from philosophy into a working design process. It provides the templates, checklists, and documentation requirements that ensure every automation entering the institution is comprehensible, observable, reversible, and replaceable.

Without concrete design principles, the restraint doctrine becomes a sentiment that erodes under operational pressure. This article provides the structural reinforcement that prevents that erosion. Every automation in the institution -- from a three-line shell script to a multi-stage monitoring daemon -- must be designed according to these principles and documented according to these templates.

## 2. Scope

This article covers:

- The practical application of the five-criterion justification test from D8-001 Section 4.1.
- Design principles for automation that serves the institution rather than replacing human understanding.
- The Agent Specification Document (ASD) template: the required documentation for every automation.
- The code review checklist for automation approval.
- Naming conventions, file organization, and header requirements.
- Inter-automation dependency rules and testing requirements.

This article does not cover:

- The philosophical case for automation restraint (see D8-001).
- Human-in-the-loop oversight procedures (see D8-003).
- Annual re-justification procedures (see D8-004).
- Emergency shutdown procedures (see D8-005).
- Assessment of inherited automation (see D8-006).

## 3. Background

### 3.1 Why Design Standards Matter

An automation without design standards is a note to yourself. It works today, in this context, with your current understanding. It is not an institutional artifact. It is personal tooling that happens to run on institutional hardware.

Design standards transform personal tooling into institutional infrastructure -- documented well enough that someone else can understand it, tested well enough that it can be trusted, and structured well enough that it can be maintained across decades and across people. The overhead of following these standards is the price of building automation that belongs to the institution rather than to the individual who wrote it.

### 3.2 The Five-Criterion Test in Practice

D8-001 Section 4.1 defines five criteria: frequency, determinism, comprehensibility, observability, and reversibility. In practice, applying these criteria requires more than reading them and nodding. Each criterion demands specific evidence, and the evidence must be documented before the automation is built, not after.

The most common failure is building the automation first and justifying it second. The justification becomes a rationalization rather than an honest assessment. This article requires the justification to be completed and reviewed before any code is written. The Agent Specification Document is the vehicle for that justification.

### 3.3 The Comprehensibility Standard

Of the five criteria, comprehensibility is the hardest to evaluate honestly. The author of an automation always finds it comprehensible -- they wrote it. The standard is not "can the author understand it?" The standard is: "can a competent generalist, reading the documentation and source code, understand what this automation does, why it does it, what happens when it fails, and how to replace it with a manual process?"

A competent generalist, for this institution's purposes, means someone with general Unix/Linux system administration skills, the ability to read shell scripts and common scripting languages, familiarity with the institution's documentation system, and no prior exposure to this specific automation. If the automation requires specialist knowledge beyond that baseline, it is either too complex or its documentation is inadequate.

## 4. System Model

### 4.1 The Agent Specification Document (ASD)

Every automation must have an Agent Specification Document, created before the automation is built and maintained throughout its life. It is Tier 2 data per D6-001. The ASD template:

**Section A -- Identity.** Agent ID (per naming convention in Section 4.4), name, version, created date, author, last modified date, automation level per D8-001, and GOV-001 approval reference for Level 3/4 agents.

**Section B -- Justification.** Each of the five criteria addressed explicitly: Frequency (show the arithmetic -- automation effort must be significantly less than manual effort over a five-year horizon), Determinism (list every decision point and whether it can be made algorithmically), Comprehensibility (plain-language description that fits on one page), Observability (what is logged, where, in what format), and Reversibility (what changes are made, can each be undone, what safeguards exist for irreversible actions).

**Section C -- Technical Specification.** Inputs and their sources. Outputs including log entries. Dependencies with minimum versions. Execution environment and required permissions per SEC-002. Schedule or trigger mechanism. Resource requirements.

**Section D -- Failure Modes and Recovery.** Every identified failure mode with cause, symptoms, impact, and recovery. The "unknown failure" case. The manual fallback procedure -- a tested, step-by-step process for performing the task by hand.

**Section E -- Kill Procedure.** How to stop this automation immediately, including exact commands. Integration with D8-005. Post-kill state description and verification steps.

**Section F -- Change Log.** Dated record of every modification for the automation's entire life.

### 4.2 Design Principles

**Principle 1: Single Responsibility.** Each automation does one thing. Multiple distinct operations are implemented as multiple automations that can be run, tested, and killed independently.

**Principle 2: Explicit Configuration.** No magic numbers, no undocumented hardcoded paths. Every configurable value is defined in a configuration section or external file, documented in ASD Section C.

**Principle 3: Graceful Failure.** On error: stop, log the error in human-readable form, and leave the system in a recoverable state. No silent retries, no proceeding with partial data, no swallowed errors.

**Principle 4: Minimal Privilege.** Each automation runs with minimum required permissions under a dedicated service account per SEC-002 Category 4.

**Principle 5: No Undocumented Inter-Agent Dependencies.** If one automation depends on another, both ASDs document the dependency. The total dependency graph must be reviewable from ASDs alone.

**Principle 6: Human-Readable Logging.** Every log entry includes timestamp, agent ID, action, result, and sufficient context for understanding without consulting source code.

**Principle 7: Idempotency Where Possible.** Design so that re-running after a partial failure produces the correct result rather than compounding the failure.

### 4.3 The Code Review Checklist

Before deployment, every automation must pass this checklist, performed as a distinct, documented act and recorded in the ASD Change Log:

1. **ASD Completeness:** All sections filled, no blank or "TBD" entries.
2. **Justification Validity:** Five-criterion arithmetic is reasonable, claims are accurate.
3. **Source Code Review:** Code matches ASD description, error handling present for every fallible operation, no undocumented behaviors.
4. **Header Compliance:** Human-readable header per R-D8-09 stating purpose, creation date, author, rationale.
5. **Configuration Review:** All configurable values explicit per Principle 2.
6. **Privilege Review:** Minimum permissions only per Principle 4.
7. **Logging Review:** Test run produces logs sufficient to determine what happened without reading source code.
8. **Failure Testing:** Introduce a failure condition. Automation fails gracefully per Principle 3.
9. **Kill Test:** Execute ASD Section E kill procedure during execution. Automation stops. System state matches documentation.
10. **Manual Fallback Test:** Execute ASD Section D fallback. It achieves the same result.
11. **Dependency Verification:** All dependencies exist at required versions. No undocumented dependencies.
12. **Inter-Agent Review:** All cross-agent dependencies documented in both ASDs.

### 4.4 Naming and Organization

All automation files reside in a designated directory:

```
/institution/automation/
  agents/[agent-id]/         -- Source code, config/, tests/
  docs/[agent-id]-asd.md     -- Agent Specification Documents
  logs/[agent-id]/            -- Log output
  retired/                    -- Decommissioned agents (preserved)
```

Agent IDs follow: `agt-[domain]-[sequence]` (e.g., `agt-bkp-001` for the first backup agent).

## 5. Rules & Constraints

- **R-D8-02-01:** Every automation must have a complete ASD before deployment. Automation without an ASD must be halted until the ASD is completed.
- **R-D8-02-02:** Every automation must pass the code review checklist before deployment, documented in the ASD Change Log.
- **R-D8-02-03:** The ASD must be updated whenever the automation is modified. Discrepancies must be corrected before the next scheduled run.
- **R-D8-02-04:** No undocumented inter-agent dependencies. The dependency graph must be reconstructable from ASDs alone.
- **R-D8-02-05:** All source code must include the header required by D8-001 R-D8-09.
- **R-D8-02-06:** Manual fallback procedures must be tested at least annually or whenever the automation is modified.
- **R-D8-02-07:** Automation source code, configurations, and ASDs are Tier 2 data per D6-001.
- **R-D8-02-08:** Partial code reviews are not permitted. If time pressure demands deployment before full review, the automation operates at Level 1 or 2 until the full review is complete.

## 6. Failure Modes

- **ASD drift.** The automation evolves but the ASD is not updated. A future maintainer forms an incorrect mental model. Mitigation: R-D8-02-03 requires updates with every modification; D8-004 annual audit catches remaining drift.

- **Checklist theater.** The review is performed as formality. Mitigation: the checklist requires documented evidence (test output, kill test results, fallback execution). Evidence absence is detectable during annual audit.

- **Over-engineering.** Simple tasks acquire disproportionate documentation. Mitigation: proportionality. A five-line script's ASD should be one or two pages. The principles are a floor, not a target for elaboration.

- **Dependency sprawl.** Individual automations are simple but their dependency graph is incomprehensible. Mitigation: Principle 5 and annual dependency graph review per D8-004. If the graph cannot be drawn on a single page, the system is too complex.

- **Configuration drift across environments.** The automation works on current hardware but fails after migration because configuration was not fully externalized. Mitigation: Principle 2 requires explicit configuration. Code review checklist item 5 specifically checks for hardcoded values. Hardware transitions in years 15-30 will test this.

- **Stale manual fallbacks.** The manual fallback procedure was written when the automation was deployed and never tested again. When the automation fails years later, the fallback references tools or paths that no longer exist. Mitigation: R-D8-02-06 requires annual testing or testing on modification, whichever is sooner.

## 7. Recovery Procedures

1. **If an ASD is out of date:** Halt the automation. Update the ASD from source code and current behavior. Run the code review checklist. Resume only after the ASD is accurate.

2. **If automation is deployed without an ASD:** Halt immediately. Create the ASD retroactively. Run the five-criterion test honestly. Redeploy only if it passes; decommission per D8-004 if it fails.

3. **If undocumented inter-agent dependencies are discovered:** Map all dependencies. Update all affected ASDs. Assess graph comprehensibility. Break excessive dependencies even at the cost of efficiency.

4. **If over-engineering blocks useful automation:** Review ASDs for proportionality. Simplify documentation for simple automations while maintaining all required sections.

## 8. Evolution Path

- **Years 0-5:** Expect friction as standards are first applied. Build the habit of writing the ASD before writing the code. This habit, once established, is the institution's strongest defense against automation creep.

- **Years 5-15:** The agent inventory should be small and well-documented. The main temptation is "just a quick script." Every quick script that skips the ASD is a future black box.

- **Years 15-30:** Hardware transitions test portability. Well-documented configurations make migration straightforward. Poor documentation makes migration a rewrite.

- **Years 30-50+:** Future maintainers use ASDs as their primary guide. ASD quality from years 0-15 determines whether the automation estate is an asset or a liability.

- **Signpost for revision:** If the ASD template consistently forces blank sections, the template needs revision. If the code review checklist consistently passes automation that later causes problems, the checklist is inadequate.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The hardest part of writing design standards for a one-person institution is the self-review problem. I am the author, the reviewer, and the approver. But the checklist exists not only for me -- it exists for the person who will review my work after I am gone. The ASD is a letter to the future, explaining what I built and why.

I have kept the ASD template structured but flexible in length. A five-line script that cleans up temporary files does not need a ten-page specification. But it does need a one-page specification: what it does, why it exists, what it deletes, and how to do the same thing by hand. The floor is non-negotiable. The ceiling is the author's judgment.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 2: Integrity Over Convenience; Principle 3: Transparency)
- CON-001 -- The Founding Mandate (comprehensibility requirement)
- GOV-001 -- Authority Model (Tier 3 approval for Level 3/4 automation)
- SEC-001 -- Threat Model (Pillar 3: The Human Is the System)
- SEC-002 -- Access Control Procedures (service accounts, privilege model)
- OPS-001 -- Operations Philosophy (documentation-first principle, complexity budget)
- D6-001 -- Data Philosophy (data tier classification for automation artifacts)
- D8-001 -- Automation Philosophy: The Restraint Doctrine (five-criterion test, automation spectrum)
- D8-003 -- Human-in-the-Loop Procedures (oversight requirements)
- D8-004 -- Automation Audit and Annual Re-justification (annual review of ASDs)
- D8-005 -- Emergency Override and Kill Switch Procedures (kill procedures)
- D8-006 -- Legacy Automation Assessment (assessment against these standards)
- D9-001 -- Education Philosophy (self-teaching requirement)
- Stage 1 Documentation Framework, Domain 8: Automation & Agents

---

---

# D8-003 -- Human-in-the-Loop Procedures

**Document ID:** D8-003
**Domain:** 8 -- Automation & Agents
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, SEC-001, OPS-001, D8-001, D8-002
**Depended Upon By:** D8-004, D8-005, D8-006, D6-001 (R-D6-06), D10-002. All articles involving automation oversight or automated data modification.

---

## 1. Purpose

This article defines the procedures for human oversight of automation. It specifies how the operator monitors, approves, intervenes in, and overrides automated processes. It establishes the approval workflows, monitoring requirements, intervention procedures, escalation matrix, and override protocols that ensure no automation operates beyond the boundaries of human understanding and human authority.

D8-001 establishes that automation is a servant, never a master. This article operationalizes that principle. In a single-operator institution, "human-in-the-loop" means one person who has other responsibilities and cannot watch automation continuously. These procedures define when to look, what to look for, and what to do when something is wrong.

## 2. Scope

This article covers:

- Oversight requirements for each automation level defined in D8-001 Section 4.2.
- Approval workflows for automation requiring human authorization.
- Monitoring requirements: what to monitor, when, and how.
- The escalation matrix for classifying anomalies.
- Override procedures for superseding automation decisions.
- Monitoring dashboard design principles.

This article does not cover:

- Automation design (see D8-002).
- Annual re-justification (see D8-004).
- Emergency shutdown (see D8-005).
- Legacy automation assessment (see D8-006).

## 3. Background

### 3.1 The Single-Operator Oversight Problem

In a staffed operations center, someone is always watching. In a single-operator institution, nobody is always watching. The operator sleeps, eats, works on other tasks, and leaves the premises. Automation that runs during these periods runs without real-time human oversight.

This is not a violation of the human-in-the-loop principle. The principle does not require continuous real-time monitoring of every automation. It requires that the human is never removed from the oversight process -- that the operator reviews what automation has done, approves what it proposes to do, and retains the ability to intervene at any time. The distinction is between synchronous oversight (approval before action) and asynchronous oversight (review after action). Both are valid. The choice depends on the automation level and the consequences of the actions involved.

### 3.2 Oversight by Automation Level

- Level 0 (Fully Manual) and Level 1 (Assisted): oversight is inherent in the process.
- Level 2 (Supervised Autonomous): the automation proposes, the human approves. Synchronous oversight.
- Level 3 (Monitored Autonomous): the automation acts, the human reviews afterward. Asynchronous oversight.
- Level 4 (Fully Autonomous): the automation self-monitors, the human is notified only of anomalies. Exception-based oversight.

### 3.3 Structured Versus Unstructured Oversight

Unstructured oversight means glancing at logs when you remember to, approving requests without reading them, and assuming everything is fine because no alarms have sounded. Structured oversight means checking specific things at specific times, approving requests only after reviewing specific information, and verifying outcomes against specific expectations. The difference is the difference between a pilot's casual awareness and a pilot's checklist. This article provides the checklist.

## 4. System Model

### 4.1 Level 2 Approval Workflow

**The Approval Sequence:**

Step 1: The automation writes a proposal to the approval queue containing: agent ID, timestamp, plain-language description of the proposed action, affected data/systems, whether the action is reversible, and estimated execution time.

Step 2: The automation enters a wait state. If it has a defined timeout, it logs a timeout event and halts cleanly if the operator does not respond.

Step 3: The operator reviews the proposal, confirming: the action matches documented purpose, affected systems are correct, no unexpected changes are proposed, and timing is appropriate.

Step 4: The operator approves (automation proceeds, approval logged), modifies (parameters adjusted, both logged), or rejects (automation halts, reason logged).

Step 5: After execution, the automation logs results for the next review cycle.

**Implementation:** The approval queue is a directory of plain-text files or entries in a structured log. The mechanism is intentionally simple -- no database, no web interface. The operator checks the queue as part of daily operations per D10-002.

### 4.2 Level 3 Monitoring Requirements

**Daily Review:** As part of D10-002, review logs of all Level 3 automation from the previous 24 hours. Check: did it run as scheduled, did it complete successfully, were there errors or anomalies, are outputs consistent with expectations, was execution time normal?

**Weekly Summary:** Examine trends: increasing error rates, growing execution times, unfamiliar log patterns, missed schedule intervals.

**Monthly Verification:** Independently verify output by manually performing a sample of automated tasks and comparing results. This catches silent failure -- the automation reporting success while producing wrong results.

**The Monitoring Dashboard:** A plain-text document showing each active agent's ID, level, last run time and result, next scheduled run, current status, error counts (7-day and 30-day), and date of last manual verification. Updated automatically by a Level 3 agent (itself subject to these monitoring requirements) or maintained manually. Reviewed as the first daily operations item.

### 4.3 Level 4 Exception Handling

Level 4 automation must generate alerts for: any execution failure, results outside expected range, resource consumption exceeding thresholds, unavailable dependencies, excessive execution duration, and inability to write to its own log. Alerts are delivered through a mechanism independent of the automation itself. The alert mechanism is tested monthly.

**Additional Level 4 Audit Requirements:** Quarterly deep log review (full logs, not just exceptions). Semi-annual manual execution of the automated task. Annual re-justification per D8-004 with heightened scrutiny.

### 4.4 The Escalation Matrix

**Severity 1 -- Informational.** Unusual but harmless log entry. Response: note in operational log, monitor for recurrence.

**Severity 2 -- Degraded.** Task completed with errors or partial failures. Response: investigate within 24 hours. Pause automation if the issue is persistent.

**Severity 3 -- Failed.** Task did not complete. Response: investigate immediately. Execute manual fallback. Repair before next scheduled run or keep paused and operate manually.

**Severity 4 -- Dangerous.** Automation took an incorrect action causing or risking harm. Response: halt immediately per ASD Section E or D8-005. Assess and reverse damage. Treat as a security incident per D10-003 if Tier 1 data or security configurations are affected.

### 4.5 Override Procedures

The operator may override any automation at any time:

- **Parameter Override:** Changing inputs for a specific run. Both original and override logged.
- **Skip Override:** Preventing the next scheduled run. Skip and reason logged.
- **Result Override:** Discarding output and substituting a manual result. Both results logged.
- **Behavioral Override:** Changing automation logic. Requires source code and ASD update per D8-002. If it persists beyond one cycle, it is a modification, not an override.

Every override is recorded in the operational log with date, agent ID, type, reason, operator identity, and expected return to normal. Overrides are reviewed weekly and during annual audit per D8-004.

## 5. Rules & Constraints

- **R-D8-03-01:** Level 2 automation must not proceed without documented operator approval including identity and timestamp.
- **R-D8-03-02:** Level 3 logs must be reviewed daily per D10-002. Missed reviews must be documented and the missed logs reviewed within 48 hours.
- **R-D8-03-03:** All Level 3 and Level 4 automation must be independently verified through manual execution at least monthly.
- **R-D8-03-04:** Level 4 alert mechanisms must be tested monthly. Failed alert tests result in downgrade to Level 3 until repaired.
- **R-D8-03-05:** No automation at any level may modify or delete Tier 1 data without Level 2 synchronous approval. This overrides the automation's normal level.
- **R-D8-03-06:** Every override must be documented. Undocumented overrides are operational discipline violations.
- **R-D8-03-07:** The monitoring dashboard must be maintained and reviewed daily.
- **R-D8-03-08:** Severity 3 or 4 anomalies must not be downgraded without documented justification.

## 6. Failure Modes

- **Approval fatigue.** The operator rubber-stamps Level 2 approvals because they are always the same. Mitigation: clear proposal formatting; periodic variation testing (deliberately submitting unusual parameters to test attention).

- **Monitoring decay.** Daily review becomes a glance. Mitigation: monthly manual verification catches divergence; weekly trend review catches accumulating patterns.

- **Dashboard blindness.** "All green" for so long that changes are not noticed. Mitigation: daily review requires writing the dashboard state in the operational log -- a deliberate act distinct from not looking.

- **Bootstrap paradox.** The monitoring dashboard automation itself fails. Mitigation: daily review includes a manual check that the dashboard timestamp is current.

- **Override normalization.** Overrides become routine instead of exceptional. Mitigation: weekly override review. An automation requiring frequent overrides needs repair, level reduction, or decommissioning.

- **Alert flood.** Excessive Level 4 alerts cause the operator to stop reading them. The signal is lost in the noise. Mitigation: if an automation generates more than three alerts per week under normal conditions, its alert thresholds must be recalibrated. Alerts must be meaningful. An alert that fires constantly is worse than no alert.

- **Oversight as overhead.** The operator begins to see monitoring as bureaucratic overhead rather than a necessary function. Corners are cut. Reviews are skipped. The oversight structure becomes a fiction maintained in documentation but not in practice. Mitigation: the procedures are designed for efficiency -- twenty minutes per day, not hours. If they take more than thirty minutes, the process or the estate needs streamlining, not the oversight.

## 7. Recovery Procedures

1. **If Level 2 approvals have been rubber-stamped:** Stop all Level 2 automation. Review the last 30 days of approvals for incorrect grants. Perform one genuine review cycle for each. Improve proposal formats if the operator cannot distinguish normal from abnormal.

2. **If daily monitoring has lapsed:** Review all missed logs immediately. Manually verify all Level 3/4 automation for the lapsed period. Document the lapse and address the root cause.

3. **If the monitoring dashboard has failed silently:** Manually review all logs for the period since the last dashboard update. Repair the dashboard. Improve its monitoring.

4. **If alert fatigue has developed:** Audit all alerts. Recalibrate or remove those that are not useful. Target zero false alarms, not zero alarms.

5. **If override normalization has occurred:** List overrides from the last 90 days. Any agent overridden more than three times must be repaired, reduced in level, or decommissioned.

## 8. Evolution Path

- **Years 0-5:** Tune the balance between thoroughness and practicality. Too thorough and monitoring costs more than automation saves. Too light and oversight is meaningless.

- **Years 5-15:** Accumulated baseline data -- what "normal" looks like for each agent -- makes anomaly detection meaningful. Without baselines, every observation is just a number.

- **Years 15-30:** If succession occurs, the new operator must learn not just procedures but baselines. Consider adding baseline profiles to ASDs.

- **Years 30-50+:** The human-in-the-loop principle endures even as technology makes deeper automation tempting. Tools may change. The principle does not.

- **Signpost for revision:** If daily monitoring consistently exceeds 30 minutes, the process needs streamlining or the automation estate needs reduction.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The honest challenge of human-in-the-loop in a single-operator institution is that the human has limited attention and unlimited responsibilities. These procedures are designed for a person with twenty minutes a day for automation monitoring, not a person with nothing else to do.

The most important safeguard is not the daily log review. It is the monthly manual verification. That is where you discover the automation has been producing subtly wrong results for three weeks while the logs showed only success. Logs tell you what the automation thinks happened. Manual verification tells you what actually happened.

I have deliberately made the approval workflow low-technology -- text files, not a web application. The approval mechanism should be simpler and more reliable than the automation it oversees. If the approval system is itself complex and fragile, you have traded one problem for two.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 2: Integrity Over Convenience; Principle 3: Transparency)
- CON-001 -- The Founding Mandate (human authority over institutional systems)
- GOV-001 -- Authority Model (Tier 3 approval for Level 3/4 automation)
- SEC-001 -- Threat Model (Pillar 3: The Human Is the System)
- OPS-001 -- Operations Philosophy (daily/weekly/monthly review cycles)
- D6-001 -- Data Philosophy (R-D6-06: no automated deletion of Tier 1/2 data without human authorization)
- D8-001 -- Automation Philosophy: The Restraint Doctrine (automation spectrum, servant principle)
- D8-002 -- Agent Design Principles and Standards (ASD template, kill procedures, manual fallback)
- D8-004 -- Automation Audit and Annual Re-justification (annual review of overrides)
- D8-005 -- Emergency Override and Kill Switch Procedures (emergency response beyond routine intervention)
- D10-002 -- Daily Operations Doctrine (integration of monitoring into daily checklist)
- D10-003 -- Incident Response Procedures (Severity 4 escalation)
- Stage 1 Documentation Framework, Domain 8: Automation & Agents

---

---

# D8-004 -- Automation Audit and Annual Re-justification

**Document ID:** D8-004
**Domain:** 8 -- Automation & Agents
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, SEC-001, OPS-001, D8-001, D8-002, D8-003
**Depended Upon By:** D8-005, D8-006, D10-002. Referenced by all articles defining automation subject to annual review.

---

## 1. Purpose

This article defines the procedures for the annual review of every active automation in the institution. D8-001 R-D8-07 mandates that the total inventory be reviewed at least annually and each automation re-justified against the five criteria. This article provides the detailed procedures, worksheets, and decision frameworks for conducting that review.

The annual re-justification exists because the world changes. A task performed twenty times a day when the automation was built may now happen twice a week. A dependency that was stable may have become fragile. Without periodic re-justification, the automation inventory grows monotonically -- automations are added but never removed, and the institution gradually depends on a layer of automation no one has questioned in years.

The annual audit is also a documentation integrity check. ASDs drift from reality. Logs reveal patterns invisible in daily review. Manual fallbacks may have gone stale. The annual audit creates the structured occasion for catching all of these.

## 2. Scope

This article covers:

- The annual automation audit process: schedule, scope, and responsibilities.
- The re-justification worksheet for evaluating each automation.
- Criteria for continued operation, modification, level reduction, and decommissioning.
- The decommissioning procedure for automation that fails review.
- Triggered reviews outside the annual cycle.

This article does not cover:

- Design of new automation (see D8-002).
- Daily or weekly monitoring (see D8-003).
- Emergency shutdown (see D8-005).
- Assessment of inherited automation (see D8-006).

## 3. Background

### 3.1 Why Annual Review Is Necessary

Automation has a natural tendency to persist. Once running, an automation becomes part of the operational background. It is easier to leave it running than to question whether it should still run. The operator adapts to its presence, routes processes around it, and forgets what the institution looked like before it existed. Over years, this accumulation creates a layer of automation that functions as the institution's nervous system -- and like a nervous system, it is poorly understood by the organism it serves.

The annual review forces a structured confrontation with every automation. Not "is it running?" but "should it still be running?" Not "does it work?" but "does it still serve the institution better than the alternative?" These are questions that never arise naturally. They must be imposed by procedure.

### 3.2 The Re-justification Standard

The re-justification standard is the same as the original justification standard: the five-criterion test from D8-001 Section 4.1. An automation that met all five criteria when it was built may no longer meet them. Frequency changes. Determinism assumptions prove wrong. Comprehensibility degrades as the original author's memory fades. Observability logs accumulate without review. Reversibility guarantees are invalidated by downstream changes.

Re-justification is not a rubber stamp. It is a genuine re-application of the five criteria to the automation as it exists today, in the operational context that exists today, judged by the operator who exists today. If the operator cannot explain what an automation does -- even one they built years ago -- it fails the comprehensibility criterion.

### 3.3 The Decommissioning Default

The default outcome of a failed re-justification is decommissioning, not remediation. This is deliberate. The burden of proof lies with the automation, not with the manual process it replaced. If an automation cannot justify its continued existence, the institution reverts to the manual process. The operator may choose to remediate and re-submit the automation for review, but remediation is a new action, not a continuation of the failed review. This prevents the annual review from becoming an exercise in finding reasons to keep automation rather than genuinely evaluating it.

## 4. System Model

### 4.1 The Annual Audit Process

The audit is conducted within a two-week window as part of the annual operations cycle per OPS-001:

**Phase 1: Inventory Verification (Day 1-2).** Generate the current inventory. Compare against the ASD directory and monitoring dashboard. Flag discrepancies: active automation without ASDs, ASDs without corresponding automation, unmonitored automation. Resolve all discrepancies before proceeding.

**Phase 2: Individual Re-justification (Day 3-10).** For each active automation, complete the Re-justification Worksheet (Section 4.2). Review in order of automation level, starting with Level 4.

**Phase 3: Decision and Action (Day 10-14).** Classify each automation:

- **Continue:** Passes all five criteria. ASD is accurate. Monitoring adequate. Manual fallback tested within the year.
- **Modify:** Passes criteria but requires changes (documentation, parameters, logging, fallback). Modifications completed within 30 days.
- **Reduce Level:** Risk profile has changed. Move to a lower automation level with additional oversight.
- **Decommission:** Fails one or more criteria. Invoke Section 4.5.

**Phase 4: Documentation (Day 14).** Record results in the decision log. File worksheets as Tier 2 data.

### 4.2 The Re-justification Worksheet

One worksheet per active automation. Completed during the annual audit:

**Header:** Agent ID, name, level, original deployment date, last review date, this review date, reviewer.

**Frequency Re-evaluation:** Current task frequency. Has it changed? Current annual time savings versus annual maintenance cost. Does savings still significantly exceed cost? Show arithmetic. Verdict: PASS / FAIL / MARGINAL.

**Determinism Re-evaluation:** New edge cases discovered? Situations requiring human intervention? Override count from D8-003 logs for the past year and reasons. Verdict.

**Comprehensibility Re-evaluation:** Can the reviewer explain the automation from memory or ASD alone? Does ASD match current source code? Are there undocumented changes? Would a generalist understand it? Verdict.

**Observability Re-evaluation:** Are recent logs human-readable and sufficiently detailed? Can the reviewer determine what happened on a specific date from logs alone? Any logging gaps? Verdict.

**Reversibility Re-evaluation:** Irreversible actions in the past year? Were they authorized per D8-003? Has the reversibility profile changed? Verdict.

**Operational Health:** Error rate and trend. Manual verification discrepancies. Dependency status. Resource consumption changes. Security posture.

**Manual Fallback Status:** Was it executed this year? When last tested? Does it reference tools that still exist?

**Overall Assessment:** CONTINUE / MODIFY / REDUCE LEVEL / DECOMMISSION, with justification.

### 4.3 Triggered Reviews

Events requiring re-justification outside the annual cycle:

- After a Severity 3 or 4 incident: review within 7 days.
- After significant source code modification: review affected criteria within 14 days.
- After a dependency change: review within 30 days.
- After a succession event: all automation reviewed within 60 days per D8-006.
- If the operator cannot explain what an automation does: immediate review; automation paused until complete.

### 4.4 The Automation Health Summary

After the annual audit, the operator produces a one-page summary: total active automations by level, review outcomes, total time saved versus maintenance time, net benefit, key concerns, and comparison with the previous year. This summary is Tier 2 data and becomes a longitudinal record.

### 4.5 The Decommissioning Procedure

Step 1: **Schedule.** Record the decision. Set a date within 30 days. Ensure the manual fallback is current and tested.

Step 2: **Disable.** On the date, disable scheduled execution. Do not delete files. Remove from monitoring dashboard.

Step 3: **Observe.** For 30 days, perform the task manually and monitor for unexpected effects of the automation's absence.

Step 4: **Archive.** After the observation period, move source code, configuration, and ASD to the retired directory. Files are preserved as Tier 2 data, clearly marked as retired.

Step 5: **Close.** Record completion in the decision log. Update the Health Summary.

If unexpected effects are detected during observation, the operator may reactivate temporarily, modify and re-submit, or address the dependency directly.

## 5. Rules & Constraints

- **R-D8-04-01:** Every active automation must be re-justified annually. No exemptions.
- **R-D8-04-02:** The audit must complete within a two-week window. Results recorded in the decision log.
- **R-D8-04-03:** A FAIL verdict on any criterion prevents a CONTINUE classification.
- **R-D8-04-04:** Decommissioning is the default for failed reviews. MODIFY requires a 30-day deadline; missed deadline triggers decommissioning.
- **R-D8-04-05:** Triggered reviews must occur within specified timeframes. Missed deadlines result in automation being paused.
- **R-D8-04-06:** The Automation Health Summary must be produced after each audit and stored as Tier 2 data.
- **R-D8-04-07:** Completed worksheets are preserved for at least 10 years as Tier 2 data.
- **R-D8-04-08:** The decommissioning procedure must be followed completely. Do not skip the 30-day observation period or delete files instead of archiving.

## 6. Failure Modes

- **Rubber-stamp reviews.** Every automation gets CONTINUE because questioning is harder than approving. Mitigation: worksheets require specific evidence. "PASS" on frequency requires arithmetic. "PASS" on comprehensibility requires explanation.

- **Audit fatigue.** A large inventory makes the audit overwhelming. Mitigation: keep the inventory small via D8-001 restraint. If the inventory cannot be audited in two weeks, it is too large.

- **Decommissioning paralysis.** Fear of breaking something prevents decommissioning. Mitigation: the 30-day observation period addresses this fear. The automation is disabled, not deleted.

- **Modification drift.** MODIFY verdicts are scheduled and forgotten. The automation continues to run in its unmodified state indefinitely. Mitigation: R-D8-04-04 imposes a 30-day deadline. Modifications not completed within the deadline trigger decommissioning. This is a hard consequence, not a suggestion.

- **Triggered review neglect.** A triggering event occurs but the operator does not conduct the triggered review because they are busy with other tasks. The automation continues running in a potentially changed context without reassessment. Mitigation: R-D8-04-05 pauses the automation if the triggered review is not conducted within the specified timeframe. The pause is the forcing function.

## 7. Recovery Procedures

1. **If the annual audit has been skipped:** Conduct it immediately, prioritizing Level 4 agents. If more than three months overdue, document this as a serious operational failure with root cause analysis.

2. **If rubber-stamp reviews are suspected:** Re-review three randomly selected automations. If new verdicts differ, re-audit all automation with genuine scrutiny.

3. **If a decommissioned automation is needed:** Reactivate from the retired directory. Conduct a full re-justification. Restore only if it passes; otherwise find an alternative approach.

4. **If the inventory has grown beyond auditable size:** Declare an automation reduction sprint. Decommission ruthlessly, starting with weakest justifications.

## 8. Evolution Path

- **Years 0-5:** First audits establish rhythm with a small inventory. Use early audits to refine the worksheet template.

- **Years 5-15:** Longitudinal data tells a story. If nothing has ever been decommissioned, the review may be too lenient.

- **Years 15-30:** A new operator's first annual audit is a critical learning event, testing the comprehensibility criterion honestly for the first time since the original author left.

- **Years 30-50+:** Accumulated worksheets and summaries become institutional history of automation decisions.

- **Signpost for revision:** If the audit consistently takes less than two days (insufficient scrutiny) or more than three weeks (unwieldy inventory), adjustment is needed.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The annual re-justification is the enforcement mechanism for the restraint doctrine. Without it, D8-001 is philosophy. With it, D8-001 is practice.

I expect the frequency re-evaluation to be the most revealing. Tasks change frequency. The automation that saved an hour a day last year may save ten minutes this year because the underlying process simplified. Without annual review, I would never notice. The automation would continue consuming maintenance attention and complexity budget, solving a problem that has mostly solved itself.

The decommissioning procedure's 30-day observation period is a safety net for courage. The hardest part is the fear that something will break. The observation period makes that fear manageable: you are not deleting, you are pausing. If you are wrong, the path back is short.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 2: Integrity Over Convenience; Principle 4: Longevity Over Novelty)
- CON-001 -- The Founding Mandate (institutional mission, comprehensibility)
- GOV-001 -- Authority Model (Tier 3 approval for Level 3/4 automation, decision log)
- OPS-001 -- Operations Philosophy (annual operations cycle, complexity budget)
- D6-001 -- Data Philosophy (Tier 2 classification for audit documents)
- D8-001 -- Automation Philosophy: The Restraint Doctrine (R-D8-07 annual review mandate, five-criterion test)
- D8-002 -- Agent Design Principles and Standards (ASD template, code review checklist)
- D8-003 -- Human-in-the-Loop Procedures (override log, monitoring data, escalation matrix)
- D8-005 -- Emergency Override and Kill Switch Procedures (emergency decommissioning)
- D8-006 -- Legacy Automation Assessment (post-succession automation review)
- D10-002 -- Daily Operations Doctrine (integration into operations cycle)
- Stage 1 Documentation Framework, Domain 8: Automation & Agents

---

---

# D8-005 -- Emergency Override and Kill Switch Procedures

**Document ID:** D8-005
**Domain:** 8 -- Automation & Agents
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, SEC-001, OPS-001, D8-001, D8-002, D8-003
**Depended Upon By:** D8-004, D8-006, D10-002, D10-003. All articles involving emergency response to automation failures.

---

## 1. Purpose

This article defines the procedures for immediately halting any or all automation in the institution. It specifies the physical and software kill switches, the emergency shutdown sequence, the post-stop checks, and the restart procedures.

D8-001 establishes the servant principle: the operator must always be able to stop any automation at any time, immediately, without consequences that cannot be remedied. This article makes that principle operational. When a runaway automation is filling a disk, or a cascading failure is propagating through multiple agents faster than the operator can diagnose, this article tells them exactly what to do.

Emergency overrides are the last line of defense. The procedures are designed for speed and certainty, sacrificing diagnostic thoroughness for immediacy. Understand what went wrong after the automation is stopped, not while it is still running.

## 2. Scope

This article covers:

- Physical kill switches: hardware-level mechanisms to halt automation.
- Software kill switches: command-level mechanisms for specific or total automation halt.
- The emergency shutdown sequence.
- Post-stop checks and damage assessment.
- Restart procedures.
- Kill switch testing requirements.

This article does not cover:

- Routine intervention (see D8-003).
- Planned decommissioning (see D8-004).
- General system emergency procedures (see D10-003).

## 3. Background

### 3.1 Why Kill Switches Must Be Pre-Planned

In an emergency, the operator does not have time to read documentation or reason about dependencies. They need a procedure they can execute from memory or from a laminated card on the wall. These procedures are designed in advance, tested regularly, and maintained as the automation estate changes.

### 3.2 Physical Versus Software Kill Switches

Software kill switches are precise, auditable, and can target specific automations. But they can fail: the process may not respond to signals, the command may hang, the system may be overloaded. Physical kill switches -- disconnecting power, unplugging storage -- are crude but reliable. Both levels must be documented, maintained, and tested.

### 3.3 The Cascading Failure Scenario

The most dangerous automation emergency is cascading failure: one automation's failure triggers another's error handler, which makes incorrect decisions, triggering a third, and so on. The speed exceeds human comprehension. The correct response is not diagnosis -- it is immediate total shutdown, followed by diagnosis at leisure.

## 4. System Model

### 4.1 Kill Switch Architecture

Three levels, each more aggressive:

**Level 1: Individual Agent Kill.** Halts a single automation.
1. Identify the agent ID.
2. Execute the kill command from ASD Section E (typically `systemctl stop [service-name]`, `kill $(cat [pidfile])`, or cron removal).
3. Verify: check process list, check log for shutdown entry, confirm PID file removal.
4. Record in operational log.

**Level 2: Category Kill.** Halts all automation in a functional category (e.g., all backup agents).
1. Execute `/institution/automation/kill-category.sh [category]` -- reads the inventory, executes individual kills in sequence.
2. If the script fails, fall back to individual Level 1 kills for each agent in the category.
3. Verify all agents stopped. Record in operational log.

**Level 3: Total Kill.** Halts all automation.

*Software Total Kill:*
1. Execute `/institution/automation/kill-all.sh` -- disables all cron jobs and systemd timers, sends SIGTERM to all automation processes (identified by service accounts per SEC-002), waits 10 seconds, sends SIGKILL to remainders, logs every action.
2. If the script fails, execute manually:
   - Disable cron: `for user in $(cat /institution/automation/service-accounts.list); do crontab -r -u $user; done`
   - Disable systemd: `systemctl stop institution-automation.target`
   - Kill processes: `for user in $(cat /institution/automation/service-accounts.list); do pkill -u $user; done`
3. Verify: `ps -u [each service account]` returns no processes.

*Physical Total Kill (last resort):*
If software methods fail or the system is unresponsive:
1. Disconnect power to the automation server (or primary server if shared). This is a hard shutdown affecting all services.
2. If UPS is present, press the UPS power button or disconnect the machine from UPS.
3. Physical kills risk data corruption and filesystem damage. Use only when software methods have failed and automation is causing active, ongoing harm exceeding the harm of a hard shutdown.

### 4.2 The Emergency Shutdown Sequence

Step 1: **Stop the harm.** If you know which automation is causing the problem, Level 1 kill immediately.

Step 2: **Assess scope.** Confined to one agent, or spreading? If unsure or spreading, proceed to Level 3 total kill immediately. Err on the side of stopping too much.

Step 3: **Execute the kill.** Level 1, 2, or 3 matching your uncertainty. When in doubt, Level 3.

Step 4: **Verify.** Confirm targeted automation has halted. Check process lists, logs, and that harmful behavior has stopped.

Step 5: **Secure.** Do not restart anything. Disable scheduled triggers to prevent automatic restart. Record the time and reason.

Step 6: **Proceed to post-stop checks.**

### 4.3 Post-Stop Checks

**Data Integrity.** Examine data stores the automation had write access to. Compare against last known-good state (recent backup or integrity verification per D6-007). Do not attempt repair until the cause is understood.

**System State.** Other services functioning? Disk space adequate? Filesystems healthy? If physical kill was used, run filesystem checks after restoration.

**Log Analysis.** Read logs from the last known-good run forward. Identify what changed and what actions the automation took during the failure.

**Scope of Damage.** What damage occurred? Confined to documented scope or wider? Reversible? Any ongoing risk?

**Root Cause.** Identify why the failure occurred before restarting. Common causes: resource exhaustion, dependency failure, unexpected input, software bug, environmental change. If root cause is unknown, the automation must not restart until identified or a Tier 2 risk-acceptance decision is made per GOV-001.

### 4.4 Restart Procedures

**Single-Agent Restart:** Verify root cause addressed. Review ASD. If automation was modified, run relevant D8-002 checklist items. Restart and actively monitor the first run (watch logs in real time). Verify success before returning to normal monitoring. If first run fails, halt and investigate further.

**Category Restart:** Restart one agent at a time in criticality order. Verify each before starting the next. Monitor for inter-agent interaction issues.

**Total Restart:** Restart in priority order, one group at a time:
- Priority 1: Data protection (backup, integrity verification).
- Priority 2: System monitoring (disk space, hardware health).
- Priority 3: Operational (log rotation, cleanup, reporting).
- Priority 4: All others.

Within each group, restart one agent at a time per the single-agent procedure. Total restart may take hours or days. During the restart period, perform all automated tasks manually using ASD fallback procedures. Simultaneous restart of all automation is prohibited.

### 4.5 Kill Switch Testing

**Individual Kill Test:** Annually during D8-004 audit, and during D8-002 code review for new automation.

**Category Kill Test:** Semi-annually. Verify correct agent identification and successful stop.

**Total Kill Test:** Annually during a scheduled maintenance window. Execute total kill, perform post-stop checks, restart all automation using the total restart procedure. This also exercises restart procedures and manual fallback.

**Physical Kill Test:** Annual review of physical access to power controls. Actual disconnection test only when data loss risk is acceptable (e.g., after verified backup).

**Laminated Quick-Reference Card:** A summary of the emergency shutdown sequence posted in Zone 1 per SEC-002. Contains essential commands and physical kill procedure, written for execution under stress. Updated whenever procedures change. Verified during annual test.

## 5. Rules & Constraints

- **R-D8-05-01:** Every automation must have a documented kill procedure in ASD Section E.
- **R-D8-05-02:** Kill scripts (`kill-all.sh`, `kill-category.sh`) must be updated whenever the automation inventory changes.
- **R-D8-05-03:** In a suspected cascading failure, execute Level 3 total kill immediately. Diagnosis follows the stop.
- **R-D8-05-04:** No restart after emergency kill until post-stop checks are complete and root cause is understood or risk-acceptance is documented.
- **R-D8-05-05:** Kill switch tests must follow the schedules in Section 4.5. A failed test means the automation is paused until the kill switch is repaired.
- **R-D8-05-06:** Physical kill must be executable without software interaction. Power connections must be accessible without moving equipment or unlocking additional enclosures.
- **R-D8-05-07:** The laminated quick-reference card must be posted in Zone 1 and updated whenever procedures change.
- **R-D8-05-08:** After total kill, restart one agent at a time in priority order. Simultaneous restart is prohibited.

## 6. Failure Modes

- **Kill switch rot.** Scripts reference agents that no longer exist or miss newly added agents. Mitigation: R-D8-05-02 ties updates to inventory changes. Annual total kill test catches drift.

- **Untested kill procedures.** ASD kill commands have never been tested during execution. Mitigation: D8-002 code review includes kill test; annual testing per Section 4.5.

- **Panic-driven diagnosis.** The operator tries to understand the problem while automation is still running. Mitigation: R-D8-05-03 and the shutdown sequence are explicit: stop first, diagnose second.

- **Premature restart.** Operator restarts without completing post-stop checks. Root cause persists. Problem recurs. Mitigation: R-D8-05-04 prohibits restart without checks.

- **Physical access obstruction.** Power cable is inaccessible. Mitigation: R-D8-05-06 requires unobstructed access, verified annually.

- **Card staleness.** Laminated card references outdated commands. Mitigation: tied to kill script updates; verified during annual test.

## 7. Recovery Procedures

1. **If the total kill script fails during an emergency:** Fall back to manual software kill immediately. If that also fails, use the physical kill. Debug the script afterward, not during the emergency.

2. **If physical kill caused data corruption:** Boot from recovery media. Run filesystem checks before mounting read-write. Compare against verified backup. Restore from backup per D6-006 if corruption is confirmed.

3. **If root cause cannot be identified:** Document everything known. Make a documented Tier 2 decision per GOV-001 about whether to restart with enhanced monitoring (reduce to Level 2 requiring human approval) or leave disabled until more information is available.

4. **If restart reveals multiple agents affected by the same root cause:** Stop the restart. The common cause may be environmental. Address it before restarting any agent.

5. **If the laminated card is missing during an emergency:** Consult this document. Replace the card immediately afterward. Consider a second card at a different location.

## 8. Evolution Path

- **Years 0-5:** Practice the shutdown sequence even when nothing is wrong. Build muscle memory. The total kill test is a brief exercise with a small inventory.

- **Years 5-15:** Larger inventory makes kill scripts more complex. Category kills become important. Annual total kill test may take several hours for the full restart.

- **Years 15-30:** Hardware changes require updates to physical kill procedures. New operators must practice until the response is reflexive.

- **Years 30-50+:** The kill switch architecture must remain simple even as the automation estate evolves. The kill switch is the one system that must work when everything else is broken.

- **Signpost for revision:** If the software total kill takes more than 5 minutes or the physical kill more than 30 seconds, the procedures or infrastructure need simplification.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The laminated card may be the most important artifact in this article. Everything else is documentation I will consult at leisure. The card is what I will read at 2 AM when something is on fire and my hands are shaking. It needs to be right, and it needs to be where I can see it.

I have made the total kill procedure aggressive. When in doubt, stop everything. The cost of stopping everything is manual labor during restart. The cost of not stopping everything is potential data loss and cascading damage. Manual labor is always cheaper.

The restart procedure is deliberately slow -- one agent at a time, verified before the next. This feels agonizing when ten agents are down. But the restart is the moment of greatest risk: if the root cause was not fully addressed, the first restart triggers the problem again. Slow restart limits the blast radius of an incomplete fix.

Physical kill switches feel anachronistic. They are not. Pulling a power cord always works. Everything else is a best effort.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 1: Sovereignty; Principle 2: Integrity Over Convenience)
- CON-001 -- The Founding Mandate (human authority over all institutional systems)
- GOV-001 -- Authority Model (Tier 2 decisions for restart with unknown root cause)
- SEC-001 -- Threat Model (Pillar 1: Assume Breach, Prevent Breach)
- SEC-002 -- Access Control Procedures (Zone 1, service accounts)
- OPS-001 -- Operations Philosophy (documentation-first principle)
- D6-001 -- Data Philosophy (data tier system, integrity requirements)
- D6-006 -- Backup Doctrine (restoration after corruption)
- D6-007 -- Data Integrity & Verification (post-stop integrity checks)
- D8-001 -- Automation Philosophy: The Restraint Doctrine (servant principle, kill switch mandate)
- D8-002 -- Agent Design Principles (ASD Section E, code review kill test)
- D8-003 -- Human-in-the-Loop Procedures (escalation matrix, Severity 4)
- D8-004 -- Automation Audit (annual kill switch testing)
- D10-002 -- Daily Operations Doctrine (manual fallback during restart)
- D10-003 -- Incident Response Procedures (broader incident response)
- Stage 1 Documentation Framework, Domain 8: Automation & Agents

---

---

# D8-006 -- Legacy Automation Assessment

**Document ID:** D8-006
**Domain:** 8 -- Automation & Agents
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, SEC-001, OPS-001, D8-001, D8-002, D8-003, D8-004, D8-005
**Depended Upon By:** D10-002. Referenced by all articles involving succession or institutional transfer.

---

## 1. Purpose

This article defines the procedures for evaluating automation inherited from a previous maintainer. It addresses the most dangerous moment in an institution's automation lifecycle: the transition from one operator to the next, when the person who built and understood the automation is no longer available.

D8-001 calls this the "black box inheritance" problem. This article provides the systematic procedure for opening those black boxes: the assessment protocol for determining what inherited automation does, the reverse engineering procedures for reconstructing understanding when documentation is insufficient, and the keep-or-kill framework for deciding whether each piece of inherited automation earns its continued existence.

This article is written for the new operator who has just assumed responsibility for an institution they did not build. It is also useful for any operator who discovers automation they no longer understand, whether because they forgot what they built years ago or because the automation was created by a process that did not leave adequate documentation.

## 2. Scope

This article covers:

- The legacy assessment trigger: when to conduct a full assessment.
- The discovery process: finding all automation running in the institution.
- The black box assessment protocol: determining what unknown automation does.
- Reverse engineering procedures: reconstructing understanding from code, logs, and artifacts.
- The keep-or-kill framework for deciding each automation's fate.
- Documentation reconstruction for automation that lacks ASDs.
- Timeline and priority structure.

This article does not cover:

- Designing replacement automation (see D8-002).
- Routine annual review of the current operator's automation (see D8-004).
- Emergency response to malfunctioning legacy automation (see D8-005; stop it first, assess later).
- The broader succession process (see GOV-001 Section 7).

## 3. Background

### 3.1 The Inevitability of Legacy

In a fifty-year institution, legacy automation is a certainty. Every current operator is a future predecessor. Every automation running today is a future legacy system. Good preparation means building comprehensible automation today (per D8-002) and having a procedure for assessing automation when comprehensibility falls short. This article provides the second.

### 3.2 Why Assessment Before Action

The natural impulse is either to leave unfamiliar automation alone (fear of breaking something) or to replace it immediately (desire for a clean start). Both are dangerous. Leaving it alone means running systems you do not understand. Replacing it immediately means discarding accumulated logic and edge-case handling without understanding what you are discarding. Assessment provides a third path: understand first, decide second.

### 3.3 The Spectrum of Legacy Quality

At the best end: complete ASDs, clean code, readable logs, tested fallbacks. The assessment is verification. In the middle: partial documentation, readable code, some gaps to fill. At the worst end: no documentation, opaque code, missing logs. The full black box assessment protocol applies. This article provides procedures for all three.

## 4. System Model

### 4.1 The Legacy Assessment Trigger

A full assessment is triggered by: succession event (all automation assessed within 60 days), extended absence exceeding six months, discovery of unknown automation, or comprehension failure (the operator cannot explain their own automation).

### 4.2 The Discovery Process

Before assessing legacy automation, find all of it:

**Step 1: Scheduled tasks.** Review all cron jobs for all users (`crontab -l -u [user]` for each), all systemd timers (`systemctl list-timers --all`), at jobs (`atq`), and anacron (`/etc/anacrontab`).

**Step 2: Running processes.** List all processes, especially under service accounts. List all systemd services. Check for screen/tmux sessions with long-running automation.

**Step 3: Startup scripts.** Review `/etc/rc.local`, `/etc/init.d/`, systemd units in `/etc/systemd/system/`, user-level autostart files, and cron `@reboot` entries.

**Step 4: Expected locations.** Review `/institution/automation/agents/` per D8-002. Search for scripts in unexpected locations.

**Step 5: ASD directory.** Compare ASDs against discovered automation. Flag mismatches.

**Step 6: Compile the discovery inventory.** Location, apparent purpose, ASD status, whether previously known, initial risk assessment.

### 4.3 The Black Box Assessment Protocol

For each automation without adequate documentation:

**Phase 1: External Observation (Do Not Touch the Automation).**

Step 1: Read any existing documentation -- even partial ASDs, cron comments, README files.

Step 2: Read the source code. Note: inputs consumed, outputs produced, external commands called, files modified, exit conditions.

Step 3: Read the logs. Determine: run frequency, last execution, success/failure history, reported actions.

Step 4: Examine the footprint. What files does it own? What directories does it write to? What permissions and service account?

Step 5: Map dependencies. Does it read files from other automations? Write files consumed by others? Depend on specific services or configurations?

**Phase 2: Controlled Interaction.**

Step 6: If purpose remains unclear, copy the automation to a sandbox environment and run with test data if feasible. Compare results against production logs.

Step 7: If sandboxing is not feasible, consider running in production with enhanced monitoring (verbose logging, strace, or a file-operation wrapper). Exercise caution -- monitoring may affect behavior.

**Phase 3: Assessment Summary.**

Step 8: Produce a Legacy Assessment Report: temporary ID, discovery location, assessed purpose with confidence level per D7-001 epistemic hierarchy, assessed risk (malfunction risk and stop risk), dependency map, documentation status, and recommendation: keep, modify, or kill.

### 4.4 The Keep-or-Kill Framework

**Keep (with documentation reconstruction):** The automation serves a clear, justified purpose. It passes or can be made to pass the five-criterion test. The operator can understand it sufficiently to maintain it. Action: reconstruct the ASD per Section 4.5, integrate into normal monitoring per D8-003 and D8-004.

Criteria: clear purpose serving the institution, comprehensible (or can be made so), observable (or can be improved), risk of continued operation lower than risk of stopping, manual fallback reconstructable.

**Modify (redesign and replace):** The purpose is necessary but the implementation is unacceptable -- incomprehensible or failing criteria that cannot be remediated without extensive changes. Action: design a replacement per D8-002, test, deploy, decommission the legacy automation per D8-004.

Criteria: clear and justified purpose, implementation fails criteria, remediation approaches a rewrite, manual fallback can sustain the function during replacement.

**Kill (decommission):** The automation does not serve a justified purpose, or the risk of continued operation exceeds the benefit, or the operator cannot achieve sufficient comprehension to maintain it safely. Action: decommission per D8-004 Section 4.5 including the 30-day observation period.

**The Default:** When the operator cannot determine whether to keep or kill, the default is kill. The burden of proof lies with the automation. The 30-day observation period provides the safety net.

### 4.5 Documentation Reconstruction

For automation receiving a "Keep" verdict:

1. Create a new ASD using the D8-002 template. Record "Unknown -- assessed by [operator] on [date]" if the original author is unknown.
2. Fill in Section B (Justification) based on current conditions, not original conditions.
3. Fill in Section C (Technical Specification) from source code analysis and dependency mapping.
4. Fill in Section D (Failure Modes) from historical logs, source code error handling, and operator judgment.
5. Fill in Section E (Kill Procedure) -- determine or design one and test it.
6. Fill in Section F with an initial entry documenting the reconstruction.
7. Write or update the manual fallback procedure. Test it.
8. Run the D8-002 code review checklist. Remediate minor failures; reclassify to "Modify" for major failures.

### 4.6 Timeline and Priority Structure

Following a succession event, 60 days:

**Days 1-7: Discovery and Triage.** Complete discovery. Classify each automation as High Priority (actively modifying data or performing critical functions), Medium Priority (monitoring, reporting, non-critical functions), or Low Priority (unclear purpose, apparently inactive, trivial). If any automation is causing active harm, invoke D8-005 immediately.

**Days 8-30: High Priority Assessment.** Black box protocol for all High Priority automations. The operator may reduce these to Level 2 (human approval for each action) as a safety measure during assessment.

**Days 31-50: Medium Priority Assessment.**

**Days 51-60: Low Priority Assessment and Close.** Produce the Legacy Assessment Summary. Record in the decision log.

For estates larger than 20 agents, the timeline may be extended with documented justification.

## 5. Rules & Constraints

- **R-D8-06-01:** Full legacy assessment within 60 days of succession. No inherited automation at Level 3 or 4 until it has passed assessment.
- **R-D8-06-02:** The discovery process must be thorough. All locations in Section 4.2 must be checked.
- **R-D8-06-03:** During assessment, all unassessed automation operates at Level 2 or lower.
- **R-D8-06-04:** The keep-or-kill default is kill.
- **R-D8-06-05:** Every "Keep" verdict requires a complete ASD per D8-002 within 30 days.
- **R-D8-06-06:** Legacy Assessment Reports are Tier 2 data per D6-001.
- **R-D8-06-07:** Do not modify legacy automation during Phase 1 (external observation). Understanding precedes modification.
- **R-D8-06-08:** Automation posing a security risk (excessive privileges, unauthorized data access, external communication) is treated as a security incident per D10-003 regardless of assessment timeline.

## 6. Failure Modes

- **Discovery incompleteness.** Automation in unexpected locations runs without oversight. Mitigation: comprehensive discovery procedure. Expand search if predecessor used non-standard locations.

- **Assessment paralysis.** Volume overwhelms the operator, who defers indefinitely. Mitigation: the timeline structure and triage prioritization. Start with one agent. Complete it. Momentum builds.

- **Premature kill.** Frustration leads to killing everything. Critical functions are disrupted. Mitigation: assessment before decision; the 30-day observation period catches mistakes; the "Modify" option preserves purpose while replacing implementation.

- **Sentimental keep.** Automation kept out of respect for the predecessor rather than merit. Mitigation: R-D8-06-04 establishes kill as the default. The five-criterion test, not sentiment, governs.

- **Security risk in legacy code.** Predecessor built automation running as root, with hardcoded credentials, or accessing data beyond its scope. Mitigation: R-D8-06-08 treats security risks as immediate incidents.

- **Dependency blindness.** Killing automation without realizing other systems depend on it. Mitigation: Phase 1 dependency mapping (Section 4.3, Step 5) and the 30-day observation period.

## 7. Recovery Procedures

1. **If undiscovered automation is found after assessment completion:** Conduct the black box protocol immediately. Apply R-D8-06-03 (Level 2 or lower). Expand the discovery procedure.

2. **If premature kill causes disruption:** Reactivate from the retired directory. Conduct the full assessment with better information. Use "Modify" if the function is needed but the implementation is unacceptable.

3. **If assessment paralysis has set in:** Acknowledge it. Pick one automation -- the highest priority -- and complete its assessment. Sequential, thorough assessment of one agent is better than attempted parallel assessment of all.

4. **If security risks are discovered:** Invoke D10-003 immediately. Halt the automation. Rotate affected credentials. Audit accessed data. Resume the legacy assessment with heightened scrutiny.

5. **If documentation reconstruction reveals unexpected complexity:** Reclassify from "Keep" to "Modify." Rebuild from scratch based on the now-understood purpose.

## 8. Evolution Path

- **Years 0-5:** No legacy automation exists if the institution is newly founded. Use this period to build automation that will be good legacy -- follow D8-002 so rigorously that a future assessment will be verification, not archaeology.

- **Years 5-15:** The founding operator may encounter their own legacy: automation from year 1 that they no longer fully remember. Apply these tools to their own work as a test of documentation quality.

- **Years 15-30:** The first real legacy assessment occurs. ASD quality from years 0-15 determines whether the assessment is a manageable transition or a crisis.

- **Years 30-50+:** Multiple succession events produce layers of assessment documentation. Legacy Assessment Reports become a historical record of automation evolution across generations.

- **Signpost for revision:** If assessments consistently find D8-002 standards were followed and ASDs are accurate, the process can be streamlined. If standards were consistently not followed, the upstream problem is D8-002 compliance.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
I am writing legacy assessment procedures for automation that does not yet exist, for a successor who may not yet be born. This feels like writing a letter to be opened after my death -- and in a real sense, it is. When someone follows these procedures, it will be because I am no longer maintaining this institution.

This awareness changes how I write documentation. Every ASD is not just for me -- it is for the person who will conduct the legacy assessment of my work. Every comment, every header, every log entry is a message to that person. The question is not "do I understand this?" The question is "will they?"

The keep-or-kill default of "kill" may seem harsh applied to one's own work. It is not. It says: I trust my successor's judgment more than I trust my own documentation. If my documentation is good enough, the automation will pass the five-criterion test and be kept. If not, it deserves to be replaced by something the successor understands. My ego is not part of the institution's architecture.

The most dangerous legacy is not bad automation. It is automation just barely comprehensible enough to keep but not well enough understood to maintain confidently -- sitting in a twilight zone where the operator is afraid to touch it and afraid to stop it. The black box assessment protocol drags every automation out of that twilight and into either full comprehension or the clarity of decommissioning.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 3: Transparency; Principle 5: Succession as a Design Constraint)
- CON-001 -- The Founding Mandate (institutional continuity, comprehensibility)
- GOV-001 -- Authority Model (succession protocol Section 7, decision tiers)
- SEC-001 -- Threat Model (Pillar 3: The Human Is the System)
- SEC-002 -- Access Control Procedures (service accounts, privilege model)
- OPS-001 -- Operations Philosophy (documentation-first principle, complexity budget)
- D6-001 -- Data Philosophy (Tier 2 classification for assessment documents)
- D7-001 -- Intelligence Philosophy (epistemic hierarchy for confidence levels)
- D8-001 -- Automation Philosophy: The Restraint Doctrine (five-criterion test, black box problem)
- D8-002 -- Agent Design Principles (ASD template for documentation reconstruction)
- D8-003 -- Human-in-the-Loop Procedures (Level 2 oversight during assessment)
- D8-004 -- Automation Audit (decommissioning procedure, Health Summary)
- D8-005 -- Emergency Override and Kill Switch Procedures (emergency response to dangerous legacy)
- D9-001 -- Education Philosophy (self-teaching requirement)
- D10-003 -- Incident Response Procedures (security incidents in legacy automation)
- Stage 1 Documentation Framework, Domain 8: Automation & Agents
