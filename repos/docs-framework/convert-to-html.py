#!/usr/bin/env python3
"""
Convert docs-framework markdown files to semantic HTML articles.
Also ingests agent-generated HTML from html/new/.
Two-pass: scan all sources first, then generate with unified sidebar navigation.
"""

import re
import json
import markdown
from pathlib import Path
from collections import OrderedDict

INPUT_DIR = Path("/Users/tim/docs-framework")
OUTPUT_DIR = INPUT_DIR / "html"
AGENT_DIR = OUTPUT_DIR / "new"

md_converter = markdown.Markdown(extensions=['tables', 'sane_lists'])

# Domain mapping for sidebar grouping
DOMAIN_MAP = OrderedDict([
    ('1', 'Constitution & Philosophy'),
    ('2', 'Governance & Authority'),
    ('3', 'Security & Integrity'),
    ('4', 'Infrastructure & Power'),
    ('5', 'Platform & Core Systems'),
    ('6', 'Data & Archives'),
    ('7', 'Intelligence & Analysis'),
    ('8', 'Automation & Agents'),
    ('9', 'Education & Training'),
    ('10', 'User Operations'),
    ('11', 'Administration'),
    ('12', 'Disaster Recovery'),
    ('13', 'Evolution & Adaptation'),
    ('14', 'Research & Theory'),
    ('15', 'Ethics & Safeguards'),
    ('16', 'Interface & Navigation'),
    ('17', 'Scaling & Federation'),
    ('18', 'Import & Quarantine'),
    ('19', 'Quality Assurance'),
    ('20', 'Institutional Memory'),
    ('META', 'Meta-Documentation'),
    ('FW', 'Framework'),
])

# Map article ID prefixes to domain numbers
PREFIX_TO_DOMAIN = {
    'ETH': '1', 'CON': '1',
    'GOV': '2',
    'SEC': '3',
    'OPS': '10',
    'META': 'META',
}


def get_domain_for_article(article_id):
    """Determine which domain an article belongs to."""
    # Standard prefix articles
    for prefix, domain in PREFIX_TO_DOMAIN.items():
        if article_id.startswith(prefix + '-'):
            return domain

    # D-number articles (D4-001, D17-007, etc.)
    m = re.match(r'^D(\d+)-', article_id)
    if m:
        return m.group(1)

    # Domain framework entries
    if article_id.startswith('DOMAIN-'):
        return article_id.split('-')[1]

    return 'FW'


def md_to_html(text):
    md_converter.reset()
    return md_converter.convert(text)


def wrap_sections(html_content):
    parts = re.split(r'(<h2>.*?</h2>)', html_content)
    result = []
    in_section = False
    for part in parts:
        if re.match(r'<h2>', part):
            if in_section:
                result.append('</section>')
            result.append('<section>')
            result.append(part)
            in_section = True
        else:
            result.append(part)
    if in_section:
        result.append('</section>')
    return '\n'.join(result)


def extract_metadata(text):
    meta = {}
    remaining_lines = []
    past_meta = False
    for line in text.strip().split('\n'):
        stripped = line.strip()
        m = re.match(r'^\*\*(.+?):\*\*\s*(.*)$', stripped)
        if m and not past_meta:
            meta[m.group(1).strip()] = m.group(2).strip()
        elif stripped == '---' and not past_meta:
            continue
        elif stripped == '' and not past_meta:
            continue
        else:
            past_meta = True
            remaining_lines.append(line)
    return meta, '\n'.join(remaining_lines)


def metadata_to_html(meta):
    if not meta:
        return ''
    html = '<aside class="metadata">\n<dl>\n'
    for key, value in meta.items():
        html += f'  <dt>{key}</dt>\n  <dd>{value}</dd>\n'
    html += '</dl>\n</aside>'
    return html


