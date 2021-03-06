package datastore

import (
	"strings"

	"github.com/golang/protobuf/proto"
	config_proto "www.velocidex.com/golang/velociraptor/config/proto"
	crypto_proto "www.velocidex.com/golang/velociraptor/crypto/proto"
)

type TestDataStore struct {
	Subjects    map[string]proto.Message
	ClientTasks map[string][]*crypto_proto.GrrMessage
}

func NewTestDataStore() *TestDataStore {
	return &TestDataStore{
		Subjects:    make(map[string]proto.Message),
		ClientTasks: make(map[string][]*crypto_proto.GrrMessage),
	}
}

func (self *TestDataStore) GetClientTasks(config_obj *config_proto.Config,
	client_id string,
	do_not_lease bool) ([]*crypto_proto.GrrMessage, error) {
	result, _ := self.ClientTasks[client_id]
	if !do_not_lease {
		delete(self.ClientTasks, client_id)
	}
	return result, nil
}

func (self *TestDataStore) QueueMessageForClient(
	config_obj *config_proto.Config,
	client_id string,
	message *crypto_proto.GrrMessage) error {
	result, pres := self.ClientTasks[client_id]
	if !pres {
		result = make([]*crypto_proto.GrrMessage, 0)
	}

	result = append(result, message)

	self.ClientTasks[client_id] = result
	return nil
}

func (self *TestDataStore) UnQueueMessageForClient(
	config_obj *config_proto.Config,
	client_id string,
	message *crypto_proto.GrrMessage) error {

	result, pres := self.ClientTasks[client_id]
	if !pres {
		result = make([]*crypto_proto.GrrMessage, 0)
	}

	new_queue := make([]*crypto_proto.GrrMessage, len(result))
	for _, item := range result {
		if message.TaskId != item.TaskId {
			new_queue = append(new_queue, item)
		}
	}

	self.ClientTasks[client_id] = result
	return nil
}

func (self *TestDataStore) GetSubject(
	config_obj *config_proto.Config,
	urn string,
	message proto.Message) error {

	result, _ := self.Subjects[urn]
	if result != nil {
		proto.Merge(message, result)
	}
	return nil
}

func (self *TestDataStore) SetSubject(
	config_obj *config_proto.Config,
	urn string,
	message proto.Message) error {

	self.Subjects[urn] = message

	return nil
}

func (self *TestDataStore) DeleteSubject(
	config_obj *config_proto.Config,
	urn string) error {

	delete(self.Subjects, urn)

	return nil
}

// Lists all the children of a URN.
func (self *TestDataStore) ListChildren(
	config_obj *config_proto.Config,
	urn string,
	offset uint64, length uint64) ([]string, error) {

	result := []string{}

	for k := range self.Subjects {
		if strings.HasPrefix(k, urn) {
			k = strings.TrimLeft(strings.TrimPrefix(k, urn), "/")
			components := strings.Split(k, "/")
			if len(components) > 0 {
				result = append(result, components[0])
			}
		}
	}

	return result, nil
}

// Update the posting list index. Searching for any of the
// keywords will return the entity urn.
func (self *TestDataStore) SetIndex(
	config_obj *config_proto.Config,
	index_urn string,
	entity string,
	keywords []string) error {
	return nil
}

func (self *TestDataStore) UnsetIndex(
	config_obj *config_proto.Config,
	index_urn string,
	entity string,
	keywords []string) error {
	return nil
}

func (self *TestDataStore) CheckIndex(
	config_obj *config_proto.Config,
	index_urn string,
	entity string,
	keywords []string) error {
	return nil
}

func (self *TestDataStore) SearchClients(
	config_obj *config_proto.Config,
	index_urn string,
	query string, query_type string,
	offset uint64, limit uint64) []string {
	return nil
}

// Called to close all db handles etc. Not thread safe.
func (self *TestDataStore) Close() {
	RegisterImplementation("Test", NewTestDataStore())
}

func init() {
	RegisterImplementation("Test", NewTestDataStore())
}
