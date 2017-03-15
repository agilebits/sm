# sm
Simple secret management tool for server configuration

## How to build

```
go get -u github.com/agilebits/sm
cd ~/go/src/github.com/agilebits/sm
go install
```
## Encrypt/decrypt data on development machines

```
cat app-config.yml | sm encrypt > app-config.sm
cat app-config.sm | sm decrypt
```

On the first run, the utility will generate a new master key and store it in `~/.sm/masterkey` file. The `masterkey` must be saved and copied across all developer machines.


## Encrypt/decrypt data with Amazon Web Service KMS

First, you have to create a master key using AWS IAM and give yourself permissions to use this key for encryption and decryption. 

```
cat app-config.yml | sm encrypt --env aws --region us-east-1 --master arn:aws:kms:us-east-1:123123123123:key/d845cfa3-0719-4631-1d00-10ab63e40ddf	> app-config.sm

cat app-config.sm | sm decrypt
```



