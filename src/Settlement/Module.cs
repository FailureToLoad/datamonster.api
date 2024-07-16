namespace Settlement;

public static class SettlementModule
{
    private static readonly SettlementController controller = new();
    public static IServiceCollection RegisterSettlementModule(this IServiceCollection services)
    {
        services.AddScoped<SettlementRepository>();
        return services;
    }

    public static IEndpointRouteBuilder MapSettlementEndpoints(this IEndpointRouteBuilder endpoints)
    {
        endpoints.MapGet("/settlement", controller.GetSettlementsForUser).RequireAuthorization("manage:settlements");
        endpoints.MapGet("/settlement/{id}", controller.GetSettlement).RequireAuthorization("manage:settlements");
        endpoints.MapPost("/settlement", controller.CreateSettlement).RequireAuthorization("manage:settlements");
        endpoints.MapPut("/settlement/{id}", controller.UpdateSettlement).RequireAuthorization("manage:settlements");
        return endpoints;
    }
}