using System.Text.RegularExpressions;
using Newtonsoft.Json;
using shortener_tg.Controllers;

namespace shortener_tg.Services.Exit;

using Telegram.Bot;
using Telegram.Bot.Types;
using Microsoft.Extensions.Configuration;
using Telegram.Bot.Types.Enums;
using Setup;
using Utils.MessageChecker;

public class ExitHandler : IUpdateHandler
{
    private static readonly IUpdateChecker UpdateChecker = new BaseUpdateChecker("Exit");
    private static readonly IConfigurationSection Config = ConfigProvider.Config.GetSection("Exit");
    public async Task HandelStart(ITelegramBotClient client, Update update, CancellationToken ct)
    {
        if (!Validate(update))
        {
            return;
        }

        var fromId = update.Message.From.Id;
        if (Handlers.Users[fromId] is not null)
        {
            Handlers.Users[fromId] = null;
            await TgMessageController.Write("<b>Вы успешно вышли</b>", client, update, ct);
        }
        else
        {
            
            await TgMessageController.Write("<b>Чтобы выйти необходимо быть авторизированным</b>", client, update, ct);
        }
        
    }


    public bool Validate(Update update)
    {
        return UpdateChecker.CheckMessage(update);
    }
    
}