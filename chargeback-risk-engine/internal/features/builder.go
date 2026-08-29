package features

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	riskv1 "github.com/maksimp027/chargeback_predict/api/proto/v1"
)

type Schema struct {
	Features   []string           `json:"features"`
	BinRiskMap map[string]float64 `json:"bin_risk_map"`
}

type FeatureBuilder struct {
	schema Schema
}

func NewFeatureBuilder(schemaPath string) (*FeatureBuilder, error) {
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	var schema Schema
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("failed to parse schema json: %w", err)
	}

	return &FeatureBuilder{schema: schema}, nil
}

func (fb *FeatureBuilder) Build(req *riskv1.RiskEvaluationRequest, vf VelocityFeatures, reqTime time.Time) []float32 {
	feats := make([]float32, 19)

	// 0: amount_log = log1p(amount_cents / 100.0)
	feats[0] = float32(math.Log1p(float64(req.AmountCents) / 100.0))

	// 1: is_round_amount = (amount_cents % 1000 == 0)
	if req.AmountCents%1000 == 0 {
		feats[1] = 1.0
	} else {
		feats[1] = 0.0
	}

	// 2: amount_cents_mod_100
	feats[2] = float32(req.AmountCents % 100)

	// 3: hour_sin
	hour := float64(reqTime.Hour())
	feats[3] = float32(math.Sin(2 * math.Pi * hour / 24.0))

	// 4: hour_cos
	feats[4] = float32(math.Cos(2 * math.Pi * hour / 24.0))

	// 5: day_of_week (0=Mon .. 6=Sun)
	dow := (int(reqTime.Weekday()) + 6) % 7
	feats[5] = float32(dow)

	// 6: is_weekend (5=Sat, 6=Sun)
	if dow >= 5 {
		feats[6] = 1.0
	} else {
		feats[6] = 0.0
	}

	// 7: cvv_match_num
	if req.CvvMatch {
		feats[7] = 1.0
	} else {
		feats[7] = 0.0
	}

	// 8: is_3ds_passed_num
	if req.Is_3DsPassed {
		feats[8] = 1.0
	} else {
		feats[8] = 0.0
	}

	// 9: risk_no_3ds_cvv_fail = (!cvv_match && !is_3ds_passed)
	if !req.CvvMatch && !req.Is_3DsPassed {
		feats[9] = 1.0
	} else {
		feats[9] = 0.0
	}

	// 10: currency_code (1.0 for USD)
	if req.Currency == "USD" {
		feats[10] = 1.0
	} else {
		feats[10] = 0.0
	}

	// 11: billing_country_code (1.0 for US)
	if req.BillingCountry == "US" {
		feats[11] = 1.0
	} else {
		feats[11] = 0.0
	}

	// 12: account_age_days = (reqTime.Unix() - req.UserCreatedAt) / 86400.0
	var ageDays float64
	if req.UserCreatedAt > 0 {
		ageSec := float64(reqTime.Unix() - req.UserCreatedAt)
		if ageSec < 0 {
			ageSec = 0
		}
		ageDays = ageSec / 86400.0
	}
	feats[12] = float32(ageDays)

	// 13: account_age_days_log = log1p(max(0, account_age_days))
	feats[13] = float32(math.Log1p(math.Max(0, ageDays)))

	// 14: card_bin_risk
	binStr := fmt.Sprintf("%d", req.CardBin)
	binRisk := 0.07
	if r, ok := fb.schema.BinRiskMap[binStr]; ok {
		binRisk = r
	}
	feats[14] = float32(binRisk)

	// 15: card_tx_cnt_5m
	feats[15] = vf.CardTxCnt5m

	// 16: card_tx_cnt_15m
	feats[16] = vf.CardTxCnt15m

	// 17: card_tx_cnt_60m
	feats[17] = vf.CardTxCnt60m

	// 18: card_sum_cents_60m
	feats[18] = vf.CardSumCents60m

	return feats
}
