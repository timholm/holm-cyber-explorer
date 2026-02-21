#!/usr/bin/env python3
"""
Static site builder for holm.chat documentation.
Converts markdown files to HTML with proper structure, navigation, and styling.
No external dependencies -- uses only Python standard library.
"""

import re
import os
import html

SITE_DIR = os.path.join(os.path.dirname(os.path.abspath(__file__)), 'site')
SRC_DIR = os.path.dirname(os.path.abspath(__file__))

# Page definitions: (output_filename, source_md, page_title, nav_label)
PAGES = [
    ('domains-1-5.html', 'stage1-domains-1-5.md', 'Domains 1-5: Constitution, Governance, Security, Infrastructure, Platform', 'Domains 1-5'),
    ('domains-6-10.html', 'stage1-domains-6-10.md', 'Domains 6-10: Data, Intelligence, Automation, Education, Operations', 'Domains 6-10'),
    ('domains-11-15.html', 'stage1-domains-11-15.md', 'Domains 11-15: Administration, Disaster Recovery, Evolution, Research, Ethics', 'Domains 11-15'),
    ('domains-16-20.html', 'stage1-domains-16-20.md', 'Domains 16-20: Interface, Scaling, Import, Quality Assurance, Memory', 'Domains 16-20'),
    ('meta-framework.html', 'stage1-meta-framework.md', 'Meta-Framework: Unifying Standards', 'Meta-Framework'),
]

# Domain data for the homepage
DOMAINS = [
    (1, 'Constitution & Philosophy', 'Define the foundational beliefs, principles, and purpose of the institution.', 'domains-1-5.html#domain-1-constitution--philosophy'),
    (2, 'Governance & Authority', 'Define who decides what, how authority flows, how disputes resolve.', 'domains-1-5.html#domain-2-governance--authority'),
    (3, 'Security & Integrity', 'Define threat models, access control, cryptographic principles, trust boundaries.', 'domains-1-5.html#domain-3-security--integrity'),
    (4, 'Infrastructure & Power', 'Define physical infrastructure, power systems, network topology, hardware lifecycle.', 'domains-1-5.html#domain-4-infrastructure--power'),
    (5, 'Platform & Core Systems', 'Define the operating system layer, core services, system architecture.', 'domains-1-5.html#domain-5-platform--core-systems'),
    (6, 'Data & Archives', 'Define data models, storage philosophy, archival strategy, format longevity.', 'domains-6-10.html#domain-6-data--archives'),
    (7, 'Intelligence & Analysis', 'Define how the institution gathers, processes, and acts on information.', 'domains-6-10.html#domain-7-intelligence--analysis'),
    (8, 'Automation & Agents', 'Define automation philosophy, agent boundaries, human-in-the-loop requirements.', 'domains-6-10.html#domain-8-automation--agents'),
    (9, 'Education & Training', 'Define how knowledge is transferred across generations.', 'domains-6-10.html#domain-9-education--training'),
    (10, 'User Operations', 'Define daily workflows, operational procedures, routine maintenance.', 'domains-6-10.html#domain-10-user-operations'),
    (11, 'Administration', 'Define resource management, scheduling, budgeting, procurement, inventory.', 'domains-11-15.html#domain-11-administration'),
    (12, 'Disaster Recovery', 'Define failure scenarios, recovery procedures, backup verification, continuity.', 'domains-11-15.html#domain-12-disaster-recovery'),
    (13, 'Evolution & Adaptation', 'Define how the institution changes deliberately over time.', 'domains-11-15.html#domain-13-evolution--adaptation'),
    (14, 'Research & Theory', 'Define how the institution creates new knowledge and tests hypotheses.', 'domains-11-15.html#domain-14-research--theory'),
    (15, 'Ethics & Safeguards', 'Define ethical boundaries, oversight structures, whistleblower protections.', 'domains-11-15.html#domain-15-ethics--safeguards'),
    (16, 'Interface & Navigation', 'Define how users interact with the system, information architecture.', 'domains-16-20.html#domain-16-interface--navigation'),
    (17, 'Scaling & Federation', 'Define multi-node architecture, federation protocols, distributed governance.', 'domains-16-20.html#domain-17-scaling--federation'),
    (18, 'Import & Quarantine', 'Define how external content enters the system, validation pipelines.', 'domains-16-20.html#domain-18-import--quarantine'),
    (19, 'Quality Assurance', 'Define quality standards, testing frameworks, audit procedures.', 'domains-19-20.html#domain-19-quality-assurance' if False else 'domains-16-20.html#domain-19-quality-assurance'),
    (20, 'Institutional Memory', 'Define how the institution remembers: decision logs, oral history, lessons learned.', 'domains-16-20.html#domain-20-institutional-memory'),
]

FAVICON_SVG = '''<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32">
  <rect width="32" height="32" rx="4" fill="#1a1a2e"/>
  <text x="16" y="22" font-family="monospace" font-size="18" font-weight="bold" fill="#e0e0e0" text-anchor="middle">H</text>
</svg>'''


