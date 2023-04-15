FROM mcr.microsoft.com/dotnet/aspnet:6.0 AS base
WORKDIR /app
EXPOSE 80

FROM mcr.microsoft.com/dotnet/sdk:6.0 AS build
WORKDIR /src
COPY ["bot/shortener-tg/shortener-tg.csproj", "shortener-tg/"]
RUN dotnet restore "shortener-tg/shortener-tg.csproj"
COPY bot/ .
WORKDIR "/src/shortener-tg"
RUN dotnet build "shortener-tg.csproj" -c Release -o /app/build

FROM build AS publish
RUN dotnet publish "shortener-tg.csproj" -c Release -o /app/publish

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
COPY bot/shortener-tg/config.json .

ENTRYPOINT ["dotnet", "shortener-tg.dll"]
