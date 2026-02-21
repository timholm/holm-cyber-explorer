"""
BookForge Configuration System

Loads configuration from multiple sources with the following priority (highest to lowest):
1. Command line arguments
2. Environment variables (prefixed with BOOKFORGE_)
3. config.yaml file
4. Default values
"""

import os
import argparse
import logging
from pathlib import Path
from dataclasses import dataclass, field, asdict
from typing import Optional, Dict, Any
import yaml


# Base directory for the application
BASE_DIR = Path(__file__).parent.resolve()


@dataclass
class DatabaseConfig:
    """Database configuration settings."""
    path: str = str(BASE_DIR / "data" / "bookforge.db")
    
    def __post_init__(self):
        # Ensure the database directory exists
        Path(self.path).parent.mkdir(parents=True, exist_ok=True)


@dataclass
class OllamaConfig:
    """Ollama LLM service configuration."""
    host: str = "http://localhost:11434"
    default_model: str = "llama3.2"
    timeout: int = 300  # Request timeout in seconds
    max_retries: int = 3


@dataclass
class PiperConfig:
    """Piper TTS configuration."""
    voice_model_path: str = str(BASE_DIR / "models" / "piper" / "en_US-lessac-medium.onnx")
    voice_config_path: str = ""  # Auto-derived from model path if empty
    speaker_id: int = 0
    length_scale: float = 1.0  # Speech speed (lower = faster)
    noise_scale: float = 0.667
    noise_w: float = 0.8
    
    def __post_init__(self):
        if not self.voice_config_path:
            self.voice_config_path = self.voice_model_path + ".json"


@dataclass
class AudioConfig:
    """Audio output configuration."""
    output_dir: str = str(BASE_DIR / "output" / "audio")
    format: str = "mp3"  # mp3, wav, ogg
    sample_rate: int = 22050
    bitrate: str = "192k"
    
    def __post_init__(self):
        Path(self.output_dir).mkdir(parents=True, exist_ok=True)


@dataclass
class WorkerConfig:
    """Parallel processing configuration."""
    num_workers: int = 4
    chunk_timeout: int = 600  # Timeout per chunk in seconds
    max_queue_size: int = 100
    retry_failed: bool = True
    max_retries: int = 3


@dataclass
class ContentConfig:
    """Default content generation settings."""
    default_chapter_count: int = 10
    words_per_chapter: int = 3000
    words_per_section: int = 500
    min_word_count: int = 100
    max_word_count: int = 10000
    default_genre: str = "fiction"
    default_style: str = "descriptive"


@dataclass
class ServerConfig:
    """Web server configuration."""
    host: str = "0.0.0.0"
    port: int = 5000
    debug: bool = False
    threaded: bool = True
    secret_key: str = ""  # Generated if empty
    max_content_length: int = 16 * 1024 * 1024  # 16MB max upload
    
    def __post_init__(self):
        if not self.secret_key:
            self.secret_key = os.urandom(24).hex()


@dataclass
class LogConfig:
    """Logging configuration."""
    level: str = "INFO"
    path: str = str(BASE_DIR / "logs" / "bookforge.log")
    format: str = "%(asctime)s - %(name)s - %(levelname)s - %(message)s"
    max_bytes: int = 10 * 1024 * 1024  # 10MB
    backup_count: int = 5
    console_output: bool = True
    
    def __post_init__(self):
        Path(self.path).parent.mkdir(parents=True, exist_ok=True)
    
    def get_level(self) -> int:
        """Convert string level to logging constant."""
        return getattr(logging, self.level.upper(), logging.INFO)


