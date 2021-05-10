### Golang File Encrypter
>This is a simple utility that use **AES 256 GCM mode** to encrypt simple txt files and later output their contents to the cmdline upon decryption

### DISCLAIMER
> The tool encrypts all the files fine but since on decryption the data is converted to a string and printed on the cmdline please only use it with txt files **only**


### Building the binary
#### using makefile
```bash:
make
```
#### using golang directly
```bash:
go build encrypter.go
```
### running the binary
#### encryption
```bash:
encrypter.exe -filename < your file here > -mode e -key < Your sample key >
```
#### decryption
>the program prints the details of the file on the cmdline
```bash:
encrypter.exe -filename < your file here > -mode d -key < Your sample key >
```

