# STAGE 3: OPERATIONAL DOCTRINE (BATCH 3)

## Content Import, Documentation Quality, Onboarding, Decision Logs, and Infrastructure

**Document ID:** STAGE3-OPS-BATCH3
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Classification:** Operational Procedures -- These articles translate Stage 2 philosophy into step-by-step actionable procedures. They are designed to be followed by a single operator with no one to ask.

---

## How to Read This Document

This document contains five operational articles that belong to Stage 3 of the holm.chat Documentation Institution. Where Stage 2 established what we believe and why, Stage 3 tells you what to do about it. These are manuals for the hands, not treatises for the mind.

Each article in this batch operationalizes a different domain philosophy: D18-001 becomes a content import pipeline you can execute step by step. D19-001 becomes audit checklists and scoring rubrics you can apply to real documents. D9-001 becomes a ninety-day onboarding curriculum a new maintainer can follow without asking questions. D20-001 becomes a daily practice for maintaining the decision log. D4-001 becomes a physical site specification you can build from.

If you have not read the Stage 2 philosophy articles for these domains, stop and read them first. These procedures implement principles, and principles misunderstood produce procedures misapplied. The philosophy tells you why. This document tells you how. You need both.

If something in these procedures does not work -- because hardware has changed, because a tool referenced here no longer exists, because the physical environment differs from what was assumed -- do not abandon the procedure. Adapt it. The principles behind each procedure are stated in the Background section. Use those principles to find the equivalent procedure for your reality. Then update this document. A procedure that no longer works but has not been updated is a trap for the next person.

---

---

# D18-002 -- Content Import Pipeline

**Document ID:** D18-002
**Domain:** 18 -- Import & Quarantine
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, SEC-001, OPS-001, D18-001
**Depended Upon By:** D19-002 (quality audit of imported content), D20-002 (provenance records in decision log), D17-002 (federation sync package handling).

---

## 1. Purpose

This article defines the complete, step-by-step procedure for importing external content into the holm.chat Documentation Institution. It takes the philosophy established in D18-001 -- the immune system analogy, the four-stage pipeline, the trust scoring model -- and turns it into a procedure that the operator can execute every time a piece of external media arrives at the institution's boundary.

This procedure is designed to be followed exactly. Not because deviation is forbidden, but because every step exists for a documented reason, and skipping a step means accepting a risk that the pipeline was designed to mitigate. If you believe a step is unnecessary for a particular import, you may still have good reason to skip it -- but you must document that decision and the reasoning behind it in the import log.

This article is written for the operator who is standing at their workstation with a USB drive in hand, wondering: "What do I do with this?" The answer is: follow this procedure.

## 2. Scope

This article covers:

- Physical media handling: inspection, labeling, and chain-of-custody documentation.
- Quarantine environment setup and operation.
- Malware scanning procedures for an air-gapped environment.
- Format validation checks and conversion requirements.
- Provenance documentation: recording where content came from and how it arrived.
- Trust scoring worksheets: the practical application of D18-001's trust model.
- The approval workflow: who decides, how, and what gets recorded.
- Post-import verification and the transition from quarantine to production.

This article does not cover:

- The philosophy of import and quarantine (see D18-001).
- Federation sync package processing (see D17-002, which extends this pipeline with federation-specific checks).
- Format conversion specifications for specific file types (see D18-ART-005).
- The quarantine system's own maintenance and rebuild procedures (see D18-ART-003).

## 3. Background

D18-001 establishes that the import process is the institution's immune system -- it must be sophisticated enough to admit what is beneficial and reject what is harmful. This article is the immune system's operating manual.

The critical constraint is the air gap. In a connected environment, malware scanning relies on regularly updated signature databases, real-time threat intelligence feeds, and cloud-based analysis. In an air-gapped environment, the operator has none of these. The scanning tools have signatures that were current the last time they were manually updated. The threat intelligence is whatever the operator has imported through this same pipeline. This means the import procedure must compensate for limited scanning capability with additional layers: format restriction, behavioral analysis, structural validation, and human judgment.

The second critical constraint is that every import is a conscious, deliberate act. Nothing enters the institution passively. This is both a strength (no drive-by downloads, no background synchronization) and a responsibility (every piece of content that enters the institution entered because a human being carried it across the air gap and executed this procedure). The import log is therefore not just a technical record. It is a record of deliberate human decisions, each of which accepted responsibility for what was admitted.

## 4. System Model

### 4.1 The Physical Import Station

The import station is a dedicated workspace where all import processing occurs. It consists of:

**The quarantine machine.** A physically separate computer that is never connected to any institutional production system. It has its own storage, its own operating system, and its own set of scanning and analysis tools. It shares no network, no cable, no wireless connection with any production system. The quarantine machine is expendable -- if it is compromised by malicious content, it can be rebuilt from a known-good image without affecting the institution.

**The transfer medium.** A set of dedicated USB drives or external storage devices used exclusively for transferring approved content from the quarantine machine to the production systems. These transfer media are never used for any other purpose. They are never connected to non-institutional systems. They are stored in a designated location when not in use.

**The physical workspace.** A clean, well-lit area where incoming media can be inspected before it is connected to anything. This workspace includes: a logbook (physical or digital, but on the quarantine machine) for recording import events, the trust scoring worksheet templates, and the provenance documentation forms.

### 4.2 The Pipeline Stages

The import pipeline has seven stages, each building on the previous:

**Stage 1: Physical Reception and Inspection.**
**Stage 2: Provenance Documentation.**
**Stage 3: Quarantine Mount and Inventory.**
**Stage 4: Malware Scanning.**
**Stage 5: Format Validation.**
**Stage 6: Trust Scoring.**
**Stage 7: Approval, Transfer, and Post-Import Verification.**

Each stage produces a documented output. If any stage fails, the pipeline halts and the content does not advance. The operator may choose to reject the content outright or to address the failure and re-enter the pipeline at the failed stage.

## 5. Rules & Constraints

### 5.1 Stage 1: Physical Reception and Inspection

**Step 1.1:** Receive the incoming medium. Do not connect it to any system yet. Place it on the physical workspace.

**Step 1.2:** Inspect the medium visually. Check for:
- Physical damage (cracked casing, bent connector, water damage, corrosion).
- Signs of tampering (broken seals, mismatched labels, evidence of re-soldering on circuit boards for external drives).
- Labeling (is the medium labeled? Does the label match what you expected to receive?).

**Step 1.3:** Record the physical inspection in the import logbook:
- Date and time of receipt.
- Description of the medium (type, brand, capacity, serial number if visible).
- Physical condition assessment (good, damaged, suspect).
- Source: who or what provided this medium.
- Expected contents: what you believe is on it, if known.

**Step 1.4:** If the medium is physically damaged or shows signs of tampering, stop. Document the finding. Decide whether to proceed with elevated caution or to reject the medium outright. Record the decision and reasoning.

### 5.2 Stage 2: Provenance Documentation

Before the medium touches any system, document its provenance. Fill out the Provenance Record:

- **Import ID:** Sequential identifier (format: IMP-YYYY-NNN, where YYYY is the year and NNN is the sequential number).
- **Date of import initiation:** Today's date.
- **Medium description:** From Step 1.3.
- **Source of the content:** The original creator or publisher of the content. Be as specific as possible. "Downloaded from the official PostgreSQL website" is good. "From the internet" is not.
- **Source URL or reference:** The specific URL, ISBN, DOI, catalog number, or other identifier for the content at its source.
- **Carrier description:** How the content traveled from its source to the institution. "Downloaded by the operator to a personal laptop on 2026-02-10, transferred to this USB drive on 2026-02-14, carried to the institution on 2026-02-16" is a good carrier description. "Found this on a drive" is not.
- **Chain of custody:** List every person and system that handled this content between its source and this institution. Include dates.
- **Operator identity:** Your name or identifier.
- **Purpose of import:** Why are you bringing this content into the institution? What need does it serve? Reference the institutional mission (CON-001) if applicable.

### 5.3 Stage 3: Quarantine Mount and Inventory

**Step 3.1:** Power on the quarantine machine. Verify it is in a known-good state (boot checksum matches the expected value, or rebuild from the known-good image if uncertain).

**Step 3.2:** Connect the incoming medium to the quarantine machine. Mount it read-only. This is critical -- read-only mounting prevents any write operation that could be triggered by autorun mechanisms or filesystem-level exploits.

Command example (Linux): `mount -o ro,noexec,nosuid,nodev /dev/sdX1 /mnt/quarantine_intake`

The `noexec` flag prevents execution of binaries from the mounted medium. The `nosuid` flag prevents setuid escalation. The `nodev` flag prevents device file interpretation.

**Step 3.3:** Generate a complete inventory of the medium's contents:
- List all files with their paths, sizes, and modification dates.
- Calculate checksums (SHA-256 minimum) for every file.
- Record the total number of files and total data volume.
- Note any unexpected content (files you did not expect, hidden files, system files, executables).

**Step 3.4:** Compare the inventory against the expected contents documented in the provenance record. If there are significant discrepancies, document them. Decide whether to proceed or reject.

**Step 3.5:** Copy all files from the incoming medium to the quarantine machine's local storage (into the quarantine workspace, not into any production area). Verify the copy by comparing checksums. Once verified, unmount and physically disconnect the incoming medium. Set it aside. All further work happens on the local copy.

### 5.4 Stage 4: Malware Scanning

**Step 4.1:** Run the institution's locally installed antivirus scanner against all files in the quarantine workspace. Note: the scanner's signature database is only as current as the last manual update. Record the scanner name, version, and signature date in the import log.

**Step 4.2:** Run a second scanning tool if available (defense in depth). Different scanners catch different threats. If only one scanner is available, document this limitation.

**Step 4.3:** For any executable files, scripts, or active content (macros, embedded scripts in documents), perform additional analysis:
- Examine scripts and executables with a hex viewer or disassembler if the operator has the skill.
- Check file headers to verify that the file type matches its extension (a .pdf that is actually an executable is a red flag).
- For documents with macros (Office files, LibreOffice files), extract and review the macro code.

**Step 4.4:** Record all scan results in the import log:
- Scanner(s) used and their versions.
- Number of files scanned.
- Any detections (malware, suspicious content, potentially unwanted programs).
- Disposition of detections (quarantined, deleted, false positive determination).

**Step 4.5:** If malware is detected, stop. Do not transfer any content from this import batch to production. Document the finding. Assess whether the malware was targeted or incidental. Rebuild the quarantine machine from its known-good image. The incoming medium should be physically marked as compromised and stored separately or destroyed. Record the incident in the security log per SEC-001.

### 5.5 Stage 5: Format Validation

**Step 5.1:** For each file, verify that its format is acceptable for institutional storage. The institution maintains a list of accepted formats (defined in D18-ART-005). The general hierarchy of trust, per D18-001:

- **Tier A (Preferred):** Plain text (UTF-8), PDF/A, PNG, JPEG, SVG, CSV, JSON, XML with published schema, Markdown, ODF (OpenDocument Format).
- **Tier B (Accepted with notation):** Standard PDF, HTML, EPUB, FLAC, WAV, TIFF, SQLite databases.
- **Tier C (Accepted with conversion required):** Proprietary document formats (DOCX, XLSX, PPTX) -- must be converted to Tier A or B equivalents before admission to production. Originals may be retained in an archive area.
- **Tier D (Rejected by default):** Executable binaries, installer packages, compressed archives containing executables, DRM-protected files, formats with no open specification. Tier D content may be admitted only with a Tier 3 governance decision (per GOV-001) documenting the justification.

**Step 5.2:** For each file, validate the internal structure:
- Open the file in an appropriate viewer on the quarantine machine. Does it render correctly?
- For structured formats (XML, JSON, CSV), validate against the schema or check for well-formedness.
- For documents, check that the content is readable and complete (not truncated, not corrupted).

**Step 5.3:** For Tier C content, perform format conversion on the quarantine machine. Verify the conversion by comparing the output against the original (visual inspection for documents, data comparison for spreadsheets). Record the conversion tool, version, and any conversion issues.

**Step 5.4:** Record the format validation results in the import log:
- For each file: format identified, tier classification, validation result (pass/fail), conversion performed (if any).

### 5.6 Stage 6: Trust Scoring

Complete the Trust Scoring Worksheet for this import. The worksheet operationalizes D18-001's four-dimensional trust model:

