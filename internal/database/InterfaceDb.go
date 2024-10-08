package database

//SOLID - D
type UserRepository interface {
	AddNewUser(ac *Account)
	DeleteUser(id uint)
	CheckAvailibleEmail(ac *Account) bool
	CheckPasswordUser(ac *Account) bool
	GetIDUser(ac *Account) uint
	SetRefreshToken(idUser, refreshToken string)
	GetRefreshTokenUser(idUser string) string
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

//CheckPasswordForLogin
func CheckPasssword(repo UserRepository, ac *Account) bool {
	if check := repo.CheckPasswordUser(ac); check {
		return true
	} else {
		return false
	}
}

//GetIDFromDB for jwt
func GetID(repo UserRepository, ac *Account) uint {
	return repo.GetIDUser(ac)
}

func SetRefToken(repo UserRepository, idUser, refToken string) {
	repo.SetRefreshToken(idUser, refToken)
}

func GetRefToken(repo UserRepository, idUser string) string {
	return repo.GetRefreshTokenUser(idUser)
}
