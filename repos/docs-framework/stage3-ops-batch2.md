# STAGE 3: OPERATIONAL DOCTRINE (BATCH 2)

## Disaster Recovery Operations, Asset Management, Technology Migration, and Decision Records

**Document ID:** STAGE3-OPS-BATCH2
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Classification:** Operational Procedures -- These articles translate Stage 2 philosophy into step-by-step actionable procedures for disaster recovery drills, system rebuilds, asset tracking, technology migration, and decision record keeping.

---

## How to Read This Document

This document contains five operational articles spanning Domains 11, 12, 13, and Governance. They are the second batch of Stage 3 operational doctrine for the holm.chat Documentation Institution. Where Stage 2 established how to think about disaster recovery, administration, evolution, and governance, these articles establish how to do the work.

These are manuals. They contain procedures, templates, checklists, schedules, and grading criteria. They are designed to be consulted during execution, not merely read for understanding. Every procedure is written for a single operator who may be stressed, tired, or unfamiliar with the system -- because those are the conditions under which operational procedures are most needed and least forgiving.

These articles assume you have read the relevant Stage 2 philosophy articles: D12-001 (Disaster Recovery Philosophy), D11-001 (Administration Philosophy), D13-001 (Evolution Philosophy), and GOV-001 (Authority Model). They also assume familiarity with the Core Charter, particularly OPS-001 (Operations Philosophy) and SEC-001 (Threat Model and Security Philosophy). If you have not read those documents, read them first. The procedures here implement principles established there.

If something in these procedures no longer matches your reality -- because hardware has changed, because software has evolved, because a referenced tool no longer exists -- do not abandon the procedure. Adapt it. The principles behind each procedure are stated in the Background section. Use those principles to find the equivalent procedure for your current context. Then update this document.

---

---

# D12-002 -- Disaster Recovery Drill Framework

**Document ID:** D12-002
**Domain:** 12 -- Disaster Recovery
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, SEC-001, OPS-001, D12-001, D6-002, D6-003, SEC-002, SEC-003
**Depended Upon By:** D12-003, D13-002, GOV-007. All articles referencing disaster preparedness verification or institutional resilience testing.

---

## 1. Purpose

This article defines how the holm.chat Documentation Institution tests its own survival. It is the operational manual for disaster recovery drills -- the scheduled, controlled exercises that verify whether the institution can actually recover from the failures it claims to be prepared for.

D12-001 established a foundational principle: drills matter more than plans. A mediocre plan that has been drilled quarterly is worth more than a perfect plan that has never been tested. This article operationalizes that principle. It defines when drills happen, what scenarios they simulate, how they are executed, how performance is evaluated, and how findings are fed back into the institution's recovery procedures.

The purpose of a drill is not to confirm that everything works. The purpose is to discover what does not work -- before a real disaster forces the discovery under catastrophic conditions. A drill that reveals no problems is either testing too easy a scenario or not looking hard enough. The institution should expect every drill to produce findings. Those findings are the drill's primary output and the institution's most valuable resilience investment.

This article is written for the operator who must plan, execute, and evaluate drills alone. Every procedure is designed for solo execution. Where a drill would benefit from a second person -- and some would -- the article notes this but provides a solo-execution alternative.

## 2. Scope

**In scope:**
- The drill schedule: frequency, rotation, and integration with the operational tempo.
- Drill types: tabletop, functional, partial recovery, full recovery, surprise.
- Five specific drill scripts for the five most likely disaster scenarios.
- Execution procedures for each drill type.
- The evaluation framework: how to grade drill performance.
- After-action report format and procedures.
- Integration of drill findings into recovery procedure updates.
- Safety controls: how to prevent a drill from becoming a real disaster.

**Out of scope:**
- The recovery procedures themselves (those are defined in D6-002, D6-003, SEC-003, and subsequent Domain 12 articles; this article tests them).
- The disaster recovery philosophy (D12-001).
- Post-incident review of real disasters (DR-013; this article covers drills, not real events).
- Specific hardware failure diagnostics (Domain 4 and Domain 5 articles).

## 3. Background

### 3.1 Why Drills Are Non-Negotiable

D12-001 Section 3.2 documents the failure of untested plans. This section does not repeat that argument. It adds a specific observation relevant to drill design: the value of a drill is proportional to the discomfort it creates.

A drill that exercises only the procedures the operator is confident about teaches nothing. A drill that forces the operator to execute a procedure they have not touched in six months, on hardware they rarely interact with, under time pressure -- that drill teaches everything. It reveals forgotten steps, outdated assumptions, missing tools, and degraded media. It reveals the gap between the operator's mental model of the institution and the institution's actual state.

This discomfort is the point. The drill schedule in this article is designed to ensure that every critical recovery procedure is exercised at least once per year, including the procedures the operator would rather not think about.

### 3.2 The Air-Gap Drill Constraint

In a networked institution, disaster recovery drills can leverage virtualization, cloud-based test environments, and automated testing frameworks. This institution has none of those luxuries. Drills must be conducted on physical hardware, using physical backup media, following physical procedures. This makes drills more expensive in time and effort -- but it also makes them more realistic, because real disasters will also involve physical hardware, physical media, and physical procedures.

The air-gap constraint also means that drill findings cannot be supplemented by searching the internet for solutions. If a drill reveals that a procedure is incomplete, the operator must solve the problem from the institution's own documentation and knowledge. This is exactly what will happen during a real disaster, which is why air-gapped drills are more valuable, not less, than their networked equivalents.

### 3.3 The Drill Safety Principle

A drill must never cause the harm it is designed to prepare for. This is the drill safety principle, and it constrains every procedure in this article. Drills are conducted on backup copies, on test partitions, on secondary hardware, or through tabletop simulation -- never on production data or production systems without explicit safeguards.

The one exception is the full recovery drill, which requires taking a system offline and rebuilding it from backup. This drill is performed only on systems that can be offline for the duration without causing data loss or service interruption, and only after verifying that current backups are valid.

## 4. System Model

### 4.1 The Drill Schedule

Drills are integrated into the operational tempo defined in OPS-001:

**Quarterly Drills (4 per year):**
Each quarter, one functional drill is conducted from the five-scenario rotation defined in Section 4.3. The rotation ensures that all five scenarios are covered within a fifteen-month cycle. The quarterly drill is scheduled during the quarterly operations day defined in OPS-001 Section 4.1.

Quarterly drill rotation:
- Q1 (January-March): Scenario 1 -- Primary storage failure.
- Q2 (April-June): Scenario 2 -- Backup corruption discovery.
- Q3 (July-September): Scenario 3 -- Complete site loss.
- Q4 (October-December): Scenario 4 -- Operator incapacitation simulation.
- Q1 (next year): Scenario 5 -- Encryption key loss.
- Q2 (next year): Return to Scenario 1 with variations.

**Annual Full Recovery Drill (1 per year):**
During the annual operations cycle (OPS-001 Section 4.1), one full recovery drill is conducted. This drill rebuilds a complete subsystem from backup, verifying the end-to-end recovery path from backup media to operational service. The subsystem selected for the annual drill should rotate each year to ensure comprehensive coverage.

**Monthly Tabletop Reviews (12 per year):**
During each monthly operations session, a 30-minute tabletop review of one recovery procedure is conducted. The operator reads through the procedure, mentally executes each step, and notes any steps that reference outdated hardware, software, or configurations. This is a documentation audit disguised as a drill. It costs minimal time and catches procedure rot before it becomes dangerous.

**Surprise Drills (at least 1 per year):**
At least once per year, the operator should conduct an unscheduled drill. The method: write the five scenario numbers on slips of paper, draw one at random on a day that is not a scheduled drill day, and execute the drill that session. The purpose of surprise drills is to test readiness under realistic conditions -- real disasters do not arrive on scheduled days.

### 4.2 Drill Types

**Type 1: Tabletop Drill.**
The operator reads through a recovery procedure step by step, without executing any commands or touching any hardware. At each step, the operator asks: do I have this tool? Do I know where this data is? Does this step still apply to the current system? Is this step clear enough for a stressed operator? Findings are recorded. Duration: 30-60 minutes.

Use tabletop drills for: monthly procedure reviews, procedures that cannot be safely executed in a drill (such as full-site evacuation), and initial familiarization with new procedures.

**Type 2: Functional Drill.**
The operator executes the recovery procedure on non-production systems or copies. Data is actually restored, keys are actually used, configurations are actually applied -- but on test targets, not on production systems. The operator follows the documented procedure exactly as written, noting any deviations required. Duration: 2-6 hours depending on scenario.

Use functional drills for: quarterly scenario rotation, testing specific recovery procedures, and validating backup integrity.

**Type 3: Partial Recovery Drill.**
The operator takes a non-critical subsystem offline and recovers it from backup. This is a real recovery on real hardware, but scoped to a subsystem whose temporary absence does not threaten institutional data or critical operations. Duration: 4-8 hours.

Use partial recovery drills for: annual full recovery drill (applied to a selected subsystem), testing end-to-end recovery paths, and building operator confidence with real recovery operations.

**Type 4: Full Recovery Drill.**
The operator rebuilds a complete system from bare metal using backup media and documented procedures. This is the most demanding drill type. It should only be attempted after multiple successful partial recovery drills and only on a system that can be offline for up to 24 hours. Duration: 8-24 hours.

Use full recovery drills for: validating D12-003 (System Rebuild Procedures), testing the rebuild kit, and succession readiness verification.

### 4.3 The Five Disaster Scenario Scripts

Each script defines the simulated disaster, the pre-drill setup, the drill execution steps, the success criteria, and the time limit.

**Scenario 1: Primary Storage Failure.**

Simulated disaster: The primary storage device has failed. All data on the primary volume is inaccessible. The operating system may or may not be affected (drill variant A: OS survives on a separate disk; variant B: OS is also lost).

Pre-drill setup:
1. Verify that current backups exist and are valid (do not begin the drill without confirmed good backups).
2. Identify the test target: a spare drive or partition that will serve as the "replacement" primary storage.
3. Document the current state of the primary storage contents (file list, total size) for post-drill verification.
4. If executing variant B, prepare bootable installation media.

Drill execution (variant A -- OS survives):
1. Simulate the loss by unmounting or logically disconnecting the primary data volume. Do not physically destroy it -- this is a drill.
2. Following the backup restoration procedure in D6-002 Section 4.4, restore the most recent Tier 1 backup to the test target.
3. Verify the restored data: compare file counts, spot-check files, verify checksums against the backup manifest.
4. Restore Tier 2 data from the most recent backup.
5. Verify service configurations that depend on the restored data.
6. Record the total time from "failure" to "verification complete."

Drill execution (variant B -- OS also lost):
1. Simulate by booting from installation media (do not wipe the actual OS disk -- mark it as "failed" and set it aside).
2. Install the operating system on the test target per D12-003 procedures.
3. Restore system configuration from backup.
4. Restore data per variant A steps 2-5.
5. Verify that all critical services start and function correctly.
6. Record total time.

Success criteria:
- All Tier 1 data restored and verified within 4 hours (variant A) or 8 hours (variant B).
- All Tier 2 data restored and verified within 8 hours (variant A) or 12 hours (variant B).
- No data loss beyond the recovery point objective defined in D6-002.
- All critical services operational after restoration.

Time limit: 8 hours (variant A), 16 hours (variant B). If the time limit is exceeded, the drill continues to completion but the time overrun is recorded as a finding.

Post-drill:
1. Reconnect the "failed" volume (it was never actually damaged).
2. Verify production systems are intact.
3. Clean up test targets.
4. Complete the after-action report.

**Scenario 2: Backup Corruption Discovery.**

Simulated disaster: During a routine test restoration, you discover that your most recent backup set is corrupted. The backup appears complete but checksum verification fails on multiple files.

Pre-drill setup:
1. Identify a backup set to use as the "corrupted" set (select the most recent monthly backup).
2. Prepare a list of files to mark as "corrupted" (select 5-10 files across Tier 1 and Tier 2 data).
3. Do not actually corrupt the files. Instead, create a drill scenario document listing the "corrupted" files that the operator must treat as unrecoverable from this backup set.

Drill execution:
1. Begin by "discovering" the corruption: review the drill scenario document that lists the corrupted files.
2. Assess the scope: what data is affected? What tier? What is the most recent known-good copy?
3. Locate the next-oldest backup set (Location A weekly, Location B monthly, or Location C off-site).
4. Restore the "corrupted" files from the older backup set.
5. Verify the restored files.
6. Assess the gap: is there any data that existed only in the corrupted backup and not in any older backup? Document the potential data loss.
7. Review the backup verification procedures: would the existing verification process have caught this corruption? If not, what additional checks are needed?

Success criteria:
- Operator correctly identifies the scope of corruption within 30 minutes.
- Operator locates and uses an alternative backup within 2 hours.
- Operator identifies any unrecoverable data gap and documents it.
- Operator proposes at least one improvement to backup verification procedures.

Time limit: 4 hours.

**Scenario 3: Complete Site Loss.**

Simulated disaster: The primary site has been destroyed (fire, flood, catastrophic event). All hardware at the primary site is gone. The only surviving assets are the off-site backups (Location C) and any physical documentation stored off-site.

Pre-drill setup:
1. This is primarily a tabletop drill with a functional component. The tabletop component addresses the full scope of site loss. The functional component tests restoration from off-site backup media.
2. Gather the off-site backup media (or a copy of it) and the printed recovery documentation.
3. If a Succession Readiness Package exists off-site, retrieve it (or a copy) for the drill.
4. Prepare a separate machine (a spare or test machine) to serve as the "new site."

Drill execution -- tabletop component:
1. Assume you are standing in front of the spare machine with only the off-site backup media and printed documentation. No other institutional resources are available.
2. Walk through the complete rebuild sequence defined in D12-003: hardware assessment, OS installation, base configuration, data restoration, service configuration, verification.
3. At each step, verify: do I have the information I need in the printed documentation? Are the credentials I need accessible from the off-site materials? Are the procedures clear enough to follow without additional context?
4. Record every gap, every missing piece of information, every assumption that would not hold at a new site.

Drill execution -- functional component:
1. On the spare machine, restore Tier 1 data from the off-site backup.
2. Verify the restoration.
3. Attempt to access the documentation corpus from the restored data.
4. Record the time from "start of restoration" to "documentation corpus accessible and readable."

Success criteria:
- Tabletop component identifies all information gaps in off-site documentation.
- Functional component successfully restores Tier 1 data from off-site backup.
- Documentation corpus is accessible and readable from restored data.
- Operator can articulate a realistic timeline for full institutional rebuild.

Time limit: 6 hours (2 hours tabletop, 4 hours functional).

**Scenario 4: Operator Incapacitation Simulation.**

Simulated disaster: The primary operator is incapacitated. A successor (simulated by the operator themselves, acting as an unfamiliar person) must assume operations using only the Succession Readiness Package and institutional documentation.

Pre-drill setup:
1. Retrieve the Succession Readiness Package.
2. Set aside all personal knowledge of the institution. For the duration of this drill, the operator pretends to be a technically competent stranger who has never seen this system before.
3. This is primarily a documentation audit. The operator attempts to perform basic operations using only the written documentation, without relying on any knowledge not contained in the documents.

Drill execution:
1. Open the Succession Readiness Package. Read the letter of introduction and the first-30-days guide per GOV-004 Section 4.6.
2. Using only the documentation, attempt to:
   a. Identify all physical hardware and its function.
   b. Log in to the operator account (credentials should be accessible through the succession provisions in SEC-003).
   c. Perform the daily operations checklist per D10-002.
   d. Locate and read the most recent decision log entries.
   e. Locate and verify the most recent backup.
   f. Identify the next scheduled maintenance task.
