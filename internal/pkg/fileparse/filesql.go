package fileparse

const (
	sql_createtable3 = `
	CREATE TABLE [dbo].[iFixsvr_JF_Info] (
		[ID] [decimal](18, 0) IDENTITY (1, 1) NOT NULL ,
		[IID] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[IField1] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[IField2] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[IField3] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[IField4] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[IField5] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[IField6] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[IField7] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[IField8] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[IModTime] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[IDesc] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL 
	) ON [PRIMARY]
	`

	sql_notfinishfile = `
	SELECT top 1 FID,FMd5,FProgress FROM [iFixsvr_JF_Files] WHERE FFinish ='0'
	 `

	sql_updateprogress = `
	UPDATE iFixsvr_JF_Files SET FProgress = ?, FFinish = ?, FModTime = ? WHERE (FMD5 = ? and FID = ?)
	 `

	sql_insertinfo = `
	 INSERT INTO [iFixsvr_JF_Info] (IID,IField1,IField2,IField3,IField4,IModTime) VALUES (?,?,?,?,?,?)
	 `
)
