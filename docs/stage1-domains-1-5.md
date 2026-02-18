# STAGE 1: DOCUMENTATION FRAMEWORK
## Domains 1 through 5 -- Complete Output
### A Lifelong, Air-Gapped, Off-Grid, Self-Built Digital Institution

**Document ID:** STAGE1-FRAMEWORK-001
**Version:** 1.0.0
**Date of Origin:** 2026-02-16
**Status:** Initial Draft -- Awaiting Review
**Intended Lifespan:** 50+ years (perpetual with maintenance)
**Audience:** Founding authors, future maintainers, successors, auditors

---

## Preamble

This document defines the documentation framework for a self-sovereign digital institution designed to operate indefinitely without dependence on external networks, cloud services, commercial vendors, or institutional continuity of any single person. Every article referenced herein will follow the mandated template structure. Every domain is written with the assumption that the original authors will eventually be unavailable, that hardware will be replaced many times over, and that the cultural context surrounding this institution will shift in ways we cannot predict.

Nothing in this framework is code. Nothing is tooling. This is the map that precedes the territory.

### Article Template (Canonical Reference)

Every article produced under this framework MUST contain the following sections, in this order:

| Section | Purpose |
|---|---|
| **Title** | Unique, descriptive, searchable name |
| **Purpose** | Why this article exists; what problem it addresses |
| **Scope** | What is covered and, critically, what is NOT covered |
| **Background** | Historical context, motivation, prior art, reasoning |
| **System Model** | How the relevant system works -- diagrams, flows, relationships |
| **Rules & Constraints** | Hard rules, invariants, non-negotiable boundaries |
| **Failure Modes** | What can go wrong; enumerated and classified by severity |
| **Recovery Procedures** | Step-by-step recovery for each identified failure mode |
| **Evolution Path** | How this article and its subject are expected to change over time |
| **Commentary Section** | Space for future maintainers to annotate, disagree, or extend |
| **References** | Cross-links to other articles, external sources (archived), standards |

### Meta-Rules Governing All Documentation

1. **Fifty-Year Horizon.** Every article must be written as though the reader has never met the author and lives in a substantially different technological and cultural context.
2. **Loss of Original Author.** No article may depend on oral tradition, tacit knowledge, or the continued availability of any specific person.
3. **Hardware Impermanence.** All references to specific hardware must include abstraction layers and migration paths. No article may assume any particular device will exist in ten years.
4. **Cultural Drift.** Terminology must be defined inline. Jargon must be explained. Assumptions about values, norms, and priorities must be made explicit rather than implicit.
5. **Self-Containment.** This institution is air-gapped. Every external reference must be archived locally. No article may depend on a URL, API, or service that is not under institutional control.
6. **Contradiction Resolution.** Where two articles conflict, the article closer to Domain 1 (Constitution) takes precedence unless a formal override is documented.
7. **Plain Language.** Prefer clarity over elegance. Prefer redundancy over ambiguity.

---
---

# DOMAIN 1: CONSTITUTION & PHILOSOPHY

**Agent Role:** Philosophy Architect
**Core Purpose:** Define the foundational beliefs, principles, and purpose of the institution. Why it exists. What it stands for. What it refuses.

---

## 1.1 Domain Map

### Scope

Domain 1 is the root of all other domains. It answers the questions that precede every technical and operational decision: Why does this institution exist? What principles are non-negotiable? What is the institution willing to sacrifice, and what will it never sacrifice? What does it mean for this institution to succeed, and what does it mean for it to fail in spirit even if it survives in form?

This domain covers:

- The institution's reason for existence (its mandate)
- Core philosophical commitments (sovereignty, privacy, durability, autonomy)
- Ethical boundaries and refusals (what the institution will never do)
- The relationship between the institution and the individuals who maintain it
- The definition of institutional identity over time (the Ship of Theseus problem)
- The theory of knowledge preservation that underpins the entire documentation system
- The institution's stance on external engagement, cooperation, and isolation
- Principles governing how the institution relates to future technologies, societies, and crises

### Boundaries

Domain 1 does NOT cover:

- Operational governance procedures (Domain 2)
- Specific security implementations (Domain 3)
- Hardware or power specifications (Domain 4)
- Software architecture decisions (Domain 5)

Domain 1 provides the WHY. Other domains provide the HOW and the WHAT. When a HOW or WHAT conflicts with a WHY from Domain 1, Domain 1 prevails.

### Relationships

- **Domain 2 (Governance):** Domain 1 defines the principles that Domain 2 must operationalize. Governance structures are instruments of the philosophy, not the reverse.
- **Domain 3 (Security):** Domain 1 defines what "integrity" and "trust" mean at a philosophical level. Domain 3 implements those definitions technically.
- **Domain 4 (Infrastructure):** Domain 1 defines the commitment to physical sovereignty and off-grid operation. Domain 4 makes it real.
- **Domain 5 (Platform):** Domain 1 defines the preference for simplicity, longevity, and independence that constrains all platform decisions.

---

## 1.2 Article List

| ID | Title | Summary |
|---|---|---|
| C-001 | **The Founding Mandate** | Why this institution was created. The original problem it was built to solve. The conditions under which it would no longer need to exist. |
| C-002 | **Core Principles and Non-Negotiable Values** | Enumeration and deep explanation of every foundational principle: sovereignty, privacy, durability, autonomy, simplicity, transparency to members, opacity to outsiders. |
| C-003 | **The Refusal Registry** | An explicit, maintained list of things the institution will never do, with reasoning. Examples: never connect to the public internet, never depend on a commercial vendor for a critical function, never store data on behalf of an external party. |
| C-004 | **Institutional Identity Over Time** | How the institution defines itself when all original members are gone, all hardware is replaced, and all software is rewritten. The Ship of Theseus problem, addressed directly. |
| C-005 | **The Theory of Knowledge Preservation** | Why documentation exists. The epistemological framework: what counts as institutional knowledge, how it is recorded, how it is validated, how it decays, and how it is renewed. |
| C-006 | **The Ethics of Isolation** | The moral and practical case for air-gapped operation. Acknowledgment of what is lost. Conditions under which limited external contact is permissible. |
| C-007 | **Relationship Between Institution and Individual** | The rights and obligations of members. What the institution owes its maintainers. What maintainers owe the institution. How personal autonomy coexists with institutional discipline. |
| C-008 | **Succession Philosophy** | The principles governing how leadership, knowledge, and authority transfer between generations. Not the mechanics (Domain 2) but the philosophy: what must be preserved in a successor, and what may be allowed to change. |
| C-009 | **The Acceptable Loss Doctrine** | What the institution is willing to lose in order to survive. Prioritization framework: which assets, capabilities, or even principles can be sacrificed under existential threat, and which cannot. |
| C-010 | **Technology Philosophy** | The institution's relationship with technology. Technology as tool, not identity. The commitment to understanding every layer of the stack. The rejection of black boxes. |
| C-011 | **The Amendment Process for Foundational Principles** | How the constitution itself can be changed. The threshold for amendment. The safeguards against erosion. The distinction between amendment and betrayal. |
| C-012 | **Glossary of Foundational Terms** | Precise definitions of every term used across Domain 1, written for readers 50 years from now who may not share current assumptions about language. |
| C-013 | **Failure of Spirit: Recognizing Institutional Drift** | How to detect when the institution has technically survived but philosophically failed. Warning signs. Diagnostic questions future maintainers should ask themselves. |

---

## 1.3 Priority Order

**Tier 1 -- Write Immediately (existential foundation):**

1. C-001: The Founding Mandate
2. C-002: Core Principles and Non-Negotiable Values
3. C-003: The Refusal Registry
4. C-012: Glossary of Foundational Terms

**Tier 2 -- Write Early (required before other domains can finalize):**

5. C-005: The Theory of Knowledge Preservation
6. C-010: Technology Philosophy
7. C-007: Relationship Between Institution and Individual
8. C-004: Institutional Identity Over Time

**Tier 3 -- Write Before Operations Begin:**

9. C-006: The Ethics of Isolation
10. C-008: Succession Philosophy
11. C-009: The Acceptable Loss Doctrine
12. C-011: The Amendment Process for Foundational Principles

**Tier 4 -- Write During Early Operations:**

13. C-013: Failure of Spirit: Recognizing Institutional Drift

