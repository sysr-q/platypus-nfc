package main

type Tag struct {
	UID          string
	Type         int
	String       string
	Authorized   bool
	Confidential bool
	Block0       [16]byte
}
