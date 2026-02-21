# STAGE 3: OPERATIONAL DOCTRINE -- INTERFACE & NAVIGATION

## Information Architecture, Content Style, Search, Accessibility, and Print/Offline Access

**Document ID:** STAGE3-INTERFACE-OPS
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Classification:** Operational Procedures -- These articles translate the D16-001 Interface Philosophy into step-by-step actionable procedures for building, maintaining, and auditing the institution's documentation interface. They are designed to be followed by a single operator with no one to ask.

---

## How to Read This Document

This document contains five operational articles that belong to Stage 3 of the holm.chat Documentation Institution, all within Domain 16 (Interface & Navigation). Stage 2 established the interface philosophy in D16-001 -- what the interface values, how it relates to sovereignty, and why clarity, accessibility, and honesty are non-negotiable. Stage 3 turns that philosophy into procedures.

These are manuals. They are written for the person who needs to organize, write, index, test, and print the institution's documentation. They assume you have read D16-001 and understood the four-layer model (philosophical, structural, interaction, presentation). They do not re-argue the principles. They implement them.

The five articles address the structural layer (D16-002: Information Architecture), the content layer (D16-003: Content Style Guide), the discovery layer (D16-004: Search and Discovery), the accessibility layer (D16-005: Accessibility Standards), and the physical output layer (D16-006: Print and Offline Access). Together they form the complete operational toolkit for the institution's documentation interface.

If something in these procedures does not work -- because the documentation platform has changed, because a tool no longer exists, because the institution's scope has grown beyond what these procedures anticipated -- do not abandon the procedure. Adapt it. The principles behind each procedure are stated in the Background section. Use those principles to find the equivalent procedure for your reality. Then update this document.

---

---

# D16-002 -- Information Architecture and Site Map

**Document ID:** D16-002
**Domain:** 16 -- Interface & Navigation
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, OPS-001, D16-001
**Depended Upon By:** D16-003, D16-004, D16-005, D16-006. All articles that create, classify, or locate content within the institution.

---

## 1. Purpose

This article defines the complete information architecture (IA) for the holm.chat Documentation Institution. It specifies how all institutional content is organized, categorized, named, and located. It provides the taxonomy, the URL structure, the breadcrumb logic, the rules for adding new sections, and the audit procedure for verifying that the architecture remains coherent over time.

Information architecture is the structural layer of the interface -- Layer 2 in the D16-001 system model. It determines whether a person looking for a specific piece of institutional knowledge can find it, and whether a person browsing without a specific goal can discover what the institution contains. A sound information architecture makes the institution navigable. A poor one makes the institution a warehouse where things are stored but nothing can be found.

This article is written for the operator who must decide where a new document belongs, what to call it, how to link it to related documents, and how to verify that the overall structure still makes sense after a hundred additions. It is also written for the future operator who inherits a structure they did not build and must understand it well enough to maintain and extend it.

## 2. Scope

This article covers:

- The institution's top-level taxonomy: what the major sections are and why they exist.
- The categorization rules: how to determine where a new piece of content belongs.
- The document identification system: how documents are numbered and named.
- The URL structure: how file paths and navigation paths are constructed.
- The breadcrumb logic: how the user's current location is communicated.
- The procedure for adding a new section or category.
- The IA audit procedure: how to verify the architecture remains sound.
- The site map: the canonical map of the institution's documentation structure.

This article does not cover:

- The content of individual documents (see D16-003 for style, individual domain articles for substance).
- The search index that makes IA discoverable (see D16-004).
- The visual presentation of navigation elements (see subsequent Domain 16 presentation articles).
- The governance process for approving new domains (see GOV-001).

## 3. Background

### 3.1 Why Information Architecture Matters More in a Closed System

In an internet-connected system, a poor information architecture is partially compensated by powerful search engines that can find content regardless of where it is stored. Users learn to bypass navigation entirely and rely on search. In an air-gapped system with a local search engine, that compensation is weaker. The local search index is only as good as its configuration and maintenance. When search fails or returns too many results, the user falls back on navigation -- on browsing the structure, following breadcrumbs, and reading the site map. If the information architecture is poor, that fallback fails too, and the user is stranded.

A sound information architecture is the institution's primary findability mechanism. Search is the secondary mechanism. Both must work. But when only one can work, the architecture must be the one that survives.

### 3.2 The Tension Between Stability and Growth

An information architecture that is perfect for fifty documents may be inadequate for five hundred. Categories that were sufficient at founding may need to be split, merged, or reorganized as the institution grows. But reorganization has a cost: it breaks bookmarks, invalidates muscle memory, confuses users who learned the old structure, and requires updating every cross-reference in every document.

The architecture must therefore be designed for growth from the beginning. It must have clear extension points -- places where new categories can be added without disrupting existing ones. It must have rules for when to split a category and when to merge. And it must have a migration procedure for when reorganization is genuinely necessary, so that the disruption is managed rather than chaotic.

### 3.3 The Domain Model as Architectural Foundation

The institution's content is organized into numbered domains established in Stage 1. This domain model is the foundation of the information architecture. It is not arbitrary. Each domain represents a coherent area of institutional concern, and the domain boundaries reflect real differences in subject matter, audience, and operational context. The information architecture builds on this foundation rather than replacing it.

## 4. System Model

### 4.1 The Top-Level Taxonomy

The institution's documentation is organized into four tiers:

**Tier 1: Stages.** The broadest organizational unit. Content is classified by its nature and maturity:
- Stage 1: Structural Framework (domain definitions, article indexes, the institutional skeleton).
- Stage 2: Philosophy (the "-001" articles for each domain; principles and values).
- Stage 3: Operational Doctrine (procedures, manuals, checklists; the "how-to" layer).
- Stage 4: Living Knowledge (ongoing records, logs, accumulated operational data).

**Tier 2: Domains.** Within each stage, content is organized by domain number (1 through 20). Each domain has a defined scope, and every document belongs to exactly one domain. A document that spans two domains is either misclassified or should be split.

**Tier 3: Documents.** Within each domain, individual documents are numbered sequentially. The "-001" document is always the domain philosophy (Stage 2). Subsequent numbers are assigned in order of creation within the domain.

**Tier 4: Sections.** Within each document, content is organized by the standard section structure (Purpose, Scope, Background, System Model, Rules & Constraints, Failure Modes, Recovery Procedures, Evolution Path, Commentary, References).

### 4.2 The Document Identification System

Every document has a unique identifier following this pattern:

`D[domain_number]-[sequence_number]`

Examples:
- `D16-001` -- Interface Philosophy (Domain 16, first document).
- `D16-002` -- This document (Domain 16, second document).
- `SEC-002` -- Access Control Procedures (Domain 3, with the SEC prefix for the Security domain).
- `ETH-001` -- Ethical Foundations (Domain 1, with the ETH prefix).

Core Charter documents (Domains 1-5) use mnemonic prefixes: ETH, CON, GOV, SEC, OPS. All other domains use the D-number format. This distinction exists because the Core Charter documents are referenced so frequently that mnemonic prefixes reduce cognitive load.

Rules for document identifiers:
- Identifiers are permanent. Once assigned, a document ID never changes, even if the document is heavily revised or partially superseded.
- Sequence numbers are never reused. If D16-007 is retired, no future document becomes D16-007.
- A document that is superseded retains its ID but receives a status change to "Superseded by D16-XXX."

### 4.3 The URL and File Path Structure

Every document maps to a predictable file path:

```
/[stage]/[domain]/[document_id].[format]
```

Examples:
- `/stage2/domain16/D16-001.md`
- `/stage3/domain16/D16-002.md`
- `/stage3/domain03/SEC-002.md`

Rules for file paths:
- Stage directories use lowercase: `stage1`, `stage2`, `stage3`, `stage4`.
- Domain directories use lowercase with zero-padded numbers: `domain01` through `domain20`.
- Document filenames use the document ID exactly as assigned.
- The canonical format is Markdown (`.md`). Rendered versions (HTML, PDF) maintain the same base path with different extensions.

Batch files (multiple articles in a single file, as used during initial authoring) are stored at the stage level with descriptive names: `/stage2/stage2-philosophy-batch4.md`. As the institution matures, individual articles should be extracted to their canonical paths for independent access. The batch files are retained as historical artifacts.

### 4.4 Breadcrumb Logic

Every page in the documentation interface displays a breadcrumb trail showing the user's location in the hierarchy:

```
Institution > Stage 3: Operational Doctrine > Domain 16: Interface & Navigation > D16-002: Information Architecture
```

Breadcrumb rules:
- The breadcrumb always begins with "Institution" as the root.
- Each level in the breadcrumb is a clickable link to that level's index page.
- The final element (the current page) is displayed but not linked.
- Breadcrumbs are generated from the file path, not hardcoded. This ensures they remain accurate as the structure evolves.
- If a page is accessed via search or a cross-reference, the breadcrumb still shows the canonical location, not the path the user took to arrive.

### 4.5 The Site Map

The institution maintains a canonical site map -- a single document that lists every document in the institution, organized by stage and domain. The site map is:
- Stored at `/stage1/sitemap.md`.
- Updated whenever a new document is added.
- Reviewed during the quarterly IA audit.
- The authoritative record of what the institution contains.

The site map lists each document with its ID, title, status (Ratified, Draft, Superseded), and date.

### 4.6 Adding a New Section or Category

When the institution needs a new document:

1. Determine the correct domain. Consult the Stage 1 domain definitions. If the content does not fit any existing domain, this is a governance decision requiring a Tier 3 decision per GOV-001 before a new domain can be created.
2. Determine the correct stage. Is this a philosophy article (Stage 2), a procedure (Stage 3), or a living record (Stage 4)?
3. Assign the next available sequence number within the domain.
4. Create the document at the canonical file path.
5. Add the document to the site map.
6. Add cross-references from related documents.
7. Rebuild the search index (see D16-004).
8. Record the addition in the operational log.

