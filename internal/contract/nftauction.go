// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// NftAuctionMetaData contains all meta data concerning the NftAuction contract.
var NftAuctionMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"auctionId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startPrice\",\"type\":\"uint256\"}],\"name\":\"AuctionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"auctionId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"winner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"finalPrice\",\"type\":\"uint256\"}],\"name\":\"AuctionEnded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"auctionId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"bidder\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NewBid\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"allowERC20Token\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedERC20Tokens\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"auctions\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"ended\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"highestBidder\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"highestBid\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"nftContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"useERC20\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"erc20Token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_duration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_startPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_nftAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_erc20Token\",\"type\":\"address\"}],\"name\":\"createAuctionERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_duration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_startPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_nftAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"createAuctionETH\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_auctionId\",\"type\":\"uint256\"}],\"name\":\"endAuction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"ethAmount\",\"type\":\"uint256\"}],\"name\":\"ethToWei\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextAuctionId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onERC721Received\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_auctionId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"placeBidERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_auctionId\",\"type\":\"uint256\"}],\"name\":\"placeBidETH\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"weiAmount\",\"type\":\"uint256\"}],\"name\":\"weiToEth\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawETH\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f80fd5b503360025f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555061279f8061005c5f395ff3fe6080604052600436106100e7575f3560e01c8063a1db978211610089578063f851a44011610058578063f851a440146102f1578063f99a910c1461031b578063fa2c236514610357578063fc5284821461037f576100e7565b8063a1db97821461023d578063b9a2de3a14610265578063d83408e21461028d578063f14210a6146102c9576100e7565b8063571a26a0116100c5578063571a26a01461016b5780635c64987e146101b157806388439bbc146101d95780639d31a14214610215576100e7565b8063150b7a02146100eb5780632bfa5d6f146101275780634992c07314610143575b5f80fd5b3480156100f6575f80fd5b50610111600480360381019061010c9190611afe565b6103a9565b60405161011e9190611bbc565b60405180910390f35b610141600480360381019061013c9190611bd5565b6103bd565b005b34801561014e575f80fd5b5061016960048036038101906101649190611c00565b6105bc565b005b348015610176575f80fd5b50610191600480360381019061018c9190611bd5565b6106ea565b6040516101a89b9a99989796959493929190611caf565b60405180910390f35b3480156101bc575f80fd5b506101d760048036038101906101d29190611d58565b6107d5565b005b3480156101e4575f80fd5b506101ff60048036038101906101fa9190611bd5565b6108bc565b60405161020c9190611d83565b60405180910390f35b348015610220575f80fd5b5061023b60048036038101906102369190611d9c565b6108d8565b005b348015610248575f80fd5b50610263600480360381019061025e9190611dda565b610b79565b005b348015610270575f80fd5b5061028b60048036038101906102869190611bd5565b610ca9565b005b348015610298575f80fd5b506102b360048036038101906102ae9190611d58565b6111a2565b6040516102c09190611e18565b60405180910390f35b3480156102d4575f80fd5b506102ef60048036038101906102ea9190611bd5565b6111bf565b005b3480156102fc575f80fd5b5061030561131b565b6040516103129190611e31565b60405180910390f35b348015610326575f80fd5b50610341600480360381019061033c9190611bd5565b611340565b60405161034e9190611d83565b60405180910390f35b348015610362575f80fd5b5061037d60048036038101906103789190611e4a565b61135c565b005b34801561038a575f80fd5b506103936113ff565b6040516103a09190611d83565b60405180910390f35b5f63150b7a0260e01b905095945050505050565b806001548110610402576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103f990611f08565b60405180910390fd5b5f805f8381526020019081526020015f209050806004015f9054906101000a900460ff1615610466576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161045d90611f70565b60405180910390fd5b8060010154816003015461047a9190611fbb565b42106104bb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104b290612038565b60405180910390fd5b5f805f8581526020019081526020015f209050806008015f9054906101000a900460ff161561051f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610516906120c6565b60405180910390fd5b80600501543411610565576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161055c9061212e565b60405180910390fd5b80600201543410156105ac576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105a390612196565b60405180910390fd5b6105b68434611405565b50505050565b60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461064b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610642906121fe565b60405180910390fd5b60035f8273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff166106d4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106cb90612266565b60405180910390fd5b6106e3858585856001866114ba565b5050505050565b5f602052805f5260405f205f91509050805f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690806001015490806002015490806003015490806004015f9054906101000a900460ff16908060040160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690806005015490806006015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690806007015490806008015f9054906101000a900460ff16908060080160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508b565b60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610864576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161085b906121fe565b60405180910390fd5b600160035f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff02191690831515021790555050565b5f670de0b6b3a7640000826108d19190612284565b9050919050565b81600154811061091d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161091490611f08565b60405180910390fd5b5f805f8381526020019081526020015f209050806004015f9054906101000a900460ff1615610981576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161097890611f70565b60405180910390fd5b806001015481600301546109959190611fbb565b42106109d6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109cd90612038565b60405180910390fd5b5f805f8681526020019081526020015f209050806008015f9054906101000a900460ff16610a39576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a3090612335565b60405180910390fd5b80600501548411610a7f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a769061212e565b60405180910390fd5b8060020154841015610ac6576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610abd90612196565b60405180910390fd5b8060080160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166323b872dd3330876040518463ffffffff1660e01b8152600401610b2793929190612353565b6020604051808303815f875af1158015610b43573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610b6791906123b2565b50610b728585611405565b5050505050565b60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610c08576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bff906121fe565b60405180910390fd5b8173ffffffffffffffffffffffffffffffffffffffff1663a9059cbb60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16836040518363ffffffff1660e01b8152600401610c649291906123dd565b6020604051808303815f875af1158015610c80573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610ca491906123b2565b505050565b806001548110610cee576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ce590611f08565b60405180910390fd5b5f805f8481526020019081526020015f209050806004015f9054906101000a900460ff1615610d52576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d499061244e565b60405180910390fd5b80600101548160030154610d669190611fbb565b421015610da8576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d9f906124b6565b60405180910390fd5b6001816004015f6101000a81548160ff0219169083151502179055505f73ffffffffffffffffffffffffffffffffffffffff168160040160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16146110e957806006015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166342842e0e308360040160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1684600701546040518463ffffffff1660e01b8152600401610ea393929190612353565b5f604051808303815f87803b158015610eba575f80fd5b505af1158015610ecc573d5f803e3d5ffd5b50505050806008015f9054906101000a900460ff1615610fb1578060080160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb825f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1683600501546040518363ffffffff1660e01b8152600401610f6b9291906123dd565b6020604051808303815f875af1158015610f87573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610fab91906123b2565b50611081565b5f815f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168260050154604051610ffc90612501565b5f6040518083038185875af1925050503d805f8114611036576040519150601f19603f3d011682016040523d82523d5f602084013e61103b565b606091505b505090508061107f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016110769061255f565b60405180910390fd5b505b7fd2aa34a4fdbbc6dff6a3e56f46e0f3ae2a31d7785ff3487aa5c95c642acea501838260040160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1683600501546040516110dc9392919061257d565b60405180910390a161119d565b806006015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166342842e0e30835f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1684600701546040518463ffffffff1660e01b815260040161116f93929190612353565b5f604051808303815f87803b158015611186575f80fd5b505af1158015611198573d5f803e3d5ffd5b505050505b505050565b6003602052805f5260405f205f915054906101000a900460ff1681565b60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461124e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611245906121fe565b60405180910390fd5b5f60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168260405161129490612501565b5f6040518083038185875af1925050503d805f81146112ce576040519150601f19603f3d011682016040523d82523d5f602084013e6112d3565b606091505b5050905080611317576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161130e9061255f565b60405180910390fd5b5050565b60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b5f670de0b6b3a76400008261135591906125df565b9050919050565b60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146113eb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016113e2906121fe565b60405180910390fd5b6113f9848484845f806114ba565b50505050565b60015481565b5f805f8481526020019081526020015f2090505f8160050154111561142e5761142d81611850565b5b818160050181905550338160040160016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507f558a0d5d5468d74b0db24c74eb348b42271c2ebb4c9e953ced38aaed95fa43618333846040516114ad9392919061257d565b60405180910390a1505050565b603c8610156114fe576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016114f590612659565b60405180910390fd5b5f8511611540576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611537906126c1565b60405180910390fd5b8373ffffffffffffffffffffffffffffffffffffffff166342842e0e3330866040518463ffffffff1660e01b815260040161157d93929190612353565b5f604051808303815f87803b158015611594575f80fd5b505af11580156115a6573d5f803e3d5ffd5b505050506040518061016001604052803373ffffffffffffffffffffffffffffffffffffffff1681526020018781526020018681526020014281526020015f151581526020015f73ffffffffffffffffffffffffffffffffffffffff1681526020015f81526020018573ffffffffffffffffffffffffffffffffffffffff16815260200184815260200183151581526020018273ffffffffffffffffffffffffffffffffffffffff168152505f8060015481526020019081526020015f205f820151815f015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506020820151816001015560408201518160020155606082015181600301556080820151816004015f6101000a81548160ff02191690831515021790555060a08201518160040160016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060c0820151816005015560e0820151816006015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101008201518160070155610120820151816008015f6101000a81548160ff0219169083151502179055506101408201518160080160016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055509050507fc9050d42180a61cb0d9ebb8ad118b62fe6eab12cf12ff752c4a0cc7da9ddf62760015433858860405161182994939291906126df565b60405180910390a160015f81548092919061184390612722565b9190505550505050505050565b806008015f9054906101000a900460ff1615611933578060080160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb8260040160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1683600501546040518363ffffffff1660e01b81526004016118ed9291906123dd565b6020604051808303815f875af1158015611909573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061192d91906123b2565b50611a05565b5f8160040160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16826005015460405161198090612501565b5f6040518083038185875af1925050503d805f81146119ba576040519150601f19603f3d011682016040523d82523d5f602084013e6119bf565b606091505b5050905080611a03576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016119fa9061255f565b60405180910390fd5b505b50565b5f80fd5b5f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f611a3982611a10565b9050919050565b611a4981611a2f565b8114611a53575f80fd5b50565b5f81359050611a6481611a40565b92915050565b5f819050919050565b611a7c81611a6a565b8114611a86575f80fd5b50565b5f81359050611a9781611a73565b92915050565b5f80fd5b5f80fd5b5f80fd5b5f8083601f840112611abe57611abd611a9d565b5b8235905067ffffffffffffffff811115611adb57611ada611aa1565b5b602083019150836001820283011115611af757611af6611aa5565b5b9250929050565b5f805f805f60808688031215611b1757611b16611a08565b5b5f611b2488828901611a56565b9550506020611b3588828901611a56565b9450506040611b4688828901611a89565b935050606086013567ffffffffffffffff811115611b6757611b66611a0c565b5b611b7388828901611aa9565b92509250509295509295909350565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b611bb681611b82565b82525050565b5f602082019050611bcf5f830184611bad565b92915050565b5f60208284031215611bea57611be9611a08565b5b5f611bf784828501611a89565b91505092915050565b5f805f805f60a08688031215611c1957611c18611a08565b5b5f611c2688828901611a89565b9550506020611c3788828901611a89565b9450506040611c4888828901611a56565b9350506060611c5988828901611a89565b9250506080611c6a88828901611a56565b9150509295509295909350565b611c8081611a2f565b82525050565b611c8f81611a6a565b82525050565b5f8115159050919050565b611ca981611c95565b82525050565b5f61016082019050611cc35f83018e611c77565b611cd0602083018d611c86565b611cdd604083018c611c86565b611cea606083018b611c86565b611cf7608083018a611ca0565b611d0460a0830189611c77565b611d1160c0830188611c86565b611d1e60e0830187611c77565b611d2c610100830186611c86565b611d3a610120830185611ca0565b611d48610140830184611c77565b9c9b505050505050505050505050565b5f60208284031215611d6d57611d6c611a08565b5b5f611d7a84828501611a56565b91505092915050565b5f602082019050611d965f830184611c86565b92915050565b5f8060408385031215611db257611db1611a08565b5b5f611dbf85828601611a89565b9250506020611dd085828601611a89565b9150509250929050565b5f8060408385031215611df057611def611a08565b5b5f611dfd85828601611a56565b9250506020611e0e85828601611a89565b9150509250929050565b5f602082019050611e2b5f830184611ca0565b92915050565b5f602082019050611e445f830184611c77565b92915050565b5f805f8060808587031215611e6257611e61611a08565b5b5f611e6f87828801611a89565b9450506020611e8087828801611a89565b9350506040611e9187828801611a56565b9250506060611ea287828801611a89565b91505092959194509250565b5f82825260208201905092915050565b7f4e6f7420657869737400000000000000000000000000000000000000000000005f82015250565b5f611ef2600983611eae565b9150611efd82611ebe565b602082019050919050565b5f6020820190508181035f830152611f1f81611ee6565b9050919050565b7f456e6465640000000000000000000000000000000000000000000000000000005f82015250565b5f611f5a600583611eae565b9150611f6582611f26565b602082019050919050565b5f6020820190508181035f830152611f8781611f4e565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f611fc582611a6a565b9150611fd083611a6a565b9250828201905080821115611fe857611fe7611f8e565b5b92915050565b7f45787069726564000000000000000000000000000000000000000000000000005f82015250565b5f612022600783611eae565b915061202d82611fee565b602082019050919050565b5f6020820190508181035f83015261204f81612016565b9050919050565b7f546869732061756374696f6e20616363657074732045524332302c206e6f74205f8201527f4554480000000000000000000000000000000000000000000000000000000000602082015250565b5f6120b0602383611eae565b91506120bb82612056565b604082019050919050565b5f6020820190508181035f8301526120dd816120a4565b9050919050565b7f42696420746f6f206c6f770000000000000000000000000000000000000000005f82015250565b5f612118600b83611eae565b9150612123826120e4565b602082019050919050565b5f6020820190508181035f8301526121458161210c565b9050919050565b7f4269642062656c6f7720737461727420707269636500000000000000000000005f82015250565b5f612180601583611eae565b915061218b8261214c565b602082019050919050565b5f6020820190508181035f8301526121ad81612174565b9050919050565b7f4f6e6c792061646d696e000000000000000000000000000000000000000000005f82015250565b5f6121e8600a83611eae565b91506121f3826121b4565b602082019050919050565b5f6020820190508181035f830152612215816121dc565b9050919050565b7f546f6b656e206e6f7420616c6c6f7765640000000000000000000000000000005f82015250565b5f612250601183611eae565b915061225b8261221c565b602082019050919050565b5f6020820190508181035f83015261227d81612244565b9050919050565b5f61228e82611a6a565b915061229983611a6a565b92508282026122a781611a6a565b915082820484148315176122be576122bd611f8e565b5b5092915050565b7f546869732061756374696f6e2061636365707473204554482c206e6f742045525f8201527f4332300000000000000000000000000000000000000000000000000000000000602082015250565b5f61231f602383611eae565b915061232a826122c5565b604082019050919050565b5f6020820190508181035f83015261234c81612313565b9050919050565b5f6060820190506123665f830186611c77565b6123736020830185611c77565b6123806040830184611c86565b949350505050565b61239181611c95565b811461239b575f80fd5b50565b5f815190506123ac81612388565b92915050565b5f602082840312156123c7576123c6611a08565b5b5f6123d48482850161239e565b91505092915050565b5f6040820190506123f05f830185611c77565b6123fd6020830184611c86565b9392505050565b7f416c726561647920656e646564000000000000000000000000000000000000005f82015250565b5f612438600d83611eae565b915061244382612404565b602082019050919050565b5f6020820190508181035f8301526124658161242c565b9050919050565b7f4e6f742079657420656e646564000000000000000000000000000000000000005f82015250565b5f6124a0600d83611eae565b91506124ab8261246c565b602082019050919050565b5f6020820190508181035f8301526124cd81612494565b9050919050565b5f81905092915050565b50565b5f6124ec5f836124d4565b91506124f7826124de565b5f82019050919050565b5f61250b826124e1565b9150819050919050565b7f455448207472616e73666572206661696c6564000000000000000000000000005f82015250565b5f612549601383611eae565b915061255482612515565b602082019050919050565b5f6020820190508181035f8301526125768161253d565b9050919050565b5f6060820190506125905f830186611c86565b61259d6020830185611c77565b6125aa6040830184611c86565b949350505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f6125e982611a6a565b91506125f483611a6a565b925082612604576126036125b2565b5b828204905092915050565b7f4d696e203630207365636f6e64730000000000000000000000000000000000005f82015250565b5f612643600e83611eae565b915061264e8261260f565b602082019050919050565b5f6020820190508181035f83015261267081612637565b9050919050565b7f5374617274207072696365203e203000000000000000000000000000000000005f82015250565b5f6126ab600f83611eae565b91506126b682612677565b602082019050919050565b5f6020820190508181035f8301526126d88161269f565b9050919050565b5f6080820190506126f25f830187611c86565b6126ff6020830186611c77565b61270c6040830185611c86565b6127196060830184611c86565b95945050505050565b5f61272c82611a6a565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361275e5761275d611f8e565b5b60018201905091905056fea26469706673582212205af3a094642bfd6e5f715de2a4db31e68e48170319e707484d629ce9be35e0ff64736f6c63430008190033",
}

