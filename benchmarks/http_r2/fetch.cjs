const readline = require('node:readline');
const { createWriteStream } = require('node:fs');
const { stat } = require('node:fs/promises');
const { join } = require('node:path');
const { pipeline } = require('node:stream/promises');

(async () => {
    const input = readline.createInterface({ input: process.stdin })[Symbol.asyncIterator]();
    const url = (await input.next()).value;
    const count = Number((await input.next()).value);
    const mode = (await input.next()).value;
    const width = Number((await input.next()).value);
    const size = Number((await input.next()).value);
    const directory = (await input.next()).value;
    async function receive(index) {
        const response = await fetch(url, { signal: AbortSignal.timeout(10000) });
        if (response.status !== 200) throw Error('unexpected response');
        let bytes;
        if (mode === 'stream') {
            const path = join(directory, `${index}.bin`);
            await pipeline(response.body, createWriteStream(path, { flags: 'wx', flush: true }));
            bytes = (await stat(path)).size;
        } else {
            bytes = (await response.text()).length;
        }
        if (bytes !== size) throw Error('unexpected response size');
        return bytes;
    }
    console.log('READY');
    await input.next();
    const start = process.hrtime.bigint();
    let total = 0;
    const samples = [];
    for (let i = 0; i < count; i += width) {
        const began = process.hrtime.bigint();
        if (width === 1) {
            total += await receive(i);
        } else {
            const batch = await Promise.all(Array.from({ length: width }, (_, offset) => receive(i + offset)));
            total += batch.reduce((sum, bytes) => sum + bytes, 0);
        }
        samples.push(process.hrtime.bigint() - began);
    }
    console.log(String(process.hrtime.bigint() - start) + ' ' + total);
    console.log(samples.join(' '));
    console.log('END');
    await input.next();
})().catch(error => { console.error(error); process.exitCode = 1; });