When the institution needs a new domain:

1. Draft a domain scope statement following the Stage 1 format.
2. Record the proposal in the decision log as a Tier 3 decision.
3. Verify that the proposed domain does not overlap with existing domains.
4. Assign the next available domain number.
5. Create the domain directory structure in all four stages.
6. Write the "-001" philosophy article for the new domain.
7. Update the site map.
8. Announce the new domain in the operational log.

## 5. Rules & Constraints

- **R-D16-02-01: One Document, One Domain.** Every document belongs to exactly one domain. Cross-domain topics are addressed by cross-references, not by placing a document in multiple domains.
- **R-D16-02-02: Identifiers Are Permanent.** Document identifiers are never changed, reused, or reassigned. They are the institution's permanent reference system.
- **R-D16-02-03: The Site Map Is Authoritative.** If a document exists in the file system but not in the site map, it is undocumented and must be added. If a document exists in the site map but not in the file system, the site map entry must be corrected or the document must be created.
- **R-D16-02-04: Paths Are Predictable.** A person who knows a document's ID and stage must be able to construct the file path without consulting any external reference. The path structure must be rule-based, not arbitrary.
- **R-D16-02-05: Breadcrumbs Are Mandatory.** Every page in the rendered documentation must display a breadcrumb trail. Pages without breadcrumbs fail the IA audit.
- **R-D16-02-06: New Domains Require Governance Approval.** Creating a new domain is a Tier 3 decision per GOV-001. Adding a document to an existing domain is an operational action that does not require governance approval.
- **R-D16-02-07: The Site Map Is Updated Synchronously.** When a new document is created, the site map is updated in the same work session. Deferred site map updates are prohibited because they are forgotten.
- **R-D16-02-08: Cross-References Use Document IDs.** All internal references between documents use the document ID (e.g., "see D16-004"), not file paths or page titles. Document IDs are stable; paths and titles may change.

## 6. Failure Modes

- **Orphaned documents.** A document exists in the file system but is not in the site map, not linked from any other document, and effectively invisible. The IA audit catches this by comparing the file system against the site map. Mitigation: R-D16-02-07 prevents orphans at creation time; the quarterly audit catches any that slip through.

- **Overcrowded domains.** A domain accumulates so many documents that browsing the domain index becomes unwieldy. The domain has forty documents and finding the right one requires reading all forty titles. Mitigation: when a domain exceeds twenty documents, evaluate whether it should be subdivided into sub-domains or whether a curated "start here" guide should be created for the domain.

- **Taxonomy drift.** Over time, the boundaries between domains become blurry. Documents are assigned to domains based on convenience rather than the domain scope statement. Mitigation: the IA audit includes a classification review that checks whether each document still belongs in its assigned domain according to the domain scope statement.

- **Path structure violations.** Documents are saved in ad hoc locations outside the canonical path structure, breaking the predictability rule. Mitigation: R-D16-02-04 is enforced during the IA audit. Any document found outside the canonical path is moved to the correct location and all references are updated.

- **Breadcrumb failure.** The breadcrumb generation breaks and displays incorrect hierarchy, or new pages are added without breadcrumb support. Mitigation: R-D16-02-05 makes breadcrumbs a pass/fail item in the IA audit.

## 7. Recovery Procedures

1. **If orphaned documents are discovered:** For each orphaned document, determine its correct domain and stage classification. Add it to the site map. Add cross-references from related documents. Rebuild the search index. If the document cannot be classified, it may be an artifact that should be archived rather than integrated.

2. **If the site map has become inaccurate:** Rebuild the site map from the file system. Walk every stage and domain directory. List every document found. Compare the rebuilt site map against the old one. Investigate discrepancies. The file system is the ground truth for what exists; the site map is the ground truth for what should exist. Where they disagree, resolve in favor of creating a complete and accurate site map.

3. **If a domain has become overcrowded:** Do not split the domain immediately. First, review all documents in the domain. Identify natural sub-groupings. Propose a subdivision plan. Record the proposal as a Tier 3 decision. If approved, create the sub-domains, migrate documents, update all cross-references, and update the site map. Migration must be completed in a single work session to avoid a partially migrated state.

4. **If path structure violations are widespread:** Inventory all documents in non-canonical locations. Create a migration plan. Execute the migration: move each document to its canonical path, update all cross-references to that document, update the site map, and rebuild the search index. Verify with a full IA audit afterward.

5. **If taxonomy drift has occurred:** Conduct a full classification review. For each document, compare its actual content against the domain scope statement of its assigned domain. Documents that no longer fit should be reclassified. Reclassification changes the domain directory but not the document ID. Update the site map and all cross-references.

## 8. Evolution Path

- **Years 0-5:** The initial architecture is established. The domain model is complete. The first hundred documents are created and classified. This is the period to establish the habits: update the site map with every new document, use canonical paths, maintain breadcrumbs. The architecture will feel oversized for the content. That is correct. It is designed for decades of growth.

- **Years 5-15:** The architecture should be absorbing new content without strain. The site map should be growing steadily. The first domain overcrowding issues may appear, requiring subdivision. The IA audit should be revealing small drift issues that are corrected before they compound.

- **Years 15-30:** A major structural review may be warranted. The original twenty-domain model may need extension. New stages may be needed for content types not anticipated at founding. The architecture must bend without breaking. Changes at this scale are Tier 3 governance decisions.

- **Years 30-50+:** The architecture has been maintained by multiple operators. Its value is proportional to how consistently it has been maintained. If the site map is current, the path structure is canonical, and the breadcrumbs work, the architecture is sound. If any of these have decayed, the recovery procedures above should be applied before further growth.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The temptation is to over-engineer the information architecture -- to create elaborate faceted taxonomies, tag systems, and cross-classification schemes. I have resisted that temptation. The architecture is deliberately simple: four stages, twenty domains, sequential numbering, predictable paths. Simple architectures survive because they can be understood by anyone. Elaborate architectures survive only as long as someone understands the elaboration.

The batch file approach used during initial authoring (multiple articles in a single Markdown file) is a practical concession to the founding period's rapid output pace. As the institution matures, I expect and encourage the extraction of individual articles to their canonical paths. The batch files should be retained as historical artifacts -- they document the sequence and context of the founding.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 3: Transparency of Operation)
- CON-001 -- The Founding Mandate (institutional scope and boundaries)
- GOV-001 -- Authority Model (Tier 3 decisions for structural changes)
- OPS-001 -- Operations Philosophy (complexity budget; documentation-first principle)
- D16-001 -- Interface Philosophy (Layer 2: Structural Layer; R-D16-10: The Interface Must Explain Itself)
- D16-004 -- Search and Discovery Procedures (search index dependency on IA)
- Stage 1 Documentation Framework (domain definitions and article indexes)

---

---

# D16-003 -- Content Style Guide

**Document ID:** D16-003
**Domain:** 16 -- Interface & Navigation
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, OPS-001, D16-001, D16-002
**Depended Upon By:** All articles in all domains. Every document in the institution is subject to this style guide.

---

## 1. Purpose

This article defines how all institutional content should be written. It establishes the tone, voice, terminology standards, formatting rules, plain language requirements, and accessibility writing standards that govern every document in the holm.chat Documentation Institution.

Style is not decoration. It is infrastructure. A consistent writing style across hundreds of documents and decades of authorship is what makes the institution feel like a single coherent body of knowledge rather than a disorganized collection of individual writings. When every document uses the same terminology for the same concepts, the same formatting for the same structures, and the same level of clarity in its explanations, the reader can focus on the content rather than deciphering the presentation.

This article is addressed to every person who will write, edit, or revise any institutional document. It is the standard against which all institutional writing is measured. It is not a suggestion. It is a constraint, in the same way that the document identification system in D16-002 is a constraint. Deviations from this style guide must be justified and documented, not silently adopted.

## 2. Scope

This article covers:

- Tone and voice: how institutional writing should sound.
- Terminology standards: how key concepts are named and used consistently.
- Formatting rules: how documents are structured visually.
- Plain language principles: how to write clearly for a broad audience.
- Accessibility writing standards: how to write for screen readers, cognitive accessibility, and translation.
- Examples of good and bad writing within the institutional context.
- The style review process: how writing is checked against these standards.

This article does not cover:

- The information architecture that determines where documents are stored (see D16-002).
- The visual presentation of rendered documents (typography, layout, colors -- see subsequent Domain 16 presentation articles).
- The content of specific domains (what to write is determined by domain experts; how to write it is determined here).
- Technical writing for code documentation (see Domain 8 articles for code-specific conventions).

## 3. Background

### 3.1 Why Style Matters for a Multi-Generational Institution

A document written today may be read by an operator in 2060 who has never met the author, never been trained by the author, and has no way to ask the author clarifying questions. That future reader must be able to understand the document from the document alone. This requirement eliminates many writing habits that work in conversational contexts: assumed knowledge, insider references, ambiguous pronouns, unstated context, and terminology that the author understands but has never defined.

Style consistency also matters across authors. If the founding operator writes in one style and a successor writes in another, the institution's documentation will develop a patchwork quality that undermines trust. The reader will wonder whether the inconsistency is merely stylistic or whether it reflects substantive disagreement between the documents. A style guide eliminates that ambiguity.

### 3.2 The Air-Gap Constraint on Writing

In a connected environment, a reader who encounters an unfamiliar term can search the internet for a definition. In an air-gapped institution, the only definitions available are those the institution itself provides. This means every technical term, every acronym, and every piece of institutional jargon must be defined within the institution's documentation. External references are useful context in the References section, but they cannot be relied upon for comprehension because the reader may not have access to them.

### 3.3 Plain Language Is Not Simple Language

Plain language is sometimes misunderstood as dumbed-down language. It is not. Plain language is precise language that does not require specialized training to understand. A plain-language explanation of cryptographic key management is not less accurate than a jargon-heavy one. It is more accessible without being less correct. The institution's ethical commitment to Sovereignty of the Individual (ETH-001, Principle 1) demands that institutional knowledge be accessible to its operator, not locked behind specialist vocabulary.

