# STAGE 4: SPECIALIZED SYSTEMS -- RESEARCH & THEORY ADVANCED

## Domain 14 Articles D14-002 through D14-006

**Document ID:** STAGE4-RESEARCH-ADVANCED
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Classification:** Stage 4 -- Specialized Systems. Advanced reference documents for the institution's research and knowledge infrastructure. These articles extend the philosophy established in D14-001 (Research Philosophy) into operational systems for experimentation, technology evaluation, knowledge organization, negative result management, and multi-generational research continuity.
**Depends On:** ETH-001, CON-001, GOV-001, OPS-001, D14-001, D13-001, D11-001
**Depended Upon By:** All subsequent Domain 14 articles. Referenced by Domain 13 (Evolution) for validated methodology. Referenced by Domain 15 (Ethics & Safeguards) for research integrity context.

---

## How to Read This Document

This document contains five advanced reference articles for Domain 14: Research & Theory. They are Stage 4 documents -- specialized systems that build upon the philosophy of D14-001 and the operational foundations of Stage 3. Where Stage 2 established why the institution conducts research, and Stage 3 established the basic procedures for doing so, Stage 4 provides the detailed systems that make research rigorous, organized, and sustainable across decades.

These articles are interconnected. D14-002 (Experimental Design) produces the data that D14-004 (Knowledge Classification) organizes. D14-005 (Negative Results Registry) captures what D14-002 discovers does not work. D14-003 (Technology Evaluation) applies the experimental methods of D14-002 to a specific and recurring research question. D14-006 (Multi-Generational Research Continuity) ensures that all of it survives the transition from one operator to the next.

If you are reading these for the first time, read them in order. The articles assume familiarity with D14-001 and will not re-argue the philosophical positions established there. They assume you understand the air-gap constraint, the echo-chamber danger, and the knowledge lifecycle. If those concepts are unfamiliar, stop here and read D14-001 first.

If you are a successor operator reading these for the first time, D14-006 is written specifically for you. Start there if you need immediate orientation on the state of the institution's research programs.

---

---

# D14-002 -- Experimental Design for Solo Researchers

**Document ID:** D14-002
**Domain:** 14 -- Research & Theory
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, D14-001, OPS-001, GOV-001
**Depended Upon By:** D14-003, D14-005, D14-006. All Domain 14 articles that involve structured investigation. Referenced by Domain 13 for experimental validation of evolution proposals.

---

## 1. Purpose

This article defines how to design and execute rigorous experiments when the institution consists of a single researcher operating in isolation. It addresses the fundamental challenge that most experimental methodology was developed for research teams -- groups of people who can divide the roles of experimenter, observer, analyst, and critic. When all of those roles collapse into one person, the standard methodology does not simply scale down. It breaks in specific, predictable ways, and those breakages must be addressed with deliberate structural countermeasures.

The purpose is not to pretend that solo research can achieve the same rigor as a well-funded laboratory with peer review. It cannot. D14-001 is honest about that limitation (per ETH-001, Principle 6). The purpose is to extract the maximum possible rigor from a single-person operation -- to identify which elements of experimental methodology still apply, which must be adapted, and which must be replaced with alternatives that serve the same epistemological function.

This article also defines the experiment log format -- the standardized documentation structure that ensures every experiment conducted by this institution is recorded completely enough to be understood, evaluated, and replicated by a future operator who was not present when it was conducted.

## 2. Scope

**In scope:**
- Experimental design principles adapted for a single researcher with no external peer review.
- Variable identification and control in small-scale, resource-constrained environments.
- Blinding and bias mitigation strategies for solo operation.
- Statistical validity at small sample sizes, including when statistics are and are not appropriate.
- The experiment log format: mandatory fields, optional fields, completion criteria.
- Pre-registration of hypotheses and success criteria before experimentation begins.
- The relationship between experiment design and the knowledge classification system (D14-004).

**Out of scope:**
- The philosophical justification for research (D14-001).
- Technology evaluation specifically (D14-003, though it uses these methods).
- Knowledge taxonomy management (D14-004).
- Negative results documentation format (D14-005, though experimental failures feed that registry).
- Multi-generational continuity of research programs (D14-006).

## 3. Background

### 3.1 Why Solo Experimentation Requires Its Own Methodology

The standard scientific method assumes, implicitly, a social infrastructure. Hypotheses are tested by experimenters who report results to peers who scrutinize the methodology, attempt replication, and challenge the conclusions. The rigor of the method depends on this social process as much as it depends on the experimental procedure itself. Remove the social process, and the methodology must compensate.

A solo researcher faces three specific challenges that team-based research does not:

First, there is no role separation. The person who designs the experiment, the person who executes it, the person who observes the results, and the person who interprets them are all the same individual. This means the biases that role separation is designed to mitigate -- observer bias, confirmation bias, expectation effects -- are all fully active in every phase of every experiment.

Second, there is no independent replication. The results of any experiment are correlated with the experimenter's equipment, environment, habits, and blind spots. D14-001 describes adapted replication strategies (temporal, methodological, environmental, documentation-based), but these are mitigations, not solutions.

Third, there is no external calibration. In a research community, the researcher learns what constitutes sufficient evidence by observing the standards applied to others' work. In isolation, the researcher must set their own evidentiary standards, which means those standards are subject to the same biases they are meant to control.

### 3.2 The Pre-Registration Principle

The single most important safeguard against confirmation bias in solo research is pre-registration: the practice of writing down the hypothesis, the experimental procedure, the variables, and the success/failure criteria before the experiment begins. Pre-registration does not prevent bias -- nothing fully prevents bias in a solo operation. But it makes bias visible. When the results arrive and the experimenter is tempted to reinterpret the success criteria to fit the data, the pre-registered criteria are there in writing, immutable, forcing a confrontation between what was expected and what was found.

Pre-registration is mandatory in this institution (see R-EXP-01). It is not a formality. It is the primary structural defense against the most common experimental failure mode: unconsciously adjusting the question to fit the answer.

### 3.3 Statistical Validity and Honesty About Sample Size

Many experiments conducted by this institution will involve small sample sizes -- often extremely small. When you are testing whether a particular backup strategy is reliable, you may run it ten times. When you are evaluating hardware durability, your sample may be a single unit. When you are comparing two approaches to a task, your sample is your own performance under two conditions.

Statistical methods designed for large samples do not apply here, and pretending they do is worse than not using statistics at all. A p-value calculated from three data points is not meaningful. A confidence interval from a sample of one is fiction.

This article takes the position that honest qualitative assessment, clearly labeled as such, is superior to dishonest quantitative analysis. When sample sizes support statistical analysis, use it. When they do not, describe what you observed, acknowledge the limitations, and assign an appropriate confidence level through the knowledge classification system. Do not dress opinion in the costume of data.

## 4. System Model

### 4.1 The Experiment Lifecycle

Every experiment in this institution follows a defined lifecycle with five phases:

**Phase 1: Question Formulation.** The experiment begins with a question that is specific enough to be testable. "Is ZFS a good filesystem?" is not a testable question. "Does ZFS scrub detect and correct single-bit errors on the institution's hardware within 24 hours?" is testable. The question formulation phase ends with a written question and a clear statement of why the answer matters to the institution.

**Phase 2: Pre-Registration.** Before any experimental work begins, the experimenter writes and files a pre-registration document containing: the hypothesis (what the experimenter expects to find), the null hypothesis (what the result looks like if the hypothesis is wrong), the experimental procedure, the variables (independent, dependent, and controlled), the sample size and justification, the success criteria, the failure criteria, and the planned analysis method. The pre-registration is timestamped and filed in the research archive. It cannot be modified after filing; if the experimenter realizes the design is flawed, a new pre-registration is filed and the old one is retained with a note explaining the revision.

**Phase 3: Execution.** The experiment is conducted according to the pre-registered procedure. Deviations from the procedure are recorded in real time in the experiment log. The experimenter records raw observations, not interpretations. "The scrub completed in 4 hours 17 minutes and reported 0 errors" is a raw observation. "The scrub worked as expected" is an interpretation, and it belongs in Phase 4.

**Phase 4: Analysis.** The raw data is analyzed according to the pre-registered analysis method. The analysis compares the observed results against the pre-registered success and failure criteria. The analysis must explicitly address: did the results match the hypothesis, the null hypothesis, or neither? Were there unexpected observations? Were there deviations from the procedure that may have affected the results? The analysis is documented in full, including the reasoning process, not just the conclusion.

**Phase 5: Classification and Archival.** The results are classified according to D14-004 (Knowledge Classification) and archived according to the research archive procedures. If the results are negative (the hypothesis was disproven or the experiment was inconclusive), an entry is created in the Negative Results Registry (D14-005).

### 4.2 Variable Control in a Solo Environment

Classical experimental design distinguishes between independent variables (what the experimenter manipulates), dependent variables (what the experimenter measures), and controlled variables (what the experimenter holds constant). In a solo environment with limited resources, controlling variables is the primary methodological challenge.

**Environmental controls.** Many experiments are conducted on the institution's production or test hardware, which means the environment cannot be fully isolated. The experimenter must document the environmental conditions under which the experiment is conducted: what else was running on the system, what the ambient temperature was, what the power state was, what the time of day was. These become part of the experiment record and allow future analysis of whether environmental factors affected the results.

**Temporal controls.** Experiments should be conducted at consistent times and under consistent conditions when possible. If an experiment must be run multiple times, the runs should be spaced to avoid temporal confounds (such as system load varying by time of day, or battery state varying by season).

**Blinding adaptations.** True blinding -- where the experimenter does not know which condition they are testing -- is generally impossible in a solo environment. The adapted alternatives are: delayed analysis (record raw data, then analyze it at least 48 hours later, after the emotional investment in the outcome has faded), pre-registered criteria (the success/failure criteria are locked before the data arrives), and quantitative measurement (where possible, use numerical measurements rather than subjective assessments, as numbers are harder to unconsciously distort).

### 4.3 The Experiment Log Format

Every experiment is documented in a standardized experiment log with the following mandatory fields:

- **Experiment ID:** A unique identifier in the format EXP-YYYY-NNN (year and sequential number).
- **Date initiated:** The date pre-registration was filed.
- **Date completed:** The date the final analysis was written.
- **Question:** The specific, testable question.
- **Hypothesis:** What the experimenter expected to find.
- **Null hypothesis:** What the result looks like if the hypothesis is wrong.
- **Variables:** Independent, dependent, and controlled.
- **Procedure:** Step-by-step instructions for executing the experiment.
- **Sample size and justification:** How many trials, and why that number.
- **Success criteria:** Pre-registered, specific, measurable.
- **Failure criteria:** Pre-registered, specific, measurable.
- **Raw data:** Observations recorded during execution, without interpretation.
- **Deviations:** Any departures from the pre-registered procedure.
- **Analysis:** Interpretation of the data against the criteria.
- **Conclusion:** The result, classified by confidence level per D14-004.
- **Negative result entry:** If applicable, reference to the D14-005 registry entry.
- **Implications:** What this result means for institutional practice.
- **Related experiments:** Cross-references to prior or subsequent experiments.
- **Adversarial review:** The documented attempt to disprove the conclusion (mandatory for verified results per D14-001, R-RES-04).

