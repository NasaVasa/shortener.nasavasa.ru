using Microsoft.AspNetCore.Mvc;
using Telegram.Bot;
using Telegram.Bot.Types;
using Telegram.Bot.Types.Enums;

namespace shortener_tg.Controllers;

public class TgMessageController : ControllerBase
{
    public static async Task Write(string message, ITelegramBotClient client, Update update, CancellationToken ct)
    {
        await client.SendTextMessageAsync(
            update.Message!.Chat.Id,
            message,
            parseMode: ParseMode.Html,
            cancellationToken: ct);
    }
}