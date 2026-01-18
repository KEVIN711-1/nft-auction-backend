// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract KevinNFT is ERC721, ERC721URIStorage, Ownable {
    
    uint256 private _nextTokenId = 1;
    uint256 public mintPrice;
    uint256 public maxSupply;
    bool public isMintingEnabled = true;
    
    event NFTMinted(address indexed owner, uint256 indexed tokenId, string uri);
    event PriceUpdated(uint256 newPrice);
    event MintingToggled(bool enabled);
    event Withdrawn(address indexed owner, uint256 amount);
    
    constructor(
        string memory name,
        string memory symbol,
        uint256 _mintPrice,
        uint256 _maxSupply
    ) ERC721(name, symbol) Ownable(msg.sender) {
        mintPrice = _mintPrice;
        maxSupply = _maxSupply;
    }
    
    // ==================== 铸造函数 ====================
    
    function mintNFT(string calldata uri) public payable returns (uint256) {
        require(isMintingEnabled, "Minting is disabled");
        require(msg.value >= mintPrice, "Insufficient payment");
        require(_nextTokenId <= maxSupply, "Max supply reached");
        
        uint256 tokenId = _nextTokenId;
        _safeMint(msg.sender, tokenId);
        _setTokenURI(tokenId, uri);
        
        emit NFTMinted(msg.sender, tokenId, uri);
        
        _nextTokenId++;
        return tokenId;
    }
    
    function ownerMint(string calldata uri, address recipient) public onlyOwner returns (uint256) {
        require(_nextTokenId <= maxSupply, "Max supply reached");
        
        uint256 tokenId = _nextTokenId;
        _safeMint(recipient, tokenId);
        _setTokenURI(tokenId, uri);
        
        emit NFTMinted(recipient, tokenId, uri);
        
        _nextTokenId++;
        return tokenId;
    }
    
    // ==================== 管理函数 ====================
    
    function setMintPrice(uint256 newPrice) public onlyOwner {
        mintPrice = newPrice;
        emit PriceUpdated(newPrice);
    }
    
    function toggleMinting(bool enabled) public onlyOwner {
        isMintingEnabled = enabled;
        emit MintingToggled(enabled);
    }
    
    function withdraw() public onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No funds to withdraw");
        
        (bool success, ) = payable(owner()).call{value: balance}("");
        require(success, "Withdrawal failed");
        
        emit Withdrawn(owner(), balance);
    }
    
    // ==================== 视图函数 ====================
    
    function getNextTokenId() public view returns (uint256) {
        return _nextTokenId;
    }
    
    function totalSupply() public view returns (uint256) {
        return _nextTokenId - 1;
    }
    
    function remainingSupply() public view returns (uint256) {
        return maxSupply - totalSupply();
    }
    
    function getStats() public view returns (
        uint256 currentTokenId,
        uint256 currentSupply,
        uint256 remaining,
        uint256 price,
        bool mintingEnabled
    ) {
        return (
            _nextTokenId,
            totalSupply(),
            remainingSupply(),
            mintPrice,
            isMintingEnabled
        );
    }
    
    // ==================== 重写函数（简化版） ====================
    
    // 只重写tokenURI和supportsInterface
    function tokenURI(uint256 tokenId)
        public
        view
        virtual
        override(ERC721, ERC721URIStorage)
        returns (string memory)
    {
        return super.tokenURI(tokenId);
    }
    
    function supportsInterface(bytes4 interfaceId)
        public
        view
        virtual
        override(ERC721, ERC721URIStorage)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
    
    // 移除_burn的重写，因为ERC721URIStorage已经处理了
}