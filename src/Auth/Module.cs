using System.Security.Claims;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.IdentityModel.Tokens;
using Microsoft.AspNetCore.Authorization;

namespace datamonster.Auth;
public static class AuthModule
{

    public static IServiceCollection RegisterAuthModule(this IServiceCollection services)
    {
        services.AddCors(options =>
        {
            options.AddDefaultPolicy(builder =>
            {
                builder.WithOrigins("http://localhost:8090");
                builder.WithMethods("HEAD", "GET", "POST");
                builder.WithHeaders("Origin", "Accept", "Authorization", "Content-Type", "X-CSRF-Token");
                builder.AllowCredentials();
                builder.SetPreflightMaxAge(TimeSpan.FromSeconds(2520));
            });
        });

        services.AddAuthentication(options =>
        {
            options.DefaultAuthenticateScheme = JwtBearerDefaults.AuthenticationScheme;
            options.DefaultChallengeScheme = JwtBearerDefaults.AuthenticationScheme;
        })
        .AddJwtBearer(options =>
        {
            options.Authority = Environment.GetEnvironmentVariable("AUTH0_DOMAIN");
            options.Audience = Environment.GetEnvironmentVariable("AUTH0_AUDIENCE");
            options.TokenValidationParameters = new TokenValidationParameters
            {
                NameClaimType = ClaimTypes.NameIdentifier
            };
            options.IncludeErrorDetails = true;
        });

        services.AddSingleton<IAuthorizationHandler, HasScopeHandler>();

        services.AddAuthorizationBuilder()
            .AddPolicy("manage:settlements", policy => policy.Requirements.Add(new HasScopeRequirement("manage:settlements", Environment.GetEnvironmentVariable("AUTH0_DOMAIN")!)));

        return services;
    }
}