**Source Trust (1-5 scale):**
- 5: Known, verified institutional or academic source with an established track record.
- 4: Known source with generally reliable reputation but not personally verified.
- 3: Source is identifiable but reputation is unknown or mixed.
- 2: Source is vaguely identified ("from a forum," "shared by a colleague").
- 1: Source is unknown or anonymous.

**Carrier Trust (1-5 scale):**
- 5: Content acquired directly from the original source by the operator, with verified integrity (checksums published by the source match).
- 4: Content acquired from a reputable mirror or secondary source, with chain of custody fully documented.
- 3: Content passed through one or two intermediaries with partial documentation.
- 2: Content passed through multiple intermediaries or chain of custody is poorly documented.
- 1: Chain of custody is unknown or the content was found on a medium of uncertain origin.

**Format Trust (1-5 scale):**
- 5: Tier A format. Open specification, well-understood, no active content capability.
- 4: Tier B format. Open or well-documented, minimal active content risk.
- 3: Tier C format before conversion (proprietary but convertible).
- 2: Complex format with active content capability (macros, scripts) that has been reviewed.
- 1: Opaque, proprietary, or executable format.

**Content Trust (1-5 scale):**
- 5: Content is internally consistent, verifiable against independent sources, and well-structured.
- 4: Content appears consistent and well-structured but has not been independently verified.
- 3: Content is plausible but contains some inconsistencies or gaps.
- 2: Content has significant quality concerns (poor structure, unverifiable claims, obvious errors).
- 1: Content is unverifiable, contradictory, or of dubious quality.

**Aggregate Trust Score:** Sum the four dimensions. Maximum 20, minimum 4.

- **16-20:** High trust. Standard validation is sufficient. Proceed to approval.
- **11-15:** Moderate trust. Review the weakest dimension. If the weakness can be mitigated (e.g., converting from Tier C to Tier A format), mitigate it and rescore. If not, proceed with the moderate score noted in the import record.
- **6-10:** Low trust. Requires enhanced review. The operator must write a specific justification for admitting this content, referencing the institutional need it serves. Consider whether the content can be obtained from a higher-trust source.
- **4-5:** Very low trust. Rejected by default. Admission requires a Tier 3 governance decision with documented justification.

### 5.7 Stage 7: Approval, Transfer, and Post-Import Verification

**Step 7.1: Approval decision.** Review the complete import record: physical inspection, provenance, inventory, scan results, format validation, and trust score. Make the admission decision:

- **Approve:** Content meets all requirements and the trust score supports admission.
- **Approve with conditions:** Content is admitted but with metadata flags (e.g., "unverified source," "format conversion performed," "low trust score -- use with caution").
- **Defer:** Content requires additional review or a governance decision before admission.
- **Reject:** Content fails one or more pipeline stages and cannot be remediated.

Record the decision, the reasoning, and any conditions in the import log. Sign and date the entry.

**Step 7.2: Transfer to production.** For approved content only:
- Copy the approved files from the quarantine workspace to the dedicated transfer medium (the clean USB drive reserved for quarantine-to-production transfers).
- Calculate checksums of the files on the transfer medium.
- Disconnect the transfer medium from the quarantine machine.
- Connect the transfer medium to the production system.
- Copy the files to the appropriate location in the institutional file structure.
- Verify the copy by comparing checksums against those calculated in the quarantine workspace.
- Disconnect and return the transfer medium to its storage location.

**Step 7.3: Post-import verification.** On the production system:
- Open each imported file and verify it renders correctly.
- Verify that the file is in the expected location within the institutional directory structure.
- Update the institutional catalog with the new content, including all provenance metadata.
- Create or update the relevant index entries so the content is discoverable.

**Step 7.4: Close the import record.** In the import log, record:
- Final disposition (approved, approved with conditions, rejected, deferred).
- Date and time of completion.
- Location of admitted content in the institutional file structure.
- Any notes for future reference.

**Step 7.5: Clean the quarantine workspace.** Delete all files from the quarantine workspace. If any malware was detected during this import, rebuild the quarantine machine from its known-good image instead of simply deleting files.

## 6. Failure Modes

- **Quarantine bypass.** The operator connects the incoming medium directly to a production system, skipping quarantine entirely. This is the most serious failure mode because it defeats the entire pipeline. Mitigation: physical separation. The production systems should not have available USB ports during normal operation (disabled in BIOS/UEFI or physically blocked). The quarantine machine should be the only system with accessible external media ports.

- **Stale malware signatures.** The quarantine machine's scanning tools have not been updated in months or years. They cannot detect recent threats. Mitigation: schedule regular updates of scanning signatures as part of the import pipeline itself (scanning tool updates are themselves imports and must go through the pipeline). Document the signature date in every import record so future reviewers know the limitation.

- **Trust score inflation.** The operator consistently scores content higher than warranted because they want to streamline the import process. Mitigation: the quarterly documentation quality audit (D19-002) includes a review of recent trust scores for consistency and calibration.

- **Provenance fabrication.** The operator records provenance that is more detailed or certain than the actual chain of custody warrants. Mitigation: ETH-001 Principle 6 (Honest Accounting of Limitations). It is better to record "provenance uncertain" than to fabricate a confident-sounding provenance that is untrue.

- **Format validation theater.** The operator checks the format validation box without actually opening and inspecting the files. Mitigation: the import log requires specific entries (file opened, rendering verified, structure validated) that are difficult to fabricate without actually doing the work. The quarterly audit reviews import logs for signs of perfunctory entries.

- **Transfer medium contamination.** The dedicated transfer medium is used for a non-approved purpose and becomes contaminated. Mitigation: the transfer medium is physically labeled and stored separately. Its use is restricted to approved quarantine-to-production transfers only. If its integrity is ever in doubt, it is reformatted or replaced.

## 7. Recovery Procedures

1. **If content was admitted that should have been rejected:** Identify all files admitted from the compromised import. Quarantine them on the production system (move to an isolated directory, remove from indexes, flag in the catalog). Assess whether the content has caused any damage (data corruption, configuration changes). If damage is suspected, follow the incident response procedure in D10-003. Trace the pipeline failure -- which stage should have caught the problem? Revise that stage's procedures to prevent recurrence. Document the incident.

2. **If the quarantine machine is compromised:** Do not use the quarantine machine for any further imports. Rebuild it from the known-good image. Review the import that caused the compromise to understand the attack vector. Update scanning tools and procedures to address the vulnerability. Document the incident in the security log.

3. **If provenance records are found to be inaccurate:** Identify all content associated with the inaccurate provenance. Re-evaluate the trust scores for that content using corrected information. If the corrected provenance significantly reduces the trust score, consider re-validating or removing the content. Correct the provenance record by appending a correction entry (do not modify the original; add a new entry referencing it).

4. **If the import log has gaps:** Reconstruct what can be reconstructed from the institutional catalog (which should contain provenance metadata for all admitted content). Create RECONSTRUCTED entries in the import log for any imports that can be identified but were not properly logged. Resume rigorous logging immediately.

5. **If a large volume of content needs emergency import:** Do not shortcut the pipeline. Instead, allocate additional time. Process imports in batches, completing the full pipeline for each batch before starting the next. If the volume is truly overwhelming, prioritize by institutional need and process the highest-priority content first. Document the emergency import decision in the governance log as a Tier 3 decision.

## 8. Evolution Path

- **Years 0-5:** The import pipeline is being exercised for the first time. Expect it to feel slow. This is intentional -- thoroughness now prevents contamination later. Keep detailed notes on which steps take the longest, which steps catch the most issues, and which steps feel unnecessary. Use these notes to refine the pipeline during annual review. Do not remove steps based on the absence of problems; the absence of problems may be the steps working.

- **Years 5-15:** The pipeline should be well-practiced. The operator can execute it from memory but should still follow the documented procedure to guard against drift. The trust scoring calibration should be validated against years of import history. Scanning tools should have been updated multiple times. Consider whether additional format tiers need to be defined for content types that did not exist in 2026.

- **Years 15-30:** Hardware and media types will have changed. USB drives may be obsolete. The pipeline steps remain valid but the physical implementation may need significant updates. The principles -- quarantine isolation, provenance documentation, trust scoring, format validation -- survive technology transitions. The implementations do not.

- **Years 30-50+:** A successor may be operating the pipeline. The import log, stretching back decades, is itself a valuable institutional resource -- a record of everything the institution has ever admitted from the outside world, and the reasoning behind each admission. Protect this log with the same care as the decision log.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The first time I executed this pipeline, it took nearly two hours for a single USB drive containing six PDF files. This will feel absurd to anyone accustomed to drag-and-drop file management. It felt absurd to me. But I also caught a mismatched file extension -- a file labeled .pdf that was actually a .html file with embedded JavaScript. In a connected environment, this would be a minor curiosity. In an air-gapped institution, it is exactly the kind of thing the pipeline is designed to catch. Two hours well spent.

The trust scoring worksheet felt subjective the first few times. It will become less subjective as the operator accumulates calibration experience -- a feel for what a "3" really means versus a "4." The first year's trust scores should be reviewed during the annual audit to check for consistency and to recalibrate if needed. Treat the first year as calibration, not as permanent judgment.

## 10. References

- D18-001 -- Import & Quarantine Philosophy (the principles this procedure implements)
- ETH-001 -- Ethical Foundations (Principle 2: Integrity Over Convenience; Principle 6: Honest Accounting of Limitations)
- CON-001 -- The Founding Mandate (Section 5.4: institutional boundaries)
- SEC-001 -- Threat Model and Security Philosophy (supply chain threats; quarantine requirements; R-SEC-03)
- OPS-001 -- Operations Philosophy (documentation-first principle; operational tempo)
- GOV-001 -- Authority Model (Tier 3 decisions for exceptional imports)
- D18-ART-003 -- Quarantine System Maintenance (rebuilding the quarantine machine)
- D18-ART-005 -- Format Conversion Specifications
- D19-002 -- Documentation Quality Audit Procedures (audit of import records)
- D20-002 -- Decision Log Operations (recording import decisions)

---

---

# D19-002 -- Documentation Quality Audit Procedures

**Document ID:** D19-002
**Domain:** 19 -- Quality Assurance
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, OPS-001, D19-001
**Depended Upon By:** All domains that produce documentation. D9-002 (onboarding includes quality expectations). D20-002 (audit results recorded in institutional memory).

---

## 1. Purpose

This article defines the practical procedures for auditing the quality of documentation within the holm.chat Documentation Institution. It takes the quality philosophy of D19-001 -- the five dimensions, the quality lifecycle, the warnings against quality tyranny -- and turns it into checklists, scoring rubrics, schedules, and templates that the operator can use to assess whether the institution's documentation is serving its purpose.

Documentation is the institution. This is not a metaphor. CON-001 states that the hardware can be replaced and the software reinstalled, but if the documentation is lost, the institution is dead. Quality assurance for documentation is therefore quality assurance for the institution itself. An audit that discovers a documentation failure has discovered an institutional vulnerability.

This article is written for the operator performing the audit -- which, in a single-operator institution, is the same person who wrote the documentation. This creates an obvious conflict of interest. The procedures in this article are designed to mitigate that conflict through structure: specific checklists, explicit scoring criteria, and a rubric that makes "good enough" a measurable standard rather than a subjective feeling.

## 2. Scope

This article covers:

- Quality checklists for each of the five quality dimensions defined in D19-001.
- Scoring rubrics with concrete criteria for each score.
- The audit schedule: what gets audited, when, and how often.
- Templates for audit reports.
- Procedures for handling audit findings.
- Calibration procedures for maintaining scoring consistency over time.

This article does not cover:

- The philosophy of quality (see D19-001).
- Quality standards for non-documentation artifacts (software, hardware configurations, operational processes are audited through their own domain procedures).
- Federation-level quality assurance (see D19-ART-010).
- Defect tracking and resolution workflows (see D19-ART-007).

## 3. Background

D19-001 warns of two equal and opposite dangers: quality neglect and quality tyranny. This article walks the line between them. The checklists and rubrics are specific enough to prevent neglect -- you cannot claim "the documentation is fine" without evaluating it against concrete criteria. But they are proportional enough to prevent tyranny -- a low-risk operational note receives a lighter audit than a root document revision.

The five quality dimensions from D19-001 are: correctness, completeness, clarity, currency, and durability. These dimensions are not equally important for every document. A time-sensitive operational note must be correct and current but need not be exceptionally durable. A root document must score highly on all five dimensions. The audit procedures account for this by weighting dimensions differently based on the document's tier.

