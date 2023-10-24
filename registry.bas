Function Initialize() Uint64
10 STORE("CEO",HEX(SCID()))
20 SEND_ASSET_TO_ADDRESS(SIGNER(),1,SCID())
30 STORE("QUORUM",0)
40 STORE("EPOCH-INIT",1691607914)
70 STORE("newCollectionFee",0)
80 STORE("OAO_VERSION","PI")
90 STORE("OAO_NAME","Dero Asset & Address Registry")
99 RETURN 0
End Function

Function NewCollection(name String, owner String, ownerType Uint64, asset1 String, asset2 String, price1 Uint64, price2 Uint64, return1 Uint64, return2 Uint64, collectionType Uint64) Uint64
10 IF EXISTS("c"+name+"Owner")==0 THEN GOTO 15
11 IF LOAD("c"+name+"OwnerType") == 0 && ASSETVALUE(HEXDECODE(LOAD("c"+name+"Owner"))) ! = 1 THEN GOTO 100
12 SEND_ASSET_TO_ADDRESS(SIGNER(),ASSETVALUE(HEXDECODE(LOAD("c"+name+"Owner"))),HEXDECODE(LOAD("c"+name+"Owner")))
13 IF LOAD("c"+name+"OwnerType") == 1 && ADDRESS_STRING(SIGNER()) != LOAD("c"+name+"Owner") THEN GOTO 100
14 GOTO 20
15 IF DEROVALUE() ! = LOAD("newCollectionFee") THEN GOTO 100
16 smrtStore("treasury0000000000000000000000000000000000000000000000000000000000000000",DEROVALUE())
20 STORE("c"+name+"Owner",owner)
30 STORE("c"+name+"OwnerType",ownerType)
40 STORE("c"+name+"Asset1",asset1)
50 STORE("c"+name+"Asset2",asset2)
60 STORE("c"+name+"Price1",price1)
70 STORE("c"+name+"Price2",price2)
80 STORE("c"+name+"Return1",return1)
90 STORE("c"+name+"Return2",return2)
95 STORE("c"+name+"Type",collectionType)
99 RETURN 0
100 RETURN 1
End Function

Function smrtStore(key String, value Uint64) Uint64
10 IF EXISTS(key) THEN GOTO 20
11 RETURN STORE(key,value)
20 RETURN STORE(key,LOAD(key) + value)
End Function
/*
 0%3 = NO ADDRESSES
1%3  = NO ASSETS
<6 - NO MULTI
<3 (MOD 6) IMMUTABLE
*/
Function handleDel(collection String, scid String,T String) Uint64
10 IF EXISTS("n"+collection+scid) == 0 THEN GOTO 20
11 IF EXISTS(T+collection+LOAD("n"+collection+scid)) == 0 THEN GOTO 20
12 DELETE(T+collection+LOAD("n"+collection+scid))
20 RETURN 0
End Function

Function RegisterAsset(collection String, name String, scid String, index Uint64) Uint64
1 IF EXISTS("a"+collection+name) THEN GOTO 100
2 IF EXISTS("n"+collection+scid) && LOAD("c"+collection+"Type")%6 >2 || LOAD("c"+collection + "Type")%3 == 1 THEN GOTO 100
3 IF ASSETVALUE(HEXDECODE(scid)) != 1 THEN GOTO 100
4 SEND_ASSET_TO_ADDRESS(SIGNER(),1,HEXDECODE(scid))
6 IF checkTokens(collection,0) THEN GOTO 100
7 STORE("a"+collection+name,scid)
8 IF LOAD("c"+collection+"Type")<6 THEN GOTO 11
9 STORE("n"+collection+scid+index,name)
10 RETURN 0
11 handleDel(collection,scid,"a")
12 STORE("n"+collection+scid,name)
99 RETURN 0
100 RETURN 1
End Function

Function RegisterAddress(collection String,name String) Uint64
5 IF EXISTS("d"+collection+name) THEN GOTO 100
10 IF EXISTS("n"+collection+ADDRESS_STRING(SIGNER())) && LOAD("c"+collection+"Type")%6 >2 || LOAD("c"+collection+"Type")%3 == 0 THEN GOTO 100
11 IF checkTokens(collection,0) THEN GOTO 100
15 STORE("d"+collection+name,ADDRESS_STRING(SIGNER()))
16 IF LOAD("c"+collection+"Type") < 6 THEN GOTO 20
17 STORE("n"+collection+ADDRESS_STRING(SIGNER())+index,name)
18 RETURN 0
20 handleDel(collection,ADDRESS_STRING(SIGNER()),"d")
40 STORE("n"+collection+ADDRESS_STRING(SIGNER()),name)
99 RETURN 0
100 RETURN 1
End Function

Function checkTokens(collection String,flag Uint64) Uint64
10 LET flag = handleToken(collection, LOAD("c"+collection+"Asset1"),LOAD("c"+collection+"Return1"),LOAD("c"+collection+"Price1"),0) 
20 RETURN flag + handleToken(collection, LOAD("c"+collection+"Asset2"),LOAD("c"+collection+"Return2"),LOAD("c"+collection+"Price2"),0)
End Function

Function handleToken(collection String, token String, refund Uint64, price Uint64, amount Uint64) Uint64
1 IF token == "" THEN GOTO 21
2 LET amount = ASSETVALUE(HEXDECODE(token))
3 IF amount != price THEN GOTO 100
10 IF refund THEN GOTO 20
15 smrtStore("c"+collection+"Treasury"+token,MAX(1,amount*9/10))
16 smrtStore("treasury"+token,amount/10)
17 RETURN 0
20 SEND_ASSET_TO_ADDRESS(SIGNER(),amount,HEXDECODE(token))
21 RETURN 0
100 RETURN 1
End Function

