package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// maxJSONBodyBytes membatasi ukuran body JSON semua endpoint.
// Endpoint aplikasi ini tidak menerima payload besar.
const maxJSONBodyBytes = 1 << 20 // 1 MiB

// decodeJSONBody melakukan decode JSON secara ketat:
//   - body dibatasi MaxBytesReader (anti memory abuse);
//   - field tak dikenal ditolak (strict contract, anti typo/ambiguity);
//   - trailing JSON token setelah object ditolak.
//
// Error yang dikembalikan aman untuk diekspos ke client (tanpa detail parser).
func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodyBytes)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return err
	}

	// Pastikan tidak ada data tambahan setelah object utama.
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return fmt.Errorf("request body must contain a single JSON object")
		}
		return err
	}
	return nil
}
