# STAGE 3: DATA & ARCHIVES OPERATIONS

## Operational Doctrine for Domain 6 -- Data & Archives

**Document ID:** STAGE3-DATA-OPS
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Classification:** Stage 3 -- Operational Doctrine. Practical manuals for the stewardship, organization, migration, verification, and lifecycle management of institutional data. These articles translate the philosophy established in D6-001 into concrete, executable procedures.
**Depends On:** ETH-001, CON-001, GOV-001, SEC-001, OPS-001, D6-001, D6-002, D6-003

---

## How to Read This Document

This document contains five operational articles for Domain 6: Data & Archives. They are numbered D6-004 through D6-008. They are not philosophical documents. They are manuals. They tell you how to organize archives, how to migrate formats, how to verify backups, how to manage metadata, and how to manage storage media across its entire lifecycle.

Each article traces its authority to D6-001 (Data Philosophy), which establishes the principles of data sovereignty, the four-tier data hierarchy, and the ethic of stewardship that governs every decision about institutional data. Where a procedure in this document conflicts with a principle in D6-001, the principle prevails and this document must be amended.

These articles were written simultaneously and cross-reference one another. Archive management (D6-004) depends on metadata standards (D6-007). Format migration (D6-005) depends on storage lifecycle awareness (D6-008). Backup verification (D6-006) depends on the archive structures defined in D6-004. Read them in order the first time. After that, consult each as needed.

If you are a future maintainer: these procedures were written for the technology available at founding. The principles behind each procedure are stated in the Background section. When the technology changes -- and it will -- adapt the procedures while preserving the principles. Then update this document. A procedure that no longer matches reality is worse than no procedure at all, because it creates false confidence.

---

---

# D6-004 -- Archive Management Procedures

**Document ID:** D6-004
**Domain:** 6 -- Data & Archives
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, SEC-001, OPS-001, D6-001, D6-002, D6-003, D6-007
**Depended Upon By:** D6-005, D6-006, D6-008, D6-009, D6-011, D6-014. All articles involving data storage, retrieval, or audit.

---

## 1. Purpose

This article defines how archives are organized, maintained, accessed, and audited within the holm.chat Documentation Institution. An archive is not merely a collection of old files. It is the physical and logical embodiment of the institution's memory -- the structure through which decades of data remain findable, readable, and meaningful long after the person who created them has forgotten the context or is no longer present to explain it.

The distinction between "storing data" and "maintaining an archive" is the distinction between a box of unsorted papers in an attic and a library with a catalogue. Both contain information. Only one allows retrieval. This institution is designed to persist for fifty years or more. Without disciplined archive management, data accumulates into an undifferentiated mass that grows less useful with every passing year. With disciplined archive management, data becomes more valuable over time because the archive structure itself encodes relationships, context, and retrievability.

This article provides the directory structures, naming conventions, metadata requirements, access procedures, and audit schedules that transform raw storage into a functioning archive. It implements the principles of data stewardship established in D6-001 and the metadata standards defined in D6-007. Every procedure here serves one goal: any competent person with access to the archive and this document should be able to find, understand, and use any piece of data the institution has preserved, without needing to ask anyone for help.

## 2. Scope

This article covers:

- The canonical directory structure for all institutional archives.
- Naming conventions for directories, files, and archive bundles.
- Metadata requirements for archived items, implementing D6-007 at the archive level.
- Access and retrieval procedures: how to find data in the archive and how to extract it for use.
- The archive audit: what is checked, how often, and what to do when problems are found.
- The transition process: how data moves from active storage to the archive.
- Archive integrity verification as it relates to the broader integrity framework in D6-006.

This article does not cover:

- The philosophy of what to keep and why (see D6-001).
- Specific file formats for archival storage (see D6-003).
- Backup procedures for the archive itself (see D6-006).
- Metadata field definitions and controlled vocabularies (see D6-007).
- Storage media selection and health monitoring (see D6-008).
- Retention schedules and disposal procedures (see D6-011).

## 3. Background

### 3.1 Why Structure Matters More Than Storage

Storage is cheap and getting cheaper. Structure is expensive and getting more expensive, because it requires human judgment. An institution that buys more disks when storage fills is spending money. An institution that maintains a coherent archive structure as data grows is spending discipline. Over fifty years, the disciplined institution will be able to find a specific document in minutes. The undisciplined institution will have petabytes of data and no idea what is in them.

The history of digital archives is a history of structural failure. Organizations that stored everything without organizing it discovered, years later, that their data was functionally inaccessible -- not because the bits had degraded, but because nobody could determine what was where or what it meant. The physical media were intact. The logical structure was absent. This article exists to prevent that outcome.

### 3.2 The Air-Gap Archival Constraint

In an internet-connected environment, search engines and indexing services compensate for poor archive structure. You can find a needle in a haystack if you have a sufficiently powerful magnet. In an air-gapped institution, the archive structure is the search engine. There is no external indexing service. There is no cloud-based full-text search. The directory hierarchy, the naming conventions, and the metadata catalogue are the only tools available for locating data. They must be good enough to work without computational brute force.

### 3.3 Archives as Institutional Infrastructure

The archive is not a passive repository. It is active infrastructure. The governance system references archived decisions. The intelligence function queries archived assessments. The training system draws on archived materials. The operational logs form an archived record of institutional activity. Every domain in this institution depends on the archive functioning correctly. Archive management is therefore not a maintenance task relegated to spare time. It is core infrastructure maintenance, equivalent in importance to keeping the servers running.

## 4. System Model

### 4.1 The Canonical Directory Structure

All institutional data resides within a single root archive path. Below this root, the structure is organized by tier, then by domain, then by year, then by category. The structure is:

```
/archive/
  tier1-institutional/
    governance/
      decisions/
        YYYY/
      amendments/
        YYYY/
      reviews/
        YYYY/
    documentation/
      domain-NN/
        YYYY/
    logs/
      operational/
        YYYY/
        YYYY-MM/
      security/
        YYYY/
      system/
        YYYY/
  tier2-operational/
    projects/
      [project-name]/
        YYYY/
    configurations/
      [system-name]/
        current/
        history/
          YYYY/
    maintenance/
      YYYY/
  tier3-reference/
    research/
      [topic]/
    imported/
      YYYY/
      YYYY-MM/
    media/
      [collection]/
  tier4-transient/
    scratch/
    drafts/
    temp/
```

Rules governing this structure:

- **The root path is a configuration parameter.** It is documented in the system configuration and may differ between machines. All procedures reference `$ARCHIVE_ROOT` rather than a hardcoded path.
- **Year directories use four-digit years.** No two-digit year abbreviations. No ambiguity about century.
- **Month directories, where used, use YYYY-MM format.** Never month names. Never locale-dependent abbreviations.
- **Domain directories use the domain number with zero-padding.** Domain 6 is `domain-06`, not `domain-6`. This ensures correct lexicographic sorting.
- **Project and topic directories use lowercase kebab-case.** Spaces, uppercase letters, and special characters other than hyphens are prohibited in directory names.
- **No directory shall exceed 1,000 immediate children.** If a directory approaches this limit, subdivide it by adding a layer of categorization. An overfull directory is as useless as an unstructured one.

### 4.2 Naming Conventions

Every file in the archive follows a naming pattern that encodes essential context without requiring the catalogue to be consulted:

**Pattern:** `YYYY-MM-DD_[source]_[descriptor]_v[version].[ext]`

- **YYYY-MM-DD:** The date the content was created (not the date it was archived).
- **source:** The originating entity, system, or process. For operator-created documents: `op`. For system-generated data: the system name. For imported data: the external source abbreviated to a maximum of 12 characters.
- **descriptor:** A brief, human-readable description of the content. Maximum 50 characters. Lowercase kebab-case. Must be meaningful to someone who has never seen the file. `meeting-notes` is acceptable. `mn` is not. `untitled` is never acceptable.
- **v[version]:** The version number. For documents that do not have formal versioning, use `v1`. For documents that are updated, increment the version. The archive retains all versions; it does not overwrite.
- **[ext]:** The file extension matching the format. Must correspond to an approved format per D6-003.

**Examples:**

- `2026-02-15_op_quarterly-archive-audit-report_v1.md`
- `2026-01-10_syslog_primary-server-january_v1.txt`
- `2025-12-01_imported_solar-panel-specifications_v1.pdf`

**Prohibited naming practices:**

- Spaces in filenames. Use hyphens.
- Uppercase letters in filenames. Use lowercase exclusively.
- Characters outside ASCII alphanumerics, hyphens, underscores, and periods.
- Generic names: `data`, `file`, `backup`, `new`, `old`, `test`, `final`, `final-v2`.
- Date formats other than ISO 8601 (YYYY-MM-DD).

### 4.3 The Archive Transition Process

Data moves from active storage to the archive when its active use frequency declines. The transition process is:

1. **Identify candidates.** During the monthly archive review (see Section 4.5), identify data in active storage that has not been accessed or modified in 90 days (Tier 2 data) or 30 days (Tier 4 data). Tier 1 data is archived on creation. Tier 3 data follows the 90-day rule.

2. **Verify classification.** Confirm that the data's tier classification per D6-001 is still correct. Data that was Tier 2 when created may now be Tier 3. Reclassify if necessary and document the change.

3. **Verify metadata.** Confirm that all required metadata fields per D6-007 are populated. Missing metadata must be filled before archival. If metadata cannot be determined, document this explicitly in a companion `.meta` file.

4. **Verify format.** Confirm that the data is stored in an approved archival format per D6-003. If not, convert it before archiving. The active-format version may be discarded after successful conversion and verification.

5. **Generate integrity checksums.** Compute SHA-256 checksums for every file being archived. Record the checksums in the archive's integrity manifest (see Section 4.4).

6. **Move to archive location.** Place the files in the correct location within the canonical directory structure. Do not copy -- move. Active storage should not retain copies of archived data unless the data is also needed for active operations.

7. **Update the catalogue.** Add entries to the metadata catalogue per D6-007. The catalogue entry must include the archive path, the checksum, the archive date, and the tier classification.

8. **Verify placement.** Confirm that the files are accessible at their archive location, that the checksums match, and that the catalogue entry is correct.

### 4.4 The Integrity Manifest

Every directory in the archive at the year level or below contains a file named `MANIFEST.sha256` that lists the SHA-256 checksum of every file in that directory. The manifest format is one line per file:

```
[checksum]  [filename]
```

This is the standard output format of `sha256sum` and can be verified with `sha256sum -c MANIFEST.sha256`.

Manifests are regenerated whenever files are added to a directory. Manifests are verified during every archive audit. A manifest that does not match the actual files indicates either corruption, unauthorized modification, or an incomplete archive transition. All three are incidents that require investigation.

### 4.5 Retrieval Procedures

To retrieve data from the archive:

1. **Identify the target.** Use the metadata catalogue (per D6-007) to locate the data. Search by date, source, descriptor, tier, or any metadata field. If the catalogue is unavailable, navigate the directory structure using the tier and date hierarchy.

2. **Verify integrity before extraction.** Run `sha256sum -c MANIFEST.sha256` in the directory containing the target file. If the checksum fails, do not use the file. Initiate the corruption response per Section 7.2.

3. **Copy, do not move.** Extract data from the archive by copying it to a working location. Never modify files in the archive directly. The archive is immutable after the transition is complete.

4. **Record the retrieval.** Log the retrieval in the operational log: what was retrieved, from where, by whom, for what purpose. This is not bureaucracy. It is the mechanism by which the institution knows which archived data is still being used, which informs retention decisions per D6-011.

5. **Return modified versions.** If the retrieved data is modified and the modification has archival value, archive the new version as a separate file with an incremented version number. The original remains untouched.

### 4.6 The Archive Audit

The archive audit is the periodic verification that the archive is intact, organized, and accessible. It occurs at three frequencies:

**Monthly Quick Audit (30 minutes):**
- Verify that the archive root is accessible and mounted correctly.
- Spot-check five randomly selected manifest files by running `sha256sum -c` on each.
- Review the archive transition log for the past month: were all transitions completed correctly?
- Check that no files exist outside the canonical directory structure (orphaned files).

**Quarterly Comprehensive Audit (2-4 hours):**
- Verify all manifests in the Tier 1 archive by running full checksum verification.
- Spot-check 20% of manifests in Tier 2 and Tier 3 archives.
- Verify that the metadata catalogue is consistent with the actual archive contents: every catalogued item exists, and every archived item is catalogued.
- Review the archive's storage consumption. Project growth rates. Flag if any storage tier is above 80% capacity.
- Verify that naming conventions are being followed. Flag and correct violations.
- Document audit results in the operational log.

**Annual Deep Audit (1-2 days):**
- Verify every manifest in the entire archive. This is the full integrity sweep.
- Verify that every archived file can be opened and read with currently available tools. This catches format rot before it becomes data loss.
- Review the directory structure for organizational problems: directories that are too full, categories that no longer make sense, structural improvements needed.
- Reconcile the archive inventory with the metadata catalogue. Resolve all discrepancies.
- Produce an annual archive health report documenting: total data volume by tier, growth trends, integrity verification results, format health, and any corrective actions taken.
- The annual audit report is itself archived as Tier 1 data.

