Function NF(H String, i Uint64, name String, image String, tagline String, desc String, Goal Uint64, Deadline Uint64, WithdrawlType Uint64, Recipient String, ICO Uint64, icoToken String) Uint64
10 IF ASSETVALUE(HEXDECODE(H)) != 1 THEN GOTO 100
20 IF EXISTS(H+i+"_G") THEN GOTO 100
30 SetMetadata(H,i,name,image,tagline,desc)
31 STORE(H+i+"_G",Goal)
32 STORE(H+i+"_D",Deadline)
33 STORE(H+i+"_F",Recipient)
34 STORE(H+i+"_R",0)
35 STORE(H+i+"_C",0)
36 STORE(H+i+"WithdrawlType",WithdrawlType)
37 STORE(H+i+"ICO",ICO)
38 IF ICO == 0 THEN GOTO 99
39 STORE(H+i+"icoAmount",ASSETVALUE(HEXDECODE(icoToken)))
40 STORE(H+i+"icoToken",icoToken)
99 RETURN 0
100 RETURN 1
End Function

Function add(key String, value Uint64) Uint64
10 IF EXISTS(key) THEN GOTO 20
11 STORE(key,value)
12 RETURN 0
20 STORE(key,LOAD(key)+value)
25 RETURN 0
End Function    

Function SG(H String,Refundable Uint64) Uint64
10 IF LOAD(H+"ICO") == 0 THEN GOTO 15
11 add(H+ADDRESS_STRING(SIGNER())+"ICO",DEROVALUE())
15 add(H+"_R",DEROVALUE())
20 IF Refundable ==1 THEN GOTO 40
21 IF LOAD(H+"WithdrawlType") == 1 THEN GOTO 30
22 SEND_DERO_TO_ADDRESS(ADDRESS_RAW(LOAD(H+"_F")),DEROVALUE()*9/10)
23 add(H+"_C",DEROVALUE())
24 add("treasuryDERO",DEROVALUE()/10)
25 RETURN 0
30 add(H+"Claimable",DEROVALUE())
35 RETURN 0
40 add(H+ADDRESS_STRING(SIGNER()),DEROVALUE())
45 RETURN 0
End Function

Function OAOWithdrawFromFundraiser(H String, amount Uint64) Uint64
1 IF LOAD(H+"WithdrawlType") == 0 THEN GOTO 100
10 IF BLOCK_TIMESTAMP() > LOAD(H+"_D") && LOAD(H+"_R")>= LOAD(H+"_G") && LOAD(H+"_C") < LOAD(H+"_R") THEN GOTO 20
11 LET amount = MIN(LOAD(H+"Claimable"),ASSETVALUE(HEXDECODE(LOAD(H+"_F"))))
12 SEND_DERO_TO_ADDRESS(SIGNER(),amount*9/10)
13 add("treasuryDERO",amount/10)
14 STORE(H+"Claimable",LOAD(H+"Claimable")-amount)
15 add(H+"_C",amount)
16 RETURN 0
20 LET amount = MIN(ASSETVALUE(HEXDECODE(LOAD(H+"_F"))),LOAD(H+"_R")-LOAD(H+"_C"))
21 SEND_DERO_TO_ADDRESS(SIGNER(),amount*9/10)
22 STORE(H+"Claimable",0)
23 add(H+"_C",amount)
24 add("treasuryDERO",amount/10)
25 RETURN 0
100 RETURN 1
End Function

Function GetTokens(H String) Uint64
10 IF BLOCK_TIMESTAMP() < LOAD(H+"_D") THEN GOTO 100
11 IF LOAD(H+"ICO") == 0 THEN GOTO 100
12 IF LOAD(H+"_R") < LOAD(H+"_G") THEN GOTO 100
20 SEND_ASSET_TO_ADDRESS(SIGNER(),LOAD(H+icoAmount)*LOAD(H+ADDRESS_STRING(SIGNER())+"ICO")/LOAD(H+"_R"),HEXDECODE(LOAD(H+"icoToken")))
30 STORE(H+ADDRESS_STRING(SIGNER())+"ICO",0)
99 RETURN 0
100 RETURN 1
End Function

