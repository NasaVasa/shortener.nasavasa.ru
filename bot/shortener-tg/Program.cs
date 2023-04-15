using shortener_tg.Services;
using Telegram.Bot;

namespace shortener_tg;

class Program
{
    static void Main(string[] args)
    {
//        var builder = WebApplication.CreateBuilder(args);
//        var app = builder.Build();

        ITelegramBotClient client =
            new TelegramBotClient(Environment.GetEnvironmentVariable("BOT_TOKEN") ??
                                  throw new ArgumentException("В переменных окружения не было обнаружено BOT_TOKEN"));
        client.StartReceiving(Handlers.HandelUpdates, Handlers.ErrorHandler);

        while(true)
        {
            Console.ReadLine();
        }

//        app.Run();
    }
}