def slugify(text):
    """Convert text to URL-safe slug for anchor IDs."""
    text = text.lower().strip()
    text = re.sub(r'[^\w\s-]', '', text)
    text = re.sub(r'[\s]+', '-', text)
    text = re.sub(r'-+', '-', text)
    text = text.strip('-')
    return text


def escape(text):
    """HTML-escape text."""
    return html.escape(text, quote=True)


class MarkdownConverter:
    """Convert markdown to HTML. Handles headings, lists, tables, code blocks,
    bold, italic, links, and horizontal rules."""

    def __init__(self):
        self.toc = []  # list of (level, id, text)

    def convert(self, md_text):
        """Convert full markdown text to HTML, collecting TOC entries."""
        self.toc = []
        lines = md_text.split('\n')
        html_parts = []
        i = 0
        in_list = False
        list_type = None  # 'ul' or 'ol'
        in_table = False
        table_html = []

        while i < len(lines):
            line = lines[i]

            # Fenced code block
            if line.strip().startswith('```'):
                lang = line.strip()[3:].strip()
                code_lines = []
                i += 1
                while i < len(lines) and not lines[i].strip().startswith('```'):
                    code_lines.append(lines[i])
                    i += 1
                i += 1  # skip closing ```
                code_content = escape('\n'.join(code_lines))
                lang_attr = f' class="language-{escape(lang)}"' if lang else ''
                html_parts.append(f'<pre><code{lang_attr}>{code_content}</code></pre>')
                continue

            # Close table if we were in one and line doesn't look like table
            if in_table and not line.strip().startswith('|'):
                html_parts.append('</tbody></table></div>')
                in_table = False
                table_html = []

            # Close list if we were in one and this line is not a list item or continuation
            if in_list and line.strip() and not re.match(r'^(\s*[-*+]\s|\s*\d+\.\s)', line) and not line.startswith('  '):
                html_parts.append(f'</{list_type}>')
                in_list = False
                list_type = None

            # Blank line
            if not line.strip():
                if in_list:
                    # Don't close list on single blank line; check next
                    pass
                i += 1
                continue

            # Horizontal rule
            if re.match(r'^---+\s*$', line.strip()) or re.match(r'^\*\*\*+\s*$', line.strip()):
                if in_list:
                    html_parts.append(f'</{list_type}>')
                    in_list = False
                    list_type = None
                html_parts.append('<hr/>')
                i += 1
                continue

            # Headings
            heading_match = re.match(r'^(#{1,6})\s+(.+)$', line)
            if heading_match:
                if in_list:
                    html_parts.append(f'</{list_type}>')
                    in_list = False
                    list_type = None
                level = len(heading_match.group(1))
                text = heading_match.group(2).strip()
                anchor_id = slugify(text)
                # Ensure unique anchors
                base_id = anchor_id
                counter = 1
                existing_ids = [t[1] for t in self.toc]
                while anchor_id in existing_ids:
                    anchor_id = f'{base_id}-{counter}'
                    counter += 1
                self.toc.append((level, anchor_id, text))
                inline_text = self._inline(text)
                html_parts.append(f'<h{level} id="{anchor_id}">{inline_text}</h{level}>')
                i += 1
                continue

            # Table
            if line.strip().startswith('|'):
                if not in_table:
                    # Start table
                    in_table = True
                    table_html = []
                    # Parse header row
                    cells = self._parse_table_row(line)
                    # Check if next line is separator
                    if i + 1 < len(lines) and re.match(r'^\|[\s\-:|]+\|', lines[i+1].strip()):
                        # Determine alignment from separator
                        sep_line = lines[i+1].strip()
                        alignments = self._parse_table_alignments(sep_line)
                        header_html = '<thead><tr>'
                        for j, cell in enumerate(cells):
                            align = alignments[j] if j < len(alignments) else ''
                            align_attr = f' style="text-align:{align}"' if align else ''
                            header_html += f'<th{align_attr}>{self._inline(cell)}</th>'
                        header_html += '</tr></thead><tbody>'
                        html_parts.append(f'<div class="table-wrap"><table>{header_html}')
                        i += 2  # skip header + separator
                        continue
                    else:
                        # No separator, treat as data
                        html_parts.append('<div class="table-wrap"><table><tbody>')
                # Data row
                cells = self._parse_table_row(line)
                row_html = '<tr>'
                for cell in cells:
                    row_html += f'<td>{self._inline(cell)}</td>'
                row_html += '</tr>'
                html_parts.append(row_html)
                i += 1
                continue

            # Unordered list
            ul_match = re.match(r'^(\s*)[-*+]\s+(.+)$', line)
            if ul_match:
                indent = len(ul_match.group(1))
                content = ul_match.group(2)
                if not in_list:
                    in_list = True
                    list_type = 'ul'
                    html_parts.append('<ul>')
                elif list_type != 'ul':
                    html_parts.append(f'</{list_type}>')
                    list_type = 'ul'
                    html_parts.append('<ul>')
                html_parts.append(f'<li>{self._inline(content)}</li>')
                i += 1
                continue

            # Ordered list
            ol_match = re.match(r'^(\s*)\d+\.\s+(.+)$', line)
            if ol_match:
                content = ol_match.group(2)
                if not in_list:
                    in_list = True
                    list_type = 'ol'
                    html_parts.append('<ol>')
                elif list_type != 'ol':
                    html_parts.append(f'</{list_type}>')
                    list_type = 'ol'
                    html_parts.append('<ol>')
                html_parts.append(f'<li>{self._inline(content)}</li>')
                i += 1
                continue

            # Checkbox list items (treat as ul)
            cb_match = re.match(r'^(\s*)[-*+]\s+\[[ x]\]\s+(.+)$', line)
            if cb_match:
                content = cb_match.group(2)
                checked = '[x]' in line[:line.index(content)]
                if not in_list:
                    in_list = True
                    list_type = 'ul'
                    html_parts.append('<ul class="checklist">')
                marker = '&#9745;' if checked else '&#9744;'
                html_parts.append(f'<li>{marker} {self._inline(content)}</li>')
                i += 1
                continue

            # Paragraph
            para_lines = [line]
            i += 1
            while i < len(lines) and lines[i].strip() and not lines[i].strip().startswith('#') \
                    and not lines[i].strip().startswith('|') and not lines[i].strip().startswith('```') \
                    and not re.match(r'^---+\s*$', lines[i].strip()) \
                    and not re.match(r'^(\s*)[-*+]\s+', lines[i]) \
                    and not re.match(r'^(\s*)\d+\.\s+', lines[i]):
                para_lines.append(lines[i])
                i += 1
            para_text = ' '.join(para_lines)
            html_parts.append(f'<p>{self._inline(para_text)}</p>')

        # Close any open structures
        if in_list:
            html_parts.append(f'</{list_type}>')
        if in_table:
            html_parts.append('</tbody></table></div>')

        return '\n'.join(html_parts)

    def _parse_table_row(self, line):
        """Parse a table row into cells."""
        line = line.strip()
        if line.startswith('|'):
            line = line[1:]
        if line.endswith('|'):
            line = line[:-1]
        return [cell.strip() for cell in line.split('|')]

    def _parse_table_alignments(self, sep_line):
        """Parse alignment from separator line."""
        cells = self._parse_table_row(sep_line)
        alignments = []
        for cell in cells:
            cell = cell.strip()
            if cell.startswith(':') and cell.endswith(':'):
                alignments.append('center')
            elif cell.endswith(':'):
                alignments.append('right')
            elif cell.startswith(':'):
                alignments.append('left')
            else:
                alignments.append('')
        return alignments

    def _inline(self, text):
        """Process inline markdown: bold, italic, code, links."""
        # Inline code first (to avoid processing markdown inside code)
        parts = []
        code_pattern = re.compile(r'`([^`]+)`')
        last = 0
        for m in code_pattern.finditer(text):
            parts.append(self._inline_formatting(text[last:m.start()]))
            parts.append(f'<code>{escape(m.group(1))}</code>')
            last = m.end()
        parts.append(self._inline_formatting(text[last:]))
        return ''.join(parts)

    def _inline_formatting(self, text):
        """Process bold, italic, and links."""
        # Bold+italic
        text = re.sub(r'\*\*\*(.+?)\*\*\*', r'<strong><em>\1</em></strong>', text)
        # Bold
        text = re.sub(r'\*\*(.+?)\*\*', r'<strong>\1</strong>', text)
        # Italic
        text = re.sub(r'\*(.+?)\*', r'<em>\1</em>', text)
        # Links [text](url)
        text = re.sub(r'\[([^\]]+)\]\(([^)]+)\)', r'<a href="\2">\1</a>', text)
        return text