## 5. Rules & Constraints

- **R-D6-004-01:** All archived data must reside within the canonical directory structure defined in Section 4.1. Data stored outside this structure is undiscoverable and effectively lost.
- **R-D6-004-02:** All archived files must follow the naming convention defined in Section 4.2. Files that do not conform must be renamed before archival.
- **R-D6-004-03:** Every archive directory at the year level or below must contain a valid `MANIFEST.sha256` file. Directories without manifests are considered unverified.
- **R-D6-004-04:** Archived data is immutable. Files in the archive may not be modified after the transition is complete. Corrections are made by archiving a new version, not by editing the existing file.
- **R-D6-004-05:** The monthly quick audit, quarterly comprehensive audit, and annual deep audit must be performed on schedule. Deferred audits are recorded as operational incidents.
- **R-D6-004-06:** No file may be archived without all required metadata fields per D6-007 being populated.
- **R-D6-004-07:** Retrievals from the archive must be logged. Unlogged retrievals defeat the usage-tracking that informs retention decisions.
- **R-D6-004-08:** The archive structure may be modified only through the GOV-001 Tier 3 decision process. Ad hoc structural changes create inconsistencies that compound over decades.

## 6. Failure Modes

- **Structural drift.** The operator begins placing files outside the canonical structure -- "just temporarily" -- and the temporary locations become permanent. Over years, a shadow archive grows alongside the real one. Data in the shadow archive is unaudited, uncatalogued, and effectively invisible. Mitigation: the monthly audit checks for orphaned files. R-D6-004-01 is absolute.

- **Naming convention erosion.** Under time pressure, the operator stops following naming conventions. Files accumulate with names like `backup-old-2.tar.gz` and `stuff-from-january`. Years later, these names are meaningless. Mitigation: the quarterly audit verifies naming compliance. R-D6-004-02 requires correction before archival.

- **Manifest staleness.** Manifests are not updated when files are added. The manifest says the directory contains ten files; it actually contains twelve. The audit cannot detect corruption in the unmanifested files. Mitigation: the archive transition process (Section 4.3) requires manifest regeneration. The audit verifies manifest completeness.

- **Catalogue-archive divergence.** The metadata catalogue says a file exists at a location where it does not, or files exist in the archive that the catalogue does not reference. Mitigation: the quarterly audit includes a reconciliation step. The annual audit performs full reconciliation.

- **Retrieval without logging.** The operator copies files from the archive without recording the retrieval. The institution loses visibility into which archived data is still relevant. Retention decisions are made blind. Mitigation: R-D6-004-07. The operational discipline established in OPS-001 reinforces this as habit.

- **Archive growth exceeding storage.** The archive grows faster than storage capacity. The institution faces a choice between triage and expansion. Mitigation: the quarterly audit tracks storage consumption and projects growth. D6-008 governs storage lifecycle and procurement planning.

## 7. Recovery Procedures

**7.1 If structural drift is discovered:**
1. Identify all files outside the canonical structure. Use a filesystem scan comparing actual paths against the expected structure.
2. For each orphaned file, determine its correct location in the canonical structure based on its tier, domain, date, and content.
3. Verify the file's metadata. If metadata is missing, reconstruct it to the extent possible using filesystem timestamps, file content, and the operational log.
4. Move the file to its correct location. Update the manifest and the catalogue.
5. Document the incident: how many files were orphaned, for how long, and what corrective steps were taken.
6. Review why the drift occurred. Adjust operational procedures to prevent recurrence.

**7.2 If integrity verification fails:**
1. Identify the specific files whose checksums do not match.
2. Do not delete or overwrite the corrupted files. They may be partially recoverable.
3. Check whether backup copies of the affected files exist per D6-006. If so, restore from the most recent verified backup.
4. If no backup copy exists, attempt to open and read the corrupted files. Assess what data can be salvaged.
5. Document the corruption: which files, when detected, probable cause (media degradation, incomplete write, hardware error), and the outcome of recovery attempts.
6. If the corruption indicates media degradation, escalate to D6-008 storage lifecycle procedures.

**7.3 If the catalogue is lost or corrupted:**
1. The archive itself is the primary record. The catalogue is a derived index, not the source of truth.
2. Rebuild the catalogue by walking the archive directory structure and reading the metadata files and manifest files in each directory.
3. This is labor-intensive but not data-threatening. The data is intact; only the index is lost.
4. Prioritize rebuilding the Tier 1 catalogue entries first.
5. Document the incident and investigate the cause. Ensure the catalogue itself is backed up per D6-006 and stored as Tier 1 data.

**7.4 If the archive becomes inaccessible due to media failure:**
1. Invoke D6-008 media failure procedures.
2. Restore from backup per D6-006.
3. After restoration, run a full integrity verification (annual deep audit scope) on the restored archive.
4. Document the failure and the recovery in the operational log and the annual archive health report.

## 8. Evolution Path

- **Years 0-5:** The archive structure is being established. Expect the category system to require refinement as real data reveals gaps in the initial taxonomy. Document every structural decision and the reasoning behind it. The first annual deep audit will likely reveal organizational problems that need correction.

- **Years 5-15:** The archive structure should be stable. The primary challenges will be maintaining discipline as the archive grows from gigabytes to terabytes, and managing the first wave of storage media transitions per D6-008. The quarterly audits become increasingly important as the archive grows beyond the size where the operator can hold its structure in memory.

- **Years 15-30:** The archive may contain data from multiple hardware generations, multiple storage media types, and potentially multiple operators. The directory structure's consistency across these transitions is its most important quality. The naming conventions and metadata standards prove their worth as the oldest data ages beyond anyone's direct memory.

- **Years 30-50+:** The archive is now a historical record spanning decades. Its value increases with age, but so does the maintenance burden. The annual deep audit may take several days. Consider whether automation can assist with integrity verification while maintaining the human-in-the-loop requirement for structural decisions.

- **Signpost for revision:** If the canonical directory structure consistently fails to accommodate new types of data, or if the archive has grown so large that the audit schedules are impractical, this article should be revised through the GOV-001 Tier 3 decision process.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The hardest part of archive management is not the technical structure. It is the discipline. Every single time I save a file, I must ask: does this follow the naming convention? Is the metadata complete? Is it going to the right directory? And every single time I am in a hurry or tired, the temptation is to dump the file somewhere convenient and "fix it later." Later never comes. The shadow archive grows. This is how digital entropy wins.

The canonical directory structure in Section 4.1 is deliberately simple -- three levels of hierarchy before you reach the year directory. More levels would allow finer categorization but would create more opportunities for misplacement and more cognitive overhead for every file operation. If future operators find they need more depth, they should add it within the existing structure rather than creating parallel structures. One tree, well-maintained, is worth more than a forest of abandoned filing systems.

The immutability rule (R-D6-004-04) will feel excessive at first. "Why can't I just fix the typo in the archived file?" Because the moment you allow modification of archived files, you lose the ability to trust that the archive reflects what actually happened. The archive is a record, not a living document. Living documents live in active storage. Records live in the archive. This distinction is the foundation of institutional memory.

## 10. References

- ETH-001 -- Ethical Foundations of the Institution (Principle 2: Integrity Over Convenience; Principle 4: Longevity Over Novelty)
- CON-001 -- The Founding Mandate (documentation completeness, institutional continuity)
- SEC-001 -- Threat Model and Security Philosophy (data integrity threats)
- OPS-001 -- Operations Philosophy (documentation-first principle, operational tempo, audit schedules)
- D6-001 -- Data Philosophy (four-tier data hierarchy, data lifecycle, data stewardship ethic)
- D6-002 -- The Canonical Data Model (entity types and relationships that the archive must accommodate)
- D6-003 -- Format Longevity Doctrine (approved archival formats)
- D6-006 -- Backup Verification and Test Restoration (backup of the archive; integrity verification)
- D6-007 -- Metadata Standards and Procedures (required metadata fields for all archived items)
- D6-008 -- Storage Lifecycle Management (media health monitoring for archive storage)
- D6-011 -- Retention Schedules & Disposal Policy (when archived data may be disposed)
- D6-014 -- Data Ingest Procedures (data entering the institution before archival)
- GOV-001 -- Authority Model (decision tiers for structural changes)

---

---

# D6-005 -- Format Migration Procedures

**Document ID:** D6-005
**Domain:** 6 -- Data & Archives
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, SEC-001, OPS-001, D6-001, D6-003, D6-004, D6-007, D6-008
**Depended Upon By:** D6-006, D6-009, D6-011, D6-013. All articles involving long-term data preservation or format changes.

---

## 1. Purpose

This article defines the procedures for migrating data between formats as technology changes over the institution's lifespan. Format migration is the process of converting data from one file format to another -- not because the data has changed, but because the format it was stored in has become obsolete, unsupported, or inferior to an available alternative. It is one of the most consequential operations the institution performs, because a failed migration can destroy data that survived decades of storage intact.

D6-003 establishes which formats are approved and why. This article establishes how to move data from a format that is leaving the approved list to a format that remains on it. It covers the assessment of when migration is needed, the testing that must precede any migration, the batch conversion procedures, the verification that confirms migration success, and the format longevity watchlist that provides early warning of formats approaching end-of-life.

In an air-gapped institution with a fifty-year horizon, format migration is not an exceptional event. It is a certainty. File formats that are ubiquitous today will be obscure in fifteen years and potentially unreadable in thirty. The institution that does not plan for format migration is the institution that will one day discover it has terrabytes of data it can no longer open. This article ensures that migration happens proactively, safely, and completely.

## 2. Scope

This article covers:

- The format longevity watchlist: how formats are monitored for declining support and viability.
- Assessment criteria for determining when migration is necessary.
- Pre-migration testing procedures.
- Batch conversion procedures for large-scale format migration.
- Post-migration verification: confirming that converted data is complete and faithful.
- Rollback procedures: what to do when a migration fails.
- The migration log: documenting every migration for institutional memory.

This article does not cover:

- Which formats are approved or why (see D6-003).
- The philosophy of format selection (see D6-001 and D6-003).
- Storage media migration, which is moving data between physical media types (see D6-008).
- Data migration in the context of disaster recovery (see D6-013).
- The daily operational handling of format conversion during data ingest (see D6-014).

## 3. Background

### 3.1 The Inevitability of Format Obsolescence

No file format lasts forever. Formats depend on software that reads and writes them, on specifications that document their structure, and on communities that maintain the tools and knowledge required to work with them. All of these can decline. When the last tool capable of reading a format ceases to function, every file in that format becomes inaccessible -- not corrupted, not deleted, but locked behind an encoding that no available software can interpret.

The history of computing is littered with format casualties. WordStar documents. Lotus 1-2-3 spreadsheets. MacPaint files. RealMedia streams. Each was dominant in its era. Each is now difficult or impossible to open without specialized recovery efforts. The institution's choice of open, well-documented formats per D6-003 reduces but does not eliminate this risk. Even open formats can decline if the community maintaining their tools disperses.

### 3.2 The Fidelity Problem

Format migration is not like copying a file. When you copy a file, you produce an identical sequence of bytes. When you convert a file between formats, you produce a new sequence of bytes that represents the same information in a different encoding. The representation may lose nuance. A rich-text document converted to plain text loses its formatting. A lossy-compressed image converted to another lossy format loses additional quality. Even lossless conversions can introduce subtle differences -- metadata handled differently, encoding variations, precision changes.

The fidelity problem means that every format migration must be verified, not merely executed. The operator must confirm that the converted data faithfully represents the original. This verification is format-specific and cannot be fully automated. It requires human judgment.

### 3.3 The Scale Problem

When a format migration is triggered, it typically affects not one file but hundreds or thousands. The institution may have ten years of data in a particular format. Converting all of it is a project, not a task. The batch conversion procedures in this article are designed to handle this scale while maintaining verification rigor.

## 4. System Model

### 4.1 The Format Longevity Watchlist

The institution maintains a format longevity watchlist -- a living document that tracks the health of every file format in active use within the archive. The watchlist is reviewed quarterly and updated during each review.

For each format on the watchlist, the following information is maintained:

- **Format name and version.** The specific format being tracked (e.g., "ODF 1.3", "PNG 1.2", "FLAC 1.3.4").
- **Primary tool(s).** The software used to read and write this format within the institution.
- **Tool version and last update.** The version of the tools and when they were last updated.
- **Specification status.** Whether the format specification is open, publicly available, and actively maintained.
- **Community health.** A qualitative assessment of the format's user and developer community: thriving, stable, declining, or moribund.
- **Institutional data volume.** How much data the institution stores in this format (approximate).
- **Migration target.** If migration becomes necessary, what format would the data be converted to.
- **Watchlist status.** One of: Green (healthy, no action needed), Yellow (declining, monitor closely), Orange (migration planning required), Red (migration in progress or overdue).

**Status transition criteria:**

- **Green to Yellow:** The primary tool has not received a meaningful update in two years, OR the community shows signs of declining activity, OR the specification body has ceased meeting.
- **Yellow to Orange:** The primary tool no longer compiles or runs on the institution's current operating system without significant workarounds, OR a viable successor format has emerged and the community is migrating to it, OR the specification has been formally deprecated.
- **Orange to Red:** The primary tool will not function after the next planned OS upgrade, OR the format has known deficiencies that create data integrity risks, OR the migration target format has been validated and the migration plan is ready.

