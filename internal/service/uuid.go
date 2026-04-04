package service

import (
	"crypto/rand"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

const (
	uuidVersion4Mask   = 0x0f
	uuidVersion4       = 0x40
	uuidVariantMask    = 0x3f
	uuidVariantRFC4122 = 0x80
	uuidSize           = 16
	uuidVersionByte    = 6
	uuidVariantByte    = 8
)

func newUUIDBytes() []byte {
	id := make([]byte, uuidSize)

	_, _ = rand.Read(id)

	id[uuidVersionByte] = (id[uuidVersionByte] & uuidVersion4Mask) | uuidVersion4
	id[uuidVariantByte] = (id[uuidVariantByte] & uuidVariantMask) | uuidVariantRFC4122

	return id
}

func uuidToString(id pgtype.UUID) string {
	b := id.Bytes

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func parseUUID(s string) (pgtype.UUID, error) {
	var id pgtype.UUID

	err := id.Scan(s)
	if err != nil {
		return id, fmt.Errorf("parsing UUID %q: %w", s, err)
	}

	return id, nil
}
