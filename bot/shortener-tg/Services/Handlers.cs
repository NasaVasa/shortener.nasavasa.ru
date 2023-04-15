using System.Text.Json;
using shortener_tg.Controllers;
using shortener_tg.Services.Authorisation;
using shortener_tg.Services.CreateUrl;
using shortener_tg.Services.DeleteUrl;
using shortener_tg.Services.Exit;
using shortener_tg.Services.GetUrlKeys;
using shortener_tg.Services.Help;
using shortener_tg.Services.InfoUrl;
using shortener_tg.Services.Registration;
using shortener_tg.Setup;
using Telegram.Bot;
using Telegram.Bot.Types;
using Telegram.Bot.Types.Enums;

namespace shortener_tg.Services;

public static class Handlers
{
    public static Dictionary<long, (Func<ITelegramBotClient, Update, CancellationToken,List<string>, Task>,List<string>)> Tasks = new();
    public static Dictionary<long, string?> Users = new();

    private static readonly List<IUpdateHandler> UpdateHandlers = new()
    {
        new StartHandler(),
        new HelpHandler(),
        new RegistrationHandler(),
        new AuthorisationHandler(),
        new GetUrlKeysHandler(),
        new CreateUrlHandler(),
        new DeleteUrlHandler(),
        new InfoUrlHandler(),
        new ExitHandler(),
    };

    public static async Task HandelUpdates(ITelegramBotClient client, Update update, CancellationToken ct)
    {
        Console.WriteLine(
            $"[{DateTime.Now}] - Message {update.Message?.Text ?? update.CallbackQuery?.Data} from @{update.Message?.From?.Username ?? update.CallbackQuery?.From.Username}"
        );
        try
        {
            var fromId = update.Message.From.Id;
            if (Tasks.ContainsKey(fromId) && Tasks[fromId] is not (null,null))
            {
                var currentTask = Tasks[fromId];
                await currentTask.Item1.Invoke(client,update,ct,currentTask.Item2);
            }
            else
            {
                await Task.WhenAll(
                    from handler in UpdateHandlers
                    select handler.HandelStart(client, update, ct));
            }
            
        }
        catch (Exception e)
        {
            if (update.Message is not null)
            {
                await TgMessageController.Write("Error while processing... Ask admin for help", client, update, ct);
            }

            Console.WriteLine(ConfigProvider.Config["AdminChatID"]);
            /*await client.SendTextMessageAsync(
                ConfigProvider.Config["AdminChatID"],
                $"Error on message <code>{update.Message?.Text ?? update.CallbackQuery?.Data}</code> from @{update.Message?.From?.Username ?? update.CallbackQuery?.From.Username}\n\n" +
                $"<code>{JsonSerializer.Serialize(e.InnerException?.Data ?? e.Data)}</code>\n\n" +
                $"<code>{e}</code>",
                parseMode: ParseMode.Html,
                cancellationToken: ct
            );*/
            Console.WriteLine(e);
        }
    }

    public static Task ErrorHandler(
        ITelegramBotClient client,
        Exception ex,
        CancellationToken ct)
    {
        Console.WriteLine(ex);
        return Task.CompletedTask;
    }
}