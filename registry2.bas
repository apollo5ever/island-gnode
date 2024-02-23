Function Initialize() Uint64
5 STORE("CEO","c5b5778d4c297e172c4cd318d5db1a3ffef70e1ec03c0051181beb19ec44b59a")
6 STORE("modify","c5b5778d4c297e172c4cd318d5db1a3ffef70e1ec03c0051181beb19ec44b59a")
7 STORE("q.modify",0)
8 STORE("q.update",0)
9 STORE("q.del",0)
10 RETURN 0
End Function

Function NewCollection(name String, id String, mutable Uint64, namesPerMember Uint64, tokens String, fees String, returns String, CEO String, board String) Uint64
10 dim c as String
20 LET c = ITOA(add("collectionCount",1))
30 IF storeNew(c+"name",name,"S") THEN GOTO 1000
40 STORE(c+"id",id)
45 STORE(c+"npm",namesPerMember)
50 STORE(c+"CEO",CEO)
55 STORE(c+"modify",CEO)
60 newBoard(c,board, STRLEN(board)/64)
70 STORE(c+"q.allowance",0)
80 STORE(c+"q.birth",0)
90 STORE(c+"q.modify",0)
100 STORE(c+"q.del",0)
110 STORE(c+".tokens",tokens)
120 STORE(c+".fees",fees)
130 STORE(c+".returns",returns)
999 RETURN 0
1000 RETURN 1
End Function

Function TransferName(c String, name String, member String) Uint64
10 IF notOwner(c+"m."+name) THEN GOTO 100
20 RETURN registerMember(c,name,member,"asset",0)
100 RETURN 1
End Function

Function notOwner(member String) Uint64
10 IF EXISTS(member+".asset") == 0 THEN GOTO 30
20 IF ASSETVALUE(HEXDECODE(LOAD(member+".asset"))) == 1 THEN GOTO 45
30 IF SIGNER() == ADDRESS_RAW(LOAD(member+".address")) THEN GOTO 50
40 RETURN 1
45 SEND_ASSET_TO_ADDRESS(SIGNER(),1,HEXDECODE(LOAD(member+".asset")))
50 RETURN 0
End Function

Function registerMember(c String, name String, ceoAddress String, ceoToken String, days Uint64) Uint64
10 STORE(c+"m."+name+".ceoToken",ceoToken)
15 STORE(c+"m."+name+".modify",ceoToken)
20 STORE(c+"m."+name+".ceoAddress",ceoAddress)
25 LET days = getExpiry(c,name,days)
30 STORE(c+"m."+name+".expiry",MAX(BLOCK_TIMESTAMP()+days,days))
90 RETURN 0
End Function

Function RenewMember(c String, name String, days Uint64) Uint64
25 getExpiry(c,name,days)
90 RETURN 0
End Function

Function RegisterMember(c String, name String, id String, ceoAddress String, ceoToken String, board String, days Uint64) Uint64
10 IF storeNew(c+"m."+name+".id", id,"S") THEN GOTO 100
20 registerMember(c,name,ceoAddress, ceoToken, days)
90 RETURN 0
100 RETURN 1
End Function

Function getVotes(c String, name String, action String) Uint64
10 IF name != "" THEN GOTO 100
20 IF c !="" THEN GOTO 50
30 RETURN LOAD("v."+action) >= LOAD("q."+action)
50 IF EXISTS(c+"q."+action) THEN GOTO 70
60 RETURN 1
70 RETURN LOAD(c+"v."+action) >= LOAD(c+"q."+action)
100 IF EXISTS(c+"m."+name+".q."+action) THEN GOTO 120
110 RETURN 1
120 RETURN LOAD(c+"m."+name+".q."+action) >= LOAD(c+"m."+name+".q."+action)
End Function

