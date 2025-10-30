package main

func BackupOne() (value interface{}) {
	defer func() {
		if p := recover(); p != nil {
            value = p
		}
	}()
	panic(1)
}