The single-operator audit problem is real and must be addressed honestly. When you audit your own work, you bring two liabilities: you know what you meant, so you may read it into text that does not actually say it; and you are invested in the quality of your work, so you may be reluctant to score it low. The mitigations are structural: time separation (audit documentation at least 30 days after writing it, so your short-term memory of intent has faded), checklist discipline (evaluate each criterion individually rather than forming a global impression), and external calibration (during the annual review, re-audit a sample of documents that were audited in previous years to check for scoring drift).

## 4. System Model

### 4.1 Document Tiers

All institutional documents fall into one of four tiers, each with different audit requirements:

**Tier 1: Root Documents.** ETH-001, CON-001, GOV-001, SEC-001, OPS-001. These are audited annually with the full five-dimension rubric. All five dimensions are weighted equally. A score below 4 on any dimension triggers mandatory remediation.

**Tier 2: Domain Philosophy Articles.** The -001 articles for each domain. Audited annually with the full rubric. All five dimensions are weighted equally, but the remediation threshold is a score below 3 on any dimension.

**Tier 3: Operational Articles.** Stage 3 articles like this one. Audited on a rotating basis -- each article is audited at least once every two years. Correctness and currency are weighted double; durability is weighted at half (operational procedures are expected to evolve).

**Tier 4: Operational Records.** Import logs, decision log entries, maintenance records, audit reports themselves. Spot-checked during quarterly reviews. Correctness and completeness are the primary dimensions; clarity and durability are evaluated but not scored.

### 4.2 The Audit Cycle

The audit cycle is integrated into the operational tempo defined in OPS-001:

**Quarterly (during quarterly operations):** Spot-check five randomly selected Tier 4 records. Review the previous quarter's audit findings for resolution status. Review the import log's trust scores for calibration consistency.

**Annually (during annual operations):** Full audit of all Tier 1 and Tier 2 documents. Audit of one-third of Tier 3 documents (rotating so that every Tier 3 document is audited at least once every three years, and at least once every two years if resources allow). Calibration check: re-audit two previously audited documents and compare scores. Produce the Annual Documentation Quality Report.

**Triggered audits:** Any time a document is significantly revised, it receives a post-revision audit within 30 days. Any time an operational failure is traced to a documentation deficiency, the implicated document receives an immediate audit.

## 5. Rules & Constraints

### 5.1 The Five-Dimension Checklist

For each document under audit, evaluate against the following checklist. Each dimension has specific criteria. Score each criterion Yes (1) or No (0), then calculate the dimension score.

**Dimension 1: Correctness**

| # | Criterion | Y/N |
|---|-----------|-----|
| C1 | All factual claims in the document are accurate as of the audit date. | |
| C2 | All procedures described in the document produce the results they claim when followed exactly. | |
| C3 | All references to other documents point to the correct document and the correct section. | |
| C4 | All technical specifications (commands, paths, configurations) are accurate for the current system. | |
| C5 | There are no internal contradictions within the document. | |

Dimension score: (Sum of Yes answers / 5) x 5, rounded to nearest 0.5. Maximum 5.0.

**Dimension 2: Completeness**

| # | Criterion | Y/N |
|---|-----------|-----|
| K1 | All sections required by the document template are present (Purpose, Scope, Background, System Model, Rules & Constraints, Failure Modes, Recovery Procedures, Evolution Path, Commentary, References). | |
| K2 | The scope section accurately describes what the document covers and does not cover. | |
| K3 | All failure modes that a reasonably experienced operator could anticipate are documented. | |
| K4 | All recovery procedures address the documented failure modes. | |
| K5 | Dependencies (Depends On, Depended Upon By) are complete and accurate. | |

Dimension score: (Sum / 5) x 5, rounded to nearest 0.5.

**Dimension 3: Clarity**

| # | Criterion | Y/N |
|---|-----------|-----|
| L1 | A technically competent reader unfamiliar with this specific document could understand its content without external explanation. | |
| L2 | All domain-specific terms are defined at first use or in a referenced glossary. | |
| L3 | The document's structure (headings, sections, lists) supports efficient navigation. | |
| L4 | Procedures are written as explicit, numbered steps that can be followed sequentially. | |
| L5 | The document does not assume knowledge that is not documented elsewhere in the institution. | |

Dimension score: (Sum / 5) x 5, rounded to nearest 0.5.

**Dimension 4: Currency**

| # | Criterion | Y/N |
|---|-----------|-----|
| U1 | The document reflects the current state of the systems, procedures, or policies it describes. | |
| U2 | All dates, version numbers, and status indicators are current. | |
| U3 | No section describes a system, procedure, or configuration that has been superseded or deprecated. | |
| U4 | The document has been reviewed or updated within the timeframe appropriate to its tier. | |
| U5 | The Commentary Section contains at least one entry from the current operational period (most recent two years). | |

Dimension score: (Sum / 5) x 5, rounded to nearest 0.5.

**Dimension 5: Durability**

| # | Criterion | Y/N |
|---|-----------|-----|
| D1 | The document explains the "why" behind decisions, not just the "what." | |
| D2 | The document avoids references to specific tools or technologies without also explaining the underlying principle that a future equivalent must satisfy. | |
| D3 | The document uses formats that will remain readable (plain text, standard Markdown, no embedded active content). | |
| D4 | The document's value does not depend on external resources that may become unavailable. | |
| D5 | A reader fifty years from now, with no contact with the current operator, could extract useful guidance from this document. | |

Dimension score: (Sum / 5) x 5, rounded to nearest 0.5.

### 5.2 The Scoring Rubric

After calculating dimension scores, compute the weighted aggregate:

For Tier 1 and Tier 2 documents:
- Aggregate = (Correctness + Completeness + Clarity + Currency + Durability) / 5

For Tier 3 documents:
- Aggregate = (Correctness x 2 + Completeness + Clarity + Currency x 2 + Durability x 0.5) / 6.5

For Tier 4 documents (spot-check only):
- Aggregate = (Correctness + Completeness) / 2

**Score interpretation:**

- **4.5-5.0: Exemplary.** The document meets or exceeds all quality standards. No action required. Record the result.
- **3.5-4.4: Acceptable.** The document meets quality standards with minor issues. Note specific criteria that scored No. Schedule remediation during routine maintenance. The document remains in service.
- **2.5-3.4: Below standard.** The document has significant quality issues. Specific remediation actions must be documented and scheduled within 30 days. The document remains in service but is flagged in the catalog as "pending quality remediation."
- **1.0-2.4: Failing.** The document has critical quality issues. It is flagged as "unreliable" in the catalog. Remediation is prioritized above routine work. If the document describes a safety-critical procedure (security, backup, recovery), an immediate interim notice is posted directing operators to verify the procedure against actual system state before following it.

### 5.3 The Audit Report Template

Each audit produces a report with the following structure:

```
DOCUMENTATION QUALITY AUDIT REPORT

Audit ID: QA-YYYY-NNN
Date: [audit date]
Auditor: [operator identifier]
Document Audited: [document ID and title]
Document Tier: [1/2/3/4]
Document Version: [version number]
Last Audit Date: [date of previous audit, or "INITIAL" if first audit]

DIMENSION SCORES:
  Correctness:   [score] / 5.0
  Completeness:  [score] / 5.0
  Clarity:       [score] / 5.0
  Currency:      [score] / 5.0
  Durability:    [score] / 5.0

WEIGHTED AGGREGATE: [score] / 5.0
RATING: [Exemplary / Acceptable / Below Standard / Failing]

FINDINGS:
  [For each criterion that scored No, describe the specific deficiency
   and the specific remediation required.]

COMPARISON TO PREVIOUS AUDIT:
  [If previously audited, note dimensions that improved, declined, or
   remained stable. Flag any declining trends.]

REMEDIATION SCHEDULE:
  [For Below Standard or Failing: specific actions with deadlines.]

AUDITOR NOTES:
  [Any observations not captured above. Include notes on whether the
   audit process itself worked well for this document type.]
```

### 5.4 Calibration Procedures

To maintain scoring consistency over time:

**Annual calibration (during annual operations):** Select two documents that were audited in the previous year. Re-audit them using the same checklist and rubric. Compare the new scores to the previous scores. If any dimension score differs by more than 1.0 point, investigate:
- Did the document change? (Check version history.)
- Did the systems it describes change? (Check operational logs.)
- Did your scoring standards drift? (Most likely if the document and systems are unchanged but the score shifted.)

If scoring drift is detected, review the affected dimension's criteria. Consider whether your understanding of the criteria has changed. Document the calibration result and any scoring adjustments.

**Cross-temporal calibration:** During the annual review, read one audit report from the earliest available year. Assess whether the standards applied then are consistent with the standards applied now. Document any observed drift and its likely cause.

## 6. Failure Modes

- **Audit avoidance.** The operator consistently postpones or skips audits because they are time-consuming and the results are predictable. Years pass without audit. Documentation quality degrades undetected. Mitigation: the audit schedule is part of the operational tempo (OPS-001). Skipped audits are documented in the operational log. Two consecutive skipped quarterly audits trigger a governance review.

- **Score inflation.** The operator consistently scores their own work higher than warranted. Mitigation: time separation (audit at least 30 days after writing), checklist discipline (evaluate each criterion individually), and annual calibration checks.

- **Checklist rigidity.** The checklist becomes a straitjacket. Documents that serve their purpose well but do not conform to the checklist template receive low scores. Mitigation: the auditor notes section allows the auditor to override a low aggregate score with a written justification. The override must explain why the document serves its purpose despite the low checklist score.

- **Remediation backlog.** Audit findings accumulate faster than they can be remediated. The operator is perpetually behind on quality fixes. Mitigation: prioritize remediation by document tier and by the severity of the finding. Not every Below Standard finding requires immediate action. Schedule remediation within the sustainable operational tempo. If the backlog grows beyond what the operational tempo can absorb, declare a documentation sprint per CON-001 Recovery Procedure 2.

- **Quality tyranny by checklist.** The checklist becomes more important than the purpose it serves. Documents are written to satisfy the checklist rather than to serve their readers. Mitigation: D19-001, R-D19-01 (Quality Serves the Mission). The checklist is a tool. The mission is the master. If the checklist and the mission conflict, revise the checklist.

- **Audit of audit failure.** The audit procedures themselves are never audited. The quality of quality assurance degrades. Mitigation: D19-001, R-D19-03 (QA Must Pass Its Own Standards). This article is itself subject to audit using the same checklist. Schedule it for audit in the second annual cycle.

## 7. Recovery Procedures

1. **If audits have been skipped for an extended period:** Do not attempt to audit everything at once. Begin with Tier 1 documents (the root documents). Complete their audits. Then proceed to Tier 2. Then resume the normal rotating schedule for Tier 3. Acknowledge the gap in the audit record and document it.

2. **If score inflation is detected during calibration:** Re-audit a sample of documents from the inflated period, applying the corrected standard. Identify documents that would have been rated Below Standard or Failing under the corrected standard. Add them to the remediation queue. Document the calibration correction.

3. **If the remediation backlog is overwhelming:** Triage. Classify all pending remediations by severity: Critical (Failing documents describing safety-critical procedures), High (Failing documents, Below Standard root documents), Medium (Below Standard operational documents), Low (Acceptable documents with minor findings). Address Critical and High first. Defer Low. Consider whether some Medium findings can be resolved by updating the document during the next routine revision rather than as a separate remediation task.

4. **If the audit process itself is found inadequate:** Revise the checklist or rubric through a Tier 3 governance decision. Document the revision rationale. Do not retroactively re-score documents against the revised criteria (per D19-001, R-D19-06: Standards Evolve Forward). Apply the revised criteria to future audits only.

5. **If a documentation failure causes an operational incident:** Conduct an immediate triggered audit of the implicated document. Identify the specific quality failure (was it incorrect? incomplete? outdated?). Remediate immediately. Then assess whether the same type of failure could exist in other documents. If so, schedule targeted audits of similar documents.

## 8. Evolution Path

- **Years 0-5:** The audit process is being established. The first few audits will be slow and the scoring will feel uncertain. This is normal. Treat the first year's audits as calibration -- the scores are real, but the operator's confidence in the scoring should increase over time. Accumulate audit reports. The reports themselves become a valuable longitudinal record of documentation quality.

- **Years 5-15:** The audit process should be routine. Patterns should be visible: which types of documents consistently score well? Which dimensions are chronically weak? Use these patterns to focus quality improvement efforts. Consider whether the checklist criteria need revision based on operational experience.