def classify_heading(line):
    stripped = line.strip()
    if not stripped.startswith('# '):
        return None
    text = stripped[2:]

    if re.match(r'^STAGE \d+', text):
        return None
    if text.startswith('===') or text.startswith('{') or text.startswith('Select ') or text.startswith('Updated:') or text.startswith('Status for '):
        return None

    m = re.match(r'^([A-Z][A-Z0-9]*-\d+[A-Za-z]*)\s*--\s*(.+)$', text)
    if m:
        return (m.group(1), m.group(2).strip())

    m = re.match(r'^DOMAIN (\d+):\s*(.+)$', text)
    if m:
        return (f"DOMAIN-{m.group(1)}", m.group(2).strip())

    m = re.match(r'^(APPENDIX [A-Z]):\s*(.+)$', text)
    if m:
        return (m.group(1).replace(' ', '-'), m.group(2).strip())

    if re.match(r'^(CROSS-DOMAIN|UNIVERSAL|Consolidated|Cross-Domain)', text):
        clean_id = re.sub(r'[^A-Za-z0-9]+', '-', text).strip('-').upper()[:40]
        return (clean_id, text)

    m = re.match(r'^(Stage 1 Meta-Framework):\s*(.+)$', text)
    if m:
        return ('META-FRAMEWORK', m.group(2).strip())

    return None


def split_into_articles(content):
    lines = content.split('\n')
    articles = []
    preamble_lines = []
    current_id = None
    current_title = None
    current_lines = []

    for line in lines:
        stripped = line.strip()
        heading_info = None
        if stripped.startswith('# ') and not stripped.startswith('## '):
            heading_info = classify_heading(stripped)

        if heading_info:
            article_id, title = heading_info
            if current_id is not None:
                articles.append((current_id, current_title, '\n'.join(current_lines)))
            elif current_lines:
                preamble_lines = current_lines[:]
            current_id = article_id
            current_title = title
            current_lines = []
        else:
            current_lines.append(line)

    if current_id is not None:
        articles.append((current_id, current_title, '\n'.join(current_lines)))
    elif current_lines:
        preamble_lines = current_lines

    return '\n'.join(preamble_lines), articles


def scan_agent_html(agent_dir):
    """Scan agent-generated HTML files and extract article ID, title, and body."""
    agent_articles = []
    if not agent_dir.exists():
        return agent_articles

    for html_file in sorted(agent_dir.glob("*.html")):
        text = html_file.read_text()

        # Extract article body: everything inside <article ...>...</article>
        # Handle both bare <article> and full-page wrapped files
        m = re.search(r'(<article[^>]*>)(.*?)(</article>)', text, re.DOTALL)
        if not m:
            continue

        article_tag = m.group(1)
        article_inner = m.group(2)
        article_close = m.group(3)
        article_body = article_tag + article_inner + article_close

        # Extract article ID from the tag
        id_match = re.search(r'id="([^"]+)"', article_tag)
        article_id = id_match.group(1).upper() if id_match else None
        if not article_id:
            continue

        # Extract title from <h1>
        h1_match = re.search(r'<h1>(.+?)</h1>', article_inner)
        if h1_match:
            raw_title = h1_match.group(1)
            # Strip "ID &mdash; " prefix to get clean title
            title = re.sub(r'^[A-Z][A-Z0-9]*-\d+\s*[&mdash;—\-]+\s*', '', raw_title).strip()
            if not title:
                title = raw_title
        else:
            title = article_id

        agent_articles.append({
            'source': 'agent',
            'id': article_id,
            'title': title,
            'article_html': article_body,  # pre-built HTML
            'filename': html_file.stem + '.html',
            'domain': get_domain_for_article(article_id),
        })

    return agent_articles


def article_to_html(article_id, title, content):
    meta, body = extract_metadata(content)
    body_html = md_to_html(body)
    body_html = wrap_sections(body_html)
    body_html = re.sub(r'<h2>(\d+)\.\s+', r'<h2>', body_html)

    safe_id = re.sub(r'[^a-z0-9-]', '-', article_id.lower()).strip('-')
    parts = []
    parts.append(f'<article id="{safe_id}">')

    if re.match(r'^[A-Z][A-Z0-9]*-\d+', article_id):
        parts.append(f'  <h1>{article_id} &mdash; {title}</h1>')
    elif article_id.startswith('DOMAIN-'):
        parts.append(f'  <h1>Domain {article_id.split("-")[1]}: {title}</h1>')
    else:
        parts.append(f'  <h1>{title}</h1>')

    if meta:
        parts.append(metadata_to_html(meta))
    parts.append(body_html)
    parts.append('</article>')
    return '\n'.join(parts)


