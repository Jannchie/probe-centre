package job

// Init is to start all jobs.
func Init() {
	go createBiliTasks()
	go RemoveIpRecord()
}