// NftAuctionABI is the input ABI used to generate the binding from.
// Deprecated: Use NftAuctionMetaData.ABI instead.
var NftAuctionABI = NftAuctionMetaData.ABI

// NftAuctionBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NftAuctionMetaData.Bin instead.
var NftAuctionBin = NftAuctionMetaData.Bin

// DeployNftAuction deploys a new Ethereum contract, binding an instance of NftAuction to it.
func DeployNftAuction(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *NftAuction, error) {
	parsed, err := NftAuctionMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NftAuctionBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NftAuction{NftAuctionCaller: NftAuctionCaller{contract: contract}, NftAuctionTransactor: NftAuctionTransactor{contract: contract}, NftAuctionFilterer: NftAuctionFilterer{contract: contract}}, nil
}

// NftAuction is an auto generated Go binding around an Ethereum contract.
type NftAuction struct {
	NftAuctionCaller     // Read-only binding to the contract
	NftAuctionTransactor // Write-only binding to the contract
	NftAuctionFilterer   // Log filterer for contract events
}

// NftAuctionCaller is an auto generated read-only Go binding around an Ethereum contract.
type NftAuctionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NftAuctionTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NftAuctionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NftAuctionFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NftAuctionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NftAuctionSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NftAuctionSession struct {
	Contract     *NftAuction       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NftAuctionCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NftAuctionCallerSession struct {
	Contract *NftAuctionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// NftAuctionTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NftAuctionTransactorSession struct {
	Contract     *NftAuctionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// NftAuctionRaw is an auto generated low-level Go binding around an Ethereum contract.
type NftAuctionRaw struct {
	Contract *NftAuction // Generic contract binding to access the raw methods on
}

// NftAuctionCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NftAuctionCallerRaw struct {
	Contract *NftAuctionCaller // Generic read-only contract binding to access the raw methods on
}

// NftAuctionTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NftAuctionTransactorRaw struct {
	Contract *NftAuctionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNftAuction creates a new instance of NftAuction, bound to a specific deployed contract.
func NewNftAuction(address common.Address, backend bind.ContractBackend) (*NftAuction, error) {
	contract, err := bindNftAuction(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NftAuction{NftAuctionCaller: NftAuctionCaller{contract: contract}, NftAuctionTransactor: NftAuctionTransactor{contract: contract}, NftAuctionFilterer: NftAuctionFilterer{contract: contract}}, nil
}

// NewNftAuctionCaller creates a new read-only instance of NftAuction, bound to a specific deployed contract.
func NewNftAuctionCaller(address common.Address, caller bind.ContractCaller) (*NftAuctionCaller, error) {
	contract, err := bindNftAuction(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NftAuctionCaller{contract: contract}, nil
}

// NewNftAuctionTransactor creates a new write-only instance of NftAuction, bound to a specific deployed contract.
func NewNftAuctionTransactor(address common.Address, transactor bind.ContractTransactor) (*NftAuctionTransactor, error) {
	contract, err := bindNftAuction(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NftAuctionTransactor{contract: contract}, nil
}

// NewNftAuctionFilterer creates a new log filterer instance of NftAuction, bound to a specific deployed contract.
func NewNftAuctionFilterer(address common.Address, filterer bind.ContractFilterer) (*NftAuctionFilterer, error) {
	contract, err := bindNftAuction(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NftAuctionFilterer{contract: contract}, nil
}

// bindNftAuction binds a generic wrapper to an already deployed contract.
func bindNftAuction(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NftAuctionMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NftAuction *NftAuctionRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NftAuction.Contract.NftAuctionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NftAuction *NftAuctionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NftAuction.Contract.NftAuctionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NftAuction *NftAuctionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NftAuction.Contract.NftAuctionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NftAuction *NftAuctionCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NftAuction.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NftAuction *NftAuctionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NftAuction.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NftAuction *NftAuctionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NftAuction.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_NftAuction *NftAuctionCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NftAuction.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_NftAuction *NftAuctionSession) Admin() (common.Address, error) {
	return _NftAuction.Contract.Admin(&_NftAuction.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_NftAuction *NftAuctionCallerSession) Admin() (common.Address, error) {
	return _NftAuction.Contract.Admin(&_NftAuction.CallOpts)
}

// AllowedERC20Tokens is a free data retrieval call binding the contract method 0xd83408e2.
//
// Solidity: function allowedERC20Tokens(address ) view returns(bool)
func (_NftAuction *NftAuctionCaller) AllowedERC20Tokens(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _NftAuction.contract.Call(opts, &out, "allowedERC20Tokens", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AllowedERC20Tokens is a free data retrieval call binding the contract method 0xd83408e2.
//
// Solidity: function allowedERC20Tokens(address ) view returns(bool)
func (_NftAuction *NftAuctionSession) AllowedERC20Tokens(arg0 common.Address) (bool, error) {
	return _NftAuction.Contract.AllowedERC20Tokens(&_NftAuction.CallOpts, arg0)
}

// AllowedERC20Tokens is a free data retrieval call binding the contract method 0xd83408e2.
//
// Solidity: function allowedERC20Tokens(address ) view returns(bool)
func (_NftAuction *NftAuctionCallerSession) AllowedERC20Tokens(arg0 common.Address) (bool, error) {
	return _NftAuction.Contract.AllowedERC20Tokens(&_NftAuction.CallOpts, arg0)
}

// Auctions is a free data retrieval call binding the contract method 0x571a26a0.
//
// Solidity: function auctions(uint256 ) view returns(address seller, uint256 duration, uint256 startPrice, uint256 startTime, bool ended, address highestBidder, uint256 highestBid, address nftContract, uint256 tokenId, bool useERC20, address erc20Token)
func (_NftAuction *NftAuctionCaller) Auctions(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Seller        common.Address
	Duration      *big.Int
	StartPrice    *big.Int
	StartTime     *big.Int
	Ended         bool
	HighestBidder common.Address
	HighestBid    *big.Int
	NftContract   common.Address
	TokenId       *big.Int
	UseERC20      bool
	Erc20Token    common.Address
}, error) {
	var out []interface{}
	err := _NftAuction.contract.Call(opts, &out, "auctions", arg0)

	outstruct := new(struct {
		Seller        common.Address
		Duration      *big.Int
		StartPrice    *big.Int
		StartTime     *big.Int
		Ended         bool
		HighestBidder common.Address
		HighestBid    *big.Int
		NftContract   common.Address
		TokenId       *big.Int
		UseERC20      bool
		Erc20Token    common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Seller = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Duration = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartPrice = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.StartTime = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Ended = *abi.ConvertType(out[4], new(bool)).(*bool)
	outstruct.HighestBidder = *abi.ConvertType(out[5], new(common.Address)).(*common.Address)
	outstruct.HighestBid = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.NftContract = *abi.ConvertType(out[7], new(common.Address)).(*common.Address)
	outstruct.TokenId = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)
	outstruct.UseERC20 = *abi.ConvertType(out[9], new(bool)).(*bool)
	outstruct.Erc20Token = *abi.ConvertType(out[10], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// Auctions is a free data retrieval call binding the contract method 0x571a26a0.
//
// Solidity: function auctions(uint256 ) view returns(address seller, uint256 duration, uint256 startPrice, uint256 startTime, bool ended, address highestBidder, uint256 highestBid, address nftContract, uint256 tokenId, bool useERC20, address erc20Token)
func (_NftAuction *NftAuctionSession) Auctions(arg0 *big.Int) (struct {
	Seller        common.Address
	Duration      *big.Int
	StartPrice    *big.Int
	StartTime     *big.Int
	Ended         bool
	HighestBidder common.Address
	HighestBid    *big.Int
	NftContract   common.Address
	TokenId       *big.Int
	UseERC20      bool
	Erc20Token    common.Address
}, error) {
	return _NftAuction.Contract.Auctions(&_NftAuction.CallOpts, arg0)
}

// Auctions is a free data retrieval call binding the contract method 0x571a26a0.
//
// Solidity: function auctions(uint256 ) view returns(address seller, uint256 duration, uint256 startPrice, uint256 startTime, bool ended, address highestBidder, uint256 highestBid, address nftContract, uint256 tokenId, bool useERC20, address erc20Token)
func (_NftAuction *NftAuctionCallerSession) Auctions(arg0 *big.Int) (struct {
	Seller        common.Address
	Duration      *big.Int
	StartPrice    *big.Int
	StartTime     *big.Int
	Ended         bool
	HighestBidder common.Address
	HighestBid    *big.Int
	NftContract   common.Address
	TokenId       *big.Int
	UseERC20      bool
	Erc20Token    common.Address
}, error) {
	return _NftAuction.Contract.Auctions(&_NftAuction.CallOpts, arg0)
}

// EthToWei is a free data retrieval call binding the contract method 0x88439bbc.
//
// Solidity: function ethToWei(uint256 ethAmount) pure returns(uint256)
func (_NftAuction *NftAuctionCaller) EthToWei(opts *bind.CallOpts, ethAmount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _NftAuction.contract.Call(opts, &out, "ethToWei", ethAmount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EthToWei is a free data retrieval call binding the contract method 0x88439bbc.
//
// Solidity: function ethToWei(uint256 ethAmount) pure returns(uint256)
func (_NftAuction *NftAuctionSession) EthToWei(ethAmount *big.Int) (*big.Int, error) {
	return _NftAuction.Contract.EthToWei(&_NftAuction.CallOpts, ethAmount)
}

// EthToWei is a free data retrieval call binding the contract method 0x88439bbc.
//
// Solidity: function ethToWei(uint256 ethAmount) pure returns(uint256)
func (_NftAuction *NftAuctionCallerSession) EthToWei(ethAmount *big.Int) (*big.Int, error) {
	return _NftAuction.Contract.EthToWei(&_NftAuction.CallOpts, ethAmount)
}

// NextAuctionId is a free data retrieval call binding the contract method 0xfc528482.
//
// Solidity: function nextAuctionId() view returns(uint256)
func (_NftAuction *NftAuctionCaller) NextAuctionId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NftAuction.contract.Call(opts, &out, "nextAuctionId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextAuctionId is a free data retrieval call binding the contract method 0xfc528482.
//
// Solidity: function nextAuctionId() view returns(uint256)
func (_NftAuction *NftAuctionSession) NextAuctionId() (*big.Int, error) {
	return _NftAuction.Contract.NextAuctionId(&_NftAuction.CallOpts)
}

// NextAuctionId is a free data retrieval call binding the contract method 0xfc528482.
//
// Solidity: function nextAuctionId() view returns(uint256)
func (_NftAuction *NftAuctionCallerSession) NextAuctionId() (*big.Int, error) {
	return _NftAuction.Contract.NextAuctionId(&_NftAuction.CallOpts)
}

// OnERC721Received is a free data retrieval call binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) pure returns(bytes4)
func (_NftAuction *NftAuctionCaller) OnERC721Received(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) ([4]byte, error) {
	var out []interface{}
	err := _NftAuction.contract.Call(opts, &out, "onERC721Received", arg0, arg1, arg2, arg3)

	if err != nil {
		return *new([4]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)

	return out0, err

}

// OnERC721Received is a free data retrieval call binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) pure returns(bytes4)
func (_NftAuction *NftAuctionSession) OnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) ([4]byte, error) {
	return _NftAuction.Contract.OnERC721Received(&_NftAuction.CallOpts, arg0, arg1, arg2, arg3)
}

// OnERC721Received is a free data retrieval call binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) pure returns(bytes4)
func (_NftAuction *NftAuctionCallerSession) OnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) ([4]byte, error) {
	return _NftAuction.Contract.OnERC721Received(&_NftAuction.CallOpts, arg0, arg1, arg2, arg3)
}

// WeiToEth is a free data retrieval call binding the contract method 0xf99a910c.
//
// Solidity: function weiToEth(uint256 weiAmount) pure returns(uint256)
func (_NftAuction *NftAuctionCaller) WeiToEth(opts *bind.CallOpts, weiAmount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _NftAuction.contract.Call(opts, &out, "weiToEth", weiAmount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WeiToEth is a free data retrieval call binding the contract method 0xf99a910c.
//
// Solidity: function weiToEth(uint256 weiAmount) pure returns(uint256)
func (_NftAuction *NftAuctionSession) WeiToEth(weiAmount *big.Int) (*big.Int, error) {
	return _NftAuction.Contract.WeiToEth(&_NftAuction.CallOpts, weiAmount)
}

// WeiToEth is a free data retrieval call binding the contract method 0xf99a910c.
//
// Solidity: function weiToEth(uint256 weiAmount) pure returns(uint256)
func (_NftAuction *NftAuctionCallerSession) WeiToEth(weiAmount *big.Int) (*big.Int, error) {
	return _NftAuction.Contract.WeiToEth(&_NftAuction.CallOpts, weiAmount)
}

// AllowERC20Token is a paid mutator transaction binding the contract method 0x5c64987e.
//
// Solidity: function allowERC20Token(address token) returns()
func (_NftAuction *NftAuctionTransactor) AllowERC20Token(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _NftAuction.contract.Transact(opts, "allowERC20Token", token)
}

// AllowERC20Token is a paid mutator transaction binding the contract method 0x5c64987e.
//
// Solidity: function allowERC20Token(address token) returns()
func (_NftAuction *NftAuctionSession) AllowERC20Token(token common.Address) (*types.Transaction, error) {
	return _NftAuction.Contract.AllowERC20Token(&_NftAuction.TransactOpts, token)
}

// AllowERC20Token is a paid mutator transaction binding the contract method 0x5c64987e.
//
// Solidity: function allowERC20Token(address token) returns()
func (_NftAuction *NftAuctionTransactorSession) AllowERC20Token(token common.Address) (*types.Transaction, error) {
	return _NftAuction.Contract.AllowERC20Token(&_NftAuction.TransactOpts, token)
}

// CreateAuctionERC20 is a paid mutator transaction binding the contract method 0x4992c073.
//
// Solidity: function createAuctionERC20(uint256 _duration, uint256 _startPrice, address _nftAddress, uint256 _tokenId, address _erc20Token) returns()
func (_NftAuction *NftAuctionTransactor) CreateAuctionERC20(opts *bind.TransactOpts, _duration *big.Int, _startPrice *big.Int, _nftAddress common.Address, _tokenId *big.Int, _erc20Token common.Address) (*types.Transaction, error) {
	return _NftAuction.contract.Transact(opts, "createAuctionERC20", _duration, _startPrice, _nftAddress, _tokenId, _erc20Token)
}

// CreateAuctionERC20 is a paid mutator transaction binding the contract method 0x4992c073.
//
// Solidity: function createAuctionERC20(uint256 _duration, uint256 _startPrice, address _nftAddress, uint256 _tokenId, address _erc20Token) returns()
func (_NftAuction *NftAuctionSession) CreateAuctionERC20(_duration *big.Int, _startPrice *big.Int, _nftAddress common.Address, _tokenId *big.Int, _erc20Token common.Address) (*types.Transaction, error) {
	return _NftAuction.Contract.CreateAuctionERC20(&_NftAuction.TransactOpts, _duration, _startPrice, _nftAddress, _tokenId, _erc20Token)
}

// CreateAuctionERC20 is a paid mutator transaction binding the contract method 0x4992c073.
//
// Solidity: function createAuctionERC20(uint256 _duration, uint256 _startPrice, address _nftAddress, uint256 _tokenId, address _erc20Token) returns()
func (_NftAuction *NftAuctionTransactorSession) CreateAuctionERC20(_duration *big.Int, _startPrice *big.Int, _nftAddress common.Address, _tokenId *big.Int, _erc20Token common.Address) (*types.Transaction, error) {
	return _NftAuction.Contract.CreateAuctionERC20(&_NftAuction.TransactOpts, _duration, _startPrice, _nftAddress, _tokenId, _erc20Token)
}

// CreateAuctionETH is a paid mutator transaction binding the contract method 0xfa2c2365.
//
// Solidity: function createAuctionETH(uint256 _duration, uint256 _startPrice, address _nftAddress, uint256 _tokenId) returns()
func (_NftAuction *NftAuctionTransactor) CreateAuctionETH(opts *bind.TransactOpts, _duration *big.Int, _startPrice *big.Int, _nftAddress common.Address, _tokenId *big.Int) (*types.Transaction, error) {
	return _NftAuction.contract.Transact(opts, "createAuctionETH", _duration, _startPrice, _nftAddress, _tokenId)
}

// CreateAuctionETH is a paid mutator transaction binding the contract method 0xfa2c2365.
//
// Solidity: function createAuctionETH(uint256 _duration, uint256 _startPrice, address _nftAddress, uint256 _tokenId) returns()
func (_NftAuction *NftAuctionSession) CreateAuctionETH(_duration *big.Int, _startPrice *big.Int, _nftAddress common.Address, _tokenId *big.Int) (*types.Transaction, error) {
	return _NftAuction.Contract.CreateAuctionETH(&_NftAuction.TransactOpts, _duration, _startPrice, _nftAddress, _tokenId)
}

// CreateAuctionETH is a paid mutator transaction binding the contract method 0xfa2c2365.
//
// Solidity: function createAuctionETH(uint256 _duration, uint256 _startPrice, address _nftAddress, uint256 _tokenId) returns()
func (_NftAuction *NftAuctionTransactorSession) CreateAuctionETH(_duration *big.Int, _startPrice *big.Int, _nftAddress common.Address, _tokenId *big.Int) (*types.Transaction, error) {
	return _NftAuction.Contract.CreateAuctionETH(&_NftAuction.TransactOpts, _duration, _startPrice, _nftAddress, _tokenId)
}

// EndAuction is a paid mutator transaction binding the contract method 0xb9a2de3a.
//
// Solidity: function endAuction(uint256 _auctionId) returns()
func (_NftAuction *NftAuctionTransactor) EndAuction(opts *bind.TransactOpts, _auctionId *big.Int) (*types.Transaction, error) {
	return _NftAuction.contract.Transact(opts, "endAuction", _auctionId)
}

// EndAuction is a paid mutator transaction binding the contract method 0xb9a2de3a.
//
// Solidity: function endAuction(uint256 _auctionId) returns()
func (_NftAuction *NftAuctionSession) EndAuction(_auctionId *big.Int) (*types.Transaction, error) {
	return _NftAuction.Contract.EndAuction(&_NftAuction.TransactOpts, _auctionId)
}

// EndAuction is a paid mutator transaction binding the contract method 0xb9a2de3a.
//
// Solidity: function endAuction(uint256 _auctionId) returns()
func (_NftAuction *NftAuctionTransactorSession) EndAuction(_auctionId *big.Int) (*types.Transaction, error) {
	return _NftAuction.Contract.EndAuction(&_NftAuction.TransactOpts, _auctionId)
}

// PlaceBidERC20 is a paid mutator transaction binding the contract method 0x9d31a142.
//
// Solidity: function placeBidERC20(uint256 _auctionId, uint256 _amount) returns()
func (_NftAuction *NftAuctionTransactor) PlaceBidERC20(opts *bind.TransactOpts, _auctionId *big.Int, _amount *big.Int) (*types.Transaction, error) {
	return _NftAuction.contract.Transact(opts, "placeBidERC20", _auctionId, _amount)
}

// PlaceBidERC20 is a paid mutator transaction binding the contract method 0x9d31a142.
//
// Solidity: function placeBidERC20(uint256 _auctionId, uint256 _amount) returns()
func (_NftAuction *NftAuctionSession) PlaceBidERC20(_auctionId *big.Int, _amount *big.Int) (*types.Transaction, error) {
	return _NftAuction.Contract.PlaceBidERC20(&_NftAuction.TransactOpts, _auctionId, _amount)
}

// PlaceBidERC20 is a paid mutator transaction binding the contract method 0x9d31a142.
//
// Solidity: function placeBidERC20(uint256 _auctionId, uint256 _amount) returns()
func (_NftAuction *NftAuctionTransactorSession) PlaceBidERC20(_auctionId *big.Int, _amount *big.Int) (*types.Transaction, error) {
	return _NftAuction.Contract.PlaceBidERC20(&_NftAuction.TransactOpts, _auctionId, _amount)
}

// PlaceBidETH is a paid mutator transaction binding the contract method 0x2bfa5d6f.
//
// Solidity: function placeBidETH(uint256 _auctionId) payable returns()
func (_NftAuction *NftAuctionTransactor) PlaceBidETH(opts *bind.TransactOpts, _auctionId *big.Int) (*types.Transaction, error) {
	return _NftAuction.contract.Transact(opts, "placeBidETH", _auctionId)
}

// PlaceBidETH is a paid mutator transaction binding the contract method 0x2bfa5d6f.
//
// Solidity: function placeBidETH(uint256 _auctionId) payable returns()
func (_NftAuction *NftAuctionSession) PlaceBidETH(_auctionId *big.Int) (*types.Transaction, error) {
	return _NftAuction.Contract.PlaceBidETH(&_NftAuction.TransactOpts, _auctionId)
}

// PlaceBidETH is a paid mutator transaction binding the contract method 0x2bfa5d6f.
//
// Solidity: function placeBidETH(uint256 _auctionId) payable returns()
func (_NftAuction *NftAuctionTransactorSession) PlaceBidETH(_auctionId *big.Int) (*types.Transaction, error) {
	return _NftAuction.Contract.PlaceBidETH(&_NftAuction.TransactOpts, _auctionId)
}

// WithdrawERC20 is a paid mutator transaction binding the contract method 0xa1db9782.
//
// Solidity: function withdrawERC20(address token, uint256 amount) returns()
func (_NftAuction *NftAuctionTransactor) WithdrawERC20(opts *bind.TransactOpts, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _NftAuction.contract.Transact(opts, "withdrawERC20", token, amount)
}

// WithdrawERC20 is a paid mutator transaction binding the contract method 0xa1db9782.
//
// Solidity: function withdrawERC20(address token, uint256 amount) returns()
func (_NftAuction *NftAuctionSession) WithdrawERC20(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _NftAuction.Contract.WithdrawERC20(&_NftAuction.TransactOpts, token, amount)
}

// WithdrawERC20 is a paid mutator transaction binding the contract method 0xa1db9782.
//
// Solidity: function withdrawERC20(address token, uint256 amount) returns()
func (_NftAuction *NftAuctionTransactorSession) WithdrawERC20(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _NftAuction.Contract.WithdrawERC20(&_NftAuction.TransactOpts, token, amount)
}

// WithdrawETH is a paid mutator transaction binding the contract method 0xf14210a6.
//
// Solidity: function withdrawETH(uint256 amount) returns()
func (_NftAuction *NftAuctionTransactor) WithdrawETH(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _NftAuction.contract.Transact(opts, "withdrawETH", amount)
}

// WithdrawETH is a paid mutator transaction binding the contract method 0xf14210a6.
//
// Solidity: function withdrawETH(uint256 amount) returns()
func (_NftAuction *NftAuctionSession) WithdrawETH(amount *big.Int) (*types.Transaction, error) {
	return _NftAuction.Contract.WithdrawETH(&_NftAuction.TransactOpts, amount)
}

// WithdrawETH is a paid mutator transaction binding the contract method 0xf14210a6.
//
// Solidity: function withdrawETH(uint256 amount) returns()
func (_NftAuction *NftAuctionTransactorSession) WithdrawETH(amount *big.Int) (*types.Transaction, error) {
	return _NftAuction.Contract.WithdrawETH(&_NftAuction.TransactOpts, amount)
}

// NftAuctionAuctionCreatedIterator is returned from FilterAuctionCreated and is used to iterate over the raw logs and unpacked data for AuctionCreated events raised by the NftAuction contract.
type NftAuctionAuctionCreatedIterator struct {
	Event *NftAuctionAuctionCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NftAuctionAuctionCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NftAuctionAuctionCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NftAuctionAuctionCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NftAuctionAuctionCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NftAuctionAuctionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NftAuctionAuctionCreated represents a AuctionCreated event raised by the NftAuction contract.
type NftAuctionAuctionCreated struct {
	AuctionId  *big.Int
	Seller     common.Address
	TokenId    *big.Int
	StartPrice *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAuctionCreated is a free log retrieval operation binding the contract event 0xc9050d42180a61cb0d9ebb8ad118b62fe6eab12cf12ff752c4a0cc7da9ddf627.
//
// Solidity: event AuctionCreated(uint256 auctionId, address seller, uint256 tokenId, uint256 startPrice)
func (_NftAuction *NftAuctionFilterer) FilterAuctionCreated(opts *bind.FilterOpts) (*NftAuctionAuctionCreatedIterator, error) {

	logs, sub, err := _NftAuction.contract.FilterLogs(opts, "AuctionCreated")
	if err != nil {
		return nil, err
	}
	return &NftAuctionAuctionCreatedIterator{contract: _NftAuction.contract, event: "AuctionCreated", logs: logs, sub: sub}, nil
}

// WatchAuctionCreated is a free log subscription operation binding the contract event 0xc9050d42180a61cb0d9ebb8ad118b62fe6eab12cf12ff752c4a0cc7da9ddf627.
//
// Solidity: event AuctionCreated(uint256 auctionId, address seller, uint256 tokenId, uint256 startPrice)
func (_NftAuction *NftAuctionFilterer) WatchAuctionCreated(opts *bind.WatchOpts, sink chan<- *NftAuctionAuctionCreated) (event.Subscription, error) {

	logs, sub, err := _NftAuction.contract.WatchLogs(opts, "AuctionCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NftAuctionAuctionCreated)
				if err := _NftAuction.contract.UnpackLog(event, "AuctionCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAuctionCreated is a log parse operation binding the contract event 0xc9050d42180a61cb0d9ebb8ad118b62fe6eab12cf12ff752c4a0cc7da9ddf627.
//
// Solidity: event AuctionCreated(uint256 auctionId, address seller, uint256 tokenId, uint256 startPrice)
func (_NftAuction *NftAuctionFilterer) ParseAuctionCreated(log types.Log) (*NftAuctionAuctionCreated, error) {
	event := new(NftAuctionAuctionCreated)
	if err := _NftAuction.contract.UnpackLog(event, "AuctionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NftAuctionAuctionEndedIterator is returned from FilterAuctionEnded and is used to iterate over the raw logs and unpacked data for AuctionEnded events raised by the NftAuction contract.
type NftAuctionAuctionEndedIterator struct {
	Event *NftAuctionAuctionEnded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NftAuctionAuctionEndedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NftAuctionAuctionEnded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NftAuctionAuctionEnded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NftAuctionAuctionEndedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NftAuctionAuctionEndedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NftAuctionAuctionEnded represents a AuctionEnded event raised by the NftAuction contract.
type NftAuctionAuctionEnded struct {
	AuctionId  *big.Int
	Winner     common.Address
	FinalPrice *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAuctionEnded is a free log retrieval operation binding the contract event 0xd2aa34a4fdbbc6dff6a3e56f46e0f3ae2a31d7785ff3487aa5c95c642acea501.
//
// Solidity: event AuctionEnded(uint256 auctionId, address winner, uint256 finalPrice)
func (_NftAuction *NftAuctionFilterer) FilterAuctionEnded(opts *bind.FilterOpts) (*NftAuctionAuctionEndedIterator, error) {

	logs, sub, err := _NftAuction.contract.FilterLogs(opts, "AuctionEnded")
	if err != nil {
		return nil, err
	}
	return &NftAuctionAuctionEndedIterator{contract: _NftAuction.contract, event: "AuctionEnded", logs: logs, sub: sub}, nil
}

// WatchAuctionEnded is a free log subscription operation binding the contract event 0xd2aa34a4fdbbc6dff6a3e56f46e0f3ae2a31d7785ff3487aa5c95c642acea501.
//
// Solidity: event AuctionEnded(uint256 auctionId, address winner, uint256 finalPrice)
func (_NftAuction *NftAuctionFilterer) WatchAuctionEnded(opts *bind.WatchOpts, sink chan<- *NftAuctionAuctionEnded) (event.Subscription, error) {

	logs, sub, err := _NftAuction.contract.WatchLogs(opts, "AuctionEnded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NftAuctionAuctionEnded)
				if err := _NftAuction.contract.UnpackLog(event, "AuctionEnded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAuctionEnded is a log parse operation binding the contract event 0xd2aa34a4fdbbc6dff6a3e56f46e0f3ae2a31d7785ff3487aa5c95c642acea501.
//
// Solidity: event AuctionEnded(uint256 auctionId, address winner, uint256 finalPrice)
func (_NftAuction *NftAuctionFilterer) ParseAuctionEnded(log types.Log) (*NftAuctionAuctionEnded, error) {
	event := new(NftAuctionAuctionEnded)
	if err := _NftAuction.contract.UnpackLog(event, "AuctionEnded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NftAuctionNewBidIterator is returned from FilterNewBid and is used to iterate over the raw logs and unpacked data for NewBid events raised by the NftAuction contract.
type NftAuctionNewBidIterator struct {
	Event *NftAuctionNewBid // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NftAuctionNewBidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NftAuctionNewBid)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NftAuctionNewBid)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NftAuctionNewBidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NftAuctionNewBidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NftAuctionNewBid represents a NewBid event raised by the NftAuction contract.
type NftAuctionNewBid struct {
	AuctionId *big.Int
	Bidder    common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNewBid is a free log retrieval operation binding the contract event 0x558a0d5d5468d74b0db24c74eb348b42271c2ebb4c9e953ced38aaed95fa4361.
//
// Solidity: event NewBid(uint256 auctionId, address bidder, uint256 amount)
func (_NftAuction *NftAuctionFilterer) FilterNewBid(opts *bind.FilterOpts) (*NftAuctionNewBidIterator, error) {

	logs, sub, err := _NftAuction.contract.FilterLogs(opts, "NewBid")
	if err != nil {
		return nil, err
	}
	return &NftAuctionNewBidIterator{contract: _NftAuction.contract, event: "NewBid", logs: logs, sub: sub}, nil
}

// WatchNewBid is a free log subscription operation binding the contract event 0x558a0d5d5468d74b0db24c74eb348b42271c2ebb4c9e953ced38aaed95fa4361.
//
// Solidity: event NewBid(uint256 auctionId, address bidder, uint256 amount)
func (_NftAuction *NftAuctionFilterer) WatchNewBid(opts *bind.WatchOpts, sink chan<- *NftAuctionNewBid) (event.Subscription, error) {

	logs, sub, err := _NftAuction.contract.WatchLogs(opts, "NewBid")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NftAuctionNewBid)
				if err := _NftAuction.contract.UnpackLog(event, "NewBid", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNewBid is a log parse operation binding the contract event 0x558a0d5d5468d74b0db24c74eb348b42271c2ebb4c9e953ced38aaed95fa4361.
//
// Solidity: event NewBid(uint256 auctionId, address bidder, uint256 amount)
func (_NftAuction *NftAuctionFilterer) ParseNewBid(log types.Log) (*NftAuctionNewBid, error) {
	event := new(NftAuctionNewBid)
	if err := _NftAuction.contract.UnpackLog(event, "NewBid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