def build_toc_html(toc_entries, max_level=3):
    """Build a table of contents HTML from toc entries."""
    if not toc_entries:
        return ''
    # Filter to relevant levels (h1-h3 for domain pages, h2-h3 for others)
    min_level = min(e[0] for e in toc_entries) if toc_entries else 1
    filtered = [(lvl, aid, txt) for lvl, aid, txt in toc_entries if lvl <= max_level]
    if not filtered:
        return ''

    html_parts = ['<nav class="toc"><h2 class="toc-title">Table of Contents</h2><ul>']
    prev_level = filtered[0][0]
    for level, anchor_id, text in filtered:
        if level > prev_level:
            html_parts.append('<ul>' * (level - prev_level))
        elif level < prev_level:
            html_parts.append('</ul>' * (prev_level - level))
        html_parts.append(f'<li><a href="#{anchor_id}">{escape(text)}</a></li>')
        prev_level = level
    # Close remaining
    html_parts.append('</ul>' * (prev_level - filtered[0][0] + 1))
    html_parts.append('</nav>')
    return '\n'.join(html_parts)


# Global list populated during build so nav and index stay in sync
ALL_NAV_ITEMS = []


def nav_html(current_page=''):
    """Build the site navigation from ALL_NAV_ITEMS."""
    parts = ['<nav class="site-nav"><ul>']
    home_active = ' class="active"' if current_page == 'index.html' else ''
    parts.append(f'<li{home_active}><a href="index.html">Home</a></li>')
    for href, label, *_ in ALL_NAV_ITEMS:
        active = ' class="active"' if href == current_page else ''
        parts.append(f'<li{active}><a href="{href}">{label}</a></li>')
    parts.append('</ul></nav>')
    return '\n'.join(parts)


