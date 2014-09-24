package main

func IsJPEG(pic []byte) bool {
	// Magic Number
	// JPEG文件，前两个字节是0xff、0xd8，最后两个字节是0xff、0xd9
	if pic[0] == 0xff && pic[1] == 0xd8 &&
		pic[len(pic)-2] == 0xff && pic[len(pic)-1] == 0xd9 {
		return true
	}
	return false
}

func IsPNG(pic []byte) bool {
	// magic number
	// PNG图片前8个字节如下
	magicNumber := [...]byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}
	for i := 0; i < len(magicNumber); i++ {
		if pic[i] != magicNumber[i] {
			return false
		}
	}
	return true
}