- **Years 15-30:** A successor may be performing audits. The rubric and checklists must be clear enough for someone who did not design them to apply them consistently. The calibration procedures become critical -- the successor must be able to score documents consistently with the institution's historical standards, not just their own intuition.

- **Years 30-50+:** The audit reports spanning decades are themselves a history of the institution's quality culture. They show not just whether quality was maintained but how the institution thought about quality at different points in its life. Protect these records as part of the institutional memory (Domain 20).

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The irony of auditing my own documentation is not lost on me. Every criterion in the clarity checklist -- "could a reader unfamiliar with this document understand it?" -- is impossible for me to evaluate objectively, because I am the most familiar possible reader. The time-separation mitigation helps but does not eliminate the problem. Future operators: if you have the opportunity to have someone else audit your documentation, even informally, even a non-technical reader checking for clarity, take it. External eyes see what internal eyes cannot.

I have deliberately set the Failing threshold quite low (below 2.5). This is not because I believe low-quality documentation is acceptable. It is because I believe the quality system must be usable before it can be aspirational. A system that flags everything as Failing will be abandoned. A system that identifies genuine failures while accepting honest imperfection will be maintained. The threshold can be raised as the institution's quality culture matures. Start achievable. Evolve toward excellent.

## 10. References

- D19-001 -- Quality Philosophy (the principles this procedure implements)
- ETH-001 -- Ethical Foundations (Principle 2: Integrity Over Convenience; Principle 6: Honest Accounting)
- CON-001 -- The Founding Mandate (documentation as the institution)
- GOV-001 -- Authority Model (decision tiers for quality standard changes)
- OPS-001 -- Operations Philosophy (operational tempo; quarterly and annual review cycles)
- D18-002 -- Content Import Pipeline (audit of import records and trust scores)
- D20-002 -- Decision Log Operations (recording audit results in institutional memory)
- D19-ART-007 -- Defect Tracking and Resolution Workflows
- D19-ART-010 -- Federation-Level Quality Assurance

---

---

# D9-002 -- New Maintainer Onboarding Guide

**Document ID:** D9-002
**Domain:** 9 -- Knowledge & Training
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, SEC-001, OPS-001, D9-001
**Depended Upon By:** GOV-001 (succession protocol references this onboarding guide). D19-002 (onboarding includes quality audit training). D20-002 (onboarding includes decision log operations).

---

## 1. Purpose

This article defines the complete onboarding curriculum for a new maintainer of the holm.chat Documentation Institution. It answers the question: if a person who has never seen this institution before is given physical access and told "this is now your responsibility," what do they need to learn, in what order, and how do they know when they are ready?

This is the succession protocol's practical companion. GOV-001 defines who can assume authority and what governance steps are required. This article defines what that person needs to learn and how they learn it. GOV-001 is the legal framework. This article is the training program.

The onboarding guide is designed for the worst case: an unplanned succession where the new maintainer has no access to the previous operator. Every lesson, exercise, and assessment in this guide can be completed using only the institution's own documentation and systems. No oral instruction is required. No prior relationship with the founder is assumed. The institution must be able to train its own successor from its documentation alone.

This is the hardest test of the documentation-first principle (OPS-001, Section 4.2). If the documentation is truly complete, onboarding from documentation alone is possible. If onboarding fails, the documentation has failed.

## 2. Scope

This article covers:

- The reading order for institutional documents.
- Skill assessments for each phase of onboarding.
- Hands-on exercises that develop operational competence.
- Milestone checkpoints that verify readiness for increasing responsibility.
- The "first 90 days" plan: a day-by-day structure for the onboarding period.
- Criteria for onboarding completion: what "ready" means.

This article does not cover:

- The governance aspects of succession (see GOV-001, Section 7.1).
- The specific security credentials and how they are transferred (see SEC-001, SEC-003).
- The emotional and psychological aspects of inheriting an institution you did not build (acknowledged here, addressed in D9-ART-005).
- Ongoing professional development after onboarding (see D9-ART-003).

## 3. Background

Succession is the institution's most vulnerable moment. The previous operator's knowledge -- the undocumented intuitions, the muscle memory, the "I just know where that file is" -- is gone. What remains is the documentation, the systems, and the new maintainer's capacity to learn.

The onboarding curriculum is structured around a simple model: comprehension before authority. The new maintainer must understand what the institution is, why it exists, how it works, and how to keep it alive before they exercise the authority to change it. This is not gatekeeping. It is the difference between a surgeon who has studied anatomy and a person with a scalpel.

The 90-day structure is borrowed from organizational onboarding best practice, adapted for a single-person institution. The first 30 days focus on reading and observation -- understanding the institution as it is. The second 30 days focus on supervised operation -- performing routine tasks with reference to documentation. The third 30 days focus on independent operation -- executing the full operational tempo without reference assistance, demonstrating readiness.

In a planned succession, the previous operator is available during onboarding. In an unplanned succession, the documentation is the only teacher. The onboarding guide is written for the unplanned case. If the previous operator is available, their presence is a bonus, not a requirement.

## 4. System Model

### 4.1 The Onboarding Phases

Onboarding proceeds through four phases:

**Phase 0: Orientation (Days 1-3).**
Understanding what this institution is, why it exists, and what the new maintainer is taking responsibility for.

**Phase 1: Comprehension (Days 4-30).**
Reading and understanding the institutional documentation, from the ethical foundations through the operational procedures.

**Phase 2: Supervised Operation (Days 31-60).**
Performing the operational tempo under guidance from the documentation, with each task verified against expected outcomes.

**Phase 3: Independent Operation (Days 61-90).**
Operating the institution without reference assistance, demonstrating the ability to handle routine operations and respond to simulated anomalies.

### 4.2 The Competency Model

The new maintainer must demonstrate competency in five domains:

**Philosophical Competency:** Understanding why the institution exists and what principles govern it. Assessed through reading comprehension and written reflection.

**Governance Competency:** Understanding how decisions are made, recorded, and reviewed. Assessed through the ability to classify decisions by tier and create proper decision log entries.

**Security Competency:** Understanding the threat model, the air gap, and the security architecture. Assessed through the ability to explain the trust model and execute security procedures.

**Operational Competency:** Ability to execute the daily, weekly, monthly, quarterly, and annual operational cycles. Assessed through observed performance of each cycle.

**Recovery Competency:** Ability to diagnose and recover from system failures using documented procedures. Assessed through simulated failure exercises.

## 5. Rules & Constraints

### 5.1 Phase 0: Orientation (Days 1-3)

**Day 1: First Contact**

Reading list (in order):
1. This document (D9-002) -- to understand the onboarding process itself.
2. CON-001 -- The Founding Mandate. Understand what the institution is and why it exists.
3. ETH-001 -- Ethical Foundations. Understand the values that govern all decisions.

Exercise 1.1: After reading CON-001 and ETH-001, write a one-page summary in your own words: What is this institution? What does it value? What does it refuse to do?

Assessment: The summary must accurately capture the institutional mission, the six ethical principles, and the air-gap mandate. It need not be eloquent. It must be correct.

**Day 2: Governance and Security Overview**

Reading list:
4. GOV-001 -- Authority Model. Understand decision tiers, the decision log, and the succession protocol.
5. SEC-001 -- Threat Model. Understand the threat categories, the trust model, and the air-gap architecture.

Exercise 2.1: List the four decision tiers and give one example of a decision that would fall into each tier. The examples should be original, not copied from GOV-001.

Exercise 2.2: In your own words, explain why the institution is air-gapped. What threats does the air gap address? What does it cost?

**Day 3: Operations Overview**

Reading list:
6. OPS-001 -- Operations Philosophy. Understand the operational tempo, the documentation-first principle, and the complexity budget.

Exercise 3.1: Write out the daily, weekly, monthly, quarterly, and annual operational cycles from memory (consulting the document is permitted, but the exercise is to internalize the rhythm). For each cycle, write one sentence explaining its purpose.

Milestone Checkpoint 0: Review your Phase 0 exercises. If you can accurately describe the institution's mission, values, governance, security posture, and operational rhythm, proceed to Phase 1. If not, re-read the relevant documents and redo the exercises.

### 5.2 Phase 1: Comprehension (Days 4-30)

**Days 4-10: Domain Philosophy Articles**

Reading list: All Stage 2 domain philosophy articles, in domain number order. Focus on understanding how each domain relates to the five root documents.

Exercise 1.1 (Days 4-10): For each domain philosophy article, create a one-paragraph summary that answers: What does this domain govern? What is its relationship to the root documents? What are its most important rules?

**Days 11-17: Stage 3 Operational Articles**

Reading list: All Stage 3 operational articles, in the order they appear. Focus on understanding the procedures themselves -- what the operator actually does.

Exercise 1.2: Select three operational articles from different domains. For each, trace the authority chain: which root document authorizes this procedure? Which domain philosophy does it implement? What happens if this procedure is not followed?

**Days 18-24: The Decision Log and Institutional Memory**

Reading list: The complete decision log, from the earliest entry to the most recent. This may take several days depending on the log's length. Read every entry, including the rationale and the alternatives considered.

Exercise 1.3: Select five decisions from the log that you find surprising, confusing, or questionable. For each, write a brief note explaining what surprised you and what you think the rationale was. Then compare your guess to the recorded rationale. Note discrepancies.

**Days 25-28: Systems Familiarization**

Conduct a guided walkthrough of the institution's physical and logical systems. Using the infrastructure documentation (Domain 4) and the platform documentation (Domain 5):

Exercise 1.4: Create a physical map of the institution: where is each piece of hardware? What does it do? How are the pieces connected (or deliberately not connected)?

Exercise 1.5: Create a logical map: what services run on each machine? What data lives where? What are the backup locations?

**Days 29-30: Phase 1 Assessment**

Complete the Phase 1 Assessment:

1. Written test (open-book, using institutional documentation):
   - Classify ten hypothetical decisions by governance tier. (Seven of ten correct to pass.)
   - Identify the appropriate recovery procedure for five hypothetical failure scenarios. (Four of five correct to pass.)
   - Explain the import pipeline stages and what each stage prevents. (All stages must be named and explained.)

2. Practical assessment:
   - Locate a specific document in the institutional corpus using only the navigation and search tools.
   - Locate a specific decision in the decision log and explain its context and implications.
   - Identify the current status of the institution's backup systems from the operational logs.

Milestone Checkpoint 1: All Phase 1 assessments must be passed before proceeding to Phase 2. If any assessment is failed, re-read the relevant material and retake after five days.

### 5.3 Phase 2: Supervised Operation (Days 31-60)

During Phase 2, the new maintainer performs the operational tempo under documentation guidance. "Supervised" means the maintainer follows the documented procedures exactly and verifies each outcome against the expected result documented in the procedures. In a planned succession, the previous operator observes and provides feedback. In an unplanned succession, the documentation is the supervisor.

**Days 31-37: Daily Operations**

Perform the daily operational cycle (OPS-001, Section 4.1) every day for seven days. After each session:
- Record what you did in the operational log.
- Note any step where the documentation was unclear or the expected outcome did not match the actual outcome.
- Note any step where you had to improvise or make a judgment call.

Exercise 2.1: After seven days, write a brief assessment: Are the daily operations sustainable? Were any procedures unclear? Did you discover any documentation gaps?

**Days 38-44: Weekly Operations**

Perform the weekly operational cycle at least twice during this period. Follow the same documentation-guided process as for daily operations.

Exercise 2.2: During a weekly cycle, deliberately verify one backup by performing a test restoration of a randomly selected file. Document the process and the result.

**Days 45-51: Monthly Operations**

Perform one complete monthly operational cycle. This includes backup integrity verification, hardware health assessment, and threat model awareness review.

Exercise 2.3: Conduct a complete hardware health assessment using the procedures documented in Domain 4. Record all SMART data, visual inspection results, and environmental readings. Compare against historical baselines if available.

**Days 52-58: Import and Audit**

Perform one complete content import using D18-002. If no external content needs to be imported, perform a practice import using a test medium prepared for this purpose (a USB drive with known-safe files of various formats, including deliberate test cases: a mismatched file extension, a file in a Tier C format, a file with incomplete provenance).

Perform one documentation quality audit using D19-002 on a Tier 3 document.

Exercise 2.4: Complete the trust scoring worksheet for the practice import. Write a brief reflection on the experience: did the scoring feel calibrated? Which dimensions were hardest to assess?