Optional fields include: photographs or diagrams, equipment serial numbers, software versions, environmental measurements, and links to raw data files.

## 5. Rules & Constraints

- **R-EXP-01:** Pre-registration is mandatory for all formal experiments. No experimental result may be classified above "speculative" without a pre-registration document filed before execution began.
- **R-EXP-02:** The experiment log format (Section 4.3) is mandatory for all formal experiments. Informal observations and quick tests do not require the full format but must still be documented in the operational log with a note indicating they were not conducted under formal experimental protocols.
- **R-EXP-03:** Raw data must be recorded during execution, not reconstructed afterward. If raw data was not recorded, the experiment's maximum confidence classification is "provisional," regardless of the result.
- **R-EXP-04:** No experiment may be conducted on production systems without a documented risk assessment approved through the appropriate governance tier per GOV-001 and D14-001 R-RES-08. Test environments must be used unless the question specifically requires production conditions, and this requirement must be documented and justified.
- **R-EXP-05:** Statistical methods may only be applied when the sample size and data characteristics justify their use. The experimenter must document why a given statistical test is appropriate. If the sample size does not support statistical analysis, the result must be reported as a qualitative observation with an explicit statement of the sample limitation.
- **R-EXP-06:** Every experiment must be documented completely enough that a future operator, unfamiliar with the experiment, could reproduce the procedure from the documentation alone. This is tested by the documentation replication method described in D14-001, Section 4.3.
- **R-EXP-07:** Experiments that produce negative results must file an entry in the Negative Results Registry (D14-005) within 30 days of experiment completion.

## 6. Failure Modes

- **Pre-registration drift.** The experimenter modifies the pre-registration after seeing preliminary results, consciously or unconsciously adjusting the criteria to fit the data. Mitigation: pre-registration documents are timestamped and filed immutably. Revisions require a new filing with an explanation. The original is never deleted.
- **Observer effect contamination.** The experimenter's expectations influence what they observe, particularly with subjective measurements. Mitigation: prefer quantitative measurement. Use delayed analysis. Pre-register success criteria in numerical terms where possible.
- **Sample size delusion.** The experimenter treats a small number of observations as statistically significant, applying formal tests that are not valid for the sample size. Mitigation: R-EXP-05 requires explicit justification of statistical methods. The default is qualitative assessment with acknowledged limitations.
- **Procedure drift.** Over the course of a multi-session experiment, the experimenter gradually deviates from the pre-registered procedure without noticing or recording the deviations. Mitigation: the experiment log requires real-time recording of deviations. The procedure should be physically present (printed or displayed) during each experimental session.
- **Analysis paralysis.** The experimenter generates so much raw data that analysis becomes overwhelming and is deferred indefinitely. Mitigation: pre-register the analysis method. The analysis method defines what data is relevant to the conclusion. Everything else is supplementary.
- **Rigor theater.** The experimenter follows the formal procedures -- pre-registration, log format, analysis -- without genuine intellectual engagement. The forms are completed, but the adversarial self-critique required by D14-001 is absent. Mitigation: the adversarial review field in the experiment log is not optional for verified results. The review must contain a genuine attempt to disprove the conclusion, and "I could not think of any objections" is not a genuine attempt.

## 7. Recovery Procedures

1. **If pre-registration has been neglected:** Do not retroactively create pre-registrations for completed experiments. That defeats their purpose. Instead, classify all un-pre-registered results as "provisional" at best. Begin pre-registering all new experiments immediately. Review the un-pre-registered results with explicit awareness that confirmation bias may have affected the conclusions.

2. **If experiment logs are incomplete:** Reconstruct what can be reconstructed from operational logs, notes, and memory. Mark reconstructed entries explicitly. Accept that incompletely documented experiments carry lower confidence. Establish the habit of recording in real time by keeping the log physically present during experimentation.

3. **If statistical methods have been misapplied:** Conduct a statistical audit. For each experiment that used formal statistical analysis, verify that the sample size and data characteristics justified the method used. Re-analyze with appropriate methods or downgrade to qualitative assessment. Update the knowledge classification accordingly.

4. **If rigor theater is occurring:** This is the hardest failure to recover from because it requires honest self-assessment. Review the last five adversarial reviews. For each, ask: did I genuinely try to disprove this conclusion? If the adversarial reviews are pro forma, re-conduct them with a specific mandate: identify at least one plausible alternative explanation or methodological weakness for each result. If none can be found, that is acceptable -- but the search must be genuine.

5. **If experiments are being avoided entirely:** Return to D14-001 R-RES-06, which requires at least two formal investigations per year. Identify the simplest, lowest-risk experimental question relevant to current institutional operations. Design and execute it. The goal is to restart the practice, not to produce significant results.

## 8. Evolution Path

- **Years 0-5:** Experimental methodology is being learned and practiced. The first experiments will likely be crude -- simple comparisons, basic performance measurements, straightforward technology tests. This is expected. The value of these early experiments is in building the habit and refining the process, not in producing groundbreaking results. The experiment log format may need adjustment as practical experience reveals fields that are missing or unnecessary.
- **Years 5-15:** The experiment archive grows. Patterns emerge across experiments. The experimenter becomes better at identifying confounds, designing controls, and recognizing their own biases. Temporal replication of early experiments becomes possible and may reveal that some early conclusions were wrong. This is a feature, not a failure.
- **Years 15-30:** If succession has occurred, the new operator brings genuinely different biases and perspectives to experimental design. This is the closest the institution comes to independent replication. The experiment archive serves as a training resource, showing the successor how the institution conducts research and what has been learned.
- **Years 30-50+:** The experiment archive is a longitudinal dataset spanning decades. It contains not just individual results but the institution's evolving understanding of its own systems. Hardware generations, software paradigms, and operational practices have changed. The archive records how those changes were evaluated and adopted.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
I am writing an experimental methodology document for a laboratory that consists of one person, one desk, and whatever hardware I have built. The absurdity is not lost on me. But the absurdity is the point. The experiments I conduct here -- on backup reliability, on hardware longevity, on software stability -- are not academic exercises. They are the basis for decisions that affect whether this institution survives the next decade. If I test a backup strategy badly and conclude it works when it does not, I will discover the error at the worst possible time: during a recovery. So I will test carefully, document thoroughly, and be honest about what the results do and do not tell me. The formality is not for show. It is for the version of me who needs to trust these results at 3 AM during an actual failure.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 2: Integrity Over Convenience; Principle 6: Honest Accounting of Limitations)
- CON-001 -- The Founding Mandate (self-contained knowledge, lifetime operation)
- GOV-001 -- Authority Model (governance tiers for experiments on production systems)
- OPS-001 -- Operations Philosophy (documentation-first, operational logging)
- D14-001 -- Research Philosophy (knowledge lifecycle, anti-echo-chamber structures, replication adaptations, R-RES-04, R-RES-06, R-RES-08)
- D14-004 -- Knowledge Classification and Taxonomy (confidence levels, classification system)
- D14-005 -- Negative Results Registry (documentation of failed experiments)
- Stage 1 Documentation Framework, Domain 14 (article list, dependencies)

---

---

# D14-003 -- Technology Evaluation Framework

**Document ID:** D14-003
**Domain:** 14 -- Research & Theory
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, SEC-001, OPS-001, D14-001, D14-002, D13-001
**Depended Upon By:** D14-004, D14-005, D14-006. All Domain 13 (Evolution) articles involving technology adoption. All domain articles involving technology selection.

---

## 1. Purpose

This article defines the framework by which this institution evaluates new technologies for potential adoption. It is one of the most consequential research activities the institution conducts, because technology adoption decisions are among the hardest to reverse. A filesystem choice may last decades. A database engine becomes the foundation of years of accumulated data. A hardware platform determines what software can run and what upgrades are possible. These decisions must be made carefully, with a structured process that resists both the excitement of novelty and the inertia of familiarity.

D14-001 establishes that the institution must avoid both stagnation and reckless adoption. This article operationalizes that balance by providing a repeatable evaluation framework -- a rubric, a process, and a documentation format that ensures every technology decision is made with full awareness of its implications and recorded thoroughly enough that future operators understand why the decision was made.

The framework is designed for the specific constraints of an air-gapped, off-grid, single-operator institution. The evaluation criteria differ substantially from those used by internet-connected organizations. Uptime SLAs, cloud integration, automatic update mechanisms, and vendor support contracts are irrelevant here. What matters is longevity, maintainability without internet access, air-gap compatibility, documentation quality, and the health of the community that produces the technology -- because that community is the source of future updates, knowledge, and fixes that must be manually transported across the air gap.

## 2. Scope

**In scope:**
- The technology evaluation rubric: criteria, weighting, and scoring.
- The evaluation process: from initial identification through proof-of-concept to adoption decision.
- Air-gap compatibility assessment: specific tests and requirements.
- Documentation quality assessment: what constitutes sufficient documentation for an isolated institution.
- Community health assessment: indicators of project longevity and sustainability.
- The proof-of-concept procedure: how to test a technology safely before committing to it.
- The Technology Decision Record (TDR): format and filing requirements.
- The relationship between technology evaluation and the evolution framework (D13-001).

**Out of scope:**
- General experimental methodology (D14-002, though technology evaluation uses those methods).
- Specific technology choices for specific domains (individual domain articles).
- The procurement process for acquiring technology across the air gap (D11-001, ADM-series).
- Security assessment of new software (SEC-001, SEC-series, though security is one evaluation criterion).

## 3. Background

### 3.1 The Technology Adoption Trap

Technology adoption in an air-gapped institution presents a double bind. On one side is the danger of stagnation: clinging to known technology until it becomes unsupportable, its community disbands, and its documentation rots. On the other side is the danger of churn: adopting every promising new technology, enduring constant migration, and never achieving the stability that a fifty-year institution requires.

Both dangers are real. Both have destroyed institutions. The history of computing is littered with organizations that bet on technologies that seemed permanent at the time -- and equally littered with organizations that refused to move until it was too late to move gracefully.

The evaluation framework exists to navigate between these dangers. It provides a structured basis for the question: should we adopt this technology? And it provides an equally structured basis for the question: should we stay with what we have?

### 3.2 The Air-Gap Filter

The air gap is the most powerful filter in the technology evaluation framework. Technologies that require constant internet connectivity are immediately disqualified. Technologies that depend on remote repositories for routine operation are disqualified. Technologies whose update mechanism assumes network access are not disqualified but receive a significant penalty in the evaluation rubric, because every update will require manual intervention to transport across the air gap.

