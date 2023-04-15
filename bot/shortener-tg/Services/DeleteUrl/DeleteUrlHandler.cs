using Newtonsoft.Json;
using shortener_tg.Controllers;
using shortener_tg.Services.CreateUrl;

namespace shortener_tg.Services.DeleteUrl;

using Telegram.Bot;
using Telegram.Bot.Types;
using Microsoft.Extensions.Configuration;
using Telegram.Bot.Types.Enums;
using Telegram.Bot.Types.ReplyMarkups;
using System.Text.Json;
using Setup;
using Utils.MessageChecker;

public class DeleteUrlHandler : IUpdateHandler
{
    private static readonly IUpdateChecker UpdateChecker = new BaseUpdateChecker("DeleteUrl");
    private static readonly IConfigurationSection Config = ConfigProvider.Config.GetSection("DeleteUrl");

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
        Handlers.Tasks[fromId] = (HandelEnterKey, null);
        await TgMessageController.Write("Введите секретный ключ от ссылки, которую вы хотите удалить", client, update,
            ct);
    }

    public async Task HandelEnterKey(ITelegramBotClient client, Update update, CancellationToken ct, List<string> list)
    {
        var messageText = update.Message?.Text;
        var fromId = update.Message.From.Id;
        Response response;
        if (messageText[0] == '/')
        {
            response = new Response()
            {
                error = "url not found"
            };
        }
        else
        {
            var request = _apiRequestPath + messageText;
            var responseMessage =
                await _httpClient.DeleteAsync(request, ct);
            var responseContent = await responseMessage.Content.ReadAsStringAsync(ct);
            response = JsonConvert.DeserializeObject<Response>(responseContent);
        }


        if (response?.error == null)
        {
            await TgMessageController.Write("<b>Ссылка успешно удалена</b>", client, update,
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

    public bool Validate(Update update)
    {
        return UpdateChecker.CheckMessage(update);
    }

    public class Response
    {
        public string? error { get; set; }
    }
}