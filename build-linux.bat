SET CGO_ENABLED=0 
SET GOOS=linux
SET GOARCH=amd64
cd main && go build -o ../linux/v2ray -ldflags "-s -w"