def safe_filename(article_id):
    return re.sub(r'[^a-z0-9-]', '-', article_id.lower()).strip('-')


def build_sidebar_html(domain_articles, current_id=None):
    """Build sidebar navigation HTML grouped by domain."""
    html = '<nav class="sidebar" id="sidebar">\n'
    html += '  <div class="sidebar-header">\n'
    html += '    <a href="index.html" class="home-link">holm.chat</a>\n'
    html += '    <button class="sidebar-toggle" onclick="document.body.classList.toggle(\'sidebar-closed\')" aria-label="Toggle menu">&times;</button>\n'
    html += '  </div>\n'
    html += '  <div class="sidebar-content">\n'

    for domain_num, domain_name in DOMAIN_MAP.items():
        if domain_num not in domain_articles:
            continue
        articles = domain_articles[domain_num]
        if not articles:
            continue

        is_current_domain = any(a['id'] == current_id for a in articles)
        open_attr = ' open' if is_current_domain else ''

        html += f'  <details{open_attr}>\n'
        html += f'    <summary>{domain_name}</summary>\n'
        html += '    <ul>\n'
        for a in articles:
            active = ' class="active"' if a['id'] == current_id else ''
            label = a['id'] if re.match(r'^[A-Z][A-Z0-9]*-\d+', a['id']) else a['title'][:30]
            html += f'      <li{active}><a href="{a["filename"]}">{label}</a></li>\n'
        html += '    </ul>\n'
        html += '  </details>\n'

    html += '  </div>\n'
    html += '</nav>\n'
    return html


SIDEBAR_CSS = """
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body { font-family: Georgia, serif; line-height: 1.6; color: #222; display: flex; min-height: 100vh; }

    /* Sidebar */
    .sidebar { width: 260px; min-width: 260px; background: #1a1a2e; color: #ccc; height: 100vh; position: fixed; top: 0; left: 0; overflow-y: auto; z-index: 100; transition: transform 0.2s; }
    .sidebar-header { padding: 1em; border-bottom: 1px solid #333; display: flex; justify-content: space-between; align-items: center; }
    .home-link { color: #fff; text-decoration: none; font-weight: bold; font-size: 1.1em; }
    .sidebar-toggle { background: none; border: none; color: #888; font-size: 1.4em; cursor: pointer; display: none; }
    .sidebar-content { padding: 0.5em 0; }
    .sidebar details { border-bottom: 1px solid #2a2a4a; }
    .sidebar summary { padding: 0.6em 1em; cursor: pointer; font-size: 0.85em; font-weight: bold; color: #aaa; text-transform: uppercase; letter-spacing: 0.03em; }
    .sidebar summary:hover { color: #fff; background: #2a2a4a; }
    .sidebar ul { list-style: none; padding: 0 0 0.4em 0; }
    .sidebar li { font-size: 0.82em; }
    .sidebar li a { display: block; padding: 0.3em 1em 0.3em 1.8em; color: #bbb; text-decoration: none; }
    .sidebar li a:hover { color: #fff; background: #2a2a4a; }
    .sidebar li.active a { color: #fff; background: #16213e; border-left: 3px solid #4a90d9; padding-left: calc(1.8em - 3px); }

    /* Main content */
    main { margin-left: 260px; flex: 1; max-width: 52em; padding: 2em 2em 4em 2em; }
    .mobile-menu-btn { display: none; position: fixed; top: 0.6em; left: 0.6em; z-index: 200; background: #1a1a2e; color: #fff; border: none; padding: 0.4em 0.7em; font-size: 1.2em; cursor: pointer; border-radius: 4px; }

    /* Typography */
    h1 { border-bottom: 2px solid #333; padding-bottom: 0.3em; margin-bottom: 0.8em; font-size: 1.6em; }
    h2 { border-bottom: 1px solid #ddd; padding-bottom: 0.2em; margin-top: 2em; margin-bottom: 0.6em; font-size: 1.25em; }
    h3 { margin-top: 1.5em; margin-bottom: 0.4em; }
    p { margin-bottom: 0.8em; }
    section { margin-bottom: 2em; }
    ul, ol { margin: 0.5em 0 0.8em 1.5em; }
    li { margin-bottom: 0.3em; }
    aside.metadata { background: #f8f8f8; border-left: 3px solid #666; padding: 0.8em 1.2em; margin: 1em 0; font-size: 0.9em; }
    aside.metadata dl { margin: 0; }
    aside.metadata dt { font-weight: bold; display: inline; }
    aside.metadata dt::after { content: ": "; }
    aside.metadata dd { display: inline; margin: 0; }
    aside.metadata dd::after { content: "\\A"; white-space: pre; }
    table { border-collapse: collapse; width: 100%; margin: 1em 0; font-size: 0.9em; }
    th, td { border: 1px solid #ccc; padding: 0.5em; text-align: left; }
    th { background: #f0f0f0; }
    pre { background: #f5f5f5; padding: 1em; overflow-x: auto; border: 1px solid #ddd; margin: 0.8em 0; }
    code { font-family: "Courier New", monospace; font-size: 0.9em; }
    blockquote { border-left: 3px solid #999; margin: 0.8em 0; padding-left: 1em; color: #555; }

    /* Mobile */
    @media (max-width: 800px) {
        .sidebar { transform: translateX(-100%); }
        .sidebar-toggle { display: block; }
        body:not(.sidebar-closed) .sidebar { transform: translateX(-100%); }
        body.sidebar-open .sidebar { transform: translateX(0); }
        main { margin-left: 0; padding: 3em 1em 2em 1em; }
        .mobile-menu-btn { display: block; }
    }
"""