### 4.2 The Migration Assessment

When a format reaches Orange status on the watchlist, the operator initiates a migration assessment. The assessment determines whether migration is necessary and, if so, what the migration plan should be.

The assessment answers five questions:

1. **Is migration truly necessary, or can the existing tools be preserved?** In some cases, maintaining an older tool (or its runtime environment) is less risky than converting large volumes of data. If the tool can be preserved indefinitely through virtualization or static compilation, migration may be deferred. However, deferral is itself a decision that must be documented and reviewed.

2. **What is the migration target?** The target format must be on the D6-003 approved list. If no approved format is suitable, the format evaluation process in D6-003 must be invoked to approve a new target.

3. **What is the fidelity expectation?** For each data type in the source format, what information must be preserved in the target format? What information may be lost? What loss is unacceptable? These criteria are format-specific and must be explicitly stated.

4. **What is the scale?** How many files, what total volume, and how long will conversion take at the available processing speed? This determines whether the migration is a single-session task or a multi-week project.

5. **What are the risks?** What can go wrong during conversion? What is the rollback plan? What is the cost of failure?

The migration assessment is documented as a Tier 3 decision per GOV-001 and archived as Tier 1 data.

### 4.3 Pre-Migration Testing

No batch migration proceeds without testing. The testing procedure is:

1. **Select a test sample.** Choose a representative sample of files from the source dataset. The sample must include: the smallest file, the largest file, files with the most complex content (maximum metadata, maximum structural complexity), files from the oldest and newest dates in the dataset, and at least 20 files selected randomly from the full dataset.

2. **Convert the test sample.** Run the conversion tool on the test sample, capturing all output, warnings, and errors.

3. **Verify the converted files.** For each converted file:
   - Confirm it opens without error in the target format's primary tool.
   - Compare it visually or structurally to the original. For text formats, use diff. For image formats, compare visually and check dimensions, color depth, and file size. For structured data, compare field-by-field.
   - Verify that all required metadata survived the conversion or was correctly re-applied.
   - Check for silent data loss: information present in the source that is absent in the target without any warning from the conversion tool.

4. **Document test results.** Record: which files were tested, what conversion tool and version were used, what issues were found, and whether the test is considered passing or failing.

5. **Iterate if needed.** If the test reveals problems, adjust the conversion process (different tool, different settings, pre-processing steps) and re-test. Do not proceed to batch conversion until the test sample passes.

### 4.4 Batch Conversion Procedures

Once testing passes, the batch conversion proceeds:

1. **Create a migration workspace.** A separate directory outside the archive where converted files will be staged before replacing the originals. The workspace is on storage with sufficient capacity for the entire converted dataset.

2. **Snapshot the source.** Create a complete backup of the source data before any conversion begins. This backup is the rollback point. It is not disposed of until the migration is fully verified and the retention period in Step 8 has elapsed.

3. **Run the batch conversion.** Process files in batches of manageable size (100-500 files per batch, depending on file size and conversion speed). After each batch:
   - Record the number of files processed, successes, warnings, and failures.
   - Spot-check three files from the batch using the verification criteria from the test phase.
   - If any batch produces a failure rate above 2%, stop the conversion and investigate.

4. **Generate manifests.** Create SHA-256 manifests for the converted files in the migration workspace, following D6-004 manifest procedures.

5. **Perform comprehensive verification.** After all batches are complete:
   - Verify the total file count matches the source.
   - Run automated comparison where possible (file count, total size within expected ratio, metadata field presence).
   - Manually verify a random sample of at least 5% of converted files or 50 files, whichever is greater.

6. **Move converted files to archive.** Replace the source files in the archive with the converted files. Update manifests. Update the metadata catalogue.

7. **Retain source backup.** Keep the pre-migration backup for a minimum of 180 days after migration completion. This is the grace period during which undetected migration problems may surface.

8. **Update the watchlist.** Change the format's watchlist status. If no more data exists in the source format, remove it from the watchlist. Add the target format if it is not already tracked.

9. **Document the migration.** Record in the migration log: source format, target format, number of files migrated, tool used, dates of migration, verification results, any issues encountered, and the disposition of the source backup.

### 4.5 The Migration Log

The migration log is a Tier 1 document that records every format migration performed in the institution's history. It is the institutional memory of what data has been converted, when, and how. It enables future operators to understand the provenance of data that may have been through multiple format transitions.

Each entry contains: migration date, source format, target format, number of files, conversion tool and version, verification method, issues encountered, and the location of the migration assessment document.

## 5. Rules & Constraints

- **R-D6-005-01:** No format migration shall proceed without a documented migration assessment approved as a Tier 3 decision per GOV-001.
- **R-D6-005-02:** No batch migration shall proceed without successful pre-migration testing on a representative sample.
- **R-D6-005-03:** A complete backup of the source data must exist before any batch conversion begins. This backup must be verified and retained for a minimum of 180 days after migration completion.
- **R-D6-005-04:** The format longevity watchlist must be reviewed at least quarterly. Deferred reviews are recorded as operational incidents.
- **R-D6-005-05:** Every migration must be recorded in the migration log with all required fields. Undocumented migrations are institutional memory gaps.
- **R-D6-005-06:** If a batch conversion produces a failure rate above 2%, conversion must be halted and the process re-evaluated before continuing.
- **R-D6-005-07:** Source data in the archive must not be deleted or overwritten until the converted data has been verified and the 180-day grace period has elapsed.
- **R-D6-005-08:** Silent data loss -- information lost during conversion without any warning from the conversion tool -- is the most dangerous migration failure. Verification procedures must specifically check for it.

## 6. Failure Modes

- **Deferred migration.** A format reaches Orange on the watchlist, but the migration is deferred because it seems like a lot of work. The format reaches Red. The tools stop working. The data is stranded. Mitigation: the quarterly watchlist review creates accountability. The assessment process at Orange status forces engagement before the situation becomes critical.

- **Untested batch conversion.** The operator skips testing because the conversion tool "should work." It produces systematic errors -- metadata stripped, encoding corrupted, truncated output -- across thousands of files before anyone notices. Mitigation: R-D6-005-02 requires testing. The batch process includes per-batch spot-checks.

- **Silent fidelity loss.** The conversion appears successful, but subtle information is lost -- floating-point precision reduced, Unicode characters replaced, color profiles discarded. The loss is not detected until years later when the data is needed and the original no longer exists. Mitigation: R-D6-005-08 identifies this as the primary threat. Verification must include format-specific fidelity checks, not just "does it open."

- **Backup disposal too early.** The source backup is deleted before the 180-day grace period, and a migration problem is discovered afterward. Mitigation: R-D6-005-07 establishes the minimum retention period. The operational log tracks backup disposition dates.

- **Conversion tool unavailability.** The tool needed for migration is itself becoming obsolete. The institution faces the problem of migrating data with a tool that is also disappearing. Mitigation: the watchlist tracks tool health, not just format health. When a tool is declining, securing a working copy of the tool (including its runtime dependencies) is part of migration planning.

- **Chain migration fragility.** Data that has been through multiple format migrations accumulates fidelity losses at each step. By the third or fourth migration, the data may have drifted significantly from the original. Mitigation: the migration log records the complete conversion history. When possible, migrations should convert from the earliest available version rather than from the most recent intermediate format.

## 7. Recovery Procedures

**7.1 If a batch conversion fails mid-process:**
1. Stop the conversion immediately. Do not continue processing additional batches.
2. Assess the scope of the failure: how many files were affected, what went wrong.
3. The source data is intact (it was backed up before conversion began per R-D6-005-03). The failure affects only the migration workspace.
4. Investigate the root cause. Common causes: insufficient disk space, tool bug triggered by specific file characteristics, corrupted input file, interrupted process.
5. Fix the root cause.
6. Clear the migration workspace of partially converted data.
7. Restart the conversion from the beginning, not from the point of failure. Partial batches cannot be trusted.
8. Document the failure and the corrective action in the migration log.

**7.2 If post-migration verification reveals systematic errors:**
1. Determine the scope: are all converted files affected, or only a subset?
2. Determine the nature of the error: data loss, corruption, metadata stripping, encoding error.
3. If the source backup exists, the recovery is straightforward: discard the converted files and restore from the source backup.
4. Re-evaluate the conversion process. A different tool, different settings, or a different target format may be needed.
5. Return to pre-migration testing with the revised process.
6. Document the failed migration attempt in the migration log.

**7.3 If format migration becomes necessary but no conversion tool exists:**
1. This is a crisis. It means the institution waited too long to migrate and the tools have disappeared.
2. Search the institution's archived software collection for older versions of tools that can read the format.
3. If a tool is found, assess whether it can be run (it may require an older OS; consider virtualization).
4. If no tool is found, obtain the format specification (which should have been archived per D6-003). Assess the feasibility of writing a custom parser.
5. If the format specification is not available, attempt to reverse-engineer the format from the data itself. This is a last resort and may result in partial data recovery.
6. Document the incident as a critical failure. Update D6-003 and the watchlist process to prevent recurrence.

**7.4 If a migration is discovered to have caused data loss after the source backup was disposed:**
1. This is a permanent data loss event. Acknowledge it honestly.
2. Determine the scope of the loss: which specific information was lost across how many files.
3. Assess whether the lost information can be reconstructed from other sources (related documents, logs, other copies).
4. Document the loss in the operational log and the annual archive health report as a Tier 1 incident.
5. Review why the verification process failed to detect the loss. Update verification procedures.
6. Extend the source backup retention period for future migrations (the 180-day minimum may need to increase).

## 8. Evolution Path

- **Years 0-5:** The format landscape is stable. The watchlist is established and populated. Most formats are Green. The procedures are tested with minor, elective migrations (moving to a newer version of the same format) to build operational familiarity.

- **Years 5-15:** The first compulsory migrations will likely occur. Formats chosen at founding will begin showing their age. This is the period where the migration procedures prove their worth or reveal their gaps. The migration log begins to accumulate history.

- **Years 15-30:** Multiple rounds of migration may have occurred. The chain migration problem becomes relevant. The watchlist process should be mature enough to provide multi-year warning before a format becomes critical.

- **Years 30-50+:** Formats that do not exist today will be the migration targets. The principles in this article -- test before you convert, verify after you convert, keep the backup until you are sure -- are format-agnostic and should remain valid. The specific tools and procedures will need updating.

- **Signpost for revision:** If the institution has experienced a format migration failure that these procedures did not adequately prevent, or if the scale of the archive makes the current verification procedures impractical, this article should be revised through the GOV-001 Tier 3 decision process.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
Format migration is the data operation I dread most. Every other data operation -- backup, verification, cataloguing -- works with the data as it is. Migration transforms the data into something different and asks you to trust that the transformation was faithful. That trust must be earned through testing and verification, not assumed.

The 180-day source backup retention period is a compromise. Ideally, you would keep the original format forever, just in case. But "just in case" conflicts with the storage discipline that D6-001 demands. One hundred eighty days is long enough for most migration problems to surface through normal use of the converted data. If future operators find that this period is insufficient -- if migration problems are discovered after 180 days -- they should extend it. Better to spend storage on insurance than to discover you are uninsured.

The format longevity watchlist is the most important artifact this article creates. It is the early warning system. A format that moves from Green to Yellow is not a crisis -- it is a signal that you have years to plan. A format that moves from Yellow to Orange is a call to action with months of lead time. If a format ever goes directly from Green to Red, something has gone badly wrong with the monitoring process, and the watchlist review procedures need immediate revision.

## 10. References

- ETH-001 -- Ethical Foundations of the Institution (Principle 2: Integrity Over Convenience; Principle 4: Longevity Over Novelty)
- CON-001 -- The Founding Mandate (air-gap mandate, institutional continuity)
- SEC-001 -- Threat Model and Security Philosophy (data integrity threats)
- OPS-001 -- Operations Philosophy (documentation-first principle, operational tempo)
- D6-001 -- Data Philosophy (format sovereignty, data stewardship, triage framework)
- D6-003 -- Format Longevity Doctrine (approved formats, format evaluation process)
- D6-004 -- Archive Management Procedures (directory structure, manifests, catalogue)
- D6-007 -- Metadata Standards and Procedures (metadata preservation during migration)
- D6-008 -- Storage Lifecycle Management (storage media context for migration planning)
- D6-009 -- Data Migration Doctrine (broader migration context)
- D6-013 -- Disaster Recovery for Data Systems (migration in disaster scenarios)
- GOV-001 -- Authority Model (Tier 3 decision process for migration assessments)

---

---

# D6-006 -- Backup Verification and Test Restoration

**Document ID:** D6-006
**Domain:** 6 -- Data & Archives
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, SEC-001, OPS-001, D6-001, D6-003, D6-004, D6-007, D6-008
**Depended Upon By:** D6-005, D6-009, D6-011, D6-013, D10-002. All articles involving data recovery, disaster planning, or backup operations.

---

## 1. Purpose

