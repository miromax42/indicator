// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package strategy_test

import (
	"os"
	"testing"

	"github.com/miromax42/indicator/v2/asset"
	"github.com/miromax42/indicator/v2/helper"
	"github.com/miromax42/indicator/v2/strategy"
	"github.com/miromax42/indicator/v2/strategy/trend"
)

func TestSplitStrategy(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/repository/brk-b.csv", true)
	if err != nil {
		t.Fatal(err)
	}

	results, err := helper.ReadFromCsvFile[strategy.Result]("testdata/split.csv", true)
	if err != nil {
		t.Fatal(err)
	}

	expected := helper.Map(results, func(r *strategy.Result) strategy.Action { return r.Action })

	split := strategy.NewSplitStrategy(
		trend.NewMacdStrategy(),
		trend.NewApoStrategy(),
	)

	actual := split.Compute(snapshots)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSplitStrategyReport(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/repository/brk-b.csv", true)
	if err != nil {
		t.Fatal(err)
	}

	split := strategy.NewSplitStrategy(
		trend.NewMacdStrategy(),
		trend.NewApoStrategy(),
	)

	report := split.Report(snapshots)

	fileName := "split.html"
	defer os.Remove(fileName)

	err = report.WriteToFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAllSplitStrategies(t *testing.T) {
	strategies := []strategy.Strategy{
		strategy.NewBuyAndHoldStrategy(),
		strategy.NewMajorityStrategyWith("", []strategy.Strategy{
			strategy.NewBuyAndHoldStrategy(),
		}),
	}

	allSplitStrategies := strategy.AllSplitStrategies(strategies)

	expected := len(strategies)*len(strategies) - len(strategies)
	actual := len(allSplitStrategies)

	if actual != expected {
		t.Fatalf("actual=%d expected=%d", actual, expected)
	}
}