Rationale: Tier 1 articles are prerequisites for every other document in every other domain. Tier 2 articles are required by Domain 2 and Domain 3 before they can finalize their own work. Tier 3 articles are needed before the institution begins operating but can be drafted in parallel with other domains. Tier 4 is reflective and benefits from operational experience.

---

## 1.4 Dependencies

| Article | Depends On | Depended On By |
|---|---|---|
| C-001 | None (root document) | All articles in all domains |
| C-002 | C-001 | G-001, G-002, S-001, I-001, P-001 (all domain root articles) |
| C-003 | C-001, C-002 | S-002, S-003, I-003, P-002 |
| C-004 | C-001, C-002 | G-005, G-006 |
| C-005 | C-001, C-002 | All documentation-related articles across all domains |
| C-006 | C-002, C-003 | S-001, S-004, I-005 |
| C-007 | C-001, C-002 | G-001, G-003, G-004 |
| C-008 | C-002, C-007 | G-005, G-006 |
| C-009 | C-002, C-003 | S-008, I-009, P-010 |
| C-010 | C-001, C-002 | I-001, P-001, P-003 |
| C-011 | C-001, C-002, C-007 | G-002 |
| C-012 | C-001 | All articles in all domains |
| C-013 | C-001, C-002, C-004 | G-008 |

---

## 1.5 Writing Schedule

| Week | Articles | Notes |
|---|---|---|
| Week 1 | C-001, C-012 | Mandate and glossary first. Glossary will expand continuously. |
| Week 2 | C-002 | Requires deep deliberation. May require multiple drafts. |
| Week 3 | C-003 | Best written immediately after C-002, while the boundaries are fresh. |
| Week 4 | C-005, C-010 | These two can be drafted in parallel by the same author. |
| Week 5 | C-007 | Requires careful thought about obligations and rights. |
| Week 6 | C-004 | Philosophical. Benefits from having C-002 and C-007 complete. |
| Week 7 | C-006, C-008 | Can be drafted in parallel. |
| Week 8 | C-009, C-011 | Can be drafted in parallel. |
| Week 10 | C-013 | Write after initial operational experience. Revisit at month 6, year 1, year 5. |

---

## 1.6 Review Process

1. **First Draft:** Written by the Philosophy Architect (or designated founding author).
2. **Self-Review:** Author revisits after 48 hours minimum. Checks for implicit assumptions, undefined terms, and emotional rather than rational justifications.
3. **Cross-Domain Review:** Each Tier 1 and Tier 2 article must be reviewed by at least the Governance Engineer and one other domain lead, to verify it is implementable and not self-contradictory.
4. **Adversarial Review:** A designated reviewer attempts to find loopholes, contradictions, or scenarios where the principles produce absurd or harmful outcomes.
5. **Future-Reader Test:** Give the article to someone unfamiliar with the project. Can they understand it without asking questions? If not, revise.
6. **Ratification:** Tier 1 articles require unanimous agreement among all founding members. Tier 2 and below require majority agreement plus no unresolved objections.
7. **Version Lock:** Once ratified, the article is versioned (1.0.0) and any subsequent changes follow the Amendment Process defined in C-011.

---

## 1.7 Maintenance Plan

- **Annual Review:** Every article in Domain 1 is re-read by the current Philosophy Architect (or equivalent) once per year. The review is documented in the Commentary Section with a date stamp and the reviewer's assessment: "Still valid," "Needs minor update," "Needs major revision," or "Potentially obsolete."
- **Trigger-Based Review:** Any of the following triggers an immediate review of all Domain 1 articles:
  - Change in the number of institutional members (addition or loss)
  - A governance crisis (as defined in Domain 2)
  - A security breach (as defined in Domain 3)
  - A major hardware transition (as defined in Domain 4)
  - Discovery of internal contradiction
- **Commentary Accumulation:** Future maintainers are encouraged to add commentary even when not performing a formal review. The Commentary Section is append-only and never edited (only annotated).
- **Glossary Maintenance:** C-012 is a living document. New terms may be added at any time. Existing definitions may only be changed through the formal review process.
- **Succession of Authorship:** If the Philosophy Architect role is vacated, the maintenance plan itself must be reviewed by the incoming author within 30 days of assuming the role.

---
---

# DOMAIN 2: GOVERNANCE & AUTHORITY

**Agent Role:** Governance Engineer
**Core Purpose:** Define who decides what, how authority flows, how disputes resolve, succession planning, decision frameworks.

---

## 2.1 Domain Map

### Scope

Domain 2 translates the principles of Domain 1 into operational decision-making structures. It answers: Who has authority to make which decisions? How is that authority granted, constrained, and revoked? How do disputes between individuals or between domains get resolved? What happens when the decision-making structure itself fails?

This domain covers:

- The authority model (roles, responsibilities, limits)
- Decision-making frameworks (consensus, delegation, emergency override)
- Dispute resolution procedures
- Succession planning and knowledge transfer
- Role definitions and their boundaries
- Meeting and deliberation structures
- Record-keeping requirements for decisions
- The relationship between individual judgment and institutional rules
- Emergency governance (when normal processes are unavailable)
- Amendment and evolution of governance structures themselves

### Boundaries

Domain 2 does NOT cover:

- The philosophical basis for governance (Domain 1, specifically C-007 and C-008)
- Technical access control implementation (Domain 3)
- Physical access to infrastructure (Domain 4)
- Software permission models (Domain 5)

Domain 2 defines the human decision-making layer. Domains 3, 4, and 5 implement technical enforcement of those decisions where appropriate.

### Relationships

- **Domain 1 (Constitution):** Domain 2 is subordinate to Domain 1. Every governance structure must be traceable to a constitutional principle. Where governance rules conflict with constitutional principles, the constitution prevails.
- **Domain 3 (Security):** Domain 2 defines who has authority; Domain 3 defines how that authority is technically verified and enforced. Changes to access control policy originate in Domain 2 and are implemented by Domain 3.
- **Domain 4 (Infrastructure):** Domain 2 defines who may authorize physical changes to infrastructure. Domain 4 defines the technical procedures.
- **Domain 5 (Platform):** Domain 2 defines the approval process for platform changes. Domain 5 defines the technical execution.

---

## 2.2 Article List

| ID | Title | Summary |
|---|---|---|
| G-001 | **The Authority Model** | Complete definition of all roles, their authorities, their limits, and the principle of least authority as applied to human governance. |
| G-002 | **Decision-Making Frameworks** | The taxonomy of decisions (routine, significant, critical, existential) and the process required for each category. Includes consensus requirements, quorum rules, and time limits. |
| G-003 | **Role Definitions and Boundaries** | Detailed specification of every named role in the institution. What each role may do, must do, and must not do. How roles interact. |
| G-004 | **Dispute Resolution Procedures** | Step-by-step processes for resolving disagreements between members, between roles, and between domains. Escalation paths. Final arbitration. |
| G-005 | **Succession Planning** | How each role is transferred. Knowledge transfer requirements. Minimum overlap periods. What happens in unplanned succession (death, incapacitation, departure). |
| G-006 | **The Continuity Protocol** | What happens when the institution is reduced to a single member. Or zero members who then return. The minimum viable governance structure. |
| G-007 | **Record-Keeping Requirements** | What decisions must be recorded, in what format, with what level of detail. Retention periods. Access rules for decision records. |
| G-008 | **Institutional Health Assessment** | Periodic self-assessment framework. How to measure whether governance is functioning. Metrics, warning signs, diagnostic questions. |
| G-009 | **Emergency Governance** | Decision-making under crisis conditions. Who has expanded authority. What constraints remain even in emergencies. How normal governance is restored. |
| G-010 | **The Veto and Override Framework** | When and how decisions can be blocked or overridden. The difference between a veto (blocking) and an override (forcing). Safeguards against abuse of both. |
| G-011 | **External Relations Policy** | Governance of any interaction with entities outside the institution. Who may authorize contact. What information may be shared. Under what circumstances. |
| G-012 | **Governance Amendment Process** | How governance structures themselves change. The threshold for change. Protection against both ossification and instability. |
| G-013 | **Meeting and Deliberation Protocols** | How formal deliberations are conducted, recorded, and archived. Quorum, agenda, documentation standards. |
| G-014 | **Accountability and Transparency Framework** | How role-holders are held accountable. What information is visible to all members. What is restricted and why. |