## 4. System Model

### 4.1 Voice and Tone

**Voice** is the consistent personality of the institution's writing. It does not change between documents or authors. The institutional voice is:

- **Direct.** Sentences state their point without hedging or burying meaning in subordinate clauses. "Back up the key to a USB drive" rather than "It is recommended that consideration be given to the backup of key material to removable media."

- **Honest.** The writing acknowledges uncertainty, difficulty, and limitation. Per D16-001 R-D16-05, if a procedure is difficult, the writing says so. If a risk is significant, the writing names it.

- **Respectful.** The writing treats the reader as an intelligent adult. Procedures explain why each step matters, not just what to do.

- **Impersonal but warm.** Third person for general statements ("The operator performs the backup"), second person for procedures ("Verify the backup completed"). First person only in the Commentary Section. The writing avoids both bureaucratic detachment and false casual friendliness.

**Tone** varies by document type while the voice remains constant:

- Philosophy articles (Stage 2): Reflective, deliberate, willing to explore complexity. Longer sentences are acceptable when exploring nuance.
- Operational articles (Stage 3): Practical, concrete, step-oriented. Shorter sentences. Numbered steps for procedures.
- Living records (Stage 4): Factual, dated, brief. Record what happened, when, and why. Do not editorialize in operational logs.

### 4.2 Terminology Standards

The institution maintains a terminology register -- a controlled vocabulary of key terms and their definitions. The terminology register is stored at `/stage1/terminology-register.md` and is updated whenever a new term is introduced.

Rules for terminology:

- **Define at first use.** Every technical term, acronym, or institutional concept is defined the first time it appears in a document. Subsequent uses within the same document may use the term without re-definition.
- **Use consistent terms.** Once a concept has an official name, use that name everywhere. Do not alternate between synonyms. The air-gapped environment is always called "the air gap" or "the air-gapped system." It is never called "the offline environment," "the disconnected system," or "the isolated network" -- because these synonyms introduce ambiguity about whether the same thing is being described.
- **Prefer common words.** Where a common English word conveys the meaning as precisely as a technical term, use the common word. "Password" is preferred over "authentication credential" in user-facing text. "Copy" is preferred over "replicate" unless the technical distinction matters.
- **Avoid jargon clusters.** A sentence with more than two technical terms should be simplified. Break "The LUKS-encrypted RAID array is backed up via rsync to the off-site NAS" into separate sentences that define each term.
- **Acronyms are spelled out at first use.** "Web Content Accessibility Guidelines (WCAG)" on first occurrence; "WCAG" thereafter.

### 4.3 Formatting Rules

All documents are authored in Markdown:

**Headings:** H1 for document title only (one per document). H2 for major sections. H3 for subsections. H4 for sub-subsections; if H5 is needed, split the document. Headings are numbered: "## 4. System Model", "### 4.1 The Account Architecture."

**Lists:** Bulleted for unordered items, numbered for sequential steps. Items begin with capitals. Complete sentences end with periods; fragments do not.

**Emphasis:** **Bold** for defined terms, rule identifiers, and warnings. *Italic* for document titles in running text and single-word emphasis. No underlining (conflicts with hyperlinks). No ALL CAPS (hostile to screen readers).

**Code:** Inline code in backticks. Multi-line blocks in triple backticks with language identifier. Commands include enough context for comprehension.

**Cross-references:** Use document IDs ("see D16-004") and rule identifiers ("per R-D16-02-01"). No page numbers or vague references ("see above").

**Paragraphs:** One idea per paragraph. Five sentences maximum before considering a split. First sentences should convey the main point.

### 4.4 Plain Language in Practice

The following plain language rules apply to all institutional writing:

1. **Active voice by default.** "The operator verifies the backup" rather than "The backup is verified by the operator." Passive voice is permitted when the actor is unknown, unimportant, or when active voice would be awkward.
2. **Short sentences for procedures.** Procedural steps should average fifteen to twenty words. Sentences over thirty words should be split.
3. **One instruction per sentence.** "Open the terminal" and "Run the backup command" are two sentences, not "Open the terminal and run the backup command."
4. **Concrete over abstract.** "Save the file to /stage3/domain16/" rather than "Persist the artifact to the appropriate location."
5. **Positive over negative.** "Keep the door locked" rather than "Do not leave the door unlocked" -- unless the negative formulation is genuinely clearer.
6. **No hidden assumptions.** If a procedure requires a prerequisite, state the prerequisite. "Before beginning: verify that you are logged in as the administrative account per SEC-002" is better than beginning a procedure that silently fails if you are in the wrong account.
7. **Explain why, not just what.** "Lock the screen when leaving the work zone (to prevent unauthorized access during your absence)" is better than "Lock the screen when leaving the work zone."

### 4.5 Accessibility Writing Standards

Writing accessibility means producing text that works for all readers, including those using assistive technology:

- **Alternative text for images.** Every image must have alt text that conveys the image's informational content. If the image is a diagram showing a four-stage pipeline, the alt text describes the pipeline, not the visual appearance. "Diagram: Four-stage import pipeline showing Reception, Quarantine, Validation, and Admission stages in sequence" rather than "Colorful flow chart."
- **Meaningful link text.** Links must describe their destination. "See the accessibility audit checklist (D16-005, Section 4.5)" rather than "Click here." Screen readers often present links out of context; "click here" is meaningless without the surrounding sentence.
- **Table headers.** Every data table must have header rows and/or header columns. Tables must not be used for layout purposes.
- **Reading level.** Aim for a Flesch-Kincaid grade level of 10-12 for procedural content and 12-14 for philosophical content. This is not a ceiling on complexity but a guide for sentence structure and word choice.
- **Consistent structure.** Screen reader users navigate by headings. The consistent section structure (Purpose, Scope, Background, etc.) allows them to jump directly to the section they need. Never skip heading levels (e.g., jumping from H2 to H4).
- **Avoid directional language.** Do not write "the button on the left" or "the section below." Screen readers and print layouts may present content differently. Use explicit references: "the Export button" or "Section 4.3."

### 4.6 Examples of Good and Bad Writing

**Example 1: Procedure Step**

Bad: "The utilization of the cryptographic key generation functionality should be undertaken subsequent to the verification of entropy availability on the system in question."

Good: "Before generating a key, verify that the system has sufficient entropy. Run `cat /proc/sys/kernel/random/entropy_avail` and confirm the value is at least 256."

Why the good version is better: It is active, concrete, specific. It tells the reader what to do, in what order, with what command, and what result to expect. The bad version is passive, abstract, and requires the reader to unpack bureaucratic language before they can act.

**Example 2: Error Communication**

Bad: "Warning: do not do this wrong or bad things will happen."

Good: "Warning: if the backup encryption key is lost and no backup copy exists, all data encrypted with that key is permanently inaccessible. There is no recovery procedure. Store backup copies per SEC-003 Section 4.3 before encrypting any data."

Why the good version is better: It specifies the exact consequence, the condition, and the prevention with a specific reference. The bad version creates anxiety without actionable information.

## 5. Rules & Constraints

- **R-D16-03-01: All Documents Follow This Guide.** Every document in the institution is subject to this style guide. No domain, stage, or document type is exempt. Deviations must be documented in the document's Commentary Section with a rationale.
- **R-D16-03-02: Terminology Is Controlled.** Terms in the terminology register must be used as defined. New terms must be added to the register before they are used in published documents.
- **R-D16-03-03: Acronyms Are Defined at First Use.** Every acronym is spelled out the first time it appears in each document.
- **R-D16-03-04: Alt Text Is Mandatory.** Every image, diagram, and non-text element must have descriptive alternative text.
- **R-D16-03-05: Cross-References Use Document IDs.** Internal references always use the document identification system, never file paths, page numbers, or vague references.
- **R-D16-03-06: The Commentary Section Is Personal.** The Commentary Section is the only place where first-person voice, personal opinion, and informal tone are permitted. It is clearly separated from the rest of the document.
- **R-D16-03-07: Formatting Is Semantic.** Bold, italic, headings, and lists are used for their semantic meaning, not for visual effect. Do not use bold for decoration. Do not use headings to make text larger.
- **R-D16-03-08: Writing Is Reviewed Before Publication.** Every document is read aloud (or processed through a text-to-speech tool) before publication to catch awkward phrasing, ambiguous references, and sentences that are too long. This review is documented in the operational log.

## 6. Failure Modes

- **Style drift.** Over years and multiple authors, the institutional voice drifts. Different operators develop different habits. Documents from different eras read differently. Mitigation: this style guide is the reference standard. The quarterly review includes a sample style check: select three documents at random and verify they conform to this guide.

- **Jargon creep.** New technical terms are introduced without definition, or casual abbreviations become standard without being added to the terminology register. Over time, the documentation becomes opaque to non-specialists. Mitigation: R-D16-03-02 and R-D16-03-03. The IA audit (D16-002) includes a terminology check.

- **Formality escalation.** Each successive author tries to sound more "professional" than the last. The writing becomes progressively more bureaucratic, passive, and impenetrable. Mitigation: the plain language rules in Section 4.4 are the antidote. When in doubt, read the sentence aloud. If it sounds like a government form, rewrite it.

- **Accessibility decay.** Alt text is omitted from new images. Link text reverts to "click here." Heading levels are skipped. Mitigation: the accessibility writing standards in Section 4.5 are checked during the quarterly review, and the accessibility audit in D16-005 catches violations.

- **Over-prescription.** This style guide becomes so detailed that it paralyzes writers. Every sentence requires checking against dozens of rules. Writing slows to a crawl. Mitigation: the rules here are principles, not a checklist for every sentence. Writers should internalize the principles and write naturally within them. The review process catches significant violations; it does not line-edit for perfection.