This filter eliminates a large fraction of modern software. That elimination is a feature, not a bug. The remaining candidates -- software that can be installed from local media, operated without network access, updated through manual package installation, and documented thoroughly enough to troubleshoot without searching the internet -- are inherently more suitable for this institution's constraints.

### 3.3 The Longevity Question

The most important evaluation criterion, and the hardest to assess, is longevity. Will this technology still be viable in ten years? Twenty? Will its community still exist? Will its file formats still be readable? Will replacement hardware still run it?

No one can predict the future with certainty. But there are indicators. Technologies with open-source licenses are more likely to survive than proprietary ones, because the code persists even if the original developer abandons it. Technologies with simple, well-documented file formats are more resilient than those with opaque binary formats. Technologies with large, diverse communities are more likely to survive the loss of any single contributor than those dependent on a single developer.

The evaluation rubric encodes these indicators. It does not predict the future. It assesses the probability of survival based on observable present characteristics.

## 4. System Model

### 4.1 The Evaluation Rubric

Every technology evaluation uses a standardized rubric with six criteria, each scored on a five-point scale (1 = unacceptable, 2 = poor, 3 = adequate, 4 = good, 5 = excellent):

**Criterion 1: Longevity (weight 3x).** How likely is this technology to remain viable over 10-20+ years? Indicators include: age of the project, size and diversity of the contributor community, adoption breadth, license type (open source strongly preferred), governance structure, historical stability, and the existence of a clear succession plan or foundation stewardship. A technology that has been stable for a decade scores higher than one released last year, regardless of feature comparison.

**Criterion 2: Maintainability (weight 3x).** How difficult is it to maintain this technology without internet access and without vendor support? Indicators include: complexity of configuration, clarity of error messages, quality of logging, dependency count and depth, build reproducibility, and the skill level required for routine maintenance. A technology that the operator can fully understand and troubleshoot is more valuable than one that is technically superior but opaque.

**Criterion 3: Air-Gap Compatibility (weight 3x).** Can this technology be installed, operated, updated, and maintained entirely offline? Specific tests include: installation from local media, operation without any network call, update application from locally-stored packages, license validation without phone-home, and documentation availability offline. Any technology that phones home for license validation receives an automatic score of 1.

**Criterion 4: Documentation Quality (weight 2x).** How good is the available documentation, and can it be stored locally? Indicators include: completeness, accuracy, currency, availability in offline-friendly formats (HTML dump, PDF, man pages, plain text), quality of error documentation, availability of architecture documents (not just tutorials), and the presence of troubleshooting guides for common failure modes.

**Criterion 5: Data Sovereignty (weight 2x).** Does the institution retain full control of its data when using this technology? Indicators include: open and documented file/storage formats, availability of export tools, absence of vendor lock-in mechanisms, ability to read data without the original software, and compatibility with the institution's backup and archive strategies.

**Criterion 6: Community Health (weight 1x).** What is the current state of the project's community? Indicators include: contributor activity (commits, releases, bug fixes), responsiveness to bug reports, quality of community discourse, diversity of contributors (not dependent on a single person or company), and the trajectory (growing, stable, or declining).

The weighted score produces a total out of 70 (maximum). The following thresholds guide -- but do not dictate -- the adoption decision:

- **56-70 (Strong Candidate):** Technology meets or exceeds institutional requirements. Proceed to proof-of-concept.
- **42-55 (Qualified Candidate):** Technology has notable strengths but also significant gaps. Proceed to proof-of-concept only if the gaps can be mitigated and the alternative is worse.
- **28-41 (Weak Candidate):** Technology has fundamental concerns. Do not proceed unless no alternative exists and the need is critical.
- **Below 28 (Disqualified):** Technology is unsuitable for this institution.

### 4.2 The Evaluation Process

Technology evaluation follows a defined sequence:

**Step 1: Identification.** A technology is identified as a potential candidate, either during a scheduled technology review cycle or in response to a specific institutional need. The identification includes a brief statement of what problem the technology would solve and why current solutions are insufficient.

**Step 2: Preliminary Assessment.** The technology is evaluated against the rubric using publicly available information (documentation, community forums captured during air-gap crossings, published reviews from the institution's reference library). This is a desk evaluation, not a hands-on test. The goal is to determine whether the technology merits the investment of a proof-of-concept.

**Step 3: Proof-of-Concept.** If the preliminary assessment scores 42 or above, a proof-of-concept is conducted. The proof-of-concept is a structured experiment following D14-002 methodology: pre-registered hypothesis, defined success criteria, controlled execution, documented results. The proof-of-concept tests the technology on non-production systems under conditions that approximate production use. Minimum duration: 30 days of active use, or the duration necessary to exercise all major features relevant to the institution's use case, whichever is longer.

**Step 4: Decision.** The proof-of-concept results, the rubric score, and the overall assessment are compiled into a Technology Decision Record (TDR). The adoption decision is made through the appropriate governance tier per GOV-001. The TDR is filed in the research archive regardless of the decision.

**Step 5: Integration or Rejection.** If adopted, the technology is integrated through the evolution framework (D13-001). If rejected, the TDR is filed in the research archive and referenced in the Negative Results Registry (D14-005) if the rejection was based on specific technical failures discovered during the proof-of-concept.

### 4.3 The Technology Decision Record (TDR)

Every technology evaluation produces a TDR with the following mandatory fields:

- **TDR ID:** Format TDR-YYYY-NNN.
- **Date:** Date of final decision.
- **Technology evaluated:** Name, version, source.
- **Problem statement:** What institutional need prompted the evaluation.
- **Current solution:** What the institution currently uses for this need (if anything).
- **Rubric scores:** All six criteria with individual scores, justifications, and weighted total.
- **Proof-of-concept summary:** Duration, scope, key findings (or "Not conducted" with justification).
- **Decision:** Adopt, Reject, or Defer (with conditions for revisiting).
- **Rationale:** The reasoning behind the decision, including dissenting considerations.
- **Migration plan reference:** If adopted, reference to the migration plan in Domain 13.
- **Review date:** When this decision should be revisited (mandatory, even for rejections).
- **Related TDRs:** Cross-references to evaluations of competing or complementary technologies.

## 5. Rules & Constraints

- **R-TEV-01:** No technology may be adopted for production use without a completed Technology Decision Record. Emergency adoptions (e.g., during disaster recovery) must have a TDR filed within 90 days of adoption.
- **R-TEV-02:** The evaluation rubric (Section 4.1) must be applied in full for every formal evaluation. Individual criteria may not be skipped because the evaluator considers them irrelevant. If a criterion is genuinely not applicable, it must be scored and the justification for the score documented.
- **R-TEV-03:** Proof-of-concept testing must be conducted on non-production systems unless the nature of the evaluation specifically requires production conditions, and this requirement is documented and approved per GOV-001.
- **R-TEV-04:** Technologies with an Air-Gap Compatibility score of 1 may not be adopted under any circumstances. The air gap is non-negotiable per CON-001 and SEC-001.
- **R-TEV-05:** All Technology Decision Records must be retained permanently in the research archive, including records for rejected technologies. Rejection records are as valuable as adoption records -- they prevent future operators from re-evaluating technologies that have already been found unsuitable.
- **R-TEV-06:** Every adopted technology must have a scheduled review date, no more than five years from adoption, at which the technology is re-evaluated against the current rubric. Technologies may be re-confirmed, scheduled for replacement, or flagged for active migration.
- **R-TEV-07:** Technology evaluations must include an exit assessment: how difficult would it be to migrate away from this technology if it proves unsuitable after adoption? Technologies with high exit costs must score correspondingly higher on all other criteria to justify the risk.

## 6. Failure Modes

- **Novelty bias.** The evaluator is excited by a new technology and unconsciously inflates rubric scores to justify adoption. Mitigation: the rubric requires written justification for each score. The adversarial review from D14-001 applies to technology evaluations. A second evaluation of the same technology after a 30-day cooling period is recommended for high-impact adoptions.
- **Familiarity bias.** The evaluator unconsciously penalizes new technologies to justify staying with current solutions. Mitigation: the rubric must also be applied to the current solution. If the current solution scores lower than the candidate, the burden of justification shifts to explaining why the current solution should be retained.
- **Proof-of-concept theater.** The proof-of-concept is conducted as a formality, using unrealistically favorable conditions that do not reflect actual institutional use. Mitigation: the proof-of-concept must follow D14-002 experimental methodology with pre-registered success criteria. The conditions must approximate production use, including failure injection where appropriate.
- **Decision record neglect.** Technology decisions are made informally without completing the TDR. Over time, no one remembers why a particular technology was chosen. Mitigation: R-TEV-01 makes the TDR mandatory. No technology is considered officially adopted without one.
- **Review date amnesia.** The scheduled review date passes without a review. Technologies that should have been replaced or re-evaluated persist on inertia. Mitigation: review dates are recorded in the institutional scheduling system per D11-001. They are treated as mandatory calendar events.
- **Community collapse blindness.** An adopted technology's community deteriorates, but the institution does not notice because it is behind the air gap. Mitigation: the air-gap crossing process (Domain 18) should include a scheduled check of adopted technology communities. Community health is a persistent monitoring obligation, not a one-time assessment.

## 7. Recovery Procedures

1. **If technologies have been adopted without TDRs:** Conduct retroactive evaluations. Apply the rubric to all currently adopted technologies. File TDRs for each. Identify any technologies that would not have passed the evaluation and schedule replacement planning for them.

2. **If review dates have been missed:** Conduct the reviews immediately. Prioritize technologies that have been in use longest without review. Accept that some reviews will reveal that the institution has been using unsuitable technologies, and begin migration planning where necessary.

3. **If a community collapse has been discovered:** Assess the severity. If the technology is still functional but unsupported, begin searching for alternatives immediately. If the technology is actively deteriorating (bugs not fixed, security vulnerabilities unpatched), escalate to a priority evaluation and migration per the evolution framework.

4. **If the evaluation process has degraded into formality:** Re-read this article and D14-001. Conduct a genuine, adversarial evaluation of the most recently adopted technology. If the evaluation reveals concerns that were missed, re-evaluate other recent adoptions with the same rigor.

5. **If exit costs are discovered to be higher than assessed:** Document the true exit costs. Update the TDR. If the technology must still be replaced, the exit cost assessment informs the migration timeline. Use this experience to improve exit cost assessment in future evaluations.

## 8. Evolution Path

- **Years 0-5:** The framework is being applied to the institution's founding technology stack. These initial evaluations establish the baseline TDR library and calibrate the rubric against real-world experience. Some rubric weights may need adjustment as the institution discovers which criteria are most predictive of actual suitability.
- **Years 5-15:** The TDR library grows into a decision history. Patterns become visible: which criteria best predicted success? Which technologies lasted? Which failed unexpectedly? These patterns should inform rubric refinements. Scheduled reviews of adopted technologies begin producing re-evaluation data.
- **Years 15-30:** Technology generations have turned over. The institution has experienced at least one major technology migration and can evaluate the framework's effectiveness based on how well it predicted success and failure. The framework itself is a technology that must be evaluated and evolved.
- **Years 30-50+:** The TDR library is a multi-decade record of the institution's technology choices. It tells the story of what was tried, what worked, what failed, and why. For a successor operator, it is an invaluable guide to the institution's technological identity and the reasoning behind it.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
Every technology decision I make now is a bet on the future. Some of those bets will be wrong. The purpose of the TDR is not to guarantee that every bet wins -- that is impossible. The purpose is to ensure that every bet is made consciously, with documented reasoning, so that the person who has to live with the consequences (whether that is future-me or a successor) knows why the bet was placed and can decide intelligently whether to let it ride or cut losses. The worst technology decision is not the one that turns out to be wrong. It is the one that nobody remembers making, for reasons nobody recorded, that nobody knows how to reverse.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 2: Integrity Over Convenience; Principle 6: Honest Accounting)
- CON-001 -- The Founding Mandate (air-gap constraint, lifetime operation, self-sovereignty)
- GOV-001 -- Authority Model (decision tiers for technology adoption)
- SEC-001 -- Threat Model and Security Philosophy (air-gap requirements, supply chain risk)
- OPS-001 -- Operations Philosophy (sustainability, complexity budget)
- D14-001 -- Research Philosophy (knowledge lifecycle, anti-echo-chamber structures)
- D14-002 -- Experimental Design for Solo Researchers (proof-of-concept methodology)
- D13-001 -- Evolution Philosophy (technology migration, change management)
- D11-001 -- Administration Philosophy (resource allocation, scheduling)
- D14-005 -- Negative Results Registry (documentation of rejected technologies)