def make_page_html(title, body_html, sidebar_html):
    return f"""<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{title} - holm.chat</title>
  <style>{SIDEBAR_CSS}
  </style>
</head>
<body>
<button class="mobile-menu-btn" onclick="document.body.classList.toggle('sidebar-open')">&#9776;</button>
{sidebar_html}
<main>
{body_html}
</main>
</body>
</html>"""


def generate_index_body(domain_articles):
    html = '<h1>holm.chat Documentation</h1>\n'
    html += '<p>Complete institutional documentation library. 20 domains, 5 stages.</p>\n'

    stage_groups = OrderedDict([
        ('Stage 1: Framework', ['FW']),
        ('Stage 2: Constitution', ['1', '2', '3', '10']),
        ('Stage 3: Operations', ['4', '5', '6', '7', '8', '9', '11', '12', '13', '14', '15', '16', '17', '18', '19', '20']),
        ('Stage 5: Meta', ['META']),
    ])

    # Collect all domains that have articles
    all_domains_with_articles = set(domain_articles.keys())

    for stage_label, stage_domains in stage_groups.items():
        domains_in_stage = [d for d in stage_domains if d in all_domains_with_articles]
        if not domains_in_stage:
            continue
        html += f'\n<section>\n<h2>{stage_label}</h2>\n'
        for domain_num in domains_in_stage:
            domain_name = DOMAIN_MAP.get(domain_num, f'Domain {domain_num}')
            articles = domain_articles[domain_num]
            html += f'<h3>{domain_name}</h3>\n<ul>\n'
            for a in articles:
                label = f'{a["id"]} &mdash; {a["title"]}' if re.match(r'^[A-Z][A-Z0-9]*-\d+', a['id']) else a['title']
                html += f'  <li><a href="{a["filename"]}">{label}</a></li>\n'
            html += '</ul>\n'
        html += '</section>\n'

    return html


