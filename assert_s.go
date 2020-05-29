package yee

func assertS(judge bool, error string) {
	if !judge {
		panic(error)
	}
}