## 7. Recovery Procedures

1. **If style drift is detected:** Identify the dimensions of drift. Determine the correct standard by reference to this guide. Revise drifted documents. If the drift is a genuine improvement, propose a guide amendment through governance rather than silently adopting the new style.

2. **If jargon has accumulated:** Extract undefined terms, define each, add to the terminology register, and revise documents. Schedule as a dedicated task.

3. **If accessibility writing has decayed:** Conduct the audit per D16-005. Prioritize missing alt text and heading hierarchy violations. Fix systematically.

4. **If the guide itself is inadequate:** Document the problem in the Commentary Section, then revise as a Tier 2 decision per GOV-001. Apply to new documents and existing documents as they are revised. Do not mandate retroactive rewrites.

## 8. Evolution Path

- **Years 0-5:** The founding operator establishes the style by example. The first fifty documents create the practical baseline. This style guide should be revisited after the first fifty documents to verify that its rules are practical and that the examples reflect actual institutional writing.

- **Years 5-15:** A successor may join. The style guide becomes essential because it bridges two different writers. The terminology register grows. The first style drift may be detected and corrected. This is the period when the guide proves its value or reveals its gaps.

- **Years 15-30:** The guide may need revision to account for changes in language conventions, new document types, or new accessibility standards. Revision should be conservative: change the guide to accommodate genuine need, not to reflect passing fashion.

- **Years 30-50+:** The guide has been maintained by multiple authors across decades. Its greatest value is continuity -- a document from 2026 and a document from 2056 should be recognizably from the same institution. If they are not, the guide has failed and should be restored to its original principles.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
Writing about how to write is inherently recursive and somewhat absurd. This style guide is itself subject to its own rules, which means I have been writing it while simultaneously checking it against itself. The result, I hope, is a document that practices what it preaches: direct, honest, specific, and accessible.

The hardest rule to follow is the one about explaining why. It is much faster to write "do this" than "do this because of that." But the "because" is what makes the instruction survivable across decades. An instruction without a reason is brittle -- it breaks the moment the circumstances change, because the reader does not know what the instruction was protecting. An instruction with a reason is adaptable -- even if the specific action becomes obsolete, the reason endures and guides the reader toward an equivalent action.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 1: Sovereignty; Principle 3: Transparency)
- CON-001 -- The Founding Mandate (documentation completeness requirement)
- OPS-001 -- Operations Philosophy (documentation-first principle; quarterly review)
- D16-001 -- Interface Philosophy (R-D16-01: Clarity Over Cleverness; R-D16-04: Comprehensibility Across Generations; R-D16-05: Honest Communication)
- D16-002 -- Information Architecture (document identification system; cross-reference conventions)
- D16-005 -- Accessibility Standards and Testing (accessibility audit that checks writing standards)
- Web Content Accessibility Guidelines (WCAG) 2.1 -- referenced for accessibility writing principles (external; not required for comprehension of this document)

---

---

# D16-004 -- Search and Discovery Procedures

**Document ID:** D16-004
**Domain:** 16 -- Interface & Navigation
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, OPS-001, D16-001, D16-002, D16-003
**Depended Upon By:** All articles that reference findability, discoverability, or navigation within the institution.

---

## 1. Purpose

This article defines how search and discovery work within the air-gapped documentation system of the holm.chat Documentation Institution. It specifies how the search index is built, how search is configured, how relevance is tuned, and what alternative discovery paths exist for users who cannot or prefer not to use search.

In a connected environment, search is a solved problem -- powerful search engines with vast indexes and sophisticated ranking algorithms handle it. In an air-gapped environment, search is a local problem that the institution must solve for itself, using only local tools and local resources. The search engine, the index, the ranking logic, and the relevance tuning are all the institution's responsibility. There is no third party to outsource this to.

This article is written for the operator who must set up the search system from scratch, maintain the index as content grows, tune the relevance when search results are poor, and provide alternative discovery paths for users and situations where search is insufficient or unavailable.

## 2. Scope

This article covers:

- The search architecture: what software runs the search and where.
- Index building: how the search index is created and maintained.
- Search configuration: what fields are indexed, how they are weighted, and how queries are processed.
- Relevance tuning: how to improve search results when they are poor.
- Alternative discovery paths: browsing, cross-references, related articles, and curated entry points.
- Search testing: how to verify that search is working correctly.
- Degraded-mode discovery: how to find content when search is completely unavailable.

This article does not cover:

- The information architecture that organizes the content being searched (see D16-002).
- The content standards that make content searchable (see D16-003).
- The visual design of the search interface (see subsequent Domain 16 presentation articles).
- Full-text search of non-documentation content such as system logs or raw data files (see domain-specific articles).

## 3. Background

### 3.1 Search in an Air-Gapped Environment

Most search systems assume network connectivity -- either to a cloud search service or to a search server on a local network. An air-gapped documentation system has neither. The search engine must run on the same hardware that hosts the documentation, using only local resources. This constrains the choice of search technology: the engine must be lightweight enough to run alongside the documentation platform, must not require network access, and must be maintainable by a single operator.

Fortunately, the institution's documentation corpus is small by search engine standards. Even at full maturity, the corpus is likely to be thousands of documents, not millions. A lightweight local search engine -- such as those built into static site generators or standalone indexing tools -- is more than adequate. The challenge is not scale. The challenge is configuration, maintenance, and relevance tuning.

### 3.2 Search Is Not the Only Discovery Path

D16-001 establishes graceful degradation as a core interface principle (R-D16-06). If search fails, the user must still be able to find content. This means the institution must maintain robust alternative discovery paths that do not depend on search at all: browsable indexes, cross-references between related documents, curated "start here" guides, and the site map.

Search is the fastest discovery path for users who know what they are looking for. Browsing is the best discovery path for users who are exploring. Cross-references are the best discovery path for users who have found one relevant document and want to find related ones. All three must work. None is sufficient alone.

### 3.3 The Index as a Maintained Artifact

A search index is not a static artifact. It must be rebuilt whenever content is added, modified, or removed. An outdated index is worse than no index at all, because it returns results that do not match reality -- documents that have been removed still appear, new documents are missing, and modified content is indexed with outdated text. The index maintenance procedure must be as routine and non-negotiable as the site map update procedure in D16-002.

## 4. System Model

### 4.1 The Search Architecture

The search system consists of three components:

**The Indexer.** A program that reads the documentation corpus, extracts searchable text and metadata, and builds a search index. The indexer runs locally on the institution's hardware. It is invoked manually by the operator or triggered automatically when the documentation platform rebuilds.

**The Index.** A data structure (typically one or more files) that contains the indexed content in a format optimized for fast searching. The index is stored alongside the documentation. It is a derived artifact -- it can always be rebuilt from the source documents.

**The Search Interface.** The component that accepts a query from the user, looks it up in the index, and presents the results. In a static site, this is typically a JavaScript-based client-side search. In a server-based system, it may be a server-side search endpoint. In both cases, it runs entirely locally.

Recommended tools (as of founding): for static sites, Lunr.js, Pagefind, or MiniSearch (client-side, no server or network required). For server-based systems, SQLite FTS5 or Xapian (lightweight, local). The specific tool may change. The principles in this article are tool-agnostic.

### 4.2 Index Building Procedure

The search index is rebuilt whenever the documentation changes. The procedure:

1. Verify the documentation corpus is in a consistent state. All pending changes have been saved. The site map is current per D16-002.
2. Run the indexer against the documentation corpus. The indexer processes all documents in the canonical paths defined by D16-002.
3. Verify the index was built successfully. Check the indexer's output for errors or warnings.
4. Verify the index size is reasonable. A dramatic increase or decrease in index size suggests that something has gone wrong (a directory was missed, or duplicate content was indexed).
5. Test the index with three standard queries (see Section 4.5).
6. Deploy the new index to the documentation platform.
7. Record the index rebuild in the operational log, noting the number of documents indexed and any issues encountered.

Index rebuilds should occur:
- After every new document is added.
- After any document is substantially modified.
- After any document is removed or its status changes.
- At minimum, monthly, even if no changes have been made, to verify the index is still current.

### 4.3 Search Configuration

The search system indexes the following fields from each document, with the following relevance weights:

| Field | Description | Weight |
|-------|-------------|--------|
| Title | The document's H1 heading | Highest |
| Document ID | The identifier (e.g., D16-004) | Highest |
| Headings | All H2 and H3 headings | High |
| Rule IDs | Rule identifiers (e.g., R-D16-02-01) | High |
| Body text | The full text of the document | Medium |
| Metadata | Domain, stage, status, date | Low |

Weight rationale: A user searching for "access control" should see SEC-002 (which has "Access Control" in its title) before a document that merely mentions access control in passing. A user searching for "R-D16-02-01" should find the exact rule immediately.

Additional configuration:
- **Stemming:** Enabled. Searching for "backup" should also find "backups" and "backed up." The stemming algorithm should be English-language.
- **Stop words:** Common English stop words (the, a, is, at, which, on) are excluded from the index to reduce noise.
- **Minimum query length:** Two characters. Single-character queries return too many results to be useful.
- **Result limit:** Display the top twenty results. If the user needs more, the query should be refined rather than paginating through hundreds of results.

### 4.4 Relevance Tuning

Relevance tuning is the ongoing process of improving search results. It begins with the initial weights in Section 4.3 and evolves based on experience.

The relevance tuning procedure: maintain a search quality log recording queries with poor results. For each, diagnose the cause (not indexed? ranked too low? ambiguous query? too many irrelevant results?). Adjust the configuration. Rebuild and re-test. Verify improvement without degrading other queries. Record the change.

Tune incrementally. Do not attempt perfect relevance in a single pass.

### 4.5 Search Testing

The institution maintains a set of standard test queries used to verify search quality after every index rebuild and every configuration change:

**Test Query Set (initial):**

