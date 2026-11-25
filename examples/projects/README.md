# Practical projects / Proyectos reales

From the repository root, use an installed TRIPOL compiler:

- `tripol run examples/projects/expense_report.tri` — integer COP expense summary. Expected total 46000, maximum 21000, four entries.
- `tripol run examples/projects/notes.tri` — interactive local notes; appends to `notes.txt` in the working directory. Use an isolated directory if you have an existing file of that name.
- `tripol run examples/origin_bot.tri` — Windows Discord bot, configured using the environment as described in [Discord documentation](../../docs/discord.md). Never commit credentials.
- `tripol run examples/nova_game.tri` — see the [graphics documentation](../../docs/graphics.md) for platform restrictions.

These are small working-project candidates; automated build/smoke checks are proposed in the ecosystem workflow, not a claim that they have already passed.

- `tripol run examples/projects/collections.tri` — reusable undo stack and FIFO work queue (EVOLUTION).
