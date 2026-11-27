package core

import (
	"encoding/json"
	"testing"
)

const memberJoined = `{"op":0,"t":"GUILD_MEMBER_ADD","s":9,"d":{"guild_id":"321","user":{"id":"18446744073709551614"}}}`
const roleChanged = `{"op":0,"t":"GUILD_ROLE_UPDATE","s":10,"d":{"guild_id":"321","role":{"id":"812","name":"Maintainer"}}}`

const eventRegistry = `
func joined(event:discord.GatewayEvent)->bool {
 show event.name,event.sequence,event.guild_id,event.user_id
 return true
}
func changed(event:discord.GatewayEvent)->bool {
 show event.name,event.guild_id,event.field("role")
 return true
}
registry=discord.events()
registry.intent("GUILD_MEMBERS")
registry.on("member_join",&joined)
registry.on("GUILD_ROLE_UPDATE",&changed)
`

func TestDiscordGatewayEventsAndIntents(t *testing.T) {
	member, _ := json.Marshal(memberJoined)
	role, _ := json.Marshal(roleChanged)
	source := "import tripol.discord\n" + eventRegistry + `
show registry.intent_mask()
show registry.preview(` + string(member) + `)
show registry.preview(` + string(role) + `)
show registry.preview("{\"op\":1,\"t\":\"GUILD_ROLE_UPDATE\",\"d\":{}}")
registry.on("GUILD_MEMBER_ADD",&discord.ignore_gateway_event)
show registry.preview(` + string(member) + `)
registry.intent("GUILD_MEMBERS")
registry.intent("DIRECT_MESSAGE_POLLS")
show registry.intent_mask()
`
	runCase(t, source, "3\nGUILD_MEMBER_ADD 9 321 18446744073709551614\ntrue\nGUILD_ROLE_UPDATE 321 {\"id\":\"812\",\"name\":\"Maintainer\"}\ntrue\nfalse\nfalse\n33554435\n", false)
	if _, err := CompileWeb(source); err != nil {
		t.Fatal(err)
	}
	negative := `import tripol.discord
registry=discord.events()
try {registry.on("member_join",&discord.ignore_gateway_event)} catch code {show code}
try {registry.on("PRESENCE_UPDATE",&discord.ignore_gateway_event)} catch code {show code}
try {registry.on("INTERACTION_CREATE",&discord.ignore_gateway_event)} catch code {show code}
try {registry.on("VOICE_SERVER_UPDATE",&discord.ignore_gateway_event)} catch code {show code}
try {registry.intent("unknown")} catch code {show code}
show registry.intent_mask()
`
	runCase(t, negative, "12\n12\n12\n12\n12\n1\n", false)
}
