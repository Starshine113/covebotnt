package crouter

import (
	"fmt"
	"math"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/Starshine113/covebotnt/structs"
	"github.com/bwmarrin/discordgo"
)

type cmdList []*Command

func (c cmdList) Len() int      { return len(c) }
func (c cmdList) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c cmdList) Less(i, j int) bool {
	return sort.StringsAreSorted([]string{c[i].Name, c[j].Name})
}

// Invite returns an invite link for the bot
func (ctx *Ctx) Invite() string {
	// perms is the list of permissions the bot will be granted by default
	var perms = discordgo.PermissionReadMessages +
		discordgo.PermissionReadMessageHistory +
		discordgo.PermissionSendMessages +
		discordgo.PermissionManageMessages +
		discordgo.PermissionManageEmojis +
		discordgo.PermissionChangeNickname +
		discordgo.PermissionEmbedLinks +
		discordgo.PermissionAttachFiles +
		discordgo.PermissionUseExternalEmojis +
		discordgo.PermissionAddReactions

	return fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%v&permissions=%v&redirect_uri=https%%3A%%2F%%2Fstarshines.xyz%%2Fcovebot%%2Fsetup.html&scope=bot", ctx.BotUser.ID, perms)
}

// Help is the help command
func (r *Router) Help(ctx *Ctx, guildSettings *structs.GuildSettings) (err error) {
	err = ctx.TriggerTyping()
	if err != nil {
		return
	}

	if len(ctx.Args) == 0 {
		permLevel := PermLevelNone

		if err = checkOwner(ctx.Author.ID, r.BotOwners); err == nil {
			permLevel = PermLevelOwner
		} else if err = checkAdmin(ctx); err == nil {
			permLevel = PermLevelAdmin
		} else if err = checkModPerm(ctx, guildSettings); err == nil {
			permLevel = PermLevelMod
		} else if err = checkHelperPerm(ctx, guildSettings); err == nil {
			permLevel = PermLevelHelper
		}

		return r.details(ctx, permLevel)
	}

	var cmd *Command
	g := r.GetGroup(ctx.Args[0])
	if g != nil {
		if len(ctx.Args) == 1 {
			_, err = ctx.Send(ctx.GroupEmbed(g))
			return
		}
		cmd = g.GetCommand(ctx.Args[1])
		if cmd != nil {
			_, err = ctx.Send(ctx.GroupCmdEmbed(g, cmd))
			return
		}
	}
	cmd = r.GetCommand(ctx.Args[0])
	if cmd != nil {
		_, err = ctx.Send(ctx.CmdEmbed(cmd))
		return
	}

	_, err = ctx.Send(fmt.Sprintf("%v Invalid command or group provided:\n> `%v` is not a known command, group or alias.", ErrorEmoji, ctx.Args[0]))

	return
}

// GroupEmbed ...
func (ctx *Ctx) GroupEmbed(g *Group) *discordgo.MessageEmbed {
	var aliases string
	if g.Aliases == nil {
		aliases = "N/A"
	} else {
		aliases = strings.Join(g.Aliases, ", ")
	}

	var subCmds []string
	for _, cmd := range g.Subcommands {
		subCmds = append(subCmds, cmd.Name)
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("```%v```", strings.ToUpper(g.Name)),
		Description: g.Description,
		Color:       0x21a1a8,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Subcommands",
				Value:  fmt.Sprintf("```%v```", strings.Join(subCmds, "\n")),
				Inline: false,
			},
			{
				Name:   "Aliases",
				Value:  fmt.Sprintf("```%v```\n** **", aliases),
				Inline: false,
			},
			{
				Name:   "Default command",
				Value:  g.Command.Description,
				Inline: false,
			},
			{
				Name:   "Usage",
				Value:  fmt.Sprintf("```%v%v %v```", ctx.GuildPrefix, strings.ToLower(g.Command.Name), g.Command.Usage),
				Inline: false,
			},
			{
				Name:   "Permission level",
				Value:  "```" + g.Command.Permissions.String() + "```",
				Inline: false,
			},
		},
	}

	return embed
}

