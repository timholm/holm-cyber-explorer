"""
Ollama Client for Book Generation

A production-ready async client for generating book content using Ollama LLMs.
Supports outline generation, chapter writing, and summarization with streaming,
retry logic, and comprehensive error handling.
"""

import asyncio
import aiohttp
import json
import logging
from dataclasses import dataclass, field
from typing import AsyncIterator, Optional, List, Dict, Any, Callable
from enum import Enum
import time
from functools import wraps

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


class OllamaError(Exception):
    """Base exception for Ollama client errors."""
    pass


class ConnectionError(OllamaError):
    """Raised when connection to Ollama fails."""
    pass


class ModelNotFoundError(OllamaError):
    """Raised when requested model is not available."""
    pass


class GenerationError(OllamaError):
    """Raised when text generation fails."""
    pass


class TimeoutError(OllamaError):
    """Raised when request times out."""
    pass


class ModelPreference(Enum):
    """Preferred models for generation."""
    MISTRAL_NEMO = "mistral-nemo"
    LLAMA3_8B = "llama3.1:8b"
    LLAMA3_LATEST = "llama3.1"


@dataclass
class RetryConfig:
    """Configuration for retry behavior."""
    max_retries: int = 3
    base_delay: float = 1.0
    max_delay: float = 30.0
    exponential_base: float = 2.0


@dataclass
class GenerationConfig:
    """Configuration for text generation."""
    temperature: float = 0.7
    top_p: float = 0.9
    top_k: int = 40
    num_predict: int = 4096
    repeat_penalty: float = 1.1
    stop: List[str] = field(default_factory=list)


@dataclass
class BookOutline:
    """Represents a book outline."""
    title: str
    genre: str
    audience: str
    chapters: List[Dict[str, str]]
    synopsis: str
    themes: List[str]
    raw_content: str


@dataclass
class ChapterContent:
    """Represents generated chapter content."""
    chapter_num: int
    title: str
    content: str
    word_count: int
    raw_response: str


@dataclass
class ChapterSummary:
    """Represents a chapter summary."""
    chapter_num: Optional[int]
    summary: str
    key_events: List[str]
    characters_introduced: List[str]
    plot_points: List[str]


def retry_async(config: RetryConfig = None):
    """Decorator for async retry logic with exponential backoff."""
    if config is None:
        config = RetryConfig()
    
    def decorator(func: Callable):
        @wraps(func)
        async def wrapper(*args, **kwargs):
            last_exception = None
            delay = config.base_delay
            
            for attempt in range(config.max_retries + 1):
                try:
                    return await func(*args, **kwargs)
                except (aiohttp.ClientError, asyncio.TimeoutError, GenerationError) as e:
                    last_exception = e
                    if attempt < config.max_retries:
                        logger.warning(
                            f"Attempt {attempt + 1}/{config.max_retries + 1} failed: {e}. "
                            f"Retrying in {delay:.1f}s..."
                        )
                        await asyncio.sleep(delay)
                        delay = min(delay * config.exponential_base, config.max_delay)
                    else:
                        logger.error(f"All {config.max_retries + 1} attempts failed")
            
            raise last_exception
        return wrapper
    return decorator


