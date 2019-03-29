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