| Query | Expected Top Result | Rationale |
|-------|-------------------|-----------|
| "access control" | SEC-002 | Exact title match |
| "D16-004" | This document | Document ID search |
| "air gap" | Multiple results | Common cross-cutting term; results should be relevant |
| "backup procedure" | D6-002 or equivalent | Core operational procedure |
| "R-D16-02-01" | D16-002 | Specific rule lookup |

The test query set is expanded over time as new documents are added and new search patterns are observed. Each test query includes the expected top result and a rationale. If a test query fails (the expected result is not in the top three), this triggers the relevance tuning procedure.

### 4.6 Alternative Discovery Paths

Search is one of several discovery mechanisms. The following alternatives must be maintained:

**Browsable indexes.** Each stage and domain has an index page listing all documents with titles and descriptions, updated whenever documents change.

**Cross-references.** Every document's References section and inline references (e.g., "see D16-002") create a web of connections that allow discovery by following related documents.

**The site map.** The canonical site map (D16-002, Section 4.5) is the discovery path of last resort -- a complete listing that can be scanned linearly.

**Curated entry points.** "Start here" documents for common tasks and roles: "New Operator Start Here," "Security Checklist," "Quarterly Review Guide."

**Related articles.** Each document may include a "Related Articles" section listing documents on related topics, maintained manually and reviewed during the IA audit.

### 4.7 Degraded-Mode Discovery

If the search system is completely unavailable (the search index is corrupted, the JavaScript fails to load, the search server crashes), the user must still be able to find content. Degraded-mode discovery relies on:

1. The site map, which can be read as a plain document even without any interactive features.
2. The browsable stage and domain indexes.
3. The document identification system, which allows finding any document by constructing its path from its ID.
4. The printed site map, if print copies are maintained per D16-006.

The search system is a convenience layer on top of an architecture that is navigable without it. This is a deliberate design choice driven by R-D16-06 (Graceful Degradation).

## 5. Rules & Constraints

- **R-D16-04-01: The Search Index Is Rebuilt After Every Content Change.** No content change is complete until the search index has been rebuilt to reflect it. An outdated index is a findability failure.
- **R-D16-04-02: Search Must Work Without Network Access.** The search system must operate entirely locally. No component of the search system may require network access at any time -- not during indexing, not during querying, not during relevance tuning.
- **R-D16-04-03: Alternative Discovery Paths Must Be Maintained.** Browsable indexes, cross-references, and the site map must be maintained regardless of whether search is working. Search failure must not prevent content discovery.
- **R-D16-04-04: Search Quality Is Tested After Every Configuration Change.** The standard test query set is run after every index rebuild and every configuration change. Failures trigger the relevance tuning procedure.
- **R-D16-04-05: The Search Quality Log Is Maintained.** Queries that produce poor results are recorded in the search quality log. The log is reviewed during the quarterly review to identify patterns.
- **R-D16-04-06: Search Results Must Be Honest.** Search results must not be manipulated, boosted, or suppressed for any reason other than relevance. Per D16-001 R-D16-07 (No Hidden State), the user must be able to trust that search results reflect the actual content of the institution.

## 6. Failure Modes

- **Stale index.** The search index has not been rebuilt after content changes. New documents are invisible to search. Modified documents return outdated snippets. Users lose trust in search and abandon it. Mitigation: R-D16-04-01 makes index rebuilds mandatory after every content change. The quarterly review verifies that the index was rebuilt at least monthly.

- **Relevance decay.** The search configuration was tuned for an early corpus and has not been updated as the corpus grew. Results that were once relevant become buried under newer, less relevant content. Mitigation: the search quality log captures these problems as they are noticed. The relevance tuning procedure addresses them.

- **Search-only navigation.** Users become so dependent on search that they never use browsable indexes or cross-references. When search fails, they are lost. Mitigation: R-D16-04-03 ensures alternative paths are maintained. Curated entry points guide users to important content without requiring search.

- **Over-tuning.** The operator spends excessive time fine-tuning relevance for edge cases, degrading results for common queries in the process. Mitigation: test every tuning change against the full standard test query set, not just the query that prompted the change. If a change improves one query but degrades others, reconsider the change.

- **Index corruption.** The search index becomes corrupted due to a failed rebuild, a disk error, or a software bug. Search returns nonsensical results or fails entirely. Mitigation: the index is a derived artifact that can always be rebuilt from source documents. If corruption is suspected, delete the index and rebuild from scratch.

## 7. Recovery Procedures

1. **If the search index is stale:** Rebuild the index immediately. Run the standard test queries. Record the rebuild in the operational log with a note about how many content changes were missed.

2. **If the search index is corrupted:** Delete the existing index entirely. Rebuild from scratch by running the indexer against the full documentation corpus. Test with the standard query set. If corruption recurs, investigate the indexer and the storage system for underlying issues.

3. **If search results are poor for many queries:** Conduct a systematic relevance review. Run every query in the standard test set plus the queries in the search quality log. For each poor result, diagnose the cause (see Section 4.4). Adjust the configuration. Rebuild and re-test. This may require multiple iterations.

4. **If the search system is completely unavailable:** Fall back to degraded-mode discovery (Section 4.7). Diagnose the search system failure separately. The institution's content remains accessible through the site map, browsable indexes, and the document identification system. Search failure is an inconvenience, not a crisis, because the architecture does not depend on it.

5. **If the search tool itself must be replaced:** Evaluate replacement tools against the requirements in this article: local operation, no network dependency, adequate performance for the corpus size, configurable relevance weights. Migrate the configuration (indexed fields, weights, stop words) to the new tool. Rebuild the index. Run the full test query set. Document the tool change in the decision log.

## 8. Evolution Path

- **Years 0-5:** The search system is established. The initial index is small and relevance is easy. The search quality log begins to accumulate. The standard test query set is established and expanded. This is the time to build good habits around index rebuilding.

- **Years 5-15:** The corpus grows significantly. Relevance tuning becomes more important as more documents compete for the same queries. The search quality log should reveal patterns that guide configuration changes. The alternative discovery paths should be well-established and actively used alongside search.

- **Years 15-30:** The search tool may need replacement as the underlying platform evolves. The principles in this article should guide tool selection. The configuration (fields, weights, test queries) should transfer to the new tool. The search quality log provides historical context for relevance decisions.

- **Years 30-50+:** Search technology will have changed significantly. The specific tools and techniques in this article may be obsolete. The principles will not be: local operation, honest results, maintained indexes, alternative discovery paths, and systematic relevance tuning. Apply these principles to whatever technology is available.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
The search system is one of the areas where the air gap imposes the most visible constraint. In a connected environment, I would use a hosted search service and never think about indexes, relevance weights, or stemming algorithms. In an air-gapped environment, all of that is my responsibility.

I have deliberately designed the architecture so that search failure is survivable. The site map, the browsable indexes, the cross-references, and the predictable path structure all work without search. Search makes the institution faster to navigate. These other paths make it navigable at all. Both matter, but if I had to choose, I would choose the paths that work without software over the software that works without paths.

The standard test query set is small at founding. It will grow. Every time I search for something and the results are poor, that query goes into the search quality log and eventually into the test set. Over time, the test set becomes a portrait of how the institution is actually used, not how I imagined it would be used.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 3: Transparency of Operation)
- CON-001 -- The Founding Mandate (air-gap mandate; self-contained operation)
- OPS-001 -- Operations Philosophy (quarterly review; documentation-first principle)
- D16-001 -- Interface Philosophy (R-D16-06: Graceful Degradation; R-D16-07: No Hidden State)
- D16-002 -- Information Architecture (site map; document identification system; file path structure)
- D16-003 -- Content Style Guide (terminology standards that affect search relevance)
- D16-005 -- Accessibility Standards (search accessibility requirements)

---

---

# D16-005 -- Accessibility Standards and Testing

**Document ID:** D16-005
**Domain:** 16 -- Interface & Navigation
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, OPS-001, D16-001, D16-002, D16-003
**Depended Upon By:** All articles that produce user-facing content or interface elements. D16-006 (print accessibility). Domain 19 (quality metrics for accessibility).

---

## 1. Purpose

This article defines the concrete accessibility requirements for the holm.chat Documentation Institution. It specifies what accessibility standards must be met, how those standards are adapted for air-gapped operation, what testing procedures are used to verify compliance, how assistive technologies are supported, and what the accessibility audit checklist contains.

Accessibility is not a feature. It is not an enhancement. It is not a "nice to have." Per D16-001 R-D16-02, accessibility is non-negotiable. Every interface must be usable by persons with visual, auditory, motor, and cognitive disabilities. This article translates that principle into specific, testable requirements.

This article is written for the operator who must build accessible interfaces, test them, and maintain accessibility as the institution evolves. It is also written for the future operator who inherits an interface and must verify that its accessibility has not degraded. The audit checklist in this article is the tool for that verification.

## 2. Scope

This article covers:

- The accessibility standard adopted by the institution and its adaptation for air-gapped operation.
- Specific requirements for visual accessibility, motor accessibility, cognitive accessibility, and auditory accessibility.
- Assistive technology support requirements.
- The accessibility testing procedure: manual testing, automated testing, and assistive technology testing.
- The accessibility audit checklist: a comprehensive list of items to verify.
- The procedure for handling accessibility failures.

This article does not cover:

- The accessibility of written content (covered by D16-003, Section 4.5 -- Accessibility Writing Standards).
- The visual design of accessible interfaces (covered by subsequent Domain 16 presentation articles).
- The accessibility of print materials (covered by D16-006).
- Accessibility metrics and reporting (covered by Domain 19).

## 3. Background

### 3.1 The Adopted Standard: WCAG 2.1 AA

The institution adopts the Web Content Accessibility Guidelines (WCAG) 2.1 at the AA conformance level as its baseline accessibility standard. WCAG was chosen because it is the most widely recognized, most thoroughly documented, and most broadly applicable accessibility standard available. The AA level was chosen because it represents a reasonable balance between comprehensive accessibility and practical achievability for a single-operator institution.

