using datamonster.Context;
using Microsoft.EntityFrameworkCore;

namespace Settlement;

public class SettlementRepository(RecordsContext context)
{
    private readonly RecordsContext context = context;

    public async Task<List<SettlementRecord>> GetSettlementsForUser(string user)
    {
        return await context.Settlements.Where(s => s.Owner == user).ToListAsync();
    }

    public async Task<SettlementRecord?> GetSettlement(int id, string user)
    {
        return await context.Settlements.Where(s => s.Id == id && s.Owner == user).FirstOrDefaultAsync();
    }

    public async Task<SettlementRecord?> CreateSettlement(SettlementRecord settlement)
    {
        context.Settlements.Add(settlement);
        await context.SaveChangesAsync();
        return settlement;
    }

    public async Task<SettlementRecord?> UpdateSettlement(SettlementRecord settlement)
    {
        context.Settlements.Update(settlement);
        await context.SaveChangesAsync();
        return settlement;
    }

    public async Task<bool> DeleteSettlement(SettlementRecord settlement)
    {
        context.Settlements.Remove(settlement);
        await context.SaveChangesAsync();
        return true;
    }
}