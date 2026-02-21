"""
BookForge Style Guide Templates

Pre-built style guide templates for different book types.
Each template provides comprehensive guidelines for generating
consistent, genre-appropriate content.
"""

from typing import Dict, List, Any

# =============================================================================
# STYLE TEMPLATES
# =============================================================================

STYLE_TEMPLATES: Dict[str, Dict[str, Any]] = {

    # =========================================================================
    # TEXTBOOK - Academic, Educational Style
    # =========================================================================
    "TEXTBOOK": {
        "name": "Textbook / Academic",
        "description": "Formal educational content designed for learning and retention",

        "tone_and_voice": {
            "primary_tone": "authoritative yet accessible",
            "voice": "third-person academic",
            "formality_level": "formal",
            "personality": "knowledgeable, patient, encouraging",
            "guidelines": [
                "Write with authority while remaining approachable",
                "Use precise, discipline-specific terminology with clear definitions",
                "Maintain objectivity and present balanced perspectives",
                "Employ a teaching voice that guides rather than lectures",
                "Be encouraging without being condescending",
                "Use 'we' to include the reader in the learning journey"
            ]
        },

        "chapter_structure": {
            "opening_elements": [
                "Learning Objectives (3-5 measurable goals)",
                "Chapter Overview (1-2 paragraphs)",
                "Key Terms Preview",
                "Opening Case Study or Scenario (optional)"
            ],
            "body_elements": [
                "Numbered sections with clear headings (3-6 main sections)",
                "Concept explanations with examples",
                "Diagrams, tables, and visual aids descriptions",
                "In-text exercises and reflection questions",
                "Real-world applications and case studies",
                "Cross-references to related chapters"
            ],
            "closing_elements": [
                "Chapter Summary (bulleted key points)",
                "Key Terms and Definitions",
                "Review Questions (mix of recall and application)",
                "Critical Thinking Exercises",
                "Further Reading and Resources",
                "Practice Problems (if applicable)"
            ],
            "recommended_sections_per_chapter": "4-6",
            "subsection_depth": "2-3 levels (e.g., 1.1, 1.1.1)"
        },

        "word_count_targets": {
            "chapter_total": "5000-8000 words",
            "introduction": "300-500 words",
            "main_sections": "800-1500 words each",
            "case_studies": "400-800 words",
            "summary": "300-500 words",
            "paragraphs": "100-200 words",
            "sentences": "15-25 words average"
        },

        "formatting_rules": {
            "headings": "Use hierarchical numbering (1, 1.1, 1.1.1)",
            "terminology": "Bold key terms on first use, define in margin or glossary",
            "lists": "Use numbered lists for sequences, bulleted for non-sequential items",
            "examples": "Indent and label clearly (Example 1.1, Figure 2.3)",
            "citations": "Use appropriate academic citation style (APA, MLA, Chicago)",
            "callouts": [
                "Definition boxes for key terms",
                "Note boxes for important points",
                "Warning boxes for common misconceptions",
                "Think About It boxes for reflection",
                "Try It Yourself boxes for exercises"
            ],
            "visual_elements": "Include placeholder descriptions for diagrams, charts, tables"
        },

        "example_phrases": {
            "introducing_concepts": [
                "In this section, we will explore...",
                "A fundamental principle of [topic] is...",
                "To understand [concept], we must first consider...",
                "Research has consistently demonstrated that...",
                "The concept of [term] refers to..."
            ],
            "transitions": [
                "Building on this foundation...",
                "Having established the basic principles...",
                "This leads us to an important consideration...",
                "In contrast to the previous approach...",
                "A related concept is..."
            ],
            "examples": [
                "Consider the following example...",
                "To illustrate this principle...",
                "This concept is exemplified by...",
                "A practical application of this is...",
                "As demonstrated in the case of..."
            ],
            "summaries": [
                "In summary, the key points are...",
                "To recapitulate the main ideas...",
                "The essential takeaways from this section include...",
                "This chapter has established that..."
            ]
        },

        "avoid": [
            "Overly casual language or slang",
            "Personal anecdotes (unless pedagogically valuable)",
            "Assumptions about prior knowledge without explanation",
            "Jargon without definition",
            "Overly complex sentences that obscure meaning",
            "Passive voice overuse (prefer active where possible)",
            "Condescending or patronizing tone",
            "Unsubstantiated claims",
            "Cultural or regional references that may not translate",
            "Gendered language (use they/their for singular)",
            "Rhetorical questions without purpose"
        ],

        "special_considerations": {
            "accessibility": "Define all acronyms, provide alt-text descriptions for visuals",
            "scaffolding": "Build complexity gradually within and across chapters",
            "assessment_alignment": "Ensure review questions align with learning objectives",
            "engagement": "Include interactive elements every 2-3 pages of content"
        }
    },

    # =========================================================================
    # NOVEL - Fiction Narrative Style
    # =========================================================================
    "NOVEL": {
        "name": "Novel / Fiction",
        "description": "Engaging narrative fiction with immersive storytelling",

        "tone_and_voice": {
            "primary_tone": "immersive and emotionally engaging",
            "voice": "character-driven (first or third person based on POV choice)",
            "formality_level": "varies by genre and character",
            "personality": "authentic to characters, evocative, compelling",
            "guidelines": [
                "Show, don't tell - use action and dialogue to reveal character",
                "Maintain consistent point of view throughout scenes",
                "Use sensory details to create vivid settings",
                "Balance dialogue, action, and introspection",
                "Let characters drive the narrative",
                "Create distinct voices for each character",
                "Match prose rhythm to emotional beats"
            ]
        },

        "chapter_structure": {
            "opening_elements": [
                "Hook - compelling first line or paragraph",
                "Scene-setting - establish time, place, situation",
                "Character grounding - whose POV, what's their state"
            ],
            "body_elements": [
                "Rising action with clear scene goals",
                "Conflict and tension (internal and/or external)",
                "Character development through action and choice",
                "Dialogue that reveals character and advances plot",
                "Sensory details and environmental immersion",
                "Subplot advancement where relevant"
            ],
            "closing_elements": [
                "Scene resolution or escalation",
                "Emotional beat or moment of change",
                "Hook or question that propels reader forward",
                "Transition to next chapter"
            ],
            "scene_structure": "Scene-sequel pattern: Goal, Conflict, Disaster; Reaction, Dilemma, Decision",
            "pacing_notes": "Alternate high-tension and breathing room chapters"
        },

        "word_count_targets": {
            "chapter_total": "2500-5000 words (genre dependent)",
            "scenes": "1000-2500 words each",
            "paragraphs": "50-150 words (shorter for action, longer for introspection)",
            "dialogue_exchanges": "Keep snappy - 1-3 sentences per speech",
            "novel_total": {
                "literary_fiction": "70000-100000 words",
                "commercial_fiction": "80000-100000 words",
                "romance": "50000-80000 words",
                "thriller": "80000-100000 words",
                "fantasy_scifi": "90000-120000 words",
                "young_adult": "50000-80000 words"
            }
        },

        "formatting_rules": {
            "paragraphs": "New paragraph for new speaker, new action, new idea",
            "dialogue": "Use quotation marks, new line for each speaker",
            "thoughts": "Italics for direct internal thoughts",
            "emphasis": "Use italics sparingly for emphasis",
            "scene_breaks": "Use # or *** centered for scene breaks within chapters",
            "chapter_openings": "May include epigraphs, dates, or location markers",
            "flashbacks": "Use clear temporal markers and tense consistency"
        },

        "example_phrases": {
            "opening_hooks": [
                "[Character] had always known this day would come.",
                "The [object] arrived on a [day] that would change everything.",
                "Three things happened the night [event].",
                "Later, [Character] would wonder if things might have been different if..."
            ],
            "action_beats": [
                "[He/She] [action verb + adverb or specific detail].",
                "The [sound/sight/smell] of [specific detail] filled the [space].",
                "[Character] moved [how], [additional action]."
            ],
            "dialogue_tags": [
                "said (primary - invisible to readers)",
                "asked, replied, answered (occasional variation)",
                "Action beats instead of tags: [Character] set down the cup. 'I can't do this anymore.'"
            ],
            "chapter_endings": [
                "And that's when [twist or revelation].",
                "[Character] didn't know it yet, but [foreshadowing].",
                "The [object/person/moment] was gone. And nothing would ever be the same."
            ]
        },

        "avoid": [
            "Purple prose - overly elaborate or flowery language",
            "Info-dumping - large blocks of exposition or backstory",
            "Said bookisms - 'he ejaculated,' 'she expostulated'",
            "Adverb overuse - 'said angrily,' 'walked slowly'",
            "Telling emotions - 'She felt sad' instead of showing",
            "Inconsistent POV - head-hopping within scenes",
            "Clichés and overused metaphors",
            "Predictable or convenient plot solutions",
            "Undifferentiated dialogue - all characters sounding the same",
            "Excessive internal monologue that stalls pacing",
            "Simultaneous actions with 'as' or 'while' overuse",
            "Starting sentences with 'Suddenly' or 'Then'"
        ],

        "special_considerations": {
            "genre_conventions": "Honor expected tropes while finding fresh approaches",
            "character_arcs": "Track emotional and goal-based progression",
            "pacing": "Vary chapter lengths for rhythm; short chapters increase tension",
            "world_building": "Integrate naturally through character experience",
            "theme": "Let themes emerge through story, don't preach"
        }
    },

    # =========================================================================
    # SELF_HELP - Motivational, Actionable Style
    # =========================================================================
    "SELF_HELP": {
        "name": "Self-Help / Personal Development",
        "description": "Motivational content focused on practical transformation",

        "tone_and_voice": {
            "primary_tone": "empowering, supportive, action-oriented",
            "voice": "direct second-person (you/your)",
            "formality_level": "conversational professional",
            "personality": "encouraging, authentic, relatable, confident",
            "guidelines": [
                "Address the reader directly as 'you'",
                "Balance empathy with challenge - acknowledge struggles, then push growth",
                "Use personal stories to build connection and credibility",
                "Be specific and actionable - vague advice frustrates readers",
                "Maintain optimism without toxic positivity",
                "Acknowledge difficulties while emphasizing agency",
                "Write as a trusted mentor, not a guru on a pedestal"
            ]
        },

        "chapter_structure": {
            "opening_elements": [
                "Compelling hook - story, question, or provocative statement",
                "The problem - what pain point does this chapter address?",
                "The promise - what will the reader gain?",
                "Why it matters - connect to larger goals"
            ],
            "body_elements": [
                "Core concept or principle explanation",
                "Supporting evidence (research, case studies, testimonials)",
                "Personal story or anecdote illustrating the principle",
                "Step-by-step framework or methodology",
                "Common obstacles and how to overcome them",
                "Exercises, prompts, or reflection questions"
            ],
            "closing_elements": [
                "Key takeaways (3-5 bullet points)",
                "Action steps - specific things to do NOW",
                "Journaling prompts or reflection questions",
                "Preview of next chapter (optional)",
                "Motivational send-off"
            ],
            "recommended_framework": "Use memorable acronyms or numbered steps (The 5 Steps, The ABC Method)"
        },

        "word_count_targets": {
            "chapter_total": "3000-5000 words",
            "introduction_hook": "200-400 words",
            "main_content": "2000-3500 words",
            "action_section": "400-600 words",
            "paragraphs": "75-150 words",
            "book_total": "40000-60000 words"
        },

        "formatting_rules": {
            "headings": "Use benefit-driven or curiosity-inducing headings",
            "lists": "Frequent use of numbered steps and bulleted lists",
            "callouts": [
                "Action Step boxes",
                "Reflection prompts",
                "Key Insight highlights",
                "Real-life example boxes",
                "Quick Win suggestions"
            ],
            "exercises": "Clearly formatted with space for writing (in print) or clear prompts",
            "stories": "Open with scene-setting, close with lesson learned",
            "emphasis": "Bold key concepts and actionable phrases"
        },

        "example_phrases": {
            "opening_hooks": [
                "What if everything you believed about [topic] was wrong?",
                "Here's a truth that might be hard to hear...",
                "Let me tell you about the moment everything changed for me.",
                "You're about to learn the one thing that separates [successful people] from everyone else."
            ],
            "empathy_statements": [
                "If you're struggling with [issue], you're not alone.",
                "I know how it feels when...",
                "This might be the hardest chapter in this book, but it's also the most important.",
                "You've probably tried [common approach] and felt frustrated when..."
            ],
            "action_language": [
                "Here's exactly what I want you to do...",
                "Starting today, commit to...",
                "Your action step: [specific behavior]",
                "Put this book down and [immediate action]",
                "This week, I challenge you to..."
            ],
            "motivational_closings": [
                "You have everything you need to begin. The only question is: will you?",
                "Every expert was once a beginner. Your journey starts now.",
                "Remember: progress, not perfection.",
                "The best time to start was yesterday. The second best time is now."
            ]
        },

        "avoid": [
            "Vague advice without specific action steps",
            "Preaching or talking down to the reader",
            "Toxic positivity - dismissing real struggles",
            "Promising quick fixes or overnight transformation",
            "Overuse of clichés ('live your best life,' 'be your authentic self')",
            "Name-dropping without substance",
            "Excessive self-promotion",
            "Shaming or guilt-tripping the reader",
            "Complex jargon or overly academic language",
            "Long tangents without clear purpose",
            "Blaming external circumstances without acknowledging reader agency"
        ],

        "special_considerations": {
            "credibility": "Back claims with research, statistics, or verifiable examples",
            "accessibility": "Make frameworks simple enough to remember and apply",
            "progressions": "Build complexity - early chapters should set foundation",
            "accountability": "Include tracking tools, checklists, or commitment devices",
            "variety": "Mix teaching styles - stories, data, exercises, reflection"
        }
    },

    # =========================================================================
    # TECHNICAL - Documentation, How-To Style
    # =========================================================================
    "TECHNICAL": {
        "name": "Technical / How-To Guide",
        "description": "Clear, precise documentation for technical procedures and systems",

        "tone_and_voice": {
            "primary_tone": "clear, precise, helpful",
            "voice": "second-person imperative (you/your) or neutral instructional",
            "formality_level": "professional but accessible",
            "personality": "knowledgeable, patient, methodical",
            "guidelines": [
                "Prioritize clarity above all else",
                "Use consistent terminology throughout",
                "Write in active voice with imperative mood for instructions",
                "Assume minimum prior knowledge, build progressively",
                "Provide context before procedures",
                "Anticipate questions and address them proactively",
                "Test all procedures before documenting"
            ]
        },

        "chapter_structure": {
            "opening_elements": [
                "Chapter objectives (what reader will be able to do)",
                "Prerequisites (skills, tools, access needed)",
                "Overview (what this chapter covers and why it matters)",
                "Time estimate (how long procedures typically take)"
            ],
            "body_elements": [
                "Conceptual overview before procedures",
                "Step-by-step numbered procedures",
                "Code samples or command examples",
                "Screenshots or diagram descriptions",
                "Expected outputs and verification steps",
                "Troubleshooting common issues",
                "Best practices and tips"
            ],
            "closing_elements": [
                "Summary of what was accomplished",
                "Verification checklist",
                "Next steps or related topics",
                "Quick reference (commands, syntax, etc.)",
                "Additional resources and documentation links"
            ],
            "procedure_format": "Numbered steps with clear start/end points"
        },

        "word_count_targets": {
            "chapter_total": "2000-4000 words",
            "concept_explanations": "200-500 words",
            "procedures": "500-1500 words per procedure",
            "troubleshooting_sections": "300-600 words",
            "steps": "20-50 words each (one action per step)",
            "paragraphs": "50-100 words"
        },

        "formatting_rules": {
            "code": "Use monospace formatting, proper syntax highlighting indicators",
            "commands": "Display on separate lines, indicate user input vs output",
            "file_paths": "Use monospace, indicate platform-specific variations",
            "ui_elements": "Bold for buttons, menus; use consistent naming",
            "variables": "Use angle brackets or italics for placeholders <your-value>",
            "callouts": [
                "Note - additional helpful information",
                "Warning - potential problems or data loss risks",
                "Tip - efficiency improvements or best practices",
                "Important - critical information that must not be missed",
                "Example - illustrative use case"
            ],
            "lists": "Use numbered lists for sequential steps, bullets for options/items"
        },

        "example_phrases": {
            "introducing_concepts": [
                "Before we begin, let's understand...",
                "[Feature] allows you to...",
                "In this section, you'll learn how to...",
                "This approach is useful when..."
            ],
            "procedure_starts": [
                "To [accomplish goal], follow these steps:",
                "Complete the following procedure to [outcome]:",
                "This section walks you through [process]."
            ],
            "steps": [
                "1. Open [application/file].",
                "2. Navigate to [location].",
                "3. Enter the following command: `[command]`",
                "4. Verify that [expected result].",
                "5. Click **Save** to apply your changes."
            ],
            "troubleshooting": [
                "If you encounter [error], try...",
                "Common causes include:",
                "To resolve this issue:",
                "Verify that [prerequisite] before proceeding."
            ],
            "notes_and_warnings": [
                "Note: This step requires administrator privileges.",
                "Warning: This action cannot be undone.",
                "Tip: You can use [shortcut] to speed up this process.",
                "Important: Back up your data before proceeding."
            ]
        },

        "avoid": [
            "Ambiguous instructions ('simply,' 'just,' 'easily')",
            "Assuming knowledge without verification",
            "Multiple actions in a single step",
            "Unclear antecedents ('click it,' 'run this')",
            "Humor that may confuse or date poorly",
            "Wordiness - every word should serve a purpose",
            "Inconsistent terminology (pick one term and stick with it)",
            "Undocumented prerequisites",
            "Steps without verification/confirmation points",
            "Screenshots without context or labels",
            "Platform-specific assumptions without noting them"
        ],

        "special_considerations": {
            "versioning": "Note software versions, dates, and update requirements",
            "accessibility": "Provide text alternatives for visual instructions",
            "searchability": "Use consistent keywords for discoverability",
            "modularity": "Write procedures that can stand alone or be combined",
            "testing": "All procedures should be verified before publication"
        }
    },

    # =========================================================================
    # CHILDREN - Simple, Engaging Style
    # =========================================================================
    "CHILDREN": {
        "name": "Children's Book",
        "description": "Age-appropriate content that engages and delights young readers",

        "tone_and_voice": {
            "primary_tone": "warm, playful, wonder-filled",
            "voice": "friendly narrator or character-driven",
            "formality_level": "casual and accessible",
            "personality": "kind, curious, imaginative, reassuring",
            "guidelines": [
                "Use simple, concrete language appropriate to age level",
                "Create rhythm and flow that's pleasing when read aloud",
                "Include repetition for engagement and predictability",
                "Balance excitement with reassurance",
                "Respect children's intelligence while meeting them at their level",
                "Make abstract concepts concrete through familiar experiences",
                "Incorporate sensory details children can relate to"
            ]
        },

        "chapter_structure": {
            "picture_book": {
                "pages": "32 pages standard (28-40 range)",
                "text_per_page": "1-3 sentences (under 100 words per spread)",
                "page_turns": "Use for reveals, surprises, pacing",
                "structure": "Beginning hook, rising action, climax, resolution"
            },
            "early_reader": {
                "chapters": "Short chapters (300-500 words)",
                "structure": "Clear beginning, middle, end per chapter",
                "cliffhangers": "Gentle hooks to encourage continued reading"
            },
            "middle_grade": {
                "chapters": "1500-3000 words",
                "structure": "Full chapter arcs with series-level progression",
                "subplots": "Age-appropriate secondary storylines"
            },
            "general_elements": [
                "Strong opening that establishes character and situation",
                "Clear problem or goal",
                "Obstacles and attempts",
                "Resolution with emotional satisfaction",
                "Character growth or lesson (show, don't preach)"
            ]
        },

        "word_count_targets": {
            "picture_book": "500-1000 words total",
            "early_reader": "1000-5000 words total",
            "chapter_book": "10000-20000 words total",
            "middle_grade": "30000-50000 words total",
            "sentences": {
                "picture_book": "5-10 words",
                "early_reader": "8-15 words",
                "chapter_book": "10-18 words",
                "middle_grade": "12-20 words"
            }
        },

        "formatting_rules": {
            "sentences": "Short and clear, one idea per sentence for younger readers",
            "paragraphs": "Brief paragraphs, more white space for younger readers",
            "dialogue": "Clear speaker identification, natural child-like speech",
            "vocabulary": "Age-appropriate with occasional stretch words in context",
            "repetition": "Intentional use for rhythm, emphasis, and memorability",
            "sound_words": "Onomatopoeia for engagement (BOOM, splish-splash, WHOOSH)",
            "page_breaks": "Consider where illustrations would appear"
        },

        "example_phrases": {
            "openings": [
                "Once upon a time, in a [place] far away...",
                "[Character name] was no ordinary [character type].",
                "On the day everything changed, [Character] woke up to...",
                "There was something different about today."
            ],
            "action_and_emotion": [
                "[Character]'s heart went thump-thump-thump.",
                "'Oh no!' [Character] cried. 'What will I do?'",
                "[Character] had an idea. A wonderful, magnificent idea!",
                "And then... [page turn for reveal]"
            ],
            "rhythmic_patterns": [
                "Over the bridge, through the trees, past the pond, if you please.",
                "Big ones, small ones, short ones, tall ones.",
                "They looked high. They looked low. They looked everywhere they could go."
            ],
            "resolutions": [
                "And from that day on, [Character] knew that...",
                "[Character] smiled. Everything was going to be okay.",
                "And they all [lived happily ever after / celebrated together / learned something important].",
                "'I did it!' [Character] cheered. And indeed, [he/she/they] had."
            ]
        },

        "avoid": [
            "Talking down to children or being condescending",
            "Heavy-handed moral lessons (show through story instead)",
            "Complex vocabulary without context clues",
            "Scary content without resolution and reassurance",
            "Adult humor that goes over children's heads",
            "Long descriptive passages that lose young attention spans",
            "Passive voice (children prefer active subjects)",
            "Abstract concepts without concrete examples",
            "Cultural references children won't understand",
            "Outdated or stereotypical character representations",
            "Complicated plots with too many characters",
            "Didactic or preachy endings"
        ],

        "special_considerations": {
            "age_appropriateness": {
                "0-3": "Concept books, simple stories, heavy repetition",
                "3-5": "Picture books, clear emotions, wish-fulfillment",
                "5-7": "Early readers, problem-solving, friendship themes",
                "7-10": "Chapter books, adventure, self-discovery",
                "10-12": "Middle grade, complexity, identity, relationships"
            },
            "read_aloud": "Test by reading aloud - rhythm matters",
            "illustrations": "Leave room for visual storytelling in picture books",
            "diversity": "Include diverse characters and experiences naturally",
            "emotional_safety": "Handle difficult topics with care and resolution"
        }
    },

    # =========================================================================
    # BUSINESS - Professional, Corporate Style
    # =========================================================================
    "BUSINESS": {
        "name": "Business / Professional",
        "description": "Corporate and professional content for business audiences",

        "tone_and_voice": {
            "primary_tone": "confident, credible, results-oriented",
            "voice": "professional first-person plural (we) or direct second-person (you)",
            "formality_level": "professional formal",
            "personality": "authoritative, strategic, pragmatic",
            "guidelines": [
                "Lead with value and business impact",
                "Use data and evidence to support claims",
                "Be concise - business readers have limited time",
                "Focus on outcomes and ROI",
                "Balance confidence with realistic assessment",
                "Use industry terminology appropriately",
                "Maintain credibility through precision and accuracy"
            ]
        },

        "chapter_structure": {
            "opening_elements": [
                "Executive summary or key takeaways upfront",
                "Business context and relevance",
                "Chapter objectives and outcomes",
                "Quick-reference framework preview"
            ],
            "body_elements": [
                "Core concepts with business application",
                "Case studies from real organizations",
                "Data, metrics, and research findings",
                "Implementation frameworks and models",
                "Common challenges and solutions",
                "Best practices and benchmarks",
                "Tools and templates"
            ],
            "closing_elements": [
                "Key takeaways (3-5 bullet points)",
                "Action checklist",
                "Implementation timeline",
                "Metrics for measuring success",
                "Resources and tools"
            ],
            "recommended_structure": "Pyramid principle - lead with conclusion, then support"
        },

        "word_count_targets": {
            "chapter_total": "4000-6000 words",
            "executive_summary": "150-300 words",
            "main_sections": "800-1200 words each",
            "case_studies": "500-800 words",
            "action_items": "200-400 words",
            "paragraphs": "75-125 words",
            "sentences": "15-25 words average",
            "book_total": "50000-70000 words"
        },

        "formatting_rules": {
            "headings": "Action-oriented or benefit-driven headings",
            "lists": "Frequent bullet points and numbered lists for scannability",
            "emphasis": "Bold for key terms and critical points",
            "data": "Present in tables, charts, or callout boxes",
            "callouts": [
                "Key Insight boxes",
                "Case Study sections",
                "Action Item checklists",
                "Expert Quote highlights",
                "Metric/KPI callouts",
                "Quick Reference summaries"
            ],
            "visuals": "Include descriptions for charts, graphs, frameworks",
            "white_space": "Use for readability and emphasis"
        },

        "example_phrases": {
            "opening_hooks": [
                "In today's competitive landscape, organizations that [action] achieve [result].",
                "Research shows that companies implementing [strategy] see [specific metric improvement].",
                "The difference between thriving organizations and struggling ones often comes down to...",
                "What separates market leaders from the competition?"
            ],
            "presenting_data": [
                "According to [Source], [statistic].",
                "Our research indicates that [finding].",
                "Organizations that [action] report [X]% improvement in [metric].",
                "The data reveals a clear pattern:"
            ],
            "case_study_intros": [
                "Consider how [Company] transformed their approach...",
                "When [Company] faced [challenge], they...",
                "A leading [industry] company discovered that..."
            ],
            "action_oriented": [
                "To implement this strategy, your team should...",
                "The first step is to assess your current [area].",
                "Begin by identifying [key factor].",
                "Measure success through [specific metrics]."
            ],
            "transitions": [
                "Building on this foundation...",
                "This principle directly applies to...",
                "From a practical standpoint...",
                "The implications for your organization are significant."
            ]
        },

        "avoid": [
            "Buzzword overload without substance",
            "Vague claims without data or evidence",
            "Overly academic or theoretical content",
            "Excessive jargon that obscures meaning",
            "Long-winded explanations - get to the point",
            "Promising unrealistic results",
            "Ignoring practical implementation challenges",
            "One-size-fits-all recommendations",
            "Dated examples or references",
            "Ignoring diverse organizational contexts",
            "Preaching without acknowledging business realities",
            "Excessive use of passive voice"
        ],

        "special_considerations": {
            "audience_segmentation": "May need variations for C-suite, managers, practitioners",
            "industry_relevance": "Adapt examples to reader's industry context",
            "actionability": "Every chapter should leave readers with concrete next steps",
            "credibility": "Cite sources, include research, reference real organizations",
            "time_sensitivity": "Business readers skim - make content scannable",
            "global_considerations": "Account for cultural and regional business differences"
        }
    }
}