def main():
    OUTPUT_DIR.mkdir(exist_ok=True)

    # === PASS 1a: Scan markdown files ===
    md_files = sorted(INPUT_DIR.glob("stage*.md"))
    print(f"Found {len(md_files)} markdown files.")

    print("Pass 1a: Scanning markdown articles...")
    md_articles = []

    for md_file in md_files:
        content = md_file.read_text()
        preamble, articles = split_into_articles(content)
        for article_id, title, article_content in articles:
            md_articles.append({
                'source': md_file.stem,
                'id': article_id,
                'title': title,
                'content': article_content,
                'filename': safe_filename(article_id) + '.html',
                'domain': get_domain_for_article(article_id),
            })

    print(f"  Found {len(md_articles)} markdown articles.")

    # === PASS 1b: Scan agent-generated HTML files ===
    print("Pass 1b: Scanning agent-generated HTML articles...")
    agent_articles = scan_agent_html(AGENT_DIR)
    print(f"  Found {len(agent_articles)} agent articles.")

    # === Merge: agent articles override markdown if same ID ===
    md_ids = {a['id'] for a in md_articles}
    agent_ids = {a['id'] for a in agent_articles}
    overlap = md_ids & agent_ids
    if overlap:
        print(f"  Overlap (agent overrides markdown): {len(overlap)} articles")

    # Build combined list: all markdown articles + agent articles not in markdown
    all_articles = list(md_articles)
    for a in agent_articles:
        if a['id'] not in md_ids:
            all_articles.append(a)

    # Sort by domain then ID for consistent ordering
    def sort_key(a):
        domain = a['domain']
        try:
            domain_num = int(domain)
        except ValueError:
            domain_num = 999
        # Extract numeric part of article ID for sorting
        num_match = re.search(r'-(\d+)', a['id'])
        art_num = int(num_match.group(1)) if num_match else 0
        return (domain_num, a['id'][:10], art_num)

    all_articles.sort(key=sort_key)

    # Deduplicate filenames (some generic IDs like APPENDIX-A appear in multiple source files)
    seen_filenames = {}
    for a in all_articles:
        fn = a['filename']
        if fn in seen_filenames:
            seen_filenames[fn] += 1
            base = fn.rsplit('.', 1)[0]
            a['filename'] = f"{base}-{seen_filenames[fn]}.html"
        else:
            seen_filenames[fn] = 1

    print(f"\n  Total articles: {len(all_articles)} ({len(md_articles)} markdown + {len(agent_articles) - len(overlap)} new agent)\n")

    # Build domain-grouped index for sidebar
    domain_articles = OrderedDict()
    for domain_num in DOMAIN_MAP:
        arts = [a for a in all_articles if a['domain'] == domain_num]
        if arts:
            domain_articles[domain_num] = arts

    # === PASS 2: Generate HTML with unified sidebar ===
    print("Pass 2: Generating HTML with unified sidebar...")

    for a in all_articles:
        # Agent articles already have article HTML; markdown articles need conversion
        if 'article_html' in a:
            article_html = a['article_html']
        else:
            article_html = article_to_html(a['id'], a['title'], a['content'])

        sidebar_html = build_sidebar_html(domain_articles, current_id=a['id'])
        page_title = f"{a['id']} — {a['title']}" if re.match(r'^[A-Z][A-Z0-9]*-\d+', a['id']) else a['title']
        full_html = make_page_html(page_title, article_html, sidebar_html)
        (OUTPUT_DIR / a['filename']).write_text(full_html)

    # Generate index page
    index_body = generate_index_body(domain_articles)
    index_sidebar = build_sidebar_html(domain_articles, current_id=None)
    index_html = make_page_html("Documentation Index", index_body, index_sidebar)
    (OUTPUT_DIR / "index.html").write_text(index_html)

    # Summary
    print(f"\n{'='*60}")
    print(f"Done.")
    print(f"  Articles:    {len(all_articles)}")
    print(f"  Domains:     {len(domain_articles)}")
    total_files = len(list(OUTPUT_DIR.glob("*.html")))
    print(f"  HTML files:  {total_files}")
    print(f"  Output:      {OUTPUT_DIR}")

    # Write article manifest
    manifest = [{
        'id': a['id'],
        'title': a['title'],
        'domain': a['domain'],
        'filename': a['filename'],
        'source': a['source'],
    } for a in all_articles]
    (OUTPUT_DIR / 'manifest.json').write_text(json.dumps(manifest, indent=2))
    print(f"  Manifest:    manifest.json ({len(manifest)} entries)")


if __name__ == '__main__':
    main()
