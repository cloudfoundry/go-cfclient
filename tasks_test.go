package cfclient

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListTasks(t *testing.T) {
	Convey("List Tasks", t, func() {
		setup(MockRoute{"GET", "/v3/tasks", []string{listTasksPayloadPage1, listTasksPayloadPage2}, "", 200, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		task, err := client.ListTasks()
		So(err, ShouldBeNil)

		So(len(task), ShouldEqual, 4)

		So(task[0].GUID, ShouldEqual, "d5cc22ec-99a3-4e6a-af91-a44b4ab7b6fa")
		So(task[0].State, ShouldEqual, "SUCCEEDED")
		So(task[0].SequenceID, ShouldEqual, 1)
		So(task[0].MemoryInMb, ShouldEqual, 512)
		So(task[0].DiskInMb, ShouldEqual, 1024)
		So(task[0].CreatedAt.String(), ShouldEqual, time.Date(2016, 05, 04, 17, 00, 41, 0, time.FixedZone("UTC", 0)).String())

		So(task[1].GUID, ShouldEqual, "63b4cd89-fd8b-4bf1-a311-7174fcc907d6")
		So(task[1].State, ShouldEqual, "FAILED")
		So(task[1].SequenceID, ShouldEqual, 2)
		So(task[1].MemoryInMb, ShouldEqual, 1024)
		So(task[1].DiskInMb, ShouldEqual, 1024)
		So(task[1].CreatedAt.String(), ShouldEqual, time.Date(2016, 05, 04, 17, 00, 43, 0, time.FixedZone("UTC", 0)).String())

		So(task[2].GUID, ShouldEqual, "abcdefc-99a3-4e6a-af91-a44b4ab7b6fa")
		So(task[2].State, ShouldEqual, "SUCCEEDED")
		So(task[2].SequenceID, ShouldEqual, 3)
		So(task[2].MemoryInMb, ShouldEqual, 512)
		So(task[2].DiskInMb, ShouldEqual, 1024)
		So(task[2].CreatedAt.String(), ShouldEqual, time.Date(2016, 05, 04, 17, 00, 44, 0, time.FixedZone("UTC", 0)).String())

		So(task[3].GUID, ShouldEqual, "hijklm9-fd8b-4bf1-a311-7174fcc907d6")
		So(task[3].State, ShouldEqual, "SUCCEEDED")
		So(task[3].SequenceID, ShouldEqual, 4)
		So(task[3].MemoryInMb, ShouldEqual, 1024)
		So(task[3].DiskInMb, ShouldEqual, 1024)
		So(task[3].CreatedAt.String(), ShouldEqual, time.Date(2016, 05, 04, 17, 00, 46, 0, time.FixedZone("UTC", 0)).String())
	})
}
func TestListTasksByQuery(t *testing.T) {
	Convey("List Tasks", t, func() {
		setup(MockRoute{"GET", "/v3/tasks", []string{listTasksPayloadPage1, listTasksPayloadPage2}, "", 200, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		task, err := client.ListTasksByQuery(nil)
		So(err, ShouldBeNil)

		So(len(task), ShouldEqual, 4)

		So(task[0].GUID, ShouldEqual, "d5cc22ec-99a3-4e6a-af91-a44b4ab7b6fa")
		So(task[0].State, ShouldEqual, "SUCCEEDED")
		So(task[0].SequenceID, ShouldEqual, 1)
		So(task[0].MemoryInMb, ShouldEqual, 512)
		So(task[0].DiskInMb, ShouldEqual, 1024)
		So(task[0].CreatedAt.String(), ShouldEqual, time.Date(2016, 05, 04, 17, 00, 41, 0, time.FixedZone("UTC", 0)).String())

		So(task[1].GUID, ShouldEqual, "63b4cd89-fd8b-4bf1-a311-7174fcc907d6")
		So(task[1].State, ShouldEqual, "FAILED")
		So(task[1].SequenceID, ShouldEqual, 2)
		So(task[1].MemoryInMb, ShouldEqual, 1024)
		So(task[1].DiskInMb, ShouldEqual, 1024)
		So(task[1].CreatedAt.String(), ShouldEqual, time.Date(2016, 05, 04, 17, 00, 43, 0, time.FixedZone("UTC", 0)).String())

		So(task[3].GUID, ShouldEqual, "hijklm9-fd8b-4bf1-a311-7174fcc907d6")
		So(task[3].State, ShouldEqual, "SUCCEEDED")
		So(task[3].SequenceID, ShouldEqual, 4)
		So(task[3].MemoryInMb, ShouldEqual, 1024)
		So(task[3].DiskInMb, ShouldEqual, 1024)
		So(task[3].CreatedAt.String(), ShouldEqual, time.Date(2016, 05, 04, 17, 00, 46, 0, time.FixedZone("UTC", 0)).String())
	})
}

func TestCreateTask(t *testing.T) {
	Convey("Create Task", t, func() {
		mocks := []MockRoute{
			{"POST", "/v3/apps/740ebd2b-162b-469a-bd72-3edb96fabd9a/tasks", []string{createTaskPayload}, "", 201, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		tr := TaskRequest{
			Command:          "rake db:migrate",
			Name:             "migrate",
			MemoryInMegabyte: 512,
			DiskInMegabyte:   1024,
			DropletGUID:      "740ebd2b-162b-469a-bd72-3edb96fabd9a",
		}
		task, err := client.CreateTask(tr)
		So(err, ShouldBeNil)

		So(task.Command, ShouldEqual, "rake db:migrate")
		So(task.Name, ShouldEqual, "migrate")
		So(task.DiskInMb, ShouldEqual, 1024)
		So(task.MemoryInMb, ShouldEqual, 512)
		So(task.DropletGUID, ShouldEqual, "740ebd2b-162b-469a-bd72-3edb96fabd9a")
	})
}

func TestCreateTaskFails(t *testing.T) {
	Convey("Create Task fails", t, func() {
		mocks := []MockRoute{
			{"POST", "/v3/apps/740ebd2b-162b-469a-bd72-3edb96fabd9a/tasks", []string{errorV3Payload}, "", 400, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		tr := TaskRequest{
			Command:          "rake db:migrate",
			Name:             "migrate",
			MemoryInMegabyte: 512,
			DiskInMegabyte:   1024,
			DropletGUID:      "740ebd2b-162b-469a-bd72-3edb96fabd9a",
		}

		task, err := client.CreateTask(tr)
		So(err.Error(), ShouldEqual, "Error creating task: cfclient error (CF-UnprocessableEntity|10008): something went wrong")
		So(task.Name, ShouldBeEmpty)
	})
}

func TestTerminateTask(t *testing.T) {
	Convey("Terminate Task", t, func() {
		mocks := []MockRoute{
			{"PUT", "/v3/tasks/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/cancel", []string{""}, "", 202, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		errTerm := client.TerminateTask("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx")
		So(errTerm, ShouldBeNil)
	})
}

func TestGetTask(t *testing.T) {
	Convey("Create Task", t, func() {
		mocks := []MockRoute{
			{"GET", "/v3/tasks/740ebd2b-162b-469a-bd72-3edb96fabd9a", []string{createTaskPayload}, "", 200, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		task, err := client.GetTaskByGuid("740ebd2b-162b-469a-bd72-3edb96fabd9a")
		So(err, ShouldBeNil)

		So(task.Command, ShouldEqual, "rake db:migrate")
		So(task.Name, ShouldEqual, "migrate")
		So(task.DiskInMb, ShouldEqual, 1024)
		So(task.MemoryInMb, ShouldEqual, 512)
		So(task.DropletGUID, ShouldEqual, "740ebd2b-162b-469a-bd72-3edb96fabd9a")
	})
}

func TestTasksByApp(t *testing.T) {
	Convey("List Tasks by App", t, func() {
		setup(MockRoute{"GET", "/v3/apps/ccc25a0f-c8f4-4b39-9f1b-de9f328d0ee5/tasks", []string{listTasksByAppPayloadPage1, listTasksByAppPayloadPage2}, "", 200, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		task, err := client.TasksByApp("ccc25a0f-c8f4-4b39-9f1b-de9f328d0ee5")
		So(err, ShouldBeNil)

		So(len(task), ShouldEqual, 4)

		So(task[0].GUID, ShouldEqual, "d5cc22ec-99a3-4e6a-af91-a44b4ab7b6fa")
		So(task[0].State, ShouldEqual, "SUCCEEDED")
		So(task[0].SequenceID, ShouldEqual, 1)
		So(task[0].MemoryInMb, ShouldEqual, 512)
		So(task[0].DiskInMb, ShouldEqual, 1024)
		So(task[0].CreatedAt.String(), ShouldEqual, time.Date(2016, 05, 04, 17, 00, 41, 0, time.FixedZone("UTC", 0)).String())

		So(task[1].GUID, ShouldEqual, "63b4cd89-fd8b-4bf1-a311-7174fcc907d6")
		So(task[1].State, ShouldEqual, "FAILED")
		So(task[1].SequenceID, ShouldEqual, 2)
		So(task[1].MemoryInMb, ShouldEqual, 1024)
		So(task[1].DiskInMb, ShouldEqual, 1024)
		So(task[1].CreatedAt.String(), ShouldEqual, time.Date(2016, 05, 04, 17, 00, 43, 0, time.FixedZone("UTC", 0)).String())

		So(task[2].GUID, ShouldEqual, "abcdefc-99a3-4e6a-af91-a44b4ab7b6fa")
		So(task[2].State, ShouldEqual, "SUCCEEDED")
		So(task[2].SequenceID, ShouldEqual, 3)
		So(task[2].MemoryInMb, ShouldEqual, 512)
		So(task[2].DiskInMb, ShouldEqual, 1024)
		So(task[2].CreatedAt.String(), ShouldEqual, time.Date(2016, 05, 04, 17, 00, 44, 0, time.FixedZone("UTC", 0)).String())

		So(task[3].GUID, ShouldEqual, "hijklm9-fd8b-4bf1-a311-7174fcc907d6")
		So(task[3].State, ShouldEqual, "SUCCEEDED")
		So(task[3].SequenceID, ShouldEqual, 4)
		So(task[3].MemoryInMb, ShouldEqual, 1024)
		So(task[3].DiskInMb, ShouldEqual, 1024)
		So(task[3].CreatedAt.String(), ShouldEqual, time.Date(2016, 05, 04, 17, 00, 46, 0, time.FixedZone("UTC", 0)).String())
	})
}

func TestTasksByAppByQuery(t *testing.T) {
	Convey("List Tasks by App", t, func() {
		setup(MockRoute{"GET", "/v3/apps/ccc25a0f-c8f4-4b39-9f1b-de9f328d0ee5/tasks", []string{listTasksByAppPayloadPage1, listTasksByAppPayloadPage2}, "", 200, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		task, err := client.TasksByAppByQuery("ccc25a0f-c8f4-4b39-9f1b-de9f328d0ee5", nil)
		So(err, ShouldBeNil)

		So(len(task), ShouldEqual, 4)

		So(task[0].GUID, ShouldEqual, "d5cc22ec-99a3-4e6a-af91-a44b4ab7b6fa")
		So(task[0].State, ShouldEqual, "SUCCEEDED")
		So(task[0].SequenceID, ShouldEqual, 1)
		So(task[0].MemoryInMb, ShouldEqual, 512)
		So(task[0].DiskInMb, ShouldEqual, 1024)
		So(task[0].CreatedAt.String(), ShouldEqual, time.Date(2016, 05, 04, 17, 00, 41, 0, time.FixedZone("UTC", 0)).String())

		So(task[1].GUID, ShouldEqual, "63b4cd89-fd8b-4bf1-a311-7174fcc907d6")
		So(task[1].State, ShouldEqual, "FAILED")
		So(task[1].SequenceID, ShouldEqual, 2)
		So(task[1].MemoryInMb, ShouldEqual, 1024)
		So(task[1].DiskInMb, ShouldEqual, 1024)
		So(task[1].CreatedAt.String(), ShouldEqual, time.Date(2016, 05, 04, 17, 00, 43, 0, time.FixedZone("UTC", 0)).String())

		So(task[2].GUID, ShouldEqual, "abcdefc-99a3-4e6a-af91-a44b4ab7b6fa")
		So(task[2].State, ShouldEqual, "SUCCEEDED")
		So(task[2].SequenceID, ShouldEqual, 3)
		So(task[2].MemoryInMb, ShouldEqual, 512)
		So(task[2].DiskInMb, ShouldEqual, 1024)
		So(task[2].CreatedAt.String(), ShouldEqual, time.Date(2016, 05, 04, 17, 00, 44, 0, time.FixedZone("UTC", 0)).String())

		So(task[3].GUID, ShouldEqual, "hijklm9-fd8b-4bf1-a311-7174fcc907d6")
		So(task[3].State, ShouldEqual, "SUCCEEDED")
		So(task[3].SequenceID, ShouldEqual, 4)
		So(task[3].MemoryInMb, ShouldEqual, 1024)
		So(task[3].DiskInMb, ShouldEqual, 1024)
		So(task[3].CreatedAt.String(), ShouldEqual, time.Date(2016, 05, 04, 17, 00, 46, 0, time.FixedZone("UTC", 0)).String())
	})
}
