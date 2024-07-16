using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;
using Microsoft.EntityFrameworkCore;
namespace datamonster.Context;

[Table("survivor", Schema = "campaign")]
[Index("Settlement", "Name", Name = "survivor_settlement_name_key", IsUnique = true)]
public partial class SurvivorRecord
{
    [Key]
    [Column("id")]
    public int Id { get; set; }

    [Column("settlement")]
    public int Settlement { get; set; }

    [Column("name")]
    [StringLength(50)]
    public string Name { get; set; } = null!;

    [Column("gender")]
    [MaxLength(1)]
    public char? Gender { get; set; }

    [Column("birth")]
    public short Birth { get; set; }

    [Column("huntxp")]
    public short Huntxp { get; set; }

    [Column("survival")]
    public short Survival { get; set; }

    [Column("courage")]
    public short Courage { get; set; }

    [Column("understanding")]
    public short Understanding { get; set; }

    [Column("movement")]
    public short Movement { get; set; }

    [Column("accuracy")]
    public short Accuracy { get; set; }

    [Column("strength")]
    public short Strength { get; set; }

    [Column("evasion")]
    public short Evasion { get; set; }

    [Column("luck")]
    public short Luck { get; set; }

    [Column("speed")]
    public short Speed { get; set; }

    [Column("insanity")]
    public short Insanity { get; set; }

    [Column("systemic_pressure")]
    public short SystemicPressure { get; set; }

    [Column("torment")]
    public short Torment { get; set; }

    [Column("lumi")]
    public short Lumi { get; set; }

    [Column("status")]
    [StringLength(50)]
    public string? Status { get; set; }

    [ForeignKey("Settlement")]
    [InverseProperty("Survivors")]
    public virtual SettlementRecord SettlementNavigation { get; set; } = null!;
}