package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	args := os.Args[1:]
	if len(args) != 4 {
		log.Fatalf("Number of arguments should equal 4")
	}
	poolId := args[0]
	fromToken := args[1]
	toToken := args[2]
	inputAmount := new(big.Int)
	inputAmount.SetString(args[3], 10)

	amount, err := CalculateOutputAmount(poolId, fromToken, toToken, inputAmount)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("Output amount is %v", amount)
}

func CalculateOutputAmount(poolId string, from string, to string, input *big.Int) (*big.Int, error) {
	pairAbi, err := parseAbi("PairAbi.json")
	if err != nil {
		return nil, fmt.Errorf("can't fetch ABI: %v", err)
	}

	// keep only for checking the solution
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/045e48c084c44c3cafa56096ce9acd42")
	if err != nil {
		return nil, fmt.Errorf("can't connect to Ethereum client: %v", err)
	}

	contractAddress := common.HexToAddress(poolId)

	r0, r1, err := getReserves(client, pairAbi, contractAddress)
	if err != nil {
		return nil, fmt.Errorf("can't get reserves: %v", err)
	}
	t0, err := getToken0(client, pairAbi, contractAddress)
	if err != nil {
		return nil, fmt.Errorf("can't get token0: %v", err)
	}
	t1, err := getToken1(client, pairAbi, contractAddress)
	if err != nil {
		return nil, fmt.Errorf("can't get token1: %v", err)
	}

	if !strings.EqualFold(from, t0) && !strings.EqualFold(from, t1) {
		return nil, fmt.Errorf("token %v is not in the pool", from)
	}
	if !strings.EqualFold(to, t0) && !strings.EqualFold(to, t1) {
		return nil, fmt.Errorf("token %v is not in the pool", to)
	}

	if strings.EqualFold(from, t0) {
		return calculate(input, r0, r1), nil
	}
	return calculate(input, r1, r0), nil
}

func parseAbi(path string) (*abi.ABI, error) {
	abiFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open ABI file: %v", err)
	}
	defer abiFile.Close()

	abiBytes, err := io.ReadAll(abiFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read ABI file: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI file: %v", err)
	}
	return &parsedABI, nil
}

func getReserves(client *ethclient.Client, abi *abi.ABI, poolAddress common.Address) (*big.Int, *big.Int, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()
	output, err := client.CallContract(ctx, ethereum.CallMsg{
		To:   &poolAddress,
		Data: abi.Methods["getReserves"].ID,
	}, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to call contract: %v", err)
	}

	reserves, err := abi.Methods["getReserves"].Outputs.Unpack(output)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unpack contract: %v", err)
	}

	reserve0 := reserves[0].(*big.Int)
	reserve1 := reserves[1].(*big.Int)
	return reserve0, reserve1, nil
}

func getToken0(client *ethclient.Client, abi *abi.ABI, poolAddress common.Address) (string, error) {
	return getToken(client, abi, poolAddress, "token0")
}

func getToken1(client *ethclient.Client, abi *abi.ABI, poolAddress common.Address) (string, error) {
	return getToken(client, abi, poolAddress, "token1")
}

func getToken(client *ethclient.Client, abi *abi.ABI, poolAddress common.Address, method string) (string, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()

	output, err := client.CallContract(ctx, ethereum.CallMsg{
		To:   &poolAddress,
		Data: abi.Methods[method].ID,
	}, nil)
	if err != nil {
		return "", fmt.Errorf("can't call contract(%v): %v", method, err)
	}

	response, err := abi.Methods[method].Outputs.Unpack(output)
	if err != nil {
		return "", fmt.Errorf("can't unpack output: %v", err)
	}
	token0 := response[0].(common.Address)
	return token0.Hex(), nil
}

func calculate(amountIn *big.Int, reserveIn *big.Int, reserveOut *big.Int) *big.Int {
	amountInWithFee := new(big.Int).Mul(amountIn, big.NewInt(997))
	numerator := new(big.Int).Mul(amountInWithFee, reserveOut)
	denominator := new(big.Int).Mul(reserveIn, big.NewInt(1000))
	denominator.Add(denominator, amountInWithFee)
	return numerator.Div(numerator, denominator)
}