This article defines the procedures for verifying that backups actually work and for conducting scheduled test restorations that prove recoverability. A backup that has never been tested is not a backup. It is a hope. This institution cannot afford to operate on hope.

The principle is simple: if you have not restored from a backup and verified the result, you do not know whether the backup works. You may believe it works. The backup tool may report success. The file sizes may look right. But until you have taken the backup, restored it to a separate location, and confirmed that the restored data is complete, readable, and uncorrupted, you have an untested assumption where you need a verified fact.

D6-001 establishes the philosophical mandate for redundancy: Tier 1 data must exist in at least three independent copies on at least two distinct media types. But copies are worthless if they are corrupted, incomplete, or in a format that cannot be restored. This article provides the verification and testing procedures that turn the D6-001 mandate from a principle into a guarantee. It defines what to verify, how to verify it, when to verify it, and what to do when verification reveals a problem.

## 2. Scope

This article covers:

- Backup integrity verification: confirming that backup files are complete and uncorrupted.
- Test restoration procedures: restoring from backup to a test location and verifying the result.
- Post-restoration verification: what to check after a test restoration to confirm success.
- The quarterly verification schedule: what is tested each quarter.
- The annual verification schedule: the comprehensive annual test.
- Verification documentation and record-keeping.
- Response procedures when verification fails.

This article does not cover:

- Backup creation procedures (see the backup creation doctrine; referenced as D6-006 in D6-001 but the creation procedures are a companion document -- this article focuses on verification).
- The philosophy of backup and redundancy (see D6-001).
- Archive structure and organization (see D6-004).
- Storage media health monitoring (see D6-008).
- Disaster recovery from a real data loss event (see D6-013).

## 3. Background

### 3.1 The Untested Backup Fallacy

The most common data loss scenario in well-administered systems is not hardware failure. It is the discovery, at the moment of crisis, that the backups are unusable. The causes are varied and instructive: the backup process silently failed months ago and nobody noticed. The backup is encrypted and the decryption key was stored on the same system that failed. The backup medium has degraded and produces read errors. The backup was made with a tool version that is no longer installed, and the format is not backward-compatible. The backup is complete but the restoration process has never been tested, and the operator discovers -- under the stress of a real emergency -- that they do not actually know how to restore.

Every one of these failures was preventable through regular test restoration. Every one of them occurred because the operator trusted the backup process without verifying the backup product. This article exists to make test restoration a routine, scheduled, documented operation rather than something that happens for the first time during a crisis.

### 3.2 The Air-Gap Backup Challenge

In an air-gapped institution, backup verification has additional constraints. There is no cloud backup to verify with a single API call. There is no offsite backup service that provides automated verification reports. Every backup exists on physical media that the operator must physically access, mount, read, and verify. Offsite backups require physical travel to the offsite location. These logistics make verification more expensive in time and effort, which makes it more tempting to skip. The scheduled verification requirements in this article exist precisely because the temptation to skip is highest when the effort is greatest.

### 3.3 Verification Is Not Validation

A distinction must be maintained between verification and validation. Verification confirms that the backup is intact -- the files exist, the checksums match, the data can be read. Validation confirms that the backup is sufficient -- if the primary data were lost, the backup would be adequate to restore operations. Verification is a technical check. Validation is a judgment call. This article addresses both, but they are performed at different frequencies and with different methods.

## 4. System Model

### 4.1 Backup Integrity Verification

Backup integrity verification confirms that backup files have not been corrupted or truncated since they were created. It is performed at every backup verification event.

**Procedure for integrity verification:**

1. **Mount the backup media.** For local backups on attached media, mount the filesystem. For removable media (USB drives, external drives), physically connect and mount. For offsite media, physically retrieve, transport, connect, and mount.

2. **Verify the backup manifest.** Every backup set must include a manifest file containing the SHA-256 checksums of all backed-up files (mirroring the archive manifest per D6-004). Run `sha256sum -c` against the manifest. Every file must match.

3. **Verify backup completeness.** Compare the file count and directory structure of the backup against the current archive. The backup should contain at least every file that existed at the time the backup was created. Use a file listing comparison rather than relying solely on checksum verification.

4. **Verify backup readability.** Select at least 10 files from the backup and open each one using its native tool. Confirm that the file is not merely present but actually readable and contains the expected content. This catches silent corruption that produces valid-looking files with garbled content.

5. **Record the verification result.** Document: date, backup set identified (creation date, media label, storage location), number of files verified, checksum results (pass/fail), completeness check result, readability spot-check result, and overall pass/fail.

### 4.2 Test Restoration Procedures

Test restoration is the process of restoring data from a backup to a separate location and confirming that the restored data is usable. It is more thorough than integrity verification because it tests the entire restoration pipeline, not just the backup files.

**Full test restoration procedure:**

1. **Prepare the test environment.** Identify a restoration target -- a separate disk, partition, or directory that is not the production archive. The test environment must have sufficient space to hold the restored data. It should be a location that can be wiped after the test without affecting any production data.

2. **Document the restoration scenario.** Before beginning, write down what you are restoring, from which backup set, to what location, and what success looks like. This is not bureaucracy -- it is the discipline that ensures you are testing a realistic scenario rather than an easy subset.

3. **Perform the restoration.** Using the same tools and procedures that would be used in a real recovery, restore the backup data to the test environment. Follow the restoration steps exactly as they are documented in the backup creation procedures. If you discover that the documented steps are incomplete or inaccurate, this is a test finding that must be corrected.

4. **Verify the restored data.** After restoration is complete, perform the following checks:

   a. **Structural verification.** Compare the directory structure of the restored data against the expected structure per D6-004. All directories must be present. The hierarchy must be correct.

   b. **Completeness verification.** Compare the file count of the restored data against the source (the current archive or the backup manifest). All files must be present.

   c. **Integrity verification.** Run `sha256sum -c` against the manifest files within the restored data. Every checksum must match.

   d. **Readability verification.** Open and examine at least 20 files across different data types, dates, and tiers. Confirm that each file is readable, complete, and contains the expected content.

   e. **Metadata verification.** Confirm that file metadata (timestamps, permissions, ownership) was correctly preserved during backup and restoration. Metadata loss is a common silent failure in backup systems.

   f. **Functional verification.** If the backup includes configuration files, databases, or other functional data, verify that the restored data is not merely present but operational -- a restored database can be queried, a restored configuration can be loaded.

5. **Document the test results.** Produce a test restoration report containing: date, backup set tested, restoration scenario, all verification results (pass/fail with details), any issues discovered, and the time required for the full restoration (this data is critical for disaster recovery planning per D6-013).

6. **Clean up the test environment.** After documentation is complete, securely wipe the test environment. Restored data in test locations is a liability -- it is an uncontrolled copy that is not subject to the access controls and audit trail of the production archive.

### 4.3 Quarterly Verification Schedule

Every quarter, the following verification activities are performed:

**Quarter 1 (January-March):**
- Full integrity verification of all Tier 1 local backups.
- Spot integrity verification (20% sample) of Tier 2 local backups.
- Test restoration of one complete Tier 1 backup set to the test environment.
- Verification that all backup media are healthy and readable per D6-008.

**Quarter 2 (April-June):**
- Full integrity verification of all Tier 1 local backups.
- Spot integrity verification (20% sample) of Tier 2 local backups.
- Test restoration of one complete Tier 2 backup set to the test environment.
- Verification of offsite backup inventory: confirm that offsite media list matches the expected inventory (physical retrieval not required unless the annual test is in this quarter).

**Quarter 3 (July-September):**
- Full integrity verification of all Tier 1 local backups.
- Spot integrity verification (20% sample) of Tier 2 local backups.
- Test restoration of one complete Tier 1 backup set, using a different backup set than Q1.
- Test of backup encryption: verify that encrypted backups can be decrypted with the current keys and that the keys are accessible per SEC-003.

**Quarter 4 (October-December):**
- Full integrity verification of all Tier 1 local backups.
- Full integrity verification of all Tier 2 local backups (annual comprehensive sweep).
- Spot integrity verification (20% sample) of Tier 3 backups.
- Test restoration of one complete Tier 2 backup set, using a different backup set than Q2.

### 4.4 Annual Verification Schedule

Once per year, in addition to the quarterly verification for that quarter, the following comprehensive test is performed:

1. **Full offsite backup verification.** Physically retrieve the offsite backup media. Perform full integrity verification on all offsite backups. Perform a test restoration of at least one complete Tier 1 offsite backup set. Return the media to the offsite location after testing.

2. **Full disaster recovery simulation.** Simulate a complete primary storage failure. Using only the backup media and the documented restoration procedures, restore the entire Tier 1 archive to the test environment. This tests not just the backups but the completeness and accuracy of the restoration documentation.

3. **Restoration timing assessment.** Measure how long the full restoration takes. Compare against the recovery time objective established in D6-013. If the actual time exceeds the objective, the gap must be addressed -- either by improving the restoration process or by revising the objective to match reality.

4. **Backup strategy review.** Using the year's quarterly verification results, assess the overall health of the backup strategy. Are backups consistently passing verification? Have any recurring issues emerged? Is the backup schedule adequate for the institution's current data growth rate? Are backup media aging out per D6-008?

5. **Annual verification report.** Produce a comprehensive report documenting all quarterly and annual verification results for the year. This report is archived as Tier 1 data and informs the backup strategy for the following year.

### 4.5 Post-Restoration Verification Checklist

This checklist is used after every test restoration, whether quarterly or annual. Every item must be explicitly marked pass or fail.

```
[ ] Directory structure matches expected structure per D6-004
[ ] Total file count matches source/manifest
[ ] SHA-256 checksums pass for all manifest-verified files
[ ] 20+ files opened and content verified as readable and correct
[ ] File timestamps preserved correctly
[ ] File permissions preserved correctly (where applicable)
[ ] Metadata files (.meta) present and intact
[ ] Manifest files (MANIFEST.sha256) present and self-consistent
[ ] No unexpected files present in restored data
[ ] No error messages or warnings during restoration process
[ ] Encrypted backups decrypted successfully
[ ] Total restoration time recorded: _____ hours/minutes
[ ] Test environment securely wiped after verification
```

## 5. Rules & Constraints

- **R-D6-006-01:** No backup shall be trusted for disaster recovery purposes until it has passed at least one test restoration. New backup sets must be test-restored within 30 days of creation.
- **R-D6-006-02:** Quarterly backup verification must be completed within the quarter. Deferred verifications are recorded as operational incidents and must be completed before the next quarterly cycle begins.
- **R-D6-006-03:** The annual offsite backup verification must include physical retrieval and testing. Remote verification is not possible in an air-gapped institution and must not be substituted.
- **R-D6-006-04:** Every verification event must be documented with all required fields. Undocumented verifications have no institutional value -- if it is not recorded, it did not happen.
- **R-D6-006-05:** When backup verification reveals a failure, the failed backup must not be relied upon for recovery until the failure is corrected and the backup re-verified.
- **R-D6-006-06:** The test environment used for restoration must be securely wiped after each test. Lingering restored data is an uncontrolled copy that violates the access control and audit trail principles of SEC-002 and D6-004.
- **R-D6-006-07:** Backup encryption keys must be verified as accessible during at least one quarterly test per year. A backup that exists but cannot be decrypted is equivalent to no backup.
- **R-D6-006-08:** The annual full disaster recovery simulation may not be waived or reduced in scope. It is the single most important verification event in the institution's operational calendar.

## 6. Failure Modes

- **Verification complacency.** After several quarters of passing verifications, the operator begins treating the verification as a formality. Checks become cursory. The spot-check of file readability is skipped. A silent corruption goes undetected. Mitigation: the post-restoration verification checklist (Section 4.5) requires explicit pass/fail for every item. The annual report must evidence full completion.

- **Test-production divergence.** The test environment does not accurately represent the production restoration scenario. The test succeeds, but a real restoration would fail because the test did not exercise the same code paths, configurations, or media. Mitigation: the annual disaster recovery simulation (Section 4.4) uses the full restoration procedure and should be as realistic as feasible.

- **Offsite backup neglect.** The offsite backups are the hardest to verify because they require physical travel. Quarters pass without offsite verification. When the offsite backups are finally tested, they are degraded or incomplete. Mitigation: R-D6-006-03 requires annual physical verification. The quarterly Q2 inventory check provides intermediate assurance.

- **Key inaccessibility.** The backup is intact but encrypted, and the decryption key is not available -- it was on the failed system, or the passphrase was forgotten, or the key has been rotated and the old key discarded. Mitigation: the Q3 quarterly test specifically verifies encryption key accessibility. SEC-003 defines key management including backup of encryption keys.

- **Restoration procedure documentation rot.** The restoration procedures were written when the backup system was established and never updated. The actual tools, paths, or configurations have changed. During a real restoration, the documented procedure fails. Mitigation: every test restoration uses the documented procedure. Discrepancies between the documented procedure and the actual required steps are a test finding that must be corrected immediately.

- **Verification infrastructure failure.** The test environment itself fails during verification -- insufficient space, hardware failure, tool unavailability. The verification is incomplete and deferred. Mitigation: preparation of the test environment is a prerequisite step in the verification procedure. If the test environment is unavailable, that is an operational issue to be resolved before the verification deadline.

