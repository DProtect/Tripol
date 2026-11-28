const readline = require('node:readline');
const { EmbedBuilder, ActionRowBuilder, ButtonBuilder, ButtonStyle, ContainerBuilder, TextDisplayBuilder } = require('discord.js');

function classic() {
    const embed = new EmbedBuilder().setTitle('TRIPOL').setDescription('Discord').setColor(0x00a0ff)
        .addFields({ name: 'Version', value: '0.28', inline: false });
    const row = new ActionRowBuilder().addComponents(
        new ButtonBuilder().setCustomId('confirm').setLabel('Confirm').setStyle(ButtonStyle.Success),
    );
    return JSON.stringify([embed.toJSON(), [row.toJSON()]]);
}

function modern() {
    const container = new ContainerBuilder().setAccentColor(0x00a0ff).addTextDisplayComponents(
        new TextDisplayBuilder().setContent('TRIPOL'),
    );
    return JSON.stringify({ flags: 32768, components: [container.toJSON()], allowed_mentions: { parse: [] } });
}

(async () => {
    const input = readline.createInterface({ input: process.stdin })[Symbol.asyncIterator]();
    console.log('READY');
    const count = Number((await input.next()).value);
    const mode = Number((await input.next()).value);
    let sum = 0;
    const work = mode === 0 ? classic : modern;
    const start = process.hrtime.bigint();
    for (let i = 0; i < count; i++) sum += work().length;
    console.log(String(process.hrtime.bigint() - start) + ' ' + sum);
    console.log(work());
    await input.next();
})().catch(error => { console.error(error); process.exitCode = 1; });
