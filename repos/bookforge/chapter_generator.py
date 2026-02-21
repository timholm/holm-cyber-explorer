#!/usr/bin/env python3
"""
Chapter Generation System for BookForge

A robust chapter generation system that uses Ollama to create complete,
well-structured book chapters with consistency tracking and streaming output.
"""

import json
import re
import time
import logging
from dataclasses import dataclass, field
from typing import Optional, Generator, Callable
from enum import Enum
import requests

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


class BookType(Enum):
    """Types of books for content generation."""
    FICTION = "fiction"
    NON_FICTION = "non_fiction"
    TECHNICAL = "technical"
    SELF_HELP = "self_help"


@dataclass
class StyleGuide:
    """Configuration for chapter styling and formatting."""
    book_type: BookType = BookType.NON_FICTION
    target_word_count_min: int = 5000
    target_word_count_max: int = 7000
    tone: str = "professional yet accessible"
    reading_level: str = "intermediate"
    include_examples: bool = True
    include_callouts: bool = True
    include_key_terms: bool = True
    include_exercises: bool = False
    chapter_structure: list = field(default_factory=lambda: [
        "opening_hook",
        "introduction",
        "main_sections",
        "examples",
        "summary",
        "key_terms"
    ])
    custom_instructions: str = ""

    def to_prompt_context(self) -> str:
        """Convert style guide to prompt context string."""
        elements = []
        if self.include_examples:
            elements.append("practical examples")
        if self.include_callouts:
            elements.append("callout boxes for important points")
        if self.include_key_terms:
            elements.append("key terms definitions")
        if self.include_exercises:
            elements.append("exercises or reflection questions")

        return f"""
Book Type: {self.book_type.value}
Target Word Count: {self.target_word_count_min}-{self.target_word_count_max} words
Tone: {self.tone}
Reading Level: {self.reading_level}
Include: {', '.join(elements) if elements else 'standard content'}
{f'Special Instructions: {self.custom_instructions}' if self.custom_instructions else ''}
"""


@dataclass
class ChapterOutline:
    """Structure for chapter outline information."""
    chapter_number: int
    title: str
    main_topic: str
    subtopics: list[str]
    key_points: list[str]
    learning_objectives: list[str] = field(default_factory=list)
    connections_to_other_chapters: list[str] = field(default_factory=list)

    def to_prompt_context(self) -> str:
        """Convert outline to prompt context string."""
        subtopics_str = "\n".join(f"  - {s}" for s in self.subtopics)
        key_points_str = "\n".join(f"  - {p}" for p in self.key_points)
        objectives_str = "\n".join(f"  - {o}" for o in self.learning_objectives) if self.learning_objectives else "N/A"
        connections_str = "\n".join(f"  - {c}" for c in self.connections_to_other_chapters) if self.connections_to_other_chapters else "N/A"

        return f"""
Chapter {self.chapter_number}: {self.title}

Main Topic: {self.main_topic}

Subtopics to Cover:
{subtopics_str}

Key Points:
{key_points_str}

Learning Objectives:
{objectives_str}

Connections to Other Chapters:
{connections_str}
"""


@dataclass
class PreviousChaptersSummary:
    """Summary of previous chapters for consistency."""
    summaries: list[dict]  # [{chapter_num, title, summary, key_concepts, characters/entities}]
    recurring_themes: list[str]
    established_terminology: dict[str, str]  # term -> definition
    narrative_threads: list[str]  # For fiction

    def to_prompt_context(self) -> str:
        """Convert to prompt context for consistency."""
        if not self.summaries:
            return "This is the first chapter - no previous content to reference."

        chapters_str = ""
        for ch in self.summaries:
            chapters_str += f"\nChapter {ch['chapter_num']}: {ch['title']}\n"
            chapters_str += f"Summary: {ch['summary']}\n"
            if ch.get('key_concepts'):
                chapters_str += f"Key Concepts: {', '.join(ch['key_concepts'])}\n"

        themes_str = ", ".join(self.recurring_themes) if self.recurring_themes else "None established yet"

        terms_str = ""
        if self.established_terminology:
            terms_str = "\n".join(f"  - {term}: {defn}" for term, defn in self.established_terminology.items())
        else:
            terms_str = "None established yet"

        threads_str = ""
        if self.narrative_threads:
            threads_str = f"\nNarrative Threads:\n" + "\n".join(f"  - {t}" for t in self.narrative_threads)

        return f"""
PREVIOUS CHAPTERS:
{chapters_str}

Recurring Themes: {themes_str}

Established Terminology:
{terms_str}
{threads_str}
"""


