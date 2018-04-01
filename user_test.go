package demux

import "testing"

func TestTopic_AddUser(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Empty name", args{""}, true},
		{"Valid name", args{"name"}, false},
		// Duplicate name has to come after Valid name!
		{"Duplicate name", args{"name"}, true},
	}
	topic, _ := NewTopic("TestTopic_AddUser")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				user *User
				err  error
			)
			if user, err = topic.AddUser(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Topic.AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if user == nil {
				t.Errorf("User is nil")
				return
			}

			if user.Name != tt.args.name {
				t.Errorf("User added with wrong name. Expected: %q, got: %q", tt.args.name, user.Name)
				return
			}

			addedUser, ok := topic.users[tt.args.name]
			if !ok {
				t.Errorf("User should have been added, wasn't")
				return
			}

			if addedUser.Name != user.Name {
				t.Errorf("Expected user to be added under name %q, got %q", user.Name, addedUser.Name)
				return
			}

			if user.devices == nil {
				t.Errorf("Devices added with nil slice")
			}
		})
	}
}
