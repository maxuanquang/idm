package configs

import "time"

type Hash struct {
	Cost int `yaml:"cost"`
}

type Token struct {
	Duration            uint32 `yaml:"duration"`
	RS512KeyPairBitSize uint16 `yaml:"rs512_key_pair_bit_size"`
}

func (t *Token) GetTokenDuration() time.Duration {
	return time.Duration(t.Duration) * time.Second
}

type Auth struct {
	Hash  Hash  `yaml:"hash"`
	Token Token `yaml:"token"`
}