---

---

# D14-004 -- Knowledge Classification and Taxonomy

**Document ID:** D14-004
**Domain:** 14 -- Research & Theory
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, OPS-001, D14-001, D14-002
**Depended Upon By:** D14-005, D14-006. All Domain 14 articles that produce or consume classified knowledge. All domains that reference institutional knowledge with confidence levels.

---

## 1. Purpose

This article defines how the institution organizes what it knows. Knowledge without organization is a pile. Knowledge with organization is a library. The difference between the two is not the knowledge itself but the structures that make it findable, assessable, and usable. This article provides those structures.

D14-001 introduced the knowledge lifecycle (observation, investigation, validation, classification, integration, archival) and established four confidence levels (verified, provisional, speculative, deprecated). This article expands those concepts into a complete taxonomy -- a system of categories, relationships, and confidence levels that allows the institution to answer three critical questions about any piece of knowledge: What kind of knowledge is this? How confident are we in it? How does it relate to other things we know?

These questions matter because knowledge that cannot be found is knowledge that does not exist, operationally speaking. Knowledge whose confidence level is unknown is either trusted too much or too little. Knowledge whose relationships are untracked creates blind spots -- the institution changes one thing without realizing it depends on another. The taxonomy is the defense against all three failures.

## 2. Scope

**In scope:**
- The knowledge taxonomy: categories, subcategories, and classification rules.
- Confidence levels: definitions, assignment criteria, promotion and demotion procedures.
- Relationship types: dependencies, derivations, contradictions, and supersessions.
- The classification process: how new knowledge is classified upon creation.
- The taxonomy review process: how the taxonomy itself evolves.
- Cross-domain knowledge integration: how knowledge from one domain relates to knowledge in others.
- The knowledge registry: the master index of all classified institutional knowledge.

**Out of scope:**
- The philosophical justification for research and knowledge management (D14-001).
- Experimental methodology for producing new knowledge (D14-002).
- Technology evaluation specifically (D14-003).
- Negative results documentation format (D14-005, though negative results are classified within this taxonomy).
- Physical storage and preservation of knowledge artifacts (RES-013, archival procedures).

## 3. Background

### 3.1 The Problem of Institutional Amnesia

Institutions forget. Not through dramatic data loss -- though that happens -- but through the slow accumulation of knowledge that is recorded but not organized, filed but not findable, known to one person but not externalized into the institutional record. When the person who knows something leaves, retires, or simply forgets, the knowledge vanishes. The institution then rediscovers it, repeating work that was already done, or worse, proceeds without it and makes decisions based on incomplete understanding.

This institution is especially vulnerable to institutional amnesia because it is operated by one person. There is no colleague to remind you that "we tried that in 2029 and it did not work." There is no institutional folklore passed through conversation. If knowledge is not in the record, it is not in the institution.

The taxonomy is the antidote. By requiring that all significant knowledge be classified, indexed, and related to other knowledge, it creates a structure that survives individual memory failure.

### 3.2 The Confidence Problem

Not all knowledge is equally reliable, and treating it as such is dangerous. A result that has been replicated through multiple methods over five years is not the same as a preliminary observation from a single experiment. A procedure that has been tested in disaster recovery drills is not the same as one that was written but never executed. Yet without a formal confidence system, both occupy the same status in the institutional record: they are written down, and they are assumed to be true.

The confidence level system introduced in D14-001 and formalized here addresses this problem. Every piece of classified knowledge carries an explicit confidence level that tells the reader: this is how much you should trust this. The levels are not grades. "Speculative" is not a failing mark -- it is an honest assessment that the evidence is preliminary. "Verified" is not a gold star -- it is a statement that specific evidentiary criteria have been met. The system values honesty over prestige.

### 3.3 The Taxonomy as Living System

A taxonomy created at founding and never revised will become useless within a decade. The institution's knowledge will outgrow the categories. New domains of inquiry will emerge that do not fit the original structure. Relationships between knowledge items will become apparent that the original taxonomy did not anticipate.

The taxonomy must therefore be a living system -- regularly reviewed, periodically revised, and designed from the start to accommodate growth. This article defines not just the taxonomy but the process for evolving it.

## 4. System Model

### 4.1 Knowledge Categories

The taxonomy organizes institutional knowledge into seven primary categories:

**Category 1: Technical Knowledge.** How things work. Hardware specifications, software behavior, system configurations, performance characteristics, failure modes, and troubleshooting procedures. This is the institution's understanding of its own infrastructure.

**Category 2: Procedural Knowledge.** How things are done. Step-by-step procedures, workflows, checklists, and protocols. This includes everything from daily maintenance routines to disaster recovery sequences. Procedural knowledge is distinct from technical knowledge: you can understand how a filesystem works (technical) without knowing the institution's specific procedure for checking filesystem integrity (procedural).

**Category 3: Analytical Knowledge.** What has been learned through investigation. Experimental results, comparative analyses, technology evaluations, performance studies, and failure analyses. This is the output of the research function defined in D14-001 and executed through D14-002.

**Category 4: Contextual Knowledge.** Why things are the way they are. Design decisions, architectural rationales, historical context, and the reasoning behind institutional choices. This is the knowledge that answers "why?" rather than "how?" -- and it is often the first knowledge lost in institutional transitions, because people document what they did but not why they did it.

**Category 5: Environmental Knowledge.** What exists outside the institution. Technology trends, community health assessments, supply chain information, regulatory changes, and external threat intelligence. This knowledge is gathered during air-gap crossings and from the institution's reference library. It is inherently less current than internal knowledge and must be dated and contextualized accordingly.

**Category 6: Negative Knowledge.** What does not work, what has been disproven, and what has been tried and abandoned. This category is the complement of the Negative Results Registry (D14-005). It is perhaps the most valuable category for preventing repeated failures, and it is the category most institutions neglect because recording failures feels unproductive.

**Category 7: Meta-Knowledge.** Knowledge about the institution's knowledge. The taxonomy itself, the classification procedures, the confidence level definitions, the research methodology, and the evolution processes. Meta-knowledge is what allows the institution to manage and improve its own knowledge infrastructure.

### 4.2 Confidence Levels

Each classified knowledge item carries one of four confidence levels, as introduced in D14-001:

**Verified.** The knowledge has been established through rigorous investigation (D14-002), has survived adversarial review (D14-001, Section 4.2), and has been replicated through at least one adapted replication method (temporal, methodological, environmental, or documentation-based). Verified knowledge may be used as a basis for operational decisions and cited as evidence in further research. Promotion to Verified requires: a completed experiment log with pre-registration, a documented adversarial review, and a documented replication.

**Provisional.** The knowledge is supported by a single investigation or a limited evidence base. It has not been replicated. It may be used cautiously in operational decisions with explicit acknowledgment of its provisional status. Most new experimental results enter the taxonomy at this level. Promotion to Provisional requires: a completed experiment log or a thorough analytical assessment with documented methodology and evidence.

**Speculative.** The knowledge is based on preliminary observation, theoretical reasoning, or incomplete evidence. It should not be used as a basis for significant operational decisions without further investigation. Speculative knowledge is valuable because it identifies areas for future research. Entry at Speculative requires: a documented observation or hypothesis with a clear statement of the evidence supporting it and the evidence lacking.

**Deprecated.** The knowledge was previously classified at a higher level but has been contradicted by subsequent evidence, superseded by better knowledge, or invalidated by changed conditions. Deprecated knowledge is not deleted from the registry. It is retained with its deprecation justification and a reference to the knowledge that supersedes it. Deprecated knowledge serves as a historical record and a warning against re-adoption.

### 4.3 Relationship Types

Knowledge items in the taxonomy are connected through four relationship types:

**Depends-on.** Knowledge item A depends on knowledge item B if A's validity requires B's validity. If B is deprecated, A must be reviewed and potentially deprecated or re-validated independently.

**Derived-from.** Knowledge item A is derived from knowledge item B if A was produced through investigation or analysis of B. The derivation chain is a provenance record that allows future operators to trace conclusions back to their evidence base.

**Contradicts.** Knowledge item A contradicts knowledge item B if they cannot both be true. Contradictions must be resolved through further investigation or by deprecating one item. Unresolved contradictions are flagged in the registry and prioritized for research.

**Supersedes.** Knowledge item A supersedes knowledge item B if A is a more recent, more complete, or more accurate treatment of the same subject. Supersession does not delete B -- it retains B as historical context and marks A as the current authoritative source.

### 4.4 The Knowledge Registry

The knowledge registry is the master index of all classified institutional knowledge. It is a structured document (or set of documents) in plain text format, compliant with the 50-year continuity rules, containing for each entry:

- **Knowledge ID:** Format KR-YYYY-NNN (year and sequential number).
- **Title:** A concise description of the knowledge item.
- **Category:** One of the seven primary categories (Section 4.1).
- **Confidence level:** Verified, Provisional, Speculative, or Deprecated.
- **Date classified:** When the item was added to the registry.
- **Date last reviewed:** When the item's classification was last confirmed or changed.
- **Source:** Reference to the experiment log, analysis, observation, or external source.
- **Relationships:** Depends-on, Derived-from, Contradicts, and Supersedes links.
- **Summary:** A brief (one to three paragraph) statement of the knowledge.
- **Full reference:** Pointer to the complete document, experiment log, or data source.

### 4.5 The Taxonomy Review Process

The taxonomy itself is reviewed on a defined cycle:

**Annual classification audit.** Once per year, a random sample of at least 10% of registry entries is reviewed for accuracy of classification. Are the confidence levels still appropriate? Are the relationships still current? Has any knowledge been superseded without the registry being updated? The audit results are documented and any corrections are applied immediately.

**Triennial taxonomy review.** Every three years, the category structure itself is reviewed. Are the seven categories still adequate? Do any categories need to be split, merged, or redefined? Are there significant bodies of knowledge that do not fit any existing category? The review may result in taxonomy modifications, which are documented as taxonomy evolution events in the registry.

**Event-triggered review.** Certain events trigger an immediate taxonomy review of affected entries: a major technology migration, a disaster recovery event, a succession event, or the deprecation of a significant body of knowledge. The triggered review covers all knowledge items related to the event and ensures the registry accurately reflects the institution's post-event knowledge state.

## 5. Rules & Constraints

- **R-KCT-01:** Every significant knowledge item produced by the institution's research function must be classified and entered in the knowledge registry within 60 days of its production. "Significant" means any knowledge that may inform an operational decision or be referenced by a future investigation.
- **R-KCT-02:** Confidence levels must be assigned based on the criteria defined in Section 4.2, not on the operator's intuitive sense of how reliable the knowledge is. The criteria are the standard, not the feeling.
- **R-KCT-03:** Relationships between knowledge items must be documented at the time of classification. The classifier must actively check for dependencies, derivations, contradictions, and supersessions with existing registry entries.
- **R-KCT-04:** Deprecated knowledge must never be deleted from the registry. Deprecation is a status change, not a deletion. The full history of a knowledge item -- including its original classification, any promotions or demotions, and its eventual deprecation -- must be preserved.
- **R-KCT-05:** The annual classification audit (Section 4.5) is mandatory and must be completed within the first quarter of each calendar year. The results must be documented in the research archive.
- **R-KCT-06:** Taxonomy modifications (changes to the category structure, confidence level definitions, or relationship types) are Tier 2 decisions under GOV-001 and require the 30-day waiting period. The taxonomy is institutional infrastructure and must not be changed casually.
- **R-KCT-07:** All knowledge registry data must be stored in plain text formats compliant with the 50-year continuity rules. The registry must be readable without any specialized software.

## 6. Failure Modes

- **Classification backlog.** Knowledge is produced faster than it is classified. The registry falls behind, and an increasing body of institutional knowledge exists outside the taxonomy. Mitigation: R-KCT-01's 60-day deadline. If the backlog grows despite this, reduce the classification effort per item (a brief classification is better than no classification) or increase the time allocated to classification in the operational tempo.
- **Confidence inflation.** Knowledge is routinely assigned higher confidence than the evidence supports, typically because the operator wants to feel certain rather than provisional. Mitigation: R-KCT-02 requires criteria-based assignment. The annual audit (R-KCT-05) provides a systematic check.
- **Relationship neglect.** Knowledge items are classified individually but their relationships are not documented. The registry becomes a flat list rather than a connected graph. Mitigation: R-KCT-03 requires active relationship checking at classification time. The audit verifies relationship documentation.
- **Taxonomy ossification.** The taxonomy is never revised, and knowledge items are forced into categories that no longer fit. Mitigation: the triennial taxonomy review (Section 4.5) is mandatory. The review explicitly asks whether the categories are still adequate.
- **Registry as bureaucracy.** The classification process becomes so burdensome that the operator avoids producing new knowledge to avoid the classification work. Mitigation: the classification process should be lightweight -- a structured entry in a plain text file, not an elaborate production. If classification takes more than 30 minutes per item on average, the process should be simplified.
- **Deprecated knowledge zombies.** Knowledge that has been deprecated is still being used as the basis for operational decisions because operators check the procedure but not the registry. Mitigation: procedures that reference knowledge items should include the knowledge ID. When a knowledge item is deprecated, a search of all documents referencing it should be conducted and references updated.

## 7. Recovery Procedures

1. **If a classification backlog has accumulated:** Conduct a classification sprint. Classify all backlogged items at the minimum viable level: category, confidence level, one-sentence summary, and source reference. Relationships can be documented in a second pass. The priority is getting everything into the registry, even if the entries are sparse.

2. **If confidence levels have been inflated:** Conduct a confidence audit. Review all "verified" items. For each, check: is there a documented replication? Is there a documented adversarial review? If not, downgrade to "provisional." Review all "provisional" items. For each, check: is there a completed experiment log? If not, downgrade to "speculative."

3. **If relationships have been neglected:** Conduct a relationship review. For each registry entry, check: does this item depend on other items? Does it derive from evidence? Does it contradict anything? Does it supersede an older entry? Document all identified relationships. This is labor-intensive but is a one-time recovery that, once completed, is maintained incrementally.

4. **If the taxonomy has ossified:** Conduct an emergency taxonomy review. Identify all knowledge items that do not fit their assigned categories. Look for patterns -- do the misfits suggest a new category is needed? Revise the taxonomy and reclassify affected items. File the revision as a Tier 2 decision per R-KCT-06.

5. **If the registry has been abandoned entirely:** Begin rebuilding from the research archive. Every experiment log, every TDR, every negative result entry is a knowledge item that should be in the registry. Classify them systematically, starting with the most recent and working backward. Accept that the early reconstruction will be incomplete and refine it over time.

## 8. Evolution Path

- **Years 0-5:** The taxonomy is new and sparsely populated. The seven categories and four confidence levels should be sufficient for the institution's initial knowledge base. The operator is learning to classify consistently. Some miscategorization is inevitable and will be corrected in early audits.
- **Years 5-15:** The registry grows to hundreds of entries. Relationships between knowledge items become increasingly important. The registry begins to function as a genuine knowledge graph -- a map of what the institution knows and how its knowledge connects. The triennial taxonomy review may reveal that one or two categories need subdivision.
- **Years 15-30:** Succession is likely during this period. The taxonomy is one of the most valuable resources the successor inherits, because it provides an organized map of everything the institution has learned. The taxonomy review following succession is critically important -- the successor may organize knowledge differently than the founder, and the taxonomy should accommodate that.
- **Years 30-50+:** The registry contains the institution's entire intellectual history. Deprecated items show what was once believed and later disproven. Relationship chains show how understanding evolved. The registry is not just an index -- it is a narrative of institutional learning.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
I resist taxonomy. My natural instinct is to throw knowledge into a wiki and rely on search. But search finds what you are looking for. Taxonomy reveals what you did not know to look for. The relationship links, especially, are the value proposition: when I deprecate a piece of knowledge, I need to know what else depends on it. When I discover a contradiction, I need to know which items are affected. A flat file system cannot provide this. A taxonomy can. So I will build the taxonomy, I will maintain it, and I will resist the temptation to let it atrophy when it feels like overhead. It is not overhead. It is the institution's nervous system.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 6: Honest Accounting of Limitations)
- CON-001 -- The Founding Mandate (knowledge preservation, 50-year continuity)
- GOV-001 -- Authority Model (Tier 2 decisions for taxonomy changes)
- OPS-001 -- Operations Philosophy (documentation-first, operational tempo)
- D14-001 -- Research Philosophy (knowledge lifecycle, confidence levels, anti-echo-chamber structures)
- D14-002 -- Experimental Design (experiment log format, pre-registration)
- D14-005 -- Negative Results Registry (Category 6: Negative Knowledge)
- D14-006 -- Multi-Generational Research Continuity (taxonomy as succession resource)
- META-00-ART-001 -- Stage 1 Meta-Framework (50-year continuity rules, plain text requirements)

---

---

# D14-005 -- Negative Results Registry

**Document ID:** D14-005
**Domain:** 14 -- Research & Theory
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, D14-001, D14-002, D14-004
**Depended Upon By:** D14-003, D14-006. All Domain 14 articles that initiate new investigations. Referenced by all domains where prior failures are relevant to current work.

---

## 1. Purpose

This article defines the Negative Results Registry -- the institution's formal record of what has been tried and failed, what has been hypothesized and disproven, and what has been evaluated and rejected. It establishes why this registry exists, how it is structured, how entries are created, and how the registry is used to prevent the institution from repeating known failures.

The concept is simple and the resistance to it is powerful. Human beings do not naturally document failure. Failure feels like wasted effort, and documenting it feels like memorializing the waste. The instinct is to move on, to try the next thing, to focus on what works rather than what does not. This instinct, in a single-operator institution with no external peer community, is catastrophic over time. It guarantees that failures are repeated, that rejected approaches are revisited without memory of the prior rejection, and that the institution spends its finite research resources rediscovering things it already knows do not work.

D14-001 established the principle (R-RES-02) that negative results must be documented with the same rigor as positive results. This article operationalizes that principle into a specific system.

## 2. Scope

**In scope:**
- The philosophical and practical justification for maintaining a negative results registry.
- The negative results entry format: mandatory fields, optional fields, completion criteria.
- Types of negative results: experimental failures, disproven hypotheses, rejected technologies, abandoned approaches.
- The registry search protocol: how to consult the registry before starting new work.
- Integration with the knowledge taxonomy (D14-004) and the experiment log format (D14-002).
- The registry maintenance process: review, consolidation, and cross-referencing.

**Out of scope:**
- General experimental methodology (D14-002).
- General knowledge classification (D14-004).
- Technology evaluation procedures (D14-003, though rejected technology evaluations feed this registry).
- Operational incident reports (Domain 10, though operational failures may generate negative result entries).

## 3. Background

### 3.1 The Cost of Forgetting Failure

In the broader scientific community, the problem of unreported negative results -- sometimes called the "file drawer problem" -- is well recognized. Journals preferentially publish positive results. Researchers preferentially submit positive results. The result is a systematic distortion of the knowledge base: the published record suggests that interventions work more often than they actually do, because the failures are sitting in file drawers, unreported.

In a single-operator institution, the file drawer is the operator's memory, and memory is the least reliable file drawer of all. The operator runs an experiment. It fails. The operator moves on. Three years later, the operator -- or worse, a successor -- faces the same question, runs the same experiment, and discovers the same failure. The time, energy, and materials spent on the second attempt are pure waste, waste that a one-line entry in a registry would have prevented.

The cost of forgetting failure compounds. Each unreported negative result is a trap waiting for a future researcher. In an institution that operates for decades, the cumulative waste of repeated failures can dwarf the cost of maintaining the registry by orders of magnitude.

