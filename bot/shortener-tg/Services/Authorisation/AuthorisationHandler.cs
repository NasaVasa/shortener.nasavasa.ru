using System.Text.RegularExpressions;
using Newtonsoft.Json;
using shortener_tg.Controllers;

namespace shortener_tg.Services.Authorisation;

using Telegram.Bot;
using Telegram.Bot.Types;
using Microsoft.Extensions.Configuration;
using Telegram.Bot.Types.Enums;
using Setup;
using Utils.MessageChecker;

public class AuthorisationHandler : IUpdateHandler
{
    private static readonly IUpdateChecker UpdateChecker = new BaseUpdateChecker("Authorisation");
    private static readonly IConfigurationSection Config = ConfigProvider.Config.GetSection("Authorisation");

    private readonly string _apiRequestPath = (Environment.GetEnvironmentVariable("API_PATH") ??
                                               throw new InvalidOperationException(
                                                   "В переменных окружения не обнаружено API_PATH")) +
                                              Config["Path"];

    private static readonly HttpClient HttpClient = new();

    public async Task HandelStart(ITelegramBotClient client, Update update, CancellationToken ct)
    {
        if (!Validate(update))
        {
            return;
        }

        var fromId = update.Message.From.Id;
        if (Handlers.Users[fromId] is not null)
        {
            await TgMessageController.Write("Вы уже вошли", client, update, ct);
        }
        else
        {
            Handlers.Tasks[fromId] = (HandelEnterLogin, new List<string>());
            await TgMessageController.Write("Введите свой login", client, update, ct);
        }
    }

    public async Task HandelEnterLogin(ITelegramBotClient client, Update update, CancellationToken ct,
        List<string> list)
    {
        var messageText = update.Message?.Text;
        var fromId = update.Message.From.Id;
        list.Add(messageText);
        Handlers.Tasks[fromId] = (HandelEnterPassword, list);
        await TgMessageController.Write("Введите свой password", client, update, ct);
    }

    public async Task HandelEnterPassword(ITelegramBotClient client, Update update, CancellationToken ct,
        List<string> list)
    {
        var messageText = update.Message?.Text;
        var fromId = update.Message.From.Id;
        var request = new StringContent(JsonConvert.SerializeObject(new Request()
        {
            login = list[0],
            password = messageText
        }));
        var responseMessage =
            await HttpClient.PostAsync(_apiRequestPath, request, ct);
        var responseContent = await responseMessage.Content.ReadAsStringAsync(ct);
        var response = JsonConvert.DeserializeObject<Response>(responseContent);
        if (response?.error is null)
        {
            Handlers.Users[fromId] = response.token;
            await TgMessageController.Write("<b>Вы успешно вошли</b>", client, update, ct);
        }
        else
        {
            await TgMessageController.Write($"Возникла ошибка: <b>{response.error}</b>\n" +
                                            "Попробуйте заново", client, update, ct);
        }
        Handlers.Tasks[fromId] = (null, null);
    }

    public bool Validate(Update update)
    {
        return UpdateChecker.CheckMessage(update);
    }

    class Request
    {
        public string login { get; set; }
        public string password { get; set; }
    }

    class Response
    {
        public string? token { get; set; }
        public string? error { get; set; }
    }
}