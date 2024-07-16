using datamonster.Context;
using Microsoft.AspNetCore.Mvc;


namespace Settlement
{
    public class SettlementController()
    {
        public async Task<IResult> GetSettlementsForUser(
            [FromServices] ILogger<SettlementController> log,
            [FromServices] SettlementRepository repo,
            HttpContext context)
        {
            var id = context?.Items["id"] as string;
            if (string.IsNullOrWhiteSpace(id))
            {
                log.LogError("No user id found in claims.");
                return Results.Unauthorized();
            }
            var settlements = await repo.GetSettlementsForUser(id);

            return Results.Json(settlements, statusCode: 200); ;
        }

        public async Task<IResult> CreateSettlement(
            [FromServices] ILogger<SettlementController> log,
            [FromServices] SettlementRepository repo,
            [FromBody] SettlementRecord settlement,
            HttpContext context)
        {
            var id = context?.Items["id"] as string;
            if (id == null)
            {
                log.LogError("No user id found in claims.");
                return Results.Unauthorized();
            }
            settlement.Owner = id;
            var newRecord = await repo.CreateSettlement(settlement);
            return Results.Ok(newRecord);
        }

        public async Task<IResult> DeleteSettlement(
            [FromServices] ILogger<SettlementController> log,
            [FromServices] SettlementRepository repo,
            [FromRoute] int id,
            HttpContext context)
        {
            var userId = context?.Items["id"] as string;
            if (userId == null)
            {
                log.LogError("No user id found in claims.");
                return Results.Unauthorized();
            }
            var settlement = await repo.GetSettlement(id, userId);
            if (settlement == null)
            {
                log.LogError("Settlement not found.");
                return Results.NotFound();
            }
            await repo.DeleteSettlement(settlement);
            return Results.NoContent();
        }

        public async Task<IResult> GetSettlement(
            [FromServices] ILogger<SettlementController> log,
            [FromServices] SettlementRepository repo,
            [FromRoute] int id,
            HttpContext context)
        {
            var userId = context?.Items["id"] as string;
            if (userId == null)
            {
                log.LogError("No user id found in claims.");
                return Results.Unauthorized();
            }
            var settlement = await repo.GetSettlement(id, userId);
            if (settlement == null)
            {
                log.LogError("Settlement not found.");
                return Results.NotFound();
            }
            return Results.Ok(settlement);
        }

        public async Task<IResult> UpdateSettlement(
            [FromServices] ILogger<SettlementController> log,
            [FromServices] SettlementRepository repo,
            [FromRoute] int id,
            [FromBody] SettlementRecord settlement,
            HttpContext context)
        {
            var userId = context?.Items["id"] as string;
            if (userId == null)
            {
                log.LogError("No user id found in claims.");
                return Results.Unauthorized();
            }
            var existing = await repo.GetSettlement(id, userId);
            if (existing == null)
            {
                var newSettlement = await repo.CreateSettlement(settlement);
                return Results.Ok(newSettlement);
            }
            var updated = await repo.UpdateSettlement(settlement);
            return Results.Ok(updated);
        }
    }
}