// GroupCmdEmbed ...
func (ctx *Ctx) GroupCmdEmbed(g *Group, cmd *Command) *discordgo.MessageEmbed {
	var aliases string

	if cmd.Aliases == nil {
		aliases = "N/A"
	} else {
		aliases = strings.Join(cmd.Aliases, ", ")
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("```%v %v```", strings.ToUpper(g.Name), strings.ToUpper(cmd.Name)),
		Description: cmd.Description,
		Color:       0x21a1a8,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Usage",
				Value:  fmt.Sprintf("```%v%v %v %v```", ctx.GuildPrefix, strings.ToLower(g.Name), strings.ToLower(cmd.Name), cmd.Usage),
				Inline: false,
			},
			{
				Name:   "Aliases",
				Value:  fmt.Sprintf("```%v```", aliases),
				Inline: false,
			},
			{
				Name:   "Permission level",
				Value:  "```" + cmd.Permissions.String() + "```",
				Inline: false,
			},
		},
	}

	return embed
}

// CmdEmbed ...
func (ctx *Ctx) CmdEmbed(cmd *Command) *discordgo.MessageEmbed {
	var aliases string

	if cmd.Aliases == nil {
		aliases = "N/A"
	} else {
		aliases = strings.Join(cmd.Aliases, ", ")
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("```%v```", strings.ToUpper(cmd.Name)),
		Description: cmd.Description,
		Color:       0x21a1a8,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Usage",
				Value:  fmt.Sprintf("```%v%v %v```", ctx.GuildPrefix, strings.ToLower(cmd.Name), cmd.Usage),
				Inline: false,
			},
			{
				Name:   "Aliases",
				Value:  fmt.Sprintf("```%v```", aliases),
				Inline: false,
			},
			{
				Name:   "Permission level",
				Value:  "```" + cmd.Permissions.String() + "```",
				Inline: false,
			},
		},
	}

	return embed
}

