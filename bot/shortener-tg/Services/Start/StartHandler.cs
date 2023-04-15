using shortener_tg.Controllers;

namespace shortener_tg.Services.Help;

using Microsoft.Extensions.Configuration;
using Telegram.Bot;
using Telegram.Bot.Types;
using Telegram.Bot.Types.Enums;
using Telegram.Bot.Types.ReplyMarkups;
using Setup;
using Utils.MessageChecker;

public class StartHandler : IUpdateHandler
{
    private static readonly IUpdateChecker UpdateChecker = new BaseUpdateChecker("Start");
    private static readonly IConfigurationSection Config = ConfigProvider.Config.GetSection("Start");

    public async Task HandelStart(ITelegramBotClient client, Update update, CancellationToken ct)
    {
        if (!Validate(update))
        {
            return;
        }

        var fromId = update.Message.From.Id;
        var answer = Config["Answer"];
        Handlers.Users.Add(fromId,null);
        Handlers.Tasks.Add(fromId,(null,null));
        await TgMessageController.Write(answer, client, update, ct);
    }

    public bool Validate(Update update)
    {
        return UpdateChecker.CheckMessage(update);
    }
}