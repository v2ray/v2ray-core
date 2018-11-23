package handshake

import (
	"errors"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/qerr"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"github.com/marten-seemann/qtls"
)

type extensionHandlerServer struct {
	ourParams  *TransportParameters
	paramsChan chan<- TransportParameters

	version           protocol.VersionNumber
	supportedVersions []protocol.VersionNumber

	logger utils.Logger
}

var _ tlsExtensionHandler = &extensionHandlerServer{}

// newExtensionHandlerServer creates a new extension handler for the server
func newExtensionHandlerServer(
	params *TransportParameters,
	supportedVersions []protocol.VersionNumber,
	version protocol.VersionNumber,
	logger utils.Logger,
) (tlsExtensionHandler, <-chan TransportParameters) {
	// Processing the ClientHello is performed statelessly (and from a single go-routine).
	// Therefore, we have to use a buffered chan to pass the transport parameters to that go routine.
	paramsChan := make(chan TransportParameters)
	return &extensionHandlerServer{
		ourParams:         params,
		paramsChan:        paramsChan,
		supportedVersions: supportedVersions,
		version:           version,
		logger:            logger,
	}, paramsChan
}

func (h *extensionHandlerServer) GetExtensions(msgType uint8) []qtls.Extension {
	if messageType(msgType) != typeEncryptedExtensions {
		return nil
	}
	h.logger.Debugf("Sending Transport Parameters: %s", h.ourParams)
	return []qtls.Extension{{
		Type: quicTLSExtensionType,
		Data: (&encryptedExtensionsTransportParameters{
			NegotiatedVersion: h.version,
			SupportedVersions: protocol.GetGreasedVersions(h.supportedVersions),
			Parameters:        *h.ourParams,
		}).Marshal(),
	}}
}

func (h *extensionHandlerServer) ReceivedExtensions(msgType uint8, exts []qtls.Extension) error {
	if messageType(msgType) != typeClientHello {
		return nil
	}
	var found bool
	chtp := &clientHelloTransportParameters{}
	for _, ext := range exts {
		if ext.Type != quicTLSExtensionType {
			continue
		}
		if err := chtp.Unmarshal(ext.Data); err != nil {
			return err
		}
		found = true
	}
	if !found {
		return errors.New("ClientHello didn't contain a QUIC extension")
	}

	// perform the stateless version negotiation validation:
	// make sure that we would have sent a Version Negotiation Packet if the client offered the initial version
	// this is the case if and only if the initial version is not contained in the supported versions
	if chtp.InitialVersion != h.version && protocol.IsSupportedVersion(h.supportedVersions, chtp.InitialVersion) {
		return qerr.Error(qerr.VersionNegotiationMismatch, "Client should have used the initial version")
	}
	h.logger.Debugf("Received Transport Parameters: %s", &chtp.Parameters)
	h.paramsChan <- chtp.Parameters
	return nil
}
