package main

type location interface {
	Scan()
	Print()
	PrintFiles(int)
	Save()
	Load()
	GetRandomFilePath() (string, error)
}
