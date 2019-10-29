package User

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type User struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Email    string `json:"Email"`
}

func (user User) GetName() string {
	return user.Username
}

func (user *User) SetName(name string) {
	user.Username = name
}

func (user User) GetPassword() string {
	return user.Password
}

func (user *User) SetPassword(password string) {
	user.Password = password
}

func (user User) GetEmail() string {
	return user.Email
}

func (user *User) SetEmail(email string) {
	user.Email = email
}

//whether the user exits
func checkUser(name string) int {
	file, err := os.OpenFile("data/User.json", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	for decoder.More() {
		var users User
		decoder.Decode(&users)
		//exits
		if users.Username == name {
			file.Close()
			return 0
		}
	}
	//not exists
	return 1
}

//Get all user infomation
func GetAllUserInfo() []User {
	var users []User
	byteIn, err := ioutil.ReadFile("data/User.json")
	check(err)
	jsonStr := string(byteIn)
	json.Unmarshal([]byte(jsonStr), &users)
	return users
}

func check(r error) {
	if r != nil {
		log.Fatal(r)
	}
}

//register an  user with name, password, email
func RegisterAnUser(user *User) {
	flog, err := os.OpenFile("data/input_output.log", os.O_APPEND|os.O_WRONLY, 0666)
	defer flog.Close()
	check(err)
	logger := log.New(flog, "", log.LstdFlags)
	logger.Printf("agenda register -u %s -p %s -e %s", user.Username, user.Password, user.Email)

	var userInfo User
	userInfo.Username = user.Username
	userInfo.Password = user.Password
	userInfo.Email = user.Email

	file, err := os.OpenFile("data/User.json", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	num := checkUser(user.Username)
	AllUserInfo := GetAllUserInfo()
	if num == 1 { //not exited
		AllUserInfo = append(AllUserInfo, userInfo)
		encoder := json.NewEncoder(file)
		encoder.Encode(AllUserInfo)
		os.Stdout.WriteString("Register succeed!\n")
		logger.Print("Register succeed!\n")
	} else {
		os.Stdout.WriteString("The userName have been registered.\n")
		logger.Print("The userName have been registered./n")
	}
	file.Close()

}

func LogIn(user *User) {
	users := GetAllUserInfo()
	flog, err := os.OpenFile("data/input_output.log", os.O_APPEND|os.O_WRONLY, 0600)
	defer flog.Close()

	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(flog, "", log.LstdFlags)
	logger.Printf("agenda login -u %s -p %s", user.Username, user.Password)

	//get current username
	fin, err1 := os.Open("data/current.txt")
	if err1 != nil {
		log.Fatal(err1)
	}
	defer fin.Close()
	reader := bufio.NewReader(fin)
	name, _ := reader.ReadString('\n')

	if name != "" {
		os.Stdout.WriteString("You have already logged in.\n")
		logger.Print("You have already logged in. \n")
		return
	}

	//find the password
	var length int = len(users)
	var flag = false
	for ii := 0; ii < length; ii++ {
		if users[ii].Username == user.Username {
			if users[ii].Password == user.Password {
				fout, _ := os.Create("data/current.txt")
				defer fout.Close()
				fout.WriteString(user.Username)
				os.Stdout.WriteString("Login successfully!\n")
				logger.Print("Login successfully!\n")
			} else {
				os.Stdout.WriteString("Password is incorrect!\n")
				logger.Print("Password is incorrect!\n")
			}
			flag = true
		}
	}
	if flag == false {
		os.Stdout.WriteString("Username is not correct.\n")
		logger.Print("Username is not correct.\n")
	}
}

func LogOut() {
	flog, err := os.OpenFile("data/input_output.log", os.O_APPEND|os.O_WRONLY, 0600)
	defer flog.Close()
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(flog, "", log.LstdFlags)
	logger.Print("Agenda logout\n")

	fin, err1 := os.Open("data/current.txt")
	if err1 != nil {
		log.Fatal(err1)
	}
	defer fin.Close()
	reader := bufio.NewReader(fin)
	name, _ := reader.ReadString('\n')

	if name != "" {
		fout, _ := os.Create("data/current.txt")
		defer fout.Close()
		fout.WriteString("")
		os.Stdout.WriteString("logout successfully!\n")
		logger.Print("Logout successfully!\n")
	} else {
		os.Stdout.WriteString("You are not logged in.\n")
		logger.Print("You are not logged in.\n")
	}

}