# =============================================================================
# UTILITY FUNCTIONS
# =============================================================================

def get_style_template(style_name: str) -> Dict[str, Any]:
    """
    Retrieve a style template by name.

    Args:
        style_name: The name of the style template (e.g., 'TEXTBOOK', 'NOVEL')

    Returns:
        The style template dictionary, or empty dict if not found.
    """
    return STYLE_TEMPLATES.get(style_name.upper(), {})


def get_all_style_names() -> List[str]:
    """
    Get a list of all available style template names.

    Returns:
        List of style template names.
    """
    return list(STYLE_TEMPLATES.keys())


def format_style_for_prompt(style_name: str, sections: List[str] = None) -> str:
    """
    Format a style template as a string suitable for inclusion in a prompt.

    Args:
        style_name: The name of the style template
        sections: Optional list of sections to include. If None, includes all.
                  Valid sections: 'tone_and_voice', 'chapter_structure',
                  'word_count_targets', 'formatting_rules', 'example_phrases', 'avoid'

    Returns:
        Formatted string representation of the style guide.
    """
    template = get_style_template(style_name)
    if not template:
        return f"Style template '{style_name}' not found."

    if sections is None:
        sections = ['tone_and_voice', 'chapter_structure', 'word_count_targets',
                    'formatting_rules', 'example_phrases', 'avoid']

    output_lines = [
        f"# {template['name']} Style Guide",
        f"\n{template['description']}\n"
    ]

    for section in sections:
        if section in template:
            output_lines.append(f"\n## {section.replace('_', ' ').title()}")
            content = template[section]

            if isinstance(content, dict):
                for key, value in content.items():
                    if isinstance(value, list):
                        output_lines.append(f"\n### {key.replace('_', ' ').title()}")
                        for item in value:
                            output_lines.append(f"  - {item}")
                    elif isinstance(value, dict):
                        output_lines.append(f"\n### {key.replace('_', ' ').title()}")
                        for k, v in value.items():
                            if isinstance(v, list):
                                output_lines.append(f"  {k}:")
                                for item in v:
                                    output_lines.append(f"    - {item}")
                            else:
                                output_lines.append(f"  {k}: {v}")
                    else:
                        output_lines.append(f"  {key.replace('_', ' ').title()}: {value}")
            elif isinstance(content, list):
                for item in content:
                    output_lines.append(f"  - {item}")

    return "\n".join(output_lines)


