package Entities

type ChatList struct {
	Name  string `json:"name"`
	Chats []int  `json:"chats"`
}

type RefreshToken struct {
	Token string `json:"token"`
}

// User ID is a primary key
type User struct {
	//Primary key
	ID uint `json:"id"`
	//Unique
	Login    string `json:"login"`
	Password string `json:"password"`
	//Unique
	Email              string  `json:"email"`
	Name               string  `json:"name"`
	Surname            string  `json:"surname"`
	Friends            []int64 `json:"friends"`
	Subscribers        []int64 `json:"subscribers"`
	Avatar             string  `json:"avatar"`
	Birthday           string  `json:"birthday"`
	FavouritesBooks    string  `json:"favourites_books"`
	FavouritesFilms    string  `json:"favourites_films"`
	FavouritesGames    string  `json:"favourites_games"`
	FavouritesMeals    string  `json:"favourites_meals"`
	Description        string  `json:"description"`
	FamilyStatus       string  `json:"family_status"`
	PlaceOfResidence   string  `json:"place_of_residence"`
	AttitudeToSmocking string  `json:"attitude_to_smocking"`
	AttitudeToSport    string  `json:"attitude_to_sport"`
	AttitudeToAlcohol  string  `json:"attitude_to_alcohol"`
	Dreams             string  `json:"dreams"`
	ChatLists          string  `json:"chat_lists"`
	AllChats           []uint  `json:"all_chats"`
}

type UserDTO struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Surname  string `json:"surname" binding:"required"`
}

type UpdateUserDTO struct {
	Name                 string `json:"name"`
	Surname              string `json:"surname"`
	Avatar               string `json:"avatar"`
	Birthday             string `json:"birthday"`
	Favourites_books     string `json:"favourites_books"`
	Favourites_films     string `json:"favourites_films"`
	Favourites_games     string `json:"favourites_games"`
	Favourites_meals     string `json:"favourites_meals"`
	Description          string `json:"description"`
	Family_status        string `json:"family_status"`
	Place_of_residence   string `json:"place_of_residence"`
	Attitude_to_smocking string `json:"attitude_to_smocking"`
	Attitude_to_sport    string `json:"attitude_to_sport"`
	Attitude_to_alcohol  string `json:"attitude_to_alcohol"`
	Dreams               string `json:"dreams"`
}

type GetUserDTO struct {
	UpdateUserDTO
	Friends     []int64 `json:"friends"`
	Subscribers []int64 `json:"subscribers"`
}

type FriendOrSubsDTO struct {
	Friends     []int64 `json:"friends"`
	Subscribers []int64 `json:"subscribers"`
}

type GetFriendsAndSubsDTOForOneUser struct {
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	Avatar      string  `json:"avatar"`
	Friends     []int64 `json:"friends"`
	Subscribers []int64 `json:"subscribers"`
}

type GetFriendsAndSubsDTO struct {
	Client GetFriendsAndSubsDTOForOneUser `json:"client"`
	User   GetFriendsAndSubsDTOForOneUser `json:"user"`
}

type SignInDTO struct {
	Password string `json:"password" binding:"required"`
	Login    string `json:"login"`
	Email    string `json:"email"`
}

type SignInReturnDTO struct {
	ID           uint
	Email, Login string
}

type UserToCheck struct {
	ID       int
	Password string
}

type MiniUser struct {
	Avatar  string
	Name    string
	Surname string
	ID      uint
}

type FriendUser struct {
	MiniUser
	Subscribers []int64
}
