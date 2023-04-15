namespace shortener_tg.Controllers;
using Microsoft.AspNetCore.Mvc;
using Telegram.Bot;
using Telegram.Bot.Types;
using Services;



[ApiController]
public class BotController : ControllerBase
{
    private static readonly ITelegramBotClient Client =
        new TelegramBotClient(Environment.GetEnvironmentVariable("BOT_TOKEN") ?? "");

    [HttpPost("tg-api")]
    public async Task<IActionResult> Post([FromBody] Update update, CancellationToken cancellationToken)
    {
        await Handlers.HandelUpdates(Client, update, cancellationToken);
        return Ok();
    }
}