---

## 2.3 Priority Order

**Tier 1 -- Write Immediately (required for operations):**

1. G-001: The Authority Model
2. G-003: Role Definitions and Boundaries
3. G-002: Decision-Making Frameworks
4. G-009: Emergency Governance

**Tier 2 -- Write Early (required before full operations):**

5. G-004: Dispute Resolution Procedures
6. G-007: Record-Keeping Requirements
7. G-005: Succession Planning
8. G-010: The Veto and Override Framework

**Tier 3 -- Write Before Year One:**

9. G-006: The Continuity Protocol
10. G-013: Meeting and Deliberation Protocols
11. G-014: Accountability and Transparency Framework
12. G-011: External Relations Policy

**Tier 4 -- Write During Early Operations:**

13. G-008: Institutional Health Assessment
14. G-012: Governance Amendment Process

Rationale: The institution cannot function without knowing who has authority (G-001), what roles exist (G-003), how decisions are made (G-002), and what to do in a crisis (G-009). These four are existential prerequisites. Dispute resolution and succession are needed before the first serious disagreement or departure.

---

## 2.4 Dependencies

| Article | Depends On | Depended On By |
|---|---|---|
| G-001 | C-001, C-002, C-007 | G-002, G-003, G-004, G-009, G-010, S-001, S-005 |
| G-002 | C-002, G-001 | G-004, G-009, G-010, G-012, G-013 |
| G-003 | C-007, G-001 | G-004, G-005, G-007, S-005, I-001, P-001 |
| G-004 | G-001, G-002, G-003 | G-010 |
| G-005 | C-008, G-001, G-003 | G-006 |
| G-006 | C-009, G-005 | S-008 |
| G-007 | C-005, G-002 | G-008, G-014 |
| G-008 | C-013, G-007 | G-012 |
| G-009 | C-009, G-001, G-002 | S-008, I-009 |
| G-010 | G-001, G-002, G-004 | G-012 |
| G-011 | C-003, C-006, G-001 | S-004 |
| G-012 | C-011, G-002, G-008, G-010 | None (terminal) |
| G-013 | G-002, G-007 | None |
| G-014 | G-003, G-007 | G-008 |

---

## 2.5 Writing Schedule

| Week | Articles | Notes |
|---|---|---|
| Week 2 | G-001 | Can begin once C-001 is drafted. Will iterate alongside C-002. |
| Week 3 | G-003 | Role definitions. Requires G-001 as input. |
| Week 4 | G-002, G-009 | Decision frameworks and emergency governance can be drafted in parallel. |
| Week 5 | G-004, G-007 | Dispute resolution and record-keeping. |
| Week 6 | G-005, G-010 | Succession planning and veto framework. |
| Week 7 | G-013, G-014 | Meeting protocols and accountability. |
| Week 8 | G-006, G-011 | Continuity protocol and external relations. |
| Week 12 | G-008 | After some operational experience is available. |
| Week 14 | G-012 | Last to write; requires understanding of all other governance structures. |

---

## 2.6 Review Process

1. **First Draft:** Written by the Governance Engineer.
2. **Constitutional Compliance Check:** The Philosophy Architect reviews every governance article to confirm alignment with Domain 1 principles. This is a mandatory gate -- no governance article is ratified without this sign-off.
3. **Implementability Review:** The Security Strategist reviews to confirm that governance decisions can be technically enforced where needed.
4. **Scenario Testing:** Each article is tested against at least three scenarios:
   - A routine situation where the process works as designed.
   - A stressed situation where the process is strained (e.g., two members disagree strongly).
   - A catastrophic situation where the process must handle extreme conditions (e.g., loss of a key member, simultaneous crises).
5. **Plain Language Review:** A non-author reads the article. Every sentence that requires clarification is flagged for rewrite.
6. **Ratification:** Tier 1 articles require agreement from all founding members. Tier 2 and below require majority plus the Philosophy Architect's constitutional compliance sign-off.

---

## 2.7 Maintenance Plan

- **Semi-Annual Review:** Every governance article is reviewed every six months during the first five years, then annually thereafter. Governance structures are more likely to need adjustment than constitutional principles.
- **Post-Incident Review:** After any dispute resolution (G-004), emergency governance activation (G-009), or succession event (G-005), the relevant articles are reviewed within 30 days. Lessons learned are recorded in the Commentary Section.
- **Role Transition Review:** Whenever a role changes hands, the incoming role-holder must read all Domain 2 articles within their first 14 days and document any questions or concerns in the Commentary Sections.
- **Governance Health Check:** G-008 (Institutional Health Assessment) defines a periodic self-evaluation. The results of each evaluation may trigger reviews of other governance articles.
- **Sunset Clause Review:** Every five years, each governance article is evaluated for whether it is still needed, or whether its purpose has been absorbed by another article or made obsolete by changed circumstances.

---
---

# DOMAIN 3: SECURITY & INTEGRITY

**Agent Role:** Security Strategist
**Core Purpose:** Define threat models, access control philosophy, cryptographic principles, trust boundaries, audit frameworks.

---

## 3.1 Domain Map

### Scope

Domain 3 defines how the institution protects itself -- its data, its processes, its physical assets, and its people -- from threats both external and internal. It establishes the principles and frameworks for: identifying what must be protected, understanding what threatens it, defining how protection is achieved, verifying that protection is maintained, and recovering when protection fails.

This domain covers:

- Threat modeling methodology and the institutional threat model
- Access control philosophy and frameworks
- Cryptographic principles, key management, and long-term cryptographic planning
- Trust boundaries (between people, between systems, between the institution and the outside)
- Audit and verification frameworks
- Incident detection, response, and post-incident analysis
- Data integrity and provenance
- Physical security principles (in coordination with Domain 4)
- Insider threat considerations
- Long-term security evolution (post-quantum, algorithm aging, key rotation over decades)

### Boundaries

Domain 3 does NOT cover:

- The philosophical definition of trust and integrity (Domain 1, specifically C-002 and C-006)
- Who has authority to grant or revoke access (Domain 2, specifically G-001 and G-003)
- Physical infrastructure specifications (Domain 4)
- Software implementation of security controls (Domain 5)

Domain 3 defines the security architecture, threat model, and policy. Domain 4 implements physical security. Domain 5 implements digital security. Domain 2 defines the human authority that governs security decisions.

### Relationships