Function getPermission(c String, name String, role String, action String) Uint64
10 IF name != "" THEN GOTO 100
20 IF c != "" THEN GOTO 50
30 RETURN roleCheck(LOAD(role), STRLEN(role)) && getVotes(c,name,action)
50 IF EXISTS(c+role) THEN GOTO 70
60 RETURN roleCheck(LOAD(c+"CEO")) && getVotes(c,name,action)
70 RETURN roleCheck(LOAD(c+role)) && getVotes(c,name,action)
100 IF EXISTS(c+"m."+name+"."+role) THEN GOTO 120
110 RETURN roleCheck(LOAD(c+"m."+name+".CEO")) && getVotes(c,name,action)
120 RETURN roleCheck(LOAD(c+"m."+name+"."+role)) && getVotes(c,name,action)
End Function

Function roleCheck(tokens String, i Uint64) Uint64
10 IF i == 0 THEN GOTO 90
20 LET i= i-1
30 IF ASSETVALUE(HEXDECODE(SUBSTR(tokens,i*64,64))) == 0 THEN GOTO 10
35 SEND_ASSET_TO_ADDRESS(SIGNER(),1,HEXDECODE(SUBSTR(tokens,i*64,64)))
40 RETURN 1
90 RETURN 0
End Function

Function UpdateAddress(c String, k String, v String) Uint64
10 STORE(c+ADDRESS_STRING(SIGNER())+"["+k+"]",v)
20 RETURN 0
End Function

Function RateMember(c String, member String, rating Uint64, review String) Uint64
10 STORE("rating."+member+BLOCK_TIMESTAMP(),rating)
20 STORE("review."+member+BLOCK_TIMESTAMP(),review)
99 RETURN 0
End Function

Function getExpiry(c String, name String,days Uint64) Uint64
10 RETURN add(c+"m."+name+".expiry",pmtCheck(LOAD(c+".tokens"),LOAD(c+".fees"),LOAD(c+".returns"),STRLEN(LOAD(c+".tokens"))/64,days)*86400)
End Function

Function storeNew(k String, v String, t String) Uint64
10 IF EXISTS(k) THEN GOTO 100
20 IF t == "U" THEN GOTO 50
30 STORE(k,v)
40 RETURN 0
50 STORE(k,ATOI(v))
60 RETURN 0
100 RETURN 1
End Function

Function newBoard(c String, board String, i Uint64) Uint64
10 IF i == 0 THEN GOTO 90
20 LET i = i -1
30 IF STORE(c+"seat"+i,SUBSTR(board,i*64,64)) THEN GOTO 10
90 RETURN 0
End Function

Function pmtCheck(tokens String, fees String, returns String, i Uint64, days Uint64) Uint64
10 IF i == 0 THEN GOTO 90
20 LET i = i -1
30 IF ASSETVALUE(HEXDECODE(SUBSTR(tokens,i*64,64))) < ATOI(SUBSTR(fees,i*18,18))*days THEN GOTO 100
//35 STORE("DEPOSIT"+SUBSTR(tokens,i*64,64),ASSETVALUE(HEXDECODE(SUBSTR(tokens,i*64,64))))
40 IF SEND_ASSET_TO_ADDRESS(SIGNER(),ASSETVALUE(HEXDECODE(SUBSTR(tokens,i*64,64)))*ATOI(SUBSTR(returns,i,1)),HEXDECODE(SUBSTR(tokens,i*64,64))) THEN GOTO 10
90 RETURN days
100 RETURN 1
End Function

Function add(k String, v Uint64) Uint64
10 IF EXISTS(k) THEN GOTO 30
15 STORE(k,v)
20 RETURN(LOAD(k))
30 STORE(k,LOAD(k)+v)
35 RETURN LOAD(k)
End Function

Function Deposit(token String, c String) Uint64
1 add(c+"treasury"+token,ASSETVALUE(HEXDECODE(token)))
2 RETURN 0
End Function

