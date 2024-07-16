using Microsoft.EntityFrameworkCore;

namespace datamonster.Context;

public partial class RecordsContext : DbContext
{
    public RecordsContext()
    {
    }

    public RecordsContext(DbContextOptions<RecordsContext> options)
        : base(options)
    {
    }

    public virtual DbSet<SettlementRecord> Settlements { get; set; }

    public virtual DbSet<SurvivorRecord> Survivors { get; set; }

    protected override void OnConfiguring(DbContextOptionsBuilder optionsBuilder)
        => optionsBuilder.UseNpgsql(Environment.GetEnvironmentVariable("PGSQL_CONN"));

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.Entity<SettlementRecord>(entity =>
        {
            entity.HasKey(e => e.Id).HasName("settlement_pkey");

            entity.Property(e => e.CollectiveCognition).HasDefaultValue((short)0);
            entity.Property(e => e.DepartingSurvival).HasDefaultValue((short)0);
            entity.Property(e => e.SurvivalLimit).HasDefaultValue((short)0);
            entity.Property(e => e.Year).HasDefaultValue((short)0);
        });

        modelBuilder.Entity<SurvivorRecord>(entity =>
        {
            entity.HasKey(e => e.Id).HasName("survivor_pkey");

            entity.Property(e => e.Accuracy).HasDefaultValue((short)0);
            entity.Property(e => e.Birth).HasDefaultValue((short)0);
            entity.Property(e => e.Courage).HasDefaultValue((short)0);
            entity.Property(e => e.Evasion).HasDefaultValue((short)0);
            entity.Property(e => e.Huntxp).HasDefaultValue((short)0);
            entity.Property(e => e.Insanity).HasDefaultValue((short)0);
            entity.Property(e => e.Luck).HasDefaultValue((short)0);
            entity.Property(e => e.Lumi).HasDefaultValue((short)0);
            entity.Property(e => e.Movement).HasDefaultValue((short)5);
            entity.Property(e => e.Speed).HasDefaultValue((short)0);
            entity.Property(e => e.Strength).HasDefaultValue((short)0);
            entity.Property(e => e.Survival).HasDefaultValue((short)1);
            entity.Property(e => e.SystemicPressure).HasDefaultValue((short)0);
            entity.Property(e => e.Torment).HasDefaultValue((short)0);
            entity.Property(e => e.Understanding).HasDefaultValue((short)0);

            entity.HasOne(d => d.SettlementNavigation).WithMany(p => p.Survivors).HasConstraintName("fk_settlement_id");
        });

        OnModelCreatingPartial(modelBuilder);
    }

    partial void OnModelCreatingPartial(ModelBuilder modelBuilder);
}