**Days 59-60: Phase 2 Assessment**

1. Log review: An auditor (the previous operator, or the new maintainer performing self-review 30 days later in unplanned succession) reviews all operational log entries from Phase 2. Assessment criteria:
   - Were daily operations performed every day?
   - Were all procedures followed as documented?
   - Were all outcomes recorded?
   - Were documentation gaps identified and noted?

2. Practical demonstration:
   - Execute the daily operational cycle while narrating each step and its purpose.
   - Respond to a simulated anomaly (provided by the previous operator, or drawn from the Anomaly Exercise Bank in D9-ART-006): a system log entry that indicates a potential problem. Demonstrate the ability to assess, investigate, and respond.

Milestone Checkpoint 2: Phase 2 assessment must be passed before proceeding to Phase 3. Key criterion: the maintainer can execute routine operations from documentation, identify anomalies, and respond appropriately.

### 5.4 Phase 3: Independent Operation (Days 61-90)

During Phase 3, the new maintainer operates the institution independently. The previous operator (if available) is accessible for questions but does not observe or guide. The maintainer is expected to handle the full operational tempo and to manage issues as they arise, using the documentation as their primary resource.

**Days 61-75: Full Operational Tempo**

Operate the institution normally. Perform daily and weekly cycles. Handle any import requests that arise. Maintain the decision log. The previous operator does not intervene unless the maintainer requests help or unless a situation arises that could cause irreversible damage.

**Days 76-82: Simulated Failure Exercises**

Complete three simulated failure exercises from the Anomaly Exercise Bank:

Exercise 3.1: **Backup restoration.** A critical file has been "lost" (moved to a hidden location by the previous operator, or simply designated as the target). Restore it from backup using the documented recovery procedures. Document the entire process.

Exercise 3.2: **Quarantine machine failure.** The quarantine machine is declared "compromised" (simulated). Rebuild it from the known-good image using the documented procedures in D18-ART-003. Verify the rebuild. Document the process.

Exercise 3.3: **Decision under uncertainty.** A situation is presented that is not clearly covered by existing documentation -- a borderline case that requires judgment. The maintainer must: classify the decision by tier, identify the most relevant principles and procedures, make a decision, and document it in the decision log with full rationale.

**Days 83-88: First Quarterly Review**

Perform a complete quarterly operational cycle (OPS-001, Section 4.1). This is the most comprehensive routine operation and tests the broadest range of competencies: documentation review, decision log review, complexity budget assessment, and governance health check.

**Days 89-90: Onboarding Completion Assessment**

Final assessment. The maintainer must demonstrate:

1. **Philosophical competency:** Write a one-page reflection on the institution's ethical foundations and how they apply to a specific recent operational decision. The reflection must demonstrate genuine understanding, not rote recitation.

2. **Governance competency:** Present the decision log entries created during Phase 3. The entries must be properly formatted, classified by tier, and include adequate rationale.

3. **Operational competency:** The operational log from Phase 3 must show consistent execution of the daily and weekly cycles, proper documentation of all activities, and appropriate responses to anomalies.

4. **Security competency:** Explain the institution's threat model and describe the specific security measures in place. Demonstrate the ability to execute the import pipeline without reference to the documentation (documentation may be available for confirmation, but the maintainer should be able to execute from memory).

5. **Recovery competency:** Present the documentation from the three simulated failure exercises. The documentation must show correct diagnosis, appropriate procedure selection, successful recovery, and thorough post-incident documentation.

Milestone Checkpoint 3 (Onboarding Complete): All five competency areas must be demonstrated satisfactorily. The completion is recorded in the decision log as a Tier 2 governance event. If this is a planned succession, the previous operator certifies the completion. If unplanned, the maintainer self-certifies, noting the absence of external verification.

## 6. Failure Modes

- **Reading without comprehension.** The new maintainer reads all the documentation but does not internalize it. They pass the reading phase but struggle in the operational phases. Mitigation: the exercises at each phase require active engagement, not passive reading. The Phase 1 assessment tests comprehension, not completion.

- **Operational mimicry.** The new maintainer can follow procedures but does not understand why the procedures exist. When a situation arises that the procedures do not cover, they are lost. Mitigation: the Phase 3 "decision under uncertainty" exercise specifically tests the ability to reason from principles when procedures are insufficient.

- **Documentation gaps revealed during onboarding.** The new maintainer encounters a system or procedure that is inadequately documented. They cannot proceed without oral explanation. Mitigation: this is a documentation failure, not an onboarding failure. The gap must be documented and flagged for remediation. If the previous operator is available, they provide the missing information and immediately document it. If not, the maintainer must reconstruct from system inspection and document the reconstruction.

- **Onboarding too slow.** The new maintainer takes significantly longer than 90 days. The operational tempo suffers. Mitigation: the 90-day timeline is a target, not a hard limit. Some maintainers will need more time, especially in unplanned successions. The critical constraint is competency, not speed. Extend the timeline if needed, but do not certify completion until competencies are demonstrated.

- **Onboarding too fast.** The new maintainer rushes through the phases, treating onboarding as a formality rather than a learning process. Mitigation: the milestone checkpoints are gates, not suggestions. Each checkpoint requires demonstrated competency. Rushing through the reading does not produce demonstrated competency in the assessment.

- **Previous operator interference.** In planned succession, the previous operator cannot let go. They continue to make decisions, override the new maintainer, or skip onboarding steps because "I can just show them." Mitigation: Phase 3 explicitly requires independent operation. The previous operator's role during Phase 3 is to be available for questions, not to operate the institution.

## 7. Recovery Procedures

1. **If onboarding stalls in Phase 1:** Identify the specific comprehension barrier. Is it a specific document that is unclear? (Flag for quality remediation.) Is it the volume of material? (Extend the timeline; break the reading into smaller segments.) Is it a fundamental mismatch between the maintainer's background and the institution's requirements? (This is a succession planning failure, not an onboarding failure. Consult GOV-001.)

2. **If onboarding reveals critical documentation gaps:** Treat the gaps as high-priority remediation items. In a planned succession, the previous operator documents the missing information immediately. In an unplanned succession, the new maintainer creates RECONSTRUCTED documentation based on system inspection, flagged for future verification.

3. **If the maintainer fails an assessment:** Identify the specific competencies that were not demonstrated. Provide targeted re-study of the relevant material. Allow a re-take after a waiting period (seven days minimum for Phase 1, fourteen days for Phases 2 and 3). Two consecutive failures in the same assessment trigger a governance review.

4. **If unplanned succession occurs with no prior onboarding preparation:** The new maintainer begins at Phase 0 with no external support. The 90-day timeline may need to be extended. The institution enters maintenance-only mode (no development, no architectural changes) until the maintainer reaches Phase 3. Prioritize: daily operations first, security verification second, comprehensive understanding third.

5. **If the previous operator is available but resistant to succession:** This is a governance issue, not a training issue. Consult GOV-001. The succession protocol has clear steps. Follow them. The onboarding curriculum is executed regardless of the previous operator's comfort with succession.

## 8. Evolution Path

- **Years 0-5:** The onboarding guide is theoretical -- no succession has occurred. But the guide should be tested. The founding operator should "onboard themselves" by following the Phase 1 reading list and completing the exercises. This is both a test of the guide and a useful exercise in re-familiarization.

- **Years 5-15:** If a succession occurs during this period, the onboarding guide receives its first real test. Expect it to need significant revision based on the experience. Every gap, confusion, and inadequacy discovered during the first real onboarding must be documented and addressed. The post-onboarding debrief is the guide's most important source of improvement data.

- **Years 15-30:** The guide should have been tested at least once. Its revisions should reflect real operational experience. The institution's documentation should be mature enough that documentation-only onboarding is genuinely possible. If it is not, the documentation-first principle has a gap that must be addressed.

- **Years 30-50+:** The guide may have been used multiple times. It should be streamlined by experience but still comprehensive. The founding operator's original assumptions about what a new maintainer needs to know should be validated or revised by the accumulated experience of actual onboarding events.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
Writing an onboarding guide for an institution that has never had a second operator feels like writing a travel guide for a country nobody has visited. Every assumption about what a new maintainer will find confusing, what they will find obvious, and what they will need to practice is, at this point, a guess. An educated guess, informed by the experience of onboarding into other organizations, but a guess nonetheless.

The most important thing I can say to the first person who uses this guide in earnest: please revise it. Your experience of actually being onboarded will reveal things I could not predict. Write them down. Update the guide. The second person to use it should have a better experience than you did. That improvement, compounded over generations of maintainers, is how institutional knowledge grows.

I have deliberately structured the guide around independence rather than mentorship. In a perfect world, every succession would be planned and the previous operator would be available for months of patient instruction. In reality, people die, become incapacitated, or simply leave without warning. The guide must work in the worst case. If it works in the worst case, it will certainly work in the best case. The reverse is not true.

## 10. References

- GOV-001 -- Authority Model (succession protocol; decision tiers; Section 7.1)
- ETH-001 -- Ethical Foundations (the principles the maintainer must internalize)
- CON-001 -- The Founding Mandate (the mission the maintainer must understand)
- SEC-001 -- Threat Model (the security posture the maintainer must maintain)
- OPS-001 -- Operations Philosophy (the operational tempo the maintainer must execute)
- D9-001 -- Knowledge & Training Philosophy (the learning principles this guide implements)
- D18-002 -- Content Import Pipeline (import procedure practiced during onboarding)
- D19-002 -- Documentation Quality Audit Procedures (audit procedure practiced during onboarding)
- D20-002 -- Decision Log Operations (decision log operations practiced during onboarding)
- D9-ART-005 -- Psychological Dimensions of Institutional Succession
- D9-ART-006 -- Anomaly Exercise Bank

---

---

# D20-002 -- Decision Log Operations

**Document ID:** D20-002
**Domain:** 20 -- Institutional Memory
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, SEC-001, OPS-001, D20-001
**Depended Upon By:** All domains that produce decisions. D9-002 (onboarding includes decision log training). D19-002 (audit of decision log quality).

---

## 1. Purpose

This article defines the day-to-day procedures for maintaining the decision log of the holm.chat Documentation Institution. It takes the philosophy of institutional memory established in D20-001 -- the sacred, append-only record, the anti-revisionism principles, the context recovery model -- and turns it into a practical manual for writing entries, classifying decisions, tagging them for retrieval, searching the log, and conducting the annual review.

The decision log is the institution's most important artifact. D20-001 makes the case for why. This article tells you how to keep it. If you follow one procedure in this entire institution with absolute fidelity, let it be this one. Every other procedure can be recovered from documentation if it is forgotten. But if the decision log is lost or corrupted, the institution loses its memory, and an institution without memory is an institution that will repeat every mistake and rediscover every lesson from scratch.

This article is written for the operator who has just made a decision and needs to record it before the context fades. The format is designed to be filled out in minutes, not hours. The classification system is designed to be intuitive after the first few uses. The search system is designed to answer the question that matters most: "Have we faced this before?"

## 2. Scope

This article covers:

- The decision log entry format: every field, what it means, and how to fill it in.
- Decision classification: how to assign a tier and a category.
- Tagging conventions: how to tag entries for future retrieval.
- Recording decisions in real-time: the practical workflow for capturing decisions as they happen.
- Searching the decision log: how to find what you need.
- The annual review procedure: how to review the log for completeness, patterns, and lessons.
- Log maintenance: integrity verification, backup, and format migration.

This article does not cover:

- The philosophy of institutional memory (see D20-001).
- Oral history capture (see D20-ART-003).
- Lessons-learned frameworks (see D20-ART-004).
- The decision log's technical storage format and schema (see D20-ART-002).
- Distributed decision logs across federated nodes (see D20-ART-010).

## 3. Background

GOV-001 established that every decision at Tier 1, 2, or 3 must be recorded in the decision log. OPS-001 established that operational logs are maintained in an append-only format. D20-001 established that the decision log is inviolable -- entries may not be modified or deleted.

This article operationalizes these requirements. The challenge is not philosophical (everyone agrees the decision log is important) but practical (recording decisions is tedious, especially when you are in the middle of doing something and would rather keep working than stop and write). The procedures here are designed to minimize the friction of recording while maximizing the value of the record.

The key insight from D20-001 is that the most valuable part of a decision record is not the decision itself -- it is the rationale. "We decided to use Format X" is useful but limited. "We decided to use Format X because Format Y had no open specification, Format Z was not supported by our archival tools, and Format X met all five of our format criteria with the sole disadvantage of larger file sizes, which we accepted as a reasonable cost" -- that is a record that serves future operators. The entry format in this article requires rationale. It is the one field that cannot be abbreviated.