Function UnRegisterAsset(name String, collection String, scid String) Uint64
10 IF EXISTS("n"+collection+scid) == 0 THEN GOTO 100
20 IF ASSETVALUE(HEXDECODE(scid)) !=1 THEN GOTO 100
30 DELETE("n"+collection+scid)
40 DELETE("a"+collection+name)
50 SEND_ASSET_TO_ADDRESS(SIGNER(),1,HEXDECODE(scid))
99 RETURN 0
100 RETURN 1
End Function

Function RateAsset(scid String, collection String, rating Uint64, feedback String) Uint64
1 IF EXISTS("n"+collection+scid) == 0 THEN GOTO 100
10 IF DEROVALUE() != 10000 THEN GOTO 100
20 STORE("r"+collection+scid+ADDRESS_STRING(SIGNER()),rating)
30 STORE("f"+collection+scid+ADDRESS_STRING(SIGNER()),comment)
40 smrtStore("treasury0000000000000000000000000000000000000000000000000000000000000000",10000)
99 RETURN 0
End Function

Function Deposit(token String) Uint64
1 STORE("treasury"+token,LOAD("treasury"+token)+ASSETVALUE(HEXDECODE(LOAD(token))))
2 RETURN 0
End Function

Function Withdraw(amount Uint64, token String, special Uint64) Uint64
1 IF ASSETVALUE(HEXDECODE(LOAD("CEO"))) != 1 THEN GOTO 99
2 SEND_ASSET_TO_ADDRESS(SIGNER(),1,HEXDECODE(LOAD("CEO")))
3 IF special ==1 THEN GOTO 20
4 IF amount > LOAD("treasury"+token) THEN GOTO 99
5 IF BLOCK_TIMESTAMP() < LOAD("allowanceRefresh"+token) THEN GOTO 8
6 STORE("allowanceRefresh"+token,BLOCK_TIMESTAMP()+LOAD("allowanceInterval"+token))
7 STORE("allowanceUsed"+token,0)
8 IF amount + LOAD("allowanceUsed"+token) > LOAD("allowance"+token) THEN GOTO 99
9 SEND_ASSET_TO_ADDRESS(SIGNER(),amount,HEXDECODE(LOAD(token)))
10 STORE("allowanceUsed"+token,LOAD("allowanceUsed"+token)+amount)
11 STORE("treasury"+token,LOAD("treasury"+token)-amount)
19 RETURN 0
20 IF LOAD("allowanceSpecial"+token) > LOAD("treasury"+token) THEN GOTO 99
21 SEND_ASSET_TO_ADDRESS(SIGNER(),LOAD("allowanceSpecial"+token),HEXDECODE(LOAD(token)))
22 STORE("treasury"+token,LOAD("treasury"+token)-LOAD("allowanceSpecial"+token))
23 DELETE("allowanceSpecial"+token)
98 RETURN 0
99 RETURN 1
End Function

Function Propose(hash String, k String, v String, t String, seat Uint64) Uint64
10 IF ASSETVALUE(HEXDECODE(LOAD("CEO"))) != 1 THEN GOTO 13
11 SEND_ASSET_TO_ADDRESS(SIGNER(),1,HEXDECODE(LOAD("CEO")))
12 GOTO 15
13 IF ASSETVALUE(HEXDECODE(LOAD("seat"+seat))) !=1 THEN GOTO 100
14 SEND_ASSET_TO_ADDRESS(SIGNER(),1,HEXDECODE(LOAD("seat"+seat)))
15 STORE("APPROVE", 0)
20 IF hash =="" THEN GOTO 40
25 STORE("HASH",hash)
30 STORE("k","")
35 RETURN 0
40 STORE("k",k)
45 STORE("HASH","")
49 STORE("t",t)
80 STORE("v",v)
90 RETURN 0
100 RETURN 1
End Function

Function Approve(seat Uint64) Uint64
10 IF ASSETVALUE(HEXDECODE(LOAD("seat"+seat)))!=1 THEN GOTO 100
20 STORE("APPROVE",LOAD("APPROVE")+1)
30 STORE("trustee"+seat,ADDRESS_STRING(SIGNER()))
99 RETURN 0
100 RETURN 1
End Function

Function ClaimSeat(seat Uint64) Uint64
10 IF ADDRESS_STRING(SIGNER())!= LOAD("trustee"+seat) THEN GOTO 100
20 SEND_ASSET_TO_ADDRESS(SIGNER(),1,HEXDECODE(LOAD("seat"+seat)))
30 IF LOAD("APPROVE") == 0 THEN GOTO 99
40 STORE("APPROVE",LOAD("APPROVE")-1)
99 RETURN 0
100 RETURN 1
End Function

Function Update(code String) Uint64
10 IF ASSETVALUE(HEXDECODE(LOAD("CEO")))!=1 THEN GOTO 100
15 SEND_ASSET_TO_ADDRESS(SIGNER(),1,HEXDECODE(LOAD("CEO")))
20 IF SHA256(code) != HEXDECODE(LOAD("HASH")) THEN GOTO 100
30 IF LOAD("APPROVE") < LOAD("QUORUM") THEN GOTO 100
40 UPDATE_SC_CODE(code)
99 RETURN 0
100 RETURN 1
End Function

Function Store() Uint64
10 IF LOAD("APPROVE") < LOAD("QUORUM") THEN GOTO 100
20 STORE("APPROVE",0)
30 IF LOAD("t") == "U" THEN GOTO 60
40 STORE(LOAD("k"), LOAD("v"))
45 STORE("k","")
50 RETURN 0
60 STORE(LOAD("k"),ATOI(LOAD("v")))
65 STORE("k","")
99 RETURN 0
100 RETURN 1
End Function