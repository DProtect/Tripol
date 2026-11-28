const readline = require('node:readline');
(async () => {
    const input = readline.createInterface({ input: process.stdin })[Symbol.asyncIterator]();
    console.log('READY');
    const url = (await input.next()).value;
    const count = Number((await input.next()).value);
    const start = process.hrtime.bigint();
    let total = 0;
    for (let i = 0; i < count; i++) {
        const response = await fetch(url, { signal: AbortSignal.timeout(10000) });
        const data = await response.text();
        if (response.status !== 200 || data.length !== 512) throw Error('unexpected response');
        total += data.length;
    }
    console.log(String(process.hrtime.bigint() - start) + ' ' + total);
    await input.next();
})().catch(error => { console.error(error); process.exitCode = 1; });