3. At each step, record: was the information sufficient? Was anything ambiguous? Was anything missing? Would a stranger be able to complete this step?
4. Attempt one recovery procedure (select the simplest one) using only the documentation.

Success criteria:
- All six sub-tasks in step 2 are completable from documentation alone.
- No critical information requires knowledge not present in the documentation.
- The selected recovery procedure is executable from documentation alone.
- The Succession Readiness Package is current (all references match current system state).

Time limit: 8 hours.

**Scenario 5: Encryption Key Loss.**

Simulated disaster: The active copy of a critical encryption key has been lost (simulated -- not actually deleted). The operator must recover the key from backup and restore access to the encrypted data.

Pre-drill setup:
1. Select a non-critical encrypted volume or a test encrypted volume for the drill.
2. Record the key fingerprint and the location of all backup copies per SEC-003 key inventory.
3. Do not actually delete the active key. Instead, simulate the loss by pretending the active key location is inaccessible.

Drill execution:
1. "Discover" the key loss: attempt to access the encrypted volume and find that the active key is unavailable.
2. Consult the key inventory (SEC-003 Section 4.1) to identify backup key locations.
3. Retrieve the Tier B backup key from Zone 1.
4. Decrypt the Tier B backup media.
5. Extract the lost key.
6. Use the recovered key to access the encrypted volume.
7. Verify data integrity on the decrypted volume.
8. Following SEC-003 Section 7.1, create a new Tier B backup (since the original was opened during recovery).

Success criteria:
- Operator locates the correct backup key within 30 minutes.
- Key is successfully recovered from Tier B backup.
- Encrypted volume is accessible with the recovered key.
- Data integrity is verified.
- New Tier B backup is created.

Time limit: 4 hours.

### 4.4 The Evaluation Framework

Drill performance is evaluated on five dimensions, each scored on a three-point scale:

**Dimension 1: Procedure Adherence.**
Did the operator follow the documented procedure, or did they improvise?
- 3 (Excellent): Procedure followed as written with no deviations.
- 2 (Adequate): Minor deviations required, all documented during drill.
- 1 (Inadequate): Major deviations required, or operator abandoned procedure and improvised.

**Dimension 2: Time Performance.**
Was the recovery completed within the defined time limit?
- 3 (Excellent): Completed within 75% of the time limit.
- 2 (Adequate): Completed within the time limit.
- 1 (Inadequate): Time limit exceeded.

**Dimension 3: Data Integrity.**
Was all data recovered intact?
- 3 (Excellent): All data verified, all checksums match.
- 2 (Adequate): All critical data recovered, minor discrepancies in non-critical data.
- 1 (Inadequate): Critical data loss or unresolvable integrity failures.

**Dimension 4: Documentation Quality.**
Were the recovery procedures adequate for the scenario?
- 3 (Excellent): Procedures were complete, accurate, and clear.
- 2 (Adequate): Procedures were usable but required clarification or minor updates.
- 1 (Inadequate): Procedures were outdated, incomplete, or misleading.

**Dimension 5: Finding Generation.**
Did the drill produce actionable findings for improving resilience?
- 3 (Excellent): Multiple specific, actionable findings documented.
- 2 (Adequate): At least one finding documented.
- 1 (Inadequate): No findings generated (which itself is a finding -- the drill was too easy or the operator was not looking critically enough).

**Overall Score:**
- 13-15: Strong resilience posture. Maintain current procedures.
- 10-12: Adequate with room for improvement. Address findings within 30 days.
- 7-9: Significant gaps. Prioritize remediation. Schedule a follow-up drill within 60 days.
- 5-6: Critical deficiencies. Declare a Tier 3 governance action to address DR gaps immediately.

### 4.5 The After-Action Report

Every drill produces an after-action report (AAR). The AAR is a Tier 1 document per D6-001 and is stored in the documentation corpus.

**After-Action Report Template:**

```
DISASTER RECOVERY DRILL -- AFTER-ACTION REPORT

Drill ID:         DR-DRILL-[YYYY]-[sequential number]
Date:             [date of drill execution]
Scenario:         [scenario number and name]
Drill Type:       [tabletop / functional / partial recovery / full recovery]
Scheduled/Surprise: [scheduled / surprise]
Duration:         [actual time from start to completion]
Time Limit:       [defined time limit for this scenario]

SCENARIO DESCRIPTION:
[Brief description of the simulated disaster and any variants used]

EXECUTION SUMMARY:
[Narrative account of what happened during the drill, step by step.
Note any deviations from the documented procedure.]

EVALUATION SCORES:
Procedure Adherence:    [1-3] -- [brief justification]
Time Performance:       [1-3] -- [brief justification]
Data Integrity:         [1-3] -- [brief justification]
Documentation Quality:  [1-3] -- [brief justification]
Finding Generation:     [1-3] -- [brief justification]
Overall Score:          [5-15]

FINDINGS:
[Numbered list of specific findings. Each finding includes:]
  F1: [Description of the gap, error, or issue discovered]
      Severity: [Critical / Major / Minor]
      Affected Procedure: [Document ID and section]
      Recommended Action: [Specific action to address the finding]
      Deadline: [Date by which the action should be completed]

LESSONS LEARNED:
[What worked well. What would the operator do differently next time.
Any observations about the institution's overall resilience posture.]

FOLLOW-UP ITEMS:
[List of all action items generated by this drill, with owners and deadlines]

Signed: [Operator identifier]
Date:   [Date of report completion]
```

The AAR must be completed within 7 days of drill execution. Findings with severity "Critical" must have corrective actions initiated within 48 hours. Findings with severity "Major" must have corrective actions initiated within 30 days. Findings with severity "Minor" must be addressed before the next quarterly drill.

## 5. Rules & Constraints

- **R-D12-02-01:** Disaster recovery drills must be conducted at least quarterly per the schedule in Section 4.1. Missed drills must be rescheduled within 30 days and the miss documented in the operational log.
- **R-D12-02-02:** No drill may be conducted on production data or production systems without explicit safeguards that prevent data loss. The drill safety principle (Section 3.3) is absolute.
- **R-D12-02-03:** Every drill must produce an after-action report per Section 4.5. The AAR is a Tier 1 document.
- **R-D12-02-04:** Critical findings from drills must have corrective actions initiated within 48 hours. Corrective actions must be tracked to completion.
- **R-D12-02-05:** All five disaster scenarios must be exercised within every 15-month period. No scenario may go more than 18 months without being drilled.
- **R-D12-02-06:** At least one surprise drill must be conducted per year.
- **R-D12-02-07:** The drill schedule, scenario scripts, and evaluation framework must be reviewed during the annual operations cycle. Scenarios must be updated to reflect changes in the institution's systems and threat model.
- **R-D12-02-08:** Drill findings that reveal gaps in recovery procedures must be integrated into the affected procedures within 90 days. The integration must be verified in the next drill of the affected scenario.

## 6. Failure Modes

- **Drill theater.** The operator conducts drills perfunctorily, following the motions without genuine engagement. The same easy scenarios are repeated. Difficult scenarios are avoided. Findings are generic rather than specific. Mitigation: the evaluation framework scores finding generation. A drill that produces no findings scores 1 on that dimension, which should prompt self-examination. The surprise drill requirement disrupts complacency. The scenario rotation prevents cherry-picking.

- **Drill avoidance.** The operator postpones drills indefinitely. "Too busy this quarter" becomes a permanent state. Mitigation: R-D12-02-01 makes drills mandatory. The quarterly operations day in OPS-001 explicitly includes drill time. If drills are consistently skipped, this is a governance issue requiring Tier 3 review.

- **Drill-caused damage.** A drill inadvertently damages production data or systems. The drill becomes the disaster. Mitigation: the drill safety principle (Section 3.3) and R-D12-02-02. All drill procedures specify the use of test targets, copies, and spare hardware. Pre-drill setup includes verification of current backups.

- **Finding accumulation without action.** Drills produce findings, but the findings are never addressed. The AAR file grows; the procedures remain unchanged. Mitigation: R-D12-02-04 and R-D12-02-08 impose deadlines on corrective actions. The quarterly governance health check (GOV-007) reviews open drill findings.

- **Scenario staleness.** The five scenarios become outdated as the institution evolves. New systems are not covered. Old systems referenced in scenarios have been decommissioned. Mitigation: R-D12-02-07 requires annual review and update of scenarios.

## 7. Recovery Procedures

1. **If drills have lapsed for more than two quarters:** Do not attempt to catch up by drilling all scenarios at once. Begin with a tabletop drill of Scenario 1 (primary storage failure) to rebuild the drill habit. Schedule the remaining scenarios at two-week intervals until the rotation is current. Document the lapse in the operational log.

2. **If a drill causes production damage:** Stop the drill immediately. Assess the damage. If data loss has occurred, initiate the actual recovery procedure (not the drill procedure). Restore from the most recent verified backup. Document the incident as both a drill failure and a real incident. Investigate the cause: which safety control failed? Update drill procedures to prevent recurrence.

3. **If drill findings are never addressed:** Conduct a finding audit. List all open findings from the past 12 months. Categorize them by severity. Address all Critical findings immediately. Schedule Major findings into the next quarterly operations day. Assess whether Minor findings are still relevant. Close findings that have been overtaken by events. Document the audit results.

4. **If scenarios no longer match institutional reality:** Conduct an immediate scenario update. For each of the five scenarios, verify that the referenced hardware, software, backup locations, and procedures still exist. Update the scenario scripts. If a scenario is no longer relevant (for example, if the institution no longer uses the referenced encryption scheme), replace it with a scenario that reflects the current threat model.

5. **If drill scores are consistently low (below 10):** This indicates systemic problems with the institution's disaster recovery posture. Escalate to a Tier 3 governance decision. Conduct a comprehensive review of all recovery procedures, backup integrity, and spare hardware availability. Consider whether the institution's complexity has exceeded the operator's ability to maintain adequate disaster preparedness.

## 8. Evolution Path

- **Years 0-5:** The drill framework is being established. Early drills will be rough. Scores will be low. This is expected and healthy -- the institution is young and its recovery procedures are untested. Use this period to refine both the drill procedures and the recovery procedures they test. Expect the scenario scripts to be revised frequently.

- **Years 5-15:** The drill framework should be routine. Scores should generally be in the 10-15 range. The focus shifts from establishing the framework to preventing complacency. Introduce more challenging variants of existing scenarios. Consider adding new scenarios for threats that have emerged since founding.

- **Years 15-30:** Hardware generations have changed. Scenarios must evolve with them. The drill history -- the accumulated AARs -- becomes a valuable record of the institution's resilience evolution. New operators brought in through succession should participate in drills as part of their training (GOV-004 Stage 2 requires participation in at least one recovery procedure).

- **Years 30-50+:** The drill framework should have been refined through decades of practice. The AAR archive tells the story of every test, every finding, every improvement. This archive is institutional knowledge of the highest order -- it documents not just what the institution can survive, but how it learned to survive it.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators. Each entry should include the date, the author's identifier, and the context for the commentary.*

**2026-02-16 -- Founding Entry:**
Writing drill scripts for disasters that have not happened is an exercise in structured imagination. I find myself torn between making the drills realistic enough to be useful and simple enough to be executable by one person in a day. I have erred on the side of realism for the quarterly functional drills and simplicity for the monthly tabletops.

The scenario I am most anxious about is Scenario 4 -- the operator incapacitation simulation. Pretending to be a stranger to my own system is difficult. My hands know where things are even when my conscious mind pretends not to. I expect this scenario to produce the most uncomfortable findings, because it tests not just the systems but the documentation, and documentation gaps are the hardest to see from the inside.

The evaluation framework uses a three-point scale deliberately. A more granular scale would invite false precision. The question is not whether the drill scored 7.3 or 7.5 -- it is whether the institution's resilience posture is strong, adequate, or deficient. Three levels are enough to answer that question honestly.

## 10. References

- ETH-001 -- Ethical Foundations of the Institution (Principle 2: Integrity Over Convenience; Principle 6: Honest Accounting of Limitations)
- CON-001 -- The Founding Mandate (self-contained operation, documentation completeness)
- GOV-001 -- Authority Model (governance tiers for DR decisions, decision logging)
- SEC-001 -- Threat Model and Security Philosophy (threat categories that inform drill scenarios)
- OPS-001 -- Operations Philosophy (operational tempo, quarterly and annual cycles, sustainability requirement)
- D12-001 -- Disaster Recovery Philosophy (three pillars, graceful degradation, rebuild mindset, drill philosophy)
- D6-002 -- Backup Doctrine (backup procedures tested by drills)
- D6-003 -- Data Integrity Verification (integrity verification procedures tested by drills)
- SEC-002 -- Access Control Procedures (physical zone model referenced in drill scenarios)
- SEC-003 -- Cryptographic Key Management (key recovery procedures tested in Scenario 5)
- GOV-004 -- Succession Planning and Execution (succession readiness verified by Scenario 4)
- GOV-007 -- Institutional Health Assessment (quarterly review of drill findings)
- D12-003 -- System Rebuild Procedures (rebuild procedures validated by full recovery drills)

---

---

# D12-003 -- System Rebuild Procedures

**Document ID:** D12-003
**Domain:** 12 -- Disaster Recovery
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, SEC-001, OPS-001, D12-001, D12-002, D6-002, SEC-002, SEC-003, D5-002
**Depended Upon By:** D13-002. All articles involving catastrophic recovery, hardware replacement, or institutional continuity after total loss.

---

## 1. Purpose

This article is the most important disaster recovery document in the institution. It defines how to rebuild the entire holm.chat Documentation Institution from scratch -- from bare metal to a fully operational system -- after catastrophic loss.

Catastrophic loss means that the primary site, primary hardware, or primary storage is gone. Not degraded, not partially damaged -- gone. The operator is standing in front of new (or different) hardware with nothing but backup media, printed documentation, and the knowledge in their head. This article tells them what to do.

D12-001 introduced the rebuild mindset: the practice of designing every system with reconstruction in mind. This article is the operational expression of that mindset. It is the procedure that the rebuild mindset was designed to make possible. If D12-001 is the philosophy, this article is the playbook.

This article is deliberately written to be understandable by a person who did not build the original system. The succession scenario (GOV-004) includes situations where a new operator must rebuild the institution from documentation alone. This article must serve that reader as well as the founding operator.

## 2. Scope

**In scope:**
- The complete rebuild sequence from bare hardware to operational institution.
- Hardware assessment and procurement for replacement systems.
- Operating system installation and base configuration.
- Storage setup and encryption.
- Data restoration from backup media.
- Service and application configuration.
- Verification and validation of the rebuilt system.
- The rebuild kit: what it contains and how to maintain it.
- Rebuild timeline estimates and prioritization.

**Out of scope:**
- The philosophy of disaster recovery (D12-001).
- Specific backup creation procedures (D6-002; this article consumes backups, it does not create them).
- Specific encryption key generation (SEC-003; this article uses existing keys from backup).
- Ongoing operations after rebuild (OPS-001 and Domain 10 articles).
- Drill procedures for testing the rebuild (D12-002).

## 3. Background

### 3.1 The Rebuild Scenario

The rebuild scenario is the worst case that the institution is designed to survive. Everything at the primary site is lost. The operator must reconstruct the institution on new hardware, potentially in a new location, using only what was stored off-site.

