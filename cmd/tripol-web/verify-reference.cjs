'use strict';

const assert = require('node:assert/strict');
const fs = require('node:fs');
const cp = require('node:child_process');
const os = require('node:os');
const path = require('node:path');
const vm = require('node:vm');
require('./wasm_exec.js');
const host = require('./host-r3.js');

const [catalogPath, compilerPath, examplesPath] = process.argv.slice(2);
const compiler = path.resolve(compilerPath);
const environment = { ...process.env, TRIPOL_NO_UPDATE: '1' };

function native(args, cwd, options = {}) {
  return cp.spawnSync(compiler, args, {
    encoding: 'utf8', cwd, env: environment, timeout: 10000, ...options,
  });
}

function error(result) {
  return String(result.stderr || result.error || result.signal || '').trim().slice(0, 500);
}

async function main() {
  const context = vm.createContext({});
  if (examplesPath) {
    vm.runInContext(fs.readFileSync(examplesPath, 'utf8'), context, { timeout: 2000 });
  }
  vm.runInContext(fs.readFileSync(catalogPath, 'utf8'), context, { timeout: 2000 });
  const catalog = context.TripolLibraries;
  const go = new Go();
  const { instance } = await WebAssembly.instantiate(
    fs.readFileSync(path.join(__dirname, 'tripol-compiler.wasm')), go.importObject,
  );
  go.run(instance);
  const dir = fs.mkdtempSync(path.join(os.tmpdir(), 'tripol-reference-'));
  const rows = [];
  const methods = [];
  const examples = [];
  function web(row, source, input, expected) {
    try {
      const result = tripolCompileWeb(source);
      if (!result.ok) throw Error(result.error);
      assert.ok(WebAssembly.validate(result.wasm), 'invalid WebAssembly');
      const actual = host.execute(result.wasm, { input: input || '' });
      row.web = 'executed';
      row.output = actual.output;
      if (expected != null && actual.output !== expected) row.web = 'output-mismatch';
    } catch (err) {
      row.web = 'error';
      row.webDiagnostic = String(err.message);
    }
  }
  try {
    const inventory = native(['libs'], dir);
    assert.equal(inventory.status, 0, error(inventory));
    const expectedCounts = new Map([...inventory.stdout.matchAll(/^tripol\.(\w+) \((\d+) functions\)$/gm)]
      .map(([, name, count]) => [name, Number(count)]));
    assert.equal(catalog.length, expectedCounts.size, 'catalog modules and compiler inventory differ');
    for (const lib of catalog) {
      assert.equal(lib.functions.length, expectedCounts.get(lib.name), lib.name);
      const module = native(['docs', lib.name], dir);
      assert.equal(module.status, 0, error(module));
      const names = [...module.stdout.matchAll(/^func\s+(\w+)\s*\(/gm)].map(m => m[1]).sort();
      assert.equal(JSON.stringify(lib.functions.map(f => f.name).sort()), JSON.stringify(names), lib.name);
      const methodNames = [...module.stdout.matchAll(/^func\s+(\w+\.\w+)\s*\(/gm)].map(m => m[1]).sort();
      const documentedMethods = Array.from(lib.methods || [], f => f.receiver + '.' + f.name).sort();
      assert.deepEqual(documentedMethods, methodNames, lib.name);
      for (const f of [...lib.functions, ...(lib.methods || [])]) {
        const name = [lib.name, f.receiver, f.name].filter(Boolean).join('.');
        assert.ok(f.en && f.es && f.signature && f.source, `missing documentation: ${name}`);
        const row = { name };
        const caseDir = fs.mkdtempSync(path.join(dir, 'case-'));
        const source = path.join(caseDir, 'example.tri');
        const output = path.join(caseDir, 'program');
        fs.writeFileSync(source, f.source);
        const target = f.native && f.target !== 'native' ? 'windows' : 'linux';
        const built = native(['build', source, '--target', target, '-o', output], caseDir);
        row.native = built.status === 0 ? 'compiled' : 'compile-error';
        if (built.status !== 0) row.nativeDiagnostic = error(built);
        if (f.native) row.web = 'native-only';
        else web(row, f.source, f.input, f.expected);
        (f.receiver ? methods : rows).push(row);
      }
    }
    for (const [index, example] of (context.TripolExamples || []).entries()) {
      const [name, , source, expected, input] = example;
      const row = { index, name };
      if (context.TripolExampleMeta[index].native) row.web = 'native-only';
      else web(row, source, input, expected);
      examples.push(row);
    }
    const allRows = [...rows, ...methods];
    const summary = allRows.reduce((counts, row) => {
      for (const [kind, status] of Object.entries({ native: row.native, web: row.web })) {
        const key = kind + ' ' + status;
        counts[key] = (counts[key] || 0) + 1;
      }
      return counts;
    }, {});
    fs.writeFileSync(path.join(__dirname, 'reference-report.json'), JSON.stringify({
      modules: catalog.length, functions: rows.length, methods: methods.length, summary, rows, methodRows: methods, examples,
    }, null, 2));
    console.log('REFERENCE', `${rows.length} functions, ${methods.length} methods`, summary);
    const webOK = row => ['executed', 'native-only'].includes(row.web);
    console.log('PLAYGROUND', examples.filter(webOK).length + '/' + examples.length);
    return allRows.some(row => row.native !== 'compiled' || !webOK(row)) || !examples.every(webOK) ? 1 : 0;
  } finally {
    fs.rmSync(dir, { recursive: true, force: true });
  }
}

main().then(code => process.exit(code), err => { console.error(err); process.exit(1); });
