"""
Book Outline Generator using Ollama (mistral-nemo)
Generates structured book outlines for various genres and book types.
"""

import json
import re
import time
from dataclasses import dataclass, field, asdict
from enum import Enum
from typing import Optional
import requests


class BookType(Enum):
    NOVEL = "novel"
    TEXTBOOK = "textbook"
    SELF_HELP = "self_help"
    MEMOIR = "memoir"
    BUSINESS = "business"
    TECHNICAL = "technical"
    CHILDREN = "children"
    MYSTERY = "mystery"
    FANTASY = "fantasy"
    SCIENCE_FICTION = "science_fiction"
    ROMANCE = "romance"
    THRILLER = "thriller"
    BIOGRAPHY = "biography"
    COOKBOOK = "cookbook"
    TRAVEL = "travel"


@dataclass
class ChapterOutline:
    """Represents a single chapter in the book outline."""
    chapter_number: int
    title: str
    key_topics: list[str] = field(default_factory=list)
    summary: str = ""
    learning_objectives: list[str] = field(default_factory=list)  # For non-fiction
    character_developments: list[str] = field(default_factory=list)  # For fiction
    plot_points: list[str] = field(default_factory=list)  # For fiction
    estimated_word_count: int = 0


@dataclass
class PartDivision:
    """Represents a part/section division in the book."""
    part_number: int
    title: str
    description: str = ""
    chapters: list[ChapterOutline] = field(default_factory=list)


@dataclass
class CharacterArc:
    """Represents a character arc for fiction books."""
    character_name: str
    role: str  # protagonist, antagonist, supporting, etc.
    arc_description: str
    key_moments: list[str] = field(default_factory=list)
    starting_state: str = ""
    ending_state: str = ""


@dataclass
class BookOutline:
    """Complete book outline structure."""
    title: str
    subtitle: str = ""
    genre: str = ""
    book_type: str = ""
    target_audience: str = ""
    description: str = ""
    premise: str = ""
    themes: list[str] = field(default_factory=list)
    parts: list[PartDivision] = field(default_factory=list)
    chapters: list[ChapterOutline] = field(default_factory=list)  # If no parts
    character_arcs: list[CharacterArc] = field(default_factory=list)  # For fiction
    learning_outcomes: list[str] = field(default_factory=list)  # For non-fiction
    total_estimated_words: int = 0
    generation_metadata: dict = field(default_factory=dict)

    def to_dict(self) -> dict:
        """Convert to dictionary for database storage."""
        return asdict(self)


class OutlineGenerationError(Exception):
    """Raised when outline generation fails."""
    pass


