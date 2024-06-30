package token

import "github.com/pkoukk/tiktoken-go"

// ModelToEncoding extend model encoding setting
var ModelToEncoding = map[string]string{
	"ep-20240603062111-s4snw": tiktoken.MODEL_O200K_BASE,
	"doubao-pro-32k-240515":   tiktoken.MODEL_O200K_BASE,
}