def page_template(title, content, toc_html, current_page=''):
    """Wrap content in full HTML page."""
    return f'''<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{escape(title)} - holm.chat</title>
<link rel="stylesheet" href="style.css">
<link rel="icon" type="image/svg+xml" href="data:image/svg+xml,{FAVICON_SVG.replace('"', '%22').replace('#', '%23').replace('<', '%3C').replace('>', '%3E').replace(' ', '%20').replace(chr(10), '')}">
</head>
<body>
<header class="site-header">
<div class="header-inner">
<a href="index.html" class="site-title">holm.chat</a>
<span class="site-subtitle">Documentation Institution</span>
</div>
{nav_html(current_page)}
</header>
<div class="layout">
<aside class="sidebar">
{toc_html}
</aside>
<main class="content">
{content}
</main>
</div>
<footer class="site-footer">
<div class="footer-inner">
<p>holm.chat Documentation Institution &mdash; Air-Gapped, Off-Grid, Self-Built</p>
<p>Stage 1: Documentation Framework &mdash; Version 1.0.0 &mdash; 2026-02-16</p>
</div>
</footer>
</body>
</html>'''


def build_index_page():
    """Build the homepage."""
    domain_cards = []
    for num, name, desc, link in DOMAINS:
        domain_cards.append(f'''<a href="{link}" class="domain-card">
<div class="domain-number">{num:02d}</div>
<div class="domain-info">
<h3>{escape(name)}</h3>
<p>{escape(desc)}</p>
</div>
</a>''')

    content = f'''<div class="hero">
<h1>Documentation Framework</h1>
<p class="hero-sub">A Lifelong, Air-Gapped, Off-Grid, Self-Built Digital Institution</p>
<p class="hero-meta">Document ID: STAGE1-FRAMEWORK &mdash; Version 1.0.0 &mdash; 2026-02-16<br>
Intended Lifespan: 50+ years &mdash; Status: Initial Framework</p>
</div>

<section class="intro">
<h2>Preamble</h2>
<p>This documentation framework defines a self-sovereign digital institution designed to operate indefinitely without dependence on external networks, cloud services, commercial vendors, or institutional continuity of any single person. Every domain is written with the assumption that the original authors will eventually be unavailable, that hardware will be replaced many times over, and that the cultural context surrounding this institution will shift in ways we cannot predict.</p>
<p>Nothing in this framework is code. Nothing is tooling. This is the map that precedes the territory.</p>
<p>The framework spans <strong>20 domains</strong> and proposes approximately <strong>239-398 articles</strong> across the full documentation corpus. It is designed to constitute the institutional memory of a self-sovereign digital institution built to outlive its creators.</p>
</section>

<section class="domain-grid-section">
<h2>All 20 Domains</h2>
<div class="domain-grid">
{''.join(domain_cards)}
</div>
</section>

<section class="meta-section">
<h2>All Documents</h2>
<div class="doc-links">
{''.join(f"""<a href="{href}" class="doc-link">
<strong>{escape(label)}</strong>
<span>{escape(subtitle or '')}</span>
</a>""" for href, label, subtitle in ALL_NAV_ITEMS)}
</div>
</section>

<section class="principles-section">
<h2>Meta-Rules Governing All Documentation</h2>
<ol>
<li><strong>Fifty-Year Horizon.</strong> Every article must be written as though the reader has never met the author and lives in a substantially different technological and cultural context.</li>
<li><strong>Loss of Original Author.</strong> No article may depend on oral tradition, tacit knowledge, or the continued availability of any specific person.</li>
<li><strong>Hardware Impermanence.</strong> All references to specific hardware must include abstraction layers and migration paths.</li>
<li><strong>Cultural Drift.</strong> Terminology must be defined inline. Jargon must be explained. Assumptions must be made explicit.</li>
<li><strong>Self-Containment.</strong> This institution is air-gapped. Every external reference must be archived locally.</li>
<li><strong>Contradiction Resolution.</strong> Where two articles conflict, the article closer to Domain 1 (Constitution) takes precedence.</li>
<li><strong>Plain Language.</strong> Prefer clarity over elegance. Prefer redundancy over ambiguity.</li>
</ol>
</section>
'''
    return page_template(
        'Documentation Framework',
        content,
        '',
        'index.html'
    )