func (r *Router) details(ctx *Ctx, p PermLevel) (err error) {
	if err = ctx.TriggerTyping(); err != nil {
		return err
	}

	var cmds cmdList
	for _, c := range r.Commands {
		if c.Permissions <= p {
			cmds = append(cmds, c)
		}
	}

	for _, g := range r.Groups {
		if g.Command.Permissions <= p {
			cmd := *g.Command
			cmd.Name = fmt.Sprintf("%v %v", g.Name, g.Command.Name)
			cmds = append(cmds, &cmd)
		}
		for _, c := range g.Subcommands {
			if c.Permissions <= p {
				cmd := *c
				cmd.Name = fmt.Sprintf("%v %v", g.Name, c.Name)
				cmds = append(cmds, &cmd)
			}
		}
	}

	sort.Sort(cmds)
	cmdSlices := make([][]*Command, 0)

	for i := 0; i < len(cmds); i += 5 {
		end := i + 5

		if end > len(cmds) {
			end = len(cmds)
		}

		cmdSlices = append(cmdSlices, cmds[i:end])
	}

	embeds := make([]*discordgo.MessageEmbed, 0)

	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "Setup",
			Value:  "For help setting up the bot, check out the [docs](https://starshines.xyz/covebot/setup.html).",
			Inline: false,
		},
		{
			Name:   "Source code",
			Value:  "CoveBotn't is licensed under the GNU AGPLv3. The source code can be found [here](https://github.com/Starshine113/covebotnt).",
			Inline: false,
		},
	}

	startTime := ctx.AdditionalParams["startTime"].(time.Time).UTC()
	c := ctx.AdditionalParams["config"].(structs.BotConfig)
	botAuthor, err := ctx.Session.User("694563574386786314")
	if err != nil {
		return err
	}

	if c.Bot.Invite != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Invite",
			Value:  "Invite the bot with [this](" + ctx.Invite() + ") link.",
			Inline: false,
		})
	}

	embeds = append(embeds, &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    ctx.BotUser.Username + " help",
			IconURL: ctx.BotUser.AvatarURL("128"),
		},
		Description: ctx.BotUser.Username + " is a general purpose bot, with a gatekeeper, moderation commands, and starboard functionality.",
		Color:       0x21a1a8,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Page 1/%v | Made with discordgo %v", len(cmdSlices)+1, discordgo.VERSION),
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Fields: append(fields, []*discordgo.MessageEmbedField{
			{
				Name:   "Bot version",
				Value:  fmt.Sprintf("v%v-%v", ctx.AdditionalParams["botVer"].(string), ctx.AdditionalParams["gitVer"].(string)),
				Inline: true,
			},
			{
				Name:   "Go version",
				Value:  runtime.Version(),
				Inline: true,
			},
			{
				Name:   "discordgo version",
				Value:  discordgo.VERSION,
				Inline: true,
			},
			{
				Name:   "Author",
				Value:  botAuthor.Mention() + " / " + botAuthor.String(),
				Inline: false,
			},
			{
				Name:   "Uptime",
				Value:  fmt.Sprintf("Up %v\n(Since %v)", PrettyDurationString(time.Since(startTime)), startTime.Format("Jan _2 2006, 15:04:05 MST")),
				Inline: false,
			},
			{
				Name:   "** **",
				Value:  "Use ⬅️ ➡️ to navigate between pages, the numbers to choose a command, and ❌ to delete this message.",
				Inline: false,
			},
		}...)})

	for i, c := range cmdSlices {
		embeds = append(embeds, ctx.detailEmbed(i, len(cmdSlices)+1, c))
	}

	msg, err := ctx.PagedEmbed(embeds)

	ctx.AdditionalParams["cmdList"] = cmdSlices

	for i, e := range emoji {
		emoji := e
		if err = ctx.Session.MessageReactionAdd(ctx.Channel.ID, msg.ID, emoji); err != nil {
			return
		}

		index := i
		ctx.AddReactionHandler(msg.ID, e, func(ctx *Ctx) {
			page := ctx.AdditionalParams["page"].(int) - 1
			if page == -1 {
				ctx.Session.MessageReactionRemove(ctx.Channel.ID, msg.ID, emoji, ctx.Author.ID)
				return
			}
			cmdList := ctx.AdditionalParams["cmdList"].([][]*Command)

			cmdSlice := cmdList[page]
			if index >= len(cmdSlice) {
				return
			}

			ctx.Session.ChannelMessageDelete(ctx.Channel.ID, msg.ID)
			ctx.Send(ctx.CmdEmbed(cmdSlice[index]))
		})
	}

	return
}

var emoji = []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣"}

func (ctx *Ctx) detailEmbed(i, pages int, cmds []*Command) *discordgo.MessageEmbed {
	if len(cmds) > 5 {
		return nil
	}

	commands := make([]string, 0)

	for i, cmd := range cmds {
		commands = append(commands, fmt.Sprintf("%v `%v`: %v", emoji[i], cmd.Name, cmd.Description))
	}
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    ctx.BotUser.Username + " help",
			IconURL: ctx.BotUser.AvatarURL("128"),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Page %v/%v", i+2, pages),
		},
		Timestamp:   time.Now().Format(time.RFC3339),
		Description: strings.Join(commands, "\n\n"),
		Fields: []*discordgo.MessageEmbedField{{
			Name:   "Usage",
			Value:  "Use ⬅️ ➡️ to navigate between pages, the numbers to choose a command, and ❌ to delete this message.",
			Inline: false,
		}},
		Color: 0x21a1a8,
	}
	return embed
}

// PrettyDurationString ...
func PrettyDurationString(duration time.Duration) (out string) {
	var days, hours, hoursFrac, minutes float64

	hours = duration.Hours()
	hours, hoursFrac = math.Modf(hours)
	minutes = hoursFrac * 60

	hoursFrac = math.Mod(hours, 24)
	days = (hours - hoursFrac) / 24
	hours = hours - (days * 24)
	minutes = minutes - math.Mod(minutes, 1)

	if days != 0 {
		out += fmt.Sprintf("%v days, ", days)
	}
	if hours != 0 {
		out += fmt.Sprintf("%v hours, ", hours)
	}
	out += fmt.Sprintf("%v minutes", minutes)

	return
}
