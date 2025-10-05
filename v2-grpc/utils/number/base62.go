package number

func generateBase62Chars() string {
	var chars []byte

	for i := '0'; i <= '9'; i++ {
		chars = append(chars, byte(i))
	}
	for i := 'a'; i <= 'z'; i++ {
		chars = append(chars, byte(i))
	}
	for i := 'A'; i <= 'Z'; i++ {
		chars = append(chars, byte(i))
	}
	return string(chars)
}

// 对id进行62进制编码，并且长度为length
func DecimalToBase62(num int64, length int) string {
	const base = 62

	base62Chars := generateBase62Chars()

	encoded := make([]byte, length)

	for i := length - 1; i >= 0; i-- {
		encoded[i] = base62Chars[num%base]
		num /= base
	}

	return string(encoded)
}