def get_style_summary(style_name: str) -> str:
    """
    Get a brief summary of a style template.

    Args:
        style_name: The name of the style template

    Returns:
        Brief summary string.
    """
    template = get_style_template(style_name)
    if not template:
        return f"Style template '{style_name}' not found."

    tone = template.get('tone_and_voice', {})
    return f"""
{template['name']}
{'-' * len(template['name'])}
{template['description']}

Primary Tone: {tone.get('primary_tone', 'N/A')}
Voice: {tone.get('voice', 'N/A')}
Formality: {tone.get('formality_level', 'N/A')}
"""


# =============================================================================
# MAIN - Display available styles when run directly
# =============================================================================

if __name__ == "__main__":
    print("=" * 60)
    print("BookForge Style Templates")
    print("=" * 60)
    print("\nAvailable Style Templates:")
    print("-" * 40)

    for name in get_all_style_names():
        template = STYLE_TEMPLATES[name]
        print(f"\n{name}")
        print(f"  {template['name']}")
        print(f"  {template['description']}")

    print("\n" + "=" * 60)
    print("Usage Examples:")
    print("-" * 40)
    print("""
from style_templates import STYLE_TEMPLATES, get_style_template, format_style_for_prompt

# Get a specific template
textbook_style = get_style_template('TEXTBOOK')

# Get tone guidelines
tone = textbook_style['tone_and_voice']['guidelines']

# Format for prompt inclusion
prompt_text = format_style_for_prompt('NOVEL', sections=['tone_and_voice', 'avoid'])
""")