### 3.2 The Stigma of Negative Results

The most powerful obstacle to maintaining a negative results registry is not logistical -- it is psychological. Recording a failure feels like admitting defeat. The experiment log for a failed experiment is, the instinct whispers, a monument to wasted effort. This instinct must be confronted directly.

A failed experiment is not wasted effort. It is knowledge. It tells you that a particular approach, under particular conditions, does not produce the expected result. That knowledge has concrete value: it narrows the search space for future work, it prevents resource expenditure on known dead ends, and it sometimes reveals unexpected insights about why something does not work that illuminate how other things do work.

Furthermore, per ETH-001 Principle 6, the institution must honestly account for its limitations. A research archive that contains only successes is not honest. It presents a distorted picture of the institution's knowledge by omitting the attempts that did not succeed. The negative results registry is, in this sense, an ethical obligation, not just a practical tool.

### 3.3 The Predecessor Problem

Perhaps the most important justification for the negative results registry is the predecessor problem: the situation where a successor operator, unfamiliar with the institution's research history, embarks on a line of investigation that their predecessor already pursued and abandoned. Without a negative results registry, the successor has no way to know this. The predecessor's failures exist only in the predecessor's memory, which is no longer available.

The registry is, in this context, a letter from the past to the future: "I tried this. Here is what happened. Here is why I stopped. If you want to try it again, at least you know where I got stuck."

## 4. System Model

### 4.1 Types of Negative Results

The registry accommodates four types of negative results:

**Type 1: Experimental Failures.** An experiment was conducted per D14-002, and the hypothesis was disproven. The experiment produced a clear negative result. This is the cleanest type of negative result and the easiest to document, because the experiment log already contains most of the necessary information.

**Type 2: Inconclusive Investigations.** An investigation was conducted but produced no clear result -- the evidence was insufficient, the methodology was flawed, or external factors prevented completion. Inconclusive results are negative results in the sense that they did not achieve their objective, even if they did not disprove the hypothesis. They are recorded because they contain valuable information about methodology limitations and environmental constraints.

**Type 3: Technology Rejections.** A technology was evaluated per D14-003 and rejected. The Technology Decision Record (TDR) documents the evaluation. The registry entry cross-references the TDR and captures the specific reasons for rejection in a form that is quickly searchable.

**Type 4: Abandoned Approaches.** A line of inquiry or an operational approach was pursued for a period, found to be unproductive or impractical, and abandoned. This type is the most subjective and the hardest to document because the abandonment often happened gradually rather than at a clear decision point. It is also the type most prone to omission, and therefore the type where the registry adds the most value.

### 4.2 The Negative Results Entry Format

Each registry entry contains the following mandatory fields:

- **NR ID:** Format NR-YYYY-NNN (year and sequential number).
- **Date:** Date the entry was filed.
- **Type:** One of the four types defined in Section 4.1.
- **Title:** A concise, descriptive title. The title should be findable: a future researcher scanning the registry for relevant prior work should be able to identify relevance from the title alone.
- **Summary:** A one-paragraph summary of what was attempted, what was expected, and what actually happened.
- **Methodology reference:** Reference to the experiment log (D14-002), TDR (D14-003), or other documentation of the investigation.
- **Root cause analysis:** Why did this fail? What was the mechanism of failure? This is the most valuable field in the entry because it allows future researchers to determine whether the failure was inherent (the approach is fundamentally flawed) or conditional (the approach might work under different conditions).
- **Conditions and constraints:** Under what specific conditions was this attempted? Hardware, software versions, environmental conditions, and any other contextual factors. A failure under specific conditions does not necessarily mean failure under all conditions.
- **Residual value:** Did the failed investigation produce any useful knowledge, even though it did not achieve its primary objective? Unexpected observations, methodology insights, and incidental discoveries are documented here.
- **Recommendations for future work:** Should this approach be revisited? Under what changed conditions might it succeed? Or is the failure fundamental enough that the approach should be considered permanently closed? This field is the predecessor's advice to the future.
- **Knowledge registry cross-reference:** The D14-004 Knowledge ID for the corresponding negative knowledge entry.
- **Related entries:** Cross-references to other negative results entries on similar topics.

### 4.3 The Registry Search Protocol

Before beginning any new formal investigation (D14-002) or technology evaluation (D14-003), the researcher must search the negative results registry for relevant prior work. This is a mandatory step, not a suggested one. The search protocol is:

**Step 1: Keyword search.** Search the registry titles and summaries for terms related to the planned investigation.

**Step 2: Category search.** Search by knowledge category (D14-004) for entries in the same domain as the planned investigation.

**Step 3: Relationship search.** Check the knowledge registry (D14-004) for knowledge items related to the planned investigation and follow any links to negative results entries.

**Step 4: Document the search.** Record the search in the pre-registration document (D14-002, Section 4.1, Phase 2). List all relevant negative results entries found and explain how the planned investigation differs from or builds upon the prior work. If no relevant entries are found, record that the search was conducted and found nothing.

The search documentation serves two purposes: it ensures the search was actually done (not just claimed), and it provides the researcher with relevant context before the investigation begins.

### 4.4 Registry Maintenance

The registry requires periodic maintenance to remain useful:

**Annual consolidation.** Once per year, review the registry for entries that can be consolidated. If five separate experiments have all failed to achieve the same objective for the same reason, a single consolidated entry may be clearer than five individual ones. The individual entries are retained but linked to the consolidation.

**Cross-reference verification.** During the annual review, verify that all cross-references (to experiment logs, TDRs, knowledge registry entries, and other negative results entries) are still valid. Update broken references.

**Relevance review.** Some negative results become irrelevant over time -- the technology they evaluated no longer exists, the hardware they tested has been retired, the conditions they operated under no longer apply. These entries are not deleted (per D14-004, R-KCT-04, nothing is deleted) but may be marked as "historically relevant only" to help current researchers focus on entries that are still operationally pertinent.

## 5. Rules & Constraints

- **R-NRR-01:** Every formal experiment (D14-002) that fails to confirm its hypothesis must have a negative results entry filed within 30 days of the experiment's completion. This is a restatement and operationalization of D14-001, R-RES-02.
- **R-NRR-02:** Every technology evaluation (D14-003) that results in a rejection must have a negative results entry filed within 30 days of the TDR being finalized.
- **R-NRR-03:** The registry search protocol (Section 4.3) is mandatory before any new formal investigation or technology evaluation begins. The search must be documented in the pre-registration document.
- **R-NRR-04:** Negative results entries must not be deleted, modified retrospectively, or made less visible. They are append-only records. If subsequent work reveals that a negative result was incorrect (the approach actually works under conditions not tested), a new positive-result entry is created and the negative result entry is updated with a forward reference, but the original entry remains.
- **R-NRR-05:** The root cause analysis field (Section 4.2) is mandatory. "It did not work" is not a root cause. The entry must explain, to the extent determinable, why it did not work. If the root cause is unknown, the entry must state that explicitly and identify what further investigation would be needed to determine it.
- **R-NRR-06:** The registry must be stored in plain text format, searchable without specialized tools, and compliant with the 50-year continuity rules. The registry must be usable by a future operator with nothing more than a text editor and the human ability to read.
- **R-NRR-07:** Negative results must be entered in the knowledge registry (D14-004) under Category 6 (Negative Knowledge) with the appropriate confidence level. The negative results registry and the knowledge registry are complementary, not redundant -- the knowledge registry provides classification and relationships; the negative results registry provides detailed documentation.

## 6. Failure Modes

- **Registry neglect.** Negative results are not entered in the registry because the operator considers them unimportant or because the 30-day deadline passes without action. Over time, the registry becomes sparse and unreliable. Mitigation: R-NRR-01 and R-NRR-02 make filing mandatory. The annual classification audit (D14-004, R-KCT-05) cross-checks experiment logs against registry entries to identify missing entries.
- **Shallow root cause analysis.** Entries are filed but the root cause analysis is superficial -- "it did not work" or "performance was insufficient." These entries are nearly useless for preventing repeated failures because they do not explain why the failure occurred. Mitigation: R-NRR-05 requires substantive root cause analysis. The annual review should flag entries with inadequate root cause analysis for revision.
- **Search protocol bypass.** The researcher begins a new investigation without searching the registry, either because they consider it unnecessary or because they forgot. The investigation may repeat prior work. Mitigation: R-NRR-03 makes the search mandatory and requires documentation in the pre-registration. A pre-registration that does not include registry search documentation is incomplete.
- **Registry bloat.** Over decades, the registry becomes so large that searching it is impractical. Mitigation: the annual consolidation process (Section 4.4). The relevance review marks obsolete entries. The category structure of D14-004 provides filtering. If the registry exceeds a manageable size, an index by topic should be created.
- **False closure.** An approach is recorded as a negative result and future researchers avoid it, even though changed conditions (new hardware, new software, new techniques) might make it viable. Mitigation: the "Recommendations for future work" field (Section 4.2) explicitly addresses this. Entries should distinguish between fundamental failures and conditional failures. The conditions field records what was true when the failure occurred, allowing future researchers to assess whether those conditions still apply.
- **Psychological avoidance.** The operator simply does not want to document failures because it feels demoralizing. Mitigation: cultural reinforcement through the Commentary Section and through institutional norms established in D14-001. Negative results documentation is framed as an act of service to the future, not as self-flagellation.

## 7. Recovery Procedures

1. **If the registry has been neglected:** Conduct a registry reconstruction sprint. Review all experiment logs (D14-002) and Technology Decision Records (D14-003) from the neglected period. Identify all negative results that should have been entered. Create entries for each, using the available documentation. Accept that some entries will be incomplete due to faded memory or insufficient records -- incomplete entries are better than no entries.

2. **If root cause analyses are shallow:** Review all entries flagged as having insufficient root cause analysis. For recent entries, re-examine the evidence and attempt a deeper analysis. For older entries where re-analysis is not possible, annotate the entry with "Root cause analysis limited by available evidence as of [date]" and document what further investigation would be needed.

3. **If the search protocol has been routinely bypassed:** Conduct a retroactive search for all investigations initiated without registry searches. Check whether any of those investigations repeated prior failures. If so, document the duplication and the cost. Use the concrete examples to reinforce the value of the search protocol. Going forward, treat the search documentation as a gating requirement for pre-registration approval.

4. **If the registry has become bloated and unsearchable:** Create a topic index. Group entries by subject area, technology domain, and failure type. Mark entries whose conditions no longer apply as "historically relevant only." Consider creating a "top 20" list of the most important negative results that every researcher should be aware of.

5. **If negative results documentation has been psychologically avoided:** Acknowledge the avoidance honestly in the Commentary Section. Commit to entering the next three negative results within one week of occurrence. Build momentum. The hardest part is starting.

## 8. Evolution Path