@dataclass
class Config:
    """Main configuration container."""
    database: DatabaseConfig = field(default_factory=DatabaseConfig)
    ollama: OllamaConfig = field(default_factory=OllamaConfig)
    piper: PiperConfig = field(default_factory=PiperConfig)
    audio: AudioConfig = field(default_factory=AudioConfig)
    worker: WorkerConfig = field(default_factory=WorkerConfig)
    content: ContentConfig = field(default_factory=ContentConfig)
    server: ServerConfig = field(default_factory=ServerConfig)
    log: LogConfig = field(default_factory=LogConfig)
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert configuration to dictionary."""
        return asdict(self)
    
    def save_yaml(self, path: str) -> None:
        """Save configuration to YAML file."""
        with open(path, "w") as f:
            yaml.dump(self.to_dict(), f, default_flow_style=False, sort_keys=False)
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Config":
        """Create configuration from dictionary."""
        config = cls()
        
        if "database" in data:
            config.database = DatabaseConfig(**data["database"])
        if "ollama" in data:
            config.ollama = OllamaConfig(**data["ollama"])
        if "piper" in data:
            config.piper = PiperConfig(**data["piper"])
        if "audio" in data:
            config.audio = AudioConfig(**data["audio"])
        if "worker" in data:
            config.worker = WorkerConfig(**data["worker"])
        if "content" in data:
            config.content = ContentConfig(**data["content"])
        if "server" in data:
            config.server = ServerConfig(**data["server"])
        if "log" in data:
            config.log = LogConfig(**data["log"])
        
        return config


def load_yaml_config(config_path: Optional[str] = None) -> Dict[str, Any]:
    """Load configuration from YAML file."""
    if config_path is None:
        config_path = os.environ.get("BOOKFORGE_CONFIG_PATH", str(BASE_DIR / "config.yaml"))
    
    path = Path(config_path)
    if path.exists():
        with open(path, "r") as f:
            return yaml.safe_load(f) or {}
    return {}


def load_env_config() -> Dict[str, Any]:
    """Load configuration from environment variables.
    
    Environment variables are prefixed with BOOKFORGE_ and use double underscores
    for nested keys. For example:
        BOOKFORGE_DATABASE__PATH=/path/to/db
        BOOKFORGE_OLLAMA__HOST=http://localhost:11434
        BOOKFORGE_SERVER__DEBUG=true
    """
    config: Dict[str, Any] = {}
    prefix = "BOOKFORGE_"
    
    for key, value in os.environ.items():
        if not key.startswith(prefix):
            continue
        
        # Remove prefix and split by double underscore
        key_path = key[len(prefix):].lower().split("__")
        
        # Convert value types
        if value.lower() in ("true", "yes", "1"):
            value = True
        elif value.lower() in ("false", "no", "0"):
            value = False
        elif value.isdigit():
            value = int(value)
        elif value.replace(".", "", 1).isdigit():
            value = float(value)
        
        # Build nested dictionary
        current = config
        for part in key_path[:-1]:
            if part not in current:
                current[part] = {}
            current = current[part]
        current[key_path[-1]] = value
    
    return config


def create_arg_parser() -> argparse.ArgumentParser:
    """Create command line argument parser."""
    parser = argparse.ArgumentParser(
        description="BookForge - AI-powered audiobook generation",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  bookforge --debug --port 8080
  bookforge --ollama-model mistral --workers 8
  bookforge --config /path/to/config.yaml
        """
    )
    
    # General options
    parser.add_argument(
        "--config", "-c",
        dest="config_path",
        help="Path to configuration file (default: config.yaml)"
    )
    
    # Database options
    parser.add_argument(
        "--db-path",
        dest="database_path",
        help="Path to SQLite database"
    )
    
    # Ollama options
    parser.add_argument(
        "--ollama-host",
        dest="ollama_host",
        help="Ollama API host URL"
    )
    parser.add_argument(
        "--ollama-model", "-m",
        dest="ollama_model",
        help="Default Ollama model to use"
    )
    
    # Piper options
    parser.add_argument(
        "--voice-model",
        dest="piper_voice_model",
        help="Path to Piper voice model"
    )
    
    # Audio options
    parser.add_argument(
        "--audio-output", "-o",
        dest="audio_output",
        help="Directory for audio output files"
    )
    
    # Worker options
    parser.add_argument(
        "--workers", "-w",
        dest="num_workers",
        type=int,
        help="Number of parallel workers"
    )
    
    # Content options
    parser.add_argument(
        "--chapters",
        dest="chapter_count",
        type=int,
        help="Default number of chapters"
    )
    parser.add_argument(
        "--words-per-chapter",
        dest="words_per_chapter",
        type=int,
        help="Target words per chapter"
    )
    
    # Server options
    parser.add_argument(
        "--host", "-H",
        dest="server_host",
        help="Server host address"
    )
    parser.add_argument(
        "--port", "-p",
        dest="server_port",
        type=int,
        help="Server port number"
    )
    parser.add_argument(
        "--debug", "-d",
        dest="debug",
        action="store_true",
        help="Enable debug mode"
    )
    
    # Logging options
    parser.add_argument(
        "--log-level", "-l",
        dest="log_level",
        choices=["DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL"],
        help="Logging level"
    )
    parser.add_argument(
        "--log-path",
        dest="log_path",
        help="Path to log file"
    )
    
    return parser


def parse_args_to_config(args: argparse.Namespace) -> Dict[str, Any]:
    """Convert parsed arguments to configuration dictionary."""
    config: Dict[str, Any] = {}
    
    # Map arguments to nested config structure
    arg_mapping = {
        "database_path": ("database", "path"),
        "ollama_host": ("ollama", "host"),
        "ollama_model": ("ollama", "default_model"),
        "piper_voice_model": ("piper", "voice_model_path"),
        "audio_output": ("audio", "output_dir"),
        "num_workers": ("worker", "num_workers"),
        "chapter_count": ("content", "default_chapter_count"),
        "words_per_chapter": ("content", "words_per_chapter"),
        "server_host": ("server", "host"),
        "server_port": ("server", "port"),
        "debug": ("server", "debug"),
        "log_level": ("log", "level"),
        "log_path": ("log", "path"),
    }
    
    for arg_name, (section, key) in arg_mapping.items():
        value = getattr(args, arg_name, None)
        if value is not None:
            if section not in config:
                config[section] = {}
            config[section][key] = value
    
    return config


