using System.Text.RegularExpressions;
using Newtonsoft.Json;
using shortener_tg.Controllers;
using shortener_tg.Services.InfoUrl;

namespace shortener_tg.Services.GetUrlKeys;

using Telegram.Bot;
using Telegram.Bot.Types;
using Microsoft.Extensions.Configuration;
using Telegram.Bot.Types.Enums;
using Telegram.Bot.Types.ReplyMarkups;
using Setup;
using Utils.MessageChecker;

public class GetUrlKeysHandler : IUpdateHandler
{
    private static readonly IUpdateChecker UpdateChecker = new BaseUpdateChecker("GetUrlKeys");
    private static readonly IConfigurationSection Config = ConfigProvider.Config.GetSection("GetUrlKeys");

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
        if (Handlers.Users[fromId] is not null)
        {
            var request = _apiRequestPath + Handlers.Users[fromId];
            var responseMessage =
                await _httpClient.GetAsync(request, ct);
            var responseContent = await responseMessage.Content.ReadAsStringAsync(ct);
            var response = JsonConvert.DeserializeObject<Response>(responseContent);
            if (response?.error == null)

                if (response?.urls == null)
                {
                    await TgMessageController.Write($"<b>Список ваших ссылок пуст</b>", client, update,
                        ct);
                }
                else
                {
                    var res = "";
                    for (var i = 0; i < response.urls.Length; i++)
                    {
                        var currentResponse = InfoUrlHandler.MakeRequest(response.urls[i], ct).Result;
                        res += $"<b>---{i + 1}---</b>\n" +
                               $"<b>Полная ссылка</b>: {currentResponse.long_url}\n" +
                               $"<b>Короткая ссылка</b>: {Environment.GetEnvironmentVariable("DOMAIN")}{currentResponse.short_url}\n" +
                               $"<b>Секретный ключ</b>: <code>{response.urls[i]}</code>\n" +
                               $"<b>Количество уникальных переходов по ссылке</b>: {currentResponse.number_of_clicks}\n" +
                               $"<b>Ссылка была создана</b>: {DateTime.Parse(currentResponse.dt_created)}\n";
                    }

                    await TgMessageController.Write(res, client, update,
                        ct);
                }

            else
            {
                await TgMessageController.Write($"Возникла ошибка: <b>{response.error}</b>\n" +
                                                "Попробуйте заново", client, update, ct);
            }
        }
        else
        {
            await TgMessageController.Write("<b>Чтобы посмотреть свои ссылки необходимо авторизироваться</b>", client,
                update,
                ct);
        }
    }

    public bool Validate(Update update)
    {
        return UpdateChecker.CheckMessage(update);
    }

    class Response
    {
        public string[]? urls { get; set; }
        public string? error { get; set; }
    }
}