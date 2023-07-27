package main

import (
	"log"
	"math/big"
	"testing"
)

func TestCalculateOutputAmount(t *testing.T) {
	poolId := "0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852"
	fromToken := "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
	toToken := "0xdac17f958d2ee523a2206206994597c13d831ec7"
	inputAmount := big.NewInt(1000000000000000000)
	expectedOutputAmount := big.NewInt(1870000000)

	outputAmount, err := CalculateOutputAmount(poolId, fromToken, toToken, inputAmount)
	if err != nil {
		t.Errorf("CalculateOutputAmount failed with error: %s", err)
	}
	log.Printf("%v", outputAmount)

	if !withinPercentage(outputAmount, expectedOutputAmount, 10) {
		t.Errorf("CalculateOutputAmount incorrect, got: %s, want close to: %s", outputAmount.String(), expectedOutputAmount.String())
	}
}

func withinPercentage(a, b *big.Int, margin int64) bool {
	diff := new(big.Int).Sub(a, b)
	diff.Abs(diff)

	percentageDiff := new(big.Int).Mul(diff, big.NewInt(100))
	percentageDiff.Div(percentageDiff, a)

	return percentageDiff.Cmp(big.NewInt(margin)) <= 0
}