def deep_merge(base: Dict[str, Any], override: Dict[str, Any]) -> Dict[str, Any]:
    """Deep merge two dictionaries, with override taking precedence."""
    result = base.copy()
    
    for key, value in override.items():
        if key in result and isinstance(result[key], dict) and isinstance(value, dict):
            result[key] = deep_merge(result[key], value)
        else:
            result[key] = value
    
    return result


def load_config(args: Optional[argparse.Namespace] = None) -> Config:
    """Load configuration from all sources.
    
    Priority (highest to lowest):
    1. Command line arguments
    2. Environment variables
    3. YAML configuration file
    4. Default values
    
    Args:
        args: Parsed command line arguments. If None, will parse sys.argv.
    
    Returns:
        Fully configured Config object.
    """
    # Start with empty config (defaults will come from dataclasses)
    merged_config: Dict[str, Any] = {}
    
    # Parse command line arguments if not provided
    if args is None:
        parser = create_arg_parser()
        args = parser.parse_args()
    
    # Get config file path from args or environment
    config_path = getattr(args, "config_path", None)
    
    # Load from YAML file (lowest priority of external sources)
    yaml_config = load_yaml_config(config_path)
    merged_config = deep_merge(merged_config, yaml_config)
    
    # Load from environment variables (medium priority)
    env_config = load_env_config()
    merged_config = deep_merge(merged_config, env_config)
    
    # Load from command line arguments (highest priority)
    args_config = parse_args_to_config(args)
    merged_config = deep_merge(merged_config, args_config)
    
    # Create Config object from merged dictionary
    return Config.from_dict(merged_config)


def setup_logging(config: Config) -> logging.Logger:
    """Configure logging based on configuration.
    
    Args:
        config: Configuration object.
    
    Returns:
        Configured root logger.
    """
    from logging.handlers import RotatingFileHandler
    
    log_config = config.log
    logger = logging.getLogger("bookforge")
    logger.setLevel(log_config.get_level())
    
    # Clear existing handlers
    logger.handlers.clear()
    
    # Create formatter
    formatter = logging.Formatter(log_config.format)
    
    # File handler with rotation
    file_handler = RotatingFileHandler(
        log_config.path,
        maxBytes=log_config.max_bytes,
        backupCount=log_config.backup_count
    )
    file_handler.setLevel(log_config.get_level())
    file_handler.setFormatter(formatter)
    logger.addHandler(file_handler)
    
    # Console handler
    if log_config.console_output:
        console_handler = logging.StreamHandler()
        console_handler.setLevel(log_config.get_level())
        console_handler.setFormatter(formatter)
        logger.addHandler(console_handler)
    
    return logger


def generate_default_config_yaml(output_path: Optional[str] = None) -> str:
    """Generate default configuration YAML file.
    
    Args:
        output_path: Path to write the file. If None, returns the YAML string.
    
    Returns:
        YAML configuration string.
    """
    config = Config()
    yaml_content = yaml.dump(config.to_dict(), default_flow_style=False, sort_keys=False)
    
    if output_path:
        with open(output_path, "w") as f:
            f.write(yaml_content)
    
    return yaml_content


# Global configuration instance (lazy-loaded)
_config: Optional[Config] = None


def get_config() -> Config:
    """Get the global configuration instance.
    
    Loads configuration on first access.
    
    Returns:
        Global Config instance.
    """
    global _config
    if _config is None:
        _config = load_config()
    return _config


def reload_config() -> Config:
    """Reload configuration from all sources.
    
    Returns:
        Newly loaded Config instance.
    """
    global _config
    _config = load_config()
    return _config


# Convenience accessors
def get_database_config() -> DatabaseConfig:
    """Get database configuration."""
    return get_config().database


def get_ollama_config() -> OllamaConfig:
    """Get Ollama configuration."""
    return get_config().ollama


def get_piper_config() -> PiperConfig:
    """Get Piper TTS configuration."""
    return get_config().piper


def get_audio_config() -> AudioConfig:
    """Get audio configuration."""
    return get_config().audio


def get_worker_config() -> WorkerConfig:
    """Get worker configuration."""
    return get_config().worker


def get_content_config() -> ContentConfig:
    """Get content generation configuration."""
    return get_config().content


def get_server_config() -> ServerConfig:
    """Get server configuration."""
    return get_config().server


def get_log_config() -> LogConfig:
    """Get logging configuration."""
    return get_config().log


if __name__ == "__main__":
    # When run directly, print current configuration
    import json
    
    config = load_config()
    print("Current BookForge Configuration:")
    print("=" * 50)
    print(json.dumps(config.to_dict(), indent=2, default=str))