Function Withdraw(amount Uint64, token String, special Uint64, c String) Uint64
1 IF ASSETVALUE(HEXDECODE(LOAD(c+"CEO"))) != 1 THEN GOTO 99
2 SEND_ASSET_TO_ADDRESS(SIGNER(),1,HEXDECODE(LOAD(c+".CEO")))
3 IF special ==1 THEN GOTO 20
4 IF amount > LOAD(c+"treasury"+token) THEN GOTO 99
5 IF BLOCK_TIMESTAMP() < LOAD(c+"allowanceRefresh"+token) THEN GOTO 8
6 STORE(c+"allowanceRefresh"+token,BLOCK_TIMESTAMP()+LOAD(c+"allowanceInterval"+token))
7 STORE(c+"allowanceUsed"+token,0)
8 IF amount + LOAD(c+"allowanceUsed"+token) > LOAD(c+"allowance"+token) THEN GOTO 99
9 SEND_ASSET_TO_ADDRESS(SIGNER(),amount,HEXDECODE(token))
10 add(c+"allowanceUsed"+token,amount)
11 STORE(c+"treasury"+token,LOAD(c+"treasury"+token)-amount)
19 RETURN 0
20 IF LOAD(c+"allowanceSpecial"+token) > LOAD(c+"treasury"+token) THEN GOTO 99
21 SEND_ASSET_TO_ADDRESS(SIGNER(),LOAD(c+"allowanceSpecial"+token),HEXDECODE(token))
22 STORE(c+"treasury"+token,LOAD(c+"treasury"+token)-LOAD(c+"allowanceSpecial"+token))
23 DELETE(c+"allowanceSpecial"+token)
98 RETURN 0
99 RETURN 1
End Function

Function Propose(hash String, k String, v String, t String, seat Uint64, c String) Uint64
10 IF roleCheck(LOAD(c+"modify"),STRLEN(LOAD(c+"modify"))/64) THEN GOTO 13
12 GOTO 15
13 IF ASSETVALUE(HEXDECODE(LOAD(c+"seat"+seat))) !=1 THEN GOTO 100
14 SEND_ASSET_TO_ADDRESS(SIGNER(),1,HEXDECODE(LOAD(c+"seat"+seat)))
15 STORE(c+"APPROVE", 0)
20 IF hash =="" THEN GOTO 40
25 STORE(c+"HASH",hash)
30 STORE(c+"k","")
35 RETURN 0
40 STORE(c+"k",k)
45 STORE(c+"HASH","")
49 STORE(c+"t",t)
80 STORE(c+"v",v)
90 RETURN 0
100 RETURN 1
End Function

Function Approve(seat Uint64,c String) Uint64
10 IF ASSETVALUE(HEXDECODE(LOAD(c+"seat"+seat)))!=1 THEN GOTO 100
20 STORE(c+"APPROVE",LOAD(c+"APPROVE")+1)
30 STORE(c+"trustee"+seat,ADDRESS_STRING(SIGNER()))
99 RETURN 0
100 RETURN 1
End Function

Function ClaimSeat(seat Uint64, c String) Uint64
10 IF ADDRESS_STRING(SIGNER())!= LOAD(c+"trustee"+seat) THEN GOTO 100
20 SEND_ASSET_TO_ADDRESS(SIGNER(),1,HEXDECODE(LOAD(c+"seat"+seat)))
30 IF LOAD(c+"APPROVE") == 0 THEN GOTO 99
40 STORE(c+"APPROVE",LOAD(c+"APPROVE")-1)
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

Function Store(c String, name String) Uint64
10 IF getPermission(c,name,"modify","modify") == 0 THEN GOTO 100
20 STORE(c+"APPROVE",0)
30 IF LOAD(c+"t") == "U" THEN GOTO 60
40 STORE(c+LOAD(c+"k"), LOAD(c+"v"))
45 STORE(c+"k","")
50 RETURN 0
60 STORE(LOAD(c+"k"),ATOI(LOAD(c+"v")))
65 STORE(c+"k","")
99 RETURN 0
100 RETURN 1
End Function

Function SetModifiers(m String, c String) Uint64
10 IF roleCheck(LOAD(c+"CEO"),1) THEN GOTO 100
20 STORE(c+"modify",m)
90 RETURN 0
100 RETURN 1
End Function