class OutlineGenerator:
    """
    Generates book outlines using Ollama with mistral-nemo model.
    """
    
    FICTION_TYPES = {
        BookType.NOVEL, BookType.MYSTERY, BookType.FANTASY,
        BookType.SCIENCE_FICTION, BookType.ROMANCE, BookType.THRILLER,
        BookType.CHILDREN
    }
    
    NON_FICTION_TYPES = {
        BookType.TEXTBOOK, BookType.SELF_HELP, BookType.BUSINESS,
        BookType.TECHNICAL, BookType.BIOGRAPHY, BookType.MEMOIR,
        BookType.COOKBOOK, BookType.TRAVEL
    }

    def __init__(
        self,
        ollama_host: str = "http://localhost:11434",
        model: str = "mistral-nemo",
        max_retries: int = 3,
        retry_delay: float = 2.0,
        timeout: int = 120
    ):
        self.ollama_host = ollama_host.rstrip("/")
        self.model = model
        self.max_retries = max_retries
        self.retry_delay = retry_delay
        self.timeout = timeout
    
    def _is_fiction(self, book_type: BookType) -> bool:
        """Determine if the book type is fiction."""
        return book_type in self.FICTION_TYPES
    
    def _build_prompt(
        self,
        title: str,
        genre: str,
        target_audience: str,
        num_chapters: int,
        description: str,
        book_type: BookType
    ) -> str:
        """Build the prompt for outline generation."""
        is_fiction = self._is_fiction(book_type)
        
        base_prompt = f"""You are an expert book outline creator. Generate a detailed, structured outline for a book with the following specifications:

BOOK DETAILS:
- Title: {title}
- Genre: {genre}
- Book Type: {book_type.value}
- Target Audience: {target_audience}
- Number of Chapters: {num_chapters}
- Description: {description}

Generate a comprehensive JSON outline with the following structure:
{{
    "title": "{title}",
    "subtitle": "<generate an appropriate subtitle>",
    "premise": "<one paragraph premise/hook>",
    "themes": ["<theme1>", "<theme2>", ...],
    "parts": [
        {{
            "part_number": 1,
            "title": "<part title>",
            "description": "<what this part covers>",
            "chapters": [
                {{
                    "chapter_number": 1,
                    "title": "<chapter title>",
                    "key_topics": ["<topic1>", "<topic2>", ...],
                    "summary": "<2-3 sentence summary of chapter content>",
                    "estimated_word_count": <number>
"""
        
        if is_fiction:
            base_prompt += """,
                    "plot_points": ["<key plot point 1>", "<key plot point 2>", ...],
                    "character_developments": ["<character development note>", ...]
"""
        else:
            base_prompt += """,
                    "learning_objectives": ["<objective 1>", "<objective 2>", ...]
"""
        
        base_prompt += """
                }
            ]
        }
    ],
"""
        
        if is_fiction:
            base_prompt += """
    "character_arcs": [
        {
            "character_name": "<name>",
            "role": "<protagonist/antagonist/supporting>",
            "arc_description": "<description of character's journey>",
            "key_moments": ["<moment1>", "<moment2>", ...],
            "starting_state": "<character's initial state>",
            "ending_state": "<character's final state>"
        }
    ],
"""
        else:
            base_prompt += """
    "learning_outcomes": [
        "<overall learning outcome 1>",
        "<overall learning outcome 2>",
        ...
    ],
"""
        
        base_prompt += f"""
    "total_estimated_words": <total word count estimate>
}}

IMPORTANT GUIDELINES:
1. Divide the {num_chapters} chapters logically into 2-4 parts/sections
2. Each chapter should have 3-5 key topics
3. Make chapter titles engaging and descriptive
4. Ensure logical flow and progression between chapters
5. Word count estimates should be realistic (average 3000-5000 words per chapter)
"""
        
        if is_fiction:
            base_prompt += """
6. Include at least 2-3 character arcs for main characters
7. Plot points should create tension and drive the narrative
8. Character developments should show growth/change
"""
        else:
            base_prompt += """
6. Learning objectives should be specific and measurable
7. Each chapter should build on previous knowledge
8. Include practical, actionable content
"""
        
        base_prompt += """

Respond ONLY with valid JSON. No additional text or explanation.
"""
        
        return base_prompt
    
    def _call_ollama(self, prompt: str) -> str:
        """Make a request to Ollama API."""
        url = f"{self.ollama_host}/api/generate"
        
        payload = {
            "model": self.model,
            "prompt": prompt,
            "stream": False,
            "options": {
                "temperature": 0.7,
                "num_predict": 4096
            }
        }
        
        response = requests.post(
            url,
            json=payload,
            timeout=self.timeout
        )
        response.raise_for_status()
        
        result = response.json()
        return result.get("response", "")
    
    def _extract_json(self, text: str) -> dict:
        """Extract and parse JSON from the response text."""
        # Try to find JSON in the response
        text = text.strip()
        
        # Try direct parse first
        try:
            return json.loads(text)
        except json.JSONDecodeError:
            pass
        
        # Try to find JSON block in markdown
        json_match = re.search(r"```(?:json)?\s*({.*?})\s*```", text, re.DOTALL)
        if json_match:
            try:
                return json.loads(json_match.group(1))
            except json.JSONDecodeError:
                pass
        
        # Try to find raw JSON object
        json_match = re.search(r"({.*})", text, re.DOTALL)
        if json_match:
            try:
                return json.loads(json_match.group(1))
            except json.JSONDecodeError:
                pass
        
        # Try to fix common JSON issues
        fixed_text = text
        # Remove trailing commas before closing brackets
        fixed_text = re.sub(r",\s*}", "}", fixed_text)
        fixed_text = re.sub(r",\s*]", "]", fixed_text)
        
        try:
            return json.loads(fixed_text)
        except json.JSONDecodeError as e:
            raise OutlineGenerationError(f"Failed to parse JSON: {e}")
    
    def _validate_outline(self, data: dict, num_chapters: int) -> bool:
        """Validate the generated outline structure."""
        required_fields = ["title", "parts"]
        
        for field in required_fields:
            if field not in data:
                return False
        
        # Count total chapters
        total_chapters = 0
        for part in data.get("parts", []):
            if "chapters" not in part:
                return False
            total_chapters += len(part["chapters"])
        
        # Allow some flexibility in chapter count
        if total_chapters < num_chapters - 2 or total_chapters > num_chapters + 2:
            return False
        
        return True
    
    def _parse_outline(
        self,
        data: dict,
        book_type: BookType,
        genre: str,
        target_audience: str,
        description: str
    ) -> BookOutline:
        """Parse the JSON data into a BookOutline object."""
        is_fiction = self._is_fiction(book_type)
        
        # Parse parts and chapters
        parts = []
        all_chapters = []
        
        for part_data in data.get("parts", []):
            chapters = []
            for ch_data in part_data.get("chapters", []):
                chapter = ChapterOutline(
                    chapter_number=ch_data.get("chapter_number", len(all_chapters) + 1),
                    title=ch_data.get("title", f"Chapter {len(all_chapters) + 1}"),
                    key_topics=ch_data.get("key_topics", []),
                    summary=ch_data.get("summary", ""),
                    learning_objectives=ch_data.get("learning_objectives", []) if not is_fiction else [],
                    character_developments=ch_data.get("character_developments", []) if is_fiction else [],
                    plot_points=ch_data.get("plot_points", []) if is_fiction else [],
                    estimated_word_count=ch_data.get("estimated_word_count", 3500)
                )
                chapters.append(chapter)
                all_chapters.append(chapter)
            
            part = PartDivision(
                part_number=part_data.get("part_number", len(parts) + 1),
                title=part_data.get("title", f"Part {len(parts) + 1}"),
                description=part_data.get("description", ""),
                chapters=chapters
            )
            parts.append(part)
        
        # Parse character arcs for fiction
        character_arcs = []
        if is_fiction:
            for arc_data in data.get("character_arcs", []):
                arc = CharacterArc(
                    character_name=arc_data.get("character_name", "Unknown"),
                    role=arc_data.get("role", "supporting"),
                    arc_description=arc_data.get("arc_description", ""),
                    key_moments=arc_data.get("key_moments", []),
                    starting_state=arc_data.get("starting_state", ""),
                    ending_state=arc_data.get("ending_state", "")
                )
                character_arcs.append(arc)
        
        # Calculate total word count
        total_words = sum(ch.estimated_word_count for ch in all_chapters)
        if total_words == 0:
            total_words = data.get("total_estimated_words", len(all_chapters) * 3500)
        
        outline = BookOutline(
            title=data.get("title", ""),
            subtitle=data.get("subtitle", ""),
            genre=genre,
            book_type=book_type.value,
            target_audience=target_audience,
            description=description,
            premise=data.get("premise", ""),
            themes=data.get("themes", []),
            parts=parts,
            chapters=all_chapters,
            character_arcs=character_arcs,
            learning_outcomes=data.get("learning_outcomes", []) if not is_fiction else [],
            total_estimated_words=total_words,
            generation_metadata={
                "model": self.model,
                "generated_at": time.strftime("%Y-%m-%d %H:%M:%S"),
                "is_fiction": is_fiction
            }
        )
        
        return outline
    
    def generate(
        self,
        title: str,
        genre: str,
        target_audience: str,
        num_chapters: int,
        description: str,
        book_type: Optional[str] = None
    ) -> dict:
        """
        Generate a book outline.
        
        Args:
            title: The book title
            genre: The book genre (e.g., "Fantasy", "Self-Help")
            target_audience: Target reader description
            num_chapters: Desired number of chapters
            description: Brief description of the book
            book_type: Type of book (from BookType enum values)
        
        Returns:
            dict: Structured outline ready for database storage
        
        Raises:
            OutlineGenerationError: If generation fails after retries
        """
        # Determine book type
        if book_type:
            try:
                bt = BookType(book_type.lower().replace("-", "_").replace(" ", "_"))
            except ValueError:
                # Try to infer from genre
                bt = self._infer_book_type(genre)
        else:
            bt = self._infer_book_type(genre)
        
        prompt = self._build_prompt(
            title=title,
            genre=genre,
            target_audience=target_audience,
            num_chapters=num_chapters,
            description=description,
            book_type=bt
        )
        
        last_error = None
        
        for attempt in range(self.max_retries):
            try:
                # Call Ollama
                response_text = self._call_ollama(prompt)
                
                # Extract and parse JSON
                data = self._extract_json(response_text)
                
                # Validate structure
                if not self._validate_outline(data, num_chapters):
                    raise OutlineGenerationError(
                        f"Generated outline structure is invalid or chapter count mismatch"
                    )
                
                # Parse into structured format
                outline = self._parse_outline(
                    data=data,
                    book_type=bt,
                    genre=genre,
                    target_audience=target_audience,
                    description=description
                )
                
                return outline.to_dict()
                
            except requests.RequestException as e:
                last_error = OutlineGenerationError(f"Ollama API error: {e}")
            except OutlineGenerationError as e:
                last_error = e
            except Exception as e:
                last_error = OutlineGenerationError(f"Unexpected error: {e}")
            
            if attempt < self.max_retries - 1:
                time.sleep(self.retry_delay * (attempt + 1))
        
        raise OutlineGenerationError(
            f"Failed to generate outline after {self.max_retries} attempts. "
            f"Last error: {last_error}"
        )
    
    def _infer_book_type(self, genre: str) -> BookType:
        """Infer book type from genre string."""
        genre_lower = genre.lower()
        
        type_mapping = {
            "fantasy": BookType.FANTASY,
            "science fiction": BookType.SCIENCE_FICTION,
            "sci-fi": BookType.SCIENCE_FICTION,
            "mystery": BookType.MYSTERY,
            "thriller": BookType.THRILLER,
            "romance": BookType.ROMANCE,
            "self-help": BookType.SELF_HELP,
            "self help": BookType.SELF_HELP,
            "business": BookType.BUSINESS,
            "technical": BookType.TECHNICAL,
            "programming": BookType.TECHNICAL,
            "textbook": BookType.TEXTBOOK,
            "academic": BookType.TEXTBOOK,
            "biography": BookType.BIOGRAPHY,
            "memoir": BookType.MEMOIR,
            "children": BookType.CHILDREN,
            "cookbook": BookType.COOKBOOK,
            "cooking": BookType.COOKBOOK,
            "travel": BookType.TRAVEL,
        }
        
        for key, book_type in type_mapping.items():
            if key in genre_lower:
                return book_type
        
        # Default to novel for unrecognized fiction-like genres
        fiction_keywords = ["fiction", "novel", "story", "adventure", "horror"]
        if any(kw in genre_lower for kw in fiction_keywords):
            return BookType.NOVEL
        
        return BookType.SELF_HELP  # Default for non-fiction


