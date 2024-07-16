using datamonster.Auth;
using datamonster.Context;
using Microsoft.AspNetCore.HttpLogging;
using Settlement;
using datamonster.Middleware;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();
builder.Services.AddLogging();
builder.Services.AddHttpLogging(logging =>
{
    logging.LoggingFields = HttpLoggingFields.All;
    logging.RequestHeaders.Add("sec-ch-ua");
    logging.MediaTypeOptions.AddText("application/javascript");
    logging.RequestBodyLogLimit = 4096;
    logging.ResponseBodyLogLimit = 4096;
    logging.CombineLogs = true;
});
builder.Services.AddDbContext<RecordsContext>();
builder.Services.RegisterSettlementModule();
builder.Services.RegisterAuthModule();
builder.Services.AddTransient<UserIdExtractor>();

var app = builder.Build();
app.UseCors();
app.UseHttpLogging();
// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
    app.UseDeveloperExceptionPage();
}

app.UseRouting();
app.UseAuthentication();
app.UseAuthorization();
app.UseMiddleware<UserIdExtractor>();
app.MapSettlementEndpoints();
app.Run();