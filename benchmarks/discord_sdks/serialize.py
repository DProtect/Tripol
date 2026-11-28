import asyncio
import json
import time

import discord


def classic():
    embed = discord.Embed(title="TRIPOL", description="Discord", color=0x00A0FF)
    embed.add_field(name="Version", value="0.28", inline=False)
    view = discord.ui.View()
    view.add_item(discord.ui.Button(label="Confirm", custom_id="confirm", style=discord.ButtonStyle.success))
    return json.dumps([embed.to_dict(), view.to_components()], separators=(",", ":"))


def modern():
    view = discord.ui.LayoutView()
    view.add_item(discord.ui.Container(discord.ui.TextDisplay("TRIPOL"), accent_colour=0x00A0FF))
    return json.dumps(
        {"flags": 32768, "components": view.to_components(), "allowed_mentions": {"parse": []}},
        separators=(",", ":"),
    )


async def main():
    print("READY", flush=True)
    count, mode = int(input()), int(input())
    work = classic if mode == 0 else modern
    start = time.perf_counter_ns()
    total = sum(len(work()) for _ in range(count))
    print(time.perf_counter_ns() - start, total)
    print(work(), flush=True)
    input()


if __name__ == "__main__":
    asyncio.run(main())