WCAG is designed for web content, and this institution's primary interface is web-based (HTML rendered from Markdown, served locally). Where the institution produces non-web content (print materials, plain text files), the WCAG principles are applied in spirit rather than to the letter. The four WCAG principles -- Perceivable, Operable, Understandable, Robust -- apply to all content in all formats.

### 3.2 Adapting WCAG for Air-Gapped Operation

Several WCAG requirements assume internet connectivity: dynamic external content, third-party widgets, online evaluation tools, and cloud-based accessibility APIs. None of these are available in an air-gapped environment. These are not exemptions but substitutions -- local equivalents must be found for each. Where no local equivalent exists, the gap is documented in the accessibility audit with a plan for future resolution.

### 3.3 Why Accessibility Matters for a Single-Operator Institution

The objection arises: "I am the only user and I do not have a disability. Why does accessibility matter?" Three answers.

First, you may develop a disability. Vision degrades with age. Injuries can impair motor function. Cognitive changes are a normal part of aging. An interface built for accessibility today is an interface you can still use in thirty years.

Second, your successor may have a disability. The institution is designed to survive succession, and you do not choose your successor's abilities. An inaccessible institution cannot be inherited by a person who needs assistive technology.

Third, accessible design is better design. Clear headings, consistent navigation, sufficient contrast, keyboard operability, and plain language benefit all users, not just those with disabilities. Accessibility constraints produce interfaces that are clearer, more navigable, and more robust than unconstrained designs. D16-001 R-D16-01 (Clarity Over Cleverness) and accessibility are natural allies.

## 4. System Model

### 4.1 Visual Accessibility Requirements

**Color contrast.** All text must meet WCAG 2.1 AA contrast ratios: 4.5:1 for normal text, 3:1 for large text (18pt or 14pt bold). This applies to all text in the documentation interface, including navigation, breadcrumbs, headers, body text, and code blocks.

**Color independence.** Information must not be conveyed by color alone. If a status indicator uses green for "good" and red for "bad," it must also use text labels ("Good," "Bad") or icons (checkmark, X). A colorblind user must be able to extract the same information as a user with full color vision.

**Text resizing.** The interface must remain functional when text is resized to 200% of the default size. Content may reflow, but it must not be cut off, overlap, or become unreadable. This is tested by zooming the browser to 200%.

**Focus indicators.** When navigating with a keyboard, the currently focused element must have a visible focus indicator (outline, highlight, or similar visual cue). The default browser focus indicator is acceptable. Custom styles that remove the focus indicator are prohibited.

**Images.** All informational images must have alt text per D16-003 R-D16-03-04. Decorative images (if any) must have empty alt attributes (`alt=""`) so screen readers skip them.

**No flashing.** Content must not flash more than three times per second. This is a seizure risk. This should never be an issue in a documentation interface, but the requirement is stated for completeness.

### 4.2 Motor Accessibility Requirements

**Keyboard operability.** All functions must be operable with a keyboard alone. Every interactive element (link, button, form control, navigation menu) must be reachable via the Tab key and activatable via Enter or Space. No function may require a mouse, touch screen, or specific pointing device.

**Focus order.** The Tab key must move focus in a logical order that follows the visual layout (typically left to right, top to bottom for left-to-right languages). Focus must not jump unpredictably or get trapped in a loop.

**Skip navigation.** The interface must provide a "skip to main content" link as the first focusable element on each page. This allows keyboard users to bypass the navigation menu and go directly to the content.

**Target size.** Interactive elements (links, buttons) must have a minimum target size of 24x24 CSS pixels per WCAG 2.1 AA. This benefits users with motor impairments and touch screen users.

**No time limits.** No interface function should impose a time limit. Documentation is read at the reader's pace. If any function does require timing (which is unlikely in a documentation system), it must be adjustable or disableable.

### 4.3 Cognitive Accessibility Requirements

**Consistent navigation.** Navigation elements must appear in the same location on every page. The site structure must follow the predictable patterns defined in D16-002. Users who learn one page's layout can apply that knowledge to every other page.

**Clear language.** All text must follow the plain language standards in D16-003 Section 4.4. Jargon must be defined at first use. Sentences must be direct and unambiguous.

**Error identification.** Form errors must state what went wrong and how to fix it. No codes alone ("Error 403") without human-readable explanation.

**Consistent identification.** The same concept uses the same name everywhere, per D16-003 terminology standards.

**Predictable behavior.** Links go somewhere. Buttons do something. Nothing happens on focus alone. Context changes are user-initiated.

### 4.4 Auditory Accessibility Requirements

The institution's documentation is primarily text-based, so auditory accessibility requirements are minimal. However:

- If any audio or video content is added to the institution, it must have text alternatives: transcripts for audio, captions for video.
- System sounds (if any) must not be the sole means of conveying information. Every sound must have a visual equivalent.

### 4.5 Assistive Technology Support

The interface must be compatible with: **screen readers** (semantically correct HTML, ARIA landmarks, meaningful link text; test with Orca on Linux or VoiceOver on macOS); **keyboard-only navigation** (per Section 4.2); **screen magnification** (functional at 200% zoom); and **high-contrast modes** (functional with OS high-contrast activated).

Install at least one screen reader during initial setup. Verify it functions with the documentation interface. Include it in the software inventory per D8-001. Update it with system software updates.

### 4.6 The Accessibility Testing Procedure

Accessibility testing uses three methods, applied in order:

**Method 1: Automated testing.**
Run a local accessibility tool (Pa11y, Axe-core, or Lighthouse) against the rendered HTML. Automated testing catches 30-40% of issues: missing alt text, contrast failures, missing form labels, heading violations. It cannot catch issues requiring judgment.

Procedure: generate the rendered HTML; run the tool against every page or a sample of at least twenty; record results; distinguish true failures from false positives; fix all true failures; re-run to verify.

**Method 2: Manual keyboard testing.**
Navigate the entire interface using only a keyboard. Verify: every link is reachable via Tab and activatable via Enter; focus order is logical; focus indicators are visible; the skip navigation link works; no focus traps exist. Test from the home page through at least five different document pages. Record any failures.

**Method 3: Screen reader testing.**
Navigate using a screen reader (Orca on Linux, VoiceOver on macOS). Verify: page titles are read correctly; headings are announced at correct levels; links have meaningful text; images are announced with alt text; breadcrumbs are navigable; reading order matches visual order. Test at least three document pages. Record any failures.

### 4.7 The Accessibility Audit Checklist

The following checklist is used during the quarterly accessibility audit and after any significant interface change:

**Perceivable:**
- [ ] All images have meaningful alt text.
- [ ] Color contrast meets 4.5:1 for normal text, 3:1 for large text.
- [ ] Information is not conveyed by color alone.
- [ ] Text can be resized to 200% without loss of content.
- [ ] Audio/video content (if any) has text alternatives.

**Operable:**
- [ ] All functions are keyboard-operable.
- [ ] Focus order is logical.
- [ ] Focus indicators are visible.
- [ ] Skip navigation link is present and functional.
- [ ] No focus traps exist.
- [ ] Interactive element target size is at least 24x24 pixels.
- [ ] No time limits are imposed.

**Understandable:**
- [ ] Navigation is consistent across pages.
- [ ] Language is clear per D16-003 standards.
- [ ] Error messages (if applicable) identify the error and suggest correction.
- [ ] Terminology is consistent per D16-003 standards.
- [ ] Interface behavior is predictable.

**Robust:**
- [ ] HTML validates without errors that affect accessibility.
- [ ] ARIA roles and labels are used correctly (not excessively).
- [ ] The interface works with at least one screen reader.
- [ ] The interface works in high-contrast mode.
- [ ] The interface works at 200% zoom.

**Documentation-specific:**
- [ ] Heading hierarchy is correct (no skipped levels).
- [ ] Breadcrumbs are present and accurate.
- [ ] Cross-references are meaningful (not "click here").
- [ ] Tables have headers.
- [ ] Code blocks are accessible to screen readers.

## 5. Rules & Constraints

- **R-D16-05-01: WCAG 2.1 AA Is the Baseline.** All institutional interfaces must meet WCAG 2.1 Level AA conformance. This is a minimum, not a target. Exceeding AA where practical is encouraged.
- **R-D16-05-02: Accessibility Is Tested, Not Assumed.** Accessibility compliance must be verified through the testing procedure in Section 4.6. A claim of accessibility without testing is not accepted.
- **R-D16-05-03: Accessibility Failures Block Deployment.** A new or modified interface that fails any item in the audit checklist must not be deployed until the failure is resolved. Accessibility failures are not "known issues" to be fixed later.
- **R-D16-05-04: A Screen Reader Must Be Available.** At least one screen reader must be installed and functional on the institution's primary workstation. It is part of the institution's required software.
- **R-D16-05-05: Accessibility Is Tested Quarterly.** The full accessibility audit checklist is run quarterly as part of the institutional review cycle. The results are recorded in the operational log.
- **R-D16-05-06: Accessibility Failures Are Tier 1 Bugs.** An accessibility failure that prevents a category of users from accessing content is treated with the same urgency as a data integrity failure. It is not a cosmetic issue.

## 6. Failure Modes

- **Accessibility erosion.** Each small interface change introduces a minor accessibility issue. Over time, the accumulated issues make the interface unusable for assistive technology users. Mitigation: R-D16-05-03 prevents individual changes from introducing issues. R-D16-05-05 catches any that slip through.

- **Testing tool dependency.** Relying solely on automated testing misses 60% of issues. Mitigation: all three testing methods must be used.

- **Assistive technology unavailability.** Screen reader not installed or operator cannot use it. Mitigation: R-D16-05-04 requires installation; operator learns basic navigation during training.

- **False compliance.** Checklist completed superficially. Mitigation: audit results must include specific evidence (e.g., "minimum contrast ratio found: 5.2:1"), not just checkmarks.

