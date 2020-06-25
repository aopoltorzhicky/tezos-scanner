package p2p

import "time"

// ScannerOption -
type ScannerOption func(*Scanner)

// WithThreadsCount -
func WithThreadsCount(threadsCount int64) ScannerOption {
	return func(scanner *Scanner) {
		scanner.threadsCount = threadsCount
	}
}

// WithAttemptsDuration -
func WithAttemptsDuration(attemptsSeconds int64) ScannerOption {
	return func(scanner *Scanner) {
		scanner.attemptsDuration = time.Duration(attemptsSeconds) * time.Second
	}
}

// WithDropAfter -
func WithDropAfter(dropAfter int64) ScannerOption {
	return func(scanner *Scanner) {
		scanner.dropAfter = dropAfter
	}
}

// WithSyncedTime -
func WithSyncedTime(syncedTime int64) ScannerOption {
	return func(scanner *Scanner) {
		scanner.syncedTime = syncedTime
	}
}