class OllamaClient:
    """
    Async client for Ollama API with book generation capabilities.
    
    Features:
    - Streaming response support
    - Automatic retry with exponential backoff
    - Model preference fallback
    - Comprehensive error handling
    - Production-ready logging
    """
    
    DEFAULT_BASE_URL = "http://localhost:11434"
    
    def __init__(
        self,
        base_url: str = None,
        model: str = None,
        timeout: float = 300.0,
        retry_config: RetryConfig = None,
        generation_config: GenerationConfig = None
    ):
        """
        Initialize the Ollama client.
        
        Args:
            base_url: Ollama API base URL (default: http://localhost:11434)
            model: Preferred model name (auto-selects if not specified)
            timeout: Request timeout in seconds
            retry_config: Retry behavior configuration
            generation_config: Text generation parameters
        """
        self.base_url = base_url or self.DEFAULT_BASE_URL
        self.preferred_model = model
        self.timeout = aiohttp.ClientTimeout(total=timeout)
        self.retry_config = retry_config or RetryConfig()
        self.generation_config = generation_config or GenerationConfig()
        self._session: Optional[aiohttp.ClientSession] = None
        self._active_model: Optional[str] = None
        
    async def __aenter__(self):
        """Async context manager entry."""
        await self.connect()
        return self
    
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        """Async context manager exit."""
        await self.close()
    
    async def connect(self) -> None:
        """Establish connection and verify Ollama is available."""
        if self._session is None:
            self._session = aiohttp.ClientSession(timeout=self.timeout)
        
        try:
            await self._health_check()
            await self._select_model()
            logger.info(f"Connected to Ollama at {self.base_url}, using model: {self._active_model}")
        except Exception as e:
            await self.close()
            raise ConnectionError(f"Failed to connect to Ollama: {e}") from e
    
    async def close(self) -> None:
        """Close the client session."""
        if self._session:
            await self._session.close()
            self._session = None
            logger.info("Ollama client connection closed")
    
    async def _health_check(self) -> bool:
        """Check if Ollama service is healthy."""
        try:
            async with self._session.get(f"{self.base_url}/api/tags") as response:
                if response.status != 200:
                    raise ConnectionError(f"Ollama health check failed: HTTP {response.status}")
                return True
        except aiohttp.ClientConnectorError as e:
            raise ConnectionError(f"Cannot connect to Ollama at {self.base_url}: {e}") from e
    
    async def _select_model(self) -> str:
        """Select the best available model based on preferences."""
        available = await self.list_models()
        available_names = [m["name"] for m in available]
        
        # If user specified a model, verify it exists
        if self.preferred_model:
            if self.preferred_model in available_names:
                self._active_model = self.preferred_model
                return self._active_model
            # Try without tag
            base_name = self.preferred_model.split(":")[0]
            matching = [n for n in available_names if n.startswith(base_name)]
            if matching:
                self._active_model = matching[0]
                logger.info(f"Using {self._active_model} instead of {self.preferred_model}")
                return self._active_model
            raise ModelNotFoundError(f"Model '{self.preferred_model}' not found. Available: {available_names}")
        
        # Auto-select from preferred models
        for pref in ModelPreference:
            if pref.value in available_names:
                self._active_model = pref.value
                return self._active_model
            # Check for partial match
            matching = [n for n in available_names if n.startswith(pref.value.split(":")[0])]
            if matching:
                self._active_model = matching[0]
                return self._active_model
        
        # Fallback to first available model
        if available_names:
            self._active_model = available_names[0]
            logger.warning(f"No preferred model found, using: {self._active_model}")
            return self._active_model
        
        raise ModelNotFoundError("No models available in Ollama")
    
    async def list_models(self) -> List[Dict[str, Any]]:
        """
        List all available models in Ollama.
        
        Returns:
            List of model information dictionaries
        """
        if self._session is None:
            async with aiohttp.ClientSession(timeout=self.timeout) as session:
                async with session.get(f"{self.base_url}/api/tags") as response:
                    if response.status != 200:
                        raise OllamaError(f"Failed to list models: HTTP {response.status}")
                    data = await response.json()
                    return data.get("models", [])
        
        async with self._session.get(f"{self.base_url}/api/tags") as response:
            if response.status != 200:
                raise OllamaError(f"Failed to list models: HTTP {response.status}")
            data = await response.json()
            return data.get("models", [])
    
    @retry_async()
    async def _generate(
        self,
        prompt: str,
        system: str = None,
        stream: bool = False,
        **kwargs
    ) -> str:
        """
        Internal method for text generation.
        
        Args:
            prompt: The prompt to send
            system: System prompt for context
            stream: Whether to use streaming
            **kwargs: Additional generation parameters
            
        Returns:
            Generated text response
        """
        if self._session is None:
            raise ConnectionError("Client not connected. Call connect() first.")
        
        payload = {
            "model": self._active_model,
            "prompt": prompt,
            "stream": stream,
            "options": {
                "temperature": kwargs.get("temperature", self.generation_config.temperature),
                "top_p": kwargs.get("top_p", self.generation_config.top_p),
                "top_k": kwargs.get("top_k", self.generation_config.top_k),
                "num_predict": kwargs.get("num_predict", self.generation_config.num_predict),
                "repeat_penalty": kwargs.get("repeat_penalty", self.generation_config.repeat_penalty),
            }
        }
        
        if system:
            payload["system"] = system
        
        stop_sequences = kwargs.get("stop", self.generation_config.stop)
        if stop_sequences:
            payload["options"]["stop"] = stop_sequences
        
        try:
            async with self._session.post(
                f"{self.base_url}/api/generate",
                json=payload
            ) as response:
                if response.status != 200:
                    error_text = await response.text()
                    raise GenerationError(f"Generation failed: HTTP {response.status} - {error_text}")
                
                if stream:
                    return await self._handle_stream(response)
                else:
                    data = await response.json()
                    return data.get("response", "")
                    
        except asyncio.TimeoutError:
            raise TimeoutError(f"Generation timed out after {self.timeout.total}s")
    
    async def _handle_stream(self, response: aiohttp.ClientResponse) -> str:
        """Handle streaming response from Ollama."""
        full_response = []
        
        async for line in response.content:
            if line:
                try:
                    data = json.loads(line.decode('utf-8'))
                    chunk = data.get("response", "")
                    full_response.append(chunk)
                    
                    if data.get("done", False):
                        break
                except json.JSONDecodeError:
                    continue
        
        return "".join(full_response)
    
    async def generate_stream(
        self,
        prompt: str,
        system: str = None,
        **kwargs
    ) -> AsyncIterator[str]:
        """
        Generate text with streaming output.
        
        Args:
            prompt: The prompt to send
            system: System prompt for context
            **kwargs: Additional generation parameters
            
        Yields:
            Text chunks as they are generated
        """
        if self._session is None:
            raise ConnectionError("Client not connected. Call connect() first.")
        
        payload = {
            "model": self._active_model,
            "prompt": prompt,
            "stream": True,
            "options": {
                "temperature": kwargs.get("temperature", self.generation_config.temperature),
                "top_p": kwargs.get("top_p", self.generation_config.top_p),
                "top_k": kwargs.get("top_k", self.generation_config.top_k),
                "num_predict": kwargs.get("num_predict", self.generation_config.num_predict),
                "repeat_penalty": kwargs.get("repeat_penalty", self.generation_config.repeat_penalty),
            }
        }
        
        if system:
            payload["system"] = system
        
        async with self._session.post(
            f"{self.base_url}/api/generate",
            json=payload
        ) as response:
            if response.status != 200:
                error_text = await response.text()
                raise GenerationError(f"Generation failed: HTTP {response.status} - {error_text}")
            
            async for line in response.content:
                if line:
                    try:
                        data = json.loads(line.decode('utf-8'))
                        chunk = data.get("response", "")
                        if chunk:
                            yield chunk
                        if data.get("done", False):
                            break
                    except json.JSONDecodeError:
                        continue
    
    async def generate_outline(
        self,
        title: str,
        genre: str,
        audience: str,
        num_chapters: int = 12
    ) -> BookOutline:
        """
        Generate a complete book outline.
        
        Args:
            title: Book title
            genre: Book genre (e.g., "fantasy", "mystery", "romance")
            audience: Target audience (e.g., "young adult", "adult", "children")
            num_chapters: Number of chapters to outline
            
        Returns:
            BookOutline object with complete structure
        """
        logger.info(f"Generating outline for '{title}' ({genre}, {audience}, {num_chapters} chapters)")
        
        system_prompt = """You are an expert book outliner and story architect. 
Create detailed, compelling book outlines that provide a strong foundation for writing.
Your outlines should include clear character arcs, plot progression, and thematic elements.
Format your response as a structured outline that can be easily followed."""

        prompt = f"""Create a detailed book outline with the following specifications:

TITLE: {title}
GENRE: {genre}
TARGET AUDIENCE: {audience}
NUMBER OF CHAPTERS: {num_chapters}

Please provide:

1. SYNOPSIS (2-3 paragraphs describing the overall story)

2. MAIN THEMES (list 3-5 central themes)

3. CHAPTER OUTLINE
For each of the {num_chapters} chapters, provide:
- Chapter number and title
- Brief summary (2-3 sentences)
- Key events
- Character development points
- How it advances the main plot

Format each chapter as:
CHAPTER [number]: [Title]
Summary: [summary]
Key Events: [events]
Character Development: [development]
Plot Advancement: [advancement]

Make sure the story has a clear beginning, middle, and end with proper pacing and tension."""

        response = await self._generate(
            prompt=prompt,
            system=system_prompt,
            stream=True,
            temperature=0.8,
            num_predict=8192
        )
        
        # Parse the response into structured outline
        outline = self._parse_outline(response, title, genre, audience, num_chapters)
        logger.info(f"Generated outline with {len(outline.chapters)} chapters")
        
        return outline
    
    def _parse_outline(
        self,
        response: str,
        title: str,
        genre: str,
        audience: str,
        num_chapters: int
    ) -> BookOutline:
        """Parse raw outline response into structured BookOutline."""
        chapters = []
        synopsis = ""
        themes = []
        
        lines = response.split('\n')
        current_chapter = None
        in_synopsis = False
        in_themes = False
        
        for line in lines:
            line_stripped = line.strip()
            
            # Detect synopsis section
            if 'SYNOPSIS' in line_stripped.upper():
                in_synopsis = True
                in_themes = False
                continue
            
            # Detect themes section
            if 'THEME' in line_stripped.upper() and 'MAIN' in line_stripped.upper():
                in_synopsis = False
                in_themes = True
                continue
            
            # Detect chapter section
            if line_stripped.upper().startswith('CHAPTER'):
                in_synopsis = False
                in_themes = False
                
                if current_chapter:
                    chapters.append(current_chapter)
                
                # Extract chapter number and title
                parts = line_stripped.split(':', 1)
                chapter_title = parts[1].strip() if len(parts) > 1 else f"Chapter {len(chapters) + 1}"
                current_chapter = {
                    'number': len(chapters) + 1,
                    'title': chapter_title,
                    'summary': '',
                    'key_events': '',
                    'character_development': '',
                    'plot_advancement': ''
                }
                continue
            
            # Parse content based on current section
            if in_synopsis and line_stripped:
                synopsis += line_stripped + " "
            
            if in_themes and line_stripped:
                # Look for theme items (usually bullet points or numbered)
                if line_stripped.startswith(('-', '*', '•')) or (line_stripped[0].isdigit() and '.' in line_stripped[:3]):
                    theme = line_stripped.lstrip('-*•0123456789. ')
                    if theme:
                        themes.append(theme)
            
            if current_chapter:
                lower_line = line_stripped.lower()
                if lower_line.startswith('summary:'):
                    current_chapter['summary'] = line_stripped[8:].strip()
                elif lower_line.startswith('key events:'):
                    current_chapter['key_events'] = line_stripped[11:].strip()
                elif lower_line.startswith('character development:'):
                    current_chapter['character_development'] = line_stripped[22:].strip()
                elif lower_line.startswith('plot advancement:'):
                    current_chapter['plot_advancement'] = line_stripped[17:].strip()
                elif line_stripped and not any(line_stripped.lower().startswith(p) for p in 
                    ['summary:', 'key events:', 'character development:', 'plot advancement:', 'chapter']):
                    # Append to summary if it's continuation text
                    if current_chapter['summary']:
                        current_chapter['summary'] += " " + line_stripped
        
        # Add last chapter
        if current_chapter:
            chapters.append(current_chapter)
        
        # Ensure we have the requested number of chapters
        while len(chapters) < num_chapters:
            chapters.append({
                'number': len(chapters) + 1,
                'title': f'Chapter {len(chapters) + 1}',
                'summary': 'To be developed',
                'key_events': '',
                'character_development': '',
                'plot_advancement': ''
            })
        
        return BookOutline(
            title=title,
            genre=genre,
            audience=audience,
            chapters=chapters[:num_chapters],
            synopsis=synopsis.strip(),
            themes=themes[:5] if themes else ['Theme to be developed'],
            raw_content=response
        )
    
    async def generate_chapter(
        self,
        outline: BookOutline,
        chapter_num: int,
        style_guide: str = None,
        previous_summary: str = None,
        target_words: int = 3000
    ) -> ChapterContent:
        """
        Generate full chapter content based on outline.
        
        Args:
            outline: The book outline to follow
            chapter_num: Which chapter to generate (1-indexed)
            style_guide: Writing style instructions
            previous_summary: Summary of previous chapters for continuity
            target_words: Approximate word count target
            
        Returns:
            ChapterContent object with full chapter text
        """
        if chapter_num < 1 or chapter_num > len(outline.chapters):
            raise ValueError(f"Invalid chapter number: {chapter_num}. Must be 1-{len(outline.chapters)}")
        
        chapter_info = outline.chapters[chapter_num - 1]
        logger.info(f"Generating Chapter {chapter_num}: {chapter_info['title']}")
        
        system_prompt = f"""You are an expert fiction writer crafting a {outline.genre} novel for {outline.audience} readers.
Write engaging, vivid prose that brings the story to life.
{style_guide or 'Use a compelling narrative voice appropriate for the genre and audience.'}
Focus on showing rather than telling, with natural dialogue and rich descriptions."""

        context_section = ""
        if previous_summary:
            context_section = f"""
PREVIOUS CHAPTERS SUMMARY:
{previous_summary}

"""

        prompt = f"""Write Chapter {chapter_num} of "{outline.title}"

BOOK SYNOPSIS:
{outline.synopsis}

MAIN THEMES: {', '.join(outline.themes)}
{context_section}
CHAPTER {chapter_num}: {chapter_info['title']}
Chapter Summary: {chapter_info['summary']}
Key Events to Include: {chapter_info['key_events']}
Character Development: {chapter_info['character_development']}
Plot Advancement: {chapter_info['plot_advancement']}

TARGET LENGTH: Approximately {target_words} words

Write the complete chapter now. Begin with the chapter title and include all necessary scenes, dialogue, and descriptions to fully realize this part of the story. Ensure smooth transitions and maintain consistency with the established story elements."""

        response = await self._generate(
            prompt=prompt,
            system=system_prompt,
            stream=True,
            temperature=0.75,
            num_predict=target_words * 2  # Allow for variance
        )
        
        word_count = len(response.split())
        logger.info(f"Generated Chapter {chapter_num} with {word_count} words")
        
        return ChapterContent(
            chapter_num=chapter_num,
            title=chapter_info['title'],
            content=response,
            word_count=word_count,
            raw_response=response
        )
    
    async def generate_summary(
        self,
        chapter_content: str,
        chapter_num: int = None
    ) -> ChapterSummary:
        """
        Generate a summary of chapter content for continuity tracking.
        
        Args:
            chapter_content: The full chapter text to summarize
            chapter_num: Optional chapter number for reference
            
        Returns:
            ChapterSummary object with key information
        """
        logger.info(f"Generating summary for chapter {chapter_num or 'unknown'}")
        
        system_prompt = """You are a skilled editor creating concise chapter summaries.
Extract the key information needed to maintain story continuity.
Be precise and include all important plot points, character developments, and setting details."""

        prompt = f"""Analyze the following chapter and provide a structured summary:

CHAPTER CONTENT:
{chapter_content[:8000]}  # Limit content to avoid token limits

Please provide:

1. SUMMARY (2-3 paragraphs covering main events)

2. KEY EVENTS (bullet list of important occurrences)

3. CHARACTERS INTRODUCED OR DEVELOPED (list with brief notes)

4. PLOT POINTS (bullet list of story developments)

5. SETTING/WORLD DETAILS (any new information about the world)

6. UNRESOLVED THREADS (elements that need continuation)

Format your response clearly with these section headers."""

        response = await self._generate(
            prompt=prompt,
            system=system_prompt,
            stream=True,
            temperature=0.3,  # Lower temperature for accuracy
            num_predict=2048
        )
        
        # Parse the summary response
        summary = self._parse_summary(response, chapter_num)
        logger.info(f"Generated summary with {len(summary.key_events)} key events")
        
        return summary
    
    def _parse_summary(self, response: str, chapter_num: int = None) -> ChapterSummary:
        """Parse raw summary response into structured ChapterSummary."""
        summary_text = ""
        key_events = []
        characters = []
        plot_points = []
        
        current_section = None
        lines = response.split('\n')
        
        for line in lines:
            line_stripped = line.strip()
            upper_line = line_stripped.upper()
            
            # Detect sections
            if 'SUMMARY' in upper_line and ('KEY' not in upper_line):
                current_section = 'summary'
                continue
            elif 'KEY EVENT' in upper_line:
                current_section = 'events'
                continue
            elif 'CHARACTER' in upper_line:
                current_section = 'characters'
                continue
            elif 'PLOT POINT' in upper_line:
                current_section = 'plot'
                continue
            elif 'SETTING' in upper_line or 'WORLD' in upper_line:
                current_section = 'setting'
                continue
            elif 'UNRESOLVED' in upper_line:
                current_section = 'unresolved'
                continue
            
            # Parse content based on section
            if line_stripped:
                if current_section == 'summary':
                    summary_text += line_stripped + " "
                elif current_section == 'events':
                    if line_stripped.startswith(('-', '*', '•')) or (len(line_stripped) > 1 and line_stripped[0].isdigit()):
                        event = line_stripped.lstrip('-*•0123456789. ')
                        if event:
                            key_events.append(event)
                elif current_section == 'characters':
                    if line_stripped.startswith(('-', '*', '•')) or (len(line_stripped) > 1 and line_stripped[0].isdigit()):
                        char = line_stripped.lstrip('-*•0123456789. ')
                        if char:
                            characters.append(char)
                elif current_section == 'plot':
                    if line_stripped.startswith(('-', '*', '•')) or (len(line_stripped) > 1 and line_stripped[0].isdigit()):
                        point = line_stripped.lstrip('-*•0123456789. ')
                        if point:
                            plot_points.append(point)
        
        return ChapterSummary(
            chapter_num=chapter_num,
            summary=summary_text.strip() or response[:500],
            key_events=key_events or ['No specific events extracted'],
            characters_introduced=characters or ['No new characters noted'],
            plot_points=plot_points or ['No specific plot points extracted']
        )
    
    @property
    def model(self) -> str:
        """Get the currently active model name."""
        return self._active_model
    
    @property
    def is_connected(self) -> bool:
        """Check if client is connected."""
        return self._session is not None and not self._session.closed


