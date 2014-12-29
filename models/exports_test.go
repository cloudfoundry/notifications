package models

func ClearDB() {
	mutex.Lock()
	defer mutex.Unlock()

	_database = nil
}