## 7. Recovery Procedures

**7.1 If backup integrity verification fails:**
1. Identify which specific files or backup sets failed verification.
2. Determine the nature of the failure: checksum mismatch (corruption), missing files (incomplete backup), unreadable media (media degradation).
3. If the failed backup is one of multiple copies, verify the other copies. If another copy passes, the institution still has a valid backup.
4. If all copies of a backup set have failed, determine whether the source data is still available in the production archive. If so, create a new backup immediately.
5. Investigate the root cause. Common causes: media degradation (see D6-008), backup tool error, interrupted backup process, storage controller error.
6. Document the failure, root cause, and corrective action in the operational log and the quarterly verification report.

**7.2 If test restoration fails:**
1. Determine the point of failure: did the restoration tool fail, did the data fail to decompress/decrypt, did the restored data fail structural/integrity/readability checks?
2. If the restoration tool failed: verify the tool is correctly installed and configured. Check for version mismatches between the backup creation tool and the restoration tool.
3. If the data is corrupt: fall back to Section 7.1.
4. If the restoration succeeded but verification failed (e.g., permissions wrong, metadata missing): investigate whether the failure is in the backup or in the restoration process. Fix the process and re-test.
5. Do not mark the backup as verified until the test restoration passes completely.

**7.3 If the annual disaster recovery simulation reveals the restoration time exceeds the recovery objective:**
1. Identify the bottleneck: is it media read speed, decompression time, network transfer (if applicable between zones), or verification time?
2. Consider whether the bottleneck can be addressed by faster media, parallel restoration, or reduced verification during emergency restoration (with full verification after).
3. If the bottleneck cannot be practically resolved, revise the recovery time objective in D6-013 to reflect reality. An unrealistic objective is worse than an honest one because it creates false confidence.
4. Document the gap and the plan to address it in the annual verification report.

**7.4 If backup encryption keys are not accessible during verification:**
1. This is a critical incident. Immediately invoke SEC-003 key recovery procedures.
2. Verify that the key exists in the Tier B (local) and Tier C (offsite) key backup locations per SEC-003.
3. If the key is recoverable, restore it and complete the verification.
4. If the key is not recoverable, the encrypted backup is lost. Determine what data is affected and whether unencrypted copies exist.
5. Document the incident as a Tier 1 security and data integrity event.
6. Review the key management procedures to prevent recurrence.

## 8. Evolution Path

- **Years 0-5:** The verification procedures are new and the operator is building the habit. The first annual disaster recovery simulation will likely reveal gaps in the restoration documentation. This is expected and valuable -- better to discover the gaps during a test than during a crisis.

- **Years 5-15:** The verification schedules should be routine. The primary challenge is maintaining rigor as the archive grows. Larger archives mean longer verification times. The institution may need to invest in faster verification infrastructure or accept longer verification windows.

- **Years 15-30:** The backup strategy will have evolved through multiple storage media generations per D6-008. Each media transition introduces a window of verification vulnerability where the new media has not yet accumulated a track record. Increase verification frequency during transitions.

- **Years 30-50+:** The annual disaster recovery simulation is the most important test. By this point, the institution may be maintained by someone who was not present when the backups were created. The simulation tests not just the backups but the comprehensibility of the restoration documentation to a person who did not write it.

- **Signpost for revision:** If the verification schedules consistently cannot be completed within their designated time windows, or if the verification procedures do not catch a backup failure that causes real data loss, this article requires immediate revision.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
I have never experienced a catastrophic data loss. I have experienced the cold sweat of realizing that my backups had not been tested and that I had no real confidence they would work. That experience is the reason this article exists.

The annual disaster recovery simulation will be the most time-consuming verification event in the operational calendar. It may take an entire day. It will be tempting to skip it when things are busy or to reduce its scope to "just check a few files." Do not. The simulation is the only test that exercises the full restoration pipeline under realistic conditions. Every shortcut in the simulation is a gap in the institution's knowledge of its own recoverability.

The quarterly rotation of test scenarios -- testing Tier 1 in Q1 and Q3, Tier 2 in Q2 and Q4, offsite in the annual test -- is designed to ensure that every backup type is tested at least once per year without making any single quarter's verification burden overwhelming. If the institution's data grows to the point where even a single-tier test restoration takes multiple days, the verification schedule should be revised to spread the work more evenly, not to reduce the total testing.

One principle I want to state explicitly: a failed verification is a success, not a failure. A test that reveals a problem has done its job. A test that always passes may not be testing hard enough. When a verification fails, the correct response is gratitude that the problem was found during a test rather than during a crisis, followed by immediate corrective action.

## 10. References

- ETH-001 -- Ethical Foundations of the Institution (Principle 2: Integrity Over Convenience)
- CON-001 -- The Founding Mandate (air-gap mandate, self-sovereign operation)
- SEC-001 -- Threat Model and Security Philosophy (data integrity threats, defense in depth)
- SEC-003 -- Cryptographic Key Management (backup encryption keys, key recovery)
- OPS-001 -- Operations Philosophy (operational tempo, documentation-first principle)
- D6-001 -- Data Philosophy (R-D6-02: minimum three copies on two media types; data tier hierarchy)
- D6-003 -- Format Longevity Doctrine (format compatibility for restored data)
- D6-004 -- Archive Management Procedures (directory structure, manifests, integrity verification)
- D6-007 -- Metadata Standards and Procedures (metadata verification during restoration)
- D6-008 -- Storage Lifecycle Management (media health context for backup media)
- D6-013 -- Disaster Recovery for Data Systems (recovery time objectives, full recovery procedures)
- D10-002 -- Daily Operations Doctrine (backup as daily operational task)
- GOV-001 -- Authority Model (incident classification for verification failures)

---

---

# D6-007 -- Metadata Standards and Procedures

**Document ID:** D6-007
**Domain:** 6 -- Data & Archives
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, SEC-001, OPS-001, D6-001, D6-002, D6-003, D6-004
**Depended Upon By:** D6-005, D6-006, D6-008, D6-009, D6-011, D6-014. All articles involving data cataloguing, search, or retrieval.

---

## 1. Purpose

This article defines how metadata is created, maintained, and used within the holm.chat Documentation Institution. Metadata is data about data -- the descriptive, structural, and administrative information that makes institutional data findable, understandable, and manageable over time. Without metadata, a file is an opaque blob of bytes. Its name hints at its contents. Its timestamp suggests when it was created. But its purpose, its relationships, its classification, its provenance, and its retention requirements are unknown. Metadata is what transforms a file from an artifact into a record.

In an air-gapped institution designed to persist for fifty years, metadata is not a convenience. It is infrastructure. The metadata is what will allow a future operator who has never seen a particular file to understand what it is, where it came from, why it was kept, and how it relates to everything else in the archive. Without metadata, institutional data degrades from knowledge to noise over a single generation.

D6-001 establishes the principle that data without context is noise. This article provides the procedures that create and maintain that context. It defines the required metadata fields for every data type, the quality checks that ensure metadata remains accurate, the search and discovery mechanisms that use metadata to make the archive useful, and the metadata audit that verifies compliance.

## 2. Scope

This article covers:

- The institutional metadata schema: what fields are required, recommended, and optional for each data type.
- Metadata creation procedures: how and when metadata is generated.
- Metadata maintenance: how metadata is updated as data moves through its lifecycle.
- Metadata quality checks: how to verify that metadata is complete, accurate, and consistent.
- Search and discovery: how metadata enables finding data in the archive.
- The metadata catalogue: the central index of all institutional metadata.
- The metadata audit: scheduled verification of metadata health.

This article does not cover:

- The canonical data model (see D6-002).
- Archive directory structure and naming (see D6-004, though metadata and naming are related).
- File format specifications (see D6-003).
- Retention decisions informed by metadata (see D6-011).
- The data ingest process that triggers initial metadata creation (see D6-014).

## 3. Background

### 3.1 The Metadata Crisis

The metadata crisis is not that metadata is missing. It is that metadata is absent precisely when it is needed most -- years after the data was created, when the original context has evaporated. The person who created the file knew what it was. They did not need metadata to understand it. But the person who discovers the file ten years later, or the successor operator who inherits the archive, has none of that implicit knowledge. They need metadata, and if the metadata was never created, the data is effectively orphaned.

This crisis is universal in digital systems. Email archives with millions of messages and no tagging. Photo libraries with thousands of images and no descriptions. Document collections where the folder name "Projects" contains three hundred undifferentiated files. The institutional answer is not to hope that future operators will intuit the meaning of old data. It is to capture the meaning now, while the context exists, in a structured and durable form.

### 3.2 Metadata in an Air-Gapped Environment

In a connected environment, metadata can be enriched automatically -- machine learning classifies images, natural language processing extracts keywords, linking to external databases provides context. In an air-gapped institution, all metadata creation is manual or generated by local tools. There is no external enrichment. The operator creates the metadata, and the quality of the metadata depends entirely on the operator's diligence.

This constraint makes metadata procedures more important, not less. The effort of creating good metadata at the time of data creation is small. The effort of retroactively adding metadata to years of untagged data is enormous and often impossible because the context has been lost.

### 3.3 Metadata Is a Promise to the Future

Every metadata record is a message from the present to the future. It says: "This is what this data is. This is where it came from. This is why it matters. This is how it relates to everything else." The integrity of that message depends on the metadata being created honestly, maintained faithfully, and preserved alongside the data it describes. Metadata that is fabricated, neglected, or allowed to drift from reality is worse than no metadata, because it creates false confidence.

## 4. System Model

### 4.1 The Metadata Schema

Every piece of data in the institution carries metadata organized in three layers:

**Layer 1: Core Metadata (Required for all data, no exceptions).**

| Field | Description | Example |
|---|---|---|
| `id` | Unique identifier. Format: `[domain]-[YYYYMMDD]-[sequence]` | `D6-20260215-0042` |
| `title` | Human-readable title of the data item | `Quarterly Archive Audit Report Q1 2026` |
| `created` | Creation date in ISO 8601 format | `2026-02-15` |
| `creator` | Who or what created the data | `operator` or `syslog-daemon` |
| `tier` | Data tier per D6-001 (1, 2, 3, or 4) | `1` |
| `format` | File format and version | `text/markdown; charset=utf-8` |
| `checksum` | SHA-256 hash of the content | `a1b2c3d4...` |
| `archived` | Date the item entered the archive | `2026-03-01` |
| `description` | One-paragraph summary of the content | Free text, 50-300 words |
| `retention` | Retention class per D6-011 | `permanent` or `5-year` or `transient` |

**Layer 2: Contextual Metadata (Required for Tier 1 and Tier 2; recommended for Tier 3).**

| Field | Description | Example |
|---|---|---|
| `domain` | The institutional domain this data belongs to | `6` |
| `project` | Associated project, if any | `archive-system-v2` |
| `related` | IDs of related data items | `[D6-20260210-0038, D6-20260212-0040]` |
| `supersedes` | ID of item this data replaces, if any | `D6-20260101-0015` |
| `provenance` | Origin chain: how this data came to exist | `Created during Q1 audit per D6-004` |
| `classification` | Security classification per SEC-001 | `institutional` or `restricted` |
| `keywords` | Controlled vocabulary terms for search | `[archive, audit, integrity, tier-1]` |
| `version` | Version number of the content | `1.0.0` |

**Layer 3: Technical Metadata (Automatically generated where possible; required for Tier 1).**

| Field | Description | Example |
|---|---|---|
| `size` | File size in bytes | `245760` |
| `encoding` | Character encoding | `UTF-8` |
| `permissions` | File permissions at archival | `640` |
| `original-filename` | Filename before any renaming | `audit-report-final.md` |
| `migration-history` | List of format migrations this item has undergone | `[]` or `[{from: "docx", to: "odt", date: "2028-06-15"}]` |
| `backup-locations` | Where backup copies of this item exist | `[local-backup-drive-1, offsite-vault]` |

### 4.2 Metadata File Format

Metadata is stored in a companion file alongside the data it describes. The companion file uses the same base name as the data file with a `.meta.yaml` extension:

```
2026-02-15_op_quarterly-archive-audit-report_v1.md
2026-02-15_op_quarterly-archive-audit-report_v1.meta.yaml
```

The YAML format is chosen for human readability and editability. YAML is a plain-text format that can be read and written with any text editor, requires no specialized tools, and is widely supported across programming languages. It is an approved archival format per D6-003.

Example metadata file:

```yaml
id: D6-20260215-0042
title: Quarterly Archive Audit Report Q1 2026
created: 2026-02-15
creator: operator
tier: 1
format: text/markdown; charset=utf-8
checksum: sha256:a1b2c3d4e5f6...
archived: 2026-03-01
description: >
  Comprehensive audit of the Tier 1 and Tier 2 archives for Q1 2026.
  Includes integrity verification results, naming convention compliance,
  storage consumption analysis, and corrective actions taken.
retention: permanent
domain: 6
project: null
related:
  - D6-20251215-0038
  - D6-20260101-0039
supersedes: null
provenance: Created during quarterly comprehensive audit per D6-004 Section 4.6
classification: institutional
keywords:
  - archive
  - audit
  - integrity
  - quarterly
version: 1.0.0
size: 245760
encoding: UTF-8
permissions: "640"
original-filename: audit-report-q1-2026.md
migration-history: []
backup-locations:
  - local-backup-drive-1
  - offsite-vault-a
```