@dataclass
class ChapterSection:
    """Represents a section of a generated chapter."""
    section_type: str
    content: str
    word_count: int


@dataclass
class GeneratedChapter:
    """Complete generated chapter with metadata."""
    chapter_number: int
    title: str
    full_content: str
    sections: list[ChapterSection]
    word_count: int
    key_terms: dict[str, str]
    summary: str
    generation_metadata: dict


class ChapterValidationError(Exception):
    """Raised when chapter validation fails."""
    pass


class OllamaConnectionError(Exception):
    """Raised when Ollama connection fails."""
    pass


class ChapterGenerator:
    """
    Main chapter generation system using Ollama.

    Features:
    - Streaming output for progress tracking
    - Retry logic with exponential backoff
    - Structure validation
    - Consistency maintenance with previous chapters
    - Markdown formatting
    """

    DEFAULT_MODEL = "llama3.2"
    DEFAULT_OLLAMA_URL = "http://localhost:11434"

    def __init__(
        self,
        ollama_url: str = DEFAULT_OLLAMA_URL,
        model: str = DEFAULT_MODEL,
        max_retries: int = 3,
        retry_delay: float = 2.0,
        timeout: int = 600
    ):
        self.ollama_url = ollama_url.rstrip('/')
        self.model = model
        self.max_retries = max_retries
        self.retry_delay = retry_delay
        self.timeout = timeout

        # Validate connection on init
        self._validate_connection()

    def _validate_connection(self) -> None:
        """Validate Ollama is accessible and model is available."""
        try:
            response = requests.get(f"{self.ollama_url}/api/tags", timeout=10)
            response.raise_for_status()

            models = response.json().get('models', [])
            model_names = [m['name'].split(':')[0] for m in models]

            if self.model.split(':')[0] not in model_names:
                logger.warning(f"Model '{self.model}' not found. Available: {model_names}")
                logger.info(f"Will attempt to pull model '{self.model}' on first use.")

        except requests.exceptions.RequestException as e:
            raise OllamaConnectionError(f"Cannot connect to Ollama at {self.ollama_url}: {e}")

    def _build_generation_prompt(
        self,
        outline: ChapterOutline,
        style_guide: StyleGuide,
        previous_summary: Optional[PreviousChaptersSummary]
    ) -> str:
        """Build the complete prompt for chapter generation."""

        prev_context = previous_summary.to_prompt_context() if previous_summary else "This is the first chapter."

        structure_instructions = self._get_structure_instructions(style_guide)

        prompt = f"""You are an expert author tasked with writing a complete book chapter.
Generate high-quality, engaging content that follows the outline precisely while maintaining consistency with previous chapters.

=== STYLE GUIDE ===
{style_guide.to_prompt_context()}

=== CHAPTER OUTLINE ===
{outline.to_prompt_context()}

=== PREVIOUS CHAPTERS CONTEXT ===
{prev_context}

=== STRUCTURE REQUIREMENTS ===
{structure_instructions}

=== FORMATTING REQUIREMENTS ===
Use proper Markdown formatting:
- Use ## for main section headers
- Use ### for subsections
- Use **bold** for emphasis and key terms
- Use *italic* for definitions or secondary emphasis
- Use > for callout boxes and important notes
- Use - or * for bullet lists
- Use 1. 2. 3. for numbered lists
- Use ``` for code blocks (if technical content)
- Use --- for section breaks when appropriate

=== OUTPUT INSTRUCTIONS ===
Generate the COMPLETE chapter now. The chapter should be {style_guide.target_word_count_min}-{style_guide.target_word_count_max} words.
Start with the chapter title as a level 1 heading (# Chapter {outline.chapter_number}: {outline.title}).

Begin writing the chapter:
"""

        return prompt

    def _get_structure_instructions(self, style_guide: StyleGuide) -> str:
        """Get structure instructions based on book type."""

        base_structure = """
The chapter MUST include these sections in order:

1. **Opening Hook** (1-2 paragraphs)
   - Start with an engaging hook: story, question, surprising fact, or scenario
   - Draw the reader in immediately
   - Connect to the chapter's main topic

2. **Introduction** (2-3 paragraphs)
   - Clearly state what this chapter will cover
   - Explain why this topic matters
   - Preview the main sections
"""

        if style_guide.book_type == BookType.NON_FICTION:
            base_structure += """
3. **Main Content Sections** (bulk of chapter)
   - Divide into 3-5 major sections with clear headings
   - Each section should have:
     * Clear explanation of concepts
     * Supporting evidence or research
     * Practical applications
   - Use transitions between sections

4. **Examples and Case Studies**
   - Include 2-3 detailed examples
   - Make examples relevant and relatable
   - Show real-world applications

5. **Callout Boxes** (throughout)
   - Use blockquotes (>) for:
     * Key insights
     * Pro tips
     * Common mistakes to avoid
     * Quick summaries

6. **Chapter Summary** (2-3 paragraphs)
   - Recap main points
   - Reinforce key takeaways
   - Connect to upcoming chapters if applicable

7. **Key Terms** (at end)
   - List and define 5-10 important terms from the chapter
   - Format as a definition list or glossary
"""

        elif style_guide.book_type == BookType.FICTION:
            base_structure += """
3. **Scene Development** (bulk of chapter)
   - Build scenes with vivid sensory details
   - Develop character through action and dialogue
   - Maintain consistent pacing
   - Include internal character thoughts where appropriate

4. **Conflict and Tension**
   - Advance the central conflict
   - Include smaller conflicts within scenes
   - Build toward chapter climax

5. **Chapter Ending**
   - End with hook or cliffhanger when appropriate
   - Or provide satisfying scene resolution
   - Set up next chapter's events
"""

        elif style_guide.book_type == BookType.TECHNICAL:
            base_structure += """
3. **Concept Explanation** (bulk of chapter)
   - Start with fundamentals
   - Build complexity gradually
   - Use analogies to explain difficult concepts
   - Include diagrams descriptions where helpful

4. **Code Examples / Technical Demonstrations**
   - Include working code snippets
   - Explain each part of the code
   - Show both simple and advanced usage

5. **Best Practices and Pitfalls**
   - Highlight recommended approaches
   - Warn about common mistakes
   - Provide debugging tips

6. **Hands-On Exercises**
   - Include 2-3 practice exercises
   - Provide hints or partial solutions
   - Scale difficulty progressively

7. **Chapter Summary**
   - Bullet-point key takeaways
   - Quick reference card

8. **Key Terms / API Reference**
   - Define technical terms
   - Include method/function signatures if applicable
"""

        elif style_guide.book_type == BookType.SELF_HELP:
            base_structure += """
3. **Core Teaching** (bulk of chapter)
   - Present the main concept or principle
   - Support with research or evidence
   - Include personal stories or client examples
   - Make it relatable to reader's life

4. **Actionable Steps**
   - Provide clear, numbered action items
   - Make steps specific and achievable
   - Include timeframes where appropriate

5. **Reflection Questions**
   - Include 3-5 thought-provoking questions
   - Help readers apply concepts to their lives
   - Encourage journaling or self-assessment

6. **Common Obstacles**
   - Address likely challenges
   - Provide solutions for each obstacle
   - Normalize difficulties

7. **Chapter Summary**
   - Recap main insights
   - List action items
   - Motivational closing
"""

        return base_structure

    def _stream_generate(
        self,
        prompt: str,
        on_token: Optional[Callable[[str], None]] = None
    ) -> Generator[str, None, None]:
        """Stream tokens from Ollama generation."""

        payload = {
            "model": self.model,
            "prompt": prompt,
            "stream": True,
            "options": {
                "num_predict": 16000,  # Allow long generation
                "temperature": 0.7,
                "top_p": 0.9,
            }
        }

        try:
            response = requests.post(
                f"{self.ollama_url}/api/generate",
                json=payload,
                stream=True,
                timeout=self.timeout
            )
            response.raise_for_status()

            for line in response.iter_lines():
                if line:
                    data = json.loads(line)
                    token = data.get('response', '')
                    if token:
                        if on_token:
                            on_token(token)
                        yield token

                    if data.get('done', False):
                        break

        except requests.exceptions.RequestException as e:
            raise OllamaConnectionError(f"Generation failed: {e}")

    def _count_words(self, text: str) -> int:
        """Count words in text."""
        return len(text.split())

    def _extract_sections(self, content: str) -> list[ChapterSection]:
        """Extract sections from generated content."""
        sections = []

        # Split by level 2 headers
        pattern = r'^## (.+?)$'
        parts = re.split(pattern, content, flags=re.MULTILINE)

        if parts[0].strip():  # Content before first ##
            sections.append(ChapterSection(
                section_type="opening",
                content=parts[0].strip(),
                word_count=self._count_words(parts[0])
            ))

        # Process paired header/content
        for i in range(1, len(parts), 2):
            if i + 1 < len(parts):
                header = parts[i].strip()
                section_content = parts[i + 1].strip()
                sections.append(ChapterSection(
                    section_type=header.lower().replace(' ', '_'),
                    content=f"## {header}\n\n{section_content}",
                    word_count=self._count_words(section_content)
                ))

        return sections

    def _extract_key_terms(self, content: str, style_guide: StyleGuide) -> dict[str, str]:
        """Extract key terms and definitions from content."""
        key_terms = {}

        if not style_guide.include_key_terms:
            return key_terms

        # Look for Key Terms section
        key_terms_match = re.search(
            r'## Key Terms.*?$(.*?)(?=^## |\Z)',
            content,
            re.MULTILINE | re.DOTALL
        )

        if key_terms_match:
            terms_section = key_terms_match.group(1)

            # Try different formats
            # Format: **term**: definition
            pattern1 = r'\*\*([^*]+)\*\*[:\s]+([^\n]+)'
            # Format: - term: definition
            pattern2 = r'^[\-\*]\s*([^:]+):\s*(.+)$'

            for match in re.finditer(pattern1, terms_section):
                key_terms[match.group(1).strip()] = match.group(2).strip()

            for match in re.finditer(pattern2, terms_section, re.MULTILINE):
                term = match.group(1).strip()
                if term not in key_terms:
                    key_terms[term] = match.group(2).strip()

        return key_terms

    def _extract_summary(self, content: str) -> str:
        """Extract the summary section from content."""
        # Look for Summary or Conclusion section
        summary_match = re.search(
            r'## (?:Summary|Conclusion|Chapter Summary).*?$(.*?)(?=^## Key Terms|^## |\Z)',
            content,
            re.MULTILINE | re.DOTALL | re.IGNORECASE
        )

        if summary_match:
            return summary_match.group(1).strip()

        # Fallback: last few paragraphs
        paragraphs = [p.strip() for p in content.split('\n\n') if p.strip()]
        return '\n\n'.join(paragraphs[-3:]) if paragraphs else ""

    def _validate_chapter(
        self,
        content: str,
        style_guide: StyleGuide,
        outline: ChapterOutline
    ) -> tuple[bool, list[str]]:
        """
        Validate the generated chapter meets requirements.

        Returns: (is_valid, list of issues)
        """
        issues = []
        word_count = self._count_words(content)

        # Check word count
        if word_count < style_guide.target_word_count_min * 0.7:
            issues.append(
                f"Word count ({word_count}) is significantly below minimum "
                f"({style_guide.target_word_count_min}). Needs expansion."
            )
        elif word_count < style_guide.target_word_count_min:
            issues.append(
                f"Word count ({word_count}) is slightly below minimum "
                f"({style_guide.target_word_count_min})."
            )

        # Check for required sections
        if style_guide.book_type != BookType.FICTION:
            if style_guide.include_key_terms and 'key term' not in content.lower():
                issues.append("Missing Key Terms section.")

        # Check for markdown formatting
        if '##' not in content:
            issues.append("Missing section headers (##).")

        if style_guide.include_callouts and '>' not in content:
            issues.append("Missing callout boxes (>).")

        # Check chapter title is present
        expected_title = f"# Chapter {outline.chapter_number}"
        if expected_title.lower() not in content.lower()[:500]:
            issues.append("Chapter title not found at beginning.")

        # Check for subtopics coverage
        content_lower = content.lower()
        missing_topics = []
        for subtopic in outline.subtopics:
            # Check if any significant words from subtopic appear
            words = [w for w in subtopic.lower().split() if len(w) > 4]
            if words and not any(w in content_lower for w in words):
                missing_topics.append(subtopic)

        if missing_topics:
            issues.append(f"Subtopics may not be fully covered: {missing_topics[:3]}")

        is_valid = len([i for i in issues if 'significantly' in i or 'Missing' in i]) == 0

        return is_valid, issues

    def _expand_chapter(
        self,
        content: str,
        issues: list[str],
        outline: ChapterOutline,
        style_guide: StyleGuide,
        on_token: Optional[Callable[[str], None]] = None
    ) -> str:
        """Expand a chapter that's too short or missing sections."""

        issues_str = "\n".join(f"- {issue}" for issue in issues)

        expansion_prompt = f"""The following chapter draft needs expansion and improvement.

CURRENT DRAFT:
{content}

ISSUES TO ADDRESS:
{issues_str}

REQUIREMENTS:
- Expand the content to reach {style_guide.target_word_count_min}-{style_guide.target_word_count_max} words
- Add missing sections if noted
- Maintain the same style and tone
- Keep all existing good content
- Add more examples, details, and explanations
- Use proper Markdown formatting

Generate the EXPANDED and IMPROVED chapter:
"""

        expanded_content = ""
        for token in self._stream_generate(expansion_prompt, on_token):
            expanded_content += token

        return expanded_content

    def generate_chapter(
        self,
        outline: ChapterOutline,
        style_guide: Optional[StyleGuide] = None,
        previous_summary: Optional[PreviousChaptersSummary] = None,
        on_token: Optional[Callable[[str], None]] = None,
        on_progress: Optional[Callable[[str, float], None]] = None
    ) -> GeneratedChapter:
        """
        Generate a complete chapter.

        Args:
            outline: Chapter outline with structure and key points
            style_guide: Styling configuration (uses defaults if None)
            previous_summary: Summary of previous chapters for consistency
            on_token: Callback for each generated token (for streaming display)
            on_progress: Callback for progress updates (stage, percentage)

        Returns:
            GeneratedChapter with full content and metadata

        Raises:
            ChapterValidationError: If chapter fails validation after retries
            OllamaConnectionError: If Ollama communication fails
        """
        if style_guide is None:
            style_guide = StyleGuide()

        if on_progress:
            on_progress("Building prompt", 0.05)

        prompt = self._build_generation_prompt(outline, style_guide, previous_summary)

        start_time = time.time()
        content = ""
        attempt = 0

        while attempt < self.max_retries:
            attempt += 1

            if on_progress:
                on_progress(f"Generating (attempt {attempt}/{self.max_retries})", 0.1)

            logger.info(f"Generation attempt {attempt}/{self.max_retries}")

            try:
                content = ""
                token_count = 0

                for token in self._stream_generate(prompt, on_token):
                    content += token
                    token_count += 1

                    # Progress updates during generation
                    if on_progress and token_count % 100 == 0:
                        # Estimate progress based on target word count
                        current_words = self._count_words(content)
                        progress = min(0.8, 0.1 + (current_words / style_guide.target_word_count_max) * 0.7)
                        on_progress(f"Generating ({current_words} words)", progress)

                if on_progress:
                    on_progress("Validating structure", 0.85)

                # Validate the generated content
                is_valid, issues = self._validate_chapter(content, style_guide, outline)

                if is_valid:
                    logger.info(f"Chapter validated successfully on attempt {attempt}")
                    if issues:
                        logger.info(f"Minor issues (non-blocking): {issues}")
                    break
                else:
                    logger.warning(f"Validation issues on attempt {attempt}: {issues}")

                    # Try to expand if too short
                    if any('word count' in issue.lower() for issue in issues):
                        if on_progress:
                            on_progress("Expanding content", 0.7)
                        content = self._expand_chapter(
                            content, issues, outline, style_guide, on_token
                        )
                        is_valid, issues = self._validate_chapter(content, style_guide, outline)
                        if is_valid:
                            break

                    if attempt < self.max_retries:
                        delay = self.retry_delay * (2 ** (attempt - 1))
                        logger.info(f"Retrying in {delay}s...")
                        time.sleep(delay)

            except OllamaConnectionError as e:
                logger.error(f"Connection error on attempt {attempt}: {e}")
                if attempt < self.max_retries:
                    delay = self.retry_delay * (2 ** (attempt - 1))
                    time.sleep(delay)
                else:
                    raise

        # Final validation
        is_valid, issues = self._validate_chapter(content, style_guide, outline)
        if not is_valid:
            logger.warning(f"Chapter generated with issues after {attempt} attempts: {issues}")

        if on_progress:
            on_progress("Extracting metadata", 0.9)

        # Extract components
        sections = self._extract_sections(content)
        key_terms = self._extract_key_terms(content, style_guide)
        summary = self._extract_summary(content)
        word_count = self._count_words(content)

        generation_time = time.time() - start_time

        if on_progress:
            on_progress("Complete", 1.0)

        return GeneratedChapter(
            chapter_number=outline.chapter_number,
            title=outline.title,
            full_content=content,
            sections=sections,
            word_count=word_count,
            key_terms=key_terms,
            summary=summary,
            generation_metadata={
                "model": self.model,
                "generation_time_seconds": round(generation_time, 2),
                "attempts": attempt,
                "validation_issues": issues,
                "style_guide": {
                    "book_type": style_guide.book_type.value,
                    "target_word_count": f"{style_guide.target_word_count_min}-{style_guide.target_word_count_max}"
                }
            }
        )

    def generate_chapter_stream(
        self,
        outline: ChapterOutline,
        style_guide: Optional[StyleGuide] = None,
        previous_summary: Optional[PreviousChaptersSummary] = None
    ) -> Generator[tuple[str, str], None, GeneratedChapter]:
        """
        Generate chapter with streaming output.

        Yields:
            Tuples of (event_type, data) where event_type is:
            - 'token': Individual token
            - 'progress': Progress update JSON
            - 'complete': Final chapter JSON

        Returns:
            GeneratedChapter when complete
        """
        tokens = []

        def token_collector(token: str):
            tokens.append(token)

        def progress_callback(stage: str, progress: float):
            pass  # Will yield progress separately

        # Generate with collection
        chapter = self.generate_chapter(
            outline=outline,
            style_guide=style_guide,
            previous_summary=previous_summary,
            on_token=token_collector,
            on_progress=progress_callback
        )

        return chapter


