using System.Text.RegularExpressions;
using Newtonsoft.Json;
using shortener_tg.Controllers;

namespace shortener_tg.Services.CreateUrl;

using Telegram.Bot;
using Telegram.Bot.Types;
using Microsoft.Extensions.Configuration;
using Telegram.Bot.Types.Enums;
using Telegram.Bot.Types.ReplyMarkups;
using Setup;
using Utils.MessageChecker;

public class CreateUrlHandler : IUpdateHandler
{
    private static readonly IUpdateChecker UpdateChecker = new BaseUpdateChecker("CreateUrl");
    private static readonly IConfigurationSection Config = ConfigProvider.Config.GetSection("CreateUrl");

    private readonly string _apiRequestPath = (Environment.GetEnvironmentVariable("API_PATH") ??
                                               throw new InvalidOperationException(
                                                   "В переменных окружения не обнаружено API_PATH")) +
                                              Config["Path"];

    private readonly HttpClient _httpClient = new();

    public async Task HandelStart(ITelegramBotClient client, Update update, CancellationToken ct)
    {
        if (!Validate(update))
        {
            return;
        }

        var fromId = update.Message.From.Id;
        Handlers.Tasks[fromId] = (HandelEnterUrl, new List<string>());
        await TgMessageController.Write("Введите ссылку которую хотите укоротить", client, update, ct);
    }

    public async Task HandelEnterUrl(ITelegramBotClient client, Update update, CancellationToken ct, List<string> list)
    {
        var messageText = update.Message?.Text;
        var fromId = update.Message.From.Id;
        var regexWithProtocol = new Regex(Config.GetSection("Regex")["WithProtocol"]);
        var regexWithoutProtocol = new Regex(Config.GetSection("Regex")["WithoutProtocol"]);
        if (messageText != null && regexWithoutProtocol.IsMatch(messageText))
        {
            messageText = "https://" + messageText;
        }

        if (!regexWithProtocol.IsMatch(messageText))
        {
            Handlers.Tasks[fromId] = (null, null);
            await TgMessageController.Write("<b>Ссылка невалидна</b>", client, update, ct);
            return;
        }

        list.Add(messageText);
        if (Handlers.Users[fromId] is not null)
        {
            Handlers.Tasks[fromId] = (HandelEnterShortUrl, list);
            await TgMessageController.Write("Введите короткую ссылку, которую вы хотите получить\n" +
                                            $"Она будет иметь вид {Environment.GetEnvironmentVariable("DOMAIN")}{{КОРОТКАЯССЫЛКА}}",
                client, update, ct);
        }
        else
        {
            var request = new StringContent(JsonConvert.SerializeObject(new Request
            {
                url = list[0]
            }));
            Console.WriteLine(_apiRequestPath);
            Console.WriteLine(list[0]);
            var responseMessage =
                await _httpClient.PostAsync(_apiRequestPath, request, ct);
            var responseContent = await responseMessage.Content.ReadAsStringAsync(ct);
            var response = JsonConvert.DeserializeObject<Response>(responseContent);
            if (response?.error == null)
            {
                await TgMessageController.Write($"<b>Новая ссылка</b>: {response.short_url}\n" +
                                                $"<b>Ваш секреный ключ к данной ссылке</b>: <code>{response.secret_key}</code>",
                    client, update, ct);
            }
            else
            {
                await TgMessageController.Write($"Возникла ошибка: <b>{response.error}</b>\n" +
                                                "Попробуйте заново", client, update, ct);
            }

            Handlers.Tasks[fromId] = (null, null);
        }
    }

    public async Task HandelEnterShortUrl(ITelegramBotClient client, Update update, CancellationToken ct,
        List<string> list)
    {
        var fromId = update.Message.From.Id;
        var messageText = update.Message?.Text;
        list.Add(messageText);
        
        var request = new StringContent(JsonConvert.SerializeObject(new Request
        {
            url = list[0],
            vip_key = list[1],
            token = Handlers.Users[fromId],
            time_to_live_unit = "DAYS",
            time_to_live = 99999999
        }));
        var responseMessage =
            await _httpClient.PostAsync(_apiRequestPath, request, ct);
        var responseContent = await responseMessage.Content.ReadAsStringAsync(ct);
        var response = JsonConvert.DeserializeObject<Response>(responseContent);
        if (response?.error == null)
        {
            await TgMessageController.Write($"<b>Новая ссылка</b>: {response.short_url}\n" +
                                            $"<b>Ваш секреный ключ к данной ссылке</b>: <code>{response.secret_key}</code>", client, update, ct);
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

    public class Request
    {
        public int? time_to_live { get; set; }
        public string? time_to_live_unit { get; set; }
        public string? token { get; set; }
        public string? url { get; set; }
        public string? vip_key { get; set; }
    }

    public class Response
    {
        public string? secret_key { get; set; }
        public string? short_url { get; set; }
        public string? error { get; set; }
    }
}