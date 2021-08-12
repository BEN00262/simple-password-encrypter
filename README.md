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
### Running the binary
#### saving a credential
```bash:
encrypter.exe -filename < your file here > -mode s -key < Your sample key > -target < the site/anything u will use as a unique identifier > -password < the actual password >
```
#### showing all the saved targets
>the program prints all the saved targets

```bash:
encrypter.exe -filename < your file here > -mode r -key < Your sample key >
```

#### reading a password
>the program copies the password to the clipboard if it is found otherwise it informs u

```bash:
encrypter.exe -filename < your file here > -mode r -key < Your sample key > -target < the unique id for the saved password > 
```

#### deleting a credential
>the program deletes the saved credential from its store

```bash:
encrypter.exe -filename < your file here > -mode d -key < Your sample key > -target < the unique id for the saved password > 
```