- **Standard obsolescence.** WCAG 2.1 AA is superseded. Mitigation: evaluate new standards for adoption during the annual review.

## 7. Recovery Procedures

1. **If the quarterly audit reveals multiple failures:** Prioritize failures by severity. Failures that completely prevent access (no keyboard operability, no alt text on critical images) are fixed first. Failures that degrade the experience (poor contrast, inconsistent focus order) are fixed second. Schedule the fixes as a dedicated accessibility remediation task in the operational plan.

2. **If a screen reader is not available:** Acquire and install one at the earliest opportunity. In the interim, conduct manual keyboard testing (Method 2) with extra diligence, and pay special attention to HTML semantic correctness (which can be verified without a screen reader using automated tools). Document the gap in the audit results.

3. **If an accessibility failure is discovered in production:** Assess its severity. If it completely blocks a category of users, fix it immediately -- this is equivalent to a system outage for that category. If it degrades the experience, fix it within the current operational cycle. Record the failure and its resolution in the operational log.

4. **If automated testing tools become unavailable:** Manual testing (keyboard and screen reader) remains available. Supplement with manual contrast checking (using locally installed color picker tools) and manual HTML validation. Automated tools are a convenience that accelerates testing; they are not required for testing to occur.

5. **If the adopted standard is updated:** Review the new standard. Identify new criteria. Assess the institution's compliance against the new criteria. Plan a remediation cycle for any gaps. Update this document to reference the new standard. This is a Tier 2 decision per GOV-001.

## 8. Evolution Path

- **Years 0-5:** The accessibility baseline is established. The first interface is built with accessibility in mind from the start, per D16-001. The screen reader is installed and the operator learns basic screen reader navigation. The quarterly audit habit is formed. This is the cheapest time to build accessibility -- retrofitting is always more expensive.

- **Years 5-15:** Accessibility standards may be updated (WCAG 2.2, WCAG 3.0, or successors). The institution evaluates and adopts updated standards. New assistive technologies may emerge that require interface adjustments. The testing procedure may need new methods.

- **Years 15-30:** A successor may have different accessibility needs. The accessibility architecture must accommodate needs that the founding operator did not have. The quarterly audit ensures that the interface has not degraded. If the successor has a disability, the interface should already support them because of the standards maintained throughout.

- **Years 30-50+:** Accessibility technology will have evolved significantly. The principles (perceivable, operable, understandable, robust) are timeless. The specific requirements will change. Future operators should adopt the accessibility standard of their era, using the principles in this document and in D16-001 as the unchanging foundation.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
I do not currently use assistive technology. This means I am building accessibility requirements based on standards, documentation, and empathy rather than lived experience. I am aware of the limitations of this approach. The standards are thorough, but they cannot capture every nuance of assistive technology use. The testing procedures are rigorous, but they are performed by someone who uses these technologies for testing, not for daily life.

This is why the quarterly audit is so important. It is not a formality. It is the mechanism that catches the issues I cannot see because I do not experience them. If a future operator uses a screen reader daily, their experience will reveal gaps that no amount of testing by a sighted operator could find. I encourage that future operator to update this article with what they learn.

I have also made the deliberate choice to install a screen reader and learn to use it, even though I do not need it personally. Using assistive technology, even briefly and clumsily, builds understanding that reading about it does not. I recommend this practice for every future operator.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 1: Sovereignty; Principle 2: Integrity Over Convenience)
- CON-001 -- The Founding Mandate (institutional durability; successor readiness)
- OPS-001 -- Operations Philosophy (quarterly review cycle; operational discipline)
- D16-001 -- Interface Philosophy (R-D16-02: Accessibility Is Non-Negotiable; R-D16-04: Comprehensibility Across Generations; R-D16-09: Input Method Neutrality)
- D16-002 -- Information Architecture (heading hierarchy; breadcrumb requirements)
- D16-003 -- Content Style Guide (Section 4.5: Accessibility Writing Standards; alt text requirements)
- D16-006 -- Print and Offline Access Procedures (print accessibility)
- WCAG 2.1 -- Web Content Accessibility Guidelines, W3C Recommendation (external; the full text should be imported into the institution per D18-001 procedures)
- Domain 19 -- Quality Assurance (accessibility quality metrics)

---

---

# D16-006 -- Print and Offline Access Procedures

**Document ID:** D16-006
**Domain:** 16 -- Interface & Navigation
**Version:** 1.0.0
**Date:** 2026-02-16
**Status:** Ratified
**Depends On:** ETH-001, CON-001, OPS-001, D16-001, D16-002, D16-003, D16-005
**Depended Upon By:** D6-002 (backup procedures for print materials). Domain 19 (quality metrics for print outputs). Domain 20 (institutional memory in physical form).

---

## 1. Purpose

This article defines the procedures for producing printed versions of institutional documentation, maintaining print stylesheets, assembling critical document print packages, creating laminated emergency references, providing offline mobile access, and producing the annual print archive.

D16-001 R-D16-08 states: "All critical content must be renderable on paper. The interface must support or enable a print path that produces readable, structured output. In a fifty-year institution, there will be times when the only functioning interface is ink on paper." This article implements that requirement.

Print and offline access are not secondary concerns. They are survival mechanisms. When the hardware fails, when the power is out, when the digital system is being rebuilt or migrated, the institution's knowledge must still be accessible. Paper does not crash. Paper does not require a boot sequence. Paper does not suffer from bit rot. Paper is the institution's most durable interface, and it must be maintained as deliberately as any digital system.

This article is written for the operator who must produce print materials, maintain the print stylesheet, decide what gets printed and how often, and store printed materials so they remain usable for decades.

## 2. Scope

This article covers:

- The print stylesheet: requirements, design principles, and maintenance.
- The critical documents print package: what is printed, how, and how often.
- Laminated emergency references: which documents are laminated and why.
- Offline mobile access: serving documentation on a local device without the primary system.
- The annual print archive: the yearly full print of the institutional documentation.
- Storage and preservation of printed materials.

This article does not cover:

- The digital rendering of documentation (see subsequent Domain 16 presentation articles).
- Backup procedures for digital content (see D6-002).
- The content of specific documents (see domain-specific articles).
- Paper and supply procurement (see operational logistics articles).

## 3. Background

### 3.1 Why Paper Still Matters

The institution is digital-first. Its primary interface is a screen. Its primary storage is a disk. But "digital-first" does not mean "digital-only." Digital systems fail. Hardware breaks. Software corrupts. Power goes out. In any of these scenarios, the institution's knowledge exists only in a form that requires a functioning digital system to access.

Paper is the institution's backup interface. A printed document can be read with no electricity, no hardware, no software, and no technical skill beyond literacy. A printed document does not require a boot sequence. It does not suffer from incompatible file formats. It does not become inaccessible because a software dependency is no longer available. It is technology-independent in a way that no digital format can match.

This does not mean everything should be printed. It means certain critical documents must be printed, and the institution must maintain the capability to print any document when needed.

### 3.2 The Spectrum of Offline Access

The institution maintains four levels of access, each a fallback for the levels above it: (1) the primary digital system, (2) an offline mobile device with a static documentation copy, (3) printed documents requiring only light and literacy, and (4) laminated emergency references surviving conditions that destroy ordinary paper.

### 3.3 Print Quality and Readability

Web pages printed without a print stylesheet produce broken layouts, cut-off text, and navigation waste. A proper print stylesheet transforms screen layout into paper layout. Print readability is an accessibility concern -- if a printed document is illegible, the institution has failed the person relying on it.

## 4. System Model

### 4.1 The Print Stylesheet

The institution maintains a CSS print stylesheet that is applied whenever a document is printed from the browser or converted to PDF. The print stylesheet ensures:

**Typography:** Body text in serif font, 11-12pt, black on white. Headings bold and proportionally larger (H1: 18pt, H2: 14pt, H3: 12pt). Code blocks in monospace, 9-10pt. Line height 1.4-1.6 for body text.

**Page layout:** A4 default, US Letter secondary. Margins at least 20mm. Headers show document title and ID on every page. Footers show page number and total ("Page 3 of 17"). Page breaks avoid splitting paragraphs, table rows, or separating headings from their content.

**Content adaptation:** Navigation elements hidden. Hyperlink URLs displayed in parentheses or as footnotes (internal document ID references need no URL expansion). Decorative images hidden; informational images included with alt text as captions. Tables repeat headers if split across pages.

**Maintenance:** Test the stylesheet after every platform update by printing at least three documents of different types. The stylesheet is version-controlled in the platform configuration.

### 4.2 The Critical Documents Print Package

The following documents are printed and maintained as a physical package. This package is the minimum documentation needed to operate and recover the institution if the digital system is unavailable.

**Contents of the Critical Documents Print Package:**

1. **ETH-001 -- Ethical Foundations.** The institution's values. Without these, nothing else makes sense.
2. **CON-001 -- The Founding Mandate.** What the institution is, what it does, and why.
3. **GOV-001 -- Authority Model.** Who decides what, especially the succession protocol.
4. **SEC-001 -- Threat Model and Security Philosophy.** What the institution protects against.
5. **SEC-002 -- Access Control Procedures.** How to get into the system.
6. **SEC-003 -- Cryptographic Key Management.** How to handle the keys (the keys themselves are stored separately per SEC-003; this is the procedure document).
7. **OPS-001 -- Operations Philosophy.** How the institution operates.
8. **D6-002 -- Backup Doctrine.** How to restore the institution from backup.
9. **D10-002 -- Daily Operations Doctrine.** The daily operational procedures.
10. **D10-003 -- Incident Response Procedures.** What to do when something goes wrong.
11. **The site map** (from D16-002). A complete listing of all institutional documents.
12. **The key inventory** (from SEC-003, metadata only -- no actual keys).
13. **The hardware inventory** (from D8-001).
14. **The software inventory** (from D8-001).