def create_sample_outline() -> ChapterOutline:
    """Create a sample outline for testing."""
    return ChapterOutline(
        chapter_number=1,
        title="Introduction to Machine Learning",
        main_topic="Fundamental concepts of machine learning and its applications",
        subtopics=[
            "What is Machine Learning?",
            "Types of Machine Learning: Supervised, Unsupervised, Reinforcement",
            "Real-world Applications",
            "The Machine Learning Pipeline",
            "Getting Started with ML Tools"
        ],
        key_points=[
            "ML enables computers to learn from data without explicit programming",
            "Different types of ML suit different problem types",
            "Data quality is crucial for ML success",
            "ML is transforming industries from healthcare to finance"
        ],
        learning_objectives=[
            "Define machine learning and explain its significance",
            "Distinguish between supervised, unsupervised, and reinforcement learning",
            "Identify appropriate ML approaches for given problems",
            "Understand the basic ML workflow"
        ],
        connections_to_other_chapters=[
            "Chapter 2 will dive deeper into supervised learning",
            "Chapter 5 covers the data preparation mentioned here"
        ]
    )


def main():
    """Example usage of the Chapter Generator."""
    import sys

    print("=" * 60)
    print("BookForge Chapter Generator")
    print("=" * 60)

    # Initialize generator
    try:
        generator = ChapterGenerator(
            model="llama3.2",
            max_retries=3
        )
        print(f"Connected to Ollama with model: {generator.model}")
    except OllamaConnectionError as e:
        print(f"Error: {e}")
        print("Make sure Ollama is running: ollama serve")
        sys.exit(1)

    # Create sample inputs
    outline = create_sample_outline()

    style_guide = StyleGuide(
        book_type=BookType.TECHNICAL,
        target_word_count_min=3000,  # Shorter for demo
        target_word_count_max=4000,
        tone="friendly and educational",
        include_examples=True,
        include_callouts=True,
        include_key_terms=True
    )

    # For first chapter, no previous summary
    previous_summary = None

    print(f"\nGenerating: Chapter {outline.chapter_number}: {outline.title}")
    print(f"Target: {style_guide.target_word_count_min}-{style_guide.target_word_count_max} words")
    print("-" * 60)

    # Progress callback
    def show_progress(stage: str, progress: float):
        bar_length = 30
        filled = int(bar_length * progress)
        bar = '#' * filled + '-' * (bar_length - filled)
        print(f"\r[{bar}] {progress*100:.0f}% - {stage}", end='', flush=True)

    # Token callback for streaming
    def show_token(token: str):
        print(token, end='', flush=True)

    try:
        print("\n")  # New line before content

        chapter = generator.generate_chapter(
            outline=outline,
            style_guide=style_guide,
            previous_summary=previous_summary,
            on_token=show_token,
            on_progress=show_progress
        )

        print("\n")
        print("=" * 60)
        print("GENERATION COMPLETE")
        print("=" * 60)
        print(f"Chapter: {chapter.chapter_number} - {chapter.title}")
        print(f"Word Count: {chapter.word_count}")
        print(f"Sections: {len(chapter.sections)}")
        print(f"Key Terms: {len(chapter.key_terms)}")
        print(f"Generation Time: {chapter.generation_metadata['generation_time_seconds']}s")
        print(f"Attempts: {chapter.generation_metadata['attempts']}")

        if chapter.generation_metadata['validation_issues']:
            print(f"Notes: {chapter.generation_metadata['validation_issues']}")

        # Save to file
        output_file = f"chapter_{chapter.chapter_number}.md"
        with open(output_file, 'w') as f:
            f.write(chapter.full_content)
        print(f"\nSaved to: {output_file}")

    except ChapterValidationError as e:
        print(f"\nValidation Error: {e}")
        sys.exit(1)
    except OllamaConnectionError as e:
        print(f"\nConnection Error: {e}")
        sys.exit(1)
    except KeyboardInterrupt:
        print("\n\nGeneration cancelled.")
        sys.exit(0)


if __name__ == "__main__":
    main()