## 4. System Model

### 4.1 The Decision Log Architecture

The decision log is a single, append-only file (or set of files organized by year) stored in the institutional memory directory. Its format is plain text (UTF-8) with structured fields, chosen for maximum durability and minimum tooling dependency. Any text editor can read it. Any text search tool can search it.

The log is organized chronologically. Entries are added to the end. Entries are never inserted in the middle, never reordered, and never modified. If an entry contains an error, a correction entry is appended that references the original.

The log is backed up with the same frequency and rigor as all Tier 1 institutional records. It is included in every backup set. It is the first file verified during backup integrity checks.

### 4.2 The Entry Lifecycle

A decision log entry goes through three states:

**Draft:** The entry is being written. It may be incomplete. It is not yet part of the official log. Drafts exist in a designated drafts directory, not in the log itself.

**Committed:** The entry is complete and has been appended to the log. It is now part of the permanent record. It cannot be modified. The timestamp of commitment is recorded.

**Referenced:** The entry has been cited by another entry (a follow-up decision, a correction, a review note). Referenced entries form chains that tell the story of how decisions evolved.

## 5. Rules & Constraints

### 5.1 The Decision Log Entry Format

Every entry in the decision log follows this format:

```
===== DECISION LOG ENTRY =====
Decision ID:    DL-YYYY-NNN
Date:           YYYY-MM-DD
Time:           HH:MM (24-hour, local time)
Operator:       [operator identifier]
Tier:           [1 / 2 / 3 / 4]
Category:       [see Section 5.2]
Tags:           [see Section 5.3]
Status:         [DECIDED / DEFERRED / SUPERSEDED / CORRECTED]
References:     [IDs of related entries, or "None"]
Waiting Period: [start date - end date, or "N/A" for Tier 4]

SUMMARY:
[One sentence. What was decided. Clear enough to be understood
 without reading the full entry.]

CONTEXT:
[What situation prompted this decision? What problem were you
 solving? What was the state of affairs that made a decision
 necessary? Write enough that a person who was not present can
 understand why this moment required a decision.]

ALTERNATIVES CONSIDERED:
[List every alternative you considered, including "do nothing."
 For each alternative, explain why it was rejected. Be specific.
 "It seemed worse" is not a sufficient rejection rationale.
 "It required a proprietary format that would violate R-ETH-04"
 is sufficient.]

DECISION:
[The decision itself. Be precise. State exactly what will be
 done, when, and by whom.]

RATIONALE:
[Why this alternative was chosen over the others. This is the
 most important field in the entry. Write it as though you are
 explaining to a future version of yourself who has forgotten
 everything about this situation. Because that is exactly what
 will happen.]

IMPACT:
[What systems, documents, procedures, or data are affected by
 this decision? What will change as a result?]

REVERSIBILITY:
[Can this decision be undone? At what cost? What would need to
 happen to reverse it? If it is irreversible, say so and explain
 why the irreversibility was accepted.]

EXPECTED OUTCOME:
[What do you expect to happen as a result of this decision?
 This is recorded so that the annual review can compare expected
 outcomes against actual outcomes.]

REVIEW DATE:
[When should this decision be reviewed for continued
 appropriateness? Default: one year from today for Tier 3,
 two years for Tier 2, five years for Tier 1.]

===== END ENTRY =====
```

### 5.2 Decision Classification Categories

Decisions are classified by category to enable filtered searching. The categories are:

- **ARCHITECTURE:** Decisions about system design, infrastructure, data structures, or technical architecture.
- **GOVERNANCE:** Decisions about rules, policies, procedures, authority, or succession.
- **SECURITY:** Decisions about security measures, threat responses, access controls, or incident handling.
- **OPERATIONS:** Decisions about daily operations, maintenance schedules, operational tempo, or workflow.
- **CONTENT:** Decisions about what content to import, create, organize, or retire.
- **QUALITY:** Decisions about quality standards, audit procedures, or quality findings.
- **DOCUMENTATION:** Decisions about documentation structure, format, conventions, or revisions.
- **PERSONNEL:** Decisions about succession, onboarding, training, or role assignment.
- **EMERGENCY:** Decisions made under emergency conditions with compressed or waived waiting periods.

A decision may belong to multiple categories. List all applicable categories, primary category first.

### 5.3 Tagging Conventions

Tags are free-form keywords that supplement the formal classification. They are the primary mechanism for finding related decisions across categories. Conventions:

- Tags are lowercase, hyphenated for multi-word terms (e.g., `backup-strategy`, `format-migration`, `quarantine-rebuild`).
- Every entry should have at least two tags and no more than ten.
- Use existing tags when possible. Before creating a new tag, search the log to see if a synonymous tag already exists.
- Maintain a tag index (a separate file listing all tags in use, with a brief description and the count of entries using each tag). Update the index when committing an entry with a new tag.
- Recommended tags for common decisions: `hardware`, `software`, `backup`, `encryption`, `format`, `import`, `succession`, `documentation`, `budget`, `schedule`, `recovery`, `audit`, `federation`.

### 5.4 Recording Decisions in Real-Time

The practical workflow for capturing a decision:

**Step 1: Recognize the decision.** Not every action is a decision. A decision exists when you choose between alternatives that affect the institution beyond the immediate task. Replacing a failed hard drive with an identical one is a maintenance action (logged in the operational log). Choosing to replace a failed hard drive with a different model that changes the storage architecture is a decision (logged in the decision log).

**Step 2: Draft immediately.** Open the drafts directory. Create a new file using the entry format. Fill in what you can while the context is fresh. At minimum, fill in: Date, Tier, Summary, Context, and Decision. The other fields can be completed within 24 hours.

**Step 3: Complete within 24 hours.** Return to the draft and fill in all remaining fields. The rationale field is the priority -- context fades fastest. If you find yourself unable to articulate the rationale, that is a signal that the decision may not have been well-considered. Re-evaluate if possible.

**Step 4: Review and commit.** Read the complete entry. Ask: would a person who was not present understand what was decided and why? If yes, append the entry to the log. If no, revise until the answer is yes.