def generate_outline(
    title: str,
    genre: str,
    target_audience: str,
    num_chapters: int,
    description: str,
    book_type: Optional[str] = None,
    ollama_host: str = "http://localhost:11434",
    model: str = "mistral-nemo"
) -> dict:
    """
    Convenience function to generate a book outline.
    
    Args:
        title: The book title
        genre: The book genre
        target_audience: Target reader description
        num_chapters: Desired number of chapters
        description: Brief description of the book
        book_type: Optional book type override
        ollama_host: Ollama server URL
        model: Model name to use
    
    Returns:
        dict: Structured outline for database storage
    """
    generator = OutlineGenerator(ollama_host=ollama_host, model=model)
    return generator.generate(
        title=title,
        genre=genre,
        target_audience=target_audience,
        num_chapters=num_chapters,
        description=description,
        book_type=book_type
    )


# Example usage
if __name__ == "__main__":
    # Example: Generate a fantasy novel outline
    try:
        outline = generate_outline(
            title="The Crystal Kingdoms",
            genre="Fantasy",
            target_audience="Young adults aged 16-25 who enjoy epic fantasy",
            num_chapters=12,
            description="A young mage discovers she is the last heir to a fallen kingdom and must unite warring factions to defeat an ancient evil awakening in the north.",
            book_type="fantasy"
        )
        
        print("Generated Outline:")
        print(json.dumps(outline, indent=2))
        
    except OutlineGenerationError as e:
        print(f"Error: {e}")
    
    # Example: Generate a self-help book outline
    try:
        outline = generate_outline(
            title="The Productivity Mindset",
            genre="Self-Help",
            target_audience="Working professionals aged 25-45 seeking work-life balance",
            num_chapters=10,
            description="A comprehensive guide to maximizing productivity while maintaining mental health and personal relationships.",
            book_type="self_help"
        )
        
        print("\nGenerated Self-Help Outline:")
        print(json.dumps(outline, indent=2))
        
    except OutlineGenerationError as e:
        print(f"Error: {e}")
