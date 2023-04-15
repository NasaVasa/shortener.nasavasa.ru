namespace shortener_tg.Setup;

public static class ConfigProvider
{
    public static IConfiguration Config { get; } = new ConfigurationBuilder()
        .SetBasePath(Environment.GetEnvironmentVariable("CONFIG_PATH") ?? "../../")
        .AddJsonFile("config.json", optional: false, reloadOnChange: true)
        .Build();
}