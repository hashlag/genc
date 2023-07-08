# GENC — CLI file encryption tool

GENC is a command line file encryption tool written in GoLang. It provides strong encryption via AES-256, integrity verification via HMAC-SHA512 and uses PBKDF2 SHA-512 (1 048 576 iterations, 384 bit salt) to derive keys from passwords.

## Screenshot

![screenshot](https://raw.githubusercontent.com/hashlag/static/main/genc/genc.gif)

## Features

- Single file encryption/decryption
- Recursive directory encryption/decryption
- Integrity verification
- Optional additional prompting to hide passwords from command line history (enabled by default)
- Optional autodeletion of original files after encryption/decryption
- Compact size `~3MiB`
- Low RAM consumption
- Сross-platform

## Installation

GENC pre-compiled binaries are stored in `bin/` and do not require installation. Just copy the suitable binary to any directory you want. You're ready to go!

Direct links to binaries:

- [Windows x64](https://github.com/hashlag/genc/raw/main/bin/genc-win-amd64.exe)
- [Windows x32](https://github.com/hashlag/genc/raw/main/bin/genc-win-386.exe)
- [Linux x64](https://github.com/hashlag/genc/raw/main/bin/genc-linux-amd64)
- [Linux x32](https://github.com/hashlag/genc/raw/main/bin/genc-linux-386)

## Compilation from source

1. Make sure you have GoLang installed
2. Clone the repository via `git clone https://github.com/hashlag/genc`
3. Navigate to the project directory
4. Build GENC using `go build .`
5. Try it by running the resulting binary with the `-h` flag!

## Usage

### Encrypting a single file

```
genc -m E -t file_to_encrypt.txt
```

Then you will be asked for a password.
After successful encryption GENC will print something like this:

```
Encryption succeeded.
file_to_encrypt.txt --> file_to_encrypt.txt.genc
```

You can also use `-d` to delete `file_to_encrypt.txt` after encryption.

### Decrypting a single file

```
genc -m D -t file_to_decrypt.txt.genc
```

You will be asked for a password.
Output will be like:

```
Decrypted.
file_to_decrypt.txt.genc --> file_to_decrypt.txt
```

### Encrypting/decrypting a directory

There are no special command line flags for directory encryption.
Just pass a directory as a target. GENC is able to distinguish a directory from a file by itself.

```
genc -m E -t directory_to_encrypt
```

and to decrypt:

```
genc -m D -t directory_to_encrypt
```

### Autodelete

You can use `-d` flag to delete original files after encryption or encrypted ones after decryption.

### Passing password as a command line flag

**Unsafe since password may remain in the command line history.**

```
genc -m E -t file.txt -p password
```

## License

This project is licensed under the MIT License. See the [LICENSE](https://raw.githubusercontent.com/hashlag/genc/main/LICENSE) file for details.

## Third-Party Libraries

This project uses:
- [prompter](https://github.com/Songmu/prompter) library by [Songmu](https://github.com/Songmu), which is licensed under the MIT License.
See [the library's LICENSE](https://raw.githubusercontent.com/Songmu/prompter/main/LICENSE) file for details.
- [go-isatty](https://github.com/mattn/go-isatty) library by [Yasuhiro MATSUMOTO](https://github.com/mattn), which is licensed under the MIT License.
See [the library's LICENSE](https://raw.githubusercontent.com/mattn/go-isatty/master/LICENSE) file for details.