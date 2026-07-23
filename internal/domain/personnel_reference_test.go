package domain

import "testing"

func TestRankOptionsIncludeBrazilianMilitaryRanks(t *testing.T) {
	options := RankOptions()

	if len(options) != 17 {
		t.Fatalf("expected 17 rank options, got %d", len(options))
	}

	expected := []Rank{
		RankBrigadierGeneral,
		RankColonel,
		RankLieutenantColonel,
		RankMajor,
		RankCaptain,
		RankFirstLieutenant,
		RankSecondLieutenant,
		RankAspirant,
		RankWarrantOfficer,
		RankFirstSergeant,
		RankSecondSergeant,
		RankThirdSergeant,
		RankCorporal,
		RankSoldier,
		RankVariableEffectiveCorporal,
		RankVariableEffectiveSoldier,
		RankStudent,
	}

	for index, rank := range expected {
		if options[index].Value != rank {
			t.Fatalf("expected rank %s at index %d, got %s", rank, index, options[index].Value)
		}
	}
}

func TestRankLabelsAndAbbreviations(t *testing.T) {
	tests := []struct {
		name         string
		rank         Rank
		label        string
		abbreviation string
		displayLabel string
	}{
		{
			name:         "brigadier general",
			rank:         RankBrigadierGeneral,
			label:        "General de Brigada",
			abbreviation: "Gen Bda",
			displayLabel: "General de Brigada (Gen Bda)",
		},
		{
			name:         "third sergeant",
			rank:         RankThirdSergeant,
			label:        "Terceiro-Sargento",
			abbreviation: "3º Sgt",
			displayLabel: "Terceiro-Sargento (3º Sgt)",
		},
		{
			name:         "soldier",
			rank:         RankSoldier,
			label:        "Soldado",
			abbreviation: "Sd",
			displayLabel: "Soldado (Sd)",
		},
		{
			name:         "student",
			rank:         RankStudent,
			label:        "Aluno",
			abbreviation: "Al",
			displayLabel: "Aluno (Al)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rank.Label(); got != tt.label {
				t.Fatalf("expected label %q, got %q", tt.label, got)
			}

			if got := tt.rank.Abbreviation(); got != tt.abbreviation {
				t.Fatalf("expected abbreviation %q, got %q", tt.abbreviation, got)
			}

			if got := tt.rank.DisplayLabel(); got != tt.displayLabel {
				t.Fatalf("expected display label %q, got %q", tt.displayLabel, got)
			}
		})
	}
}

func TestLegacyRankAliasesRemainValid(t *testing.T) {
	if !IsValidRank(RankPrivate) {
		t.Fatal("expected RankPrivate alias to remain valid")
	}

	if !IsValidRank(RankSergeant) {
		t.Fatal("expected RankSergeant alias to remain valid")
	}

	if !IsValidRank(RankLieutenant) {
		t.Fatal("expected RankLieutenant alias to remain valid")
	}
}
