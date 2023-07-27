# Task

You need to write a program (in golang) that accepts following params:
address of uniswap_v2 pool in ethereum (ex https://etherscan.io/address/0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852),
inputToken address, outputToken address, inputAmount.

Program should return outputAmount that corresponding uniswap_v2 pool  will return
if you try to swap inputAmount of  fromToken.
All math calculations should be done inside your program (not calling external services for results).

Example:

PoolID: 0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852
FromToken: 0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2
ToToken: 0xdac17f958d2ee523a2206206994597c13d831ec7
InputAmount: 1e18

You need to calculate outputAmount

# How to run

```
go build main.go

// from usdt to weth(with decimals)
./main 0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852 0xdac17f958d2ee523a2206206994597c13d831ec7 0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2 1868492640

// from weth to usdt(with decimals)
./main 0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852 0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2 0xdac17f958d2ee523a2206206994597c13d831ec7 1000000000000000000
```
