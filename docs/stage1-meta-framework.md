# Stage 1 Meta-Framework: Unifying Standards for the Holm Documentation Institution

**Document ID:** META-00-ART-001-v1.0
**Status:** Published
**Created:** 2026-02-16
**Author:** Founder (Tim)
**Scope:** All 20 documentation domains
**Purpose:** Establish the canonical standards that bind every domain, article, workflow, and publishing decision across the entire documentation institution. This document is the root authority. When any domain-specific convention conflicts with this framework, this framework wins.

---

## Table of Contents

1. [Naming Conventions](#1-naming-conventions)
2. [Versioning Rules](#2-versioning-rules)
3. [Article Template (Canonical)](#3-article-template-canonical)
4. [Publishing Lifecycle](#4-publishing-lifecycle)
5. [Review Workflow](#5-review-workflow)
6. [Commentary Workflow](#6-commentary-workflow)
7. [Master Domain Index](#7-master-domain-index)
8. [Dependency Graph](#8-dependency-graph)
9. [Website Structure](#9-website-structure)
10. [50-Year Continuity Rules](#10-50-year-continuity-rules)

---

## 1. NAMING CONVENTIONS

### 1.1 Article ID Format

Every article in the institution receives a globally unique identifier with the following structure:

```
{DOMAIN}-{SEQUENCE}-ART-{NUMBER}-v{MAJOR}.{MINOR}
```

**Anatomy of an Article ID:**

| Segment       | Format    | Description                                              | Example   |
|---------------|-----------|----------------------------------------------------------|-----------|
| DOMAIN        | 3-4 chars | Domain code (see Section 1.2)                            | `INFR`    |
| SEQUENCE      | 2 digits  | Domain-internal ordering bucket (01-99)                  | `01`      |
| ART           | literal   | Fixed string. Identifies this as an article.             | `ART`     |
| NUMBER        | 3 digits  | Sequential article number within the domain+sequence     | `001`     |
| v{MAJOR}.{MINOR} | digits | Version (see Section 2). Patch level is omitted in IDs.| `v1.0`    |

**Full example:** `INFR-01-ART-001-v1.0`
This is Infrastructure domain, bucket 01, article 001, version 1.0.

**Rules:**

- Article IDs are **immutable once assigned**. An article never changes its base ID (`INFR-01-ART-001`); only the version suffix changes.
- Sequence buckets (the `01` segment) group articles within a domain by sub-topic. Each domain defines its own bucket numbering in its domain index. Bucket `00` is reserved for domain-level meta-articles (the domain's own index, conventions, and scope definition).
- Article numbers are **never reused**. If article 007 is permanently archived or deleted, number 007 is retired. The next article is 008.
- The version suffix in the ID always reflects the **current published version**. Previous versions are referenced by appending the full semver: `INFR-01-ART-001-v1.0.3`.

### 1.2 Domain Codes

All 20 domains and their assigned codes:

| Code   | Domain Name                        |
|--------|------------------------------------|
| `META` | Meta-Documentation & Standards     |
| `INFR` | Infrastructure & Hardware          |
| `NETW` | Networking & Communications        |
| `POWR` | Power Systems & Energy             |
| `STOR` | Data Storage & Backup              |
| `SECR` | Security & Access Control          |
| `OPSYS`| Operating Systems & Software       |
| `APPL` | Applications & Services            |
| `SITE` | Website & Publishing (holm.chat)   |
| `GOVR` | Governance & Decision-Making       |
| `FINC` | Finance & Resource Management      |
| `LRNG` | Learning & Knowledge Management    |
| `WTHR` | Weather & Environmental Systems    |
| `FOOD` | Food Production & Preservation     |
| `SHLT` | Shelter & Physical Plant           |
| `WATR` | Water Systems                      |
| `TOOL` | Tools & Workshop                   |
| `HLTH` | Health & Medical                   |
| `COMM` | Community & Social Systems         |
| `LGCY` | Legacy & Succession Planning       |

**Rules for domain codes:**

- Codes are 3-4 uppercase ASCII characters. No digits, no special characters.
- Codes are permanent. Once assigned, a domain code never changes, even if the domain is renamed.
- If a 21st domain is ever needed, it must be proposed through the Governance domain (`GOVR`) and ratified by adding it to this document. The new code must not collide with any existing code or any plausible abbreviation of an existing domain.

### 1.3 Version Numbering Scheme

See Section 2 for full versioning rules. Summary:

```
v{MAJOR}.{MINOR}.{PATCH}
```

- **MAJOR**: Structural or conceptual rewrite. The article's thesis, scope, or fundamental approach changed.
- **MINOR**: Substantive content addition or modification. New sections, corrected procedures, updated data.
- **PATCH**: Typo fixes, formatting corrections, clarified wording with no change in meaning.

In Article IDs, the patch level is omitted (shown as `v1.0`). In version history and file names, the full semver is used (`v1.0.0`).

### 1.4 File Naming Rules

Every article is stored as a single Markdown file. The file name mirrors the Article ID with the following conventions:

```
{domain}-{sequence}-art-{number}-v{major}.{minor}.{patch}.md
```

**Rules:**

- All lowercase.
- Hyphens as separators (never underscores, never spaces).
- The version in the filename is the **full semver** including patch level.
- Example: `infr-01-art-001-v1.0.0.md`

**Directory structure:**

```
/docs/
  {domain-code}/
    {sequence}/
      {filename}.md
```

Example path: `/docs/infr/01/infr-01-art-001-v1.0.0.md`

**Supplementary files** (diagrams, images, data tables) are stored alongside the article in a subdirectory:

```
/docs/infr/01/infr-01-art-001-assets/
  diagram-power-topology.png
  table-hardware-inventory.csv
```

Asset filenames are descriptive, lowercase, hyphen-separated, and prefixed with their type (`diagram-`, `table-`, `photo-`, `schematic-`).

**Previous versions** are stored in a `_versions/` subdirectory:

```
/docs/infr/01/_versions/
  infr-01-art-001-v1.0.0.md
  infr-01-art-001-v1.1.0.md
```

The current version lives at the top level. All prior versions are moved into `_versions/` when superseded.

### 1.5 URL Slug Conventions for holm.chat

Public-facing URLs on the website follow a human-readable hierarchy:

```
https://holm.chat/docs/{domain-slug}/{article-slug}/
```

**Domain slugs** are lowercase, full-word versions of the domain name:

| Domain Code | URL Domain Slug          |
|-------------|--------------------------|
| `META`      | `meta`                   |
| `INFR`      | `infrastructure`         |
| `NETW`      | `networking`             |
| `POWR`      | `power`                  |
| `STOR`      | `storage`                |
| `SECR`      | `security`               |
| `OPSYS`     | `operating-systems`      |
| `APPL`      | `applications`           |
| `SITE`      | `website`                |
| `GOVR`      | `governance`             |
| `FINC`      | `finance`                |
| `LRNG`      | `learning`               |
| `WTHR`      | `weather`                |
| `FOOD`      | `food`                   |
| `SHLT`      | `shelter`                |
| `WATR`      | `water`                  |
| `TOOL`      | `tools`                  |
| `HLTH`      | `health`                 |
| `COMM`      | `community`              |
| `LGCY`      | `legacy`                 |

**Article slugs** are derived from the article title:

- Lowercase.
- Spaces replaced with hyphens.
- Strip all punctuation except hyphens.
- Maximum 60 characters. Truncate at a word boundary.
- Must be unique within the domain. If a collision occurs, append `-2`, `-3`, etc.

**Example URLs:**

```
https://holm.chat/docs/infrastructure/solar-panel-array-maintenance/
https://holm.chat/docs/power/battery-bank-rotation-schedule/
https://holm.chat/docs/meta/naming-conventions-and-standards/
```

**Version-specific URLs** append the version:

```
https://holm.chat/docs/infrastructure/solar-panel-array-maintenance/v1.0.0/
```

The bare URL (without version) always points to the current published version.

**Additional URL paths:**

```
https://holm.chat/docs/                          -> Master index
https://holm.chat/docs/{domain-slug}/            -> Domain index
https://holm.chat/docs/{domain-slug}/{slug}/history/  -> Version history
https://holm.chat/docs/{domain-slug}/{slug}/comments/ -> Commentary thread
https://holm.chat/changelog/                     -> Global changelog
https://holm.chat/about/                         -> About the institution
https://holm.chat/graph/                         -> Dependency graph visualization
```

---

## 2. VERSIONING RULES

### 2.1 Semantic Versioning for Articles

Every article uses semantic versioning: `v{MAJOR}.{MINOR}.{PATCH}`

All articles begin life at `v0.1.0` (initial draft). The first published version is `v1.0.0`.

### 2.2 When to Increment Each Level

**MAJOR (v1.0.0 -> v2.0.0):**

Increment the major version when:

- The article's core thesis, conclusion, or recommendation changes.
- The scope of the article is fundamentally redefined (broader or narrower).
- The article is restructured so heavily that someone familiar with the old version would not recognize the flow.
- A procedure is replaced with an entirely different procedure (not just updated steps).
- The article is rewritten for a different audience.

A major version increment resets minor and patch to zero.

**MINOR (v1.0.0 -> v1.1.0):**

Increment the minor version when:

- A new section is added.
- An existing section is substantially rewritten (but the overall thesis is unchanged).
- Procedures are updated with new steps, changed parameters, or corrected sequences.
- New failure modes, recovery procedures, or constraints are documented.
- Data tables, diagrams, or references are meaningfully updated.
- Cross-references to other articles are added or corrected.

A minor version increment resets patch to zero.

**PATCH (v1.0.0 -> v1.0.1):**

Increment the patch version when:

- Typographical errors are corrected.
- Formatting is improved without changing content.
- Wording is clarified without changing meaning.
- Broken links are fixed.
- Metadata is corrected (dates, codes, tags).

### 2.3 Article States

Every article exists in exactly one of six states at any time:

| State        | Code | Description                                                                 |
|--------------|------|-----------------------------------------------------------------------------|
| **Draft**    | `D`  | Article is being written. Not visible on the public website. Version is `v0.x.y`. |
| **Review**   | `R`  | Article is complete and undergoing self-review or cross-domain review.       |
| **Approved** | `A`  | Article has passed review and is cleared for publishing.                     |
| **Published**| `P`  | Article is live on holm.chat. Version is `v1.0.0` or higher.                |
| **Revised**  | `V`  | A published article is being revised. A new draft version exists alongside the current published version. |
| **Archived** | `X`  | Article is no longer current. Remains accessible but marked as superseded or obsolete. |

**State encoding in metadata:**

The article's front-matter includes a `status` field with the state code. Example:

```yaml
status: P
```

### 2.4 Version History Requirements

Every article must maintain a version history table in its front-matter or a dedicated section. The table includes:

| Version  | Date       | State | Summary of Changes                        |
|----------|------------|-------|-------------------------------------------|
| v0.1.0   | 2026-02-16 | D     | Initial draft created.                    |
| v0.2.0   | 2026-02-18 | D     | Added failure modes section.              |
| v1.0.0   | 2026-02-20 | P     | First published version.                  |
| v1.1.0   | 2026-03-05 | P     | Added winter maintenance procedures.      |
| v1.1.1   | 2026-03-06 | P     | Fixed typo in step 4 of procedure.        |

**Rules:**

- Every version ever created must appear in the history. No version is omitted.
- The "Summary of Changes" field must be a human-readable sentence, not a commit hash or code reference.
- Date format is always `YYYY-MM-DD` (ISO 8601).
- The version history is append-only. Previous entries are never modified (except to correct an error in the history record itself, which should be noted).

### 2.5 Branching and Forking Rules for Articles

Because this is a manual, air-gapped system with no Git-like tooling assumed, "branching" is handled through file conventions:

**Revision branching:**

When a published article (`v1.1.0`) needs revision:

1. Copy the current version to `_versions/`.
2. Change the article's status to `V` (Revised).
3. Work on the file in place, incrementing the version in the front-matter as you go (e.g., `v1.2.0-draft`).
4. When revision is complete, go through the review and publish cycle. Update the version to `v1.2.0` and status to `P`.

**Forking:**

If an article needs to split into two articles (scope grew too large):

1. The original article retains its ID and is revised with a narrowed scope.
2. A new article is created with a new ID, covering the extracted content.
3. Both articles cross-reference each other in their References section.
4. The original article's version history notes the fork: "v2.0.0 - Scope narrowed. Content on [topic] forked to [NEW-ID]."

**Merging:**

If two articles should be combined:

1. One article (the "surviving" article) absorbs the content and increments its major version.
2. The other article is archived with a note: "Merged into [SURVIVING-ID] as of [date]."
3. The surviving article's version history notes the merge.

---

## 3. ARTICLE TEMPLATE (CANONICAL)

### 3.1 The Canonical Template

Every article must follow this structure. Fields marked **(REQUIRED)** must be present in every article. Fields marked **(OPTIONAL)** may be omitted if genuinely not applicable, but their omission should be a conscious decision, not an oversight.

```markdown
---
# === FRONT-MATTER (YAML) ===
id: "{DOMAIN}-{SEQ}-ART-{NUM}-v{MAJ}.{MIN}"
title: "{Article Title}"
domain: "{DOMAIN}"
domain_name: "{Full Domain Name}"
sequence: "{SEQ}"
version: "v{MAJ}.{MIN}.{PAT}"
status: "{D|R|A|P|V|X}"
created: "YYYY-MM-DD"
last_modified: "YYYY-MM-DD"
author: "{Name or role}"
reviewers: []
tags: []
depends_on: []
depended_on_by: []
supersedes: ""
superseded_by: ""
slug: "{url-slug}"
summary: "{One-sentence summary, max 200 characters}"
---

# {Article Title}

## 1. Purpose
(REQUIRED)

{Why does this article exist? What question does it answer? What decision does
it support? Write 2-5 sentences. This is the single most important section for
a reader deciding whether to continue reading.}

## 2. Scope
(REQUIRED)

{What is covered and -- equally important -- what is NOT covered. Explicit
boundaries prevent scope creep and help readers find the right article. Use
bullet points for clarity.}

**In scope:**
- ...

**Out of scope:**
- ...

## 3. Background
(REQUIRED)

{Context a reader needs before the main content makes sense. Historical
decisions, environmental constraints, prior attempts, foundational concepts.
This section answers "what do I need to know first?" Write for a reader who is
intelligent but unfamiliar with this specific topic. Aim for 1-4 paragraphs.}

## 4. System Model
(OPTIONAL -- required for technical/procedural articles)

{A description of the system, process, or structure this article documents.
Include diagrams where possible (reference asset files). Describe components,
their relationships, inputs, outputs, and boundaries. This is the "what it is"
section before you get to "what to do with it."}

## 5. Rules & Constraints
(REQUIRED)

{The non-negotiable rules that govern this domain or topic. These are the hard
boundaries. Write them as numbered rules with clear, unambiguous language.
Each rule should be testable -- a reader should be able to determine whether
the rule is being followed or violated.}

1. ...
2. ...
3. ...

## 6. Failure Modes
(REQUIRED for procedural/technical articles, OPTIONAL for policy articles)

{What can go wrong? List known failure modes with their symptoms, likelihood
(if known), severity, and immediate response. Be honest about what you do not
know -- "unknown failure modes may exist" is a valid and valuable statement.}

| Failure Mode | Symptoms | Severity | Immediate Response |
|---|---|---|---|
| ... | ... | ... | ... |

## 7. Recovery Procedures
(REQUIRED if Failure Modes section is present)

{Step-by-step recovery instructions for each failure mode listed above.
Number each procedure. Write in imperative mood ("Do X", not "You should do
X"). Assume the reader is stressed and in a hurry -- be clear and sequential.}

### 7.1 Recovery from {Failure Mode 1}

1. ...
2. ...
3. ...

## 8. Evolution Path
(REQUIRED)

{How do you expect this system, process, or policy to change over time?
What are the known limitations of the current approach? What would trigger
a revision? This section is a gift to your future self and to future
maintainers. Write at least 2-3 bullet points.}

- ...
- ...
- ...

## 9. Commentary
(OPTIONAL -- added over time)

{Founder commentary, reflections, lessons learned, and observations that
do not fit in the main body. See Section 6 of the Meta-Framework for the
commentary workflow. Each comment is timestamped and appended, never
inserted into the middle of existing commentary.}

### Comment — {YYYY-MM-DD}

{Commentary text.}

### Comment — {YYYY-MM-DD}

{Commentary text.}

## 10. References
(REQUIRED)

{Cross-references to other articles in this institution, external sources,
and any material that informed this article. Use the Article ID for internal
references. For external sources, provide enough information to locate the
source even if the URL is dead (author, title, date, publication).}

**Internal references:**
- [{ARTICLE-ID}] {Article Title}
- ...

**External references:**
- {Author}. "{Title}." {Publication/Source}, {Date}. {URL if applicable}.
- ...

## Version History

| Version | Date | State | Summary of Changes |
|---|---|---|---|
| v0.1.0 | YYYY-MM-DD | D | Initial draft. |
```

### 3.2 Field Length Guidelines

| Section             | Guideline                                                    |
|---------------------|--------------------------------------------------------------|
| Purpose             | 2-5 sentences. 50-150 words.                                 |
| Scope               | 3-10 bullet points per sub-list. 50-200 words total.         |
| Background          | 1-4 paragraphs. 100-500 words.                               |
| System Model        | As long as needed. Prefer diagrams with brief prose. 100-1000 words. |
| Rules & Constraints | 3-20 numbered rules. Each rule: 1-3 sentences.               |
| Failure Modes       | 2-15 rows in the table. Each cell: 1-3 sentences.            |
| Recovery Procedures | 3-20 numbered steps per procedure. Each step: 1-2 sentences. |
| Evolution Path      | 2-10 bullet points. Each: 1-3 sentences.                     |
| Commentary          | No limit. Each comment entry: 1-10 paragraphs.               |
| References          | No limit. Minimum 1 internal or external reference.           |
| Summary (metadata)  | Maximum 200 characters.                                      |

These are guidelines, not hard limits. An article that needs 25 rules should have 25 rules. But an article with 200 rules probably needs to be split.

### 3.3 Required vs. Optional Fields

**Always required (every article, no exceptions):**

- All front-matter fields (ID, title, domain, version, status, created, last_modified, author, slug, summary)
- Purpose
- Scope
- Background
- Rules & Constraints
- Evolution Path
- References
- Version History

**Conditionally required:**

- System Model: Required for any article describing a physical system, software system, procedure, or process. Optional for pure policy or philosophical articles.
- Failure Modes: Required for any article with a System Model or procedures. Optional for policy articles.
- Recovery Procedures: Required whenever Failure Modes is present.

**Always optional (added over time):**

- Commentary: Added as the founder accumulates experience and reflections.

### 3.4 Formatting Standards (Markdown Conventions)

**General rules:**

1. Use ATX-style headers (`#`, `##`, `###`). Never setext-style (underlines).
2. One blank line before and after every header.
3. One blank line between paragraphs.
4. Use fenced code blocks (triple backticks) with language identifiers for all code, commands, and configuration snippets.
5. Use tables for structured data. Align pipes for readability in the source file.
6. Use numbered lists for sequential procedures or rules. Use bullet lists for unordered items.
7. Bold (`**text**`) for emphasis on key terms on first introduction. Italic (`*text*`) for titles of external works or for gentle emphasis.
8. Never use HTML in Markdown files. If Markdown cannot express it, describe it in prose or use a diagram.
9. Line length: Wrap prose at approximately 80 characters in the source file for readability in plain text editors. This is a soft guideline, not enforced mechanically.
10. Internal links use the Article ID format: `[INFR-01-ART-003-v1.0](link)`.
11. Images are referenced with descriptive alt text: `![Diagram of power distribution topology](infr-01-art-001-assets/diagram-power-topology.png)`.
12. Front-matter uses YAML between `---` fences at the top of the file.
13. No trailing whitespace.
14. Files end with a single newline character.
15. UTF-8 encoding exclusively.

---

## 4. PUBLISHING LIFECYCLE

### 4.1 State Transition Diagram

```
                  +-----------+
                  |   DRAFT   |  (v0.x.y)
                  |    (D)    |
                  +-----+-----+
                        |
                        | Author completes self-review checklist
                        v
                  +-----------+
                  |  REVIEW   |
                  |    (R)    |
                  +-----+-----+
                        |
                  +-----+-----+
                  |           |
                  v           v
            +-----------+  +-----------+
            | APPROVED  |  |   DRAFT   |  (returned with feedback)
            |    (A)    |  |    (D)    |
            +-----+-----+  +-----------+
                  |
                  | Founder publishes manually
                  v
            +-----------+
            | PUBLISHED |  (v1.0.0+)
            |    (P)    |
            +-----+-----+
                  |
            +-----+-----+
            |           |
            v           v
      +-----------+  +-----------+
      |  REVISED  |  | ARCHIVED  |
      |    (V)    |  |    (X)    |
      +-----+-----+  +-----------+
            |
            | (goes through R -> A -> P again)
            v
      +-----------+
      | PUBLISHED |  (v1.x.0+ or v2.0.0+)
      |    (P)    |
      +-----------+
```

### 4.2 State Transition Rules

| From       | To        | Condition                                                              |
|------------|-----------|------------------------------------------------------------------------|
| Draft      | Review    | Self-review checklist complete (see Section 5). All required sections present. |
| Review     | Draft     | Review reveals issues requiring significant rework.                    |
| Review     | Approved  | All quality gates passed. No unresolved issues.                        |
| Approved   | Published | Founder manually publishes. Version set to `v1.0.0` (or next major/minor). |
| Published  | Revised   | Founder begins a revision. Current published version remains live.     |
| Revised    | Review    | Revision complete, enters review again.                                |
| Published  | Archived  | Article is obsolete, superseded, or no longer relevant.                |
| Archived   | Revised   | Archived article is being brought back (rare). Must go through full review cycle. |

**Rules:**

1. No state may be skipped. An article cannot go from Draft directly to Published.
2. An article in Review cannot be modified. If changes are needed, it returns to Draft.
3. Only one version of an article may be in Published state at a time.
4. A Revised article does not replace the Published version until it completes the full Review -> Approved -> Published cycle.
5. Archived articles remain accessible on the website with a clear "ARCHIVED" banner and a link to any superseding article.

### 4.3 Who Can Transition (Founder-Only Model)

In the initial phase of the institution, **all state transitions are performed by the founder (Tim)**. There is no delegation, no committee, no voting.

**Rationale:** This is a single-person institution in its founding phase. Introducing approval workflows with multiple people would be premature complexity. The governance model can evolve later (see `GOVR` domain).

**What "founder-only" means in practice:**

- The founder writes all articles.
- The founder performs all self-reviews.
- The founder approves all articles.
- The founder publishes all articles.
- The founder decides when to archive.

**Succession provision:** The `LGCY` (Legacy) domain documents how this authority transfers. If the founder is incapacitated, the succession plan activates. See Section 10 of this document and the Legacy domain.

### 4.4 Rollback Procedures

If a published article contains a serious error discovered after publication:

**Immediate rollback (within 48 hours of publication):**

1. Revert to the previous version by copying it from `_versions/` back to the main directory.
2. Update the front-matter status to `P` and the version to the previous version number.
3. Add a version history entry: "v{X.Y.Z} - ROLLED BACK. Reason: {description}."
4. The rolled-back version is moved to `_versions/` with a filename suffix: `infr-01-art-001-v1.2.0-ROLLEDBACK.md`.
5. Add a note to the rolled-back file's front-matter: `rolled_back: true` and `rollback_reason: "{description}"`.

**Corrective revision (after 48 hours or for less severe issues):**

1. Enter the normal Revised -> Review -> Approved -> Published cycle.
2. Increment the version appropriately (minor for content fix, patch for typo).
3. Note in the version history what was corrected and why.

### 4.5 Manual Publishing Process

Publishing is a deliberate, manual act. There is no continuous deployment, no automated pipeline. This is intentional -- every publication is a conscious decision.

**Publishing checklist:**

1. Verify the article is in `Approved` state.
2. Set the version number in the front-matter to the target version.
3. Set the status to `P`.
4. Set `last_modified` to today's date.
5. Move any previous version to the `_versions/` directory.
6. Update the domain index article to include or update the article's entry.
7. Update the master index (if this is a new article).
8. Copy the file to the website's content directory.
9. Rebuild or update the website (static site generation or manual HTML update).
10. Verify the article is accessible at its URL.
11. Verify all internal links in the article resolve correctly.
12. Log the publication in the global changelog.

---

## 5. REVIEW WORKFLOW

### 5.1 Self-Review Checklist

Before any article moves from Draft to Review, the author must complete this checklist. Every item must be checked. If an item does not apply, mark it "N/A" with a reason.

**Content completeness:**

- [ ] Purpose section clearly states why the article exists.
- [ ] Scope section explicitly lists what is in and out of scope.
- [ ] Background section provides enough context for an intelligent stranger.
- [ ] All required sections are present and substantive (not placeholder text).
- [ ] Rules are numbered, unambiguous, and testable.
- [ ] Failure modes are honest (including "unknown" where applicable).
- [ ] Recovery procedures are step-by-step and written in imperative mood.
- [ ] Evolution path includes at least 2 forward-looking bullet points.
- [ ] References include at least 1 internal or external reference.

**Accuracy:**

- [ ] All factual claims are verified or marked as assumptions.
- [ ] All procedures have been tested or are clearly marked "untested."
- [ ] Numbers, measurements, and specifications are double-checked.
- [ ] Dates are correct and in ISO 8601 format.

**Consistency:**

- [ ] Article ID follows the naming convention (Section 1.1).
- [ ] File name matches the article ID (Section 1.4).
- [ ] Front-matter is complete and all fields are populated.
- [ ] Cross-references use correct Article IDs.
- [ ] Terminology matches the institution's glossary (maintained in `LRNG` domain).
- [ ] No contradictions with other published articles in the same domain.

**Formatting:**

- [ ] Markdown follows the conventions in Section 3.4.
- [ ] Tables render correctly.
- [ ] Images have alt text and are stored in the correct assets directory.
- [ ] Code blocks have language identifiers.
- [ ] No trailing whitespace or mixed line endings.

**Meta:**

- [ ] The article is necessary. It does not duplicate another article.
- [ ] The article is at the right scope. It is neither too broad nor too narrow.
- [ ] A competent person unfamiliar with this system could understand and act on this article.

### 5.2 Cross-Domain Review Triggers

An article must undergo cross-domain review (not just self-review) when:

1. **The article references another domain's systems or procedures.** The referenced domain's articles must be checked for consistency. Example: A `POWR` article about battery maintenance that references `INFR` hardware specifications.

2. **The article establishes a rule or constraint that affects another domain.** Example: A `SECR` article that mandates access controls affecting `OPSYS` system administration.

3. **The article documents a failure mode whose recovery involves another domain.** Example: A `WATR` article whose recovery procedure requires `TOOL` domain equipment.

4. **The article is in the `META`, `GOVR`, or `LGCY` domains.** These domains affect all other domains and require heightened scrutiny.

5. **The article supersedes or archives an article that other articles depend on.** All dependent articles must be checked for broken references.

**Cross-domain review procedure:**

1. Identify all affected domains using the `depends_on` front-matter field and a manual search for Article ID references.
2. Re-read the relevant articles in the affected domains.
3. Verify no contradictions, broken references, or procedural conflicts.
4. Note any required updates to other articles and create revision tasks for them.
5. Document the cross-domain review in the article's version history: "Cross-domain review completed: checked against {list of article IDs}."

### 5.3 Consistency Checks

During review, verify the following cross-cutting concerns:

**Terminology consistency:**

- The same concept is called the same thing in every article. No synonyms for technical terms.
- If a term is defined in the `LRNG` domain glossary, use that exact term.

**Measurement consistency:**

- Use metric units as primary. Imperial units may be included in parentheses if useful for practical reasons.
- Temperatures in Celsius.
- Dates in ISO 8601 (`YYYY-MM-DD`).
- Time in 24-hour format with timezone if relevant.

**Procedural consistency:**

- Procedures use the same voice (imperative) and structure (numbered steps) across all domains.
- Safety warnings use a consistent format: `**WARNING:** {text}` for danger to persons, `**CAUTION:** {text}` for danger to equipment, `**NOTE:** {text}` for important but non-dangerous information.

**Reference consistency:**

- All Article IDs referenced in the text exist and are in Published or Archived state.
- No circular "see also" loops shorter than 3 articles (A references B, B references A is acceptable only if genuinely bidirectional).

### 5.4 Quality Gates

An article cannot pass review unless all of the following gates are clear:

| Gate | Criterion | Verification Method |
|---|---|---|
| **G1: Complete** | All required sections present and substantive. | Self-review checklist. |
| **G2: Accurate** | All claims verified or flagged as assumptions. | Author attestation. |
| **G3: Consistent** | No contradictions with other published articles. | Cross-reference check. |
| **G4: Clear** | A competent stranger could understand and act on the article. | Read-aloud test: read the article aloud and note any confusion. |
| **G5: Scoped** | Article covers exactly what it says it covers, no more, no less. | Compare content against Scope section. |
| **G6: Referenced** | All sources cited. All internal references valid. | Link check. |
| **G7: Formatted** | Markdown is clean, tables render, images load. | Visual inspection. |
| **G8: Necessary** | Article does not duplicate existing coverage. | Search for overlapping articles. |

### 5.5 Approval Criteria

An article is approved when:

1. All 8 quality gates are passed.
2. The self-review checklist is complete with all items checked or marked N/A with reasons.
3. Cross-domain review is complete (if triggered).
4. The founder is satisfied that the article meets the institution's standards.
5. The article's version history documents the review.

Approval is recorded by changing the article's status to `A` and adding a version history entry: "Passed review. Ready for publication."

---

## 6. COMMENTARY WORKFLOW

### 6.1 How Founder Commentary Is Added

Commentary is the founder's ongoing dialogue with the documentation. It captures reflections, lessons learned, changed opinions, real-world observations, and contextual notes that do not belong in the article's main body.

**When to add commentary (rather than revising):**

- You have a new observation, but the article's procedures and rules are still correct.
- You want to note a real-world experience that validates or challenges the article.
- You want to flag something for future revision without revising now.
- You want to record the reasoning behind a decision that the article documents but does not explain.
- You want to preserve a historical perspective that may be valuable later.

**When to revise (rather than adding commentary):**

- The article's procedures, rules, or factual claims need to change.
- The article is missing important information that readers need.
- The article is wrong.

**The bright line:** Commentary is additive and reflective. Revision is corrective and structural. If a reader would be misled by reading the article without the commentary, the article needs a revision, not a comment.

### 6.2 Comment Threading Model

Comments are flat and chronological. There are no threaded replies, no nesting, no "reply to comment #3" structures.

**Rationale:** Threaded comments create complexity that degrades over decades. A flat, chronological log is the most durable structure. If a comment relates to a previous comment, it says so explicitly in its text.

**Comment format:**

```markdown
### Comment -- {YYYY-MM-DD}

{Commentary text. Any length. May include sub-headers (####), lists, code
blocks, and references to other articles. May reference previous comments
by date.}
```

**Rules:**

1. Comments are appended to the Commentary section in chronological order. Never inserted between existing comments.
2. Each comment includes a date header. If multiple comments are added on the same day, append a sequence letter: `2026-03-15a`, `2026-03-15b`.
3. Comments are never deleted. If a comment is factually wrong, a subsequent comment corrects it: "Correction to comment dated {date}: ..."
4. Comments are never edited after the day they are written. Corrections are made via new comments.
5. Adding a comment does not change the article's version number. Comments are metadata, not content.
6. Adding a comment updates the `last_modified` date in the front-matter.

### 6.3 Commentary vs. Revision Distinction

| Dimension          | Commentary                              | Revision                                    |
|--------------------|-----------------------------------------|---------------------------------------------|
| **Purpose**        | Reflect, observe, note                  | Correct, expand, restructure                |
| **Changes body?**  | No                                      | Yes                                         |
| **Changes version?** | No                                   | Yes (minor or major increment)              |
| **Requires review?** | No                                   | Yes (full review cycle)                     |
| **Reader impact**  | Enriches understanding                  | Changes instructions or rules               |
| **Placement**      | Commentary section only                 | Any section of the article body             |
| **Reversible?**    | No (append-only)                        | Yes (via rollback or further revision)      |

### 6.4 Historical Commentary Preservation

Commentary is a historical record. It is preserved under the following rules:

1. **When an article is revised**, all existing commentary carries over to the new version unchanged. Commentary that is no longer relevant (because the revision addressed the issue) is not removed -- it is historical context.

2. **When an article is archived**, its commentary is archived with it.

3. **When an article is forked**, commentary is copied to both resulting articles. Each article may then note in a new comment which portions of the inherited commentary are most relevant to its narrower scope.

4. **When an article is merged**, commentary from both articles is combined chronologically in the surviving article. A new comment notes the merge and the source of inherited comments.

5. **Commentary is included in the website display** of the article, in its own clearly delineated section, visually distinct from the article body.

### 6.5 Commentary-Driven Revision Triggers

Commentary may trigger a revision when:

1. **Three or more comments on the same article address the same issue.** This indicates the article's body should be updated to incorporate the accumulated knowledge.

2. **A comment explicitly states "this article needs revision."** The founder may write a comment that serves as a revision flag.

3. **A comment contradicts the article's main body.** If commentary says "I now believe step 5 is wrong," the article needs revision.

4. **External circumstances referenced in commentary invalidate the article.** Example: a comment notes that a referenced tool is no longer available, and the article's procedures depend on it.

When a revision is triggered by commentary, the revision's version history entry should note: "Revision triggered by commentary dated {date}."

---

## 7. MASTER DOMAIN INDEX

### 7.1 All 20 Domains

---

#### 1. META -- Meta-Documentation & Standards

| Attribute              | Value                                                                                     |
|------------------------|-------------------------------------------------------------------------------------------|
| **Domain Code**        | `META`                                                                                    |
| **Domain Name**        | Meta-Documentation & Standards                                                            |
| **Agent Role**         | Meta-Documentation Coordinator                                                            |
| **Core Purpose**       | Defines the standards, templates, and processes that govern all other domains.             |
| **Key Dependencies**   | `GOVR` (governance authority), `LGCY` (succession of standards), `SITE` (publishing)      |
| **Estimated Articles** | 8-12                                                                                      |

---

#### 2. INFR -- Infrastructure & Hardware

| Attribute              | Value                                                                                     |
|------------------------|-------------------------------------------------------------------------------------------|
| **Domain Code**        | `INFR`                                                                                    |
| **Domain Name**        | Infrastructure & Hardware                                                                 |
| **Agent Role**         | Infrastructure Architect                                                                  |
| **Core Purpose**       | Documents all physical and digital hardware: servers, storage devices, networking equipment, and their configurations. |
| **Key Dependencies**   | `POWR` (power supply), `NETW` (connectivity), `SECR` (physical security), `SHLT` (housing)|
| **Estimated Articles** | 20-35                                                                                     |

---

#### 3. NETW -- Networking & Communications

| Attribute              | Value                                                                                     |
|------------------------|-------------------------------------------------------------------------------------------|
| **Domain Code**        | `NETW`                                                                                    |
| **Domain Name**        | Networking & Communications                                                               |
| **Agent Role**         | Network Systems Specialist                                                                |
| **Core Purpose**       | Documents network topology, protocols, communication methods, and connectivity systems for the air-gapped and off-grid environment. |
| **Key Dependencies**   | `INFR` (hardware), `SECR` (network security), `OPSYS` (network services), `POWR` (power for networking equipment) |
| **Estimated Articles** | 12-20                                                                                     |

---

#### 4. POWR -- Power Systems & Energy

| Attribute              | Value                                                                                     |
|------------------------|-------------------------------------------------------------------------------------------|
| **Domain Code**        | `POWR`                                                                                    |
| **Domain Name**        | Power Systems & Energy                                                                    |
| **Agent Role**         | Power Systems Engineer                                                                    |
| **Core Purpose**       | Documents all energy generation, storage, distribution, and management systems including solar, battery, and backup power. |
| **Key Dependencies**   | `INFR` (hardware powered), `WTHR` (solar/wind conditions), `SHLT` (installation locations), `TOOL` (maintenance tools) |
| **Estimated Articles** | 15-25                                                                                     |

---

#### 5. STOR -- Data Storage & Backup

| Attribute              | Value                                                                                     |
|------------------------|-------------------------------------------------------------------------------------------|
| **Domain Code**        | `STOR`                                                                                    |
| **Domain Name**        | Data Storage & Backup                                                                     |
| **Agent Role**         | Data Storage Custodian                                                                    |
| **Core Purpose**       | Documents data storage architecture, backup procedures, media rotation, integrity verification, and disaster recovery for all digital assets. |
| **Key Dependencies**   | `INFR` (storage hardware), `SECR` (encryption, access), `OPSYS` (filesystem, software), `POWR` (power for storage) |
| **Estimated Articles** | 12-18                                                                                     |

---

#### 6. SECR -- Security & Access Control

| Attribute              | Value                                                                                     |
|------------------------|-------------------------------------------------------------------------------------------|
| **Domain Code**        | `SECR`                                                                                    |
| **Domain Name**        | Security & Access Control                                                                 |
| **Agent Role**         | Security Officer                                                                          |
| **Core Purpose**       | Documents all security policies, physical access controls, digital access controls, threat models, and incident response for the entire institution. |
| **Key Dependencies**   | `INFR` (secured assets), `NETW` (network security), `OPSYS` (system security), `GOVR` (security policy authority), `SHLT` (physical security) |
| **Estimated Articles** | 15-25                                                                                     |

---

#### 7. OPSYS -- Operating Systems & Software

| Attribute              | Value                                                                                     |
|------------------------|-------------------------------------------------------------------------------------------|
| **Domain Code**        | `OPSYS`                                                                                   |
| **Domain Name**        | Operating Systems & Software                                                              |
| **Agent Role**         | Systems Administrator                                                                     |
| **Core Purpose**       | Documents operating system configurations, software inventory, update procedures, and system administration practices. |
| **Key Dependencies**   | `INFR` (hardware platform), `SECR` (system hardening), `STOR` (data management), `NETW` (network services) |
| **Estimated Articles** | 18-30                                                                                     |

---

#### 8. APPL -- Applications & Services

| Attribute              | Value                                                                                     |
|------------------------|-------------------------------------------------------------------------------------------|
| **Domain Code**        | `APPL`                                                                                    |
| **Domain Name**        | Applications & Services                                                                   |
| **Agent Role**         | Applications Manager                                                                      |
| **Core Purpose**       | Documents all application-layer software, services, and tools used across the institution, including configuration, maintenance, and usage guidelines. |
| **Key Dependencies**   | `OPSYS` (platform), `SECR` (app security), `STOR` (app data), `SITE` (web applications)  |
| **Estimated Articles** | 15-25                                                                                     |

---

#### 9. SITE -- Website & Publishing (holm.chat)

| Attribute              | Value                                                                                     |
|------------------------|-------------------------------------------------------------------------------------------|
| **Domain Code**        | `SITE`                                                                                    |
| **Domain Name**        | Website & Publishing (holm.chat)                                                          |
| **Agent Role**         | Web Publisher                                                                              |
| **Core Purpose**       | Documents the design, structure, publishing process, and maintenance of the holm.chat website as the public face of the documentation institution. |
| **Key Dependencies**   | `META` (content standards), `INFR` (hosting), `SECR` (web security), `OPSYS` (web server) |
| **Estimated Articles** | 10-15                                                                                     |

---

#### 10. GOVR -- Governance & Decision-Making

| Attribute              | Value                                                                                     |
|------------------------|-------------------------------------------------------------------------------------------|
| **Domain Code**        | `GOVR`                                                                                    |
| **Domain Name**        | Governance & Decision-Making                                                              |
| **Agent Role**         | Governance Architect                                                                      |
| **Core Purpose**       | Documents the decision-making framework, authority structures, policy-making procedures, and institutional rules that govern the institution itself. |
| **Key Dependencies**   | `META` (documentation standards), `LGCY` (succession), `FINC` (resource allocation), `SECR` (access authority) |
| **Estimated Articles** | 8-15                                                                                      |

---

#### 11. FINC -- Finance & Resource Management

| Attribute              | Value                                                                                     |
|------------------------|-------------------------------------------------------------------------------------------|
| **Domain Code**        | `FINC`                                                                                    |
| **Domain Name**        | Finance & Resource Management                                                             |
| **Agent Role**         | Resource Manager                                                                          |
| **Core Purpose**       | Documents financial planning, budgeting, procurement, resource tracking, and the economics of maintaining the off-grid institution over decades. |
| **Key Dependencies**   | `GOVR` (spending authority), `INFR` (capital assets), `POWR` (energy costs), `FOOD` (food costs) |
| **Estimated Articles** | 10-18                                                                                     |

---

#### 12. LRNG -- Learning & Knowledge Management

| Attribute              | Value                                                                                     |
|------------------------|-------------------------------------------------------------------------------------------|
| **Domain Code**        | `LRNG`                                                                                    |
| **Domain Name**        | Learning & Knowledge Management                                                           |
| **Agent Role**         | Knowledge Curator                                                                         |
| **Core Purpose**       | Documents the institution's approach to learning, skill development, knowledge preservation, glossary of terms, and intellectual continuity. |
| **Key Dependencies**   | `META` (documentation standards), `LGCY` (knowledge transfer), `GOVR` (learning priorities) |
| **Estimated Articles** | 10-15                                                                                     |

---

#### 13. WTHR -- Weather & Environmental Systems

| Attribute              | Value                                                                                     |
|------------------------|-------------------------------------------------------------------------------------------|
| **Domain Code**        | `WTHR`                                                                                    |
| **Domain Name**        | Weather & Environmental Systems                                                           |
| **Agent Role**         | Environmental Monitor                                                                     |
| **Core Purpose**       | Documents weather monitoring, climate patterns, environmental data collection, and how weather affects all other systems. |
| **Key Dependencies**   | `POWR` (solar/wind generation), `FOOD` (growing conditions), `WATR` (precipitation), `SHLT` (weather protection), `INFR` (sensor hardware) |
| **Estimated Articles** | 8-12                                                                                      |

---

#### 14. FOOD -- Food Production & Preservation

| Attribute              | Value                                                                                     |
|------------------------|-------------------------------------------------------------------------------------------|
| **Domain Code**        | `FOOD`                                                                                    |
| **Domain Name**        | Food Production & Preservation                                                            |
| **Agent Role**         | Food Systems Manager                                                                      |
| **Core Purpose**       | Documents food growing, harvesting, preservation, storage, nutrition planning, and long-term food security for the off-grid institution. |
| **Key Dependencies**   | `WATR` (irrigation), `WTHR` (growing conditions), `HLTH` (nutrition), `TOOL` (farming tools), `POWR` (food preservation energy) |
| **Estimated Articles** | 15-25                                                                                     |

---

#### 15. SHLT -- Shelter & Physical Plant

| Attribute              | Value                                                                                     |
|------------------------|-------------------------------------------------------------------------------------------|
| **Domain Code**        | `SHLT`                                                                                    |
| **Domain Name**        | Shelter & Physical Plant                                                                  |
| **Agent Role**         | Facilities Manager                                                                        |
| **Core Purpose**       | Documents all buildings, structures, physical infrastructure, maintenance schedules, and construction knowledge for the institution's physical plant. |
| **Key Dependencies**   | `TOOL` (construction/repair tools), `POWR` (building power), `WATR` (plumbing), `SECR` (physical security), `WTHR` (weather exposure) |
| **Estimated Articles** | 15-25                                                                                     |

---

#### 16. WATR -- Water Systems

| Attribute              | Value                                                                                     |
|------------------------|-------------------------------------------------------------------------------------------|
| **Domain Code**        | `WATR`                                                                                    |
| **Domain Name**        | Water Systems                                                                             |
| **Agent Role**         | Water Systems Engineer                                                                    |
| **Core Purpose**       | Documents water sourcing, treatment, storage, distribution, wastewater handling, and water quality monitoring. |
| **Key Dependencies**   | `INFR` (pumps, pipes hardware), `POWR` (pump power), `HLTH` (water quality), `SHLT` (plumbing integration), `WTHR` (rainfall) |
| **Estimated Articles** | 10-18                                                                                     |

---

#### 17. TOOL -- Tools & Workshop

| Attribute              | Value                                                                                     |
|------------------------|-------------------------------------------------------------------------------------------|
| **Domain Code**        | `TOOL`                                                                                    |
| **Domain Name**        | Tools & Workshop                                                                          |
| **Agent Role**         | Workshop Manager                                                                          |
| **Core Purpose**       | Documents all tools, workshop equipment, maintenance procedures, safety protocols, and fabrication capabilities. |
| **Key Dependencies**   | `SECR` (tool access/safety), `SHLT` (workshop space), `FINC` (tool procurement), `POWR` (power tools) |
| **Estimated Articles** | 10-18                                                                                     |

---

#### 18. HLTH -- Health & Medical

| Attribute              | Value                                                                                     |
|------------------------|-------------------------------------------------------------------------------------------|
| **Domain Code**        | `HLTH`                                                                                    |
| **Domain Name**        | Health & Medical                                                                          |
| **Agent Role**         | Health Officer                                                                            |
| **Core Purpose**       | Documents health monitoring, first aid, medical supplies, emergency medical procedures, nutrition, and long-term health maintenance for an off-grid environment. |
| **Key Dependencies**   | `FOOD` (nutrition), `WATR` (water quality), `SECR` (medical records security), `COMM` (emergency communication) |
| **Estimated Articles** | 12-20                                                                                     |

---

#### 19. COMM -- Community & Social Systems

| Attribute              | Value                                                                                     |
|------------------------|-------------------------------------------------------------------------------------------|
| **Domain Code**        | `COMM`                                                                                    |
| **Domain Name**        | Community & Social Systems                                                                |
| **Agent Role**         | Community Coordinator                                                                     |
| **Core Purpose**       | Documents social structures, communication norms, conflict resolution, community events, and the human systems that sustain the institution beyond its technical infrastructure. |
| **Key Dependencies**   | `GOVR` (authority structures), `HLTH` (wellbeing), `LRNG` (education), `LGCY` (cultural continuity) |
| **Estimated Articles** | 8-15                                                                                      |

---

#### 20. LGCY -- Legacy & Succession Planning

| Attribute              | Value                                                                                     |
|------------------------|-------------------------------------------------------------------------------------------|
| **Domain Code**        | `LGCY`                                                                                    |
| **Domain Name**        | Legacy & Succession Planning                                                              |
| **Agent Role**         | Legacy Steward                                                                            |
| **Core Purpose**       | Documents succession plans, knowledge transfer procedures, institutional memory, and the long-term survival strategy for the institution beyond any single person's lifetime. |
| **Key Dependencies**   | `GOVR` (authority transfer), `META` (standards continuity), `STOR` (archive durability), `COMM` (cultural transmission), `FINC` (financial continuity) |
| **Estimated Articles** | 8-12                                                                                      |

---

### 7.2 Article Count Summary

| Domain | Code   | Estimated Articles | Cumulative |
|--------|--------|-------------------|------------|
| Meta-Documentation | `META` | 8-12 | 8-12 |
| Infrastructure | `INFR` | 20-35 | 28-47 |
| Networking | `NETW` | 12-20 | 40-67 |
| Power Systems | `POWR` | 15-25 | 55-92 |
| Data Storage | `STOR` | 12-18 | 67-110 |
| Security | `SECR` | 15-25 | 82-135 |
| Operating Systems | `OPSYS` | 18-30 | 100-165 |
| Applications | `APPL` | 15-25 | 115-190 |
| Website | `SITE` | 10-15 | 125-205 |
| Governance | `GOVR` | 8-15 | 133-220 |
| Finance | `FINC` | 10-18 | 143-238 |
| Learning | `LRNG` | 10-15 | 153-253 |
| Weather | `WTHR` | 8-12 | 161-265 |
| Food | `FOOD` | 15-25 | 176-290 |
| Shelter | `SHLT` | 15-25 | 191-315 |
| Water | `WATR` | 10-18 | 201-333 |
| Tools | `TOOL` | 10-18 | 211-351 |
| Health | `HLTH` | 12-20 | 223-371 |
| Community | `COMM` | 8-15 | 231-386 |
| Legacy | `LGCY` | 8-12 | 239-398 |

**Total estimated articles across all domains: 239-398**

---

## 8. DEPENDENCY GRAPH

### 8.1 Domain Dependencies

The following table shows which domains each domain depends on (reads from) and which domains depend on it (consumed by). A dependency means that articles in the dependent domain reference, require knowledge of, or cannot be fully understood without articles in the dependency domain.

```
                          +---------+
                          |  META   |
                          | (root)  |
                          +----+----+
                               |
              +----------------+----------------+
              |                |                |
         +----+----+     +----+----+      +----+----+
         |  GOVR   |     |  LRNG   |      |  SITE   |
         +----+----+     +----+----+      +----+----+
              |                |                |
    +---------+---------+     |           +----+----+
    |         |         |     |           |  APPL   |
+---+---+ +--+--+ +---+---+  |           +----+----+
| FINC  | | SECR| | LGCY  |  |                |
+---+---+ +--+--+ +---+---+  |           +----+----+
    |         |        |      |           | OPSYS   |
    |    +----+----+   |      |           +----+----+
    |    |         |   |      |                |
    |  +-+-+   +---+---+--+--+-------+--------+
    |  |   |   |          |          |
    | ++---++ ++---------++-+  +----+----+
    | | INFR | | NETW       |  |  STOR   |
    | +--+---+ +-----+------+  +----+----+
    |    |           |              |
    +----+-----+-----+----+--------+
         |     |          |
    +----+--+--+---+ +---+---+
    | POWR  | SHLT | | WATR  |
    +---+---+--+---+ +---+---+
        |      |          |
    +---+---+  |     +----+----+
    | WTHR  |  |     |  HLTH   |
    +-------+  |     +----+----+
               |          |
          +----+----+ +---+---+
          |  TOOL   | | FOOD  |
          +---------+ +---+---+
                          |
                     +----+----+
                     |  COMM   |
                     +---------+
```

*Note: This diagram is simplified. The full dependency matrix is below.*

### 8.2 Full Dependency Matrix

In this matrix, a mark in row X, column Y means "domain X depends on domain Y."

| Depends on -> | META | INFR | NETW | POWR | STOR | SECR | OPSYS | APPL | SITE | GOVR | FINC | LRNG | WTHR | FOOD | SHLT | WATR | TOOL | HLTH | COMM | LGCY |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| **META** | -- | | | | | | | | | x | | | | | | | | | | x |
| **INFR** | x | -- | | x | | x | | | | | | | | | x | | | | | |
| **NETW** | x | x | -- | x | | x | x | | | | | | | | | | | | | |
| **POWR** | x | x | | -- | | | | | | | | | x | | x | | x | | | |
| **STOR** | x | x | | x | -- | x | x | | | | | | | | | | | | | |
| **SECR** | x | x | x | | | -- | x | | | x | | | | | x | | | | | |
| **OPSYS** | x | x | x | | x | x | -- | | | | | | | | | | | | | |
| **APPL** | x | | | | x | x | x | -- | x | | | | | | | | | | | |
| **SITE** | x | x | | | | x | x | | -- | | | | | | | | | | | |
| **GOVR** | x | | | | | x | | | | -- | x | | | | | | | | | x |
| **FINC** | x | x | | x | | | | | | x | -- | | | x | | | | | | |
| **LRNG** | x | | | | | | | | | x | | -- | | | | | | | | x |
| **WTHR** | x | x | | x | | | | | | | | | -- | | | | | | | |
| **FOOD** | x | | | x | | | | | | | | | x | -- | | x | x | x | | |
| **SHLT** | x | | | x | | x | | | | | | | x | | -- | x | x | | | |
| **WATR** | x | x | | x | | | | | | | | | x | | x | -- | | x | | |
| **TOOL** | x | | | x | | x | | | | | x | | | | x | | -- | | | |
| **HLTH** | x | | | | | x | | | | | | | | x | | x | | -- | x | |
| **COMM** | x | | | | | | | | | x | | x | | | | | | x | -- | x |
| **LGCY** | x | | | | x | | | | | x | x | | | | | | | | x | -- |

### 8.3 Reading Order for New Maintainers

A person inheriting stewardship of this institution should read the domains in this order. The order ensures that each domain's dependencies have been covered before it is read.

**Tier 1 -- Foundational (read first, in order):**

1. `META` -- You must understand the standards before reading anything else.
2. `GOVR` -- Understand who has authority and how decisions are made.
3. `LGCY` -- Understand the succession plan (this is why you are reading).

**Tier 2 -- Core Infrastructure (read next, in any order within the tier):**

4. `SECR` -- Security context is needed before touching any system.
5. `INFR` -- Understand what hardware exists.
6. `POWR` -- Understand how it is powered.
7. `NETW` -- Understand how it is connected.

**Tier 3 -- Digital Systems (read next, in order):**

8. `STOR` -- Understand how data is stored and backed up.
9. `OPSYS` -- Understand the operating systems and software stack.
10. `APPL` -- Understand the applications running on that stack.
11. `SITE` -- Understand how the website works.

**Tier 4 -- Knowledge & Resources (read next, in any order within the tier):**

12. `LRNG` -- Understand how knowledge is managed.
13. `FINC` -- Understand the financial picture.

**Tier 5 -- Physical Sustenance (read next, in any order within the tier):**

14. `WTHR` -- Understand the environment.
15. `WATR` -- Understand the water systems.
16. `SHLT` -- Understand the physical structures.
17. `POWR` -- (re-read in context of physical systems, if needed).
18. `FOOD` -- Understand food production.
19. `TOOL` -- Understand tools and workshop capabilities.

**Tier 6 -- Human Systems (read last):**

20. `HLTH` -- Understand health and medical systems.
21. `COMM` -- Understand community and social structures.

*Note: `POWR` appears in both Tier 2 and Tier 5 because it has both digital and physical aspects. On second reading, focus on the physical/energy sections.*

### 8.4 Critical Path (What Must Be Documented First)

Not all 20 domains are equally urgent. The critical path identifies which domains must be documented first because other domains cannot be meaningfully documented without them.

**Phase 1 -- Must exist before any other domain can be documented:**

- `META` -- This document. Without naming conventions, templates, and standards, no other domain can produce conforming articles.

**Phase 2 -- Must exist before technical domains:**

- `GOVR` -- Authority structure must be established before policies are written.
- `SECR` -- Security policies must exist before systems are documented (to avoid documenting insecure configurations).

**Phase 3 -- Must exist before dependent technical domains:**

- `INFR` -- Hardware documentation is referenced by nearly everything.
- `POWR` -- Power is referenced by every physical system.

**Phase 4 -- Can proceed in parallel once Phase 3 is complete:**

- All remaining domains can be developed in parallel, though each should ensure its dependencies in Phases 1-3 are at least partially documented.

### 8.5 Circular Dependency Resolution

Several domains have mutual dependencies (A depends on B and B depends on A). These are resolved by the following principle:

**Principle:** When two domains depend on each other, the domain that is more foundational writes its articles first, using placeholders for the other domain's content. The less foundational domain then writes its articles with real references. Finally, the foundational domain updates its placeholders with real references.

**Known circular dependencies and their resolution order:**

| Domain A | Domain B | Resolution: Write First | Rationale |
|---|---|---|---|
| `INFR` | `POWR` | `INFR` | Hardware exists before it is powered. |
| `INFR` | `SHLT` | `SHLT` | Buildings exist before hardware is installed in them. |
| `SECR` | `OPSYS` | `SECR` | Security policy is defined before system hardening. |
| `SECR` | `INFR` | `SECR` | Security policy is defined before hardware is documented. |
| `POWR` | `WTHR` | `WTHR` | Weather conditions exist before solar/wind systems are designed. |
| `WATR` | `SHLT` | `SHLT` | Structures exist before plumbing is installed. |
| `FOOD` | `WATR` | `WATR` | Water systems exist before irrigation is documented. |
| `GOVR` | `META` | `META` | Standards exist before governance formalizes them. |
| `LGCY` | `GOVR` | `GOVR` | Authority structure exists before succession is planned. |

**Placeholder convention:** When an article references a dependency that has not yet been documented, use the format:

```
[PENDING: {DOMAIN}-{SEQ}-ART-{NUM} -- {expected topic}]
```

All `[PENDING: ...]` references must be resolved before the referencing article can reach Published state. A published article may not contain pending references.

---

## 9. WEBSITE STRUCTURE

### 9.1 URL Hierarchy for holm.chat

```
holm.chat/
|
+-- /                               Home page: mission statement, recent updates
+-- /docs/                           Master documentation index
|   +-- /docs/{domain-slug}/         Domain index page
|   +-- /docs/{domain-slug}/{slug}/  Individual article
|   +-- /docs/{domain-slug}/{slug}/v{x.y.z}/  Historical version
|   +-- /docs/{domain-slug}/{slug}/history/    Version history page
|   +-- /docs/{domain-slug}/{slug}/comments/   Commentary thread
|
+-- /changelog/                      Global changelog (all domains)
+-- /graph/                          Interactive dependency graph
+-- /search/                         Search interface
+-- /about/                          About the institution
+-- /about/mission/                  Mission statement (long form)
+-- /about/principles/               Core principles
+-- /about/founder/                  Founder's note
+-- /letter/                         Letter to the future maintainer (see Section 10)
+-- /status/                         System status dashboard (if applicable)
```

### 9.2 Navigation Model

The website uses a three-level navigation hierarchy:

**Level 1 -- Global navigation (always visible):**

- Home
- Documentation (-> /docs/)
- Changelog
- Graph
- Search
- About

**Level 2 -- Domain navigation (visible on /docs/ pages):**

- Sidebar or top-navigation listing all 20 domains, grouped by tier (see Section 8.3).
- Current domain is highlighted.
- Each domain shows its article count.

**Level 3 -- Article navigation (visible on individual article pages):**

- Table of contents for the current article (generated from headers).
- Previous/Next article links within the domain.
- "Related articles" section showing cross-referenced articles from `depends_on` and `depended_on_by` metadata.
- Version selector (dropdown or links to previous versions).
- Commentary link.

**Breadcrumb navigation:**

Every page displays a breadcrumb trail:

```
Home > Documentation > {Domain Name} > {Article Title}
```

### 9.3 Search Considerations

**Requirements:**

1. Full-text search across all published articles.
2. Search results must show: article title, domain, version, and a content snippet with the search term highlighted.
3. Search must support filtering by domain, status (Published/Archived), and date range.
4. Search must work without JavaScript for the basic case (form submission with server-side results).
5. Search index must be rebuildable from the Markdown source files.

**Implementation guidance (technology-agnostic):**

- If the site is static, use a pre-built search index (e.g., a JSON file generated at build time).
- If the site has a server component, use a lightweight full-text search engine.
- The search index must be storable as a flat file that can be regenerated from source.
- Do not depend on any external search service (no cloud APIs). This is an air-gapped institution.

**Search should index:**

- Article title (highest weight)
- Article summary (high weight)
- Section headers (medium weight)
- Body text (standard weight)
- Tags (high weight)
- Article ID (exact match)

### 9.4 Comment System Requirements

The "comment system" on the website is a **display system**, not an interactive comment form. Comments are added to articles by the founder through the editorial process (see Section 6) and displayed on the website.

**Requirements:**

1. Commentary is displayed on the article page in its own section, visually distinct from the article body (different background color, border, or typography).
2. Each comment shows its date prominently.
3. Comments are displayed in chronological order.
4. A dedicated `/comments/` URL for each article shows only the commentary, for readers who want to browse comments across articles.
5. There is no public comment submission form. The founder adds comments through the editorial process. If public comments are ever desired, that decision goes through the `GOVR` domain.
6. Commentary does not appear in search results by default, but can be included via a search filter toggle.

### 9.5 Revision History Display

**Requirements:**

1. Every article page includes a "Version History" link or section.
2. The version history page (`/history/`) shows the full version history table from the article's front-matter.
3. For each version, a link to the full text of that version is provided (if preserved in `_versions/`).
4. A diff view between any two versions is desirable but not required in the initial implementation. If implemented, it should be a plain text diff, not a rich diff.
5. The current version is clearly marked.
6. Archived articles show "ARCHIVED" status prominently, with a link to the superseding article if one exists.
7. The revision history includes the reason for each change (from the "Summary of Changes" field).

---

## 10. 50-YEAR CONTINUITY RULES

This institution is designed to outlast any single technology, platform, company, or person. The following rules are designed to keep the documentation accessible, understandable, and maintainable for at least 50 years (approximately 2076 and beyond).

### 10.1 Format Longevity Requirements

1. **All articles must be stored as plain text Markdown files.** Markdown has survived since 2004 and is the closest thing to a future-proof rich text format. Even if Markdown parsers cease to exist, the files are human-readable plain text.

2. **No proprietary formats.** No `.docx`, `.pages`, `.odt`, or any format requiring specific software to read. If a diagram must be stored as an image, use PNG (raster) or SVG (vector). Never use proprietary image formats.

3. **No database dependencies.** The documentation must be fully functional as a directory of text files. Any database (search index, metadata cache) must be regenerable from the source files.

4. **ASCII-safe content.** Article text should be writeable in pure ASCII wherever possible. UTF-8 is the encoding standard, but avoid Unicode characters that could be confused with ASCII characters (e.g., smart quotes, em-dashes represented as characters rather than `--`). Exception: proper names, units, and technical symbols that require Unicode.

5. **No embedded binary content.** Images, diagrams, and data files are stored as separate files, referenced by relative path. An article must be understandable even if all images are missing (alt text and prose descriptions must stand alone).

6. **Printable.** Every article must be readable when printed on paper. This means no reliance on hyperlinks for critical content (links are supplementary), no interactive elements in the content, and no content that only makes sense on a screen.

7. **Self-contained.** Every article must be understandable without access to any external system (no "see our wiki at..." or "log into the dashboard to..."). External references are for attribution and further reading, not for essential content.

8. **Multiple physical copies.** The complete documentation must exist in at least three physical locations at all times: the primary system, a local backup, and an off-site backup. At least one copy should be on a medium that does not require electricity to read (printed paper or microfiche for critical articles).

### 10.2 Technology-Agnostic Writing Guidelines

Write for a reader who may not share your technological context. In 50 years, the specific technologies you use today may be forgotten.

1. **Explain acronyms on first use, every time.** Do not assume the reader knows what "SSH," "DNS," or "USB" means. Write: "SSH (Secure Shell, a protocol for encrypted remote access to computers)."

2. **Describe physical objects by function, not brand.** Write "the photovoltaic panels on the south-facing roof" not "the Sunpower Maxeon 6 panels." Brand names go in parentheses or footnotes for specificity, but the article must be understandable without them.

3. **Include units of measurement with context.** Write "the battery bank stores 48 kilowatt-hours (kWh) of energy, enough to power the entire facility for approximately 3 days at average consumption" rather than just "48 kWh."

4. **Date everything.** Every factual claim that may become outdated should include the date it was verified. Write "as of 2026-02, the well produces approximately 20 liters per minute."

5. **Explain the "why" as much as the "what."** A future maintainer needs to know why a decision was made, not just what the decision was. Without the "why," they cannot evaluate whether the decision is still appropriate.

6. **Avoid jargon that is era-specific.** Technical jargon is acceptable (and often necessary), but slang, memes, and era-specific references should be avoided. Write "a rapid, iterative process" not "moving fast and breaking things."

7. **Assume the reader is intelligent but uninformed.** Write at the level of a literate adult who has never seen your systems. Do not write down to the reader, but do not assume shared context.

8. **Describe what you see, not what you assume.** When documenting physical systems, describe observable characteristics. "The pipe is copper, approximately 2 centimeters in outer diameter, running along the east wall at a height of approximately 2 meters" is better than "the standard plumbing."

### 10.3 Succession Documentation Requirements

The institution must be transferable to a new maintainer who has never met the founder. The following succession documentation is maintained in the `LGCY` domain:

1. **Maintainer's Orientation Guide.** A single article that a new maintainer reads first. It explains what this institution is, why it exists, how it is organized, and where to start. This is the "if you read nothing else, read this" document.

2. **Access Credentials Package.** A secure document (physical, not digital) containing all passwords, keys, and access codes. Stored separately from the documentation with its own access procedure. Referenced in `SECR` domain but stored physically.

3. **Financial Continuity Plan.** How the institution is funded, where the money is, what recurring costs exist, and how to maintain financial viability. Maintained in `FINC` domain.

4. **Skills Inventory.** What skills does the maintainer need? Where can they learn each skill? What can be hired out? Maintained in `LRNG` domain.

5. **Emergency Procedures Quick Reference.** A laminated, printed document with the 10 most critical procedures (power failure, water failure, data backup, medical emergency, etc.). Stored in a known physical location. Referenced in relevant domains but existing as a physical artifact.

6. **Annual Succession Drill.** Once per year, the founder reviews all succession documentation by imagining they must hand over the institution tomorrow. Any gaps are documented and scheduled for resolution.

7. **Relationship Map.** Who are the external contacts (neighbors, suppliers, medical providers, legal contacts) that the institution depends on? Names, roles, contact information, and the nature of the relationship.

### 10.4 Cultural Context Preservation

Documentation is written in a cultural context that will not be shared by all future readers. Preserve that context:

1. **Record the founding motivation.** Why does this institution exist? What problem does it solve? What values drive it? This goes in the `GOVR` domain and on the website's About page.

2. **Record the era.** What year is it? What is happening in the world? What is the state of technology? A brief "state of the world" article, updated annually, provides context for all other articles. This lives in the `LRNG` domain.

3. **Record the environment.** What is the physical location? What is the climate? What is the surrounding community like? This goes in the `WTHR` and `COMM` domains.

4. **Record the assumptions.** Every major decision rests on assumptions about the world. "We use solar power because grid power is unavailable" is an assumption. If grid power becomes available, the assumption changes and the decision may need revisiting. Record assumptions explicitly.

5. **Record the failures.** What was tried and did not work? Failed experiments and abandoned approaches are as valuable as successes. They prevent future maintainers from repeating mistakes. Use the Commentary section to record failures.

6. **Preserve the voice.** The founder's writing voice, opinions, and personality are part of the institution's cultural identity. Do not sanitize them out. The Commentary system exists specifically to preserve the human voice alongside the technical documentation.

### 10.5 Letter to the Future Maintainer Framework

Every domain must include, as its final article, a "Letter to the Future Maintainer." This is a personal, direct communication from the current maintainer to whoever takes over.

**Template for the Letter:**

```markdown
---
id: "{DOMAIN}-99-ART-001-v1.0"
title: "Letter to the Future Maintainer: {Domain Name}"
...
---

# Letter to the Future Maintainer: {Domain Name}

If you are reading this, you have inherited responsibility for the
{Domain Name} domain of this institution. This letter is written to
help you understand what I was thinking, what I was worried about, and
what I wish I had done differently.

## What This Domain Is Really About

{In your own words, beyond the formal scope statement, what is this
domain trying to accomplish? What is the animating spirit?}

## What I Got Right

{What are you most confident about? What decisions have been validated
by experience?}

## What I Got Wrong (or Might Have)

{What are you least confident about? What keeps you up at night? What
would you do differently if starting over?}

## What You Should Do First

{If the new maintainer can only do three things, what should they be?}

## What You Should Never Do

{Are there mistakes that would be catastrophic? Decisions that should
never be reversed without extreme caution?}

## What I Wish I Had Done

{What is on the to-do list that never got done? What would you do with
more time, money, or knowledge?}

## People and Resources

{Who should the new maintainer talk to? What books, references, or
resources were most valuable?}

## A Personal Note

{Anything else. This is your space to be human.}
```

**Rules for the Letter:**

1. The letter is written in first person.
2. The letter is honest. It includes doubts, mistakes, and uncertainties.
3. The letter is updated at least annually (or whenever a major change occurs).
4. The letter is the last article in the domain (sequence bucket `99`).
5. The letter is published -- it is not hidden or secret. Transparency about limitations is a strength.
6. The letter may include commentary like any other article.

---

## Appendix A: Quick Reference Card

**Article ID format:** `{DOMAIN}-{SEQ}-ART-{NUM}-v{MAJ}.{MIN}`
**File naming:** `{domain}-{seq}-art-{num}-v{maj}.{min}.{pat}.md` (all lowercase)
**Version format:** `v{MAJOR}.{MINOR}.{PATCH}` (semver)
**First draft version:** `v0.1.0`
**First published version:** `v1.0.0`
**Date format:** `YYYY-MM-DD` (ISO 8601)
**Encoding:** UTF-8
**States:** Draft (D) -> Review (R) -> Approved (A) -> Published (P) -> Revised (V) | Archived (X)
**Domain count:** 20
**Estimated total articles:** 239-398
**Website:** holm.chat

---

## Appendix B: Glossary of Terms Used in This Framework

| Term | Definition |
|---|---|
| **Article** | A single documentation unit. One Markdown file. The atomic unit of the institution. |
| **Domain** | One of 20 major subject areas. Each domain contains multiple articles. |
| **Sequence bucket** | A sub-grouping within a domain (the two-digit number in the Article ID). Defined by each domain. |
| **Front-matter** | YAML metadata at the top of each Markdown file, between `---` fences. |
| **Commentary** | Append-only reflections and observations by the founder. Not part of the article body. |
| **State** | The lifecycle stage of an article: Draft, Review, Approved, Published, Revised, or Archived. |
| **Rollback** | Reverting a published article to its previous version due to a serious error. |
| **Fork** | Splitting one article into two when its scope is too broad. |
| **Merge** | Combining two articles into one when they overlap too much. |
| **Pending reference** | A placeholder for an article that has not yet been written: `[PENDING: ...]`. |
| **Critical path** | The sequence of domains that must be documented first because others depend on them. |
| **Succession drill** | An annual exercise where the founder reviews all succession documentation. |
| **Air-gapped** | Not connected to the public internet. All systems are self-contained. |
| **Off-grid** | Not connected to public utilities (power, water, sewage). All resources are self-generated. |
| **Institution** | The entire documentation system and the physical/digital infrastructure it describes. |

---

## Appendix C: Change Log for This Document

| Version | Date | Summary |
|---|---|---|
| v1.0.0 | 2026-02-16 | Initial publication. Establishes all 10 sections of the meta-framework. |

---

*This document is the root of the documentation tree. If it is wrong, everything built on it inherits the error. Maintain it with care.*

*Document ID: META-00-ART-001-v1.0*
*Status: Published*
*Last modified: 2026-02-16*
