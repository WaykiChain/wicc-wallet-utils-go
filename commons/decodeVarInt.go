package commons

func EncodeInOldWay(value int64) []byte {
   size:=size(value)
   tmp:=make([]byte,((size*8+6)/7))
   ret:=[]byte{}
   len:=0
   n:=value
  // h:=byte(0x00)
	for{
		h:=byte(0)
		if len==0{
			h=0x00
		}else {
			h=0x80
		}

		tmp[len] = byte((byte(n) & 0x7f)|h)

		if n<= 0x7f{
			break
		}

		n=(n>>7)-1
		len+=1
	}

	for {
	   ret =append(ret,tmp[len])
		len-=1
		if len< 0 {
			break
		}
	}
	return ret
}

func size(value int64) int {
	ret := 0
	n := value

	for {
		ret += 1
		if n <= 0x7f {
			break
		}
		n = (n >> 7) - 1
	}

	return ret
}
