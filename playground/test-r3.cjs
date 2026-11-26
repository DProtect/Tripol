const assert=require('node:assert/strict');const {run}=require('./runtime-r2.1.js');const examples=require('./examples-r3.js');const editor=require('./editor-r2.1.js');assert.equal(examples.length,100);for(const [name,es,source,want,input] of examples){assert.equal(run(source,input).output,want,name);console.log('PASS',name)}
const errors=[["show 9223372036854775807 + 1","overflow"],["show 1 / 0","arithmetic"],["a=ints(1);show a[1]","bounds"],["while true {}","limit"],["const x=1;x=2","constant"],["import tripol.math\nshow 1","Imports"],["show missing","Unknown variable"],["show (1+","Expected expression"],["show f\"{\"","Unclosed interpolation"],["show f\"{missing}\"","Unknown variable"],["for n in 1 {show n}","Expected array"],["for i,i in [1] {show i}","Duplicate loop variable"]];for(const [s,r] of errors)assert.throws(()=>run(s),new RegExp(r));
assert.equal(run('show f"{{x}} {true} {1+2}"').output,'{x} true 3\n');
assert.equal(editor.tokenize('show "<script>" // hi').map(t=>t.v).join(''),'show "<script>" // hi');
assert.equal(editor.tokenize('// show')[0].kind,'comment');
assert.equal(editor.tokenize('"# <img>"')[0].kind,'string');
assert.equal(editor.tokenize('show')[0].kind,'keyword');
assert.equal(editor.tokenize('input')[0].kind,'function');
assert.ok(editor.suggestions('sho',3).items.includes('show'));
assert.ok(editor.suggestions('score = 4\nsc',12).items.includes('score'));
assert.equal(editor.format('if true {\nshow 1\n} else {\nshow 2\n}'),'if true {\n    show 1\n} else {\n    show 2\n}');
assert.equal(editor.format('show "{" // {\nshow 2'),'show "{" // {\nshow 2');
assert.deepEqual(editor.location('a\nbc',4),{line:2,column:3});
assert.deepEqual(editor.lineRange('a\nbc\nz',2),[2,4]);
console.log('PASS 12 error cases, interpolation escape and 11 editor checks');