# Convenience functions for standalone usage
async def list_models(base_url: str = None) -> List[Dict[str, Any]]:
    """List available Ollama models without establishing a full connection."""
    client = OllamaClient(base_url=base_url)
    return await client.list_models()


async def quick_generate(
    prompt: str,
    model: str = None,
    base_url: str = None
) -> str:
    """Quick text generation without full client setup."""
    async with OllamaClient(base_url=base_url, model=model) as client:
        return await client._generate(prompt)


# Example usage and testing
async def main():
    """Example usage of the Ollama client for book generation."""
    
    print("Ollama Book Generation Client")
    print("=" * 50)
    
    # List available models
    try:
        models = await list_models()
        print(f"\nAvailable models: {[m['name'] for m in models]}")
    except Exception as e:
        print(f"Error listing models: {e}")
        return
    
    # Create client and generate content
    async with OllamaClient() as client:
        print(f"\nUsing model: {client.model}")
        
        # Generate an outline
        print("\nGenerating book outline...")
        outline = await client.generate_outline(
            title="The Clockwork Conspiracy",
            genre="Steampunk Mystery",
            audience="Young Adult",
            num_chapters=10
        )
        
        print(f"\nOutline Generated:")
        print(f"Title: {outline.title}")
        print(f"Synopsis: {outline.synopsis[:200]}...")
        print(f"Themes: {outline.themes}")
        print(f"Chapters: {len(outline.chapters)}")
        
        for ch in outline.chapters[:3]:
            print(f"  - Chapter {ch['number']}: {ch['title']}")
        
        # Generate first chapter
        print("\nGenerating Chapter 1...")
        chapter = await client.generate_chapter(
            outline=outline,
            chapter_num=1,
            style_guide="Use vivid sensory descriptions and maintain a sense of mystery",
            target_words=2000
        )
        
        print(f"\nChapter 1: {chapter.title}")
        print(f"Word count: {chapter.word_count}")
        print(f"Preview: {chapter.content[:500]}...")
        
        # Generate summary
        print("\nGenerating chapter summary...")
        summary = await client.generate_summary(
            chapter_content=chapter.content,
            chapter_num=1
        )
        
        print(f"\nSummary: {summary.summary[:300]}...")
        print(f"Key events: {summary.key_events[:3]}")


if __name__ == "__main__":
    asyncio.run(main())
