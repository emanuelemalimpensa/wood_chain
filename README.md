



SI RACCOMANDA DI COPIARE LA DIRECTORY "listati" NELLA DIRECTORY "fabric-sample/test-network"
PER ESEGUIRE IL DEPLOY SEGUENDO I LISTATI 

$ tree -L 1
.
├── chaincode-go-contracts
├── chaincode-go-transport
├── chaincode-go-trunk
├── listati
├── scenario1
└── scenario2

La directory chaincode-go-contracts contiene lo smart contract per la gestione delle identita'
certificate. E' deployable sulla test-network.
La directory chaincode-go-trunk contiene lo smart contract per la gestione degli asset fisici
sulla blockchain. Per alcune funzioni e' necessario aver eseguito precedentemente il deploy
dello smart contract per la gestione delle identita' certficiate
E' deployable sulla test-network.
La directory chaincode-go-transport contiene lo smart contract per la gestione dei contratti
per il trasporto degli asset fisici.
E' deployable sulla test-network.

La directory listati contiene tre subfolder 
.
├── contracts
├── transport
└── tt

ognuna delle quali contiene i listati per l'invocazione delle funzioni degli smart contract utilizzando
la test-network.

La directory scenario1 
.
├── forest_corp
│   └── rest-api-go-id-and-trunks
└── id_corp
    ├── api_db
    ├── end_point_id_emitter
    └── rest-api-go-ids


contiene due subfolder ognuna delle quali presenta un servizio rest per l'interazione con la 
blockchain. id_corp richiede l'interazione con un db per la generazione degli id certificati,
e' composta da 3 rest server interagenti.
forest_corp contiene un server rest che espone i metodi per interagire con blockchain e 
richiedere un audit per l'ottenimento degli id certificati.

La directory scenario2 
.
├── rest-api-go-transport-forest
├── rest-api-go-transport-sawmill
└── rest-api-go-transport-shipping

contiene i tre server rest che permettono la creazione di un accordo per il trasferimento
di prodotti tra le organizzazioni. 




+++++++++++++++++++++++++++++++++++++++++++++++++++++++
	DEPLOY SU BLOCKCHAIN DI SCENARIO 1
+++++++++++++++++++++++++++++++++++++++++++++++++++++++
 ./network.sh down
 ./network.sh up createChannel -c mychannel
 source ./org1.sh
 peer channel list
 ./install_script/package.sh basic.tar.gz ../asset-transfer-basic/chaincode-go golang basic_1.0 **##AGGIUSTARE PATH 
 ./install_script/install.sh basic.tar.gz 
 source ./org2.sh 
 ./install_script/install.sh basic.tar.gz 
 source ./install_script/esporta_id_cc.sh basic
 ./install_script/check_commit.sh basic 1.0 mychannel
 ./install_script/approve_for_my_org.sh basic 1.0 mychannel
 source ./org1.sh 
 ./install_script/approve_for_my_org.sh basic 1.0 mychannel
 ./install_script/check_commit.sh basic 1.0 mychannel
 ./install_script/commit.sh basic 1.0 mychannel 1 	
 ./install_script/invoke.sh basic mychannel 1		
 ./install_script/get_all_assets.sh basic mychannel
 
 
 
 
+++++++++++++++++++++++++++++++++++++++++++++++++++++++
	DEPLOY SU BLOCKCHAIN DI SCENARIO 2
+++++++++++++++++++++++++++++++++++++++++++++++++++++++
SCENARIO 2
 ./network.sh down
 ./network.sh up createChannel -c mychannel
 cd addOrg3/
 ./addOrg3.sh up -c mychannel
 cd ..
 ./install_script/package.sh transport.tar.gz ../../smart_contract\&server/chaincode-go-transport/ golang transport_1.0 **##AGGIUSTARE PATH 
 source ./org1.sh
 ./install_script/install.sh transport.tar.gz 
 source ./install_script/esporta_id_cc.sh transport
 source ./org2.sh 
 ./install_script/install.sh transport.tar.gz
 source ./install_script/esporta_id_cc.sh transport  
 source ./org3.sh 
 ./install_script/install.sh transport.tar.gz 
 source ./install_script/esporta_id_cc.sh transport
 source ./org1.sh 
 ./install_script/approve_for_my_org.sh transport 1.0 mychannel
 source ./org3.sh 
 ./install_script/approve_for_my_org.sh transport 1.0 mychannel
 source ./org2.sh 
 ./install_script/approve_for_my_org.sh transport 1.0 mychannel
 ./install_script/commit.sh transport 1.0 mychannel 2
 ./install_script/invoke.sh transport mychannel 2
 ./install_script/get_all_assets.sh transport mychannel

 source ./org1.sh
 ./listati/transport/get_all_assets.sh transport mychannel
 ./listati/transport/create_asset.sh transport mychannel
 ./listati/transport/get_all_assets.sh transport mychannel
 ./listati/transport/update_seller.sh transport mychannel
 ./listati/transport/get_all_assets.sh transport mychannel
 ./listati/transport/update_buyer.sh transport mychannel
 ./listati/transport/get_all_assets.sh transport mychannel
 ./listati/transport/approva_seller.sh transport mychannel ##ERROR1***
 ./listati/transport/approva_buyer.sh transport mychannel  ##ERROR1***
 ./listati/transport/get_all_assets.sh transport mychannel ##NESSUNA MODIFICA APPORTATA
 source ./org3.sh
 ./listati/transport/approva_seller.sh transport mychannel ##ERROR2***
 source ./org2.sh
 ./listati/transport/approva_buyer.sh transport mychannel
 ./listati/transport/get_all_assets.sh transport mychannel ##MODIFICHE APPORTATE
 source ./org3.sh
 ./listati/transport/approva_seller.sh transport mychannel
 ./listati/transport/get_all_assets.sh transport mychannel ##MODIFICHE APPORTATE

		***ERROR1 = Current Org cannot express objective Org will
		***ERROR2 = Current Org Org3MSP trying to approve but Buyer has not approved yet




APIs:

	SCENARIO 1
	
	Forest:
	localhost:3001
		/assets			[GET]
		/table			[GET]
		/trunk			[GET]
		/table/{id}		[GET]
		/trunk/{id}		[GET]
		/trunk/table		[POST]
		/trunk/{id}/newowner	[POST]
		/table/{id}/newowner	[POST]
		/history/trunk/{id}	[GET]
		/history/table/{id}	[GET]
		/audit			[POST] ---------|
		/trunk/new		[POST]		|
							|
	Id corp:					|
	localhost:27000	<-------------------------------|
		/api/product/audit   [POST] --------------------|
						  	        |
	localhost:5000 <----------------------------------------o
		/api/product/audit		[POST]		|
		/api/product/findall		[GET]		|
		/api/product/inserted/{id}	[POST]		|
								|
	localhost:3000 <----------------------------------------|
		/mock	[POST]
		/mock	[GET]
		
		
	SCENARIO 2
	
	Forest:
	localhost:3001 <-------------------------|
		/sell/approval/{id} 	[PUT]  --|--------------|
		/sell		    	[POST] --|------|	|
						 |	|	|
	Sawmill:				 |	|	|
	localhost:3002				 |	|	|
		/buy/approval/{id} 	[PUT]  --|------|-------|
		/buy		    	[POST] --|	|	|
							|	|
	Transport-shipping:				|	|
	localhost:3003 <--------------------------------|-------|
		/sell/new		[POST]	
		/buyer/{id}/{org}	[PUT]
		/seller/{id}/{org}	[PUT]
