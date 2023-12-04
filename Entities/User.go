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
	FavoritesBooks     string  `json:"favorites_books"`
	FavoritesFilms     string  `json:"favorites_films"`
	FavoritesGames     string  `json:"favorites_games"`
	FavoritesMeals     string  `json:"favorites_meals"`
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
	Favorites_books      string `json:"favorites_books"`
	Favorites_films      string `json:"favorites_films"`
	Favorites_games      string `json:"favorites_games"`
	Favorites_meals      string `json:"favorites_meals"`
	Description          string `json:"description"`
	Family_status        int8   `json:"family_status"`
	Place_of_residence   string `json:"place_of_residence"`
	Attitude_to_smocking int8   `json:"attitude_to_smocking"`
	Attitude_to_sport    int8   `json:"attitude_to_sport"`
	Attitude_to_alcohol  int8   `json:"attitude_to_alcohol"`
	Dreams               string `json:"dreams"`
}

type GetUserDTO struct {
	UpdateUserDTO
	Friends     []int64 `json:"friends"`
	Subscribers []int64 `json:"subscribers"`
}

type Check struct {
	Email string `json:"email"`
	Login string `json:"login"`
}

type RefreshDTO struct {
	Id uint `json:"id"`
	// Token has a password
	Token string `json:"token"`
}

type RefreshResponseDTO struct {
	Avatar       string `json:"avatar"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type FriendOrSubsDTO struct {
	Friends     []int64 `json:"friends"`
	Subscribers []int64 `json:"subscribers"`
}

type GetFriendsAndSubsDTOForUser struct {
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	Avatar      string  `json:"avatar"`
	Friends     []int64 `json:"friends"`
	Subscribers []int64 `json:"subscribers"`
}

type GetFriendsAndSubsDTOForClient struct {
	Friends     []int64 `json:"friends"`
	Subscribers []int64 `json:"subscribers"`
}

type GetFriendsAndSubsDTO struct {
	Client GetFriendsAndSubsDTOForClient `json:"client"`
	User   GetFriendsAndSubsDTOForUser   `json:"user"`
}

type SignInDTO struct {
	Password string `json:"password" binding:"required"`
	Login    string `json:"login"`
	Email    string `json:"email"`
}

type SignInReturnDTO struct {
	ID                    uint
	Email, Login          string
	Name, Surname, Avatar string
}

type UserToCheck struct {
	ID       int
	Password string
}

type MiniUser struct {
	Avatar  string `json:"avatar"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	ID      uint   `json:"id"`
}

type FriendUser struct {
	MiniUser
	IsClientSub bool `json:"is_client_sub"`
}
