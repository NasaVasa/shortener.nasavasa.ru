using Newtonsoft.Json;
using shortener_tg.Controllers;
using shortener_tg.Services.CreateUrl;

namespace shortener_tg.Services.InfoUrl;

using Telegram.Bot;
using Telegram.Bot.Types;
using Microsoft.Extensions.Configuration;
using Telegram.Bot.Types.Enums;
using Telegram.Bot.Types.ReplyMarkups;
using System.Text.Json;
using Setup;
using Utils.MessageChecker;

public class InfoUrlHandler : IUpdateHandler
{
    private static readonly IUpdateChecker UpdateChecker = new BaseUpdateChecker("InfoUrl");
    private static readonly IConfigurationSection Config = ConfigProvider.Config.GetSection("InfoUrl");

    private static readonly string _apiRequestPath = (Environment.GetEnvironmentVariable("API_PATH") ??
                                                      throw new InvalidOperationException(
                                                          "В переменных окружения не обнаружено API_PATH")) +
                                                     Config["Path"];

    private static readonly HttpClient _httpClient = new();

    public async Task HandelStart(ITelegramBotClient client, Update update, CancellationToken ct)
    {
        if (!Validate(update))
        {
            return;
        }

        var fromId = update.Message.From.Id;
        Handlers.Tasks[fromId] = (HandelEnterKey, null);
        await TgMessageController.Write("Введите секретный ключ от ссылки, информацию о которой вы хотите получить",
            client, update,
            ct);
    }

    public async Task HandelEnterKey(ITelegramBotClient client, Update update, CancellationToken ct, List<string> list)
    {
        var messageText = update.Message?.Text;
        var fromId = update.Message.From.Id;
        var response = MakeRequest(messageText, ct).Result;
        if (response.error == null)
        {
            await TgMessageController.Write("<b>Информация о ссылке</b>:\n" +
                                            $"<b>Полная ссылка</b>: {response.long_url}\n" +
                                            $"<b>Короткая ссылка</b>: {Environment.GetEnvironmentVariable("DOMAIN")}{response.short_url}\n" +
                                            $"<b>Количество уникальных переходов по ссылке</b>: {response.number_of_clicks}\n" +
                                            $"<b>Ссылка была создана</b>: {DateTime.Parse(response.dt_created)}\n",
                client, update,
                ct);
        }
        else
        {
            await TgMessageController.Write($"Возникла ошибка: <b>{response.error}</b>\n" +
                                            "Попробуйте заново", client, update,
                ct);
        }

        Handlers.Tasks[fromId] = (null, null);
    }

    public static async Task<Response> MakeRequest(string secretKey, CancellationToken ct)
    {
        if (secretKey[0] == '/')
        {
            return new Response()
            {
                error = "url not found"
            };
        }
        var request = _apiRequestPath + secretKey;
        var responseMessage =
            await _httpClient.GetAsync(request, ct);
        var responseContent = await responseMessage.Content.ReadAsStringAsync(ct);
        var response = JsonConvert.DeserializeObject<Response>(responseContent);
        return response;
    }


    public bool Validate(Update update)
    {
        return UpdateChecker.CheckMessage(update);
    }

    public class Response
    {
        public string[]? all_redirect_times { get; set; }
        public string? dt_created { get; set; }
        public string? dt_will_delete { get; set; }
        public string? long_url { get; set; }
        public int? number_of_clicks { get; set; }
        public string? short_url { get; set; }
        public string? error { get; set; }
    }
}