using shortener_tg.Controllers;

namespace shortener_tg.Services.Help;

using Microsoft.Extensions.Configuration;
using Telegram.Bot;
using Telegram.Bot.Types;
using Telegram.Bot.Types.Enums;
using Telegram.Bot.Types.ReplyMarkups;
using Setup;
using Utils.MessageChecker;

public class HelpHandler : IUpdateHandler
{
    private static readonly IUpdateChecker UpdateChecker = new BaseUpdateChecker("Help");
    private static readonly IConfigurationSection Config = ConfigProvider.Config.GetSection("Help");

    public async Task HandelStart(ITelegramBotClient client, Update update, CancellationToken ct)
    {
        if (!Validate(update))
        {
            return;
        }
        var answer = Config["Answer"];
        await TgMessageController.Write(answer, client, update, ct);
    }

    public bool Validate(Update update)
    {
        return UpdateChecker.CheckMessage(update);
    }
}