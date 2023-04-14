package certmagic

type ReqeustParam struct {
	Email       string `noste:"邮箱"`
	Domain      string `noste:"域名"`
	CaType      string `noste:"Ca 类型"`
	Provider    string `noste:"Ca 提供商"`
	SecretId    string `noste:"密钥 ID"`
	SecretKey   string `noste:"密钥"`
	EabKeyId    string `noste:"EAB密钥 ID"`
	EabMacKey   string `noste:"EAB密钥"`
	StoragePath string `noste:"存储目录"`
}

type Certificate struct {
	Names       []string
	NotAfter    int64
	NotBefore   int64
	Certificate [][]byte
	PrivateKey  []byte
	Issuer      map[string]any
}
