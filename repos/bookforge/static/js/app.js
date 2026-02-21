console.log('BookForge loaded');
fetch('/health').then(r => r.json()).then(d => console.log('Health:', d));
