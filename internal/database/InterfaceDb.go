package database

//SOLID - D
type UserRepository interface {
	AddNewUser(ac *Account)
	DeleteUser(id uint)
	CheckAvailibleEmail(ac *Account) bool
}

func RegisterUser(repo UserRepository, ac *Account) {
	repo.AddNewUser(ac)
}

func RemoveUser(repo UserRepository, id uint) {
	repo.DeleteUser(id)
}

func EmailAvailible(repo UserRepository, ac *Account) bool {
	if check := repo.CheckAvailibleEmail(ac); check { //if true
		return true
	} else {
		return false
	}
}