- **Years 0-5:** The registry is small. Most entries are from technology evaluations and initial system configuration experiments. The registry search protocol may feel unnecessary because the operator remembers all the failures. Follow it anyway -- the habit matters more than the current utility.
- **Years 5-15:** The registry begins to earn its keep. The operator encounters a question they investigated before but had forgotten. The registry reminds them. This is the moment the registry transitions from a discipline to a tool. The accumulation of entries starts revealing patterns -- certain categories of approach that consistently fail, certain conditions that consistently cause problems.
- **Years 15-30:** If succession occurs, the registry is one of the most valuable resources the predecessor can leave. The successor, facing an institution they did not build, needs to know not just what works but what has been tried and does not work. Without the registry, the successor will spend years rediscovering the predecessor's failures.
- **Years 30-50+:** The registry is a multi-decade record of institutional learning from failure. It contains wisdom that no single person remembers -- the accumulated knowledge of what does not work, why, and under what conditions. It is, in a sense, the institution's immune memory: the record of every pathogen it has encountered and defeated.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
I am going to fail at things. I am going to try configurations that do not work, evaluate technologies that turn out to be unsuitable, run experiments that disprove my hypotheses, and pursue lines of inquiry that lead nowhere. This is not pessimism. It is realism. The question is not whether I will fail but whether my failures will be wasted.

If I record them -- honestly, with root cause analysis, with context, with recommendations -- then my failures are converted from waste into knowledge. They become the institution's scar tissue: not pretty, but tougher than what was there before.

The hardest entries to write will be the ones about things I really wanted to work. The backup strategy I was proud of but that failed under load. The technology I evangelized but that turned out to be unsuitable. The approach I invested weeks in before discovering it was fundamentally flawed. Those entries will sting. I will write them anyway. The person who needs to read them -- whether that is future-me or a successor -- deserves the truth, not a curated highlight reel.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 2: Integrity Over Convenience; Principle 6: Honest Accounting of Limitations)
- CON-001 -- The Founding Mandate (knowledge preservation, institutional continuity)
- D14-001 -- Research Philosophy (R-RES-02: negative result documentation, Section 4.2: negative result preservation)
- D14-002 -- Experimental Design for Solo Researchers (experiment log format, pre-registration)
- D14-003 -- Technology Evaluation Framework (Technology Decision Records, rejected technology documentation)
- D14-004 -- Knowledge Classification and Taxonomy (Category 6: Negative Knowledge, confidence levels)
- D14-006 -- Multi-Generational Research Continuity (negative results as succession resource)
- META-00-ART-001 -- Stage 1 Meta-Framework (50-year continuity rules)

---

---

# D14-006 -- Multi-Generational Research Continuity

**Document ID:** D14-006
**Domain:** 14 -- Research & Theory
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, OPS-001, D14-001, D14-002, D14-003, D14-004, D14-005, D13-001, D11-001
**Depended Upon By:** All Domain 14 articles involving long-term research programs. Referenced by Domain 13 (Evolution) for continuity of research-informed evolution decisions. Referenced by Domain 9 (Documentation & Knowledge) for knowledge transfer standards.

---

## 1. Purpose

This article addresses the most difficult challenge in the research domain: how to maintain research programs across the transition from one operator to the next. Everything else in Domain 14 -- experimental design, technology evaluation, knowledge classification, negative results -- exists within the working life of a single operator. This article addresses the boundary between operators, the moment when one person's research practice must become another person's inheritance.

The challenge is not merely logistical. It is not enough to hand over a well-organized archive (though that is necessary). The challenge is epistemological: how does a successor operator understand not just what was known, but why it was investigated, what questions were left open, what approaches were abandoned and why, and what the predecessor's unwritten intuitions and suspicions were? How does the successor distinguish between knowledge that the predecessor was confident about and knowledge that the predecessor accepted provisionally because they ran out of time to investigate further?

CON-001 declares that this institution must be transferable through documentation alone. D14-006 is the research domain's answer to that declaration. It defines the systems that make research continuity possible across the most disruptive event a single-operator institution can face: the replacement of the single operator.

## 2. Scope

**In scope:**
- The research handoff protocol: what must be documented, when, and in what format.
- Research state documentation: how to capture the current state of all active and suspended research programs.
- Open questions registry: how to document questions that remain unanswered.
- Abandoned leads documentation: how to record lines of inquiry that were started but not completed, and why they were abandoned.
- The research continuity audit: a periodic verification that the research state is documented well enough for handoff.
- Successor orientation: how a new operator gets up to speed on the research program.
- Long-term research programs: how to design and maintain investigations that span years or decades.

**Out of scope:**
- General succession planning (Domain 9, Domain 13, GOV-001).
- Operational handoff procedures (OPS-001, Domain 12 for emergency succession).
- Physical asset transfer (D11-001, ADM-series).
- Governance transition (GOV-001).

## 3. Background

### 3.1 The Successor's Dilemma

Imagine a successor operator arriving at this institution for the first time. They have access to all the documentation. They can read the charter, the policies, the operational procedures. They can follow the maintenance checklists and the backup schedules. But when they encounter the research domain, they face a different kind of challenge.

The predecessor was in the middle of evaluating a new storage technology. The evaluation was 70% complete. The preliminary results were promising but inconclusive. The predecessor had a hunch -- based on years of experience with similar technologies -- that a specific failure mode might emerge under sustained load, but they had not yet designed the test. None of this was in the TDR because the evaluation was not finished.

The successor now has three options, all bad. They can complete the evaluation without the predecessor's context, potentially missing the suspected failure mode. They can abandon the evaluation and start fresh, wasting all prior work. Or they can try to reconstruct the predecessor's thinking from incomplete records, which is possible only if the predecessor documented their in-progress research state.

This article ensures that the third option is available and reliable.

### 3.2 Research State as Institutional Infrastructure

Most institutions treat research state as ephemeral -- it exists in the researcher's head, in their lab notebook, in their half-finished draft papers. When the researcher leaves, the state evaporates. This is acceptable in a research university where multiple researchers work in overlapping areas and institutional knowledge is distributed. It is not acceptable in a single-operator institution where all research knowledge is concentrated in one mind.

This article treats research state as institutional infrastructure, subject to the same documentation and preservation requirements as any other critical institutional asset. The state of every active research program, every open question, and every abandoned lead must be documented in a form that survives the operator's departure.

### 3.3 The Living Research State Document

The key mechanism of research continuity is the Living Research State Document (LRSD) -- a continuously maintained document that captures the current state of all research activity. The LRSD is not a summary written at the end of a research program. It is updated as research progresses, capturing the current state of understanding, the current hypotheses, the current plan, and -- critically -- the current uncertainties and intuitions.

The LRSD is the research equivalent of a ship's log. It records not just what has been decided and done, but what is being considered, what seems promising, what seems suspicious, and what the researcher plans to do next. If the researcher is replaced tomorrow, the LRSD tells the successor where they are in the journey.

## 4. System Model

### 4.1 The Living Research State Document (LRSD)

The LRSD is a plain-text document, updated at least monthly, containing the following sections:

**Section 1: Active Research Programs.** For each active research program or investigation:
- Program ID and title.
- Start date and expected duration.
- Current phase (formulation, pre-registration, execution, analysis, classification).
- Current status summary: what has been done, what has been found so far.
- Current hypothesis and how it has evolved since the program began.
- Current plan: what the next steps are, in what order.
- Current concerns: what the researcher is worried about, what might go wrong, what does not feel right. This is the section where intuition and informal suspicion are recorded. It is explicitly permitted -- and encouraged -- to record hunches that are not yet supported by evidence, as long as they are labeled as such.
- Dependencies: what other research or operational activities this program depends on.
- Resources: what equipment, time, and materials the program requires.

**Section 2: Open Questions.** A registry of questions that the institution has identified as important but has not yet investigated, or has investigated inconclusively:
- Question ID and statement.
- Date identified.
- Why it matters: what institutional decision or understanding depends on the answer.
- Current state of knowledge: what is known, what is suspected, what is unknown.
- Barriers to investigation: why the question has not been answered yet (lack of time, lack of equipment, prerequisite knowledge missing, etc.).
- Priority: high (answer needed within 1 year), medium (answer needed within 3 years), low (answer would be valuable but is not urgent).

**Section 3: Abandoned Leads.** A registry of lines of inquiry that were started but discontinued:
- Lead ID and description.
- Date started and date abandoned.
- Reason for abandonment: why the lead was discontinued. The reasons are classified as: resource constraint (could not afford the time or materials), dead end (investigation revealed the question was not meaningful or answerable), superseded (a better approach was found), or deferred (still potentially valuable but deprioritized).
- What was learned before abandonment: even incomplete investigations typically produce some knowledge. Record it.
- Conditions for revival: under what circumstances should this lead be revisited? Is there a specific trigger (new technology, new equipment, more time) that would make the investigation viable?
- Cross-reference to negative results entry (D14-005) if applicable.

**Section 4: Research Environment.** A snapshot of the current research infrastructure:
- Available test equipment and its condition.
- Current software tools used for research.
- Current limitations and known gaps in research capability.
- Planned improvements to research infrastructure.

**Section 5: Researcher's Notes.** An unstructured section where the operator records observations, speculations, patterns they have noticed, and connections they suspect but have not validated. This section is explicitly informal. Its purpose is to capture the kind of knowledge that is usually lost when a researcher leaves -- the tacit understanding that is too preliminary for formal documentation but too valuable to discard.

### 4.2 The Research Handoff Protocol

When operator succession occurs, the research handoff follows a defined protocol:

**Phase 1: LRSD Finalization (pre-succession).** The departing operator updates the LRSD to reflect the current state of all research activity as of the handoff date. Each active program receives a "handoff summary" addendum that provides the departing operator's honest assessment of the program's viability, the quality of the evidence so far, and their personal recommendation for whether the successor should continue, modify, or abandon the program.

**Phase 2: Archive Verification (pre-succession).** The departing operator verifies that all research artifacts -- experiment logs, TDRs, negative results entries, knowledge registry entries, raw data -- are present, organized, and findable. Any gaps in the archive are documented. This verification is conducted using the research continuity audit checklist (Section 4.3).

**Phase 3: Orientation (post-succession).** The successor reads the LRSD in its entirety. They may -- and should -- annotate it with questions, disagreements, and alternative perspectives. The LRSD is designed to be the starting point for the successor's understanding, not a set of orders. The successor is free to continue, modify, or abandon any research program based on their own judgment, but they must document their reasoning for departures from the predecessor's plan.

**Phase 4: Confirmation (post-succession).** Within 90 days of assuming responsibility, the successor files a Research Continuity Confirmation in the research archive. This document records: which active programs they intend to continue, which they intend to modify (and how), which they intend to suspend or abandon (and why), which open questions they consider priorities, and what new research directions they plan to pursue. This confirmation is not a commitment that binds the successor permanently -- it is a snapshot of their initial assessment, subject to revision as they gain experience.

