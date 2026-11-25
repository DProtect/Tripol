package core

import "testing"

func TestOriginAutomaticChannels(t *testing.T) {
	src := `import tripol.discord
import tripol.memory
 discord.login("dummy.test")
 show discord.configured(),discord.has_channel()
 discord.reply("!ping","Pong!")
 a="{\"op\":0,\"t\":\"MESSAGE_CREATE\",\"d\":{\"channel_id\":\"123\",\"content\":\"!ping\",\"author\":{}}}"
 b="{\"op\":0,\"t\":\"MESSAGE_CREATE\",\"d\":{\"channel_id\":\"456\",\"content\":\"!ping\",\"author\":{}}}"
 bad="{\"op\":0,\"t\":\"MESSAGE_CREATE\",\"d\":{\"channel_id\":\"12/3\",\"content\":\"!ping\",\"author\":{}}}"
 bot="{\"op\":0,\"t\":\"MESSAGE_CREATE\",\"d\":{\"channel_id\":\"456\",\"content\":\"!ping\",\"author\":{\"bot\":true}}}"
 show discord.event_channel(a),discord.gateway_reply(a)
 show discord.event_channel(b),discord.gateway_reply(b)
 show discord.event_channel(bad)=="",discord.gateway_reply(bad)=="",discord.gateway_reply(bot)==""
 discord.channel("123")
 show discord.has_channel(),discord.gateway_reply(b)=="",discord.gateway_reply(a)
 discord.all_channels()
 show discord.has_channel(),discord.gateway_reply(b)
 before=memory.used()
 repeat 2000 {scratch {assert discord.event_channel(a)=="123";assert discord.gateway_reply(b)=="Pong!"}}
 show before==memory.used()
 try {discord.say("hello")} catch code {show code}
 try {discord.run_rest()} catch code {show code}
 discord.welcome("Hi")
 try {discord.run()} catch code {show code}
 discord.reset()
 show discord.configured(),discord.has_channel()`
	want := "true false\n123 Pong!\n456 Pong!\ntrue true true\ntrue true Pong!\nfalse Pong!\ntrue\n12\n12\n12\nfalse false\n"
	for _, p := range []string{"nano", "balanced", "turbo", "debug"} {
		t.Run(p, func(t *testing.T) {
			got, bad := hyperRun(t, src, p)
			if bad || got != want {
				t.Fatalf("%q %v", got, bad)
			}
		})
	}
}

func TestOriginBinary(t *testing.T) {
	src := `import tripol.binary
 show binary.u16_le(4660),binary.u16_be(4660),binary.u32_le(4294967295)
 a=[52,18];b=[18,52];c=[255,255,255,255]
 show binary.read_u16_le(a,0),binary.read_u16_be(b,0),binary.read_u32_le(c,0)
 ok=[0,255];bad=[256];short=[0,1];wrong=[0,256]
 show binary.valid(ok),binary.valid(bad),binary.valid(ints(0))
 try {binary.read_u32_le(short,0)} catch code {show code}
 try {binary.read_u16_le(wrong,0)} catch code {show code}
 try {binary.u16_le(-1)} catch code {show code}`
	runCase(t, src, "[52, 18] [18, 52] [255, 255, 255, 255]\n4660 4660 4294967295\ntrue false true\n12\n12\n12\n", false)
}
