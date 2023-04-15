namespace shortener_tg.Utils.MessageChecker;

using System.Text.RegularExpressions;
using Telegram.Bot.Types;
using Telegram.Bot.Types.Enums;
using Setup;



public class BaseUpdateChecker : IUpdateChecker
{
    private readonly IConfigurationSection _config;

    public BaseUpdateChecker(string configSection)
    {
        _config = ConfigProvider.Config.GetSection(configSection);
    }

    public bool CheckMessage(Update update)
    {
        var updateTypes = _config.GetSection("updateType").Get<string[]>();
        if (!updateTypes.Contains(update.Type.ToString()))
        {
            return false;
        }

        if (update.Type is UpdateType.CallbackQuery)
        {
            var callbackQueryData = _config.GetSection("CallbackQueryData").Get<string[]>();
            return callbackQueryData.Contains(update.CallbackQuery!.Data);
        }
        

        if (update.Message?.Text is null)
        {
            return false;
        }

        var txt = update.Message.Text!.ToLower();

        var flag = false;
        
        var keywords = _config.GetSection("KeyWords").Get<string[]>();

        if (keywords is not null)
        {
            var words = txt.Split().Select(x => x.Trim());
            flag |= keywords.Any(keyword => words.Contains(keyword.ToLower()));
        }

        var regExp = _config.GetValue<string>("regExp");
        if (regExp is not null)
        {
            var t = new Regex(regExp);
            flag |= t.Match(txt).Success;
        }

        return flag;
    }
}