// Copyright 2022 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package pack

import (
	"fish_net/cmd/user/domain"
	"fish_net/kitex_gen/user"
)

// User pack user info
func User(u *domain.User) *user.User {
	if u == nil {
		return nil
	}

	return &user.User{
		UserId:     int64(u.ID),
		Username:   u.Username,
		Nickname:   u.Nickname,
		Avater:     u.Avater,
		CreateTime: u.CreatedAt.Unix(),
		UpdateTime: u.CreatedAt.Unix(),
	}
}

// Users pack list of user info
func Users(us []*domain.User) []*user.User {
	users := make([]*user.User, 0)
	for _, u := range us {
		if temp := User(u); temp != nil {
			users = append(users, temp)
		}
	}
	return users
}

func UserNames(us []*domain.User) []string {
	userNames := make([]string, len(us))
	for _, u := range us {
		userNames = append(userNames, u.Username)
	}
	return userNames
}
