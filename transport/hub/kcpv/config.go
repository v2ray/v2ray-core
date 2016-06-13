package kcpv

/*AdvancedConfig define behavior of KCP in detail

MaximumTransmissionUnit:
Largest protocol data unit that the layer can pass onwards
can be discovered by running tracepath

SendingWindowSize , ReceivingWindowSize:
the size of Tx/Rx window, by packet

ForwardErrorCorrectionGroupSize:
The the size of packet to generate a Forward Error Correction
packet, this is used to counteract packet loss.

AcknowledgeNoDelay:
Do not wait a certain of time before sending the ACK packet,
increase bandwich cost and might increase performance

Dscp:
Differentiated services code point,
be used to reconized to discriminate packet.
It is recommanded to keep it 0 to avoid being detected.

ReadTimeout,WriteTimeout:
Close the Socket if it have been silent for too long,
Small value can recycle zombie socket faster but
can cause v2ray to kill the proxy connection it is relaying,
Higher value can prevent server from closing zombie socket and
waste resources.
*/
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

/*Config define basic behavior of KCP
Mode:
can be one of these values:
fast3,fast2,fast,normal
<<<<<<- less delay
->>>>>> less bandwich wasted

EncryptionKey:
a string that will be the EncryptionKey of
All KCP connection we Listen-Accpet or
Dial, We are not very sure about how this
encryption hehave and DO use a unique randomly
generated key.
*/
type Config struct {
	Mode            string          `json:"Mode"`
	Key             string          `json:"EncryptionKey"`
	AdvancedConfigs *AdvancedConfig `json:"AdvancedConfig,omitempty"`
}

var DefaultAdvancedConfigs = &AdvancedConfig{
	Mtu: 1350, Sndwnd: 1024, Rcvwnd: 1024, Fec: 4, Dscp: 0, ReadTimeout: 600, WriteTimeout: 500, Acknodelay: false,
}
