package entity

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type User struct {
	Name     string
	Password string
	Email    string
	Phone    string
}

type Meeting struct {
	Sponsor       string
	Participators []string
	StartDate     time.Time
	EndDate       time.Time
	Title         string
}

const meetingPath = "data/meetings.json"
const userPath = "data/users.json"
const loginPath = "data/login"

func (m Meeting) IsParticipator(username string) bool {
	for i := 0; i < len(m.Participators); i++ {
		if username == m.Participators[i] {
			return true
		}
	}
	return false
}

type Storage struct {
	MeetingList []Meeting
	UserList    []User
}

func (s *Storage) WriteUsers() {
	f, err := os.Create(userPath)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	data, err := json.MarshalIndent(s.UserList, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	f.Write(data)
}

func (s *Storage) WriteMeetings() {
	f, err := os.Create(meetingPath)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	data, err := json.MarshalIndent(s.MeetingList, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	f.Write(data)
}

func (s *Storage) ReadMeetings() {
	if f, err := os.Open(meetingPath); err != nil {
		fmt.Println(err)
	} else {
		decoder := json.NewDecoder(f)
		decoder.Decode(&s.MeetingList)
	}
}

func (s *Storage) ReadUsers() {
	if f, err := os.Open(userPath); err != nil {
		fmt.Println(err)
	} else {
		decoder := json.NewDecoder(f)
		decoder.Decode(&s.UserList)
	}
}

func (s *Storage) CreateUser(u User) {
	s.UserList = append(s.UserList, u)
}

func (s *Storage) QueryUser(filter func(u User) bool) []User {
	users := []User{}
	for i := 0; i < len(s.UserList); i++ {
		if filter(s.UserList[i]) == true {
			users = append(users, s.UserList[i])
		}
	}
	return users
}

func (s *Storage) UserLogin(userName string, password string) error {
	s.ReadUsers()
	for i := 0; i < len(s.UserList); i++ {
		if s.UserList[i].Name == userName && s.UserList[i].Password == password {
			f, err := os.Create(loginPath)
			if err != nil {
				return err
			}
			f.Write([]byte(userName))
			return nil
		}
	}
	err := fmt.Errorf("failed: userName or password error!\n")
	return err
}