### 4.3 Metadata Creation Procedures

Metadata is created at three points in the data lifecycle:

**At creation:** When the operator creates a new document, record, or artifact, they create the companion metadata file simultaneously. The core metadata fields are populated immediately. Contextual metadata fields are populated to the extent known. Technical metadata fields are generated automatically where tools exist or populated manually.

**At ingest:** When data enters the institution from external sources per D6-014, the ingest process requires metadata creation as a mandatory step. The operator classifies the data, assigns a tier, writes a description, and populates all core fields. Contextual fields are populated based on the ingest assessment.

**At archival transition:** When data moves from active storage to the archive per D6-004, the metadata is reviewed and completed. Any fields that were marked as unknown at creation must be resolved or explicitly documented as unresolvable. The `archived` date field is set.

### 4.4 Metadata Maintenance

Metadata is not static. It must be updated when:

- The data item is modified (new version): update `version`, `checksum`, `size`, and add a `related` reference to the previous version.
- The data item undergoes format migration per D6-005: update `format`, `checksum`, `size`, and add an entry to `migration-history`.
- The data item's tier classification changes: update `tier` and `retention`.
- Related items are created or discovered: update `related`.
- The data item is backed up to a new location: update `backup-locations`.
- An error in the metadata is discovered: correct the error and note the correction in the metadata audit log.

Every metadata modification must be reflected in the metadata catalogue (Section 4.6).

### 4.5 The Controlled Vocabulary

The `keywords` field uses a controlled vocabulary -- a defined set of terms that standardize how data is described. The controlled vocabulary prevents the problem of synonymous but different tags: one item tagged `backup` and another tagged `backups` and another tagged `data-backup`, all meaning the same thing but invisible to each other's searches.

The controlled vocabulary is maintained as a Tier 1 document in the archive. It is organized by domain and may be extended through the GOV-001 Tier 4 decision process (operational decisions). New terms are added when existing terms are insufficient. Existing terms are not renamed or removed; they are marked as deprecated and mapped to their replacement term.

The initial controlled vocabulary is seeded from the domain names, the data tier names, the operational process names, and the common institutional concepts. It will grow over time as the institution's data diversifies.

### 4.6 The Metadata Catalogue

The metadata catalogue is the central searchable index of all institutional metadata. It aggregates the individual `.meta.yaml` files into a single queryable repository. The catalogue is implemented as a flat-file database (a consolidated YAML or CSV file) that can be searched with standard text tools (grep, awk) or a simple script.

The catalogue is rebuilt from the individual metadata files during each quarterly metadata audit. Between audits, it is updated incrementally as new items are archived or existing metadata is modified.

The catalogue is stored as Tier 1 data and backed up per D6-006. However, per D6-004 Section 7.3, the catalogue is a derived index. If it is lost, it can be rebuilt from the individual metadata files. The individual files are the source of truth; the catalogue is a convenience for search and discovery.

### 4.7 Search and Discovery

Metadata enables finding data in the archive through several search methods:

**By field value.** Search for all items with a specific tier, domain, creator, date range, or retention class. Example: "Find all Tier 1 items in Domain 6 created in 2026." This is a structured query against the catalogue.

**By keyword.** Search the `keywords` field using the controlled vocabulary. Example: "Find all items tagged with `archive` and `audit`." The controlled vocabulary ensures consistent terminology.

**By relationship.** Follow the `related` and `supersedes` links to navigate from one item to its related items. This enables tracing the history and context of a specific piece of data.

**By full text.** Search the `description` and `title` fields for free-text terms. This is a fallback for when structured fields and keywords do not capture the needed concept.

**By provenance.** Search the `provenance` field to find all data items that originated from a specific process, event, or source. Example: "Find all items created during the 2027 annual audit."

### 4.8 The Metadata Audit

The metadata audit verifies that metadata is complete, accurate, and consistent. It occurs at two frequencies:

**Quarterly Metadata Audit (1-2 hours):**
- Select 50 items from the archive at random (stratified by tier: 20 Tier 1, 20 Tier 2, 10 Tier 3).
- For each item, verify that the companion `.meta.yaml` file exists and is well-formed.
- Verify that all core metadata fields are populated and non-empty.
- Verify that the `checksum` field matches the actual checksum of the data file.
- Verify that the `tier` and `retention` fields are consistent with each other per D6-011.
- Verify that `keywords` use only terms from the controlled vocabulary.
- Rebuild the metadata catalogue from individual files. Compare against the existing catalogue. Resolve discrepancies.
- Document audit results in the operational log.

**Annual Metadata Audit (4-8 hours):**
- Perform the quarterly audit at full scope: verify all Tier 1 metadata files, 50% of Tier 2, and 20% of Tier 3.
- Verify that every item in the archive has a companion metadata file. Items without metadata are flagged for remediation.
- Verify that the controlled vocabulary is current and that no items use deprecated terms without the mapping.
- Assess metadata quality: are descriptions informative or perfunctory? Are provenance chains complete? Are relationship links accurate?
- Identify metadata patterns that suggest systemic problems: many items with identical descriptions (copy-paste metadata), missing contextual fields on Tier 1 items, stale `backup-locations` fields.
- Produce an annual metadata health report. Archive it as Tier 1 data.

## 5. Rules & Constraints

- **R-D6-007-01:** Every data item in the archive must have a companion `.meta.yaml` file containing all core metadata fields populated. Items without metadata are non-compliant and must be remediated.
- **R-D6-007-02:** Metadata is created at the time of data creation or ingest. Retroactive metadata creation is a remediation activity, not a normal workflow. The goal is to never need it.
- **R-D6-007-03:** The `description` field must contain a substantive summary of the data item's content and purpose. Descriptions shorter than 50 words for Tier 1 items or 25 words for other tiers are flagged for review during the metadata audit.
- **R-D6-007-04:** The `keywords` field must use only terms from the controlled vocabulary. New terms are added to the vocabulary before being used in metadata.
- **R-D6-007-05:** The `checksum` field must be verified during every metadata audit. A checksum mismatch indicates either data corruption or metadata staleness. Both require investigation.
- **R-D6-007-06:** The metadata catalogue must be rebuilt from source files at least quarterly. The catalogue is never the sole record of any metadata.
- **R-D6-007-07:** Metadata modifications must be reflected in both the companion file and the catalogue. Inconsistencies between the two are detected during the quarterly audit.
- **R-D6-007-08:** The controlled vocabulary must be reviewed and updated at least annually. Deprecated terms must be mapped to replacement terms, not simply deleted.

## 6. Failure Modes

- **Metadata neglect.** The operator stops creating metadata because it feels like overhead. Files accumulate without companion metadata. Over years, large portions of the archive become undescribed and unfindable. Mitigation: R-D6-007-01 and R-D6-007-02 establish metadata creation as a mandatory, immediate activity. The quarterly audit detects gaps.

- **Perfunctory metadata.** Metadata is created but without care. Descriptions read "data file" or "see filename." Keywords are generic or absent. The metadata exists but provides no value. Mitigation: R-D6-007-03 establishes minimum quality standards. The annual audit assesses metadata quality, not just presence.

- **Vocabulary drift.** The controlled vocabulary is not maintained. Different operators use different terms for the same concept. Search returns inconsistent results. Mitigation: R-D6-007-04 requires vocabulary compliance. R-D6-007-08 requires annual vocabulary review.

- **Checksum staleness.** The data file is modified (perhaps during a format migration) but the checksum in the metadata is not updated. The metadata claims a checksum that does not match the file. This could mask real corruption because the checksum cannot be trusted. Mitigation: R-D6-007-05 requires checksum verification during every audit. D6-005 migration procedures require metadata updates.

- **Catalogue-source divergence.** The metadata catalogue drifts from the individual `.meta.yaml` files. Search against the catalogue returns stale or incorrect results. Mitigation: R-D6-007-06 requires quarterly catalogue rebuilds. R-D6-007-07 requires synchronized updates.

- **Over-engineering.** The metadata schema becomes so complex that creating metadata for a single file takes longer than creating the file itself. The operator rebels and stops creating metadata altogether. Mitigation: the three-layer schema is designed to be progressive. Only Layer 1 is mandatory for all items. Layer 2 is mandatory only for Tier 1 and Tier 2. Layer 3 is largely automatable. If the schema needs simplification in practice, revise it -- but do not abandon it.

## 7. Recovery Procedures

**7.1 If a large body of data is discovered without metadata:**
1. Declare a metadata remediation sprint. Prioritize by tier: Tier 1 items first, then Tier 2.
2. For each item without metadata, attempt to reconstruct the core fields from available evidence: filesystem timestamps (for `created`), file content (for `title`, `description`), directory location (for `tier`, `domain`), file format (for `format`).
3. For fields that cannot be reconstructed, use explicit unknowns: `description: "Content unknown; metadata created during remediation sprint [date]."` Honest ignorance is better than fabricated metadata.
4. Mark all remediated metadata files with a provenance note: `provenance: "Metadata created retroactively during remediation on [date]. Original context unavailable."`
5. Document the remediation effort: how many items were affected, how much metadata could be reconstructed, and what remains unknown.
6. Investigate why the metadata gap occurred and correct the process that allowed it.

**7.2 If the controlled vocabulary becomes inconsistent:**
1. Export all unique `keywords` values from all metadata files.
2. Compare against the current controlled vocabulary. Identify unauthorized terms.
3. For each unauthorized term, determine whether it should be added to the vocabulary or mapped to an existing term.
4. Update all metadata files that use unauthorized terms to use the correct vocabulary terms.
5. Rebuild the controlled vocabulary document and the metadata catalogue.

**7.3 If the metadata catalogue is lost or corrupted:**
1. The individual `.meta.yaml` files are the source of truth. The catalogue is derived.
2. Walk the archive directory structure. Collect all `.meta.yaml` files.
3. Parse each file and rebuild the catalogue.
4. Verify the rebuilt catalogue against a sample of individual files.
5. Investigate the cause of the loss. Ensure the catalogue is backed up per D6-006.

**7.4 If metadata quality has degraded systemically:**
1. This is a discipline problem, not a data problem. The metadata exists but is inadequate.
2. Conduct the annual metadata audit with enhanced scrutiny.
3. Identify the pattern: are descriptions inadequate? Keywords missing? Relationships not tracked?
4. Address the root cause: if the operator has insufficient time for metadata, evaluate whether the schema can be simplified. If the operator has insufficient motivation, revisit D6-001 Section 4.2 on why metadata matters.
5. Remediate the worst examples first: Tier 1 items with inadequate metadata are the highest priority.

## 8. Evolution Path

- **Years 0-5:** The metadata schema is new. The controlled vocabulary is growing rapidly as new types of data are catalogued. Expect frequent additions to the vocabulary. Expect the operator to develop shortcuts and templates for common metadata patterns. Document these templates -- they are operational knowledge.

- **Years 5-15:** The metadata corpus becomes a valuable asset in its own right. The ability to search and navigate the archive through metadata transforms how the institution accesses its own knowledge. The controlled vocabulary should stabilize. The metadata audit should be routine.

- **Years 15-30:** The metadata created in the early years proves its value as the original context fades from memory. Items that seemed self-explanatory when they were created now depend entirely on their metadata for comprehensibility. This is when the quality of early metadata decisions becomes apparent.

- **Years 30-50+:** The metadata may need its own migration as the YAML format ages. The principles of metadata -- structured description, controlled vocabulary, companion files, audit -- are format-independent. If YAML is superseded, migrate the metadata using the same procedures in D6-005 that govern any format migration.

- **Signpost for revision:** If the metadata schema consistently fails to capture important distinctions, if the controlled vocabulary is growing unmanageably, or if the audit burden is disproportionate to the value provided, this article should be revised through the GOV-001 Tier 3 decision process.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
Metadata is the part of data management that nobody wants to do. It is the digital equivalent of filing -- necessary, unglamorous, and immediately punished when neglected. I write this article knowing that the biggest threat to metadata quality is my own laziness.

The three-layer schema is a compromise between thoroughness and practicality. I considered a simpler schema (just title, date, tier, and description) and rejected it because it does not support search or relationship tracking. I considered a more elaborate schema (Dublin Core, METS, or PREMIS-style) and rejected it because the overhead would guarantee non-compliance. Three layers, with only the first universally mandatory, is the balance I can realistically maintain.

The companion file approach (`.meta.yaml` next to every data file) is chosen over a centralized database for one critical reason: survivability. If the catalogue is lost, the metadata survives because it is distributed alongside the data. If a directory is copied or moved, its metadata travels with it. If the institution's tools change, the YAML files are readable with any text editor. Centralized metadata databases are more searchable but are single points of failure. In a fifty-year institution, distributed resilience wins over centralized elegance.

The description field is the most important field in the schema. It is the field that a future operator will read when they do not understand what a file is. It is the field that transforms an opaque artifact into a comprehensible record. Every time I am tempted to write "see filename" in the description field, I will remember that I am writing for someone who has not seen the filename's context. The description must stand alone.

