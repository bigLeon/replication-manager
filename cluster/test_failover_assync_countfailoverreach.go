package cluster

func (cluster *Cluster) testNumberFailOverLimitReach(conf string, test string) bool {
	if cluster.initTestCluster(conf, test) == false {
		return false
	}
	cluster.conf.MaxDelay = 0
	err := cluster.disableSemisync()
	if err != nil {
		cluster.LogPrintf("ERROR : %s", err)
		cluster.closeTestCluster(conf, test)
		return false
	}

	SaveMaster := cluster.master
	SaveMasterURL := cluster.master.URL

	cluster.LogPrintf("INFO :  Master is %s", cluster.master.URL)
	cluster.master.State = stateFailed
	cluster.conf.Interactive = false
	cluster.master.FailCount = cluster.conf.MaxFail
	cluster.conf.FailLimit = 3
	cluster.conf.FailTime = 0
	cluster.failoverCtr = 3
	cluster.conf.CheckFalsePositiveHeartbeat = false
	cluster.conf.RplChecks = false
	cluster.conf.MaxDelay = 20
	cluster.checkfailed()

	cluster.waitFailoverEnd()
	cluster.LogPrintf("INFO : New Master  %s ", cluster.master.URL)
	if cluster.master.URL != SaveMasterURL {
		cluster.LogPrintf("INFO : Old master %s ==  Next master %s  ", SaveMasterURL, cluster.master.URL)
		cluster.closeTestCluster(conf, test)
		SaveMaster.FailCount = 0
		return false
	}
	SaveMaster.FailCount = 0
	cluster.closeTestCluster(conf, test)
	return true
}