**Print schedule:** The critical documents print package is reprinted:
- Annually, during the annual review cycle.
- Whenever any document in the package is substantively revised.
- Whenever a succession event occurs, so the new operator has a current physical copy.

**Storage:** The print package is stored in a binder or folder in Zone 3 (the Storage Zone per SEC-002). A second copy is stored at the off-site backup location if one exists. The binder is labeled with the print date and the version of each document included.

### 4.3 Laminated Emergency References

Certain documents are printed on laminated stock for durability. Laminated documents survive water, dirt, handling, and the conditions that would destroy ordinary paper. They are the institution's last-resort documentation.

**Documents that are laminated:**

1. **Emergency Recovery Card.** A single-page reference (front and back) containing:
   - The institution's name and physical address.
   - The location of the primary hardware.
   - The location of backup media.
   - The boot sequence for the primary system.
   - The location of the succession packet.
   - The key inventory summary (which keys exist and where they are stored -- no actual keys).
   - Contact information for the designated successor.
   - A brief list of first-response procedures for common emergencies (power failure, hardware failure, suspected compromise).

2. **Succession Summary Card.** A single-page reference containing:
   - The succession protocol summary from GOV-001.
   - The location of the succession packet.
   - The accounts available to the successor and how to activate them.
   - The location of the critical documents print package.
   - The first ten steps a new operator should take.

3. **Daily Operations Quick Reference.** A single-page reference containing:
   - The daily checklist from D10-002.
   - The weekly checklist from D10-002.
   - The monthly checklist from D10-002.
   - The location of the full D10-002 document for detailed procedures.

**Lamination procedure:** Print at 100% scale, verify content is current, laminate with heavy-weight pouches (at least 5 mil / 125 micron), label with date. Store one copy with the critical documents print package and one copy adjacent to the primary hardware.

**Reprint schedule:** Annually or whenever content changes.

### 4.4 Offline Mobile Access

The institution can produce a portable copy of its documentation for use on a mobile device (tablet, laptop, e-reader) that is separate from the primary system.

**Production procedure:** Generate a static HTML version of the entire documentation site including the search index. Copy to a USB drive or mobile device. The static site must function from the local file system with relative links, local resources, and client-side search. Verify loading and search on the target device. Label with the copy date.

**Update schedule:** Quarterly or when significant content changes occur.

**Use cases:** Consulting documentation during primary system maintenance; carrying documentation to off-site locations; providing a successor with documentation before primary system access.

**Security note:** The offline copy must not contain cryptographic keys, passwords, or credential material. It is a documentation copy, not a system backup.

### 4.5 The Annual Print Archive

Once per year, during the annual review cycle, the institution produces a complete print of all documentation. This is not a selective print like the critical documents package. It is an exhaustive print of every document in the institution.

**Purpose:** The annual print archive serves as a point-in-time snapshot of the institution's entire knowledge base. If the digital system and all digital backups were lost, the most recent annual print archive would be the starting point for reconstruction.

**Procedure:** Generate print-ready versions of every document, organized by stage and domain. Print using the stylesheet (Section 4.1). Assemble in a binder with divider tabs, a cover page (institution name, "Annual Print Archive," date, document count), a table of contents, and a printed site map. Store in Zone 3 per SEC-002. Retain all previous years' archives as historical records.

**Practical considerations:** If a full print is impractical, print at minimum the critical documents package plus documents created or revised since the last full print, and record the partial print in the operational log. Use double-sided printing and archival-quality (acid-free) paper when available.

### 4.6 Storage and Preservation of Printed Materials

Store printed materials at room temperature, low-to-moderate humidity (use desiccant packets in humid locations), minimal light exposure, and in binders or closed cabinets for physical protection. Label every container with contents and date. Record the inventory of printed materials in the operational log.

## 5. Rules & Constraints

- **R-D16-06-01: The Print Stylesheet Is Maintained.** The print stylesheet must be tested after every platform update and every significant template change. A broken print stylesheet is a print access failure.
- **R-D16-06-02: The Critical Documents Print Package Is Current.** The physical print package must reflect the current versions of all included documents. An outdated print package is worse than none at all because it may be relied upon and mislead.
- **R-D16-06-03: Laminated References Are Updated Annually.** Laminated emergency references must be reprinted annually or whenever their content changes, whichever comes first.
- **R-D16-06-04: The Annual Print Archive Is Produced.** A complete or substantial print archive is produced annually during the annual review cycle. If a complete print is impractical, a partial print plus the critical documents package is the minimum.
- **R-D16-06-05: Offline Mobile Copies Exclude Credentials.** No offline mobile copy may contain cryptographic keys, passwords, or other credential material. Documentation content only.
- **R-D16-06-06: Previous Print Archives Are Retained.** Old print archives are not discarded. They are historical records of the institution's knowledge at a point in time.
- **R-D16-06-07: Print Materials Are Stored Properly.** Printed materials must be stored in conditions that preserve their readability per Section 4.6. Materials stored in damaging conditions must be relocated.

## 6. Failure Modes

- **Print stylesheet neglect.** Platform updates break the stylesheet silently; printing later produces unusable output. Mitigation: R-D16-06-01 requires testing after every update.

- **Critical documents package staleness.** Revised documents make the printed package inaccurate. Mitigation: R-D16-06-02 requires currency; the annual review checks print dates against revision dates.

- **Laminated card omission.** Cards reference obsolete hardware or personnel. Mitigation: R-D16-06-03 mandates annual reprinting.

- **Print archive impracticality.** Corpus grows too large for full annual printing. Mitigation: R-D16-06-04 allows partial prints as a fallback.

- **Paper degradation.** Poor storage makes materials unreadable within a decade. Mitigation: R-D16-06-07 requires proper storage; the annual review includes physical inspection.

- **Offline copy with credentials.** Key material inadvertently included on a lost device. Mitigation: R-D16-06-05 prohibits credentials; the annual review verifies compliance.

## 7. Recovery Procedures

1. **If the print stylesheet is broken:** Identify what changed (usually a platform update). Update the stylesheet. Test with three document types. Use "export to PDF" as a temporary workaround if needed.

2. **If the critical documents print package is outdated:** Reprint revised documents and replace outdated pages. Update the cover page date. Reprint the entire package if more than half the documents need updating.

3. **If laminated references are outdated:** Reprint with current content, laminate, replace old cards, and destroy old versions to prevent confusion.

4. **If annual print archives have been missed:** Produce a current-year archive. Do not retroactively print missed years. Note the gap in the operational log.

5. **If printed materials have degraded:** Relocate readable materials to better storage. Reprint unreadable ones from digital sources, or consult D6-002 if digital sources are also lost.

6. **If credentials are in an offline copy:** Retrieve the device, delete the copy, assess whether the device was outside physical control (if so, rotate credentials per SEC-003), remove credentials from the source, and rebuild.

## 8. Evolution Path

- **Years 0-5:** The print infrastructure is established. The print stylesheet is written and tested. The first critical documents package and laminated cards are produced. The first annual print archive is created. The habit of maintaining print alongside digital is formed. This is the hardest period because the print effort feels redundant when the digital system is working perfectly.

- **Years 5-15:** The annual print archive becomes routine. The critical documents package is updated naturally as documents are revised. The print stylesheet may need updating as the documentation platform evolves. The stored archives begin to have historical value -- they show how the institution's knowledge evolved over time.

- **Years 15-30:** Print materials may be needed for the first time in a real emergency or succession event. Their value becomes concrete rather than theoretical. The storage conditions should be inspected and, if necessary, improved. The accumulated print archives should be reviewed for physical condition.

- **Years 30-50+:** The oldest print archives are now decades old. Archival-quality paper should still be readable. Standard paper may have degraded. This is the period that validates the paper quality and storage decisions made at founding. If the materials have survived, the institution has proven its physical durability alongside its digital durability.

## 9. Commentary Section

*This section is reserved for dated entries by current and future operators.*

**2026-02-16 -- Founding Entry:**
I expect the print procedures to feel like the most old-fashioned part of this entire framework. We live in a digital world. Printing documentation feels like keeping a horse in case the car breaks down. But the analogy is not quite right. A car has many failure modes. A horse has many failure modes. Paper has almost none. It does not require maintenance. It does not require power. It does not require software updates. It sits on a shelf, unchanging, until the moment you need it.

The critical documents print package is the part of this article I feel most strongly about. If the digital system fails catastrophically -- disk failure, cryptographic key loss, fire, flood -- the print package is what allows the next person to understand what this institution was, what it contained, and how to rebuild it. That is worth a binder on a shelf and an hour of printing once a year.

The laminated emergency recovery card is inspired by the laminated checklists used in aviation. Pilots do not rely on memory during emergencies. They rely on checklists that have been printed, laminated, and placed within arm's reach. An institution should extend the same courtesy to its operator.

## 10. References

- ETH-001 -- Ethical Foundations (Principle 4: Longevity Over Novelty)
- CON-001 -- The Founding Mandate (institutional durability; self-contained operation)
- OPS-001 -- Operations Philosophy (annual review cycle; operational discipline)
- D16-001 -- Interface Philosophy (R-D16-08: Printability; R-D16-06: Graceful Degradation)
- D16-002 -- Information Architecture (site map; document identification system)
- D16-003 -- Content Style Guide (formatting standards that affect print output)
- D16-005 -- Accessibility Standards (print accessibility considerations)
- SEC-002 -- Access Control Procedures (Zone 3: Storage Zone)
- SEC-003 -- Cryptographic Key Management (key inventory; credential exclusion from print)
- D6-002 -- Backup Doctrine (relationship between digital backup and print archive)
- D8-001 -- Platform Doctrine (hardware and software inventories)
- D10-002 -- Daily Operations Doctrine (operational checklists)
- D10-003 -- Incident Response Procedures (emergency procedures referenced in laminated cards)
- GOV-001 -- Authority Model (succession protocol referenced in succession summary card)

---
