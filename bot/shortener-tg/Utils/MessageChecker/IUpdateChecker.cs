using Telegram.Bot.Types;

namespace shortener_tg.Utils.MessageChecker;

public interface IUpdateChecker
{
    public bool CheckMessage(Update update);
}