This scenario is unlikely in any given year. Over fifty years, it is not unlikely at all. Fires, floods, earthquakes, theft, and structural failures are all plausible mechanisms. SEC-001 Section 4.2 rates physical threats as "near certain for at least one instance" over the fifty-year horizon. The rebuild procedure is the institution's response to that certainty.

### 3.2 The Minimum Viable Rebuild

Not everything needs to be rebuilt simultaneously. The institution can operate at reduced capability (D12-001 Section 4.2, graceful degradation) while the rebuild progresses. The rebuild sequence is therefore ordered by criticality: the most essential functions are restored first, and less critical functions are added as time and resources permit.

The minimum viable rebuild is: one functioning machine, with encrypted storage, containing the restored documentation corpus, the decision log, the operational log, and the backup system. From this minimum, the operator can consult the documentation to rebuild everything else. The documentation is the seed. Everything else grows from it.

### 3.3 The Hardware Question

The original hardware will not always be available for a rebuild. If the primary site is destroyed, the hardware is destroyed with it. Even if the loss is less dramatic -- a single server failure, for instance -- the exact same hardware may no longer be manufactured.

This article is therefore hardware-agnostic at the architectural level. The procedures describe what capabilities the hardware must provide, not which specific products to buy. The specific hardware choices are documented in the institutional asset registry (D11-002) and the procurement plan (subsequent Domain 11 articles), which are updated as hardware generations change.

## 4. System Model

### 4.1 The Rebuild Kit

The rebuild kit is the minimum set of artifacts needed to reconstruct the institution on new hardware. It is stored off-site (Zone 3 per SEC-002) and maintained per D12-001 R-DR-04. It contains:

**Physical components:**
- Bootable installation media for the institution's operating system (USB drive or optical disc, updated annually).
- Off-site backup media (Location C per D6-002) containing the most recent Tier 1, Tier 2, and Tier 3 backups.
- Encryption key backup (Tier C per SEC-003) with the sealed passphrase envelope.
- Printed copy of this article (D12-003).
- Printed copy of the rebuild checklist (Section 4.7 of this article).
- Printed copy of the system configuration summary (a condensed version of the institutional architecture maintained as part of the rebuild kit, updated after every significant system change).

**Documentation components (stored on the backup media):**
- The complete documentation corpus.
- The complete decision log and operational log.
- System configuration files and service definitions.
- The asset registry (D11-002).
- All recovery procedures.

**Maintenance requirements:**
- The rebuild kit must be updated after every significant system change (Tier 2 or Tier 3 decision that modifies system architecture).
- The bootable installation media must be refreshed at least annually with the current OS version.
- The printed materials must be reprinted whenever the underlying documents change.
- The rebuild kit must be verified annually during the annual operations cycle by confirming that all components are present, readable, and current.

### 4.2 Phase 1: Assessment and Procurement (Day 1-7)

Before rebuilding, assess what is available and what is needed.

**Step 1: Assess surviving assets.**
- What hardware survived, if any? Is it functional?
- Is the rebuild kit intact? Verify all components are present.
- Are the off-site backups accessible? Can the backup media be read?
- Are the encryption keys accessible? Can the Tier C backup be decrypted?
- Is the printed documentation available?

**Step 2: Assess the rebuild site.**
- Where will the rebuilt institution be located?
- Is adequate power available (considering off-grid requirements per CON-001)?
- Is the physical space suitable for equipment (temperature, humidity, physical security)?
- Can Zone 1, Zone 2, and Zone 3 (SEC-002) be established at the new site?

**Step 3: Determine hardware requirements.**
The rebuilt institution requires, at minimum:
- One general-purpose computer capable of running the institution's operating system. Minimum specifications: 64-bit processor, 16GB RAM (32GB preferred), sufficient storage for the institution's data (consult the asset registry for current data volume plus 50% growth margin).
- Primary storage: one or more drives with capacity for the full institutional data set. SSDs preferred for primary storage; HDDs acceptable.
- Backup storage: at least one separate drive for local backups (Location A per D6-002).
- Removable media reader compatible with the backup media format (USB, optical, etc.).
- UPS (uninterruptible power supply) or battery system sufficient for orderly shutdown.
- Monitor, keyboard, and input devices.

**Step 4: Procure missing hardware.**
If hardware must be purchased, prioritize:
1. The primary compute platform (the machine itself).
2. Primary storage.
3. Backup storage.
4. UPS/power protection.
5. Peripheral devices.

Use the procurement procedures in subsequent Domain 11 articles. If those procedures are not available (because the documentation is being restored), apply these principles: buy from reputable vendors, prefer widely-available hardware over exotic hardware, prefer hardware with published specifications and open-source driver support, and retain all receipts and documentation for the asset registry.

### 4.3 Phase 2: Operating System Installation (Day 2-3)

**Step 1: Prepare the installation media.**
If the rebuild kit contains bootable installation media, use it. If the media is damaged or outdated, a new installation image must be obtained. In an air-gapped institution, this means physically transporting installation media from an external source. Follow the quarantine procedures in Domain 18 for any media entering the institution from external sources.

**Step 2: Boot from installation media.**
- Configure the new hardware's BIOS/UEFI to boot from the installation media (USB or optical).
- Disable all network interfaces in BIOS/UEFI (enforce the air gap at the hardware level per SEC-001 R-SEC-01).
- Disable wireless radios if present (Wi-Fi, Bluetooth).
- Boot the installation media.

**Step 3: Partition and encrypt the storage.**
- Create the partition layout. The recommended layout is:
  - /boot -- 1GB, unencrypted (contains the bootloader and kernel).
  - / (root) -- remainder of the primary drive, LUKS-encrypted.
  - If a separate data drive exists: /data -- entire drive, LUKS-encrypted with a separate key.
- Initialize LUKS encryption on the root partition:
  ```
  cryptsetup luksFormat /dev/[root-partition]
  ```
  Use the disk encryption key from the rebuild kit (Tier C backup). If the original key is being reused, enter it when prompted. If a new key is being generated (because the original is compromised or unavailable), follow SEC-003 Section 4.2.
- Open the encrypted partition:
  ```
  cryptsetup open /dev/[root-partition] crypt-root
  ```
- Format the partitions:
  ```
  mkfs.ext4 /dev/[boot-partition]
  mkfs.ext4 /dev/mapper/crypt-root
  ```

**Step 4: Install the operating system.**
- Follow the installer prompts, targeting the encrypted partitions.
- Set the hostname to the institution's standard hostname (documented in the system configuration summary in the rebuild kit).
- Set the timezone.
- Create the root account with a strong password per SEC-003.
- Do not enable any network services. The institution is air-gapped.

**Step 5: Post-installation base configuration.**
- Boot into the newly installed system.
- Verify the air gap: `ip link show` -- all network interfaces should be DOWN or absent.
- Set the default umask to 077 per SEC-002 R-SEC2-04:
  Add `umask 077` to /etc/profile and /etc/bash.bashrc (or equivalent for the installed shell).
- Configure the system clock. If no NTP is available (air-gapped), set the time manually and document the time source.
- Update the bootloader configuration if necessary.

### 4.4 Phase 3: Account Setup (Day 3)

Create the account architecture defined in SEC-002 Section 4.1:

**Step 1: Create the administrative account.**
```
useradd -m -s /bin/bash -G sudo [admin-username]
passwd [admin-username]
```
Set a strong, unique password per SEC-003.

**Step 2: Create the operator account.**
```
useradd -m -s /bin/bash [operator-username]
passwd [operator-username]
```
Set a strong, unique password.

**Step 3: Configure sudo.**
Ensure the administrative account requires password authentication for every sudo invocation (do not configure NOPASSWD).

**Step 4: Create service accounts as needed.**
Service accounts are created during Phase 5 (service configuration) as each service is restored.

