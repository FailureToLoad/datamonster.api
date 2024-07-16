using System.Security.Claims;
using Castle.Components.DictionaryAdapter.Xml;
using datamonster.Context;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using Moq;
using Newtonsoft.Json;
using Settlement;

namespace datamonster.test;

public class SettlementEndpointTests
{
    private SettlementController target;

    [SetUp]
    public void Setup()
    {
        target = new SettlementController();
    }

    [Test]
    public async Task GetSettlementsForUser_ReturnsListWhenGivenId()
    {
        var owner = "owner";
        var mockDb = GetMockDb(owner);
        var logger = new Mock<ILogger<SettlementController>>();
        var context = new DefaultHttpContext
        {
            // RequestServices needs to be set so the IResult implementation can log.
            RequestServices = new ServiceCollection().AddLogging().BuildServiceProvider(),
            Response =
            {
                // The default response body is Stream.Null which throws away anything that is written to it.
                Body = new MemoryStream(),
            },
        };

        context.Request.HttpContext.User = new ClaimsPrincipal(new ClaimsIdentity(new Claim[]
        {
            new("id", owner)
        }));
        var result = await target.GetSettlementsForUser(logger.Object, new SettlementRepository(mockDb.Object), context);
        Assert.That(context.Response.StatusCode, Is.EqualTo(200));


        var returnedValue = await GetResult(result, context);
        Assert.That(returnedValue.Count, Is.EqualTo(3));
        Assert.Multiple(() =>
        {
            Assert.That(returnedValue.Settlements[0].Id, Is.EqualTo(1));
            Assert.That(returnedValue.Settlements[0].Owner, Is.EqualTo("owner"));
        });
    }

    private async Task<ListSettlementsResponse> GetResult(IResult result, HttpContext context)
    {
        await result.ExecuteAsync(context);
        var stream = context.Response.Body;
        stream.Position = 0;
        using var reader = new StreamReader(stream);
        var body = reader.ReadToEnd();
        var settlementResponse = JsonConvert.DeserializeObject<ListSettlementsResponse>(body);
        return settlementResponse;
    }

    private static Mock<RecordsContext> GetMockDb(string owner)
    {
        var data = new List<SettlementRecord>
            {
                new() { Id = 1,  Owner = owner },
                new() { Id = 2,  Owner = "other"},
                new() { Id = 3,  Owner = owner},
                new() { Id = 4,  Owner = "other"},
                new() { Id = 5,  Owner = owner},
                new() { Id = 6, Owner = "other"},
            }.AsAsyncQueryable();

        var mockSet = new Mock<DbSet<SettlementRecord>>();
        mockSet.As<IQueryable<SettlementRecord>>().Setup(m => m.Provider).Returns(data.Provider);
        mockSet.As<IQueryable<SettlementRecord>>().Setup(m => m.Expression).Returns(data.Expression);
        mockSet.As<IQueryable<SettlementRecord>>().Setup(m => m.ElementType).Returns(data.ElementType);
        mockSet.As<IQueryable<SettlementRecord>>().Setup(m => m.GetEnumerator()).Returns(() => data.GetEnumerator());
        var dbContext = new Mock<RecordsContext>();
        dbContext.Setup(x => x.Settlements).Returns(mockSet.Object);
        return dbContext;
    }
}