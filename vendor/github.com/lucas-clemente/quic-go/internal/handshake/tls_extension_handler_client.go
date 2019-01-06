package handshake

import (
	"errors"
	"fmt"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/qerr"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"github.com/marten-seemann/qtls"
)

type extensionHandlerClient struct {
	ourParams  *TransportParameters
	paramsChan chan<- TransportParameters

	origConnID        protocol.ConnectionID
	initialVersion    protocol.VersionNumber
	supportedVersions []protocol.VersionNumber
	version           protocol.VersionNumber

	logger utils.Logger
}

var _ tlsExtensionHandler = &extensionHandlerClient{}

// newExtensionHandlerClient creates a new extension handler for the client.
func newExtensionHandlerClient(
	params *TransportParameters,
	origConnID protocol.ConnectionID,
	initialVersion protocol.VersionNumber,
	supportedVersions []protocol.VersionNumber,
	version protocol.VersionNumber,
	logger utils.Logger,
) (tlsExtensionHandler, <-chan TransportParameters) {
	// The client reads the transport parameters from the Encrypted Extensions message.
	// The paramsChan is used in the session's run loop's select statement.
	// We have to use an unbuffered channel here to make sure that the session actually processes the transport parameters immediately.
	paramsChan := make(chan TransportParameters)
	return &extensionHandlerClient{
		ourParams:         params,
		paramsChan:        paramsChan,
		origConnID:        origConnID,
		initialVersion:    initialVersion,
		supportedVersions: supportedVersions,
		version:           version,
		logger:            logger,
	}, paramsChan
}

func (h *extensionHandlerClient) GetExtensions(msgType uint8) []qtls.Extension {
	if messageType(msgType) != typeClientHello {
		return nil
	}
	h.logger.Debugf("Sending Transport Parameters: %s", h.ourParams)
	return []qtls.Extension{{
		Type: quicTLSExtensionType,
		Data: (&clientHelloTransportParameters{
			InitialVersion: h.initialVersion,
			Parameters:     *h.ourParams,
		}).Marshal(),
	}}
}

func (h *extensionHandlerClient) ReceivedExtensions(msgType uint8, exts []qtls.Extension) error {
	if messageType(msgType) != typeEncryptedExtensions {
		return nil
	}

	var found bool
	eetp := &encryptedExtensionsTransportParameters{}
	for _, ext := range exts {
		if ext.Type != quicTLSExtensionType {
			continue
		}
		if err := eetp.Unmarshal(ext.Data); err != nil {
			return err
		}
		found = true
	}
	if !found {
		return errors.New("EncryptedExtensions message didn't contain a QUIC extension")
	}

	// check that the negotiated_version is the current version
	if eetp.NegotiatedVersion != h.version {
		return qerr.Error(qerr.VersionNegotiationMismatch, "current version doesn't match negotiated_version")
	}
	// check that the current version is included in the supported versions
	if !protocol.IsSupportedVersion(eetp.SupportedVersions, h.version) {
		return qerr.Error(qerr.VersionNegotiationMismatch, "current version not included in the supported versions")
	}
	// if version negotiation was performed, check that we would have selected the current version based on the supported versions sent by the server
	if h.version != h.initialVersion {
		negotiatedVersion, ok := protocol.ChooseSupportedVersion(h.supportedVersions, eetp.SupportedVersions)
		if !ok || h.version != negotiatedVersion {
			return qerr.Error(qerr.VersionNegotiationMismatch, "would have picked a different version")
		}
	}

	params := eetp.Parameters
	// check that the server sent a stateless reset token
	if len(params.StatelessResetToken) == 0 {
		return errors.New("server didn't sent stateless_reset_token")
	}
	// check the Retry token
	if !h.origConnID.Equal(params.OriginalConnectionID) {
		return fmt.Errorf("expected original_connection_id to equal %s, is %s", h.origConnID, params.OriginalConnectionID)
	}
	h.logger.Debugf("Received Transport Parameters: %s", &params)
	h.paramsChan <- params
	return nil
}
