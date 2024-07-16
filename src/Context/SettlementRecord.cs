using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;
namespace datamonster.Context;

[Table("settlement", Schema = "campaign")]
public partial class SettlementRecord
{
    [Key]
    [Column("id")]
    public int Id { get; set; }

    [Column("owner")]
    [StringLength(50)]
    public string Owner { get; set; } = null!;

    [Column("name")]
    [StringLength(50)]
    public string Name { get; set; } = null!;

    [Column("survival_limit")]
    public short? SurvivalLimit { get; set; }

    [Column("departing_survival")]
    public short? DepartingSurvival { get; set; }

    [Column("collective_cognition")]
    public short? CollectiveCognition { get; set; }

    [Column("year")]
    public short? Year { get; set; }

    [InverseProperty("SettlementNavigation")]
    public virtual ICollection<SurvivorRecord> Survivors { get; set; } = new List<SurvivorRecord>();
}