- **Domain 1 (Constitution):** Domain 1 defines what "security" means for this institution (not confidentiality-availability-integrity in the abstract, but those concepts as they apply to this institution's specific values and refusals). Domain 3 operationalizes that definition.
- **Domain 2 (Governance):** Domain 2 defines who authorizes access changes. Domain 3 defines the technical policy for how those changes are validated and enforced.
- **Domain 4 (Infrastructure):** Domain 3 defines physical security requirements. Domain 4 implements them. Domain 3 audits Domain 4's implementation.
- **Domain 5 (Platform):** Domain 3 defines software security requirements, cryptographic standards, and audit logging requirements. Domain 5 implements them. Domain 3 audits Domain 5's implementation.

---

## 3.2 Article List

| ID | Title | Summary |
|---|---|---|
| S-001 | **The Institutional Threat Model** | Comprehensive enumeration of threats to the institution: natural disaster, hardware failure, insider compromise, external intrusion attempt, social engineering, supply chain compromise, cryptographic obsolescence, knowledge loss, and institutional capture. |
| S-002 | **Access Control Philosophy** | The principles governing who and what may access institutional resources. Least privilege. Need to know. Default deny. The distinction between authentication (proving identity) and authorization (proving permission). |
| S-003 | **Trust Boundary Definitions** | Formal definition of every trust boundary in the institution: between members, between roles, between systems, between the institution and the physical environment, between the institution and the outside world. |
| S-004 | **The Air-Gap Security Model** | Detailed analysis of what air-gapping provides and what it does not provide. Threat vectors that survive air-gapping (insider threat, supply chain, electromagnetic emanation, physical theft). Compensating controls for each. |
| S-005 | **Cryptographic Principles and Standards** | The institution's cryptographic framework: algorithm selection criteria, key lengths, rotation schedules, the rejection of security through obscurity, and the plan for post-quantum migration. |
| S-006 | **Key Management Lifecycle** | How cryptographic keys are generated, stored, distributed, rotated, and destroyed. Custodial responsibilities. Escrow and recovery. Multi-party key ceremonies. |
| S-007 | **Audit and Verification Framework** | What is audited, how often, by whom, and what constitutes a finding. The distinction between routine audits, triggered audits, and forensic investigations. |
| S-008 | **Incident Response Framework** | Detection, classification, containment, eradication, recovery, and lessons-learned procedures. Roles during an incident. Communication protocols. |
| S-009 | **Data Integrity and Provenance** | How the institution ensures that data has not been tampered with, that its origin is known, and that its chain of custody is documented. Hash chains, signatures, and integrity verification procedures. |
| S-010 | **Insider Threat Considerations** | Acknowledging and addressing the risk posed by authorized members. Separation of duties. Mutual oversight. The balance between trust and verification. |
| S-011 | **Physical Security Principles** | Security requirements for the physical environment: access control to facilities, tamper evidence, environmental monitoring, and destruction procedures for decommissioned equipment. |
| S-012 | **Long-Term Cryptographic Evolution** | A 50-year view of cryptographic risk. How the institution plans for algorithm obsolescence, increasing computational power, and potential breakthroughs in cryptanalysis. Migration strategies. |
| S-013 | **Supply Chain Security** | How the institution evaluates, accepts, and monitors hardware and software entering the institution from external sources. Verification procedures. Quarantine processes. |
| S-014 | **Security Documentation and Classification** | How security-sensitive information within the institution's own documentation is handled. What is restricted. What is available to all members. How classification decisions are made and reviewed. |

---

## 3.3 Priority Order

**Tier 1 -- Write Immediately (existential security foundation):**

1. S-001: The Institutional Threat Model
2. S-002: Access Control Philosophy
3. S-003: Trust Boundary Definitions
4. S-004: The Air-Gap Security Model

**Tier 2 -- Write Early (required before systems are deployed):**

5. S-005: Cryptographic Principles and Standards
6. S-006: Key Management Lifecycle
7. S-009: Data Integrity and Provenance
8. S-011: Physical Security Principles

**Tier 3 -- Write Before Full Operations:**

9. S-007: Audit and Verification Framework
10. S-008: Incident Response Framework
11. S-010: Insider Threat Considerations
12. S-013: Supply Chain Security

**Tier 4 -- Write During Early Operations:**

13. S-014: Security Documentation and Classification
14. S-012: Long-Term Cryptographic Evolution

Rationale: The threat model (S-001) must exist before any security decisions can be justified. Access control and trust boundaries (S-002, S-003) are prerequisites for Domains 4 and 5. The air-gap model (S-004) is specific to this institution's architecture and must be well-understood early. Cryptographic standards and key management (S-005, S-006) must be established before any cryptographic operations begin, because retroactive key management is orders of magnitude harder than doing it right from the start.

---

## 3.4 Dependencies

| Article | Depends On | Depended On By |
|---|---|---|
| S-001 | C-002, C-003, C-006 | S-002, S-003, S-004, S-007, S-008, S-010 |
| S-002 | C-002, G-001, S-001 | S-003, S-005, S-014, P-005 |
| S-003 | S-001, S-002, G-001 | S-004, S-010, S-011, I-005, P-004 |
| S-004 | C-006, S-001, S-003 | S-013, I-005, P-004 |
| S-005 | C-010, S-002 | S-006, S-009, S-012, P-006 |
| S-006 | S-005, G-003 | S-009, S-012 |
| S-007 | S-001, S-002, G-007 | S-008, S-010, S-014 |
| S-008 | C-009, G-009, S-001, S-007 | None (operational) |
| S-009 | C-005, S-005, S-006 | S-007, P-006 |
| S-010 | S-001, S-003, G-003 | S-007 |
| S-011 | S-003, S-004 | I-002, I-005 |
| S-012 | S-005, S-006 | None (long-term planning) |
| S-013 | S-004, C-003 | I-003, P-002 |
| S-014 | S-002, S-007 | None (internal process) |

---

## 3.5 Writing Schedule

| Week | Articles | Notes |
|---|---|---|
| Week 3 | S-001 | Threat model. Can begin once C-002 and C-003 are drafted. |
| Week 4 | S-002, S-003 | Access control philosophy and trust boundaries. Can be drafted in parallel. |
| Week 5 | S-004 | Air-gap model. Requires S-003. |
| Week 6 | S-005, S-011 | Cryptographic principles and physical security. Can be drafted in parallel. |
| Week 7 | S-006, S-009 | Key management and data integrity. Sequential -- S-006 first, then S-009. |
| Week 8 | S-007, S-010 | Audit framework and insider threat. Can be drafted in parallel. |
| Week 9 | S-008, S-013 | Incident response and supply chain. Can be drafted in parallel. |
| Week 11 | S-014 | Security documentation classification. |
| Week 16 | S-012 | Long-term cryptographic evolution. Benefits from having all other crypto articles finalized. |

---

## 3.6 Review Process

1. **First Draft:** Written by the Security Strategist.
2. **Threat Model Validation:** For each article, the reviewer asks: "Does this address the threats identified in S-001?" Any gap is a finding.
3. **Attack Surface Review:** A designated reviewer (ideally someone other than the author) attempts to identify attack vectors that the article fails to address. This is documented even if no gaps are found.
4. **Constitutional Alignment Check:** The Philosophy Architect verifies that security measures do not violate constitutional principles (e.g., that insider threat controls do not create a surveillance culture that violates C-007's principles about the relationship between individuals and the institution).
5. **Implementability Check:** The Infrastructure Lead (Domain 4) and Platform Architect (Domain 5) review to confirm that the security requirements can actually be implemented with available or planned resources.
6. **Red Team Exercise:** For Tier 1 articles, conduct a tabletop exercise where reviewers role-play as adversaries attempting to defeat the described controls.
7. **Ratification:** All security articles require sign-off from the Security Strategist, the Governance Engineer (to confirm authority alignment), and at least one other domain lead.

---

## 3.7 Maintenance Plan

- **Quarterly Review of Threat Model:** S-001 is reviewed every three months during the first two years, then semi-annually. The threat landscape changes; the threat model must keep pace.
- **Annual Cryptographic Review:** S-005, S-006, and S-012 are reviewed annually against current cryptographic research (maintained in an archived local library of relevant publications).
- **Post-Incident Mandatory Review:** After any security incident (as defined in S-008), all articles referenced during the incident are reviewed within 14 days. Findings are recorded in Commentary Sections and may trigger formal revisions.
- **Penetration Testing Correlation:** Whenever a security test is conducted (physical or digital), the results are mapped back to the relevant articles. Gaps between documented controls and actual effectiveness are treated as high-priority findings.
- **Key Lifecycle Events:** Every key rotation, generation, or destruction event triggers a review of S-006 to confirm procedures were followed and the article remains accurate.
- **Five-Year Deep Review:** Every five years, the entire domain undergoes a comprehensive review with the assumption that the threat landscape has fundamentally changed. This review may trigger rewrites rather than amendments.

---
---

# DOMAIN 4: INFRASTRUCTURE & POWER

**Agent Role:** Power & Infrastructure Lead
**Core Purpose:** Define physical infrastructure, power systems, network topology, hardware lifecycle, environmental controls.

---

## 4.1 Domain Map

### Scope

Domain 4 addresses everything physical: the power that runs the systems, the hardware that computes and stores, the environment that protects the hardware, and the physical network that connects components. This domain must plan for a world where every piece of hardware will eventually fail, every power source will eventually degrade, and every physical environment will eventually change.

This domain covers:

- Power generation, storage, distribution, and redundancy
- Hardware selection criteria, lifecycle management, and replacement planning
- Physical network topology (internal only -- this is an air-gapped institution)
- Environmental controls (temperature, humidity, dust, vibration, water)
- Physical site selection, preparation, and long-term suitability
- Backup power and graceful degradation
- Hardware decommissioning and secure destruction
- Spare parts philosophy and inventory management
- Physical maintenance schedules and procedures
- Cable management and physical network documentation

### Boundaries

Domain 4 does NOT cover:

- The philosophical commitment to off-grid operation (Domain 1, specifically C-006)
- Who authorizes infrastructure changes (Domain 2, specifically G-001 and G-003)
- Physical security policy (Domain 3, specifically S-011 -- though Domain 4 implements it)
- Operating system or software that runs on the hardware (Domain 5)

Domain 4 defines the physical substrate. Domain 5 defines what runs on that substrate. Domain 3 defines the security requirements that constrain Domain 4's implementation choices.

### Relationships

- **Domain 1 (Constitution):** Domain 1's commitment to sovereignty, off-grid operation, and long-term durability directly constrains every infrastructure decision. No infrastructure choice may create an external dependency.
- **Domain 2 (Governance):** Domain 2 defines who may authorize infrastructure purchases, changes, and decommissions. Domain 4 executes those decisions.
- **Domain 3 (Security):** Domain 3 defines physical security requirements (S-011), supply chain security requirements (S-013), and trust boundaries that include physical boundaries (S-003). Domain 4 implements all of these.
- **Domain 5 (Platform):** Domain 4 provides the hardware platform. Domain 5 must be designed to run on what Domain 4 provides, and Domain 4 must provide what Domain 5 requires. This is a bidirectional dependency that requires close coordination.

---

## 4.2 Article List

| ID | Title | Summary |
|---|---|---|
| I-001 | **Infrastructure Philosophy and Constraints** | The overarching principles governing all infrastructure decisions: repairability over performance, simplicity over features, availability over cutting-edge, local sourcing where possible, understanding every component. |
| I-002 | **Physical Site Requirements** | What the institution requires from its physical location: environmental stability, security characteristics, accessibility for maintenance, isolation from external threats, and long-term suitability. |
| I-003 | **Hardware Selection Criteria** | How hardware is chosen: preference for open designs, documented interfaces, proven reliability, availability of spare parts, repairability, and resistance to obsolescence. The rejection of black-box hardware. |
| I-004 | **Power Generation and Storage** | Primary and backup power systems. Solar, wind, micro-hydro, generator, and battery bank design principles. Sizing methodology. Redundancy requirements. Fuel storage and rotation (if applicable). |
| I-005 | **Power Distribution and Management** | How power flows from generation to consumption. Circuit design principles. Over-current protection. Monitoring. Load prioritization during power scarcity. Graceful shutdown procedures. |
| I-006 | **Environmental Control Systems** | Temperature, humidity, dust, and vibration management. Passive cooling strategies. Active cooling when necessary. Monitoring and alerting. Seasonal planning. |
| I-007 | **Internal Network Topology** | The physical network connecting institutional systems. Topology design. Cable specifications. Segmentation principles. Documentation requirements. The absolute prohibition on wireless networking (or conditions under which it is permitted, if any). |
| I-008 | **Hardware Lifecycle Management** | How hardware ages, how its health is monitored, when it is replaced, and how transitions are managed. Proactive replacement versus reactive replacement. The spare parts inventory. |
| I-009 | **Graceful Degradation and Minimum Viable Infrastructure** | What happens when components fail. The priority order for systems. The minimum infrastructure required to maintain institutional continuity. Load shedding procedures. |
| I-010 | **Hardware Decommissioning and Secure Destruction** | How hardware is retired. Data sanitization procedures. Physical destruction requirements for storage media. Documentation of decommissioned assets. Environmental responsibility. |
| I-011 | **Maintenance Schedules and Procedures** | Routine maintenance tasks, their frequencies, their procedures, and their documentation requirements. Preventive maintenance philosophy. |
| I-012 | **Spare Parts and Inventory Management** | What spare parts are kept on hand, in what quantities, how they are stored, and how inventory is tracked. The balance between preparedness and waste. |
| I-013 | **Infrastructure Documentation Standards** | How physical infrastructure is documented: diagrams, labels, cable maps, power flow charts, equipment registers. The requirement that any competent person can understand the physical layout from documentation alone. |
| I-014 | **Physical Infrastructure Evolution Path** | How the institution plans for major infrastructure transitions: moving to a new site, upgrading power systems, transitioning to new hardware architectures. Migration planning principles. |

---

## 4.3 Priority Order

**Tier 1 -- Write Immediately (required for any physical setup):**

1. I-001: Infrastructure Philosophy and Constraints
2. I-002: Physical Site Requirements
3. I-004: Power Generation and Storage
4. I-003: Hardware Selection Criteria

**Tier 2 -- Write Before First Hardware Deployment:**

5. I-005: Power Distribution and Management
6. I-007: Internal Network Topology
7. I-006: Environmental Control Systems
8. I-013: Infrastructure Documentation Standards

**Tier 3 -- Write Before Full Operations:**

9. I-008: Hardware Lifecycle Management
10. I-009: Graceful Degradation and Minimum Viable Infrastructure
11. I-010: Hardware Decommissioning and Secure Destruction
12. I-011: Maintenance Schedules and Procedures

**Tier 4 -- Write During Early Operations:**

13. I-012: Spare Parts and Inventory Management
14. I-014: Physical Infrastructure Evolution Path

Rationale: You cannot select a site without I-002, you cannot select hardware without I-003, and you cannot power anything without I-004. The philosophy (I-001) must come first to constrain all subsequent decisions. Documentation standards (I-013) are in Tier 2 because all physical work must be documented from the very beginning -- retrofitting documentation is far harder than maintaining it from day one.

---

## 4.4 Dependencies

| Article | Depends On | Depended On By |
|---|---|---|
| I-001 | C-002, C-006, C-010 | I-002 through I-014 (all Domain 4 articles) |
| I-002 | I-001, S-011 | I-004, I-006, I-007 |
| I-003 | I-001, C-010, S-013 | I-008, I-012, P-001, P-003 |
| I-004 | I-001, I-002 | I-005, I-009 |
| I-005 | I-004, S-003 | I-006, I-009, P-008 |
| I-006 | I-002, I-005 | I-011 |
| I-007 | I-001, I-002, S-003, S-004 | P-004, P-008 |
| I-008 | I-003, I-001 | I-010, I-012 |
| I-009 | C-009, I-004, I-005, G-009 | I-012 |
| I-010 | I-008, S-011 | None (terminal process) |
| I-011 | I-004, I-005, I-006, I-007, I-008 | None (operational) |
| I-012 | I-003, I-008, I-009 | None (operational) |
| I-013 | C-005, I-001 | All Domain 4 articles (documentation standard) |
| I-014 | I-001, I-003, I-008 | None (long-term planning) |

---

## 4.5 Writing Schedule

| Week | Articles | Notes |
|---|---|---|
| Week 3 | I-001, I-013 | Philosophy and documentation standards. Write together to ensure consistency. |
| Week 4 | I-002, I-003 | Site requirements and hardware criteria. Can be drafted in parallel. |
| Week 5 | I-004 | Power generation. Requires I-002 for site context. |
| Week 6 | I-005, I-007 | Power distribution and network topology. Can be drafted in parallel. |
| Week 7 | I-006 | Environmental controls. Requires I-002 and I-005. |
| Week 8 | I-008, I-009 | Hardware lifecycle and graceful degradation. Can be drafted in parallel. |
| Week 9 | I-010, I-011 | Decommissioning and maintenance. Can be drafted in parallel. |
| Week 11 | I-012 | Spare parts inventory. Benefits from I-008 and I-009 being complete. |
| Week 14 | I-014 | Long-term evolution. Benefits from operational experience and all prior articles. |

---

## 4.6 Review Process

1. **First Draft:** Written by the Power & Infrastructure Lead.
2. **Constitutional Compliance Check:** The Philosophy Architect reviews to confirm that no infrastructure decision creates an external dependency or violates sovereignty principles.
3. **Security Review:** The Security Strategist reviews every infrastructure article against S-011 (Physical Security), S-013 (Supply Chain Security), and S-003 (Trust Boundaries).
4. **Platform Compatibility Review:** The Platform Architect reviews I-003, I-007, and I-005 to confirm that the hardware, network, and power meet software platform requirements.
5. **Practical Feasibility Review:** At least one reviewer with hands-on experience in the relevant area (electrical, networking, mechanical) reviews for practical errors, safety issues, and unrealistic assumptions.
6. **Safety Review:** All articles involving power (I-004, I-005) or environmental systems (I-006) undergo a specific safety review. Electrical safety, fire risk, and environmental hazards are explicitly evaluated.
7. **Ratification:** Tier 1 articles require agreement from all founding members. Tier 2 and below require the Infrastructure Lead, the Security Strategist, and the Platform Architect.

---

## 4.7 Maintenance Plan

- **Monthly Physical Inspection:** All physical infrastructure is inspected monthly. Inspection findings are compared against I-011 (Maintenance Schedules) and recorded.
- **Quarterly Power System Review:** I-004 and I-005 are reviewed quarterly against actual power generation and consumption data. Discrepancies between documented capacity and actual performance are flagged.
- **Annual Hardware Inventory:** A complete physical inventory is conducted annually and compared against I-012 (Spare Parts) and I-008 (Hardware Lifecycle). Any discrepancies are investigated.
- **Post-Failure Review:** After any hardware failure, the relevant articles (I-008, I-009, and the specific system article) are reviewed within 7 days. Root cause analysis is documented in the Commentary Section.
- **Pre-Procurement Review:** Before any new hardware is acquired, I-003 (Hardware Selection Criteria) and S-013 (Supply Chain Security) are reviewed to confirm the procurement meets documented standards.
- **Five-Year Infrastructure Assessment:** Every five years, the entire Domain 4 is reviewed with the question: "If we were starting fresh today with what we know now, what would we do differently?" Findings feed into I-014 (Evolution Path).
- **Environmental Monitoring Continuous:** I-006 implicitly requires continuous environmental monitoring. Alerts from monitoring systems trigger immediate review of the relevant article.

---
---

# DOMAIN 5: PLATFORM & CORE SYSTEMS

**Agent Role:** Platform Architect
**Core Purpose:** Define the operating system layer, core services, system architecture, platform decisions, abstraction layers.

---

## 5.1 Domain Map

### Scope

Domain 5 defines the software and systems layer that sits on top of Domain 4's physical infrastructure. It answers: What operating system runs? What core services are provided? How are systems architected for longevity, simplicity, and maintainability? How are platform decisions made and documented? How does the software layer abstract away hardware specifics so that Domain 4 can evolve independently?

This domain covers:

- Operating system selection, configuration, and maintenance philosophy
- Core services definition (storage, compute, identity, time, logging)
- System architecture principles (simplicity, modularity, replaceability)
- Abstraction layers between hardware and application concerns
- Data storage architecture and long-term data format decisions
- Backup and restore systems
- System monitoring and health checks
- Configuration management and reproducibility
- Software supply chain (how software enters the institution, given air-gap)
- Platform evolution and migration planning

### Boundaries

Domain 5 does NOT cover:

- The philosophy of technology (Domain 1, specifically C-010)
- Who authorizes platform changes (Domain 2)
- Security implementation details at the cryptographic level (Domain 3, specifically S-005, S-006)
- Physical hardware or power (Domain 4)

Domain 5 sits between Domain 4 (below, providing hardware) and any future application-layer domains (above, consuming platform services). Domain 3 defines security requirements that Domain 5 must implement.

### Relationships

- **Domain 1 (Constitution):** Domain 1's technology philosophy (C-010) and knowledge preservation theory (C-005) directly constrain platform decisions. The commitment to understanding every layer, rejecting black boxes, and preferring simplicity originates in Domain 1.
- **Domain 2 (Governance):** Domain 2 defines the approval process for platform changes. No platform decision that affects institutional capability may be made without the appropriate governance process.
- **Domain 3 (Security):** Domain 3 defines security requirements that Domain 5 must implement: access control (S-002), cryptographic standards (S-005), audit logging (S-007), and data integrity (S-009). Domain 5 is the primary implementer of Domain 3's digital security requirements.
- **Domain 4 (Infrastructure):** Domain 5 runs on Domain 4's hardware. This creates a bidirectional dependency: Domain 5 must work within Domain 4's constraints, and Domain 4 must provide what Domain 5 requires. The abstraction layer between them is critical and is defined in this domain.

---

## 5.2 Article List

| ID | Title | Summary |
|---|---|---|
| P-001 | **Platform Philosophy and Architectural Principles** | The overarching principles for all platform decisions: prefer boring technology, understand every dependency, minimize the stack, design for replacement, separate mechanism from policy. |
| P-002 | **Software Acquisition and Verification** | How software enters the air-gapped institution. Verification procedures, provenance tracking, quarantine, building from source where feasible. The trust model for external software. |
| P-003 | **Operating System Selection and Configuration** | Criteria for OS selection. Configuration principles: minimal installation, hardened baseline, documented divergence from defaults. The OS as a long-term commitment and migration planning. |
| P-004 | **Network Services Architecture** | Internal network services: DNS (or equivalent), time synchronization, internal certificate authority, service discovery. Design for an air-gapped environment with no external dependencies. |
| P-005 | **Identity and Access Management** | How users and services are identified, authenticated, and authorized at the platform level. Integration with Domain 3's access control philosophy (S-002) and Domain 2's authority model (G-001). |
| P-006 | **Data Storage Architecture** | How data is stored: filesystem choices, database decisions (if any), data format standards for long-term preservation, storage redundancy, and the preference for human-readable formats. |
| P-007 | **Backup and Restore Architecture** | How system state and data are backed up, verified, stored, and restored. Backup frequency, retention, verification testing, and the principle that untested backups are not backups. |
| P-008 | **System Monitoring and Health** | How system health is observed: logging, metrics, alerting. What is monitored, at what thresholds alerts are raised, and how monitoring data is stored and reviewed. |
| P-009 | **Configuration Management and Reproducibility** | How system configuration is documented, versioned, and reproducible. The goal: any authorized person can rebuild any system from documentation alone, without relying on the original builder. |
| P-010 | **Graceful Degradation at the Platform Level** | How the software platform handles hardware failure, resource scarcity, or partial system loss. Service priority ordering. Minimum viable platform. Controlled shutdown procedures. |
| P-011 | **Data Format and Long-Term Preservation Standards** | The institution's standards for data formats: preference for open, documented, text-based formats. Migration plans for binary formats. The 50-year readability test. |
| P-012 | **Platform Evolution and Migration Planning** | How the platform changes over time. The process for evaluating, testing, and executing platform migrations (OS upgrades, architecture changes, new hardware support). |
| P-013 | **The Hardware Abstraction Layer** | How software is insulated from specific hardware. Virtualization, containerization, or neither -- the decision framework. The contract between Domain 4 and Domain 5. |
| P-014 | **Platform Documentation Standards** | How the platform itself is documented: architecture diagrams, service inventories, configuration registries, runbooks. The requirement for documentation that outlives its author. |
| P-015 | **Dependency Management and Minimization** | How software dependencies are tracked, evaluated, and minimized. The institutional preference for fewer dependencies. Vendoring strategy. The risk model for each dependency. |

---

## 5.3 Priority Order

**Tier 1 -- Write Immediately (required for any platform decision):**

1. P-001: Platform Philosophy and Architectural Principles
2. P-003: Operating System Selection and Configuration
3. P-002: Software Acquisition and Verification
4. P-014: Platform Documentation Standards

**Tier 2 -- Write Before First System Deployment:**

5. P-006: Data Storage Architecture
6. P-005: Identity and Access Management
7. P-004: Network Services Architecture
8. P-013: The Hardware Abstraction Layer

**Tier 3 -- Write Before Full Operations:**

9. P-007: Backup and Restore Architecture
10. P-008: System Monitoring and Health
11. P-009: Configuration Management and Reproducibility
12. P-011: Data Format and Long-Term Preservation Standards

**Tier 4 -- Write During Early Operations:**

13. P-010: Graceful Degradation at the Platform Level
14. P-015: Dependency Management and Minimization
15. P-012: Platform Evolution and Migration Planning

Rationale: Platform philosophy (P-001) must precede all decisions. OS selection (P-003) and software acquisition (P-002) are the very first concrete platform decisions. Documentation standards (P-014) must be established before the first system is documented. Data storage (P-006) and identity (P-005) are foundational services that everything else depends on. The hardware abstraction layer (P-013) must be defined before hardware-specific decisions lock in assumptions.

---

## 5.4 Dependencies

| Article | Depends On | Depended On By |
|---|---|---|
| P-001 | C-002, C-010, I-001 | All Domain 5 articles |
| P-002 | C-003, S-013, P-001 | P-003, P-015 |
| P-003 | P-001, P-002, I-003 | P-004, P-005, P-006, P-008, P-009 |
| P-004 | P-003, S-003, S-004, I-007 | P-005, P-008 |
| P-005 | S-002, G-001, P-003, P-004 | P-007, P-008 |
| P-006 | C-005, P-001, P-003, S-005, S-009 | P-007, P-011 |
| P-007 | P-005, P-006 | P-010 |
| P-008 | P-003, P-004, P-005, I-005 | P-010 |
| P-009 | P-001, P-003, P-014 | P-012 |
| P-010 | C-009, I-009, P-007, P-008 | P-012 |
| P-011 | C-005, P-006 | P-012 |
| P-012 | P-009, P-010, P-011, I-014 | None (long-term planning) |
| P-013 | P-001, I-003, I-007 | P-003, P-012 |
| P-014 | C-005, I-013, P-001 | All Domain 5 articles (documentation standard) |
| P-015 | P-001, P-002 | P-003, P-012 |

**Note on Circular Dependencies:** P-003 depends on P-013 (the hardware abstraction layer informs OS selection), and P-013 depends on P-001 which informs P-003. This is resolved by writing P-001 and P-013 first as principles, then making the concrete OS selection in P-003 informed by those principles. P-003 and P-013 should be reviewed together to ensure consistency.

---

## 5.5 Writing Schedule

| Week | Articles | Notes |
|---|---|---|
| Week 3 | P-001, P-014 | Platform philosophy and documentation standards. Write together. |
| Week 5 | P-002 | Software acquisition. Requires S-013 to be in progress. |
| Week 6 | P-013, P-003 | Hardware abstraction layer first, then OS selection. Sequential within the week. |
| Week 7 | P-006, P-005 | Data storage and identity. Can be drafted in parallel. |
| Week 8 | P-004 | Network services. Requires P-003 and I-007. |
| Week 9 | P-007, P-008 | Backup and monitoring. Can be drafted in parallel. |
| Week 10 | P-009, P-011 | Configuration management and data format standards. Can be drafted in parallel. |
| Week 12 | P-010, P-015 | Graceful degradation and dependency management. |
| Week 16 | P-012 | Platform evolution. Last to write; requires all other articles as input. |

---

## 5.6 Review Process

1. **First Draft:** Written by the Platform Architect.
2. **Philosophy Alignment Check:** The Philosophy Architect reviews P-001, P-011, and any article that makes technology selection decisions, to confirm alignment with C-010 (Technology Philosophy) and C-005 (Knowledge Preservation).
3. **Security Review:** The Security Strategist reviews every article for compliance with Domain 3 requirements. Specific focus areas: P-005 (identity must meet S-002), P-006 (storage must meet S-009), P-002 (acquisition must meet S-013).
4. **Infrastructure Compatibility Review:** The Infrastructure Lead reviews P-003, P-004, P-013, and P-008 to confirm compatibility with Domain 4's physical infrastructure.
5. **Reproducibility Test:** For articles that describe system configuration (P-003, P-004, P-005, P-009), a reviewer other than the author attempts to understand the described configuration well enough to replicate it. Gaps are flagged.
6. **Long-Term Viability Assessment:** Each article that selects a specific technology (P-003, P-006) is evaluated against a 10-year viability horizon: Is this technology likely to be maintained? Is migration feasible if it is not?
7. **Ratification:** Tier 1 articles require agreement from all founding members. Tier 2 and below require the Platform Architect, the Security Strategist, and the Infrastructure Lead.

---

## 5.7 Maintenance Plan

- **Monthly System Health Correlation:** Monthly, actual system state is compared against P-008 (Monitoring) expectations and P-009 (Configuration Management) documentation. Drift is documented and corrected.
- **Quarterly Backup Verification:** P-007 (Backup and Restore) requires quarterly restore tests. Results are documented. Failures trigger immediate article review.
- **Semi-Annual Dependency Audit:** P-015 (Dependency Management) requires a semi-annual review of all software dependencies: Are they still maintained upstream? Are there known vulnerabilities? Is a migration needed?
- **Annual Platform Review:** The entire Domain 5 is reviewed annually against the question: "Does our platform still serve the institution's needs, and is it still the simplest adequate solution?"
- **Pre-Change Review:** Before any platform change (OS update, configuration change, new service deployment), the relevant articles are reviewed to confirm the change is consistent with documented principles and procedures.
- **Post-Change Documentation Update:** After any platform change, all affected articles are updated within 7 days. Documentation must never lag behind reality by more than 7 days.
- **Five-Year Migration Assessment:** Every five years, P-012 (Platform Evolution) is used to conduct a comprehensive assessment: Should the institution plan a major platform migration? What technologies on the current platform are approaching end-of-life?

---
---

# CROSS-DOMAIN INTEGRATION

## Master Dependency Graph (Key Cross-Domain Links)

The following table identifies the most critical cross-domain dependencies -- the links where a failure to coordinate between domains will cause systemic problems.

| From | To | Nature of Dependency |
|---|---|---|
| C-002 (Core Principles) | All domain root articles | Every domain must trace decisions back to constitutional principles |
| C-005 (Knowledge Preservation) | All documentation standards (I-013, P-014) | The theory of knowledge preservation governs how all documentation works |
| C-010 (Technology Philosophy) | I-001, P-001 | Infrastructure and platform philosophy must implement constitutional technology principles |
| G-001 (Authority Model) | S-002, P-005 | Access control and identity implement the authority model |
| G-009 (Emergency Governance) | S-008, I-009, P-010 | Emergency procedures across all domains must be consistent with emergency governance |
| S-003 (Trust Boundaries) | I-007, P-004, P-013 | Physical and logical architecture must respect trust boundaries |
| S-005 (Cryptographic Principles) | P-006, S-009 | Data storage and integrity depend on cryptographic standards |
| S-013 (Supply Chain Security) | I-003, P-002 | Hardware and software acquisition must follow supply chain security |
| I-003 (Hardware Selection) | P-003, P-013 | OS and abstraction layer must be compatible with selected hardware |
| I-007 (Network Topology) | P-004, P-008 | Network services and monitoring depend on physical network |

## Cross-Domain Review Requirements

Every article that appears in the "Depended On By" column with cross-domain references MUST be reviewed by the lead of each dependent domain before ratification. This is non-negotiable.

## Master Writing Schedule (All Domains Combined)

| Week | Domain 1 | Domain 2 | Domain 3 | Domain 4 | Domain 5 |
|---|---|---|---|---|---|
| 1 | C-001, C-012 | -- | -- | -- | -- |
| 2 | C-002 | G-001 | -- | -- | -- |
| 3 | C-003 | G-003 | S-001 | I-001, I-013 | P-001, P-014 |
| 4 | C-005, C-010 | G-002, G-009 | S-002, S-003 | I-002, I-003 | -- |
| 5 | C-007 | G-004, G-007 | S-004 | I-004 | P-002 |
| 6 | C-004 | G-005, G-010 | S-005, S-011 | I-005, I-007 | P-013, P-003 |
| 7 | C-006, C-008 | G-013, G-014 | S-006, S-009 | I-006 | P-006, P-005 |
| 8 | C-009, C-011 | G-006, G-011 | S-007, S-010 | I-008, I-009 | P-004 |
| 9 | -- | -- | S-008, S-013 | I-010, I-011 | P-007, P-008 |
| 10 | C-013 | -- | -- | -- | P-009, P-011 |
| 11 | -- | -- | S-014 | I-012 | -- |
| 12 | -- | G-008 | -- | -- | P-010, P-015 |
| 14 | -- | G-012 | -- | I-014 | -- |
| 16 | -- | -- | S-012 | -- | P-012 |

**Total articles across all domains: 13 + 14 + 14 + 14 + 15 = 70 articles**

**Estimated completion of all first drafts: Week 16 (4 months)**

---

# UNIVERSAL REVIEW PROCESS

In addition to domain-specific review processes documented above, all articles across all domains share the following review requirements:

## Mandatory Checks for Every Article

1. **Template Compliance:** Does the article contain all 11 required sections (Title through References)?
2. **Term Definition:** Is every technical or institutional term either defined inline or referenced to a glossary (C-012 or domain-specific)?
3. **Assumption Explicitness:** Are all assumptions stated explicitly? Could a reader in 2076 identify what was assumed and evaluate whether those assumptions still hold?
4. **Cross-Reference Integrity:** Do all references to other articles point to articles that exist (or are scheduled to exist)? Are referenced articles consistent with this article?
5. **Failure Mode Coverage:** Does the Failure Modes section address at least: human error, equipment failure, knowledge loss, and malicious action?
6. **Recovery Completeness:** Does every identified failure mode have a corresponding recovery procedure?
7. **Evolution Path Realism:** Does the Evolution Path section acknowledge specific, plausible ways the subject matter will change?
8. **Plain Language Test:** Can a literate non-specialist understand the article without external help?

## Review Tracking

Every review produces a record containing:

- Reviewer name and role
- Date of review
- Article ID and version reviewed
- Finding classification: Pass, Minor Finding, Major Finding, Blocker
- Description of each finding
- Resolution status

Review records are maintained alongside the articles they reference and are never deleted.

---

# UNIVERSAL MAINTENANCE PLAN

## Ongoing Maintenance Principles

1. **Documentation Is Never Done.** Every article is a living document. The Commentary Section ensures that institutional knowledge accumulates even between formal revisions.
2. **Reality Governs.** If actual practice diverges from documented practice, the divergence must be resolved -- either by changing practice or changing documentation. Undocumented divergence is a finding.
3. **Maintainer Succession.** Every domain must have a documented succession plan for its lead author role. The maintenance plan for each domain includes the responsibility to identify and prepare a successor.
4. **Periodic Deep Review.** Every five years, the entire documentation corpus undergoes a comprehensive review with fresh eyes -- ideally including at least one person who was not involved in the original writing.
5. **Version History.** Every article maintains a version history in its metadata. No change is made without incrementing the version. Major changes increment the major version. Minor changes increment the minor version. Commentary additions do not change the version.

## Maintenance Calendar Summary

| Frequency | Action |
|---|---|
| Continuous | Commentary additions to any article |
| Monthly | System state compared against Domain 4 and Domain 5 documentation |
| Quarterly | Domain 3 threat model review; Domain 5 backup verification |
| Semi-annually | Domain 2 governance review (first 5 years); Domain 5 dependency audit |
| Annually | All domains: full article-by-article review; Domain 1 constitutional re-reading |
| Every 5 years | All domains: deep review with fresh perspective; infrastructure and platform migration assessment |

---

# APPENDIX A: ARTICLE INDEX (ALL DOMAINS)

For reference, the complete list of all 70 proposed articles:

## Domain 1: Constitution & Philosophy (13 articles)

| ID | Title |
|---|---|
| C-001 | The Founding Mandate |
| C-002 | Core Principles and Non-Negotiable Values |
| C-003 | The Refusal Registry |
| C-004 | Institutional Identity Over Time |
| C-005 | The Theory of Knowledge Preservation |
| C-006 | The Ethics of Isolation |
| C-007 | Relationship Between Institution and Individual |
| C-008 | Succession Philosophy |
| C-009 | The Acceptable Loss Doctrine |
| C-010 | Technology Philosophy |
| C-011 | The Amendment Process for Foundational Principles |
| C-012 | Glossary of Foundational Terms |
| C-013 | Failure of Spirit: Recognizing Institutional Drift |

## Domain 2: Governance & Authority (14 articles)

| ID | Title |
|---|---|
| G-001 | The Authority Model |
| G-002 | Decision-Making Frameworks |
| G-003 | Role Definitions and Boundaries |
| G-004 | Dispute Resolution Procedures |
| G-005 | Succession Planning |
| G-006 | The Continuity Protocol |
| G-007 | Record-Keeping Requirements |
| G-008 | Institutional Health Assessment |
| G-009 | Emergency Governance |
| G-010 | The Veto and Override Framework |
| G-011 | External Relations Policy |
| G-012 | Governance Amendment Process |
| G-013 | Meeting and Deliberation Protocols |
| G-014 | Accountability and Transparency Framework |

## Domain 3: Security & Integrity (14 articles)

| ID | Title |
|---|---|
| S-001 | The Institutional Threat Model |
| S-002 | Access Control Philosophy |
| S-003 | Trust Boundary Definitions |
| S-004 | The Air-Gap Security Model |
| S-005 | Cryptographic Principles and Standards |
| S-006 | Key Management Lifecycle |
| S-007 | Audit and Verification Framework |
| S-008 | Incident Response Framework |
| S-009 | Data Integrity and Provenance |
| S-010 | Insider Threat Considerations |
| S-011 | Physical Security Principles |
| S-012 | Long-Term Cryptographic Evolution |
| S-013 | Supply Chain Security |
| S-014 | Security Documentation and Classification |

## Domain 4: Infrastructure & Power (14 articles)

| ID | Title |
|---|---|
| I-001 | Infrastructure Philosophy and Constraints |
| I-002 | Physical Site Requirements |
| I-003 | Hardware Selection Criteria |
| I-004 | Power Generation and Storage |
| I-005 | Power Distribution and Management |
| I-006 | Environmental Control Systems |
| I-007 | Internal Network Topology |
| I-008 | Hardware Lifecycle Management |
| I-009 | Graceful Degradation and Minimum Viable Infrastructure |
| I-010 | Hardware Decommissioning and Secure Destruction |
| I-011 | Maintenance Schedules and Procedures |
| I-012 | Spare Parts and Inventory Management |
| I-013 | Infrastructure Documentation Standards |
| I-014 | Physical Infrastructure Evolution Path |

## Domain 5: Platform & Core Systems (15 articles)

| ID | Title |
|---|---|
| P-001 | Platform Philosophy and Architectural Principles |
| P-002 | Software Acquisition and Verification |
| P-003 | Operating System Selection and Configuration |
| P-004 | Network Services Architecture |
| P-005 | Identity and Access Management |
| P-006 | Data Storage Architecture |
| P-007 | Backup and Restore Architecture |
| P-008 | System Monitoring and Health |
| P-009 | Configuration Management and Reproducibility |
| P-010 | Graceful Degradation at the Platform Level |
| P-011 | Data Format and Long-Term Preservation Standards |
| P-012 | Platform Evolution and Migration Planning |
| P-013 | The Hardware Abstraction Layer |
| P-014 | Platform Documentation Standards |
| P-015 | Dependency Management and Minimization |

---

# APPENDIX B: DOCUMENT VERSIONING SCHEME

All articles follow semantic versioning: **MAJOR.MINOR.PATCH**

- **MAJOR** (e.g., 1.0.0 to 2.0.0): Fundamental change to the article's scope, purpose, or conclusions. Requires full re-review.
- **MINOR** (e.g., 1.0.0 to 1.1.0): Significant addition, clarification, or correction that changes the article's guidance. Requires domain-lead review.
- **PATCH** (e.g., 1.0.0 to 1.0.1): Typographical correction, formatting fix, or minor clarification that does not change the article's guidance. Requires author self-review.
- **Commentary additions** do not change the version number. Commentary is append-only and is always attributed with date and author.

---

# APPENDIX C: HOW TO USE THIS FRAMEWORK

This document is the map. The articles are the territory. To use this framework:

1. **Start with Domain 1, Tier 1.** Write C-001 and C-012. Everything else depends on them.
2. **Follow the Master Writing Schedule.** It is sequenced to respect dependencies. Deviating from it will create articles that reference unwritten foundations.
3. **Use the Article Template.** Every article must contain all 11 sections. An incomplete article is a draft, not a document.
4. **Follow the Review Process.** No article is ratified without review. Unreviewed articles are drafts.
5. **Maintain from Day One.** The maintenance plan is not something that starts after all articles are written. Commentary sections should be used from the moment an article is first drafted.
6. **When in doubt, refer to Domain 1.** The constitution is the root. When two articles conflict, the one closer to Domain 1 prevails unless an explicit override is documented.

This framework expects to produce 70 articles over approximately 16 weeks. Each article will be a substantial document. The total documentation corpus, once complete, will constitute the institutional memory of a self-sovereign digital institution designed to outlive its creators.

---

**End of Stage 1: Documentation Framework -- Domains 1 through 5**
