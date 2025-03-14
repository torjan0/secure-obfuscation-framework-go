Secure Obfuscation Framework for Go
# Secure Obfuscation Framework for Go
A powerful command-line tool designed to obfuscate Go source code, enh
## Features
- **Variable Renaming**: Transforms local variable names into cryptic, rand
- **String Encryption**: Encrypts string literals using AES and decrypts the
- **Control Flow Flattening**: Restructures function logic into a state mach
- **Dead Code Injection**: Inserts meaningless statements to mislead rever
- **Customizable Obfuscation Levels**: Choose from `none`, `light`, `mediu
- **Verbose Logging**: Provides detailed output for debugging and transpa
## Installation
### Prerequisites
- Go 1.16 or later
- `garble` installed:
```bash
go install mvdan.cc/garble@latest
```
### Steps
1. Clone the repository:
```bash
git clone https://github.com/yourusername/secure-obfuscation-framewor
cd secure-obfuscation-framework-go
```
2. Install dependencies:
```bash
go get github.com/urfave/cli/v2
```
3. Build the tool (optional):
```bash
go build -o go-obfuscate cmd/main.go
```
## Usage
Run the tool with the build command:
```bash
go run cmd/main.go build [flags]
# OR, if built:
./go-obfuscate build [flags]
```
### Flags
- `--level, -l`: Obfuscation level (none, light, medium, heavy). Default: light.
- `--verbose, -v`: Enable detailed logging.
### Examples
Basic obfuscation:
```bash
go run cmd/main.go build --level=light
```
Medium obfuscation with logging:
```bash
go run cmd/main.go build --level=medium --verbose
```
Full obfuscation:
```bash
go run cmd/main.go build --level=heavy
```
The obfuscated binary is saved as `obfuscated_main`.
## Obfuscation Levels
| Level | Features Applied |
|--------|-----------------|
| none | No obfuscation; code remains unchanged. |
| light | Variable renaming only. |
| medium | Variable renaming + string encryption. |
| heavy | All above + control flow flattening + dead code injection. |
## How It Works
1. **Copy**: Copies the source code to a temporary directory to preserve th
2. **Transform**: Applies selected obfuscation techniques using the `go/as
3. **Build**: Compiles the transformed code with `garble` for additional bin
4. **Output**: Moves the resulting binary (`obfuscated_main`) to the curren
## Example Output
### Original Code
```go
package main
import "fmt"
func main() {
message := "Hello, World!"
if true {
fmt.Println(message)
}
}
```
### After Heavy Obfuscation (Simplified)
```go
package main
import (
"crypto/aes"
"crypto/cipher"
"encoding/base64"
"fmt"
)
func main() {
fmt.Println("deadcode")
xQwRtY8jK9pL2mN := decryptString("encrypted_string")
state := 0
for state < 1 {
switch state {
case 0:
if true {
fmt.Println(xQwRtY8jK9pL2mN)
}
state++
}
}
}
func decryptString(s string) string { /* ... */ }
func decryptAES(data []byte, keyStr string) string { /* ... */ }
```
## Troubleshooting
- `"garble: command not found"`: Ensure `garble` is installed and in your P
- **Permission errors**: Run with sufficient permissions for file operations.
- **Verbose mode**: Use `--verbose` to diagnose transformation issues.

## License
This project is licensed under the MIT License. See the LICENSE file for de