**Step 5: Verify account architecture.**
- Log in as each account and verify access levels.
- Verify that the operator account cannot perform administrative actions without sudo.
- Verify that the root account prompt is visually distinctive (configure PS1 in root's .bashrc).

### 4.5 Phase 4: Data Restoration (Day 3-5)

**Step 1: Prepare the backup media.**
- Connect the off-site backup media (Location C).
- Decrypt the backup media using the Tier C passphrase from the sealed envelope:
  ```
  cryptsetup open /dev/[backup-device] crypt-backup
  mount /dev/mapper/crypt-backup /mnt/backup
  ```

**Step 2: Verify backup integrity.**
- Check the backup manifest:
  ```
  cat /mnt/backup/MANIFEST.txt
  ```
- Spot-check checksums on several files to confirm the backup is readable and uncorrupted.

**Step 3: Restore Tier 1 data first.**
Tier 1 data is the institutional memory: the documentation corpus, decision log, operational log, key inventory, and governance records.
```
mkdir -p /data/tier1
tar -xzf /mnt/backup/tier1-[date].tar.gz -C /data/tier1/
```
Verify:
- Compare file count against the manifest.
- Spot-check critical files (open the documentation index, read the most recent decision log entry).
- If checksum manifests are available, verify:
  ```
  cd /data/tier1 && sha256sum -c /mnt/backup/tier1-[date].sha256
  ```

**Step 4: Restore Tier 2 data.**
Tier 2 data includes system configurations, active project data, and automation scripts.
```
mkdir -p /data/tier2
tar -xzf /mnt/backup/tier2-[date].tar.gz -C /data/tier2/
```
Verify as with Tier 1.

**Step 5: Restore Tier 3 data.**
Tier 3 data includes reference materials, research archives, and imported knowledge bases.
```
mkdir -p /data/tier3
tar -xzf /mnt/backup/tier3-[date].tar.gz -C /data/tier3/
```
Verify as with Tier 1.

**Step 6: Set ownership and permissions.**
```
chown -R [operator-username]:[operator-username] /data/tier1 /data/tier2 /data/tier3
chmod -R 700 /data/tier1 /data/tier2 /data/tier3
```
Adjust group ownership for files that need to be accessible to service accounts per SEC-002 Section 4.2.

**Step 7: Unmount and secure the backup media.**
```
umount /mnt/backup
cryptsetup close crypt-backup
```
Return the backup media to secure storage.

### 4.6 Phase 5: Service Configuration (Day 5-10)

With the OS installed and data restored, configure the services that make the institution operational.

**Step 1: Consult the restored documentation.**
Open the documentation corpus (now restored in /data/tier1). Locate the system configuration documentation. This documentation, maintained as part of the institutional corpus, describes every service running on the institution's systems, its purpose, its configuration, and its dependencies.

**Step 2: Prioritize service restoration.**
Restore services in this order:
1. Backup system (D6-002) -- so that the rebuilt system is immediately protected.
2. Integrity verification (D6-003) -- so that restored data can be verified.
3. Documentation access tools -- so that the documentation corpus is navigable.
4. Operational logging -- so that the rebuild itself is documented.
5. All remaining services in order of criticality as defined in the system configuration documentation.

**Step 3: For each service:**
a. Create the service account if required per SEC-002 Section 4.1.
b. Install the service software from the rebuild kit or from packages on the installation media.
c. Apply the service configuration from the restored Tier 2 data or from the documented configuration.
d. Test the service.
e. Document any deviations from the original configuration (new hardware, different paths, etc.).

**Step 4: Restore automation.**
Restore cron jobs, scheduled tasks, and automated monitoring per the documentation. Verify each automated task by running it manually once before enabling the schedule.

### 4.7 Phase 6: Verification and Validation (Day 10-14)

**Step 1: System verification checklist.**
Execute the following checks and record the results:

```
SYSTEM REBUILD VERIFICATION CHECKLIST

Date: _______________
Operator: _______________
Rebuild start date: _______________
Rebuild complete date: _______________

HARDWARE:
[ ] All hardware functional and documented in asset registry
[ ] Air gap verified -- no active network interfaces
[ ] UPS/power protection functional
[ ] Physical security established (Zone 1, Zone 2, Zone 3)

OPERATING SYSTEM:
[ ] OS installed and bootable
[ ] Full-disk encryption active and verified
[ ] Account architecture matches SEC-002 (root, admin, operator, service accounts)
[ ] Default umask set to 077
[ ] System clock accurate

DATA:
[ ] Tier 1 data restored and verified
[ ] Tier 2 data restored and verified
[ ] Tier 3 data restored and verified
[ ] File permissions correct
[ ] Checksums verified against backup manifests

SERVICES:
[ ] Backup system operational -- first backup created on rebuilt system
[ ] Integrity verification operational
[ ] Documentation access functional
[ ] Operational logging functional
[ ] All critical services operational

SECURITY:
[ ] Encryption keys installed and functional
[ ] Key inventory accurate for rebuilt system
[ ] Access controls verified
[ ] Physical security zones established

DOCUMENTATION:
[ ] Rebuild process documented in operational log
[ ] Any deviations from this procedure documented
[ ] Asset registry updated for new hardware
[ ] Decision log updated with rebuild decision (Tier 2)
```

**Step 2: First backup of rebuilt system.**
Immediately after verification, create a full backup of the rebuilt system per D6-002. This backup becomes the new baseline. Store copies at Locations A and B. Schedule the next off-site transport.

**Step 3: Operational log entry.**
Record the rebuild in the operational log:
- Date of the disaster or event that triggered the rebuild.
- Date rebuild started and completed.
- Hardware used (make, model, specifications).
- Backup media used (date of backup, location source).
- Any data that could not be restored (documenting the loss per ETH-001 Principle 6).
- Any deviations from this procedure and the reasons.
- Total rebuild time.

**Step 4: Decision log entry.**
Record the rebuild as a Tier 2 decision in the decision log per GOV-001.

### 4.8 Rebuild Timeline Summary

| Phase | Activities | Estimated Duration |
|-------|-----------|-------------------|
| 1. Assessment & Procurement | Assess surviving assets, evaluate rebuild site, determine and procure hardware | 1-7 days (hardware-dependent) |
| 2. OS Installation | Prepare media, partition, encrypt, install, base configuration | 1 day |
| 3. Account Setup | Create account architecture per SEC-002 | 0.5 day |
| 4. Data Restoration | Restore Tier 1, 2, 3 from backup, verify | 2-3 days |
| 5. Service Configuration | Restore and configure all services | 3-5 days |
| 6. Verification | Checklist, first backup, logging | 1-2 days |
| **Total** | | **8-18 days** |

The timeline assumes hardware is available. If hardware must be procured, add procurement lead time. The institution is operational at minimum viable capability (documentation accessible, backup system running) by the end of Phase 4 -- approximately 4-6 days after hardware is available.

## 5. Rules & Constraints

- **R-D12-03-01:** The rebuild kit must be maintained at an off-site location and updated after every significant system change per D12-001 R-DR-04.
- **R-D12-03-02:** The rebuild kit must be verified annually by confirming all components are present, readable, and current.
- **R-D12-03-03:** This article (D12-003) must be stored in printed form in the rebuild kit. The printed copy must be updated whenever this article is amended.
- **R-D12-03-04:** A full or partial rebuild drill per D12-002 must be conducted at least annually to validate these procedures.
- **R-D12-03-05:** The system configuration summary in the rebuild kit must be updated after every Tier 2 or Tier 3 decision that modifies the institutional system architecture.
- **R-D12-03-06:** Data restoration must follow tier priority: Tier 1 first, then Tier 2, then Tier 3. No Tier 2 or 3 data restoration may delay or compromise Tier 1 restoration.
- **R-D12-03-07:** The rebuilt system must pass the verification checklist (Section 4.7) before being declared operational.

## 6. Failure Modes

- **Rebuild kit decay.** The rebuild kit is created once and never updated. When needed, it contains outdated OS installation media, stale backup media, and printed documentation that references systems and configurations that no longer exist. Mitigation: R-D12-03-01 and R-D12-03-02 require regular updates and annual verification.

- **Hardware incompatibility.** The available replacement hardware is incompatible with the OS version on the installation media, or with the system configurations in the backup. Drivers are missing. Hardware-specific configurations do not apply. Mitigation: the rebuild kit should include the most recent stable OS version, which will have the broadest hardware support. The procedures in Section 4.3 are hardware-agnostic. The annual rebuild kit refresh keeps the installation media current.

- **Backup media failure.** The off-site backup media is unreadable -- physical damage, media degradation, or encryption key problems. Mitigation: D6-002 requires periodic verification of backup media. SEC-003 ensures encryption keys exist in multiple copies. The annual rebuild kit verification includes a spot-check of backup media readability.

- **Documentation dependency loop.** The rebuild procedure references documentation that is stored in the backup that has not yet been restored. The operator cannot follow the procedure because the information it references is not yet accessible. Mitigation: this article and the rebuild checklist are printed and stored in the rebuild kit. The most critical information for the rebuild is in these printed materials, not in the backup. Phase 4 (data restoration) makes the full documentation available.

- **Skill atrophy.** The operator has not performed a bare-metal installation in years. The procedures are documented but the operator's hands are unpracticed. Steps that should take minutes take hours. Mitigation: D12-002 requires annual rebuild drills that exercise these procedures.

- **Scope paralysis.** The operator, facing the enormity of rebuilding everything, is paralyzed by the scale of the task. They do not know where to start, or they try to rebuild everything simultaneously. Mitigation: the phased approach in this article provides a clear sequence. Phase 1 assesses. Phase 2 installs the OS. Phase 3 creates accounts. Phase 4 restores data. Each phase has a defined scope and a defined deliverable. Follow the phases in order.

## 7. Recovery Procedures

The recovery procedures for this article address failures of the rebuild process itself.

1. **If the rebuild kit is missing or destroyed:** The institution must be rebuilt from whatever survives. If any backup media exists at any location (Locations A, B, or C per D6-002), the data can be restored. If no backup media survives, the institution has suffered total data loss. In this case, the documentation that exists only in human memory and in any printed materials must be used to rebuild the institutional framework from scratch. This is the worst-case scenario. It is survivable only if the operator (or successor) remembers or has access to the fundamental architecture of the institution. Begin with a fresh OS installation and start rebuilding the documentation corpus from memory and available references.

2. **If the OS installation fails:** Troubleshoot the installation. Common causes: incompatible hardware (try different hardware), corrupt installation media (obtain new media), insufficient storage (use a larger drive). If the institution's standard OS cannot be installed on available hardware, any compatible Unix-like operating system can serve as a platform. The institutional procedures are designed to be OS-independent at the application level. Document the OS substitution and update the system configuration documentation accordingly.

3. **If data restoration from backup fails:** Attempt restoration from an alternate backup (a different date, a different location). If the backup archive is partially corrupt, extract what can be extracted and document the loss. If the backup encryption key does not work, follow SEC-003 recovery procedures for key loss. If all backup copies are unreadable, the data protected by those backups is lost. Document the loss honestly per ETH-001 Principle 6.

4. **If service configuration fails:** Consult the restored documentation for the service's configuration guide. If the configuration assumes hardware or software that is no longer available, adapt the configuration to the current platform. Document all adaptations. If a service cannot be restored on the new platform, assess whether it is critical to the institution's mission. Non-critical services can be deferred; critical services must be replaced with functional equivalents.

## 8. Evolution Path

- **Years 0-5:** These procedures are being written against the institution's founding hardware and software. The first rebuild drill will likely reveal gaps -- steps that are too vague, assumptions that are undocumented, tools that are referenced but not included in the rebuild kit. Expect significant revision after the first annual drill.

- **Years 5-15:** Hardware generations begin to change. The rebuild procedures must evolve to accommodate new hardware while maintaining the same architectural principles. Each hardware change triggers a rebuild kit refresh and a procedure review.

- **Years 15-30:** The institution may have undergone one or more real rebuilds by this point. The lessons from those rebuilds -- captured in the operational log and the after-action reports -- are the most valuable updates to these procedures. Real experience always surpasses planned procedures.

- **Years 30-50+:** The technology landscape is unrecognizable compared to founding. The specific commands in this article are almost certainly outdated. The sequence -- assess, install, configure, restore, verify -- is not. Future operators should follow the sequence, adapt the commands, and update this article.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
This article is the institution's answer to the question: what survives the fire? Not the hardware. Not the building. Not even the operator, necessarily. What survives is what was stored off-site: the backups, the keys, the printed documentation, and the knowledge encoded in all three.

I have designed the rebuild process to be executable by someone who is not me. This is deliberate. If I am the one rebuilding after a disaster, I will have knowledge beyond what is in these documents and the rebuild will go faster. But if a successor is rebuilding -- if I am the disaster the institution must survive -- these documents must be sufficient on their own.

The eight-to-eighteen-day timeline may seem slow. In a world of cloud computing, where a server can be provisioned in minutes, two weeks to rebuild a personal institution seems archaic. But this institution is not a cloud server. It is an air-gapped, encrypted, documented, self-sovereign system built to last fifty years. Rebuilding it carefully in two weeks is preferable to rebuilding it hastily in two days and discovering three months later that something critical was missed.

## 10. References

- ETH-001 -- Ethical Foundations of the Institution (Principle 2: Integrity; Principle 6: Honest Accounting of Limitations)
- CON-001 -- The Founding Mandate (air-gap, off-grid, self-built requirements)
- GOV-001 -- Authority Model (Tier 2 decisions for rebuild, decision logging)
- SEC-001 -- Threat Model and Security Philosophy (air-gap enforcement, physical threats)
- SEC-002 -- Access Control Procedures (account architecture, physical zones)
- SEC-003 -- Cryptographic Key Management (key recovery from Tier C backup)
- OPS-001 -- Operations Philosophy (documentation-first principle, operational tempo post-rebuild)
- D12-001 -- Disaster Recovery Philosophy (rebuild mindset, rebuild kit concept)
- D12-002 -- Disaster Recovery Drill Framework (rebuild drills)
- D6-002 -- Backup Doctrine (backup locations, backup media, restoration procedures)
- D6-003 -- Data Integrity Verification (post-restoration verification)
- D5-002 -- Operating System Maintenance Procedures (OS installation and configuration)
- D11-002 -- Asset Inventory and Tracking (asset registry for replacement hardware)
- GOV-004 -- Succession Planning and Execution (successor rebuild scenario)

---

---

# D11-002 -- Asset Inventory and Tracking

**Document ID:** D11-002
**Domain:** 11 -- Administration
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, OPS-001, D11-001, SEC-002, D6-002, D12-001
**Depended Upon By:** D12-002, D12-003, D13-002. All articles involving hardware replacement, procurement planning, disaster recovery asset assessment, or technology migration.

---

## 1. Purpose

This article defines the complete asset inventory system for the holm.chat Documentation Institution. It establishes how every physical and significant digital asset is catalogued, tracked through its lifecycle, audited for accuracy, and planned for eventual replacement.

D11-001 established that administration is the practice of knowing what you have, where it is, what condition it is in, and what it needs. This article operationalizes that principle into a specific system: the asset registry. The asset registry is the single authoritative record of every asset the institution owns, has owned, or plans to acquire. It is the material memory of the institution -- the counterpart to the decision log's governance memory and the documentation corpus's knowledge memory.

Without the asset registry, disaster recovery is guesswork (D12-003 cannot specify replacement hardware if the current hardware is undocumented), procurement is reactive rather than planned (you discover you need a replacement when the current component fails), and succession is incomplete (the successor inherits hardware they cannot identify or assess).

This article is written for the operator who must build the registry from scratch, maintain it over years, and hand it to a successor who has never seen the equipment.

## 2. Scope

**In scope:**
- The asset registry: its structure, fields, and storage.
- Asset categories: hardware, storage media, consumables, tools, infrastructure.
- The asset lifecycle: acquisition, registration, deployment, maintenance, depreciation, retirement, disposal.
- Audit procedures: how to verify the registry matches physical reality.
- Replacement planning: how to anticipate and prepare for end-of-life.
- Templates for the asset registry and related documents.
- Integration with disaster recovery, procurement, and evolution planning.

**Out of scope:**
- Budget and financial accounting for assets (subsequent Domain 11 articles).
- Procurement procedures for new assets (subsequent Domain 11 articles).
- Digital asset management for data and documents (Domain 6).
- Software inventory (Domain 5).

## 3. Background

### 3.1 Why Asset Tracking Matters in an Air-Gapped Institution

In a conventional technology environment, asset management is a convenience -- it helps organizations optimize spending and plan refreshes, but a missing entry in the asset database rarely causes a crisis. In an air-gapped, off-grid institution, asset management is a survival function.

The air gap means that replacement parts cannot be ordered for next-day delivery. The off-grid nature means that power and environmental systems require consumables (batteries, filters, fuel) that must be stockpiled and tracked. The fifty-year horizon means that hardware will go through multiple generations, and the history of what was replaced, when, and why becomes essential context for future replacement decisions.

When the primary storage array begins showing SMART warnings, the operator needs to know immediately: what drives are installed, when were they purchased, what is their warranty status, where are the spare drives, and are those spares compatible with the current array? The asset registry answers all of these questions. Without it, the operator is searching shelves and guessing.

### 3.2 The Lifecycle Perspective

Assets are not static. They move through a predictable lifecycle: they are acquired, they are deployed, they perform, they age, they degrade, and they are eventually retired and replaced. Effective asset management tracks this entire lifecycle, not just the current state.

The lifecycle perspective transforms asset management from a static inventory ("what do I have?") into a dynamic planning tool ("what will I need, and when?"). If the registry records not just what drives are installed but when they were installed and what their expected lifespan is, the operator can predict -- years in advance -- when replacements will be needed and begin procurement planning accordingly.

In an institution where procurement may require significant lead time (because of the air gap, because of off-grid logistics, because of budget constraints), the difference between planned procurement and emergency procurement is often the difference between a smooth transition and a crisis.

### 3.3 The Relationship to Other Domains

The asset registry is a dependency for multiple domains:
- **Domain 12 (Disaster Recovery):** D12-003 needs the registry to determine replacement hardware specifications. D12-002 uses it to verify spare hardware availability for drills.
- **Domain 13 (Evolution):** D13-002 uses the registry to plan technology migrations. The registry's lifecycle data informs sunset and replacement timing.
- **Domain 3 (Security):** SEC-002 references physical zones. The registry documents what hardware is in which zone.
- **Domain 6 (Data):** D6-002 tracks backup media. The registry tracks the physical media assets.

## 4. System Model

### 4.1 The Asset Registry Structure

The asset registry is a structured document (plain text or simple structured data per R-ADM-05) containing one entry per asset. Each entry contains the following fields:

**Required fields (all assets):**

```
ASSET REGISTRY ENTRY

Asset ID:           [Unique identifier: AST-[YYYY]-[sequential number]]
Category:           [Hardware / Storage Media / Consumable / Tool / Infrastructure]
Description:        [Plain-language description of the asset]
Make/Model:         [Manufacturer and model number, if applicable]
Serial Number:      [Manufacturer serial number, if applicable]
Acquisition Date:   [Date the asset was acquired]
Acquisition Source: [Where/from whom the asset was obtained]
Acquisition Cost:   [Cost at acquisition, if applicable]
Location:           Zone __ / [specific placement]
Status:             [Active / Spare / Maintenance / Deprecated / Retired / Disposed]
Condition:          [Excellent / Good / Fair / Degraded / Failed]
Purpose:            [What the asset is used for in the institution]
Dependencies:       [What other assets or systems depend on this asset]
Expected Lifespan:  [Estimated useful life from acquisition date]
Replacement Date:   [Projected date for replacement based on expected lifespan]
Maintenance Schedule: [What maintenance this asset requires and how often]
Last Maintenance:   [Date of the most recent maintenance action]
Notes:              [Any additional information relevant to this asset]
```

**Additional fields for specific categories:**

For **hardware assets** (servers, drives, peripherals):
```
Specifications:     [CPU, RAM, storage capacity, interfaces, etc.]
Firmware Version:   [Current firmware version, if applicable]
Warranty Expiry:    [Date warranty expires, if applicable]
SMART Data:         [Most recent SMART status for storage devices]
Power Requirements: [Wattage, voltage, connector type]
```

For **storage media** (backup drives, optical discs, USB drives):
```
Capacity:           [Storage capacity]
Format:             [Filesystem or encryption format]
Contents:           [Brief description of what is stored]
Write Count:        [Number of times written, for flash media]
Last Verified:      [Date media was last verified readable]
Encryption:         [Yes/No, key reference if encrypted]
```

For **consumables** (batteries, filters, cables, cleaning supplies):
```
Quantity On Hand:   [Current count]
Minimum Stock Level: [Reorder point -- when to procure more]
Consumption Rate:   [Average usage per month/quarter/year]
Shelf Life:         [Expiration date or storage limit, if applicable]
```

For **infrastructure** (power systems, cooling, furniture, enclosures):
```
Installation Date:  [When the infrastructure was installed]
Capacity:           [Power capacity, cooling capacity, load rating, etc.]
Certification:      [Any safety certifications or ratings]
Inspection Schedule: [Frequency of required inspections]
Last Inspection:    [Date of most recent inspection]
```

### 4.2 Asset Categories

**Category 1: Hardware.**
Computing equipment: servers, workstations, monitors, keyboards, mice, USB hubs, display adapters. Storage devices that are permanently installed: internal drives (HDD, SSD), RAID controllers, drive bays.

**Category 2: Storage Media.**
Removable or semi-permanent storage: external HDDs, USB flash drives, optical discs (CD/DVD/Blu-ray), tape media. Backup media per D6-002. Media used for the rebuild kit per D12-003.

**Category 3: Consumables.**
Items that are used up or wear out through normal operation and must be replaced: batteries (UPS, portable), cables (data, power), thermal paste, cleaning supplies, printer consumables (if applicable), optical media (blank discs), pens and paper for printed documentation.

**Category 4: Tools.**
Items used for maintenance and repair: screwdrivers, anti-static equipment, multimeters, crimping tools, label makers, storage containers and organizers.

**Category 5: Infrastructure.**
Physical plant and environmental: UPS/battery systems, power distribution, cooling equipment, server racks or enclosures, physical security hardware (locks, keys, cabinets), furniture (desks, chairs), environmental monitoring (thermometers, hygrometers).

### 4.3 The Asset Lifecycle

**Stage 1: Acquisition.**
When a new asset enters the institution:
1. Assign an Asset ID using the format AST-[YYYY]-[sequential number].
2. Complete all required fields in the registry entry.
3. Physically label the asset with its Asset ID. Use durable labels (printed, not handwritten) that will remain legible for the asset's expected lifespan.
4. Photograph the asset (optional but recommended for hardware and infrastructure).
5. Store any documentation, manuals, or receipts that came with the asset. Reference their storage location in the Notes field.
6. Record the acquisition in the operational log.

**Stage 2: Deployment.**
When an asset is put into active use:
1. Update the Status field from "Spare" to "Active."
2. Update the Location field to the asset's deployed position.
3. Update the Dependencies field to reflect what depends on this asset.
4. If the asset replaces a previous asset, update the previous asset's status to "Retired" and cross-reference the replacement in both entries.

**Stage 3: Maintenance.**
During the asset's active life:
1. Follow the maintenance schedule recorded in the registry.
2. After each maintenance action, update the Last Maintenance field.
3. Update the Condition field based on the maintenance assessment.
4. For storage devices, update SMART Data after each check.
5. For consumables, update Quantity On Hand after each use.

**Stage 4: Depreciation and Monitoring.**
As the asset ages:
1. When the asset reaches 75% of its expected lifespan, flag it for replacement planning.
2. Begin procurement research for a replacement per subsequent Domain 11 procurement articles.
3. Update the Condition field to reflect aging: "Good" may transition to "Fair" as the asset shows wear.
4. If the asset shows signs of impending failure (increasing SMART errors, reduced battery capacity, physical wear), update Status to "Deprecated" and accelerate replacement planning.

**Stage 5: Retirement.**
When an asset is removed from active service:
1. Update Status to "Retired."
2. Record the retirement date and reason in the Notes field.
3. If the asset contains data, ensure all data is migrated or backed up before retirement.
4. If the asset is being replaced, cross-reference the replacement asset.
5. Move the asset to a designated retired-assets storage area or proceed to disposal.

**Stage 6: Disposal.**
When a retired asset is permanently removed from the institution:
1. If the asset contained data at any point (drives, media), securely erase it per D6-002 Section 4.5.
2. Update Status to "Disposed."
3. Record the disposal method and date in the Notes field.
4. Retain the registry entry permanently. Do not delete registry entries for disposed assets -- they are part of the institutional record.

### 4.4 Audit Procedures

The asset registry must be periodically reconciled with physical reality. Audits verify that the registry is accurate and that no assets are missing, mislocated, or undocumented.

**Quarterly Quick Audit (during quarterly operations day per OPS-001):**
1. Select one asset category (rotate through all five categories over five quarters).
2. For every asset in that category with Status "Active" or "Spare":
   a. Physically locate the asset.
   b. Verify it is in the location recorded in the registry.
   c. Verify its condition matches the recorded condition.
   d. For storage media: verify the Last Verified date is within policy.
   e. For consumables: verify the Quantity On Hand is accurate and above the minimum stock level.
3. Record the audit results: date, category audited, discrepancies found, corrective actions taken.
4. Duration: 1-2 hours.

**Annual Comprehensive Audit (during annual operations cycle per OPS-001):**
1. Audit every asset in the registry, all categories.
2. For each asset:
   a. Physically locate and inspect the asset.
   b. Verify all registry fields are current and accurate.
   c. Assess condition and update if needed.
   d. Verify maintenance schedule is being followed.
   e. Assess whether the asset is approaching end-of-life.
3. Walk through every physical space (Zone 1, Zone 2, Zone 3) and verify that every physical asset is in the registry. Any unregistered asset found during the walkthrough must be registered immediately.
4. Review the replacement planning queue: which assets are flagged for replacement within the next 24 months?
5. Generate an Annual Asset Report summarizing: total assets by category, assets acquired this year, assets retired this year, assets flagged for replacement, consumables inventory status, overall inventory health assessment.
6. Duration: half day to full day.

### 4.5 Replacement Planning

Replacement planning is the forward-looking function of asset management. It uses lifecycle data to predict when assets will need replacement and initiates procurement before the need becomes urgent.

**The Replacement Planning Queue:**
Maintain a list of assets sorted by projected replacement date. This list is derived from the registry's Expected Lifespan and Replacement Date fields. Review the queue during each quarterly audit and the annual comprehensive audit.

**Replacement planning thresholds:**

| Time to Replacement | Action Required |
|---------------------|----------------|
| 24+ months | Awareness only. Asset is on the queue. |
| 12-24 months | Research phase. Identify suitable replacements. Document options. |
| 6-12 months | Procurement phase. Acquire the replacement per procurement procedures. |
| 0-6 months | Transition phase. Schedule the replacement. Test the new asset. Execute the swap. |
| Past due | Emergency. The asset has exceeded its expected lifespan. Prioritize replacement immediately. |

**Replacement procedure:**
1. Identify the replacement asset (already procured per the planning thresholds above).
2. Register the replacement in the asset registry (Stage 1: Acquisition).
3. Test the replacement asset to confirm it functions correctly.
4. Schedule the swap during a maintenance window.
5. If the swap involves data migration (storage devices), follow D13-002 migration procedures.
6. Execute the swap: deploy the new asset, retire the old asset.
7. Update the registry for both assets.
8. Verify all dependent systems function correctly with the new asset.
9. Document the replacement in the operational log.

### 4.6 Templates

**Asset Registry Entry Template:**

```
---
ASSET REGISTRY ENTRY
Asset ID:           AST-____-____
Category:           [Hardware / Storage Media / Consumable / Tool / Infrastructure]
Description:
Make/Model:
Serial Number:
Acquisition Date:   ____-__-__
Acquisition Source:
Acquisition Cost:
Location:           Zone __ / [specific placement]
Status:             [Active / Spare / Maintenance / Deprecated / Retired / Disposed]
Condition:          [Excellent / Good / Fair / Degraded / Failed]
Purpose:
Dependencies:
Expected Lifespan:
Replacement Date:   ____-__-__
Maintenance Schedule:
Last Maintenance:   ____-__-__
Notes:

[Category-specific fields as applicable]
---
```

**Quarterly Audit Report Template:**

```
ASSET INVENTORY -- QUARTERLY AUDIT REPORT

Date:               ____-__-__
Category Audited:
Total Assets in Category: ____
Assets Verified:    ____
Discrepancies Found: ____

DISCREPANCY DETAILS:
[For each discrepancy:]
  Asset ID:
  Issue:            [missing / mislocated / condition mismatch / data stale]
  Corrective Action:
  Status:           [resolved / pending]

CONSUMABLES STATUS (if auditing consumables):
  Items Below Minimum Stock:
  [List items and current quantities]

NOTES:

Signed: _______________
```

**Annual Asset Report Template:**

```
ASSET INVENTORY -- ANNUAL REPORT

Reporting Period:   ____-__-__ to ____-__-__
Total Assets:       ____
  Hardware:         ____
  Storage Media:    ____
  Consumables:      ____ (line items; total quantities vary)
  Tools:            ____
  Infrastructure:   ____

ACQUISITIONS THIS PERIOD:
  [List: Asset ID, Description, Date]

RETIREMENTS THIS PERIOD:
  [List: Asset ID, Description, Date, Reason]

DISPOSALS THIS PERIOD:
  [List: Asset ID, Description, Date, Method]

REPLACEMENT PLANNING QUEUE:
  Assets due for replacement within 12 months:
  [List: Asset ID, Description, Replacement Date, Status of procurement]

  Assets due for replacement within 12-24 months:
  [List: Asset ID, Description, Replacement Date]

CONSUMABLES INVENTORY HEALTH:
  Items below minimum stock level:
  [List: Item, Current quantity, Minimum level]

  Items approaching shelf life expiration:
  [List: Item, Expiration date, Quantity]

OVERALL ASSESSMENT:
  [Narrative: Is the institution's material inventory healthy?
  Are replacement plans on track? Are there any concerns?]

Signed: _______________
Date:   ____-__-__
```

## 5. Rules & Constraints

- **R-D11-02-01:** Every physical asset in the institution must be registered in the asset registry before it is deployed. No asset is considered operational until it is registered, per D11-001 R-ADM-01.
- **R-D11-02-02:** Every registered asset must be physically labeled with its Asset ID. Labels must be durable and legible for the expected lifespan of the asset.
- **R-D11-02-03:** The asset registry must be maintained as a Tier 1 document per D6-001. It is backed up with the same frequency and verification as all Tier 1 data.
- **R-D11-02-04:** Quarterly audits of one asset category must be conducted during each quarterly operations day per OPS-001. The annual comprehensive audit must be conducted during the annual operations cycle.
- **R-D11-02-05:** Asset registry entries are permanent. Entries for retired or disposed assets must not be deleted. They are part of the institutional record and provide historical context for future procurement and replacement decisions.
- **R-D11-02-06:** Replacement planning must begin at least 12 months before an asset's projected end-of-life, per the planning thresholds in Section 4.5. Assets that exceed their expected lifespan without a replacement plan must be flagged as a governance finding.
- **R-D11-02-07:** The asset registry must be stored in plain text or simple structured data formats per D11-001 R-ADM-05. No proprietary software may be required to read the registry.
- **R-D11-02-08:** Any change to an asset's status, location, or condition must be recorded in the registry within 24 hours of the change occurring.

## 6. Failure Modes

- **Registry-reality divergence.** The registry says one thing; the shelves say another. Assets are moved without updating the registry. New assets are deployed without registration. Over time, the registry becomes fiction. Mitigation: quarterly audits catch divergence early. R-D11-02-08 requires prompt updates. The physical labeling requirement (R-D11-02-02) creates a visible link between assets and their registry entries.

- **Registration fatigue.** The operator grows tired of registering every cable, every USB drive, every pack of batteries. Registration is skipped for "minor" assets. Over time, the definition of "minor" expands until only the most expensive hardware is tracked. Mitigation: the asset categories in Section 4.2 define what must be tracked. Consumables tracking is lighter-weight than hardware tracking (quantity counts rather than individual entries). The Lightness Principle from D11-001 applies: track what matters for decisions, not every paperclip.

- **Replacement planning neglect.** The planning thresholds are ignored. Assets are replaced only when they fail. Every replacement is an emergency. Mitigation: the replacement planning queue is reviewed during every quarterly audit. R-D11-02-06 requires planning to begin at least 12 months before projected end-of-life.

- **Historical amnesia.** Retired and disposed asset entries are deleted to "keep the registry clean." The institution loses the historical record of what hardware it has used, what failed, and what was replaced. Future procurement decisions are made without the benefit of past experience. Mitigation: R-D11-02-05 prohibits deletion of registry entries. The registry is append-only for entries, though individual entry fields can be updated.

- **Over-tracking.** The registry becomes so detailed that maintaining it consumes disproportionate time. Every USB cable has a serial number entry with twelve fields. Mitigation: the Lightness Principle (D11-001 Section 4.2). The test for every tracked field is: does this field inform a decision? If a cable's serial number never informs any decision, do not track serial numbers for cables.

## 7. Recovery Procedures

1. **If the asset registry has been lost:** This is a significant institutional loss but not a catastrophic one. Conduct a physical audit of all assets currently in the institution. Create new registry entries for every asset found. Mark all entries as "RECONSTRUCTED -- original acquisition data unavailable" where acquisition details cannot be determined. Estimate lifespan and replacement dates based on current condition assessment. Going forward, maintain the registry per these procedures.

2. **If the registry has diverged from reality:** Declare a reconciliation sprint per D11-001 Recovery Procedure 3. Trust physical reality over the registry. Walk through every zone. For every physical asset found, verify its registry entry. For every registry entry marked "Active," verify the physical asset exists. Correct all discrepancies. Investigate the cause of the divergence and add a procedural safeguard.

3. **If replacement planning has been neglected:** Review the registry for every asset with Status "Active." Calculate the remaining lifespan for each asset (Expected Lifespan minus time since Acquisition Date). Identify all assets that have less than 12 months of remaining lifespan. For those assets, initiate immediate procurement research. Prioritize by criticality: assets whose failure would cause data loss or service interruption are replaced first.

4. **If consumables have fallen below minimum stock levels:** Initiate emergency procurement per the most critical items first. Update consumption rate estimates based on actual usage. Adjust minimum stock levels if they were set too low. Review whether the procurement lead time assumptions are still valid.

5. **If unregistered assets are discovered during an audit:** Register them immediately. Attempt to determine their origin (purchase records, operator memory, labeling). If origin cannot be determined, record "Unknown -- discovered during audit on [date]" in the Acquisition Source field. For hardware assets of unknown origin, consider whether they should be quarantined per SEC-001 trust model (hardware that has been outside the operator's physical control is not trusted).

## 8. Evolution Path

- **Years 0-5:** The asset registry is being built for the first time. Every asset currently in the institution must be registered. This is a substantial initial effort. Expect the first comprehensive audit to take a full day or more. The registry format may be revised as practical experience reveals which fields are useful and which are overhead.

- **Years 5-15:** The registry becomes genuinely useful for replacement planning. The lifecycle data accumulated over years reveals patterns: which manufacturers' drives last longer, which batteries degrade faster, which cables fail most often. These patterns inform smarter procurement decisions.

- **Years 15-30:** Multiple hardware generations have been tracked. The registry's historical entries tell the story of the institution's material evolution. This history is valuable context for future procurement and for understanding why the institution is configured as it is.

- **Years 30-50+:** The registry is a comprehensive material history of the institution. A successor inheriting the institution can trace every piece of hardware back to its origin, understand why it was chosen, and see the full history of what it replaced. This continuity of material knowledge is as important as the continuity of governance knowledge in the decision log.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
I am going to be honest: I do not enjoy inventory management. Counting spare drives and labeling USB cables is not why I built this institution. But D11-001 is right -- the absence of administration is not freedom, it is amnesia. And I have already experienced the cost of amnesia: the hour spent searching for a cable I knew I had somewhere, the surprise when a drive failed six months before I expected it to, the discovery that I had been running on a battery backup that was past its replacement date.

The Lightness Principle is my defense against making this a full-time job. I will track what matters for decisions. I will not track the serial number of every Ethernet cable. But I will track every drive, every battery, every backup medium, and every piece of infrastructure that my data depends on. That is the line, and I think it is the right one.

The replacement planning thresholds are aggressive -- beginning procurement research twelve months before projected end-of-life. In an internet-connected world, this would be excessive. In an air-gapped institution where procurement requires physical trips and advance planning, it may not be aggressive enough. I will revisit these thresholds after the first hardware replacement cycle.

## 10. References

- ETH-001 -- Ethical Foundations of the Institution (Principle 2: Integrity Over Convenience; Principle 6: Honest Accounting)
- CON-001 -- The Founding Mandate (air-gap and off-grid constraints on procurement, lifetime operation)
- GOV-001 -- Authority Model (Tier 3 decisions for cross-domain resource allocation)
- OPS-001 -- Operations Philosophy (operational tempo for audits, quarterly and annual cycles)
- D11-001 -- Administration Philosophy (administrative principles, Lightness Principle, R-ADM-01)
- SEC-002 -- Access Control Procedures (physical zone model for asset location tracking)
- D6-001 -- Data Philosophy (Tier 1 classification for the asset registry)
- D6-002 -- Backup Doctrine (backup media tracking, secure erasure procedures)
- D12-001 -- Disaster Recovery Philosophy (rebuild mindset, spare parts as resilience investment)
- D12-003 -- System Rebuild Procedures (asset registry as input for hardware replacement)
- D13-002 -- Technology Migration Playbook (asset lifecycle data informs migration timing)

---

---

# D13-002 -- Technology Migration Playbook

**Document ID:** D13-002
**Domain:** 13 -- Evolution & Adaptation
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, OPS-001, D13-001, D11-001, D11-002, D12-001, D12-002, D6-002, SEC-003
**Depended Upon By:** All articles involving planned technology transitions, hardware swaps, OS upgrades, data format conversions, or software migrations.

---

## 1. Purpose

This article is the step-by-step operational manual for migrating between technologies within the holm.chat Documentation Institution. It covers the three most common categories of migration -- operating system upgrades, hardware swaps, and data format conversions -- and provides a universal migration framework applicable to any planned technology change.

D13-001 established the evolution philosophy: the conservatism imperative, the four-stage Evolution Decision Framework, the legacy coexistence principle, and the design-for-the-unpredictable mandate. This article operationalizes that philosophy into executable procedures. Where D13-001 says "every migration must have a documented rollback plan," this article defines what that rollback plan looks like. Where D13-001 says "verification is not optional," this article defines the specific verification steps.

Technology migrations are the planned version of what D12-003 addresses in crisis: replacing one technology with another. The difference is time. In disaster recovery, you rebuild under pressure with whatever is available. In technology migration, you plan, test, execute, and verify at a sustainable pace. This difference in tempo does not make migrations less risky -- it makes them differently risky. The danger in migration is not urgency but complacency. Because there is time, the operator may skip pre-migration testing. Because the old system still works, the operator may defer rollback planning. Because the new technology seems straightforward, the operator may underestimate the blast radius of the change.

This article is the antidote to that complacency. It imposes structure on every migration, regardless of how simple it appears.

## 2. Scope

**In scope:**
- The universal migration framework: applicable to all technology changes.
- Operating system upgrades: minor versions, major versions, and distribution changes.
- Hardware swaps: replacing individual components and full platform migrations.
- Data format conversions: migrating data from one format to another.
- Pre-migration testing procedures.
- Rollback planning and execution.
- Post-migration verification and stabilization.
- Migration documentation requirements.

**Out of scope:**
- The philosophy of evolution and when to migrate (D13-001).
- Emergency hardware replacement during disaster recovery (D12-003).
- Asset procurement (subsequent Domain 11 articles; this article assumes the replacement technology is already available).
- Specific encryption algorithm migrations (SEC-003 handles key rotation; this article handles the platform aspects).

## 3. Background

### 3.1 Migration as the Institution's Primary Survival Mechanism

Over fifty years, the institution will undergo dozens of technology migrations. Drives will be replaced multiple times. The operating system will be upgraded repeatedly and may change distributions or even paradigms. Data formats will be converted as standards evolve and old formats become unsupported. Each of these migrations carries risk: data loss, configuration errors, compatibility failures, and the subtle corruption that comes from assumptions that held in the old system but do not hold in the new one.

The institution's ability to execute migrations reliably is, in the long run, its most important survival skill. A single botched migration can cause more damage than most disasters, because migration errors often affect the entire data set rather than a single component.

### 3.2 The Air-Gap Migration Constraint

In a networked environment, migrations can leverage online repositories, vendor support, automated migration tools, and cloud-based test environments. This institution has none of these. Every migration is executed with locally available resources: the current system, the replacement technology (physically transported across the air gap), the institution's documentation, and the operator's skills.

This constraint makes pre-migration testing especially critical. There is no support hotline to call when a migration goes wrong. The operator must be able to diagnose and resolve migration problems using only the resources within the institution. Testing must be thorough enough to surface problems before the production migration, not after.

### 3.3 The Three Migration Categories

While this article provides a universal framework applicable to any migration, three categories account for the vast majority of migrations the institution will face:

**Operating system upgrades** occur every 2-5 years. They range from minor version updates (security patches, bug fixes) to major version upgrades (new kernel, new package versions, new default configurations) to distribution changes (switching from one Linux distribution to another).

**Hardware swaps** occur every 5-10 years for major components (servers, storage arrays) and more frequently for peripherals and consumables. They range from drop-in replacements (same model, same interface) to platform migrations (different architecture, different interfaces, different capabilities).

**Data format conversions** occur whenever a format becomes unsupported, whenever a better format becomes available, or whenever a software change requires a different format. They are the most dangerous category because they touch the institution's data directly.

## 4. System Model

### 4.1 The Universal Migration Framework

Every technology migration, regardless of category or scope, follows this seven-phase framework:

**Phase 1: Justification and Approval.**
1. Document the justification for the migration per D13-001 Section 4.1, Stage 1.
2. Classify the migration by governance tier per GOV-001:
   - OS major version upgrade or distribution change: Tier 2.
   - OS minor version update: Tier 3.
   - Hardware swap (same interface): Tier 3.
   - Hardware platform migration: Tier 2.
   - Data format conversion: Tier 2.
3. Submit the migration for governance approval with the required waiting period.
4. Document the approval in the decision log.

**Phase 2: Impact Assessment.**
1. List every system, service, configuration, and data set affected by the migration.
2. List every dependency that the migrating component has on other components, and every dependency other components have on it.
3. Assess the blast radius: if the migration fails completely, what is affected?
4. Assess data risk: does the migration touch any Tier 1 data? If yes, this is a high-risk migration requiring enhanced safeguards.
5. Document the assessment in the migration plan.

**Phase 3: Pre-Migration Testing.**
1. Create a test environment. This may be:
   - A spare machine or partition designated for testing.
   - A virtual environment if the institution operates any virtualization.
   - A non-production subsystem that can be temporarily repurposed.
2. Execute the migration in the test environment.
3. Record the results: what worked, what failed, what required adjustment.
4. Test all dependent systems in the test environment.
5. If the test reveals problems, resolve them before proceeding. If the problems cannot be resolved, abort the migration and document why.
6. The pre-migration test must demonstrate that the migration can succeed before the production migration is attempted.

**Phase 4: Rollback Planning.**
1. Define the rollback trigger: what conditions will cause the migration to be aborted and rolled back? Be specific. "Something goes wrong" is not a trigger. "Data verification fails on more than 1% of files" is a trigger.
2. Define the rollback procedure: step-by-step instructions for returning to the pre-migration state.
3. Verify that rollback is possible. If the migration is irreversible (for example, a data format conversion that destroys the original), this must be explicitly acknowledged and approved as a risk per D13-001 R-EVL-02.
4. Verify that a complete, verified backup exists per D6-002 before the migration begins. This backup is the ultimate rollback mechanism.
5. Document the rollback plan in the migration plan.

**Phase 5: Production Migration.**
1. Verify the pre-migration backup is complete and verified.
2. Record the migration start in the operational log: date, time, migration ID, scope.
3. Execute the migration following the procedure validated in Phase 3.
4. At each major step, verify the result before proceeding to the next step.
5. If a rollback trigger is activated, stop the migration and execute the rollback procedure.
6. Record the migration completion (or rollback) in the operational log.

**Phase 6: Post-Migration Verification.**
1. Verify all migrated data: checksums, file counts, content spot-checks.
2. Verify all dependent services: start each service, test its core function, verify its output.
3. Verify system performance: is the system performing comparably to pre-migration baseline? Significant performance degradation is a finding requiring investigation.
4. Verify security posture: is the air gap intact? Are encryption keys functional? Are access controls correct?
5. Run a reduced version of the daily operations checklist to confirm the institution is functional.
6. Document the verification results.

**Phase 7: Stabilization and Old System Decommissioning.**
1. Enter the stabilization period. During this period (minimum 30 days for Tier 2 migrations, 7 days for Tier 3 migrations), the old system is kept intact but inactive. It serves as a safety net.
2. Monitor the new system for issues that only appear under sustained operation.
3. At the end of the stabilization period, assess: is the new system performing satisfactorily? Have any issues emerged?
4. If the new system is satisfactory, decommission the old system per D13-001 legacy coexistence principle. For hardware: update the asset registry (D11-002) to mark the old hardware as retired. For data: verify the old copy is no longer needed, then securely erase per D6-002.
5. If the new system is not satisfactory, extend the stabilization period or execute the rollback.
6. Record the final outcome in the decision log.

### 4.2 Operating System Upgrades

**Minor Version Update (security patches, bug fixes):**

Pre-migration:
1. Obtain the update packages. In an air-gapped institution, this means physically transporting update media from an external source. Follow Domain 18 quarantine procedures.
2. Verify package integrity (GPG signatures, checksums) before applying.
3. Review the update changelog for any changes that could affect institutional services.
4. Create a full system backup per D6-002.

Migration:
1. Apply updates on a test system first, if available.
2. If no test system is available, apply updates to the production system during a maintenance window.
3. Using the administrative account:
   ```
   sudo apt update  (if using local repository mirror)
   sudo apt upgrade
   ```
   Or the equivalent for the installed package manager.
4. Review any prompts carefully. Do not accept default configuration changes without understanding them.
5. Reboot if required by the update (kernel updates require reboot).
6. Verify all services start correctly after reboot.

Post-migration:
- Verify system functionality per Phase 6.
- Stabilization period: 7 days.
- Record the update in the operational log.

**Major Version Upgrade (new major OS version):**

Pre-migration:
1. Obtain the upgrade media and all required packages.
2. Read the upgrade release notes thoroughly. Identify all breaking changes.
3. Assess impact on all institutional services: do any services require version-specific configurations?
4. Create a full system backup and verify it.
5. If possible, perform the upgrade on a test system first and exercise all critical services.

Migration:
1. Schedule a full maintenance day. Major upgrades can take several hours.
2. Follow the distribution's documented upgrade procedure.
3. At each major step, verify before proceeding.
4. After the upgrade completes, verify each service individually.
5. Verify the air gap: confirm no network interfaces were enabled by the upgrade.

Post-migration:
- Full Phase 6 verification.
- Stabilization period: 30 days.
- During stabilization, run all quarterly operations checks to exercise the full operational surface.
- Record as a Tier 2 decision in the decision log.

**Distribution Change (switching to a different OS distribution):**

This is the most complex OS migration. Treat it as a rebuild per D12-003, with the addition of service migration from the old distribution to the new one.

Pre-migration:
1. All steps of the major version upgrade pre-migration, plus:
2. Identify all distribution-specific configurations (init system, package naming, default paths, default tools).
3. Create a mapping document: for every institutional service, document its configuration on the old distribution and its equivalent on the new distribution.
4. Perform a full test migration on separate hardware.

Migration:
1. Install the new distribution on the target hardware per D12-003 Phase 2-3.
2. Migrate services one at a time, verifying each before proceeding to the next.
3. Restore data from backup to the new system.
4. Verify all data and services.

Post-migration:
- Full Phase 6 verification.
- Stabilization period: 90 days (the old system remains available for the full period).
- Record as a Tier 2 decision.

### 4.3 Hardware Swaps

**Drop-in Replacement (same model or compatible equivalent):**

Pre-migration:
1. Verify the replacement hardware is compatible (same interface, same capacity or greater).
2. Register the new hardware in the asset registry per D11-002.
3. Test the new hardware independently: for drives, run SMART self-test and write/read verification; for other components, run appropriate diagnostics.
4. Create a full backup of any data on the hardware being replaced.

Migration:
1. Schedule the swap during a maintenance window.
2. Shut down the system (or the affected subsystem) cleanly.
3. Physically swap the component.
4. If the swap involves storage: restore data from backup to the new device, or use the cloning procedure:
   ```
   dd if=/dev/[old-device] of=/dev/[new-device] bs=4M status=progress
   ```
   Then verify the clone:
   ```
   sha256sum /dev/[old-device] /dev/[new-device]
   ```
   Hashes must match.
5. Boot the system and verify functionality.

Post-migration:
- Verify data integrity on the new device.
- Stabilization period: 7 days (retain the old device as a backup).
- Update the asset registry: retire the old device, activate the new device.

**Platform Migration (different hardware architecture):**

Treat this as a rebuild per D12-003 on the new hardware, with additional service migration validation. All data is restored from backup rather than cloned, because the old system's disk layout may not be compatible with the new hardware.

Pre-migration:
1. All steps of a distribution change pre-migration (the new hardware may require a different kernel, different drivers, or a different distribution).
2. Perform a complete test rebuild on the new hardware before decommissioning the old hardware.

Migration:
1. Follow D12-003 Phases 2-6 on the new hardware.
2. Verify all services on the new platform.
3. Run at least one DR drill scenario (D12-002) on the new platform to verify resilience.

Post-migration:
- Stabilization period: 90 days.
- The old platform remains available as a fallback for the full stabilization period.
- Record as a Tier 2 decision.

### 4.4 Data Format Conversions

Data format conversions are the highest-risk category of migration because they touch the institution's data directly. A failed format conversion can corrupt data permanently.

**The Format Conversion Procedure:**

Pre-migration:
1. Identify all files in the source format. Generate a complete manifest:
   ```
   find /data -name "*.[old-extension]" > /tmp/conversion-manifest.txt
   wc -l /tmp/conversion-manifest.txt
   ```
2. Create a full backup of all files to be converted. Verify the backup.
3. Select the target format per D13-001 R-EVL-04 (must be an open format readable without proprietary software).
4. Identify or obtain the conversion tool. Verify the tool produces correct output by testing on sample files.
5. Define the verification criteria: how will you confirm that the converted file is equivalent to the original? This may involve content comparison, visual inspection, structural validation, or round-trip conversion (convert A to B, then B back to A, and verify A matches the original).

Migration:
1. Convert a test batch (10% of files or 100 files, whichever is smaller) first.
2. Verify the test batch against the verification criteria.
3. If the test batch passes verification, proceed with the full conversion.
4. For each file:
   a. Convert from old format to new format.
   b. Verify the converted file.
   c. Record the conversion result (success/failure) in a conversion log.
5. After all files are converted, generate a summary: total files, successful conversions, failed conversions, files requiring manual attention.

Post-migration:
1. Do not delete the original files immediately. Retain them for the full stabilization period.
2. Stabilization period: 90 days for Tier 1 data, 30 days for Tier 2 data, 7 days for Tier 3 data.
3. During stabilization, use the converted files in normal operations to verify they function correctly in practice.
4. At the end of stabilization, if all converted files are verified, archive the originals (do not delete them -- archive them as historical artifacts per D13-001 R-EVL-07).
5. Record the conversion as a Tier 2 decision.

### 4.5 The Migration Plan Template

Every migration produces a migration plan document. This document is a Tier 2 data artifact.

```
TECHNOLOGY MIGRATION PLAN

Migration ID:       MIG-[YYYY]-[sequential number]
Date Proposed:      ____-__-__
Proposed By:        [operator identifier]
Governance Tier:    [Tier 2 / Tier 3]
Approval Date:      ____-__-__

SUMMARY:
[One-paragraph description of what is being migrated and why]

JUSTIFICATION:
[Why this migration is necessary. What happens if it is not done.
Reference to D13-001 Evolution Decision Framework Stage 1.]

CATEGORY:
[OS Upgrade / Hardware Swap / Format Conversion / Other]

SCOPE:
[What systems, data, and services are affected]

IMPACT ASSESSMENT:
[Blast radius. Dependencies. Data risk level.]

PRE-MIGRATION TEST RESULTS:
[Date of test. Outcome. Issues discovered. Resolutions.]

ROLLBACK PLAN:
  Rollback Triggers:
  [Specific conditions that trigger rollback]

  Rollback Procedure:
  [Step-by-step rollback instructions]

  Rollback Feasibility:
  [Feasible / Partially feasible / Irreversible]
  [If irreversible, governance approval reference]

PRE-MIGRATION CHECKLIST:
[ ] Backup verified
[ ] Test migration completed successfully
[ ] Rollback plan documented and reviewed
[ ] Governance approval obtained
[ ] Maintenance window scheduled
[ ] All stakeholders notified (successor, if applicable)

MIGRATION PROCEDURE:
[Step-by-step production migration procedure]

POST-MIGRATION VERIFICATION:
[Specific verification steps and success criteria]

STABILIZATION PERIOD:
[Duration. Monitoring activities. Decommissioning criteria.]

OUTCOME:
[Completed after migration]
  Migration Date:     ____-__-__
  Result:             [Success / Partial / Rolled back]
  Issues Encountered:
  Lessons Learned:
  Decommissioning Date: ____-__-__

Signed: _______________
```

## 5. Rules & Constraints

- **R-D13-02-01:** Every technology migration must follow the universal migration framework (Section 4.1). No phase may be skipped without documented justification approved at the appropriate governance tier.
- **R-D13-02-02:** Every migration must include pre-migration testing (Phase 3) unless no test environment exists and the migration is classified as low-risk (Tier 3, no Tier 1 data affected). If testing is skipped, the skip must be documented with justification.
- **R-D13-02-03:** Every migration must have a documented rollback plan (Phase 4). If rollback is not feasible, this must be acknowledged as a risk and approved per D13-001 R-EVL-02.
- **R-D13-02-04:** A complete, verified backup must exist before any production migration begins. No exceptions.
- **R-D13-02-05:** The old system must be retained intact for the full stabilization period. Premature decommissioning of the old system before stabilization is complete requires Tier 2 governance approval regardless of the migration's original tier.
- **R-D13-02-06:** Data format conversions affecting Tier 1 data must retain the original files permanently (archived, not deleted) per D13-001 R-EVL-07.
- **R-D13-02-07:** Every migration must produce a migration plan document per Section 4.5. The completed plan (including outcome) is a Tier 2 data artifact.
- **R-D13-02-08:** Post-migration verification (Phase 6) must include data integrity verification (checksums), service functionality verification, and security posture verification. All three must pass before the migration is considered successful.

## 6. Failure Modes

- **Testing omission.** The operator skips pre-migration testing because "this is a simple upgrade" or "I have done this before." The production migration fails in a way that testing would have caught. Mitigation: R-D13-02-02 makes testing mandatory for all migrations involving Tier 1 data. The migration plan template includes a pre-migration checklist that explicitly requires test completion.

- **Rollback impossibility.** The migration proceeds without a viable rollback plan. When problems emerge, the operator cannot return to the pre-migration state. Mitigation: R-D13-02-03 requires a documented rollback plan. R-D13-02-04 ensures a backup exists as the ultimate rollback mechanism.

- **Premature decommissioning.** The old system is removed before the stabilization period ends. A problem emerges with the new system that the old system could have resolved, but it is no longer available. Mitigation: R-D13-02-05 requires retention of the old system for the full stabilization period.

- **Data corruption during conversion.** A format conversion tool produces silently incorrect output. The converted files appear valid but contain subtle errors. The original files are deleted. The errors are discovered months later. Mitigation: the verification criteria in Section 4.4 must be designed to detect subtle errors, not just gross failures. R-D13-02-06 retains original files permanently for Tier 1 data. The 90-day stabilization period for Tier 1 data provides extended time to discover problems in practice.

- **Migration fatigue.** The institution faces multiple migrations simultaneously (OS upgrade, drive replacement, format conversion). The operator attempts all of them at once. Interactions between migrations create unexpected failures. Mitigation: stagger migrations. Only one migration should be in Phase 5 (production) at any time. Other migrations wait in Phase 4 (rollback planning complete, ready to execute) until the current migration reaches Phase 7 (stabilization).

- **Documentation lag.** The migration is executed successfully but the documentation is never updated to reflect the new system state. Procedures, configurations, and references still describe the pre-migration system. Mitigation: the migration plan template includes documentation update as part of Phase 7. The post-migration verification checklist should include a documentation currency check.

## 7. Recovery Procedures

1. **If a migration has failed and rollback is triggered:** Execute the rollback plan documented in Phase 4. If the rollback plan fails, restore from the pre-migration backup per D6-002. Document the failure and the rollback in the migration plan and the operational log. Conduct a post-mortem: why did the migration fail? Why did the rollback plan fail (if applicable)? What needs to change before attempting the migration again?

2. **If a migration has partially succeeded:** Assess the scope: what parts of the migration succeeded and what parts failed? If the succeeded parts are stable and independent, consider completing the migration incrementally rather than rolling back entirely. If the failed parts are entangled with the succeeded parts, rollback entirely and retry.

3. **If post-migration verification reveals problems:** During stabilization, determine whether the problems are fixable in the new system or whether rollback is necessary. Minor issues (configuration adjustments, permission corrections) can be fixed in place. Major issues (data integrity failures, service incompatibilities) should trigger rollback to the old system. Document all findings.

4. **If a data format conversion has produced corrupt output:** Stop the conversion immediately. Assess how many files are affected. Restore the original files from the pre-conversion backup. Investigate the conversion tool: was it defective, misconfigured, or used on incompatible files? Fix the cause. Re-test on a small batch. Only resume the conversion when the test batch passes verification.

5. **If multiple migrations need to happen simultaneously:** Do not attempt simultaneous production migrations. Prioritize by urgency and risk: migrations that address active threats (failing hardware, deprecated security) go first; migrations that add capabilities go later. Schedule migrations at least one stabilization period apart.

## 8. Evolution Path

- **Years 0-5:** The first migrations are likely small: OS security patches, minor hardware swaps. Use these as practice for the framework. The migration plan template may be refined based on practical experience. The pre-migration testing procedures will become more efficient as the operator gains familiarity.

- **Years 5-15:** The first major migrations occur: OS major version upgrades, storage array replacements, possibly the first data format conversions. The framework is tested against real complexity. Expect the rollback and stabilization procedures to be revised based on experience.

- **Years 15-30:** Multiple hardware generations have been navigated. The migration history (accumulated migration plans) becomes a valuable record. Patterns emerge: which types of migrations are most problematic, what testing catches the most issues, what stabilization period is actually needed. Use these patterns to refine the framework.

- **Years 30-50+:** Technology migrations are routine. The framework has been refined through decades of practice. The migration archive tells the story of the institution's technological evolution. New operators can study past migrations to understand why the institution is configured as it is and what challenges previous transitions presented.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The seven-phase framework may seem heavy for a single-person institution. It was designed to be heavy. Every migration I have ever witnessed go wrong -- in professional environments with teams of engineers -- went wrong because a phase was skipped. Testing was skipped because "it is a simple change." Rollback planning was skipped because "what could go wrong?" Stabilization was cut short because "we need the old hardware for another project."

In a single-operator institution, these failures are amplified. There is no team to catch the error the operator missed. There is no second pair of eyes on the verification. The framework compensates for the absence of a team by imposing structure that a team would provide naturally through division of labor and peer review.

The one thing I am most certain will need revision is the stabilization periods. Thirty days for a Tier 2 migration and seven days for a Tier 3 migration are my best guesses. Experience will reveal whether these are too long (creating unnecessary overhead from maintaining two systems) or too short (allowing problems to slip through). I will revise after the first few migrations.

## 10. References

- ETH-001 -- Ethical Foundations of the Institution (Principle 4: Longevity Over Novelty)
- CON-001 -- The Founding Mandate (air-gap constraint on migration logistics)
- GOV-001 -- Authority Model (governance tiers for migration approval, decision logging)
- OPS-001 -- Operations Philosophy (documentation-first principle, maintenance windows, operational tempo)
- D13-001 -- Evolution Philosophy (Evolution Decision Framework, conservatism imperative, legacy coexistence, rollback mandate)
- D11-001 -- Administration Philosophy (resource planning for migrations)
- D11-002 -- Asset Inventory and Tracking (asset registry for hardware tracking during swaps)
- D12-001 -- Disaster Recovery Philosophy (rebuild mindset as migration foundation)
- D12-002 -- Disaster Recovery Drill Framework (DR drills on new platform after migration)
- D12-003 -- System Rebuild Procedures (rebuild as platform migration method)
- D6-002 -- Backup Doctrine (pre-migration backup requirement, backup as rollback mechanism)
- SEC-003 -- Cryptographic Key Management (key handling during platform migration)
- D5-002 -- Operating System Maintenance Procedures (OS-level upgrade procedures)

---

---

# GOV-002 -- Decision Record Keeping

**Document ID:** GOV-002
**Domain:** 2 -- Governance & Authority
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, OPS-001
**Depended Upon By:** GOV-003, GOV-004, GOV-005, GOV-006, GOV-007. All articles that produce or consume decision records. Referenced by all domains that record governance decisions.

---

## 1. Purpose

This article defines how decisions are recorded, stored, retrieved, and used for institutional learning within the holm.chat Documentation Institution. It operationalizes the decision record concept introduced in GOV-001 Section 4.2 into a complete, practical system.

GOV-001 established that every decision at Tier 1, 2, or 3 must be recorded in the institutional decision log with specific fields. This article defines the detailed format of those records, the physical and logical structure of the decision log, the procedures for creating and indexing entries, the classification system that makes records retrievable, and the methods by which the decision log serves as a tool for institutional learning.

The decision log is not bureaucratic overhead. It is the institution's governance memory. When the operator wonders, five years from now, why a specific system is configured a particular way, the decision log provides the answer. When a successor inherits the institution and encounters a policy they do not understand, the decision log explains the reasoning. When a dispute arises about whether a past decision was correct, the decision log provides the evidence. And when the institution needs to learn from its own history -- to identify patterns in its decision-making, to detect systematic biases, to improve the quality of future decisions -- the decision log provides the data.

Without the decision log, the institution suffers from what CON-001 calls "purpose amnesia": the gradual loss of understanding of why things are the way they are. With the decision log, every configuration, every policy, every architectural choice is traceable to a specific decision, made on a specific date, for specific reasons. This traceability is the foundation of institutional coherence over decades.

## 2. Scope

**In scope:**
- The decision record format: all fields, their definitions, and their usage rules.
- The decision log structure: how records are organized and stored.
- Classification of decisions by tier, domain, and topic.
- The decision record creation procedure.
- Indexing and search procedures for retrieving past decisions.
- The decision log as a tool for institutional learning.
- Integration with Commentary Sections, the operational log, and other institutional records.
- Maintenance and preservation of the decision log.

**Out of scope:**
- The governance tier system itself (GOV-001 Section 4.1).
- Specific decision-making procedures for specific domains (addressed in each domain's operational articles).
- Dispute resolution involving past decisions (GOV-003).
- The amendment process for root documents (GOV-001).

## 3. Background

### 3.1 The Decision Record as Institutional DNA

Every institution is the sum of its decisions. The hardware, the software, the policies, the procedures -- all of these are the current expression of decisions made at specific points in time by specific people for specific reasons. In most personal technology projects, these decisions exist only in the operator's memory, which means they exist approximately and temporarily. Memory fades, distorts, and eventually disappears.

The decision record is the mechanism by which this institution preserves its decisions in durable, exact form. Each record captures not just what was decided but why it was decided, what alternatives were considered, and what trade-offs were accepted. This "why" is the most valuable part of the record. A system configuration tells you what the institution does. A decision record tells you why it does it. Without the why, future operators can only guess at intent, and guesses compound into misunderstanding over decades.

### 3.2 The Append-Only Principle

GOV-001 R-GOV-03 establishes that the decision log is append-only. Entries may not be modified or deleted. Corrections must be recorded as new entries that reference the original. This principle is the decision log's most important property.

The append-only principle serves three functions. First, it creates an indelible audit trail. Any future reader -- including a successor, an auditor, or the operator's own future self -- can see the complete history of institutional governance, including mistakes, corrections, and changes of direction. Second, it prevents historical revisionism. The temptation to clean up past decisions, to make them look more logical than they were, to erase the reasoning that turned out to be wrong -- this temptation is real, and the append-only principle removes it. Third, it supports institutional learning. You cannot learn from decisions you have erased. The record of decisions that were later reversed is often more instructive than the record of decisions that stood.

### 3.3 The Relationship to Other Records

The institution maintains several types of records, each serving a different purpose:

- **The decision log** (this article): records governance decisions (Tier 1-3).
- **The operational log** (OPS-001): records daily operations, maintenance actions, and Tier 4 decisions.
- **The Commentary Sections** (embedded in each article): records interpretive context, lessons learned, and evolving understanding.
- **The asset registry** (D11-002): records the institution's material state.
- **The backup log** (D6-002): records backup operations and verifications.
- **The key inventory** (SEC-003): records cryptographic key metadata.

These records are complementary, not competing. A decision recorded in the decision log may reference an operational log entry that prompted the decision. A Commentary Section entry may reference a decision log entry that it interprets. The decision log is the governance spine that connects the other records into a coherent institutional narrative.

## 4. System Model

### 4.1 The Decision Record Format

Each decision record contains the following fields, expanding on the format defined in GOV-001 Section 4.2:

```
DECISION RECORD

Decision ID:        DEC-[YYYY]-[sequential number]
Date:               [Date the decision was made -- not the proposal date, the decision date]
Author:             [Identifier of the person who made the decision]
Tier:               [1 / 2 / 3]
Status:             [Proposed / Waiting / Ratified / Superseded / Reversed]

SUMMARY:
[One sentence describing the decision. This is the record's headline.
It should be clear enough that someone scanning the log can understand
the decision without reading the full record.]

DOMAIN(S):
[Which domain(s) this decision affects. Use domain numbers and names.]

TOPIC TAGS:
[2-5 keywords for indexing. Examples: backup, encryption, hardware,
migration, governance, succession, security, format, procedure.]

CONTEXT:
[What prompted this decision? What situation, problem, observation,
or scheduled review led to this decision being considered?
Reference operational log entries, drill findings, audit results,
or other records that provide context.]

ALTERNATIVES CONSIDERED:
[What other options were evaluated? For each alternative:]
  Alternative 1: [Description]
    Pros: [...]
    Cons: [...]
    Why rejected: [...]

  Alternative 2: [Description]
    Pros: [...]
    Cons: [...]
    Why rejected: [...]

  [Additional alternatives as needed]

RATIONALE:
[Why this decision was made. What principles, evidence, or reasoning
led to choosing this option over the alternatives. Reference root
documents (ETH-001, CON-001, etc.) where applicable.
This is the most important section of the record.]

IMPACT:
[What systems, processes, documents, or data are affected by this
decision. Be specific. This section helps future readers understand
the scope of the decision's effects.]

IMPLEMENTATION:
[Brief description of how the decision will be implemented.
Reference the specific articles or procedures that will be
updated or created as a result of this decision.]

REVERSIBILITY:
  Reversible:       [Yes / Partially / No]
  Reversal Cost:    [Low / Medium / High / Catastrophic]
  Reversal Procedure: [Brief description of how this decision
                       could be undone, if applicable]

WAITING PERIOD:
  Start Date:       [Date the waiting period began]
  End Date:         [Date the waiting period ends]
  Review Date(s):   [Date(s) when the decision was reviewed during
                     the waiting period, per R-GOV-02]

REVIEW SCHEDULE:
  Next Review:      [Date when this decision should be reviewed
                     for continued appropriateness]
  Review Frequency: [How often this decision should be reviewed --
                     annually, biennially, or on a specific trigger]

RELATED DECISIONS:
[References to other decision records that are related to this one.
Include: decisions this one supersedes, decisions this one depends on,
decisions that depend on this one.]

AMENDMENTS:
[This section is populated only after the decision is made.
If the decision is later modified, superseded, or reversed,
record the reference to the new decision here.]
  [DEC-YYYY-NNN: brief description of amendment]
```

### 4.2 The Decision Log Structure

The decision log is a Tier 1 data artifact per D6-001. It is stored as a collection of plain-text files (per D11-001 R-ADM-05 and D13-001 R-EVL-04), one file per decision record, organized in a directory structure:

```
/data/tier1/decision-log/
  index.txt                    -- Master index of all decisions
  by-year/
    2026/
      DEC-2026-001.txt
      DEC-2026-002.txt
      ...
    2027/
      DEC-2027-001.txt
      ...
  by-tier/
    tier1/                     -- Symbolic links to Tier 1 decisions
    tier2/                     -- Symbolic links to Tier 2 decisions
    tier3/                     -- Symbolic links to Tier 3 decisions
  by-domain/
    domain-01/                 -- Symbolic links to decisions in each domain
    domain-02/
    ...
    domain-20/
    cross-domain/              -- Decisions affecting multiple domains
```

The canonical copy of each decision is in the `by-year/` directory. The `by-tier/` and `by-domain/` directories contain symbolic links that provide alternative access paths for searching.

### 4.3 The Master Index

The master index (`index.txt`) is a summary listing of every decision in the log. It provides a scannable overview of the institution's governance history.

**Index format:**

```
DECISION LOG -- MASTER INDEX

Last Updated: [date]
Total Decisions: [count]

[Decision ID] | [Date] | [Tier] | [Status] | [Summary]
----------------------------------------------------------------------
DEC-2026-001  | 2026-02-16 | T1 | Ratified  | Founding of the institution...
DEC-2026-002  | 2026-02-16 | T2 | Ratified  | Adoption of LUKS full-disk...
DEC-2026-003  | 2026-03-01 | T3 | Ratified  | Set daily backup schedule...
...
```

The index is updated whenever a new decision is recorded or an existing decision's status changes. It is a convenience document, not the authoritative record -- the individual decision files are authoritative. If the index and a decision file disagree, the decision file prevails.

### 4.4 The Decision Record Creation Procedure

**Step 1: Identify the decision.**
Determine that a governance decision is being made (as opposed to a routine operational action). If the action involves changing policy, architecture, security posture, or any item that a future operator would benefit from understanding, it is a governance decision requiring a record.

**Step 2: Classify the tier.**
Using GOV-001 Section 4.1:
- Does it amend a root document? Tier 1.
- Does it change the security, data, or infrastructure architecture? Tier 2.
- Does it change operational procedures or policies? Tier 3.
- Is it routine? Tier 4 (operational log only, no decision record required).

**Step 3: Draft the record.**
Create a new file using the format in Section 4.1. Fill in all fields. The most important fields to draft carefully are:
- **Summary:** Must be clear to someone who reads only this line.
- **Alternatives Considered:** Must be honest. If no alternatives were considered, say so -- and explain why.
- **Rationale:** Must explain the reasoning, not just state the conclusion. "We chose option A" is not rationale. "We chose option A because it satisfies ETH-001 Principle 4 while option B introduces a dependency that violates CON-001 Section 5.3" is rationale.

**Step 4: Initiate the waiting period.**
Set the Status to "Waiting." Record the waiting period start date. Set a calendar reminder for the review date(s) and the end date.

**Step 5: Review during the waiting period.**
Per GOV-001 R-GOV-02, review the decision at least once during the waiting period. During review:
- Re-read the rationale. Does it still hold?
- Re-read the alternatives. Has new information changed the evaluation?
- Re-read the impact assessment. Is it complete?
- Record the review date in the Waiting Period section.

**Step 6: Ratify or withdraw.**
At the end of the waiting period:
- If the decision still stands, update Status to "Ratified."
- If the decision has been reconsidered, update Status to "Withdrawn" and record the reason in the Amendments section.
- Update the master index.

**Step 7: Implement.**
Proceed with implementation per the Implementation field. Record the implementation in the operational log and reference the Decision ID.

### 4.5 Search and Retrieval Procedures

The decision log must be searchable. Without effective retrieval, the log is a write-only archive -- decisions go in but never come out for consultation. The following search methods are supported:

**Search by ID:** If you know the Decision ID (referenced from another document or log entry), navigate directly to the file:
```
cat /data/tier1/decision-log/by-year/[year]/DEC-[year]-[number].txt
```

**Search by keyword:** Use text search across all decision records:
```
grep -r -l "[keyword]" /data/tier1/decision-log/by-year/
```
For more specific searches:
```
grep -r -l "TOPIC TAGS:.*encryption" /data/tier1/decision-log/by-year/
```

**Search by tier:** Browse the symbolic links in the tier directory:
```
ls /data/tier1/decision-log/by-tier/tier2/
```

**Search by domain:** Browse the symbolic links in the domain directory:
```
ls /data/tier1/decision-log/by-domain/domain-12/
```

**Search by date range:** Browse the year directories and filter by date:
```
grep "Date:" /data/tier1/decision-log/by-year/2028/*.txt | grep "2028-0[1-6]"
```

**Search by status:** Find all active (non-superseded, non-reversed) decisions:
```
grep -r -l "Status:.*Ratified" /data/tier1/decision-log/by-year/
```

**Browse the master index:** For a high-level scan:
```
cat /data/tier1/decision-log/index.txt
```

### 4.6 The Decision Log as a Learning Tool

The decision log's highest-value function is not record-keeping -- it is institutional learning. The log contains the raw material for understanding how the institution makes decisions, where its reasoning is strong, where it is weak, and how it can improve.

**Quarterly Decision Review (during quarterly operations per OPS-001):**
1. Read all decision records created in the past quarter.
2. For each decision, ask: has the rationale held up? Have the predicted impacts materialized? Have any unforeseen consequences appeared?
3. If a decision's rationale has not held up, initiate a review per GOV-001 (potentially leading to a superseding decision).
4. Record the quarterly review findings in the operational log.

**Annual Decision Audit (during annual operations per OPS-001):**
1. Review all decisions created in the past year.
2. Generate statistics: how many decisions per tier? Per domain? What is the average number of alternatives considered? How many decisions have been superseded or reversed?
3. Identify patterns:
   - Are certain domains generating disproportionate decision volume? (May indicate instability in that domain.)
   - Are Tier 3 decisions frequently being reclassified upward? (May indicate the operator is underestimating decision significance.)
   - Are decisions frequently being reversed? (May indicate inadequate analysis or changed circumstances.)
   - Are the same types of decisions recurring? (May indicate a structural issue that should be addressed at a higher tier.)
4. Review a random sample of decisions from previous years: are the rationales still relevant? Do the decisions reflect the institution's principles as currently understood?
5. Record the audit findings in the operational log and, if significant patterns are discovered, in the Commentary Section of this article.

**Succession Context:**
During succession (GOV-004), the decision log is one of the most important documents the successor reads. The decision log tells the successor not just what the institution does but why it does it. For the successor:
- Read all Tier 1 and Tier 2 decisions in full.
- Read the summaries of all Tier 3 decisions.
- Pay particular attention to decisions that were reversed or superseded -- these reveal how the institution's understanding evolved.
- Record your own reactions in your first-impressions journal (GOV-004 Section 4.5).

## 5. Rules & Constraints

- **R-GOV-02-01:** Every Tier 1, Tier 2, and Tier 3 decision must be recorded in the decision log using the format defined in Section 4.1. No exceptions.
- **R-GOV-02-02:** The decision log is append-only per GOV-001 R-GOV-03. Entries may not be modified or deleted. Corrections, amendments, supersessions, and reversals must be recorded as new entries that reference the original.
- **R-GOV-02-03:** The Summary field of every decision record must be comprehensible to a reader who reads only that field. It must convey what was decided, not merely that a decision was made.
- **R-GOV-02-04:** The Rationale field must explain why the decision was made, not merely state what was decided. A record without rationale is incomplete and must be corrected (by adding a supplementary record with the rationale).
- **R-GOV-02-05:** The Alternatives Considered field must list at least one alternative for Tier 1 and Tier 2 decisions. If no alternatives were considered, this must be explicitly stated with an explanation of why.
- **R-GOV-02-06:** The master index must be updated within 24 hours of any new decision being recorded or any decision's status changing.
- **R-GOV-02-07:** The decision log must be maintained as Tier 1 data per D6-001. It is backed up daily per D6-002.
- **R-GOV-02-08:** The quarterly decision review and annual decision audit defined in Section 4.6 are mandatory components of the operational tempo per OPS-001.
- **R-GOV-02-09:** Every decision record must include a Review Schedule with a next review date. No decision may be recorded without a plan for eventual reassessment.

## 6. Failure Modes

- **Decision log neglect.** Decisions are made without being recorded. The log falls silent for weeks or months. When it is consulted, it describes an institution that existed in the past, not the institution that exists now. Mitigation: R-GOV-02-01 makes recording mandatory. The quarterly decision review (Section 4.6) will reveal gaps. The operational tempo in OPS-001 includes time for governance record-keeping.

- **Rationale omission.** Decisions are recorded with what was decided but not why. Five years later, the record says "we switched from format A to format B" but not why. Was it performance? Compatibility? A philosophical preference? Without the rationale, the decision is nearly useless for institutional learning. Mitigation: R-GOV-02-04 requires rationale. The decision creation procedure (Section 4.4) explicitly prompts for rationale and distinguishes it from summary.

- **Index decay.** The master index is not updated. New decisions exist as files in the directory but are not reflected in the index. The index becomes an unreliable search tool. Mitigation: R-GOV-02-06 requires index updates within 24 hours. The decision creation procedure includes index update as a step.

- **Search failure.** The decision log exists but is effectively unsearchable because the topic tags are inconsistent, the summaries are vague, or the directory structure has degraded. Decisions that were made are not found when searched for. Mitigation: the search procedures in Section 4.5 provide multiple access paths. The annual decision audit includes an assessment of search effectiveness.

- **Learning deficit.** The decision log is maintained but never consulted for learning. Decisions are made without checking whether similar decisions were made in the past. The same reasoning errors recur because no one reviews the history. Mitigation: the quarterly and annual reviews in Section 4.6 are mandatory. The self-dispute protocol in GOV-003 includes consulting the decision log for precedent (Level 3 escalation).

- **Append-only erosion.** The operator begins editing past decision records -- correcting typos at first, then clarifying language, then subtly adjusting the rationale to make past decisions look better. The append-only property is lost without any visible violation. Mitigation: the append-only principle is a structural safeguard (D15-002 Section 4.1). Technical enforcement (checksums on historical records, write-protection) should be implemented where possible. Cultural enforcement (the practice of creating amendment records rather than editing) must be maintained through discipline.

## 7. Recovery Procedures

1. **If the decision log has been neglected:** Reconstruct what can be reconstructed. Review the operational log for the neglect period. For every operational log entry that describes a Tier 1, 2, or 3 decision, create a decision record retroactively. Mark these records as "RECONSTRUCTED -- recorded retroactively on [date], original decision made on [estimated date]." The rationale for reconstructed records will necessarily be less detailed than for real-time records. Record what can be remembered, and honestly note where memory is uncertain. Going forward, adhere strictly to the recording requirements.

2. **If the decision log has been corrupted or lost:** Restore from the most recent backup per D6-002. If the most recent backup is also corrupted, use the next-oldest backup. If no backup contains an intact decision log, begin reconstruction from the operational log, Commentary Sections, and the operator's memory. This is a significant institutional loss. Record the loss and the reconstruction effort as a Tier 2 decision.

3. **If the master index is out of date:** Regenerate the index from the actual decision record files. The files are authoritative; the index is derived. Walk through every file in the `by-year/` directories and rebuild the index entries. Update the symbolic links in `by-tier/` and `by-domain/` directories.

4. **If search is ineffective:** Audit the topic tags across all records. Standardize the tagging vocabulary. Update records that have vague or missing tags (this is an exception to the append-only rule: metadata tags may be added for indexing purposes, but the substantive content of the record must not change). Consider creating a topic tag glossary that defines the standard tags for the institution.

5. **If the append-only principle has been violated:** Detect the violation through checksums or version comparison against backups. Restore the original records from backup. Create a new decision record documenting the violation: what was changed, when, and (if known) why. Add an entry to the Commentary Section of this article documenting the violation and the corrective action. Implement technical controls to prevent recurrence (file permissions, write-protection, integrity monitoring).

## 8. Evolution Path

- **Years 0-5:** The decision log is being established. The first records will feel formal and slow. This is intentional. The habits formed now -- thorough rationale, honest alternatives, specific summaries -- will define the quality of the institution's governance record for decades. Expect the record format to be refined as practical experience reveals which fields are useful and which are overhead.

- **Years 5-15:** The decision log contains hundreds of records. The annual audit begins to reveal patterns. The log becomes genuinely useful for institutional learning -- not just record-keeping but actual improvement in decision quality based on reviewing past decisions. The search procedures become critical as the volume grows.

- **Years 15-30:** The decision log is a substantial historical document. Succession planning incorporates the decision log as a primary training resource. The accumulated rationale and alternatives provide rich context that no other document can offer. Consider whether the index and search mechanisms are adequate for the volume of records.

- **Years 30-50+:** The decision log spans generations of operators. It tells the complete governance story of the institution. It is one of the most valuable documents in the corpus -- not because any individual record is transformative, but because the aggregate record reveals how an institution thinks, learns, and evolves. Preserve it with the same care as any other Tier 1 data.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators. Each entry should include the date, the author's identifier, and the context for the commentary.*

**2026-02-16 -- Founding Entry:**
I have kept informal decision records before -- notes in text files, comments in configuration files, half-remembered conversations with myself about why I chose one approach over another. They were better than nothing, but only slightly. The problem was always the rationale. I would record what I decided but not why, because the why seemed obvious at the time. Six months later, the why was gone, and I was left with a configuration I could not explain to myself.

This formal system is designed to prevent that specific failure. The rationale field is the heart of the decision record. Everything else -- the summary, the alternatives, the impact assessment -- serves the rationale. If I can only fill in one field, it should be the rationale.

I am also honest about the append-only principle: it will be tempting to violate. Not for malicious reasons -- I will not want to cover up mistakes -- but for aesthetic ones. I will want to fix a typo in a three-year-old decision record. I will want to clarify a rationale that now seems poorly worded. I will want to update an impact assessment that turned out to be wrong. The answer to all of these temptations is: create a new record. Reference the old one. Let the history stand as it was. The imperfection is the point. It is the record of a real institution making real decisions with incomplete information, which is the only kind of institution there is.

## 10. References

- ETH-001 -- Ethical Foundations of the Institution (Principle 2: Integrity; Principle 3: Transparency; Principle 6: Honest Accounting)
- CON-001 -- The Founding Mandate (purpose amnesia as institutional failure mode)
- GOV-001 -- Authority Model (decision record definition, tier system, waiting periods, append-only rule R-GOV-03, review requirement R-GOV-02)
- OPS-001 -- Operations Philosophy (operational tempo for quarterly and annual reviews, documentation-first principle)
- D15-002 -- Safeguards Architecture (append-only logs as structural safeguard)
- GOV-003 -- Dispute Resolution Procedures (decision log consultation for precedent)
- GOV-004 -- Succession Planning and Execution (decision log as successor training resource)
- GOV-007 -- Institutional Health Assessment (quarterly governance health check includes decision log review)
- D6-001 -- Data Philosophy (Tier 1 classification for the decision log)
- D6-002 -- Backup Doctrine (daily backup of Tier 1 data including decision log)
- D11-001 -- Administration Philosophy (record-keeping as institutional memory, Lightness Principle)

---

---

*End of Stage 3 Operational Doctrine -- Batch 2*

**Document Total:** 5 articles
**Article IDs:** D12-002, D12-003, D11-002, D13-002, GOV-002
**Status:** All five articles ratified as of 2026-02-16.
**Next Stage:** Continue with additional Stage 3 operational articles across remaining domains.