def build_404_page():
    """Build the 404 page."""
    content = '''<div class="error-page">
<h1>404</h1>
<p class="error-subtitle">Document Not Found</p>
<p>The requested page does not exist in this documentation set. This may indicate a broken cross-reference, a deprecated article, or an incorrect URL.</p>
<h2>Recovery Procedures</h2>
<ol>
<li>Return to the <a href="index.html">homepage</a> and navigate from the domain index.</li>
<li>Check the URL for typographical errors.</li>
<li>If you followed a link from within the documentation, record this as a defect in the Commentary Section of the referring article.</li>
</ol>
<p><a href="index.html" class="btn">Return to Homepage</a></p>
</div>'''
    return page_template('404 - Document Not Found', content, '', '')


def build_css():
    """Build the stylesheet."""
    return '''/* holm.chat Documentation Institution -- Stylesheet
   Design: Dark theme, monospace-friendly, optimized for long-form reading.
   Air-gapped compatible: no external fonts, no CDNs, no JavaScript required. */

/* === RESET & BASE === */
*, *::before, *::after {
    box-sizing: border-box;
    margin: 0;
    padding: 0;
}

:root {
    --bg-primary: #0f0f1a;
    --bg-secondary: #161625;
    --bg-tertiary: #1c1c30;
    --bg-card: #1a1a2e;
    --bg-code: #12121f;
    --text-primary: #d4d4e0;
    --text-secondary: #9a9ab0;
    --text-muted: #6a6a80;
    --accent: #7b8cde;
    --accent-dim: #4a5490;
    --border: #2a2a40;
    --border-light: #353550;
    --link: #8b9cf0;
    --link-hover: #aab4f5;
    --success: #5fa85f;
    --warning: #c49a3a;
    --danger: #c45a5a;
    --font-mono: "SF Mono", "Fira Code", "Fira Mono", "Roboto Mono", "Cascadia Code", Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
    --font-body: "Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
    --content-width: 820px;
    --sidebar-width: 280px;
}

html {
    font-size: 16px;
    scroll-behavior: smooth;
}

body {
    font-family: var(--font-body);
    background: var(--bg-primary);
    color: var(--text-primary);
    line-height: 1.7;
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    -webkit-font-smoothing: antialiased;
}

/* === HEADER === */
.site-header {
    background: var(--bg-secondary);
    border-bottom: 1px solid var(--border);
    position: sticky;
    top: 0;
    z-index: 100;
}

.header-inner {
    max-width: 1200px;
    margin: 0 auto;
    padding: 1rem 2rem 0.5rem;
    display: flex;
    align-items: baseline;
    gap: 1rem;
}

.site-title {
    font-family: var(--font-mono);
    font-size: 1.4rem;
    font-weight: 700;
    color: var(--text-primary);
    text-decoration: none;
    letter-spacing: -0.02em;
}

.site-title:hover {
    color: var(--accent);
}

.site-subtitle {
    font-size: 0.85rem;
    color: var(--text-muted);
    font-family: var(--font-mono);
}

/* === NAVIGATION === */
.site-nav {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 2rem;
    overflow-x: auto;
}

.site-nav ul {
    display: flex;
    list-style: none;
    gap: 0;
    border-bottom: none;
}

.site-nav li {
    flex-shrink: 0;
}

.site-nav a {
    display: block;
    padding: 0.6rem 1rem;
    color: var(--text-secondary);
    text-decoration: none;
    font-size: 0.85rem;
    font-family: var(--font-mono);
    border-bottom: 2px solid transparent;
    transition: color 0.15s, border-color 0.15s;
    white-space: nowrap;
}

.site-nav a:hover {
    color: var(--text-primary);
    border-bottom-color: var(--accent-dim);
}

.site-nav li.active a {
    color: var(--accent);
    border-bottom-color: var(--accent);
}

/* === LAYOUT === */
.layout {
    display: flex;
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
    gap: 2rem;
    flex: 1;
    width: 100%;
}

.sidebar {
    width: var(--sidebar-width);
    flex-shrink: 0;
    position: sticky;
    top: 80px;
    max-height: calc(100vh - 100px);
    overflow-y: auto;
    padding-right: 1rem;
}

.content {
    flex: 1;
    min-width: 0;
    max-width: var(--content-width);
}

/* === TABLE OF CONTENTS === */
.toc {
    font-size: 0.82rem;
    line-height: 1.5;
}

.toc-title {
    font-family: var(--font-mono);
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: var(--text-muted);
    margin-bottom: 0.75rem;
    padding-bottom: 0.5rem;
    border-bottom: 1px solid var(--border);
}

.toc ul {
    list-style: none;
    padding-left: 0;
}

.toc ul ul {
    padding-left: 1rem;
}

.toc li {
    margin: 0.2rem 0;
}

.toc a {
    color: var(--text-secondary);
    text-decoration: none;
    display: block;
    padding: 0.15rem 0.5rem;
    border-radius: 3px;
    transition: color 0.15s, background 0.15s;
}

.toc a:hover {
    color: var(--text-primary);
    background: var(--bg-tertiary);
}

/* === TYPOGRAPHY === */
h1, h2, h3, h4, h5, h6 {
    font-family: var(--font-body);
    font-weight: 700;
    line-height: 1.3;
    margin-top: 2.5rem;
    margin-bottom: 1rem;
    color: var(--text-primary);
}

h1 {
    font-size: 2rem;
    margin-top: 0;
    padding-bottom: 0.5rem;
    border-bottom: 2px solid var(--border);
}

h2 {
    font-size: 1.5rem;
    padding-bottom: 0.3rem;
    border-bottom: 1px solid var(--border);
}

h3 { font-size: 1.2rem; }
h4 { font-size: 1.05rem; }
h5 { font-size: 0.95rem; }
h6 { font-size: 0.9rem; color: var(--text-secondary); }

p {
    margin-bottom: 1rem;
}

a {
    color: var(--link);
    text-decoration: none;
}

a:hover {
    color: var(--link-hover);
    text-decoration: underline;
}

strong {
    font-weight: 700;
    color: var(--text-primary);
}

em {
    font-style: italic;
    color: var(--text-secondary);
}

hr {
    border: none;
    border-top: 1px solid var(--border);
    margin: 2rem 0;
}

/* === LISTS === */
ul, ol {
    margin-bottom: 1rem;
    padding-left: 1.5rem;
}

li {
    margin-bottom: 0.3rem;
}

li > ul, li > ol {
    margin-top: 0.3rem;
    margin-bottom: 0.3rem;
}

.checklist {
    list-style: none;
    padding-left: 0;
}

/* === CODE === */
code {
    font-family: var(--font-mono);
    font-size: 0.88em;
    background: var(--bg-code);
    color: var(--accent);
    padding: 0.15em 0.4em;
    border-radius: 3px;
    border: 1px solid var(--border);
}

pre {
    background: var(--bg-code);
    border: 1px solid var(--border);
    border-radius: 6px;
    padding: 1.2rem;
    overflow-x: auto;
    margin-bottom: 1.5rem;
    line-height: 1.5;
}

pre code {
    background: none;
    border: none;
    padding: 0;
    font-size: 0.85rem;
    color: var(--text-primary);
}

/* === TABLES === */
.table-wrap {
    overflow-x: auto;
    margin-bottom: 1.5rem;
}

table {
    width: 100%;
    border-collapse: collapse;
    font-size: 0.9rem;
}

th, td {
    padding: 0.6rem 0.8rem;
    text-align: left;
    border: 1px solid var(--border);
    vertical-align: top;
}

th {
    background: var(--bg-tertiary);
    font-weight: 700;
    color: var(--text-primary);
    white-space: nowrap;
}

td {
    background: var(--bg-secondary);
}

tr:hover td {
    background: var(--bg-tertiary);
}

/* === HOMEPAGE === */
.hero {
    text-align: center;
    padding: 3rem 0 2rem;
    border-bottom: 1px solid var(--border);
    margin-bottom: 2rem;
}

.hero h1 {
    font-size: 2.4rem;
    border: none;
    margin-bottom: 0.5rem;
}

.hero-sub {
    font-size: 1.1rem;
    color: var(--text-secondary);
    margin-bottom: 1rem;
    font-style: italic;
}

.hero-meta {
    font-family: var(--font-mono);
    font-size: 0.8rem;
    color: var(--text-muted);
}

.intro {
    margin-bottom: 3rem;
}

.intro h2 {
    margin-top: 0;
}

/* Domain Grid */
.domain-grid-section h2 {
    margin-bottom: 1.5rem;
}

.domain-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
    gap: 0.75rem;
    margin-bottom: 3rem;
}

.domain-card {
    display: flex;
    align-items: flex-start;
    gap: 1rem;
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: 6px;
    padding: 1rem 1.2rem;
    text-decoration: none;
    transition: border-color 0.15s, background 0.15s;
}

.domain-card:hover {
    border-color: var(--accent-dim);
    background: var(--bg-tertiary);
    text-decoration: none;
}

.domain-number {
    font-family: var(--font-mono);
    font-size: 1.4rem;
    font-weight: 700;
    color: var(--accent-dim);
    min-width: 2.5rem;
    text-align: center;
    line-height: 1.3;
}

.domain-info h3 {
    font-size: 0.95rem;
    margin: 0 0 0.25rem;
    color: var(--text-primary);
}

.domain-info p {
    font-size: 0.82rem;
    color: var(--text-secondary);
    margin: 0;
    line-height: 1.4;
}

/* Document Links */
.doc-links {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-bottom: 2rem;
}

.doc-link {
    display: flex;
    align-items: center;
    gap: 1.5rem;
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: 6px;
    padding: 1rem 1.5rem;
    text-decoration: none;
    transition: border-color 0.15s, background 0.15s;
}

.doc-link:hover {
    border-color: var(--accent-dim);
    background: var(--bg-tertiary);
    text-decoration: none;
}

.doc-link strong {
    font-size: 1rem;
    min-width: 160px;
    color: var(--accent);
}

.doc-link span {
    font-size: 0.85rem;
    color: var(--text-secondary);
}

.doc-link .article-count {
    margin-left: auto;
    font-family: var(--font-mono);
    font-size: 0.78rem;
    color: var(--text-muted);
    white-space: nowrap;
}

/* Principles */
.principles-section {
    margin-bottom: 2rem;
}

.principles-section ol {
    padding-left: 1.5rem;
}

.principles-section li {
    margin-bottom: 0.75rem;
}

/* === 404 PAGE === */
.error-page {
    text-align: center;
    padding: 4rem 0;
}

.error-page h1 {
    font-family: var(--font-mono);
    font-size: 6rem;
    border: none;
    color: var(--accent-dim);
    margin-bottom: 0;
}

.error-subtitle {
    font-size: 1.3rem;
    color: var(--text-secondary);
    margin-bottom: 2rem;
}

.error-page ol {
    display: inline-block;
    text-align: left;
    margin: 1rem auto 2rem;
}

.btn {
    display: inline-block;
    padding: 0.6rem 1.5rem;
    background: var(--accent-dim);
    color: var(--text-primary);
    border-radius: 4px;
    text-decoration: none;
    font-family: var(--font-mono);
    font-size: 0.9rem;
    transition: background 0.15s;
}

.btn:hover {
    background: var(--accent);
    text-decoration: none;
}

/* === COMMENT SECTION STYLING (placeholder for future use) === */
.commentary-section {
    margin-top: 3rem;
    padding-top: 2rem;
    border-top: 2px solid var(--border);
}

.commentary-section h2 {
    color: var(--text-muted);
    font-size: 1rem;
    text-transform: uppercase;
    letter-spacing: 0.1em;
}

.comment {
    background: var(--bg-tertiary);
    border-left: 3px solid var(--accent-dim);
    padding: 1rem 1.5rem;
    margin: 1rem 0;
    border-radius: 0 4px 4px 0;
}

.comment-date {
    font-family: var(--font-mono);
    font-size: 0.8rem;
    color: var(--text-muted);
    margin-bottom: 0.5rem;
}

.comment-body {
    font-size: 0.92rem;
    color: var(--text-secondary);
}

/* === FOOTER === */
.site-footer {
    background: var(--bg-secondary);
    border-top: 1px solid var(--border);
    padding: 1.5rem 2rem;
    margin-top: auto;
}

.footer-inner {
    max-width: 1200px;
    margin: 0 auto;
    text-align: center;
}

.site-footer p {
    font-size: 0.8rem;
    color: var(--text-muted);
    font-family: var(--font-mono);
    margin-bottom: 0.25rem;
}

/* === RESPONSIVE === */
@media (max-width: 1024px) {
    .layout {
        flex-direction: column;
    }
    .sidebar {
        width: 100%;
        position: static;
        max-height: none;
        padding-right: 0;
        margin-bottom: 1rem;
        padding-bottom: 1rem;
        border-bottom: 1px solid var(--border);
    }
    .toc ul ul {
        padding-left: 0.75rem;
    }
    .content {
        max-width: 100%;
    }
}

@media (max-width: 768px) {
    .header-inner {
        padding: 0.75rem 1rem 0.25rem;
        flex-direction: column;
        gap: 0.25rem;
    }
    .site-nav {
        padding: 0 0.5rem;
    }
    .site-nav a {
        padding: 0.4rem 0.6rem;
        font-size: 0.78rem;
    }
    .layout {
        padding: 1rem;
    }
    h1 { font-size: 1.5rem; }
    h2 { font-size: 1.25rem; }
    .hero h1 { font-size: 1.7rem; }
    .domain-grid {
        grid-template-columns: 1fr;
    }
    .doc-link {
        flex-direction: column;
        gap: 0.5rem;
    }
    .doc-link .article-count {
        margin-left: 0;
    }
}

@media (max-width: 480px) {
    .site-nav ul {
        flex-wrap: wrap;
    }
}

/* === PRINT === */
@media print {
    :root {
        --bg-primary: #ffffff;
        --bg-secondary: #f8f8f8;
        --bg-tertiary: #f0f0f0;
        --bg-card: #f5f5f5;
        --bg-code: #f0f0f0;
        --text-primary: #111111;
        --text-secondary: #333333;
        --text-muted: #666666;
        --accent: #2244aa;
        --accent-dim: #4466cc;
        --border: #cccccc;
        --border-light: #dddddd;
        --link: #2244aa;
    }

    body {
        font-size: 11pt;
        line-height: 1.5;
    }

    .site-header,
    .site-nav,
    .sidebar,
    .site-footer {
        display: none !important;
    }

    .layout {
        display: block;
        max-width: 100%;
        padding: 0;
    }

    .content {
        max-width: 100%;
    }

    a {
        text-decoration: underline;
    }

    a[href]::after {
        content: " (" attr(href) ")";
        font-size: 0.8em;
        color: var(--text-muted);
    }

    a[href^="#"]::after {
        content: "";
    }

    h1, h2, h3 {
        page-break-after: avoid;
    }

    table, pre {
        page-break-inside: avoid;
    }

    .domain-card, .doc-link {
        border: 1px solid #ccc;
        break-inside: avoid;
    }

    .hero {
        border-bottom: 2px solid #000;
    }
}
'''


