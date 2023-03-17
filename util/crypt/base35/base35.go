package base35

//仿照base64编码。此编码，生成仅有小写字母和数字的编码。
//参考 url https://www.cnblogs.com/Mishell/p/12241872.html
import (
	"bytes"
	"math/big"
	"strconv"
)

var baseBets = []byte("123456789abcdefghijklmnopqrstuvwxyz")

// Encode 编码
func Encode(input string) ([]byte, string) {
	x := big.NewInt(0).SetBytes([]byte(input))
	base := big.NewInt(35)
	//fmt.Println(base)
	zero := big.NewInt(0)
	mod := &big.Int{}
	var result []byte
	// 被除数/除数=商……余数
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, baseBets[mod.Int64()])
		//fmt.Println(result)
	}
	ReverseBytes(result)
	return result, string(result)
}

//EncodeStr EncodeStr
func EncodeStr(input string) string {
	_, str := Encode(input)
	return str
}

//EncodeU32Str EncodeStr
func EncodeU32Str(input uint32) string {
	i := string(strconv.Itoa(int(input)))
	_, str := Encode(i)
	return str
}

// Decode 解码
func Decode(input []byte) ([]byte, string) {
	result := big.NewInt(0)
	for _, b := range input {
		charIndex := bytes.IndexByte(baseBets, b)
		result.Mul(result, big.NewInt(35))
		result.Add(result, big.NewInt(int64(charIndex)))
	}
	decoded := result.Bytes()
	if input[0] == baseBets[0] {
		decoded = append([]byte{0x00}, decoded...)
	}
	return decoded, string(decoded)
}

// ReverseBytes 翻转字节
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
