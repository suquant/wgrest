package usecase

import (
	"sort"
	"strings"

	"github.com/suquant/wgrest/internal/domain/entity"
	"github.com/suquant/wgrest/internal/infrastructure/wgquick"
	"github.com/suquant/wgrest/internal/infrastructure/wireguard"
)

// PeerUseCase handles business logic for peer operations.
type PeerUseCase struct {
	wgClient   *wireguard.Client
	wgquickSvc *wgquick.Service
}

// NewPeerUseCase creates a new peer use case.
func NewPeerUseCase(
	wgClient *wireguard.Client,
	wgquickSvc *wgquick.Service,
) *PeerUseCase {
	return &PeerUseCase{
		wgClient:   wgClient,
		wgquickSvc: wgquickSvc,
	}
}

// ListPeers returns all peers for a device with pagination, filtering, and sorting.
func (uc *PeerUseCase) ListPeers(deviceName string, page, perPage int, query, sortField string) ([]entity.Peer, int, error) {
	peers, err := uc.wgClient.ListPeers(deviceName)
	if err != nil {
		return nil, 0, err
	}

	// Apply search filter
	if query != "" {
		filtered := make([]entity.Peer, 0)
		queryLower := strings.ToLower(query)
		for _, p := range peers {
			if strings.Contains(strings.ToLower(p.PublicKey), queryLower) ||
				strings.Contains(strings.ToLower(p.Endpoint), queryLower) ||
				containsIP(p.AllowedIPs, queryLower) {
				filtered = append(filtered, p)
			}
		}
		peers = filtered
	}

	// Apply sorting
	sortPeers(peers, sortField)

	total := len(peers)

	// Apply pagination
	if perPage <= 0 {
		perPage = 100
	}
	if page < 0 {
		page = 0
	}

	start := page * perPage
	if start >= len(peers) {
		return []entity.Peer{}, total, nil
	}

	end := start + perPage
	if end > len(peers) {
		end = len(peers)
	}

	return peers[start:end], total, nil
}

// GetPeer returns a specific peer.
func (uc *PeerUseCase) GetPeer(deviceName string, urlSafePubKey string) (*entity.Peer, error) {
	return uc.wgClient.GetPeer(deviceName, urlSafePubKey)
}

// CreatePeer creates a new peer.
func (uc *PeerUseCase) CreatePeer(deviceName string, req entity.PeerCreateOrUpdateRequest) (*entity.Peer, error) {
	peer, err := uc.wgClient.CreatePeer(deviceName, req)
	if err != nil {
		return nil, err
	}

	// Trigger config save
	uc.saveDeviceConfig(deviceName)

	return peer, nil
}

// UpdatePeer updates a peer.
func (uc *PeerUseCase) UpdatePeer(deviceName string, urlSafePubKey string, req entity.PeerCreateOrUpdateRequest) (*entity.Peer, error) {
	peer, err := uc.wgClient.UpdatePeer(deviceName, urlSafePubKey, req)
	if err != nil {
		return nil, err
	}

	// Trigger config save
	uc.saveDeviceConfig(deviceName)

	return peer, nil
}

// DeletePeer deletes a peer.
func (uc *PeerUseCase) DeletePeer(deviceName string, urlSafePubKey string) (*entity.Peer, error) {
	peer, err := uc.wgClient.DeletePeer(deviceName, urlSafePubKey)
	if err != nil {
		return nil, err
	}

	// Trigger config save
	uc.saveDeviceConfig(deviceName)

	return peer, nil
}

func (uc *PeerUseCase) saveDeviceConfig(deviceName string) {
	// Use wg showconf to save current state
	if err := uc.wgquickSvc.SaveFromShowconf(deviceName); err != nil {
		// Log but don't fail
	}
}


func containsIP(ips []string, query string) bool {
	for _, ip := range ips {
		if strings.Contains(strings.ToLower(ip), query) {
			return true
		}
	}
	return false
}

func sortPeers(peers []entity.Peer, sortField string) {
	if sortField == "" {
		return
	}

	desc := strings.HasPrefix(sortField, "-")
	if desc {
		sortField = sortField[1:]
	}

	sort.Slice(peers, func(i, j int) bool {
		var less bool
		switch sortField {
		case "pub_key":
			less = peers[i].PublicKey < peers[j].PublicKey
		case "receive_bytes":
			less = peers[i].ReceiveBytes < peers[j].ReceiveBytes
		case "transmit_bytes":
			less = peers[i].TransmitBytes < peers[j].TransmitBytes
		case "total_bytes":
			totalI := peers[i].ReceiveBytes + peers[i].TransmitBytes
			totalJ := peers[j].ReceiveBytes + peers[j].TransmitBytes
			less = totalI < totalJ
		case "last_handshake_time":
			less = peers[i].LastHandshakeTime.Before(peers[j].LastHandshakeTime)
		default:
			return false
		}

		if desc {
			return !less
		}
		return less
	})
}