def extract_md_title(md_path):
    """Extract the first # heading and ## subtitle from a markdown file."""
    title = None
    subtitle = None
    with open(md_path, 'r') as f:
        for line in f:
            line = line.strip()
            if title is None and line.startswith('# ') and not line.startswith('## '):
                title = line[2:].strip()
            elif title and subtitle is None and line.startswith('## '):
                subtitle = line[3:].strip()
                break
            elif title and line and not line.startswith('*') and not line.startswith('#'):
                break  # stop if we hit content without finding subtitle
    return title, subtitle


def discover_extra_pages():
    """Auto-discover markdown files not in the hardcoded PAGES list."""
    known_md = {md_file for _, md_file, _, _ in PAGES}
    extra = []
    import glob
    for md_path in sorted(glob.glob(os.path.join(SRC_DIR, 'stage*.md'))):
        md_file = os.path.basename(md_path)
        if md_file in known_md:
            continue
        stem = md_file.replace('.md', '')
        name_part = re.sub(r'^stage\d+-', '', stem)
        out_file = f'{name_part}.html'
        # Extract real title from the markdown content
        md_title, md_subtitle = extract_md_title(md_path)
        title = md_title or name_part.replace('-', ' ').replace('_', ' ').title()
        # Build a short nav label from the name part
        nav_label = name_part.replace('-', ' ').replace('_', ' ').title()
        extra.append((out_file, md_file, title, nav_label, md_subtitle))
    return extra


