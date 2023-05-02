export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="peer-org-1"
export CORE_PEER_TLS_ROOTCERT_FILE=${FABRIC_HOME}/organizations/peer/peer-org-1/nodes/peer-org-1-peer-1/tls/msp/cacerts/root-tls-ca-7054.pem
export CORE_PEER_MSPCONFIGPATH=${FABRIC_HOME}/organizations/peer/peer-org-1/users/peer-org-1-org-ca-admin/msp
export CORE_PEER_ADDRESS=peer-org-1-peer-1:7051

peer lifecycle chaincode install token.tar.gz