## 10. References

- ETH-001 -- Ethical Foundations of the Institution (Principle 2: Integrity Over Convenience; Principle 4: Longevity Over Novelty)
- CON-001 -- The Founding Mandate (documentation completeness, institutional continuity)
- SEC-001 -- Threat Model and Security Philosophy (data classification)
- OPS-001 -- Operations Philosophy (documentation-first principle, operational tempo)
- D6-001 -- Data Philosophy (data tiers, data lifecycle, "data without context is noise")
- D6-002 -- The Canonical Data Model (entity types that metadata describes)
- D6-003 -- Format Longevity Doctrine (YAML as approved format; format metadata)
- D6-004 -- Archive Management Procedures (directory structure, naming conventions, manifests)
- D6-005 -- Format Migration Procedures (metadata updates during migration)
- D6-006 -- Backup Verification and Test Restoration (metadata verification during restoration)
- D6-008 -- Storage Lifecycle Management (backup-locations metadata field)
- D6-011 -- Retention Schedules & Disposal Policy (retention field in metadata)
- D6-014 -- Data Ingest Procedures (metadata creation at ingest)
- GOV-001 -- Authority Model (decision tiers for vocabulary and schema changes)

---

---

# D6-008 -- Storage Lifecycle Management

**Document ID:** D6-008
**Domain:** 6 -- Data & Archives
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, SEC-001, OPS-001, D6-001, D6-003, D6-004, D6-005, D6-006, D6-007
**Depended Upon By:** D6-009, D6-011, D6-013. All articles involving hardware procurement, media replacement, or physical storage.

---

## 1. Purpose

This article defines how the institution manages the lifecycle of its physical storage media -- from procurement through active use, degradation monitoring, and eventual retirement. Storage media are the physical substrate on which all institutional data exists. Every file, every archive, every backup depends on media that is subject to physical degradation, mechanical failure, electronic wear, and technological obsolescence. A storage lifecycle management program ensures that no data is ever lost because the media it was stored on failed without warning.

D6-001 establishes the principle of physical sovereignty over data. This article operationalizes that sovereignty by providing the procedures for understanding how each type of storage media ages, how to detect degradation before it causes data loss, when to migrate data from aging media to fresh media, and how to plan the procurement of replacement media before it is urgently needed.

The institution uses multiple types of storage media, each with distinct characteristics, advantages, and failure modes. This article addresses the lifecycle of hard disk drives (HDD), solid-state drives (SSD), optical media (Blu-ray, DVD), and magnetic tape. It provides specific guidance for each type while establishing a unified management framework. The operator does not need to be a storage engineer. They need to follow the monitoring schedules, recognize the warning signs, and act on the replacement timelines defined here.

## 2. Scope

This article covers:

- Lifecycle characteristics of HDD, SSD, optical media, and magnetic tape.
- Degradation detection methods for each media type.
- Health monitoring schedules: what to check, how often, and what the results mean.
- Migration triggers: when to move data from aging media to replacement media.
- Procurement planning: how to ensure replacement media is available before it is needed.
- Media retirement procedures: how to decommission and securely dispose of retired media.
- The media inventory: tracking all storage media in the institution.

This article does not cover:

- Data format migration between file formats (see D6-005).
- Logical data organization on storage media (see D6-004).
- Backup procedures and verification (see D6-006).
- Disaster recovery from complete media failure (see D6-013).
- The philosophy of data persistence (see D6-001).
- Encryption of storage media (see SEC-003).

## 3. Background

### 3.1 No Storage Medium Is Permanent

Every storage medium degrades. Hard drives have mechanical components that wear. SSDs have flash cells that can only endure a finite number of write cycles. Optical discs degrade through chemical reactions with air and light. Magnetic tape loses signal strength over time. The question is not whether media will fail but when, and whether the institution will be prepared.

The popular belief that "digital storage is forever" is one of the most dangerous myths in information technology. Analog formats -- stone tablets, papyrus, vellum -- have survived millennia. No digital storage medium has demonstrated a lifespan beyond a few decades under real-world conditions. The institution's fifty-year horizon exceeds the reliable lifespan of every currently available storage technology. Multiple generations of media will be used and retired over the institution's life.

### 3.2 Failure Patterns Differ by Media Type

Different storage technologies fail in different ways. HDDs tend to fail suddenly after a period of increasing SMART errors. SSDs degrade gradually as write endurance is consumed, then fail suddenly when critical thresholds are crossed. Optical media degrades slowly through oxidation and delamination, often without any warning until the disc becomes unreadable. Magnetic tape loses signal strength continuously, with acceleration in hostile environments.

Understanding these patterns is essential because the monitoring strategy must match the failure mode. Monitoring an SSD for HDD-style mechanical warnings will miss the actual degradation indicators. This article provides media-specific monitoring guidance.

### 3.3 Procurement Lead Time

In an air-gapped, off-grid institution, replacement storage media cannot be ordered from an online retailer and delivered the next day. Procurement requires identifying the need, specifying the media, obtaining it from a physical retail location or through a mail-order process, and importing it through the quarantine process. This takes time -- days to weeks. If a storage medium fails before its replacement has been procured, the institution faces a redundancy gap during which data exists in fewer copies than the minimum required by D6-001 R-D6-02. Procurement planning is therefore a proactive, scheduled activity, not a reactive one.

### 3.4 Technological Obsolescence

Storage interfaces change. IDE gave way to SATA. SATA may give way to NVMe or its successors. The institution must plan for the possibility that the interfaces used by its current media will not be available on future hardware. Media migration per this article addresses physical media replacement within a compatible technology. When the technology itself changes, the migration becomes a larger project involving D6-005 (format), this article (media), and potentially Domain 4 and 5 articles (hardware and platform).

## 4. System Model

### 4.1 Media Lifecycle Characteristics

**Hard Disk Drives (HDD):**
- **Expected lifespan:** 5-7 years under continuous operation; up to 10 years under light or intermittent use. Archival HDDs stored powered-off should be powered on and verified at least every 12 months to prevent stiction (heads sticking to platters).
- **Failure indicators:** Increasing reallocated sector count, rising pending sector count, increasing seek error rate, audible clicking or grinding, increasing spin-up time. These are reported by SMART (Self-Monitoring, Analysis, and Reporting Technology).
- **Degradation pattern:** Often predictable through SMART data. A drive with zero reallocated sectors that suddenly develops five reallocated sectors in one month is on a degradation curve. The time between the first SMART warning and catastrophic failure varies but is often 3-12 months.
- **Best suited for:** Active operational storage, large backup sets, archival storage with regular power-on verification.
- **Not suited for:** Only copy of critical data without SMART monitoring.

**Solid-State Drives (SSD):**
- **Expected lifespan:** 5-10 years depending on write volume. The primary wear mechanism is flash cell degradation from write/erase cycles. Read operations do not significantly contribute to wear. Consumer SSDs typically endure 150-600 TBW (terabytes written); enterprise SSDs significantly more.
- **Failure indicators:** Increasing media wearout indicator, rising reallocated NAND block count, increasing uncorrectable error count. Reported by SMART with SSD-specific attributes. Also monitor the drive's reported remaining endurance percentage where available.
- **Degradation pattern:** Gradual write endurance consumption followed by abrupt failure. SSDs often provide less warning than HDDs before catastrophic failure. The SMART endurance indicator is the primary warning.
- **Best suited for:** Active operational storage where speed matters, OS boot drives, frequently accessed databases.
- **Not suited for:** Long-term archival storage, because flash cells can lose charge over years of powered-off storage (data retention without power is typically 1-2 years for consumer SSDs at room temperature, longer for enterprise).

**Optical Media (Blu-ray, DVD):**
- **Expected lifespan:** 5-20 years for standard recordable media (DVD-R, BD-R). Potentially 50-100+ years for archival-grade media (M-DISC or equivalent). Lifespan is highly dependent on storage conditions.
- **Failure indicators:** Increasing read error rates during verification, visible discoloration or cloudiness, delamination at the edges. Optical media degrades from the outside edge inward.
- **Degradation pattern:** Slow, continuous chemical degradation. No electronic warning system. Degradation is detected only through periodic read verification. By the time a disc is unreadable, the degradation is complete.
- **Best suited for:** Write-once archival storage, offsite backup copies, permanent preservation of Tier 1 data.
- **Not suited for:** Active data that requires frequent modification. Any use where the inability to monitor degradation electronically is unacceptable.

**Magnetic Tape (LTO or equivalent):**
- **Expected lifespan:** 15-30 years for modern LTO tape stored in controlled conditions. Lifespan is significantly affected by temperature and humidity.
- **Failure indicators:** Increasing read error rates, tape drive cleaning requests becoming more frequent, visible damage to tape edges, sticky-shed syndrome (binder deterioration causing tape to stick to itself or to the read head).
- **Degradation pattern:** Slow, continuous signal strength loss. Accelerated by poor storage conditions. Periodic verification by reading the tape is the primary detection method.
- **Best suited for:** Large-volume archival storage, deep backup, high-capacity offsite storage.
- **Not suited for:** Active data requiring random access. Environments without climate control.

### 4.2 The Media Inventory

The institution maintains a media inventory -- a Tier 1 document that records every storage device in use, in reserve, and retired. Each entry contains:

- **Media ID:** A unique physical identifier (written on the device with permanent marker or a label).
- **Media type:** HDD, SSD, optical, tape.
- **Manufacturer and model:** Specific make and model.
- **Capacity:** Total storage capacity.
- **Serial number:** Manufacturer's serial number.
- **Date acquired:** When the media was obtained.
- **Date placed in service:** When the media began storing institutional data.
- **Service role:** What the media is used for (primary archive, backup, offsite, operational).
- **Data stored:** Summary of what data resides on this media (not a file list, but a category description: "Tier 1 archive 2026-2028," "Monthly backups Q1 2026").
- **Health status:** Current health assessment (healthy, degrading, replacement scheduled, retired).
- **Last health check:** Date and results of the most recent health check.
- **Next scheduled check:** Date of the next health check.
- **Retirement date:** Estimated or actual date of retirement.

### 4.3 Health Monitoring Schedule

**HDDs -- Monthly:**
- Run `smartctl -a /dev/[device]` (or equivalent) and record key metrics: Reallocated_Sector_Ct, Current_Pending_Sector, Offline_Uncorrectable, Spin_Retry_Count, Command_Timeout.
- Compare current values against the previous month. Any increase in reallocated or pending sectors is a warning.
- Record temperature: drives consistently above 45C have accelerated failure rates.
- For powered-off archival HDDs: power on, run SMART extended self-test, verify data readability by running checksums on a sample of files, then power off.

**SSDs -- Monthly:**
- Run `smartctl -a /dev/[device]` and record key metrics: Wear_Leveling_Count (or equivalent endurance indicator), Reallocated_Nand_Blk_Cnt, Uncorrectable_Error_Cnt.
- Note the reported percentage of write endurance remaining.
- Compare against previous months to determine write endurance consumption rate.
- Project the date at which write endurance will reach 10% remaining. This is the migration planning trigger.

**Optical Media -- Semi-annually:**
- For each optical disc in the archive, perform a read verification: read the entire disc and compare against the stored checksums.
- Visually inspect the disc for discoloration, cloudiness, or delamination.
- Note the read error rate. An increasing error rate across consecutive checks indicates degradation.
- For archival-grade media (M-DISC or equivalent): annual verification may be sufficient if the media is stored in controlled conditions.

**Magnetic Tape -- Annually:**
- For each tape, perform a full read pass and compare against stored checksums.
- Monitor the tape drive's error counter during the read. Increasing error rates indicate tape degradation.
- Inspect the tape cartridge for visible damage, dust, or humidity exposure.
- For tapes in controlled storage: the annual verification is sufficient. For tapes in uncontrolled storage: increase to semi-annual.

### 4.4 Migration Triggers

Data must be migrated from a storage medium to a replacement when any of the following conditions is met:

1. **SMART warning threshold.** Any increase in reallocated sectors (HDD) or uncorrectable errors (SSD) beyond baseline. This indicates active degradation. Begin migration planning immediately. Complete migration within 30 days.

2. **Endurance exhaustion approaching.** SSD write endurance below 20% remaining. Plan migration. Execute when endurance reaches 10%.

3. **Age threshold.** HDD in continuous service for 5 years. SSD in continuous service for 7 years. These are conservative thresholds. Media may last longer, but the risk-reward balance shifts against continued use.

4. **Optical degradation.** Optical media showing measurably increased read error rates compared to the first verification after burning. Begin migration. Complete before the next semi-annual check.

5. **Tape degradation.** Tape showing increased read error rates or sticky-shed symptoms. Begin migration immediately.

6. **Technology obsolescence.** The media's interface or technology is no longer supported by available hardware. Plan migration to current technology before the last compatible hardware fails or becomes unobtainable.

7. **Checksum failure.** Any data on the medium fails integrity verification. This is an immediate migration trigger. The data must be recovered from this medium or from a backup and written to a healthy medium.

### 4.5 Media Migration Procedure

