namespace datamonster.Middleware;

public class UserIdExtractor : IMiddleware
{
    public async Task InvokeAsync(HttpContext context, RequestDelegate next)
    {
        var name = context.User?.Identity?.Name;
        if (!string.IsNullOrEmpty(name))
        {
            var items = name.Split('|');
            context.Items["id"] = items[1];
        }

        await next(context);
    }
}