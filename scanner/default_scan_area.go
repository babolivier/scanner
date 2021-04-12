package scanner

import (
	"github.com/babolivier/scanner/common"
)

// storeCurrentScanArea retrieves the current scan area and stores it in memory.
func (s *Scanner) storeCurrentScanArea() (err error) {
	s.defaultScanArea = new(common.ScanArea)
	if s.defaultScanArea.TLX, err = s.conn.GetOption("tl-x"); err != nil {
		return err
	}
	if s.defaultScanArea.TLY, err = s.conn.GetOption("tl-y"); err != nil {
		return err
	}
	if s.defaultScanArea.BRX, err = s.conn.GetOption("br-x"); err != nil {
		return err
	}
	if s.defaultScanArea.BRY, err = s.conn.GetOption("br-y"); err != nil {
		return err
	}

	return nil
}

// resetScanArea resets the scan area parameters on the SANE connection using the values
// retrieved when the SANE connection was established.
func (s *Scanner) resetScanArea() error {
	if _, err := s.conn.SetOption("tl-x", s.defaultScanArea.TLX); err != nil {
		return err
	}

	if _, err := s.conn.SetOption("tl-y", s.defaultScanArea.TLY); err != nil {
		return err
	}

	if _, err := s.conn.SetOption("br-x", s.defaultScanArea.BRX); err != nil {
		return err
	}

	if _, err := s.conn.SetOption("br-y", s.defaultScanArea.BRY); err != nil {
		return err
	}

	return nil
}
