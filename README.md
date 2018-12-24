Capi v0.1
=============

Important!  Your coin.conf file needs to have `txindex=1`

capi will look for a config file in ./config/config.json

A sample config is provided below:
If you have the capi binary in ~/home/Username/capiv1.0/capi
1 - create a new config directory under  ~/home/Username/capiv1.0/config
2 - create a new file called `config.json` and paste the below contents and change it as you see fit.

```
{
  "Coin"        :         "Feathercoin",
  "Ticker"      :         "FTC",
  "Daemon"      :         "nomp.pools.ac:19337",
  "RPCUser"     :         "admin",
  "RPCPassword" :         "123456",
  "HTTPPostMode":         true,
  "DisableTLS"  :         true,
  "EnableCoinCodexAPI" :  true,
  "capi_port"          :  ":8000"
}

```
```
Coin =  Coin name as a string
Ticker = Coin symbol / price ticker as a string
Daemon =  Coin daemon URL+Port as a string
RPCUser =  Coin daemon RPC user as a string
RPCPassword = Coin daemon RPC Password as a string
HTTPPostMode = Enable HTTP only posts to coin daemon, bool true or false
DisableTLS = DisableTLS connection to coin daemon, bool true or false
EnableCoinCodexAPI = Enable or Disable the CoinCodexApi data which is used for prices, bool true or false
capi_port = The port capi will listen on, string of ":PORT" . eg :8000 for capi to listen on port 8000
```

Once capi is running you can hit the following api endpoints:

`/tx/{TransactionID}`  
example call: http://127.0.0.1:8000/tx/e473d72f183f25e0f1cb97ab2b977bc98039cb0ac31e91fbb96a8257ff622bd5
example response:
```
{
"hex": "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff03510108ffffffff020050d6dc010000002321021a9f576cf7e5c0e3ad0c7c02e60aec4bb205362d240f6e5601c1c2bad0fc9374ac0000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf900000000",
"txid": "e473d72f183f25e0f1cb97ab2b977bc98039cb0ac31e91fbb96a8257ff622bd5",
"hash": "e473d72f183f25e0f1cb97ab2b977bc98039cb0ac31e91fbb96a8257ff622bd5",
"size": 145,
"vsize": 145,
"version": 1,
"locktime": 0,
"vin": [
{
"coinbase": "510108",
"sequence": 4294967295
}
],
"vout": [
{
"value": 80,
"n": 0,
"scriptPubKey": {
"asm": "021a9f576cf7e5c0e3ad0c7c02e60aec4bb205362d240f6e5601c1c2bad0fc9374 OP_CHECKSIG",
"hex": "21021a9f576cf7e5c0e3ad0c7c02e60aec4bb205362d240f6e5601c1c2bad0fc9374ac",
"reqSigs": 1,
"type": "pubkey",
"addresses": [
"mojQnATqdQhzKMGS6LnvBxdo1t2zTCTB2u"
]
}
},
{
"value": 0,
"n": 1,
"scriptPubKey": {
"asm": "OP_RETURN aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf9",
"hex": "6a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf9",
"type": "nulldata"
}
}
],
"blockhash": "afece97d0118541252714663e8a7c719a0baac65c0588dcd28fe83238d61ab16",
"confirmations": 5611,
"time": 1544779888,
"blocktime": 1544779888
}
```


`/block/{BlockHashorBlockHeight}`
example call 1: http://127.0.0.1:8000/block/0
example call 2: http://127.0.0.1:8000/block/4b6c3362e2f2a6b6317c85ecaa0f5415167e2bb333d2bf3d3699d73df613b91f
example response:
```
[
{
"hash": "4b6c3362e2f2a6b6317c85ecaa0f5415167e2bb333d2bf3d3699d73df613b91f",
"confirmations": 5612,
"size": 280,
"strippedSize": 280,
"weight": 1120,
"height": 0,
"version": 1,
"versionHex": "00000001",
"merkleRoot": "97ddfbbae6be97fd6cdf3e7ca13232a3afff2353e29badfab7f73011edd4ced9",
"tx": [
"97ddfbbae6be97fd6cdf3e7ca13232a3afff2353e29badfab7f73011edd4ced9"
],
"time": 1536656597,
"nonce": 529517,
"bits": "1e0ffff0",
"difficulty": 0.000244140625,
"previousBlockHash": "",
"nextBlockHash": "afece97d0118541252714663e8a7c719a0baac65c0588dcd28fe83238d61ab16"
}
]
```


`/market`
example call: http://127.0.0.1:8000/market
example response:
```
{
"symbol": "FTC",
"coin_name": "Feathercoin",
"today_open": 0.018926488,
"price_high_24_usd": 0.020271526,
"price_low_24_usd": 0.018661281,
"volume_24_usd": 7017.345357014,
"data_provider": "CoinCodex.com"
}
```
