'use strict';

const assert = require('node:assert/strict');

function publicSignatures(source, module) {
  const signatures = new Map();
  for (const match of source.matchAll(/^func\s+(\w+(?:\.\w+)?)\s*\(/gm)) {
    const name = match[1];
    if (name.split('.').at(-1).startsWith('__')) continue;
    const start = match.index + match[0].length;
    let cursor = start;
    let depth = 1;
    for (; cursor < source.length && depth > 0; cursor++) {
      if (source[cursor] === '"') {
        cursor++;
        while (cursor < source.length && source[cursor] !== '"') {
          cursor += source[cursor] === '\\' ? 2 : 1;
        }
      } else if (source[cursor] === '(') depth++;
      else if (source[cursor] === ')') depth--;
    }
    assert.equal(depth, 0, `unterminated signature ${module}.${name}`);
    let parameters = source.slice(start, cursor - 1);
    const returns = source.slice(cursor, source.indexOf('{', cursor)).trim();
    if (name.includes('.')) parameters = parameters.replace(/^\s*[^,]+,?\s*/, '');
    const key = module + '.' + name;
    const signature = `func ${key}(${parameters}) ${returns}`.trimEnd();
    if (!signatures.has(key)) signatures.set(key, []);
    signatures.get(key).push(signature);
  }
  return signatures;
}

function normalized(signatures) {
  return signatures.trim().split(/\r?\n/).filter(line => line.trim()).map(
    line => line.replace(/"(?:\\.|[^"\\])*"|\s+/g, part => part.startsWith('"') ? part : ''),
  ).sort();
}

function referenceSignatures(markdown) {
  const output = new Map();
  for (const section of markdown.split(/^#{3,4} +/m).slice(1)) {
    const [heading] = section.split(/\r?\n/, 1);
    const escaped = heading.trim();
    if (!/^\w+\.\w+(?:\.\w+)?$/.test(escaped)) continue;
    const fenced = section.match(/^```(?:text|tripol)\r?\n(func [\s\S]*?)^```/m);
    const inline = section.match(/^`(func [^`]+)`/m);
    if (fenced || inline) output.set(escaped, (fenced || inline)[1]);
  }
  return output;
}

function htmlSignatures(html, module) {
  const output = new Map();
  const entities = { '&amp;': '&', '&quot;': '"', '&#x27;': "'", '&#39;': "'", '&lt;': '<', '&gt;': '>' };
  for (const card of html.matchAll(/<section class="function-card" id="([^"]+)">([\s\S]*?)<\/section>/g)) {
    const signature = card[2].match(/<pre class="signature"><code>([\s\S]*?)<\/code><\/pre>/);
    if (signature) output.set(module + '.' + card[1], signature[1].replace(/&(?:amp|quot|lt|gt|#39|#x27);/g, key => entities[key]));
  }
  return output;
}

function verifySignature(expected, documented, name, source) {
  assert.ok(documented, `missing ${source} signature for ${name}`);
  assert.deepEqual(normalized(documented), normalized(expected.join('\n')), `${source} signature mismatch for ${name}`);
}

module.exports = { publicSignatures, referenceSignatures, htmlSignatures, verifySignature };
