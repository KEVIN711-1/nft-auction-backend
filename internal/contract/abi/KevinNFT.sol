// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";


// 编写 NFT 合约
// 使用 OpenZeppelin 的 ERC721 库编写一个 NFT 合约。
// 合约应包含以下功能：
// 构造函数：设置 NFT 的名称和符号。
// mintNFT 函数：允许用户铸造 NFT，并关联元数据链接（tokenURI）。
// 在 Remix IDE 中编译合约。
contract KevinNFT is ERC721, ERC721URIStorage, Ownable {
    uint256 private newTokenId;
    uint256 private minPrice=0.01 ether;

    event nftTransfer( uint256 value, address to);
    event nftMint( address owner, uint256 tokenId, string uri);

    constructor ( string memory _name, string memory  _symbol )
        ERC721( _name, _symbol)
        Ownable(msg.sender) {
    }

    function mintNFT ( string calldata uri ) public payable {
        require(msg.value >= minPrice, "you should pay over 0.01 ether" );

        uint256  tokenId = newTokenId;
        _safeMint( msg.sender, tokenId );

        _setTokenURI( tokenId, uri );
        newTokenId++;

        emit nftMint( msg.sender, tokenId, uri);
    }

    function withdraw() public onlyOwner {
        uint256 balance = address(this).balance;
        require( balance > 0, "No balance to withdraw");

        (bool success, ) = address(msg.sender).call{value: balance}("");
        require(success, "transfer failed");

        emit nftTransfer( balance, msg.sender);
    }
    
    function setMinPrice(uint256 price)public onlyOwner {
        minPrice = price;
    }
    /*
    * 由于同时继承 ERC721 和 ERC721URI 其中同时存在tokenURI 和 supportsInterface 函数
    * 为了解决冲突需要重写override(ERC721, ERC721URIStorage) 这两个函数
    * super 按照继承顺序查找父合约 调用找到的第一个父合约的 tokenURI 函数
    */
    function tokenURI(uint256 tokenId) public view override(ERC721, ERC721URIStorage) returns (string memory) {
        return super.tokenURI(tokenId);
    }

    function supportsInterface(bytes4 interfaceId) public view virtual override(ERC721, ERC721URIStorage) returns (bool) {
        return super.supportsInterface(interfaceId);
    }
}