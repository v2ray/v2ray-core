package v2board

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"v2ray.com/core/app/proxyman/command"
	"v2ray.com/core/app/stats"
	feature_stats "v2ray.com/core/features/stats"

	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/vmess"
)

/*
   {
       "id": 1,
       "email": "contact@battleroa.ch",
       "t": 0,
       "u": 0,
       "d": 0,
       "transfer_enable": 1073741824000,
       "v2ray_user": {
           "uuid": "eddcd953-ea5d-406c-8f60-6516f5294acf",
           "email": "eddcd953-ea5d-406c-8f60-6516f5294acf@v2board.user",
           "alter_id": 2,
           "level": 0
       }
   }
*/
type V2BoardUser struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	T              int64     `json:"t"`
	U              int64     `json:"u"`
	D              int64     `json:"d"`
	TransferEnable int64     `json:"transfer_enable`
	V2RayUser      V2RayUser `json:"v2ray_user"`
}

type V2BoardTrafficLog struct {
	Uplink   int64  `json:"u"`
	Downlink int64  `json:"d"`
	UserID   string `json:"user_id"`
}

/*
{
	"uuid": "eddcd953-ea5d-406c-8f60-6516f5294acf",
	"email": "eddcd953-ea5d-406c-8f60-6516f5294acf@v2board.user",
	"alter_id": 2,
	"level": 0
}
*/
type V2RayUser struct {
	UUID    string `json:"uuid"`
	Email   string `json:"email"`
	AlterID uint32 `json:"alter_id"`
	Level   uint32 `json:"level"`
}

var TAG string = "proxy"

func (v *V2Board) Loop() {
	fetchUserTicker := time.NewTicker(1 * time.Minute)
	v.FetchUserList()
	statsTicker := time.NewTicker(10 * time.Second)
	v.GetStats()

	for {
		select {
		case <-fetchUserTicker.C:
			v.FetchUserList()

		case <-statsTicker.C:
			v.GetStats()
		}
	}
}

func (v *V2Board) FetchUserList() {
	resp, err := http.Get(v.UserListUri())
	if err != nil {
		newError("Failed to Fetch user list").Base(err).AtError().WriteToLog()
		return
	}
	defer resp.Body.Close()

	ret := []V2BoardUser{}
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		newError("Failed to Decode user list").Base(err).AtError().WriteToLog()
		return
	}
	v.maintain(&ret)
}

func (v *V2Board) GetStats() {
	trafficStats := map[string]*V2BoardTrafficLog{}
	manager, ok := v.stats.(*stats.Manager)
	if !ok {
		newError("stats.Manager is not set").AtError().WriteToLog()
	}
	counters := []feature_stats.Counter{}

	manager.Visit(func(name string, c feature_stats.Counter) bool {
		var value int64
		value = c.Value()
		if value == 0 {
			return true
		}
		counters = append(counters, c)
		username, direction := v.GetUserName(name)
		if username == "" {
			return true
		}
		if val, ok := trafficStats[username]; ok {
			if direction == "uplink" {
				val.Uplink += value
			} else {
				val.Downlink += value
			}
		} else {
			trafficStats[username] = &V2BoardTrafficLog{
				UserID: username,
			}
			if direction == "uplink" {
				trafficStats[username].Uplink += value
			} else {
				trafficStats[username].Downlink += value
			}
		}

		return true
	})

	jsonString, jsonErr := json.Marshal(trafficStats)
	newError("Submitting json record: ", trafficStats).AtInfo().WriteToLog()
	if jsonErr != nil {
		newError("Cannot marshal json object").AtError().WriteToLog()
		return
	}
	rsp, postErr := http.Post(v.ReportTrafficUri(), "application/json", bytes.NewReader(jsonString))
	if postErr != nil {
		newError("Cannot upload json traffic record").Base(postErr).AtWarning().WriteToLog()
		return
	}
	newError("Submitted status code is ", rsp.StatusCode).AtDebug().WriteToLog()

	if rsp.StatusCode != 200 {
		newError("Cannot upload json traffic record, status code is ", rsp.StatusCode).AtWarning().WriteToLog()
		return
	}
	for _, c := range counters {
		c.Set(0)
	}

}

var pattern *regexp.Regexp

//	pattern = regexp.MustCompile("user>>>([\\w\\-_])+>>>traffic>>>(uplink|downlink)")

func (v *V2Board) GetUserName(UserDetail string) (name, direction string) {
	matches := pattern.FindStringSubmatch(UserDetail)
	if matches == nil {
		newError("pattern '%s' not match\n", UserDetail).AtWarning().WriteToLog()
		return "", ""
	}
	name, direction = matches[1], matches[2]
	return
}

func (v *V2Board) maintain(list *[]V2BoardUser) {
	newset := NewUserSet()
	for _, user := range *list {
		newset.Add(user.V2RayUser)
	}

	waitfordelete := make([]V2RayUser, 0)
	waitforadd := make([]V2RayUser, 0)

	for _, user := range v.userset.List() {
		if !newset.Has(user) {
			waitfordelete = append(waitfordelete, user)
		}
	}

	for _, user := range newset.List() {
		if !v.userset.Has(user) {
			waitforadd = append(waitforadd, user)
		}
	}

	ctx := context.Background()

	for _, user := range waitfordelete {
		err := v.DeleteUser(ctx, &user)
		if err != nil {
			newError("Failed to delete user ", user).Base(err).AtWarning().WriteToLog()
		} else {
			v.userset.Remove(user)
			newError("Deleted user ", user).Base(err).AtInfo().WriteToLog()
		}

	}

	for _, user := range waitforadd {
		err := v.AddUser(ctx, &user)
		if err != nil {
			newError("Failed to add user ", user).Base(err).AtWarning().WriteToLog()
		} else {
			v.userset.Add(user)
			newError("Added user ", user).Base(err).AtInfo().WriteToLog()
		}
	}
}

func (v *V2Board) AddUser(ctx context.Context, u *V2RayUser) error {
	// type.CreateObject(context.Background(), )
	newError("Add user :", u.Email).AtInfo().WriteToLog()
	op := &command.AddUserOperation{
		User: &protocol.User{
			Level: u.Level,
			Email: u.Email,
			Account: serial.ToTypedMessage(&vmess.Account{
				Id:               u.UUID,
				AlterId:          u.AlterID,
				SecuritySettings: &protocol.SecurityConfig{Type: protocol.SecurityType_AUTO},
			}),
		},
	}

	handler, err := v.im.GetHandler(ctx, TAG)
	if err != nil {
		return newError("failed to get handler: ", TAG).Base(err)
	}

	return op.ApplyInbound(context.Background(), handler)
}

func (v *V2Board) DeleteUser(ctx context.Context, u *V2RayUser) error {
	newError("Delete user :", u.Email).AtInfo().WriteToLog()
	op := &command.RemoveUserOperation{
		Email: u.Email,
	}

	handler, err := v.im.GetHandler(ctx, TAG)
	if err != nil {
		return newError("failed to get handler: ", TAG).Base(err)
	}

	return op.ApplyInbound(context.Background(), handler)
}
func init() {
	pattern = regexp.MustCompile("user>>>(.+)>>>traffic>>>(uplink|downlink)")
}
