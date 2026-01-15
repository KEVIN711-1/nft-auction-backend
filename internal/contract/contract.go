package contract

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// /internal
// ├── contract/          # 类比为 include/ 目录
// │   └── contract.go  # 🔥 核心头文件！定义所有接口（类似 contract.h）
// ├── client/
// │   ├── nft_client.go      # 具体实现1（类似 nft_impl.c）
// │   └── auction_client.go  # 具体实现2（类似 auction_impl.c）
// └── service/
//     └── nft_service.

// ==================== NFT合约接口（不变）====================
type NFTContract interface {
	// 基本信息
	GetName(ctx context.Context) (string, error)
	GetSymbol(ctx context.Context) (string, error)

	// NFT 查询
	GetOwner(ctx context.Context, tokenID *big.Int) (common.Address, error)
	GetTokenURI(ctx context.Context, tokenID *big.Int) (string, error)
	GetBalanceOf(ctx context.Context, address common.Address) (*big.Int, error)
	CheckIfMinted(ctx context.Context, tokenID *big.Int) (bool, error)

	// 验证
	CheckOwner(ctx context.Context, tokenID *big.Int, address string) (bool, error)

	// 转账
	TransferFrom(ctx context.Context, from, to common.Address, tokenID *big.Int) error
}

// ==================== 拍卖合约接口（新增）====================
type AuctionContract interface {
	// 查询方法
	GetAuctionInfo(ctx context.Context, auctionID *big.Int) (
		common.Address, // seller
		*big.Int, // duration
		*big.Int, // startPrice
		*big.Int, // startTime
		bool, // ended
		common.Address, // highestBidder
		*big.Int, // highestBid
		common.Address, // nftContract
		*big.Int, // tokenId
		common.Address, // tokenAddress
		*big.Int, // bidTokenAmount
		*big.Int, // timeRemaining
		error,
	)

	GetAuctionCount(ctx context.Context) (*big.Int, error)
	GetAdmin(ctx context.Context) (common.Address, error)
	IsTokenAllowed(ctx context.Context, tokenAddress common.Address) (bool, error)

	// 交易方法（需要签名）
	PlaceBidETH(ctx context.Context, auctionID *big.Int, amount *big.Int) error
	PlaceBidERC20(ctx context.Context, auctionID *big.Int, amount *big.Int) error
	EndAuction(ctx context.Context, auctionID *big.Int) error
	CreateAuctionETH(ctx context.Context, duration *big.Int, startPrice *big.Int,
		nftAddress common.Address, tokenID *big.Int) error
	CreateAuctionERC20(ctx context.Context, duration *big.Int, startPrice *big.Int,
		nftAddress common.Address, tokenID *big.Int, erc20Token common.Address) error

	// 状态检查
	IsActive() bool
	GetContractAddress() common.Address
}
