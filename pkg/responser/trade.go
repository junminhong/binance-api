package responser

type Data struct {
	E  string `json:"e"`
	E1 int    `json:"E"`
	S  string `json:"s"`
	A  int    `json:"a"`
	P  string `json:"p"`
	Q  string `json:"q"`
	F  int    `json:"f"`
	L  int    `json:"l"`
	T  int    `json:"T"`
	M  bool   `json:"m"`
	M1 bool   `json:"M"`
}

type AggTradeData struct {
	Stream string
	Data   Data
}
