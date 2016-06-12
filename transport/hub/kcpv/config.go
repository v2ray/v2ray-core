package kcpv

type AdvancedConfig struct {
	Mtu          int  `json:"MaximumTransmissionUnit"`
	Sndwnd       int  `json:"SendingWindowSize"`
	Rcvwnd       int  `json:"ReceivingWindowSize"`
	Fec          int  `json:"ForwardErrorCorrectionGroupSize"`
	Acknodelay   bool `json:"AcknowledgeNoDelay"`
	Dscp         int  `json:"Dscp"`
	ReadTimeout  int  `json:"ReadTimeout"`
	WriteTimeout int  `json:"WriteTimeout"`
}

type Config struct {
	Mode            string          `json:"Mode"`
	Key             string          `json:"EncryptionKey"`
	AdvancedConfigs *AdvancedConfig `json:"AdvancedConfig,omitempty"`
}

var DefaultAdvancedConfigs = &AdvancedConfig{
	Mtu: 1350, Sndwnd: 1024, Rcvwnd: 1024, Fec: 4, Dscp: 0, ReadTimeout: 18600, WriteTimeout: 18500, Acknodelay: false,
}
