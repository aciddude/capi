#### The below shows the functionality of the coind package


```
	// Load config file
	configFile := coind.LoadConfig("./coind/config/config.json")
	
	// Create new coin daemon
	coinDaemon, err := coind.New(configFile)
	if err != nil {
		fmt.Print(err)
	}

//////------------------- Basic Calls ---------------------
    
    	///// How to call getblockchaininfo
    
    	getbcInfo, err := coinDaemon.GetBlockchainInfo()
    	if err != nil {
    		fmt.Printf("E0005: Could not get info from coinDaemon \n Daemon Error: %s ", err)
    	}
    	fmt.Println("\nGetBlockchainInfo \n")
    	jsongetbcinfo, _ := json.MarshalIndent(getbcInfo, "", " ")
    	fmt.Printf("%v \n", string(jsongetbcinfo))
    
    	///// How to call getinfo
    
    	getInfo, err := coinDaemon.GetInfo()
    	if err != nil {
    		fmt.Printf("E0005: Could not get info from coinDaemon \n Daemon Error: %s ", err)
    	}
    	fmt.Println("GetInfo\n")
    	jsongetinfo, _ := json.MarshalIndent(getInfo, "", " ")
    	fmt.Printf("%v \n", string(jsongetinfo))
    
    	///// How to call getmininginfo
    
    	//getminingInfo, err := coinDaemon.GetMiningInfo()
    	if err != nil {
    		fmt.Printf("E0005: Could not get info from coinDaemon \n Daemon Error: %s ", err)
    	}
    	fmt.Println("\nGetMiningInfo")
    	jsongetmininginfo, _ := json.MarshalIndent(getminingInfo, "", " ")
        fmt.Printf("%v \n", string(jsongetmininginfo))
    
    	///// How to call  getnetworkinfo
    
    	getnetworkInfo, err := coinDaemon.GetNetworkInfo()
    	if err != nil {
    		fmt.Printf("E0005: Could not get info from coinDaemon \n Daemon Error: %s ", err)
    	}
    	fmt.Println("\nGetNetworkInfo")
    	jsongetnetworkinfo, _ := json.MarshalIndent(getnetworkInfo, "", " ")
    	fmt.Printf("%v \n", string(jsongetnetworkinfo))
    
    	///// How to get the block count / current height
    
    	blockcount, err := coinDaemon.GetBlockCount()
    	if err != nil {
    		fmt.Printf("ERROR", err)
    	}
    	fmt.Printf("%d", blockcount)
    
    	///// How to get a single raw transaction
    	rawtx, err := coinDaemon.GetRawTransaction("fe94a3a70d888a491395d446c67aad546c756cc3ba86511dfc1e63f93ed97cc6", true)
    	if err != nil {
    		fmt.Printf("ERROR %v", err)
    	}
    	jsonrawtx, err := json.MarshalIndent(rawtx, "", " ")
    	if err != nil {
    		fmt.Printf("Marshal Error: %v", err)
    	}
    	fmt.Printf("RAW TX \n%s", jsonrawtx)



//////------------------- Advance Example 1: Get an array of blocks ---------------------

	//// 1 -  specify the start and end height of the blocks array you want back

	startHeight := 25
	endHeight := 50 

	//// 2 - listsize is the size of the array that will be returned, if startHeight = 25 and EndHeight = 50  your listsize = 25

	listsize := endHeight - startHeight

	//// 3 - make an array of getblockhash requests using the startHeight and endHeight

	getblockhashrequest, _ := coind.MakeBlockHashListRequest(startHeight, endHeight)

	//// 4 -  calls the coind daemon and provides the array above

	getblockreponse, _ := coinDaemon.GetBlockHashList(getblockhashrequest)

	//// 5 -  parse the getblockhash reponse,  removes quotes and returns an array []string

	hashlist, err := coind.Parsehashlist(getblockreponse)
	if err != nil {
		log.Println(err)
	}

	//// 6 - make an array of getblock requests using the listsize and the hashlist

	getblockrequest, err := coinDaemon.MakeGetBlockListRequest(listsize, hashlist)
	if err != nil {
		fmt.Printf("ERROR :", err)
	}

	//// 7 - call the daemon and get an array of blocks back using your array of getblock requests

	blocklist, err := coinDaemon.GetBlockList(getblockrequest)
	if err != nil {
		fmt.Errorf("ERROR :", err)
	
	}
	//// Marshall and indent your list to make it look pretty, and print it

	jsonblocklist, _ := json.MarshalIndent(blocklist, "", " ")
	fmt.Printf("%s", jsonblocklist)
	
	
 //////------------------- Advance Example 2: Get an array of transactions ---------------------
	
	//// 8 -  parse the blocklist for tx hashes, returns an array []string
    	txlist, err := coind.ParseBlockTX(jsonblocklist)
    
        //// 9 - build an array of getrawtransaction requests using the txlist/transaction hashes from above
    	makerequest, err := coinDaemon.MakeRawTxListRequest(txlist)
    	if err != nil {
    		fmt.Printf("ERROR :", err)
    	}
    
        //// 10 - make a getrawtransaction request using the array from above
    	rawtxlist, err := coinDaemon.GetRawTransactionList(makerequest)
    	if err != nil {
    		fmt.Printf("ERROR:\nRaw Transacaction List Request %v ", err)
    	}
        
        //// 11 indent the json to make it look pretty and print
    	prettyRawtx, _ := json.MarshalIndent(rawtxlist, "", " ")
    
    	fmt.Printf("%s", prettyRawtx)
```
