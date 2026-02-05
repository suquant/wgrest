package usecase

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/suquant/wgrest/internal/domain/entity"
)

func TestSortPeers_ByPublicKey(t *testing.T) {
	peers := []entity.Peer{
		{PublicKey: "charlie"},
		{PublicKey: "alpha"},
		{PublicKey: "bravo"},
	}

	sortPeers(peers, "pub_key")

	assert.Equal(t, "alpha", peers[0].PublicKey)
	assert.Equal(t, "bravo", peers[1].PublicKey)
	assert.Equal(t, "charlie", peers[2].PublicKey)
}

func TestSortPeers_ByPublicKeyDesc(t *testing.T) {
	peers := []entity.Peer{
		{PublicKey: "charlie"},
		{PublicKey: "alpha"},
		{PublicKey: "bravo"},
	}

	sortPeers(peers, "-pub_key")

	assert.Equal(t, "charlie", peers[0].PublicKey)
	assert.Equal(t, "bravo", peers[1].PublicKey)
	assert.Equal(t, "alpha", peers[2].PublicKey)
}

func TestSortPeers_ByReceiveBytes(t *testing.T) {
	peers := []entity.Peer{
		{PublicKey: "a", ReceiveBytes: 300},
		{PublicKey: "b", ReceiveBytes: 100},
		{PublicKey: "c", ReceiveBytes: 200},
	}

	sortPeers(peers, "receive_bytes")

	assert.Equal(t, int64(100), peers[0].ReceiveBytes)
	assert.Equal(t, int64(200), peers[1].ReceiveBytes)
	assert.Equal(t, int64(300), peers[2].ReceiveBytes)
}

func TestSortPeers_ByTotalBytes(t *testing.T) {
	peers := []entity.Peer{
		{PublicKey: "a", ReceiveBytes: 100, TransmitBytes: 200},  // 300
		{PublicKey: "b", ReceiveBytes: 50, TransmitBytes: 50},    // 100
		{PublicKey: "c", ReceiveBytes: 100, TransmitBytes: 100},  // 200
	}

	sortPeers(peers, "total_bytes")

	assert.Equal(t, "b", peers[0].PublicKey) // 100
	assert.Equal(t, "c", peers[1].PublicKey) // 200
	assert.Equal(t, "a", peers[2].PublicKey) // 300
}

func TestSortPeers_ByLastHandshakeTime(t *testing.T) {
	now := time.Now()
	peers := []entity.Peer{
		{PublicKey: "a", LastHandshakeTime: now},
		{PublicKey: "b", LastHandshakeTime: now.Add(-2 * time.Hour)},
		{PublicKey: "c", LastHandshakeTime: now.Add(-1 * time.Hour)},
	}

	sortPeers(peers, "last_handshake_time")

	assert.Equal(t, "b", peers[0].PublicKey)
	assert.Equal(t, "c", peers[1].PublicKey)
	assert.Equal(t, "a", peers[2].PublicKey)
}

func TestSortPeers_InvalidField(t *testing.T) {
	peers := []entity.Peer{
		{PublicKey: "charlie"},
		{PublicKey: "alpha"},
		{PublicKey: "bravo"},
	}

	// Should not change order for invalid field
	sortPeers(peers, "invalid_field")

	assert.Equal(t, "charlie", peers[0].PublicKey)
	assert.Equal(t, "alpha", peers[1].PublicKey)
	assert.Equal(t, "bravo", peers[2].PublicKey)
}

func TestSortPeers_EmptyField(t *testing.T) {
	peers := []entity.Peer{
		{PublicKey: "charlie"},
		{PublicKey: "alpha"},
	}

	// Should not change order for empty field
	sortPeers(peers, "")

	assert.Equal(t, "charlie", peers[0].PublicKey)
	assert.Equal(t, "alpha", peers[1].PublicKey)
}

func TestContainsIP(t *testing.T) {
	testCases := []struct {
		name     string
		ips      []string
		query    string
		expected bool
	}{
		{"match first", []string{"10.0.0.1/32", "10.0.0.2/32"}, "10.0.0.1", true},
		{"match second", []string{"10.0.0.1/32", "10.0.0.2/32"}, "10.0.0.2", true},
		{"partial match", []string{"10.0.0.1/32"}, "10.0", true},
		{"no match", []string{"10.0.0.1/32"}, "192.168", false},
		{"empty ips", []string{}, "10.0", false},
		{"empty query", []string{"10.0.0.1/32"}, "", true},
		{"case insensitive", []string{"FD00::1/64"}, "fd00", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := containsIP(tc.ips, tc.query)
			assert.Equal(t, tc.expected, result)
		})
	}
}