def main():
    global ALL_NAV_ITEMS
    os.makedirs(SITE_DIR, exist_ok=True)

    # Write CSS
    css_path = os.path.join(SITE_DIR, 'style.css')
    with open(css_path, 'w') as f:
        f.write(build_css())
    print(f'  Wrote {css_path}')

    # Combine hardcoded pages with auto-discovered ones
    extra_pages = discover_extra_pages()
    # Hardcoded pages: (out_file, md_file, title, nav_label) - add empty subtitle
    all_pages = [(o, m, t, n, None) for o, m, t, n in PAGES] + extra_pages
    built = 0

    # Populate global nav items for nav_html and index page
    ALL_NAV_ITEMS = []
    for out_file, md_file, title, nav_label, subtitle in all_pages:
        md_path = os.path.join(SRC_DIR, md_file)
        if not os.path.exists(md_path):
            continue
        # For hardcoded pages without subtitle, extract it
        if subtitle is None:
            _, subtitle = extract_md_title(md_path)
        ALL_NAV_ITEMS.append((out_file, nav_label, subtitle or title))

    # Build content pages from markdown
    for out_file, md_file, title, nav_label, _subtitle in all_pages:
        md_path = os.path.join(SRC_DIR, md_file)
        if not os.path.exists(md_path):
            print(f'  Skipping {md_file} (not found)')
            continue
        with open(md_path, 'r') as f:
            md_content = f.read()

        converter = MarkdownConverter()
        html_content = converter.convert(md_content)
        toc = build_toc_html(converter.toc, max_level=3)
        page = page_template(title, html_content, toc, out_file)

        out_path = os.path.join(SITE_DIR, out_file)
        with open(out_path, 'w') as f:
            f.write(page)
        print(f'  Wrote {out_path} ({len(converter.toc)} headings)')
        built += 1

    # Build index
    index_path = os.path.join(SITE_DIR, 'index.html')
    with open(index_path, 'w') as f:
        f.write(build_index_page())
    print(f'  Wrote {index_path}')

    # Build 404
    four_path = os.path.join(SITE_DIR, '404.html')
    with open(four_path, 'w') as f:
        f.write(build_404_page())
    print(f'  Wrote {four_path}')

    print(f'\nBuild complete. {built + 2} pages generated in {SITE_DIR}/')


if __name__ == '__main__':
    main()