### 4.3 The Research Continuity Audit

The research continuity audit is a periodic check, conducted annually, that verifies the institution's research state is documented well enough for handoff at any time. The audit is not conducted in anticipation of succession. It is conducted routinely, because the best time to prepare for succession is before you know it is coming.

The audit checklist:

1. Is the LRSD current? Has it been updated within the last 30 days?
2. Are all active research programs documented in the LRSD with current status?
3. Does each active program have a current plan with next steps?
4. Are all experiment logs (D14-002) complete and filed in the research archive?
5. Are all TDRs (D14-003) complete and filed?
6. Is the negative results registry (D14-005) current? Are there experiments or evaluations that produced negative results but lack registry entries?
7. Is the knowledge registry (D14-004) current? Are there classified knowledge items that have not been entered?
8. Are all open questions documented in the LRSD?
9. Are all abandoned leads documented in the LRSD with reasons and conditions for revival?
10. Is the Researcher's Notes section current? Has the operator recorded their current thinking, hunches, and informal observations?
11. Could a competent successor, reading only the documentation, understand the state of every research program well enough to make informed decisions about continuation?

The last question is the critical one. It is subjective, and the operator will be biased toward answering "yes." The mitigation is documentation replication (D14-001, Section 4.3): the operator selects one active research program and attempts to understand its state from the documentation alone, without relying on memory. If the documentation is insufficient, it is improved.

### 4.4 Long-Term Research Programs

Some investigations span years or decades. Monitoring hardware degradation over time. Tracking the long-term stability of a storage format. Assessing whether a particular technology's community remains healthy over a full product generation. These long-term programs require specific continuity practices:

**Milestone documentation.** Long-term programs define milestones -- points at which the accumulated data is analyzed and interim conclusions are drawn. Each milestone produces a formal report that is filed in the research archive and summarized in the LRSD. If the operator changes between milestones, the milestone reports provide a structured history of the program's progress.

**Data continuity.** The raw data from long-term programs must be stored in formats that will remain readable for the duration of the program. Per the 50-year continuity rules, this means plain text, CSV, or other simple structured formats. Data stored in proprietary formats that may become unreadable is a continuity risk.

**Hypothesis evolution tracking.** Over a long-term program, the initial hypothesis may evolve substantially as data accumulates. The LRSD tracks this evolution, recording not just the current hypothesis but the history of how it changed and what evidence prompted each change. This evolution record is essential for a successor who needs to understand not just what is currently believed but how the understanding developed.

**Abandonment criteria.** Long-term programs define, at inception, the conditions under which they should be abandoned. Without explicit abandonment criteria, long-term programs tend to persist indefinitely on institutional inertia, consuming resources without producing value. The criteria should be revisited at each milestone.

## 5. Rules & Constraints

- **R-MRC-01:** The Living Research State Document must be updated at least monthly. An LRSD that has not been updated in more than 60 days triggers an automatic audit requirement.
- **R-MRC-02:** The research continuity audit (Section 4.3) must be conducted annually, with results documented in the research archive. The audit may not be deferred for more than one quarter.
- **R-MRC-03:** Upon succession, the research handoff protocol (Section 4.2) is mandatory. The successor must file a Research Continuity Confirmation within 90 days of assuming responsibility.
- **R-MRC-04:** All research artifacts (experiment logs, TDRs, negative results entries, knowledge registry entries, raw data, LRSD versions) must be preserved in the research archive in formats compliant with the 50-year continuity rules.
- **R-MRC-05:** Long-term research programs (expected duration exceeding two years) must define milestones, data continuity requirements, hypothesis evolution tracking, and abandonment criteria at inception. These definitions are filed in the research archive and referenced in the LRSD.
- **R-MRC-06:** The Researcher's Notes section of the LRSD is explicitly permitted to contain informal, speculative, and subjective content. The operator must not self-censor this section out of concern that informal observations are not "rigorous enough." The purpose of this section is to capture exactly the kind of knowledge that is lost when only formal results are documented.
- **R-MRC-07:** Departing operators must complete LRSD Finalization and Archive Verification (Phases 1 and 2 of the handoff protocol) before succession takes effect. If succession is unplanned (emergency, incapacitation), the successor conducts the archive verification and documents any gaps discovered.

## 6. Failure Modes

- **LRSD staleness.** The LRSD is not updated regularly and drifts out of sync with actual research activity. If succession occurs during a stale period, the successor inherits an inaccurate picture of the research state. Mitigation: R-MRC-01's monthly update requirement. The annual audit (R-MRC-02) checks LRSD currency.
- **Formality suppression.** The operator keeps the formal sections of the LRSD current but neglects the Researcher's Notes section, losing the informal knowledge that is often the most valuable for a successor. Mitigation: R-MRC-06 explicitly legitimizes informal content. The audit checklist (Section 4.3, item 10) specifically checks for current notes.
- **Handoff document overload.** The LRSD becomes so long and detailed that a successor cannot absorb it in a reasonable time. Mitigation: the LRSD should be structured for progressive reading. The active programs section provides the immediate context. The open questions and abandoned leads provide deeper history. The researcher's notes provide the most intimate context. A successor can start with the first section and deepen their reading as needed.
- **Successor divergence without documentation.** The successor departs from the predecessor's research directions without documenting their reasoning. Over time, the institution loses track of why certain programs were continued and others were not. Mitigation: R-MRC-03 requires a Research Continuity Confirmation that documents the successor's decisions.
- **Long-term program inertia.** A long-term research program persists past its usefulness because no one reviews the abandonment criteria. Mitigation: R-MRC-05 requires abandonment criteria at inception. Milestone reviews include an explicit check against those criteria.
- **Emergency succession gap.** The operator is incapacitated or departs unexpectedly without completing the handoff protocol. Mitigation: the continuous maintenance of the LRSD (updated monthly, audited annually) means that even without a formal handoff, the successor has a reasonably current snapshot of the research state. The gap between the last LRSD update and the succession event is at most 60 days (per R-MRC-01), and the content of that gap can often be partially reconstructed from operational logs and experiment logs.

## 7. Recovery Procedures

1. **If the LRSD has gone stale:** Conduct an immediate update. For each research program, determine its current state from memory, experiment logs, and operational logs. Update the LRSD. File a note in the Commentary Section acknowledging the period of staleness and any information that may have been lost during that period.

2. **If research continuity audits have been skipped:** Conduct the audit immediately using the checklist in Section 4.3. Accept that some items may reveal significant gaps. Prioritize closing the gaps based on the severity of the consequence if succession occurred today.

3. **If a succession has occurred without a formal handoff:** The successor conducts an archive inventory, identifying all available research artifacts. They construct a best-effort LRSD from the available materials. Gaps are documented explicitly. The successor files a Research Continuity Confirmation that acknowledges the limitations of the informal handoff and identifies areas where the predecessor's knowledge is irrecoverably lost.

4. **If long-term research programs have persisted without milestone reviews:** Conduct milestone reviews for all long-term programs immediately. For each, assess: is the program still producing value? Has the original question been answered (even partially)? Do the abandonment criteria suggest the program should be closed? Document the review findings and make continuation decisions explicitly.

5. **If the Researcher's Notes section has been neglected:** Set aside one hour. Write down everything you currently think, suspect, wonder about, and worry about regarding the institution's research programs. Do not edit for formality. Do not censor speculation. The goal is to externalize the tacit knowledge that has been accumulating in your head. Then commit to adding to this section weekly, even if the additions are brief.

## 8. Evolution Path

- **Years 0-5:** The LRSD is new and the operator is the founder. The handoff protocol exists in documentation but has not been tested. The value of the LRSD in this period is primarily as a discipline -- the monthly update forces the operator to reflect on research state and direction. The Researcher's Notes section, if maintained honestly, will become the most valuable section in a decade.
- **Years 5-15:** The LRSD and associated research artifacts accumulate substance. If a succession occurs in this period, the handoff protocol is tested for the first time. The first succession will almost certainly reveal gaps in the documentation -- things the predecessor took for granted that the successor needs explained. These gaps should be documented and used to improve the LRSD template and the handoff protocol for future successions.
- **Years 15-30:** The institution may have experienced one or more successions. The LRSD contains layers of research state from different operators, each with their own perspectives, priorities, and notes. The Researcher's Notes section, spanning multiple operators, becomes a conversation across time -- each operator reading the notes of their predecessors and adding their own.
- **Years 30-50+:** The LRSD, experiment archive, knowledge registry, and negative results registry together form a multi-decade record of the institution's intellectual life. This record is one of the institution's most valuable assets. It contains not just what was known but how it was discovered, not just what worked but what was tried, not just what was decided but what was wondered. The continuity mechanisms defined in this article are what make this record possible.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
I am writing a succession document for a position that only I hold, for a successor who may not exist for years. This feels premature. It is not.

The truth is that I do not know when succession will happen. It could be planned, decades from now, when I choose to hand off the institution. It could be unplanned, next year, through accident or illness. If I wait until succession is imminent to prepare the research handoff, I will either rush the documentation or leave it undone.

The LRSD is not just a succession tool. It is a thinking tool. The monthly update forces me to step back from the details of individual experiments and ask: where is the research program going? What do I think I know? What am I unsure about? What have I been meaning to investigate but have not gotten to? These are valuable questions regardless of whether a successor ever reads my answers.

The section I am most committed to maintaining is the Researcher's Notes. I know from experience in other contexts that the most valuable knowledge a departing person has is the stuff that never makes it into formal documentation -- the hunches, the suspicions, the patterns they noticed but never proved, the approaches they discarded for reasons they never wrote down. I want to capture as much of that as possible, not because it is rigorous (it is not) but because it is real, and because the person who comes after me will need every advantage I can give them.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 6: Honest Accounting of Limitations)
- CON-001 -- The Founding Mandate (transferability through documentation, lifetime operation)
- GOV-001 -- Authority Model (succession procedures, decision documentation)
- OPS-001 -- Operations Philosophy (documentation-first, operational logging)
- D14-001 -- Research Philosophy (knowledge lifecycle, replication, anti-echo-chamber structures)
- D14-002 -- Experimental Design for Solo Researchers (experiment log format)
- D14-003 -- Technology Evaluation Framework (Technology Decision Records)
- D14-004 -- Knowledge Classification and Taxonomy (knowledge registry, confidence levels)
- D14-005 -- Negative Results Registry (negative result documentation)
- D13-001 -- Evolution Philosophy (research as input to evolution, change management)
- D11-001 -- Administration Philosophy (resource management, scheduling)
- D9-001 -- Documentation Philosophy (documentation standards, competence levels)
- META-00-ART-001 -- Stage 1 Meta-Framework (50-year continuity rules)

---

*End of Stage 4: Research & Theory Advanced Articles*