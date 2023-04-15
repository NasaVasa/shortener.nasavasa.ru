using Telegram.Bot;
using Telegram.Bot.Types;

namespace shortener_tg.Services;

public interface IUpdateHandler
{
    public Task HandelStart(ITelegramBotClient client, Update update, CancellationToken ct);
    public  bool Validate(Update update);
}