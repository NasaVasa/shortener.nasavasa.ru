using Telegram.Bot;
using shortener_tg.Services;

namespace shortener_tg.Setup;

public static class PollingSetup
{
    public static void Setup()
    {
        ITelegramBotClient client = new TelegramBotClient(Environment.GetEnvironmentVariable("BOT_TOKEN") ?? "");
        client.DeleteWebhookAsync();
        client.StartReceiving(Handlers.HandelUpdates, Handlers.ErrorHandler);
       
        Console.ReadLine();
        client.SetWebhookAsync(Environment.GetEnvironmentVariable("webhook_url") ?? "");
    }
}