**Step 5: Update references.** If this decision relates to a previous decision, add the reference to both entries (the new entry's References field, and a new correction/annotation entry referencing the old one if the old entry did not anticipate this follow-up).

**Step 6: Update the tag index.** If you used any new tags, add them to the tag index with descriptions.

### 5.5 Searching the Decision Log

The decision log is designed to be searched. The three primary search strategies:

**Chronological browsing.** When you need to understand the sequence of decisions around a particular time period, browse the log chronologically. The yearly file organization makes this straightforward.

**Tag-based search.** When you need to find all decisions related to a topic, search by tag. The tag index provides a starting point. Then search the log for the tag string to find all entries.

**Full-text search.** When you need to find a specific decision and are not sure of its tags, search the full text. The structured format means you can narrow searches to specific fields (e.g., search only within RATIONALE fields for a specific concept).

**Decision chain tracing.** When you find a decision that references other decisions, follow the chain. Read the referenced entries. Then search for entries that reference the decision you found. This builds the complete context around a topic.

### 5.6 The Annual Decision Log Review

The annual review of the decision log (part of the annual operational cycle per OPS-001) follows this procedure:

**Step 1: Completeness check.** Review the operational log for the past year. For each significant event or change recorded in the operational log, verify that a corresponding decision log entry exists. If a decision was made but not logged, create a RECONSTRUCTED entry now, marking it clearly.

**Step 2: Review date check.** List all entries whose review date falls within the past year. For each:
- Is the decision still appropriate? If yes, extend the review date.
- Has the decision been superseded? If yes, create a SUPERSEDED entry referencing the new decision.
- Has the expected outcome matched the actual outcome? Record the comparison in an annotation entry.

**Step 3: Pattern analysis.** Read through the year's entries looking for patterns:
- Are the same types of decisions being made repeatedly? (This may indicate a need for a standing policy rather than repeated ad hoc decisions.)
- Are decisions being made at the wrong tier? (Too many Tier 4 decisions that should have been Tier 3?)
- Are rationales getting shorter or more perfunctory? (Sign of governance fatigue.)

**Step 4: Tag maintenance.** Review the tag index. Remove tags that are no longer in use (or mark them as historical). Consolidate synonymous tags. Ensure all entries from the past year are adequately tagged.

**Step 5: Produce the Annual Decision Log Report.** Summarize:
- Total entries for the year, broken down by tier and category.
- Key decisions and their outcomes.
- Patterns observed.
- Completeness assessment (estimated percentage of significant decisions that were captured).
- Recommendations for the coming year (any changes to the decision-making process itself).

### 5.7 Log Maintenance

**Integrity verification.** At least monthly (during monthly operations), verify the log's integrity:
- Confirm the log file's checksum matches the recorded checksum from the last verification.
- Confirm the entry count matches the expected count.
- Confirm the most recent entry's Decision ID is sequential with no gaps.

**Backup.** The decision log is included in every backup. During backup verification, the decision log is one of the files tested by restoration and checksum comparison.

**Format migration.** If the log's storage format needs to change (e.g., moving from flat text files to a different structure), the migration is a Tier 2 governance decision. The migration must preserve every entry exactly. The old format files are retained as archival copies. The migration process is documented in the decision log itself.

**Year-end rollover.** At the end of each calendar year, close the current year's log file and open a new one for the coming year. The previous year's file becomes read-only. The Decision ID sequence continues across years without resetting (DL-2026-001, DL-2026-002, ..., DL-2027-001, ...).

## 6. Failure Modes

- **Logging fatigue.** The operator grows tired of writing decision log entries. Entries become perfunctory -- short summaries with no rationale, no alternatives, no expected outcomes. The log exists but is useless for context recovery. Mitigation: the entry format is designed to be completable in 15-20 minutes for a typical Tier 3 decision. If entries consistently take longer, the format may need simplification. The annual review checks for declining entry quality.

- **Tier misclassification.** Decisions are systematically classified at the wrong tier -- usually too low, to avoid the waiting period requirements of higher tiers. Mitigation: the annual review includes a tier audit. If a pattern of misclassification is detected, document it and recalibrate. The tier definitions in GOV-001 are the authoritative reference.

- **Tag entropy.** Tags proliferate without discipline. Multiple tags mean the same thing. The tag index is not maintained. Searching by tag becomes unreliable. Mitigation: the tag maintenance step in the annual review. Consolidate synonymous tags. Maintain the tag index actively.

- **Decision log corruption.** The log file is accidentally modified (a stray keystroke while reading it) or corrupted by a storage failure. Mitigation: integrity verification during monthly operations detects modifications. Backup provides recovery. The log file should be stored with write-protect attributes on the operating system level and only temporarily made writable during the commit operation.

- **Over-logging.** Everything is logged, including trivial decisions that clutter the log and bury significant decisions in noise. Mitigation: the criterion in Step 1 of Section 5.4 -- a decision exists when you choose between alternatives that affect the institution beyond the immediate task. Routine actions that follow established procedures are operational log entries, not decision log entries.

- **Under-logging.** Significant decisions are not logged because the operator did not recognize them as decisions at the time. Mitigation: the annual completeness check (Step 1 of the annual review) compares the operational log against the decision log to detect gaps. Train yourself to recognize decisions by their characteristics: you chose between alternatives, the choice has lasting implications, and a future operator would benefit from knowing why you chose as you did.

## 7. Recovery Procedures

1. **If entries have been missed:** Conduct a gap analysis by comparing the operational log against the decision log for the period in question. For each missing decision, create a RECONSTRUCTED entry using the best available information. Mark these entries with Status: RECONSTRUCTED and note that they were created after the fact, with the date of reconstruction.

2. **If the decision log is corrupted:** Restore from the most recent verified backup. Compare the restored log against the backup copy. If the corruption affected entries that were created after the most recent backup, recreate them from memory and operational log records. Mark recreated entries as RECONSTRUCTED. Investigate and document the cause of corruption.

3. **If the tag index is lost or corrupted:** Rebuild it by scanning all decision log entries and extracting all tags. This is tedious but mechanical. The log itself is the authoritative source; the tag index is a convenience that can always be regenerated.

4. **If logging fatigue has led to a period of low-quality entries:** Do not retroactively revise the low-quality entries (that would violate the append-only rule). Instead, create annotation entries for the most significant decisions from the low-quality period, adding the rationale and context that the original entries lacked. These annotations reference the original entries and provide the missing context. Resume rigorous logging going forward.

5. **If the decision log format needs to change:** Follow the format migration procedure in Section 5.7. This is a Tier 2 governance decision. Document the migration extensively. Retain the old format as an archive. Verify that no entries were lost or altered during migration.

## 8. Evolution Path

- **Years 0-5:** The decision log is being populated for the first time. The first hundred entries will establish the operator's logging habits and the practical utility of the format. Expect the format to feel awkward at first and natural after several months. Resist the temptation to abbreviate entries during this establishment period -- the full-format habit is easier to relax later than to establish later.

- **Years 5-15:** The log contains years of decisions. Its value is becoming apparent -- the operator can search for precedents, trace decision chains, and compare expected outcomes against actual outcomes. The annual review is substantive, not pro forma. The tag system should be well-established and the tag index actively maintained.

- **Years 15-30:** The decision log is a substantial historical document. A successor can read it and understand not just what the institution is, but how it became what it is. The log's durability is tested as storage media are migrated and file formats potentially change. The plain-text format choice pays dividends here -- it survives media and format transitions with minimal risk.

- **Years 30-50+:** The decision log may span thousands of entries across decades. Search becomes important. The tag system and the annual review summaries serve as indices into the full log. The log is the institution's autobiography, written one decision at a time. Protect it accordingly.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The first decision log entry I wrote took forty-five minutes. The tenth took fifteen. The format felt bureaucratic at first. By the twentieth entry, it felt natural -- like a conversation with my future self. "Here is what I decided. Here is why. Here is what I expected. Check back in a year and see if I was right."

The hardest field to fill in honestly is ALTERNATIVES CONSIDERED. The temptation is to list only the alternative you chose and one obvious bad alternative, making the decision look inevitable. Resist this. Most real decisions involve multiple viable alternatives, and the reasons for choosing one over another are precisely what future operators need to understand. If you cannot articulate why you rejected an alternative, you may not have actually considered it -- and that itself is worth noting.

I have set the default review dates conservatively (one year for Tier 3, two years for Tier 2, five years for Tier 1). These may need adjustment as operational experience accumulates. If the annual review consistently finds that decisions are still appropriate at their review dates, the intervals can be extended. If decisions are frequently found to be outdated at their review dates, the intervals should be shortened.

## 10. References

- D20-001 -- Institutional Memory Philosophy (the principles this procedure implements)
- GOV-001 -- Authority Model (decision tiers; decision record requirements; Section 4.2)
- ETH-001 -- Ethical Foundations (Principle 3: Transparency of Operation; Principle 6: Honest Accounting)
- CON-001 -- The Founding Mandate (purpose amnesia as failure mode)
- OPS-001 -- Operations Philosophy (operational tempo; quarterly and annual review cycles)
- SEC-001 -- Threat Model (data integrity threats; operational security)
- D20-ART-002 -- Decision Log Technical Format and Schema
- D20-ART-003 -- Oral History Capture Procedures
- D20-ART-004 -- Lessons-Learned Framework
- D19-002 -- Documentation Quality Audit Procedures (audit of decision log entries)

---

---

# D4-001 -- Infrastructure Philosophy and Physical Site Requirements

**Document ID:** D4-001
**Domain:** 4 -- Infrastructure & Power
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, GOV-001, SEC-001, OPS-001
**Depended Upon By:** All articles in Domain 4. All articles in Domain 5 (platform runs on infrastructure). D12-001 (disaster recovery depends on infrastructure resilience). SEC-002 (physical access controls).

---

## 1. Purpose

This article defines the infrastructure philosophy and physical site requirements for the holm.chat Documentation Institution. It specifies the physical foundation upon which the entire institution rests: how power is generated and delivered, how computing systems are interconnected (or deliberately not interconnected), how the physical environment is controlled, how the site is physically secured, and how all of these systems are maintained over a lifetime.

The institution is software running on hardware in a physical place. The Stage 2 philosophy articles and Stage 3 operational manuals address the software and the procedures. This article addresses the physical reality: the building, the power, the cooling, the cables, the locks, and the concrete. Without a sound physical foundation, every other article in the institution is theoretical.

This article is written for the person who will select a site, install power systems, wire networks, configure environmental controls, and maintain all of these for decades. It assumes technical competence but not specialist expertise in electrical engineering, HVAC, or physical security. Where specialist knowledge is required, this article says so and provides enough context for the operator to engage a specialist effectively.

## 2. Scope

This article covers:

- Power systems: primary, secondary, and emergency power. Generation, storage, conditioning, and distribution.
- Network topology: the internal network architecture of the air-gapped institution. Cabling, switching, segmentation.
- Environmental controls: temperature, humidity, dust, and water management.
- Physical security: locks, barriers, surveillance, and access control at the physical level.
- Site selection criteria: what to look for and what to avoid in a physical location.
- Redundancy principles: single points of failure and how to eliminate them.
- Maintenance schedules: what to inspect, test, and replace, and when.

This article does not cover:

- Software platform and operating system choices (see Domain 5).
- Logical network configuration (there is no network to the outside; internal logical configuration is in Domain 5).
- Backup media storage and rotation (see Domain 12).
- Disaster recovery procedures (see D12-001; this article provides the infrastructure that D12-001 depends on).
- Specific hardware procurement and lifecycle management (see D4-ART-003).

## 3. Background

CON-001 establishes three architectural mandates: air-gapped, off-grid capable, and self-built. This article is where those mandates meet physical reality.

**Air-gapped** means no network connection to the outside world. But the institution still has an internal network -- the machines that compose it must communicate with each other for services like backup, monitoring, and data access. The network topology in this article defines an internal-only network that is physically incapable of reaching the outside world. This is not a firewall rule. It is a physical fact: there is no cable, no wireless radio, no pathway by which a packet could leave the institutional network.

**Off-grid capable** does not mean the institution never uses grid power. It means the institution can survive the loss of grid power for a sustained period. The power architecture in this article provides for primary power (which may be grid, solar, or other), secondary power (which covers primary power failures), and emergency power (which covers the failure of both primary and secondary). The goal is not perpetual off-grid operation (though that is possible with the right site and power systems). The goal is resilience: the institution does not die when the grid fails.

**Self-built** means the operator understands every physical system. This does not mean the operator personally welded every rack or crimped every cable. It means the operator understands the design of every system, can explain why it is configured the way it is, and can maintain or replace it. Where a specialist was hired to install a system, the operator watched, learned, and documented. The operator is not dependent on any specialist's continued availability.

## 4. System Model

### 4.1 The Power Architecture

Power is delivered through a three-tier architecture:

**Tier 1: Primary Power.**
The main power source for the institution during normal operations. Options include:

- **Grid power:** The simplest option. Reliable in most locations. The institution connects to the local electrical grid through a standard service entrance. Advantages: continuous, high capacity, low cost per kilowatt-hour. Disadvantages: the operator does not control it; outages are outside the operator's control; dependency on an external entity.
- **Solar power with battery storage:** The most common off-grid option. Photovoltaic panels charge a battery bank, which provides power to the institution. Advantages: no external dependency; renewable; low ongoing cost. Disadvantages: high upfront cost; dependent on sunlight; battery banks have finite lifespans (typically 10-15 years for lithium, 5-8 years for lead-acid); requires more physical space.
- **Hybrid (grid + solar):** Grid power as primary with solar as supplemental, or solar as primary with grid as backup. This is the recommended configuration where grid power is available, because it provides the reliability of grid power with the resilience of solar.

Regardless of source, primary power must be conditioned before reaching computing equipment:
- Surge protection at the service entrance.
- Uninterruptible power supply (UPS) at the equipment level to handle brief interruptions and provide clean power.
- Power distribution that allows individual circuits to be isolated for maintenance.

**Tier 2: Secondary Power.**
Activated when primary power fails. The secondary source must be capable of powering the institution for at least 72 hours without intervention. Options include:

- **Battery bank (if primary is grid):** A dedicated battery bank, separate from any solar storage, that is kept charged by primary power and activated on primary failure. Sizing: calculate the institution's total power draw and multiply by 72 hours. Add a 25% margin.
- **Generator (diesel, propane, or natural gas):** Automatically or manually started on primary failure. Fuel storage must be sufficient for at least 72 hours of operation. The generator requires regular testing and maintenance (see Section 5.7).
- **Solar (if primary is grid):** Solar panels and battery storage serve as secondary power. This works well in combination with grid primary power.

**Tier 3: Emergency Power.**
Activated when both primary and secondary power fail. Emergency power is not intended to run the institution normally. It is intended to power a controlled shutdown -- saving all data, completing in-progress writes, and bringing systems to a safe halt state. Emergency power is typically a dedicated UPS or battery that powers only the essential systems (main server, primary storage, and the shutdown orchestration system) for 15-30 minutes.

### 4.2 The Network Topology

The institutional network is internal only. It connects the institution's machines to each other and to nothing else.

**Physical topology:**

- All network connections are wired (Ethernet). No wireless access points. No Bluetooth. No infrared. Wireless radios are physically disabled (removed from hardware where possible, disconnected at the hardware level where removal is not possible). This is per SEC-001, R-SEC-01: physical enforcement, not software disabling.
- The network uses a single flat switch or a small number of switches, depending on the number of machines. There is no router (there is nothing to route to). There is no gateway (there is no outside).
- Cable runs are documented: which port connects to which machine, the cable path, the cable type, and the cable length. A network diagram is maintained and updated whenever connections change.

**Logical topology:**

- The network uses a private address space (e.g., 10.0.0.0/24).
- Each machine has a static IP address. No DHCP (in a static, known environment, dynamic addressing adds complexity without benefit).
- DNS, if used, is provided by a local DNS server with no external forwarders. Alternatively, a static hosts file on each machine maps hostnames to IP addresses.
- Network segmentation between functional zones is recommended but not required for small installations. If the institution grows beyond five machines, segment the network into at least two zones: operational (daily-use machines and services) and storage (backup servers, archive storage). The quarantine machine is never on this network.

**The quarantine machine exception:**

The quarantine machine (D18-002) is not on the institutional network. It is not connected to any other machine by any means. It is a standalone system. Content moves from the quarantine machine to the institutional network only via the dedicated transfer medium (a USB drive carried by the operator). This is the physical implementation of the quarantine isolation requirement in D18-001.

### 4.3 Environmental Controls

Computing equipment operates reliably within specific environmental parameters. Exceeding these parameters causes immediate failure or accelerated degradation. The institution must maintain:

**Temperature:** 18-27 degrees Celsius (64-80 degrees Fahrenheit) in the equipment area. The ideal target is 20-22 degrees Celsius. Higher temperatures reduce equipment lifespan and increase fan noise. Lower temperatures are acceptable but waste energy.

**Humidity:** 40-60% relative humidity in the equipment area. Below 30% risks static discharge, which can damage equipment. Above 70% risks condensation and corrosion.

**Dust:** Minimized. Computing equipment draws air through fans, and that air carries dust. Dust accumulates on heat sinks, reduces cooling efficiency, and accelerates wear. Mitigation: air filtration on any ventilation intake serving the equipment area; regular cleaning of equipment (quarterly, using compressed air and appropriate cleaning tools); sealed cable pathways to reduce dust intrusion.

**Water:** None. No water pipes should run through or above the equipment area. If the building's plumbing makes this impossible, install drip trays above any equipment under overhead pipes, and install a water detection sensor that alerts the operator (an audible alarm, since there is no network alert pathway to the outside).

Environmental monitoring is continuous. At minimum, the equipment area has:
- A temperature and humidity sensor with logging capability.
- A water detection sensor near the floor.
- A smoke detection sensor.
- Readings are logged and reviewed during the daily operational check.

### 4.4 Physical Security

Physical security is the institution's perimeter defense. In an air-gapped institution, physical access equals total access (given time and skill). The physical security model has three layers:

**Layer 1: The Building.**
The building should be structurally sound, with controlled entry points. Doors and windows should have locks. The entry to the equipment area should not be visible or obvious from outside the building. The building should be in a location where unusual activity (unfamiliar vehicles, unfamiliar people) is noticeable.

**Layer 2: The Equipment Room.**
The room containing the institution's hardware should have:
- A reinforced door with a deadbolt lock.
- No windows, or windows with security film and locks.
- A separate key or combination from the building's general entry.
- An access log: a physical logbook at the door recording every entry (date, time, person, purpose). In a single-operator institution, this feels redundant. It is not. It creates a record that detects unauthorized access, and it creates a habit that persists through succession.

**Layer 3: The Equipment.**
Individual equipment items should be physically secured:
- Server chassis should be locked or require a tool to open.
- Storage media (backup drives, archive media) should be stored in a locked cabinet or safe.
- The quarantine machine should be physically separate from production equipment (ideally in a different part of the room or a different room entirely).

### 4.5 Site Selection Criteria

When selecting or evaluating a site for the institution, consider:

**Must-have criteria:**
- Structurally sound building with no known flood risk, landslide risk, or severe geological instability.
- Reliable primary power available (grid, solar capacity, or both).
- A room suitable for the equipment area: climate-controllable, no overhead water pipes, sufficient electrical capacity, secure.
- Physical security: the site can be locked and monitored. Unauthorized entry is detectable.
- Sufficient physical space for current equipment plus at least 50% growth.

**Strong preferences:**
- Low humidity climate or effective dehumidification available.
- Distance from major electromagnetic interference sources (high-voltage power lines, radio transmitters, industrial equipment).
- Geographic separation from backup storage locations (see Domain 12).
- Stable temperature -- mild climates reduce the energy cost of environmental control.
- Low crime area to reduce physical security risk.

**Avoid:**
- Flood plains. Underground installations in areas with high water tables.
- Areas with frequent, prolonged power outages (if grid power is planned as primary).
- Shared spaces where non-institutional personnel have regular access.
- Buildings with known structural issues (asbestos, compromised foundations, inadequate electrical wiring).

## 5. Rules & Constraints

- **R-D4-01: No Single Point of Power Failure.** The institution must have at least two independent power sources (primary and secondary). The loss of any single power source must not cause data loss. Emergency power must ensure graceful shutdown even when primary and secondary both fail.

- **R-D4-02: Physical Air Gap Enforcement.** No physical cable, wireless radio, or other communication pathway may connect any institutional system to any non-institutional system. The quarantine machine is not on the institutional network. Wireless hardware is physically disabled, not software-disabled.

- **R-D4-03: Environmental Monitoring Is Continuous.** Temperature, humidity, and water detection sensors must be operational at all times in the equipment area. Sensor readings must be logged. Anomalous readings must trigger investigation within 24 hours.

- **R-D4-04: Physical Access Is Logged.** Every entry to the equipment room must be recorded (date, time, person, purpose). The access log is a Tier 2 institutional record (active for as long as the site is in use, then retained for the institution's lifetime).

- **R-D4-05: Network Topology Is Documented.** A current network diagram showing all physical connections, all IP addresses, and all cable paths must be maintained. The diagram is updated within 24 hours of any network change.

- **R-D4-06: Redundancy Where Failure Is Catastrophic.** Any infrastructure component whose failure would cause data loss or extended downtime must have a documented redundancy plan: either a hot spare, a cold spare, or a documented procurement path that can be executed within an acceptable downtime window.

- **R-D4-07: Maintenance Is Preventive.** Infrastructure maintenance follows a preventive schedule (see Section 5.7), not a reactive approach. Components are inspected and tested before they fail, not after.

### 5.7 Maintenance Schedules

**Daily (during daily operations):**
- Check environmental sensor readings (temperature, humidity).
- Visual inspection of the equipment area (no water leaks, no unusual lights on equipment, no unusual sounds).
- Verify UPS status indicators (online, battery health).

**Monthly (during monthly operations):**
- Clean equipment air filters if equipped.
- Check all cable connections for physical integrity (no loose connectors, no cable damage).
- Test the water detection sensor (apply a small amount of water to the sensor pad; verify alarm triggers; dry immediately).
- Review environmental logs for the past month. Flag any readings outside the acceptable range.
- Check UPS battery health via the UPS management interface. Record battery voltage, charge level, and estimated runtime.

**Quarterly (during quarterly operations):**
- Clean internal equipment (compressed air to remove dust from heat sinks, fans, and circuit boards). Power down equipment before cleaning.
- Inspect physical security: test all locks, check door frames and hinges, inspect window security.
- Test secondary power: simulate a primary power failure (disconnect primary power) and verify that secondary power activates and supports the institution for at least 15 minutes. Restore primary power. Document the test.
- Review and update the network diagram if any changes occurred.

**Annually (during annual operations):**
- Full infrastructure review:
  - Assess the health and remaining useful life of all hardware components. Create a replacement timeline for components approaching end-of-life (see D4-ART-003 for hardware lifecycle management).
  - Assess the capacity of power systems against current and projected demand.
  - Assess the physical site for any new risks (new construction nearby, changes in drainage, changes in neighborhood, structural wear).
- Test emergency power: simulate the failure of both primary and secondary power. Verify that emergency power activates and provides sufficient time for a controlled shutdown. Perform the controlled shutdown. Restart from the shutdown state. Document the entire test.
- Test all environmental sensors by exposing them to out-of-range conditions and verifying they respond correctly.
- Review the physical access log for the past year. Flag any anomalies.
- Solar power systems (if applicable): inspect panels for physical damage, check inverter operation, review annual energy production against expectations, clean panels.
- Battery systems: test battery capacity under load. Compare against original rated capacity. Batteries that have degraded below 70% of rated capacity should be scheduled for replacement.

## 6. Failure Modes

- **Power cascade failure.** Primary power fails. Secondary power activates but fails due to a defect that was not caught by testing (dead batteries, fuel leak in generator, tripped breaker). Emergency power activates but provides only minutes. If the controlled shutdown does not complete in time, data corruption may result. Mitigation: quarterly secondary power tests. Annual emergency power tests. Monthly UPS battery checks. Defense in depth -- each tier is independent.

- **Environmental control failure.** The cooling system fails during a heat wave. Equipment temperatures exceed safe limits. If not detected quickly, hardware damage occurs. Mitigation: continuous temperature monitoring with alarms. Automated thermal shutdown thresholds configured in the BIOS/UEFI of all machines (shutdown at 85 degrees Celsius CPU temperature). Manual procedures for emergency ventilation (open doors, portable fans).

- **Physical security breach.** Unauthorized access to the equipment room. In the worst case, hardware is stolen or tampered with. Mitigation: layered physical security (building, room, equipment). Access logging. If a breach is detected, follow SEC-001 incident response: assume all systems that were physically accessible are compromised until proven otherwise.

- **Network isolation breach.** A wireless radio that was supposed to be disabled is re-enabled (by a BIOS update, a firmware change, or an operator error). The air gap is silently compromised. Mitigation: R-D4-02 requires physical disabling, not software disabling. During the quarterly infrastructure review, verify that all wireless hardware remains physically disabled.

- **Cable degradation.** Over years, cables degrade -- insulation cracks, connectors corrode, strain relieves fail. Network performance degrades silently until connections become intermittent. Mitigation: monthly cable inspection. Replace cables at the first sign of physical degradation. Maintain spare cables of every type in use.

- **Site becomes unsuitable.** Over decades, the site's suitability may change: new flood risk due to climate or development changes, structural deterioration, neighborhood changes affecting security. Mitigation: annual site assessment. If the site's risk profile changes significantly, begin planning for site migration (a Tier 1 governance decision given the magnitude).

## 7. Recovery Procedures

1. **If power is lost at all three tiers simultaneously:** This is a worst-case scenario. If equipment was running, data corruption is possible. When power is restored: do not boot systems immediately. Inspect equipment for physical damage (melted components, tripped breakers, error indicators). Boot systems one at a time, starting with the primary storage server. Run filesystem integrity checks before mounting any data volumes. If corruption is detected, restore from verified backup. Document the incident completely, including the cause of the three-tier failure and the measures taken to prevent recurrence.

2. **If environmental controls fail and equipment overheats:** If equipment has not yet shut down: perform an immediate controlled shutdown of all systems. Open doors and windows for ventilation. Do not restart until the room temperature is within the acceptable range and the cooling system has been repaired or replaced. If equipment has already shut down (thermal protection): allow the room to cool. Inspect equipment before restarting. Test all systems before returning to normal operations.

3. **If a physical security breach is detected:** Secure the site immediately. Photograph the scene. Do not touch or modify anything until an assessment is complete. Check all equipment for signs of tampering (opened chassis, disconnected cables, unfamiliar devices connected). If tampering is suspected, follow SEC-001 incident response: assume compromised until proven otherwise. If hardware was removed, report the theft to appropriate authorities and activate the disaster recovery procedures in D12-001.

4. **If the network is found to have an unintended external connection:** Disconnect the connection immediately. Determine how it occurred. Assess what data may have been exposed. Follow SEC-001 air-gap breach procedures: treat all connected systems as potentially compromised. Rebuild from known-good state if necessary.

5. **If a site migration is required:** This is a major undertaking. Follow GOV-001 Tier 1 decision process. Plan the migration in phases: select the new site (using the criteria in Section 4.5), prepare the new site's infrastructure, migrate backup copies first, verify at the new site, then migrate production systems. Maintain the old site as a backup until the new site is fully verified. Document the migration extensively in the decision log.

## 8. Evolution Path

- **Years 0-5:** The infrastructure is being established. Expect to discover inadequacies in the power system, the cooling, or the physical layout. The first year is a learning year for the site -- you will learn its quirks, its failure patterns, and its strengths. Document everything. The maintenance schedule is being calibrated against real experience.

- **Years 5-15:** The infrastructure should be stable. Hardware replacements will begin as the first generation of equipment ages. The power system's performance characteristics should be well-understood. Battery replacements will be needed at least once during this period. The focus shifts from building infrastructure to maintaining it -- which, per OPS-001, is the harder and more important task.

- **Years 15-30:** Major infrastructure evolution may be needed. Hardware generations will have changed significantly. Power technology may have evolved (battery technology, in particular, improves rapidly). The network topology may need expansion if the institution's scope grows. Site suitability should be reassessed given any changes in the local environment.

- **Years 30-50+:** The site may need to be migrated to a new location (buildings age, neighborhoods change, the operator's life circumstances change). The infrastructure philosophy -- three-tier power, air-gapped internal network, environmental monitoring, physical security -- survives any site transition. The specific implementations will change entirely. The principles will not.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The physical foundation is the part of this institution that is easiest to neglect because it is the least intellectually interesting. Writing code is more engaging than crimping Ethernet cables. Designing data structures is more stimulating than sizing a battery bank. But the institution runs on power, not on ideas. If the power fails and the cooling fails and the hardware fails, every elegant article in this corpus is nothing more than text on a dead disk.

I have structured the power architecture as three tiers because I have seen too many "backup power" systems that consist of a single UPS that has not been tested since it was installed. A UPS that has not been tested is not backup power. It is a heavy box that makes you feel safe. Test your power. Test your cooling. Test your locks. Test everything. Then test it again next quarter.

The air gap's physical enforcement is worth emphasizing. A software-disabled WiFi radio is one firmware update away from being re-enabled. A WiFi card that has been physically removed from the motherboard cannot be re-enabled by any software. Physical enforcement is more work during initial setup. It is absolute thereafter. In security, absoluteness has a value that convenience does not.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 2: Integrity Over Convenience; Principle 4: Longevity Over Novelty)
- CON-001 -- The Founding Mandate (air-gapped, off-grid, self-built mandates)
- GOV-001 -- Authority Model (decision tiers for infrastructure changes)
- SEC-001 -- Threat Model and Security Philosophy (physical threats; air-gap enforcement; R-SEC-01)
- OPS-001 -- Operations Philosophy (maintenance schedules; operational tempo; sustainability requirement)
- D18-002 -- Content Import Pipeline (quarantine machine isolation requirements)
- D12-001 -- Disaster Recovery Philosophy (infrastructure resilience as a prerequisite)
- SEC-002 -- Access Control Procedures (physical access controls)
- D4-ART-003 -- Hardware Procurement and Lifecycle Management
- Stage 1 Documentation Framework, Domain 4: Infrastructure & Power

---

---

*End of Stage 3 Operational Doctrine -- Batch 3*

**Document Total:** 5 articles
**Domains Covered:** 18 (Import & Quarantine), 19 (Quality Assurance), 9 (Knowledge & Training), 20 (Institutional Memory), 4 (Infrastructure & Power)
**Combined Estimated Word Count:** ~15,000 words
**Status:** All five articles ratified as of 2026-02-16.
**Next Stage:** Additional Stage 3 operational articles for remaining domains, and Stage 4 reference specifications.
