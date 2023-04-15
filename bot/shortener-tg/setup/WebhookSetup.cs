using Telegram.Bot;

namespace shortener_tg.Setup;

public class WebhookSetup
{
    public static void Setup(string[] args)
    {
        ITelegramBotClient client = new TelegramBotClient(Environment.GetEnvironmentVariable("BOT_TOKEN") ?? "");

        client.SetWebhookAsync(Environment.GetEnvironmentVariable("webhook_url")!);
        Console.WriteLine($"Webhook set to {Environment.GetEnvironmentVariable("webhook_url")}");

        var builder = WebApplication.CreateBuilder(args);
        builder.Services
            .AddControllers()
            .AddNewtonsoftJson();

        var app = builder.Build();
        app.MapControllers();

        app.Run();
    }
}