When a migration trigger is met:

1. **Procure replacement media.** If replacement media is not already in reserve, procure it per Section 4.6. Do not begin migration until the replacement is available and verified as healthy.

2. **Verify source data.** Before migration, verify the integrity of the data on the source medium using checksums. This ensures you are copying good data, not propagating corruption.

3. **Copy data to replacement.** Use a tool that verifies each read and write (e.g., `rsync --checksum` or `dd` with subsequent verification). Do not use tools that report success without verification.

4. **Verify destination data.** After copying, run full checksum verification on the destination medium. Compare against the source manifests per D6-004.

5. **Update the archive.** If the migration changes the physical location of archive data, update the metadata catalogue and the media inventory.

6. **Retain the source.** Keep the source medium available for 90 days after migration in case errors are discovered in the destination copy. After 90 days, the source may be retired.

7. **Retire the source medium.** Per Section 4.7.

8. **Document the migration.** Record in the operational log and the media inventory: source medium ID, destination medium ID, data migrated, migration date, verification results.

### 4.6 Procurement Planning

The institution maintains a reserve of at least one replacement unit for each media type in active use. This reserve is the buffer that prevents emergency procurement under time pressure.

**Procurement planning procedure:**

1. **Annual procurement assessment.** During the annual storage review (coinciding with the annual archive audit per D6-004), assess:
   - How many media units of each type are in active service?
   - How many are in reserve?
   - How many are approaching migration triggers?
   - What is the projected media consumption for the next 12-24 months based on data growth and media aging?

2. **Procurement list.** Based on the assessment, create a procurement list for the next 12 months. The list specifies: media type, capacity, quantity, and target acquisition date.

3. **Technology review.** Before procuring, review whether the current media types are still the best choice. Has a new technology matured sufficiently for institutional use? Is the current technology approaching obsolescence? This review informs whether to procure more of the same or to begin a technology transition.

4. **Acquisition.** Procure media from the procurement list. All media entering the institution passes through the quarantine process. New media is tested before being placed in the reserve: for HDDs and SSDs, run a full surface scan and SMART extended self-test. For optical media, verify blank readiness. For tape, verify the drive can read and write the tape without errors.

5. **Reserve maintenance.** The reserve is stored in appropriate conditions: cool, dry, and protected from physical damage. HDDs in reserve are powered on and tested semi-annually. SSDs in reserve are powered on and verified semi-annually (to prevent charge loss in flash cells). Optical media in reserve is stored in protective cases away from light. Tape in reserve is stored vertically in a climate-controlled environment.

### 4.7 Media Retirement Procedures

When a storage medium is retired:

1. **Verify data migration.** Confirm that all data from the retiring medium has been successfully migrated to replacement media and verified per Section 4.5.

2. **Secure erasure.** If the medium contained any data classified above Tier 4, perform secure erasure:
   - For HDDs: use `shred` or `nwipe` to perform a multi-pass overwrite. For maximum assurance, physically destroy the platters.
   - For SSDs: use the manufacturer's secure erase command (ATA Secure Erase). Note that standard file deletion and even multi-pass overwrites are unreliable for SSDs due to wear leveling. For maximum assurance, physically destroy the flash chips.
   - For optical media: physically destroy the disc (shred, break, or scratch the data surface thoroughly).
   - For magnetic tape: degauss the tape, then physically destroy if warranted by the data's classification.

3. **Update the inventory.** Mark the medium as retired in the media inventory. Record the retirement date, the reason, and the disposal method.

4. **Physical disposal.** Dispose of the destroyed media in accordance with local waste regulations. Do not donate, sell, or recycle storage media that contained institutional data unless destruction is verified.

## 5. Rules & Constraints

- **R-D6-008-01:** The media inventory must be maintained as a Tier 1 document. Every storage medium in the institution must be inventoried. Uninventoried media must not store institutional data.
- **R-D6-008-02:** Health monitoring must be performed on the schedule defined in Section 4.3. Deferred health checks are recorded as operational incidents.
- **R-D6-008-03:** When a migration trigger defined in Section 4.4 is met, migration must be initiated within 30 days. Migration may not be indefinitely deferred.
- **R-D6-008-04:** The institution must maintain a reserve of at least one replacement unit for each active media type. Operating without reserve is an operational risk that must be documented and remediated at the earliest opportunity.
- **R-D6-008-05:** All retired storage media that contained data above Tier 4 must be securely erased before disposal. Disposal without secure erasure is a security violation per SEC-001.
- **R-D6-008-06:** Data must not be stored solely on SSD media that will be powered off for extended periods without periodic verification. SSD flash charge loss can cause data loss in as little as one year for consumer SSDs at elevated temperatures.
- **R-D6-008-07:** The annual procurement assessment must be completed and documented. Operating without procurement planning creates emergency dependencies on timely acquisition.
- **R-D6-008-08:** When a media type used by the institution is approaching technological obsolescence, migration planning to a current technology must begin before the last compatible hardware becomes unavailable.

## 6. Failure Modes

- **Monitoring neglect.** The operator stops performing health checks. A drive develops SMART warnings that go unnoticed. The drive fails. Data exists only on this drive because backup verification has also lapsed. Mitigation: R-D6-008-02 mandates monitoring. The operational tempo in OPS-001 schedules it. Cross-reference with D6-006 backup verification prevents single-copy exposure.

- **Procurement delay.** A migration trigger fires, but replacement media is not in reserve and cannot be obtained quickly. The institution must continue using degrading media while waiting for the replacement. Mitigation: R-D6-008-04 requires a reserve. The annual procurement assessment projects need.

- **SSD charge loss.** An SSD used for archival storage is powered off for two years. When powered on, portions of the data are corrupted because flash cells lost charge. Mitigation: R-D6-008-06 prohibits relying on powered-off SSDs without verification. The semi-annual reserve maintenance schedule in Section 4.6 catches this for reserve drives.

- **Optical media silent degradation.** Optical discs stored in a hot or humid environment degrade faster than expected. The semi-annual check reveals that multiple discs are unreadable. Mitigation: the semi-annual optical verification catches degradation. Proper storage conditions (cool, dry, dark) extend lifespan. Using archival-grade media (M-DISC) provides additional resilience.

- **Technology cliff.** A storage interface (e.g., SATA) is discontinued. The institution's drives still work, but no new hardware supports them. When the last compatible controller fails, all drives become inaccessible. Mitigation: R-D6-008-08 requires proactive migration when obsolescence approaches. The annual technology review in procurement planning monitors this.

- **Incomplete retirement.** A drive is retired but not securely erased. It is disposed of casually. Someone recovers the drive and accesses institutional data. Mitigation: R-D6-008-05 requires secure erasure. The retirement procedure in Section 4.7 includes physical destruction for high-sensitivity media.

## 7. Recovery Procedures

**7.1 If a drive fails without warning:**
1. Assess the scope: what data was on the failed drive? Is it the only copy?
2. If other copies exist (backup per D6-006, other media per D6-001 R-D6-02): restore from the verified copy. Proceed to media replacement per Section 4.5.
3. If no other copy exists: attempt data recovery. For HDDs, professional data recovery services may retrieve data from failed drives but this requires breaking the air gap for the recovery service. Assess whether the data's value justifies the security implications. For SSDs with controller failure, some recovery is possible by reading flash chips directly, but this is specialized work.
4. Document the failure as a Tier 1 incident. Investigate why the data existed on only one medium (violation of R-D6-02 per D6-001) and correct the process.

**7.2 If health monitoring reveals a drive in early degradation:**
1. This is a success -- degradation detected before failure. This is why we monitor.
2. Procure replacement media if not already in reserve.
3. Begin migration per Section 4.5 at the earliest opportunity.
4. Continue monitoring the degrading drive at increased frequency (weekly instead of monthly) until migration is complete.
5. After migration, retire the degrading drive per Section 4.7.

**7.3 If optical media is discovered to be degraded during verification:**
1. Determine the extent: is one disc affected or multiple?
2. If the data on the degraded disc exists elsewhere (another copy, on-disk backup): verify the other copy and burn a replacement disc.
3. If partial reads are possible: recover as much data as possible. Use tools designed for error-tolerant reading (e.g., `ddrescue` adapted for optical drives).
4. Investigate the storage conditions. If environmental factors caused accelerated degradation, improve storage conditions for all remaining optical media.
5. Document the incident and update the media inventory.

**7.4 If a technology transition becomes necessary:**
1. This is a planned project, not a crisis, if detected early through the annual technology review.
2. Identify the scope: how much data is on the outgoing technology?
3. Procure sufficient new-technology media to accommodate all data plus growth.
4. Test the new media thoroughly before trusting it with production data.
5. Migrate data in phases, starting with Tier 1 data (highest priority for protection).
6. Verify each phase before proceeding to the next.
7. Retain the old-technology media for 180 days after complete migration (same grace period as format migration per D6-005).
8. Document the technology transition comprehensively. This documentation will be valuable when the next transition occurs, potentially decades later.

## 8. Evolution Path

- **Years 0-5:** The initial media inventory is established. The monitoring routine is being built into operational habits. The first media procurements fill out the reserve. Health monitoring data accumulates but does not yet show long-term trends.

- **Years 5-15:** The first media retirements occur as the founding-era drives reach their age thresholds. The monitoring data now covers a meaningful span and can inform predictions. The first technology transition may occur as storage interfaces evolve. The procurement planning process proves its value or reveals its gaps.

- **Years 15-30:** Multiple generations of media have been used and retired. The media inventory becomes a historical record of the institution's storage evolution. The monitoring procedures may need adaptation for media types that do not exist today. The principles -- monitor, detect, migrate before failure -- remain constant.

- **Years 30-50+:** The institution's data has survived across numerous media generations. The lifecycle management procedures are the reason. The specific media types in Section 4.1 will be obsolete -- replaced by technologies that are not yet invented. The replacement technology should be evaluated against the same criteria: expected lifespan, failure indicators, degradation pattern, monitoring method. This article provides the framework; future operators provide the specifics.

- **Signpost for revision:** If a new storage technology emerges that does not fit the lifecycle model in this article (e.g., a technology with radically different degradation characteristics), or if the monitoring schedules prove inadequate based on actual failure experience, this article should be revised through the GOV-001 Tier 3 decision process.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
Storage media is the most physical, most tangible part of data management. You can hold a hard drive. You can see a disc. You can weigh a tape cartridge. And yet, the data on these physical objects is invisible and fragile in a way that the objects themselves are not. A hard drive can sit on a shelf looking perfectly fine while its data rots away inside. This disconnect between appearance and reality is why monitoring matters. You cannot tell by looking at a drive whether its data is intact. You have to ask the drive, systematically and regularly.

The SSD charge-loss issue is the one that most concerns me for this institution. SSDs are fast, silent, and reliable during active use. They are terrible for powered-off long-term storage because the electrical charge in flash cells dissipates over time. This is not a theoretical risk -- it is a documented physical phenomenon. Consumer SSDs at elevated temperatures can lose data in as little as a year without power. Enterprise SSDs are better but not immune. The institution must treat SSDs as active-use media, not archival media. For long-term storage, HDDs (with periodic power-on verification) and archival-grade optical media are safer choices.

The procurement reserve requirement may seem excessive when budgets are tight. One spare drive of each type, sitting on a shelf, doing nothing. But that spare drive is the difference between a scheduled migration and an emergency. When a drive starts failing, the spare means you can begin migration immediately instead of scrambling to acquire a replacement. The spare has already paid for itself the first time you need it and it is there.

I want future operators to understand one thing above all: storage media is consumable. It is not an asset you buy once and own forever. It is a consumable that you buy, use, monitor, replace, and discard on a schedule measured in years. Budgeting for storage media is like budgeting for electricity or food. It is a recurring, permanent cost of institutional operation. The moment the institution treats storage media as a one-time purchase, it has begun the countdown to data loss.

## 10. References

- ETH-001 -- Ethical Foundations of the Institution (Principle 1: Sovereignty; Principle 2: Integrity Over Convenience; Principle 4: Longevity Over Novelty)
- CON-001 -- The Founding Mandate (air-gap mandate, physical sovereignty, institutional continuity)
- SEC-001 -- Threat Model and Security Philosophy (physical threats, data integrity threats, secure disposal)
- SEC-003 -- Cryptographic Key Management (encryption of storage media)
- OPS-001 -- Operations Philosophy (operational tempo, maintenance schedules, documentation-first principle)
- D6-001 -- Data Philosophy (R-D6-02: minimum three copies on two media types; physical sovereignty; data tiers)
- D6-003 -- Format Longevity Doctrine (format considerations during media migration)
- D6-004 -- Archive Management Procedures (archive integrity manifests, directory structure)
- D6-005 -- Format Migration Procedures (coordination when media migration coincides with format migration)
- D6-006 -- Backup Verification and Test Restoration (backup media health verification)
- D6-007 -- Metadata Standards and Procedures (backup-locations metadata field updates during migration)
- D6-009 -- Data Migration Doctrine (broader migration context)
- D6-013 -- Disaster Recovery for Data Systems (recovery from complete media failure)
- GOV-001 -- Authority Model (decision tiers for technology transitions and procurement)

---