Function SetName(H String,i Uint64, Name String) Uint64
10 IF ASSETVALUE(HEXDECODE(H)) != 1 THEN GOTO 100
40 SEND_ASSET_TO_ADDRESS(SIGNER(),1,HEXDECODE(H))
50 STORE(H+i+"Name",Name)
99 RETURN 0
100 RETURN 1
End Function

Function SetImage(H String,i Uint64, Image String) Uint64
10 IF ASSETVALUE(HEXDECODE(H)) != 1 THEN GOTO 100
40 SEND_ASSET_TO_ADDRESS(SIGNER(),1,HEXDECODE(H))
50 STORE(H+i+"Image",Image)
99 RETURN 0
100 RETURN 1
End Function

Function SetTagline(H String, i Uint64, Tagline String) Uint64
10 IF ASSETVALUE(HEXDECODE(H)) != 1 THEN GOTO 100
40 SEND_ASSET_TO_ADDRESS(SIGNER(),1,HEXDECODE(H))
50 STORE(H+i+"Tagline",Tagline)
99 RETURN 0
100 RETURN 1
End Function

Function SetDescription(H String, i Uint64, Description String) Uint64
10 IF ASSETVALUE(HEXDECODE(H)) != 1 THEN GOTO 100
40 SEND_ASSET_TO_ADDRESS(SIGNER(),1,HEXDECODE(H))
50 STORE(H+i+"Desc",Description)
99 RETURN 0
100 RETURN 1
End Function

Function SetMetadata(H String, i Uint64, Name String, Image String, Tagline String, Description String) Uint64
10 IF ASSETVALUE(HEXDECODE(H)) != 1 THEN GOTO 100
40 STORE(H+i+"Image",Image)
50 STORE(H+i+"Tagline",Tagline)
60 STORE(H+i+"Desc",Description)
70 STORE(H+i+"Name",Name)
80 SEND_ASSET_TO_ADDRESS(SIGNER(),1,HEXDECODE(H))
99 RETURN 0
100 RETURN 1
End Function

Function WFF(H String, i Uint64) Uint64
10 IF EXISTS(H+i+"_D") == 0 THEN GOTO 100
15 IF LOAD(H+i+"withdrawlType") == 1 THEN GOTO 100
20 IF LOAD(H+i+"_D") > BLOCK_TIMESTAMP() THEN GOTO 65
30 IF LOAD(H+i+"_R") >= LOAD(H+i+"_G") THEN GOTO 70
40 IF EXISTS(H+i+ADDRESS_STRING(SIGNER())) == 0 THEN GOTO 100 
50 SEND_DERO_TO_ADDRESS(SIGNER(),LOAD(H+i+ADDRESS_STRING(SIGNER()))) 
56 DELETE(H+i+ADDRESS_STRING(SIGNER())) 
60 RETURN 0
65 IF LOAD(H+i+"_R") < LOAD(H+i+"_G") THEN GOTO 100
70 SEND_DERO_TO_ADDRESS(ADDRESS_RAW(LOAD(H+i+"_F")), (LOAD(H+i+"_R")-LOAD(H+i+"_C"))*9/10) 
71 add("treasuryDERO",(LOAD(H+i+"_R")-LOAD(H+i+"_C"))/10)
75 STORE(H+i+"_C",LOAD(H+i+"_R")) 
99 RETURN 0
100 RETURN 1
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

Function SS(shares Uint64) Uint64
10 IF EXISTS(ADDRESS_STRING(SIGNER())+"_SHARES") == 0 THEN GOTO 100
20 IF LOAD(ADDRESS_STRING(SIGNER())+"_SHARES") < shares THEN GOTO 100
30 STORE(ADDRESS_STRING(SIGNER())+"_SHARES",LOAD(ADDRESS_STRING(SIGNER())+"_SHARES")-shares)
40 SEND_ASSET_TO_ADDRESS(SIGNER(),shares*10000,HEXDECODE(LOAD("COCO")))
50 STORE("T_COCO",LOAD("T_COCO")-shares*10000)
99 RETURN 0
100 RETURN 1
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