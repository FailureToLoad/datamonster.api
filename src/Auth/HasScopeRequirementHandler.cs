using Microsoft.AspNetCore.Authorization;

namespace datamonster.Auth;

public class HasScopeHandler : AuthorizationHandler<HasScopeRequirement>
{
    protected override Task HandleRequirementAsync(AuthorizationHandlerContext context, HasScopeRequirement requirement)
    {
        Console.WriteLine("hit scope requirement");
        if (!context.User.HasClaim(c => c.Type == "scope" && c.Issuer == requirement.Issuer))
            return Task.CompletedTask;

        // Split the scopes string into an array
        var scopes = context.User.FindFirst(c => c.Type! == "scope" && c.Issuer == requirement.Issuer)?.Value.Split(' ');

        // Succeed if the scope array contains the required scope
        if (scopes!.Any(s => s == requirement.Scope))
            context.Succeed(requirement);

        return Task.CompletedTask;
    }
}