package internet

import "net"

var v2AlternativeDialer *V2AlternativeDialerT

type V2AlternativeDialerT interface {
	Dial(nw string, ad string) (net.Conn, error)
}

func SubstituteDialer(substituteWith V2AlternativeDialerT) error {
	v2AlternativeDialer = &substituteWith
	return nil
}

func isDefaultDialerSubstituted() bool {
	return (v2AlternativeDialer != nil)
}
